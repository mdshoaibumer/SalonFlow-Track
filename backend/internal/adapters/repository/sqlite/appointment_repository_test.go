package sqlite

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/pkg/uid"
)

func setupAppointmentTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE appointments (
			id TEXT PRIMARY KEY,
			customer_id TEXT NOT NULL,
			staff_id TEXT NOT NULL,
			appointment_date TEXT NOT NULL DEFAULT '',
			start_time TEXT NOT NULL,
			end_time TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'booked',
			notes TEXT NOT NULL DEFAULT '',
			is_walkin INTEGER NOT NULL DEFAULT 0,
			total_amount REAL NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE TABLE appointment_services (
			id TEXT PRIMARY KEY,
			appointment_id TEXT NOT NULL,
			service_id TEXT NOT NULL,
			service_name TEXT NOT NULL DEFAULT '',
			duration_minutes INTEGER NOT NULL DEFAULT 0,
			price REAL NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL,
			FOREIGN KEY (appointment_id) REFERENCES appointments(id)
		);
		CREATE TABLE appointment_history (
			id TEXT PRIMARY KEY,
			appointment_id TEXT NOT NULL,
			old_status TEXT NOT NULL DEFAULT '',
			new_status TEXT NOT NULL DEFAULT '',
			changed_by TEXT NOT NULL DEFAULT '',
			note TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			FOREIGN KEY (appointment_id) REFERENCES appointments(id)
		);
	`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestAppointmentRepository_Create(t *testing.T) {
	db := setupAppointmentTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewAppointmentRepository(db, log)
	ctx := context.Background()

	appt := &domain.Appointment{
		ID:              uid.New(),
		CustomerID:      uid.New().String(),
		StaffID:         uid.New().String(),
		AppointmentDate: "2024-06-15",
		StartTime:       "10:00",
		EndTime:         "11:00",
		Status:          domain.AppointmentStatusBooked,
		Notes:           "Test appointment",
		TotalAmount:     500,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}

	err := repo.Create(ctx, appt)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	got, err := repo.GetByID(ctx, appt.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got.ID != appt.ID {
		t.Errorf("got ID %v, want %v", got.ID, appt.ID)
	}
	if got.Status != domain.AppointmentStatusBooked {
		t.Errorf("got status %v, want %v", got.Status, domain.AppointmentStatusBooked)
	}
}

func TestAppointmentRepository_Update(t *testing.T) {
	db := setupAppointmentTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewAppointmentRepository(db, log)
	ctx := context.Background()

	appt := &domain.Appointment{
		ID:              uid.New(),
		CustomerID:      uid.New().String(),
		StaffID:         uid.New().String(),
		AppointmentDate: "2024-06-15",
		StartTime:       "10:00",
		EndTime:         "11:00",
		Status:          domain.AppointmentStatusBooked,
		Notes:           "Original",
		TotalAmount:     300,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}
	_ = repo.Create(ctx, appt)

	appt.Notes = "Updated"
	appt.TotalAmount = 600
	err := repo.Update(ctx, appt)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	got, _ := repo.GetByID(ctx, appt.ID)
	if got.Notes != "Updated" {
		t.Errorf("got notes %q, want %q", got.Notes, "Updated")
	}
	if got.TotalAmount != 600 {
		t.Errorf("got amount %v, want 600", got.TotalAmount)
	}
}

func TestAppointmentRepository_List(t *testing.T) {
	db := setupAppointmentTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewAppointmentRepository(db, log)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		appt := &domain.Appointment{
			ID:              uid.New(),
			CustomerID:      uid.New().String(),
			StaffID:         uid.New().String(),
			AppointmentDate: "2024-06-15",
			StartTime:       "10:00",
			EndTime:         "11:00",
			Status:          domain.AppointmentStatusBooked,
			CreatedAt:       time.Now().UTC(),
			UpdatedAt:       time.Now().UTC(),
		}
		_ = repo.Create(ctx, appt)
	}

	list, total, err := repo.List(ctx, domain.AppointmentFilter{})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if total != 3 {
		t.Errorf("got total %d, want 3", total)
	}
	if len(list) != 3 {
		t.Errorf("got %d appointments, want 3", len(list))
	}
}

func TestAppointmentRepository_Delete(t *testing.T) {
	db := setupAppointmentTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewAppointmentRepository(db, log)
	ctx := context.Background()

	appt := &domain.Appointment{
		ID:              uid.New(),
		CustomerID:      uid.New().String(),
		StaffID:         uid.New().String(),
		AppointmentDate: "2024-06-15",
		StartTime:       "10:00",
		EndTime:         "11:00",
		Status:          domain.AppointmentStatusBooked,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}
	_ = repo.Create(ctx, appt)

	err := repo.Delete(ctx, appt.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = repo.GetByID(ctx, appt.ID)
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestAppointmentRepository_Services(t *testing.T) {
	db := setupAppointmentTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewAppointmentRepository(db, log)
	ctx := context.Background()

	apptID := uid.New()
	appt := &domain.Appointment{
		ID:              apptID,
		CustomerID:      uid.New().String(),
		StaffID:         uid.New().String(),
		AppointmentDate: "2024-06-15",
		StartTime:       "10:00",
		EndTime:         "11:00",
		Status:          domain.AppointmentStatusBooked,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}
	_ = repo.Create(ctx, appt)

	services := []domain.AppointmentService{
		{
			ID:              uid.New(),
			AppointmentID:   apptID,
			ServiceID:       uid.New().String(),
			ServiceName:     "Haircut",
			DurationMinutes: 30,
			Price:           200,
			CreatedAt:       time.Now().UTC(),
		},
	}
	err := repo.AddServices(ctx, services)
	if err != nil {
		t.Fatalf("AddServices failed: %v", err)
	}

	got, err := repo.GetServices(ctx, apptID)
	if err != nil {
		t.Fatalf("GetServices failed: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("got %d services, want 1", len(got))
	}
}

func TestAppointmentRepository_History(t *testing.T) {
	db := setupAppointmentTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewAppointmentRepository(db, log)
	ctx := context.Background()

	apptID := uid.New()
	appt := &domain.Appointment{
		ID:              apptID,
		CustomerID:      uid.New().String(),
		StaffID:         uid.New().String(),
		AppointmentDate: "2024-06-15",
		StartTime:       "10:00",
		EndTime:         "11:00",
		Status:          domain.AppointmentStatusBooked,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}
	_ = repo.Create(ctx, appt)

	history := &domain.AppointmentHistory{
		ID:            uid.New(),
		AppointmentID: apptID,
		OldStatus:     "booked",
		NewStatus:     "confirmed",
		ChangedBy:     "system",
		Note:          "Auto confirmed",
		CreatedAt:     time.Now().UTC(),
	}
	err := repo.AddHistory(ctx, history)
	if err != nil {
		t.Fatalf("AddHistory failed: %v", err)
	}

	entries, err := repo.GetHistory(ctx, apptID)
	if err != nil {
		t.Fatalf("GetHistory failed: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("got %d history entries, want 1", len(entries))
	}
	if entries[0].NewStatus != "confirmed" {
		t.Errorf("got new_status %q, want %q", entries[0].NewStatus, "confirmed")
	}
}
