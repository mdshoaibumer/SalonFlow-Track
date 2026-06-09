package domain

import (
	"testing"
)

func TestNewCustomer(t *testing.T) {
	c := NewCustomer("John Doe", "9876543210")

	if c.FullName != "John Doe" {
		t.Errorf("expected 'John Doe', got %q", c.FullName)
	}
	if c.Phone != "9876543210" {
		t.Errorf("expected '9876543210', got %q", c.Phone)
	}
	if c.Status != CustomerStatusActive {
		t.Errorf("expected active, got %q", c.Status)
	}
	if c.CustomerCode == "" {
		t.Error("expected non-empty customer code")
	}
}

func TestCustomerValidate_Valid(t *testing.T) {
	c := NewCustomer("Jane Doe", "9876543210")
	if err := c.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestCustomerValidate_NameRequired(t *testing.T) {
	c := NewCustomer("", "9876543210")
	if err := c.Validate(); err != ErrCustomerNameRequired {
		t.Errorf("expected ErrCustomerNameRequired, got %v", err)
	}
}

func TestCustomerValidate_PhoneRequired(t *testing.T) {
	c := NewCustomer("John", "")
	if err := c.Validate(); err != ErrCustomerPhoneRequired {
		t.Errorf("expected ErrCustomerPhoneRequired, got %v", err)
	}
}

func TestCustomerValidate_PhoneInvalid(t *testing.T) {
	c := NewCustomer("John", "1234567890")
	if err := c.Validate(); err != ErrCustomerPhoneInvalid {
		t.Errorf("expected ErrCustomerPhoneInvalid, got %v", err)
	}
}

func TestCustomerRecordVisit(t *testing.T) {
	c := NewCustomer("John", "9876543210")
	c.RecordVisit(500, c.CreatedAt)

	if c.TotalVisits != 1 {
		t.Errorf("expected 1 visit, got %d", c.TotalVisits)
	}
	if c.TotalSpent != 500 {
		t.Errorf("expected 500 spent, got %f", c.TotalSpent)
	}
	if c.LastVisitDate == nil {
		t.Error("expected last visit date to be set")
	}
}

func TestCustomerDeactivate(t *testing.T) {
	c := NewCustomer("John", "9876543210")
	c.Deactivate()
	if c.Status != CustomerStatusInactive {
		t.Errorf("expected inactive, got %q", c.Status)
	}
}
