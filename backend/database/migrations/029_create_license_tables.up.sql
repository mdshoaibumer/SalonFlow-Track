-- License & Subscription tables

CREATE TABLE IF NOT EXISTS licenses (
    id TEXT PRIMARY KEY,
    license_key TEXT NOT NULL UNIQUE,
    customer_name TEXT NOT NULL,
    salon_name TEXT NOT NULL,
    device_id TEXT NOT NULL DEFAULT '',
    issued_date TEXT NOT NULL,
    expiry_date TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','grace_period','expired','suspended')),
    signature TEXT NOT NULL DEFAULT '',
    last_validation TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS license_events (
    id TEXT PRIMARY KEY,
    license_id TEXT NOT NULL REFERENCES licenses(id),
    event_type TEXT NOT NULL CHECK (event_type IN ('activated','renewed','expired','validated','suspended','grace_started','restricted')),
    event_date TEXT NOT NULL,
    notes TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL
);

CREATE INDEX idx_licenses_key ON licenses(license_key);
CREATE INDEX idx_licenses_status ON licenses(status);
CREATE INDEX idx_licenses_expiry ON licenses(expiry_date);
CREATE INDEX idx_license_events_license ON license_events(license_id);
CREATE INDEX idx_license_events_type ON license_events(event_type);
CREATE INDEX idx_license_events_date ON license_events(event_date);
