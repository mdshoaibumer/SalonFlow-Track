package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/internal/core/security"
)

const (
	maxFailedAttempts    = 5
	lockDurationMinutes  = 30
	sessionDurationHours = 8
	rememberMeDays       = 30
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrAccountLocked      = errors.New("account is locked due to too many failed attempts")
	ErrAccountInactive    = errors.New("account is inactive")
	ErrSessionExpired     = errors.New("session has expired")
	ErrSessionNotFound    = errors.New("session not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserExists         = errors.New("username or email already exists")
	ErrPermissionDenied   = errors.New("permission denied")
)

// AuthUseCase handles authentication and authorization logic.
type AuthUseCase struct {
	repo       ports.AuthRepository
	log        *slog.Logger
	appVersion string
}

// NewAuthUseCase creates a new AuthUseCase.
func NewAuthUseCase(repo ports.AuthRepository, log *slog.Logger, appVersion string) *AuthUseCase {
	return &AuthUseCase{repo: repo, log: log, appVersion: appVersion}
}

// LoginInput holds login request data.
type LoginInput struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RememberMe bool   `json:"remember_me"`
	DeviceID   string `json:"device_id"`
}

// LoginOutput holds login response data.
type LoginOutput struct {
	Token     string             `json:"token"`
	User      domain.SessionInfo `json:"user"`
	ExpiresAt time.Time          `json:"expires_at"`
}

// Login authenticates a user and creates a session.
func (uc *AuthUseCase) Login(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	// Try username first, then email
	user, err := uc.repo.GetUserByUsername(ctx, input.Username)
	if err != nil {
		user, err = uc.repo.GetUserByEmail(ctx, input.Username)
		if err != nil {
			uc.log.Warn("login failed: user not found", "username", input.Username)
			return nil, ErrInvalidCredentials
		}
	}

	// Check if account is locked
	if user.IsLocked {
		if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
			uc.log.Warn("login attempt on locked account", "user_id", user.ID, "username", user.Username)
			return nil, ErrAccountLocked
		}
		// Lock expired, unlock
		if err := uc.repo.UnlockUser(ctx, user.ID); err != nil {
			return nil, fmt.Errorf("unlock user: %w", err)
		}
		if err := uc.repo.ResetFailedAttempts(ctx, user.ID); err != nil {
			return nil, fmt.Errorf("reset attempts: %w", err)
		}
	}

	// Check if active
	if !user.IsActive {
		return nil, ErrAccountInactive
	}

	// Verify password
	valid, err := security.VerifyPassword(input.Password, user.PasswordHash)
	if err != nil || !valid {
		// Increment failed attempts
		if err := uc.repo.IncrementFailedAttempts(ctx, user.ID); err != nil {
			uc.log.Error("failed to increment failed attempts", "error", err)
		}

		// Lock if exceeded max
		if user.FailedAttempts+1 >= maxFailedAttempts {
			lockUntil := time.Now().Add(lockDurationMinutes * time.Minute).Format(time.RFC3339)
			if err := uc.repo.LockUser(ctx, user.ID, &lockUntil); err != nil {
				uc.log.Error("failed to lock user", "error", err)
			}
			uc.log.Warn("account locked due to failed attempts", "user_id", user.ID, "username", user.Username)

			// Audit: account locked
			uc.audit(ctx, user, "account_locked", "auth", "", "Account locked after too many failed attempts", "critical")
			return nil, ErrAccountLocked
		}

		uc.log.Warn("login failed: invalid password", "username", input.Username, "attempts", user.FailedAttempts+1)

		// Audit: failed login
		uc.audit(ctx, user, "login_failed", "auth", "", "Invalid password", "warning")
		return nil, ErrInvalidCredentials
	}

	// Successful login - reset failed attempts
	if err := uc.repo.ResetFailedAttempts(ctx, user.ID); err != nil {
		uc.log.Error("failed to reset attempts", "error", err)
	}
	if err := uc.repo.UpdateLastLogin(ctx, user.ID); err != nil {
		uc.log.Error("failed to update last login", "error", err)
	}

	// Generate session token
	token, err := security.GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	// Calculate expiry
	var expiresAt time.Time
	if input.RememberMe {
		expiresAt = time.Now().Add(time.Duration(rememberMeDays) * 24 * time.Hour)
	} else {
		expiresAt = time.Now().Add(time.Duration(sessionDurationHours) * time.Hour)
	}

	// Create session
	session := &domain.Session{
		ID:           uuid.New().String(),
		UserID:       user.ID,
		TokenHash:    security.HashToken(token),
		DeviceID:     input.DeviceID,
		RememberMe:   input.RememberMe,
		ExpiresAt:    expiresAt,
		LastActiveAt: time.Now(),
		CreatedAt:    time.Now(),
	}

	if err := uc.repo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	// Get user permissions
	permissions, err := uc.repo.GetUserPermissions(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("get permissions: %w", err)
	}

	// Get user roles
	roles, err := uc.repo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("get roles: %w", err)
	}
	roleNames := make([]string, len(roles))
	for i, r := range roles {
		roleNames[i] = r.Name
	}

	// Audit: successful login
	uc.audit(ctx, user, "login", "auth", "", "Successful login", "info")

	uc.log.Info("user logged in", "user_id", user.ID, "username", user.Username)

	return &LoginOutput{
		Token:     token,
		ExpiresAt: expiresAt,
		User: domain.SessionInfo{
			UserID:      user.ID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Email:       user.Email,
			Roles:       roleNames,
			Permissions: permissions,
		},
	}, nil
}

