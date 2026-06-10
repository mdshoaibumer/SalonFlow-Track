-- Rollback: 010_create_payments_table

DROP INDEX IF EXISTS idx_payments_deleted_at;
DROP INDEX IF EXISTS idx_payments_paid_at;
DROP INDEX IF EXISTS idx_payments_method;
DROP INDEX IF EXISTS idx_payments_invoice_id;
DROP TABLE IF EXISTS payments;
