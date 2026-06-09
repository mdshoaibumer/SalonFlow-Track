-- Rollback: 011_create_expense_tables

DROP INDEX IF EXISTS idx_expenses_deleted_at;
DROP INDEX IF EXISTS idx_expenses_date;
DROP INDEX IF EXISTS idx_expenses_category_id;
DROP TABLE IF EXISTS expenses;

DROP INDEX IF EXISTS idx_expense_categories_is_active;
DROP TABLE IF EXISTS expense_categories;
