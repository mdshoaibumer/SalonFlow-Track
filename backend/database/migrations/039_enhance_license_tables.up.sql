-- Enhance license tables for offline licensing & subscription management

-- Add grace_until and last_verified_at columns to licenses
ALTER TABLE licenses ADD COLUMN grace_until TEXT NOT NULL DEFAULT '';
ALTER TABLE licenses ADD COLUMN last_verified_at TEXT NOT NULL DEFAULT '';

-- License notifications table
CREATE TABLE IF NOT EXISTS license_notifications (
    id TEXT PRIMARY KEY,
    license_id TEXT NOT NULL REFERENCES licenses(id),
    notification_type TEXT NOT NULL CHECK (notification_type IN ('7_days_remaining','3_days_remaining','1_day_remaining','expired','grace_period_remaining')),
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    is_read INTEGER NOT NULL DEFAULT 0,
    is_dismissed INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL
);

CREATE INDEX idx_license_notifications_license ON license_notifications(license_id);
CREATE INDEX idx_license_notifications_read ON license_notifications(is_read);
CREATE INDEX idx_license_notifications_type ON license_notifications(notification_type);

-- Update existing licenses to populate grace_until from expiry_date + 30 days
UPDATE licenses SET grace_until = date(expiry_date, '+30 days') WHERE grace_until = '';
UPDATE licenses SET last_verified_at = last_validation WHERE last_verified_at = '' AND last_validation != '';
