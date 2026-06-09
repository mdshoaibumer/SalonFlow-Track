-- Migration: 026_create_expense_management (DOWN)
DROP TABLE IF EXISTS expense_number_seq;
DROP TABLE IF EXISTS expenses;
DROP TABLE IF EXISTS expense_categories;
