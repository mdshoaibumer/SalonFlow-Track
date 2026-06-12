package usecase

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/internal/core/security"
)

// mockAuthRepo is a test double for AuthRepository
type mockAuthRepo struct {
	users    map[string]*domain.User
	sessions map[string]*domain.Session
	roles    []domain.Role
	perms    []domain.Permission
	rolePerm map[string][]string // roleID -> permCodes
	userRole map[string][]string // userID -> roleIDs
	audits   []domain.AuditLog
}

func newMockAuthRepo() *mockAuthRepo {
	return &mockAuthRepo{
		users:    make(map[string]*domain.User),
		sessions: make(map[string]*domain.Session),
		roles: []domain.Role{
			{ID: "role-owner", Name: "owner", IsSystem: true},
			{ID: "role-manager", Name: "manager", IsSystem: true},
		},
		perms: []domain.Permission{
			{ID: "p1", Code: "customers.read", Module: "customers", Action: "read"},
			{ID: "p2", Code: "customers.create", Module: "customers", Action: "create"},
			{ID: "p3", Code: "billing.view", Module: "billing", Action: "view"},
		},
		rolePerm: map[string][]string{
			"role-owner":   {"customers.read", "customers.create", "billing.view"},
			"role-manager": {"customers.read", "billing.view"},
		},
		userRole: make(map[string][]string),
	}
}

func (m *mockAuthRepo) CreateUser(_ context.Context, user *domain.User) error {
	m.users[user.ID] = user
	return nil
}
func (m *mockAuthRepo) GetUserByID(_ context.Context, id string) (*domain.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	copy := *u
	return &copy, nil
}
func (m *mockAuthRepo) GetUserByUsername(_ context.Context, username string) (*domain.User, error) {
	for _, u := range m.users {
		if u.Username == username {
			copy := *u
			return &copy, nil
		}
	}
	return nil, ErrUserNotFound
}
func (m *mockAuthRepo) GetUserByEmail(_ context.Context, email string) (*domain.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			copy := *u
			return &copy, nil
		}
	}
	return nil, ErrUserNotFound
}
func (m *mockAuthRepo) UpdateUser(_ context.Context, user *domain.User) error {
	m.users[user.ID] = user
	return nil
}
func (m *mockAuthRepo) ListUsers(_ context.Context) ([]domain.User, error) {
	var result []domain.User
	for _, u := range m.users {
		result = append(result, *u)
	}
	return result, nil
}
func (m *mockAuthRepo) DeleteUser(_ context.Context, id string) error {
	delete(m.users, id)
	return nil
}
func (m *mockAuthRepo) UpdatePassword(_ context.Context, userID string, hash string) error {
	if u, ok := m.users[userID]; ok {
		u.PasswordHash = hash
	}
	return nil
}
func (m *mockAuthRepo) IncrementFailedAttempts(_ context.Context, userID string) error {
	if u, ok := m.users[userID]; ok {
		u.FailedAttempts++
	}
	return nil
}
func (m *mockAuthRepo) LockUser(_ context.Context, userID string, until *string) error {
	if u, ok := m.users[userID]; ok {
		u.IsLocked = true
		if until != nil {
			t, _ := time.Parse(time.RFC3339, *until)
			u.LockedUntil = &t
		}
	}
	return nil
}
func (m *mockAuthRepo) UnlockUser(_ context.Context, userID string) error {
	if u, ok := m.users[userID]; ok {
		u.IsLocked = false
	}
	return nil
}
func (m *mockAuthRepo) ResetFailedAttempts(_ context.Context, userID string) error {
	if u, ok := m.users[userID]; ok {
		u.FailedAttempts = 0
	}
	return nil
}
func (m *mockAuthRepo) UpdateLastLogin(_ context.Context, _ string) error { return nil }

func (m *mockAuthRepo) CreateSession(_ context.Context, session *domain.Session) error {
	m.sessions[session.TokenHash] = session
	return nil
}
func (m *mockAuthRepo) GetSessionByToken(_ context.Context, tokenHash string) (*domain.Session, error) {
	s, ok := m.sessions[tokenHash]
	if !ok {
		return nil, ErrSessionNotFound
	}
	return s, nil
}
func (m *mockAuthRepo) DeleteSession(_ context.Context, id string) error {
	for k, s := range m.sessions {
		if s.ID == id {
			delete(m.sessions, k)
		}
	}
	return nil
}
func (m *mockAuthRepo) DeleteUserSessions(_ context.Context, userID string) error {
	for k, s := range m.sessions {
		if s.UserID == userID {
			delete(m.sessions, k)
		}
	}
	return nil
}
func (m *mockAuthRepo) DeleteExpiredSessions(_ context.Context) error  { return nil }
func (m *mockAuthRepo) TouchSession(_ context.Context, _ string) error { return nil }

