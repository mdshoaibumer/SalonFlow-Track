package license

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// ============================================================
// Key Generation Tests
// ============================================================

func TestGenerateKey_Format(t *testing.T) {
	e := NewEngine()
	key := e.GenerateKey()

	if !strings.HasPrefix(key, "SALONFLOW-") {
		t.Errorf("key should start with SALONFLOW-, got %q", key)
	}

	parts := strings.Split(key, "-")
	if len(parts) != 4 {
		t.Errorf("key should have 4 parts (SALONFLOW-XXXX-XXXX-XXXX), got %d", len(parts))
	}

	for i := 1; i < len(parts); i++ {
		if len(parts[i]) != 4 {
			t.Errorf("segment %d length = %d, want 4", i, len(parts[i]))
		}
	}
}

func TestGenerateKey_Unique(t *testing.T) {
	e := NewEngine()
	keys := make(map[string]bool)
	for i := 0; i < 100; i++ {
		k := e.GenerateKey()
		if keys[k] {
			t.Fatalf("duplicate key generated: %s", k)
		}
		keys[k] = true
	}
}

func TestGenerateKey_ValidCharacters(t *testing.T) {
	e := NewEngine()
	const validChars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	for i := 0; i < 50; i++ {
		key := e.GenerateKey()
		parts := strings.Split(key, "-")
		for _, seg := range parts[1:] {
			for _, ch := range seg {
				if !strings.ContainsRune(validChars, ch) {
					t.Errorf("invalid char %c in key segment %q", ch, seg)
				}
			}
		}
	}
}

// ============================================================
// Device ID Tests
// ============================================================

func TestGenerateDeviceID_Length(t *testing.T) {
	e := NewEngine()
	id := e.GenerateDeviceID()

	if len(id) != 32 {
		t.Errorf("device ID length = %d, want 32 (16 bytes hex)", len(id))
	}
}

func TestGenerateDeviceID_Deterministic(t *testing.T) {
	e := NewEngine()
	id1 := e.GenerateDeviceID()
	id2 := e.GenerateDeviceID()

	if id1 != id2 {
		t.Error("device ID should be deterministic on same machine")
	}
}

func TestGenerateDeviceID_HexEncoded(t *testing.T) {
	e := NewEngine()
	id := e.GenerateDeviceID()

	const hexChars = "0123456789abcdef"
	for _, ch := range id {
		if !strings.ContainsRune(hexChars, ch) {
			t.Errorf("device ID should be hex, found char %c", ch)
			break
		}
	}
}

// ============================================================
// Signature Tests
// ============================================================

func TestSignLicense_NotEmpty(t *testing.T) {
	e := NewEngine()
	sig := e.SignLicense("KEY", "2026-12-31", "device")
	if sig == "" {
		t.Fatal("signature should not be empty")
	}
}

func TestSignLicense_Deterministic(t *testing.T) {
	e := NewEngine()
	sig1 := e.SignLicense("KEY", "2026-12-31", "device")
	sig2 := e.SignLicense("KEY", "2026-12-31", "device")
	if sig1 != sig2 {
		t.Error("same inputs should produce same signature")
	}
}

func TestSignLicense_DifferentInputs_DifferentSignatures(t *testing.T) {
	e := NewEngine()
	sig1 := e.SignLicense("KEY-A", "2026-12-31", "device")
	sig2 := e.SignLicense("KEY-B", "2026-12-31", "device")
	if sig1 == sig2 {
		t.Error("different keys should produce different signatures")
	}
}

func TestValidateSignature_Valid(t *testing.T) {
	e := NewEngine()
	key := "SALONFLOW-ABCD-EFGH-IJKL"
	expiry := "2026-12-31"
	device := "dev-001"

	sig := e.SignLicense(key, expiry, device)
	if !e.ValidateSignature(key, expiry, device, sig) {
		t.Error("valid signature should validate")
	}
}

func TestValidateSignature_TamperedKey(t *testing.T) {
	e := NewEngine()
	sig := e.SignLicense("REAL-KEY", "2026-12-31", "device")
	if e.ValidateSignature("FAKE-KEY", "2026-12-31", "device", sig) {
		t.Error("tampered key should fail validation")
	}
}

func TestValidateSignature_TamperedExpiry(t *testing.T) {
	e := NewEngine()
	sig := e.SignLicense("KEY", "2026-12-31", "device")
	if e.ValidateSignature("KEY", "2030-12-31", "device", sig) {
		t.Error("tampered expiry should fail validation")
	}
}

func TestValidateSignature_TamperedDevice(t *testing.T) {
	e := NewEngine()
	sig := e.SignLicense("KEY", "2026-12-31", "device-A")
	if e.ValidateSignature("KEY", "2026-12-31", "device-B", sig) {
		t.Error("tampered device should fail validation")
	}
}

func TestValidateSignature_InvalidSignature(t *testing.T) {
	e := NewEngine()
	if e.ValidateSignature("KEY", "2026-12-31", "device", "totally-wrong") {
		t.Error("invalid signature should fail validation")
	}
}

func TestValidateSignature_EmptySignature(t *testing.T) {
	e := NewEngine()
	if e.ValidateSignature("KEY", "2026-12-31", "device", "") {
		t.Error("empty signature should fail validation")
	}
}

