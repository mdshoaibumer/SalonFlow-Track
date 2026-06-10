-- Migration: 014_create_salary_tables
-- Description: Create salaries and salary_line_items tables
-- Created: 2026-06-08

CREATE TABLE IF NOT EXISTS salaries (
    id               TEXT PRIMARY KEY,
    staff_id         TEXT NOT NULL REFERENCES staff(id),
    month            INTEGER NOT NULL CHECK(month BETWEEN 1 AND 12),
    year             INTEGER NOT NULL CHECK(year >= 2020),
    base_salary      REAL NOT NULL DEFAULT 0,
    total_earnings   REAL NOT NULL DEFAULT 0,
    total_deductions REAL NOT NULL DEFAULT 0,
    net_salary       REAL NOT NULL DEFAULT 0,
    status           TEXT NOT NULL DEFAULT 'draft' CHECK(status IN ('draft', 'approved', 'paid')),
    paid_at          DATETIME DEFAULT NULL,
    paid_via         TEXT DEFAULT '' CHECK(paid_via IN ('', 'cash', 'bank_transfer', 'upi')),
    notes            TEXT DEFAULT '',
    generated_by     TEXT DEFAULT '',
    approved_by      TEXT DEFAULT '',
    created_at       DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at       DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at       DATETIME DEFAULT NULL
);

-- One salary per staff per month
CREATE UNIQUE INDEX idx_salaries_staff_month_year ON salaries(staff_id, month, year) WHERE deleted_at IS NULL;
CREATE INDEX idx_salaries_staff_id ON salaries(staff_id);
CREATE INDEX idx_salaries_status ON salaries(status);
CREATE INDEX idx_salaries_period ON salaries(year, month);
CREATE INDEX idx_salaries_deleted_at ON salaries(deleted_at);

CREATE TABLE IF NOT EXISTS salary_line_items (
    id           TEXT PRIMARY KEY,
    salary_id    TEXT NOT NULL REFERENCES salaries(id) ON DELETE CASCADE,
    type         TEXT NOT NULL CHECK(type IN ('base_pay', 'commission', 'bonus', 'advance_deduction', 'other_deduction', 'other_earning')),
    description  TEXT NOT NULL,
    amount       REAL NOT NULL,
    reference_id TEXT DEFAULT NULL,
    sort_order   INTEGER NOT NULL DEFAULT 0,
    created_at   DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at   DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_salary_line_items_salary_id ON salary_line_items(salary_id);
CREATE INDEX idx_salary_line_items_type ON salary_line_items(type);
