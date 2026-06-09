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

func setupWhatsAppHandler(t *testing.T) chi.Router {
	t.Helper()
	db := testutil.TestDB(t)
	log := testutil.TestLogger(t)
	repo := sqlite.NewWhatsAppRepository(db, log)
	uc := usecase.NewWhatsAppUseCase(repo, log)
	h := handler.NewWhatsAppHandler(uc)
	r := chi.NewRouter()
	r.Mount("/api/v1/whatsapp", h.Routes())
	return r
}

func TestWhatsAppAPI_CreateTemplate(t *testing.T) {
	router := setupWhatsAppHandler(t)

	body := `{
		"name":"Welcome Message",
		"category":"general",
		"body":"Hello {{name}}, welcome to our salon!",
		"variables":"[\"name\"]"
	}`
	w := testutil.DoRequest(router, http.MethodPost, "/api/v1/whatsapp/templates", body)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["success"] != true {
		t.Fatal("expected success true")
	}
}

func TestWhatsAppAPI_CreateTemplate_Validation(t *testing.T) {
	router := setupWhatsAppHandler(t)

	// Empty name
	body := `{"name":"","body":"Hello"}`
	w := testutil.DoRequest(router, http.MethodPost, "/api/v1/whatsapp/templates", body)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestWhatsAppAPI_ListTemplates(t *testing.T) {
	router := setupWhatsAppHandler(t)

	// Create
	body := `{"name":"Reminder","category":"appointment","body":"Hi {{name}}, reminder for {{date}}.","variables":"[\"name\",\"date\"]"}`
	testutil.DoRequest(router, http.MethodPost, "/api/v1/whatsapp/templates", body)

	// List
	w := testutil.DoRequest(router, http.MethodGet, "/api/v1/whatsapp/templates", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["success"] != true {
		t.Fatal("expected success true")
	}
}

func TestWhatsAppAPI_SendMessage(t *testing.T) {
	router := setupWhatsAppHandler(t)

	// Create template first
	tplBody := `{"name":"Booking","category":"appointment","body":"Hi {{name}}, confirmed for {{date}}.","variables":"[\"name\",\"date\"]"}`
	w := testutil.DoRequest(router, http.MethodPost, "/api/v1/whatsapp/templates", tplBody)
	if w.Code != http.StatusCreated {
		t.Fatalf("template create failed: %d: %s", w.Code, w.Body.String())
	}
	var tplResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &tplResp)
	data := tplResp["data"].(map[string]interface{})
	tplID := data["id"].(string)

	// Send message via /send endpoint
	msgBody := `{"template_id":"` + tplID + `","phone":"9876543210","name":"Priya","variables":{"name":"Priya","date":"Dec 20"}}`
	w = testutil.DoRequest(router, http.MethodPost, "/api/v1/whatsapp/send", msgBody)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestWhatsAppAPI_ListMessages(t *testing.T) {
	router := setupWhatsAppHandler(t)

	w := testutil.DoRequest(router, http.MethodGet, "/api/v1/whatsapp/messages", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestWhatsAppAPI_Stats(t *testing.T) {
	router := setupWhatsAppHandler(t)

	w := testutil.DoRequest(router, http.MethodGet, "/api/v1/whatsapp/stats", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestWhatsAppAPI_CreateAutomation(t *testing.T) {
	router := setupWhatsAppHandler(t)

	// Create template first
	tplBody := `{"name":"Auto Reminder","category":"appointment","body":"Hi {{name}}, see you!","variables":"[\"name\"]"}`
	w := testutil.DoRequest(router, http.MethodPost, "/api/v1/whatsapp/templates", tplBody)
	if w.Code != http.StatusCreated {
		t.Fatalf("template create failed: %d: %s", w.Code, w.Body.String())
	}
	var tplResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &tplResp)
	data := tplResp["data"].(map[string]interface{})
	tplID := data["id"].(string)

	// Create automation rule
	body := `{"name":"Pre-visit reminder","trigger_type":"appointment_reminder","template_id":"` + tplID + `","delay_minutes":60}`
	w = testutil.DoRequest(router, http.MethodPost, "/api/v1/whatsapp/rules", body)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}
