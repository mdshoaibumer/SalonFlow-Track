package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/pkg/uid"
)

// Cloud providers.
const (
	CloudProviderNone     = "none"
	CloudProviderGDrive   = "google_drive"
	CloudProviderS3       = "aws_s3"
	CloudProviderDOSpaces = "digitalocean_spaces"
)

// Cloud backup statuses.
const (
	CloudBackupPending   = "pending"
	CloudBackupUploading = "uploading"
	CloudBackupCompleted = "completed"
	CloudBackupFailed    = "failed"
	CloudBackupRestoring = "restoring"
	CloudBackupRestored  = "restored"
)

// CloudBackupConfig stores cloud backup configuration.
type CloudBackupConfig struct {
	ID                      uuid.UUID `json:"id"`
	Provider                string    `json:"provider"`
	BucketName              string    `json:"bucket_name"`
	Region                  string    `json:"region"`
	AccessKey               string    `json:"access_key"`
	Endpoint                string    `json:"endpoint"`
	EncryptBackups          bool      `json:"encrypt_backups"`
	AutoBackup              bool      `json:"auto_backup"`
	AutoBackupIntervalHours int       `json:"auto_backup_interval_hours"`
	MaxVersions             int       `json:"max_versions"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

// NewCloudBackupConfig creates a default config.
func NewCloudBackupConfig() *CloudBackupConfig {
	now := time.Now().UTC()
	return &CloudBackupConfig{
		ID:                      uid.New(),
		Provider:                CloudProviderNone,
		EncryptBackups:          true,
		AutoBackupIntervalHours: 24,
		MaxVersions:             10,
		CreatedAt:               now,
		UpdatedAt:               now,
	}
}

// CloudBackupHistory records a cloud backup event.
type CloudBackupHistory struct {
	ID           uuid.UUID `json:"id"`
	Provider     string    `json:"provider"`
	FileName     string    `json:"file_name"`
	FileSize     int64     `json:"file_size"`
	RemotePath   string    `json:"remote_path"`
	Status       string    `json:"status"`
	IsEncrypted  bool      `json:"is_encrypted"`
	ErrorMessage string    `json:"error_message"`
	StartedAt    string    `json:"started_at"`
	CompletedAt  string    `json:"completed_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// NewCloudBackupHistory creates a new history entry.
func NewCloudBackupHistory(provider, fileName string, fileSize int64, encrypted bool) *CloudBackupHistory {
	now := time.Now().UTC()
	return &CloudBackupHistory{
		ID:          uid.New(),
		Provider:    provider,
		FileName:    fileName,
		FileSize:    fileSize,
		Status:      CloudBackupPending,
		IsEncrypted: encrypted,
		StartedAt:   now.Format(time.RFC3339),
		CreatedAt:   now,
	}
}

// CloudBackupStats holds cloud backup dashboard stats.
type CloudBackupStats struct {
	LastBackupAt   string `json:"last_backup_at"`
	TotalBackups   int    `json:"total_backups"`
	TotalSizeBytes int64  `json:"total_size_bytes"`
	Provider       string `json:"provider"`
	AutoEnabled    bool   `json:"auto_enabled"`
}
