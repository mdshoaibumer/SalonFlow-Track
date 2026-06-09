-- Printer Settings
CREATE TABLE IF NOT EXISTS printer_settings (
    id TEXT PRIMARY KEY,
    default_printer TEXT NOT NULL DEFAULT '',
    paper_width TEXT NOT NULL DEFAULT '80mm' CHECK (paper_width IN ('58mm','80mm','A4')),
    margin_top INTEGER NOT NULL DEFAULT 5,
    margin_bottom INTEGER NOT NULL DEFAULT 5,
    margin_left INTEGER NOT NULL DEFAULT 5,
    margin_right INTEGER NOT NULL DEFAULT 5,
    header_text TEXT NOT NULL DEFAULT '',
    footer_text TEXT NOT NULL DEFAULT 'Thank you for visiting!',
    show_logo INTEGER NOT NULL DEFAULT 0,
    show_qr INTEGER NOT NULL DEFAULT 0,
    upi_id TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

-- Print Jobs (history)
CREATE TABLE IF NOT EXISTS print_jobs (
    id TEXT PRIMARY KEY,
    document_type TEXT NOT NULL CHECK (document_type IN ('invoice','receipt','expense','salary_slip','purchase')),
    document_id TEXT NOT NULL DEFAULT '',
    printer_name TEXT NOT NULL DEFAULT '',
    paper_width TEXT NOT NULL DEFAULT '80mm',
    status TEXT NOT NULL DEFAULT 'queued' CHECK (status IN ('queued','printing','completed','failed')),
    copies INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL
);

CREATE INDEX idx_print_jobs_document ON print_jobs(document_type, document_id);
CREATE INDEX idx_print_jobs_created ON print_jobs(created_at);
