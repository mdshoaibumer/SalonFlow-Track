-- Migration: 027_create_inventory_tables
-- Description: Create products, stock_transactions, purchase_entries, purchase_items tables
-- Created: 2026-06-09

CREATE TABLE IF NOT EXISTS products (
    id             TEXT PRIMARY KEY,
    product_code   TEXT NOT NULL UNIQUE,
    name           TEXT NOT NULL,
    category       TEXT NOT NULL CHECK(category IN ('hair_care','facial','spa','coloring','treatment','retail','equipment','other')),
    brand          TEXT NOT NULL DEFAULT '',
    unit           TEXT NOT NULL DEFAULT 'pcs',
    sku            TEXT NOT NULL DEFAULT '',
    purchase_price REAL NOT NULL DEFAULT 0 CHECK(purchase_price >= 0),
    selling_price  REAL NOT NULL DEFAULT 0 CHECK(selling_price >= 0),
    current_stock  REAL NOT NULL DEFAULT 0,
    minimum_stock  REAL NOT NULL DEFAULT 0,
    maximum_stock  REAL NOT NULL DEFAULT 0,
    status         TEXT NOT NULL DEFAULT 'active' CHECK(status IN ('active','inactive','discontinued')),
    created_at     DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at     DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_product_code ON products(product_code);
CREATE INDEX idx_products_low_stock ON products(current_stock, minimum_stock);

CREATE TABLE IF NOT EXISTS stock_transactions (
    id               TEXT PRIMARY KEY,
    product_id       TEXT NOT NULL REFERENCES products(id),
    transaction_type TEXT NOT NULL CHECK(transaction_type IN ('purchase','consumption','sale','adjustment','return','damage')),
    quantity         REAL NOT NULL,
    unit_cost        REAL NOT NULL DEFAULT 0,
    reference_type   TEXT NOT NULL DEFAULT '',
    reference_id     TEXT NOT NULL DEFAULT '',
    notes            TEXT NOT NULL DEFAULT '',
    transaction_date DATE NOT NULL,
    created_at       DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at       DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_stock_transactions_product_id ON stock_transactions(product_id);
CREATE INDEX idx_stock_transactions_type ON stock_transactions(transaction_type);
CREATE INDEX idx_stock_transactions_date ON stock_transactions(transaction_date);

CREATE TABLE IF NOT EXISTS purchase_entries (
    id              TEXT PRIMARY KEY,
    purchase_number TEXT NOT NULL UNIQUE,
    vendor_name     TEXT NOT NULL,
    invoice_number  TEXT NOT NULL DEFAULT '',
    purchase_date   DATE NOT NULL,
    total_amount    REAL NOT NULL DEFAULT 0,
    notes           TEXT NOT NULL DEFAULT '',
    created_at      DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at      DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_purchase_entries_date ON purchase_entries(purchase_date);
CREATE INDEX idx_purchase_entries_vendor ON purchase_entries(vendor_name);

CREATE TABLE IF NOT EXISTS purchase_items (
    id                TEXT PRIMARY KEY,
    purchase_entry_id TEXT NOT NULL REFERENCES purchase_entries(id),
    product_id        TEXT NOT NULL REFERENCES products(id),
    quantity          REAL NOT NULL CHECK(quantity > 0),
    unit_price        REAL NOT NULL CHECK(unit_price >= 0),
    line_total        REAL NOT NULL DEFAULT 0,
    created_at        DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at        DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_purchase_items_entry ON purchase_items(purchase_entry_id);
CREATE INDEX idx_purchase_items_product ON purchase_items(product_id);

-- Sequence table for product codes and purchase numbers
CREATE TABLE IF NOT EXISTS product_code_seq (
    prefix TEXT NOT NULL,
    seq    INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (prefix)
);

INSERT INTO product_code_seq (prefix, seq) VALUES ('PRD', 0);
INSERT INTO product_code_seq (prefix, seq) VALUES ('PUR', 0);
