package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// ExpenseUseCase handles expense business logic.
type ExpenseUseCase struct {
	expenseRepo ports.ExpenseRepository
	invoiceRepo ports.InvoiceRepository
	log         *slog.Logger
}

// NewExpenseUseCase creates a new ExpenseUseCase.
func NewExpenseUseCase(expenseRepo ports.ExpenseRepository, invoiceRepo ports.InvoiceRepository, log *slog.Logger) *ExpenseUseCase {
	return &ExpenseUseCase{
		expenseRepo: expenseRepo,
		invoiceRepo: invoiceRepo,
		log:         log,
	}
}

// --- Input/Output DTOs ---

type CreateExpenseInput struct {
	CategoryID       string  `json:"category_id"`
	Amount           float64 `json:"amount"`
	ExpenseDate      string  `json:"expense_date"`
	PaymentMethod    string  `json:"payment_method"`
	VendorName       string  `json:"vendor_name"`
	InvoiceReference string  `json:"invoice_reference"`
	Description      string  `json:"description"`
}

type UpdateExpenseInput struct {
	CategoryID       string  `json:"category_id"`
	Amount           float64 `json:"amount"`
	ExpenseDate      string  `json:"expense_date"`
	PaymentMethod    string  `json:"payment_method"`
	VendorName       string  `json:"vendor_name"`
	InvoiceReference string  `json:"invoice_reference"`
	Description      string  `json:"description"`
	Status           string  `json:"status"`
}

type ListExpensesInput struct {
	CategoryID    string
	Status        string
	PaymentMethod string
	DateFrom      string
	DateTo        string
	Search        string
	Page          int
	PerPage       int
}

