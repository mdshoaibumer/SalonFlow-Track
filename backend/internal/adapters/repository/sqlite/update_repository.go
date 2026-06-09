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

// UpdateRepository implements ports.UpdateRepository.
type UpdateRepository struct {
	db  *sql.DB
	log *slog.Logger
}

// NewUpdateRepository creates a new UpdateRepository.
func NewUpdateRepository(db *sql.DB, log *slog.Logger) *UpdateRepository {
	return &UpdateRepository{db: db, log: log}
}

func (r *UpdateRepository) CreateVersion(ctx context.Context, version *domain.AppVersion) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO app_versions (id, version, release_date, release_notes, installed_at, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		version.ID.String(), version.Version, version.ReleaseDate, version.ReleaseNotes,
		version.InstalledAt, version.Status, version.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("create version", err)
	}
	return nil
}

func (r *UpdateRepository) UpdateVersion(ctx context.Context, version *domain.AppVersion) error {
	_, err := r.db.ExecContext(ctx, `UPDATE app_versions SET installed_at = ?, status = ? WHERE id = ?`,
		version.InstalledAt, version.Status, version.ID.String())
	if err != nil {
		return apperror.Database("update version", err)
	}
	return nil
}

func (r *UpdateRepository) GetVersionByName(ctx context.Context, version string) (*domain.AppVersion, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, version, release_date, release_notes, installed_at, status, created_at
		FROM app_versions WHERE version = ?`, version)
	v, err := r.scanVersion(row)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("version", version)
	}
	if err != nil {
		return nil, apperror.Database("get version", err)
	}
	return v, nil
}

func (r *UpdateRepository) GetInstalledVersion(ctx context.Context) (*domain.AppVersion, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, version, release_date, release_notes, installed_at, status, created_at
		FROM app_versions WHERE status = 'installed' ORDER BY installed_at DESC LIMIT 1`)
	v, err := r.scanVersion(row)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("version", "installed")
	}
	if err != nil {
		return nil, apperror.Database("get installed version", err)
	}
	return v, nil
}

func (r *UpdateRepository) ListVersions(ctx context.Context, limit, offset int) ([]domain.AppVersion, int, error) {
	var total int
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM app_versions`).Scan(&total)

	rows, err := r.db.QueryContext(ctx, `SELECT id, version, release_date, release_notes, installed_at, status, created_at
		FROM app_versions ORDER BY created_at DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, 0, apperror.Database("list versions", err)
	}
	defer rows.Close()

	var versions []domain.AppVersion
	for rows.Next() {
		var v domain.AppVersion
		var idStr, createdAt string
		err := rows.Scan(&idStr, &v.Version, &v.ReleaseDate, &v.ReleaseNotes, &v.InstalledAt, &v.Status, &createdAt)
		if err != nil {
			return nil, 0, apperror.Database("scan version", err)
		}
		v.ID, _ = uuid.Parse(idStr)
		v.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		versions = append(versions, v)
	}
	return versions, total, nil
}

func (r *UpdateRepository) CreateUpdateRecord(ctx context.Context, record *domain.UpdateRecord) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO update_history (id, from_version, to_version, update_date, status, error_message, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		record.ID.String(), record.FromVersion, record.ToVersion, record.UpdateDate,
		record.Status, record.ErrorMessage, record.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("create update record", err)
	}
	return nil
}

func (r *UpdateRepository) UpdateUpdateRecord(ctx context.Context, record *domain.UpdateRecord) error {
	_, err := r.db.ExecContext(ctx, `UPDATE update_history SET status = ?, error_message = ? WHERE id = ?`,
		record.Status, record.ErrorMessage, record.ID.String())
	if err != nil {
		return apperror.Database("update update record", err)
	}
	return nil
}

func (r *UpdateRepository) ListUpdateHistory(ctx context.Context, limit, offset int) ([]domain.UpdateRecord, int, error) {
	var total int
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM update_history`).Scan(&total)

	rows, err := r.db.QueryContext(ctx, `SELECT id, from_version, to_version, update_date, status, error_message, created_at
		FROM update_history ORDER BY created_at DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, 0, apperror.Database("list update history", err)
	}
	defer rows.Close()

	var records []domain.UpdateRecord
	for rows.Next() {
		var rec domain.UpdateRecord
		var idStr, createdAt string
		err := rows.Scan(&idStr, &rec.FromVersion, &rec.ToVersion, &rec.UpdateDate, &rec.Status, &rec.ErrorMessage, &createdAt)
		if err != nil {
			return nil, 0, apperror.Database("scan update record", err)
		}
		rec.ID, _ = uuid.Parse(idStr)
		rec.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		records = append(records, rec)
	}
	return records, total, nil
}

func (r *UpdateRepository) scanVersion(row *sql.Row) (*domain.AppVersion, error) {
	var v domain.AppVersion
	var idStr, createdAt string
	err := row.Scan(&idStr, &v.Version, &v.ReleaseDate, &v.ReleaseNotes, &v.InstalledAt, &v.Status, &createdAt)
	if err != nil {
		return nil, err
	}
	v.ID, _ = uuid.Parse(idStr)
	v.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return &v, nil
}
