package sqlite

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// WhatsAppRepository is the SQLite implementation.
type WhatsAppRepository struct {
	db  *sql.DB
	log *slog.Logger
}

// NewWhatsAppRepository creates a new WhatsAppRepository.
func NewWhatsAppRepository(db *sql.DB, log *slog.Logger) *WhatsAppRepository {
	return &WhatsAppRepository{db: db, log: log}
}

// CreateTemplate inserts a template.
func (r *WhatsAppRepository) CreateTemplate(ctx context.Context, tmpl *domain.WhatsAppTemplate) error {
	isActive := 0
	if tmpl.IsActive {
		isActive = 1
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO whatsapp_templates (id, name, category, body, variables, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		tmpl.ID, tmpl.Name, tmpl.Category, tmpl.Body, tmpl.Variables, isActive,
		tmpl.CreatedAt.Format(time.RFC3339), tmpl.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("create_wa_template", err)
	}
	return nil
}

// UpdateTemplate updates a template.
func (r *WhatsAppRepository) UpdateTemplate(ctx context.Context, tmpl *domain.WhatsAppTemplate) error {
	tmpl.UpdatedAt = time.Now().UTC()
	isActive := 0
	if tmpl.IsActive {
		isActive = 1
	}
	res, err := r.db.ExecContext(ctx, `
		UPDATE whatsapp_templates SET name=?, category=?, body=?, variables=?, is_active=?, updated_at=? WHERE id=?`,
		tmpl.Name, tmpl.Category, tmpl.Body, tmpl.Variables, isActive, tmpl.UpdatedAt.Format(time.RFC3339), tmpl.ID)
	if err != nil {
		return apperror.Database("update_wa_template", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperror.NotFound("whatsapp_template", tmpl.ID.String())
	}
	return nil
}

// DeleteTemplate deletes a template.
func (r *WhatsAppRepository) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM whatsapp_templates WHERE id=?`, id)
	if err != nil {
		return apperror.Database("delete_wa_template", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperror.NotFound("whatsapp_template", id.String())
	}
	return nil
}

// GetTemplate retrieves a template by ID.
func (r *WhatsAppRepository) GetTemplate(ctx context.Context, id uuid.UUID) (*domain.WhatsAppTemplate, error) {
	var tmpl domain.WhatsAppTemplate
	var isActive int
	var createdAt, updatedAt string
	err := r.db.QueryRowContext(ctx, `SELECT id, name, category, body, variables, is_active, created_at, updated_at FROM whatsapp_templates WHERE id=?`, id).
		Scan(&tmpl.ID, &tmpl.Name, &tmpl.Category, &tmpl.Body, &tmpl.Variables, &isActive, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("whatsapp_template", id.String())
	}
	if err != nil {
		return nil, apperror.Database("get_wa_template", err)
	}
	tmpl.IsActive = isActive == 1
	tmpl.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	tmpl.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return &tmpl, nil
}

// ListTemplates lists templates, optionally by category.
func (r *WhatsAppRepository) ListTemplates(ctx context.Context, category string) ([]domain.WhatsAppTemplate, error) {
	query := `SELECT id, name, category, body, variables, is_active, created_at, updated_at FROM whatsapp_templates`
	var args []interface{}
	if category != "" {
		query += ` WHERE category=?`
		args = append(args, category)
	}
	query += ` ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.Database("list_wa_templates", err)
	}
	defer rows.Close()

	var templates []domain.WhatsAppTemplate
	for rows.Next() {
		var tmpl domain.WhatsAppTemplate
		var isActive int
		var createdAt, updatedAt string
		if err := rows.Scan(&tmpl.ID, &tmpl.Name, &tmpl.Category, &tmpl.Body, &tmpl.Variables, &isActive, &createdAt, &updatedAt); err != nil {
			return nil, apperror.Database("list_wa_templates_scan", err)
		}
		tmpl.IsActive = isActive == 1
		tmpl.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		tmpl.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		templates = append(templates, tmpl)
	}
	return templates, nil
}

// CreateMessage inserts a message.
func (r *WhatsAppRepository) CreateMessage(ctx context.Context, msg *domain.WhatsAppMessage) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO whatsapp_messages (id, template_id, recipient_phone, recipient_name, message_body, status, provider, provider_message_id, error_message, sent_at, delivered_at, read_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		msg.ID, msg.TemplateID, msg.RecipientPhone, msg.RecipientName, msg.MessageBody,
		msg.Status, msg.Provider, msg.ProviderMessageID, msg.ErrorMessage,
		msg.SentAt, msg.DeliveredAt, msg.ReadAt, msg.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("create_wa_message", err)
	}
	return nil
}

// UpdateMessageStatus updates a message's status.
func (r *WhatsAppRepository) UpdateMessageStatus(ctx context.Context, id uuid.UUID, status, providerMsgID string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	var setExtra string
	switch status {
	case domain.WAStatusSent:
		setExtra = ", sent_at='" + now + "'"
	case domain.WAStatusDelivered:
		setExtra = ", delivered_at='" + now + "'"
	case domain.WAStatusRead:
		setExtra = ", read_at='" + now + "'"
	}
	_, err := r.db.ExecContext(ctx, `UPDATE whatsapp_messages SET status=?, provider_message_id=?`+setExtra+` WHERE id=?`, status, providerMsgID, id)
	if err != nil {
		return apperror.Database("update_wa_message_status", err)
	}
	return nil
}

