package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// LicenseUseCase handles license business logic.
type LicenseUseCase struct {
	repo   ports.LicenseRepository
	engine ports.LicenseEngine
	log    *slog.Logger
}

// NewLicenseUseCase creates a new LicenseUseCase.
func NewLicenseUseCase(repo ports.LicenseRepository, engine ports.LicenseEngine, log *slog.Logger) *LicenseUseCase {
	return &LicenseUseCase{repo: repo, engine: engine, log: log}
}

// GetStatus returns the current license status for the dashboard.
func (uc *LicenseUseCase) GetStatus(ctx context.Context) (*domain.LicenseStatus, error) {
	lic, err := uc.repo.GetActiveLicense(ctx)
	if err != nil {
		return &domain.LicenseStatus{
			IsRestricted:  true,
			NeedsRenewal:  true,
			DaysRemaining: -999,
		}, nil
	}
	return &domain.LicenseStatus{
		License:            lic,
		DaysRemaining:      lic.DaysRemaining(),
		GraceDaysRemaining: lic.GraceDaysRemaining(),
		IsRestricted:       lic.IsRestricted(),
		NeedsRenewal:       lic.DaysRemaining() <= 30,
	}, nil
}

// Validate validates the current license and updates its status.
func (uc *LicenseUseCase) Validate(ctx context.Context) (*domain.LicenseValidation, error) {
	lic, err := uc.repo.GetActiveLicense(ctx)
	if err != nil {
		return &domain.LicenseValidation{
			Valid:        false,
			Status:       domain.LicenseStatusExpired,
			IsRestricted: true,
			Message:      "No license found. Please activate.",
		}, nil
	}

	// Check signature integrity
	if !uc.engine.ValidateSignature(lic.LicenseKey, lic.ExpiryDate, lic.DeviceID, lic.Signature) {
		lic.Status = domain.LicenseStatusSuspended
		uc.repo.UpdateLicense(ctx, lic)
		uc.logEvent(ctx, lic.ID, domain.LicenseEventSuspended, "Tamper detected: invalid signature")
		return &domain.LicenseValidation{
			Valid:        false,
			Status:       domain.LicenseStatusSuspended,
			IsRestricted: true,
			Message:      "License tampered. Contact support.",
		}, nil
	}

	// Check device binding
	currentDevice := uc.engine.GenerateDeviceID()
	if lic.DeviceID != "" && lic.DeviceID != currentDevice {
		return &domain.LicenseValidation{
			Valid:        false,
			Status:       domain.LicenseStatusSuspended,
			IsRestricted: true,
			Message:      "License bound to different device. Contact support.",
		}, nil
	}

	// Check expiry
	days := lic.DaysRemaining()
	now := time.Now().UTC()
	lic.LastValidation = now.Format(time.RFC3339)

	switch {
	case days > 0:
		lic.Status = domain.LicenseStatusActive
	case days >= -domain.GracePeriodDays:
		if lic.Status != domain.LicenseStatusGracePeriod {
			lic.Status = domain.LicenseStatusGracePeriod
			uc.logEvent(ctx, lic.ID, domain.LicenseEventGraceStarted, fmt.Sprintf("Grace period started, %d days remaining", domain.GracePeriodDays+days))
		}
	default:
		if lic.Status != domain.LicenseStatusExpired {
			lic.Status = domain.LicenseStatusExpired
			uc.logEvent(ctx, lic.ID, domain.LicenseEventExpired, "License expired, restricted mode")
			uc.logEvent(ctx, lic.ID, domain.LicenseEventRestricted, "Restricted mode activated")
		}
	}

	uc.repo.UpdateLicense(ctx, lic)
	uc.logEvent(ctx, lic.ID, domain.LicenseEventValidated, "")

	return &domain.LicenseValidation{
		Valid:         lic.Status == domain.LicenseStatusActive || lic.Status == domain.LicenseStatusGracePeriod,
		Status:        lic.Status,
		DaysRemaining: days,
		IsRestricted:  lic.IsRestricted(),
		Message:       uc.statusMessage(lic.Status, days),
	}, nil
}

