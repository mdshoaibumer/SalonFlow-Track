package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/pkg/uid"
)

// Import target entities.
const (
	ImportEntityStaff     = "staff"
	ImportEntityCustomers = "customers"
	ImportEntityServices  = "services"
	ImportEntityProducts  = "products"
	ImportEntityExpenses  = "expenses"
	ImportEntityAdvances  = "advances"
	ImportEntitySalary    = "salary"
)

// Import job statuses.
const (
	ImportStatusPending    = "pending"
	ImportStatusValidating = "validating"
	ImportStatusValidated  = "validated"
	ImportStatusImporting  = "importing"
	ImportStatusCompleted  = "completed"
	ImportStatusFailed     = "failed"
)

// Import log statuses.
const (
	ImportLogSuccess = "success"
	ImportLogError   = "error"
	ImportLogWarning = "warning"
	ImportLogSkipped = "skipped"
)

// ImportJob represents an import job.
type ImportJob struct {
	ID            uuid.UUID `json:"id"`
	TemplateID    string    `json:"template_id,omitempty"`
	FileName      string    `json:"file_name"`
	FilePath      string    `json:"file_path"`
	TargetEntity  string    `json:"target_entity"`
	Status        string    `json:"status"`
	TotalRows     int       `json:"total_rows"`
	ValidRows     int       `json:"valid_rows"`
	InvalidRows   int       `json:"invalid_rows"`
	ImportedRows  int       `json:"imported_rows"`
	ColumnMapping string    `json:"column_mapping"`
	ErrorMessage  string    `json:"error_message,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// NewImportJob creates a new import job.
func NewImportJob(fileName, filePath, targetEntity string) *ImportJob {
	now := time.Now().UTC()
	return &ImportJob{
		ID:            uid.New(),
		FileName:      fileName,
		FilePath:      filePath,
		TargetEntity:  targetEntity,
		Status:        ImportStatusPending,
		ColumnMapping: "{}",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// ImportLog represents a single row import log entry.
type ImportLog struct {
	ID        uuid.UUID `json:"id"`
	JobID     uuid.UUID `json:"job_id"`
	RowNumber int       `json:"row_number"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	RowData   string    `json:"row_data"`
	CreatedAt time.Time `json:"created_at"`
}

// NewImportLog creates a new import log.
func NewImportLog(jobID uuid.UUID, rowNumber int, status, message, rowData string) *ImportLog {
	return &ImportLog{
		ID:        uid.New(),
		JobID:     jobID,
		RowNumber: rowNumber,
		Status:    status,
		Message:   message,
		RowData:   rowData,
		CreatedAt: time.Now().UTC(),
	}
}

// ImportTemplate represents a saved column mapping template.
type ImportTemplate struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	TargetEntity  string    `json:"target_entity"`
	ColumnMapping string    `json:"column_mapping"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// NewImportTemplate creates a new import template.
func NewImportTemplate(name, targetEntity, columnMapping string) *ImportTemplate {
	now := time.Now().UTC()
	return &ImportTemplate{
		ID:            uid.New(),
		Name:          name,
		TargetEntity:  targetEntity,
		ColumnMapping: columnMapping,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// ImportPreview represents the validation preview before importing.
type ImportPreview struct {
	JobID       uuid.UUID      `json:"job_id"`
	TotalRows   int            `json:"total_rows"`
	ValidRows   int            `json:"valid_rows"`
	InvalidRows int            `json:"invalid_rows"`
	Warnings    int            `json:"warnings"`
	Headers     []string       `json:"headers"`
	SampleRows  [][]string     `json:"sample_rows"`
	Errors      []ImportLogRow `json:"errors,omitempty"`
}

// ImportLogRow is a summary of a log row for preview.
type ImportLogRow struct {
	RowNumber int    `json:"row_number"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

// ColumnMapping defines how source columns map to target fields.
type ColumnMapping struct {
	SourceColumn string `json:"source_column"`
	TargetField  string `json:"target_field"`
}
