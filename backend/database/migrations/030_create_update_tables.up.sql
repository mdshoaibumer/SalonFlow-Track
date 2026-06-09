-- App versions and update history tables

CREATE TABLE IF NOT EXISTS app_versions (
    id TEXT PRIMARY KEY,
    version TEXT NOT NULL UNIQUE,
    release_date TEXT NOT NULL,
    release_notes TEXT NOT NULL DEFAULT '',
    installed_at TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'available' CHECK (status IN ('available','downloading','downloaded','installing','installed','failed','rolled_back')),
    created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS update_history (
    id TEXT PRIMARY KEY,
    from_version TEXT NOT NULL,
    to_version TEXT NOT NULL,
    update_date TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','downloading','downloaded','installing','completed','failed','rolled_back')),
    error_message TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL
);

CREATE INDEX idx_app_versions_version ON app_versions(version);
CREATE INDEX idx_app_versions_status ON app_versions(status);
CREATE INDEX idx_update_history_status ON update_history(status);
CREATE INDEX idx_update_history_date ON update_history(update_date);
