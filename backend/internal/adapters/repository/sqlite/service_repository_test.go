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

func setupServiceTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`
		CREATE TABLE services (
			id TEXT PRIMARY KEY,
			service_code TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			category TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			duration_minutes INTEGER NOT NULL,
			price REAL NOT NULL,
			cost_price REAL NOT NULL DEFAULT 0,
			commission_type TEXT NOT NULL DEFAULT 'percentage',
			commission_value REAL NOT NULL DEFAULT 0,
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

func TestServiceRepository_Create(t *testing.T) {
	db := setupServiceTestDB(t)
	defer db.Close()
	repo := NewServiceRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	svc := domain.NewService("Hair Cut", "hair", 30, 300)
	err := repo.Create(context.Background(), svc)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	found, err := repo.GetByID(context.Background(), svc.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.Name != "Hair Cut" {
		t.Errorf("expected 'Hair Cut', got %q", found.Name)
	}
}

func TestServiceRepository_GetByName(t *testing.T) {
	db := setupServiceTestDB(t)
	defer db.Close()
	repo := NewServiceRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	svc := domain.NewService("Facial Treatment", "facial", 45, 800)
	_ = repo.Create(context.Background(), svc)

	found, err := repo.GetByName(context.Background(), "facial treatment")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.ID != svc.ID {
		t.Error("expected same service")
	}
}

func TestServiceRepository_List(t *testing.T) {
	db := setupServiceTestDB(t)
	defer db.Close()
	repo := NewServiceRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	_ = repo.Create(context.Background(), domain.NewService("Hair Cut", "hair", 30, 300))
	_ = repo.Create(context.Background(), domain.NewService("Facial", "facial", 45, 800))
	_ = repo.Create(context.Background(), domain.NewService("Head Massage", "massage", 20, 200))

	services, total, err := repo.List(context.Background(), ports.ServiceFilter{Limit: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 3 {
		t.Errorf("expected 3 total, got %d", total)
	}
	if len(services) != 3 {
		t.Errorf("expected 3 services, got %d", len(services))
	}
}

func TestServiceRepository_ListByCategory(t *testing.T) {
	db := setupServiceTestDB(t)
	defer db.Close()
	repo := NewServiceRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	_ = repo.Create(context.Background(), domain.NewService("Hair Cut", "hair", 30, 300))
	_ = repo.Create(context.Background(), domain.NewService("Hair Spa", "hair", 60, 1200))
	_ = repo.Create(context.Background(), domain.NewService("Facial", "facial", 45, 800))

	services, total, err := repo.List(context.Background(), ports.ServiceFilter{Category: "hair", Limit: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 2 {
		t.Errorf("expected 2 hair services, got %d", total)
	}
	if len(services) != 2 {
		t.Errorf("expected 2 services, got %d", len(services))
	}
}

func TestServiceRepository_Update(t *testing.T) {
	db := setupServiceTestDB(t)
	defer db.Close()
	repo := NewServiceRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	svc := domain.NewService("Hair Cut", "hair", 30, 300)
	_ = repo.Create(context.Background(), svc)

	svc.Price = 400
	svc.DurationMinutes = 45
	err := repo.Update(context.Background(), svc)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	found, _ := repo.GetByID(context.Background(), svc.ID)
	if found.Price != 400 {
		t.Errorf("expected price 400, got %f", found.Price)
	}
	if found.DurationMinutes != 45 {
		t.Errorf("expected duration 45, got %d", found.DurationMinutes)
	}
}

func TestServiceRepository_Delete(t *testing.T) {
	db := setupServiceTestDB(t)
	defer db.Close()
	repo := NewServiceRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	svc := domain.NewService("Hair Cut", "hair", 30, 300)
	_ = repo.Create(context.Background(), svc)

	err := repo.Delete(context.Background(), svc.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = repo.GetByID(context.Background(), svc.ID)
	if err == nil {
		t.Error("expected error after deletion")
	}
}

func TestServiceRepository_CountByStatus(t *testing.T) {
	db := setupServiceTestDB(t)
	defer db.Close()
	repo := NewServiceRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	svc1 := domain.NewService("Hair Cut", "hair", 30, 300)
	_ = repo.Create(context.Background(), svc1)

	svc2 := domain.NewService("Facial", "facial", 45, 800)
	_ = repo.Create(context.Background(), svc2)
	svc2.Status = "inactive"
	_ = repo.Update(context.Background(), svc2)

	total, active, inactive, err := repo.CountByStatus(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if active != 1 {
		t.Errorf("expected active 1, got %d", active)
	}
	if inactive != 1 {
		t.Errorf("expected inactive 1, got %d", inactive)
	}
}

func TestServiceRepository_Search(t *testing.T) {
	db := setupServiceTestDB(t)
	defer db.Close()
	repo := NewServiceRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	_ = repo.Create(context.Background(), domain.NewService("Hair Cut", "hair", 30, 300))
	_ = repo.Create(context.Background(), domain.NewService("Hair Spa", "hair", 60, 1200))
	_ = repo.Create(context.Background(), domain.NewService("Facial", "facial", 45, 800))

	services, total, err := repo.List(context.Background(), ports.ServiceFilter{Search: "hair", Limit: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 2 {
		t.Errorf("expected 2 results, got %d", total)
	}
	if len(services) != 2 {
		t.Errorf("expected 2 services, got %d", len(services))
	}
}
