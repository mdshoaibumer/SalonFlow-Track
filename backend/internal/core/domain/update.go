package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/pkg/uid"
)

// Update statuses.
const (
	UpdateStatusAvailable   = "available"
	UpdateStatusDownloading = "downloading"
	UpdateStatusDownloaded  = "downloaded"
	UpdateStatusInstalling  = "installing"
	UpdateStatusInstalled   = "installed"
	UpdateStatusFailed      = "failed"
	UpdateStatusRolledBack  = "rolled_back"
)

// Update history statuses.
const (
	UpdateHistoryPending     = "pending"
	UpdateHistoryDownloading = "downloading"
	UpdateHistoryDownloaded  = "downloaded"
	UpdateHistoryInstalling  = "installing"
	UpdateHistoryCompleted   = "completed"
	UpdateHistoryFailed      = "failed"
	UpdateHistoryRolledBack  = "rolled_back"
)

// AppVersion represents a software version record.
type AppVersion struct {
	ID           uuid.UUID `json:"id"`
	Version      string    `json:"version"`
	ReleaseDate  string    `json:"release_date"`
	ReleaseNotes string    `json:"release_notes"`
	InstalledAt  string    `json:"installed_at,omitempty"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

// NewAppVersion creates a new app version record.
func NewAppVersion(version, releaseDate, releaseNotes string) *AppVersion {
	return &AppVersion{
		ID:           uid.New(),
		Version:      version,
		ReleaseDate:  releaseDate,
		ReleaseNotes: releaseNotes,
		Status:       UpdateStatusAvailable,
		CreatedAt:    time.Now().UTC(),
	}
}

// UpdateRecord represents an update history entry.
type UpdateRecord struct {
	ID           uuid.UUID `json:"id"`
	FromVersion  string    `json:"from_version"`
	ToVersion    string    `json:"to_version"`
	UpdateDate   string    `json:"update_date"`
	Status       string    `json:"status"`
	ErrorMessage string    `json:"error_message,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// NewUpdateRecord creates a new update record.
func NewUpdateRecord(fromVersion, toVersion string) *UpdateRecord {
	now := time.Now().UTC()
	return &UpdateRecord{
		ID:          uid.New(),
		FromVersion: fromVersion,
		ToVersion:   toVersion,
		UpdateDate:  now.Format(time.RFC3339),
		Status:      UpdateHistoryPending,
		CreatedAt:   now,
	}
}

// VersionInfo is what the update server returns.
type VersionInfo struct {
	Version      string `json:"version"`
	ReleaseDate  string `json:"release_date"`
	DownloadURL  string `json:"download_url"`
	Checksum     string `json:"checksum"`
	Mandatory    bool   `json:"mandatory"`
	ReleaseNotes string `json:"release_notes"`
}

// UpdateStatus represents the current update state for the dashboard.
type UpdateStatus struct {
	CurrentVersion  string `json:"current_version"`
	LatestVersion   string `json:"latest_version,omitempty"`
	UpdateAvailable bool   `json:"update_available"`
	Status          string `json:"status"`
	ReleaseNotes    string `json:"release_notes,omitempty"`
}
