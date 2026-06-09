-- Migration: 008_create_service_categories_table
-- Description: Create service_categories for grouping salon services
-- Created: 2026-06-08

CREATE TABLE IF NOT EXISTS service_categories (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT DEFAULT '',
    sort_order  INTEGER NOT NULL DEFAULT 0,
    is_active   INTEGER NOT NULL DEFAULT 1,
    created_at  DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at  DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at  DATETIME DEFAULT NULL
);

CREATE INDEX idx_service_categories_is_active ON service_categories(is_active);

-- Add category_id FK to existing services table
ALTER TABLE services ADD COLUMN category_id TEXT DEFAULT NULL REFERENCES service_categories(id);

-- Seed default categories
INSERT INTO service_categories (id, name, description, sort_order) VALUES
    ('cat-hair-001', 'Hair', 'Haircuts, styling, coloring', 1),
    ('cat-skin-001', 'Skin', 'Facials, cleanup, skincare', 2),
    ('cat-spa-001', 'Spa', 'Massage, body treatments', 3),
    ('cat-nails-001', 'Nails', 'Manicure, pedicure, nail art', 4),
    ('cat-bridal-001', 'Bridal', 'Bridal packages, makeup', 5),
    ('cat-other-001', 'Other', 'Miscellaneous services', 6);
