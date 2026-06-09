package sqlite

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// CloudBackupRepository is the SQLite implementation.
type CloudBackupRepository struct {
	db  *sql.DB
	log *slog.Logger
}

// NewCloudBackupRepository creates a new CloudBackupRepository.
func NewCloudBackupRepository(db *sql.DB, log *slog.Logger) *CloudBackupRepository {
	return &CloudBackupRepository{db: db, log: log}
}

// GetConfig retrieves cloud backup config.
func (r *CloudBackupRepository) GetConfig(ctx context.Context) (*domain.CloudBackupConfig, error) {
	var cfg domain.CloudBackupConfig
	var encrypt, auto int
	var createdAt, updatedAt string

	err := r.db.QueryRowContext(ctx, `
		SELECT id, provider, bucket_name, region, access_key, endpoint, encrypt_backups, auto_backup, auto_backup_interval_hours, max_versions, created_at, updated_at
		FROM cloud_backup_config LIMIT 1`).
		Scan(&cfg.ID, &cfg.Provider, &cfg.BucketName, &cfg.Region, &cfg.AccessKey, &cfg.Endpoint,
			&encrypt, &auto, &cfg.AutoBackupIntervalHours, &cfg.MaxVersions, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("cloud_backup_config", "default")
	}
	if err != nil {
		return nil, apperror.Database("get_cloud_config", err)
	}

	cfg.EncryptBackups = encrypt == 1
	cfg.AutoBackup = auto == 1
	cfg.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	cfg.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return &cfg, nil
}

// SaveConfig creates or updates the config.
func (r *CloudBackupRepository) SaveConfig(ctx context.Context, cfg *domain.CloudBackupConfig) error {
	cfg.UpdatedAt = time.Now().UTC()
	encrypt := 0
	if cfg.EncryptBackups {
		encrypt = 1
	}
	auto := 0
	if cfg.AutoBackup {
		auto = 1
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO cloud_backup_config (id, provider, bucket_name, region, access_key, endpoint, encrypt_backups, auto_backup, auto_backup_interval_hours, max_versions, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			provider=excluded.provider, bucket_name=excluded.bucket_name, region=excluded.region,
			access_key=excluded.access_key, endpoint=excluded.endpoint,
			encrypt_backups=excluded.encrypt_backups, auto_backup=excluded.auto_backup,
			auto_backup_interval_hours=excluded.auto_backup_interval_hours,
			max_versions=excluded.max_versions, updated_at=excluded.updated_at`,
		cfg.ID, cfg.Provider, cfg.BucketName, cfg.Region, cfg.AccessKey, cfg.Endpoint,
		encrypt, auto, cfg.AutoBackupIntervalHours, cfg.MaxVersions,
		cfg.CreatedAt.Format(time.RFC3339), cfg.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("save_cloud_config", err)
	}
	return nil
}

// CreateHistory inserts a history entry.
func (r *CloudBackupRepository) CreateHistory(ctx context.Context, h *domain.CloudBackupHistory) error {
	isEncrypted := 0
	if h.IsEncrypted {
		isEncrypted = 1
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO cloud_backup_history (id, provider, file_name, file_size, remote_path, status, is_encrypted, error_message, started_at, completed_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		h.ID, h.Provider, h.FileName, h.FileSize, h.RemotePath, h.Status, isEncrypted,
		h.ErrorMessage, h.StartedAt, h.CompletedAt, h.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("create_cloud_history", err)
	}
	return nil
}

// UpdateHistoryStatus updates a history entry's status.
func (r *CloudBackupRepository) UpdateHistoryStatus(ctx context.Context, id uuid.UUID, status, errorMsg string) error {
	completedAt := ""
	if status == domain.CloudBackupCompleted || status == domain.CloudBackupRestored {
		completedAt = time.Now().UTC().Format(time.RFC3339)
	}
	_, err := r.db.ExecContext(ctx, `UPDATE cloud_backup_history SET status=?, error_message=?, completed_at=? WHERE id=?`,
		status, errorMsg, completedAt, id)
	if err != nil {
		return apperror.Database("update_cloud_history_status", err)
	}
	return nil
}

// ListHistory lists cloud backup history.
func (r *CloudBackupRepository) ListHistory(ctx context.Context, limit, offset int) ([]domain.CloudBackupHistory, int, error) {
	var total int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM cloud_backup_history`).Scan(&total)
	if err != nil {
		return nil, 0, apperror.Database("list_cloud_history_count", err)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, provider, file_name, file_size, remote_path, status, is_encrypted, error_message, started_at, completed_at, created_at
		FROM cloud_backup_history ORDER BY created_at DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, 0, apperror.Database("list_cloud_history", err)
	}
	defer rows.Close()

	var history []domain.CloudBackupHistory
	for rows.Next() {
		var h domain.CloudBackupHistory
		var isEncrypted int
		var createdAt string
		if err := rows.Scan(&h.ID, &h.Provider, &h.FileName, &h.FileSize, &h.RemotePath, &h.Status, &isEncrypted, &h.ErrorMessage, &h.StartedAt, &h.CompletedAt, &createdAt); err != nil {
			return nil, 0, apperror.Database("list_cloud_history_scan", err)
		}
		h.IsEncrypted = isEncrypted == 1
		h.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		history = append(history, h)
	}
	return history, total, nil
}

// GetStats gets cloud backup stats.
func (r *CloudBackupRepository) GetStats(ctx context.Context) (*domain.CloudBackupStats, error) {
	var stats domain.CloudBackupStats

	// Get config info
	cfg, err := r.GetConfig(ctx)
	if err == nil {
		stats.Provider = cfg.Provider
		stats.AutoEnabled = cfg.AutoBackup
	}

	// Get totals
	r.db.QueryRowContext(ctx, `
		SELECT COALESCE(COUNT(*), 0), COALESCE(SUM(file_size), 0)
		FROM cloud_backup_history WHERE status='completed'`).
		Scan(&stats.TotalBackups, &stats.TotalSizeBytes)

	// Get last backup time
	r.db.QueryRowContext(ctx, `
		SELECT COALESCE(completed_at, '') FROM cloud_backup_history WHERE status='completed' ORDER BY created_at DESC LIMIT 1`).
		Scan(&stats.LastBackupAt)

	return &stats, nil
}
