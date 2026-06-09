package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/pkg/uid"
)

// Customer status constants.
const (
	CustomerStatusActive   = "active"
	CustomerStatusInactive = "inactive"
)

// Customer represents a salon customer.
type Customer struct {
	ID              uuid.UUID  `json:"id"`
	CustomerCode    string     `json:"customer_code"`
	FullName        string     `json:"full_name"`
	Phone           string     `json:"phone"`
	Email           string     `json:"email"`
	Gender          string     `json:"gender"`
	DateOfBirth     *time.Time `json:"date_of_birth,omitempty"`
	AnniversaryDate *time.Time `json:"anniversary_date,omitempty"`
	Address         string     `json:"address"`
	Notes           string     `json:"notes"`
	TotalVisits     int        `json:"total_visits"`
	TotalSpent      float64    `json:"total_spent"`
	LastVisitDate   *time.Time `json:"last_visit_date,omitempty"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// NewCustomer creates a new Customer with UUIDv7 and auto-generated code.
func NewCustomer(fullName, phone string) *Customer {
	now := time.Now().UTC()
	return &Customer{
		ID:           uid.New(),
		CustomerCode: generateCustomerCode(),
		FullName:     strings.TrimSpace(fullName),
		Phone:        strings.TrimSpace(phone),
		Gender:       GenderOther,
		Status:       CustomerStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// generateCustomerCode creates a unique customer code like "CUS-XXXXXXXX".
func generateCustomerCode() string {
	id := uid.New()
	s := strings.ReplaceAll(id.String(), "-", "")
	short := strings.ToUpper(s[len(s)-8:])
	return fmt.Sprintf("CUS-%s", short)
}

// Validate checks customer business rules.
func (c *Customer) Validate() error {
	if strings.TrimSpace(c.FullName) == "" {
		return ErrCustomerNameRequired
	}
	if strings.TrimSpace(c.Phone) == "" {
		return ErrCustomerPhoneRequired
	}
	if !phoneRegex.MatchString(c.Phone) {
		return ErrCustomerPhoneInvalid
	}
	if !validGenders[c.Gender] {
		return ErrCustomerInvalidGender
	}
	return nil
}

// RecordVisit increments visit count and updates last visit.
func (c *Customer) RecordVisit(amount float64, visitDate time.Time) {
	c.TotalVisits++
	c.TotalSpent += amount
	c.LastVisitDate = &visitDate
	c.UpdatedAt = time.Now().UTC()
}

// Deactivate sets the customer status to inactive.
func (c *Customer) Deactivate() {
	c.Status = CustomerStatusInactive
	c.UpdatedAt = time.Now().UTC()
}

// Activate sets the customer status to active.
func (c *Customer) Activate() {
	c.Status = CustomerStatusActive
	c.UpdatedAt = time.Now().UTC()
}

// IsActive returns true if the customer is active.
func (c *Customer) IsActive() bool {
	return c.Status == CustomerStatusActive
}
