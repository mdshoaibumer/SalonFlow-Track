package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/pkg/uid"
)

// Payment methods for expenses.
var ValidExpensePaymentMethods = map[string]bool{
	"cash":          true,
	"upi":           true,
	"bank_transfer": true,
	"card":          true,
	"cheque":        true,
}

// Expense statuses.
var ValidExpenseStatuses = map[string]bool{
	"pending":  true,
	"approved": true,
	"paid":     true,
	"rejected": true,
}

// ExpenseCategory represents an expense category.
type ExpenseCategory struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewExpenseCategory creates a new ExpenseCategory.
func NewExpenseCategory(name, description string) *ExpenseCategory {
	now := time.Now().UTC()
	return &ExpenseCategory{
		ID:          uid.New(),
		Name:        name,
		Description: description,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Validate validates an ExpenseCategory.
func (c *ExpenseCategory) Validate() error {
	if c.Name == "" {
		return ErrExpenseCategoryNameRequired
	}
	return nil
}

// Expense represents a single expense record.
type Expense struct {
	ID               uuid.UUID `json:"id"`
	ExpenseNumber    string    `json:"expense_number"`
	CategoryID       uuid.UUID `json:"category_id"`
	CategoryName     string    `json:"category_name,omitempty"`
	Amount           float64   `json:"amount"`
	ExpenseDate      string    `json:"expense_date"`
	PaymentMethod    string    `json:"payment_method"`
	VendorName       string    `json:"vendor_name"`
	InvoiceReference string    `json:"invoice_reference"`
	Description      string    `json:"description"`
	AttachmentPath   string    `json:"attachment_path"`
	Status           string    `json:"status"`
	CreatedBy        string    `json:"created_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// NewExpense creates a new Expense with generated ID and defaults.
func NewExpense(categoryID uuid.UUID, amount float64, expenseDate, paymentMethod, vendorName, invoiceRef, description, createdBy string) *Expense {
	now := time.Now().UTC()
	return &Expense{
		ID:               uid.New(),
		CategoryID:       categoryID,
		Amount:           amount,
		ExpenseDate:      expenseDate,
		PaymentMethod:    paymentMethod,
		VendorName:       vendorName,
		InvoiceReference: invoiceRef,
		Description:      description,
		Status:           "pending",
		CreatedBy:        createdBy,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// Validate validates an Expense.
func (e *Expense) Validate() error {
	if e.CategoryID == uuid.Nil {
		return ErrExpenseCategoryRequired
	}
	if e.Amount <= 0 {
		return ErrExpenseInvalidAmount
	}
	if e.ExpenseDate == "" {
		return ErrExpenseDateRequired
	}
	if !ValidExpensePaymentMethods[e.PaymentMethod] {
		return ErrExpenseInvalidPaymentMethod
	}
	return nil
}

// Approve marks the expense as approved.
func (e *Expense) Approve() {
	e.Status = "approved"
	e.UpdatedAt = time.Now().UTC()
}

// MarkPaid marks the expense as paid.
func (e *Expense) MarkPaid() {
	e.Status = "paid"
	e.UpdatedAt = time.Now().UTC()
}

// Reject marks the expense as rejected.
func (e *Expense) Reject() {
	e.Status = "rejected"
	e.UpdatedAt = time.Now().UTC()
}

// ProfitLoss represents a profit & loss summary for a period.
type ProfitLoss struct {
	Period             string            `json:"period"`
	TotalRevenue       float64           `json:"total_revenue"`
	TotalExpenses      float64           `json:"total_expenses"`
	GrossProfit        float64           `json:"gross_profit"`
	ProfitMargin       float64           `json:"profit_margin"`
	ExpensesByCategory []CategoryExpense `json:"expenses_by_category"`
}

// CategoryExpense represents expense total for a category.
type CategoryExpense struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Amount       float64 `json:"amount"`
	Percentage   float64 `json:"percentage"`
}

// ExpenseStats represents dashboard stats.
type ExpenseStats struct {
	TodayExpenses   float64 `json:"today_expenses"`
	MonthlyExpenses float64 `json:"monthly_expenses"`
	TodayRevenue    float64 `json:"today_revenue"`
	MonthlyRevenue  float64 `json:"monthly_revenue"`
	MonthlyProfit   float64 `json:"monthly_profit"`
	ProfitMargin    float64 `json:"profit_margin"`
}

// MonthlyTrend represents a monthly data point for charts.
type MonthlyTrend struct {
	Month    string  `json:"month"`
	Revenue  float64 `json:"revenue"`
	Expenses float64 `json:"expenses"`
	Profit   float64 `json:"profit"`
}
