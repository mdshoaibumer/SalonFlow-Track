package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// InvoiceRepository is the SQLite implementation of ports.InvoiceRepository.
type InvoiceRepository struct {
	db  *sql.DB
	log *slog.Logger
}

// NewInvoiceRepository creates a new InvoiceRepository.
func NewInvoiceRepository(db *sql.DB, log *slog.Logger) *InvoiceRepository {
	return &InvoiceRepository{db: db, log: log}
}

// Create inserts a new invoice with its items in a transaction.
func (r *InvoiceRepository) Create(ctx context.Context, invoice *domain.Invoice) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return apperror.Database("begin invoice tx", err)
	}
	defer tx.Rollback()

	invoiceQuery := `
		INSERT INTO invoices (id, invoice_number, customer_id, staff_id, subtotal, discount, tax, grand_total, payment_status, payment_method, notes, invoice_date, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = tx.ExecContext(ctx, invoiceQuery,
		invoice.ID.String(),
		invoice.InvoiceNumber,
		invoice.CustomerID.String(),
		invoice.StaffID.String(),
		invoice.Subtotal,
		invoice.Discount,
		invoice.Tax,
		invoice.GrandTotal,
		invoice.PaymentStatus,
		invoice.PaymentMethod,
		invoice.Notes,
		invoice.InvoiceDate.Format(time.RFC3339),
		invoice.CreatedAt.Format(time.RFC3339),
		invoice.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return apperror.Conflict("invoice number already exists")
		}
		return apperror.Database("create invoice", err)
	}

	// Insert items
	itemQuery := `
		INSERT INTO invoice_items (id, invoice_id, service_id, service_name_snapshot, quantity, unit_price, discount, line_total, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	for _, item := range invoice.Items {
		_, err = tx.ExecContext(ctx, itemQuery,
			item.ID.String(),
			invoice.ID.String(),
			item.ServiceID.String(),
			item.ServiceNameSnapshot,
			item.Quantity,
			item.UnitPrice,
			item.Discount,
			item.LineTotal,
			item.CreatedAt.Format(time.RFC3339),
			item.UpdatedAt.Format(time.RFC3339),
		)
		if err != nil {
			return apperror.Database("create invoice item", err)
		}
	}

	return tx.Commit()
}

// GetByID retrieves an invoice by ID with its items.
func (r *InvoiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Invoice, error) {
	query := `
		SELECT id, invoice_number, customer_id, staff_id, subtotal, discount, tax, grand_total, payment_status, payment_method, notes, invoice_date, created_at, updated_at
		FROM invoices WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id.String())
	inv, err := r.scanInvoice(row)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("invoice", id.String())
	}
	if err != nil {
		return nil, apperror.Database("get invoice by id", err)
	}

	// Fetch items
	itemsQuery := `
		SELECT id, invoice_id, service_id, service_name_snapshot, quantity, unit_price, discount, line_total, created_at, updated_at
		FROM invoice_items WHERE invoice_id = ?`

	rows, err := r.db.QueryContext(ctx, itemsQuery, id.String())
	if err != nil {
		return nil, apperror.Database("get invoice items", err)
	}
	defer rows.Close()

	for rows.Next() {
		item, err := r.scanInvoiceItem(rows)
		if err != nil {
			return nil, apperror.Database("scan invoice item", err)
		}
		inv.Items = append(inv.Items, *item)
	}

	return inv, nil
}

// List returns invoices matching the filter.
func (r *InvoiceRepository) List(ctx context.Context, filter ports.InvoiceFilter) ([]domain.Invoice, int, error) {
	var conditions []string
	var args []interface{}

	if filter.CustomerID != "" {
		conditions = append(conditions, "customer_id = ?")
		args = append(args, filter.CustomerID)
	}
	if filter.StaffID != "" {
		conditions = append(conditions, "staff_id = ?")
		args = append(args, filter.StaffID)
	}
	if filter.PaymentStatus != "" {
		conditions = append(conditions, "payment_status = ?")
		args = append(args, filter.PaymentStatus)
	}
	if filter.DateFrom != "" {
		conditions = append(conditions, "invoice_date >= ?")
		args = append(args, filter.DateFrom)
	}
	if filter.DateTo != "" {
		conditions = append(conditions, "invoice_date <= ?")
		args = append(args, filter.DateTo)
	}
	if filter.Search != "" {
		conditions = append(conditions, "(invoice_number LIKE ?)")
		search := "%" + filter.Search + "%"
		args = append(args, search)
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM invoices %s", where)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, apperror.Database("count invoices", err)
	}

	// Fetch rows
	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	dataQuery := fmt.Sprintf(`
		SELECT id, invoice_number, customer_id, staff_id, subtotal, discount, tax, grand_total, payment_status, payment_method, notes, invoice_date, created_at, updated_at
		FROM invoices %s
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`, where)

	dataArgs := append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, apperror.Database("list invoices", err)
	}
	defer rows.Close()

	var results []domain.Invoice
	for rows.Next() {
		inv, err := r.scanInvoiceRows(rows)
		if err != nil {
			return nil, 0, apperror.Database("scan invoice row", err)
		}
		results = append(results, *inv)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, apperror.Database("iterate invoice rows", err)
	}

	return results, total, nil
}

// GetNextSequence returns the next invoice sequence number for the year.
func (r *InvoiceRepository) GetNextSequence(ctx context.Context, year int) (int, error) {
	prefix := fmt.Sprintf("INV-%d-", year)
	query := `SELECT COUNT(*) FROM invoices WHERE invoice_number LIKE ?`
	var count int
	if err := r.db.QueryRowContext(ctx, query, prefix+"%").Scan(&count); err != nil {
		return 0, apperror.Database("get next sequence", err)
	}
	return count + 1, nil
}

// RecordPayment inserts a payment record.
func (r *InvoiceRepository) RecordPayment(ctx context.Context, payment *domain.Payment) error {
	query := `
		INSERT INTO payments (id, invoice_id, amount, payment_method, reference_number, payment_date, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		payment.ID.String(),
		payment.InvoiceID.String(),
		payment.Amount,
		payment.PaymentMethod,
		payment.ReferenceNumber,
		payment.PaymentDate.Format(time.RFC3339),
		payment.CreatedAt.Format(time.RFC3339),
		payment.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return apperror.Database("record payment", err)
	}
	return nil
}

