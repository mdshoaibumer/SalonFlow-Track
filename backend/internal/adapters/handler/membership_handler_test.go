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

func setupMembershipHandler(t *testing.T) chi.Router {
	t.Helper()
	db := testutil.TestDB(t)
	log := testutil.TestLogger(t)
	repo := sqlite.NewMembershipRepository(db, log)
	uc := usecase.NewMembershipUseCase(repo, log)
	h := handler.NewMembershipHandler(uc)
	r := chi.NewRouter()
	r.Mount("/api/v1/memberships", h.Routes())
	return r
}

func TestMembershipAPI_CreatePlan(t *testing.T) {
	router := setupMembershipHandler(t)

	body := `{
		"name":"Gold Package",
		"plan_type":"package",
		"price":5000,
		"duration_days":90,
		"max_sessions":12,
		"discount_percentage":10,
		"priority_booking":true,
		"services":[{"service_id":"svc1","service_name":"Haircut","sessions_included":4}]
	}`
	w := testutil.DoRequest(router, http.MethodPost, "/api/v1/memberships/plans", body)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["success"] != true {
		t.Fatal("expected success true")
	}
}

func TestMembershipAPI_CreatePlan_Validation(t *testing.T) {
	router := setupMembershipHandler(t)

	// Missing name
	body := `{"plan_type":"package","price":5000}`
	w := testutil.DoRequest(router, http.MethodPost, "/api/v1/memberships/plans", body)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestMembershipAPI_ListPlans(t *testing.T) {
	router := setupMembershipHandler(t)

	// Create
	body := `{"name":"Silver","plan_type":"membership","price":3000,"duration_days":30,"max_sessions":8,"services":[]}`
	testutil.DoRequest(router, http.MethodPost, "/api/v1/memberships/plans", body)

	// List
	w := testutil.DoRequest(router, http.MethodGet, "/api/v1/memberships/plans", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["success"] != true {
		t.Fatal("expected success true")
	}
}

func TestMembershipAPI_GetPlan_NotFound(t *testing.T) {
	router := setupMembershipHandler(t)

	w := testutil.DoRequest(router, http.MethodGet, "/api/v1/memberships/plans/01912345-6789-7abc-def0-000000000099", "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestMembershipAPI_SellPlan(t *testing.T) {
	router := setupMembershipHandler(t)

	// Create a plan first
	planBody := `{"name":"Premium","plan_type":"package","price":8000,"duration_days":60,"max_sessions":10,"services":[]}`
	w := testutil.DoRequest(router, http.MethodPost, "/api/v1/memberships/plans", planBody)
	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	data := createResp["data"].(map[string]interface{})
	planID := data["id"].(string)

	// Sell
	sellBody := `{"customer_id":"01912345-6789-7abc-def0-123456789001","plan_id":"` + planID + `","amount_paid":8000}`
	w = testutil.DoRequest(router, http.MethodPost, "/api/v1/memberships/sell", sellBody)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestMembershipAPI_Stats(t *testing.T) {
	router := setupMembershipHandler(t)

	w := testutil.DoRequest(router, http.MethodGet, "/api/v1/memberships/stats", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestMembershipAPI_DeletePlan(t *testing.T) {
	router := setupMembershipHandler(t)

	// Create
	body := `{"name":"ToDelete","plan_type":"package","price":1000,"duration_days":30,"max_sessions":5,"services":[]}`
	w := testutil.DoRequest(router, http.MethodPost, "/api/v1/memberships/plans", body)
	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	data := createResp["data"].(map[string]interface{})
	planID := data["id"].(string)

	// Delete
	w = testutil.DoRequest(router, http.MethodDelete, "/api/v1/memberships/plans/"+planID, "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// Verify
	w = testutil.DoRequest(router, http.MethodGet, "/api/v1/memberships/plans/"+planID, "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", w.Code)
	}
}
