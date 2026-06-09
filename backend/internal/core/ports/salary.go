package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// SalaryRepository defines persistence operations for salary and advance management.
type SalaryRepository interface {
	// Salary Cycles
	CreateCycle(ctx context.Context, cycle *domain.SalaryCycle) error
	GetCycleByID(ctx context.Context, id uuid.UUID) (*domain.SalaryCycle, error)
	GetCycleByMonthYear(ctx context.Context, month, year int) (*domain.SalaryCycle, error)
	ListCycles(ctx context.Context, filter SalaryCycleFilter) ([]domain.SalaryCycle, int, error)
	UpdateCycleStatus(ctx context.Context, id uuid.UUID, status, generatedAt, generatedBy string) error

	// Salary Records
	CreateRecord(ctx context.Context, record *domain.SalaryRecord) error
	GetRecordByID(ctx context.Context, id uuid.UUID) (*domain.SalaryRecord, error)
	ListRecordsByCycle(ctx context.Context, cycleID uuid.UUID) ([]domain.SalaryRecord, error)
	UpdateRecordPayment(ctx context.Context, id uuid.UUID, status, paymentDate string) error

	// Advances
	CreateAdvance(ctx context.Context, advance *domain.Advance) error
	GetAdvanceByID(ctx context.Context, id uuid.UUID) (*domain.Advance, error)
	ListAdvances(ctx context.Context, filter AdvanceFilter) ([]domain.Advance, int, error)
	UpdateAdvance(ctx context.Context, advance *domain.Advance) error
	GetPendingAdvances(ctx context.Context, staffID uuid.UUID) ([]domain.Advance, error)
	GetTotalOutstandingAdvances(ctx context.Context, staffID uuid.UUID) (float64, error)

	// Stats
	GetPayrollStats(ctx context.Context, month, year int) (*PayrollStats, error)
}

// SalaryCycleFilter holds query params for listing cycles.
type SalaryCycleFilter struct {
	Year   int
	Status string
	Limit  int
	Offset int
}

// AdvanceFilter holds query params for listing advances.
type AdvanceFilter struct {
	StaffID string
	Status  string
	Limit   int
	Offset  int
}

// PayrollStats holds payroll dashboard statistics.
type PayrollStats struct {
	TotalPayroll        float64 `json:"total_payroll"`
	PendingPayments     int     `json:"pending_payments"`
	PaidSalaries        int     `json:"paid_salaries"`
	OutstandingAdvances float64 `json:"outstanding_advances"`
}
