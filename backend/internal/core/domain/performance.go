package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/pkg/uid"
)

// StaffPerformanceDaily represents daily performance metrics for a staff member.
type StaffPerformanceDaily struct {
	ID               uuid.UUID `json:"id"`
	StaffID          uuid.UUID `json:"staff_id"`
	BusinessDate     string    `json:"business_date"` // YYYY-MM-DD
	InvoiceCount     int       `json:"invoice_count"`
	CustomerCount    int       `json:"customer_count"`
	ServiceCount     int       `json:"service_count"`
	Revenue          float64   `json:"revenue"`
	CommissionAmount float64   `json:"commission_amount"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// StaffPerformanceSummary is an aggregated performance view.
type StaffPerformanceSummary struct {
	StaffID       uuid.UUID `json:"staff_id"`
	StaffName     string    `json:"staff_name"`
	Revenue       float64   `json:"revenue"`
	CustomerCount int       `json:"customer_count"`
	InvoiceCount  int       `json:"invoice_count"`
	ServiceCount  int       `json:"service_count"`
	AvgBill       float64   `json:"avg_bill"`
	Commission    float64   `json:"commission"`
	Rank          int       `json:"rank"`
}

// RevenueTrendPoint represents a single point in a revenue trend.
type RevenueTrendPoint struct {
	Date    string  `json:"date"`
	Revenue float64 `json:"revenue"`
}

// NewStaffPerformanceDaily creates a new performance record.
func NewStaffPerformanceDaily(staffID uuid.UUID, businessDate string) *StaffPerformanceDaily {
	now := time.Now().UTC()
	return &StaffPerformanceDaily{
		ID:           uid.New(),
		StaffID:      staffID,
		BusinessDate: businessDate,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// AddInvoice adds invoice data to the daily performance record.
func (p *StaffPerformanceDaily) AddInvoice(revenue float64, serviceCount int, commission float64) {
	p.InvoiceCount++
	p.CustomerCount++
	p.ServiceCount += serviceCount
	p.Revenue += revenue
	p.CommissionAmount += commission
	p.UpdatedAt = time.Now().UTC()
}
