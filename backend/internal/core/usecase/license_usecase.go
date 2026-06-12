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
	"github.com/salonflow/salonflow-track/pkg/apperror"
	"github.com/salonflow/salonflow-track/pkg/uid"
)

// ErrLicenseRestricted is returned when an operation is blocked due to license restrictions.
var ErrLicenseRestricted = errors.New("operation blocked: license expired or invalid, please renew")

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
		_ = uc.repo.UpdateLicense(ctx, lic)
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
	lic.LastVerifiedAt = now.Format(time.RFC3339)

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

	_ = uc.repo.UpdateLicense(ctx, lic)
	uc.logEvent(ctx, lic.ID, domain.LicenseEventValidated, "")

	// Check and create notifications
	uc.checkAndCreateNotifications(ctx, lic, days)

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

// ImportLicenseFile imports a license from an encrypted license file.
func (uc *LicenseUseCase) ImportLicenseFile(ctx context.Context, fileData []byte) (*domain.License, error) {
	lfd, err := uc.engine.ParseLicenseFile(fileData)
	if err != nil {
		return nil, apperror.Business("INVALID_LICENSE_FILE", fmt.Sprintf("Failed to parse license file: %s", err.Error()))
	}

	if lfd.LicenseKey == "" {
		return nil, apperror.Validation("license_key", "License file missing license key")
	}
	if lfd.ExpiryDate == "" {
		return nil, apperror.Validation("expiry_date", "License file missing expiry date")
	}

	// Validate signature from file
	deviceID := lfd.DeviceID
	if deviceID == "" {
		// Bind to current device
		deviceID = uc.engine.GenerateDeviceID()
	}

	// If file has a signature, validate it
	if lfd.Signature != "" {
		if !uc.engine.ValidateSignature(lfd.LicenseKey, lfd.ExpiryDate, deviceID, lfd.Signature) {
			return nil, apperror.Business("INVALID_SIGNATURE", "License file signature is invalid")
		}
	} else {
		// Generate signature for this device
		lfd.Signature = uc.engine.SignLicense(lfd.LicenseKey, lfd.ExpiryDate, deviceID)
	}

	lic := domain.NewLicense(lfd.LicenseKey, lfd.CustomerName, lfd.SalonName, deviceID, lfd.IssuedDate, lfd.ExpiryDate, lfd.Signature)

	if err := uc.repo.CreateLicense(ctx, lic); err != nil {
		return nil, err
	}

	uc.logEvent(ctx, lic.ID, domain.LicenseEventImported, fmt.Sprintf("Imported license file for %s", lfd.SalonName))
	return lic, nil
}

// ExportLicenseFile exports the current license as an encrypted file.
func (uc *LicenseUseCase) ExportLicenseFile(ctx context.Context) ([]byte, error) {
	lic, err := uc.repo.GetActiveLicense(ctx)
	if err != nil {
		return nil, apperror.Business("NO_LICENSE", "No license found to export")
	}
	return uc.engine.ExportLicenseFile(lic)
}

