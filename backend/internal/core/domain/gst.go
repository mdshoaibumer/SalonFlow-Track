package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/pkg/uid"
)

// GSTSettings represents the business GST configuration.
type GSTSettings struct {
	ID           uuid.UUID `json:"id"`
	BusinessName string    `json:"business_name"`
	GSTIN        string    `json:"gstin"`
	State        string    `json:"state"`
	Address      string    `json:"address"`
	HSNCode      string    `json:"hsn_code"`
	CGSTRate     float64   `json:"cgst_rate"`
	SGSTRate     float64   `json:"sgst_rate"`
	IGSTRate     float64   `json:"igst_rate"`
	IsGSTEnabled bool      `json:"is_gst_enabled"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewGSTSettings creates a default GST settings entry.
func NewGSTSettings() *GSTSettings {
	now := time.Now().UTC()
	return &GSTSettings{
		ID:        uid.New(),
		CGSTRate:  9.0,
		SGSTRate:  9.0,
		IGSTRate:  18.0,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// TaxRate represents a specific tax rate (can override per category).
type TaxRate struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	HSNCode   string    `json:"hsn_code"`
	CGSTRate  float64   `json:"cgst_rate"`
	SGSTRate  float64   `json:"sgst_rate"`
	IGSTRate  float64   `json:"igst_rate"`
	Category  string    `json:"category"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewTaxRate creates a new tax rate.
func NewTaxRate(name, hsnCode, category string, cgst, sgst, igst float64) *TaxRate {
	now := time.Now().UTC()
	return &TaxRate{
		ID:        uid.New(),
		Name:      name,
		HSNCode:   hsnCode,
		CGSTRate:  cgst,
		SGSTRate:  sgst,
		IGSTRate:  igst,
		Category:  category,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// InvoiceTaxLine represents a tax breakdown for an invoice line item.
type InvoiceTaxLine struct {
	ID            uuid.UUID `json:"id"`
	InvoiceID     uuid.UUID `json:"invoice_id"`
	ItemID        string    `json:"item_id"`
	TaxableAmount float64   `json:"taxable_amount"`
	CGSTRate      float64   `json:"cgst_rate"`
	CGSTAmount    float64   `json:"cgst_amount"`
	SGSTRate      float64   `json:"sgst_rate"`
	SGSTAmount    float64   `json:"sgst_amount"`
	IGSTRate      float64   `json:"igst_rate"`
	IGSTAmount    float64   `json:"igst_amount"`
	TotalTax      float64   `json:"total_tax"`
	IsInterstate  bool      `json:"is_interstate"`
	HSNCode       string    `json:"hsn_code"`
	CreatedAt     time.Time `json:"created_at"`
}

// NewInvoiceTaxLine creates a new invoice tax line.
func NewInvoiceTaxLine(invoiceID uuid.UUID, itemID string, taxableAmount float64, isInterstate bool, cgstRate, sgstRate, igstRate float64, hsnCode string) *InvoiceTaxLine {
	line := &InvoiceTaxLine{
		ID:            uid.New(),
		InvoiceID:     invoiceID,
		ItemID:        itemID,
		TaxableAmount: taxableAmount,
		IsInterstate:  isInterstate,
		HSNCode:       hsnCode,
		CreatedAt:     time.Now().UTC(),
	}

	if isInterstate {
		line.IGSTRate = igstRate
		line.IGSTAmount = taxableAmount * igstRate / 100
		line.TotalTax = line.IGSTAmount
	} else {
		line.CGSTRate = cgstRate
		line.CGSTAmount = taxableAmount * cgstRate / 100
		line.SGSTRate = sgstRate
		line.SGSTAmount = taxableAmount * sgstRate / 100
		line.TotalTax = line.CGSTAmount + line.SGSTAmount
	}

	return line
}

// GSTInvoiceSummary represents the tax summary for an invoice.
type GSTInvoiceSummary struct {
	InvoiceID     uuid.UUID        `json:"invoice_id"`
	TaxableAmount float64          `json:"taxable_amount"`
	TotalCGST     float64          `json:"total_cgst"`
	TotalSGST     float64          `json:"total_sgst"`
	TotalIGST     float64          `json:"total_igst"`
	TotalTax      float64          `json:"total_tax"`
	GrandTotal    float64          `json:"grand_total"`
	TaxLines      []InvoiceTaxLine `json:"tax_lines"`
}

// GSTReport represents a GST report (daily/monthly).
type GSTReport struct {
	Period        string  `json:"period"`
	StartDate     string  `json:"start_date"`
	EndDate       string  `json:"end_date"`
	TotalInvoices int     `json:"total_invoices"`
	TaxableAmount float64 `json:"taxable_amount"`
	TotalCGST     float64 `json:"total_cgst"`
	TotalSGST     float64 `json:"total_sgst"`
	TotalIGST     float64 `json:"total_igst"`
	TotalTax      float64 `json:"total_tax"`
	GrandTotal    float64 `json:"grand_total"`
}

// GSTReportFilter defines filters for GST report generation.
type GSTReportFilter struct {
	Period    string `json:"period"` // daily, monthly
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
