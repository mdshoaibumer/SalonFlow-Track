package sqlite_test

import (
	"context"
	"database/sql"
	"testing"

	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/salonflow/salonflow-track/internal/adapters/repository/sqlite"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/pkg/apperror"
	"github.com/salonflow/salonflow-track/pkg/uid"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:?_foreign_keys=ON")
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	_, err = db.Exec(`
		CREATE TABLE staff (
			id                      TEXT PRIMARY KEY,
			staff_code              TEXT NOT NULL UNIQUE,
			full_name               TEXT NOT NULL,
			phone                   TEXT NOT NULL,
			email                   TEXT DEFAULT '',
			gender                  TEXT NOT NULL DEFAULT 'male',
			designation             TEXT NOT NULL DEFAULT 'stylist',
			joining_date            TEXT NOT NULL,
			base_salary             REAL NOT NULL DEFAULT 0,
			commission_percentage   REAL NOT NULL DEFAULT 0,
			status                  TEXT NOT NULL DEFAULT 'active',
			created_at              TEXT NOT NULL,
			updated_at              TEXT NOT NULL
		);
		CREATE UNIQUE INDEX idx_staff_phone ON staff(phone);
	`)
	if err != nil {
		t.Fatalf("create test table: %v", err)
	}
	return db
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
}

func TestStaffRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewStaffRepository(db, testLogger())
	ctx := context.Background()

	staff := domain.NewStaff("Nazim Khan", "9876543210", "stylist", 15000, 10)

	err := repo.Create(ctx, staff)
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	got, err := repo.GetByID(ctx, staff.ID)
	if err != nil {
		t.Fatalf("GetByID() error: %v", err)
	}
	if got.FullName != "Nazim Khan" {
		t.Errorf("expected name Nazim Khan, got %s", got.FullName)
	}
	if got.Phone != "9876543210" {
		t.Errorf("expected phone 9876543210, got %s", got.Phone)
	}
	if got.BaseSalary != 15000 {
		t.Errorf("expected salary 15000, got %f", got.BaseSalary)
	}
	if got.CommissionPercentage != 10 {
		t.Errorf("expected commission 10, got %f", got.CommissionPercentage)
	}
	if got.Status != "active" {
		t.Errorf("expected status active, got %s", got.Status)
	}
}

func TestStaffRepository_Create_DuplicatePhone(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewStaffRepository(db, testLogger())
	ctx := context.Background()

	s1 := domain.NewStaff("Staff 1", "9876543210", "stylist", 10000, 5)
	repo.Create(ctx, s1)

	s2 := domain.NewStaff("Staff 2", "9876543210", "assistant", 8000, 0)
	err := repo.Create(ctx, s2)
	if err == nil {
		t.Fatal("expected error for duplicate phone")
	}
	if !apperror.Is(err, apperror.KindConflict) {
		t.Errorf("expected Conflict error, got: %v", err)
	}
}

func TestStaffRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewStaffRepository(db, testLogger())
	ctx := context.Background()

	_, err := repo.GetByID(ctx, uid.New())
	if err == nil {
		t.Fatal("expected error for non-existent staff")
	}
	if !apperror.Is(err, apperror.KindNotFound) {
		t.Errorf("expected NotFound error, got: %v", err)
	}
}

func TestStaffRepository_GetByPhone(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewStaffRepository(db, testLogger())
	ctx := context.Background()

	staff := domain.NewStaff("Ravi", "9111111111", "stylist", 12000, 8)
	repo.Create(ctx, staff)

	got, err := repo.GetByPhone(ctx, "9111111111")
	if err != nil {
		t.Fatalf("GetByPhone() error: %v", err)
	}
	if got.FullName != "Ravi" {
		t.Errorf("expected name Ravi, got %s", got.FullName)
	}
}

func TestStaffRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewStaffRepository(db, testLogger())
	ctx := context.Background()

	phones := []string{"9000000001", "9000000002", "9000000003"}
	for i, name := range []string{"Alice", "Bob", "Charlie"} {
		s := domain.NewStaff(name, phones[i], "stylist", 10000, 5)
		if err := repo.Create(ctx, s); err != nil {
			t.Fatalf("Create(%s) error: %v", name, err)
		}
	}

	results, total, err := repo.List(ctx, ports.StaffFilter{Limit: 10})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}
}

