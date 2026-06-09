-- Rollback: 017_create_performance_views

DROP VIEW IF EXISTS v_pending_advances;
DROP VIEW IF EXISTS v_low_stock_products;
DROP VIEW IF EXISTS v_daily_revenue_summary;
DROP VIEW IF EXISTS v_monthly_staff_performance;
DROP VIEW IF EXISTS v_daily_staff_performance;
