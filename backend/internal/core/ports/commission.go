package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// CommissionRepository defines persistence operations for commission rules and transactions.
type CommissionRepository interface {
	// Rules
	CreateRule(ctx context.Context, rule *domain.CommissionRule) error
	GetRuleByID(ctx context.Context, id uuid.UUID) (*domain.CommissionRule, error)
	ListRules(ctx context.Context, filter CommissionRuleFilter) ([]domain.CommissionRule, int, error)
	UpdateRule(ctx context.Context, rule *domain.CommissionRule) error
	DeleteRule(ctx context.Context, id uuid.UUID) error
	GetActiveRules(ctx context.Context) ([]domain.CommissionRule, error)

	// Transactions
	CreateTransaction(ctx context.Context, tx *domain.CommissionTransaction) error
	ListTransactions(ctx context.Context, filter CommissionTxFilter) ([]domain.CommissionTransaction, int, error)
	GetStaffCommission(ctx context.Context, staffID uuid.UUID, dateFrom, dateTo string) (float64, error)
	GetMonthlyCommission(ctx context.Context, month string) ([]CommissionStaffSummary, error)
}

// CommissionRuleFilter holds query params for listing rules.
type CommissionRuleFilter struct {
	RuleType   string
	TargetType string
	IsActive   *bool
	Limit      int
	Offset     int
}

// CommissionTxFilter holds query params for listing transactions.
type CommissionTxFilter struct {
	StaffID  string
	DateFrom string
	DateTo   string
	Status   string
	Limit    int
	Offset   int
}

// CommissionStaffSummary aggregates commission for a staff member.
type CommissionStaffSummary struct {
	StaffID    uuid.UUID `json:"staff_id"`
	StaffName  string    `json:"staff_name"`
	Revenue    float64   `json:"revenue"`
	Commission float64   `json:"commission"`
}
