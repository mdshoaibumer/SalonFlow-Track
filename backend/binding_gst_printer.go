package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// GSTService exposes GST/tax operations to the Wails frontend.
type GSTService struct {
	ctx   context.Context
	uc    *usecase.GSTUseCase
	guard *PermissionGuard
}

func NewGSTService(uc *usecase.GSTUseCase) *GSTService {
	return &GSTService{uc: uc}
}

func (s *GSTService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *GSTService) GetSettings() (*domain.GSTSettings, error) {
	return s.uc.GetSettings(s.ctx)
}

func (s *GSTService) SaveSettings(settings *domain.GSTSettings) error {
	return s.uc.SaveSettings(s.ctx, settings)
}

func (s *GSTService) ListTaxRates(category string) ([]domain.TaxRate, error) {
	return s.uc.ListTaxRates(s.ctx, category)
}

func (s *GSTService) CreateTaxRate(rate *domain.TaxRate) error {
	return s.uc.CreateTaxRate(s.ctx, rate)
}

func (s *GSTService) UpdateTaxRate(rate *domain.TaxRate) error {
	return s.uc.UpdateTaxRate(s.ctx, rate)
}

func (s *GSTService) DeleteTaxRate(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.uc.DeleteTaxRate(s.ctx, uid)
}

func (s *GSTService) GetReport(filter domain.GSTReportFilter) (*domain.GSTReport, error) {
	return s.uc.GetReport(s.ctx, filter)
}

// PrinterService exposes printing operations to the Wails frontend.
type PrinterService struct {
	ctx   context.Context
	uc    *usecase.PrinterUseCase
	guard *PermissionGuard
}

func NewPrinterService(uc *usecase.PrinterUseCase) *PrinterService {
	return &PrinterService{uc: uc}
}

func (s *PrinterService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *PrinterService) GetSettings() (*domain.PrinterSettings, error) {
	return s.uc.GetSettings(s.ctx)
}

func (s *PrinterService) SaveSettings(settings *domain.PrinterSettings) error {
	return s.uc.SaveSettings(s.ctx, settings)
}

func (s *PrinterService) PrintInvoice(data *domain.ReceiptData) (*domain.PrintJob, string, error) {
	return s.uc.PrintInvoice(s.ctx, data)
}

func (s *PrinterService) PrintReceipt(data *domain.ReceiptData) (*domain.PrintJob, []byte, error) {
	return s.uc.PrintReceipt(s.ctx, data)
}

func (s *PrinterService) PrintTest() (*domain.PrintJob, []byte, error) {
	return s.uc.PrintTest(s.ctx)
}

func (s *PrinterService) ListPrintJobs(limit, offset int) ([]domain.PrintJob, int, error) {
	return s.uc.ListPrintJobs(s.ctx, limit, offset)
}

func (s *PrinterService) GetPrintJob(id string) (*domain.PrintJob, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.uc.GetPrintJob(s.ctx, uid)
}