type ListExpensesOutput struct {
	Expenses   []domain.Expense `json:"expenses"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	PerPage    int              `json:"per_page"`
	TotalPages int              `json:"total_pages"`
}

type ProfitLossInput struct {
	DateFrom string
	DateTo   string
}

type ExpenseReportInput struct {
	DateFrom string
	DateTo   string
}

type ExpenseReport struct {
	DateFrom           string                   `json:"date_from"`
	DateTo             string                   `json:"date_to"`
	TotalExpenses      float64                  `json:"total_expenses"`
	ExpensesByCategory []domain.CategoryExpense `json:"expenses_by_category"`
	ExpenseCount       int                      `json:"expense_count"`
}

// --- Category Operations ---

func (uc *ExpenseUseCase) CreateCategory(ctx context.Context, name, description string) (*domain.ExpenseCategory, error) {
	cat := domain.NewExpenseCategory(name, description)
	if err := cat.Validate(); err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: err.Error()}
	}
	if err := uc.expenseRepo.CreateCategory(ctx, cat); err != nil {
		return nil, err
	}
	uc.log.Info("expense category created", "id", cat.ID, "name", cat.Name)
	return cat, nil
}

func (uc *ExpenseUseCase) ListCategories(ctx context.Context, activeOnly bool) ([]domain.ExpenseCategory, error) {
	return uc.expenseRepo.ListCategories(ctx, activeOnly)
}

func (uc *ExpenseUseCase) UpdateCategory(ctx context.Context, id uuid.UUID, name, description string, isActive bool) (*domain.ExpenseCategory, error) {
	cat, err := uc.expenseRepo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	cat.Name = name
	cat.Description = description
	cat.IsActive = isActive
	if err := cat.Validate(); err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: err.Error()}
	}
	if err := uc.expenseRepo.UpdateCategory(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

// --- Expense CRUD ---

func (uc *ExpenseUseCase) CreateExpense(ctx context.Context, input CreateExpenseInput) (*domain.Expense, error) {
	categoryID, err := uuid.Parse(input.CategoryID)
	if err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid category ID"}
	}

	// Verify category exists
	_, err = uc.expenseRepo.GetCategoryByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	exp := domain.NewExpense(categoryID, input.Amount, input.ExpenseDate, input.PaymentMethod,
		input.VendorName, input.InvoiceReference, input.Description, "")

	if err := exp.Validate(); err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: err.Error()}
	}

	// Generate expense number
	year := time.Now().Year()
	if input.ExpenseDate != "" {
		t, err := time.Parse("2006-01-02", input.ExpenseDate)
		if err == nil {
			year = t.Year()
		}
	}
	expNum, err := uc.expenseRepo.NextExpenseNumber(ctx, year)
	if err != nil {
		return nil, err
	}
	exp.ExpenseNumber = expNum

	if err := uc.expenseRepo.CreateExpense(ctx, exp); err != nil {
		return nil, err
	}

	uc.log.Info("expense created", "id", exp.ID, "number", exp.ExpenseNumber, "amount", exp.Amount)
	return exp, nil
}

func (uc *ExpenseUseCase) GetExpense(ctx context.Context, id uuid.UUID) (*domain.Expense, error) {
	return uc.expenseRepo.GetExpenseByID(ctx, id)
}

func (uc *ExpenseUseCase) UpdateExpense(ctx context.Context, id uuid.UUID, input UpdateExpenseInput) (*domain.Expense, error) {
	exp, err := uc.expenseRepo.GetExpenseByID(ctx, id)
	if err != nil {
		return nil, err
	}

	categoryID, err := uuid.Parse(input.CategoryID)
	if err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid category ID"}
	}

	exp.CategoryID = categoryID
	exp.Amount = input.Amount
	exp.ExpenseDate = input.ExpenseDate
	exp.PaymentMethod = input.PaymentMethod
	exp.VendorName = input.VendorName
	exp.InvoiceReference = input.InvoiceReference
	exp.Description = input.Description
	if input.Status != "" && domain.ValidExpenseStatuses[input.Status] {
		exp.Status = input.Status
	}

	if err := exp.Validate(); err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: err.Error()}
	}

	if err := uc.expenseRepo.UpdateExpense(ctx, exp); err != nil {
		return nil, err
	}

	uc.log.Info("expense updated", "id", exp.ID)
	return exp, nil
}

func (uc *ExpenseUseCase) DeleteExpense(ctx context.Context, id uuid.UUID) error {
	// Verify it exists
	_, err := uc.expenseRepo.GetExpenseByID(ctx, id)
	if err != nil {
		return err
	}
	if err := uc.expenseRepo.DeleteExpense(ctx, id); err != nil {
		return err
	}
	uc.log.Info("expense deleted", "id", id)
	return nil
}

func (uc *ExpenseUseCase) ListExpenses(ctx context.Context, input ListExpensesInput) (*ListExpensesOutput, error) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PerPage <= 0 {
		input.PerPage = 20
	}
	if input.PerPage > 100 {
		input.PerPage = 100
	}

	filter := ports.ExpenseFilter{
		CategoryID:    input.CategoryID,
		Status:        input.Status,
		PaymentMethod: input.PaymentMethod,
		DateFrom:      input.DateFrom,
		DateTo:        input.DateTo,
		Search:        input.Search,
		Limit:         input.PerPage,
		Offset:        (input.Page - 1) * input.PerPage,
	}

	expenses, total, err := uc.expenseRepo.ListExpenses(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := total / input.PerPage
	if total%input.PerPage > 0 {
		totalPages++
	}

	return &ListExpensesOutput{
		Expenses:   expenses,
		Total:      total,
		Page:       input.Page,
		PerPage:    input.PerPage,
		TotalPages: totalPages,
	}, nil
}

// --- Profit & Loss ---

func (uc *ExpenseUseCase) GetProfitLoss(ctx context.Context, input ProfitLossInput) (*domain.ProfitLoss, error) {
	if input.DateFrom == "" || input.DateTo == "" {
		now := time.Now().UTC()
		input.DateFrom = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
		input.DateTo = now.Format("2006-01-02")
	}

	// Get revenue from invoices
	revenue, err := uc.getRevenueByDateRange(ctx, input.DateFrom, input.DateTo)
	if err != nil {
		return nil, err
	}

	// Get total expenses
	totalExpenses, err := uc.expenseRepo.GetTotalExpensesByDateRange(ctx, input.DateFrom, input.DateTo)
	if err != nil {
		return nil, err
	}

	// Get expense breakdown by category
	byCategory, err := uc.expenseRepo.GetExpensesByCategory(ctx, input.DateFrom, input.DateTo)
	if err != nil {
		return nil, err
	}

	grossProfit := revenue - totalExpenses
	var profitMargin float64
	if revenue > 0 {
		profitMargin = (grossProfit / revenue) * 100
	}

	period := fmt.Sprintf("%s to %s", input.DateFrom, input.DateTo)

	return &domain.ProfitLoss{
		Period:             period,
		TotalRevenue:       revenue,
		TotalExpenses:      totalExpenses,
		GrossProfit:        grossProfit,
		ProfitMargin:       profitMargin,
		ExpensesByCategory: byCategory,
	}, nil
}

func (uc *ExpenseUseCase) GetExpenseStats(ctx context.Context) (*domain.ExpenseStats, error) {
	today := time.Now().UTC().Format("2006-01-02")
	now := time.Now().UTC()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")

	todayExpenses, err := uc.expenseRepo.GetTotalExpensesByDateRange(ctx, today, today)
	if err != nil {
		return nil, err
	}

	monthlyExpenses, err := uc.expenseRepo.GetTotalExpensesByDateRange(ctx, monthStart, today)
	if err != nil {
		return nil, err
	}

	todayRevenue, err := uc.getRevenueByDateRange(ctx, today, today)
	if err != nil {
		return nil, err
	}

	monthlyRevenue, err := uc.getRevenueByDateRange(ctx, monthStart, today)
	if err != nil {
		return nil, err
	}

	monthlyProfit := monthlyRevenue - monthlyExpenses
	var profitMargin float64
	if monthlyRevenue > 0 {
		profitMargin = (monthlyProfit / monthlyRevenue) * 100
	}

	return &domain.ExpenseStats{
		TodayExpenses:   todayExpenses,
		MonthlyExpenses: monthlyExpenses,
		TodayRevenue:    todayRevenue,
		MonthlyRevenue:  monthlyRevenue,
		MonthlyProfit:   monthlyProfit,
		ProfitMargin:    profitMargin,
	}, nil
}

func (uc *ExpenseUseCase) GetExpenseReport(ctx context.Context, input ExpenseReportInput) (*ExpenseReport, error) {
	if input.DateFrom == "" || input.DateTo == "" {
		now := time.Now().UTC()
		input.DateFrom = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
		input.DateTo = now.Format("2006-01-02")
	}

	total, err := uc.expenseRepo.GetTotalExpensesByDateRange(ctx, input.DateFrom, input.DateTo)
	if err != nil {
		return nil, err
	}

	byCategory, err := uc.expenseRepo.GetExpensesByCategory(ctx, input.DateFrom, input.DateTo)
	if err != nil {
		return nil, err
	}

	// Get count
	expenses, count, err := uc.expenseRepo.ListExpenses(ctx, ports.ExpenseFilter{
		DateFrom: input.DateFrom,
		DateTo:   input.DateTo,
		Limit:    1,
	})
	_ = expenses

	return &ExpenseReport{
		DateFrom:           input.DateFrom,
		DateTo:             input.DateTo,
		TotalExpenses:      total,
		ExpensesByCategory: byCategory,
		ExpenseCount:       count,
	}, err
}

func (uc *ExpenseUseCase) GetMonthlyTrend(ctx context.Context, months int) ([]domain.MonthlyTrend, error) {
	if months <= 0 {
		months = 6
	}

	trends, err := uc.expenseRepo.GetMonthlyExpenseTrend(ctx, months)
	if err != nil {
		return nil, err
	}

	// Enrich with revenue data from invoices
	for i := range trends {
		monthStart := trends[i].Month + "-01"
		// Parse end of month
		t, parseErr := time.Parse("2006-01-02", monthStart)
		if parseErr != nil {
			continue
		}
		monthEnd := t.AddDate(0, 1, -1).Format("2006-01-02")

		rev, revErr := uc.getRevenueByDateRange(ctx, monthStart, monthEnd)
		if revErr == nil {
			trends[i].Revenue = rev
			trends[i].Profit = rev - trends[i].Expenses
		}
	}

	return trends, nil
}

// getRevenueByDateRange calculates revenue from paid invoices.
func (uc *ExpenseUseCase) getRevenueByDateRange(ctx context.Context, dateFrom, dateTo string) (float64, error) {
	// Use invoice list with date filter and sum grand_total for paid/partial invoices
	invoices, _, err := uc.invoiceRepo.List(ctx, ports.InvoiceFilter{
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Limit:    10000,
		Offset:   0,
	})
	if err != nil {
		return 0, err
	}

	var total float64
	for _, inv := range invoices {
		total += inv.GrandTotal
	}
	return total, nil
}
