package usecase

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// PrinterUseCase handles printing business logic.
type PrinterUseCase struct {
	repo   ports.PrinterRepository
	engine ports.PrintEngine
	log    *slog.Logger
}

// NewPrinterUseCase creates a new PrinterUseCase.
func NewPrinterUseCase(repo ports.PrinterRepository, engine ports.PrintEngine, log *slog.Logger) *PrinterUseCase {
	return &PrinterUseCase{repo: repo, engine: engine, log: log}
}

// GetSettings retrieves printer settings.
func (uc *PrinterUseCase) GetSettings(ctx context.Context) (*domain.PrinterSettings, error) {
	settings, err := uc.repo.GetSettings(ctx)
	if err != nil {
		if apperror.Is(err, apperror.KindNotFound) {
			return domain.NewPrinterSettings(), nil
		}
		return nil, err
	}
	return settings, nil
}

// SaveSettings saves printer settings.
func (uc *PrinterUseCase) SaveSettings(ctx context.Context, settings *domain.PrinterSettings) error {
	if settings.ID == uuid.Nil {
		settings.ID = domain.NewPrinterSettings().ID
	}
	return uc.repo.SaveSettings(ctx, settings)
}

// PrintInvoice creates a print job and generates receipt content.
func (uc *PrinterUseCase) PrintInvoice(ctx context.Context, data *domain.ReceiptData) (*domain.PrintJob, string, error) {
	settings, err := uc.GetSettings(ctx)
	if err != nil {
		return nil, "", err
	}

	job := domain.NewPrintJob(domain.PrintDocInvoice, data.InvoiceNumber, settings.DefaultPrinter, settings.PaperWidth, 1)
	if err := uc.repo.CreatePrintJob(ctx, job); err != nil {
		return nil, "", err
	}

	// Apply settings
	data.FooterText = settings.FooterText
	data.UPIID = settings.UPIID

	receipt := uc.engine.FormatReceipt(data, settings.PaperWidth)

	// Mark as completed (in real system, would send to printer)
	job.Status = domain.PrintStatusCompleted
	_ = uc.repo.UpdatePrintJobStatus(ctx, job.ID, domain.PrintStatusCompleted)

	return job, receipt, nil
}

// PrintReceipt creates a thermal receipt print job.
func (uc *PrinterUseCase) PrintReceipt(ctx context.Context, data *domain.ReceiptData) (*domain.PrintJob, []byte, error) {
	settings, err := uc.GetSettings(ctx)
	if err != nil {
		return nil, nil, err
	}

	job := domain.NewPrintJob(domain.PrintDocReceipt, data.InvoiceNumber, settings.DefaultPrinter, settings.PaperWidth, 1)
	if err := uc.repo.CreatePrintJob(ctx, job); err != nil {
		return nil, nil, err
	}

	data.FooterText = settings.FooterText
	data.UPIID = settings.UPIID

	escpos := uc.engine.GenerateESCPOS(data, settings.PaperWidth)

	job.Status = domain.PrintStatusCompleted
	_ = uc.repo.UpdatePrintJobStatus(ctx, job.ID, domain.PrintStatusCompleted)

	return job, escpos.Commands, nil
}

// PrintTest generates a test page.
func (uc *PrinterUseCase) PrintTest(ctx context.Context) (*domain.PrintJob, []byte, error) {
	settings, err := uc.GetSettings(ctx)
	if err != nil {
		return nil, nil, err
	}

	job := domain.NewPrintJob("receipt", "test", settings.DefaultPrinter, settings.PaperWidth, 1)
	if err := uc.repo.CreatePrintJob(ctx, job); err != nil {
		return nil, nil, err
	}

	testPage := uc.engine.FormatTestPage(settings.DefaultPrinter, settings.PaperWidth)

	job.Status = domain.PrintStatusCompleted
	_ = uc.repo.UpdatePrintJobStatus(ctx, job.ID, domain.PrintStatusCompleted)

	return job, testPage.Commands, nil
}

// ListPrintJobs lists print history.
func (uc *PrinterUseCase) ListPrintJobs(ctx context.Context, limit, offset int) ([]domain.PrintJob, int, error) {
	return uc.repo.ListPrintJobs(ctx, limit, offset)
}

// GetPrintJob gets a print job by ID.
func (uc *PrinterUseCase) GetPrintJob(ctx context.Context, id uuid.UUID) (*domain.PrintJob, error) {
	return uc.repo.GetPrintJob(ctx, id)
}
