package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// LicenseRepository manages license data in the database.
type LicenseRepository interface {
	CreateLicense(ctx context.Context, license *domain.License) error
	UpdateLicense(ctx context.Context, license *domain.License) error
	GetActiveLicense(ctx context.Context) (*domain.License, error)
	GetLicenseByKey(ctx context.Context, key string) (*domain.License, error)

	CreateEvent(ctx context.Context, event *domain.LicenseEvent) error
	ListEvents(ctx context.Context, licenseID uuid.UUID, limit, offset int) ([]domain.LicenseEvent, int, error)
}

// LicenseEngine handles license key generation, signing, validation, and device binding.
type LicenseEngine interface {
	GenerateKey() string
	GenerateDeviceID() string
	SignLicense(key, expiryDate, deviceID string) string
	ValidateSignature(key, expiryDate, deviceID, signature string) bool
}
