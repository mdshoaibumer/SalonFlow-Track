-- Migration: 003_create_services_table
-- Description: Create the services table for salon service offerings
-- Created: 2026-06-08

CREATE TABLE IF NOT EXISTS services (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT DEFAULT '',
    duration    INTEGER NOT NULL DEFAULT 30, -- minutes
    price       REAL NOT NULL DEFAULT 0.0,
    category    TEXT NOT NULL DEFAULT 'general',
    is_active   INTEGER NOT NULL DEFAULT 1,
    created_at  DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at  DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_services_category ON services(category);
CREATE INDEX idx_services_is_active ON services(is_active);
