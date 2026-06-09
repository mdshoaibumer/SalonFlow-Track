package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/pkg/uid"
)

// Staff status constants.
const (
	StaffStatusActive   = "active"
	StaffStatusInactive = "inactive"
)

// Staff designation constants.
const (
	DesignationStylist      = "stylist"
	DesignationAssistant    = "assistant"
	DesignationReceptionist = "receptionist"
	DesignationManager      = "manager"
)

// Staff gender constants.
const (
	GenderMale   = "male"
	GenderFemale = "female"
	GenderOther  = "other"
)

var validDesignations = map[string]bool{
	DesignationStylist:      true,
	DesignationAssistant:    true,
	DesignationReceptionist: true,
	DesignationManager:      true,
}

var validGenders = map[string]bool{
	GenderMale:   true,
	GenderFemale: true,
	GenderOther:  true,
}

var phoneRegex = regexp.MustCompile(`^[6-9]\d{9}$`)

// Staff represents a salon staff member.
type Staff struct {
	ID                   uuid.UUID `json:"id"`
	StaffCode            string    `json:"staff_code"`
	FullName             string    `json:"full_name"`
	Phone                string    `json:"phone"`
	Email                string    `json:"email"`
	Gender               string    `json:"gender"`
	Designation          string    `json:"designation"`
	JoiningDate          time.Time `json:"joining_date"`
	BaseSalary           float64   `json:"base_salary"`
	CommissionPercentage float64   `json:"commission_percentage"`
	Status               string    `json:"status"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// NewStaff creates a new Staff with a UUIDv7 and auto-generated staff code.
func NewStaff(fullName, phone, designation string, baseSalary, commission float64) *Staff {
	now := time.Now().UTC()
	return &Staff{
		ID:                   uid.New(),
		StaffCode:            generateStaffCode(),
		FullName:             strings.TrimSpace(fullName),
		Phone:                strings.TrimSpace(phone),
		Gender:               GenderMale,
		Designation:          designation,
		JoiningDate:          now,
		BaseSalary:           baseSalary,
		CommissionPercentage: commission,
		Status:               StaffStatusActive,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
}

// generateStaffCode creates a unique staff code like "STF-XXXXXXXX".
func generateStaffCode() string {
	id := uid.New()
	// Use last 8 chars of UUID (random portion of UUIDv7) for uniqueness
	s := strings.ReplaceAll(id.String(), "-", "")
	short := strings.ToUpper(s[len(s)-8:])
	return fmt.Sprintf("STF-%s", short)
}

// Validate checks staff business rules.
func (s *Staff) Validate() error {
	if strings.TrimSpace(s.FullName) == "" {
		return ErrStaffNameRequired
	}
	if strings.TrimSpace(s.Phone) == "" {
		return ErrStaffPhoneRequired
	}
	if !phoneRegex.MatchString(s.Phone) {
		return ErrStaffPhoneInvalid
	}
	if !validDesignations[s.Designation] {
		return ErrStaffInvalidDesignation
	}
	if !validGenders[s.Gender] {
		return ErrStaffInvalidGender
	}
	if s.BaseSalary < 0 {
		return ErrStaffInvalidSalary
	}
	if s.CommissionPercentage < 0 || s.CommissionPercentage > 100 {
		return ErrStaffInvalidCommission
	}
	return nil
}

// Deactivate sets the staff member status to inactive.
func (s *Staff) Deactivate() {
	s.Status = StaffStatusInactive
	s.UpdatedAt = time.Now().UTC()
}

// Activate sets the staff member status to active.
func (s *Staff) Activate() {
	s.Status = StaffStatusActive
	s.UpdatedAt = time.Now().UTC()
}

// IsActive returns whether the staff member is active.
func (s *Staff) IsActive() bool {
	return s.Status == StaffStatusActive
}
