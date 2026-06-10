-- Migration: 002_create_staff_table
-- Description: Create the staff table for salon employees
-- Created: 2026-06-08

CREATE TABLE IF NOT EXISTS staff (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    phone       TEXT,
    role        TEXT NOT NULL DEFAULT 'stylist',
    is_active   INTEGER NOT NULL DEFAULT 1,
    joined_at   DATETIME NOT NULL DEFAULT (datetime('now')),
    specialties TEXT DEFAULT '[]', -- JSON array
    created_at  DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at  DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_staff_is_active ON staff(is_active);
CREATE INDEX idx_staff_role ON staff(role);
