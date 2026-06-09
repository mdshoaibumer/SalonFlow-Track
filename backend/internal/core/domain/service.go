package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/pkg/uid"
)

// Service status constants.
const (
	ServiceStatusActive   = "active"
	ServiceStatusInactive = "inactive"
)

// Service category constants.
const (
	CategoryHair      = "hair"
	CategoryFacial    = "facial"
	CategorySkin      = "skin"
	CategorySpa       = "spa"
	CategoryMassage   = "massage"
	CategoryColoring  = "coloring"
	CategoryTreatment = "treatment"
	CategoryOther     = "other"
)

// Commission type constants.
const (
	CommissionTypeFixed      = "fixed"
	CommissionTypePercentage = "percentage"
)

var validCategories = map[string]bool{
	CategoryHair:      true,
	CategoryFacial:    true,
	CategorySkin:      true,
	CategorySpa:       true,
	CategoryMassage:   true,
	CategoryColoring:  true,
	CategoryTreatment: true,
	CategoryOther:     true,
}

var validCommissionTypes = map[string]bool{
	CommissionTypeFixed:      true,
	CommissionTypePercentage: true,
}

// Service represents a salon service offering.
type Service struct {
	ID              uuid.UUID `json:"id"`
	ServiceCode     string    `json:"service_code"`
	Name            string    `json:"name"`
	Category        string    `json:"category"`
	Description     string    `json:"description"`
	DurationMinutes int       `json:"duration_minutes"`
	Price           float64   `json:"price"`
	CostPrice       float64   `json:"cost_price"`
	CommissionType  string    `json:"commission_type"`
	CommissionValue float64   `json:"commission_value"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// NewService creates a new Service with a UUIDv7 and auto-generated code.
func NewService(name, category string, durationMinutes int, price float64) *Service {
	now := time.Now().UTC()
	return &Service{
		ID:              uid.New(),
		ServiceCode:     generateServiceCode(),
		Name:            strings.TrimSpace(name),
		Category:        category,
		DurationMinutes: durationMinutes,
		Price:           price,
		CommissionType:  CommissionTypePercentage,
		CommissionValue: 0,
		Status:          ServiceStatusActive,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// generateServiceCode creates a unique service code like "SVC-XXXXXXXX".
func generateServiceCode() string {
	id := uid.New()
	s := strings.ReplaceAll(id.String(), "-", "")
	short := strings.ToUpper(s[len(s)-8:])
	return fmt.Sprintf("SVC-%s", short)
}

// Validate checks service business rules.
func (s *Service) Validate() error {
	if strings.TrimSpace(s.Name) == "" {
		return ErrServiceNameRequired
	}
	if !validCategories[s.Category] {
		return ErrServiceInvalidCategory
	}
	if s.DurationMinutes <= 0 {
		return ErrServiceInvalidDuration
	}
	if s.Price <= 0 {
		return ErrServiceInvalidPrice
	}
	if s.CostPrice < 0 {
		return ErrServiceInvalidCostPrice
	}
	if !validCommissionTypes[s.CommissionType] {
		return ErrServiceInvalidCommissionType
	}
	if s.CommissionValue < 0 {
		return ErrServiceInvalidCommissionValue
	}
	if s.CommissionType == CommissionTypePercentage && s.CommissionValue > 100 {
		return ErrServiceInvalidCommissionValue
	}
	return nil
}

// Deactivate sets the service status to inactive.
func (s *Service) Deactivate() {
	s.Status = ServiceStatusInactive
	s.UpdatedAt = time.Now().UTC()
}

// Activate sets the service status to active.
func (s *Service) Activate() {
	s.Status = ServiceStatusActive
	s.UpdatedAt = time.Now().UTC()
}

// IsActive returns true if the service is active.
func (s *Service) IsActive() bool {
	return s.Status == ServiceStatusActive
}
