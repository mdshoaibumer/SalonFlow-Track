package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// ImportRepository manages import jobs and logs in the database.
type ImportRepository interface {
	CreateJob(ctx context.Context, job *domain.ImportJob) error
	UpdateJob(ctx context.Context, job *domain.ImportJob) error
	GetJob(ctx context.Context, id uuid.UUID) (*domain.ImportJob, error)
	ListJobs(ctx context.Context, limit, offset int) ([]domain.ImportJob, int, error)

	CreateLog(ctx context.Context, log *domain.ImportLog) error
	CreateLogBatch(ctx context.Context, logs []domain.ImportLog) error
	ListLogs(ctx context.Context, jobID uuid.UUID, status string, limit, offset int) ([]domain.ImportLog, int, error)

	CreateTemplate(ctx context.Context, tmpl *domain.ImportTemplate) error
	ListTemplates(ctx context.Context, entity string) ([]domain.ImportTemplate, error)
}

// ImportEngine handles file parsing, column detection, and data extraction.
type ImportEngine interface {
	ParseFile(filePath string) (headers []string, rows [][]string, err error)
	DetectEntity(headers []string) string
	SuggestMapping(headers []string, targetEntity string) []domain.ColumnMapping
	ValidateRow(row map[string]string, targetEntity string) (bool, []string)
	UploadDir() string
}
