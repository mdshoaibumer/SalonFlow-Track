-- Cloud Backup Configuration
CREATE TABLE IF NOT EXISTS cloud_backup_config (
    id TEXT PRIMARY KEY,
    provider TEXT NOT NULL DEFAULT 'none' CHECK (provider IN ('none','google_drive','aws_s3','digitalocean_spaces')),
    bucket_name TEXT NOT NULL DEFAULT '',
    region TEXT NOT NULL DEFAULT '',
    access_key TEXT NOT NULL DEFAULT '',
    endpoint TEXT NOT NULL DEFAULT '',
    encrypt_backups INTEGER NOT NULL DEFAULT 1,
    auto_backup INTEGER NOT NULL DEFAULT 0,
    auto_backup_interval_hours INTEGER NOT NULL DEFAULT 24,
    max_versions INTEGER NOT NULL DEFAULT 10,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

-- Cloud Backup History
CREATE TABLE IF NOT EXISTS cloud_backup_history (
    id TEXT PRIMARY KEY,
    provider TEXT NOT NULL,
    file_name TEXT NOT NULL,
    file_size INTEGER NOT NULL DEFAULT 0,
    remote_path TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','uploading','completed','failed','restoring','restored')),
    is_encrypted INTEGER NOT NULL DEFAULT 0,
    error_message TEXT NOT NULL DEFAULT '',
    started_at TEXT NOT NULL DEFAULT '',
    completed_at TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL
);

CREATE INDEX idx_cloud_backup_history_status ON cloud_backup_history(status);
CREATE INDEX idx_cloud_backup_history_created ON cloud_backup_history(created_at);
