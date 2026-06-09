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

func setupCommissionTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`
		CREATE TABLE staff (
			id TEXT PRIMARY KEY,
			staff_code TEXT NOT NULL,
			full_name TEXT NOT NULL,
			phone TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'active',
			created_at TEXT NOT NULL DEFAULT '',
			updated_at TEXT NOT NULL DEFAULT ''
		);
		CREATE TABLE invoices (
			id TEXT PRIMARY KEY
		);
		CREATE TABLE commission_rules (
			id TEXT PRIMARY KEY,
			rule_name TEXT NOT NULL,
			rule_type TEXT NOT NULL,
			target_type TEXT NOT NULL,
			target_id TEXT,
			calculation_type TEXT NOT NULL,
			calculation_value REAL NOT NULL DEFAULT 0,
			minimum_target REAL NOT NULL DEFAULT 0,
			maximum_target REAL NOT NULL DEFAULT 0,
			is_active INTEGER NOT NULL DEFAULT 1,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE TABLE commission_transactions (
			id TEXT PRIMARY KEY,
			staff_id TEXT NOT NULL,
			invoice_id TEXT NOT NULL,
			rule_id TEXT,
			revenue_amount REAL NOT NULL DEFAULT 0,
			commission_amount REAL NOT NULL DEFAULT 0,
			business_date TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
	`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestCommissionRepository_CreateRule(t *testing.T) {
	db := setupCommissionTestDB(t)
	defer db.Close()
	repo := NewCommissionRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	rule := domain.NewCommissionRule("Revenue 10%", domain.RuleTypeRevenueBased, domain.TargetTypeGlobal, "", domain.CalcTypePercentage, 10, 0, 0)
	err := repo.CreateRule(context.Background(), rule)
	if err != nil {
		t.Fatalf("create rule failed: %v", err)
	}

	found, err := repo.GetRuleByID(context.Background(), rule.ID)
	if err != nil {
		t.Fatalf("get rule failed: %v", err)
	}
	if found.RuleName != "Revenue 10%" {
		t.Errorf("expected 'Revenue 10%%', got %q", found.RuleName)
	}
	if found.CalculationValue != 10 {
		t.Errorf("expected calc value 10, got %.2f", found.CalculationValue)
	}
	if !found.IsActive {
		t.Error("expected rule to be active")
	}
}

func TestCommissionRepository_ListRules(t *testing.T) {
	db := setupCommissionTestDB(t)
	defer db.Close()
	repo := NewCommissionRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	rule1 := domain.NewCommissionRule("Revenue 10%", domain.RuleTypeRevenueBased, domain.TargetTypeGlobal, "", domain.CalcTypePercentage, 10, 0, 0)
	rule2 := domain.NewCommissionRule("Fixed ₹100", domain.RuleTypeFixed, domain.TargetTypeGlobal, "", domain.CalcTypeFixedAmount, 100, 0, 0)
	_ = repo.CreateRule(context.Background(), rule1)
	_ = repo.CreateRule(context.Background(), rule2)

	rules, total, err := repo.ListRules(context.Background(), ports.CommissionRuleFilter{})
	if err != nil {
		t.Fatalf("list rules failed: %v", err)
	}
	if total != 2 {
		t.Fatalf("expected 2 rules, got %d", total)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rule records, got %d", len(rules))
	}
}

func TestCommissionRepository_UpdateRule(t *testing.T) {
	db := setupCommissionTestDB(t)
	defer db.Close()
	repo := NewCommissionRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	rule := domain.NewCommissionRule("Old Name", domain.RuleTypeRevenueBased, domain.TargetTypeGlobal, "", domain.CalcTypePercentage, 5, 0, 0)
	_ = repo.CreateRule(context.Background(), rule)

	rule.RuleName = "New Name"
	rule.CalculationValue = 15
	err := repo.UpdateRule(context.Background(), rule)
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}

	found, _ := repo.GetRuleByID(context.Background(), rule.ID)
	if found.RuleName != "New Name" {
		t.Errorf("expected 'New Name', got %q", found.RuleName)
	}
	if found.CalculationValue != 15 {
		t.Errorf("expected 15, got %.2f", found.CalculationValue)
	}
}