func TestStaffRepository_List_Search(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewStaffRepository(db, testLogger())
	ctx := context.Background()

	repo.Create(ctx, domain.NewStaff("Alice Smith", "9100000001", "stylist", 10000, 5))
	repo.Create(ctx, domain.NewStaff("Bob Jones", "9100000002", "assistant", 8000, 0))

	results, total, err := repo.List(ctx, ports.StaffFilter{Search: "Alice", Limit: 10})
	if err != nil {
		t.Fatalf("List(search=Alice) error: %v", err)
	}
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
	if len(results) != 1 || results[0].FullName != "Alice Smith" {
		t.Error("expected Alice in results")
	}
}

func TestStaffRepository_List_FilterStatus(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewStaffRepository(db, testLogger())
	ctx := context.Background()

	active := domain.NewStaff("Active", "9200000001", "stylist", 10000, 5)
	repo.Create(ctx, active)

	inactive := domain.NewStaff("Inactive", "9200000002", "stylist", 10000, 5)
	repo.Create(ctx, inactive)
	inactive.Status = "inactive"
	repo.Update(ctx, inactive)

	results, total, err := repo.List(ctx, ports.StaffFilter{Status: "active", Limit: 10})
	if err != nil {
		t.Fatalf("List(status=active) error: %v", err)
	}
	if total != 1 {
		t.Errorf("expected 1 active, got %d", total)
	}
	if len(results) != 1 || results[0].FullName != "Active" {
		t.Error("expected Active staff in results")
	}
}

func TestStaffRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewStaffRepository(db, testLogger())
	ctx := context.Background()

	staff := domain.NewStaff("Ravi", "9300000001", "stylist", 12000, 10)
	repo.Create(ctx, staff)

	staff.FullName = "Ravi Kumar"
	staff.BaseSalary = 18000
	staff.CommissionPercentage = 15
	err := repo.Update(ctx, staff)
	if err != nil {
		t.Fatalf("Update() error: %v", err)
	}

	got, _ := repo.GetByID(ctx, staff.ID)
	if got.FullName != "Ravi Kumar" {
		t.Errorf("expected updated name, got %s", got.FullName)
	}
	if got.BaseSalary != 18000 {
		t.Errorf("expected 18000, got %f", got.BaseSalary)
	}
	if got.CommissionPercentage != 15 {
		t.Errorf("expected 15, got %f", got.CommissionPercentage)
	}
}

func TestStaffRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewStaffRepository(db, testLogger())
	ctx := context.Background()

	staff := domain.NewStaff("ToDelete", "9400000001", "assistant", 8000, 0)
	repo.Create(ctx, staff)

	err := repo.Delete(ctx, staff.ID)
	if err != nil {
		t.Fatalf("Delete() error: %v", err)
	}

	_, err = repo.GetByID(ctx, staff.ID)
	if !apperror.Is(err, apperror.KindNotFound) {
		t.Errorf("expected NotFound after delete, got: %v", err)
	}
}

func TestStaffRepository_CountByStatus(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewStaffRepository(db, testLogger())
	ctx := context.Background()

	active := domain.NewStaff("Active1", "9500000001", "stylist", 10000, 5)
	repo.Create(ctx, active)

	active2 := domain.NewStaff("Active2", "9500000002", "stylist", 10000, 5)
	repo.Create(ctx, active2)

	inactive := domain.NewStaff("Inactive1", "9500000003", "stylist", 10000, 5)
	repo.Create(ctx, inactive)
	inactive.Status = "inactive"
	repo.Update(ctx, inactive)

	total, act, inact, err := repo.CountByStatus(ctx)
	if err != nil {
		t.Fatalf("CountByStatus() error: %v", err)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if act != 2 {
		t.Errorf("expected active 2, got %d", act)
	}
	if inact != 1 {
		t.Errorf("expected inactive 1, got %d", inact)
	}
}
