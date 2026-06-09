package license

import (
	"strings"
	"testing"
)

func TestEngine_GenerateKey(t *testing.T) {
	e := NewEngine()
	key := e.GenerateKey()

	if !strings.HasPrefix(key, "SALONFLOW-") {
		t.Errorf("key should start with SALONFLOW-, got %q", key)
	}

	parts := strings.Split(key, "-")
	if len(parts) != 4 {
		t.Errorf("key should have 4 parts (SALONFLOW-XXXX-XXXX-XXXX), got %d parts", len(parts))
	}

	// Each segment after SALONFLOW should be 4 chars
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) != 4 {
			t.Errorf("segment %d length = %d, want 4", i, len(parts[i]))
		}
	}

	// Keys should be unique
	key2 := e.GenerateKey()
	if key == key2 {
		t.Error("two generated keys should not be identical")
	}
}

func TestEngine_GenerateDeviceID(t *testing.T) {
	e := NewEngine()
	id := e.GenerateDeviceID()

	if len(id) != 32 { // 16 bytes hex encoded
		t.Errorf("device ID length = %d, want 32", len(id))
	}

	// Should be deterministic for same machine
	id2 := e.GenerateDeviceID()
	if id != id2 {
		t.Error("device ID should be deterministic")
	}
}

func TestEngine_SignAndValidate(t *testing.T) {
	e := NewEngine()

	key := "SALONFLOW-ABCD-EFGH-IJKL"
	expiry := "2025-12-31"
	device := "abc123"

	sig := e.SignLicense(key, expiry, device)
	if sig == "" {
		t.Fatal("signature should not be empty")
	}

	// Valid signature
	if !e.ValidateSignature(key, expiry, device, sig) {
		t.Error("valid signature should validate")
	}

	// Tampered key
	if e.ValidateSignature("SALONFLOW-XXXX-XXXX-XXXX", expiry, device, sig) {
		t.Error("tampered key should not validate")
	}

	// Tampered expiry
	if e.ValidateSignature(key, "2030-12-31", device, sig) {
		t.Error("tampered expiry should not validate")
	}

	// Tampered device
	if e.ValidateSignature(key, expiry, "different_device", sig) {
		t.Error("tampered device should not validate")
	}

	// Tampered signature
	if e.ValidateSignature(key, expiry, device, "invalid_sig") {
		t.Error("invalid signature should not validate")
	}
}
