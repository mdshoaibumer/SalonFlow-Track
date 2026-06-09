package license

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"strings"
)

// secretKey is used for HMAC signing. In production this should be embedded at build time.
var secretKey = []byte("salonflow-license-signing-key-v1-prod")

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
