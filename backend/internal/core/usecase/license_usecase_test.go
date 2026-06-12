package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"log/slog"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// ============================================================
// Mock: LicenseRepository
// ============================================================

type mockLicenseRepo struct {
	license       *domain.License
	events        []domain.LicenseEvent
	notifications []domain.LicenseNotification
	createErr     error
	updateErr     error
}

func (m *mockLicenseRepo) CreateLicense(_ context.Context, lic *domain.License) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.license = lic
	return nil
}

func (m *mockLicenseRepo) UpdateLicense(_ context.Context, lic *domain.License) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.license = lic
	return nil
}

func (m *mockLicenseRepo) GetActiveLicense(_ context.Context) (*domain.License, error) {
	if m.license == nil {
		return nil, apperror.NotFound("license", "active")
	}
	return m.license, nil
}

func (m *mockLicenseRepo) GetLicenseByKey(_ context.Context, key string) (*domain.License, error) {
	if m.license != nil && m.license.LicenseKey == key {
		return m.license, nil
	}
	return nil, apperror.NotFound("license", key)
}

func (m *mockLicenseRepo) CreateEvent(_ context.Context, event *domain.LicenseEvent) error {
	m.events = append(m.events, *event)
	return nil
}

func (m *mockLicenseRepo) ListEvents(_ context.Context, _ uuid.UUID, limit, offset int) ([]domain.LicenseEvent, int, error) {
	total := len(m.events)
	end := offset + limit
	if end > total {
		end = total
	}
	if offset >= total {
		return nil, total, nil
	}
	return m.events[offset:end], total, nil
}

func (m *mockLicenseRepo) CreateNotification(_ context.Context, n *domain.LicenseNotification) error {
	m.notifications = append(m.notifications, *n)
	return nil
}

func (m *mockLicenseRepo) ListNotifications(_ context.Context, _ string, unreadOnly bool) ([]domain.LicenseNotification, error) {
	if unreadOnly {
		var result []domain.LicenseNotification
		for _, n := range m.notifications {
			if !n.IsRead && !n.IsDismissed {
				result = append(result, n)
			}
		}
		return result, nil
	}
	return m.notifications, nil
}

func (m *mockLicenseRepo) MarkNotificationRead(_ context.Context, id string) error {
	for i := range m.notifications {
		if m.notifications[i].ID == id {
			m.notifications[i].IsRead = true
			return nil
		}
	}
	return nil
}

func (m *mockLicenseRepo) DismissNotification(_ context.Context, id string) error {
	for i := range m.notifications {
		if m.notifications[i].ID == id {
			m.notifications[i].IsDismissed = true
			return nil
		}
	}
	return nil
}

func (m *mockLicenseRepo) HasNotificationType(_ context.Context, _ string, notifType string) (bool, error) {
	for _, n := range m.notifications {
		if n.NotificationType == notifType && !n.IsDismissed {
			return true, nil
		}
	}
	return false, nil
}

// Ensure mockLicenseRepo satisfies the interface
var _ ports.LicenseRepository = (*mockLicenseRepo)(nil)

// ============================================================
// Mock: LicenseEngine
// ============================================================

type mockLicenseEngine struct {
	deviceID   string
	signature  string
	validSig   bool
	parseData  *domain.LicenseFileData
	parseErr   error
	exportData []byte
	exportErr  error
}

func (m *mockLicenseEngine) GenerateKey() string {
	return "SALONFLOW-TEST-XXXX-YYYY"
}

func (m *mockLicenseEngine) GenerateDeviceID() string {
	return m.deviceID
}

func (m *mockLicenseEngine) SignLicense(_, _, _ string) string {
	return m.signature
}

func (m *mockLicenseEngine) ValidateSignature(_, _, _, _ string) bool {
	return m.validSig
}

func (m *mockLicenseEngine) ParseLicenseFile(_ []byte) (*domain.LicenseFileData, error) {
	if m.parseErr != nil {
		return nil, m.parseErr
	}
	return m.parseData, nil
}

func (m *mockLicenseEngine) ExportLicenseFile(_ *domain.License) ([]byte, error) {
	if m.exportErr != nil {
		return nil, m.exportErr
	}
	if m.exportData != nil {
		return m.exportData, nil
	}
	return []byte(`{"license_key":"exported"}`), nil
}

var _ ports.LicenseEngine = (*mockLicenseEngine)(nil)

// ============================================================
// Helper
// ============================================================

func newTestLicenseUC(repo *mockLicenseRepo, engine *mockLicenseEngine) *LicenseUseCase {
	return NewLicenseUseCase(repo, engine, slog.Default())
}

func daysFromNow(days int) string {
	return time.Now().AddDate(0, 0, days).Format("2006-01-02")
}