func (m *mockAuthRepo) GetRoles(_ context.Context) ([]domain.Role, error) { return m.roles, nil }
func (m *mockAuthRepo) GetRoleByID(_ context.Context, id string) (*domain.Role, error) {
	for _, r := range m.roles {
		if r.ID == id {
			return &r, nil
		}
	}
	return nil, nil
}
func (m *mockAuthRepo) GetRoleByName(_ context.Context, name string) (*domain.Role, error) {
	for _, r := range m.roles {
		if r.Name == name {
			return &r, nil
		}
	}
	return nil, nil
}
func (m *mockAuthRepo) CreateRole(_ context.Context, role *domain.Role) error {
	m.roles = append(m.roles, *role)
	return nil
}
func (m *mockAuthRepo) UpdateRole(_ context.Context, _ *domain.Role) error { return nil }
func (m *mockAuthRepo) DeleteRole(_ context.Context, _ string) error       { return nil }

func (m *mockAuthRepo) AssignRole(_ context.Context, userID, roleID, _ string) error {
	m.userRole[userID] = append(m.userRole[userID], roleID)
	return nil
}
func (m *mockAuthRepo) RemoveRole(_ context.Context, userID, roleID string) error {
	roles := m.userRole[userID]
	for i, r := range roles {
		if r == roleID {
			m.userRole[userID] = append(roles[:i], roles[i+1:]...)
			break
		}
	}
	return nil
}
func (m *mockAuthRepo) GetUserRoles(_ context.Context, userID string) ([]domain.Role, error) {
	var result []domain.Role
	for _, rid := range m.userRole[userID] {
		for _, r := range m.roles {
			if r.ID == rid {
				result = append(result, r)
			}
		}
	}
	return result, nil
}

func (m *mockAuthRepo) GetPermissions(_ context.Context) ([]domain.Permission, error) {
	return m.perms, nil
}
func (m *mockAuthRepo) GetRolePermissions(_ context.Context, roleID string) ([]domain.Permission, error) {
	codes := m.rolePerm[roleID]
	var result []domain.Permission
	for _, c := range codes {
		for _, p := range m.perms {
			if p.Code == c {
				result = append(result, p)
			}
		}
	}
	return result, nil
}
func (m *mockAuthRepo) GetUserPermissions(_ context.Context, userID string) ([]string, error) {
	var perms []string
	seen := make(map[string]bool)
	for _, rid := range m.userRole[userID] {
		for _, code := range m.rolePerm[rid] {
			if !seen[code] {
				perms = append(perms, code)
				seen[code] = true
			}
		}
	}
	return perms, nil
}
func (m *mockAuthRepo) AssignPermissionToRole(_ context.Context, roleID, permID string) error {
	for _, p := range m.perms {
		if p.ID == permID {
			m.rolePerm[roleID] = append(m.rolePerm[roleID], p.Code)
		}
	}
	return nil
}
func (m *mockAuthRepo) RemovePermissionFromRole(_ context.Context, _, _ string) error { return nil }

func (m *mockAuthRepo) CreateAuditLog(_ context.Context, entry *domain.AuditLog) error {
	m.audits = append(m.audits, *entry)
	return nil
}
func (m *mockAuthRepo) ListAuditLogs(_ context.Context, filter ports.AuditFilter) ([]domain.AuditLog, int, error) {
	return m.audits, len(m.audits), nil
}

// ===========
// TESTS
// ===========

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

func TestAuthUseCase_CreateUser(t *testing.T) {
	repo := newMockAuthRepo()
	uc := NewAuthUseCase(repo, testLogger(), "1.0.0")

	user, err := uc.CreateUser(context.Background(), CreateUserInput{
		Username:    "testuser",
		Email:       "test@example.com",
		Password:    "Test@1234",
		DisplayName: "Test User",
		Phone:       "1234567890",
		RoleID:      "role-owner",
	}, "system")

	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}
	if user.Username != "testuser" {
		t.Errorf("expected username 'testuser', got %q", user.Username)
	}
	if user.PasswordHash == "" {
		t.Error("password hash should not be empty")
	}
}

