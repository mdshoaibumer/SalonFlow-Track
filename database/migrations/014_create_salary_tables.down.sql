-- Rollback: 014_create_salary_tables

DROP INDEX IF EXISTS idx_salary_line_items_type;
DROP INDEX IF EXISTS idx_salary_line_items_salary_id;
DROP TABLE IF EXISTS salary_line_items;

DROP INDEX IF EXISTS idx_salaries_deleted_at;
DROP INDEX IF EXISTS idx_salaries_period;
DROP INDEX IF EXISTS idx_salaries_status;
DROP INDEX IF EXISTS idx_salaries_staff_id;
DROP INDEX IF EXISTS idx_salaries_staff_month_year;
DROP TABLE IF EXISTS salaries;
