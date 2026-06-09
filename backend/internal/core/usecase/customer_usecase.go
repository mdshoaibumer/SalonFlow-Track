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

// CustomerUseCase handles customer business logic.
type CustomerUseCase struct {
	repo ports.CustomerRepository
	log  *slog.Logger
}

// NewCustomerUseCase creates a new CustomerUseCase.
func NewCustomerUseCase(repo ports.CustomerRepository, log *slog.Logger) *CustomerUseCase {
	return &CustomerUseCase{repo: repo, log: log}
}

// CreateCustomerInput is the input DTO for creating a customer.
type CreateCustomerInput struct {
	FullName        string `json:"full_name"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	Gender          string `json:"gender"`
	DateOfBirth     string `json:"date_of_birth"`
	AnniversaryDate string `json:"anniversary_date"`
	Address         string `json:"address"`
	Notes           string `json:"notes"`
}

// UpdateCustomerInput is the input DTO for updating a customer.
type UpdateCustomerInput struct {
	FullName        string `json:"full_name"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	Gender          string `json:"gender"`
	DateOfBirth     string `json:"date_of_birth"`
	AnniversaryDate string `json:"anniversary_date"`
	Address         string `json:"address"`
	Notes           string `json:"notes"`
	Status          string `json:"status"`
}

// ListCustomerInput is the input DTO for listing customers.
type ListCustomerInput struct {
	Search  string `json:"search"`
	Status  string `json:"status"`
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
}

// ListCustomerOutput is the output DTO for listing customers.
type ListCustomerOutput struct {
	Customers  []domain.Customer `json:"customers"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	PerPage    int               `json:"per_page"`
	TotalPages int               `json:"total_pages"`
}

// CustomerStats holds customer statistics.
type CustomerStats struct {
	Total         int `json:"total"`
	Active        int `json:"active"`
	Inactive      int `json:"inactive"`
	NewThisMonth  int `json:"new_this_month"`
	BirthdayToday int `json:"birthday_today"`
}

// Create creates a new customer.
func (uc *CustomerUseCase) Create(ctx context.Context, input CreateCustomerInput) (*domain.Customer, error) {
	customer := domain.NewCustomer(input.FullName, input.Phone)
	customer.Email = strings.TrimSpace(input.Email)
	customer.Address = strings.TrimSpace(input.Address)
	customer.Notes = strings.TrimSpace(input.Notes)

	if input.Gender != "" {
		customer.Gender = input.Gender
	}
	if input.DateOfBirth != "" {
		if t, err := time.Parse("2006-01-02", input.DateOfBirth); err == nil {
			customer.DateOfBirth = &t
		}
	}
	if input.AnniversaryDate != "" {
		if t, err := time.Parse("2006-01-02", input.AnniversaryDate); err == nil {
			customer.AnniversaryDate = &t
		}
	}

	if err := customer.Validate(); err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: err.Error()}
	}

	// Check phone uniqueness
	existing, err := uc.repo.GetByPhone(ctx, customer.Phone)
	if err == nil && existing != nil {
		return nil, &apperror.Error{Kind: apperror.KindConflict, Message: domain.ErrCustomerPhoneDuplicate.Error(), Field: "phone"}
	}

	if err := uc.repo.Create(ctx, customer); err != nil {
		return nil, err
	}

	uc.log.Info("customer created", "id", customer.ID, "name", customer.FullName)
	return customer, nil
}

// GetByID retrieves a customer by ID.
func (uc *CustomerUseCase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	return uc.repo.GetByID(ctx, id)
}

// List returns paginated customer list.
func (uc *CustomerUseCase) List(ctx context.Context, input ListCustomerInput) (*ListCustomerOutput, error) {
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

	filter := ports.CustomerFilter{
		Status: input.Status,
		Search: input.Search,
		Limit:  input.PerPage,
		Offset: offset,
	}

	customers, total, err := uc.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := total / input.PerPage
	if total%input.PerPage > 0 {
		totalPages++
	}

	return &ListCustomerOutput{
		Customers:  customers,
		Total:      total,
		Page:       input.Page,
		PerPage:    input.PerPage,
		TotalPages: totalPages,
	}, nil
}

// Update updates an existing customer.
func (uc *CustomerUseCase) Update(ctx context.Context, id uuid.UUID, input UpdateCustomerInput) (*domain.Customer, error) {
	customer, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	customer.FullName = strings.TrimSpace(input.FullName)
	customer.Phone = strings.TrimSpace(input.Phone)
	customer.Email = strings.TrimSpace(input.Email)
	customer.Gender = input.Gender
	customer.Address = strings.TrimSpace(input.Address)
	customer.Notes = strings.TrimSpace(input.Notes)

	if input.Status != "" {
		customer.Status = input.Status
	}
	if input.DateOfBirth != "" {
		if t, err := time.Parse("2006-01-02", input.DateOfBirth); err == nil {
			customer.DateOfBirth = &t
		}
	} else {
		customer.DateOfBirth = nil
	}
	if input.AnniversaryDate != "" {
		if t, err := time.Parse("2006-01-02", input.AnniversaryDate); err == nil {
			customer.AnniversaryDate = &t
		}
	} else {
		customer.AnniversaryDate = nil
	}

	if err := customer.Validate(); err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: err.Error()}
	}

	// Check phone uniqueness (exclude self)
	existing, err := uc.repo.GetByPhone(ctx, customer.Phone)
	if err == nil && existing != nil && existing.ID != customer.ID {
		return nil, &apperror.Error{Kind: apperror.KindConflict, Message: domain.ErrCustomerPhoneDuplicate.Error(), Field: "phone"}
	}

	if err := uc.repo.Update(ctx, customer); err != nil {
		return nil, err
	}

	uc.log.Info("customer updated", "id", customer.ID, "name", customer.FullName)
	return customer, nil
}

// Delete removes a customer.
func (uc *CustomerUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return err
	}
	uc.log.Info("customer deleted", "id", id)
	return nil
}

// Stats returns customer statistics.
func (uc *CustomerUseCase) Stats(ctx context.Context) (*CustomerStats, error) {
	total, active, inactive, err := uc.repo.CountByStatus(ctx)
	if err != nil {
		return nil, err
	}

	newThisMonth, err := uc.repo.CountNewThisMonth(ctx)
	if err != nil {
		return nil, err
	}

	birthdayToday, err := uc.repo.CountBirthdayToday(ctx)
	if err != nil {
		return nil, err
	}

	return &CustomerStats{
		Total:         total,
		Active:        active,
		Inactive:      inactive,
		NewThisMonth:  newThisMonth,
		BirthdayToday: birthdayToday,
	}, nil
}