func TestCommissionRepository_DeleteRule(t *testing.T) {
	db := setupCommissionTestDB(t)
	defer db.Close()
	repo := NewCommissionRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	rule := domain.NewCommissionRule("To Delete", domain.RuleTypeFixed, domain.TargetTypeGlobal, "", domain.CalcTypeFixedAmount, 50, 0, 0)
	_ = repo.CreateRule(context.Background(), rule)

	err := repo.DeleteRule(context.Background(), rule.ID)
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	_, err = repo.GetRuleByID(context.Background(), rule.ID)
	if err == nil {
		t.Fatal("expected not found error after delete")
	}
}

func TestCommissionRepository_CreateTransaction(t *testing.T) {
	db := setupCommissionTestDB(t)
	defer db.Close()
	repo := NewCommissionRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	staffID := uuid.New()
	invoiceID := uuid.New()
	ruleID := uuid.New()

	// Insert stub staff/invoice
	_, _ = db.Exec("INSERT INTO staff (id, staff_code, full_name, phone) VALUES (?, 'STF-1', 'Nazim', '9999')", staffID.String())
	_, _ = db.Exec("INSERT INTO invoices (id) VALUES (?)", invoiceID.String())

	tx := domain.NewCommissionTransaction(staffID, invoiceID, ruleID, 5000, 500, "2026-06-09")
	err := repo.CreateTransaction(context.Background(), tx)
	if err != nil {
		t.Fatalf("create transaction failed: %v", err)
	}

	total, err := repo.GetStaffCommission(context.Background(), staffID, "2026-06-01", "2026-06-30")
	if err != nil {
		t.Fatalf("get staff commission failed: %v", err)
	}
	if total != 500 {
		t.Errorf("expected 500, got %.2f", total)
	}
}

func TestCommissionRepository_GetActiveRules(t *testing.T) {
	db := setupCommissionTestDB(t)
	defer db.Close()
	repo := NewCommissionRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	active := domain.NewCommissionRule("Active", domain.RuleTypeRevenueBased, domain.TargetTypeGlobal, "", domain.CalcTypePercentage, 10, 0, 0)
	inactive := domain.NewCommissionRule("Inactive", domain.RuleTypeFixed, domain.TargetTypeGlobal, "", domain.CalcTypeFixedAmount, 50, 0, 0)
	inactive.IsActive = false

	_ = repo.CreateRule(context.Background(), active)
	_ = repo.CreateRule(context.Background(), inactive)

	rules, err := repo.GetActiveRules(context.Background())
	if err != nil {
		t.Fatalf("get active rules failed: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 active rule, got %d", len(rules))
	}
	if rules[0].RuleName != "Active" {
		t.Errorf("expected 'Active', got %q", rules[0].RuleName)
	}
}

func TestCommissionRepository_GetMonthlyCommission(t *testing.T) {
	db := setupCommissionTestDB(t)
	defer db.Close()
	repo := NewCommissionRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	staffID := uuid.New()
	_, _ = db.Exec("INSERT INTO staff (id, staff_code, full_name, phone) VALUES (?, 'STF-1', 'Nazim', '9999')", staffID.String())
	_, _ = db.Exec("INSERT INTO invoices (id) VALUES (?)", uuid.New().String())

	// Create multiple transactions
	for i := 0; i < 3; i++ {
		tx := domain.NewCommissionTransaction(staffID, uuid.New(), uuid.New(), 5000, 500, "2026-06-09")
		_ = repo.CreateTransaction(context.Background(), tx)
	}

	summaries, err := repo.GetMonthlyCommission(context.Background(), "2026-06")
	if err != nil {
		t.Fatalf("get monthly commission failed: %v", err)
	}
	if len(summaries) != 1 {
		t.Fatalf("expected 1 summary, got %d", len(summaries))
	}
	if summaries[0].Commission != 1500 {
		t.Errorf("expected commission 1500, got %.2f", summaries[0].Commission)
	}
}
