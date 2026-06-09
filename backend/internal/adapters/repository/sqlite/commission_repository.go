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

// CommissionRepository is the SQLite implementation of ports.CommissionRepository.
type CommissionRepository struct {
	db  *sql.DB
	log *slog.Logger
}

// NewCommissionRepository creates a new CommissionRepository.
func NewCommissionRepository(db *sql.DB, log *slog.Logger) *CommissionRepository {
	return &CommissionRepository{db: db, log: log}
}

// CreateRule inserts a new commission rule.
func (r *CommissionRepository) CreateRule(ctx context.Context, rule *domain.CommissionRule) error {
	query := `
		INSERT INTO commission_rules (id, rule_name, rule_type, target_type, target_id, calculation_type, calculation_value, minimum_target, maximum_target, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		rule.ID.String(),
		rule.RuleName,
		rule.RuleType,
		rule.TargetType,
		rule.TargetID,
		rule.CalculationType,
		rule.CalculationValue,
		rule.MinimumTarget,
		rule.MaximumTarget,
		boolToInt(rule.IsActive),
		rule.CreatedAt.Format(time.RFC3339),
		rule.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return apperror.Database("create commission rule", err)
	}
	return nil
}

// GetRuleByID retrieves a commission rule by ID.
func (r *CommissionRepository) GetRuleByID(ctx context.Context, id uuid.UUID) (*domain.CommissionRule, error) {
	query := `
		SELECT id, rule_name, rule_type, target_type, target_id, calculation_type, calculation_value, minimum_target, maximum_target, is_active, created_at, updated_at
		FROM commission_rules WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id.String())
	rule, err := r.scanRule(row)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("commission_rule", id.String())
	}
	if err != nil {
		return nil, apperror.Database("get commission rule", err)
	}
	return rule, nil
}

// ListRules returns commission rules with filtering.
func (r *CommissionRepository) ListRules(ctx context.Context, filter ports.CommissionRuleFilter) ([]domain.CommissionRule, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}

	if filter.RuleType != "" {
		where = append(where, "rule_type = ?")
		args = append(args, filter.RuleType)
	}
	if filter.TargetType != "" {
		where = append(where, "target_type = ?")
		args = append(args, filter.TargetType)
	}
	if filter.IsActive != nil {
		where = append(where, "is_active = ?")
		args = append(args, boolToInt(*filter.IsActive))
	}

	whereClause := strings.Join(where, " AND ")

	// Count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM commission_rules WHERE %s", whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, apperror.Database("count commission rules", err)
	}

	// Query
	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	query := fmt.Sprintf(`
		SELECT id, rule_name, rule_type, target_type, target_id, calculation_type, calculation_value, minimum_target, maximum_target, is_active, created_at, updated_at
		FROM commission_rules WHERE %s ORDER BY created_at DESC LIMIT ? OFFSET ?`, whereClause)

	args = append(args, limit, filter.Offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, apperror.Database("list commission rules", err)
	}
	defer rows.Close()

	var rules []domain.CommissionRule
	for rows.Next() {
		rule, err := r.scanRuleRow(rows)
		if err != nil {
			return nil, 0, err
		}
		rules = append(rules, *rule)
	}
	return rules, total, nil
}

// UpdateRule updates a commission rule.
func (r *CommissionRepository) UpdateRule(ctx context.Context, rule *domain.CommissionRule) error {
	query := `
		UPDATE commission_rules SET
			rule_name = ?, rule_type = ?, target_type = ?, target_id = ?,
			calculation_type = ?, calculation_value = ?, minimum_target = ?,
			maximum_target = ?, is_active = ?, updated_at = ?
		WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query,
		rule.RuleName, rule.RuleType, rule.TargetType, rule.TargetID,
		rule.CalculationType, rule.CalculationValue, rule.MinimumTarget,
		rule.MaximumTarget, boolToInt(rule.IsActive), time.Now().UTC().Format(time.RFC3339),
		rule.ID.String(),
	)
	if err != nil {
		return apperror.Database("update commission rule", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apperror.NotFound("commission_rule", rule.ID.String())
	}
	return nil
}

// DeleteRule deletes a commission rule.
func (r *CommissionRepository) DeleteRule(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM commission_rules WHERE id = ?", id.String())
	if err != nil {
		return apperror.Database("delete commission rule", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apperror.NotFound("commission_rule", id.String())
	}
	return nil
}

// GetActiveRules returns all active commission rules.
func (r *CommissionRepository) GetActiveRules(ctx context.Context) ([]domain.CommissionRule, error) {
	query := `
		SELECT id, rule_name, rule_type, target_type, target_id, calculation_type, calculation_value, minimum_target, maximum_target, is_active, created_at, updated_at
		FROM commission_rules WHERE is_active = 1 ORDER BY rule_type, created_at`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, apperror.Database("get active rules", err)
	}
	defer rows.Close()

	var rules []domain.CommissionRule
	for rows.Next() {
		rule, err := r.scanRuleRow(rows)
		if err != nil {
			return nil, err
		}
		rules = append(rules, *rule)
	}
	return rules, nil
}

// CreateTransaction inserts a new commission transaction.
func (r *CommissionRepository) CreateTransaction(ctx context.Context, tx *domain.CommissionTransaction) error {
	query := `
		INSERT INTO commission_transactions (id, staff_id, invoice_id, rule_id, revenue_amount, commission_amount, business_date, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		tx.ID.String(),
		tx.StaffID.String(),
		tx.InvoiceID.String(),
		tx.RuleID.String(),
		tx.RevenueAmount,
		tx.CommissionAmount,
		tx.BusinessDate,
		tx.Status,
		tx.CreatedAt.Format(time.RFC3339),
		tx.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return apperror.Database("create commission transaction", err)
	}
	return nil
}

