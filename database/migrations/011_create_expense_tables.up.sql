-- Migration: 011_create_expense_tables
-- Description: Create expense_categories and expenses tables
-- Created: 2026-06-08

CREATE TABLE IF NOT EXISTS expense_categories (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT DEFAULT '',
    is_default  INTEGER NOT NULL DEFAULT 0,
    is_active   INTEGER NOT NULL DEFAULT 1,
    created_at  DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at  DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at  DATETIME DEFAULT NULL
);

CREATE INDEX idx_expense_categories_is_active ON expense_categories(is_active);

-- Seed default expense categories
INSERT INTO expense_categories (id, name, description, is_default) VALUES
    ('exp-cat-rent-001', 'Rent', 'Monthly shop rent', 1),
    ('exp-cat-elec-001', 'Electricity', 'Power bills', 1),
    ('exp-cat-supp-001', 'Supplies', 'Salon supplies and consumables', 1),
    ('exp-cat-mktg-001', 'Marketing', 'Advertising and promotions', 1),
    ('exp-cat-maint-001', 'Maintenance', 'Equipment repair and upkeep', 1),
    ('exp-cat-misc-001', 'Miscellaneous', 'Other expenses', 1);

CREATE TABLE IF NOT EXISTS expenses (
    id            TEXT PRIMARY KEY,
    category_id   TEXT NOT NULL REFERENCES expense_categories(id),
    amount        REAL NOT NULL CHECK(amount > 0),
    date          DATE NOT NULL DEFAULT (date('now')),
    description   TEXT NOT NULL,
    paid_to       TEXT DEFAULT '',
    paid_by       TEXT DEFAULT '',
    receipt_path  TEXT DEFAULT '',
    created_at    DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at    DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at    DATETIME DEFAULT NULL
);

CREATE INDEX idx_expenses_category_id ON expenses(category_id);
CREATE INDEX idx_expenses_date ON expenses(date);
CREATE INDEX idx_expenses_deleted_at ON expenses(deleted_at);
