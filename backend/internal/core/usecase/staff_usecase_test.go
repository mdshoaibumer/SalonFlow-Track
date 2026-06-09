package usecase_test

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/salonflow/salonflow-track/internal/adapters/repository/sqlite"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
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

func newTestUseCase(t *testing.T) *usecase.StaffUseCase {
	db := setupTestDB(t)
	repo := sqlite.NewStaffRepository(db, testLogger())
	return usecase.NewStaffUseCase(repo, testLogger())
}

func TestCreateStaff_Success(t *testing.T) {
	uc := newTestUseCase(t)
	ctx := context.Background()

	input := usecase.CreateStaffInput{
		FullName:             "Nazim Khan",
		Phone:                "9876543210",
		Gender:               "male",
		Designation:          "stylist",
		JoiningDate:          "2024-01-15",
		BaseSalary:           15000,
		CommissionPercentage: 10,
	}

	staff, err := uc.Create(ctx, input)
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}
	if staff.FullName != "Nazim Khan" {
		t.Errorf("expected name, got %s", staff.FullName)
	}
	if staff.StaffCode == "" {
		t.Error("expected staff code to be generated")
	}
}

func TestCreateStaff_ValidationError(t *testing.T) {
	uc := newTestUseCase(t)
	ctx := context.Background()

	input := usecase.CreateStaffInput{
		FullName:    "",
		Phone:       "9876543210",
		Designation: "stylist",
	}

	_, err := uc.Create(ctx, input)
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !apperror.Is(err, apperror.KindValidation) {
		t.Errorf("expected Validation error kind, got: %v", err)
	}
}

func TestCreateStaff_DuplicatePhone(t *testing.T) {
	uc := newTestUseCase(t)
	ctx := context.Background()

	input := usecase.CreateStaffInput{
		FullName:    "Staff One",
		Phone:       "9876543210",
		Gender:      "male",
		Designation: "stylist",
		JoiningDate: "2024-01-15",
		BaseSalary:  10000,
	}
	uc.Create(ctx, input)

	input.FullName = "Staff Two"
	_, err := uc.Create(ctx, input)
	if err == nil {
		t.Fatal("expected conflict error for duplicate phone")
	}
	if !apperror.Is(err, apperror.KindConflict) {
		t.Errorf("expected Conflict error, got: %v", err)
	}
}

func TestListStaff_Pagination(t *testing.T) {
	uc := newTestUseCase(t)
	ctx := context.Background()

	phones := []string{"9100000001", "9100000002", "9100000003", "9100000004", "9100000005"}
	for i := 0; i < 5; i++ {
		input := usecase.CreateStaffInput{
			FullName:    "Staff " + string(rune('A'+i)),
			Phone:       phones[i],
			Gender:      "male",
			Designation: "stylist",
			JoiningDate: "2024-01-15",
			BaseSalary:  10000,
		}
		_, err := uc.Create(ctx, input)
		if err != nil {
			t.Fatalf("Create staff %d error: %v", i, err)
		}
	}

	output, err := uc.List(ctx, usecase.ListStaffInput{Page: 1, PerPage: 2})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if output.Total != 5 {
		t.Errorf("expected total 5, got %d", output.Total)
	}
	if len(output.Staff) != 2 {
		t.Errorf("expected 2 items on page 1, got %d", len(output.Staff))
	}
	if output.TotalPages != 3 {
		t.Errorf("expected 3 total pages, got %d", output.TotalPages)
	}
}

func TestUpdateStaff_Success(t *testing.T) {
	uc := newTestUseCase(t)
	ctx := context.Background()

	staff, _ := uc.Create(ctx, usecase.CreateStaffInput{
		FullName:    "Original Name",
		Phone:       "9876543210",
		Gender:      "male",
		Designation: "stylist",
		JoiningDate: "2024-01-15",
		BaseSalary:  10000,
	})

	updated, err := uc.Update(ctx, staff.ID, usecase.UpdateStaffInput{
		FullName:             "New Name",
		Phone:                "9876543210",
		Gender:               "male",
		Designation:          "manager",
		JoiningDate:          "2024-01-15",
		BaseSalary:           20000,
		CommissionPercentage: 15,
		Status:               "active",
	})
	if err != nil {
		t.Fatalf("Update() error: %v", err)
	}
	if updated.FullName != "New Name" {
		t.Errorf("expected updated name, got %s", updated.FullName)
	}
	if updated.BaseSalary != 20000 {
		t.Errorf("expected salary 20000, got %f", updated.BaseSalary)
	}
}

func TestDeleteStaff_Success(t *testing.T) {
	uc := newTestUseCase(t)
	ctx := context.Background()

	staff, _ := uc.Create(ctx, usecase.CreateStaffInput{
		FullName:    "To Delete",
		Phone:       "9876543210",
		Gender:      "male",
		Designation: "stylist",
		JoiningDate: "2024-01-15",
		BaseSalary:  10000,
	})

	err := uc.Delete(ctx, staff.ID)
	if err != nil {
		t.Fatalf("Delete() error: %v", err)
	}

	_, err = uc.GetByID(ctx, staff.ID)
	if !apperror.Is(err, apperror.KindNotFound) {
		t.Errorf("expected NotFound after delete, got: %v", err)
	}
}

func TestDeleteStaff_NotFound(t *testing.T) {
	uc := newTestUseCase(t)
	ctx := context.Background()

	err := uc.Delete(ctx, uid.New())
	if !apperror.Is(err, apperror.KindNotFound) {
		t.Errorf("expected NotFound error, got: %v", err)
	}
}

func TestStats(t *testing.T) {
	uc := newTestUseCase(t)
	ctx := context.Background()

	_, err := uc.Create(ctx, usecase.CreateStaffInput{
		FullName: "Active1", Phone: "9100000001", Gender: "male", Designation: "stylist", JoiningDate: "2024-01-15", BaseSalary: 10000,
	})
	if err != nil {
		t.Fatalf("Create Active1 error: %v", err)
	}
	_, err = uc.Create(ctx, usecase.CreateStaffInput{
		FullName: "Active2", Phone: "9100000002", Gender: "female", Designation: "assistant", JoiningDate: "2024-02-01", BaseSalary: 8000,
	})
	if err != nil {
		t.Fatalf("Create Active2 error: %v", err)
	}

	// Create and deactivate one
	staff, err := uc.Create(ctx, usecase.CreateStaffInput{
		FullName: "WillDeactivate", Phone: "9100000003", Gender: "male", Designation: "receptionist", JoiningDate: "2024-03-01", BaseSalary: 7000,
	})
	if err != nil {
		t.Fatalf("Create WillDeactivate error: %v", err)
	}
	_, err = uc.Update(ctx, staff.ID, usecase.UpdateStaffInput{
		FullName: "WillDeactivate", Phone: "9100000003", Gender: "male", Designation: "receptionist", JoiningDate: "2024-03-01", BaseSalary: 7000, Status: "inactive",
	})
	if err != nil {
		t.Fatalf("Update to inactive error: %v", err)
	}

	stats, err := uc.Stats(ctx)
	if err != nil {
		t.Fatalf("Stats() error: %v", err)
	}
	if stats.Total != 3 {
		t.Errorf("expected total 3, got %d", stats.Total)
	}
	if stats.Active != 2 {
		t.Errorf("expected active 2, got %d", stats.Active)
	}
	if stats.Inactive != 1 {
		t.Errorf("expected inactive 1, got %d", stats.Inactive)
	}
}
