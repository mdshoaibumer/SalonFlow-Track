-- Staff Performance Daily table
CREATE TABLE IF NOT EXISTS staff_performance_daily (
    id TEXT PRIMARY KEY,
    staff_id TEXT NOT NULL REFERENCES staff(id),
    business_date TEXT NOT NULL,
    invoice_count INTEGER NOT NULL DEFAULT 0,
    customer_count INTEGER NOT NULL DEFAULT 0,
    service_count INTEGER NOT NULL DEFAULT 0,
    revenue REAL NOT NULL DEFAULT 0,
    commission_amount REAL NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_staff_performance_daily_staff_id ON staff_performance_daily(staff_id);
CREATE INDEX IF NOT EXISTS idx_staff_performance_daily_business_date ON staff_performance_daily(business_date);
CREATE UNIQUE INDEX IF NOT EXISTS idx_staff_performance_daily_staff_date ON staff_performance_daily(staff_id, business_date);