// Logout invalidates a session.
func (uc *AuthUseCase) Logout(ctx context.Context, token string) error {
	tokenHash := security.HashToken(token)
	session, err := uc.repo.GetSessionByToken(ctx, tokenHash)
	if err != nil {
		return ErrSessionNotFound
	}

	if err := uc.repo.DeleteSession(ctx, session.ID); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	user, _ := uc.repo.GetUserByID(ctx, session.UserID)
	if user != nil {
		uc.audit(ctx, user, "logout", "auth", "", "User logged out", "info")
	}

	uc.log.Info("user logged out", "user_id", session.UserID)
	return nil
}

// ValidateSession checks if a session token is valid and returns user info.
func (uc *AuthUseCase) ValidateSession(ctx context.Context, token string) (*domain.SessionInfo, error) {
	tokenHash := security.HashToken(token)
	session, err := uc.repo.GetSessionByToken(ctx, tokenHash)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	// Check expiry
	if time.Now().After(session.ExpiresAt) {
		_ = uc.repo.DeleteSession(ctx, session.ID)
		return nil, ErrSessionExpired
	}

	// Touch session (update last_active_at)
	_ = uc.repo.TouchSession(ctx, session.ID)

	// Get user
	user, err := uc.repo.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !user.IsActive || user.IsLocked {
		_ = uc.repo.DeleteSession(ctx, session.ID)
		return nil, ErrAccountInactive
	}

	// Get permissions
	permissions, err := uc.repo.GetUserPermissions(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("get permissions: %w", err)
	}

	// Get roles
	roles, err := uc.repo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("get roles: %w", err)
	}
	roleNames := make([]string, len(roles))
	for i, r := range roles {
		roleNames[i] = r.Name
	}

	return &domain.SessionInfo{
		UserID:      user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Roles:       roleNames,
		Permissions: permissions,
	}, nil
}

// CheckPermission verifies a user has a specific permission.
func (uc *AuthUseCase) CheckPermission(ctx context.Context, token string, permission string) error {
	info, err := uc.ValidateSession(ctx, token)
	if err != nil {
		return err
	}

	for _, p := range info.Permissions {
		if p == permission {
			return nil
		}
	}

	uc.log.Warn("permission denied", "user_id", info.UserID, "permission", permission)
	return ErrPermissionDenied
}

