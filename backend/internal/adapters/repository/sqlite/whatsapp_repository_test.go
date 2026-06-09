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

func setupWhatsAppTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE whatsapp_templates (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			category TEXT NOT NULL DEFAULT 'general',
			body TEXT NOT NULL,
			variables TEXT NOT NULL DEFAULT '[]',
			is_active INTEGER NOT NULL DEFAULT 1,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE TABLE whatsapp_messages (
			id TEXT PRIMARY KEY,
			template_id TEXT NOT NULL,
			recipient_phone TEXT NOT NULL,
			recipient_name TEXT NOT NULL DEFAULT '',
			message_body TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'queued',
			provider TEXT NOT NULL DEFAULT '',
			provider_message_id TEXT NOT NULL DEFAULT '',
			error_message TEXT NOT NULL DEFAULT '',
			sent_at TEXT NOT NULL DEFAULT '',
			delivered_at TEXT NOT NULL DEFAULT '',
			read_at TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			FOREIGN KEY (template_id) REFERENCES whatsapp_templates(id)
		);
		CREATE TABLE automation_rules (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			trigger_type TEXT NOT NULL,
			template_id TEXT NOT NULL,
			delay_minutes INTEGER NOT NULL DEFAULT 0,
			is_active INTEGER NOT NULL DEFAULT 1,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			FOREIGN KEY (template_id) REFERENCES whatsapp_templates(id)
		);
	`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestWhatsAppRepository_Templates(t *testing.T) {
	db := setupWhatsAppTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewWhatsAppRepository(db, log)
	ctx := context.Background()

	tmpl := &domain.WhatsAppTemplate{
		ID:        uid.New(),
		Name:      "Booking Confirmation",
		Category:  domain.WACategoryAppointment,
		Body:      "Hi {{name}}, your appointment is confirmed for {{date}}",
		Variables: `["name","date"]`,
		IsActive:  true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err := repo.CreateTemplate(ctx, tmpl)
	if err != nil {
		t.Fatalf("CreateTemplate failed: %v", err)
	}

	got, err := repo.GetTemplate(ctx, tmpl.ID)
	if err != nil {
		t.Fatalf("GetTemplate failed: %v", err)
	}
	if got.Name != tmpl.Name {
		t.Errorf("got name %q, want %q", got.Name, tmpl.Name)
	}
}

func TestWhatsAppRepository_ListTemplates(t *testing.T) {
	db := setupWhatsAppTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewWhatsAppRepository(db, log)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		tmpl := &domain.WhatsAppTemplate{
			ID:        uid.New(),
			Name:      "Template",
			Category:  domain.WACategoryGeneral,
			Body:      "Hello",
			Variables: "[]",
			IsActive:  true,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		_ = repo.CreateTemplate(ctx, tmpl)
	}

	list, err := repo.ListTemplates(ctx, "")
	if err != nil {
		t.Fatalf("ListTemplates failed: %v", err)
	}
	if len(list) != 3 {
		t.Errorf("got %d templates, want 3", len(list))
	}
}

func TestWhatsAppRepository_Messages(t *testing.T) {
	db := setupWhatsAppTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewWhatsAppRepository(db, log)
	ctx := context.Background()

	tmplID := uid.New()
	tmpl := &domain.WhatsAppTemplate{
		ID:        tmplID,
		Name:      "Test",
		Category:  domain.WACategoryGeneral,
		Body:      "Hi",
		Variables: "[]",
		IsActive:  true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_ = repo.CreateTemplate(ctx, tmpl)

	msg := &domain.WhatsAppMessage{
		ID:             uid.New(),
		TemplateID:     tmplID.String(),
		RecipientPhone: "+919876543210",
		RecipientName:  "Test User",
		MessageBody:    "Hi",
		Status:         domain.WAStatusSent,
		SentAt:         time.Now().UTC().Format(time.RFC3339),
		CreatedAt:      time.Now().UTC(),
	}
	err := repo.CreateMessage(ctx, msg)
	if err != nil {
		t.Fatalf("CreateMessage failed: %v", err)
	}

	list, total, err := repo.ListMessages(ctx, 10, 0, "")
	if err != nil {
		t.Fatalf("ListMessages failed: %v", err)
	}
	if total != 1 {
		t.Errorf("got total %d, want 1", total)
	}
	if len(list) != 1 {
		t.Errorf("got %d messages, want 1", len(list))
	}
}

func TestWhatsAppRepository_Rules(t *testing.T) {
	db := setupWhatsAppTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewWhatsAppRepository(db, log)
	ctx := context.Background()

	tmplID := uid.New()
	tmpl := &domain.WhatsAppTemplate{
		ID:        tmplID,
		Name:      "Test",
		Category:  domain.WACategoryGeneral,
		Body:      "Hi",
		Variables: "[]",
		IsActive:  true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_ = repo.CreateTemplate(ctx, tmpl)

	rule := &domain.AutomationRule{
		ID:           uid.New(),
		Name:         "Booking Confirmation",
		TriggerType:  domain.WATriggerAppointmentConfirmed,
		TemplateID:   tmplID.String(),
		DelayMinutes: 0,
		IsActive:     true,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	err := repo.CreateRule(ctx, rule)
	if err != nil {
		t.Fatalf("CreateRule failed: %v", err)
	}

	rules, err := repo.ListRules(ctx)
	if err != nil {
		t.Fatalf("ListRules failed: %v", err)
	}
	if len(rules) != 1 {
		t.Errorf("got %d rules, want 1", len(rules))
	}
}

func TestWhatsAppRepository_Stats(t *testing.T) {
	db := setupWhatsAppTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewWhatsAppRepository(db, log)
	ctx := context.Background()

	tmplID := uid.New()
	tmpl := &domain.WhatsAppTemplate{
		ID: tmplID, Name: "Test", Category: domain.WACategoryGeneral, Body: "Hi",
		Variables: "[]", IsActive: true,
		CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
	}
	_ = repo.CreateTemplate(ctx, tmpl)

	for _, status := range []string{domain.WAStatusSent, domain.WAStatusDelivered, domain.WAStatusFailed} {
		msg := &domain.WhatsAppMessage{
			ID: uid.New(), TemplateID: tmplID.String(), RecipientPhone: "+91123",
			RecipientName: "User", MessageBody: "Test", Status: status,
			CreatedAt: time.Now().UTC(),
		}
		_ = repo.CreateMessage(ctx, msg)
	}

	stats, err := repo.GetStats(ctx)
	if err != nil {
		t.Fatalf("GetStats failed: %v", err)
	}
	// "sent" counts toward TotalSent, "delivered" also counts toward TotalSent
	if stats.TotalSent != 2 {
		t.Errorf("got sent %d, want 2", stats.TotalSent)
	}
	if stats.Delivered != 1 {
		t.Errorf("got delivered %d, want 1", stats.Delivered)
	}
	if stats.Failed != 1 {
		t.Errorf("got failed %d, want 1", stats.Failed)
	}
}
