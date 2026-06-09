package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// ExpenseFilter holds filtering options for listing expenses.
type ExpenseFilter struct {
	CategoryID    string
	Status        string
	PaymentMethod string
	DateFrom      string
	DateTo        string
	Search        string
	Limit         int
	Offset        int
}

// ExpenseRepository defines the contract for expense persistence.
type ExpenseRepository interface {
	// Categories
	CreateCategory(ctx context.Context, cat *domain.ExpenseCategory) error
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*domain.ExpenseCategory, error)
	ListCategories(ctx context.Context, activeOnly bool) ([]domain.ExpenseCategory, error)
	UpdateCategory(ctx context.Context, cat *domain.ExpenseCategory) error

	// Expenses
	CreateExpense(ctx context.Context, exp *domain.Expense) error
	GetExpenseByID(ctx context.Context, id uuid.UUID) (*domain.Expense, error)
	UpdateExpense(ctx context.Context, exp *domain.Expense) error
	DeleteExpense(ctx context.Context, id uuid.UUID) error
	ListExpenses(ctx context.Context, filter ExpenseFilter) ([]domain.Expense, int, error)

	// Number generation
	NextExpenseNumber(ctx context.Context, year int) (string, error)

	// Reporting
	GetTotalExpensesByDateRange(ctx context.Context, dateFrom, dateTo string) (float64, error)
	GetExpensesByCategory(ctx context.Context, dateFrom, dateTo string) ([]domain.CategoryExpense, error)
	GetMonthlyExpenseTrend(ctx context.Context, months int) ([]domain.MonthlyTrend, error)
}