// Renew extends the license by one month.
func (uc *LicenseUseCase) Renew(ctx context.Context, key string) (*domain.License, error) {
	lic, err := uc.repo.GetActiveLicense(ctx)
	if err != nil {
		return nil, apperror.Business("NO_LICENSE", "No license found to renew")
	}

	// If a new key is provided, validate it matches or is a renewal key
	if key != "" && key != lic.LicenseKey {
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
	newGrace := baseDate.AddDate(0, 1, domain.GracePeriodDays).Format("2006-01-02")
	lic.ExpiryDate = newExpiry
	lic.GraceUntil = newGrace
	lic.Status = domain.LicenseStatusActive
	lic.Signature = uc.engine.SignLicense(lic.LicenseKey, newExpiry, lic.DeviceID)
	lic.LastValidation = time.Now().UTC().Format(time.RFC3339)
	lic.LastVerifiedAt = time.Now().UTC().Format(time.RFC3339)

	if err := uc.repo.UpdateLicense(ctx, lic); err != nil {
		return nil, err
	}

	uc.logEvent(ctx, lic.ID, domain.LicenseEventRenewed, fmt.Sprintf("Renewed until %s", newExpiry))
	return lic, nil
}

// IsOperationAllowed checks if an operation is allowed under the current license.
func (uc *LicenseUseCase) IsOperationAllowed(ctx context.Context, operation string) error {
	lic, err := uc.repo.GetActiveLicense(ctx)
	if err != nil {
		// No license at all - restrict
		return ErrLicenseRestricted
	}

	if !lic.IsRestricted() {
		return nil
	}

	// Check if this operation is in the restricted list
	for _, op := range domain.RestrictedOperations {
		if op == operation {
			return ErrLicenseRestricted
		}
	}
	return nil
}

// GetDeviceID returns the current device's unique identifier.
func (uc *LicenseUseCase) GetDeviceID() string {
	return uc.engine.GenerateDeviceID()
}

// GetNotifications returns license notifications.
func (uc *LicenseUseCase) GetNotifications(ctx context.Context, unreadOnly bool) ([]domain.LicenseNotification, error) {
	lic, err := uc.repo.GetActiveLicense(ctx)
	if err != nil {
		return nil, nil
	}
	return uc.repo.ListNotifications(ctx, lic.ID.String(), unreadOnly)
}

// MarkNotificationRead marks a notification as read.
func (uc *LicenseUseCase) MarkNotificationRead(ctx context.Context, id string) error {
	return uc.repo.MarkNotificationRead(ctx, id)
}

// DismissNotification dismisses a notification.
func (uc *LicenseUseCase) DismissNotification(ctx context.Context, id string) error {
	return uc.repo.DismissNotification(ctx, id)
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

// checkAndCreateNotifications creates notifications based on days remaining.
func (uc *LicenseUseCase) checkAndCreateNotifications(ctx context.Context, lic *domain.License, days int) {
	licID := lic.ID.String()

	var notifType, title, message string

	switch {
	case days == 7 || (days > 0 && days <= 7 && days > 3):
		notifType = domain.NotifySevenDaysRemaining
		title = "License Expiring Soon"
		message = fmt.Sprintf("Your license expires in %d days. Please renew to avoid service interruption.", days)
	case days == 3 || (days > 0 && days <= 3 && days > 1):
		notifType = domain.NotifyThreeDaysRemaining
		title = "License Expiring in 3 Days"
		message = fmt.Sprintf("Your license expires in %d days. Renew now to continue uninterrupted access.", days)
	case days == 1 || (days > 0 && days <= 1):
		notifType = domain.NotifyOneDayRemaining
		title = "License Expires Tomorrow"
		message = "Your license expires tomorrow. Renew immediately to avoid restricted mode."
	case days == 0:
		notifType = domain.NotifyExpired
		title = "License Expired"
		message = "Your license has expired. You are now in a 30-day grace period. Renew to restore full access."
	case days < 0 && days >= -domain.GracePeriodDays:
		notifType = domain.NotifyGracePeriodRemaining
		title = "Grace Period Active"
		graceDays := domain.GracePeriodDays + days
		message = fmt.Sprintf("Grace period: %d days remaining. After this, the system will enter restricted mode.", graceDays)
	default:
		return
	}

	// Don't create duplicate notifications of the same type
	exists, _ := uc.repo.HasNotificationType(ctx, licID, notifType)
	if exists {
		return
	}

	n := &domain.LicenseNotification{
		ID:               uid.New().String(),
		LicenseID:        licID,
		NotificationType: notifType,
		Title:            title,
		Message:          message,
		IsRead:           false,
		IsDismissed:      false,
		CreatedAt:        time.Now().UTC(),
	}
	if err := uc.repo.CreateNotification(ctx, n); err != nil {
		uc.log.Error("failed to create license notification", "error", err)
	}
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
