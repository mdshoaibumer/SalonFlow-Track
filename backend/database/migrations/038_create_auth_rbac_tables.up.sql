-- Migration: 038_create_auth_rbac_tables
-- Description: Create authentication, RBAC, session, and enhanced audit tables
-- Created: 2026-06-12

-- ============================================================
-- ROLES TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS roles (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    is_system   INTEGER NOT NULL DEFAULT 0,
    created_at  DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at  DATETIME NOT NULL DEFAULT (datetime('now'))
);

-- ============================================================
-- PERMISSIONS TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS permissions (
    id          TEXT PRIMARY KEY,
    code        TEXT NOT NULL UNIQUE,
    module      TEXT NOT NULL,
    action      TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at  DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_permissions_module ON permissions(module);
CREATE INDEX idx_permissions_code ON permissions(code);

-- ============================================================
-- ROLE_PERMISSIONS TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id       TEXT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id TEXT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at    DATETIME NOT NULL DEFAULT (datetime('now')),
    PRIMARY KEY (role_id, permission_id)
);

CREATE INDEX idx_role_permissions_role ON role_permissions(role_id);

-- ============================================================
-- ENHANCED USERS TABLE (add auth fields to existing users)
-- ============================================================
-- Drop the old users table and recreate with auth fields
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id              TEXT PRIMARY KEY,
    username        TEXT NOT NULL UNIQUE,
    email           TEXT UNIQUE,
    password_hash   TEXT NOT NULL,
    display_name    TEXT NOT NULL,
    phone           TEXT DEFAULT '',
    is_active       INTEGER NOT NULL DEFAULT 1,
    is_locked       INTEGER NOT NULL DEFAULT 0,
    failed_attempts INTEGER NOT NULL DEFAULT 0,
    locked_until    DATETIME,
    last_login_at   DATETIME,
    password_changed_at DATETIME NOT NULL DEFAULT (datetime('now')),
    must_change_password INTEGER NOT NULL DEFAULT 0,
    created_at      DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at      DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_is_active ON users(is_active);

-- ============================================================
-- USER_ROLES TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS user_roles (
    user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id    TEXT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    assigned_at DATETIME NOT NULL DEFAULT (datetime('now')),
    assigned_by TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (user_id, role_id)
);

CREATE INDEX idx_user_roles_user ON user_roles(user_id);
CREATE INDEX idx_user_roles_role ON user_roles(role_id);

