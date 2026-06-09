-- Phase 6: Rebuild customers table for production use
DROP TABLE IF EXISTS customers;

CREATE TABLE customers (
    id TEXT PRIMARY KEY,
    customer_code TEXT NOT NULL UNIQUE,
    full_name TEXT NOT NULL,
    phone TEXT NOT NULL,
    email TEXT NOT NULL DEFAULT '',
    gender TEXT NOT NULL DEFAULT 'other' CHECK(gender IN ('male', 'female', 'other')),
    date_of_birth TEXT,
    anniversary_date TEXT,
    address TEXT NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT '',
    total_visits INTEGER NOT NULL DEFAULT 0,
    total_spent REAL NOT NULL DEFAULT 0,
    last_visit_date TEXT,
    status TEXT NOT NULL DEFAULT 'active' CHECK(status IN ('active', 'inactive')),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE UNIQUE INDEX idx_customers_phone ON customers(phone);
CREATE UNIQUE INDEX idx_customers_customer_code ON customers(customer_code);
CREATE INDEX idx_customers_full_name ON customers(full_name);
CREATE INDEX idx_customers_status ON customers(status);
