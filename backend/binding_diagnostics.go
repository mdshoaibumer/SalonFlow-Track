package main

import (
	"context"

	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// DiagnosticsService exposes diagnostics operations to the Wails frontend.
type DiagnosticsService struct {
	ctx context.Context
	uc  *usecase.DiagnosticsUseCase
}

func NewDiagnosticsService(uc *usecase.DiagnosticsUseCase) *DiagnosticsService {
	return &DiagnosticsService{uc: uc}
}

func (s *DiagnosticsService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

// GetDiagnostics returns system health and configuration info.
func (s *DiagnosticsService) GetDiagnostics() (*usecase.DiagnosticsInfo, error) {
	return s.uc.GetDiagnostics(s.ctx)
}

// ExportDiagnosticsBundle creates a ZIP with logs and metadata, returns the path.
func (s *DiagnosticsService) ExportDiagnosticsBundle() (string, error) {
	return s.uc.ExportDiagnosticsBundle(s.ctx)
}
