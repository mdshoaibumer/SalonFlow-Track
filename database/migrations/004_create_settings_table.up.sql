-- Migration: 004_create_settings_table
-- Description: Create the settings table for application configuration
-- Created: 2026-06-08

CREATE TABLE IF NOT EXISTS settings (
    id          TEXT PRIMARY KEY,
    key         TEXT NOT NULL UNIQUE,
    value       TEXT NOT NULL DEFAULT '',
    description TEXT DEFAULT '',
    category    TEXT NOT NULL DEFAULT 'general',
    created_at  DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at  DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE UNIQUE INDEX idx_settings_key ON settings(key);
CREATE INDEX idx_settings_category ON settings(category);

-- Seed default settings
INSERT INTO settings (id, key, value, description, category) VALUES
    ('00000000-0000-0000-0000-000000000001', 'salon.name', 'My Salon', 'Salon business name', 'salon'),
    ('00000000-0000-0000-0000-000000000002', 'salon.phone', '', 'Salon contact phone', 'salon'),
    ('00000000-0000-0000-0000-000000000003', 'salon.address', '', 'Salon physical address', 'salon'),
    ('00000000-0000-0000-0000-000000000004', 'salon.opening_time', '09:00', 'Daily opening time', 'salon'),
    ('00000000-0000-0000-0000-000000000005', 'salon.closing_time', '21:00', 'Daily closing time', 'salon'),
    ('00000000-0000-0000-0000-000000000006', 'app.theme', 'system', 'UI theme preference', 'app'),
    ('00000000-0000-0000-0000-000000000007', 'app.language', 'en', 'Application language', 'app');
