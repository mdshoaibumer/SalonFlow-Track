package sqlite

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
)

func setupCustomerTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`
		CREATE TABLE customers (
			id TEXT PRIMARY KEY,
			customer_code TEXT NOT NULL UNIQUE,
			full_name TEXT NOT NULL,
			phone TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL DEFAULT '',
			gender TEXT NOT NULL DEFAULT 'other',
			date_of_birth TEXT,
			anniversary_date TEXT,
			address TEXT NOT NULL DEFAULT '',
			notes TEXT NOT NULL DEFAULT '',
			total_visits INTEGER NOT NULL DEFAULT 0,
			total_spent REAL NOT NULL DEFAULT 0,
			last_visit_date TEXT,
			status TEXT NOT NULL DEFAULT 'active',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)
	`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestCustomerRepository_Create(t *testing.T) {
	db := setupCustomerTestDB(t)
	defer db.Close()
	repo := NewCustomerRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	c := domain.NewCustomer("John Doe", "9876543210")
	err := repo.Create(context.Background(), c)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	found, err := repo.GetByID(context.Background(), c.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.FullName != "John Doe" {
		t.Errorf("expected 'John Doe', got %q", found.FullName)
	}
}

func TestCustomerRepository_GetByPhone(t *testing.T) {
	db := setupCustomerTestDB(t)
	defer db.Close()
	repo := NewCustomerRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	c := domain.NewCustomer("Jane Doe", "9876543211")
	_ = repo.Create(context.Background(), c)

	found, err := repo.GetByPhone(context.Background(), "9876543211")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.FullName != "Jane Doe" {
		t.Errorf("expected 'Jane Doe', got %q", found.FullName)
	}
}

func TestCustomerRepository_DuplicatePhone(t *testing.T) {
	db := setupCustomerTestDB(t)
	defer db.Close()
	repo := NewCustomerRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	c1 := domain.NewCustomer("John", "9876543210")
	_ = repo.Create(context.Background(), c1)

	c2 := domain.NewCustomer("Jane", "9876543210")
	err := repo.Create(context.Background(), c2)
	if err == nil {
		t.Error("expected duplicate phone error")
	}
}

func TestCustomerRepository_List(t *testing.T) {
	db := setupCustomerTestDB(t)
	defer db.Close()
	repo := NewCustomerRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	_ = repo.Create(context.Background(), domain.NewCustomer("Alice", "9876543210"))
	_ = repo.Create(context.Background(), domain.NewCustomer("Bob", "9876543211"))
	_ = repo.Create(context.Background(), domain.NewCustomer("Charlie", "9876543212"))

	customers, total, err := repo.List(context.Background(), ports.CustomerFilter{Limit: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 3 {
		t.Errorf("expected 3 total, got %d", total)
	}
	if len(customers) != 3 {
		t.Errorf("expected 3 customers, got %d", len(customers))
	}
}

func TestCustomerRepository_Search(t *testing.T) {
	db := setupCustomerTestDB(t)
	defer db.Close()
	repo := NewCustomerRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	_ = repo.Create(context.Background(), domain.NewCustomer("Alice Smith", "9876543210"))
	_ = repo.Create(context.Background(), domain.NewCustomer("Bob Jones", "9876543211"))

	customers, total, err := repo.List(context.Background(), ports.CustomerFilter{Search: "alice", Limit: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 1 {
		t.Errorf("expected 1 result, got %d", total)
	}
	if len(customers) != 1 {
		t.Errorf("expected 1 customer, got %d", len(customers))
	}
}

func TestCustomerRepository_Update(t *testing.T) {
	db := setupCustomerTestDB(t)
	defer db.Close()
	repo := NewCustomerRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	c := domain.NewCustomer("John", "9876543210")
	_ = repo.Create(context.Background(), c)

	c.FullName = "John Updated"
	c.Email = "john@test.com"
	err := repo.Update(context.Background(), c)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	found, _ := repo.GetByID(context.Background(), c.ID)
	if found.FullName != "John Updated" {
		t.Errorf("expected 'John Updated', got %q", found.FullName)
	}
	if found.Email != "john@test.com" {
		t.Errorf("expected email, got %q", found.Email)
	}
}

func TestCustomerRepository_Delete(t *testing.T) {
	db := setupCustomerTestDB(t)
	defer db.Close()
	repo := NewCustomerRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	c := domain.NewCustomer("John", "9876543210")
	_ = repo.Create(context.Background(), c)

	err := repo.Delete(context.Background(), c.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = repo.GetByID(context.Background(), c.ID)
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestCustomerRepository_CountByStatus(t *testing.T) {
	db := setupCustomerTestDB(t)
	defer db.Close()
	repo := NewCustomerRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	c1 := domain.NewCustomer("Active", "9876543210")
	_ = repo.Create(context.Background(), c1)

	c2 := domain.NewCustomer("Inactive", "9876543211")
	_ = repo.Create(context.Background(), c2)
	c2.Status = "inactive"
	_ = repo.Update(context.Background(), c2)

	total, active, inactive, err := repo.CountByStatus(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 2 {
		t.Errorf("expected 2, got %d", total)
	}
	if active != 1 {
		t.Errorf("expected 1 active, got %d", active)
	}
	if inactive != 1 {
		t.Errorf("expected 1 inactive, got %d", inactive)
	}
}
