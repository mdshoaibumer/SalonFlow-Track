package license

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// secretKey is used for HMAC signing. In production this should be embedded at build time.
var secretKey = []byte("salonflow-license-signing-key-v1-prod")

// encryptionKey is used for license file encryption (AES-256).
var encryptionKey = sha256.Sum256([]byte("salonflow-license-file-encryption-key-v1"))

// Engine implements ports.LicenseEngine.
type Engine struct{}

// NewEngine creates a new license engine.
func NewEngine() *Engine {
	return &Engine{}
}

// GenerateKey generates a license key in format SALONFLOW-XXXX-XXXX-XXXX.
func (e *Engine) GenerateKey() string {
	return fmt.Sprintf("SALONFLOW-%s-%s-%s", randomSegment(4), randomSegment(4), randomSegment(4))
}

// GenerateDeviceID generates a unique device identifier based on machine characteristics.
func (e *Engine) GenerateDeviceID() string {
	hostname, _ := os.Hostname()
	data := fmt.Sprintf("%s|%s|%s|%d", hostname, runtime.GOOS, runtime.GOARCH, runtime.NumCPU())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16])
}

// SignLicense creates an HMAC-SHA256 signature for the license data.
func (e *Engine) SignLicense(key, expiryDate, deviceID string) string {
	payload := fmt.Sprintf("%s|%s|%s", key, expiryDate, deviceID)
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

// ValidateSignature verifies the HMAC-SHA256 signature for the license data.
func (e *Engine) ValidateSignature(key, expiryDate, deviceID, signature string) bool {
	expected := e.SignLicense(key, expiryDate, deviceID)
	return hmac.Equal([]byte(expected), []byte(signature))
}

// ParseLicenseFile decrypts and parses a license file into LicenseFileData.
func (e *Engine) ParseLicenseFile(data []byte) (*domain.LicenseFileData, error) {
	if len(data) == 0 {
		return nil, errors.New("empty license file")
	}

	// Try decrypting first (encrypted file)
	decrypted, err := e.decrypt(data)
	if err != nil {
		// Fallback: try as plain JSON
		var lfd domain.LicenseFileData
		if jsonErr := json.Unmarshal(data, &lfd); jsonErr != nil {
			return nil, errors.New("invalid license file format")
		}
		return &lfd, nil
	}

	var lfd domain.LicenseFileData
	if err := json.Unmarshal(decrypted, &lfd); err != nil {
		return nil, errors.New("corrupted license file data")
	}
	return &lfd, nil
}

// ExportLicenseFile serializes and encrypts a license into a file format.
func (e *Engine) ExportLicenseFile(lic *domain.License) ([]byte, error) {
	lfd := domain.LicenseFileData{
		LicenseKey:   lic.LicenseKey,
		CustomerName: lic.CustomerName,
		SalonName:    lic.SalonName,
		IssuedDate:   lic.IssuedDate,
		ExpiryDate:   lic.ExpiryDate,
		DeviceID:     lic.DeviceID,
		Signature:    lic.Signature,
	}
	plaintext, err := json.Marshal(lfd)
	if err != nil {
		return nil, fmt.Errorf("marshal license: %w", err)
	}
	return e.encrypt(plaintext)
}

func (e *Engine) encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(encryptionKey[:])
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return aesGCM.Seal(nonce, nonce, plaintext, nil), nil
}

func (e *Engine) decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(encryptionKey[:])
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return aesGCM.Open(nil, nonce, ciphertext, nil)
}

// randomSegment generates a random alphanumeric segment of given length (uppercase).
func randomSegment(length int) string {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		// Fallback - should never happen
		return strings.Repeat("X", length)
	}
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}