// GetPayments retrieves all payments for an invoice.
func (r *InvoiceRepository) GetPayments(ctx context.Context, invoiceID uuid.UUID) ([]domain.Payment, error) {
	query := `
		SELECT id, invoice_id, amount, payment_method, reference_number, payment_date, created_at, updated_at
		FROM payments WHERE invoice_id = ? ORDER BY payment_date ASC`

	rows, err := r.db.QueryContext(ctx, query, invoiceID.String())
	if err != nil {
		return nil, apperror.Database("get payments", err)
	}
	defer rows.Close()

	var results []domain.Payment
	for rows.Next() {
		var p domain.Payment
		var id, invoiceID, paymentDate, createdAt, updatedAt string
		err := rows.Scan(&id, &invoiceID, &p.Amount, &p.PaymentMethod, &p.ReferenceNumber, &paymentDate, &createdAt, &updatedAt)
		if err != nil {
			return nil, apperror.Database("scan payment", err)
		}
		p.ID, _ = uuid.Parse(id)
		p.InvoiceID, _ = uuid.Parse(invoiceID)
		if t, err := time.Parse(time.RFC3339, paymentDate); err == nil {
			p.PaymentDate = t
		}
		if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
			p.CreatedAt = t
		}
		if t, err := time.Parse(time.RFC3339, updatedAt); err == nil {
			p.UpdatedAt = t
		}
		results = append(results, p)
	}

	return results, nil
}

// UpdateStatus updates the payment status of an invoice.
func (r *InvoiceRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status, method string) error {
	query := `UPDATE invoices SET payment_status = ?, payment_method = ?, updated_at = ? WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, status, method, time.Now().UTC().Format(time.RFC3339), id.String())
	if err != nil {
		return apperror.Database("update invoice status", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return apperror.NotFound("invoice", id.String())
	}
	return nil
}

// GetTodayStats returns today's revenue stats.
func (r *InvoiceRepository) GetTodayStats(ctx context.Context) (todayRevenue float64, todayCount int, avgBill float64, err error) {
	today := time.Now().UTC().Format("2006-01-02")
	query := `SELECT COALESCE(SUM(grand_total), 0), COUNT(*), COALESCE(AVG(grand_total), 0) FROM invoices WHERE date(invoice_date) = ?`
	err = r.db.QueryRowContext(ctx, query, today).Scan(&todayRevenue, &todayCount, &avgBill)
	if err != nil {
		return 0, 0, 0, apperror.Database("get today stats", err)
	}
	return todayRevenue, todayCount, avgBill, nil
}

// --- scan helpers ---

func (r *InvoiceRepository) scanInvoice(row *sql.Row) (*domain.Invoice, error) {
	var inv domain.Invoice
	var id, customerID, staffID string
	var invoiceDate, createdAt, updatedAt string

	err := row.Scan(&id, &inv.InvoiceNumber, &customerID, &staffID, &inv.Subtotal, &inv.Discount, &inv.Tax, &inv.GrandTotal, &inv.PaymentStatus, &inv.PaymentMethod, &inv.Notes, &invoiceDate, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	inv.ID, _ = uuid.Parse(id)
	inv.CustomerID, _ = uuid.Parse(customerID)
	inv.StaffID, _ = uuid.Parse(staffID)
	if t, err := time.Parse(time.RFC3339, invoiceDate); err == nil {
		inv.InvoiceDate = t
	}
	if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
		inv.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, updatedAt); err == nil {
		inv.UpdatedAt = t
	}

	return &inv, nil
}

func (r *InvoiceRepository) scanInvoiceRows(rows *sql.Rows) (*domain.Invoice, error) {
	var inv domain.Invoice
	var id, customerID, staffID string
	var invoiceDate, createdAt, updatedAt string

	err := rows.Scan(&id, &inv.InvoiceNumber, &customerID, &staffID, &inv.Subtotal, &inv.Discount, &inv.Tax, &inv.GrandTotal, &inv.PaymentStatus, &inv.PaymentMethod, &inv.Notes, &invoiceDate, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	inv.ID, _ = uuid.Parse(id)
	inv.CustomerID, _ = uuid.Parse(customerID)
	inv.StaffID, _ = uuid.Parse(staffID)
	if t, err := time.Parse(time.RFC3339, invoiceDate); err == nil {
		inv.InvoiceDate = t
	}
	if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
		inv.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, updatedAt); err == nil {
		inv.UpdatedAt = t
	}

	return &inv, nil
}

func (r *InvoiceRepository) scanInvoiceItem(rows *sql.Rows) (*domain.InvoiceItem, error) {
	var item domain.InvoiceItem
	var id, invoiceID, serviceID string
	var createdAt, updatedAt string

	err := rows.Scan(&id, &invoiceID, &serviceID, &item.ServiceNameSnapshot, &item.Quantity, &item.UnitPrice, &item.Discount, &item.LineTotal, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	item.ID, _ = uuid.Parse(id)
	item.InvoiceID, _ = uuid.Parse(invoiceID)
	item.ServiceID, _ = uuid.Parse(serviceID)
	if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
		item.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, updatedAt); err == nil {
		item.UpdatedAt = t
	}

	return &item, nil
}
