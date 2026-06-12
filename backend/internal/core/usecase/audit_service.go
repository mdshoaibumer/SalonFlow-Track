package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
)

// AuditService provides audit logging capabilities to all modules.
type AuditService struct {
	repo       ports.AuthRepository
	log        *slog.Logger
	appVersion string
}

// NewAuditService creates a new AuditService.
func NewAuditService(repo ports.AuthRepository, log *slog.Logger, appVersion string) *AuditService {
	return &AuditService{repo: repo, log: log, appVersion: appVersion}
}

// LogEntry records an audit event.
func (s *AuditService) LogEntry(ctx context.Context, userID, username, action, module, entityType, entityID, description, oldValue, newValue, severity string) {
	entry := &domain.AuditLog{
		ID:          uuid.New().String(),
		Timestamp:   time.Now(),
		UserID:      userID,
		Username:    username,
		Action:      action,
		Module:      module,
		EntityType:  entityType,
		EntityID:    entityID,
		Description: description,
		OldValue:    oldValue,
		NewValue:    newValue,
		AppVersion:  s.appVersion,
		Severity:    severity,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.CreateAuditLog(ctx, entry); err != nil {
		s.log.Error("failed to write audit log",
			"error", err,
			"action", action,
			"module", module,
			"entity_id", entityID,
		)
	}
}

// LogAction is a convenience method for simple audit entries.
func (s *AuditService) LogAction(ctx context.Context, userID, username, action, module, description string) {
	s.LogEntry(ctx, userID, username, action, module, "", "", description, "", "", "info")
}

// LogCritical records a critical severity audit event.
func (s *AuditService) LogCritical(ctx context.Context, userID, username, action, module, entityID, description string) {
	s.LogEntry(ctx, userID, username, action, module, "", entityID, description, "", "", "critical")
}
