-- Commission Rules table
CREATE TABLE IF NOT EXISTS commission_rules (
    id TEXT PRIMARY KEY,
    rule_name TEXT NOT NULL,
    rule_type TEXT NOT NULL CHECK (rule_type IN ('revenue_based', 'service_based', 'fixed')),
    target_type TEXT NOT NULL CHECK (target_type IN ('global', 'staff', 'service')),
    target_id TEXT,
    calculation_type TEXT NOT NULL CHECK (calculation_type IN ('percentage', 'fixed_amount', 'tiered')),
    calculation_value REAL NOT NULL DEFAULT 0,
    minimum_target REAL NOT NULL DEFAULT 0,
    maximum_target REAL NOT NULL DEFAULT 0,
    is_active INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_commission_rules_rule_type ON commission_rules(rule_type);
CREATE INDEX IF NOT EXISTS idx_commission_rules_target ON commission_rules(target_type, target_id);
CREATE INDEX IF NOT EXISTS idx_commission_rules_active ON commission_rules(is_active);

-- Commission Transactions table
CREATE TABLE IF NOT EXISTS commission_transactions (
    id TEXT PRIMARY KEY,
    staff_id TEXT NOT NULL REFERENCES staff(id),
    invoice_id TEXT NOT NULL REFERENCES invoices(id),
    rule_id TEXT REFERENCES commission_rules(id),
    revenue_amount REAL NOT NULL DEFAULT 0,
    commission_amount REAL NOT NULL DEFAULT 0,
    business_date TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'paid')),
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_commission_transactions_staff_id ON commission_transactions(staff_id);
CREATE INDEX IF NOT EXISTS idx_commission_transactions_invoice_id ON commission_transactions(invoice_id);
CREATE INDEX IF NOT EXISTS idx_commission_transactions_business_date ON commission_transactions(business_date);
CREATE INDEX IF NOT EXISTS idx_commission_transactions_status ON commission_transactions(status);