// ListTransactions returns commission transactions with filtering.
func (r *CommissionRepository) ListTransactions(ctx context.Context, filter ports.CommissionTxFilter) ([]domain.CommissionTransaction, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}

	if filter.StaffID != "" {
		where = append(where, "staff_id = ?")
		args = append(args, filter.StaffID)
	}
	if filter.DateFrom != "" {
		where = append(where, "business_date >= ?")
		args = append(args, filter.DateFrom)
	}
	if filter.DateTo != "" {
		where = append(where, "business_date <= ?")
		args = append(args, filter.DateTo)
	}
	if filter.Status != "" {
		where = append(where, "status = ?")
		args = append(args, filter.Status)
	}

	whereClause := strings.Join(where, " AND ")

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM commission_transactions WHERE %s", whereClause)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, apperror.Database("count commission transactions", err)
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	query := fmt.Sprintf(`
		SELECT id, staff_id, invoice_id, rule_id, revenue_amount, commission_amount, business_date, status, created_at, updated_at
		FROM commission_transactions WHERE %s ORDER BY business_date DESC, created_at DESC LIMIT ? OFFSET ?`, whereClause)

	args = append(args, limit, filter.Offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, apperror.Database("list commission transactions", err)
	}
	defer rows.Close()

	var txns []domain.CommissionTransaction
	for rows.Next() {
		var t domain.CommissionTransaction
		var idStr, staffStr, invoiceStr, ruleStr, createdStr, updatedStr string
		err := rows.Scan(&idStr, &staffStr, &invoiceStr, &ruleStr, &t.RevenueAmount, &t.CommissionAmount, &t.BusinessDate, &t.Status, &createdStr, &updatedStr)
		if err != nil {
			return nil, 0, apperror.Database("scan commission transaction", err)
		}
		t.ID, _ = uuid.Parse(idStr)
		t.StaffID, _ = uuid.Parse(staffStr)
		t.InvoiceID, _ = uuid.Parse(invoiceStr)
		t.RuleID, _ = uuid.Parse(ruleStr)
		t.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
		t.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
		txns = append(txns, t)
	}
	return txns, total, nil
}

// GetStaffCommission returns total commission for a staff member in a date range.
func (r *CommissionRepository) GetStaffCommission(ctx context.Context, staffID uuid.UUID, dateFrom, dateTo string) (float64, error) {
	query := `SELECT COALESCE(SUM(commission_amount), 0) FROM commission_transactions WHERE staff_id = ? AND business_date BETWEEN ? AND ?`
	var total float64
	err := r.db.QueryRowContext(ctx, query, staffID.String(), dateFrom, dateTo).Scan(&total)
	if err != nil {
		return 0, apperror.Database("get staff commission", err)
	}
	return total, nil
}

// GetMonthlyCommission returns commission summary for each staff member for a given month prefix (YYYY-MM).
func (r *CommissionRepository) GetMonthlyCommission(ctx context.Context, month string) ([]ports.CommissionStaffSummary, error) {
	query := `
		SELECT ct.staff_id, s.full_name, COALESCE(SUM(ct.revenue_amount), 0), COALESCE(SUM(ct.commission_amount), 0)
		FROM commission_transactions ct
		JOIN staff s ON ct.staff_id = s.id
		WHERE ct.business_date LIKE ?
		GROUP BY ct.staff_id, s.full_name
		ORDER BY SUM(ct.commission_amount) DESC`

	rows, err := r.db.QueryContext(ctx, query, month+"%")
	if err != nil {
		return nil, apperror.Database("get monthly commission", err)
	}
	defer rows.Close()

	var results []ports.CommissionStaffSummary
	for rows.Next() {
		var s ports.CommissionStaffSummary
		var staffIDStr string
		err := rows.Scan(&staffIDStr, &s.StaffName, &s.Revenue, &s.Commission)
		if err != nil {
			return nil, apperror.Database("scan monthly commission", err)
		}
		s.StaffID, _ = uuid.Parse(staffIDStr)
		results = append(results, s)
	}
	return results, nil
}

func (r *CommissionRepository) scanRule(row *sql.Row) (*domain.CommissionRule, error) {
	var rule domain.CommissionRule
	var idStr, createdStr, updatedStr string
	var isActive int
	err := row.Scan(&idStr, &rule.RuleName, &rule.RuleType, &rule.TargetType, &rule.TargetID,
		&rule.CalculationType, &rule.CalculationValue, &rule.MinimumTarget, &rule.MaximumTarget,
		&isActive, &createdStr, &updatedStr)
	if err != nil {
		return nil, err
	}
	rule.ID, _ = uuid.Parse(idStr)
	rule.IsActive = isActive == 1
	rule.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	rule.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
	return &rule, nil
}

func (r *CommissionRepository) scanRuleRow(rows *sql.Rows) (*domain.CommissionRule, error) {
	var rule domain.CommissionRule
	var idStr, createdStr, updatedStr string
	var isActive int
	err := rows.Scan(&idStr, &rule.RuleName, &rule.RuleType, &rule.TargetType, &rule.TargetID,
		&rule.CalculationType, &rule.CalculationValue, &rule.MinimumTarget, &rule.MaximumTarget,
		&isActive, &createdStr, &updatedStr)
	if err != nil {
		return nil, apperror.Database("scan commission rule", err)
	}
	rule.ID, _ = uuid.Parse(idStr)
	rule.IsActive = isActive == 1
	rule.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	rule.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
	return &rule, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
