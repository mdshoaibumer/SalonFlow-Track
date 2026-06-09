package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewCommissionRule(t *testing.T) {
	rule := NewCommissionRule("Revenue 10%", RuleTypeRevenueBased, TargetTypeGlobal, "", CalcTypePercentage, 10, 0, 0)

	if rule.ID == uuid.Nil {
		t.Fatal("expected non-nil ID")
	}
	if rule.RuleName != "Revenue 10%" {
		t.Fatalf("expected name 'Revenue 10%%', got %s", rule.RuleName)
	}
	if !rule.IsActive {
		t.Fatal("expected rule to be active by default")
	}
}

func TestCommissionRule_Validate(t *testing.T) {
	tests := []struct {
		name    string
		rule    *CommissionRule
		wantErr bool
	}{
		{
			name:    "valid rule",
			rule:    NewCommissionRule("Test", RuleTypeRevenueBased, TargetTypeGlobal, "", CalcTypePercentage, 10, 0, 0),
			wantErr: false,
		},
		{
			name:    "empty name",
			rule:    NewCommissionRule("", RuleTypeRevenueBased, TargetTypeGlobal, "", CalcTypePercentage, 10, 0, 0),
			wantErr: true,
		},
		{
			name:    "invalid rule type",
			rule:    NewCommissionRule("Test", "invalid", TargetTypeGlobal, "", CalcTypePercentage, 10, 0, 0),
			wantErr: true,
		},
		{
			name:    "invalid target type",
			rule:    NewCommissionRule("Test", RuleTypeFixed, "invalid", "", CalcTypeFixedAmount, 100, 0, 0),
			wantErr: true,
		},
		{
			name:    "invalid calc type",
			rule:    NewCommissionRule("Test", RuleTypeFixed, TargetTypeGlobal, "", "invalid", 100, 0, 0),
			wantErr: true,
		},
		{
			name:    "negative calc value",
			rule:    NewCommissionRule("Test", RuleTypeFixed, TargetTypeGlobal, "", CalcTypeFixedAmount, -100, 0, 0),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rule.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommissionRule_CalculateCommission(t *testing.T) {
	tests := []struct {
		name     string
		rule     *CommissionRule
		revenue  float64
		expected float64
	}{
		{
			name:     "percentage 10% of 10000",
			rule:     NewCommissionRule("Test", RuleTypeRevenueBased, TargetTypeGlobal, "", CalcTypePercentage, 10, 0, 0),
			revenue:  10000,
			expected: 1000,
		},
		{
			name:     "fixed amount",
			rule:     NewCommissionRule("Test", RuleTypeFixed, TargetTypeGlobal, "", CalcTypeFixedAmount, 500, 0, 0),
			revenue:  10000,
			expected: 500,
		},
		{
			name:     "below minimum target",
			rule:     NewCommissionRule("Test", RuleTypeRevenueBased, TargetTypeGlobal, "", CalcTypePercentage, 10, 5000, 0),
			revenue:  3000,
			expected: 0,
		},
		{
			name:     "above minimum target",
			rule:     NewCommissionRule("Test", RuleTypeRevenueBased, TargetTypeGlobal, "", CalcTypePercentage, 10, 5000, 0),
			revenue:  8000,
			expected: 800,
		},
		{
			name:     "capped at maximum target",
			rule:     NewCommissionRule("Test", RuleTypeRevenueBased, TargetTypeGlobal, "", CalcTypePercentage, 10, 0, 10000),
			revenue:  20000,
			expected: 1000, // 10% of capped 10000
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rule.CalculateCommission(tt.revenue)
			if result != tt.expected {
				t.Errorf("CalculateCommission(%v) = %v, want %v", tt.revenue, result, tt.expected)
			}
		})
	}
}

func TestNewCommissionTransaction(t *testing.T) {
	staffID := uuid.New()
	invoiceID := uuid.New()
	ruleID := uuid.New()

	tx := NewCommissionTransaction(staffID, invoiceID, ruleID, 5000, 500, "2026-06-09")

	if tx.ID == uuid.Nil {
		t.Fatal("expected non-nil ID")
	}
	if tx.StaffID != staffID {
		t.Fatalf("expected staff ID %s, got %s", staffID, tx.StaffID)
	}
	if tx.InvoiceID != invoiceID {
		t.Fatalf("expected invoice ID %s, got %s", invoiceID, tx.InvoiceID)
	}
	if tx.Status != CommissionStatusPending {
		t.Fatalf("expected pending status, got %s", tx.Status)
	}
	if tx.CommissionAmount != 500 {
		t.Fatalf("expected commission 500, got %.2f", tx.CommissionAmount)
	}
}