func makeLicense(key, deviceID string, daysUntilExpiry int, status string) *domain.License {
	expiry := daysFromNow(daysUntilExpiry)
	grace := daysFromNow(daysUntilExpiry + domain.GracePeriodDays)
	return &domain.License{
		ID:             uuid.New(),
		LicenseKey:     key,
		CustomerName:   "Test Customer",
		SalonName:      "Test Salon",
		DeviceID:       deviceID,
		IssuedDate:     daysFromNow(-30),
		ExpiryDate:     expiry,
		GraceUntil:     grace,
		Signature:      "valid-sig",
		Status:         status,
		LastValidation: time.Now().Format(time.RFC3339),
		LastVerifiedAt: time.Now().Format(time.RFC3339),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// ============================================================
// SECTION 1: Activation Tests
// ============================================================

func TestActivate_Success(t *testing.T) {
	repo := &mockLicenseRepo{}
	engine := &mockLicenseEngine{deviceID: "dev-001", signature: "signed-abc"}
	uc := newTestLicenseUC(repo, engine)

	lic, err := uc.Activate(context.Background(), "SALONFLOW-AAAA-BBBB-CCCC", "John Doe", "Glamour Salon")
	if err != nil {
		t.Fatalf("Activate failed: %v", err)
	}

	// Verify all fields
	if lic.LicenseKey != "SALONFLOW-AAAA-BBBB-CCCC" {
		t.Errorf("key = %q, want SALONFLOW-AAAA-BBBB-CCCC", lic.LicenseKey)
	}
	if lic.CustomerName != "John Doe" {
		t.Errorf("customer = %q, want John Doe", lic.CustomerName)
	}
	if lic.SalonName != "Glamour Salon" {
		t.Errorf("salon = %q, want Glamour Salon", lic.SalonName)
	}
	if lic.DeviceID != "dev-001" {
		t.Errorf("device = %q, want dev-001", lic.DeviceID)
	}
	if lic.Signature != "signed-abc" {
		t.Errorf("signature = %q, want signed-abc", lic.Signature)
	}
	if lic.Status != domain.LicenseStatusActive {
		t.Errorf("status = %q, want active", lic.Status)
	}
	if lic.GraceUntil == "" {
		t.Error("GraceUntil should be set")
	}
	if lic.IssuedDate == "" {
		t.Error("IssuedDate should be set")
	}

	// Verify expiry is ~1 month from now
	expiry, _ := time.Parse("2006-01-02", lic.ExpiryDate)
	diff := time.Until(expiry).Hours() / 24
	if diff < 27 || diff > 32 {
		t.Errorf("expiry should be ~30 days out, got %.0f days", diff)
	}

	// Verify event logged
	if len(repo.events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(repo.events))
	}
	if repo.events[0].EventType != domain.LicenseEventActivated {
		t.Errorf("event type = %q, want activated", repo.events[0].EventType)
	}

	// Verify stored in repo
	if repo.license == nil {
		t.Fatal("license not stored in repo")
	}
}

func TestActivate_EmptyKey(t *testing.T) {
	uc := newTestLicenseUC(&mockLicenseRepo{}, &mockLicenseEngine{deviceID: "d"})
	_, err := uc.Activate(context.Background(), "", "Customer", "Salon")
	if err == nil {
		t.Fatal("expected validation error for empty key")
	}
}

func TestActivate_EmptyCustomerName(t *testing.T) {
	uc := newTestLicenseUC(&mockLicenseRepo{}, &mockLicenseEngine{deviceID: "d"})
	_, err := uc.Activate(context.Background(), "KEY", "", "Salon")
	if err == nil {
		t.Fatal("expected validation error for empty customer name")
	}
}

func TestActivate_EmptySalonName(t *testing.T) {
	uc := newTestLicenseUC(&mockLicenseRepo{}, &mockLicenseEngine{deviceID: "d"})
	_, err := uc.Activate(context.Background(), "KEY", "Customer", "")
	if err == nil {
		t.Fatal("expected validation error for empty salon name")
	}
}

func TestActivate_RepoError(t *testing.T) {
	repo := &mockLicenseRepo{createErr: errors.New("db error")}
	engine := &mockLicenseEngine{deviceID: "d", signature: "s"}
	uc := newTestLicenseUC(repo, engine)

	_, err := uc.Activate(context.Background(), "KEY", "C", "S")
	if err == nil {
		t.Fatal("expected error when repo fails")
	}
}

// ============================================================
// SECTION 2: Validation Tests
// ============================================================

func TestValidate_ActiveLicense(t *testing.T) {
	repo := &mockLicenseRepo{license: makeLicense("KEY", "dev-001", 20, domain.LicenseStatusActive)}
	engine := &mockLicenseEngine{deviceID: "dev-001", validSig: true}
	uc := newTestLicenseUC(repo, engine)

	result, err := uc.Validate(context.Background())
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}
	if !result.Valid {
		t.Error("active license should be valid")
	}
	if result.IsRestricted {
		t.Error("active license should not be restricted")
	}
	if result.Status != domain.LicenseStatusActive {
		t.Errorf("status = %q, want active", result.Status)
	}
	if result.DaysRemaining < 19 || result.DaysRemaining > 21 {
		t.Errorf("days remaining = %d, expected ~20", result.DaysRemaining)
	}
	if result.Message != "License active" {
		t.Errorf("unexpected message: %q", result.Message)
	}
}

func TestValidate_GracePeriod(t *testing.T) {
	// License expired 10 days ago -> within 30-day grace
	repo := &mockLicenseRepo{license: makeLicense("KEY", "dev-001", -10, domain.LicenseStatusActive)}
	engine := &mockLicenseEngine{deviceID: "dev-001", validSig: true}
	uc := newTestLicenseUC(repo, engine)

	result, err := uc.Validate(context.Background())
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}
	if !result.Valid {
		t.Error("grace period license should still be valid")
	}
	if result.IsRestricted {
		t.Error("grace period license should not be restricted")
	}
	if result.Status != domain.LicenseStatusGracePeriod {
		t.Errorf("status = %q, want grace_period", result.Status)
	}
}

func TestValidate_Expired(t *testing.T) {
	// License expired beyond grace period
	repo := &mockLicenseRepo{license: makeLicense("KEY", "dev-001", -(domain.GracePeriodDays + 5), domain.LicenseStatusActive)}
	engine := &mockLicenseEngine{deviceID: "dev-001", validSig: true}
	uc := newTestLicenseUC(repo, engine)

	result, err := uc.Validate(context.Background())
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}
	if result.Valid {
		t.Error("expired license should not be valid")
	}
	if !result.IsRestricted {
		t.Error("expired license should be restricted")
	}
	if result.Status != domain.LicenseStatusExpired {
		t.Errorf("status = %q, want expired", result.Status)
	}

	// Should have logged expired + restricted events
	foundExpired, foundRestricted := false, false
	for _, ev := range repo.events {
		if ev.EventType == domain.LicenseEventExpired {
			foundExpired = true
		}
		if ev.EventType == domain.LicenseEventRestricted {
			foundRestricted = true
		}
	}
	if !foundExpired {
		t.Error("expected expired event")
	}
	if !foundRestricted {
		t.Error("expected restricted event")
	}
}

func TestValidate_NoLicense(t *testing.T) {
	repo := &mockLicenseRepo{}
	engine := &mockLicenseEngine{}
	uc := newTestLicenseUC(repo, engine)

	result, err := uc.Validate(context.Background())
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}
	if result.Valid {
		t.Error("no license should be invalid")
	}
	if !result.IsRestricted {
		t.Error("no license should be restricted")
	}
	if result.Status != domain.LicenseStatusExpired {
		t.Errorf("status = %q, want expired", result.Status)
	}
}

