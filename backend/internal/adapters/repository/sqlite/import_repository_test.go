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
)

func setupImportTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE import_templates (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			target_entity TEXT NOT NULL CHECK (target_entity IN ('staff','customers','services','products','expenses','advances','salary')),
			column_mapping TEXT NOT NULL DEFAULT '{}',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE TABLE import_jobs (
			id TEXT PRIMARY KEY,
			template_id TEXT DEFAULT '',
			file_name TEXT NOT NULL,
			file_path TEXT NOT NULL,
			target_entity TEXT NOT NULL CHECK (target_entity IN ('staff','customers','services','products','expenses','advances','salary')),
			status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','validating','validated','importing','completed','failed')),
			total_rows INTEGER NOT NULL DEFAULT 0,
			valid_rows INTEGER NOT NULL DEFAULT 0,
			invalid_rows INTEGER NOT NULL DEFAULT 0,
			imported_rows INTEGER NOT NULL DEFAULT 0,
			column_mapping TEXT NOT NULL DEFAULT '{}',
			error_message TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE TABLE import_logs (
			id TEXT PRIMARY KEY,
			job_id TEXT NOT NULL REFERENCES import_jobs(id),
			row_number INTEGER NOT NULL,
			status TEXT NOT NULL CHECK (status IN ('success','error','warning','skipped')),
			message TEXT NOT NULL DEFAULT '',
			row_data TEXT NOT NULL DEFAULT '{}',
			created_at TEXT NOT NULL
		);
	`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestImportRepository_CreateAndGetJob(t *testing.T) {
	db := setupImportTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewImportRepository(db, log)
	ctx := context.Background()

	job := domain.NewImportJob("test.xlsx", "/tmp/test.xlsx", domain.ImportEntityStaff)
	job.TotalRows = 100

	err := repo.CreateJob(ctx, job)
	if err != nil {
		t.Fatalf("CreateJob: %v", err)
	}

	got, err := repo.GetJob(ctx, job.ID)
	if err != nil {
		t.Fatalf("GetJob: %v", err)
	}
	if got.FileName != "test.xlsx" {
		t.Errorf("FileName = %q, want test.xlsx", got.FileName)
	}
	if got.TotalRows != 100 {
		t.Errorf("TotalRows = %d, want 100", got.TotalRows)
	}
}

func TestImportRepository_UpdateJob(t *testing.T) {
	db := setupImportTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewImportRepository(db, log)
	ctx := context.Background()

	job := domain.NewImportJob("update.csv", "/tmp/update.csv", domain.ImportEntityCustomers)
	repo.CreateJob(ctx, job)

	job.Status = domain.ImportStatusCompleted
	job.ValidRows = 50
	job.ImportedRows = 50
	err := repo.UpdateJob(ctx, job)
	if err != nil {
		t.Fatalf("UpdateJob: %v", err)
	}

	got, _ := repo.GetJob(ctx, job.ID)
	if got.Status != domain.ImportStatusCompleted {
		t.Errorf("Status = %q, want completed", got.Status)
	}
	if got.ImportedRows != 50 {
		t.Errorf("ImportedRows = %d, want 50", got.ImportedRows)
	}
}

func TestImportRepository_ListJobs(t *testing.T) {
	db := setupImportTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewImportRepository(db, log)
	ctx := context.Background()

	repo.CreateJob(ctx, domain.NewImportJob("a.xlsx", "/tmp/a.xlsx", domain.ImportEntityStaff))
	repo.CreateJob(ctx, domain.NewImportJob("b.csv", "/tmp/b.csv", domain.ImportEntityProducts))

	jobs, total, err := repo.ListJobs(ctx, 10, 0)
	if err != nil {
		t.Fatalf("ListJobs: %v", err)
	}
	if total != 2 {
		t.Errorf("total = %d, want 2", total)
	}
	if len(jobs) != 2 {
		t.Errorf("len = %d, want 2", len(jobs))
	}
}

func TestImportRepository_Logs(t *testing.T) {
	db := setupImportTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewImportRepository(db, log)
	ctx := context.Background()

	job := domain.NewImportJob("logs.xlsx", "/tmp/logs.xlsx", domain.ImportEntityServices)
	repo.CreateJob(ctx, job)

	// Create individual log
	l := domain.NewImportLog(job.ID, 2, domain.ImportLogError, "missing name", "{}")
	err := repo.CreateLog(ctx, l)
	if err != nil {
		t.Fatalf("CreateLog: %v", err)
	}

	// Create batch
	batch := []domain.ImportLog{
		*domain.NewImportLog(job.ID, 3, domain.ImportLogSuccess, "", "{}"),
		*domain.NewImportLog(job.ID, 4, domain.ImportLogSuccess, "", "{}"),
		*domain.NewImportLog(job.ID, 5, domain.ImportLogWarning, "duplicate phone", "{}"),
	}
	err = repo.CreateLogBatch(ctx, batch)
	if err != nil {
		t.Fatalf("CreateLogBatch: %v", err)
	}

	// List all
	logs, total, err := repo.ListLogs(ctx, job.ID, "", 10, 0)
	if err != nil {
		t.Fatalf("ListLogs: %v", err)
	}
	if total != 4 {
		t.Errorf("total = %d, want 4", total)
	}
	if len(logs) != 4 {
		t.Errorf("len = %d, want 4", len(logs))
	}

	// List errors only
	logs, total, _ = repo.ListLogs(ctx, job.ID, domain.ImportLogError, 10, 0)
	if total != 1 {
		t.Errorf("error total = %d, want 1", total)
	}
}

func TestImportRepository_Templates(t *testing.T) {
	db := setupImportTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewImportRepository(db, log)
	ctx := context.Background()

	tmpl := domain.NewImportTemplate("Staff Import", domain.ImportEntityStaff, `[{"source_column":"Name","target_field":"name"}]`)
	err := repo.CreateTemplate(ctx, tmpl)
	if err != nil {
		t.Fatalf("CreateTemplate: %v", err)
	}

	templates, err := repo.ListTemplates(ctx, domain.ImportEntityStaff)
	if err != nil {
		t.Fatalf("ListTemplates: %v", err)
	}
	if len(templates) != 1 {
		t.Errorf("len = %d, want 1", len(templates))
	}
	if templates[0].Name != "Staff Import" {
		t.Errorf("Name = %q", templates[0].Name)
	}

	// List all
	templates, _ = repo.ListTemplates(ctx, "")
	if len(templates) != 1 {
		t.Errorf("all len = %d, want 1", len(templates))
	}
}

func TestImportRepository_NotFound(t *testing.T) {
	db := setupImportTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewImportRepository(db, log)
	ctx := context.Background()

	_, err := repo.GetJob(ctx, uuid.New())
	if err == nil {
		t.Error("expected error for non-existent job")
	}
}
