package sqlite

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupAnalyticsTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	// Create minimal schema needed for analytics queries
	_, err = db.Exec(`
		CREATE TABLE staff (
			id TEXT PRIMARY KEY, name TEXT NOT NULL, status TEXT NOT NULL DEFAULT 'active',
			designation TEXT, phone TEXT, email TEXT, date_of_joining TEXT,
			base_salary REAL DEFAULT 0, created_at TEXT, updated_at TEXT
		);
		CREATE TABLE customers (
			id TEXT PRIMARY KEY, name TEXT NOT NULL, phone TEXT, email TEXT,
			date_of_birth TEXT, status TEXT NOT NULL DEFAULT 'active',
			created_at TEXT NOT NULL, updated_at TEXT NOT NULL
		);
		CREATE TABLE invoices (
			id TEXT PRIMARY KEY, invoice_number TEXT, customer_id TEXT,
			total_amount REAL NOT NULL DEFAULT 0, discount_amount REAL DEFAULT 0,
			net_amount REAL DEFAULT 0, payment_status TEXT NOT NULL DEFAULT 'pending',
			payment_method TEXT DEFAULT '', notes TEXT DEFAULT '',
			created_at TEXT NOT NULL, updated_at TEXT NOT NULL
		);
		CREATE TABLE invoice_items (
			id TEXT PRIMARY KEY, invoice_id TEXT NOT NULL, service_id TEXT,
			staff_id TEXT NOT NULL, service_name TEXT NOT NULL,
			quantity INTEGER DEFAULT 1, unit_price REAL NOT NULL, total REAL NOT NULL,
			created_at TEXT NOT NULL, updated_at TEXT NOT NULL
		);
		CREATE TABLE expenses (
			id TEXT PRIMARY KEY, expense_number TEXT, category_id TEXT NOT NULL,
			amount REAL NOT NULL, expense_date DATE NOT NULL,
			payment_method TEXT, vendor_name TEXT, invoice_reference TEXT,
			description TEXT, attachment_path TEXT, status TEXT NOT NULL DEFAULT 'approved',
			created_by TEXT, created_at TEXT NOT NULL, updated_at TEXT NOT NULL
		);
		CREATE TABLE expense_categories (
			id TEXT PRIMARY KEY, name TEXT NOT NULL, description TEXT,
			is_active INTEGER DEFAULT 1, created_at TEXT, updated_at TEXT
		);
		CREATE TABLE salary_cycles (
			id TEXT PRIMARY KEY, month INTEGER NOT NULL, year INTEGER NOT NULL,
			start_date TEXT NOT NULL, end_date TEXT NOT NULL, status TEXT NOT NULL,
			generated_at TEXT, generated_by TEXT, created_at TEXT, updated_at TEXT
		);
		CREATE TABLE salary_records (
			id TEXT PRIMARY KEY, cycle_id TEXT NOT NULL, staff_id TEXT NOT NULL,
			base_salary REAL DEFAULT 0, commission_amount REAL DEFAULT 0,
			advance_deduction REAL DEFAULT 0, net_salary REAL NOT NULL,
			payment_status TEXT NOT NULL DEFAULT 'pending', payment_date TEXT,
			created_at TEXT, updated_at TEXT
		);
		CREATE TABLE advances (
			id TEXT PRIMARY KEY, staff_id TEXT NOT NULL, amount REAL NOT NULL,
			repaid_amount REAL NOT NULL DEFAULT 0, status TEXT NOT NULL DEFAULT 'approved',
			purpose TEXT, approved_date TEXT, created_at TEXT, updated_at TEXT
		);
		CREATE TABLE products (
			id TEXT PRIMARY KEY, product_code TEXT, name TEXT NOT NULL,
			category TEXT NOT NULL, brand TEXT, unit TEXT, sku TEXT,
			purchase_price REAL DEFAULT 0, selling_price REAL DEFAULT 0,
			current_stock REAL DEFAULT 0, minimum_stock REAL DEFAULT 5,
			maximum_stock REAL DEFAULT 50, status TEXT NOT NULL DEFAULT 'active',
			created_at TEXT, updated_at TEXT
		);
		CREATE TABLE stock_transactions (
			id TEXT PRIMARY KEY, product_id TEXT NOT NULL,
			transaction_type TEXT NOT NULL, quantity REAL NOT NULL,
			unit_cost REAL DEFAULT 0, reference_type TEXT, reference_id TEXT,
			notes TEXT, transaction_date DATE NOT NULL, created_at TEXT, updated_at TEXT
		);
		CREATE TABLE purchase_entries (
			id TEXT PRIMARY KEY, purchase_number TEXT, vendor_name TEXT NOT NULL,
			invoice_number TEXT, purchase_date DATE NOT NULL,
			total_amount REAL NOT NULL DEFAULT 0, notes TEXT,
			created_at TEXT, updated_at TEXT
		);
		CREATE TABLE commission_transactions (
			id TEXT PRIMARY KEY, staff_id TEXT NOT NULL, invoice_id TEXT,
			rule_id TEXT, amount REAL NOT NULL, period TEXT,
			created_at TEXT NOT NULL, updated_at TEXT NOT NULL
		);
	`)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func seedAnalyticsData(t *testing.T, db *sql.DB) {
	t.Helper()
	// Insert staff
	db.Exec(`INSERT INTO staff (id, name, status, created_at, updated_at) VALUES
		('staff-1', 'Priya', 'active', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
		('staff-2', 'Rahul', 'active', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z')`)

	// Insert customers
	db.Exec(`INSERT INTO customers (id, name, phone, status, date_of_birth, created_at, updated_at) VALUES
		('cust-1', 'Meena', '9000000001', 'active', '1990-06-09', '2026-05-01T00:00:00Z', '2026-05-01T00:00:00Z'),
		('cust-2', 'Sunita', '9000000002', 'active', '1985-03-15', '2026-06-01T00:00:00Z', '2026-06-01T00:00:00Z'),
		('cust-3', 'Ravi', '9000000003', 'inactive', '1992-11-20', '2026-04-01T00:00:00Z', '2026-04-01T00:00:00Z')`)

	// Insert invoices (today = 2026-06-09 in test context, we use this month)
	db.Exec(`INSERT INTO invoices (id, invoice_number, customer_id, total_amount, payment_status, created_at, updated_at) VALUES
		('inv-1', 'INV-001', 'cust-1', 1500, 'paid', '2026-06-09T10:00:00Z', '2026-06-09T10:00:00Z'),
		('inv-2', 'INV-002', 'cust-2', 2000, 'paid', '2026-06-09T11:00:00Z', '2026-06-09T11:00:00Z'),
		('inv-3', 'INV-003', 'cust-1', 800, 'paid', '2026-06-05T10:00:00Z', '2026-06-05T10:00:00Z'),
		('inv-4', 'INV-004', 'cust-2', 1200, 'pending', '2026-06-08T10:00:00Z', '2026-06-08T10:00:00Z')`)

	// Invoice items
	db.Exec(`INSERT INTO invoice_items (id, invoice_id, staff_id, service_name, unit_price, total, created_at, updated_at) VALUES
		('ii-1', 'inv-1', 'staff-1', 'Haircut', 500, 500, '2026-06-09T10:00:00Z', '2026-06-09T10:00:00Z'),
		('ii-2', 'inv-1', 'staff-1', 'Color', 1000, 1000, '2026-06-09T10:00:00Z', '2026-06-09T10:00:00Z'),
		('ii-3', 'inv-2', 'staff-2', 'Facial', 2000, 2000, '2026-06-09T11:00:00Z', '2026-06-09T11:00:00Z'),
		('ii-4', 'inv-3', 'staff-1', 'Haircut', 800, 800, '2026-06-05T10:00:00Z', '2026-06-05T10:00:00Z')`)

	// Expenses
	db.Exec(`INSERT INTO expense_categories (id, name, created_at, updated_at) VALUES
		('cat-1', 'Rent', '2026-01-01', '2026-01-01'),
		('cat-2', 'Utilities', '2026-01-01', '2026-01-01')`)
	db.Exec(`INSERT INTO expenses (id, category_id, amount, expense_date, status, created_at, updated_at) VALUES
		('exp-1', 'cat-1', 15000, '2026-06-01', 'approved', '2026-06-01T00:00:00Z', '2026-06-01T00:00:00Z'),
		('exp-2', 'cat-2', 3000, '2026-06-05', 'approved', '2026-06-05T00:00:00Z', '2026-06-05T00:00:00Z')`)

	// Products
	db.Exec(`INSERT INTO products (id, product_code, name, category, purchase_price, current_stock, minimum_stock, status, created_at, updated_at) VALUES
		('prod-1', 'PRD-001', 'Shampoo', 'hair_care', 200, 3, 10, 'active', '2026-01-01', '2026-01-01'),
		('prod-2', 'PRD-002', 'Conditioner', 'hair_care', 150, 20, 5, 'active', '2026-01-01', '2026-01-01')`)

	// Salary
	db.Exec(`INSERT INTO salary_cycles (id, month, year, start_date, end_date, status, created_at, updated_at) VALUES
		('cyc-1', 6, 2026, '2026-06-01', '2026-06-30', 'generated', '2026-06-01', '2026-06-01')`)
	db.Exec(`INSERT INTO salary_records (id, cycle_id, staff_id, net_salary, payment_status, created_at, updated_at) VALUES
		('sr-1', 'cyc-1', 'staff-1', 25000, 'pending', '2026-06-01', '2026-06-01'),
		('sr-2', 'cyc-1', 'staff-2', 20000, 'paid', '2026-06-01', '2026-06-01')`)

	// Advances
	db.Exec(`INSERT INTO advances (id, staff_id, amount, repaid_amount, status, created_at, updated_at) VALUES
		('adv-1', 'staff-1', 5000, 2000, 'partially_repaid', '2026-05-01', '2026-06-01')`)
}

