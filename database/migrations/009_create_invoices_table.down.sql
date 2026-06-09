-- Rollback: 009_create_invoices_table

DROP INDEX IF EXISTS idx_invoice_items_item_type;
DROP INDEX IF EXISTS idx_invoice_items_service_id;
DROP INDEX IF EXISTS idx_invoice_items_staff_id;
DROP INDEX IF EXISTS idx_invoice_items_invoice_id;
DROP TABLE IF EXISTS invoice_items;

DROP INDEX IF EXISTS idx_invoices_deleted_at;
DROP INDEX IF EXISTS idx_invoices_invoice_number;
DROP INDEX IF EXISTS idx_invoices_payment_status;
DROP INDEX IF EXISTS idx_invoices_status;
DROP INDEX IF EXISTS idx_invoices_date;
DROP INDEX IF EXISTS idx_invoices_customer_id;
DROP TABLE IF EXISTS invoices;

DELETE FROM settings WHERE key = 'invoice.last_sequence';
