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

// LicenseRepository implements ports.LicenseRepository.
type LicenseRepository struct {
	db  *sql.DB
	log *slog.Logger
}

// NewLicenseRepository creates a new LicenseRepository.
func NewLicenseRepository(db *sql.DB, log *slog.Logger) *LicenseRepository {
	return &LicenseRepository{db: db, log: log}
}

func (r *LicenseRepository) CreateLicense(ctx context.Context, lic *domain.License) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO licenses (id, license_key, customer_name, salon_name, device_id, issued_date, expiry_date, status, signature, last_validation, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		lic.ID.String(), lic.LicenseKey, lic.CustomerName, lic.SalonName, lic.DeviceID,
		lic.IssuedDate, lic.ExpiryDate, lic.Status, lic.Signature, lic.LastValidation,
		lic.CreatedAt.Format(time.RFC3339), lic.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("create license", err)
	}
	return nil
}

func (r *LicenseRepository) UpdateLicense(ctx context.Context, lic *domain.License) error {
	lic.UpdatedAt = time.Now().UTC()
	_, err := r.db.ExecContext(ctx, `UPDATE licenses SET expiry_date = ?, status = ?, signature = ?, last_validation = ?, updated_at = ? WHERE id = ?`,
		lic.ExpiryDate, lic.Status, lic.Signature, lic.LastValidation, lic.UpdatedAt.Format(time.RFC3339), lic.ID.String())
	if err != nil {
		return apperror.Database("update license", err)
	}
	return nil
}

func (r *LicenseRepository) GetActiveLicense(ctx context.Context) (*domain.License, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, license_key, customer_name, salon_name, device_id, issued_date, expiry_date, status, signature, last_validation, created_at, updated_at
		FROM licenses ORDER BY created_at DESC LIMIT 1`)
	lic, err := r.scanLicense(row)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("license", "active")
	}
	if err != nil {
		return nil, apperror.Database("get active license", err)
	}
	return lic, nil
}

func (r *LicenseRepository) GetLicenseByKey(ctx context.Context, key string) (*domain.License, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, license_key, customer_name, salon_name, device_id, issued_date, expiry_date, status, signature, last_validation, created_at, updated_at
		FROM licenses WHERE license_key = ?`, key)
	lic, err := r.scanLicense(row)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("license", key)
	}
	if err != nil {
		return nil, apperror.Database("get license by key", err)
	}
	return lic, nil
}

func (r *LicenseRepository) CreateEvent(ctx context.Context, event *domain.LicenseEvent) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO license_events (id, license_id, event_type, event_date, notes, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		event.ID.String(), event.LicenseID.String(), event.EventType, event.EventDate,
		event.Notes, event.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("create license event", err)
	}
	return nil
}

func (r *LicenseRepository) ListEvents(ctx context.Context, licenseID uuid.UUID, limit, offset int) ([]domain.LicenseEvent, int, error) {
	var total int
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM license_events WHERE license_id = ?`, licenseID.String()).Scan(&total)

	rows, err := r.db.QueryContext(ctx, `SELECT id, license_id, event_type, event_date, notes, created_at
		FROM license_events WHERE license_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`, licenseID.String(), limit, offset)
	if err != nil {
		return nil, 0, apperror.Database("list license events", err)
	}
	defer rows.Close()

	var events []domain.LicenseEvent
	for rows.Next() {
		var ev domain.LicenseEvent
		var idStr, licIDStr, createdAt string
		err := rows.Scan(&idStr, &licIDStr, &ev.EventType, &ev.EventDate, &ev.Notes, &createdAt)
		if err != nil {
			return nil, 0, apperror.Database("scan license event", err)
		}
		ev.ID, _ = uuid.Parse(idStr)
		ev.LicenseID, _ = uuid.Parse(licIDStr)
		ev.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		events = append(events, ev)
	}
	return events, total, nil
}

func (r *LicenseRepository) scanLicense(row *sql.Row) (*domain.License, error) {
	var lic domain.License
	var idStr, createdAt, updatedAt string
	err := row.Scan(&idStr, &lic.LicenseKey, &lic.CustomerName, &lic.SalonName, &lic.DeviceID,
		&lic.IssuedDate, &lic.ExpiryDate, &lic.Status, &lic.Signature, &lic.LastValidation,
		&createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	lic.ID, _ = uuid.Parse(idStr)
	lic.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	lic.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return &lic, nil
}
