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

// GSTRepository is the SQLite implementation of ports.GSTRepository.
type GSTRepository struct {
	db  *sql.DB
	log *slog.Logger
}

// NewGSTRepository creates a new GSTRepository.
func NewGSTRepository(db *sql.DB, log *slog.Logger) *GSTRepository {
	return &GSTRepository{db: db, log: log}
}

// GetSettings retrieves the GST settings (single row).
func (r *GSTRepository) GetSettings(ctx context.Context) (*domain.GSTSettings, error) {
	var s domain.GSTSettings
	var isEnabled int
	var createdAt, updatedAt string

	err := r.db.QueryRowContext(ctx, `SELECT id, business_name, gstin, state, address, hsn_code, cgst_rate, sgst_rate, igst_rate, is_gst_enabled, created_at, updated_at FROM gst_settings LIMIT 1`).
		Scan(&s.ID, &s.BusinessName, &s.GSTIN, &s.State, &s.Address, &s.HSNCode, &s.CGSTRate, &s.SGSTRate, &s.IGSTRate, &isEnabled, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("gst_settings", "default")
	}
	if err != nil {
		return nil, apperror.Database("get_gst_settings", err)
	}

	s.IsGSTEnabled = isEnabled == 1
	s.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	s.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return &s, nil
}

// SaveSettings creates or updates GST settings (upsert).
func (r *GSTRepository) SaveSettings(ctx context.Context, settings *domain.GSTSettings) error {
	settings.UpdatedAt = time.Now().UTC()
	isEnabled := 0
	if settings.IsGSTEnabled {
		isEnabled = 1
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO gst_settings (id, business_name, gstin, state, address, hsn_code, cgst_rate, sgst_rate, igst_rate, is_gst_enabled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			business_name = excluded.business_name,
			gstin = excluded.gstin,
			state = excluded.state,
			address = excluded.address,
			hsn_code = excluded.hsn_code,
			cgst_rate = excluded.cgst_rate,
			sgst_rate = excluded.sgst_rate,
			igst_rate = excluded.igst_rate,
			is_gst_enabled = excluded.is_gst_enabled,
			updated_at = excluded.updated_at`,
		settings.ID, settings.BusinessName, settings.GSTIN, settings.State, settings.Address,
		settings.HSNCode, settings.CGSTRate, settings.SGSTRate, settings.IGSTRate, isEnabled,
		settings.CreatedAt.Format(time.RFC3339), settings.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("save_gst_settings", err)
	}
	return nil
}

// CreateTaxRate inserts a new tax rate.
func (r *GSTRepository) CreateTaxRate(ctx context.Context, rate *domain.TaxRate) error {
	isActive := 0
	if rate.IsActive {
		isActive = 1
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO tax_rates (id, name, hsn_code, cgst_rate, sgst_rate, igst_rate, category, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		rate.ID, rate.Name, rate.HSNCode, rate.CGSTRate, rate.SGSTRate, rate.IGSTRate,
		rate.Category, isActive, rate.CreatedAt.Format(time.RFC3339), rate.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("create_tax_rate", err)
	}
	return nil
}

// UpdateTaxRate updates an existing tax rate.
func (r *GSTRepository) UpdateTaxRate(ctx context.Context, rate *domain.TaxRate) error {
	rate.UpdatedAt = time.Now().UTC()
	isActive := 0
	if rate.IsActive {
		isActive = 1
	}
	res, err := r.db.ExecContext(ctx, `
		UPDATE tax_rates SET name=?, hsn_code=?, cgst_rate=?, sgst_rate=?, igst_rate=?, category=?, is_active=?, updated_at=? WHERE id=?`,
		rate.Name, rate.HSNCode, rate.CGSTRate, rate.SGSTRate, rate.IGSTRate,
		rate.Category, isActive, rate.UpdatedAt.Format(time.RFC3339), rate.ID)
	if err != nil {
		return apperror.Database("update_tax_rate", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperror.NotFound("tax_rate", rate.ID.String())
	}
	return nil
}

// DeleteTaxRate removes a tax rate.
func (r *GSTRepository) DeleteTaxRate(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM tax_rates WHERE id=?`, id)
	if err != nil {
		return apperror.Database("delete_tax_rate", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperror.NotFound("tax_rate", id.String())
	}
	return nil
}

// GetTaxRate retrieves a tax rate by ID.
func (r *GSTRepository) GetTaxRate(ctx context.Context, id uuid.UUID) (*domain.TaxRate, error) {
	var rate domain.TaxRate
	var isActive int
	var createdAt, updatedAt string

	err := r.db.QueryRowContext(ctx, `SELECT id, name, hsn_code, cgst_rate, sgst_rate, igst_rate, category, is_active, created_at, updated_at FROM tax_rates WHERE id=?`, id).
		Scan(&rate.ID, &rate.Name, &rate.HSNCode, &rate.CGSTRate, &rate.SGSTRate, &rate.IGSTRate, &rate.Category, &isActive, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("tax_rate", id.String())
	}
	if err != nil {
		return nil, apperror.Database("get_tax_rate", err)
	}

	rate.IsActive = isActive == 1
	rate.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	rate.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return &rate, nil
}

// ListTaxRates lists all tax rates, optionally filtered by category.
func (r *GSTRepository) ListTaxRates(ctx context.Context, category string) ([]domain.TaxRate, error) {
	query := `SELECT id, name, hsn_code, cgst_rate, sgst_rate, igst_rate, category, is_active, created_at, updated_at FROM tax_rates`
	var args []interface{}
	if category != "" {
		query += ` WHERE category=?`
		args = append(args, category)
	}
	query += ` ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.Database("list_tax_rates", err)
	}
	defer rows.Close()

	var rates []domain.TaxRate
	for rows.Next() {
		var rate domain.TaxRate
		var isActive int
		var createdAt, updatedAt string
		if err := rows.Scan(&rate.ID, &rate.Name, &rate.HSNCode, &rate.CGSTRate, &rate.SGSTRate, &rate.IGSTRate, &rate.Category, &isActive, &createdAt, &updatedAt); err != nil {
			return nil, apperror.Database("list_tax_rates_scan", err)
		}
		rate.IsActive = isActive == 1
		rate.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		rate.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		rates = append(rates, rate)
	}
	return rates, nil
}

// GetTaxRateByCategory returns the first active tax rate for a category.
func (r *GSTRepository) GetTaxRateByCategory(ctx context.Context, category string) (*domain.TaxRate, error) {
	var rate domain.TaxRate
	var isActive int
	var createdAt, updatedAt string

	err := r.db.QueryRowContext(ctx, `SELECT id, name, hsn_code, cgst_rate, sgst_rate, igst_rate, category, is_active, created_at, updated_at FROM tax_rates WHERE category=? AND is_active=1 LIMIT 1`, category).
		Scan(&rate.ID, &rate.Name, &rate.HSNCode, &rate.CGSTRate, &rate.SGSTRate, &rate.IGSTRate, &rate.Category, &isActive, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("tax_rate", category)
	}
	if err != nil {
		return nil, apperror.Database("get_tax_rate_by_category", err)
	}

	rate.IsActive = isActive == 1
	rate.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	rate.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return &rate, nil
}

// CreateTaxLines inserts invoice tax lines in a batch.
func (r *GSTRepository) CreateTaxLines(ctx context.Context, lines []domain.InvoiceTaxLine) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return apperror.Database("create_tax_lines_begin", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO invoice_tax_lines (id, invoice_id, item_id, taxable_amount, cgst_rate, cgst_amount, sgst_rate, sgst_amount, igst_rate, igst_amount, total_tax, is_interstate, hsn_code, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return apperror.Database("create_tax_lines_prepare", err)
	}
	defer stmt.Close()

	for _, line := range lines {
		isInterstate := 0
		if line.IsInterstate {
			isInterstate = 1
		}
		_, err := stmt.ExecContext(ctx,
			line.ID, line.InvoiceID, line.ItemID, line.TaxableAmount,
			line.CGSTRate, line.CGSTAmount, line.SGSTRate, line.SGSTAmount,
			line.IGSTRate, line.IGSTAmount, line.TotalTax, isInterstate,
			line.HSNCode, line.CreatedAt.Format(time.RFC3339))
		if err != nil {
			return apperror.Database("create_tax_lines_exec", err)
		}
	}

	return tx.Commit()
}

// GetTaxLinesByInvoice retrieves all tax lines for an invoice.
func (r *GSTRepository) GetTaxLinesByInvoice(ctx context.Context, invoiceID uuid.UUID) ([]domain.InvoiceTaxLine, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, invoice_id, item_id, taxable_amount, cgst_rate, cgst_amount, sgst_rate, sgst_amount, igst_rate, igst_amount, total_tax, is_interstate, hsn_code, created_at
		FROM invoice_tax_lines WHERE invoice_id=? ORDER BY created_at`, invoiceID)
	if err != nil {
		return nil, apperror.Database("get_tax_lines", err)
	}
	defer rows.Close()

	var lines []domain.InvoiceTaxLine
	for rows.Next() {
		var line domain.InvoiceTaxLine
		var isInterstate int
		var createdAt string
		if err := rows.Scan(&line.ID, &line.InvoiceID, &line.ItemID, &line.TaxableAmount,
			&line.CGSTRate, &line.CGSTAmount, &line.SGSTRate, &line.SGSTAmount,
			&line.IGSTRate, &line.IGSTAmount, &line.TotalTax, &isInterstate,
			&line.HSNCode, &createdAt); err != nil {
			return nil, apperror.Database("get_tax_lines_scan", err)
		}
		line.IsInterstate = isInterstate == 1
		line.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		lines = append(lines, line)
	}
	return lines, nil
}

// GetGSTReport generates a GST report for the given filter period.
func (r *GSTRepository) GetGSTReport(ctx context.Context, filter domain.GSTReportFilter) (*domain.GSTReport, error) {
	report := &domain.GSTReport{
		Period:    filter.Period,
		StartDate: filter.StartDate,
		EndDate:   filter.EndDate,
	}

	err := r.db.QueryRowContext(ctx, `
		SELECT COALESCE(COUNT(DISTINCT invoice_id), 0),
			   COALESCE(SUM(taxable_amount), 0),
			   COALESCE(SUM(cgst_amount), 0),
			   COALESCE(SUM(sgst_amount), 0),
			   COALESCE(SUM(igst_amount), 0),
			   COALESCE(SUM(total_tax), 0)
		FROM invoice_tax_lines
		WHERE created_at >= ? AND created_at <= ?`,
		filter.StartDate, filter.EndDate+"T23:59:59Z").
		Scan(&report.TotalInvoices, &report.TaxableAmount, &report.TotalCGST, &report.TotalSGST, &report.TotalIGST, &report.TotalTax)
	if err != nil {
		return nil, apperror.Database("get_gst_report", err)
	}

	report.GrandTotal = report.TaxableAmount + report.TotalTax
	return report, nil
}