func TestValidate_TamperedSignature(t *testing.T) {
	repo := &mockLicenseRepo{license: makeLicense("KEY", "dev-001", 20, domain.LicenseStatusActive)}
	engine := &mockLicenseEngine{deviceID: "dev-001", validSig: false} // Signature fails
	uc := newTestLicenseUC(repo, engine)

	result, err := uc.Validate(context.Background())
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}
	if result.Valid {
		t.Error("tampered license should be invalid")
	}
	if result.Status != domain.LicenseStatusSuspended {
		t.Errorf("status = %q, want suspended", result.Status)
	}
	if !result.IsRestricted {
		t.Error("suspended license should be restricted")
	}

	// Verify suspend event logged
	found := false
	for _, ev := range repo.events {
		if ev.EventType == domain.LicenseEventSuspended {
			found = true
		}
	}
	if !found {
		t.Error("expected suspended event")
	}
}

func TestValidate_DeviceMismatch(t *testing.T) {
	repo := &mockLicenseRepo{license: makeLicense("KEY", "device-A", 20, domain.LicenseStatusActive)}
	engine := &mockLicenseEngine{deviceID: "device-B", validSig: true} // Different device
	uc := newTestLicenseUC(repo, engine)

	result, err := uc.Validate(context.Background())
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}
	if result.Valid {
		t.Error("device mismatch should be invalid")
	}
	if !result.IsRestricted {
		t.Error("device mismatch should be restricted")
	}
	if result.Status != domain.LicenseStatusSuspended {
		t.Errorf("status = %q, want suspended", result.Status)
	}
}

func TestValidate_EmptyDeviceID_SkipsBinding(t *testing.T) {
	lic := makeLicense("KEY", "", 20, domain.LicenseStatusActive) // No device bound
	repo := &mockLicenseRepo{license: lic}
	engine := &mockLicenseEngine{deviceID: "any-device", validSig: true}
	uc := newTestLicenseUC(repo, engine)

	result, err := uc.Validate(context.Background())
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}
	if !result.Valid {
		t.Error("empty device should skip binding check")
	}
}

func TestValidate_UpdatesLastVerifiedAt(t *testing.T) {
	lic := makeLicense("KEY", "dev-001", 20, domain.LicenseStatusActive)
	lic.LastVerifiedAt = "2020-01-01T00:00:00Z" // Old value
	repo := &mockLicenseRepo{license: lic}
	engine := &mockLicenseEngine{deviceID: "dev-001", validSig: true}
	uc := newTestLicenseUC(repo, engine)

	_, _ = uc.Validate(context.Background())

	if repo.license.LastVerifiedAt == "2020-01-01T00:00:00Z" {
		t.Error("LastVerifiedAt should be updated after validation")
	}
}

func TestValidate_GracePeriodTransition_LogsEvent(t *testing.T) {
	// Expired 1 day ago, not yet marked as grace_period
	lic := makeLicense("KEY", "dev-001", -1, domain.LicenseStatusActive)
	repo := &mockLicenseRepo{license: lic}
	engine := &mockLicenseEngine{deviceID: "dev-001", validSig: true}
	uc := newTestLicenseUC(repo, engine)

	_, _ = uc.Validate(context.Background())

	found := false
	for _, ev := range repo.events {
		if ev.EventType == domain.LicenseEventGraceStarted {
			found = true
		}
	}
	if !found {
		t.Error("expected grace_started event on transition")
	}
}

func TestValidate_AlreadyGracePeriod_NoDoubleEvent(t *testing.T) {
	// Already marked as grace_period
	lic := makeLicense("KEY", "dev-001", -5, domain.LicenseStatusGracePeriod)
	repo := &mockLicenseRepo{license: lic}
	engine := &mockLicenseEngine{deviceID: "dev-001", validSig: true}
	uc := newTestLicenseUC(repo, engine)

	_, _ = uc.Validate(context.Background())

	for _, ev := range repo.events {
		if ev.EventType == domain.LicenseEventGraceStarted {
			t.Error("should not log grace_started again if already in grace_period")
		}
	}
}

func TestValidate_AlreadyExpired_NoDoubleEvent(t *testing.T) {
	// Already expired
	lic := makeLicense("KEY", "dev-001", -(domain.GracePeriodDays + 5), domain.LicenseStatusExpired)
	repo := &mockLicenseRepo{license: lic}
	engine := &mockLicenseEngine{deviceID: "dev-001", validSig: true}
	uc := newTestLicenseUC(repo, engine)

	_, _ = uc.Validate(context.Background())

	for _, ev := range repo.events {
		if ev.EventType == domain.LicenseEventExpired || ev.EventType == domain.LicenseEventRestricted {
			t.Errorf("should not log %q again if already expired", ev.EventType)
		}
	}
}

// ============================================================
// SECTION 3: Renewal Tests
// ============================================================

func TestRenew_ExpiredLicense(t *testing.T) {
	// Expired 5 days ago
	repo := &mockLicenseRepo{license: makeLicense("KEY", "dev-001", -5, domain.LicenseStatusGracePeriod)}
	engine := &mockLicenseEngine{deviceID: "dev-001", signature: "new-sig"}
	uc := newTestLicenseUC(repo, engine)

	lic, err := uc.Renew(context.Background(), "")
	if err != nil {
		t.Fatalf("Renew failed: %v", err)
	}
	if lic.Status != domain.LicenseStatusActive {
		t.Errorf("status after renewal = %q, want active", lic.Status)
	}
	// Should extend from today since expired
	expectedExpiry := time.Now().AddDate(0, 1, 0).Format("2006-01-02")
	if lic.ExpiryDate != expectedExpiry {
		t.Errorf("expiry = %q, want %q", lic.ExpiryDate, expectedExpiry)
	}
	if lic.GraceUntil == "" {
		t.Error("GraceUntil should be set after renewal")
	}
	if lic.Signature != "new-sig" {
		t.Errorf("signature = %q, want new-sig", lic.Signature)
	}
}

