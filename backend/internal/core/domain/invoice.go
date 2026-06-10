package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/pkg/uid"
)

// Invoice payment status constants.
const (
	PaymentStatusPending = "pending"
	PaymentStatusPaid    = "paid"
	PaymentStatusPartial = "partial"
)

// Payment method constants.
const (
	PaymentMethodCash         = "cash"
	PaymentMethodUPI          = "upi"
	PaymentMethodCard         = "card"
	PaymentMethodBankTransfer = "bank_transfer"
)

var validPaymentMethods = map[string]bool{
	PaymentMethodCash:         true,
	PaymentMethodUPI:          true,
	PaymentMethodCard:         true,
	PaymentMethodBankTransfer: true,
}

// Invoice represents a billing invoice.
type Invoice struct {
	ID            uuid.UUID     `json:"id"`
	InvoiceNumber string        `json:"invoice_number"`
	CustomerID    uuid.UUID     `json:"customer_id"`
	StaffID       uuid.UUID     `json:"staff_id"`
	Items         []InvoiceItem `json:"items"`
	Subtotal      float64       `json:"subtotal"`
	Discount      float64       `json:"discount"`
	Tax           float64       `json:"tax"`
	GrandTotal    float64       `json:"grand_total"`
	PaymentStatus string        `json:"payment_status"`
	PaymentMethod string        `json:"payment_method"`
	Notes         string        `json:"notes"`
	InvoiceDate   time.Time     `json:"invoice_date"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// InvoiceItem represents a line item on an invoice.
type InvoiceItem struct {
	ID                  uuid.UUID `json:"id"`
	InvoiceID           uuid.UUID `json:"invoice_id"`
	ServiceID           uuid.UUID `json:"service_id"`
	ServiceNameSnapshot string    `json:"service_name_snapshot"`
	Quantity            int       `json:"quantity"`
	UnitPrice           float64   `json:"unit_price"`
	Discount            float64   `json:"discount"`
	LineTotal           float64   `json:"line_total"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// Payment represents a payment entry for an invoice.
type Payment struct {
	ID              uuid.UUID `json:"id"`
	InvoiceID       uuid.UUID `json:"invoice_id"`
	Amount          float64   `json:"amount"`
	PaymentMethod   string    `json:"payment_method"`
	ReferenceNumber string    `json:"reference_number"`
	PaymentDate     time.Time `json:"payment_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// NewInvoice creates a new Invoice with UUIDv7.
func NewInvoice(customerID, staffID uuid.UUID, invoiceNumber string) *Invoice {
	now := time.Now().UTC()
	return &Invoice{
		ID:            uid.New(),
		InvoiceNumber: invoiceNumber,
		CustomerID:    customerID,
		StaffID:       staffID,
		PaymentStatus: PaymentStatusPending,
		InvoiceDate:   now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// NewInvoiceItem creates a new invoice item.
func NewInvoiceItem(invoiceID, serviceID uuid.UUID, serviceName string, quantity int, unitPrice, discount float64) *InvoiceItem {
	now := time.Now().UTC()
	lineTotal := (unitPrice * float64(quantity)) - discount
	if lineTotal < 0 {
		lineTotal = 0
	}
	return &InvoiceItem{
		ID:                  uid.New(),
		InvoiceID:           invoiceID,
		ServiceID:           serviceID,
		ServiceNameSnapshot: serviceName,
		Quantity:            quantity,
		UnitPrice:           unitPrice,
		Discount:            discount,
		LineTotal:           lineTotal,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
}

// NewPayment creates a new payment record.
func NewPayment(invoiceID uuid.UUID, amount float64, method, reference string) *Payment {
	now := time.Now().UTC()
	return &Payment{
		ID:              uid.New(),
		InvoiceID:       invoiceID,
		Amount:          amount,
		PaymentMethod:   method,
		ReferenceNumber: strings.TrimSpace(reference),
		PaymentDate:     now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// AddItem adds an item and recalculates totals.
func (inv *Invoice) AddItem(item *InvoiceItem) {
	inv.Items = append(inv.Items, *item)
	inv.Recalculate()
}

// Recalculate recalculates the invoice totals from items.
func (inv *Invoice) Recalculate() {
	var subtotal float64
	for _, item := range inv.Items {
		subtotal += item.LineTotal
	}
	inv.Subtotal = subtotal
	inv.GrandTotal = subtotal - inv.Discount + inv.Tax
	if inv.GrandTotal < 0 {
		inv.GrandTotal = 0
	}
	inv.UpdatedAt = time.Now().UTC()
}

// Validate checks invoice business rules.
func (inv *Invoice) Validate() error {
	if inv.CustomerID == uuid.Nil {
		return ErrInvoiceCustomerRequired
	}
	if inv.StaffID == uuid.Nil {
		return ErrInvoiceStaffRequired
	}
	if len(inv.Items) == 0 {
		return ErrInvoiceItemsRequired
	}
	if inv.Discount < 0 {
		return ErrInvoiceInvalidDiscount
	}
	if inv.Tax < 0 {
		return ErrInvoiceInvalidTax
	}
	return nil
}

// GenerateInvoiceNumber creates a unique invoice number: INV-YYYY-NNNNNN
func GenerateInvoiceNumber(year int, sequence int) string {
	return fmt.Sprintf("INV-%d-%06d", year, sequence)
}

// ValidatePayment checks payment business rules.
func (p *Payment) Validate() error {
	if p.Amount <= 0 {
		return ErrPaymentInvalidAmount
	}
	if !validPaymentMethods[p.PaymentMethod] {
		return ErrPaymentInvalidMethod
	}
	return nil
}

// MarkPaid marks the invoice as fully paid.
func (inv *Invoice) MarkPaid(method string) {
	inv.PaymentStatus = PaymentStatusPaid
	inv.PaymentMethod = method
	inv.UpdatedAt = time.Now().UTC()
}

// MarkPartial marks the invoice as partially paid.
func (inv *Invoice) MarkPartial() {
	inv.PaymentStatus = PaymentStatusPartial
	inv.UpdatedAt = time.Now().UTC()
}

// IsPaid returns true if fully paid.
func (inv *Invoice) IsPaid() bool {
	return inv.PaymentStatus == PaymentStatusPaid
}
