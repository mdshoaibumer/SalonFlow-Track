-- Migration: 015_create_products_table
-- Description: Create products and stock_transactions tables for inventory
-- Created: 2026-06-08

CREATE TABLE IF NOT EXISTS products (
    id              TEXT PRIMARY KEY,
    name            TEXT NOT NULL,
    brand           TEXT DEFAULT '',
    sku             TEXT DEFAULT '',
    category        TEXT NOT NULL DEFAULT 'general',
    unit            TEXT NOT NULL DEFAULT 'piece' CHECK(unit IN ('ml', 'g', 'kg', 'l', 'piece', 'bottle', 'tube', 'pack')),
    cost_price      REAL NOT NULL DEFAULT 0 CHECK(cost_price >= 0),
    selling_price   REAL NOT NULL DEFAULT 0 CHECK(selling_price >= 0),
    current_stock   REAL NOT NULL DEFAULT 0,
    min_stock_level REAL NOT NULL DEFAULT 1,
    is_active       INTEGER NOT NULL DEFAULT 1,
    is_for_sale     INTEGER NOT NULL DEFAULT 0,
    created_at      DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at      DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at      DATETIME DEFAULT NULL
);

CREATE UNIQUE INDEX idx_products_sku ON products(sku) WHERE sku != '' AND deleted_at IS NULL;
CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_is_active ON products(is_active);
CREATE INDEX idx_products_low_stock ON products(current_stock, min_stock_level) WHERE is_active = 1;
CREATE INDEX idx_products_deleted_at ON products(deleted_at);

CREATE TABLE IF NOT EXISTS stock_transactions (
    id               TEXT PRIMARY KEY,
    product_id       TEXT NOT NULL REFERENCES products(id),
    type             TEXT NOT NULL CHECK(type IN ('purchase', 'consumption', 'sale', 'adjustment', 'return', 'damage')),
    quantity         REAL NOT NULL,
    unit_cost        REAL NOT NULL DEFAULT 0,
    total_cost       REAL NOT NULL DEFAULT 0,
    balance_after    REAL NOT NULL DEFAULT 0,
    reference_type   TEXT DEFAULT '',
    reference_id     TEXT DEFAULT NULL,
    notes            TEXT DEFAULT '',
    performed_by     TEXT DEFAULT '',
    transaction_date DATE NOT NULL DEFAULT (date('now')),
    created_at       DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at       DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_stock_transactions_product_id ON stock_transactions(product_id);
CREATE INDEX idx_stock_transactions_type ON stock_transactions(type);
CREATE INDEX idx_stock_transactions_date ON stock_transactions(transaction_date);
CREATE INDEX idx_stock_transactions_reference ON stock_transactions(reference_type, reference_id);
