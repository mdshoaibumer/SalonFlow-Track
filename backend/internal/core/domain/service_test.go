package domain

import (
	"testing"
)

func TestNewService(t *testing.T) {
	svc := NewService("Hair Cut", CategoryHair, 30, 300)

	if svc.Name != "Hair Cut" {
		t.Errorf("expected name 'Hair Cut', got %q", svc.Name)
	}
	if svc.Category != CategoryHair {
		t.Errorf("expected category 'hair', got %q", svc.Category)
	}
	if svc.DurationMinutes != 30 {
		t.Errorf("expected duration 30, got %d", svc.DurationMinutes)
	}
	if svc.Price != 300 {
		t.Errorf("expected price 300, got %f", svc.Price)
	}
	if svc.Status != ServiceStatusActive {
		t.Errorf("expected active status, got %q", svc.Status)
	}
	if svc.ServiceCode == "" {
		t.Error("expected non-empty service code")
	}
}

func TestServiceValidate_Valid(t *testing.T) {
	svc := NewService("Hair Cut", CategoryHair, 30, 300)
	if err := svc.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestServiceValidate_NameRequired(t *testing.T) {
	svc := NewService("", CategoryHair, 30, 300)
	if err := svc.Validate(); err != ErrServiceNameRequired {
		t.Errorf("expected ErrServiceNameRequired, got %v", err)
	}
}

func TestServiceValidate_InvalidCategory(t *testing.T) {
	svc := NewService("Test", "invalid", 30, 300)
	if err := svc.Validate(); err != ErrServiceInvalidCategory {
		t.Errorf("expected ErrServiceInvalidCategory, got %v", err)
	}
}

func TestServiceValidate_InvalidDuration(t *testing.T) {
	svc := NewService("Test", CategoryHair, 0, 300)
	if err := svc.Validate(); err != ErrServiceInvalidDuration {
		t.Errorf("expected ErrServiceInvalidDuration, got %v", err)
	}
}

func TestServiceValidate_InvalidPrice(t *testing.T) {
	svc := NewService("Test", CategoryHair, 30, 0)
	if err := svc.Validate(); err != ErrServiceInvalidPrice {
		t.Errorf("expected ErrServiceInvalidPrice, got %v", err)
	}
}

func TestServiceValidate_InvalidCommissionPercentage(t *testing.T) {
	svc := NewService("Test", CategoryHair, 30, 300)
	svc.CommissionType = CommissionTypePercentage
	svc.CommissionValue = 150
	if err := svc.Validate(); err != ErrServiceInvalidCommissionValue {
		t.Errorf("expected ErrServiceInvalidCommissionValue, got %v", err)
	}
}

func TestServiceDeactivate(t *testing.T) {
	svc := NewService("Test", CategoryHair, 30, 300)
	svc.Deactivate()
	if svc.Status != ServiceStatusInactive {
		t.Errorf("expected inactive, got %q", svc.Status)
	}
	if svc.IsActive() {
		t.Error("expected IsActive to return false")
	}
}

func TestServiceActivate(t *testing.T) {
	svc := NewService("Test", CategoryHair, 30, 300)
	svc.Deactivate()
	svc.Activate()
	if svc.Status != ServiceStatusActive {
		t.Errorf("expected active, got %q", svc.Status)
	}
}
