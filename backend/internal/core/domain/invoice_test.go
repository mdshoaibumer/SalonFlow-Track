package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewInvoice(t *testing.T) {
	customerID := uuid.New()
	staffID := uuid.New()
	inv := NewInvoice(customerID, staffID, "INV-2026-000001")

	if inv.CustomerID != customerID {
		t.Error("customer ID mismatch")
	}
	if inv.StaffID != staffID {
		t.Error("staff ID mismatch")
	}
	if inv.InvoiceNumber != "INV-2026-000001" {
		t.Errorf("expected INV-2026-000001, got %q", inv.InvoiceNumber)
	}
	if inv.PaymentStatus != PaymentStatusPending {
		t.Errorf("expected pending, got %q", inv.PaymentStatus)
	}
}

func TestInvoiceAddItem(t *testing.T) {
	customerID := uuid.New()
	staffID := uuid.New()
	inv := NewInvoice(customerID, staffID, "INV-2026-000001")

	item := NewInvoiceItem(inv.ID, uuid.New(), "Hair Cut", 1, 300, 0)
	inv.AddItem(item)

	if len(inv.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(inv.Items))
	}
	if inv.Subtotal != 300 {
		t.Errorf("expected subtotal 300, got %f", inv.Subtotal)
	}
	if inv.GrandTotal != 300 {
		t.Errorf("expected grand total 300, got %f", inv.GrandTotal)
	}
}

func TestInvoiceRecalculate(t *testing.T) {
	customerID := uuid.New()
	staffID := uuid.New()
	inv := NewInvoice(customerID, staffID, "INV-2026-000001")

	inv.AddItem(NewInvoiceItem(inv.ID, uuid.New(), "Hair Cut", 1, 300, 0))
	inv.AddItem(NewInvoiceItem(inv.ID, uuid.New(), "Facial", 1, 800, 0))

	inv.Discount = 100
	inv.Tax = 50
	inv.Recalculate()

	if inv.Subtotal != 1100 {
		t.Errorf("expected subtotal 1100, got %f", inv.Subtotal)
	}
	// GrandTotal = 1100 - 100 + 50 = 1050
	if inv.GrandTotal != 1050 {
		t.Errorf("expected grand total 1050, got %f", inv.GrandTotal)
	}
}

func TestInvoiceValidate_MissingCustomer(t *testing.T) {
	inv := NewInvoice(uuid.Nil, uuid.New(), "INV-2026-000001")
	inv.Items = []InvoiceItem{{}}
	if err := inv.Validate(); err != ErrInvoiceCustomerRequired {
		t.Errorf("expected ErrInvoiceCustomerRequired, got %v", err)
	}
}

func TestInvoiceValidate_MissingStaff(t *testing.T) {
	inv := NewInvoice(uuid.New(), uuid.Nil, "INV-2026-000001")
	inv.Items = []InvoiceItem{{}}
	if err := inv.Validate(); err != ErrInvoiceStaffRequired {
		t.Errorf("expected ErrInvoiceStaffRequired, got %v", err)
	}
}

func TestInvoiceValidate_NoItems(t *testing.T) {
	inv := NewInvoice(uuid.New(), uuid.New(), "INV-2026-000001")
	if err := inv.Validate(); err != ErrInvoiceItemsRequired {
		t.Errorf("expected ErrInvoiceItemsRequired, got %v", err)
	}
}

func TestInvoiceMarkPaid(t *testing.T) {
	inv := NewInvoice(uuid.New(), uuid.New(), "INV-2026-000001")
	inv.MarkPaid(PaymentMethodCash)
	if inv.PaymentStatus != PaymentStatusPaid {
		t.Errorf("expected paid, got %q", inv.PaymentStatus)
	}
	if !inv.IsPaid() {
		t.Error("expected IsPaid true")
	}
}

func TestGenerateInvoiceNumber(t *testing.T) {
	num := GenerateInvoiceNumber(2026, 1)
	if num != "INV-2026-000001" {
		t.Errorf("expected INV-2026-000001, got %q", num)
	}

	num = GenerateInvoiceNumber(2026, 999)
	if num != "INV-2026-000999" {
		t.Errorf("expected INV-2026-000999, got %q", num)
	}
}

func TestPaymentValidate(t *testing.T) {
	p := NewPayment(uuid.New(), 100, PaymentMethodCash, "")
	if err := p.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestPaymentValidate_InvalidAmount(t *testing.T) {
	p := NewPayment(uuid.New(), 0, PaymentMethodCash, "")
	if err := p.Validate(); err != ErrPaymentInvalidAmount {
		t.Errorf("expected ErrPaymentInvalidAmount, got %v", err)
	}
}

func TestPaymentValidate_InvalidMethod(t *testing.T) {
	p := NewPayment(uuid.New(), 100, "bitcoin", "")
	if err := p.Validate(); err != ErrPaymentInvalidMethod {
		t.Errorf("expected ErrPaymentInvalidMethod, got %v", err)
	}
}
