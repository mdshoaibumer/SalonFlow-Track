package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/pkg/uid"
)

// Backup types.
const (
	BackupTypeManual        = "manual"
	BackupTypeDaily         = "daily"
	BackupTypeBeforeUpdate  = "before_update"
	BackupTypeBeforeRestore = "before_restore"
)

// Backup statuses.
const (
	BackupStatusPending   = "pending"
	BackupStatusCompleted = "completed"
	BackupStatusFailed    = "failed"
	BackupStatusCorrupted = "corrupted"
	BackupStatusVerified  = "verified"
)

// Restore statuses.
const (
	RestoreStatusPending   = "pending"
	RestoreStatusCompleted = "completed"
	RestoreStatusFailed    = "failed"
)

// BackupRecord represents a backup entry in history.
type BackupRecord struct {
	ID           uuid.UUID `json:"id"`
	BackupName   string    `json:"backup_name"`
	BackupType   string    `json:"backup_type"`
	BackupPath   string    `json:"backup_path"`
	FileSize     int64     `json:"file_size"`
	Checksum     string    `json:"checksum"`
	Status       string    `json:"status"`
	ErrorMessage string    `json:"error_message,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// NewBackupRecord creates a new backup record.
func NewBackupRecord(name, backupType, path string) *BackupRecord {
	return &BackupRecord{
		ID:         uid.New(),
		BackupName: name,
		BackupType: backupType,
		BackupPath: path,
		Status:     BackupStatusPending,
		CreatedAt:  time.Now().UTC(),
	}
}

// RestoreRecord represents a restore entry in history.
type RestoreRecord struct {
	ID           uuid.UUID `json:"id"`
	BackupID     uuid.UUID `json:"backup_id"`
	BackupName   string    `json:"backup_name"`
	RestoreDate  string    `json:"restore_date"`
	Status       string    `json:"status"`
	Notes        string    `json:"notes"`
	ErrorMessage string    `json:"error_message,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// NewRestoreRecord creates a new restore record.
func NewRestoreRecord(backupID uuid.UUID, backupName, notes string) *RestoreRecord {
	return &RestoreRecord{
		ID:          uid.New(),
		BackupID:    backupID,
		BackupName:  backupName,
		RestoreDate: time.Now().Format("2006-01-02 15:04:05"),
		Status:      RestoreStatusPending,
		Notes:       notes,
		CreatedAt:   time.Now().UTC(),
	}
}

// BackupStats holds summary info for the backup dashboard.
type BackupStats struct {
	TotalBackups   int    `json:"total_backups"`
	LastBackupName string `json:"last_backup_name"`
	LastBackupDate string `json:"last_backup_date"`
	LastBackupSize int64  `json:"last_backup_size"`
	LastStatus     string `json:"last_status"`
	TotalRestores  int    `json:"total_restores"`
}

// BackupVerification holds the result of a backup verification.
type BackupVerification struct {
	BackupID     string `json:"backup_id"`
	FileExists   bool   `json:"file_exists"`
	CanOpen      bool   `json:"can_open"`
	IntegrityOK  bool   `json:"integrity_ok"`
	ChecksumOK   bool   `json:"checksum_ok"`
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message,omitempty"`
}
