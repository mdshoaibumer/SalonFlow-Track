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

// mockServiceRepo is a test double for ServiceRepository.
type mockServiceRepo struct {
	services map[uuid.UUID]*domain.Service
}

func newMockServiceRepo() *mockServiceRepo {
	return &mockServiceRepo{services: make(map[uuid.UUID]*domain.Service)}
}

func (m *mockServiceRepo) Create(_ context.Context, svc *domain.Service) error {
	m.services[svc.ID] = svc
	return nil
}
func (m *mockServiceRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Service, error) {
	if s, ok := m.services[id]; ok {
		return s, nil
	}
	return nil, apperror.NotFound("service", id.String())
}
func (m *mockServiceRepo) GetByName(_ context.Context, name string) (*domain.Service, error) {
	for _, s := range m.services {
		if s.Name == name {
			return s, nil
		}
	}
	return nil, apperror.NotFound("service", name)
}
func (m *mockServiceRepo) List(_ context.Context, filter ports.ServiceFilter) ([]domain.Service, int, error) {
	var result []domain.Service
	for _, s := range m.services {
		if filter.Status != "" && s.Status != filter.Status {
			continue
		}
		if filter.Category != "" && s.Category != filter.Category {
			continue
		}
		result = append(result, *s)
	}
	return result, len(result), nil
}
func (m *mockServiceRepo) Update(_ context.Context, svc *domain.Service) error {
	m.services[svc.ID] = svc
	return nil
}
func (m *mockServiceRepo) Delete(_ context.Context, id uuid.UUID) error {
	delete(m.services, id)
	return nil
}
func (m *mockServiceRepo) CountByStatus(_ context.Context) (int, int, int, error) {
	var total, active, inactive int
	for _, s := range m.services {
		total++
		if s.Status == "active" {
			active++
		} else {
			inactive++
		}
	}
	return total, active, inactive, nil
}

func TestServiceUseCase_Create(t *testing.T) {
	repo := newMockServiceRepo()
	uc := NewServiceUseCase(repo, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	svc, err := uc.Create(context.Background(), CreateServiceInput{
		Name:            "Hair Cut",
		Category:        "hair",
		DurationMinutes: 30,
		Price:           300,
		CommissionType:  "percentage",
		CommissionValue: 10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if svc.Name != "Hair Cut" {
		t.Errorf("expected 'Hair Cut', got %q", svc.Name)
	}
}

func TestServiceUseCase_CreateValidationError(t *testing.T) {
	repo := newMockServiceRepo()
	uc := NewServiceUseCase(repo, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	_, err := uc.Create(context.Background(), CreateServiceInput{
		Name:            "",
		Category:        "hair",
		DurationMinutes: 30,
		Price:           300,
	})
	if err == nil {
		t.Error("expected validation error")
	}
}

func TestServiceUseCase_List(t *testing.T) {
	repo := newMockServiceRepo()
	uc := NewServiceUseCase(repo, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	_, _ = uc.Create(context.Background(), CreateServiceInput{Name: "A", Category: "hair", DurationMinutes: 30, Price: 300})
	_, _ = uc.Create(context.Background(), CreateServiceInput{Name: "B", Category: "facial", DurationMinutes: 45, Price: 800})

	out, err := uc.List(context.Background(), ListServiceInput{Page: 1, PerPage: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out.Total != 2 {
		t.Errorf("expected 2, got %d", out.Total)
	}
}

func TestServiceUseCase_Update(t *testing.T) {
	repo := newMockServiceRepo()
	uc := NewServiceUseCase(repo, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	svc, _ := uc.Create(context.Background(), CreateServiceInput{
		Name:            "Hair Cut",
		Category:        "hair",
		DurationMinutes: 30,
		Price:           300,
		CommissionType:  "percentage",
	})

	updated, err := uc.Update(context.Background(), svc.ID, UpdateServiceInput{
		Name:            "Premium Hair Cut",
		Category:        "hair",
		DurationMinutes: 45,
		Price:           500,
		CommissionType:  "fixed",
		CommissionValue: 50,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updated.Name != "Premium Hair Cut" {
		t.Errorf("expected 'Premium Hair Cut', got %q", updated.Name)
	}
	if updated.Price != 500 {
		t.Errorf("expected 500, got %f", updated.Price)
	}
}

func TestServiceUseCase_Delete(t *testing.T) {
	repo := newMockServiceRepo()
	uc := NewServiceUseCase(repo, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	svc, _ := uc.Create(context.Background(), CreateServiceInput{
		Name:            "Hair Cut",
		Category:        "hair",
		DurationMinutes: 30,
		Price:           300,
		CommissionType:  "percentage",
	})

	err := uc.Delete(context.Background(), svc.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(repo.services) != 0 {
		t.Error("expected empty repo")
	}
}

func TestServiceUseCase_Stats(t *testing.T) {
	repo := newMockServiceRepo()
	uc := NewServiceUseCase(repo, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	_, _ = uc.Create(context.Background(), CreateServiceInput{Name: "A", Category: "hair", DurationMinutes: 30, Price: 300, CommissionType: "percentage"})
	_, _ = uc.Create(context.Background(), CreateServiceInput{Name: "B", Category: "facial", DurationMinutes: 45, Price: 800, CommissionType: "percentage"})

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
