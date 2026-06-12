package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
)

// AuthRepository is the SQLite implementation of ports.AuthRepository.
type AuthRepository struct {
	db  *sql.DB
	log *slog.Logger
}

// NewAuthRepository creates a new AuthRepository.
func NewAuthRepository(db *sql.DB, log *slog.Logger) *AuthRepository {
	return &AuthRepository{db: db, log: log}
}

// ============================================================
// USERS
// ============================================================

func (r *AuthRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, display_name, phone, is_active, is_locked,
			failed_attempts, locked_until, last_login_at, password_changed_at, must_change_password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	var lockedUntil, lastLogin *string
	if user.LockedUntil != nil {
		t := user.LockedUntil.Format(time.RFC3339)
		lockedUntil = &t
	}
	if user.LastLoginAt != nil {
		t := user.LastLoginAt.Format(time.RFC3339)
		lastLogin = &t
	}

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.DisplayName,
		user.Phone,
		user.IsActive,
		user.IsLocked,
		user.FailedAttempts,
		lockedUntil,
		lastLogin,
		user.PasswordChangedAt.Format(time.RFC3339),
		user.MustChangePassword,
		user.CreatedAt.Format(time.RFC3339),
		user.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *AuthRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return r.getUser(ctx, "id = ?", id)
}

func (r *AuthRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	return r.getUser(ctx, "username = ?", username)
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.getUser(ctx, "email = ?", email)
}

func (r *AuthRepository) getUser(ctx context.Context, where string, arg interface{}) (*domain.User, error) {
	query := fmt.Sprintf(`
		SELECT id, username, email, password_hash, display_name, phone, is_active, is_locked,
			failed_attempts, locked_until, last_login_at, password_changed_at, must_change_password, created_at, updated_at
		FROM users WHERE %s`, where)

	var user domain.User
	var lockedUntil, lastLogin sql.NullString
	var email sql.NullString
	var pwChangedAt string

	err := r.db.QueryRowContext(ctx, query, arg).Scan(
		&user.ID,
		&user.Username,
		&email,
		&user.PasswordHash,
		&user.DisplayName,
		&user.Phone,
		&user.IsActive,
		&user.IsLocked,
		&user.FailedAttempts,
		&lockedUntil,
		&lastLogin,
		&pwChangedAt,
		&user.MustChangePassword,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	if email.Valid {
		user.Email = email.String
	}
	if lockedUntil.Valid {
		t, _ := time.Parse(time.RFC3339, lockedUntil.String)
		user.LockedUntil = &t
	}
	if lastLogin.Valid {
		t, _ := time.Parse(time.RFC3339, lastLogin.String)
		user.LastLoginAt = &t
	}
	user.PasswordChangedAt, _ = time.Parse(time.RFC3339, pwChangedAt)

	return &user, nil
}

func (r *AuthRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users SET email = ?, display_name = ?, phone = ?, is_active = ?, is_locked = ?,
			must_change_password = ?, updated_at = ?
		WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.DisplayName,
		user.Phone,
		user.IsActive,
		user.IsLocked,
		user.MustChangePassword,
		time.Now().Format(time.RFC3339),
		user.ID,
	)
	return err
}

func (r *AuthRepository) ListUsers(ctx context.Context) ([]domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, display_name, phone, is_active, is_locked,
			failed_attempts, locked_until, last_login_at, password_changed_at, must_change_password, created_at, updated_at
		FROM users ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		var lockedUntil, lastLogin, email sql.NullString
		var pwChangedAt string

		if err := rows.Scan(
			&user.ID, &user.Username, &email, &user.PasswordHash, &user.DisplayName,
			&user.Phone, &user.IsActive, &user.IsLocked, &user.FailedAttempts,
			&lockedUntil, &lastLogin, &pwChangedAt, &user.MustChangePassword,
			&user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}

		if email.Valid {
			user.Email = email.String
		}
		if lockedUntil.Valid {
			t, _ := time.Parse(time.RFC3339, lockedUntil.String)
			user.LockedUntil = &t
		}
		if lastLogin.Valid {
			t, _ := time.Parse(time.RFC3339, lastLogin.String)
			user.LastLoginAt = &t
		}
		user.PasswordChangedAt, _ = time.Parse(time.RFC3339, pwChangedAt)
		users = append(users, user)
	}
	return users, nil
}

func (r *AuthRepository) DeleteUser(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = ?", id)
	return err
}

func (r *AuthRepository) UpdatePassword(ctx context.Context, userID string, passwordHash string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE users SET password_hash = ?, password_changed_at = ?, must_change_password = 0, updated_at = ? WHERE id = ?",
		passwordHash, time.Now().Format(time.RFC3339), time.Now().Format(time.RFC3339), userID)
	return err
}

