-- Rollback: 004_create_settings_table

DROP INDEX IF EXISTS idx_settings_category;
DROP INDEX IF EXISTS idx_settings_key;
DROP TABLE IF EXISTS settings;