func TestRenew_ActiveLicense_ExtendsFromExpiry(t *testing.T) {
	// Still 10 days left -> extends from expiry not from today
	repo := &mockLicenseRepo{license: makeLicense("KEY", "dev-001", 10, domain.LicenseStatusActive)}
	engine := &mockLicenseEngine{deviceID: "dev-001", signature: "sig"}
	uc := newTestLicenseUC(repo, engine)

	lic, err := uc.Renew(context.Background(), "")
	if err != nil {
		t.Fatalf("Renew failed: %v", err)
	}
	// Should extend from current expiry date
	expectedBase := time.Now().AddDate(0, 0, 10)
	expectedExpiry := expectedBase.AddDate(0, 1, 0).Format("2006-01-02")
	if lic.ExpiryDate != expectedExpiry {
		t.Errorf("expiry = %q, want %q", lic.ExpiryDate, expectedExpiry)
	}
}

func TestRenew_WithNewKey(t *testing.T) {
	repo := &mockLicenseRepo{license: makeLicense("OLD-KEY", "dev-001", 10, domain.LicenseStatusActive)}
	engine := &mockLicenseEngine{deviceID: "dev-001", signature: "sig"}
	uc := newTestLicenseUC(repo, engine)

	lic, err := uc.Renew(context.Background(), "NEW-KEY-123")
	if err != nil {
		t.Fatalf("Renew failed: %v", err)
	}
	if lic.LicenseKey != "NEW-KEY-123" {
		t.Errorf("key = %q, want NEW-KEY-123", lic.LicenseKey)
	}
}

func TestRenew_SameKey_NoChange(t *testing.T) {
	repo := &mockLicenseRepo{license: makeLicense("SAME-KEY", "dev-001", 10, domain.LicenseStatusActive)}
	engine := &mockLicenseEngine{deviceID: "dev-001", signature: "sig"}
	uc := newTestLicenseUC(repo, engine)

	lic, err := uc.Renew(context.Background(), "SAME-KEY")
	if err != nil {
		t.Fatalf("Renew failed: %v", err)
	}
	if lic.LicenseKey != "SAME-KEY" {
		t.Errorf("key changed unexpectedly to %q", lic.LicenseKey)
	}
}

func TestRenew_NoLicense(t *testing.T) {
	repo := &mockLicenseRepo{}
	uc := newTestLicenseUC(repo, &mockLicenseEngine{})

	_, err := uc.Renew(context.Background(), "KEY")
	if err == nil {
		t.Fatal("expected error when no license exists")
	}
}

func TestRenew_LogsEvent(t *testing.T) {
	repo := &mockLicenseRepo{license: makeLicense("KEY", "dev-001", 10, domain.LicenseStatusActive)}
	engine := &mockLicenseEngine{deviceID: "dev-001", signature: "sig"}
	uc := newTestLicenseUC(repo, engine)

	_, _ = uc.Renew(context.Background(), "")

	found := false
	for _, ev := range repo.events {
		if ev.EventType == domain.LicenseEventRenewed {
			found = true
		}
	}
	if !found {
		t.Error("expected renewed event")
	}
}

// ============================================================
// SECTION 4: Device Binding Tests
// ============================================================

func TestDeviceBinding_MatchingDevice_OK(t *testing.T) {
	repo := &mockLicenseRepo{license: makeLicense("KEY", "dev-001", 20, domain.LicenseStatusActive)}
	engine := &mockLicenseEngine{deviceID: "dev-001", validSig: true}
	uc := newTestLicenseUC(repo, engine)

	result, _ := uc.Validate(context.Background())
	if !result.Valid {
		t.Error("matching device should validate")
	}
}

func TestDeviceBinding_MismatchedDevice_Fails(t *testing.T) {
	repo := &mockLicenseRepo{license: makeLicense("KEY", "dev-001", 20, domain.LicenseStatusActive)}
	engine := &mockLicenseEngine{deviceID: "dev-002", validSig: true}
	uc := newTestLicenseUC(repo, engine)

	result, _ := uc.Validate(context.Background())
	if result.Valid {
		t.Error("mismatched device should fail")
	}
}

func TestDeviceBinding_EmptyDeviceID_SkipsCheck(t *testing.T) {
	lic := makeLicense("KEY", "", 20, domain.LicenseStatusActive) // No device bound
	repo := &mockLicenseRepo{license: lic}
	engine := &mockLicenseEngine{deviceID: "anything", validSig: true}
	uc := newTestLicenseUC(repo, engine)

	result, _ := uc.Validate(context.Background())
	if !result.Valid {
		t.Error("empty device should pass binding check")
	}
}

func TestGetDeviceID(t *testing.T) {
	engine := &mockLicenseEngine{deviceID: "my-unique-device-id"}
	uc := newTestLicenseUC(&mockLicenseRepo{}, engine)

	id := uc.GetDeviceID()
	if id != "my-unique-device-id" {
		t.Errorf("device ID = %q, want my-unique-device-id", id)
	}
}

// ============================================================
// SECTION 5: Grace Period Tests
// ============================================================

func TestGracePeriod_Day1_StillValid(t *testing.T) {
	repo := &mockLicenseRepo{license: makeLicense("KEY", "dev-001", -1, domain.LicenseStatusActive)}
	engine := &mockLicenseEngine{deviceID: "dev-001", validSig: true}
	uc := newTestLicenseUC(repo, engine)

	result, _ := uc.Validate(context.Background())
	if !result.Valid {
		t.Error("day 1 of grace should still be valid")
	}
	if result.Status != domain.LicenseStatusGracePeriod {
		t.Errorf("status = %q, want grace_period", result.Status)
	}
}

func TestGracePeriod_Day29_StillValid(t *testing.T) {
	repo := &mockLicenseRepo{license: makeLicense("KEY", "dev-001", -29, domain.LicenseStatusActive)}
	engine := &mockLicenseEngine{deviceID: "dev-001", validSig: true}
	uc := newTestLicenseUC(repo, engine)

	result, _ := uc.Validate(context.Background())
	if !result.Valid {
		t.Error("day 29 of grace should still be valid")
	}
}

