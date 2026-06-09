package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// PrinterRepository manages printer settings and print job history.
type PrinterRepository interface {
	GetSettings(ctx context.Context) (*domain.PrinterSettings, error)
	SaveSettings(ctx context.Context, settings *domain.PrinterSettings) error

	CreatePrintJob(ctx context.Context, job *domain.PrintJob) error
	UpdatePrintJobStatus(ctx context.Context, id uuid.UUID, status string) error
	ListPrintJobs(ctx context.Context, limit, offset int) ([]domain.PrintJob, int, error)
	GetPrintJob(ctx context.Context, id uuid.UUID) (*domain.PrintJob, error)
}

// PrintEngine handles receipt formatting and ESC/POS command generation.
type PrintEngine interface {
	FormatReceipt(data *domain.ReceiptData, paperWidth string) string
	GenerateESCPOS(data *domain.ReceiptData, paperWidth string) *domain.ESCPOSCommand
	FormatTestPage(printerName, paperWidth string) *domain.ESCPOSCommand
}
