-- Drop old tables from migration 012/014 if they exist (schema rebuilt here)
DROP TABLE IF EXISTS salary_line_items;
DROP TABLE IF EXISTS salaries;
DROP TABLE IF EXISTS advances;

-- Salary Cycles table
CREATE TABLE IF NOT EXISTS salary_cycles (
    id TEXT PRIMARY KEY,
    month INTEGER NOT NULL CHECK (month BETWEEN 1 AND 12),
    year INTEGER NOT NULL CHECK (year >= 2020),
    status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'generated', 'finalized')),
    generated_at TEXT,
    generated_by TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_salary_cycles_month_year ON salary_cycles(month, year);
CREATE INDEX IF NOT EXISTS idx_salary_cycles_status ON salary_cycles(status);

-- Salary Records table
CREATE TABLE IF NOT EXISTS salary_records (
    id TEXT PRIMARY KEY,
    salary_cycle_id TEXT NOT NULL REFERENCES salary_cycles(id),
    staff_id TEXT NOT NULL REFERENCES staff(id),
    base_salary REAL NOT NULL DEFAULT 0,
    commission_amount REAL NOT NULL DEFAULT 0,
    bonus_amount REAL NOT NULL DEFAULT 0,
    advance_amount REAL NOT NULL DEFAULT 0,
    deduction_amount REAL NOT NULL DEFAULT 0,
    gross_salary REAL NOT NULL DEFAULT 0,
    net_salary REAL NOT NULL DEFAULT 0,
    payment_status TEXT NOT NULL DEFAULT 'pending' CHECK (payment_status IN ('pending', 'partial', 'paid')),
    payment_date TEXT,
    notes TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_salary_records_cycle ON salary_records(salary_cycle_id);
CREATE INDEX IF NOT EXISTS idx_salary_records_staff ON salary_records(staff_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_salary_records_cycle_staff ON salary_records(salary_cycle_id, staff_id);
CREATE INDEX IF NOT EXISTS idx_salary_records_status ON salary_records(payment_status);

-- Advances table
CREATE TABLE IF NOT EXISTS advances (
    id TEXT PRIMARY KEY,
    staff_id TEXT NOT NULL REFERENCES staff(id),
    amount REAL NOT NULL CHECK (amount > 0),
    advance_date TEXT NOT NULL,
    reason TEXT NOT NULL DEFAULT '',
    recovered_amount REAL NOT NULL DEFAULT 0,
    remaining_amount REAL NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'recovering', 'recovered', 'rejected')),
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_advances_staff ON advances(staff_id);
CREATE INDEX IF NOT EXISTS idx_advances_status ON advances(status);
CREATE INDEX IF NOT EXISTS idx_advances_date ON advances(advance_date);
