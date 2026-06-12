package ports

import (
	"context"

	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// AuthRepository handles user and session persistence.
type AuthRepository interface {
	// Users
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	ListUsers(ctx context.Context) ([]domain.User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdatePassword(ctx context.Context, userID string, passwordHash string) error
	IncrementFailedAttempts(ctx context.Context, userID string) error
	LockUser(ctx context.Context, userID string, until *string) error
	UnlockUser(ctx context.Context, userID string) error
	ResetFailedAttempts(ctx context.Context, userID string) error
	UpdateLastLogin(ctx context.Context, userID string) error

	// Sessions
	CreateSession(ctx context.Context, session *domain.Session) error
	GetSessionByToken(ctx context.Context, tokenHash string) (*domain.Session, error)
	DeleteSession(ctx context.Context, id string) error
	DeleteUserSessions(ctx context.Context, userID string) error
	DeleteExpiredSessions(ctx context.Context) error
	TouchSession(ctx context.Context, id string) error

	// Roles
	GetRoles(ctx context.Context) ([]domain.Role, error)
	GetRoleByID(ctx context.Context, id string) (*domain.Role, error)
	GetRoleByName(ctx context.Context, name string) (*domain.Role, error)
	CreateRole(ctx context.Context, role *domain.Role) error
	UpdateRole(ctx context.Context, role *domain.Role) error
	DeleteRole(ctx context.Context, id string) error

	// User-Role assignments
	AssignRole(ctx context.Context, userID, roleID, assignedBy string) error
	RemoveRole(ctx context.Context, userID, roleID string) error
	GetUserRoles(ctx context.Context, userID string) ([]domain.Role, error)

	// Permissions
	GetPermissions(ctx context.Context) ([]domain.Permission, error)
	GetRolePermissions(ctx context.Context, roleID string) ([]domain.Permission, error)
	GetUserPermissions(ctx context.Context, userID string) ([]string, error)
	AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error

	// Audit
	CreateAuditLog(ctx context.Context, entry *domain.AuditLog) error
	ListAuditLogs(ctx context.Context, filter AuditFilter) ([]domain.AuditLog, int, error)
}

// AuditFilter defines query parameters for audit log searches.
type AuditFilter struct {
	UserID     string `json:"user_id"`
	Module     string `json:"module"`
	Action     string `json:"action"`
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
	Severity   string `json:"severity"`
	FromDate   string `json:"from_date"`
	ToDate     string `json:"to_date"`
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
}
