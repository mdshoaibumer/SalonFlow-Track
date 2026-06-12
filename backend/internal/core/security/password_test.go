package security

import (
	"testing"
)

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{"too short", "Ab1!", ErrPasswordTooShort},
		{"no uppercase", "abcdefg1!", ErrPasswordNoUpper},
		{"no lowercase", "ABCDEFG1!", ErrPasswordNoLower},
		{"no digit", "Abcdefg!@", ErrPasswordNoDigit},
		{"no special", "Abcdefg12", ErrPasswordNoSpecial},
		{"valid", "Admin@123", nil},
		{"valid complex", "P@ssw0rd!XYZ", nil},
		{"valid min length", "Ab1!efgh", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswordStrength(tt.password)
			if err != tt.wantErr {
				t.Errorf("ValidatePasswordStrength(%q) = %v, want %v", tt.password, err, tt.wantErr)
			}
		})
	}
}

func TestHashAndVerifyPassword(t *testing.T) {
	password := "Admin@123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if hash == "" {
		t.Fatal("HashPassword() returned empty hash")
	}

	// Verify correct password
	valid, err := VerifyPassword(password, hash)
	if err != nil {
		t.Fatalf("VerifyPassword() error = %v", err)
	}
	if !valid {
		t.Error("VerifyPassword() returned false for correct password")
	}

	// Verify incorrect password
	valid, err = VerifyPassword("WrongPass1!", hash)
	if err != nil {
		t.Fatalf("VerifyPassword() error = %v", err)
	}
	if valid {
		t.Error("VerifyPassword() returned true for incorrect password")
	}
}

func TestHashPasswordUniqueness(t *testing.T) {
	password := "Admin@123"
	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)

	if hash1 == hash2 {
		t.Error("HashPassword() should produce unique hashes due to random salt")
	}
}

func TestVerifyPasswordInvalidHash(t *testing.T) {
	_, err := VerifyPassword("test", "invalid-hash")
	if err != ErrInvalidHash {
		t.Errorf("VerifyPassword() with invalid hash = %v, want ErrInvalidHash", err)
	}
}

func TestGenerateToken(t *testing.T) {
	token1, err := GenerateToken()
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}
	if token1 == "" {
		t.Fatal("GenerateToken() returned empty token")
	}

	token2, _ := GenerateToken()
	if token1 == token2 {
		t.Error("GenerateToken() should produce unique tokens")
	}
}

func TestHashToken(t *testing.T) {
	token := "some-random-token-value"
	hash1 := HashToken(token)
	hash2 := HashToken(token)

	if hash1 != hash2 {
		t.Error("HashToken() should be deterministic")
	}

	hash3 := HashToken("different-token")
	if hash1 == hash3 {
		t.Error("HashToken() should produce different hashes for different tokens")
	}
}

func BenchmarkHashPassword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HashPassword("Admin@123")
	}
}

func BenchmarkVerifyPassword(b *testing.B) {
	hash, _ := HashPassword("Admin@123")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		VerifyPassword("Admin@123", hash)
	}
}
