package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// InvoiceUseCase handles invoice/billing business logic.
type InvoiceUseCase struct {
	invoiceRepo  ports.InvoiceRepository
	customerRepo ports.CustomerRepository
	serviceRepo  ports.ServiceRepository
	staffRepo    ports.StaffRepository
	perfUC       *PerformanceUseCase
	commUC       *CommissionUseCase
	log          *slog.Logger
}

// NewInvoiceUseCase creates a new InvoiceUseCase.
func NewInvoiceUseCase(
	invoiceRepo ports.InvoiceRepository,
	customerRepo ports.CustomerRepository,
	serviceRepo ports.ServiceRepository,
	staffRepo ports.StaffRepository,
	perfUC *PerformanceUseCase,
	commUC *CommissionUseCase,
	log *slog.Logger,
) *InvoiceUseCase {
	return &InvoiceUseCase{
		invoiceRepo:  invoiceRepo,
		customerRepo: customerRepo,
		serviceRepo:  serviceRepo,
		staffRepo:    staffRepo,
		perfUC:       perfUC,
		commUC:       commUC,
		log:          log,
	}
}

// CreateInvoiceItemInput is the input for an invoice line item.
type CreateInvoiceItemInput struct {
	ServiceID string  `json:"service_id"`
	Quantity  int     `json:"quantity"`
	Discount  float64 `json:"discount"`
}

// CreateInvoiceInput is the input DTO for creating an invoice.
type CreateInvoiceInput struct {
	CustomerID    string                   `json:"customer_id"`
	StaffID       string                   `json:"staff_id"`
	Items         []CreateInvoiceItemInput `json:"items"`
	Discount      float64                  `json:"discount"`
	Tax           float64                  `json:"tax"`
	PaymentMethod string                   `json:"payment_method"`
	Notes         string                   `json:"notes"`
}

// ListInvoiceInput is the input DTO for listing invoices.
type ListInvoiceInput struct {
	CustomerID    string `json:"customer_id"`
	StaffID       string `json:"staff_id"`
	PaymentStatus string `json:"payment_status"`
	DateFrom      string `json:"date_from"`
	DateTo        string `json:"date_to"`
	Search        string `json:"search"`
	Page          int    `json:"page"`
	PerPage       int    `json:"per_page"`
}

