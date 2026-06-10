-- Migration: 026_create_expense_management
-- Description: Create expense_categories and expenses tables for Phase 11
-- Created: 2026-06-09

-- Drop old tables from migration 011 if they exist (schema rebuilt here)
DROP TABLE IF EXISTS expenses;
DROP TABLE IF EXISTS expense_categories;

CREATE TABLE IF NOT EXISTS expense_categories (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    is_active   INTEGER NOT NULL DEFAULT 1,
    created_at  DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at  DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_expense_categories_active ON expense_categories(is_active);

-- Seed default categories
INSERT INTO expense_categories (id, name, description) VALUES
    ('cat-rent-0001', 'Rent', 'Monthly shop/salon rent'),
    ('cat-elec-0001', 'Electricity', 'Monthly electricity bills'),
    ('cat-watr-0001', 'Water', 'Monthly water bills'),
    ('cat-inet-0001', 'Internet', 'Internet and phone bills'),
    ('cat-slry-0001', 'Salary', 'Staff salary payments'),
    ('cat-invp-0001', 'Inventory Purchase', 'Products and inventory purchases'),
    ('cat-mktg-0001', 'Marketing', 'Advertising and promotional expenses'),
    ('cat-mntc-0001', 'Maintenance', 'Repair and maintenance costs'),
    ('cat-eqpt-0001', 'Equipment', 'Equipment purchases and leases'),
    ('cat-misc-0001', 'Miscellaneous', 'Other miscellaneous expenses');

CREATE TABLE IF NOT EXISTS expenses (
    id                TEXT PRIMARY KEY,
    expense_number    TEXT NOT NULL UNIQUE,
    category_id       TEXT NOT NULL REFERENCES expense_categories(id),
    amount            REAL NOT NULL CHECK(amount > 0),
    expense_date      DATE NOT NULL,
    payment_method    TEXT NOT NULL CHECK(payment_method IN ('cash','upi','bank_transfer','card','cheque')),
    vendor_name       TEXT NOT NULL DEFAULT '',
    invoice_reference TEXT NOT NULL DEFAULT '',
    description       TEXT NOT NULL DEFAULT '',
    attachment_path   TEXT NOT NULL DEFAULT '',
    status            TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending','approved','paid','rejected')),
    created_by        TEXT NOT NULL DEFAULT '',
    created_at        DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at        DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_expenses_category_id ON expenses(category_id);
CREATE INDEX idx_expenses_expense_date ON expenses(expense_date);
CREATE INDEX idx_expenses_status ON expenses(status);
CREATE INDEX idx_expenses_payment_method ON expenses(payment_method);
CREATE INDEX idx_expenses_expense_number ON expenses(expense_number);

-- Sequence table for expense numbers
CREATE TABLE IF NOT EXISTS expense_number_seq (
    year    INTEGER NOT NULL,
    seq     INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (year)
);
