package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/pkg/uid"
)

// Salary cycle statuses.
const (
	SalaryCycleStatusDraft     = "draft"
	SalaryCycleStatusGenerated = "generated"
	SalaryCycleStatusFinalized = "finalized"
)

// Salary payment statuses.
const (
	SalaryPaymentPending = "pending"
	SalaryPaymentPartial = "partial"
	SalaryPaymentPaid    = "paid"
)

// Advance statuses.
const (
	AdvanceStatusPending    = "pending"
	AdvanceStatusApproved   = "approved"
	AdvanceStatusRecovering = "recovering"
	AdvanceStatusRecovered  = "recovered"
	AdvanceStatusRejected   = "rejected"
)

// SalaryCycle represents a monthly salary generation cycle.
type SalaryCycle struct {
	ID          uuid.UUID `json:"id"`
	Month       int       `json:"month"`
	Year        int       `json:"year"`
	Status      string    `json:"status"`
	GeneratedAt string    `json:"generated_at"`
	GeneratedBy string    `json:"generated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SalaryRecord represents a staff salary for a given cycle.
type SalaryRecord struct {
	ID               uuid.UUID `json:"id"`
	SalaryCycleID    uuid.UUID `json:"salary_cycle_id"`
	StaffID          uuid.UUID `json:"staff_id"`
	StaffName        string    `json:"staff_name,omitempty"`
	BaseSalary       float64   `json:"base_salary"`
	CommissionAmount float64   `json:"commission_amount"`
	BonusAmount      float64   `json:"bonus_amount"`
	AdvanceAmount    float64   `json:"advance_amount"`
	DeductionAmount  float64   `json:"deduction_amount"`
	GrossSalary      float64   `json:"gross_salary"`
	NetSalary        float64   `json:"net_salary"`
	PaymentStatus    string    `json:"payment_status"`
	PaymentDate      string    `json:"payment_date"`
	Notes            string    `json:"notes"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Advance represents a salary advance given to staff.
type Advance struct {
	ID              uuid.UUID `json:"id"`
	StaffID         uuid.UUID `json:"staff_id"`
	StaffName       string    `json:"staff_name,omitempty"`
	Amount          float64   `json:"amount"`
	AdvanceDate     string    `json:"advance_date"`
	Reason          string    `json:"reason"`
	RecoveredAmount float64   `json:"recovered_amount"`
	RemainingAmount float64   `json:"remaining_amount"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// NewSalaryCycle creates a new salary cycle.
func NewSalaryCycle(month, year int) *SalaryCycle {
	now := time.Now().UTC()
	return &SalaryCycle{
		ID:        uid.New(),
		Month:     month,
		Year:      year,
		Status:    SalaryCycleStatusDraft,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewSalaryRecord creates a new salary record.
func NewSalaryRecord(cycleID, staffID uuid.UUID, baseSalary, commission, bonus, advance, deduction float64) *SalaryRecord {
	gross := baseSalary + commission + bonus
	net := gross - advance - deduction

	now := time.Now().UTC()
	return &SalaryRecord{
		ID:               uid.New(),
		SalaryCycleID:    cycleID,
		StaffID:          staffID,
		BaseSalary:       baseSalary,
		CommissionAmount: commission,
		BonusAmount:      bonus,
		AdvanceAmount:    advance,
		DeductionAmount:  deduction,
		GrossSalary:      gross,
		NetSalary:        net,
		PaymentStatus:    SalaryPaymentPending,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// NewAdvance creates a new advance request.
func NewAdvance(staffID uuid.UUID, amount float64, date, reason string) *Advance {
	now := time.Now().UTC()
	return &Advance{
		ID:              uid.New(),
		StaffID:         staffID,
		Amount:          amount,
		AdvanceDate:     date,
		Reason:          reason,
		RecoveredAmount: 0,
		RemainingAmount: amount,
		Status:          AdvanceStatusPending,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// Validate validates the advance.
func (a *Advance) Validate() error {
	if a.Amount <= 0 {
		return ErrAdvanceInvalidAmount
	}
	if a.AdvanceDate == "" {
		return ErrAdvanceDateRequired
	}
	return nil
}

// Approve approves the advance.
func (a *Advance) Approve() {
	a.Status = AdvanceStatusApproved
	a.UpdatedAt = time.Now().UTC()
}

// Reject rejects the advance.
func (a *Advance) Reject() {
	a.Status = AdvanceStatusRejected
	a.UpdatedAt = time.Now().UTC()
}

// Recover records a partial recovery against the advance.
func (a *Advance) Recover(amount float64) {
	a.RecoveredAmount += amount
	a.RemainingAmount = a.Amount - a.RecoveredAmount
	if a.RemainingAmount <= 0 {
		a.RemainingAmount = 0
		a.Status = AdvanceStatusRecovered
	} else {
		a.Status = AdvanceStatusRecovering
	}
	a.UpdatedAt = time.Now().UTC()
}

// CalculateNetSalary recalculates net salary.
func (r *SalaryRecord) CalculateNetSalary() {
	r.GrossSalary = r.BaseSalary + r.CommissionAmount + r.BonusAmount
	r.NetSalary = r.GrossSalary - r.AdvanceAmount - r.DeductionAmount
	r.UpdatedAt = time.Now().UTC()
}

// MarkPaid marks the salary record as paid.
func (r *SalaryRecord) MarkPaid() {
	r.PaymentStatus = SalaryPaymentPaid
	r.PaymentDate = time.Now().UTC().Format("2006-01-02")
	r.UpdatedAt = time.Now().UTC()
}
