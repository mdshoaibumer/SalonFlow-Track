package main

import (
	"context"

	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// LicenseGuard enforces license restrictions on write operations.
// When the license is expired/suspended, only read/report/backup/export operations are allowed.
type LicenseGuard struct {
	licenseUC *usecase.LicenseUseCase
	ctx       context.Context
}

// NewLicenseGuard creates a new LicenseGuard.
func NewLicenseGuard(licenseUC *usecase.LicenseUseCase) *LicenseGuard {
	return &LicenseGuard{licenseUC: licenseUC}
}

// SetContext sets the Wails context for the guard.
func (g *LicenseGuard) SetContext(ctx context.Context) {
	g.ctx = ctx
}

// RequireActive checks that the license allows the given operation.
// Returns nil if allowed, ErrLicenseRestricted if blocked.
func (g *LicenseGuard) RequireActive(operation string) error {
	if g.licenseUC == nil {
		return nil // Guard not configured, allow
	}
	return g.licenseUC.IsOperationAllowed(g.ctx, operation)
}
