-- Rollback: 015_create_products_table

DROP INDEX IF EXISTS idx_stock_transactions_reference;
DROP INDEX IF EXISTS idx_stock_transactions_date;
DROP INDEX IF EXISTS idx_stock_transactions_type;
DROP INDEX IF EXISTS idx_stock_transactions_product_id;
DROP TABLE IF EXISTS stock_transactions;

DROP INDEX IF EXISTS idx_products_deleted_at;
DROP INDEX IF EXISTS idx_products_low_stock;
DROP INDEX IF EXISTS idx_products_is_active;
DROP INDEX IF EXISTS idx_products_category;
DROP INDEX IF EXISTS idx_products_sku;
DROP TABLE IF EXISTS products;
