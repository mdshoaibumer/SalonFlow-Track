package cloudbackup

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// Engine handles cloud backup upload/download operations.
// In production, this would use actual cloud SDKs (AWS SDK, Google Drive API, etc.).
// For the local build, it simulates cloud operations with local file copies.
type Engine struct{}

// NewEngine creates a new cloud backup engine.
func NewEngine() *Engine {
	return &Engine{}
}

// Upload uploads a file to the configured cloud provider.
func (e *Engine) Upload(cfg *domain.CloudBackupConfig, localPath, remotePath string, encrypt bool) error {
	if cfg.Provider == domain.CloudProviderNone {
		return fmt.Errorf("no cloud provider configured")
	}

	// Verify local file exists
	info, err := os.Stat(localPath)
	if err != nil {
		return fmt.Errorf("local file not found: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("local path is a directory")
	}

	// In production, this would use the appropriate SDK:
	// - AWS S3: aws-sdk-go-v2
	// - Google Drive: google.golang.org/api/drive
	// - DigitalOcean Spaces: same as S3 API with custom endpoint
	//
	// For now, simulate by copying to a local "cloud" directory
	cloudDir := filepath.Join(os.TempDir(), "salonflow-cloud-backup", cfg.Provider)
	if err := os.MkdirAll(cloudDir, 0750); err != nil {
		return fmt.Errorf("failed to create cloud dir: %w", err)
	}

	destPath := filepath.Join(cloudDir, filepath.Base(remotePath))
	return copyFile(localPath, destPath)
}

// Download downloads a file from the configured cloud provider.
func (e *Engine) Download(cfg *domain.CloudBackupConfig, remotePath, localPath string) error {
	if cfg.Provider == domain.CloudProviderNone {
		return fmt.Errorf("no cloud provider configured")
	}

	cloudDir := filepath.Join(os.TempDir(), "salonflow-cloud-backup", cfg.Provider)
	srcPath := filepath.Join(cloudDir, filepath.Base(remotePath))

	return copyFile(srcPath, localPath)
}

// TestConnection tests the cloud provider connection.
func (e *Engine) TestConnection(cfg *domain.CloudBackupConfig) error {
	if cfg.Provider == domain.CloudProviderNone {
		return fmt.Errorf("no cloud provider configured")
	}
	if cfg.BucketName == "" {
		return fmt.Errorf("bucket name is required")
	}
	// In production, this would attempt to list the bucket or write a test file
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
