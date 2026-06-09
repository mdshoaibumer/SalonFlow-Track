package usecase

import (
	"context"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// ServiceUseCase handles service business logic.
type ServiceUseCase struct {
	repo ports.ServiceRepository
	log  *slog.Logger
}

// NewServiceUseCase creates a new ServiceUseCase.
func NewServiceUseCase(repo ports.ServiceRepository, log *slog.Logger) *ServiceUseCase {
	return &ServiceUseCase{repo: repo, log: log}
}

// CreateServiceInput is the input DTO for creating a service.
type CreateServiceInput struct {
	Name            string  `json:"name"`
	Category        string  `json:"category"`
	Description     string  `json:"description"`
	DurationMinutes int     `json:"duration_minutes"`
	Price           float64 `json:"price"`
	CostPrice       float64 `json:"cost_price"`
	CommissionType  string  `json:"commission_type"`
	CommissionValue float64 `json:"commission_value"`
}

// UpdateServiceInput is the input DTO for updating a service.
type UpdateServiceInput struct {
	Name            string  `json:"name"`
	Category        string  `json:"category"`
	Description     string  `json:"description"`
	DurationMinutes int     `json:"duration_minutes"`
	Price           float64 `json:"price"`
	CostPrice       float64 `json:"cost_price"`
	CommissionType  string  `json:"commission_type"`
	CommissionValue float64 `json:"commission_value"`
	Status          string  `json:"status"`
}

// ListServiceInput is the input DTO for listing services.
type ListServiceInput struct {
	Search   string `json:"search"`
	Status   string `json:"status"`
	Category string `json:"category"`
	Page     int    `json:"page"`
	PerPage  int    `json:"per_page"`
}

// ListServiceOutput is the output DTO for listing services.
type ListServiceOutput struct {
	Services   []domain.Service `json:"services"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	PerPage    int              `json:"per_page"`
	TotalPages int              `json:"total_pages"`
}

// ServiceStats holds service count statistics.
type ServiceStats struct {
	Total    int     `json:"total"`
	Active   int     `json:"active"`
	Inactive int     `json:"inactive"`
	AvgPrice float64 `json:"avg_price"`
}

// Create creates a new service.
func (uc *ServiceUseCase) Create(ctx context.Context, input CreateServiceInput) (*domain.Service, error) {
	svc := domain.NewService(input.Name, input.Category, input.DurationMinutes, input.Price)
	svc.Description = strings.TrimSpace(input.Description)
	svc.CostPrice = input.CostPrice

	if input.CommissionType != "" {
		svc.CommissionType = input.CommissionType
	}
	svc.CommissionValue = input.CommissionValue

	if err := svc.Validate(); err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: err.Error()}
	}

	if err := uc.repo.Create(ctx, svc); err != nil {
		return nil, err
	}

	uc.log.Info("service created", "id", svc.ID, "name", svc.Name)
	return svc, nil
}

// GetByID retrieves a service by ID.
func (uc *ServiceUseCase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Service, error) {
	return uc.repo.GetByID(ctx, id)
}

// List returns paginated service list.
func (uc *ServiceUseCase) List(ctx context.Context, input ListServiceInput) (*ListServiceOutput, error) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PerPage <= 0 {
		input.PerPage = 20
	}
	if input.PerPage > 100 {
		input.PerPage = 100
	}

	offset := (input.Page - 1) * input.PerPage

	filter := ports.ServiceFilter{
		Status:   input.Status,
		Category: input.Category,
		Search:   input.Search,
		Limit:    input.PerPage,
		Offset:   offset,
	}

	services, total, err := uc.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := total / input.PerPage
	if total%input.PerPage > 0 {
		totalPages++
	}

	return &ListServiceOutput{
		Services:   services,
		Total:      total,
		Page:       input.Page,
		PerPage:    input.PerPage,
		TotalPages: totalPages,
	}, nil
}

// Update updates an existing service.
func (uc *ServiceUseCase) Update(ctx context.Context, id uuid.UUID, input UpdateServiceInput) (*domain.Service, error) {
	svc, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	svc.Name = strings.TrimSpace(input.Name)
	svc.Category = input.Category
	svc.Description = strings.TrimSpace(input.Description)
	svc.DurationMinutes = input.DurationMinutes
	svc.Price = input.Price
	svc.CostPrice = input.CostPrice
	svc.CommissionType = input.CommissionType
	svc.CommissionValue = input.CommissionValue

	if input.Status != "" {
		svc.Status = input.Status
	}

	if err := svc.Validate(); err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: err.Error()}
	}

	if err := uc.repo.Update(ctx, svc); err != nil {
		return nil, err
	}

	uc.log.Info("service updated", "id", svc.ID, "name", svc.Name)
	return svc, nil
}

// Delete removes a service.
func (uc *ServiceUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return err
	}
	uc.log.Info("service deleted", "id", id)
	return nil
}

// Stats returns service statistics.
func (uc *ServiceUseCase) Stats(ctx context.Context) (*ServiceStats, error) {
	total, active, inactive, err := uc.repo.CountByStatus(ctx)
	if err != nil {
		return nil, err
	}

	// Calculate average price from active services
	var avgPrice float64
	if active > 0 {
		services, _, err := uc.repo.List(ctx, ports.ServiceFilter{Status: "active", Limit: 1000})
		if err == nil && len(services) > 0 {
			var sum float64
			for _, s := range services {
				sum += s.Price
			}
			avgPrice = sum / float64(len(services))
		}
	}

	return &ServiceStats{
		Total:    total,
		Active:   active,
		Inactive: inactive,
		AvgPrice: avgPrice,
	}, nil
}
