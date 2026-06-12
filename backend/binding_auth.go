package main

import (
	"context"

	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// AuthService exposes authentication operations to the Wails frontend.
type AuthService struct {
	ctx   context.Context
	uc    *usecase.AuthUseCase
	audit *usecase.AuditService
	token string // Current session token stored in memory
}

func NewAuthService(uc *usecase.AuthUseCase, audit *usecase.AuditService) *AuthService {
	return &AuthService{uc: uc, audit: audit}
}

func (s *AuthService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

// Login authenticates a user and stores the session token.
func (s *AuthService) Login(input usecase.LoginInput) (*usecase.LoginOutput, error) {
	result, err := s.uc.Login(s.ctx, input)
	if err != nil {
		return nil, err
	}
	s.token = result.Token
	return result, nil
}

// Logout invalidates the current session.
func (s *AuthService) Logout() error {
	if s.token == "" {
		return nil
	}
	err := s.uc.Logout(s.ctx, s.token)
	s.token = ""
	return err
}

// GetCurrentSession returns the current session info if valid.
func (s *AuthService) GetCurrentSession() (*domain.SessionInfo, error) {
	if s.token == "" {
		return nil, usecase.ErrSessionNotFound
	}
	return s.uc.ValidateSession(s.ctx, s.token)
}

// CheckPermission checks if the current user has a specific permission.
func (s *AuthService) CheckPermission(permission string) error {
	if s.token == "" {
		return usecase.ErrSessionNotFound
	}
	return s.uc.CheckPermission(s.ctx, s.token, permission)
}

// HasPermission returns true if the current user has the permission.
func (s *AuthService) HasPermission(permission string) bool {
	return s.CheckPermission(permission) == nil
}

// HasAnyPermission returns true if the current user has any of the given permissions.
func (s *AuthService) HasAnyPermission(permissions []string) bool {
	for _, p := range permissions {
		if s.HasPermission(p) {
			return true
		}
	}
	return false
}

// ChangePassword changes the current user's password.
func (s *AuthService) ChangePassword(oldPassword, newPassword string) error {
	session, err := s.GetCurrentSession()
	if err != nil {
		return err
	}
	return s.uc.ChangePassword(s.ctx, usecase.ChangePasswordInput{
		UserID:      session.UserID,
		OldPassword: oldPassword,
		NewPassword: newPassword,
	})
}

// IsAuthenticated returns whether there is a valid session.
func (s *AuthService) IsAuthenticated() bool {
	_, err := s.GetCurrentSession()
	return err == nil
}

// SetToken sets the session token (for Remember Me functionality).
func (s *AuthService) SetToken(token string) {
	s.token = token
}

// GetToken returns the current token (for storing in localStorage).
func (s *AuthService) GetToken() string {
	return s.token
}

// ============================================================
// USER MANAGEMENT (requires users.manage permission)
// ============================================================

// CreateUser creates a new user.
func (s *AuthService) CreateUser(input usecase.CreateUserInput) (*domain.User, error) {
	if err := s.requirePermission("users.manage"); err != nil {
		return nil, err
	}
	session, _ := s.GetCurrentSession()
	createdBy := ""
	if session != nil {
		createdBy = session.UserID
	}
	return s.uc.CreateUser(s.ctx, input, createdBy)
}

// UpdateUser updates a user.
func (s *AuthService) UpdateUser(input usecase.UpdateUserInput) (*domain.User, error) {
	if err := s.requirePermission("users.manage"); err != nil {
		return nil, err
	}
	return s.uc.UpdateUser(s.ctx, input)
}

// DeleteUser deactivates a user.
func (s *AuthService) DeleteUser(userID string) error {
	if err := s.requirePermission("users.manage"); err != nil {
		return err
	}
	return s.uc.DeleteUser(s.ctx, userID)
}

// GetUsers returns all users.
func (s *AuthService) GetUsers() ([]domain.User, error) {
	if err := s.requirePermission("users.view"); err != nil {
		return nil, err
	}
	return s.uc.GetUsers(s.ctx)
}

// GetUser returns a single user.
func (s *AuthService) GetUser(id string) (*domain.User, error) {
	if err := s.requirePermission("users.view"); err != nil {
		return nil, err
	}
	return s.uc.GetUser(s.ctx, id)
}

// ResetUserPassword resets a user's password (admin action).
func (s *AuthService) ResetUserPassword(userID, newPassword string) error {
	if err := s.requirePermission("users.manage"); err != nil {
		return err
	}
	return s.uc.ResetPassword(s.ctx, usecase.ResetPasswordInput{
		UserID:      userID,
		NewPassword: newPassword,
	})
}

// AssignUserRole assigns a role to a user.
func (s *AuthService) AssignUserRole(userID, roleID string) error {
	if err := s.requirePermission("users.roles"); err != nil {
		return err
	}
	session, _ := s.GetCurrentSession()
	assignedBy := ""
	if session != nil {
		assignedBy = session.UserID
	}
	return s.uc.AssignUserRole(s.ctx, userID, roleID, assignedBy)
}

// RemoveUserRole removes a role from a user.
func (s *AuthService) RemoveUserRole(userID, roleID string) error {
	if err := s.requirePermission("users.roles"); err != nil {
		return err
	}
	return s.uc.RemoveUserRole(s.ctx, userID, roleID)
}

// ============================================================
// ROLES & PERMISSIONS
// ============================================================

// GetRoles returns all roles.
func (s *AuthService) GetRoles() ([]domain.Role, error) {
	return s.uc.GetRoles(s.ctx)
}

// GetPermissions returns all permissions.
func (s *AuthService) GetPermissions() ([]domain.Permission, error) {
	return s.uc.GetPermissions(s.ctx)
}

// GetRolePermissions returns permissions for a specific role.
func (s *AuthService) GetRolePermissions(roleID string) ([]domain.Permission, error) {
	return s.uc.GetRolePermissions(s.ctx, roleID)
}

// ============================================================
// AUDIT LOGS
// ============================================================

// AuditLogListOutput holds paginated audit log results.
type AuditLogListOutput struct {
	Logs  []domain.AuditLog `json:"logs"`
	Total int               `json:"total"`
	Page  int               `json:"page"`
}

// GetAuditLogs returns filtered audit logs.
func (s *AuthService) GetAuditLogs(filter ports.AuditFilter) (*AuditLogListOutput, error) {
	if err := s.requirePermission("audit.view"); err != nil {
		return nil, err
	}
	logs, total, err := s.uc.GetAuditLogs(s.ctx, filter)
	if err != nil {
		return nil, err
	}
	return &AuditLogListOutput{
		Logs:  logs,
		Total: total,
		Page:  filter.Page,
	}, nil
}

// requirePermission checks the current session has the required permission.
func (s *AuthService) requirePermission(permission string) error {
	if s.token == "" {
		return usecase.ErrSessionNotFound
	}
	return s.uc.CheckPermission(s.ctx, s.token, permission)
}
