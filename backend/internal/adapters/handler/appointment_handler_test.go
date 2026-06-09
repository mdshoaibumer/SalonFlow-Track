package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/salonflow/salonflow-track/internal/adapters/handler"
	"github.com/salonflow/salonflow-track/internal/adapters/repository/sqlite"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
	"github.com/salonflow/salonflow-track/internal/testutil"
)

func setupAppointmentHandler(t *testing.T) chi.Router {
	t.Helper()
	db := testutil.TestDB(t)
	log := testutil.TestLogger(t)
	repo := sqlite.NewAppointmentRepository(db, log)
	uc := usecase.NewAppointmentUseCase(repo, log)
	h := handler.NewAppointmentHandler(uc)
	r := chi.NewRouter()
	r.Mount("/api/v1/appointments", h.Routes())
	return r
}

func TestAppointmentAPI_Create(t *testing.T) {
	router := setupAppointmentHandler(t)

	body := `{
		"customer_id":"01912345-6789-7abc-def0-123456789001",
		"staff_id":"01912345-6789-7abc-def0-123456789002",
		"appointment_date":"2024-12-20",
		"start_time":"10:00",
		"end_time":"11:00",
		"notes":"First visit",
		"is_walkin":false,
		"services":[{"service_id":"svc1","service_name":"Haircut","duration_minutes":45,"price":500}]
	}`
	w := testutil.DoRequest(router, http.MethodPost, "/api/v1/appointments", body)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["success"] != true {
		t.Fatal("expected success true")
	}
}

func TestAppointmentAPI_Create_Validation(t *testing.T) {
	router := setupAppointmentHandler(t)

	// Missing required fields
	body := `{"notes":"incomplete"}`
	w := testutil.DoRequest(router, http.MethodPost, "/api/v1/appointments", body)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAppointmentAPI_List(t *testing.T) {
	router := setupAppointmentHandler(t)

	// Create one
	body := `{
		"customer_id":"01912345-6789-7abc-def0-123456789001",
		"staff_id":"01912345-6789-7abc-def0-123456789002",
		"appointment_date":"2024-12-20",
		"start_time":"10:00","end_time":"11:00","is_walkin":false,
		"services":[]
	}`
	testutil.DoRequest(router, http.MethodPost, "/api/v1/appointments", body)

	// List
	w := testutil.DoRequest(router, http.MethodGet, "/api/v1/appointments", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["success"] != true {
		t.Fatal("expected success true")
	}
}

func TestAppointmentAPI_GetNotFound(t *testing.T) {
	router := setupAppointmentHandler(t)

	w := testutil.DoRequest(router, http.MethodGet, "/api/v1/appointments/01912345-6789-7abc-def0-000000000099", "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAppointmentAPI_UpdateStatus(t *testing.T) {
	router := setupAppointmentHandler(t)

	// Create
	body := `{
		"customer_id":"01912345-6789-7abc-def0-123456789001",
		"staff_id":"01912345-6789-7abc-def0-123456789002",
		"appointment_date":"2024-12-20",
		"start_time":"10:00","end_time":"11:00","is_walkin":false,
		"services":[]
	}`
	w := testutil.DoRequest(router, http.MethodPost, "/api/v1/appointments", body)
	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	data := createResp["data"].(map[string]interface{})
	id := data["id"].(string)

	// Update status
	statusBody := `{"status":"confirmed","changed_by":"admin","note":"Confirmed by phone"}`
	w = testutil.DoRequest(router, http.MethodPut, "/api/v1/appointments/"+id+"/status", statusBody)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAppointmentAPI_Delete(t *testing.T) {
	router := setupAppointmentHandler(t)

	// Create without services
	body := `{
		"customer_id":"01912345-6789-7abc-def0-123456789001",
		"staff_id":"01912345-6789-7abc-def0-123456789002",
		"appointment_date":"2024-12-20",
		"start_time":"10:00","end_time":"11:00","is_walkin":false
	}`
	w := testutil.DoRequest(router, http.MethodPost, "/api/v1/appointments", body)
	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	data := createResp["data"].(map[string]interface{})
	id := data["id"].(string)

	// Delete
	w = testutil.DoRequest(router, http.MethodDelete, "/api/v1/appointments/"+id, "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// Verify deleted
	w = testutil.DoRequest(router, http.MethodGet, "/api/v1/appointments/"+id, "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", w.Code)
	}
}
