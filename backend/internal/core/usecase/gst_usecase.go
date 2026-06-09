package usecase

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// GSTUseCase handles GST business logic.
type GSTUseCase struct {
	repo   ports.GSTRepository
	engine ports.GSTEngine
	log    *slog.Logger
}

// NewGSTUseCase creates a new GSTUseCase.
func NewGSTUseCase(repo ports.GSTRepository, engine ports.GSTEngine, log *slog.Logger) *GSTUseCase {
	return &GSTUseCase{repo: repo, engine: engine, log: log}
}

// GetSettings retrieves GST settings.
func (uc *GSTUseCase) GetSettings(ctx context.Context) (*domain.GSTSettings, error) {
	settings, err := uc.repo.GetSettings(ctx)
	if err != nil {
		// If not found, return default
		if apperror.Is(err, apperror.KindNotFound) {
			return domain.NewGSTSettings(), nil
		}
		return nil, err
	}
	return settings, nil
}

// SaveSettings saves GST settings.
func (uc *GSTUseCase) SaveSettings(ctx context.Context, settings *domain.GSTSettings) error {
	if settings.ID == uuid.Nil {
		settings.ID = domain.NewGSTSettings().ID
	}
	return uc.repo.SaveSettings(ctx, settings)
}

// CreateTaxRate creates a new tax rate.
func (uc *GSTUseCase) CreateTaxRate(ctx context.Context, rate *domain.TaxRate) error {
	if rate.Name == "" {
		return apperror.Validation("name", "Name is required")
	}
	return uc.repo.CreateTaxRate(ctx, rate)
}

// UpdateTaxRate updates a tax rate.
func (uc *GSTUseCase) UpdateTaxRate(ctx context.Context, rate *domain.TaxRate) error {
	if rate.Name == "" {
		return apperror.Validation("name", "Name is required")
	}
	return uc.repo.UpdateTaxRate(ctx, rate)
}

// DeleteTaxRate deletes a tax rate.
func (uc *GSTUseCase) DeleteTaxRate(ctx context.Context, id uuid.UUID) error {
	return uc.repo.DeleteTaxRate(ctx, id)
}

// ListTaxRates lists tax rates.
func (uc *GSTUseCase) ListTaxRates(ctx context.Context, category string) ([]domain.TaxRate, error) {
	return uc.repo.ListTaxRates(ctx, category)
}

// CalculateInvoiceTax computes GST for an invoice.
func (uc *GSTUseCase) CalculateInvoiceTax(ctx context.Context, invoiceID uuid.UUID, items []domain.InvoiceItem, isInterstate bool) (*domain.GSTInvoiceSummary, error) {
	settings, err := uc.GetSettings(ctx)
	if err != nil {
		return nil, err
	}
	if !settings.IsGSTEnabled {
		return &domain.GSTInvoiceSummary{InvoiceID: invoiceID}, nil
	}

	summary := uc.engine.CalculateInvoiceTax(invoiceID, items, settings, isInterstate)

	// Persist tax lines
	if len(summary.TaxLines) > 0 {
		if err := uc.repo.CreateTaxLines(ctx, summary.TaxLines); err != nil {
			return nil, err
		}
	}

	return summary, nil
}

// GetInvoiceTaxLines retrieves tax lines for an invoice.
func (uc *GSTUseCase) GetInvoiceTaxLines(ctx context.Context, invoiceID uuid.UUID) ([]domain.InvoiceTaxLine, error) {
	return uc.repo.GetTaxLinesByInvoice(ctx, invoiceID)
}

// GetReport generates a GST report.
func (uc *GSTUseCase) GetReport(ctx context.Context, filter domain.GSTReportFilter) (*domain.GSTReport, error) {
	if filter.StartDate == "" || filter.EndDate == "" {
		return nil, apperror.Validation("dates", "Start date and end date are required")
	}
	return uc.repo.GetGSTReport(ctx, filter)
}
