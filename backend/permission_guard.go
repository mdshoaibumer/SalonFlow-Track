package main

import (
	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// PermissionGuard provides permission checking for all binding services.
// It holds a reference to the AuthService to validate the current session.
type PermissionGuard struct {
	auth *AuthService
}

// NewPermissionGuard creates a new PermissionGuard.
func NewPermissionGuard(auth *AuthService) *PermissionGuard {
	return &PermissionGuard{auth: auth}
}

// Require checks that the current user has the given permission.
// Returns nil if authorized, ErrPermissionDenied otherwise.
func (g *PermissionGuard) Require(permission string) error {
	if g.auth == nil {
		return nil // Guard not configured, allow (for backward compat during transition)
	}
	if g.auth.token == "" {
		return usecase.ErrSessionNotFound
	}
	return g.auth.uc.CheckPermission(g.auth.ctx, g.auth.token, permission)
}

// RequireAny checks that the current user has at least one of the given permissions.
func (g *PermissionGuard) RequireAny(permissions ...string) error {
	if g.auth == nil {
		return nil
	}
	if g.auth.token == "" {
		return usecase.ErrSessionNotFound
	}
	for _, p := range permissions {
		if g.auth.uc.CheckPermission(g.auth.ctx, g.auth.token, p) == nil {
			return nil
		}
	}
	return usecase.ErrPermissionDenied
}

// CurrentUserID returns the current user's ID, or empty string if not authenticated.
func (g *PermissionGuard) CurrentUserID() string {
	if g.auth == nil || g.auth.token == "" {
		return ""
	}
	session, err := g.auth.GetCurrentSession()
	if err != nil {
		return ""
	}
	return session.UserID
}

// CurrentUsername returns the current user's username, or empty string.
func (g *PermissionGuard) CurrentUsername() string {
	if g.auth == nil || g.auth.token == "" {
		return ""
	}
	session, err := g.auth.GetCurrentSession()
	if err != nil {
		return ""
	}
	return session.Username
}