func TestAuthUseCase_CreateUserDuplicate(t *testing.T) {
	repo := newMockAuthRepo()
	uc := NewAuthUseCase(repo, testLogger(), "1.0.0")

	_, _ = uc.CreateUser(context.Background(), CreateUserInput{
		Username: "admin", Password: "Test@1234", DisplayName: "Admin",
	}, "system")

	_, err := uc.CreateUser(context.Background(), CreateUserInput{
		Username: "admin", Password: "Test@1234", DisplayName: "Admin 2",
	}, "system")

	if err != ErrUserExists {
		t.Errorf("expected ErrUserExists, got %v", err)
	}
}

func TestAuthUseCase_LoginSuccess(t *testing.T) {
	repo := newMockAuthRepo()
	uc := NewAuthUseCase(repo, testLogger(), "1.0.0")

	// Create user
	_, _ = uc.CreateUser(context.Background(), CreateUserInput{
		Username: "loginuser", Password: "Test@1234", DisplayName: "Login User", RoleID: "role-owner",
	}, "system")

	// Login
	output, err := uc.Login(context.Background(), LoginInput{
		Username: "loginuser", Password: "Test@1234",
	})
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if output.Token == "" {
		t.Error("expected non-empty token")
	}
	if output.User.Username != "loginuser" {
		t.Errorf("expected username 'loginuser', got %q", output.User.Username)
	}
	if len(output.User.Permissions) == 0 {
		t.Error("expected permissions to be populated")
	}
}

