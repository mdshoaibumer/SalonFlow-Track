package usecase

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// mockCustomerRepo is a test double for CustomerRepository.
type mockCustomerRepo struct {
	customers map[uuid.UUID]*domain.Customer
}

func newMockCustomerRepo() *mockCustomerRepo {
	return &mockCustomerRepo{customers: make(map[uuid.UUID]*domain.Customer)}
}

func (m *mockCustomerRepo) Create(_ context.Context, c *domain.Customer) error {
	m.customers[c.ID] = c
	return nil
}
func (m *mockCustomerRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Customer, error) {
	if c, ok := m.customers[id]; ok {
		return c, nil
	}
	return nil, apperror.NotFound("customer", id.String())
}
func (m *mockCustomerRepo) GetByPhone(_ context.Context, phone string) (*domain.Customer, error) {
	for _, c := range m.customers {
		if c.Phone == phone {
			return c, nil
		}
	}
	return nil, apperror.NotFound("customer", phone)
}
func (m *mockCustomerRepo) List(_ context.Context, _ ports.CustomerFilter) ([]domain.Customer, int, error) {
	var result []domain.Customer
	for _, c := range m.customers {
		result = append(result, *c)
	}
	return result, len(result), nil
}
func (m *mockCustomerRepo) Update(_ context.Context, c *domain.Customer) error {
	m.customers[c.ID] = c
	return nil
}
func (m *mockCustomerRepo) Delete(_ context.Context, id uuid.UUID) error {
	delete(m.customers, id)
	return nil
}
func (m *mockCustomerRepo) CountByStatus(_ context.Context) (int, int, int, error) {
	var total, active, inactive int
	for _, c := range m.customers {
		total++
		if c.Status == "active" {
			active++
		} else {
			inactive++
		}
	}
	return total, active, inactive, nil
}
func (m *mockCustomerRepo) CountNewThisMonth(_ context.Context) (int, error)  { return 0, nil }
func (m *mockCustomerRepo) CountBirthdayToday(_ context.Context) (int, error) { return 0, nil }

func TestCustomerUseCase_Create(t *testing.T) {
	repo := newMockCustomerRepo()
	uc := NewCustomerUseCase(repo, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	c, err := uc.Create(context.Background(), CreateCustomerInput{
		FullName: "John Doe",
		Phone:    "9876543210",
		Gender:   "male",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if c.FullName != "John Doe" {
		t.Errorf("expected 'John Doe', got %q", c.FullName)
	}
}

func TestCustomerUseCase_CreateValidationError(t *testing.T) {
	repo := newMockCustomerRepo()
	uc := NewCustomerUseCase(repo, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	_, err := uc.Create(context.Background(), CreateCustomerInput{
		FullName: "",
		Phone:    "9876543210",
	})
	if err == nil {
		t.Error("expected validation error")
	}
}

func TestCustomerUseCase_CreateDuplicatePhone(t *testing.T) {
	repo := newMockCustomerRepo()
	uc := NewCustomerUseCase(repo, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	_, _ = uc.Create(context.Background(), CreateCustomerInput{FullName: "John", Phone: "9876543210", Gender: "male"})
	_, err := uc.Create(context.Background(), CreateCustomerInput{FullName: "Jane", Phone: "9876543210", Gender: "female"})
	if err == nil {
		t.Error("expected duplicate phone error")
	}
}

func TestCustomerUseCase_Update(t *testing.T) {
	repo := newMockCustomerRepo()
	uc := NewCustomerUseCase(repo, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	c, _ := uc.Create(context.Background(), CreateCustomerInput{FullName: "John", Phone: "9876543210", Gender: "male"})

	updated, err := uc.Update(context.Background(), c.ID, UpdateCustomerInput{
		FullName: "John Updated",
		Phone:    "9876543210",
		Gender:   "male",
		Email:    "john@test.com",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updated.FullName != "John Updated" {
		t.Errorf("expected 'John Updated', got %q", updated.FullName)
	}
}

func TestCustomerUseCase_Delete(t *testing.T) {
	repo := newMockCustomerRepo()
	uc := NewCustomerUseCase(repo, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	c, _ := uc.Create(context.Background(), CreateCustomerInput{FullName: "John", Phone: "9876543210", Gender: "male"})
	err := uc.Delete(context.Background(), c.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(repo.customers) != 0 {
		t.Error("expected empty repo")
	}
}

func TestCustomerUseCase_Stats(t *testing.T) {
	repo := newMockCustomerRepo()
	uc := NewCustomerUseCase(repo, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	_, _ = uc.Create(context.Background(), CreateCustomerInput{FullName: "A", Phone: "9876543210", Gender: "male"})
	_, _ = uc.Create(context.Background(), CreateCustomerInput{FullName: "B", Phone: "9876543211", Gender: "female"})

	stats, err := uc.Stats(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if stats.Total != 2 {
		t.Errorf("expected 2, got %d", stats.Total)
	}
	if stats.Active != 2 {
		t.Errorf("expected 2 active, got %d", stats.Active)
	}
}
