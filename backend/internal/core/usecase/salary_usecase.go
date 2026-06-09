package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// SalaryUseCase handles salary and advance business logic.
type SalaryUseCase struct {
	salaryRepo     ports.SalaryRepository
	staffRepo      ports.StaffRepository
	commissionRepo ports.CommissionRepository
	log            *slog.Logger
}

// NewSalaryUseCase creates a new SalaryUseCase.
func NewSalaryUseCase(
	salaryRepo ports.SalaryRepository,
	staffRepo ports.StaffRepository,
	commissionRepo ports.CommissionRepository,
	log *slog.Logger,
) *SalaryUseCase {
	return &SalaryUseCase{
		salaryRepo:     salaryRepo,
		staffRepo:      staffRepo,
		commissionRepo: commissionRepo,
		log:            log,
	}
}

// --- Advance Management ---

// CreateAdvanceInput is the input DTO for creating an advance.
type CreateAdvanceInput struct {
	StaffID     string  `json:"staff_id"`
	Amount      float64 `json:"amount"`
	AdvanceDate string  `json:"advance_date"`
	Reason      string  `json:"reason"`
}

// CreateAdvance creates a new advance request.
func (uc *SalaryUseCase) CreateAdvance(ctx context.Context, input CreateAdvanceInput) (*domain.Advance, error) {
	staffID, err := uuid.Parse(input.StaffID)
	if err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid staff ID"}
	}

	// Verify staff exists
	_, err = uc.staffRepo.GetByID(ctx, staffID)
	if err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: "staff not found"}
	}

	advance := domain.NewAdvance(staffID, input.Amount, input.AdvanceDate, input.Reason)
	if err := advance.Validate(); err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: err.Error()}
	}

	if err := uc.salaryRepo.CreateAdvance(ctx, advance); err != nil {
		return nil, err
	}

	uc.log.Info("advance created", "id", advance.ID, "staff_id", staffID, "amount", input.Amount)
	return advance, nil
}

// ApproveAdvance approves a pending advance.
func (uc *SalaryUseCase) ApproveAdvance(ctx context.Context, id uuid.UUID) (*domain.Advance, error) {
	advance, err := uc.salaryRepo.GetAdvanceByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if advance.Status != domain.AdvanceStatusPending {
		return nil, &apperror.Error{Kind: apperror.KindBusiness, Message: "advance is not in pending status"}
	}

	advance.Approve()
	if err := uc.salaryRepo.UpdateAdvance(ctx, advance); err != nil {
		return nil, err
	}

	uc.log.Info("advance approved", "id", id)
	return advance, nil
}

// RejectAdvance rejects a pending advance.
func (uc *SalaryUseCase) RejectAdvance(ctx context.Context, id uuid.UUID) (*domain.Advance, error) {
	advance, err := uc.salaryRepo.GetAdvanceByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if advance.Status != domain.AdvanceStatusPending {
		return nil, &apperror.Error{Kind: apperror.KindBusiness, Message: "advance is not in pending status"}
	}

	advance.Reject()
	if err := uc.salaryRepo.UpdateAdvance(ctx, advance); err != nil {
		return nil, err
	}

	uc.log.Info("advance rejected", "id", id)
	return advance, nil
}

// ListAdvancesInput is the input DTO for listing advances.
type ListAdvancesInput struct {
	StaffID string `json:"staff_id"`
	Status  string `json:"status"`
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
}

