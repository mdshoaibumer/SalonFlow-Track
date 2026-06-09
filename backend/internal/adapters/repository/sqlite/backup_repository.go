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

// BackupRepository implements ports.BackupRepository.
type BackupRepository struct {
	db  *sql.DB
	log *slog.Logger
}

// NewBackupRepository creates a new BackupRepository.
func NewBackupRepository(db *sql.DB, log *slog.Logger) *BackupRepository {
	return &BackupRepository{db: db, log: log}
}

func (r *BackupRepository) CreateBackupRecord(ctx context.Context, record *domain.BackupRecord) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO backup_history (id, backup_name, backup_type, backup_path, file_size, checksum, status, error_message, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		record.ID.String(), record.BackupName, record.BackupType, record.BackupPath,
		record.FileSize, record.Checksum, record.Status, record.ErrorMessage,
		record.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("create backup record", err)
	}
	return nil
}

func (r *BackupRepository) UpdateBackupRecord(ctx context.Context, record *domain.BackupRecord) error {
	_, err := r.db.ExecContext(ctx, `UPDATE backup_history SET file_size = ?, checksum = ?, status = ?, error_message = ? WHERE id = ?`,
		record.FileSize, record.Checksum, record.Status, record.ErrorMessage, record.ID.String())
	if err != nil {
		return apperror.Database("update backup record", err)
	}
	return nil
}

func (r *BackupRepository) GetBackupByID(ctx context.Context, id uuid.UUID) (*domain.BackupRecord, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, backup_name, backup_type, backup_path, file_size, checksum, status, error_message, created_at
		FROM backup_history WHERE id = ?`, id.String())
	rec, err := r.scanBackup(row)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("backup", id.String())
	}
	if err != nil {
		return nil, apperror.Database("get backup", err)
	}
	return rec, nil
}

func (r *BackupRepository) ListBackups(ctx context.Context, limit, offset int) ([]domain.BackupRecord, int, error) {
	var total int
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM backup_history`).Scan(&total)

	rows, err := r.db.QueryContext(ctx, `SELECT id, backup_name, backup_type, backup_path, file_size, checksum, status, error_message, created_at
		FROM backup_history ORDER BY created_at DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, 0, apperror.Database("list backups", err)
	}
	defer rows.Close()

	var records []domain.BackupRecord
	for rows.Next() {
		var rec domain.BackupRecord
		var idStr, createdAt string
		err := rows.Scan(&idStr, &rec.BackupName, &rec.BackupType, &rec.BackupPath,
			&rec.FileSize, &rec.Checksum, &rec.Status, &rec.ErrorMessage, &createdAt)
		if err != nil {
			return nil, 0, apperror.Database("scan backup", err)
		}
		rec.ID, _ = uuid.Parse(idStr)
		rec.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		records = append(records, rec)
	}
	return records, total, nil
}

func (r *BackupRepository) DeleteBackupRecord(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM backup_history WHERE id = ?`, id.String())
	if err != nil {
		return apperror.Database("delete backup", err)
	}
	return nil
}

func (r *BackupRepository) GetBackupStats(ctx context.Context) (*domain.BackupStats, error) {
	stats := &domain.BackupStats{}
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM backup_history`).Scan(&stats.TotalBackups)
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM restore_history`).Scan(&stats.TotalRestores)

	row := r.db.QueryRowContext(ctx, `SELECT backup_name, created_at, file_size, status FROM backup_history ORDER BY created_at DESC LIMIT 1`)
	var createdAt string
	err := row.Scan(&stats.LastBackupName, &createdAt, &stats.LastBackupSize, &stats.LastStatus)
	if err == nil {
		stats.LastBackupDate = createdAt
	}
	return stats, nil
}

func (r *BackupRepository) CreateRestoreRecord(ctx context.Context, record *domain.RestoreRecord) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO restore_history (id, backup_id, backup_name, restore_date, status, notes, error_message, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		record.ID.String(), record.BackupID.String(), record.BackupName, record.RestoreDate,
		record.Status, record.Notes, record.ErrorMessage, record.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("create restore record", err)
	}
	return nil
}

func (r *BackupRepository) UpdateRestoreRecord(ctx context.Context, record *domain.RestoreRecord) error {
	_, err := r.db.ExecContext(ctx, `UPDATE restore_history SET status = ?, error_message = ? WHERE id = ?`,
		record.Status, record.ErrorMessage, record.ID.String())
	if err != nil {
		return apperror.Database("update restore record", err)
	}
	return nil
}

func (r *BackupRepository) ListRestores(ctx context.Context, limit, offset int) ([]domain.RestoreRecord, int, error) {
	var total int
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM restore_history`).Scan(&total)

	rows, err := r.db.QueryContext(ctx, `SELECT id, backup_id, backup_name, restore_date, status, notes, error_message, created_at
		FROM restore_history ORDER BY created_at DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, 0, apperror.Database("list restores", err)
	}
	defer rows.Close()

	var records []domain.RestoreRecord
	for rows.Next() {
		var rec domain.RestoreRecord
		var idStr, backupIDStr, createdAt string
		err := rows.Scan(&idStr, &backupIDStr, &rec.BackupName, &rec.RestoreDate,
			&rec.Status, &rec.Notes, &rec.ErrorMessage, &createdAt)
		if err != nil {
			return nil, 0, apperror.Database("scan restore", err)
		}
		rec.ID, _ = uuid.Parse(idStr)
		rec.BackupID, _ = uuid.Parse(backupIDStr)
		rec.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		records = append(records, rec)
	}
	return records, total, nil
}

func (r *BackupRepository) scanBackup(row *sql.Row) (*domain.BackupRecord, error) {
	var rec domain.BackupRecord
	var idStr, createdAt string
	err := row.Scan(&idStr, &rec.BackupName, &rec.BackupType, &rec.BackupPath,
		&rec.FileSize, &rec.Checksum, &rec.Status, &rec.ErrorMessage, &createdAt)
	if err != nil {
		return nil, err
	}
	rec.ID, _ = uuid.Parse(idStr)
	rec.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return &rec, nil
}
