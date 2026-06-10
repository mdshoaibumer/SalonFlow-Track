-- Rollback: 005_create_license_table

DROP INDEX IF EXISTS idx_license_status;
DROP TABLE IF EXISTS license;
