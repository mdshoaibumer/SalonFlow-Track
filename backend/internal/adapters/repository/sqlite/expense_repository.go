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

// ExpenseRepository implements ports.ExpenseRepository.
type ExpenseRepository struct {
	db  *sql.DB
	log *slog.Logger
}

// NewExpenseRepository creates a new ExpenseRepository.
func NewExpenseRepository(db *sql.DB, log *slog.Logger) *ExpenseRepository {
	return &ExpenseRepository{db: db, log: log}
}

// --- Categories ---

func (r *ExpenseRepository) CreateCategory(ctx context.Context, cat *domain.ExpenseCategory) error {
	query := `INSERT INTO expense_categories (id, name, description, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		cat.ID.String(), cat.Name, cat.Description,
		boolToInt(cat.IsActive), cat.CreatedAt.Format(time.RFC3339), cat.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("create expense category", err)
	}
	return nil
}

func (r *ExpenseRepository) GetCategoryByID(ctx context.Context, id uuid.UUID) (*domain.ExpenseCategory, error) {
	query := `SELECT id, name, description, is_active, created_at, updated_at FROM expense_categories WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id.String())
	cat, err := r.scanCategory(row)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("expense_category", id.String())
	}
	if err != nil {
		return nil, apperror.Database("get expense category", err)
	}
	return cat, nil
}

func (r *ExpenseRepository) ListCategories(ctx context.Context, activeOnly bool) ([]domain.ExpenseCategory, error) {
	query := `SELECT id, name, description, is_active, created_at, updated_at FROM expense_categories`
	if activeOnly {
		query += ` WHERE is_active = 1`
	}
	query += ` ORDER BY name ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, apperror.Database("list expense categories", err)
	}
	defer rows.Close()

	var categories []domain.ExpenseCategory
	for rows.Next() {
		cat, err := r.scanCategoryRow(rows)
		if err != nil {
			return nil, apperror.Database("scan expense category", err)
		}
		categories = append(categories, *cat)
	}
	return categories, nil
}

func (r *ExpenseRepository) UpdateCategory(ctx context.Context, cat *domain.ExpenseCategory) error {
	query := `UPDATE expense_categories SET name = ?, description = ?, is_active = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		cat.Name, cat.Description, boolToInt(cat.IsActive),
		time.Now().UTC().Format(time.RFC3339), cat.ID.String())
	if err != nil {
		return apperror.Database("update expense category", err)
	}
	return nil
}

// --- Expenses ---

func (r *ExpenseRepository) CreateExpense(ctx context.Context, exp *domain.Expense) error {
	query := `INSERT INTO expenses (id, expense_number, category_id, amount, expense_date, payment_method,
		vendor_name, invoice_reference, description, attachment_path, status, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		exp.ID.String(), exp.ExpenseNumber, exp.CategoryID.String(), exp.Amount,
		exp.ExpenseDate, exp.PaymentMethod, exp.VendorName, exp.InvoiceReference,
		exp.Description, exp.AttachmentPath, exp.Status, exp.CreatedBy,
		exp.CreatedAt.Format(time.RFC3339), exp.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("create expense", err)
	}
	return nil
}

func (r *ExpenseRepository) GetExpenseByID(ctx context.Context, id uuid.UUID) (*domain.Expense, error) {
	query := `SELECT e.id, e.expense_number, e.category_id, ec.name, e.amount, e.expense_date,
		e.payment_method, e.vendor_name, e.invoice_reference, e.description, e.attachment_path,
		e.status, e.created_by, e.created_at, e.updated_at
		FROM expenses e
		JOIN expense_categories ec ON ec.id = e.category_id
		WHERE e.id = ?`
	row := r.db.QueryRowContext(ctx, query, id.String())
	exp, err := r.scanExpense(row)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("expense", id.String())
	}
	if err != nil {
		return nil, apperror.Database("get expense", err)
	}
	return exp, nil
}

func (r *ExpenseRepository) UpdateExpense(ctx context.Context, exp *domain.Expense) error {
	query := `UPDATE expenses SET category_id = ?, amount = ?, expense_date = ?, payment_method = ?,
		vendor_name = ?, invoice_reference = ?, description = ?, attachment_path = ?, status = ?,
		updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		exp.CategoryID.String(), exp.Amount, exp.ExpenseDate, exp.PaymentMethod,
		exp.VendorName, exp.InvoiceReference, exp.Description, exp.AttachmentPath,
		exp.Status, time.Now().UTC().Format(time.RFC3339), exp.ID.String())
	if err != nil {
		return apperror.Database("update expense", err)
	}
	return nil
}

