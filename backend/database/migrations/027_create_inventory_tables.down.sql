-- Migration: 027_create_inventory_tables (DOWN)
DROP TABLE IF EXISTS purchase_items;
DROP TABLE IF EXISTS purchase_entries;
DROP TABLE IF EXISTS stock_transactions;
DROP TABLE IF EXISTS products;
DELETE FROM product_code_seq WHERE prefix IN ('PRD', 'PUR');
