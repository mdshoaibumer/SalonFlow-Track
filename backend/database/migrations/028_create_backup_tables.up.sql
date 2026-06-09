-- Backup & Restore history tables

CREATE TABLE IF NOT EXISTS backup_history (
    id TEXT PRIMARY KEY,
    backup_name TEXT NOT NULL,
    backup_type TEXT NOT NULL CHECK (backup_type IN ('manual','daily','before_update','before_restore')),
    backup_path TEXT NOT NULL,
    file_size INTEGER NOT NULL DEFAULT 0,
    checksum TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','completed','failed','corrupted','verified')),
    error_message TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS restore_history (
    id TEXT PRIMARY KEY,
    backup_id TEXT NOT NULL REFERENCES backup_history(id),
    backup_name TEXT NOT NULL,
    restore_date TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','completed','failed')),
    notes TEXT NOT NULL DEFAULT '',
    error_message TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL
);

CREATE INDEX idx_backup_history_status ON backup_history(status);
CREATE INDEX idx_backup_history_type ON backup_history(backup_type);
CREATE INDEX idx_backup_history_created ON backup_history(created_at);
CREATE INDEX idx_restore_history_created ON restore_history(created_at);
