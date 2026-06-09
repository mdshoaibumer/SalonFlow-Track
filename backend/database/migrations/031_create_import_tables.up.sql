-- Excel Import & Data Migration tables

CREATE TABLE IF NOT EXISTS import_templates (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    target_entity TEXT NOT NULL CHECK (target_entity IN ('staff','customers','services','products','expenses','advances','salary')),
    column_mapping TEXT NOT NULL DEFAULT '{}',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS import_jobs (
    id TEXT PRIMARY KEY,
    template_id TEXT DEFAULT '',
    file_name TEXT NOT NULL,
    file_path TEXT NOT NULL,
    target_entity TEXT NOT NULL CHECK (target_entity IN ('staff','customers','services','products','expenses','advances','salary')),
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','validating','validated','importing','completed','failed')),
    total_rows INTEGER NOT NULL DEFAULT 0,
    valid_rows INTEGER NOT NULL DEFAULT 0,
    invalid_rows INTEGER NOT NULL DEFAULT 0,
    imported_rows INTEGER NOT NULL DEFAULT 0,
    column_mapping TEXT NOT NULL DEFAULT '{}',
    error_message TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS import_logs (
    id TEXT PRIMARY KEY,
    job_id TEXT NOT NULL REFERENCES import_jobs(id),
    row_number INTEGER NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('success','error','warning','skipped')),
    message TEXT NOT NULL DEFAULT '',
    row_data TEXT NOT NULL DEFAULT '{}',
    created_at TEXT NOT NULL
);

CREATE INDEX idx_import_jobs_status ON import_jobs(status);
CREATE INDEX idx_import_jobs_entity ON import_jobs(target_entity);
CREATE INDEX idx_import_jobs_created ON import_jobs(created_at);
CREATE INDEX idx_import_logs_job ON import_logs(job_id);
CREATE INDEX idx_import_logs_status ON import_logs(status);
CREATE INDEX idx_import_templates_entity ON import_templates(target_entity);
