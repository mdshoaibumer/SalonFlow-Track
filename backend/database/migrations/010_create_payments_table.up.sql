-- Migration: 010_create_payments_table
-- Description: Create payments table for invoice payment tracking
-- Created: 2026-06-08

CREATE TABLE IF NOT EXISTS payments (
    id            TEXT PRIMARY KEY,
    invoice_id    TEXT NOT NULL REFERENCES invoices(id),
    amount        REAL NOT NULL CHECK(amount > 0),
    method        TEXT NOT NULL CHECK(method IN ('cash', 'card', 'upi')),
    reference_no  TEXT DEFAULT '',
    paid_at       DATETIME NOT NULL DEFAULT (datetime('now')),
    notes         TEXT DEFAULT '',
    is_refund     INTEGER NOT NULL DEFAULT 0,
    created_at    DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at    DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at    DATETIME DEFAULT NULL
);

CREATE INDEX idx_payments_invoice_id ON payments(invoice_id);
CREATE INDEX idx_payments_method ON payments(method);
CREATE INDEX idx_payments_paid_at ON payments(paid_at);
CREATE INDEX idx_payments_deleted_at ON payments(deleted_at);