func (r *AuthRepository) IncrementFailedAttempts(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE users SET failed_attempts = failed_attempts + 1 WHERE id = ?", userID)
	return err
}

func (r *AuthRepository) LockUser(ctx context.Context, userID string, until *string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE users SET is_locked = 1, locked_until = ? WHERE id = ?", until, userID)
	return err
}

func (r *AuthRepository) UnlockUser(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE users SET is_locked = 0, locked_until = NULL WHERE id = ?", userID)
	return err
}

func (r *AuthRepository) ResetFailedAttempts(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE users SET failed_attempts = 0 WHERE id = ?", userID)
	return err
}

func (r *AuthRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE users SET last_login_at = ? WHERE id = ?", time.Now().Format(time.RFC3339), userID)
	return err
}

// ============================================================
// SESSIONS
// ============================================================

func (r *AuthRepository) CreateSession(ctx context.Context, session *domain.Session) error {
	query := `
		INSERT INTO sessions (id, user_id, token_hash, device_id, ip_address, user_agent, remember_me, expires_at, last_active_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		session.ID,
		session.UserID,
		session.TokenHash,
		session.DeviceID,
		session.IPAddress,
		session.UserAgent,
		session.RememberMe,
		session.ExpiresAt.Format(time.RFC3339),
		session.LastActiveAt.Format(time.RFC3339),
		session.CreatedAt.Format(time.RFC3339),
	)
	return err
}

func (r *AuthRepository) GetSessionByToken(ctx context.Context, tokenHash string) (*domain.Session, error) {
	query := `SELECT id, user_id, token_hash, device_id, ip_address, user_agent, remember_me, expires_at, last_active_at, created_at
		FROM sessions WHERE token_hash = ?`

	var s domain.Session
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&s.ID, &s.UserID, &s.TokenHash, &s.DeviceID, &s.IPAddress,
		&s.UserAgent, &s.RememberMe, &s.ExpiresAt, &s.LastActiveAt, &s.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, err
	}
	return &s, nil
}

func (r *AuthRepository) DeleteSession(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM sessions WHERE id = ?", id)
	return err
}

func (r *AuthRepository) DeleteUserSessions(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM sessions WHERE user_id = ?", userID)
	return err
}

func (r *AuthRepository) DeleteExpiredSessions(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM sessions WHERE expires_at < ?", time.Now().Format(time.RFC3339))
	return err
}

func (r *AuthRepository) TouchSession(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE sessions SET last_active_at = ? WHERE id = ?",
		time.Now().Format(time.RFC3339), id)
	return err
}

// ============================================================
// ROLES
// ============================================================

func (r *AuthRepository) GetRoles(ctx context.Context) ([]domain.Role, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, description, is_system, created_at, updated_at FROM roles ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		var role domain.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.IsSystem, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *AuthRepository) GetRoleByID(ctx context.Context, id string) (*domain.Role, error) {
	var role domain.Role
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name, description, is_system, created_at, updated_at FROM roles WHERE id = ?", id).
		Scan(&role.ID, &role.Name, &role.Description, &role.IsSystem, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *AuthRepository) GetRoleByName(ctx context.Context, name string) (*domain.Role, error) {
	var role domain.Role
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name, description, is_system, created_at, updated_at FROM roles WHERE name = ?", name).
		Scan(&role.ID, &role.Name, &role.Description, &role.IsSystem, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *AuthRepository) CreateRole(ctx context.Context, role *domain.Role) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO roles (id, name, description, is_system, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		role.ID, role.Name, role.Description, role.IsSystem, role.CreatedAt.Format(time.RFC3339), role.UpdatedAt.Format(time.RFC3339))
	return err
}

func (r *AuthRepository) UpdateRole(ctx context.Context, role *domain.Role) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE roles SET name = ?, description = ?, updated_at = ? WHERE id = ?",
		role.Name, role.Description, time.Now().Format(time.RFC3339), role.ID)
	return err
}

func (r *AuthRepository) DeleteRole(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM roles WHERE id = ? AND is_system = 0", id)
	return err
}

// ============================================================
// USER-ROLE ASSIGNMENTS
// ============================================================

func (r *AuthRepository) AssignRole(ctx context.Context, userID, roleID, assignedBy string) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT OR IGNORE INTO user_roles (user_id, role_id, assigned_by) VALUES (?, ?, ?)",
		userID, roleID, assignedBy)
	return err
}

func (r *AuthRepository) RemoveRole(ctx context.Context, userID, roleID string) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM user_roles WHERE user_id = ? AND role_id = ?", userID, roleID)
	return err
}

func (r *AuthRepository) GetUserRoles(ctx context.Context, userID string) ([]domain.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.is_system, r.created_at, r.updated_at
		FROM roles r
		INNER JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = ?`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		var role domain.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.IsSystem, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

// ============================================================
// PERMISSIONS
// ============================================================

func (r *AuthRepository) GetPermissions(ctx context.Context) ([]domain.Permission, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, code, module, action, description, created_at FROM permissions ORDER BY module, action")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []domain.Permission
	for rows.Next() {
		var p domain.Permission
		if err := rows.Scan(&p.ID, &p.Code, &p.Module, &p.Action, &p.Description, &p.CreatedAt); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	return perms, nil
}

func (r *AuthRepository) GetRolePermissions(ctx context.Context, roleID string) ([]domain.Permission, error) {
	query := `
		SELECT p.id, p.code, p.module, p.action, p.description, p.created_at
		FROM permissions p
		INNER JOIN role_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id = ?
		ORDER BY p.module, p.action`

	rows, err := r.db.QueryContext(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []domain.Permission
	for rows.Next() {
		var p domain.Permission
		if err := rows.Scan(&p.ID, &p.Code, &p.Module, &p.Action, &p.Description, &p.CreatedAt); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	return perms, nil
}

func (r *AuthRepository) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	query := `
		SELECT DISTINCT p.code
		FROM permissions p
		INNER JOIN role_permissions rp ON rp.permission_id = p.id
		INNER JOIN user_roles ur ON ur.role_id = rp.role_id
		WHERE ur.user_id = ?
		ORDER BY p.code`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		perms = append(perms, code)
	}
	return perms, nil
}

func (r *AuthRepository) AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT OR IGNORE INTO role_permissions (role_id, permission_id) VALUES (?, ?)",
		roleID, permissionID)
	return err
}

