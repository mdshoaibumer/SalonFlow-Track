package handler_test

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/salonflow/salonflow-track/internal/adapters/handler"
	"github.com/salonflow/salonflow-track/internal/adapters/repository/sqlite"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
	"github.com/salonflow/salonflow-track/internal/testutil"
)

// mockCloudEngine is a test double for CloudBackupEngine.
type mockCloudEngine struct{}

func (m *mockCloudEngine) Upload(_ *domain.CloudBackupConfig, _, _ string, _ bool) error {
	return nil
}
func (m *mockCloudEngine) Download(_ *domain.CloudBackupConfig, _, _ string) error {
	return nil
}
func (m *mockCloudEngine) TestConnection(_ *domain.CloudBackupConfig) error {
	return nil
}

func setupCloudBackupHandler(t *testing.T) chi.Router {
	t.Helper()
	db := testutil.TestDB(t)
	log := testutil.TestLogger(t)
	repo := sqlite.NewCloudBackupRepository(db, log)
	engine := &mockCloudEngine{}
	uc := usecase.NewCloudBackupUseCase(repo, engine, ":memory:", log)
	h := handler.NewCloudBackupHandler(uc)
	r := chi.NewRouter()
	r.Mount("/api/v1/cloud-backup", h.Routes())
	return r
}

func TestCloudBackupAPI_GetConfig(t *testing.T) {
	router := setupCloudBackupHandler(t)

	w := testutil.DoRequest(router, http.MethodGet, "/api/v1/cloud-backup/config", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCloudBackupAPI_SaveConfig(t *testing.T) {
	router := setupCloudBackupHandler(t)

	body := `{
		"provider":"aws_s3",
		"bucket_name":"salon-backups",
		"region":"ap-south-1",
		"access_key":"test-key",
		"encrypt_backups":true,
		"auto_backup":true,
		"auto_backup_interval_hours":12,
		"max_versions":5
	}`
	w := testutil.DoRequest(router, http.MethodPost, "/api/v1/cloud-backup/config", body)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCloudBackupAPI_TestConnection(t *testing.T) {
	router := setupCloudBackupHandler(t)

	// Save config first
	body := `{"provider":"aws_s3","bucket_name":"test-bucket","region":"us-east-1","access_key":"key"}`
	testutil.DoRequest(router, http.MethodPost, "/api/v1/cloud-backup/config", body)

	// Test connection
	w := testutil.DoRequest(router, http.MethodPost, "/api/v1/cloud-backup/test", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCloudBackupAPI_ListHistory(t *testing.T) {
	router := setupCloudBackupHandler(t)

	w := testutil.DoRequest(router, http.MethodGet, "/api/v1/cloud-backup/history", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCloudBackupAPI_Stats(t *testing.T) {
	router := setupCloudBackupHandler(t)

	w := testutil.DoRequest(router, http.MethodGet, "/api/v1/cloud-backup/stats", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}
