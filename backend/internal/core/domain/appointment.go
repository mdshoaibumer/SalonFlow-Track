package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/pkg/uid"
)

// Appointment statuses.
const (
	AppointmentStatusBooked     = "booked"
	AppointmentStatusConfirmed  = "confirmed"
	AppointmentStatusInProgress = "in_progress"
	AppointmentStatusCompleted  = "completed"
	AppointmentStatusCancelled  = "cancelled"
	AppointmentStatusNoShow     = "no_show"
)

// Appointment represents a customer appointment.
type Appointment struct {
	ID              uuid.UUID            `json:"id"`
	CustomerID      string               `json:"customer_id"`
	StaffID         string               `json:"staff_id"`
	AppointmentDate string               `json:"appointment_date"`
	StartTime       string               `json:"start_time"`
	EndTime         string               `json:"end_time"`
	Status          string               `json:"status"`
	Notes           string               `json:"notes"`
	IsWalkin        bool                 `json:"is_walkin"`
	TotalAmount     float64              `json:"total_amount"`
	Services        []AppointmentService `json:"services,omitempty"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
}

// NewAppointment creates a new appointment.
func NewAppointment(customerID, staffID, date, startTime, endTime string, isWalkin bool) *Appointment {
	now := time.Now().UTC()
	return &Appointment{
		ID:              uid.New(),
		CustomerID:      customerID,
		StaffID:         staffID,
		AppointmentDate: date,
		StartTime:       startTime,
		EndTime:         endTime,
		Status:          AppointmentStatusBooked,
		IsWalkin:        isWalkin,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// AppointmentService is a service linked to an appointment.
type AppointmentService struct {
	ID              uuid.UUID `json:"id"`
	AppointmentID   uuid.UUID `json:"appointment_id"`
	ServiceID       string    `json:"service_id"`
	ServiceName     string    `json:"service_name"`
	DurationMinutes int       `json:"duration_minutes"`
	Price           float64   `json:"price"`
	CreatedAt       time.Time `json:"created_at"`
}

// NewAppointmentService creates a new appointment service link.
func NewAppointmentService(appointmentID uuid.UUID, serviceID, serviceName string, duration int, price float64) *AppointmentService {
	return &AppointmentService{
		ID:              uid.New(),
		AppointmentID:   appointmentID,
		ServiceID:       serviceID,
		ServiceName:     serviceName,
		DurationMinutes: duration,
		Price:           price,
		CreatedAt:       time.Now().UTC(),
	}
}

// AppointmentHistory records a status change.
type AppointmentHistory struct {
	ID            uuid.UUID `json:"id"`
	AppointmentID uuid.UUID `json:"appointment_id"`
	OldStatus     string    `json:"old_status"`
	NewStatus     string    `json:"new_status"`
	ChangedBy     string    `json:"changed_by"`
	Note          string    `json:"note"`
	CreatedAt     time.Time `json:"created_at"`
}

// NewAppointmentHistory creates a history entry.
func NewAppointmentHistory(appointmentID uuid.UUID, oldStatus, newStatus, changedBy, note string) *AppointmentHistory {
	return &AppointmentHistory{
		ID:            uid.New(),
		AppointmentID: appointmentID,
		OldStatus:     oldStatus,
		NewStatus:     newStatus,
		ChangedBy:     changedBy,
		Note:          note,
		CreatedAt:     time.Now().UTC(),
	}
}

// AppointmentFilter defines filters for listing appointments.
type AppointmentFilter struct {
	Date       string `json:"date"`
	StaffID    string `json:"staff_id"`
	CustomerID string `json:"customer_id"`
	Status     string `json:"status"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
}
