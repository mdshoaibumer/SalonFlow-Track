package handler_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/salonflow/salonflow-track/internal/adapters/handler"
	"github.com/salonflow/salonflow-track/internal/adapters/repository/sqlite"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

func setupTestHandler(t *testing.T) (*handler.StaffHandler, *chi.Mux) {
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
		t.Fatalf("create table: %v", err)
	}

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
	repo := sqlite.NewStaffRepository(db, log)
	uc := usecase.NewStaffUseCase(repo, log)
	h := handler.NewStaffHandler(uc)

	r := chi.NewRouter()
	r.Mount("/api/v1/staff", h.Routes())
	return h, r
}

func TestAPI_CreateStaff(t *testing.T) {
	_, router := setupTestHandler(t)

	body := `{"full_name":"Nazim Khan","phone":"9876543210","gender":"male","designation":"stylist","joining_date":"2024-01-15","base_salary":15000,"commission_percentage":10}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/staff", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp["success"] != true {
		t.Error("expected success true")
	}
	data := resp["data"].(map[string]interface{})
	if data["full_name"] != "Nazim Khan" {
		t.Errorf("expected name in response, got %v", data["full_name"])
	}
}

func TestAPI_CreateStaff_Validation(t *testing.T) {
	_, router := setupTestHandler(t)

	body := `{"full_name":"","phone":"9876543210","designation":"stylist"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/staff", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestAPI_ListStaff(t *testing.T) {
	_, router := setupTestHandler(t)

	// Create staff first
	body := `{"full_name":"Alice","phone":"9100000001","gender":"female","designation":"stylist","joining_date":"2024-01-15","base_salary":10000,"commission_percentage":5}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/staff", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// List
	req = httptest.NewRequest(http.MethodGet, "/api/v1/staff", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp["success"] != true {
		t.Error("expected success true")
	}
	meta := resp["meta"].(map[string]interface{})
	if meta["total"].(float64) != 1 {
		t.Errorf("expected total 1, got %v", meta["total"])
	}
}

func TestAPI_GetStaff_NotFound(t *testing.T) {
	_, router := setupTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/staff/00000000-0000-0000-0000-000000000001", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestAPI_DeleteStaff(t *testing.T) {
	_, router := setupTestHandler(t)

	// Create first
	body := `{"full_name":"ToDelete","phone":"9200000001","gender":"male","designation":"assistant","joining_date":"2024-06-01","base_salary":8000,"commission_percentage":0}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/staff", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	var createResp map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &createResp)
	data := createResp["data"].(map[string]interface{})
	id := data["id"].(string)

	// Delete
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/staff/"+id, nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	// Verify deleted
	req = httptest.NewRequest(http.MethodGet, "/api/v1/staff/"+id, nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", rec.Code)
	}
}

func TestAPI_Stats(t *testing.T) {
	_, router := setupTestHandler(t)

	body := `{"full_name":"Staff1","phone":"9300000001","gender":"male","designation":"stylist","joining_date":"2024-01-15","base_salary":10000,"commission_percentage":5}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/staff", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Stats
	req = httptest.NewRequest(http.MethodGet, "/api/v1/staff/stats", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["total"].(float64) != 1 {
		t.Errorf("expected total 1, got %v", data["total"])
	}
	if data["active"].(float64) != 1 {
		t.Errorf("expected active 1, got %v", data["active"])
	}
}
