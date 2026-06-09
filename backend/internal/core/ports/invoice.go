package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// InvoiceRepository defines persistence operations for Invoice.
type InvoiceRepository interface {
	Create(ctx context.Context, invoice *domain.Invoice) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Invoice, error)
	List(ctx context.Context, filter InvoiceFilter) ([]domain.Invoice, int, error)
	GetNextSequence(ctx context.Context, year int) (int, error)
	RecordPayment(ctx context.Context, payment *domain.Payment) error
	GetPayments(ctx context.Context, invoiceID uuid.UUID) ([]domain.Payment, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status, method string) error
	GetTodayStats(ctx context.Context) (todayRevenue float64, todayCount int, avgBill float64, err error)
}

// InvoiceFilter holds query parameters for listing invoices.
type InvoiceFilter struct {
	CustomerID    string
	StaffID       string
	PaymentStatus string
	DateFrom      string
	DateTo        string
	Search        string
	Limit         int
	Offset        int
}
