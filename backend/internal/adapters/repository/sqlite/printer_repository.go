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

// PrinterRepository is the SQLite implementation of ports.PrinterRepository.
type PrinterRepository struct {
	db  *sql.DB
	log *slog.Logger
}

// NewPrinterRepository creates a new PrinterRepository.
func NewPrinterRepository(db *sql.DB, log *slog.Logger) *PrinterRepository {
	return &PrinterRepository{db: db, log: log}
}

// GetSettings retrieves printer settings (single row).
func (r *PrinterRepository) GetSettings(ctx context.Context) (*domain.PrinterSettings, error) {
	var s domain.PrinterSettings
	var showLogo, showQR int
	var createdAt, updatedAt string

	err := r.db.QueryRowContext(ctx, `
		SELECT id, default_printer, paper_width, margin_top, margin_bottom, margin_left, margin_right,
		       header_text, footer_text, show_logo, show_qr, upi_id, created_at, updated_at
		FROM printer_settings LIMIT 1`).
		Scan(&s.ID, &s.DefaultPrinter, &s.PaperWidth, &s.MarginTop, &s.MarginBottom,
			&s.MarginLeft, &s.MarginRight, &s.HeaderText, &s.FooterText,
			&showLogo, &showQR, &s.UPIID, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("printer_settings", "default")
	}
	if err != nil {
		return nil, apperror.Database("get_printer_settings", err)
	}

	s.ShowLogo = showLogo == 1
	s.ShowQR = showQR == 1
	s.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	s.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return &s, nil
}

// SaveSettings creates or updates printer settings (upsert).
func (r *PrinterRepository) SaveSettings(ctx context.Context, settings *domain.PrinterSettings) error {
	settings.UpdatedAt = time.Now().UTC()
	showLogo := 0
	if settings.ShowLogo {
		showLogo = 1
	}
	showQR := 0
	if settings.ShowQR {
		showQR = 1
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO printer_settings (id, default_printer, paper_width, margin_top, margin_bottom, margin_left, margin_right, header_text, footer_text, show_logo, show_qr, upi_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			default_printer = excluded.default_printer,
			paper_width = excluded.paper_width,
			margin_top = excluded.margin_top,
			margin_bottom = excluded.margin_bottom,
			margin_left = excluded.margin_left,
			margin_right = excluded.margin_right,
			header_text = excluded.header_text,
			footer_text = excluded.footer_text,
			show_logo = excluded.show_logo,
			show_qr = excluded.show_qr,
			upi_id = excluded.upi_id,
			updated_at = excluded.updated_at`,
		settings.ID, settings.DefaultPrinter, settings.PaperWidth,
		settings.MarginTop, settings.MarginBottom, settings.MarginLeft, settings.MarginRight,
		settings.HeaderText, settings.FooterText, showLogo, showQR, settings.UPIID,
		settings.CreatedAt.Format(time.RFC3339), settings.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("save_printer_settings", err)
	}
	return nil
}

// CreatePrintJob inserts a new print job.
func (r *PrinterRepository) CreatePrintJob(ctx context.Context, job *domain.PrintJob) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO print_jobs (id, document_type, document_id, printer_name, paper_width, status, copies, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		job.ID, job.DocumentType, job.DocumentID, job.PrinterName, job.PaperWidth,
		job.Status, job.Copies, job.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("create_print_job", err)
	}
	return nil
}

// UpdatePrintJobStatus updates a print job's status.
func (r *PrinterRepository) UpdatePrintJobStatus(ctx context.Context, id uuid.UUID, status string) error {
	res, err := r.db.ExecContext(ctx, `UPDATE print_jobs SET status=? WHERE id=?`, status, id)
	if err != nil {
		return apperror.Database("update_print_job_status", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperror.NotFound("print_job", id.String())
	}
	return nil
}

// ListPrintJobs lists print jobs with pagination.
func (r *PrinterRepository) ListPrintJobs(ctx context.Context, limit, offset int) ([]domain.PrintJob, int, error) {
	var total int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM print_jobs`).Scan(&total)
	if err != nil {
		return nil, 0, apperror.Database("list_print_jobs_count", err)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, document_type, document_id, printer_name, paper_width, status, copies, created_at
		FROM print_jobs ORDER BY created_at DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, 0, apperror.Database("list_print_jobs", err)
	}
	defer rows.Close()

	var jobs []domain.PrintJob
	for rows.Next() {
		var job domain.PrintJob
		var createdAt string
		if err := rows.Scan(&job.ID, &job.DocumentType, &job.DocumentID, &job.PrinterName, &job.PaperWidth, &job.Status, &job.Copies, &createdAt); err != nil {
			return nil, 0, apperror.Database("list_print_jobs_scan", err)
		}
		job.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		jobs = append(jobs, job)
	}
	return jobs, total, nil
}

// GetPrintJob retrieves a print job by ID.
func (r *PrinterRepository) GetPrintJob(ctx context.Context, id uuid.UUID) (*domain.PrintJob, error) {
	var job domain.PrintJob
	var createdAt string

	err := r.db.QueryRowContext(ctx, `
		SELECT id, document_type, document_id, printer_name, paper_width, status, copies, created_at
		FROM print_jobs WHERE id=?`, id).
		Scan(&job.ID, &job.DocumentType, &job.DocumentID, &job.PrinterName, &job.PaperWidth, &job.Status, &job.Copies, &createdAt)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("print_job", id.String())
	}
	if err != nil {
		return nil, apperror.Database("get_print_job", err)
	}

	job.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return &job, nil
}
