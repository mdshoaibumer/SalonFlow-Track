-- GST Settings
CREATE TABLE IF NOT EXISTS gst_settings (
    id TEXT PRIMARY KEY,
    business_name TEXT NOT NULL DEFAULT '',
    gstin TEXT NOT NULL DEFAULT '',
    state TEXT NOT NULL DEFAULT '',
    address TEXT NOT NULL DEFAULT '',
    hsn_code TEXT NOT NULL DEFAULT '',
    cgst_rate REAL NOT NULL DEFAULT 9.0,
    sgst_rate REAL NOT NULL DEFAULT 9.0,
    igst_rate REAL NOT NULL DEFAULT 18.0,
    is_gst_enabled INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

-- Tax Rates (per-service/category overrides)
CREATE TABLE IF NOT EXISTS tax_rates (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    hsn_code TEXT NOT NULL DEFAULT '',
    cgst_rate REAL NOT NULL DEFAULT 9.0,
    sgst_rate REAL NOT NULL DEFAULT 9.0,
    igst_rate REAL NOT NULL DEFAULT 18.0,
    category TEXT NOT NULL DEFAULT '',
    is_active INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

-- Invoice Tax Lines
CREATE TABLE IF NOT EXISTS invoice_tax_lines (
    id TEXT PRIMARY KEY,
    invoice_id TEXT NOT NULL REFERENCES invoices(id),
    item_id TEXT NOT NULL DEFAULT '',
    taxable_amount REAL NOT NULL DEFAULT 0,
    cgst_rate REAL NOT NULL DEFAULT 0,
    cgst_amount REAL NOT NULL DEFAULT 0,
    sgst_rate REAL NOT NULL DEFAULT 0,
    sgst_amount REAL NOT NULL DEFAULT 0,
    igst_rate REAL NOT NULL DEFAULT 0,
    igst_amount REAL NOT NULL DEFAULT 0,
    total_tax REAL NOT NULL DEFAULT 0,
    is_interstate INTEGER NOT NULL DEFAULT 0,
    hsn_code TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL
);

CREATE INDEX idx_invoice_tax_lines_invoice ON invoice_tax_lines(invoice_id);
CREATE INDEX idx_tax_rates_category ON tax_rates(category);