// ============================================================
// License File Export/Import Tests (Encrypted)
// ============================================================

func TestExportAndParseLicenseFile_RoundTrip(t *testing.T) {
	e := NewEngine()

	lic := &domain.License{
		ID:           uuid.New(),
		LicenseKey:   "SALONFLOW-RNDTRP-TEST",
		CustomerName: "Round Trip Customer",
		SalonName:    "Round Trip Salon",
		DeviceID:     "device-xyz",
		IssuedDate:   "2026-01-15",
		ExpiryDate:   "2026-02-15",
		Signature:    e.SignLicense("SALONFLOW-RNDTRP-TEST", "2026-02-15", "device-xyz"),
	}

	// Export (encrypt)
	data, err := e.ExportLicenseFile(lic)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("exported data is empty")
	}

	// Should NOT be plain JSON
	var probe domain.LicenseFileData
	if json.Unmarshal(data, &probe) == nil && probe.LicenseKey != "" {
		t.Error("exported data should be encrypted, not plain JSON")
	}

	// Parse (decrypt)
	parsed, err := e.ParseLicenseFile(data)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if parsed.LicenseKey != "SALONFLOW-RNDTRP-TEST" {
		t.Errorf("key = %q, want SALONFLOW-RNDTRP-TEST", parsed.LicenseKey)
	}
	if parsed.CustomerName != "Round Trip Customer" {
		t.Errorf("customer = %q", parsed.CustomerName)
	}
	if parsed.SalonName != "Round Trip Salon" {
		t.Errorf("salon = %q", parsed.SalonName)
	}
	if parsed.DeviceID != "device-xyz" {
		t.Errorf("device = %q", parsed.DeviceID)
	}
	if parsed.IssuedDate != "2026-01-15" {
		t.Errorf("issued = %q", parsed.IssuedDate)
	}
	if parsed.ExpiryDate != "2026-02-15" {
		t.Errorf("expiry = %q", parsed.ExpiryDate)
	}
	if parsed.Signature != lic.Signature {
		t.Error("signature mismatch after round trip")
	}
}

func TestExportLicenseFile_DifferentOutputEachTime(t *testing.T) {
	e := NewEngine()
	lic := &domain.License{
		ID:         uuid.New(),
		LicenseKey: "KEY",
		ExpiryDate: "2026-12-31",
		DeviceID:   "dev",
		Signature:  "sig",
	}

	data1, _ := e.ExportLicenseFile(lic)
	data2, _ := e.ExportLicenseFile(lic)

	// Due to random nonce, encrypted outputs should differ
	if string(data1) == string(data2) {
		t.Error("two exports should produce different ciphertext (different nonce)")
	}
}

// ============================================================
// License File Parse Tests (Plain JSON fallback)
// ============================================================

func TestParseLicenseFile_PlainJSON(t *testing.T) {
	e := NewEngine()

	jsonData := []byte(`{
		"license_key": "PLAIN-KEY",
		"customer_name": "Plain Customer",
		"salon_name": "Plain Salon",
		"issued_date": "2026-03-01",
		"expiry_date": "2026-04-01",
		"device_id": "plain-dev",
		"signature": "plain-sig"
	}`)

	parsed, err := e.ParseLicenseFile(jsonData)
	if err != nil {
		t.Fatalf("parse plain JSON failed: %v", err)
	}
	if parsed.LicenseKey != "PLAIN-KEY" {
		t.Errorf("key = %q", parsed.LicenseKey)
	}
	if parsed.CustomerName != "Plain Customer" {
		t.Errorf("customer = %q", parsed.CustomerName)
	}
	if parsed.Signature != "plain-sig" {
		t.Errorf("signature = %q", parsed.Signature)
	}
}

func TestParseLicenseFile_EmptyInput(t *testing.T) {
	e := NewEngine()
	_, err := e.ParseLicenseFile([]byte{})
	if err == nil {
		t.Fatal("expected error for empty input")
	}
}

func TestParseLicenseFile_GarbageInput(t *testing.T) {
	e := NewEngine()
	_, err := e.ParseLicenseFile([]byte("this is not json and not encrypted"))
	if err == nil {
		t.Fatal("expected error for garbage input")
	}
}

func TestParseLicenseFile_TruncatedEncrypted(t *testing.T) {
	e := NewEngine()
	// Export to get valid encrypted data, then truncate it
	lic := &domain.License{
		ID:         uuid.New(),
		LicenseKey: "KEY",
		ExpiryDate: "2026-12-31",
	}
	data, _ := e.ExportLicenseFile(lic)
	truncated := data[:5] // Too short to be valid

	_, err := e.ParseLicenseFile(truncated)
	if err == nil {
		t.Fatal("expected error for truncated data")
	}
}

func TestParseLicenseFile_CorruptedEncrypted(t *testing.T) {
	e := NewEngine()
	lic := &domain.License{
		ID:         uuid.New(),
		LicenseKey: "KEY",
		ExpiryDate: "2026-12-31",
	}
	data, _ := e.ExportLicenseFile(lic)

	// Corrupt a byte in the middle
	if len(data) > 20 {
		data[15] ^= 0xFF
	}

	_, err := e.ParseLicenseFile(data)
	if err == nil {
		t.Fatal("expected error for corrupted encrypted data")
	}
}