func (r *AuthRepository) RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM role_permissions WHERE role_id = ? AND permission_id = ?", roleID, permissionID)
	return err
}

// ============================================================
// AUDIT LOGS
// ============================================================

func (r *AuthRepository) CreateAuditLog(ctx context.Context, entry *domain.AuditLog) error {
	query := `
		INSERT INTO audit_logs (id, timestamp, user_id, username, action, module, entity_type, entity_id,
			description, old_value, new_value, device_id, ip_address, user_agent, app_version, severity)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		entry.ID,
		entry.Timestamp.Format(time.RFC3339),
		entry.UserID,
		entry.Username,
		entry.Action,
		entry.Module,
		entry.EntityType,
		entry.EntityID,
		entry.Description,
		entry.OldValue,
		entry.NewValue,
		entry.DeviceID,
		entry.IPAddress,
		entry.UserAgent,
		entry.AppVersion,
		entry.Severity,
	)
	return err
}

func (r *AuthRepository) ListAuditLogs(ctx context.Context, filter ports.AuditFilter) ([]domain.AuditLog, int, error) {
	var conditions []string
	var args []interface{}

	if filter.UserID != "" {
		conditions = append(conditions, "user_id = ?")
		args = append(args, filter.UserID)
	}
	if filter.Module != "" {
		conditions = append(conditions, "module = ?")
		args = append(args, filter.Module)
	}
	if filter.Action != "" {
		conditions = append(conditions, "action = ?")
		args = append(args, filter.Action)
	}
	if filter.EntityType != "" {
		conditions = append(conditions, "entity_type = ?")
		args = append(args, filter.EntityType)
	}
	if filter.EntityID != "" {
		conditions = append(conditions, "entity_id = ?")
		args = append(args, filter.EntityID)
	}
	if filter.Severity != "" {
		conditions = append(conditions, "severity = ?")
		args = append(args, filter.Severity)
	}
	if filter.FromDate != "" {
		conditions = append(conditions, "timestamp >= ?")
		args = append(args, filter.FromDate)
	}
	if filter.ToDate != "" {
		conditions = append(conditions, "timestamp <= ?")
		args = append(args, filter.ToDate)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM audit_logs %s", whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Paginated query
	offset := (filter.Page - 1) * filter.PerPage
	query := fmt.Sprintf(`
		SELECT id, timestamp, user_id, username, action, module, entity_type, entity_id,
			description, old_value, new_value, device_id, ip_address, user_agent, app_version, severity, created_at
		FROM audit_logs %s
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?`, whereClause)

	args = append(args, filter.PerPage, offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []domain.AuditLog
	for rows.Next() {
		var entry domain.AuditLog
		var oldVal, newVal sql.NullString
		if err := rows.Scan(
			&entry.ID, &entry.Timestamp, &entry.UserID, &entry.Username,
			&entry.Action, &entry.Module, &entry.EntityType, &entry.EntityID,
			&entry.Description, &oldVal, &newVal,
			&entry.DeviceID, &entry.IPAddress, &entry.UserAgent, &entry.AppVersion,
			&entry.Severity, &entry.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		if oldVal.Valid {
			entry.OldValue = oldVal.String
		}
		if newVal.Valid {
			entry.NewValue = newVal.String
		}
		logs = append(logs, entry)
	}
	return logs, total, nil
}
