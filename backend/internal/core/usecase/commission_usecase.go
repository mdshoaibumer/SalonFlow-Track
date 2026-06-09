package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// CommissionUseCase handles commission business logic.
type CommissionUseCase struct {
	commissionRepo ports.CommissionRepository
	log            *slog.Logger
}

// NewCommissionUseCase creates a new CommissionUseCase.
func NewCommissionUseCase(commissionRepo ports.CommissionRepository, log *slog.Logger) *CommissionUseCase {
	return &CommissionUseCase{commissionRepo: commissionRepo, log: log}
}

// CreateRuleInput is the input DTO for creating a commission rule.
type CreateRuleInput struct {
	RuleName         string  `json:"rule_name"`
	RuleType         string  `json:"rule_type"`
	TargetType       string  `json:"target_type"`
	TargetID         string  `json:"target_id"`
	CalculationType  string  `json:"calculation_type"`
	CalculationValue float64 `json:"calculation_value"`
	MinimumTarget    float64 `json:"minimum_target"`
	MaximumTarget    float64 `json:"maximum_target"`
}

// UpdateRuleInput is the input DTO for updating a commission rule.
type UpdateRuleInput struct {
	RuleName         string  `json:"rule_name"`
	RuleType         string  `json:"rule_type"`
	TargetType       string  `json:"target_type"`
	TargetID         string  `json:"target_id"`
	CalculationType  string  `json:"calculation_type"`
	CalculationValue float64 `json:"calculation_value"`
	MinimumTarget    float64 `json:"minimum_target"`
	MaximumTarget    float64 `json:"maximum_target"`
	IsActive         bool    `json:"is_active"`
}

// ListRulesInput is the input DTO for listing rules.
type ListRulesInput struct {
	RuleType   string `json:"rule_type"`
	TargetType string `json:"target_type"`
	IsActive   *bool  `json:"is_active"`
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
}

// ListRulesOutput is the output DTO for listing rules.
type ListRulesOutput struct {
	Rules      []domain.CommissionRule `json:"rules"`
	Total      int                     `json:"total"`
	Page       int                     `json:"page"`
	PerPage    int                     `json:"per_page"`
	TotalPages int                     `json:"total_pages"`
}

// CreateRule creates a new commission rule.
func (uc *CommissionUseCase) CreateRule(ctx context.Context, input CreateRuleInput) (*domain.CommissionRule, error) {
	rule := domain.NewCommissionRule(
		input.RuleName, input.RuleType, input.TargetType, input.TargetID,
		input.CalculationType, input.CalculationValue, input.MinimumTarget, input.MaximumTarget,
	)

	if err := rule.Validate(); err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: err.Error()}
	}

	if err := uc.commissionRepo.CreateRule(ctx, rule); err != nil {
		return nil, err
	}

	uc.log.Info("commission rule created", "id", rule.ID, "name", rule.RuleName)
	return rule, nil
}

// GetRuleByID retrieves a rule by ID.
func (uc *CommissionUseCase) GetRuleByID(ctx context.Context, id uuid.UUID) (*domain.CommissionRule, error) {
	return uc.commissionRepo.GetRuleByID(ctx, id)
}

// ListRules returns paginated commission rules.
func (uc *CommissionUseCase) ListRules(ctx context.Context, input ListRulesInput) (*ListRulesOutput, error) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PerPage <= 0 {
		input.PerPage = 20
	}
	if input.PerPage > 100 {
		input.PerPage = 100
	}

	filter := ports.CommissionRuleFilter{
		RuleType:   input.RuleType,
		TargetType: input.TargetType,
		IsActive:   input.IsActive,
		Limit:      input.PerPage,
		Offset:     (input.Page - 1) * input.PerPage,
	}

	rules, total, err := uc.commissionRepo.ListRules(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := total / input.PerPage
	if total%input.PerPage > 0 {
		totalPages++
	}

	return &ListRulesOutput{
		Rules:      rules,
		Total:      total,
		Page:       input.Page,
		PerPage:    input.PerPage,
		TotalPages: totalPages,
	}, nil
}

// UpdateRule updates a commission rule.
func (uc *CommissionUseCase) UpdateRule(ctx context.Context, id uuid.UUID, input UpdateRuleInput) (*domain.CommissionRule, error) {
	rule, err := uc.commissionRepo.GetRuleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	rule.RuleName = input.RuleName
	rule.RuleType = input.RuleType
	rule.TargetType = input.TargetType
	rule.TargetID = input.TargetID
	rule.CalculationType = input.CalculationType
	rule.CalculationValue = input.CalculationValue
	rule.MinimumTarget = input.MinimumTarget
	rule.MaximumTarget = input.MaximumTarget
	rule.IsActive = input.IsActive
	rule.UpdatedAt = time.Now().UTC()

	if err := rule.Validate(); err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: err.Error()}
	}

	if err := uc.commissionRepo.UpdateRule(ctx, rule); err != nil {
		return nil, err
	}

	uc.log.Info("commission rule updated", "id", rule.ID)
	return rule, nil
}

// DeleteRule deletes a commission rule.
func (uc *CommissionUseCase) DeleteRule(ctx context.Context, id uuid.UUID) error {
	return uc.commissionRepo.DeleteRule(ctx, id)
}

