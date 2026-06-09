-- Migration: 005_create_license_table
-- Description: Create the license table for application licensing
-- Created: 2026-06-08

CREATE TABLE IF NOT EXISTS license (
    id          TEXT PRIMARY KEY,
    license_key TEXT NOT NULL,
    expiry_date DATETIME NOT NULL,
    status      TEXT NOT NULL DEFAULT 'active', -- active, expired, revoked
    issued_to   TEXT NOT NULL DEFAULT '',
    issued_at   DATETIME NOT NULL DEFAULT (datetime('now')),
    created_at  DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at  DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_license_status ON license(status);
