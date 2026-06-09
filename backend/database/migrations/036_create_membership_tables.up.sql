-- Membership Plans
CREATE TABLE IF NOT EXISTS membership_plans (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    plan_type TEXT NOT NULL DEFAULT 'package' CHECK (plan_type IN ('package','membership')),
    price REAL NOT NULL DEFAULT 0,
    duration_days INTEGER NOT NULL DEFAULT 365,
    max_sessions INTEGER NOT NULL DEFAULT 0,
    discount_percentage REAL NOT NULL DEFAULT 0,
    priority_booking INTEGER NOT NULL DEFAULT 0,
    is_active INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

-- Package Services (services included in a plan)
CREATE TABLE IF NOT EXISTS package_services (
    id TEXT PRIMARY KEY,
    plan_id TEXT NOT NULL REFERENCES membership_plans(id) ON DELETE CASCADE,
    service_id TEXT NOT NULL DEFAULT '',
    service_name TEXT NOT NULL DEFAULT '',
    sessions_included INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL
);

-- Member Subscriptions (customer purchases)
CREATE TABLE IF NOT EXISTS member_subscriptions (
    id TEXT PRIMARY KEY,
    customer_id TEXT NOT NULL DEFAULT '',
    plan_id TEXT NOT NULL REFERENCES membership_plans(id),
    plan_name TEXT NOT NULL DEFAULT '',
    start_date TEXT NOT NULL,
    end_date TEXT NOT NULL,
    total_sessions INTEGER NOT NULL DEFAULT 0,
    used_sessions INTEGER NOT NULL DEFAULT 0,
    amount_paid REAL NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','expired','cancelled','paused')),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE INDEX idx_membership_plans_type ON membership_plans(plan_type);
CREATE INDEX idx_package_services_plan ON package_services(plan_id);
CREATE INDEX idx_member_subscriptions_customer ON member_subscriptions(customer_id);
CREATE INDEX idx_member_subscriptions_status ON member_subscriptions(status);
