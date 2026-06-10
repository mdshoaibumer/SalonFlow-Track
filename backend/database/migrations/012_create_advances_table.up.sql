-- Migration: 012_create_advances_table
-- Description: Create advances table for staff salary advances
-- Created: 2026-06-08

CREATE TABLE IF NOT EXISTS advances (
    id               TEXT PRIMARY KEY,
    staff_id         TEXT NOT NULL REFERENCES staff(id),
    amount           REAL NOT NULL CHECK(amount > 0),
    recovered_amount REAL NOT NULL DEFAULT 0,
    balance_amount   REAL NOT NULL DEFAULT 0,
    date             DATE NOT NULL DEFAULT (date('now')),
    reason           TEXT NOT NULL,
    status           TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending', 'partial', 'recovered')),
    approved_by      TEXT DEFAULT '',
    created_at       DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at       DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at       DATETIME DEFAULT NULL
);

CREATE INDEX idx_advances_staff_id ON advances(staff_id);
CREATE INDEX idx_advances_status ON advances(status);
CREATE INDEX idx_advances_date ON advances(date);
CREATE INDEX idx_advances_deleted_at ON advances(deleted_at);
