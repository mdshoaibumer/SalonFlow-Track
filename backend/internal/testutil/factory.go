// Package testutil provides test factories, helpers, and shared infrastructure
// for backend testing.
package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/pkg/uid"
)

// TestDB creates an in-memory SQLite database with all tables.
func TestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:?_foreign_keys=on")
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	if err := CreateAllTables(db); err != nil {
		t.Fatalf("create tables: %v", err)
	}
	return db
}

// TestLogger returns a quiet logger for tests.
func TestLogger(_ *testing.T) *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

// TestContext returns a background context.
func TestContext() context.Context {
	return context.Background()
}

// TestRouter creates a chi router for handler tests.
func TestRouter() chi.Router {
	return chi.NewRouter()
}

// DoRequest sends an HTTP request to a handler and returns the recorder.
func DoRequest(handler http.Handler, method, path string, body string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

// CreateAllTables creates all application tables in the given DB.
func CreateAllTables(db *sql.DB) error {
	for _, ddl := range allDDL {
		if _, err := db.Exec(ddl); err != nil {
			return fmt.Errorf("exec DDL: %w", err)
		}
	}
	return nil
}

// ---------- Factories ----------

// StaffFactory creates a valid Staff using the domain constructor.
func StaffFactory() *domain.Staff {
	s := domain.NewStaff("Priya Sharma", fmt.Sprintf("+91%010d", time.Now().UnixNano()%10000000000), "stylist", 25000, 10)
	return s
}

// ServiceFactory creates a valid Service using the domain constructor.
func ServiceFactory() *domain.Service {
	return domain.NewService("Haircut - Ladies", "hair", 45, 500)
}

// CustomerFactory creates a valid Customer using the domain constructor.
func CustomerFactory() *domain.Customer {
	return domain.NewCustomer("Anjali Desai", fmt.Sprintf("+91%010d", time.Now().UnixNano()%10000000000))
}

// AppointmentFactory creates a valid Appointment.
func AppointmentFactory() *domain.Appointment {
	return domain.NewAppointment(
		uid.New().String(),
		uid.New().String(),
		time.Now().Format("2006-01-02"),
		"10:00", "11:00", false,
	)
}

// MembershipPlanFactory creates a valid MembershipPlan.
func MembershipPlanFactory() *domain.MembershipPlan {
	return domain.NewMembershipPlan("Gold Package", domain.PlanTypePackage, 5000, 90, 12)
}

// WhatsAppTemplateFactory creates a valid WhatsAppTemplate.
func WhatsAppTemplateFactory() *domain.WhatsAppTemplate {
	return domain.NewWhatsAppTemplate("Booking Confirm", domain.WACategoryAppointment, "Hi {{name}}, confirmed for {{date}}")
}

// CloudBackupConfigFactory creates a default cloud backup config.
func CloudBackupConfigFactory() *domain.CloudBackupConfig {
	cfg := domain.NewCloudBackupConfig()
	cfg.Provider = domain.CloudProviderS3
	cfg.BucketName = "test-bucket"
	cfg.Region = "ap-south-1"
	return cfg
}

// ---------- Schema DDL ----------

var allDDL = []string{
	`CREATE TABLE IF NOT EXISTS staff (
		id TEXT PRIMARY KEY, staff_code TEXT NOT NULL DEFAULT '',
		full_name TEXT NOT NULL, phone TEXT NOT NULL UNIQUE, email TEXT NOT NULL DEFAULT '',
		gender TEXT NOT NULL DEFAULT 'male', designation TEXT NOT NULL DEFAULT 'stylist',
		joining_date TEXT NOT NULL DEFAULT '', base_salary REAL NOT NULL DEFAULT 0,
		commission_percentage REAL NOT NULL DEFAULT 0, status TEXT NOT NULL DEFAULT 'active',
		created_at TEXT NOT NULL, updated_at TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS services (
		id TEXT PRIMARY KEY, service_code TEXT NOT NULL DEFAULT '',
		name TEXT NOT NULL, category TEXT NOT NULL DEFAULT 'general',
		description TEXT NOT NULL DEFAULT '', duration_minutes INTEGER NOT NULL DEFAULT 30,
		price REAL NOT NULL DEFAULT 0, cost_price REAL NOT NULL DEFAULT 0,
		commission_type TEXT NOT NULL DEFAULT 'percentage', commission_value REAL NOT NULL DEFAULT 0,
		status TEXT NOT NULL DEFAULT 'active', created_at TEXT NOT NULL, updated_at TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS customers (
		id TEXT PRIMARY KEY, customer_code TEXT NOT NULL DEFAULT '',
		full_name TEXT NOT NULL, phone TEXT NOT NULL UNIQUE, email TEXT NOT NULL DEFAULT '',
		gender TEXT NOT NULL DEFAULT 'female', date_of_birth TEXT, anniversary_date TEXT,
		address TEXT NOT NULL DEFAULT '', notes TEXT NOT NULL DEFAULT '',
		total_visits INTEGER NOT NULL DEFAULT 0, total_spent REAL NOT NULL DEFAULT 0,
		last_visit_date TEXT, status TEXT NOT NULL DEFAULT 'active',
		created_at TEXT NOT NULL, updated_at TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS invoices (
		id TEXT PRIMARY KEY, invoice_number TEXT NOT NULL UNIQUE,
		customer_id TEXT NOT NULL DEFAULT '', customer_name TEXT NOT NULL DEFAULT '',
		staff_id TEXT NOT NULL DEFAULT '', staff_name TEXT NOT NULL DEFAULT '',
		subtotal REAL NOT NULL DEFAULT 0, discount REAL NOT NULL DEFAULT 0,
		tax REAL NOT NULL DEFAULT 0, grand_total REAL NOT NULL DEFAULT 0,
		payment_method TEXT NOT NULL DEFAULT 'cash', status TEXT NOT NULL DEFAULT 'completed',
		notes TEXT NOT NULL DEFAULT '', created_at TEXT NOT NULL, updated_at TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS invoice_items (
		id TEXT PRIMARY KEY, invoice_id TEXT NOT NULL,
		service_id TEXT NOT NULL DEFAULT '', service_name TEXT NOT NULL,
		quantity INTEGER NOT NULL DEFAULT 1, unit_price REAL NOT NULL DEFAULT 0,
		total REAL NOT NULL DEFAULT 0, staff_id TEXT NOT NULL DEFAULT '',
		created_at TEXT NOT NULL, FOREIGN KEY (invoice_id) REFERENCES invoices(id)
	)`,
	`CREATE TABLE IF NOT EXISTS appointments (
		id TEXT PRIMARY KEY, customer_id TEXT NOT NULL, staff_id TEXT NOT NULL,
		appointment_date TEXT NOT NULL DEFAULT '', start_time TEXT NOT NULL, end_time TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'booked', notes TEXT NOT NULL DEFAULT '',
		is_walkin INTEGER NOT NULL DEFAULT 0, total_amount REAL NOT NULL DEFAULT 0,
		created_at TEXT NOT NULL, updated_at TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS appointment_services (
		id TEXT PRIMARY KEY, appointment_id TEXT NOT NULL,
		service_id TEXT NOT NULL, service_name TEXT NOT NULL DEFAULT '',
		duration_minutes INTEGER NOT NULL DEFAULT 0, price REAL NOT NULL DEFAULT 0,
		created_at TEXT NOT NULL, FOREIGN KEY (appointment_id) REFERENCES appointments(id) ON DELETE CASCADE
	)`,
	`CREATE TABLE IF NOT EXISTS appointment_history (
		id TEXT PRIMARY KEY, appointment_id TEXT NOT NULL,
		old_status TEXT NOT NULL DEFAULT '', new_status TEXT NOT NULL DEFAULT '',
		changed_by TEXT NOT NULL DEFAULT '', note TEXT NOT NULL DEFAULT '',
		created_at TEXT NOT NULL, FOREIGN KEY (appointment_id) REFERENCES appointments(id) ON DELETE CASCADE
	)`,
	`CREATE TABLE IF NOT EXISTS expenses (
		id TEXT PRIMARY KEY, category TEXT NOT NULL, amount REAL NOT NULL DEFAULT 0,
		description TEXT NOT NULL DEFAULT '', expense_date TEXT NOT NULL,
		payment_method TEXT NOT NULL DEFAULT 'cash', vendor TEXT NOT NULL DEFAULT '',
		receipt_number TEXT NOT NULL DEFAULT '', is_recurring INTEGER NOT NULL DEFAULT 0,
		created_at TEXT NOT NULL, updated_at TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS products (
		id TEXT PRIMARY KEY, name TEXT NOT NULL, sku TEXT NOT NULL DEFAULT '',
		category TEXT NOT NULL DEFAULT 'retail', brand TEXT NOT NULL DEFAULT '',
		purchase_price REAL NOT NULL DEFAULT 0, selling_price REAL NOT NULL DEFAULT 0,
		current_stock INTEGER NOT NULL DEFAULT 0, min_stock INTEGER NOT NULL DEFAULT 5,
		unit TEXT NOT NULL DEFAULT 'piece', is_active INTEGER NOT NULL DEFAULT 1,
		created_at TEXT NOT NULL, updated_at TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS stock_transactions (
		id TEXT PRIMARY KEY, product_id TEXT NOT NULL, transaction_type TEXT NOT NULL,
		quantity INTEGER NOT NULL, unit_price REAL NOT NULL DEFAULT 0,
		reference TEXT NOT NULL DEFAULT '', notes TEXT NOT NULL DEFAULT '',
		transaction_date TEXT NOT NULL, created_at TEXT NOT NULL,
		FOREIGN KEY (product_id) REFERENCES products(id)
	)`,
	`CREATE TABLE IF NOT EXISTS membership_plans (
		id TEXT PRIMARY KEY, name TEXT NOT NULL, description TEXT NOT NULL DEFAULT '',
		plan_type TEXT NOT NULL DEFAULT 'package', price REAL NOT NULL DEFAULT 0,
		duration_days INTEGER NOT NULL DEFAULT 30, max_sessions INTEGER NOT NULL DEFAULT 0,
		discount_percentage REAL NOT NULL DEFAULT 0, priority_booking INTEGER NOT NULL DEFAULT 0,
		is_active INTEGER NOT NULL DEFAULT 1, created_at TEXT NOT NULL, updated_at TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS package_services (
		id TEXT PRIMARY KEY, plan_id TEXT NOT NULL, service_id TEXT NOT NULL,
		service_name TEXT NOT NULL DEFAULT '', sessions_included INTEGER NOT NULL DEFAULT 1,
		created_at TEXT NOT NULL, FOREIGN KEY (plan_id) REFERENCES membership_plans(id)
	)`,
	`CREATE TABLE IF NOT EXISTS member_subscriptions (
		id TEXT PRIMARY KEY, customer_id TEXT NOT NULL, plan_id TEXT NOT NULL,
		plan_name TEXT NOT NULL DEFAULT '', start_date TEXT NOT NULL, end_date TEXT NOT NULL,
		total_sessions INTEGER NOT NULL DEFAULT 0, used_sessions INTEGER NOT NULL DEFAULT 0,
		amount_paid REAL NOT NULL DEFAULT 0, status TEXT NOT NULL DEFAULT 'active',
		created_at TEXT NOT NULL, updated_at TEXT NOT NULL,
		FOREIGN KEY (plan_id) REFERENCES membership_plans(id)
	)`,
	`CREATE TABLE IF NOT EXISTS whatsapp_templates (
		id TEXT PRIMARY KEY, name TEXT NOT NULL, category TEXT NOT NULL DEFAULT 'general',
		body TEXT NOT NULL, variables TEXT NOT NULL DEFAULT '[]',
		is_active INTEGER NOT NULL DEFAULT 1, created_at TEXT NOT NULL, updated_at TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS whatsapp_messages (
		id TEXT PRIMARY KEY, template_id TEXT NOT NULL,
		recipient_phone TEXT NOT NULL, recipient_name TEXT NOT NULL DEFAULT '',
		message_body TEXT NOT NULL, status TEXT NOT NULL DEFAULT 'queued',
		provider TEXT NOT NULL DEFAULT '', provider_message_id TEXT NOT NULL DEFAULT '',
		error_message TEXT NOT NULL DEFAULT '', sent_at TEXT NOT NULL DEFAULT '',
		delivered_at TEXT NOT NULL DEFAULT '', read_at TEXT NOT NULL DEFAULT '',
		created_at TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS automation_rules (
		id TEXT PRIMARY KEY, name TEXT NOT NULL, trigger_type TEXT NOT NULL,
		template_id TEXT NOT NULL, delay_minutes INTEGER NOT NULL DEFAULT 0,
		is_active INTEGER NOT NULL DEFAULT 1, created_at TEXT NOT NULL, updated_at TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS cloud_backup_config (
		id TEXT PRIMARY KEY, provider TEXT NOT NULL DEFAULT 'none',
		bucket_name TEXT NOT NULL DEFAULT '', region TEXT NOT NULL DEFAULT '',
		access_key TEXT NOT NULL DEFAULT '', endpoint TEXT NOT NULL DEFAULT '',
		encrypt_backups INTEGER NOT NULL DEFAULT 1, auto_backup INTEGER NOT NULL DEFAULT 0,
		auto_backup_interval_hours INTEGER NOT NULL DEFAULT 24, max_versions INTEGER NOT NULL DEFAULT 10,
		created_at TEXT NOT NULL, updated_at TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS cloud_backup_history (
		id TEXT PRIMARY KEY, provider TEXT NOT NULL, file_name TEXT NOT NULL,
		file_size INTEGER NOT NULL DEFAULT 0, remote_path TEXT NOT NULL DEFAULT '',
		status TEXT NOT NULL DEFAULT 'pending', is_encrypted INTEGER NOT NULL DEFAULT 0,
		error_message TEXT NOT NULL DEFAULT '', started_at TEXT NOT NULL DEFAULT '',
		completed_at TEXT NOT NULL DEFAULT '', created_at TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS gst_settings (
		id TEXT PRIMARY KEY, business_name TEXT NOT NULL DEFAULT '', gstin TEXT NOT NULL DEFAULT '',
		state TEXT NOT NULL DEFAULT '', address TEXT NOT NULL DEFAULT '', hsn_code TEXT NOT NULL DEFAULT '',
		cgst_rate REAL NOT NULL DEFAULT 9, sgst_rate REAL NOT NULL DEFAULT 9,
		igst_rate REAL NOT NULL DEFAULT 18, is_gst_enabled INTEGER NOT NULL DEFAULT 0,
		created_at TEXT NOT NULL, updated_at TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS tax_rates (
		id TEXT PRIMARY KEY, name TEXT NOT NULL, hsn_code TEXT NOT NULL DEFAULT '',
		cgst_rate REAL NOT NULL DEFAULT 0, sgst_rate REAL NOT NULL DEFAULT 0,
		igst_rate REAL NOT NULL DEFAULT 0, category TEXT NOT NULL DEFAULT 'service',
		is_active INTEGER NOT NULL DEFAULT 1, created_at TEXT NOT NULL, updated_at TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS printer_settings (
		id TEXT PRIMARY KEY, default_printer TEXT NOT NULL DEFAULT '',
		paper_width TEXT NOT NULL DEFAULT '80mm',
		margin_top INTEGER NOT NULL DEFAULT 5, margin_bottom INTEGER NOT NULL DEFAULT 5,
		margin_left INTEGER NOT NULL DEFAULT 5, margin_right INTEGER NOT NULL DEFAULT 5,
		header_text TEXT NOT NULL DEFAULT '', footer_text TEXT NOT NULL DEFAULT 'Thank you!',
		show_logo INTEGER NOT NULL DEFAULT 0, show_qr INTEGER NOT NULL DEFAULT 0,
		upi_id TEXT NOT NULL DEFAULT '', created_at TEXT NOT NULL, updated_at TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS print_jobs (
		id TEXT PRIMARY KEY, document_type TEXT NOT NULL, document_id TEXT NOT NULL DEFAULT '',
		printer_name TEXT NOT NULL DEFAULT '', paper_width TEXT NOT NULL DEFAULT '80mm',
		status TEXT NOT NULL DEFAULT 'queued', copies INTEGER NOT NULL DEFAULT 1,
		created_at TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS commission_rules (
		id TEXT PRIMARY KEY, staff_id TEXT NOT NULL, service_id TEXT NOT NULL DEFAULT '',
		rule_type TEXT NOT NULL DEFAULT 'percentage', value REAL NOT NULL DEFAULT 0,
		min_revenue REAL NOT NULL DEFAULT 0, is_active INTEGER NOT NULL DEFAULT 1,
		created_at TEXT NOT NULL, updated_at TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS commission_entries (
		id TEXT PRIMARY KEY, staff_id TEXT NOT NULL, invoice_id TEXT NOT NULL DEFAULT '',
		service_name TEXT NOT NULL DEFAULT '', invoice_amount REAL NOT NULL DEFAULT 0,
		commission_amount REAL NOT NULL DEFAULT 0, rule_id TEXT NOT NULL DEFAULT '',
		commission_date TEXT NOT NULL, status TEXT NOT NULL DEFAULT 'pending', created_at TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS salary_records (
		id TEXT PRIMARY KEY, staff_id TEXT NOT NULL, staff_name TEXT NOT NULL DEFAULT '',
		month TEXT NOT NULL, base_salary REAL NOT NULL DEFAULT 0,
		total_commission REAL NOT NULL DEFAULT 0, bonus REAL NOT NULL DEFAULT 0,
		deductions REAL NOT NULL DEFAULT 0, advances REAL NOT NULL DEFAULT 0,
		net_salary REAL NOT NULL DEFAULT 0, status TEXT NOT NULL DEFAULT 'draft',
		paid_date TEXT NOT NULL DEFAULT '', notes TEXT NOT NULL DEFAULT '',
		created_at TEXT NOT NULL, updated_at TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS advances (
		id TEXT PRIMARY KEY, staff_id TEXT NOT NULL, staff_name TEXT NOT NULL DEFAULT '',
		amount REAL NOT NULL DEFAULT 0, advance_date TEXT NOT NULL,
		reason TEXT NOT NULL DEFAULT '', status TEXT NOT NULL DEFAULT 'pending',
		deducted_in_month TEXT NOT NULL DEFAULT '', created_at TEXT NOT NULL, updated_at TEXT NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS daily_performance (
		id TEXT PRIMARY KEY, staff_id TEXT NOT NULL, staff_name TEXT NOT NULL DEFAULT '',
		performance_date TEXT NOT NULL, total_services INTEGER NOT NULL DEFAULT 0,
		total_revenue REAL NOT NULL DEFAULT 0, total_customers INTEGER NOT NULL DEFAULT 0,
		avg_rating REAL NOT NULL DEFAULT 0, created_at TEXT NOT NULL
	)`,
}
