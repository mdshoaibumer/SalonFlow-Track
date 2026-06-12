package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/pkg/uid"
)

// License statuses.
const (
	LicenseStatusActive      = "active"
	LicenseStatusGracePeriod = "grace_period"
	LicenseStatusExpired     = "expired"
	LicenseStatusSuspended   = "suspended"
)

// License event types.
const (
	LicenseEventActivated    = "activated"
	LicenseEventRenewed      = "renewed"
	LicenseEventExpired      = "expired"
	LicenseEventValidated    = "validated"
	LicenseEventSuspended    = "suspended"
	LicenseEventGraceStarted = "grace_started"
	LicenseEventRestricted   = "restricted"
	LicenseEventImported     = "imported"
)

// Notification types.
const (
	NotifySevenDaysRemaining   = "7_days_remaining"
	NotifyThreeDaysRemaining   = "3_days_remaining"
	NotifyOneDayRemaining      = "1_day_remaining"
	NotifyExpired              = "expired"
	NotifyGracePeriodRemaining = "grace_period_remaining"
)

// GracePeriodDays is the number of days after expiry before restricted mode.
const GracePeriodDays = 30

// License represents a software license.
type License struct {
	ID             uuid.UUID `json:"id"`
	LicenseKey     string    `json:"license_key"`
	CustomerName   string    `json:"customer_name"`
	SalonName      string    `json:"salon_name"`
	DeviceID       string    `json:"device_id"`
	IssuedDate     string    `json:"issued_date"`
	ExpiryDate     string    `json:"expiry_date"`
	GraceUntil     string    `json:"grace_until"`
	Status         string    `json:"status"`
	Signature      string    `json:"signature"`
	LastValidation string    `json:"last_validation"`
	LastVerifiedAt string    `json:"last_verified_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// NewLicense creates a new license.
func NewLicense(key, customerName, salonName, deviceID, issuedDate, expiryDate, signature string) *License {
	now := time.Now().UTC()
	expiry, _ := time.Parse("2006-01-02", expiryDate)
	graceUntil := expiry.AddDate(0, 0, GracePeriodDays).Format("2006-01-02")
	return &License{
		ID:             uid.New(),
		LicenseKey:     key,
		CustomerName:   customerName,
		SalonName:      salonName,
		DeviceID:       deviceID,
		IssuedDate:     issuedDate,
		ExpiryDate:     expiryDate,
		GraceUntil:     graceUntil,
		Status:         LicenseStatusActive,
		Signature:      signature,
		LastValidation: now.Format(time.RFC3339),
		LastVerifiedAt: now.Format(time.RFC3339),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// DaysRemaining returns the number of days until expiry (negative if expired).
func (l *License) DaysRemaining() int {
	expiry, err := time.Parse("2006-01-02", l.ExpiryDate)
	if err != nil {
		return -999
	}
	return int(time.Until(expiry).Hours() / 24)
}

// GraceDaysRemaining returns the number of grace days remaining (0 if not in grace).
func (l *License) GraceDaysRemaining() int {
	days := l.DaysRemaining()
	if days >= 0 {
		return GracePeriodDays
	}
	remaining := GracePeriodDays + days
	if remaining < 0 {
		return 0
	}
	return remaining
}

// IsRestricted returns true if the license is in restricted mode.
func (l *License) IsRestricted() bool {
	return l.Status == LicenseStatusExpired || l.Status == LicenseStatusSuspended
}

// LicenseEvent represents an audit event for a license.
type LicenseEvent struct {
	ID        uuid.UUID `json:"id"`
	LicenseID uuid.UUID `json:"license_id"`
	EventType string    `json:"event_type"`
	EventDate string    `json:"event_date"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
}

// NewLicenseEvent creates a new license event.
func NewLicenseEvent(licenseID uuid.UUID, eventType, notes string) *LicenseEvent {
	now := time.Now().UTC()
	return &LicenseEvent{
		ID:        uid.New(),
		LicenseID: licenseID,
		EventType: eventType,
		EventDate: now.Format(time.RFC3339),
		Notes:     notes,
		CreatedAt: now,
	}
}

// LicenseStatus represents the current state for the dashboard.
type LicenseStatus struct {
	License            *License `json:"license"`
	DaysRemaining      int      `json:"days_remaining"`
	GraceDaysRemaining int      `json:"grace_days_remaining"`
	IsRestricted       bool     `json:"is_restricted"`
	NeedsRenewal       bool     `json:"needs_renewal"`
}

// LicenseValidation represents a validation result.
type LicenseValidation struct {
	Valid         bool   `json:"valid"`
	Status        string `json:"status"`
	DaysRemaining int    `json:"days_remaining"`
	IsRestricted  bool   `json:"is_restricted"`
	Message       string `json:"message"`
}

// LicenseNotification represents a notification about license expiry.
type LicenseNotification struct {
	ID               string    `json:"id"`
	LicenseID        string    `json:"license_id"`
	NotificationType string    `json:"notification_type"`
	Title            string    `json:"title"`
	Message          string    `json:"message"`
	IsRead           bool      `json:"is_read"`
	IsDismissed      bool      `json:"is_dismissed"`
	CreatedAt        time.Time `json:"created_at"`
}

// LicenseFileData represents the data contained in an imported license file.
type LicenseFileData struct {
	LicenseKey   string `json:"license_key"`
	CustomerName string `json:"customer_name"`
	SalonName    string `json:"salon_name"`
	IssuedDate   string `json:"issued_date"`
	ExpiryDate   string `json:"expiry_date"`
	DeviceID     string `json:"device_id"`
	Signature    string `json:"signature"`
}

// RestrictedOperation defines operations that are blocked in restricted mode.
const (
	OpInvoiceCreate   = "invoice.create"
	OpCustomerCreate  = "customer.create"
	OpSalaryGenerate  = "salary.generate"
	OpExpenseCreate   = "expense.create"
	OpInventoryChange = "inventory.change"
)

// RestrictedOperations is the list of all operations blocked in restricted mode.
var RestrictedOperations = []string{
	OpInvoiceCreate,
	OpCustomerCreate,
	OpSalaryGenerate,
	OpExpenseCreate,
	OpInventoryChange,
}
