package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/pkg/uid"
)

// WhatsApp message statuses.
const (
	WAStatusQueued    = "queued"
	WAStatusSent      = "sent"
	WAStatusDelivered = "delivered"
	WAStatusRead      = "read"
	WAStatusFailed    = "failed"
)

// Template categories.
const (
	WACategoryAppointment = "appointment"
	WACategoryReminder    = "reminder"
	WACategoryBirthday    = "birthday"
	WACategoryPayment     = "payment"
	WACategoryMembership  = "membership"
	WACategoryInvoice     = "invoice"
	WACategoryGeneral     = "general"
)

// Automation trigger types.
const (
	WATriggerAppointmentConfirmed = "appointment_confirmed"
	WATriggerAppointmentReminder  = "appointment_reminder"
	WATriggerBirthday             = "birthday"
	WATriggerPaymentDue           = "payment_due"
	WATriggerMembershipExpiry     = "membership_expiry"
	WATriggerInvoiceCreated       = "invoice_created"
)

// WhatsAppTemplate is a message template.
type WhatsAppTemplate struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	Body      string    `json:"body"`
	Variables string    `json:"variables"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewWhatsAppTemplate creates a new template.
func NewWhatsAppTemplate(name, category, body string) *WhatsAppTemplate {
	now := time.Now().UTC()
	return &WhatsAppTemplate{
		ID:        uid.New(),
		Name:      name,
		Category:  category,
		Body:      body,
		Variables: "[]",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// WhatsAppMessage is a sent/queued message.
type WhatsAppMessage struct {
	ID                uuid.UUID `json:"id"`
	TemplateID        string    `json:"template_id"`
	RecipientPhone    string    `json:"recipient_phone"`
	RecipientName     string    `json:"recipient_name"`
	MessageBody       string    `json:"message_body"`
	Status            string    `json:"status"`
	Provider          string    `json:"provider"`
	ProviderMessageID string    `json:"provider_message_id"`
	ErrorMessage      string    `json:"error_message"`
	SentAt            string    `json:"sent_at"`
	DeliveredAt       string    `json:"delivered_at"`
	ReadAt            string    `json:"read_at"`
	CreatedAt         time.Time `json:"created_at"`
}

// NewWhatsAppMessage creates a new message.
func NewWhatsAppMessage(templateID, phone, name, body string) *WhatsAppMessage {
	return &WhatsAppMessage{
		ID:             uid.New(),
		TemplateID:     templateID,
		RecipientPhone: phone,
		RecipientName:  name,
		MessageBody:    body,
		Status:         WAStatusQueued,
		CreatedAt:      time.Now().UTC(),
	}
}

// AutomationRule defines when to auto-send messages.
type AutomationRule struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	TriggerType  string    `json:"trigger_type"`
	TemplateID   string    `json:"template_id"`
	DelayMinutes int       `json:"delay_minutes"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewAutomationRule creates a new automation rule.
func NewAutomationRule(name, triggerType, templateID string, delayMinutes int) *AutomationRule {
	now := time.Now().UTC()
	return &AutomationRule{
		ID:           uid.New(),
		Name:         name,
		TriggerType:  triggerType,
		TemplateID:   templateID,
		DelayMinutes: delayMinutes,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// WAMessageStats holds message delivery stats.
type WAMessageStats struct {
	TotalSent int `json:"total_sent"`
	Delivered int `json:"delivered"`
	Read      int `json:"read"`
	Failed    int `json:"failed"`
	Queued    int `json:"queued"`
}