// CalculateInvoiceCommission calculates and records commission for an invoice.
func (uc *CommissionUseCase) CalculateInvoiceCommission(ctx context.Context, staffID, invoiceID uuid.UUID, revenue float64, serviceIDs []uuid.UUID) (float64, error) {
	businessDate := time.Now().UTC().Format("2006-01-02")

	rules, err := uc.commissionRepo.GetActiveRules(ctx)
	if err != nil {
		return 0, err
	}

	var totalCommission float64

	for _, rule := range rules {
		var applicable bool
		var commAmount float64

		switch rule.RuleType {
		case domain.RuleTypeRevenueBased:
			// Revenue-based rules apply to total invoice revenue
			if rule.TargetType == domain.TargetTypeGlobal || (rule.TargetType == domain.TargetTypeStaff && rule.TargetID == staffID.String()) {
				applicable = true
				commAmount = rule.CalculateCommission(revenue)
			}
		case domain.RuleTypeServiceBased:
			// Service-based rules apply if service matches
			if rule.TargetType == domain.TargetTypeService {
				for _, svcID := range serviceIDs {
					if rule.TargetID == svcID.String() {
						applicable = true
						commAmount += rule.CalculateCommission(revenue / float64(len(serviceIDs)))
					}
				}
			}
		case domain.RuleTypeFixed:
			// Fixed rules apply per invoice
			if rule.TargetType == domain.TargetTypeGlobal || (rule.TargetType == domain.TargetTypeStaff && rule.TargetID == staffID.String()) {
				applicable = true
				commAmount = rule.CalculationValue
			}
		}

		if applicable && commAmount > 0 {
			tx := domain.NewCommissionTransaction(staffID, invoiceID, rule.ID, revenue, commAmount, businessDate)
			if err := uc.commissionRepo.CreateTransaction(ctx, tx); err != nil {
				uc.log.Error("failed to create commission transaction", "error", err)
				continue
			}
			totalCommission += commAmount
		}
	}

	return totalCommission, nil
}

// GetStaffCommissionInput is the input for staff commission query.
type GetStaffCommissionInput struct {
	StaffID  string `json:"staff_id"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

// StaffCommissionOutput is the output for staff commission query.
type StaffCommissionOutput struct {
	StaffID      uuid.UUID                      `json:"staff_id"`
	TotalRevenue float64                        `json:"total_revenue"`
	Commission   float64                        `json:"commission"`
	Transactions []domain.CommissionTransaction `json:"transactions"`
}

// GetStaffCommission returns commission details for a staff member.
func (uc *CommissionUseCase) GetStaffCommission(ctx context.Context, input GetStaffCommissionInput) (*StaffCommissionOutput, error) {
	staffID, err := uuid.Parse(input.StaffID)
	if err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid staff ID"}
	}

	if input.DateFrom == "" {
		now := time.Now().UTC()
		input.DateFrom = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
		input.DateTo = now.Format("2006-01-02")
	}

	commission, err := uc.commissionRepo.GetStaffCommission(ctx, staffID, input.DateFrom, input.DateTo)
	if err != nil {
		return nil, err
	}

	txFilter := ports.CommissionTxFilter{
		StaffID:  staffID.String(),
		DateFrom: input.DateFrom,
		DateTo:   input.DateTo,
		Limit:    100,
	}
	txns, _, err := uc.commissionRepo.ListTransactions(ctx, txFilter)
	if err != nil {
		return nil, err
	}

	var totalRevenue float64
	for _, tx := range txns {
		totalRevenue += tx.RevenueAmount
	}

	return &StaffCommissionOutput{
		StaffID:      staffID,
		TotalRevenue: totalRevenue,
		Commission:   commission,
		Transactions: txns,
	}, nil
}

// MonthlyCommissionInput is the input for monthly commission query.
type MonthlyCommissionInput struct {
	Month string `json:"month"` // YYYY-MM
}

// GetMonthlyCommission returns commission summary for a month.
func (uc *CommissionUseCase) GetMonthlyCommission(ctx context.Context, input MonthlyCommissionInput) ([]ports.CommissionStaffSummary, error) {
	if input.Month == "" {
		input.Month = time.Now().UTC().Format("2006-01")
	}
	return uc.commissionRepo.GetMonthlyCommission(ctx, input.Month)
}

// CommissionStats holds dashboard commission statistics.
type CommissionStats struct {
	TotalCommissionThisMonth float64                       `json:"total_commission_this_month"`
	TopEarner                *ports.CommissionStaffSummary `json:"top_earner"`
	AvgCommission            float64                       `json:"avg_commission"`
}

// GetStats returns commission dashboard stats.
func (uc *CommissionUseCase) GetStats(ctx context.Context) (*CommissionStats, error) {
	month := time.Now().UTC().Format("2006-01")
	summaries, err := uc.commissionRepo.GetMonthlyCommission(ctx, month)
	if err != nil {
		return nil, err
	}

	stats := &CommissionStats{}
	for _, s := range summaries {
		stats.TotalCommissionThisMonth += s.Commission
	}
	if len(summaries) > 0 {
		stats.TopEarner = &summaries[0]
		stats.AvgCommission = stats.TotalCommissionThisMonth / float64(len(summaries))
	}
	return stats, nil
}
