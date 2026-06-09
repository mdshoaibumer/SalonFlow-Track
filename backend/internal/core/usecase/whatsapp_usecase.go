package usecase

import (
	"context"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// WhatsAppUseCase handles WhatsApp business logic.
type WhatsAppUseCase struct {
	repo ports.WhatsAppRepository
	log  *slog.Logger
}

// NewWhatsAppUseCase creates a new WhatsAppUseCase.
func NewWhatsAppUseCase(repo ports.WhatsAppRepository, log *slog.Logger) *WhatsAppUseCase {
	return &WhatsAppUseCase{repo: repo, log: log}
}

// CreateTemplate creates a new template.
func (uc *WhatsAppUseCase) CreateTemplate(ctx context.Context, tmpl *domain.WhatsAppTemplate) error {
	if tmpl.Name == "" {
		return apperror.Validation("name", "Name is required")
	}
	if tmpl.Body == "" {
		return apperror.Validation("body", "Body is required")
	}
	return uc.repo.CreateTemplate(ctx, tmpl)
}

// UpdateTemplate updates a template.
func (uc *WhatsAppUseCase) UpdateTemplate(ctx context.Context, tmpl *domain.WhatsAppTemplate) error {
	return uc.repo.UpdateTemplate(ctx, tmpl)
}

// DeleteTemplate deletes a template.
func (uc *WhatsAppUseCase) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	return uc.repo.DeleteTemplate(ctx, id)
}

// ListTemplates lists templates.
func (uc *WhatsAppUseCase) ListTemplates(ctx context.Context, category string) ([]domain.WhatsAppTemplate, error) {
	return uc.repo.ListTemplates(ctx, category)
}

// SendMessage sends a WhatsApp message.
func (uc *WhatsAppUseCase) SendMessage(ctx context.Context, templateID, phone, name string, variables map[string]string) (*domain.WhatsAppMessage, error) {
	if phone == "" {
		return nil, apperror.Validation("phone", "Phone number is required")
	}

	// Get template body
	body := ""
	if templateID != "" {
		tmplID, err := uuid.Parse(templateID)
		if err == nil {
			tmpl, err := uc.repo.GetTemplate(ctx, tmplID)
			if err == nil {
				body = tmpl.Body
				// Replace variables
				for k, v := range variables {
					body = strings.ReplaceAll(body, "{{"+k+"}}", v)
				}
			}
		}
	}

	msg := domain.NewWhatsAppMessage(templateID, phone, name, body)
	if err := uc.repo.CreateMessage(ctx, msg); err != nil {
		return nil, err
	}

	// In production, this would call the WhatsApp API provider
	// For now, mark as sent
	_ = uc.repo.UpdateMessageStatus(ctx, msg.ID, domain.WAStatusSent, "")
	msg.Status = domain.WAStatusSent

	return msg, nil
}

// ListMessages lists messages.
func (uc *WhatsAppUseCase) ListMessages(ctx context.Context, limit, offset int, status string) ([]domain.WhatsAppMessage, int, error) {
	return uc.repo.ListMessages(ctx, limit, offset, status)
}

// GetStats gets message stats.
func (uc *WhatsAppUseCase) GetStats(ctx context.Context) (*domain.WAMessageStats, error) {
	return uc.repo.GetStats(ctx)
}

// CreateRule creates an automation rule.
func (uc *WhatsAppUseCase) CreateRule(ctx context.Context, rule *domain.AutomationRule) error {
	if rule.Name == "" {
		return apperror.Validation("name", "Name is required")
	}
	if rule.TriggerType == "" {
		return apperror.Validation("trigger_type", "Trigger type is required")
	}
	return uc.repo.CreateRule(ctx, rule)
}

// UpdateRule updates an automation rule.
func (uc *WhatsAppUseCase) UpdateRule(ctx context.Context, rule *domain.AutomationRule) error {
	return uc.repo.UpdateRule(ctx, rule)
}

// DeleteRule deletes an automation rule.
func (uc *WhatsAppUseCase) DeleteRule(ctx context.Context, id uuid.UUID) error {
	return uc.repo.DeleteRule(ctx, id)
}

// ListRules lists all rules.
func (uc *WhatsAppUseCase) ListRules(ctx context.Context) ([]domain.AutomationRule, error) {
	return uc.repo.ListRules(ctx)
}
