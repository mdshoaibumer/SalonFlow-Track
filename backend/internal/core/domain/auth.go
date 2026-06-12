package domain

import "time"

// User represents an application user with authentication credentials.
type User struct {
	ID                 string     `json:"id"`
	Username           string     `json:"username"`
	Email              string     `json:"email"`
	PasswordHash       string     `json:"-"`
	DisplayName        string     `json:"display_name"`
	Phone              string     `json:"phone"`
	IsActive           bool       `json:"is_active"`
	IsLocked           bool       `json:"is_locked"`
	FailedAttempts     int        `json:"failed_attempts"`
	LockedUntil        *time.Time `json:"locked_until,omitempty"`
	LastLoginAt        *time.Time `json:"last_login_at,omitempty"`
	PasswordChangedAt  time.Time  `json:"password_changed_at"`
	MustChangePassword bool       `json:"must_change_password"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	Roles              []Role     `json:"roles,omitempty"`
	Permissions        []string   `json:"permissions,omitempty"`
}

// Role represents an RBAC role.
type Role struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsSystem    bool      `json:"is_system"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Permission represents a granular permission.
type Permission struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Module      string    `json:"module"`
	Action      string    `json:"action"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// Session represents an active user session.
type Session struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	TokenHash    string    `json:"-"`
	DeviceID     string    `json:"device_id"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	RememberMe   bool      `json:"remember_me"`
	ExpiresAt    time.Time `json:"expires_at"`
	LastActiveAt time.Time `json:"last_active_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// AuditLog represents an audit trail entry.
type AuditLog struct {
	ID          string    `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	UserID      string    `json:"user_id"`
	Username    string    `json:"username"`
	Action      string    `json:"action"`
	Module      string    `json:"module"`
	EntityType  string    `json:"entity_type"`
	EntityID    string    `json:"entity_id"`
	Description string    `json:"description"`
	OldValue    string    `json:"old_value,omitempty"`
	NewValue    string    `json:"new_value,omitempty"`
	DeviceID    string    `json:"device_id"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	AppVersion  string    `json:"app_version"`
	Severity    string    `json:"severity"`
	CreatedAt   time.Time `json:"created_at"`
}

// SessionInfo is the public session data returned to the frontend.
type SessionInfo struct {
	UserID      string   `json:"user_id"`
	Username    string   `json:"username"`
	DisplayName string   `json:"display_name"`
	Email       string   `json:"email"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}