// ListInvoiceOutput is the output DTO for listing invoices.
type ListInvoiceOutput struct {
	Invoices   []domain.Invoice `json:"invoices"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	PerPage    int              `json:"per_page"`
	TotalPages int              `json:"total_pages"`
}

// RecordPaymentInput is the input DTO for recording a payment.
type RecordPaymentInput struct {
	Amount          float64 `json:"amount"`
	PaymentMethod   string  `json:"payment_method"`
	ReferenceNumber string  `json:"reference_number"`
}

// InvoiceStats holds dashboard billing statistics.
type InvoiceStats struct {
	TodayRevenue  float64 `json:"today_revenue"`
	TodayInvoices int     `json:"today_invoices"`
	AvgBillValue  float64 `json:"avg_bill_value"`
}

// Create creates a new invoice.
func (uc *InvoiceUseCase) Create(ctx context.Context, input CreateInvoiceInput) (*domain.Invoice, error) {
	// Validate customer
	customerID, err := uuid.Parse(input.CustomerID)
	if err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid customer ID"}
	}
	_, err = uc.customerRepo.GetByID(ctx, customerID)
	if err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: "customer not found"}
	}

	// Validate staff
	staffID, err := uuid.Parse(input.StaffID)
	if err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid staff ID"}
	}
	_, err = uc.staffRepo.GetByID(ctx, staffID)
	if err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: "staff not found"}
	}

	// Validate items
	if len(input.Items) == 0 {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: domain.ErrInvoiceItemsRequired.Error()}
	}

	// Generate invoice number
	year := time.Now().UTC().Year()
	seq, err := uc.invoiceRepo.GetNextSequence(ctx, year)
	if err != nil {
		return nil, err
	}
	invoiceNumber := domain.GenerateInvoiceNumber(year, seq)

	// Build invoice
	invoice := domain.NewInvoice(customerID, staffID, invoiceNumber)
	invoice.Discount = input.Discount
	invoice.Tax = input.Tax
	invoice.Notes = input.Notes

	// Build items
	for _, itemInput := range input.Items {
		serviceID, err := uuid.Parse(itemInput.ServiceID)
		if err != nil {
			return nil, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid service ID"}
		}

		svc, err := uc.serviceRepo.GetByID(ctx, serviceID)
		if err != nil {
			return nil, &apperror.Error{Kind: apperror.KindValidation, Message: "service not found: " + itemInput.ServiceID}
		}

		qty := itemInput.Quantity
		if qty <= 0 {
			qty = 1
		}

		item := domain.NewInvoiceItem(invoice.ID, serviceID, svc.Name, qty, svc.Price, itemInput.Discount)
		invoice.AddItem(item)
	}

	// Validate the invoice
	if err := invoice.Validate(); err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: err.Error()}
	}

	// Determine payment status
	if input.PaymentMethod != "" {
		invoice.MarkPaid(input.PaymentMethod)
	}

	// Persist
	if err := uc.invoiceRepo.Create(ctx, invoice); err != nil {
		return nil, err
	}

	// Update customer stats
	customer, _ := uc.customerRepo.GetByID(ctx, customerID)
	if customer != nil {
		customer.RecordVisit(invoice.GrandTotal, invoice.InvoiceDate)
		_ = uc.customerRepo.Update(ctx, customer)
	}

	// Record staff performance automatically
	if uc.perfUC != nil {
		serviceCount := len(invoice.Items)
		// Calculate commission from service settings
		var commission float64
		if uc.commUC != nil {
			var serviceIDs []uuid.UUID
			for _, item := range invoice.Items {
				serviceIDs = append(serviceIDs, item.ServiceID)
			}
			commission, _ = uc.commUC.CalculateInvoiceCommission(ctx, staffID, invoice.ID, invoice.GrandTotal, serviceIDs)
		}
		_ = uc.perfUC.RecordInvoicePerformance(ctx, staffID, invoice.GrandTotal, serviceCount, commission)
	}

	uc.log.Info("invoice created", "id", invoice.ID, "number", invoice.InvoiceNumber, "total", invoice.GrandTotal)
	return invoice, nil
}

// GetByID retrieves an invoice by ID with items.
func (uc *InvoiceUseCase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Invoice, error) {
	return uc.invoiceRepo.GetByID(ctx, id)
}

// List returns paginated invoice list.
func (uc *InvoiceUseCase) List(ctx context.Context, input ListInvoiceInput) (*ListInvoiceOutput, error) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PerPage <= 0 {
		input.PerPage = 20
	}
	if input.PerPage > 100 {
		input.PerPage = 100
	}

	offset := (input.Page - 1) * input.PerPage

	filter := ports.InvoiceFilter{
		CustomerID:    input.CustomerID,
		StaffID:       input.StaffID,
		PaymentStatus: input.PaymentStatus,
		DateFrom:      input.DateFrom,
		DateTo:        input.DateTo,
		Search:        input.Search,
		Limit:         input.PerPage,
		Offset:        offset,
	}

	invoices, total, err := uc.invoiceRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := total / input.PerPage
	if total%input.PerPage > 0 {
		totalPages++
	}

	return &ListInvoiceOutput{
		Invoices:   invoices,
		Total:      total,
		Page:       input.Page,
		PerPage:    input.PerPage,
		TotalPages: totalPages,
	}, nil
}

// RecordPayment records a payment for an invoice.
func (uc *InvoiceUseCase) RecordPayment(ctx context.Context, invoiceID uuid.UUID, input RecordPaymentInput) (*domain.Payment, error) {
	invoice, err := uc.invoiceRepo.GetByID(ctx, invoiceID)
	if err != nil {
		return nil, err
	}

	if invoice.IsPaid() {
		return nil, &apperror.Error{Kind: apperror.KindBusiness, Message: domain.ErrInvoiceAlreadyPaid.Error()}
	}

	payment := domain.NewPayment(invoiceID, input.Amount, input.PaymentMethod, input.ReferenceNumber)
	if err := payment.Validate(); err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: err.Error()}
	}

	// Calculate total paid
	existingPayments, _ := uc.invoiceRepo.GetPayments(ctx, invoiceID)
	var totalPaid float64
	for _, p := range existingPayments {
		totalPaid += p.Amount
	}
	totalPaid += input.Amount

	if totalPaid > invoice.GrandTotal {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: domain.ErrPaymentExceedsBalance.Error()}
	}

	// Record payment
	if err := uc.invoiceRepo.RecordPayment(ctx, payment); err != nil {
		return nil, err
	}

	// Update invoice status
	if totalPaid >= invoice.GrandTotal {
		_ = uc.invoiceRepo.UpdateStatus(ctx, invoiceID, domain.PaymentStatusPaid, input.PaymentMethod)
	} else {
		_ = uc.invoiceRepo.UpdateStatus(ctx, invoiceID, domain.PaymentStatusPartial, input.PaymentMethod)
	}

	uc.log.Info("payment recorded", "invoice_id", invoiceID, "amount", input.Amount)
	return payment, nil
}

// Stats returns billing dashboard statistics.
func (uc *InvoiceUseCase) Stats(ctx context.Context) (*InvoiceStats, error) {
	revenue, count, avg, err := uc.invoiceRepo.GetTodayStats(ctx)
	if err != nil {
		return nil, err
	}

	return &InvoiceStats{
		TodayRevenue:  revenue,
		TodayInvoices: count,
		AvgBillValue:  avg,
	}, nil
}
