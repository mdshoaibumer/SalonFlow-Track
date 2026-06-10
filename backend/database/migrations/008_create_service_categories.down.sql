-- Rollback: 008_create_service_categories_table

DROP INDEX IF EXISTS idx_service_categories_is_active;
DROP TABLE IF EXISTS service_categories;