// ChangePasswordInput holds password change data.
type ChangePasswordInput struct {
	UserID      string `json:"user_id"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// ChangePassword changes a user's password.
func (uc *AuthUseCase) ChangePassword(ctx context.Context, input ChangePasswordInput) error {
	user, err := uc.repo.GetUserByID(ctx, input.UserID)
	if err != nil {
		return ErrUserNotFound
	}

	// Verify old password
	valid, err := security.VerifyPassword(input.OldPassword, user.PasswordHash)
	if err != nil || !valid {
		return ErrInvalidCredentials
	}

	// Validate new password
	if err := security.ValidatePasswordStrength(input.NewPassword); err != nil {
		return err
	}

	// Hash new password
	hash, err := security.HashPassword(input.NewPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	if err := uc.repo.UpdatePassword(ctx, user.ID, hash); err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	// Invalidate all sessions except current
	_ = uc.repo.DeleteUserSessions(ctx, user.ID)

	uc.audit(ctx, user, "password_changed", "auth", "", "Password changed", "info")
	uc.log.Info("password changed", "user_id", user.ID)
	return nil
}

// ResetPasswordInput holds admin password reset data.
type ResetPasswordInput struct {
	UserID      string `json:"user_id"`
	NewPassword string `json:"new_password"`
}

// ResetPassword allows an admin to reset another user's password.
func (uc *AuthUseCase) ResetPassword(ctx context.Context, input ResetPasswordInput) error {
	user, err := uc.repo.GetUserByID(ctx, input.UserID)
	if err != nil {
		return ErrUserNotFound
	}

	if err := security.ValidatePasswordStrength(input.NewPassword); err != nil {
		return err
	}

	hash, err := security.HashPassword(input.NewPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	if err := uc.repo.UpdatePassword(ctx, user.ID, hash); err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	// Force password change on next login
	user.MustChangePassword = true
	if err := uc.repo.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	_ = uc.repo.DeleteUserSessions(ctx, user.ID)

	uc.audit(ctx, user, "password_reset", "auth", "", "Password reset by admin", "warning")
	uc.log.Info("password reset by admin", "user_id", user.ID)
	return nil
}

// CreateUserInput holds data for creating a new user.
type CreateUserInput struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
	Phone       string `json:"phone"`
	RoleID      string `json:"role_id"`
}

// CreateUser creates a new user with a role.
func (uc *AuthUseCase) CreateUser(ctx context.Context, input CreateUserInput, createdBy string) (*domain.User, error) {
	// Check uniqueness
	existing, _ := uc.repo.GetUserByUsername(ctx, input.Username)
	if existing != nil {
		return nil, ErrUserExists
	}
	if input.Email != "" {
		existing, _ = uc.repo.GetUserByEmail(ctx, input.Email)
		if existing != nil {
			return nil, ErrUserExists
		}
	}

	// Validate password
	if err := security.ValidatePasswordStrength(input.Password); err != nil {
		return nil, err
	}

	// Hash password
	hash, err := security.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	now := time.Now()
	user := &domain.User{
		ID:                 uuid.New().String(),
		Username:           input.Username,
		Email:              input.Email,
		PasswordHash:       hash,
		DisplayName:        input.DisplayName,
		Phone:              input.Phone,
		IsActive:           true,
		MustChangePassword: true,
		PasswordChangedAt:  now,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := uc.repo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	// Assign role
	if input.RoleID != "" {
		if err := uc.repo.AssignRole(ctx, user.ID, input.RoleID, createdBy); err != nil {
			return nil, fmt.Errorf("assign role: %w", err)
		}
	}

	uc.audit(ctx, user, "user_created", "users", user.ID, fmt.Sprintf("User '%s' created", user.Username), "info")
	uc.log.Info("user created", "user_id", user.ID, "username", user.Username)
	return user, nil
}

// UpdateUserInput holds data for updating a user.
type UpdateUserInput struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Phone       string `json:"phone"`
	IsActive    bool   `json:"is_active"`
}

// UpdateUser updates user profile fields.
func (uc *AuthUseCase) UpdateUser(ctx context.Context, input UpdateUserInput) (*domain.User, error) {
	user, err := uc.repo.GetUserByID(ctx, input.ID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	user.Email = input.Email
	user.DisplayName = input.DisplayName
	user.Phone = input.Phone
	user.IsActive = input.IsActive
	user.UpdatedAt = time.Now()

	if err := uc.repo.UpdateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	uc.log.Info("user updated", "user_id", user.ID)
	return user, nil
}

// DeleteUser deactivates a user.
func (uc *AuthUseCase) DeleteUser(ctx context.Context, userID string) error {
	user, err := uc.repo.GetUserByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	if err := uc.repo.DeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	_ = uc.repo.DeleteUserSessions(ctx, userID)

	uc.audit(ctx, user, "user_deleted", "users", userID, fmt.Sprintf("User '%s' deleted", user.Username), "warning")
	uc.log.Info("user deleted", "user_id", userID)
	return nil
}

// GetUsers returns all users.
func (uc *AuthUseCase) GetUsers(ctx context.Context) ([]domain.User, error) {
	users, err := uc.repo.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	// Populate roles for each user
	for i := range users {
		roles, err := uc.repo.GetUserRoles(ctx, users[i].ID)
		if err == nil {
			users[i].Roles = roles
		}
	}
	return users, nil
}

// GetUser returns a single user by ID.
func (uc *AuthUseCase) GetUser(ctx context.Context, id string) (*domain.User, error) {
	user, err := uc.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	roles, err := uc.repo.GetUserRoles(ctx, user.ID)
	if err == nil {
		user.Roles = roles
	}

	permissions, err := uc.repo.GetUserPermissions(ctx, user.ID)
	if err == nil {
		user.Permissions = permissions
	}
	return user, nil
}

// AssignUserRole assigns a role to a user.
func (uc *AuthUseCase) AssignUserRole(ctx context.Context, userID, roleID, assignedBy string) error {
	if err := uc.repo.AssignRole(ctx, userID, roleID, assignedBy); err != nil {
		return fmt.Errorf("assign role: %w", err)
	}
	uc.log.Info("role assigned", "user_id", userID, "role_id", roleID)
	return nil
}

// RemoveUserRole removes a role from a user.
func (uc *AuthUseCase) RemoveUserRole(ctx context.Context, userID, roleID string) error {
	if err := uc.repo.RemoveRole(ctx, userID, roleID); err != nil {
		return fmt.Errorf("remove role: %w", err)
	}
	uc.log.Info("role removed", "user_id", userID, "role_id", roleID)
	return nil
}

// GetRoles returns all roles.
func (uc *AuthUseCase) GetRoles(ctx context.Context) ([]domain.Role, error) {
	return uc.repo.GetRoles(ctx)
}

// GetPermissions returns all permissions.
func (uc *AuthUseCase) GetPermissions(ctx context.Context) ([]domain.Permission, error) {
	return uc.repo.GetPermissions(ctx)
}

// GetRolePermissions returns permissions for a role.
func (uc *AuthUseCase) GetRolePermissions(ctx context.Context, roleID string) ([]domain.Permission, error) {
	return uc.repo.GetRolePermissions(ctx, roleID)
}

// GetAuditLogs returns filtered audit logs.
func (uc *AuthUseCase) GetAuditLogs(ctx context.Context, filter ports.AuditFilter) ([]domain.AuditLog, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 || filter.PerPage > 100 {
		filter.PerPage = 50
	}
	return uc.repo.ListAuditLogs(ctx, filter)
}

// CleanupExpiredSessions removes all expired sessions.
func (uc *AuthUseCase) CleanupExpiredSessions(ctx context.Context) error {
	return uc.repo.DeleteExpiredSessions(ctx)
}

// EnsureDefaultAdmin creates the default admin user if no users exist.
func (uc *AuthUseCase) EnsureDefaultAdmin(ctx context.Context) error {
	users, err := uc.repo.ListUsers(ctx)
	if err != nil {
		return fmt.Errorf("list users: %w", err)
	}

	if len(users) > 0 {
		return nil // Users already exist
	}

	// Create default admin
	hash, err := security.HashPassword("Admin@123")
	if err != nil {
		return fmt.Errorf("hash default password: %w", err)
	}

	now := time.Now()
	admin := &domain.User{
		ID:                 uuid.New().String(),
		Username:           "admin",
		Email:              "admin@salonflow.local",
		PasswordHash:       hash,
		DisplayName:        "Administrator",
		IsActive:           true,
		MustChangePassword: true,
		PasswordChangedAt:  now,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := uc.repo.CreateUser(ctx, admin); err != nil {
		return fmt.Errorf("create admin: %w", err)
	}

	if err := uc.repo.AssignRole(ctx, admin.ID, "role-owner", "system"); err != nil {
		return fmt.Errorf("assign owner role: %w", err)
	}

	uc.log.Info("default admin user created", "username", "admin")
	return nil
}

// audit is a helper to create audit log entries.
func (uc *AuthUseCase) audit(ctx context.Context, user *domain.User, action, module, entityID, description, severity string) {
	entry := &domain.AuditLog{
		ID:          uuid.New().String(),
		Timestamp:   time.Now(),
		UserID:      user.ID,
		Username:    user.Username,
		Action:      action,
		Module:      module,
		EntityType:  "user",
		EntityID:    entityID,
		Description: description,
		AppVersion:  uc.appVersion,
		Severity:    severity,
		CreatedAt:   time.Now(),
	}
	if err := uc.repo.CreateAuditLog(ctx, entry); err != nil {
		uc.log.Error("failed to create audit log", "error", err)
	}
}