func TestAnalyticsRepository_GetDashboardStats(t *testing.T) {
	db := setupAnalyticsTestDB(t)
	defer db.Close()
	seedAnalyticsData(t, db)

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	repo := NewAnalyticsRepository(db, log)
	ctx := context.Background()

	stats, err := repo.GetDashboardStats(ctx)
	if err != nil {
		t.Fatalf("GetDashboardStats: %v", err)
	}

	// Monthly revenue = inv-1 (1500) + inv-2 (2000) + inv-3 (800) = 4300
	if stats.MonthlyRevenue != 4300 {
		t.Errorf("expected monthly revenue 4300, got %f", stats.MonthlyRevenue)
	}
	// Monthly expenses = 15000 + 3000 = 18000
	if stats.MonthlyExpenses != 18000 {
		t.Errorf("expected monthly expenses 18000, got %f", stats.MonthlyExpenses)
	}
	// Inventory value = 3*200 + 20*150 = 600 + 3000 = 3600
	if stats.InventoryValue != 3600 {
		t.Errorf("expected inventory value 3600, got %f", stats.InventoryValue)
	}
	// Outstanding salary = 25000 (pending)
	if stats.OutstandingSalary != 25000 {
		t.Errorf("expected outstanding salary 25000, got %f", stats.OutstandingSalary)
	}
	// Outstanding advances = 5000-2000 = 3000
	if stats.OutstandingAdvance != 3000 {
		t.Errorf("expected outstanding advances 3000, got %f", stats.OutstandingAdvance)
	}
	// Low stock = prod-1 (3 < 10)
	if stats.LowStockCount != 1 {
		t.Errorf("expected low stock 1, got %d", stats.LowStockCount)
	}
}

