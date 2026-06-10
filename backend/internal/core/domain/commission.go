package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/pkg/uid"
)

// Commission rule types.
const (
	RuleTypeRevenueBased = "revenue_based"
	RuleTypeServiceBased = "service_based"
	RuleTypeFixed        = "fixed"
)

// Commission target types.
const (
	TargetTypeGlobal  = "global"
	TargetTypeStaff   = "staff"
	TargetTypeService = "service"
)

// Commission calculation types.
const (
	CalcTypePercentage  = "percentage"
	CalcTypeFixedAmount = "fixed_amount"
	CalcTypeTiered      = "tiered"
)

// Commission transaction statuses.
const (
	CommissionStatusPending  = "pending"
	CommissionStatusApproved = "approved"
	CommissionStatusPaid     = "paid"
)

var validRuleTypes = map[string]bool{
	RuleTypeRevenueBased: true,
	RuleTypeServiceBased: true,
	RuleTypeFixed:        true,
}

var validTargetTypes = map[string]bool{
	TargetTypeGlobal:  true,
	TargetTypeStaff:   true,
	TargetTypeService: true,
}

var validCalcTypes = map[string]bool{
	CalcTypePercentage:  true,
	CalcTypeFixedAmount: true,
	CalcTypeTiered:      true,
}

// CommissionRule defines a commission calculation rule.
type CommissionRule struct {
	ID               uuid.UUID `json:"id"`
	RuleName         string    `json:"rule_name"`
	RuleType         string    `json:"rule_type"`
	TargetType       string    `json:"target_type"`
	TargetID         string    `json:"target_id"`
	CalculationType  string    `json:"calculation_type"`
	CalculationValue float64   `json:"calculation_value"`
	MinimumTarget    float64   `json:"minimum_target"`
	MaximumTarget    float64   `json:"maximum_target"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// CommissionTransaction records a commission earned on an invoice.
type CommissionTransaction struct {
	ID               uuid.UUID `json:"id"`
	StaffID          uuid.UUID `json:"staff_id"`
	InvoiceID        uuid.UUID `json:"invoice_id"`
	RuleID           uuid.UUID `json:"rule_id"`
	RevenueAmount    float64   `json:"revenue_amount"`
	CommissionAmount float64   `json:"commission_amount"`
	BusinessDate     string    `json:"business_date"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// NewCommissionRule creates a new commission rule.
func NewCommissionRule(name, ruleType, targetType, targetID, calcType string, calcValue, minTarget, maxTarget float64) *CommissionRule {
	now := time.Now().UTC()
	return &CommissionRule{
		ID:               uid.New(),
		RuleName:         name,
		RuleType:         ruleType,
		TargetType:       targetType,
		TargetID:         targetID,
		CalculationType:  calcType,
		CalculationValue: calcValue,
		MinimumTarget:    minTarget,
		MaximumTarget:    maxTarget,
		IsActive:         true,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// Validate validates the commission rule.
func (r *CommissionRule) Validate() error {
	if r.RuleName == "" {
		return ErrCommissionRuleNameRequired
	}
	if !validRuleTypes[r.RuleType] {
		return ErrCommissionInvalidRuleType
	}
	if !validTargetTypes[r.TargetType] {
		return ErrCommissionInvalidTargetType
	}
	if !validCalcTypes[r.CalculationType] {
		return ErrCommissionInvalidCalcType
	}
	if r.CalculationValue < 0 {
		return ErrCommissionInvalidCalcValue
	}
	return nil
}

// CalculateCommission calculates commission for a given revenue amount.
func (r *CommissionRule) CalculateCommission(revenue float64) float64 {
	// Check if revenue meets minimum target
	if r.MinimumTarget > 0 && revenue < r.MinimumTarget {
		return 0
	}
	// Check if revenue exceeds maximum target (cap it)
	if r.MaximumTarget > 0 && revenue > r.MaximumTarget {
		revenue = r.MaximumTarget
	}

	switch r.CalculationType {
	case CalcTypePercentage:
		return revenue * r.CalculationValue / 100
	case CalcTypeFixedAmount:
		return r.CalculationValue
	default:
		return 0
	}
}

// NewCommissionTransaction creates a new commission transaction.
func NewCommissionTransaction(staffID, invoiceID, ruleID uuid.UUID, revenueAmount, commissionAmount float64, businessDate string) *CommissionTransaction {
	now := time.Now().UTC()
	return &CommissionTransaction{
		ID:               uid.New(),
		StaffID:          staffID,
		InvoiceID:        invoiceID,
		RuleID:           ruleID,
		RevenueAmount:    revenueAmount,
		CommissionAmount: commissionAmount,
		BusinessDate:     businessDate,
		Status:           CommissionStatusPending,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}
