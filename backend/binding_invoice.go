package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// InvoiceService exposes invoice operations to the Wails frontend.
type InvoiceService struct {
	ctx      context.Context
	uc       *usecase.InvoiceUseCase
	guard    *PermissionGuard
	licGuard *LicenseGuard
}

func NewInvoiceService(uc *usecase.InvoiceUseCase) *InvoiceService {
	return &InvoiceService{uc: uc}
}

func (s *InvoiceService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *InvoiceService) ListInvoices(input usecase.ListInvoiceInput) (*usecase.ListInvoiceOutput, error) {
	if err := s.guard.Require("billing.view"); err != nil {
		return nil, err
	}
	return s.uc.List(s.ctx, input)
}

func (s *InvoiceService) GetInvoice(id string) (*domain.Invoice, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.uc.GetByID(s.ctx, uid)
}

func (s *InvoiceService) CreateInvoice(input usecase.CreateInvoiceInput) (*domain.Invoice, error) {
	if err := s.guard.Require("billing.create"); err != nil {
		return nil, err
	}
	if err := s.licGuard.RequireActive(domain.OpInvoiceCreate); err != nil {
		return nil, err
	}
	return s.uc.Create(s.ctx, input)
}

func (s *InvoiceService) RecordPayment(invoiceID string, input usecase.RecordPaymentInput) (*domain.Payment, error) {
	uid, err := uuid.Parse(invoiceID)
	if err != nil {
		return nil, err
	}
	return s.uc.RecordPayment(s.ctx, uid, input)
}

func (s *InvoiceService) GetInvoiceStats() (*usecase.InvoiceStats, error) {
	return s.uc.Stats(s.ctx)
}