func TestAnalyticsRepository_GetKPIMetrics(t *testing.T) {
	db := setupAnalyticsTestDB(t)
	defer db.Close()
	seedAnalyticsData(t, db)

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	repo := NewAnalyticsRepository(db, log)
	ctx := context.Background()

	kpis, err := repo.GetKPIMetrics(ctx, "2026-06-01", "2026-06-30")
	if err != nil {
		t.Fatalf("GetKPIMetrics: %v", err)
	}
	// Average bill value = 4300/3 = 1433.33
	if kpis.AverageBillValue < 1433 || kpis.AverageBillValue > 1434 {
		t.Errorf("expected avg bill ~1433, got %f", kpis.AverageBillValue)
	}
	// Staff productivity: 2 active staff, both produced revenue
	if kpis.StaffProductivity != 100 {
		t.Errorf("expected staff productivity 100%%, got %f", kpis.StaffProductivity)
	}
}

func TestAnalyticsRepository_GetRevenueReport(t *testing.T) {
	db := setupAnalyticsTestDB(t)
	defer db.Close()
	seedAnalyticsData(t, db)

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	repo := NewAnalyticsRepository(db, log)
	ctx := context.Background()

	report, err := repo.GetRevenueReport(ctx, "2026-06-01", "2026-06-30", "day")
	if err != nil {
		t.Fatalf("GetRevenueReport: %v", err)
	}
	if report.TotalRevenue != 4300 {
		t.Errorf("expected total revenue 4300, got %f", report.TotalRevenue)
	}
	if report.InvoiceCount != 3 {
		t.Errorf("expected 3 invoices, got %d", report.InvoiceCount)
	}
	if len(report.ByService) == 0 {
		t.Error("expected by_service data")
	}
	if len(report.ByStaff) == 0 {
		t.Error("expected by_staff data")
	}
}

