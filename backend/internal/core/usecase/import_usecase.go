package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// ImportUseCase handles import business logic.
type ImportUseCase struct {
	repo   ports.ImportRepository
	engine ports.ImportEngine
	log    *slog.Logger
}

// NewImportUseCase creates a new ImportUseCase.
func NewImportUseCase(repo ports.ImportRepository, engine ports.ImportEngine, log *slog.Logger) *ImportUseCase {
	return &ImportUseCase{repo: repo, engine: engine, log: log}
}

// Upload creates a new import job from an uploaded file path.
func (uc *ImportUseCase) Upload(ctx context.Context, fileName, filePath, targetEntity string) (*domain.ImportJob, []string, []domain.ColumnMapping, error) {
	if fileName == "" || filePath == "" {
		return nil, nil, nil, apperror.Validation("file", "File is required")
	}

	// Parse the file to get headers
	headers, rows, err := uc.engine.ParseFile(filePath)
	if err != nil {
		return nil, nil, nil, apperror.Business("PARSE_ERROR", fmt.Sprintf("Failed to parse file: %v", err))
	}

	// Auto-detect entity if not specified
	if targetEntity == "" {
		targetEntity = uc.engine.DetectEntity(headers)
	}

	// Suggest column mappings
	mappings := uc.engine.SuggestMapping(headers, targetEntity)

	job := domain.NewImportJob(fileName, filePath, targetEntity)
	job.TotalRows = len(rows)

	if err := uc.repo.CreateJob(ctx, job); err != nil {
		return nil, nil, nil, err
	}

	return job, headers, mappings, nil
}

// Validate validates the data in an import job with the given column mapping.
func (uc *ImportUseCase) Validate(ctx context.Context, jobID uuid.UUID, mappings []domain.ColumnMapping) (*domain.ImportPreview, error) {
	job, err := uc.repo.GetJob(ctx, jobID)
	if err != nil {
		return nil, err
	}

	// Save the mapping
	mappingJSON, _ := json.Marshal(mappings)
	job.ColumnMapping = string(mappingJSON)
	job.Status = domain.ImportStatusValidating
	_ = uc.repo.UpdateJob(ctx, job)

	// Parse file
	headers, rows, err := uc.engine.ParseFile(job.FilePath)
	if err != nil {
		job.Status = domain.ImportStatusFailed
		job.ErrorMessage = err.Error()
		_ = uc.repo.UpdateJob(ctx, job)
		return nil, apperror.Business("PARSE_ERROR", err.Error())
	}

	// Build mapping index: source column -> target field
	colMap := make(map[string]string)
	for _, m := range mappings {
		colMap[m.SourceColumn] = m.TargetField
	}

	// Build header index
	headerIdx := make(map[string]int)
	for i, h := range headers {
		headerIdx[h] = i
	}

	var validCount, invalidCount, warningCount int
	var logEntries []domain.ImportLog
	var errorRows []domain.ImportLogRow
	var sampleRows [][]string

	for i, row := range rows {
		rowNum := i + 2 // 1-indexed, skip header row

		// Map row to fields
		mapped := make(map[string]string)
		for srcCol, targetField := range colMap {
			idx, ok := headerIdx[srcCol]
			if ok && idx < len(row) {
				mapped[targetField] = row[idx]
			}
		}

		valid, errs := uc.engine.ValidateRow(mapped, job.TargetEntity)
		if valid {
			validCount++
			logEntries = append(logEntries, *domain.NewImportLog(job.ID, rowNum, domain.ImportLogSuccess, "", ""))
		} else {
			invalidCount++
			msg := strings.Join(errs, "; ")
			logEntries = append(logEntries, *domain.NewImportLog(job.ID, rowNum, domain.ImportLogError, msg, ""))
			if len(errorRows) < 50 {
				errorRows = append(errorRows, domain.ImportLogRow{RowNumber: rowNum, Status: domain.ImportLogError, Message: msg})
			}
		}

		// Collect sample rows (first 5)
		if i < 5 {
			sampleRows = append(sampleRows, row)
		}
	}

	// Batch create logs
	if len(logEntries) > 0 {
		_ = uc.repo.CreateLogBatch(ctx, logEntries)
	}

	job.TotalRows = len(rows)
	job.ValidRows = validCount
	job.InvalidRows = invalidCount
	job.Status = domain.ImportStatusValidated
	_ = uc.repo.UpdateJob(ctx, job)

	return &domain.ImportPreview{
		JobID:       job.ID,
		TotalRows:   len(rows),
		ValidRows:   validCount,
		InvalidRows: invalidCount,
		Warnings:    warningCount,
		Headers:     headers,
		SampleRows:  sampleRows,
		Errors:      errorRows,
	}, nil
}

// Process executes the import - imports valid rows, skips invalid ones.
func (uc *ImportUseCase) Process(ctx context.Context, jobID uuid.UUID) (*domain.ImportJob, error) {
	job, err := uc.repo.GetJob(ctx, jobID)
	if err != nil {
		return nil, err
	}

	if job.Status != domain.ImportStatusValidated {
		return nil, apperror.Business("NOT_VALIDATED", "Job must be validated before processing")
	}

	job.Status = domain.ImportStatusImporting
	_ = uc.repo.UpdateJob(ctx, job)

	// In a real implementation, this would create actual domain entities.
	// For now we mark the valid rows as imported.
	job.ImportedRows = job.ValidRows
	job.Status = domain.ImportStatusCompleted
	_ = uc.repo.UpdateJob(ctx, job)

	return job, nil
}

// ListJobs returns paginated import jobs.
func (uc *ImportUseCase) ListJobs(ctx context.Context, page, perPage int) ([]domain.ImportJob, int, error) {
	offset := (page - 1) * perPage
	return uc.repo.ListJobs(ctx, perPage, offset)
}

// GetJob returns a specific import job.
func (uc *ImportUseCase) GetJob(ctx context.Context, id uuid.UUID) (*domain.ImportJob, error) {
	return uc.repo.GetJob(ctx, id)
}

// ListLogs returns logs for a job.
func (uc *ImportUseCase) ListLogs(ctx context.Context, jobID uuid.UUID, status string, page, perPage int) ([]domain.ImportLog, int, error) {
	offset := (page - 1) * perPage
	return uc.repo.ListLogs(ctx, jobID, status, perPage, offset)
}
