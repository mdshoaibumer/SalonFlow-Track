package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// WhatsAppRepository manages WhatsApp data.
type WhatsAppRepository interface {
	// Templates
	CreateTemplate(ctx context.Context, tmpl *domain.WhatsAppTemplate) error
	UpdateTemplate(ctx context.Context, tmpl *domain.WhatsAppTemplate) error
	DeleteTemplate(ctx context.Context, id uuid.UUID) error
	GetTemplate(ctx context.Context, id uuid.UUID) (*domain.WhatsAppTemplate, error)
	ListTemplates(ctx context.Context, category string) ([]domain.WhatsAppTemplate, error)

	// Messages
	CreateMessage(ctx context.Context, msg *domain.WhatsAppMessage) error
	UpdateMessageStatus(ctx context.Context, id uuid.UUID, status, providerMsgID string) error
	ListMessages(ctx context.Context, limit, offset int, status string) ([]domain.WhatsAppMessage, int, error)
	GetStats(ctx context.Context) (*domain.WAMessageStats, error)

	// Automation
	CreateRule(ctx context.Context, rule *domain.AutomationRule) error
	UpdateRule(ctx context.Context, rule *domain.AutomationRule) error
	DeleteRule(ctx context.Context, id uuid.UUID) error
	ListRules(ctx context.Context) ([]domain.AutomationRule, error)
	GetRulesByTrigger(ctx context.Context, triggerType string) ([]domain.AutomationRule, error)
}
