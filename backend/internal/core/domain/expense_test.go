package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestExpense_Validate(t *testing.T) {
	tests := []struct {
		name    string
		exp     Expense
		wantErr error
	}{
		{
			name: "valid expense",
			exp: Expense{
				CategoryID:    uuid.New(),
				Amount:        5000,
				ExpenseDate:   "2026-06-01",
				PaymentMethod: "cash",
			},
			wantErr: nil,
		},
		{
			name: "missing category",
			exp: Expense{
				Amount:        5000,
				ExpenseDate:   "2026-06-01",
				PaymentMethod: "cash",
			},
			wantErr: ErrExpenseCategoryRequired,
		},
		{
			name: "zero amount",
			exp: Expense{
				CategoryID:    uuid.New(),
				Amount:        0,
				ExpenseDate:   "2026-06-01",
				PaymentMethod: "cash",
			},
			wantErr: ErrExpenseInvalidAmount,
		},
		{
			name: "negative amount",
			exp: Expense{
				CategoryID:    uuid.New(),
				Amount:        -100,
				ExpenseDate:   "2026-06-01",
				PaymentMethod: "cash",
			},
			wantErr: ErrExpenseInvalidAmount,
		},
		{
			name: "missing date",
			exp: Expense{
				CategoryID:    uuid.New(),
				Amount:        5000,
				PaymentMethod: "cash",
			},
			wantErr: ErrExpenseDateRequired,
		},
		{
			name: "invalid payment method",
			exp: Expense{
				CategoryID:    uuid.New(),
				Amount:        5000,
				ExpenseDate:   "2026-06-01",
				PaymentMethod: "bitcoin",
			},
			wantErr: ErrExpenseInvalidPaymentMethod,
		},
		{
			name: "valid UPI",
			exp: Expense{
				CategoryID:    uuid.New(),
				Amount:        1000,
				ExpenseDate:   "2026-06-01",
				PaymentMethod: "upi",
			},
			wantErr: nil,
		},
		{
			name: "valid bank transfer",
			exp: Expense{
				CategoryID:    uuid.New(),
				Amount:        50000,
				ExpenseDate:   "2026-06-01",
				PaymentMethod: "bank_transfer",
			},
			wantErr: nil,
		},
		{
			name: "valid cheque",
			exp: Expense{
				CategoryID:    uuid.New(),
				Amount:        25000,
				ExpenseDate:   "2026-06-01",
				PaymentMethod: "cheque",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.exp.Validate()
			if err != tt.wantErr {
				t.Errorf("Validate() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestExpenseCategory_Validate(t *testing.T) {
	t.Run("valid category", func(t *testing.T) {
		cat := NewExpenseCategory("Rent", "Monthly rent")
		if err := cat.Validate(); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("empty name", func(t *testing.T) {
		cat := NewExpenseCategory("", "No name")
		if err := cat.Validate(); err != ErrExpenseCategoryNameRequired {
			t.Errorf("expected ErrExpenseCategoryNameRequired, got %v", err)
		}
	})
}

func TestExpense_StatusTransitions(t *testing.T) {
	catID := uuid.New()
	exp := NewExpense(catID, 5000, "2026-06-01", "cash", "Vendor", "", "Test", "")

	if exp.Status != "pending" {
		t.Errorf("new expense should be pending, got %s", exp.Status)
	}

	exp.Approve()
	if exp.Status != "approved" {
		t.Errorf("after Approve() should be approved, got %s", exp.Status)
	}

	exp.MarkPaid()
	if exp.Status != "paid" {
		t.Errorf("after MarkPaid() should be paid, got %s", exp.Status)
	}

	// Test reject path
	exp2 := NewExpense(catID, 1000, "2026-06-02", "upi", "V2", "", "Test2", "")
	exp2.Reject()
	if exp2.Status != "rejected" {
		t.Errorf("after Reject() should be rejected, got %s", exp2.Status)
	}
}

func TestNewExpense_GeneratesIDAndDefaults(t *testing.T) {
	catID := uuid.New()
	exp := NewExpense(catID, 10000, "2026-06-01", "card", "Amazon", "INV-001", "Office supplies", "admin")

	if exp.ID == uuid.Nil {
		t.Error("expected non-nil ID")
	}
	if exp.CategoryID != catID {
		t.Errorf("expected category ID %s, got %s", catID, exp.CategoryID)
	}
	if exp.Amount != 10000 {
		t.Errorf("expected 10000, got %f", exp.Amount)
	}
	if exp.Status != "pending" {
		t.Errorf("expected pending, got %s", exp.Status)
	}
	if exp.CreatedBy != "admin" {
		t.Errorf("expected admin, got %s", exp.CreatedBy)
	}
	if exp.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}