// ListAdvancesOutput is the output DTO for listing advances.
type ListAdvancesOutput struct {
	Advances   []domain.Advance `json:"advances"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	PerPage    int              `json:"per_page"`
	TotalPages int              `json:"total_pages"`
}

// ListAdvances returns paginated advances.
func (uc *SalaryUseCase) ListAdvances(ctx context.Context, input ListAdvancesInput) (*ListAdvancesOutput, error) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PerPage <= 0 {
		input.PerPage = 20
	}
	if input.PerPage > 100 {
		input.PerPage = 100
	}

	filter := ports.AdvanceFilter{
		StaffID: input.StaffID,
		Status:  input.Status,
		Limit:   input.PerPage,
		Offset:  (input.Page - 1) * input.PerPage,
	}

	advances, total, err := uc.salaryRepo.ListAdvances(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := total / input.PerPage
	if total%input.PerPage > 0 {
		totalPages++
	}

	return &ListAdvancesOutput{
		Advances:   advances,
		Total:      total,
		Page:       input.Page,
		PerPage:    input.PerPage,
		TotalPages: totalPages,
	}, nil
}

// --- Salary Generation ---

// GenerateSalaryInput is the input for generating monthly salaries.
type GenerateSalaryInput struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

// GenerateSalaryOutput is the output of salary generation.
type GenerateSalaryOutput struct {
	Cycle   *domain.SalaryCycle   `json:"cycle"`
	Records []domain.SalaryRecord `json:"records"`
}

// GenerateMonthlySalary generates salary records for all active staff.
func (uc *SalaryUseCase) GenerateMonthlySalary(ctx context.Context, input GenerateSalaryInput) (*GenerateSalaryOutput, error) {
	if input.Month < 1 || input.Month > 12 {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: domain.ErrSalaryInvalidMonth.Error()}
	}
	if input.Year < 2020 {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: domain.ErrSalaryInvalidYear.Error()}
	}

	// Check if cycle already exists
	existing, _ := uc.salaryRepo.GetCycleByMonthYear(ctx, input.Month, input.Year)
	if existing != nil {
		return nil, &apperror.Error{Kind: apperror.KindConflict, Message: domain.ErrSalaryCycleExists.Error()}
	}

	// Create cycle
	cycle := domain.NewSalaryCycle(input.Month, input.Year)
	if err := uc.salaryRepo.CreateCycle(ctx, cycle); err != nil {
		return nil, err
	}

	// Get all active staff
	staffList, _, err := uc.staffRepo.List(ctx, ports.StaffFilter{Status: "active", Limit: 500})
	if err != nil {
		return nil, err
	}

	// Calculate date range for commission lookup
	monthStr := fmt.Sprintf("%d-%02d", input.Year, input.Month)

	var records []domain.SalaryRecord
	for _, staff := range staffList {
		// Get commission for the month
		commission, _ := uc.commissionRepo.GetStaffCommission(ctx, staff.ID,
			fmt.Sprintf("%s-01", monthStr),
			fmt.Sprintf("%s-31", monthStr),
		)

		// Get outstanding advances for recovery
		totalAdvance, _ := uc.salaryRepo.GetTotalOutstandingAdvances(ctx, staff.ID)

		// Create salary record
		record := domain.NewSalaryRecord(
			cycle.ID, staff.ID,
			staff.BaseSalary,
			commission,
			0, // bonus - can be added manually later
			totalAdvance,
			0, // other deductions
		)
		record.StaffName = staff.FullName

		if err := uc.salaryRepo.CreateRecord(ctx, record); err != nil {
			uc.log.Error("failed to create salary record", "staff_id", staff.ID, "error", err)
			continue
		}

		// Mark advances as recovering
		pendingAdvances, _ := uc.salaryRepo.GetPendingAdvances(ctx, staff.ID)
		for i := range pendingAdvances {
			adv := &pendingAdvances[i]
			adv.Recover(adv.RemainingAmount)
			_ = uc.salaryRepo.UpdateAdvance(ctx, adv)
		}

		records = append(records, *record)
	}

	// Update cycle status
	now := time.Now().UTC().Format(time.RFC3339)
	_ = uc.salaryRepo.UpdateCycleStatus(ctx, cycle.ID, domain.SalaryCycleStatusGenerated, now, "system")

	cycle.Status = domain.SalaryCycleStatusGenerated
	cycle.GeneratedAt = now
	cycle.GeneratedBy = "system"

	uc.log.Info("salary generated", "month", input.Month, "year", input.Year, "records", len(records))
	return &GenerateSalaryOutput{Cycle: cycle, Records: records}, nil
}

// --- Salary Queries ---

// GetSalaryByID returns a salary record by ID.
func (uc *SalaryUseCase) GetSalaryByID(ctx context.Context, id uuid.UUID) (*domain.SalaryRecord, error) {
	return uc.salaryRepo.GetRecordByID(ctx, id)
}

// ListSalariesInput is the input for listing salary records.
type ListSalariesInput struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

// ListSalaries returns salary records for a cycle.
func (uc *SalaryUseCase) ListSalaries(ctx context.Context, input ListSalariesInput) ([]domain.SalaryRecord, error) {
	cycle, err := uc.salaryRepo.GetCycleByMonthYear(ctx, input.Month, input.Year)
	if err != nil {
		return nil, err
	}
	return uc.salaryRepo.ListRecordsByCycle(ctx, cycle.ID)
}

// PaySalary marks a salary record as paid.
func (uc *SalaryUseCase) PaySalary(ctx context.Context, id uuid.UUID) error {
	record, err := uc.salaryRepo.GetRecordByID(ctx, id)
	if err != nil {
		return err
	}

	if record.PaymentStatus == domain.SalaryPaymentPaid {
		return &apperror.Error{Kind: apperror.KindBusiness, Message: domain.ErrSalaryAlreadyPaid.Error()}
	}

	paymentDate := time.Now().UTC().Format("2006-01-02")
	return uc.salaryRepo.UpdateRecordPayment(ctx, id, domain.SalaryPaymentPaid, paymentDate)
}

// ListCycles returns salary cycles.
func (uc *SalaryUseCase) ListCycles(ctx context.Context, year int) ([]domain.SalaryCycle, error) {
	cycles, _, err := uc.salaryRepo.ListCycles(ctx, ports.SalaryCycleFilter{Year: year, Limit: 50})
	if err != nil {
		return nil, err
	}
	return cycles, nil
}

// SalaryStats holds salary dashboard statistics.
type SalaryStats struct {
	TotalPayroll        float64 `json:"total_payroll"`
	PendingPayments     int     `json:"pending_payments"`
	PaidSalaries        int     `json:"paid_salaries"`
	OutstandingAdvances float64 `json:"outstanding_advances"`
}

// GetStats returns salary dashboard statistics.
func (uc *SalaryUseCase) GetStats(ctx context.Context) (*SalaryStats, error) {
	now := time.Now().UTC()
	ps, err := uc.salaryRepo.GetPayrollStats(ctx, int(now.Month()), now.Year())
	if err != nil {
		return nil, err
	}

	return &SalaryStats{
		TotalPayroll:        ps.TotalPayroll,
		PendingPayments:     ps.PendingPayments,
		PaidSalaries:        ps.PaidSalaries,
		OutstandingAdvances: ps.OutstandingAdvances,
	}, nil
}