func TestGracePeriod_Day30_StillValid(t *testing.T) {
	repo := &mockLicenseRepo{license: makeLicense("KEY", "dev-001", -30, domain.LicenseStatusActive)}
	engine := &mockLicenseEngine{deviceID: "dev-001", validSig: true}
	uc := newTestLicenseUC(repo, engine)

	result, _ := uc.Validate(context.Background())
	if !result.Valid {
		t.Error("day 30 of grace (boundary) should still be valid")
	}
}

func TestGracePeriod_Day31_Expired(t *testing.T) {
	repo := &mockLicenseRepo{license: makeLicense("KEY", "dev-001", -31, domain.LicenseStatusActive)}
	engine := &mockLicenseEngine{deviceID: "dev-001", validSig: true}
	uc := newTestLicenseUC(repo, engine)

	result, _ := uc.Validate(context.Background())
	if result.Valid {
		t.Error("day 31 past expiry should NOT be valid")
	}
	if !result.IsRestricted {
		t.Error("day 31 should be restricted")
	}
	if result.Status != domain.LicenseStatusExpired {
		t.Errorf("status = %q, want expired", result.Status)
	}
}

func TestGracePeriod_RenewalRestoresActive(t *testing.T) {
	repo := &mockLicenseRepo{license: makeLicense("KEY", "dev-001", -15, domain.LicenseStatusGracePeriod)}
	engine := &mockLicenseEngine{deviceID: "dev-001", signature: "sig"}
	uc := newTestLicenseUC(repo, engine)

	lic, err := uc.Renew(context.Background(), "")
	if err != nil {
		t.Fatalf("Renew failed: %v", err)
	}
	if lic.Status != domain.LicenseStatusActive {
		t.Errorf("status after renewal = %q, want active", lic.Status)
	}
}

// ============================================================
// SECTION 6: Restricted Mode Tests
// ============================================================

func TestRestrictedMode_ActiveLicense_AllAllowed(t *testing.T) {
	repo := &mockLicenseRepo{license: makeLicense("KEY", "dev-001", 20, domain.LicenseStatusActive)}
	uc := newTestLicenseUC(repo, &mockLicenseEngine{})

	for _, op := range domain.RestrictedOperations {
		err := uc.IsOperationAllowed(context.Background(), op)
		if err != nil {
			t.Errorf("active license should allow %q, got: %v", op, err)
		}
	}
}

func TestRestrictedMode_GracePeriod_AllAllowed(t *testing.T) {
	repo := &mockLicenseRepo{license: makeLicense("KEY", "dev-001", -10, domain.LicenseStatusGracePeriod)}
	uc := newTestLicenseUC(repo, &mockLicenseEngine{})

	for _, op := range domain.RestrictedOperations {
		err := uc.IsOperationAllowed(context.Background(), op)
		if err != nil {
			t.Errorf("grace period should allow %q, got: %v", op, err)
		}
	}
}

func TestRestrictedMode_Expired_AllBlocked(t *testing.T) {
	repo := &mockLicenseRepo{license: makeLicense("KEY", "dev-001", -40, domain.LicenseStatusExpired)}
	uc := newTestLicenseUC(repo, &mockLicenseEngine{})

	for _, op := range domain.RestrictedOperations {
		err := uc.IsOperationAllowed(context.Background(), op)
		if !errors.Is(err, ErrLicenseRestricted) {
			t.Errorf("expired license should block %q, got: %v", op, err)
		}
	}
}

func TestRestrictedMode_Suspended_AllBlocked(t *testing.T) {
	repo := &mockLicenseRepo{license: makeLicense("KEY", "dev-001", 20, domain.LicenseStatusSuspended)}
	uc := newTestLicenseUC(repo, &mockLicenseEngine{})

	for _, op := range domain.RestrictedOperations {
		err := uc.IsOperationAllowed(context.Background(), op)
		if !errors.Is(err, ErrLicenseRestricted) {
			t.Errorf("suspended license should block %q, got: %v", op, err)
		}
	}
}

func TestRestrictedMode_NoLicense_AllBlocked(t *testing.T) {
	repo := &mockLicenseRepo{} // No license
	uc := newTestLicenseUC(repo, &mockLicenseEngine{})

	for _, op := range domain.RestrictedOperations {
		err := uc.IsOperationAllowed(context.Background(), op)
		if !errors.Is(err, ErrLicenseRestricted) {
			t.Errorf("no license should block %q, got: %v", op, err)
		}
	}
}

func TestRestrictedMode_Expired_UnknownOp_Allowed(t *testing.T) {
	// An operation NOT in RestrictedOperations should still be allowed
	repo := &mockLicenseRepo{license: makeLicense("KEY", "dev-001", -40, domain.LicenseStatusExpired)}
	uc := newTestLicenseUC(repo, &mockLicenseEngine{})

	err := uc.IsOperationAllowed(context.Background(), "report.view")
	if err != nil {
		t.Errorf("non-restricted operation should be allowed even when expired, got: %v", err)
	}
}

func TestRestrictedMode_InvoiceCreate_Blocked(t *testing.T) {
	repo := &mockLicenseRepo{license: makeLicense("KEY", "dev-001", -40, domain.LicenseStatusExpired)}
	uc := newTestLicenseUC(repo, &mockLicenseEngine{})

	err := uc.IsOperationAllowed(context.Background(), domain.OpInvoiceCreate)
	if !errors.Is(err, ErrLicenseRestricted) {
		t.Errorf("invoice creation should be blocked, got: %v", err)
	}
}

func TestRestrictedMode_CustomerCreate_Blocked(t *testing.T) {
	repo := &mockLicenseRepo{license: makeLicense("KEY", "dev-001", -40, domain.LicenseStatusExpired)}
	uc := newTestLicenseUC(repo, &mockLicenseEngine{})

	err := uc.IsOperationAllowed(context.Background(), domain.OpCustomerCreate)
	if !errors.Is(err, ErrLicenseRestricted) {
		t.Errorf("customer creation should be blocked, got: %v", err)
	}
}

