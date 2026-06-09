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

func setupMembershipTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE membership_plans (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			plan_type TEXT NOT NULL DEFAULT 'package',
			price REAL NOT NULL DEFAULT 0,
			duration_days INTEGER NOT NULL DEFAULT 30,
			max_sessions INTEGER NOT NULL DEFAULT 0,
			discount_percentage REAL NOT NULL DEFAULT 0,
			priority_booking INTEGER NOT NULL DEFAULT 0,
			is_active INTEGER NOT NULL DEFAULT 1,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE TABLE package_services (
			id TEXT PRIMARY KEY,
			plan_id TEXT NOT NULL,
			service_id TEXT NOT NULL,
			service_name TEXT NOT NULL DEFAULT '',
			sessions_included INTEGER NOT NULL DEFAULT 1,
			created_at TEXT NOT NULL,
			FOREIGN KEY (plan_id) REFERENCES membership_plans(id)
		);
		CREATE TABLE member_subscriptions (
			id TEXT PRIMARY KEY,
			customer_id TEXT NOT NULL,
			plan_id TEXT NOT NULL,
			plan_name TEXT NOT NULL DEFAULT '',
			start_date TEXT NOT NULL,
			end_date TEXT NOT NULL,
			total_sessions INTEGER NOT NULL DEFAULT 0,
			used_sessions INTEGER NOT NULL DEFAULT 0,
			amount_paid REAL NOT NULL DEFAULT 0,
			status TEXT NOT NULL DEFAULT 'active',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			FOREIGN KEY (plan_id) REFERENCES membership_plans(id)
		);
	`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestMembershipRepository_Plans(t *testing.T) {
	db := setupMembershipTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewMembershipRepository(db, log)
	ctx := context.Background()

	plan := &domain.MembershipPlan{
		ID:           uid.New(),
		Name:         "Gold Package",
		PlanType:     domain.PlanTypePackage,
		Price:        5000,
		DurationDays: 90,
		MaxSessions:  12,
		Description:  "12 sessions in 90 days",
		IsActive:     true,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	err := repo.CreatePlan(ctx, plan)
	if err != nil {
		t.Fatalf("CreatePlan failed: %v", err)
	}

	got, err := repo.GetPlan(ctx, plan.ID)
	if err != nil {
		t.Fatalf("GetPlan failed: %v", err)
	}
	if got.Name != plan.Name {
		t.Errorf("got name %q, want %q", got.Name, plan.Name)
	}
	if got.Price != 5000 {
		t.Errorf("got price %v, want 5000", got.Price)
	}
}

func TestMembershipRepository_ListPlans(t *testing.T) {
	db := setupMembershipTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewMembershipRepository(db, log)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		plan := &domain.MembershipPlan{
			ID: uid.New(), Name: "Plan", PlanType: domain.PlanTypePackage,
			Price: 1000, DurationDays: 30, IsActive: true,
			CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
		}
		_ = repo.CreatePlan(ctx, plan)
	}

	list, err := repo.ListPlans(ctx, "")
	if err != nil {
		t.Fatalf("ListPlans failed: %v", err)
	}
	if len(list) != 3 {
		t.Errorf("got %d plans, want 3", len(list))
	}
}

func TestMembershipRepository_Subscriptions(t *testing.T) {
	db := setupMembershipTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewMembershipRepository(db, log)
	ctx := context.Background()

	planID := uid.New()
	plan := &domain.MembershipPlan{
		ID: planID, Name: "Test Plan", PlanType: domain.PlanTypePackage,
		Price: 2000, DurationDays: 30, MaxSessions: 10, IsActive: true,
		CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
	}
	_ = repo.CreatePlan(ctx, plan)

	now := time.Now().UTC()
	sub := &domain.MemberSubscription{
		ID:            uid.New(),
		PlanID:        planID,
		PlanName:      "Test Plan",
		CustomerID:    uid.New().String(),
		StartDate:     now.Format("2006-01-02"),
		EndDate:       now.Add(30 * 24 * time.Hour).Format("2006-01-02"),
		TotalSessions: 10,
		UsedSessions:  0,
		AmountPaid:    2000,
		Status:        domain.SubscriptionActive,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	err := repo.CreateSubscription(ctx, sub)
	if err != nil {
		t.Fatalf("CreateSubscription failed: %v", err)
	}

	subs, total, err := repo.ListSubscriptions(ctx, "", "", 10, 0)
	if err != nil {
		t.Fatalf("ListSubscriptions failed: %v", err)
	}
	if total != 1 {
		t.Errorf("got total %d, want 1", total)
	}
	if len(subs) != 1 {
		t.Errorf("got %d subs, want 1", len(subs))
	}
}

func TestMembershipRepository_IncrementSessions(t *testing.T) {
	db := setupMembershipTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewMembershipRepository(db, log)
	ctx := context.Background()

	planID := uid.New()
	plan := &domain.MembershipPlan{
		ID: planID, Name: "Test", PlanType: domain.PlanTypePackage,
		Price: 1000, DurationDays: 30, MaxSessions: 5, IsActive: true,
		CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
	}
	_ = repo.CreatePlan(ctx, plan)

	now := time.Now().UTC()
	sub := &domain.MemberSubscription{
		ID:            uid.New(),
		PlanID:        planID,
		PlanName:      "Test",
		CustomerID:    uid.New().String(),
		StartDate:     now.Format("2006-01-02"),
		EndDate:       now.Add(30 * 24 * time.Hour).Format("2006-01-02"),
		TotalSessions: 5,
		UsedSessions:  0,
		AmountPaid:    1000,
		Status:        domain.SubscriptionActive,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	_ = repo.CreateSubscription(ctx, sub)

	err := repo.IncrementUsedSessions(ctx, sub.ID)
	if err != nil {
		t.Fatalf("IncrementUsedSessions failed: %v", err)
	}

	got, err := repo.GetSubscription(ctx, sub.ID)
	if err != nil {
		t.Fatalf("GetSubscription failed: %v", err)
	}
	if got.UsedSessions != 1 {
		t.Errorf("got used %d, want 1", got.UsedSessions)
	}
}

func TestMembershipRepository_Stats(t *testing.T) {
	db := setupMembershipTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewMembershipRepository(db, log)
	ctx := context.Background()

	planID := uid.New()
	plan := &domain.MembershipPlan{
		ID: planID, Name: "Gold", PlanType: domain.PlanTypePackage,
		Price: 3000, DurationDays: 30, MaxSessions: 10, IsActive: true,
		CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
	}
	_ = repo.CreatePlan(ctx, plan)

	now := time.Now().UTC()
	sub := &domain.MemberSubscription{
		ID: uid.New(), PlanID: planID, PlanName: "Gold",
		CustomerID:    uid.New().String(),
		StartDate:     now.Format("2006-01-02"),
		EndDate:       now.Add(30 * 24 * time.Hour).Format("2006-01-02"),
		TotalSessions: 10, UsedSessions: 0,
		Status: domain.SubscriptionActive, AmountPaid: 3000,
		CreatedAt: now, UpdatedAt: now,
	}
	_ = repo.CreateSubscription(ctx, sub)

	stats, err := repo.GetStats(ctx)
	if err != nil {
		t.Fatalf("GetStats failed: %v", err)
	}
	if stats.ActiveMembers != 1 {
		t.Errorf("got active %d, want 1", stats.ActiveMembers)
	}
	if stats.TotalRevenue != 3000 {
		t.Errorf("got revenue %v, want 3000", stats.TotalRevenue)
	}
}
