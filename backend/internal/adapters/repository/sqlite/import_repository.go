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

// ImportRepository implements ports.ImportRepository.
type ImportRepository struct {
	db  *sql.DB
	log *slog.Logger
}

// NewImportRepository creates a new ImportRepository.
func NewImportRepository(db *sql.DB, log *slog.Logger) *ImportRepository {
	return &ImportRepository{db: db, log: log}
}

func (r *ImportRepository) CreateJob(ctx context.Context, job *domain.ImportJob) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO import_jobs (id, template_id, file_name, file_path, target_entity, status, total_rows, valid_rows, invalid_rows, imported_rows, column_mapping, error_message, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		job.ID.String(), job.TemplateID, job.FileName, job.FilePath, job.TargetEntity,
		job.Status, job.TotalRows, job.ValidRows, job.InvalidRows, job.ImportedRows,
		job.ColumnMapping, job.ErrorMessage,
		job.CreatedAt.Format(time.RFC3339), job.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("create import job", err)
	}
	return nil
}

func (r *ImportRepository) UpdateJob(ctx context.Context, job *domain.ImportJob) error {
	job.UpdatedAt = time.Now().UTC()
	_, err := r.db.ExecContext(ctx, `UPDATE import_jobs SET status = ?, total_rows = ?, valid_rows = ?, invalid_rows = ?, imported_rows = ?, column_mapping = ?, error_message = ?, updated_at = ? WHERE id = ?`,
		job.Status, job.TotalRows, job.ValidRows, job.InvalidRows, job.ImportedRows,
		job.ColumnMapping, job.ErrorMessage, job.UpdatedAt.Format(time.RFC3339), job.ID.String())
	if err != nil {
		return apperror.Database("update import job", err)
	}
	return nil
}

