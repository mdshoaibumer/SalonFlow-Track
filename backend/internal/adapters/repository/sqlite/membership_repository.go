package sqlite

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// MembershipRepository is the SQLite implementation.
type MembershipRepository struct {
	db  *sql.DB
	log *slog.Logger
}

// NewMembershipRepository creates a new MembershipRepository.
func NewMembershipRepository(db *sql.DB, log *slog.Logger) *MembershipRepository {
	return &MembershipRepository{db: db, log: log}
}

// CreatePlan inserts a plan.
func (r *MembershipRepository) CreatePlan(ctx context.Context, plan *domain.MembershipPlan) error {
	isActive := 0
	if plan.IsActive {
		isActive = 1
	}
	priority := 0
	if plan.PriorityBooking {
		priority = 1
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO membership_plans (id, name, description, plan_type, price, duration_days, max_sessions, discount_percentage, priority_booking, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		plan.ID, plan.Name, plan.Description, plan.PlanType, plan.Price,
		plan.DurationDays, plan.MaxSessions, plan.DiscountPercentage, priority, isActive,
		plan.CreatedAt.Format(time.RFC3339), plan.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("create_plan", err)
	}
	return nil
}

// UpdatePlan updates a plan.
func (r *MembershipRepository) UpdatePlan(ctx context.Context, plan *domain.MembershipPlan) error {
	plan.UpdatedAt = time.Now().UTC()
	isActive := 0
	if plan.IsActive {
		isActive = 1
	}
	priority := 0
	if plan.PriorityBooking {
		priority = 1
	}
	res, err := r.db.ExecContext(ctx, `
		UPDATE membership_plans SET name=?, description=?, plan_type=?, price=?, duration_days=?, max_sessions=?, discount_percentage=?, priority_booking=?, is_active=?, updated_at=? WHERE id=?`,
		plan.Name, plan.Description, plan.PlanType, plan.Price, plan.DurationDays,
		plan.MaxSessions, plan.DiscountPercentage, priority, isActive, plan.UpdatedAt.Format(time.RFC3339), plan.ID)
	if err != nil {
		return apperror.Database("update_plan", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperror.NotFound("membership_plan", plan.ID.String())
	}
	return nil
}

// DeletePlan deletes a plan.
func (r *MembershipRepository) DeletePlan(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM membership_plans WHERE id=?`, id)
	if err != nil {
		return apperror.Database("delete_plan", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperror.NotFound("membership_plan", id.String())
	}
	return nil
}

// GetPlan retrieves a plan by ID.
func (r *MembershipRepository) GetPlan(ctx context.Context, id uuid.UUID) (*domain.MembershipPlan, error) {
	var plan domain.MembershipPlan
	var isActive, priority int
	var createdAt, updatedAt string
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, description, plan_type, price, duration_days, max_sessions, discount_percentage, priority_booking, is_active, created_at, updated_at
		FROM membership_plans WHERE id=?`, id).
		Scan(&plan.ID, &plan.Name, &plan.Description, &plan.PlanType, &plan.Price,
			&plan.DurationDays, &plan.MaxSessions, &plan.DiscountPercentage, &priority, &isActive, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("membership_plan", id.String())
	}
	if err != nil {
		return nil, apperror.Database("get_plan", err)
	}
	plan.IsActive = isActive == 1
	plan.PriorityBooking = priority == 1
	plan.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	plan.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return &plan, nil
}

// ListPlans lists plans.
func (r *MembershipRepository) ListPlans(ctx context.Context, planType string) ([]domain.MembershipPlan, error) {
	query := `SELECT id, name, description, plan_type, price, duration_days, max_sessions, discount_percentage, priority_booking, is_active, created_at, updated_at FROM membership_plans`
	var args []interface{}
	if planType != "" {
		query += ` WHERE plan_type=?`
		args = append(args, planType)
	}
	query += ` ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.Database("list_plans", err)
	}
	defer rows.Close()

	var plans []domain.MembershipPlan
	for rows.Next() {
		var plan domain.MembershipPlan
		var isActive, priority int
		var createdAt, updatedAt string
		if err := rows.Scan(&plan.ID, &plan.Name, &plan.Description, &plan.PlanType, &plan.Price,
			&plan.DurationDays, &plan.MaxSessions, &plan.DiscountPercentage, &priority, &isActive, &createdAt, &updatedAt); err != nil {
			return nil, apperror.Database("list_plans_scan", err)
		}
		plan.IsActive = isActive == 1
		plan.PriorityBooking = priority == 1
		plan.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		plan.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		plans = append(plans, plan)
	}
	return plans, nil
}

// AddPlanServices inserts services for a plan.
func (r *MembershipRepository) AddPlanServices(ctx context.Context, services []domain.PackageService) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return apperror.Database("add_plan_services_begin", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO package_services (id, plan_id, service_id, service_name, sessions_included, created_at) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return apperror.Database("add_plan_services_prepare", err)
	}
	defer stmt.Close()

	for _, s := range services {
		_, err := stmt.ExecContext(ctx, s.ID, s.PlanID, s.ServiceID, s.ServiceName, s.SessionsIncluded, s.CreatedAt.Format(time.RFC3339))
		if err != nil {
			return apperror.Database("add_plan_services_exec", err)
		}
	}
	return tx.Commit()
}

// GetPlanServices retrieves services for a plan.
func (r *MembershipRepository) GetPlanServices(ctx context.Context, planID uuid.UUID) ([]domain.PackageService, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, plan_id, service_id, service_name, sessions_included, created_at FROM package_services WHERE plan_id=?`, planID)
	if err != nil {
		return nil, apperror.Database("get_plan_services", err)
	}
	defer rows.Close()

	var services []domain.PackageService
	for rows.Next() {
		var s domain.PackageService
		var createdAt string
		if err := rows.Scan(&s.ID, &s.PlanID, &s.ServiceID, &s.ServiceName, &s.SessionsIncluded, &createdAt); err != nil {
			return nil, apperror.Database("get_plan_services_scan", err)
		}
		s.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		services = append(services, s)
	}
	return services, nil
}

