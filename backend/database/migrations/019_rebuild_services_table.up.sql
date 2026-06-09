-- Phase 5: Rebuild services table for production use
DROP TABLE IF EXISTS services;

CREATE TABLE services (
    id TEXT PRIMARY KEY,
    service_code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    category TEXT NOT NULL CHECK(category IN ('hair', 'facial', 'skin', 'spa', 'massage', 'coloring', 'treatment', 'other')),
    description TEXT NOT NULL DEFAULT '',
    duration_minutes INTEGER NOT NULL CHECK(duration_minutes > 0),
    price REAL NOT NULL CHECK(price > 0),
    cost_price REAL NOT NULL DEFAULT 0 CHECK(cost_price >= 0),
    commission_type TEXT NOT NULL DEFAULT 'percentage' CHECK(commission_type IN ('fixed', 'percentage')),
    commission_value REAL NOT NULL DEFAULT 0 CHECK(commission_value >= 0),
    status TEXT NOT NULL DEFAULT 'active' CHECK(status IN ('active', 'inactive')),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE INDEX idx_services_category ON services(category);
CREATE INDEX idx_services_status ON services(status);
CREATE INDEX idx_services_name ON services(name);
CREATE UNIQUE INDEX idx_services_service_code ON services(service_code);
