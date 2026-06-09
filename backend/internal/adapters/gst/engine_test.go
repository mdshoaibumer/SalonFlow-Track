package gst

import (
	"testing"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
)

func TestEngine_CalculateTax_IntraState(t *testing.T) {
	e := NewEngine()
	cgst, sgst, igst, total := e.CalculateTax(1000, false, 9, 9, 18)

	if cgst != 90 {
		t.Errorf("CGST = %.2f, want 90", cgst)
	}
	if sgst != 90 {
		t.Errorf("SGST = %.2f, want 90", sgst)
	}
	if igst != 0 {
		t.Errorf("IGST = %.2f, want 0", igst)
	}
	if total != 180 {
		t.Errorf("Total = %.2f, want 180", total)
	}
}

func TestEngine_CalculateTax_InterState(t *testing.T) {
	e := NewEngine()
	cgst, sgst, igst, total := e.CalculateTax(1000, true, 9, 9, 18)

	if cgst != 0 {
		t.Errorf("CGST = %.2f, want 0", cgst)
	}
	if sgst != 0 {
		t.Errorf("SGST = %.2f, want 0", sgst)
	}
	if igst != 180 {
		t.Errorf("IGST = %.2f, want 180", igst)
	}
	if total != 180 {
		t.Errorf("Total = %.2f, want 180", total)
	}
}

func TestEngine_CalculateInvoiceTax(t *testing.T) {
	e := NewEngine()

	invoiceID := uuid.New()
	items := []domain.InvoiceItem{
		{ID: uuid.New(), LineTotal: 500},
		{ID: uuid.New(), LineTotal: 1000},
		{ID: uuid.New(), LineTotal: 300},
	}

	settings := &domain.GSTSettings{
		CGSTRate: 9,
		SGSTRate: 9,
		IGSTRate: 18,
		HSNCode:  "9902",
	}

	summary := e.CalculateInvoiceTax(invoiceID, items, settings, false)

	if summary.TaxableAmount != 1800 {
		t.Errorf("TaxableAmount = %.2f, want 1800", summary.TaxableAmount)
	}
	if summary.TotalCGST != 162 {
		t.Errorf("TotalCGST = %.2f, want 162", summary.TotalCGST)
	}
	if summary.TotalSGST != 162 {
		t.Errorf("TotalSGST = %.2f, want 162", summary.TotalSGST)
	}
	if summary.TotalIGST != 0 {
		t.Errorf("TotalIGST = %.2f, want 0", summary.TotalIGST)
	}
	if summary.GrandTotal != 2124 {
		t.Errorf("GrandTotal = %.2f, want 2124", summary.GrandTotal)
	}
	if len(summary.TaxLines) != 3 {
		t.Errorf("TaxLines len = %d, want 3", len(summary.TaxLines))
	}
}

func TestEngine_CalculateInvoiceTax_Interstate(t *testing.T) {
	e := NewEngine()

	invoiceID := uuid.New()
	items := []domain.InvoiceItem{
		{ID: uuid.New(), LineTotal: 1000},
	}

	settings := &domain.GSTSettings{
		CGSTRate: 9,
		SGSTRate: 9,
		IGSTRate: 18,
		HSNCode:  "9902",
	}

	summary := e.CalculateInvoiceTax(invoiceID, items, settings, true)

	if summary.TotalCGST != 0 {
		t.Errorf("TotalCGST = %.2f, want 0", summary.TotalCGST)
	}
	if summary.TotalIGST != 180 {
		t.Errorf("TotalIGST = %.2f, want 180", summary.TotalIGST)
	}
	if summary.GrandTotal != 1180 {
		t.Errorf("GrandTotal = %.2f, want 1180", summary.GrandTotal)
	}
}