func TestAuthUseCase_LoginFailure(t *testing.T) {
	repo := newMockAuthRepo()
	uc := NewAuthUseCase(repo, testLogger(), "1.0.0")

	_, _ = uc.CreateUser(context.Background(), CreateUserInput{
		Username: "user1", Password: "Test@1234", DisplayName: "User 1",
	}, "system")

	_, err := uc.Login(context.Background(), LoginInput{
		Username: "user1", Password: "WrongPassword1!",
	})
	if err != ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthUseCase_AccountLocking(t *testing.T) {
	repo := newMockAuthRepo()
	uc := NewAuthUseCase(repo, testLogger(), "1.0.0")

	_, _ = uc.CreateUser(context.Background(), CreateUserInput{
		Username: "locktest", Password: "Test@1234", DisplayName: "Lock Test",
	}, "system")

	// Fail 5 times to trigger lock
	for i := 0; i < 5; i++ {
		_, err := uc.Login(context.Background(), LoginInput{Username: "locktest", Password: "Wrong1!@"})
		if i < 4 {
			// First 4 attempts should return invalid credentials
			if err != ErrInvalidCredentials {
				t.Fatalf("attempt %d: expected ErrInvalidCredentials, got %v", i+1, err)
			}
		} else {
			// 5th attempt should lock the account
			if err != ErrAccountLocked {
				t.Fatalf("attempt %d: expected ErrAccountLocked, got %v", i+1, err)
			}
		}
	}

	// Even with correct password, account should be locked
	_, err := uc.Login(context.Background(), LoginInput{Username: "locktest", Password: "Test@1234"})
	if err != ErrAccountLocked {
		t.Errorf("expected ErrAccountLocked after lockout, got %v", err)
	}
}

func TestAuthUseCase_ValidateSession(t *testing.T) {
	repo := newMockAuthRepo()
	uc := NewAuthUseCase(repo, testLogger(), "1.0.0")

	_, _ = uc.CreateUser(context.Background(), CreateUserInput{
		Username: "sessuser", Password: "Test@1234", DisplayName: "Session User", RoleID: "role-owner",
	}, "system")

	output, _ := uc.Login(context.Background(), LoginInput{
		Username: "sessuser", Password: "Test@1234",
	})

	info, err := uc.ValidateSession(context.Background(), output.Token)
	if err != nil {
		t.Fatalf("ValidateSession() error = %v", err)
	}
	if info.Username != "sessuser" {
		t.Errorf("expected username 'sessuser', got %q", info.Username)
	}
}

func TestAuthUseCase_CheckPermission(t *testing.T) {
	repo := newMockAuthRepo()
	uc := NewAuthUseCase(repo, testLogger(), "1.0.0")

	_, _ = uc.CreateUser(context.Background(), CreateUserInput{
		Username: "permuser", Password: "Test@1234", DisplayName: "Perm User", RoleID: "role-manager",
	}, "system")

	output, _ := uc.Login(context.Background(), LoginInput{
		Username: "permuser", Password: "Test@1234",
	})

	// Manager has customers.read
	err := uc.CheckPermission(context.Background(), output.Token, "customers.read")
	if err != nil {
		t.Errorf("expected permission granted for customers.read, got %v", err)
	}

	// Manager does NOT have customers.create
	err = uc.CheckPermission(context.Background(), output.Token, "customers.create")
	if err != ErrPermissionDenied {
		t.Errorf("expected ErrPermissionDenied for customers.create, got %v", err)
	}
}

func TestAuthUseCase_ChangePassword(t *testing.T) {
	repo := newMockAuthRepo()
	uc := NewAuthUseCase(repo, testLogger(), "1.0.0")

	user, _ := uc.CreateUser(context.Background(), CreateUserInput{
		Username: "pwuser", Password: "Test@1234", DisplayName: "PW User",
	}, "system")

	// Change password
	err := uc.ChangePassword(context.Background(), ChangePasswordInput{
		UserID:      user.ID,
		OldPassword: "Test@1234",
		NewPassword: "NewP@ss123",
	})
	if err != nil {
		t.Fatalf("ChangePassword() error = %v", err)
	}

	// Verify old password no longer works
	valid, _ := security.VerifyPassword("Test@1234", repo.users[user.ID].PasswordHash)
	if valid {
		t.Error("old password should not work after change")
	}

	// Verify new password works
	valid, _ = security.VerifyPassword("NewP@ss123", repo.users[user.ID].PasswordHash)
	if !valid {
		t.Error("new password should work after change")
	}
}

func TestAuthUseCase_ChangePasswordWrongOld(t *testing.T) {
	repo := newMockAuthRepo()
	uc := NewAuthUseCase(repo, testLogger(), "1.0.0")

	user, _ := uc.CreateUser(context.Background(), CreateUserInput{
		Username: "pwuser2", Password: "Test@1234", DisplayName: "PW User 2",
	}, "system")

	err := uc.ChangePassword(context.Background(), ChangePasswordInput{
		UserID:      user.ID,
		OldPassword: "Wrong@123",
		NewPassword: "NewP@ss123",
	})
	if err != ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthUseCase_WeakPasswordRejected(t *testing.T) {
	repo := newMockAuthRepo()
	uc := NewAuthUseCase(repo, testLogger(), "1.0.0")

	_, err := uc.CreateUser(context.Background(), CreateUserInput{
		Username: "weakuser", Password: "weak", DisplayName: "Weak User",
	}, "system")
	if err == nil {
		t.Error("expected error for weak password")
	}
}

func TestAuthUseCase_Logout(t *testing.T) {
	repo := newMockAuthRepo()
	uc := NewAuthUseCase(repo, testLogger(), "1.0.0")

	_, _ = uc.CreateUser(context.Background(), CreateUserInput{
		Username: "logoutuser", Password: "Test@1234", DisplayName: "Logout User", RoleID: "role-owner",
	}, "system")

	output, _ := uc.Login(context.Background(), LoginInput{
		Username: "logoutuser", Password: "Test@1234",
	})

	err := uc.Logout(context.Background(), output.Token)
	if err != nil {
		t.Fatalf("Logout() error = %v", err)
	}

	// Session should be invalidated
	_, err = uc.ValidateSession(context.Background(), output.Token)
	if err != ErrSessionNotFound {
		t.Errorf("expected ErrSessionNotFound after logout, got %v", err)
	}
}

func TestAuthUseCase_EnsureDefaultAdmin(t *testing.T) {
	repo := newMockAuthRepo()
	uc := NewAuthUseCase(repo, testLogger(), "1.0.0")

	err := uc.EnsureDefaultAdmin(context.Background())
	if err != nil {
		t.Fatalf("EnsureDefaultAdmin() error = %v", err)
	}

	// Should be able to login as admin
	_, err = uc.Login(context.Background(), LoginInput{
		Username: "admin", Password: "Admin@123",
	})
	if err != nil {
		t.Fatalf("Login as default admin failed: %v", err)
	}

	// Calling again should be no-op
	err = uc.EnsureDefaultAdmin(context.Background())
	if err != nil {
		t.Fatalf("EnsureDefaultAdmin() second call error = %v", err)
	}
}

func TestAuthUseCase_AuditLogging(t *testing.T) {
	repo := newMockAuthRepo()
	uc := NewAuthUseCase(repo, testLogger(), "1.0.0")

	_, _ = uc.CreateUser(context.Background(), CreateUserInput{
		Username: "audituser", Password: "Test@1234", DisplayName: "Audit User", RoleID: "role-owner",
	}, "system")

	_, _ = uc.Login(context.Background(), LoginInput{
		Username: "audituser", Password: "Test@1234",
	})

	// Should have audit entries for create and login
	if len(repo.audits) < 2 {
		t.Errorf("expected at least 2 audit entries, got %d", len(repo.audits))
	}
}
