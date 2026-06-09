-- Migration: 009_create_invoices_table
-- Description: Create invoices and invoice_items tables for billing
-- Created: 2026-06-08

CREATE TABLE IF NOT EXISTS invoices (
    id              TEXT PRIMARY KEY,
    invoice_number  TEXT NOT NULL UNIQUE,
    customer_id     TEXT NOT NULL REFERENCES customers(id),
    date            DATE NOT NULL DEFAULT (date('now')),
    subtotal        REAL NOT NULL DEFAULT 0,
    discount_type   TEXT NOT NULL DEFAULT 'none' CHECK(discount_type IN ('none', 'percentage', 'fixed')),
    discount_value  REAL NOT NULL DEFAULT 0,
    discount_amount REAL NOT NULL DEFAULT 0,
    tax_amount      REAL NOT NULL DEFAULT 0,
    total_amount    REAL NOT NULL DEFAULT 0,
    paid_amount     REAL NOT NULL DEFAULT 0,
    status          TEXT NOT NULL DEFAULT 'draft' CHECK(status IN ('draft', 'completed', 'cancelled')),
    payment_status  TEXT NOT NULL DEFAULT 'unpaid' CHECK(payment_status IN ('unpaid', 'partial', 'paid')),
    notes           TEXT DEFAULT '',
    created_at      DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at      DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at      DATETIME DEFAULT NULL
);

CREATE INDEX idx_invoices_customer_id ON invoices(customer_id);
CREATE INDEX idx_invoices_date ON invoices(date);
CREATE INDEX idx_invoices_status ON invoices(status);
CREATE INDEX idx_invoices_payment_status ON invoices(payment_status);
CREATE INDEX idx_invoices_invoice_number ON invoices(invoice_number);
CREATE INDEX idx_invoices_deleted_at ON invoices(deleted_at);

-- Invoice sequence counter (for generating invoice numbers)
INSERT INTO settings (id, key, value, description, category)
VALUES ('00000000-0000-0000-0000-000000000010', 'invoice.last_sequence', '0', 'Last invoice sequence number', 'system');

CREATE TABLE IF NOT EXISTS invoice_items (
    id          TEXT PRIMARY KEY,
    invoice_id  TEXT NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    item_type   TEXT NOT NULL DEFAULT 'service' CHECK(item_type IN ('service', 'product')),
    service_id  TEXT DEFAULT NULL REFERENCES services(id),
    product_id  TEXT DEFAULT NULL REFERENCES products(id),
    staff_id    TEXT NOT NULL REFERENCES staff(id),
    name        TEXT NOT NULL,
    quantity    INTEGER NOT NULL DEFAULT 1,
    unit_price  REAL NOT NULL DEFAULT 0,
    discount    REAL NOT NULL DEFAULT 0,
    total_price REAL NOT NULL DEFAULT 0,
    created_at  DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at  DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_invoice_items_invoice_id ON invoice_items(invoice_id);
CREATE INDEX idx_invoice_items_staff_id ON invoice_items(staff_id);
CREATE INDEX idx_invoice_items_service_id ON invoice_items(service_id);
CREATE INDEX idx_invoice_items_item_type ON invoice_items(item_type);
