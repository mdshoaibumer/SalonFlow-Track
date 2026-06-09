package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// GSTRepository manages GST settings, tax rates, and invoice tax lines.
type GSTRepository interface {
	// Settings
	GetSettings(ctx context.Context) (*domain.GSTSettings, error)
	SaveSettings(ctx context.Context, settings *domain.GSTSettings) error

	// Tax Rates
	CreateTaxRate(ctx context.Context, rate *domain.TaxRate) error
	UpdateTaxRate(ctx context.Context, rate *domain.TaxRate) error
	DeleteTaxRate(ctx context.Context, id uuid.UUID) error
	GetTaxRate(ctx context.Context, id uuid.UUID) (*domain.TaxRate, error)
	ListTaxRates(ctx context.Context, category string) ([]domain.TaxRate, error)
	GetTaxRateByCategory(ctx context.Context, category string) (*domain.TaxRate, error)

	// Invoice Tax Lines
	CreateTaxLines(ctx context.Context, lines []domain.InvoiceTaxLine) error
	GetTaxLinesByInvoice(ctx context.Context, invoiceID uuid.UUID) ([]domain.InvoiceTaxLine, error)

	// Reports
	GetGSTReport(ctx context.Context, filter domain.GSTReportFilter) (*domain.GSTReport, error)
}

// GSTEngine handles GST calculation logic.
type GSTEngine interface {
	CalculateTax(taxableAmount float64, isInterstate bool, cgstRate, sgstRate, igstRate float64) (cgst, sgst, igst, total float64)
	CalculateInvoiceTax(invoiceID uuid.UUID, items []domain.InvoiceItem, settings *domain.GSTSettings, isInterstate bool) *domain.GSTInvoiceSummary
}
