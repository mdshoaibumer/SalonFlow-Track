-- Migration: 018_rebuild_staff_table
-- Description: Rebuild staff table with Phase 4 schema for Staff Management module
-- Created: 2026-06-09

DROP INDEX IF EXISTS idx_staff_is_active;
DROP INDEX IF EXISTS idx_staff_role;
DROP TABLE IF EXISTS staff;

CREATE TABLE staff (
    id                      TEXT PRIMARY KEY,
    staff_code              TEXT NOT NULL UNIQUE,
    full_name               TEXT NOT NULL,
    phone                   TEXT NOT NULL,
    email                   TEXT DEFAULT '',
    gender                  TEXT NOT NULL DEFAULT 'male' CHECK (gender IN ('male', 'female', 'other')),
    designation             TEXT NOT NULL DEFAULT 'stylist' CHECK (designation IN ('stylist', 'assistant', 'receptionist', 'manager')),
    joining_date            TEXT NOT NULL,
    base_salary             REAL NOT NULL DEFAULT 0 CHECK (base_salary >= 0),
    commission_percentage   REAL NOT NULL DEFAULT 0 CHECK (commission_percentage >= 0 AND commission_percentage <= 100),
    status                  TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    created_at              TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at              TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE UNIQUE INDEX idx_staff_phone ON staff(phone);
CREATE INDEX idx_staff_status ON staff(status);
CREATE INDEX idx_staff_designation ON staff(designation);
CREATE INDEX idx_staff_staff_code ON staff(staff_code);
CREATE INDEX idx_staff_full_name ON staff(full_name);