// Activate activates a new license key.
func (uc *LicenseUseCase) Activate(ctx context.Context, key, customerName, salonName string) (*domain.License, error) {
	if key == "" {
		return nil, apperror.Validation("license_key", "License key is required")
	}
	if customerName == "" {
		return nil, apperror.Validation("customer_name", "Customer name is required")
	}
	if salonName == "" {
		return nil, apperror.Validation("salon_name", "Salon name is required")
	}

	deviceID := uc.engine.GenerateDeviceID()
	issuedDate := time.Now().Format("2006-01-02")
	expiryDate := time.Now().AddDate(0, 1, 0).Format("2006-01-02")

	signature := uc.engine.SignLicense(key, expiryDate, deviceID)
	lic := domain.NewLicense(key, customerName, salonName, deviceID, issuedDate, expiryDate, signature)

	if err := uc.repo.CreateLicense(ctx, lic); err != nil {
		return nil, err
	}

	uc.logEvent(ctx, lic.ID, domain.LicenseEventActivated, fmt.Sprintf("Activated for %s", salonName))
	return lic, nil
}

// Renew extends the license by one month.
func (uc *LicenseUseCase) Renew(ctx context.Context, key string) (*domain.License, error) {
	lic, err := uc.repo.GetActiveLicense(ctx)
	if err != nil {
		return nil, apperror.Business("NO_LICENSE", "No license found to renew")
	}

	// If a new key is provided, validate it matches or is a renewal key
	if key != "" && key != lic.LicenseKey {
		// Could be a new key for same license - update
		lic.LicenseKey = key
	}

	// Extend from today or from expiry, whichever is later
	var baseDate time.Time
	expiry, err := time.Parse("2006-01-02", lic.ExpiryDate)
	if err != nil || expiry.Before(time.Now()) {
		baseDate = time.Now()
	} else {
		baseDate = expiry
	}
	newExpiry := baseDate.AddDate(0, 1, 0).Format("2006-01-02")
	lic.ExpiryDate = newExpiry
	lic.Status = domain.LicenseStatusActive
	lic.Signature = uc.engine.SignLicense(lic.LicenseKey, newExpiry, lic.DeviceID)
	lic.LastValidation = time.Now().UTC().Format(time.RFC3339)

	if err := uc.repo.UpdateLicense(ctx, lic); err != nil {
		return nil, err
	}

	uc.logEvent(ctx, lic.ID, domain.LicenseEventRenewed, fmt.Sprintf("Renewed until %s", newExpiry))
	return lic, nil
}

// ListEvents returns audit events for the active license.
func (uc *LicenseUseCase) ListEvents(ctx context.Context, page, perPage int) ([]domain.LicenseEvent, int, error) {
	lic, err := uc.repo.GetActiveLicense(ctx)
	if err != nil {
		return nil, 0, nil
	}
	offset := (page - 1) * perPage
	return uc.repo.ListEvents(ctx, lic.ID, perPage, offset)
}

func (uc *LicenseUseCase) logEvent(ctx context.Context, licenseID uuid.UUID, eventType, notes string) {
	event := domain.NewLicenseEvent(licenseID, eventType, notes)
	if err := uc.repo.CreateEvent(ctx, event); err != nil {
		uc.log.Error("failed to create license event", "error", err)
	}
}

func (uc *LicenseUseCase) statusMessage(status string, days int) string {
	switch status {
	case domain.LicenseStatusActive:
		if days <= 7 {
			return fmt.Sprintf("License expires in %d days. Please renew.", days)
		}
		return "License active"
	case domain.LicenseStatusGracePeriod:
		grace := domain.GracePeriodDays + days
		return fmt.Sprintf("License expired. Grace period: %d days remaining. Please renew.", grace)
	case domain.LicenseStatusExpired:
		return "License expired. Restricted mode. Please renew to restore full access."
	case domain.LicenseStatusSuspended:
		return "License suspended. Contact support."
	default:
		return "Unknown license status"
	}
}