// DeletePlanServices deletes all services for a plan.
func (r *MembershipRepository) DeletePlanServices(ctx context.Context, planID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM package_services WHERE plan_id=?`, planID)
	if err != nil {
		return apperror.Database("delete_plan_services", err)
	}
	return nil
}

// CreateSubscription inserts a subscription.
func (r *MembershipRepository) CreateSubscription(ctx context.Context, sub *domain.MemberSubscription) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO member_subscriptions (id, customer_id, plan_id, plan_name, start_date, end_date, total_sessions, used_sessions, amount_paid, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		sub.ID, sub.CustomerID, sub.PlanID, sub.PlanName, sub.StartDate, sub.EndDate,
		sub.TotalSessions, sub.UsedSessions, sub.AmountPaid, sub.Status,
		sub.CreatedAt.Format(time.RFC3339), sub.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("create_subscription", err)
	}
	return nil
}

// UpdateSubscription updates a subscription.
func (r *MembershipRepository) UpdateSubscription(ctx context.Context, sub *domain.MemberSubscription) error {
	sub.UpdatedAt = time.Now().UTC()
	res, err := r.db.ExecContext(ctx, `
		UPDATE member_subscriptions SET status=?, used_sessions=?, updated_at=? WHERE id=?`,
		sub.Status, sub.UsedSessions, sub.UpdatedAt.Format(time.RFC3339), sub.ID)
	if err != nil {
		return apperror.Database("update_subscription", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperror.NotFound("member_subscription", sub.ID.String())
	}
	return nil
}

// GetSubscription gets a subscription by ID.
func (r *MembershipRepository) GetSubscription(ctx context.Context, id uuid.UUID) (*domain.MemberSubscription, error) {
	var sub domain.MemberSubscription
	var createdAt, updatedAt string
	err := r.db.QueryRowContext(ctx, `
		SELECT id, customer_id, plan_id, plan_name, start_date, end_date, total_sessions, used_sessions, amount_paid, status, created_at, updated_at
		FROM member_subscriptions WHERE id=?`, id).
		Scan(&sub.ID, &sub.CustomerID, &sub.PlanID, &sub.PlanName, &sub.StartDate, &sub.EndDate,
			&sub.TotalSessions, &sub.UsedSessions, &sub.AmountPaid, &sub.Status, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("member_subscription", id.String())
	}
	if err != nil {
		return nil, apperror.Database("get_subscription", err)
	}
	sub.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	sub.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return &sub, nil
}

// ListSubscriptions lists subscriptions.
func (r *MembershipRepository) ListSubscriptions(ctx context.Context, customerID, status string, limit, offset int) ([]domain.MemberSubscription, int, error) {
	where := "1=1"
	var args []interface{}
	if customerID != "" {
		where += " AND customer_id=?"
		args = append(args, customerID)
	}
	if status != "" {
		where += " AND status=?"
		args = append(args, status)
	}

	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM member_subscriptions WHERE "+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, apperror.Database("list_subscriptions_count", err)
	}

	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, `SELECT id, customer_id, plan_id, plan_name, start_date, end_date, total_sessions, used_sessions, amount_paid, status, created_at, updated_at FROM member_subscriptions WHERE `+where+` ORDER BY created_at DESC LIMIT ? OFFSET ?`, args...)
	if err != nil {
		return nil, 0, apperror.Database("list_subscriptions", err)
	}
	defer rows.Close()

	var subs []domain.MemberSubscription
	for rows.Next() {
		var sub domain.MemberSubscription
		var createdAt, updatedAt string
		if err := rows.Scan(&sub.ID, &sub.CustomerID, &sub.PlanID, &sub.PlanName, &sub.StartDate, &sub.EndDate,
			&sub.TotalSessions, &sub.UsedSessions, &sub.AmountPaid, &sub.Status, &createdAt, &updatedAt); err != nil {
			return nil, 0, apperror.Database("list_subscriptions_scan", err)
		}
		sub.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		sub.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		subs = append(subs, sub)
	}
	return subs, total, nil
}

// IncrementUsedSessions increments the used sessions count.
func (r *MembershipRepository) IncrementUsedSessions(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `UPDATE member_subscriptions SET used_sessions = used_sessions + 1, updated_at=? WHERE id=?`,
		time.Now().UTC().Format(time.RFC3339), id)
	if err != nil {
		return apperror.Database("increment_sessions", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperror.NotFound("member_subscription", id.String())
	}
	return nil
}

// GetStats gets membership statistics.
func (r *MembershipRepository) GetStats(ctx context.Context) (*domain.MembershipStats, error) {
	var stats domain.MembershipStats
	err := r.db.QueryRowContext(ctx, `
		SELECT
			COALESCE(SUM(CASE WHEN status='active' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status='expired' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(amount_paid), 0)
		FROM member_subscriptions`).
		Scan(&stats.ActiveMembers, &stats.ExpiredMembers, &stats.TotalRevenue)
	if err != nil {
		return nil, apperror.Database("get_membership_stats", err)
	}
	return &stats, nil
}
