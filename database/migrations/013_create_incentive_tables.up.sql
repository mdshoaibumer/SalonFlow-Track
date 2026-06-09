-- Migration: 013_create_incentive_tables
-- Description: Create incentive_rules, incentive_rule_slabs, and staff_incentives tables
-- Created: 2026-06-08

CREATE TABLE IF NOT EXISTS incentive_rules (
    id             TEXT PRIMARY KEY,
    name           TEXT NOT NULL,
    type           TEXT NOT NULL CHECK(type IN ('revenue_slab', 'service_based', 'product_sale')),
    period         TEXT NOT NULL DEFAULT 'monthly' CHECK(period IN ('daily', 'weekly', 'monthly')),
    service_id     TEXT DEFAULT NULL REFERENCES services(id),
    category_id    TEXT DEFAULT NULL REFERENCES service_categories(id),
    staff_id       TEXT DEFAULT NULL REFERENCES staff(id),
    is_active      INTEGER NOT NULL DEFAULT 1,
    effective_from DATE NOT NULL DEFAULT (date('now')),
    effective_to   DATE DEFAULT NULL,
    created_at     DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at     DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at     DATETIME DEFAULT NULL
);

CREATE INDEX idx_incentive_rules_type ON incentive_rules(type);
CREATE INDEX idx_incentive_rules_is_active ON incentive_rules(is_active);
CREATE INDEX idx_incentive_rules_staff_id ON incentive_rules(staff_id);
CREATE INDEX idx_incentive_rules_deleted_at ON incentive_rules(deleted_at);

CREATE TABLE IF NOT EXISTS incentive_rule_slabs (
    id               TEXT PRIMARY KEY,
    rule_id          TEXT NOT NULL REFERENCES incentive_rules(id) ON DELETE CASCADE,
    min_amount       REAL NOT NULL DEFAULT 0,
    max_amount       REAL DEFAULT NULL,
    commission_type  TEXT NOT NULL CHECK(commission_type IN ('percentage', 'fixed')),
    commission_value REAL NOT NULL DEFAULT 0,
    sort_order       INTEGER NOT NULL DEFAULT 0,
    created_at       DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at       DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_incentive_rule_slabs_rule_id ON incentive_rule_slabs(rule_id);

CREATE TABLE IF NOT EXISTS staff_incentives (
    id               TEXT PRIMARY KEY,
    staff_id         TEXT NOT NULL REFERENCES staff(id),
    rule_id          TEXT NOT NULL REFERENCES incentive_rules(id),
    period_start     DATE NOT NULL,
    period_end       DATE NOT NULL,
    revenue_amount   REAL NOT NULL DEFAULT 0,
    incentive_amount REAL NOT NULL DEFAULT 0,
    status           TEXT NOT NULL DEFAULT 'calculated' CHECK(status IN ('calculated', 'approved', 'paid')),
    approved_by      TEXT DEFAULT '',
    approved_at      DATETIME DEFAULT NULL,
    created_at       DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at       DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_staff_incentives_staff_id ON staff_incentives(staff_id);
CREATE INDEX idx_staff_incentives_rule_id ON staff_incentives(rule_id);
CREATE INDEX idx_staff_incentives_period ON staff_incentives(period_start, period_end);
CREATE INDEX idx_staff_incentives_status ON staff_incentives(status);
