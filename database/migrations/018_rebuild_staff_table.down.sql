-- Migration: 018_rebuild_staff_table (rollback)
-- Restore original Phase 1 staff table schema

DROP INDEX IF EXISTS idx_staff_phone;
DROP INDEX IF EXISTS idx_staff_status;
DROP INDEX IF EXISTS idx_staff_designation;
DROP INDEX IF EXISTS idx_staff_staff_code;
DROP INDEX IF EXISTS idx_staff_full_name;
DROP TABLE IF EXISTS staff;

CREATE TABLE IF NOT EXISTS staff (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    phone       TEXT,
    role        TEXT NOT NULL DEFAULT 'stylist',
    is_active   INTEGER NOT NULL DEFAULT 1,
    joined_at   DATETIME NOT NULL DEFAULT (datetime('now')),
    specialties TEXT DEFAULT '[]',
    created_at  DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at  DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_staff_is_active ON staff(is_active);
CREATE INDEX idx_staff_role ON staff(role);