func (r *ImportRepository) GetJob(ctx context.Context, id uuid.UUID) (*domain.ImportJob, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, template_id, file_name, file_path, target_entity, status, total_rows, valid_rows, invalid_rows, imported_rows, column_mapping, error_message, created_at, updated_at
		FROM import_jobs WHERE id = ?`, id.String())

	var job domain.ImportJob
	var idStr, createdAt, updatedAt string
	err := row.Scan(&idStr, &job.TemplateID, &job.FileName, &job.FilePath, &job.TargetEntity,
		&job.Status, &job.TotalRows, &job.ValidRows, &job.InvalidRows, &job.ImportedRows,
		&job.ColumnMapping, &job.ErrorMessage, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("import_job", id.String())
	}
	if err != nil {
		return nil, apperror.Database("get import job", err)
	}
	job.ID, _ = uuid.Parse(idStr)
	job.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	job.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return &job, nil
}

func (r *ImportRepository) ListJobs(ctx context.Context, limit, offset int) ([]domain.ImportJob, int, error) {
	var total int
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM import_jobs`).Scan(&total)

	rows, err := r.db.QueryContext(ctx, `SELECT id, template_id, file_name, file_path, target_entity, status, total_rows, valid_rows, invalid_rows, imported_rows, column_mapping, error_message, created_at, updated_at
		FROM import_jobs ORDER BY created_at DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, 0, apperror.Database("list import jobs", err)
	}
	defer rows.Close()

	var jobs []domain.ImportJob
	for rows.Next() {
		var job domain.ImportJob
		var idStr, createdAt, updatedAt string
		err := rows.Scan(&idStr, &job.TemplateID, &job.FileName, &job.FilePath, &job.TargetEntity,
			&job.Status, &job.TotalRows, &job.ValidRows, &job.InvalidRows, &job.ImportedRows,
			&job.ColumnMapping, &job.ErrorMessage, &createdAt, &updatedAt)
		if err != nil {
			return nil, 0, apperror.Database("scan import job", err)
		}
		job.ID, _ = uuid.Parse(idStr)
		job.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		job.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		jobs = append(jobs, job)
	}
	return jobs, total, nil
}

func (r *ImportRepository) CreateLog(ctx context.Context, log *domain.ImportLog) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO import_logs (id, job_id, row_number, status, message, row_data, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		log.ID.String(), log.JobID.String(), log.RowNumber, log.Status,
		log.Message, log.RowData, log.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("create import log", err)
	}
	return nil
}

func (r *ImportRepository) CreateLogBatch(ctx context.Context, logs []domain.ImportLog) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return apperror.Database("begin tx", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO import_logs (id, job_id, row_number, status, message, row_data, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return apperror.Database("prepare log insert", err)
	}
	defer stmt.Close()

	for _, l := range logs {
		_, err := stmt.ExecContext(ctx, l.ID.String(), l.JobID.String(), l.RowNumber, l.Status, l.Message, l.RowData, l.CreatedAt.Format(time.RFC3339))
		if err != nil {
			return apperror.Database("insert log batch", err)
		}
	}

	return tx.Commit()
}

func (r *ImportRepository) ListLogs(ctx context.Context, jobID uuid.UUID, status string, limit, offset int) ([]domain.ImportLog, int, error) {
	var total int
	query := `SELECT COUNT(*) FROM import_logs WHERE job_id = ?`
	args := []interface{}{jobID.String()}
	if status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}
	r.db.QueryRowContext(ctx, query, args...).Scan(&total)

	selectQuery := `SELECT id, job_id, row_number, status, message, row_data, created_at FROM import_logs WHERE job_id = ?`
	selectArgs := []interface{}{jobID.String()}
	if status != "" {
		selectQuery += ` AND status = ?`
		selectArgs = append(selectArgs, status)
	}
	selectQuery += ` ORDER BY row_number ASC LIMIT ? OFFSET ?`
	selectArgs = append(selectArgs, limit, offset)

	rows, err := r.db.QueryContext(ctx, selectQuery, selectArgs...)
	if err != nil {
		return nil, 0, apperror.Database("list import logs", err)
	}
	defer rows.Close()

	var logs []domain.ImportLog
	for rows.Next() {
		var l domain.ImportLog
		var idStr, jobIDStr, createdAt string
		err := rows.Scan(&idStr, &jobIDStr, &l.RowNumber, &l.Status, &l.Message, &l.RowData, &createdAt)
		if err != nil {
			return nil, 0, apperror.Database("scan import log", err)
		}
		l.ID, _ = uuid.Parse(idStr)
		l.JobID, _ = uuid.Parse(jobIDStr)
		l.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		logs = append(logs, l)
	}
	return logs, total, nil
}

func (r *ImportRepository) CreateTemplate(ctx context.Context, tmpl *domain.ImportTemplate) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO import_templates (id, name, target_entity, column_mapping, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		tmpl.ID.String(), tmpl.Name, tmpl.TargetEntity, tmpl.ColumnMapping,
		tmpl.CreatedAt.Format(time.RFC3339), tmpl.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("create import template", err)
	}
	return nil
}

func (r *ImportRepository) ListTemplates(ctx context.Context, entity string) ([]domain.ImportTemplate, error) {
	query := `SELECT id, name, target_entity, column_mapping, created_at, updated_at FROM import_templates`
	var args []interface{}
	if entity != "" {
		query += ` WHERE target_entity = ?`
		args = append(args, entity)
	}
	query += ` ORDER BY name ASC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.Database("list import templates", err)
	}
	defer rows.Close()

	var templates []domain.ImportTemplate
	for rows.Next() {
		var t domain.ImportTemplate
		var idStr, createdAt, updatedAt string
		err := rows.Scan(&idStr, &t.Name, &t.TargetEntity, &t.ColumnMapping, &createdAt, &updatedAt)
		if err != nil {
			return nil, apperror.Database("scan import template", err)
		}
		t.ID, _ = uuid.Parse(idStr)
		t.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		t.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		templates = append(templates, t)
	}
	return templates, nil
}
