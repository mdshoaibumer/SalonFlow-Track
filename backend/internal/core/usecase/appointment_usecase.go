package usecase

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// AppointmentUseCase handles appointment business logic.
type AppointmentUseCase struct {
	repo ports.AppointmentRepository
	log  *slog.Logger
}

// NewAppointmentUseCase creates a new AppointmentUseCase.
func NewAppointmentUseCase(repo ports.AppointmentRepository, log *slog.Logger) *AppointmentUseCase {
	return &AppointmentUseCase{repo: repo, log: log}
}

// Create creates a new appointment with services.
func (uc *AppointmentUseCase) Create(ctx context.Context, appt *domain.Appointment, services []domain.AppointmentService) error {
	if appt.AppointmentDate == "" {
		return apperror.Validation("appointment_date", "Date is required")
	}
	if appt.StartTime == "" {
		return apperror.Validation("start_time", "Start time is required")
	}

	// Calculate total from services
	var total float64
	for i := range services {
		services[i].AppointmentID = appt.ID
		if services[i].ID == uuid.Nil {
			services[i].ID = domain.NewAppointmentService(appt.ID, "", "", 0, 0).ID
		}
		total += services[i].Price
	}
	appt.TotalAmount = total

	if err := uc.repo.Create(ctx, appt); err != nil {
		return err
	}

	if len(services) > 0 {
		if err := uc.repo.AddServices(ctx, services); err != nil {
			return err
		}
	}

	// Record history
	history := domain.NewAppointmentHistory(appt.ID, "", appt.Status, "", "Appointment created")
	return uc.repo.AddHistory(ctx, history)
}

// Update updates an appointment.
func (uc *AppointmentUseCase) Update(ctx context.Context, appt *domain.Appointment, services []domain.AppointmentService) error {
	existing, err := uc.repo.GetByID(ctx, appt.ID)
	if err != nil {
		return err
	}

	// Calculate total
	var total float64
	for i := range services {
		services[i].AppointmentID = appt.ID
		if services[i].ID == uuid.Nil {
			services[i].ID = domain.NewAppointmentService(appt.ID, "", "", 0, 0).ID
		}
		total += services[i].Price
	}
	appt.TotalAmount = total

	if err := uc.repo.Update(ctx, appt); err != nil {
		return err
	}

	// Replace services
	if len(services) > 0 {
		_ = uc.repo.DeleteServices(ctx, appt.ID)
		if err := uc.repo.AddServices(ctx, services); err != nil {
			return err
		}
	}

	// Record status change if changed
	if existing.Status != appt.Status {
		history := domain.NewAppointmentHistory(appt.ID, existing.Status, appt.Status, "", "Status updated")
		_ = uc.repo.AddHistory(ctx, history)
	}

	return nil
}

// UpdateStatus changes an appointment's status.
func (uc *AppointmentUseCase) UpdateStatus(ctx context.Context, id uuid.UUID, status, note string) error {
	appt, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	oldStatus := appt.Status
	appt.Status = status
	if err := uc.repo.Update(ctx, appt); err != nil {
		return err
	}

	history := domain.NewAppointmentHistory(id, oldStatus, status, "", note)
	return uc.repo.AddHistory(ctx, history)
}

// Delete deletes an appointment.
func (uc *AppointmentUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}

// GetByID gets an appointment with services.
func (uc *AppointmentUseCase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Appointment, error) {
	appt, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	services, _ := uc.repo.GetServices(ctx, id)
	appt.Services = services
	return appt, nil
}

// List lists appointments with filters.
func (uc *AppointmentUseCase) List(ctx context.Context, filter domain.AppointmentFilter) ([]domain.Appointment, int, error) {
	return uc.repo.List(ctx, filter)
}

// GetHistory gets appointment history.
func (uc *AppointmentUseCase) GetHistory(ctx context.Context, id uuid.UUID) ([]domain.AppointmentHistory, error) {
	return uc.repo.GetHistory(ctx, id)
}