func TestRestrictedMode_SalaryGenerate_Blocked(t *testing.T) {
	repo := &mockLicenseRepo{license: makeLicense("KEY", "dev-001", -40, domain.LicenseStatusExpired)}
	uc := newTestLicenseUC(repo, &mockLicenseEngine{})

	err := uc.IsOperationAllowed(context.Background(), domain.OpSalaryGenerate)
	if !errors.Is(err, ErrLicenseRestricted) {
		t.Errorf("salary generation should be blocked, got: %v", err)
	}
}

func TestRestrictedMode_ExpenseCreate_Blocked(t *testing.T) {
	repo := &mockLicenseRepo{license: makeLicense("KEY", "dev-001", -40, domain.LicenseStatusExpired)}
	uc := newTestLicenseUC(repo, &mockLicenseEngine{})

	err := uc.IsOperationAllowed(context.Background(), domain.OpExpenseCreate)
	if !errors.Is(err, ErrLicenseRestricted) {
		t.Errorf("expense creation should be blocked, got: %v", err)
	}
}

func TestRestrictedMode_InventoryChange_Blocked(t *testing.T) {
	repo := &mockLicenseRepo{license: makeLicense("KEY", "dev-001", -40, domain.LicenseStatusExpired)}
	uc := newTestLicenseUC(repo, &mockLicenseEngine{})

	err := uc.IsOperationAllowed(context.Background(), domain.OpInventoryChange)
	if !errors.Is(err, ErrLicenseRestricted) {
		t.Errorf("inventory change should be blocked, got: %v", err)
	}
}

// ============================================================
// SECTION 7: License File Import Tests
// ============================================================

func TestImportLicenseFile_Success(t *testing.T) {
	repo := &mockLicenseRepo{}
	engine := &mockLicenseEngine{
		deviceID:  "device-abc",
		signature: "import-sig",
		validSig:  true,
		parseData: &domain.LicenseFileData{
			LicenseKey:   "SALONFLOW-IMPORT-TEST",
			CustomerName: "Import Customer",
			SalonName:    "Import Salon",
			IssuedDate:   "2026-06-01",
			ExpiryDate:   "2026-07-01",
			DeviceID:     "",
			Signature:    "",
		},
	}
	uc := newTestLicenseUC(repo, engine)

	lic, err := uc.ImportLicenseFile(context.Background(), []byte("encrypted-data"))
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}
	if lic.LicenseKey != "SALONFLOW-IMPORT-TEST" {
		t.Errorf("key = %q, want SALONFLOW-IMPORT-TEST", lic.LicenseKey)
	}
	if lic.DeviceID != "device-abc" {
		t.Errorf("device = %q, want device-abc (bound to current)", lic.DeviceID)
	}
	if lic.CustomerName != "Import Customer" {
		t.Errorf("customer = %q, want Import Customer", lic.CustomerName)
	}
	if lic.Status != domain.LicenseStatusActive {
		t.Errorf("imported license status = %q, want active", lic.Status)
	}

	// Should log imported event
	found := false
	for _, ev := range repo.events {
		if ev.EventType == domain.LicenseEventImported {
			found = true
		}
	}
	if !found {
		t.Error("expected imported event")
	}
}

func TestImportLicenseFile_WithExistingDeviceID(t *testing.T) {
	repo := &mockLicenseRepo{}
	engine := &mockLicenseEngine{
		deviceID:  "my-device",
		signature: "sig",
		validSig:  true,
		parseData: &domain.LicenseFileData{
			LicenseKey:   "KEY",
			CustomerName: "C",
			SalonName:    "S",
			IssuedDate:   "2026-01-01",
			ExpiryDate:   "2026-02-01",
			DeviceID:     "file-device", // Pre-bound in file
			Signature:    "file-sig",
		},
	}
	uc := newTestLicenseUC(repo, engine)

	lic, err := uc.ImportLicenseFile(context.Background(), []byte("data"))
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}
	// Should use the device from file
	if lic.DeviceID != "file-device" {
		t.Errorf("device = %q, want file-device (from file)", lic.DeviceID)
	}
}

func TestImportLicenseFile_InvalidFileData(t *testing.T) {
	engine := &mockLicenseEngine{
		parseErr: errors.New("corrupt file"),
	}
	uc := newTestLicenseUC(&mockLicenseRepo{}, engine)

	_, err := uc.ImportLicenseFile(context.Background(), []byte("garbage"))
	if err == nil {
		t.Fatal("expected error for invalid file")
	}
}

func TestImportLicenseFile_MissingKey(t *testing.T) {
	engine := &mockLicenseEngine{
		parseData: &domain.LicenseFileData{
			LicenseKey: "", // Missing
			ExpiryDate: "2026-07-01",
		},
	}
	uc := newTestLicenseUC(&mockLicenseRepo{}, engine)

	_, err := uc.ImportLicenseFile(context.Background(), []byte("data"))
	if err == nil {
		t.Fatal("expected error for missing license key")
	}
}

func TestImportLicenseFile_MissingExpiryDate(t *testing.T) {
	engine := &mockLicenseEngine{
		parseData: &domain.LicenseFileData{
			LicenseKey: "KEY",
			ExpiryDate: "", // Missing
		},
	}
	uc := newTestLicenseUC(&mockLicenseRepo{}, engine)

	_, err := uc.ImportLicenseFile(context.Background(), []byte("data"))
	if err == nil {
		t.Fatal("expected error for missing expiry date")
	}
}

func TestImportLicenseFile_InvalidSignature(t *testing.T) {
	engine := &mockLicenseEngine{
		deviceID: "dev",
		validSig: false, // Signature validation fails
		parseData: &domain.LicenseFileData{
			LicenseKey: "KEY",
			ExpiryDate: "2026-07-01",
			Signature:  "bad-sig", // Has a signature but it's invalid
			DeviceID:   "dev",
		},
	}
	uc := newTestLicenseUC(&mockLicenseRepo{}, engine)

	_, err := uc.ImportLicenseFile(context.Background(), []byte("data"))
	if err == nil {
		t.Fatal("expected error for invalid signature")
	}
}

// ============================================================
// SECTION 8: Export License File Tests
// ============================================================

