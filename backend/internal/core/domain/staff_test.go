package domain

import "testing"

func TestNewStaff(t *testing.T) {
	s := NewStaff("Nazim Khan", "9876543210", "stylist", 15000, 10)

	if s.FullName != "Nazim Khan" {
		t.Errorf("expected name Nazim Khan, got %s", s.FullName)
	}
	if s.Phone != "9876543210" {
		t.Errorf("expected phone, got %s", s.Phone)
	}
	if s.Designation != "stylist" {
		t.Errorf("expected designation stylist, got %s", s.Designation)
	}
	if s.BaseSalary != 15000 {
		t.Errorf("expected salary 15000, got %f", s.BaseSalary)
	}
	if s.CommissionPercentage != 10 {
		t.Errorf("expected commission 10, got %f", s.CommissionPercentage)
	}
	if s.Status != StaffStatusActive {
		t.Error("expected status active")
	}
	if s.StaffCode == "" {
		t.Error("expected staff code to be generated")
	}
	if s.ID.String() == "00000000-0000-0000-0000-000000000000" {
		t.Error("expected non-nil UUID")
	}
}

func TestStaff_Validate_Valid(t *testing.T) {
	s := NewStaff("Ravi Kumar", "9111111111", "stylist", 10000, 15)
	if err := s.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestStaff_Validate_EmptyName(t *testing.T) {
	s := NewStaff("", "9111111111", "stylist", 10000, 0)
	if err := s.Validate(); err != ErrStaffNameRequired {
		t.Errorf("expected ErrStaffNameRequired, got %v", err)
	}
}

func TestStaff_Validate_EmptyPhone(t *testing.T) {
	s := NewStaff("Ravi", "", "stylist", 10000, 0)
	if err := s.Validate(); err != ErrStaffPhoneRequired {
		t.Errorf("expected ErrStaffPhoneRequired, got %v", err)
	}
}

func TestStaff_Validate_InvalidPhone(t *testing.T) {
	s := NewStaff("Ravi", "1234567890", "stylist", 10000, 0)
	if err := s.Validate(); err != ErrStaffPhoneInvalid {
		t.Errorf("expected ErrStaffPhoneInvalid, got %v", err)
	}
}

func TestStaff_Validate_InvalidDesignation(t *testing.T) {
	s := NewStaff("Ravi", "9111111111", "owner", 10000, 0)
	if err := s.Validate(); err != ErrStaffInvalidDesignation {
		t.Errorf("expected ErrStaffInvalidDesignation, got %v", err)
	}
}

func TestStaff_Validate_NegativeSalary(t *testing.T) {
	s := NewStaff("Ravi", "9111111111", "stylist", -500, 0)
	if err := s.Validate(); err != ErrStaffInvalidSalary {
		t.Errorf("expected ErrStaffInvalidSalary, got %v", err)
	}
}

func TestStaff_Validate_InvalidCommission(t *testing.T) {
	s := NewStaff("Ravi", "9111111111", "stylist", 10000, 101)
	if err := s.Validate(); err != ErrStaffInvalidCommission {
		t.Errorf("expected ErrStaffInvalidCommission, got %v", err)
	}
}

func TestStaff_Deactivate(t *testing.T) {
	s := NewStaff("Test", "9999999999", "assistant", 8000, 5)
	s.Deactivate()

	if s.Status != StaffStatusInactive {
		t.Error("expected status inactive after deactivate")
	}
	if s.IsActive() {
		t.Error("expected IsActive() to return false")
	}
}

func TestStaff_Activate(t *testing.T) {
	s := NewStaff("Test", "9999999999", "assistant", 8000, 5)
	s.Deactivate()
	s.Activate()

	if s.Status != StaffStatusActive {
		t.Error("expected status active after activate")
	}
}
