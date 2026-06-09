package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewStaffPerformanceDaily(t *testing.T) {
	staffID := uuid.New()
	perf := NewStaffPerformanceDaily(staffID, "2026-06-09")

	if perf.ID == uuid.Nil {
		t.Fatal("expected non-nil ID")
	}
	if perf.StaffID != staffID {
		t.Fatalf("expected staff ID %s, got %s", staffID, perf.StaffID)
	}
	if perf.BusinessDate != "2026-06-09" {
		t.Fatalf("expected business date 2026-06-09, got %s", perf.BusinessDate)
	}
}

func TestStaffPerformanceDaily_AddInvoice(t *testing.T) {
	perf := NewStaffPerformanceDaily(uuid.New(), "2026-06-09")

	perf.AddInvoice(2000, 3, 200)
	perf.AddInvoice(1500, 2, 150)
	perf.AddInvoice(500, 1, 50)

	if perf.InvoiceCount != 3 {
		t.Fatalf("expected 3 invoices, got %d", perf.InvoiceCount)
	}
	if perf.CustomerCount != 3 {
		t.Fatalf("expected 3 customers, got %d", perf.CustomerCount)
	}
	if perf.ServiceCount != 6 {
		t.Fatalf("expected 6 services, got %d", perf.ServiceCount)
	}
	if perf.Revenue != 4000 {
		t.Fatalf("expected revenue 4000, got %.2f", perf.Revenue)
	}
	if perf.CommissionAmount != 400 {
		t.Fatalf("expected commission 400, got %.2f", perf.CommissionAmount)
	}
}