-- ============================================================
-- SESSIONS TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS sessions (
    id              TEXT PRIMARY KEY,
    user_id         TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash      TEXT NOT NULL UNIQUE,
    device_id       TEXT NOT NULL DEFAULT '',
    ip_address      TEXT NOT NULL DEFAULT '',
    user_agent      TEXT NOT NULL DEFAULT '',
    remember_me     INTEGER NOT NULL DEFAULT 0,
    expires_at      DATETIME NOT NULL,
    last_active_at  DATETIME NOT NULL DEFAULT (datetime('now')),
    created_at      DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_sessions_token ON sessions(token_hash);
CREATE INDEX idx_sessions_expires ON sessions(expires_at);

-- ============================================================
-- ENHANCED AUDIT LOGS TABLE (replace existing)
-- ============================================================
DROP TABLE IF EXISTS audit_logs;

CREATE TABLE audit_logs (
    id              TEXT PRIMARY KEY,
    timestamp       DATETIME NOT NULL DEFAULT (datetime('now')),
    user_id         TEXT NOT NULL DEFAULT '',
    username        TEXT NOT NULL DEFAULT '',
    action          TEXT NOT NULL,
    module          TEXT NOT NULL,
    entity_type     TEXT NOT NULL DEFAULT '',
    entity_id       TEXT NOT NULL DEFAULT '',
    description     TEXT NOT NULL DEFAULT '',
    old_value       TEXT DEFAULT '',
    new_value       TEXT DEFAULT '',
    device_id       TEXT NOT NULL DEFAULT '',
    ip_address      TEXT NOT NULL DEFAULT '',
    user_agent      TEXT NOT NULL DEFAULT '',
    app_version     TEXT NOT NULL DEFAULT '',
    severity        TEXT NOT NULL DEFAULT 'info' CHECK(severity IN ('info', 'warning', 'critical')),
    created_at      DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_audit_logs_timestamp ON audit_logs(timestamp);
CREATE INDEX idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_module ON audit_logs(module);
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_severity ON audit_logs(severity);

-- ============================================================
-- SEED DEFAULT ROLES
-- ============================================================
INSERT INTO roles (id, name, description, is_system) VALUES
    ('role-owner', 'owner', 'Full access to all features', 1),
    ('role-manager', 'manager', 'Manage daily operations and staff', 1),
    ('role-receptionist', 'receptionist', 'Front desk operations', 1),
    ('role-staff', 'staff', 'Basic staff access', 1);

-- ============================================================
-- SEED PERMISSIONS
-- ============================================================
INSERT INTO permissions (id, code, module, action, description) VALUES
    -- Dashboard
    ('perm-dashboard-view', 'dashboard.view', 'dashboard', 'view', 'View dashboard'),
    -- Customers
    ('perm-customers-read', 'customers.read', 'customers', 'read', 'View customers'),
    ('perm-customers-create', 'customers.create', 'customers', 'create', 'Create customers'),
    ('perm-customers-update', 'customers.update', 'customers', 'update', 'Update customers'),
    ('perm-customers-delete', 'customers.delete', 'customers', 'delete', 'Delete customers'),
    -- Staff
    ('perm-staff-read', 'staff.read', 'staff', 'read', 'View staff'),
    ('perm-staff-create', 'staff.create', 'staff', 'create', 'Create staff'),
    ('perm-staff-update', 'staff.update', 'staff', 'update', 'Update staff'),
    ('perm-staff-delete', 'staff.delete', 'staff', 'delete', 'Delete staff'),
    -- Services
    ('perm-services-read', 'services.read', 'services', 'read', 'View services'),
    ('perm-services-create', 'services.create', 'services', 'create', 'Create services'),
    ('perm-services-update', 'services.update', 'services', 'update', 'Update services'),
    ('perm-services-delete', 'services.delete', 'services', 'delete', 'Delete services'),
    -- Billing
    ('perm-billing-view', 'billing.view', 'billing', 'view', 'View billing'),
    ('perm-billing-create', 'billing.create', 'billing', 'create', 'Create invoices'),
    ('perm-billing-void', 'billing.void', 'billing', 'void', 'Void invoices'),
    ('perm-billing-delete', 'billing.delete', 'billing', 'delete', 'Delete invoices'),
    ('perm-billing-payment', 'billing.payment', 'billing', 'payment', 'Record payments'),
    -- Salary
    ('perm-salary-view', 'salary.view', 'salary', 'view', 'View salary'),
    ('perm-salary-generate', 'salary.generate', 'salary', 'generate', 'Generate salary'),
    ('perm-salary-pay', 'salary.pay', 'salary', 'pay', 'Pay salary'),
    -- Advances
    ('perm-advances-view', 'advances.view', 'advances', 'view', 'View advances'),
    ('perm-advances-create', 'advances.create', 'advances', 'create', 'Create advances'),
    -- Commissions
    ('perm-commissions-view', 'commissions.view', 'commissions', 'view', 'View commissions'),
    ('perm-commissions-manage', 'commissions.manage', 'commissions', 'manage', 'Manage commissions'),
    -- Expenses
    ('perm-expenses-view', 'expenses.view', 'expenses', 'view', 'View expenses'),
    ('perm-expenses-create', 'expenses.create', 'expenses', 'create', 'Create expenses'),
    ('perm-expenses-delete', 'expenses.delete', 'expenses', 'delete', 'Delete expenses'),
    -- Inventory
    ('perm-inventory-view', 'inventory.view', 'inventory', 'view', 'View inventory'),
    ('perm-inventory-adjust', 'inventory.adjust', 'inventory', 'adjust', 'Adjust stock'),
    ('perm-inventory-purchase', 'inventory.purchase', 'inventory', 'purchase', 'Create purchases'),
    ('perm-inventory-manage', 'inventory.manage', 'inventory', 'manage', 'Manage products'),
    -- Reports
    ('perm-reports-view', 'reports.view', 'reports', 'view', 'View reports'),
    ('perm-reports-export', 'reports.export', 'reports', 'export', 'Export reports'),
    -- Analytics
    ('perm-analytics-view', 'analytics.view', 'analytics', 'view', 'View analytics'),
    -- Performance
    ('perm-performance-view', 'performance.view', 'performance', 'view', 'View performance'),
    -- GST
    ('perm-gst-view', 'gst.view', 'gst', 'view', 'View GST'),
    ('perm-gst-manage', 'gst.manage', 'gst', 'manage', 'Manage GST settings'),
    -- Printer
    ('perm-printer-use', 'printer.use', 'printer', 'use', 'Use printer'),
    ('perm-printer-manage', 'printer.manage', 'printer', 'manage', 'Manage printer settings'),
    -- Appointments
    ('perm-appointments-view', 'appointments.view', 'appointments', 'view', 'View appointments'),
    ('perm-appointments-create', 'appointments.create', 'appointments', 'create', 'Create appointments'),
    ('perm-appointments-update', 'appointments.update', 'appointments', 'update', 'Update appointments'),
    ('perm-appointments-delete', 'appointments.delete', 'appointments', 'delete', 'Delete appointments'),
    -- Memberships
    ('perm-memberships-view', 'memberships.view', 'memberships', 'view', 'View memberships'),
    ('perm-memberships-manage', 'memberships.manage', 'memberships', 'manage', 'Manage memberships'),
    -- WhatsApp
    ('perm-whatsapp-view', 'whatsapp.view', 'whatsapp', 'view', 'View WhatsApp'),
    ('perm-whatsapp-send', 'whatsapp.send', 'whatsapp', 'send', 'Send WhatsApp messages'),
    -- Backup
    ('perm-backup-create', 'backup.create', 'backup', 'create', 'Create backups'),
    ('perm-backup-restore', 'backup.restore', 'backup', 'restore', 'Restore backups'),
    -- Import
    ('perm-import-execute', 'import.execute', 'import', 'execute', 'Execute imports'),
    -- License
    ('perm-license-view', 'license.view', 'license', 'view', 'View license'),
    ('perm-license-manage', 'license.manage', 'license', 'manage', 'Manage license'),
    -- Updates
    ('perm-updates-view', 'updates.view', 'updates', 'view', 'View updates'),
    ('perm-updates-install', 'updates.install', 'updates', 'install', 'Install updates'),
    -- Users
    ('perm-users-view', 'users.view', 'users', 'view', 'View users'),
    ('perm-users-manage', 'users.manage', 'users', 'manage', 'Manage users'),
    ('perm-users-roles', 'users.roles', 'users', 'roles', 'Manage user roles'),
    -- Audit
    ('perm-audit-view', 'audit.view', 'audit', 'view', 'View audit logs'),
    -- Diagnostics
    ('perm-diagnostics-view', 'diagnostics.view', 'diagnostics', 'view', 'View diagnostics'),
    ('perm-diagnostics-export', 'diagnostics.export', 'diagnostics', 'export', 'Export diagnostics'),
    -- Settings
    ('perm-settings-view', 'settings.view', 'settings', 'view', 'View settings'),
    ('perm-settings-manage', 'settings.manage', 'settings', 'manage', 'Manage settings'),
    -- Profit & Loss
    ('perm-profitloss-view', 'profitloss.view', 'profitloss', 'view', 'View profit & loss');

-- ============================================================
-- ASSIGN PERMISSIONS TO ROLES
-- ============================================================

-- Owner gets ALL permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role-owner', id FROM permissions;

-- Manager gets most permissions except user management, license, updates, diagnostics export
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role-manager', id FROM permissions
WHERE code NOT IN ('users.manage', 'users.roles', 'license.manage', 'updates.install', 'diagnostics.export');

-- Receptionist gets front-desk permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role-receptionist', id FROM permissions
WHERE code IN (
    'dashboard.view',
    'customers.read', 'customers.create', 'customers.update',
    'billing.view', 'billing.create', 'billing.payment',
    'services.read',
    'staff.read',
    'appointments.view', 'appointments.create', 'appointments.update',
    'memberships.view', 'memberships.manage',
    'whatsapp.view', 'whatsapp.send',
    'printer.use',
    'inventory.view'
);

-- Staff gets minimal permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role-staff', id FROM permissions
WHERE code IN (
    'dashboard.view',
    'customers.read',
    'services.read',
    'appointments.view',
    'performance.view'
);