func TestAnalyticsRepository_GetCustomerReport(t *testing.T) {
	db := setupAnalyticsTestDB(t)
	defer db.Close()
	seedAnalyticsData(t, db)

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	repo := NewAnalyticsRepository(db, log)
	ctx := context.Background()

	report, err := repo.GetCustomerReport(ctx, "2026-06-01", "2026-06-30")
	if err != nil {
		t.Fatalf("GetCustomerReport: %v", err)
	}
	if report.TotalCustomers != 3 {
		t.Errorf("expected 3 total customers, got %d", report.TotalCustomers)
	}
	if report.InactiveCount != 1 {
		t.Errorf("expected 1 inactive, got %d", report.InactiveCount)
	}
}

func TestAnalyticsRepository_GetExpenseReport(t *testing.T) {
	db := setupAnalyticsTestDB(t)
	defer db.Close()
	seedAnalyticsData(t, db)

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	repo := NewAnalyticsRepository(db, log)
	ctx := context.Background()

	report, err := repo.GetExpenseReport(ctx, "2026-06-01", "2026-06-30")
	if err != nil {
		t.Fatalf("GetExpenseReport: %v", err)
	}
	if report.TotalExpenses != 18000 {
		t.Errorf("expected 18000, got %f", report.TotalExpenses)
	}
	if len(report.ByCategory) != 2 {
		t.Errorf("expected 2 categories, got %d", len(report.ByCategory))
	}
}

func TestAnalyticsRepository_GetInventoryReport(t *testing.T) {
	db := setupAnalyticsTestDB(t)
	defer db.Close()
	seedAnalyticsData(t, db)

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	repo := NewAnalyticsRepository(db, log)
	ctx := context.Background()

	report, err := repo.GetInventoryReport(ctx)
	if err != nil {
		t.Fatalf("GetInventoryReport: %v", err)
	}
	// 3*200 + 20*150 = 3600
	if report.TotalValue != 3600 {
		t.Errorf("expected value 3600, got %f", report.TotalValue)
	}
	if report.LowStockCount != 1 {
		t.Errorf("expected 1 low stock, got %d", report.LowStockCount)
	}
}

func TestAnalyticsRepository_GetProfitLossReport(t *testing.T) {
	db := setupAnalyticsTestDB(t)
	defer db.Close()
	seedAnalyticsData(t, db)

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	repo := NewAnalyticsRepository(db, log)
	ctx := context.Background()

	report, err := repo.GetProfitLossReport(ctx, "2026-06-01", "2026-06-30", "month")
	if err != nil {
		t.Fatalf("GetProfitLossReport: %v", err)
	}
	if report.Revenue != 4300 {
		t.Errorf("expected revenue 4300, got %f", report.Revenue)
	}
	if report.Expenses != 18000 {
		t.Errorf("expected expenses 18000, got %f", report.Expenses)
	}
	// Net profit = 4300 - 18000 - salary_cost
	// Salary: cycle dates 2026-06-01 to 2026-06-30, both records sum = 45000
	expectedProfit := 4300.0 - 18000.0 - 45000.0
	if report.NetProfit != expectedProfit {
		t.Errorf("expected net profit %f, got %f", expectedProfit, report.NetProfit)
	}
}
