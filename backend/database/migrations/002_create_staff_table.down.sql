-- Rollback: 002_create_staff_table

DROP INDEX IF EXISTS idx_staff_role;
DROP INDEX IF EXISTS idx_staff_is_active;
DROP TABLE IF EXISTS staff;
