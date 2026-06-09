package gst

import (
	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// Engine handles GST tax calculation logic.
type Engine struct{}

// NewEngine creates a new GST calculation engine.
func NewEngine() *Engine {
	return &Engine{}
}

// CalculateTax computes tax amounts for a given taxable amount.
func (e *Engine) CalculateTax(taxableAmount float64, isInterstate bool, cgstRate, sgstRate, igstRate float64) (cgst, sgst, igst, total float64) {
	if isInterstate {
		igst = roundTo2(taxableAmount * igstRate / 100)
		total = igst
		return 0, 0, igst, total
	}
	cgst = roundTo2(taxableAmount * cgstRate / 100)
	sgst = roundTo2(taxableAmount * sgstRate / 100)
	total = cgst + sgst
	return cgst, sgst, 0, total
}

// CalculateInvoiceTax computes tax for all items in an invoice.
func (e *Engine) CalculateInvoiceTax(invoiceID uuid.UUID, items []domain.InvoiceItem, settings *domain.GSTSettings, isInterstate bool) *domain.GSTInvoiceSummary {
	summary := &domain.GSTInvoiceSummary{
		InvoiceID: invoiceID,
		TaxLines:  make([]domain.InvoiceTaxLine, 0, len(items)),
	}

	for _, item := range items {
		taxableAmount := item.LineTotal
		line := domain.NewInvoiceTaxLine(
			invoiceID,
			item.ID.String(),
			taxableAmount,
			isInterstate,
			settings.CGSTRate,
			settings.SGSTRate,
			settings.IGSTRate,
			settings.HSNCode,
		)

		summary.TaxableAmount += taxableAmount
		summary.TotalCGST += line.CGSTAmount
		summary.TotalSGST += line.SGSTAmount
		summary.TotalIGST += line.IGSTAmount
		summary.TotalTax += line.TotalTax
		summary.TaxLines = append(summary.TaxLines, *line)
	}

	summary.GrandTotal = summary.TaxableAmount + summary.TotalTax
	return summary
}

func roundTo2(val float64) float64 {
	return float64(int(val*100+0.5)) / 100
}
