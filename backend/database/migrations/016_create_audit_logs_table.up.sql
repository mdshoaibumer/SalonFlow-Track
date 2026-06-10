-- Migration: 016_create_audit_logs_table
-- Description: Create audit_logs table for tracking important actions
-- Created: 2026-06-08

CREATE TABLE IF NOT EXISTS audit_logs (
    id            TEXT PRIMARY KEY,
    entity_type   TEXT NOT NULL,
    entity_id     TEXT NOT NULL,
    action        TEXT NOT NULL CHECK(action IN ('create', 'update', 'delete', 'approve', 'cancel', 'pay')),
    performed_by  TEXT NOT NULL DEFAULT '',
    old_value     TEXT DEFAULT '',
    new_value     TEXT DEFAULT '',
    ip_address    TEXT DEFAULT '',
    user_agent    TEXT DEFAULT '',
    created_at    DATETIME NOT NULL DEFAULT (datetime('now'))
);

-- Audit logs are append-only, optimized for reads by entity and time
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_performed_by ON audit_logs(performed_by);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
