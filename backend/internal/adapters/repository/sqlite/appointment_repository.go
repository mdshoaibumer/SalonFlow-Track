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

// AppointmentRepository is the SQLite implementation.
type AppointmentRepository struct {
	db  *sql.DB
	log *slog.Logger
}

// NewAppointmentRepository creates a new AppointmentRepository.
func NewAppointmentRepository(db *sql.DB, log *slog.Logger) *AppointmentRepository {
	return &AppointmentRepository{db: db, log: log}
}

// Create inserts a new appointment.
func (r *AppointmentRepository) Create(ctx context.Context, appt *domain.Appointment) error {
	isWalkin := 0
	if appt.IsWalkin {
		isWalkin = 1
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO appointments (id, customer_id, staff_id, appointment_date, start_time, end_time, status, notes, is_walkin, total_amount, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		appt.ID, appt.CustomerID, appt.StaffID, appt.AppointmentDate,
		appt.StartTime, appt.EndTime, appt.Status, appt.Notes, isWalkin,
		appt.TotalAmount, appt.CreatedAt.Format(time.RFC3339), appt.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("create_appointment", err)
	}
	return nil
}

// Update updates an appointment.
func (r *AppointmentRepository) Update(ctx context.Context, appt *domain.Appointment) error {
	appt.UpdatedAt = time.Now().UTC()
	isWalkin := 0
	if appt.IsWalkin {
		isWalkin = 1
	}
	res, err := r.db.ExecContext(ctx, `
		UPDATE appointments SET customer_id=?, staff_id=?, appointment_date=?, start_time=?, end_time=?,
		status=?, notes=?, is_walkin=?, total_amount=?, updated_at=? WHERE id=?`,
		appt.CustomerID, appt.StaffID, appt.AppointmentDate, appt.StartTime, appt.EndTime,
		appt.Status, appt.Notes, isWalkin, appt.TotalAmount, appt.UpdatedAt.Format(time.RFC3339), appt.ID)
	if err != nil {
		return apperror.Database("update_appointment", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperror.NotFound("appointment", appt.ID.String())
	}
	return nil
}

// Delete removes an appointment.
func (r *AppointmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM appointments WHERE id=?`, id)
	if err != nil {
		return apperror.Database("delete_appointment", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperror.NotFound("appointment", id.String())
	}
	return nil
}

// GetByID retrieves an appointment by ID.
func (r *AppointmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Appointment, error) {
	var appt domain.Appointment
	var isWalkin int
	var createdAt, updatedAt string

	err := r.db.QueryRowContext(ctx, `
		SELECT id, customer_id, staff_id, appointment_date, start_time, end_time, status, notes, is_walkin, total_amount, created_at, updated_at
		FROM appointments WHERE id=?`, id).
		Scan(&appt.ID, &appt.CustomerID, &appt.StaffID, &appt.AppointmentDate,
			&appt.StartTime, &appt.EndTime, &appt.Status, &appt.Notes, &isWalkin,
			&appt.TotalAmount, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("appointment", id.String())
	}
	if err != nil {
		return nil, apperror.Database("get_appointment", err)
	}

	appt.IsWalkin = isWalkin == 1
	appt.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	appt.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return &appt, nil
}

// List lists appointments with filters.
func (r *AppointmentRepository) List(ctx context.Context, filter domain.AppointmentFilter) ([]domain.Appointment, int, error) {
	where := "1=1"
	var args []interface{}

	if filter.Date != "" {
		where += " AND appointment_date=?"
		args = append(args, filter.Date)
	}
	if filter.StaffID != "" {
		where += " AND staff_id=?"
		args = append(args, filter.StaffID)
	}
	if filter.CustomerID != "" {
		where += " AND customer_id=?"
		args = append(args, filter.CustomerID)
	}
	if filter.Status != "" {
		where += " AND status=?"
		args = append(args, filter.Status)
	}
	if filter.StartDate != "" {
		where += " AND appointment_date>=?"
		args = append(args, filter.StartDate)
	}
	if filter.EndDate != "" {
		where += " AND appointment_date<=?"
		args = append(args, filter.EndDate)
	}

	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM appointments WHERE "+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, apperror.Database("list_appointments_count", err)
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	offset := filter.Offset

	query := "SELECT id, customer_id, staff_id, appointment_date, start_time, end_time, status, notes, is_walkin, total_amount, created_at, updated_at FROM appointments WHERE " + where + " ORDER BY appointment_date, start_time LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, apperror.Database("list_appointments", err)
	}
	defer rows.Close()

	var appts []domain.Appointment
	for rows.Next() {
		var appt domain.Appointment
		var isWalkin int
		var createdAt, updatedAt string
		if err := rows.Scan(&appt.ID, &appt.CustomerID, &appt.StaffID, &appt.AppointmentDate,
			&appt.StartTime, &appt.EndTime, &appt.Status, &appt.Notes, &isWalkin,
			&appt.TotalAmount, &createdAt, &updatedAt); err != nil {
			return nil, 0, apperror.Database("list_appointments_scan", err)
		}
		appt.IsWalkin = isWalkin == 1
		appt.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		appt.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		appts = append(appts, appt)
	}
	return appts, total, nil
}

// AddServices inserts appointment services.
func (r *AppointmentRepository) AddServices(ctx context.Context, services []domain.AppointmentService) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return apperror.Database("add_services_begin", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO appointment_services (id, appointment_id, service_id, service_name, duration_minutes, price, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return apperror.Database("add_services_prepare", err)
	}
	defer stmt.Close()

	for _, s := range services {
		_, err := stmt.ExecContext(ctx, s.ID, s.AppointmentID, s.ServiceID, s.ServiceName, s.DurationMinutes, s.Price, s.CreatedAt.Format(time.RFC3339))
		if err != nil {
			return apperror.Database("add_services_exec", err)
		}
	}
	return tx.Commit()
}

// GetServices retrieves services for an appointment.
func (r *AppointmentRepository) GetServices(ctx context.Context, appointmentID uuid.UUID) ([]domain.AppointmentService, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, appointment_id, service_id, service_name, duration_minutes, price, created_at
		FROM appointment_services WHERE appointment_id=? ORDER BY created_at`, appointmentID)
	if err != nil {
		return nil, apperror.Database("get_services", err)
	}
	defer rows.Close()

	var services []domain.AppointmentService
	for rows.Next() {
		var s domain.AppointmentService
		var createdAt string
		if err := rows.Scan(&s.ID, &s.AppointmentID, &s.ServiceID, &s.ServiceName, &s.DurationMinutes, &s.Price, &createdAt); err != nil {
			return nil, apperror.Database("get_services_scan", err)
		}
		s.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		services = append(services, s)
	}
	return services, nil
}

// DeleteServices removes all services for an appointment.
func (r *AppointmentRepository) DeleteServices(ctx context.Context, appointmentID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM appointment_services WHERE appointment_id=?`, appointmentID)
	if err != nil {
		return apperror.Database("delete_services", err)
	}
	return nil
}

// AddHistory inserts a history entry.
func (r *AppointmentRepository) AddHistory(ctx context.Context, history *domain.AppointmentHistory) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO appointment_history (id, appointment_id, old_status, new_status, changed_by, note, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		history.ID, history.AppointmentID, history.OldStatus, history.NewStatus, history.ChangedBy, history.Note, history.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("add_history", err)
	}
	return nil
}

// GetHistory retrieves history for an appointment.
func (r *AppointmentRepository) GetHistory(ctx context.Context, appointmentID uuid.UUID) ([]domain.AppointmentHistory, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, appointment_id, old_status, new_status, changed_by, note, created_at
		FROM appointment_history WHERE appointment_id=? ORDER BY created_at`, appointmentID)
	if err != nil {
		return nil, apperror.Database("get_history", err)
	}
	defer rows.Close()

	var history []domain.AppointmentHistory
	for rows.Next() {
		var h domain.AppointmentHistory
		var createdAt string
		if err := rows.Scan(&h.ID, &h.AppointmentID, &h.OldStatus, &h.NewStatus, &h.ChangedBy, &h.Note, &createdAt); err != nil {
			return nil, apperror.Database("get_history_scan", err)
		}
		h.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		history = append(history, h)
	}
	return history, nil
}