func TestExportLicenseFile_Success(t *testing.T) {
	repo := &mockLicenseRepo{license: makeLicense("KEY", "dev", 20, domain.LicenseStatusActive)}
	engine := &mockLicenseEngine{exportData: []byte("encrypted-output")}
	uc := newTestLicenseUC(repo, engine)

	data, err := uc.ExportLicenseFile(context.Background())
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}
	if string(data) != "encrypted-output" {
		t.Errorf("data = %q, want encrypted-output", string(data))
	}
}

func TestExportLicenseFile_NoLicense(t *testing.T) {
	uc := newTestLicenseUC(&mockLicenseRepo{}, &mockLicenseEngine{})

	_, err := uc.ExportLicenseFile(context.Background())
	if err == nil {
		t.Fatal("expected error when no license")
	}
}

// ============================================================
// SECTION 9: Notification Tests
// ============================================================

func TestNotification_7DaysRemaining(t *testing.T) {
	lic := makeLicense("KEY", "dev-001", 5, domain.LicenseStatusActive) // 5 days = within 7-day window
	repo := &mockLicenseRepo{license: lic}
	engine := &mockLicenseEngine{deviceID: "dev-001", validSig: true}
	uc := newTestLicenseUC(repo, engine)

	_, _ = uc.Validate(context.Background())

	if len(repo.notifications) == 0 {
		t.Fatal("expected notification for 5 days remaining")
	}
	notifType := repo.notifications[0].NotificationType
	if notifType != domain.NotifySevenDaysRemaining && notifType != domain.NotifyThreeDaysRemaining {
		t.Errorf("notification type = %q, expected 7-day or 3-day", notifType)
	}
}

func TestNotification_3DaysRemaining(t *testing.T) {
	// daysFromNow(3) → DaysRemaining() truncates to ~2 which is in (1, 3] → 3_days_remaining
	lic := makeLicense("KEY", "dev-001", 3, domain.LicenseStatusActive)
	repo := &mockLicenseRepo{license: lic}
	engine := &mockLicenseEngine{deviceID: "dev-001", validSig: true}
	uc := newTestLicenseUC(repo, engine)

	_, _ = uc.Validate(context.Background())

	if len(repo.notifications) == 0 {
		t.Fatal("expected notification for ~2 days remaining")
	}
	notifType := repo.notifications[0].NotificationType
	if notifType != domain.NotifyThreeDaysRemaining {
		t.Errorf("type = %q, want 3_days_remaining", notifType)
	}
}

func TestNotification_1DayRemaining(t *testing.T) {
	// daysFromNow(2) → DaysRemaining() truncates to 1 which hits days==1 → 1_day_remaining
	lic := makeLicense("KEY", "dev-001", 2, domain.LicenseStatusActive)
	repo := &mockLicenseRepo{license: lic}
	engine := &mockLicenseEngine{deviceID: "dev-001", validSig: true}
	uc := newTestLicenseUC(repo, engine)

	_, _ = uc.Validate(context.Background())

	if len(repo.notifications) == 0 {
		t.Fatal("expected notification for 1 day remaining")
	}
	if repo.notifications[0].NotificationType != domain.NotifyOneDayRemaining {
		t.Errorf("type = %q, want 1_day_remaining", repo.notifications[0].NotificationType)
	}
}

func TestNotification_GracePeriodRemaining(t *testing.T) {
	// In grace period (expired 5 days ago)
	lic := makeLicense("KEY", "dev-001", -5, domain.LicenseStatusActive)
	repo := &mockLicenseRepo{license: lic}
	engine := &mockLicenseEngine{deviceID: "dev-001", validSig: true}
	uc := newTestLicenseUC(repo, engine)

	_, _ = uc.Validate(context.Background())

	found := false
	for _, n := range repo.notifications {
		if n.NotificationType == domain.NotifyGracePeriodRemaining {
			found = true
		}
	}
	if !found {
		t.Error("expected grace_period_remaining notification")
	}
}

func TestNotification_NoDuplicates(t *testing.T) {
	lic := makeLicense("KEY", "dev-001", 5, domain.LicenseStatusActive)
	repo := &mockLicenseRepo{license: lic}
	engine := &mockLicenseEngine{deviceID: "dev-001", validSig: true}
	uc := newTestLicenseUC(repo, engine)

	// Validate multiple times
	_, _ = uc.Validate(context.Background())
	_, _ = uc.Validate(context.Background())
	_, _ = uc.Validate(context.Background())

	// Count each notification type
	typeCounts := make(map[string]int)
	for _, n := range repo.notifications {
		typeCounts[n.NotificationType]++
	}
	for nType, count := range typeCounts {
		if count > 1 {
			t.Errorf("duplicate notification %q: count=%d", nType, count)
		}
	}
}

func TestNotification_NoNotificationFor20Days(t *testing.T) {
	lic := makeLicense("KEY", "dev-001", 20, domain.LicenseStatusActive) // 20 days left
	repo := &mockLicenseRepo{license: lic}
	engine := &mockLicenseEngine{deviceID: "dev-001", validSig: true}
	uc := newTestLicenseUC(repo, engine)

	_, _ = uc.Validate(context.Background())

	if len(repo.notifications) != 0 {
		t.Errorf("expected no notifications for 20 days remaining, got %d", len(repo.notifications))
	}
}

func TestGetNotifications(t *testing.T) {
	lic := makeLicense("KEY", "dev-001", 5, domain.LicenseStatusActive)
	repo := &mockLicenseRepo{license: lic}
	engine := &mockLicenseEngine{deviceID: "dev-001", validSig: true}
	uc := newTestLicenseUC(repo, engine)

	// Create a notification via validation
	_, _ = uc.Validate(context.Background())

	notifs, err := uc.GetNotifications(context.Background(), true)
	if err != nil {
		t.Fatalf("GetNotifications failed: %v", err)
	}
	if len(notifs) == 0 {
		t.Error("expected at least one notification")
	}
}

