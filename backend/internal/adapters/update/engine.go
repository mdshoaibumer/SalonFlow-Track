package update

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/salonflow/salonflow-track/internal/core/domain"
)

const (
	defaultUpdateURL = "https://updates.salonflow.in/api/v1/version.json"
	updateDirName    = "Updates"
)

// Engine implements ports.UpdateEngine.
type Engine struct {
	updateURL string
	updateDir string
	client    *http.Client
}

// NewEngine creates a new update engine with default settings.
func NewEngine() *Engine {
	return &Engine{
		updateURL: defaultUpdateURL,
		updateDir: defaultUpdateDir(),
		client:    &http.Client{Timeout: 30 * time.Second},
	}
}

// NewEngineWithConfig creates an update engine with custom settings (for testing).
func NewEngineWithConfig(updateURL, updateDir string) *Engine {
	return &Engine{
		updateURL: updateURL,
		updateDir: updateDir,
		client:    &http.Client{Timeout: 30 * time.Second},
	}
}

// CheckForUpdate checks the update server for a new version.
func (e *Engine) CheckForUpdate(currentVersion string) (*domain.VersionInfo, error) {
	resp, err := e.client.Get(e.updateURL)
	if err != nil {
		return nil, fmt.Errorf("check update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("update server returned %d", resp.StatusCode)
	}

	var info domain.VersionInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("decode version info: %w", err)
	}

	// Compare versions - simple string comparison works for semver
	if info.Version <= currentVersion {
		return nil, nil // No update available
	}

	return &info, nil
}

// DownloadUpdate downloads the update file to the specified path.
func (e *Engine) DownloadUpdate(info *domain.VersionInfo, destPath string) error {
	if err := os.MkdirAll(filepath.Dir(destPath), 0750); err != nil {
		return fmt.Errorf("create update dir: %w", err)
	}

	resp, err := e.client.Get(info.DownloadURL)
	if err != nil {
		return fmt.Errorf("download update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned %d", resp.StatusCode)
	}

	f, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		os.Remove(destPath)
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

// VerifyDownload verifies the SHA-256 checksum of a downloaded file.
func (e *Engine) VerifyDownload(filePath, expectedChecksum string) (bool, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return false, fmt.Errorf("hash file: %w", err)
	}

	actual := hex.EncodeToString(h.Sum(nil))
	return actual == expectedChecksum, nil
}

// UpdateDir returns the directory where updates are stored.
func (e *Engine) UpdateDir() string {
	return e.updateDir
}

func defaultUpdateDir() string {
	var base string
	switch runtime.GOOS {
	case "windows":
		base = os.Getenv("LOCALAPPDATA")
		if base == "" {
			base = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local")
		}
	default:
		base, _ = os.UserHomeDir()
	}
	return filepath.Join(base, "SalonFlowTrack", updateDirName)
}