// ListMessages lists messages with pagination.
func (r *WhatsAppRepository) ListMessages(ctx context.Context, limit, offset int, status string) ([]domain.WhatsAppMessage, int, error) {
	where := "1=1"
	var args []interface{}
	if status != "" {
		where += " AND status=?"
		args = append(args, status)
	}

	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM whatsapp_messages WHERE "+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, apperror.Database("list_wa_messages_count", err)
	}

	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, `SELECT id, template_id, recipient_phone, recipient_name, message_body, status, provider, provider_message_id, error_message, sent_at, delivered_at, read_at, created_at FROM whatsapp_messages WHERE `+where+` ORDER BY created_at DESC LIMIT ? OFFSET ?`, args...)
	if err != nil {
		return nil, 0, apperror.Database("list_wa_messages", err)
	}
	defer rows.Close()

	var messages []domain.WhatsAppMessage
	for rows.Next() {
		var msg domain.WhatsAppMessage
		var createdAt string
		if err := rows.Scan(&msg.ID, &msg.TemplateID, &msg.RecipientPhone, &msg.RecipientName, &msg.MessageBody, &msg.Status, &msg.Provider, &msg.ProviderMessageID, &msg.ErrorMessage, &msg.SentAt, &msg.DeliveredAt, &msg.ReadAt, &createdAt); err != nil {
			return nil, 0, apperror.Database("list_wa_messages_scan", err)
		}
		msg.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		messages = append(messages, msg)
	}
	return messages, total, nil
}

// GetStats retrieves message delivery statistics.
func (r *WhatsAppRepository) GetStats(ctx context.Context) (*domain.WAMessageStats, error) {
	var stats domain.WAMessageStats
	err := r.db.QueryRowContext(ctx, `
		SELECT
			COALESCE(SUM(CASE WHEN status IN ('sent','delivered','read') THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status='delivered' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status='read' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status='failed' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status='queued' THEN 1 ELSE 0 END), 0)
		FROM whatsapp_messages`).
		Scan(&stats.TotalSent, &stats.Delivered, &stats.Read, &stats.Failed, &stats.Queued)
	if err != nil {
		return nil, apperror.Database("get_wa_stats", err)
	}
	return &stats, nil
}

// CreateRule inserts an automation rule.
func (r *WhatsAppRepository) CreateRule(ctx context.Context, rule *domain.AutomationRule) error {
	isActive := 0
	if rule.IsActive {
		isActive = 1
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO automation_rules (id, name, trigger_type, template_id, delay_minutes, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		rule.ID, rule.Name, rule.TriggerType, rule.TemplateID, rule.DelayMinutes, isActive,
		rule.CreatedAt.Format(time.RFC3339), rule.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("create_automation_rule", err)
	}
	return nil
}

// UpdateRule updates an automation rule.
func (r *WhatsAppRepository) UpdateRule(ctx context.Context, rule *domain.AutomationRule) error {
	rule.UpdatedAt = time.Now().UTC()
	isActive := 0
	if rule.IsActive {
		isActive = 1
	}
	res, err := r.db.ExecContext(ctx, `
		UPDATE automation_rules SET name=?, trigger_type=?, template_id=?, delay_minutes=?, is_active=?, updated_at=? WHERE id=?`,
		rule.Name, rule.TriggerType, rule.TemplateID, rule.DelayMinutes, isActive, rule.UpdatedAt.Format(time.RFC3339), rule.ID)
	if err != nil {
		return apperror.Database("update_automation_rule", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperror.NotFound("automation_rule", rule.ID.String())
	}
	return nil
}

// DeleteRule deletes an automation rule.
func (r *WhatsAppRepository) DeleteRule(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM automation_rules WHERE id=?`, id)
	if err != nil {
		return apperror.Database("delete_automation_rule", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperror.NotFound("automation_rule", id.String())
	}
	return nil
}

// ListRules lists all automation rules.
func (r *WhatsAppRepository) ListRules(ctx context.Context) ([]domain.AutomationRule, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, trigger_type, template_id, delay_minutes, is_active, created_at, updated_at FROM automation_rules ORDER BY name`)
	if err != nil {
		return nil, apperror.Database("list_automation_rules", err)
	}
	defer rows.Close()

	var rules []domain.AutomationRule
	for rows.Next() {
		var rule domain.AutomationRule
		var isActive int
		var createdAt, updatedAt string
		if err := rows.Scan(&rule.ID, &rule.Name, &rule.TriggerType, &rule.TemplateID, &rule.DelayMinutes, &isActive, &createdAt, &updatedAt); err != nil {
			return nil, apperror.Database("list_automation_rules_scan", err)
		}
		rule.IsActive = isActive == 1
		rule.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		rule.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		rules = append(rules, rule)
	}
	return rules, nil
}

// GetRulesByTrigger retrieves active rules for a trigger type.
func (r *WhatsAppRepository) GetRulesByTrigger(ctx context.Context, triggerType string) ([]domain.AutomationRule, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, trigger_type, template_id, delay_minutes, is_active, created_at, updated_at FROM automation_rules WHERE trigger_type=? AND is_active=1`, triggerType)
	if err != nil {
		return nil, apperror.Database("get_rules_by_trigger", err)
	}
	defer rows.Close()

	var rules []domain.AutomationRule
	for rows.Next() {
		var rule domain.AutomationRule
		var isActive int
		var createdAt, updatedAt string
		if err := rows.Scan(&rule.ID, &rule.Name, &rule.TriggerType, &rule.TemplateID, &rule.DelayMinutes, &isActive, &createdAt, &updatedAt); err != nil {
			return nil, apperror.Database("get_rules_by_trigger_scan", err)
		}
		rule.IsActive = isActive == 1
		rule.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		rule.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		rules = append(rules, rule)
	}
	return rules, nil
}
