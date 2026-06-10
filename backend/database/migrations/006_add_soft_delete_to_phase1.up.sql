-- Migration: 006_add_soft_delete_to_phase1_tables
-- Description: Add deleted_at column for soft delete support to Phase 1 tables
-- Created: 2026-06-08

ALTER TABLE users ADD COLUMN deleted_at DATETIME DEFAULT NULL;
ALTER TABLE staff ADD COLUMN deleted_at DATETIME DEFAULT NULL;
ALTER TABLE staff ADD COLUMN base_salary REAL NOT NULL DEFAULT 0;
ALTER TABLE services ADD COLUMN deleted_at DATETIME DEFAULT NULL;
