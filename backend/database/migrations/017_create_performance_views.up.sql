-- Migration: 017_create_performance_views
-- Description: Create views for staff performance reporting (DSR)
-- Created: 2026-06-08

-- Daily Staff Performance View (DSR)
CREATE VIEW IF NOT EXISTS v_daily_staff_performance AS
SELECT 
    ii.staff_id,
    s.name AS staff_name,
    DATE(i.date) AS work_date,
    COUNT(DISTINCT i.customer_id) AS customer_count,
    COUNT(ii.id) AS item_count,
    SUM(CASE WHEN ii.item_type = 'service' THEN ii.total_price ELSE 0 END) AS service_revenue,
    SUM(CASE WHEN ii.item_type = 'product' THEN ii.total_price ELSE 0 END) AS product_revenue,
    SUM(ii.total_price) AS total_revenue
FROM invoice_items ii
JOIN invoices i ON i.id = ii.invoice_id
JOIN staff s ON s.id = ii.staff_id
WHERE i.status = 'completed'
  AND i.deleted_at IS NULL
GROUP BY ii.staff_id, s.name, DATE(i.date);

-- Monthly Staff Performance View
CREATE VIEW IF NOT EXISTS v_monthly_staff_performance AS
SELECT 
    ii.staff_id,
    s.name AS staff_name,
    strftime('%Y', i.date) AS year,
    strftime('%m', i.date) AS month,
    COUNT(DISTINCT i.customer_id) AS customer_count,
    COUNT(DISTINCT DATE(i.date)) AS working_days,
    COUNT(ii.id) AS item_count,
    SUM(CASE WHEN ii.item_type = 'service' THEN ii.total_price ELSE 0 END) AS service_revenue,
    SUM(CASE WHEN ii.item_type = 'product' THEN ii.total_price ELSE 0 END) AS product_revenue,
    SUM(ii.total_price) AS total_revenue,
    ROUND(SUM(ii.total_price) / COUNT(DISTINCT DATE(i.date)), 2) AS avg_daily_revenue
FROM invoice_items ii
JOIN invoices i ON i.id = ii.invoice_id
JOIN staff s ON s.id = ii.staff_id
WHERE i.status = 'completed'
  AND i.deleted_at IS NULL
GROUP BY ii.staff_id, s.name, strftime('%Y', i.date), strftime('%m', i.date);

-- Daily Revenue Summary (All Staff Combined)
CREATE VIEW IF NOT EXISTS v_daily_revenue_summary AS
SELECT 
    DATE(i.date) AS date,
    COUNT(DISTINCT i.id) AS invoice_count,
    COUNT(DISTINCT i.customer_id) AS customer_count,
    SUM(i.total_amount) AS total_revenue,
    SUM(CASE WHEN p.method = 'cash' THEN p.amount ELSE 0 END) AS cash_collected,
    SUM(CASE WHEN p.method = 'card' THEN p.amount ELSE 0 END) AS card_collected,
    SUM(CASE WHEN p.method = 'upi' THEN p.amount ELSE 0 END) AS upi_collected
FROM invoices i
LEFT JOIN payments p ON p.invoice_id = i.id AND p.deleted_at IS NULL AND p.is_refund = 0
WHERE i.status = 'completed'
  AND i.deleted_at IS NULL
GROUP BY DATE(i.date);

-- Low Stock Alert View
CREATE VIEW IF NOT EXISTS v_low_stock_products AS
SELECT 
    id,
    name,
    brand,
    category,
    unit,
    current_stock,
    min_stock_level,
    cost_price,
    (current_stock * cost_price) AS stock_value
FROM products
WHERE is_active = 1
  AND deleted_at IS NULL
  AND current_stock <= min_stock_level
ORDER BY current_stock ASC;

-- Pending Advances View
CREATE VIEW IF NOT EXISTS v_pending_advances AS
SELECT 
    a.id,
    a.staff_id,
    s.name AS staff_name,
    a.amount,
    a.recovered_amount,
    a.balance_amount,
    a.date,
    a.reason,
    a.status
FROM advances a
JOIN staff s ON s.id = a.staff_id
WHERE a.status IN ('pending', 'partial')
  AND a.deleted_at IS NULL
ORDER BY a.date ASC;
