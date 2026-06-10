-- Migration: 007_create_customers_table
-- Description: Create customers table for salon clients
-- Created: 2026-06-08

CREATE TABLE IF NOT EXISTS customers (
    id            TEXT PRIMARY KEY,
    name          TEXT NOT NULL,
    phone         TEXT NOT NULL,
    email         TEXT DEFAULT '',
    gender        TEXT NOT NULL DEFAULT 'other' CHECK(gender IN ('male', 'female', 'other')),
    date_of_birth DATE DEFAULT NULL,
    notes         TEXT DEFAULT '',
    visit_count   INTEGER NOT NULL DEFAULT 0,
    last_visit_at DATETIME DEFAULT NULL,
    is_active     INTEGER NOT NULL DEFAULT 1,
    created_at    DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at    DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at    DATETIME DEFAULT NULL
);

CREATE UNIQUE INDEX idx_customers_phone ON customers(phone) WHERE deleted_at IS NULL;
CREATE INDEX idx_customers_name ON customers(name);
CREATE INDEX idx_customers_is_active ON customers(is_active);
CREATE INDEX idx_customers_deleted_at ON customers(deleted_at);
