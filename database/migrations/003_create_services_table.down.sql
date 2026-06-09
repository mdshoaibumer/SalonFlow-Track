-- Rollback: 003_create_services_table

DROP INDEX IF EXISTS idx_services_is_active;
DROP INDEX IF EXISTS idx_services_category;
DROP TABLE IF EXISTS services;