func TestMarkNotificationRead(t *testing.T) {
	lic := makeLicense("KEY", "dev-001", 5, domain.LicenseStatusActive)
	repo := &mockLicenseRepo{license: lic}
	engine := &mockLicenseEngine{deviceID: "dev-001", validSig: true}
	uc := newTestLicenseUC(repo, engine)

	_, _ = uc.Validate(context.Background())
	if len(repo.notifications) == 0 {
		t.Fatal("need at least one notification")
	}

	id := repo.notifications[0].ID
	err := uc.MarkNotificationRead(context.Background(), id)
	if err != nil {
		t.Fatalf("MarkNotificationRead failed: %v", err)
	}

	// Unread only should exclude it
	notifs, _ := uc.GetNotifications(context.Background(), true)
	for _, n := range notifs {
		if n.ID == id {
			t.Error("marked notification should not appear in unread list")
		}
	}
}

func TestDismissNotification(t *testing.T) {
	lic := makeLicense("KEY", "dev-001", 5, domain.LicenseStatusActive)
	repo := &mockLicenseRepo{license: lic}
	engine := &mockLicenseEngine{deviceID: "dev-001", validSig: true}
	uc := newTestLicenseUC(repo, engine)

	_, _ = uc.Validate(context.Background())
	if len(repo.notifications) == 0 {
		t.Fatal("need at least one notification")
	}

	id := repo.notifications[0].ID
	err := uc.DismissNotification(context.Background(), id)
	if err != nil {
		t.Fatalf("DismissNotification failed: %v", err)
	}

	// Should not appear in unread list
	notifs, _ := uc.GetNotifications(context.Background(), true)
	for _, n := range notifs {
		if n.ID == id {
			t.Error("dismissed notification should not appear")
		}
	}
}

// ============================================================
// SECTION 10: GetStatus Tests
// ============================================================

func TestGetStatus_NoLicense(t *testing.T) {
	uc := newTestLicenseUC(&mockLicenseRepo{}, &mockLicenseEngine{})

	status, err := uc.GetStatus(context.Background())
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if !status.IsRestricted {
		t.Error("no license should be restricted")
	}
	if !status.NeedsRenewal {
		t.Error("no license should need renewal")
	}
	if status.DaysRemaining != -999 {
		t.Errorf("days = %d, want -999", status.DaysRemaining)
	}
	if status.License != nil {
		t.Error("license should be nil")
	}
}

func TestGetStatus_ActiveLicense(t *testing.T) {
	lic := makeLicense("KEY", "dev-001", 20, domain.LicenseStatusActive)
	repo := &mockLicenseRepo{license: lic}
	uc := newTestLicenseUC(repo, &mockLicenseEngine{})

	status, err := uc.GetStatus(context.Background())
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.IsRestricted {
		t.Error("active license should not be restricted")
	}
	if status.License == nil {
		t.Fatal("license should not be nil")
	}
	if status.DaysRemaining < 19 || status.DaysRemaining > 21 {
		t.Errorf("days = %d, expected ~20", status.DaysRemaining)
	}
	if !status.NeedsRenewal {
		t.Error("20 days remaining should need renewal (<=30)")
	}
}

func TestGetStatus_ExpiredLicense(t *testing.T) {
	lic := makeLicense("KEY", "dev-001", -40, domain.LicenseStatusExpired)
	repo := &mockLicenseRepo{license: lic}
	uc := newTestLicenseUC(repo, &mockLicenseEngine{})

	status, err := uc.GetStatus(context.Background())
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if !status.IsRestricted {
		t.Error("expired license should be restricted")
	}
	if status.GraceDaysRemaining != 0 {
		t.Errorf("grace days = %d, want 0", status.GraceDaysRemaining)
	}
}

// ============================================================
// SECTION 11: ListEvents Tests
// ============================================================

func TestListEvents_Empty(t *testing.T) {
	lic := makeLicense("KEY", "dev-001", 20, domain.LicenseStatusActive)
	repo := &mockLicenseRepo{license: lic}
	uc := newTestLicenseUC(repo, &mockLicenseEngine{})

	events, total, err := uc.ListEvents(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("ListEvents failed: %v", err)
	}
	if total != 0 {
		t.Errorf("total = %d, want 0", total)
	}
	if len(events) != 0 {
		t.Errorf("events = %d, want 0", len(events))
	}
}

func TestListEvents_WithEvents(t *testing.T) {
	lic := makeLicense("KEY", "dev-001", 20, domain.LicenseStatusActive)
	repo := &mockLicenseRepo{license: lic}
	engine := &mockLicenseEngine{deviceID: "dev-001", validSig: true}
	uc := newTestLicenseUC(repo, engine)

	// Generate some events via validation
	_, _ = uc.Validate(context.Background())

	events, total, err := uc.ListEvents(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("ListEvents failed: %v", err)
	}
	if total == 0 {
		t.Error("expected at least one event after validation")
	}
	if len(events) == 0 {
		t.Error("expected events in result")
	}
}

func TestListEvents_NoLicense(t *testing.T) {
	uc := newTestLicenseUC(&mockLicenseRepo{}, &mockLicenseEngine{})

	events, total, err := uc.ListEvents(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("ListEvents failed: %v", err)
	}
	if total != 0 || len(events) != 0 {
		t.Error("expected empty result when no license")
	}
}

func TestListEvents_Pagination(t *testing.T) {
	lic := makeLicense("KEY", "dev-001", 20, domain.LicenseStatusActive)
	repo := &mockLicenseRepo{license: lic}

	// Add 5 events manually
	for i := 0; i < 5; i++ {
		repo.events = append(repo.events, domain.LicenseEvent{
			ID:        uuid.New(),
			LicenseID: lic.ID,
			EventType: "validated",
		})
	}

	uc := newTestLicenseUC(repo, &mockLicenseEngine{})

	// Page 1, 2 per page
	events, total, _ := uc.ListEvents(context.Background(), 1, 2)
	if total != 5 {
		t.Errorf("total = %d, want 5", total)
	}
	if len(events) != 2 {
		t.Errorf("page 1 events = %d, want 2", len(events))
	}

	// Page 3, 2 per page -> should get 1 event
	events, _, _ = uc.ListEvents(context.Background(), 3, 2)
	if len(events) != 1 {
		t.Errorf("page 3 events = %d, want 1", len(events))
	}
}
