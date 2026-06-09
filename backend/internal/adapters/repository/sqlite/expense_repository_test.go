package sqlite

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
)

func setupExpenseTestDB(t *testing.T) (*sql.DB, uuid.UUID, uuid.UUID) {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	catRentID := uuid.New()
	catElecID := uuid.New()

	_, err = db.Exec(`
		CREATE TABLE expense_categories (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			description TEXT NOT NULL DEFAULT '',
			is_active INTEGER NOT NULL DEFAULT 1,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE TABLE expenses (
			id TEXT PRIMARY KEY,
			expense_number TEXT NOT NULL UNIQUE,
			category_id TEXT NOT NULL REFERENCES expense_categories(id),
			amount REAL NOT NULL,
			expense_date DATE NOT NULL,
			payment_method TEXT NOT NULL,
			vendor_name TEXT NOT NULL DEFAULT '',
			invoice_reference TEXT NOT NULL DEFAULT '',
			description TEXT NOT NULL DEFAULT '',
			attachment_path TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL DEFAULT 'pending',
			created_by TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE TABLE expense_number_seq (
			year INTEGER NOT NULL,
			seq INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY (year)
		);
	`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO expense_categories (id, name, description, is_active, created_at, updated_at)
		VALUES (?, 'Rent', 'Monthly rent', 1, '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z')`, catRentID.String())
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`INSERT INTO expense_categories (id, name, description, is_active, created_at, updated_at)
		VALUES (?, 'Electricity', 'Power bills', 1, '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z')`, catElecID.String())
	if err != nil {
		t.Fatal(err)
	}

	return db, catRentID, catElecID
}

func TestExpenseRepository_CreateAndGet(t *testing.T) {
	db, catRentID, _ := setupExpenseTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	repo := NewExpenseRepository(db, log)
	ctx := context.Background()

	exp := domain.NewExpense(catRentID, 15000, "2026-06-01", "cash", "Landlord", "RENT-JUN-2026", "Monthly shop rent", "")
	exp.ExpenseNumber = "EXP-2026-000001"

	err := repo.CreateExpense(ctx, exp)
	if err != nil {
		t.Fatalf("CreateExpense: %v", err)
	}

	got, err := repo.GetExpenseByID(ctx, exp.ID)
	if err != nil {
		t.Fatalf("GetExpenseByID: %v", err)
	}
	if got.ExpenseNumber != "EXP-2026-000001" {
		t.Errorf("expected expense number EXP-2026-000001, got %s", got.ExpenseNumber)
	}
	if got.Amount != 15000 {
		t.Errorf("expected amount 15000, got %f", got.Amount)
	}
	if got.CategoryName != "Rent" {
		t.Errorf("expected category name Rent, got %s", got.CategoryName)
	}
	if got.VendorName != "Landlord" {
		t.Errorf("expected vendor Landlord, got %s", got.VendorName)
	}
}

func TestExpenseRepository_ListWithFilters(t *testing.T) {
	db, catRentID, catElecID := setupExpenseTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	repo := NewExpenseRepository(db, log)
	ctx := context.Background()

	exp1 := domain.NewExpense(catRentID, 15000, "2026-06-01", "cash", "Landlord", "", "Rent", "")
	exp1.ExpenseNumber = "EXP-2026-000001"
	exp2 := domain.NewExpense(catElecID, 3000, "2026-06-05", "upi", "Power Co", "", "Electricity", "")
	exp2.ExpenseNumber = "EXP-2026-000002"
	exp3 := domain.NewExpense(catRentID, 15000, "2026-07-01", "cash", "Landlord", "", "Rent Jul", "")
	exp3.ExpenseNumber = "EXP-2026-000003"

	repo.CreateExpense(ctx, exp1)
	repo.CreateExpense(ctx, exp2)
	repo.CreateExpense(ctx, exp3)

	// Filter by category
	results, total, err := repo.ListExpenses(ctx, ports.ExpenseFilter{CategoryID: catRentID.String(), Limit: 10})
	if err != nil {
		t.Fatalf("ListExpenses by category: %v", err)
	}
	if total != 2 {
		t.Errorf("expected 2 rent expenses, got %d", total)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}

	// Filter by date range
	results, total, err = repo.ListExpenses(ctx, ports.ExpenseFilter{DateFrom: "2026-06-01", DateTo: "2026-06-30", Limit: 10})
	if err != nil {
		t.Fatalf("ListExpenses by date: %v", err)
	}
	if total != 2 {
		t.Errorf("expected 2 June expenses, got %d", total)
	}
	_ = results

	// Filter by payment method
	_, total, err = repo.ListExpenses(ctx, ports.ExpenseFilter{PaymentMethod: "upi", Limit: 10})
	if err != nil {
		t.Fatalf("ListExpenses by payment: %v", err)
	}
	if total != 1 {
		t.Errorf("expected 1 UPI expense, got %d", total)
	}

	// Search
	_, total, err = repo.ListExpenses(ctx, ports.ExpenseFilter{Search: "Landlord", Limit: 10})
	if err != nil {
		t.Fatalf("ListExpenses search: %v", err)
	}
	if total != 2 {
		t.Errorf("expected 2 Landlord expenses, got %d", total)
	}
}

func TestExpenseRepository_NextExpenseNumber(t *testing.T) {
	db, _, _ := setupExpenseTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	repo := NewExpenseRepository(db, log)
	ctx := context.Background()

	num1, err := repo.NextExpenseNumber(ctx, 2026)
	if err != nil {
		t.Fatalf("NextExpenseNumber: %v", err)
	}
	if num1 != "EXP-2026-000001" {
		t.Errorf("expected EXP-2026-000001, got %s", num1)
	}

	num2, err := repo.NextExpenseNumber(ctx, 2026)
	if err != nil {
		t.Fatalf("NextExpenseNumber 2: %v", err)
	}
	if num2 != "EXP-2026-000002" {
		t.Errorf("expected EXP-2026-000002, got %s", num2)
	}

	// Different year resets
	num3, err := repo.NextExpenseNumber(ctx, 2027)
	if err != nil {
		t.Fatalf("NextExpenseNumber 2027: %v", err)
	}
	if num3 != "EXP-2027-000001" {
		t.Errorf("expected EXP-2027-000001, got %s", num3)
	}
}

func TestExpenseRepository_Delete(t *testing.T) {
	db, catRentID, _ := setupExpenseTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	repo := NewExpenseRepository(db, log)
	ctx := context.Background()

	exp := domain.NewExpense(catRentID, 5000, "2026-06-10", "card", "Shop", "", "Supplies", "")
	exp.ExpenseNumber = "EXP-2026-000010"
	repo.CreateExpense(ctx, exp)

	err := repo.DeleteExpense(ctx, exp.ID)
	if err != nil {
		t.Fatalf("DeleteExpense: %v", err)
	}

	_, err = repo.GetExpenseByID(ctx, exp.ID)
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
}

func TestExpenseRepository_TotalByDateRange(t *testing.T) {
	db, catRentID, _ := setupExpenseTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	repo := NewExpenseRepository(db, log)
	ctx := context.Background()

	exp1 := domain.NewExpense(catRentID, 10000, "2026-06-01", "cash", "V1", "", "E1", "")
	exp1.ExpenseNumber = "EXP-2026-000001"
	exp2 := domain.NewExpense(catRentID, 5000, "2026-06-15", "upi", "V2", "", "E2", "")
	exp2.ExpenseNumber = "EXP-2026-000002"
	exp3 := domain.NewExpense(catRentID, 3000, "2026-06-20", "card", "V3", "", "E3", "")
	exp3.ExpenseNumber = "EXP-2026-000003"
	exp3.Status = "rejected" // rejected should be excluded

	repo.CreateExpense(ctx, exp1)
	repo.CreateExpense(ctx, exp2)

	// Manually insert rejected one
	db.Exec(`INSERT INTO expenses (id, expense_number, category_id, amount, expense_date, payment_method, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, 'rejected', '2026-06-20T00:00:00Z', '2026-06-20T00:00:00Z')`,
		exp3.ID.String(), exp3.ExpenseNumber, catRentID.String(), 3000, "2026-06-20", "card")

	total, err := repo.GetTotalExpensesByDateRange(ctx, "2026-06-01", "2026-06-30")
	if err != nil {
		t.Fatalf("GetTotalExpensesByDateRange: %v", err)
	}
	if total != 15000 {
		t.Errorf("expected total 15000 (excluding rejected), got %f", total)
	}
}

func TestExpenseRepository_CategoryCRUD(t *testing.T) {
	db, _, _ := setupExpenseTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	repo := NewExpenseRepository(db, log)
	ctx := context.Background()

	// Create
	cat := domain.NewExpenseCategory("Marketing", "Ad spend")
	err := repo.CreateCategory(ctx, cat)
	if err != nil {
		t.Fatalf("CreateCategory: %v", err)
	}

	// Get
	got, err := repo.GetCategoryByID(ctx, cat.ID)
	if err != nil {
		t.Fatalf("GetCategoryByID: %v", err)
	}
	if got.Name != "Marketing" {
		t.Errorf("expected Marketing, got %s", got.Name)
	}

	// List (all)
	cats, err := repo.ListCategories(ctx, false)
	if err != nil {
		t.Fatalf("ListCategories: %v", err)
	}
	if len(cats) != 3 { // 2 seeded + 1 new
		t.Errorf("expected 3 categories, got %d", len(cats))
	}

	// Update
	got.IsActive = false
	err = repo.UpdateCategory(ctx, got)
	if err != nil {
		t.Fatalf("UpdateCategory: %v", err)
	}

	// List active only
	activeCats, err := repo.ListCategories(ctx, true)
	if err != nil {
		t.Fatalf("ListCategories active: %v", err)
	}
	if len(activeCats) != 2 {
		t.Errorf("expected 2 active categories, got %d", len(activeCats))
	}
}
