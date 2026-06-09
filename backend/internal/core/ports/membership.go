package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// MembershipRepository manages membership data.
type MembershipRepository interface {
	// Plans
	CreatePlan(ctx context.Context, plan *domain.MembershipPlan) error
	UpdatePlan(ctx context.Context, plan *domain.MembershipPlan) error
	DeletePlan(ctx context.Context, id uuid.UUID) error
	GetPlan(ctx context.Context, id uuid.UUID) (*domain.MembershipPlan, error)
	ListPlans(ctx context.Context, planType string) ([]domain.MembershipPlan, error)

	// Package Services
	AddPlanServices(ctx context.Context, services []domain.PackageService) error
	GetPlanServices(ctx context.Context, planID uuid.UUID) ([]domain.PackageService, error)
	DeletePlanServices(ctx context.Context, planID uuid.UUID) error

	// Subscriptions
	CreateSubscription(ctx context.Context, sub *domain.MemberSubscription) error
	UpdateSubscription(ctx context.Context, sub *domain.MemberSubscription) error
	GetSubscription(ctx context.Context, id uuid.UUID) (*domain.MemberSubscription, error)
	ListSubscriptions(ctx context.Context, customerID, status string, limit, offset int) ([]domain.MemberSubscription, int, error)
	IncrementUsedSessions(ctx context.Context, id uuid.UUID) error

	// Stats
	GetStats(ctx context.Context) (*domain.MembershipStats, error)
}
