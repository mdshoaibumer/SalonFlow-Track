package usecase

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// MembershipUseCase handles membership business logic.
type MembershipUseCase struct {
	repo ports.MembershipRepository
	log  *slog.Logger
}

// NewMembershipUseCase creates a new MembershipUseCase.
func NewMembershipUseCase(repo ports.MembershipRepository, log *slog.Logger) *MembershipUseCase {
	return &MembershipUseCase{repo: repo, log: log}
}

// CreatePlan creates a plan with services.
func (uc *MembershipUseCase) CreatePlan(ctx context.Context, plan *domain.MembershipPlan, services []domain.PackageService) error {
	if plan.Name == "" {
		return apperror.Validation("name", "Name is required")
	}
	if plan.Price <= 0 {
		return apperror.Validation("price", "Price must be positive")
	}

	if err := uc.repo.CreatePlan(ctx, plan); err != nil {
		return err
	}

	if len(services) > 0 {
		for i := range services {
			services[i].PlanID = plan.ID
			if services[i].ID == uuid.Nil {
				services[i].ID = domain.NewPackageService(plan.ID, "", "", 0).ID
			}
		}
		return uc.repo.AddPlanServices(ctx, services)
	}
	return nil
}

// UpdatePlan updates a plan.
func (uc *MembershipUseCase) UpdatePlan(ctx context.Context, plan *domain.MembershipPlan, services []domain.PackageService) error {
	if err := uc.repo.UpdatePlan(ctx, plan); err != nil {
		return err
	}
	if len(services) > 0 {
		_ = uc.repo.DeletePlanServices(ctx, plan.ID)
		for i := range services {
			services[i].PlanID = plan.ID
			if services[i].ID == uuid.Nil {
				services[i].ID = domain.NewPackageService(plan.ID, "", "", 0).ID
			}
		}
		return uc.repo.AddPlanServices(ctx, services)
	}
	return nil
}

// DeletePlan deletes a plan.
func (uc *MembershipUseCase) DeletePlan(ctx context.Context, id uuid.UUID) error {
	return uc.repo.DeletePlan(ctx, id)
}

// GetPlan gets a plan with services.
func (uc *MembershipUseCase) GetPlan(ctx context.Context, id uuid.UUID) (*domain.MembershipPlan, error) {
	plan, err := uc.repo.GetPlan(ctx, id)
	if err != nil {
		return nil, err
	}
	services, _ := uc.repo.GetPlanServices(ctx, id)
	plan.Services = services
	return plan, nil
}

// ListPlans lists plans.
func (uc *MembershipUseCase) ListPlans(ctx context.Context, planType string) ([]domain.MembershipPlan, error) {
	return uc.repo.ListPlans(ctx, planType)
}

// SellPlan creates a subscription for a customer.
func (uc *MembershipUseCase) SellPlan(ctx context.Context, customerID string, planID uuid.UUID, amountPaid float64) (*domain.MemberSubscription, error) {
	if customerID == "" {
		return nil, apperror.Validation("customer_id", "Customer is required")
	}

	plan, err := uc.repo.GetPlan(ctx, planID)
	if err != nil {
		return nil, err
	}

	sub := domain.NewMemberSubscription(customerID, plan, amountPaid)
	if err := uc.repo.CreateSubscription(ctx, sub); err != nil {
		return nil, err
	}
	return sub, nil
}

// UseSession increments used sessions.
func (uc *MembershipUseCase) UseSession(ctx context.Context, subscriptionID uuid.UUID) error {
	sub, err := uc.repo.GetSubscription(ctx, subscriptionID)
	if err != nil {
		return err
	}
	if sub.Status != domain.SubscriptionActive {
		return apperror.Business("INACTIVE", "Subscription is not active")
	}
	if sub.TotalSessions > 0 && sub.UsedSessions >= sub.TotalSessions {
		return apperror.Business("EXHAUSTED", "All sessions have been used")
	}
	return uc.repo.IncrementUsedSessions(ctx, subscriptionID)
}

// ListSubscriptions lists subscriptions.
func (uc *MembershipUseCase) ListSubscriptions(ctx context.Context, customerID, status string, limit, offset int) ([]domain.MemberSubscription, int, error) {
	return uc.repo.ListSubscriptions(ctx, customerID, status, limit, offset)
}

// GetStats gets membership stats.
func (uc *MembershipUseCase) GetStats(ctx context.Context) (*domain.MembershipStats, error) {
	return uc.repo.GetStats(ctx)
}
