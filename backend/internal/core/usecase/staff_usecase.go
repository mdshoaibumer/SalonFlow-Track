package usecase

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// StaffUseCase handles staff business logic.
type StaffUseCase struct {
	repo ports.StaffRepository
	log  *slog.Logger
}

// NewStaffUseCase creates a new StaffUseCase.
func NewStaffUseCase(repo ports.StaffRepository, log *slog.Logger) *StaffUseCase {
	return &StaffUseCase{repo: repo, log: log}
}

// CreateStaffInput is the input DTO for creating a staff member.
type CreateStaffInput struct {
	FullName             string  `json:"full_name"`
	Phone                string  `json:"phone"`
	Email                string  `json:"email"`
	Gender               string  `json:"gender"`
	Designation          string  `json:"designation"`
	JoiningDate          string  `json:"joining_date"`
	BaseSalary           float64 `json:"base_salary"`
	CommissionPercentage float64 `json:"commission_percentage"`
}

// UpdateStaffInput is the input DTO for updating a staff member.
type UpdateStaffInput struct {
	FullName             string  `json:"full_name"`
	Phone                string  `json:"phone"`
	Email                string  `json:"email"`
	Gender               string  `json:"gender"`
	Designation          string  `json:"designation"`
	JoiningDate          string  `json:"joining_date"`
	BaseSalary           float64 `json:"base_salary"`
	CommissionPercentage float64 `json:"commission_percentage"`
	Status               string  `json:"status"`
}

// ListStaffInput is the input DTO for listing staff.
type ListStaffInput struct {
	Search      string `json:"search"`
	Status      string `json:"status"`
	Designation string `json:"designation"`
	Page        int    `json:"page"`
	PerPage     int    `json:"per_page"`
}

// ListStaffOutput is the output DTO for listing staff.
type ListStaffOutput struct {
	Staff      []domain.Staff `json:"staff"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
	TotalPages int            `json:"total_pages"`
}

// StaffStats holds staff count statistics.
type StaffStats struct {
	Total    int `json:"total"`
	Active   int `json:"active"`
	Inactive int `json:"inactive"`
}

// Create creates a new staff member.
func (uc *StaffUseCase) Create(ctx context.Context, input CreateStaffInput) (*domain.Staff, error) {
	staff := domain.NewStaff(input.FullName, input.Phone, input.Designation, input.BaseSalary, input.CommissionPercentage)
	staff.Email = strings.TrimSpace(input.Email)

	if input.Gender != "" {
		staff.Gender = input.Gender
	}
	if input.JoiningDate != "" {
		t, err := time.Parse("2006-01-02", input.JoiningDate)
		if err == nil {
			staff.JoiningDate = t
		}
	}

	if err := staff.Validate(); err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: err.Error()}
	}

	// Check phone uniqueness
	existing, err := uc.repo.GetByPhone(ctx, staff.Phone)
	if err == nil && existing != nil {
		return nil, &apperror.Error{Kind: apperror.KindConflict, Message: domain.ErrStaffPhoneDuplicate.Error(), Field: "phone"}
	}

	if err := uc.repo.Create(ctx, staff); err != nil {
		return nil, err
	}

	uc.log.Info("staff created", "id", staff.ID, "name", staff.FullName)
	return staff, nil
}

// GetByID retrieves a staff member by ID.
func (uc *StaffUseCase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Staff, error) {
	return uc.repo.GetByID(ctx, id)
}

// List returns paginated staff list.
func (uc *StaffUseCase) List(ctx context.Context, input ListStaffInput) (*ListStaffOutput, error) {
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

	filter := ports.StaffFilter{
		Status:      input.Status,
		Designation: input.Designation,
		Search:      input.Search,
		Limit:       input.PerPage,
		Offset:      offset,
	}

	staff, total, err := uc.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := total / input.PerPage
	if total%input.PerPage > 0 {
		totalPages++
	}

	return &ListStaffOutput{
		Staff:      staff,
		Total:      total,
		Page:       input.Page,
		PerPage:    input.PerPage,
		TotalPages: totalPages,
	}, nil
}

// Update updates an existing staff member.
func (uc *StaffUseCase) Update(ctx context.Context, id uuid.UUID, input UpdateStaffInput) (*domain.Staff, error) {
	staff, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	staff.FullName = strings.TrimSpace(input.FullName)
	staff.Phone = strings.TrimSpace(input.Phone)
	staff.Email = strings.TrimSpace(input.Email)
	staff.Gender = input.Gender
	staff.Designation = input.Designation
	staff.BaseSalary = input.BaseSalary
	staff.CommissionPercentage = input.CommissionPercentage

	if input.Status != "" {
		staff.Status = input.Status
	}
	if input.JoiningDate != "" {
		t, err := time.Parse("2006-01-02", input.JoiningDate)
		if err == nil {
			staff.JoiningDate = t
		}
	}

	if err := staff.Validate(); err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: err.Error()}
	}

	// Check phone uniqueness (exclude self)
	existing, err := uc.repo.GetByPhone(ctx, staff.Phone)
	if err == nil && existing != nil && existing.ID != staff.ID {
		return nil, &apperror.Error{Kind: apperror.KindConflict, Message: domain.ErrStaffPhoneDuplicate.Error(), Field: "phone"}
	}

	if err := uc.repo.Update(ctx, staff); err != nil {
		return nil, err
	}

	uc.log.Info("staff updated", "id", staff.ID, "name", staff.FullName)
	return staff, nil
}

// Delete removes a staff member.
func (uc *StaffUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return err
	}
	uc.log.Info("staff deleted", "id", id)
	return nil
}

// Stats returns staff count statistics.
func (uc *StaffUseCase) Stats(ctx context.Context) (*StaffStats, error) {
	total, active, inactive, err := uc.repo.CountByStatus(ctx)
	if err != nil {
		return nil, err
	}
	return &StaffStats{Total: total, Active: active, Inactive: inactive}, nil
}
