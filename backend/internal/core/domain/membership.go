package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/pkg/uid"
)

// Plan types.
const (
	PlanTypePackage    = "package"
	PlanTypeMembership = "membership"
)

// Subscription statuses.
const (
	SubscriptionActive    = "active"
	SubscriptionExpired   = "expired"
	SubscriptionCancelled = "cancelled"
	SubscriptionPaused    = "paused"
)

// MembershipPlan is a package or membership plan.
type MembershipPlan struct {
	ID                 uuid.UUID        `json:"id"`
	Name               string           `json:"name"`
	Description        string           `json:"description"`
	PlanType           string           `json:"plan_type"`
	Price              float64          `json:"price"`
	DurationDays       int              `json:"duration_days"`
	MaxSessions        int              `json:"max_sessions"`
	DiscountPercentage float64          `json:"discount_percentage"`
	PriorityBooking    bool             `json:"priority_booking"`
	IsActive           bool             `json:"is_active"`
	Services           []PackageService `json:"services,omitempty"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
}

// NewMembershipPlan creates a new plan.
func NewMembershipPlan(name, planType string, price float64, durationDays, maxSessions int) *MembershipPlan {
	now := time.Now().UTC()
	return &MembershipPlan{
		ID:           uid.New(),
		Name:         name,
		PlanType:     planType,
		Price:        price,
		DurationDays: durationDays,
		MaxSessions:  maxSessions,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// PackageService links a service to a plan.
type PackageService struct {
	ID               uuid.UUID `json:"id"`
	PlanID           uuid.UUID `json:"plan_id"`
	ServiceID        string    `json:"service_id"`
	ServiceName      string    `json:"service_name"`
	SessionsIncluded int       `json:"sessions_included"`
	CreatedAt        time.Time `json:"created_at"`
}

// NewPackageService creates a new package service.
func NewPackageService(planID uuid.UUID, serviceID, serviceName string, sessions int) *PackageService {
	return &PackageService{
		ID:               uid.New(),
		PlanID:           planID,
		ServiceID:        serviceID,
		ServiceName:      serviceName,
		SessionsIncluded: sessions,
		CreatedAt:        time.Now().UTC(),
	}
}

// MemberSubscription is a customer's active subscription.
type MemberSubscription struct {
	ID            uuid.UUID `json:"id"`
	CustomerID    string    `json:"customer_id"`
	PlanID        uuid.UUID `json:"plan_id"`
	PlanName      string    `json:"plan_name"`
	StartDate     string    `json:"start_date"`
	EndDate       string    `json:"end_date"`
	TotalSessions int       `json:"total_sessions"`
	UsedSessions  int       `json:"used_sessions"`
	AmountPaid    float64   `json:"amount_paid"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// NewMemberSubscription creates a new subscription.
func NewMemberSubscription(customerID string, plan *MembershipPlan, amountPaid float64) *MemberSubscription {
	now := time.Now().UTC()
	startDate := now.Format("2006-01-02")
	endDate := now.AddDate(0, 0, plan.DurationDays).Format("2006-01-02")
	return &MemberSubscription{
		ID:            uid.New(),
		CustomerID:    customerID,
		PlanID:        plan.ID,
		PlanName:      plan.Name,
		StartDate:     startDate,
		EndDate:       endDate,
		TotalSessions: plan.MaxSessions,
		UsedSessions:  0,
		AmountPaid:    amountPaid,
		Status:        SubscriptionActive,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// MembershipStats holds membership statistics.
type MembershipStats struct {
	ActiveMembers  int     `json:"active_members"`
	ExpiredMembers int     `json:"expired_members"`
	TotalRevenue   float64 `json:"total_revenue"`
}