func (r *ExpenseRepository) DeleteExpense(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM expenses WHERE id = ?`, id.String())
	if err != nil {
		return apperror.Database("delete expense", err)
	}
	return nil
}

func (r *ExpenseRepository) ListExpenses(ctx context.Context, filter ports.ExpenseFilter) ([]domain.Expense, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}

	if filter.CategoryID != "" {
		where = append(where, "e.category_id = ?")
		args = append(args, filter.CategoryID)
	}
	if filter.Status != "" {
		where = append(where, "e.status = ?")
		args = append(args, filter.Status)
	}
	if filter.PaymentMethod != "" {
		where = append(where, "e.payment_method = ?")
		args = append(args, filter.PaymentMethod)
	}
	if filter.DateFrom != "" {
		where = append(where, "e.expense_date >= ?")
		args = append(args, filter.DateFrom)
	}
	if filter.DateTo != "" {
		where = append(where, "e.expense_date <= ?")
		args = append(args, filter.DateTo)
	}
	if filter.Search != "" {
		where = append(where, "(e.description LIKE ? OR e.vendor_name LIKE ? OR e.expense_number LIKE ?)")
		s := "%" + filter.Search + "%"
		args = append(args, s, s, s)
	}

	whereClause := strings.Join(where, " AND ")

	// Count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM expenses e WHERE %s", whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, apperror.Database("count expenses", err)
	}

	// Data
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	query := fmt.Sprintf(`SELECT e.id, e.expense_number, e.category_id, ec.name, e.amount, e.expense_date,
		e.payment_method, e.vendor_name, e.invoice_reference, e.description, e.attachment_path,
		e.status, e.created_by, e.created_at, e.updated_at
		FROM expenses e
		JOIN expense_categories ec ON ec.id = e.category_id
		WHERE %s ORDER BY e.expense_date DESC, e.created_at DESC LIMIT ? OFFSET ?`, whereClause)
	args = append(args, limit, filter.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, apperror.Database("list expenses", err)
	}
	defer rows.Close()

	var expenses []domain.Expense
	for rows.Next() {
		exp, err := r.scanExpenseRow(rows)
		if err != nil {
			return nil, 0, apperror.Database("scan expense", err)
		}
		expenses = append(expenses, *exp)
	}
	return expenses, total, nil
}

func (r *ExpenseRepository) NextExpenseNumber(ctx context.Context, year int) (string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", apperror.Database("begin tx for expense number", err)
	}
	defer tx.Rollback()

	// Upsert the sequence
	_, err = tx.ExecContext(ctx, `INSERT INTO expense_number_seq (year, seq) VALUES (?, 0) ON CONFLICT(year) DO NOTHING`, year)
	if err != nil {
		return "", apperror.Database("init expense seq", err)
	}

	_, err = tx.ExecContext(ctx, `UPDATE expense_number_seq SET seq = seq + 1 WHERE year = ?`, year)
	if err != nil {
		return "", apperror.Database("increment expense seq", err)
	}

	var seq int
	err = tx.QueryRowContext(ctx, `SELECT seq FROM expense_number_seq WHERE year = ?`, year).Scan(&seq)
	if err != nil {
		return "", apperror.Database("read expense seq", err)
	}

	if err := tx.Commit(); err != nil {
		return "", apperror.Database("commit expense number", err)
	}

	return fmt.Sprintf("EXP-%d-%06d", year, seq), nil
}

// --- Reporting ---

func (r *ExpenseRepository) GetTotalExpensesByDateRange(ctx context.Context, dateFrom, dateTo string) (float64, error) {
	query := `SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE expense_date BETWEEN ? AND ? AND status != 'rejected'`
	var total float64
	err := r.db.QueryRowContext(ctx, query, dateFrom, dateTo).Scan(&total)
	if err != nil {
		return 0, apperror.Database("get total expenses", err)
	}
	return total, nil
}

func (r *ExpenseRepository) GetExpensesByCategory(ctx context.Context, dateFrom, dateTo string) ([]domain.CategoryExpense, error) {
	query := `SELECT e.category_id, ec.name, COALESCE(SUM(e.amount), 0) as total
		FROM expenses e
		JOIN expense_categories ec ON ec.id = e.category_id
		WHERE e.expense_date BETWEEN ? AND ? AND e.status != 'rejected'
		GROUP BY e.category_id, ec.name
		ORDER BY total DESC`

	rows, err := r.db.QueryContext(ctx, query, dateFrom, dateTo)
	if err != nil {
		return nil, apperror.Database("get expenses by category", err)
	}
	defer rows.Close()

	var results []domain.CategoryExpense
	var grandTotal float64
	for rows.Next() {
		var ce domain.CategoryExpense
		if err := rows.Scan(&ce.CategoryID, &ce.CategoryName, &ce.Amount); err != nil {
			return nil, apperror.Database("scan category expense", err)
		}
		grandTotal += ce.Amount
		results = append(results, ce)
	}

	// Calculate percentages
	for i := range results {
		if grandTotal > 0 {
			results[i].Percentage = (results[i].Amount / grandTotal) * 100
		}
	}

	return results, nil
}

func (r *ExpenseRepository) GetMonthlyExpenseTrend(ctx context.Context, months int) ([]domain.MonthlyTrend, error) {
	query := `SELECT strftime('%Y-%m', expense_date) as month, COALESCE(SUM(amount), 0) as total
		FROM expenses
		WHERE expense_date >= date('now', ? || ' months') AND status != 'rejected'
		GROUP BY month
		ORDER BY month ASC`

	rows, err := r.db.QueryContext(ctx, query, fmt.Sprintf("-%d", months))
	if err != nil {
		return nil, apperror.Database("get monthly expense trend", err)
	}
	defer rows.Close()

	var results []domain.MonthlyTrend
	for rows.Next() {
		var mt domain.MonthlyTrend
		if err := rows.Scan(&mt.Month, &mt.Expenses); err != nil {
			return nil, apperror.Database("scan monthly trend", err)
		}
		results = append(results, mt)
	}
	return results, nil
}

// --- Scan helpers ---

func (r *ExpenseRepository) scanCategory(row *sql.Row) (*domain.ExpenseCategory, error) {
	var cat domain.ExpenseCategory
	var idStr, createdStr, updatedStr string
	var isActive int
	err := row.Scan(&idStr, &cat.Name, &cat.Description, &isActive, &createdStr, &updatedStr)
	if err != nil {
		return nil, err
	}
	cat.ID, _ = uuid.Parse(idStr)
	cat.IsActive = isActive == 1
	cat.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	cat.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
	return &cat, nil
}

func (r *ExpenseRepository) scanCategoryRow(rows *sql.Rows) (*domain.ExpenseCategory, error) {
	var cat domain.ExpenseCategory
	var idStr, createdStr, updatedStr string
	var isActive int
	err := rows.Scan(&idStr, &cat.Name, &cat.Description, &isActive, &createdStr, &updatedStr)
	if err != nil {
		return nil, err
	}
	cat.ID, _ = uuid.Parse(idStr)
	cat.IsActive = isActive == 1
	cat.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	cat.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
	return &cat, nil
}

func (r *ExpenseRepository) scanExpense(row *sql.Row) (*domain.Expense, error) {
	var exp domain.Expense
	var idStr, catIDStr, createdStr, updatedStr string
	err := row.Scan(&idStr, &exp.ExpenseNumber, &catIDStr, &exp.CategoryName, &exp.Amount,
		&exp.ExpenseDate, &exp.PaymentMethod, &exp.VendorName, &exp.InvoiceReference,
		&exp.Description, &exp.AttachmentPath, &exp.Status, &exp.CreatedBy, &createdStr, &updatedStr)
	if err != nil {
		return nil, err
	}
	exp.ID, _ = uuid.Parse(idStr)
	exp.CategoryID, _ = uuid.Parse(catIDStr)
	exp.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	exp.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
	return &exp, nil
}

func (r *ExpenseRepository) scanExpenseRow(rows *sql.Rows) (*domain.Expense, error) {
	var exp domain.Expense
	var idStr, catIDStr, createdStr, updatedStr string
	err := rows.Scan(&idStr, &exp.ExpenseNumber, &catIDStr, &exp.CategoryName, &exp.Amount,
		&exp.ExpenseDate, &exp.PaymentMethod, &exp.VendorName, &exp.InvoiceReference,
		&exp.Description, &exp.AttachmentPath, &exp.Status, &exp.CreatedBy, &createdStr, &updatedStr)
	if err != nil {
		return nil, err
	}
	exp.ID, _ = uuid.Parse(idStr)
	exp.CategoryID, _ = uuid.Parse(catIDStr)
	exp.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	exp.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
	return &exp, nil
}
