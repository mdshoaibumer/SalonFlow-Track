-- Rollback: 012_create_advances_table

DROP INDEX IF EXISTS idx_advances_deleted_at;
DROP INDEX IF EXISTS idx_advances_date;
DROP INDEX IF EXISTS idx_advances_status;
DROP INDEX IF EXISTS idx_advances_staff_id;
DROP TABLE IF EXISTS advances;
