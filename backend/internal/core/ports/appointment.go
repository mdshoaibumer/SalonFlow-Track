package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// AppointmentRepository manages appointments in the database.
type AppointmentRepository interface {
	Create(ctx context.Context, appt *domain.Appointment) error
	Update(ctx context.Context, appt *domain.Appointment) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Appointment, error)
	List(ctx context.Context, filter domain.AppointmentFilter) ([]domain.Appointment, int, error)

	// Services
	AddServices(ctx context.Context, services []domain.AppointmentService) error
	GetServices(ctx context.Context, appointmentID uuid.UUID) ([]domain.AppointmentService, error)
	DeleteServices(ctx context.Context, appointmentID uuid.UUID) error

	// History
	AddHistory(ctx context.Context, history *domain.AppointmentHistory) error
	GetHistory(ctx context.Context, appointmentID uuid.UUID) ([]domain.AppointmentHistory, error)
}
