# Migration Plan

## Complete Migration Sequence

| # | File | Description | Phase |
|---|------|-------------|-------|
| 001 | `001_create_users_table` | Application operators/users | Phase 1 |
| 002 | `002_create_staff_table` | Salon staff members | Phase 1 |
| 003 | `003_create_services_table` | Service catalog | Phase 1 |
| 004 | `004_create_settings_table` | Key-value settings + seeds | Phase 1 |
| 005 | `005_create_license_table` | License management | Phase 1 |
| 006 | `006_add_soft_delete_to_phase1` | Add deleted_at, base_salary to Phase 1 tables | Phase 2 |
| 007 | `007_create_customers_table` | Customer/client management | Phase 2 |
| 008 | `008_create_service_categories` | Service category grouping + FK | Phase 2 |
| 009 | `009_create_invoices_table` | Invoices + invoice_items | Phase 2 |
| 010 | `010_create_payments_table` | Payment records | Phase 2 |
| 011 | `011_create_expense_tables` | Expense categories + expenses | Phase 2 |
| 012 | `012_create_advances_table` | Staff salary advances | Phase 2 |
| 013 | `013_create_incentive_tables` | Rules, slabs, staff_incentives | Phase 2 |
| 014 | `014_create_salary_tables` | Salaries + line items | Phase 2 |
| 015 | `015_create_products_table` | Products + stock_transactions | Phase 2 |
| 016 | `016_create_audit_logs_table` | Audit logging | Phase 2 |
| 017 | `017_create_performance_views` | Reporting views (DSR, monthly, low stock) | Phase 2 |

## Naming Convention

```
NNN_verb_noun[_qualifier].{up|down}.sql

NNN         → 3-digit sequential number (001, 002, ...)
verb        → create, add, alter, drop, seed, create
noun        → table name or logical group
qualifier   → optional context
up/down     → forward migration / rollback
```

## Versioning Strategy

- Migrations are **sequential and immutable** once applied
- Never modify an applied migration — create a new one instead
- Each migration runs in a single transaction
- The `schema_migrations` table tracks applied versions
- Rollbacks (`*.down.sql`) exist for development; production rolls forward only

## Dependency Order

```
001 users          (independent)
002 staff          (independent)
003 services       (independent)
004 settings       (independent)
005 license        (independent)
006 soft_delete    (depends on 001, 002, 003)
007 customers      (independent)
008 categories     (depends on 003 - adds FK to services)
009 invoices       (depends on 007 customers, 003 services, 002 staff)
010 payments       (depends on 009 invoices)
011 expenses       (independent)
012 advances       (depends on 002 staff)
013 incentives     (depends on 002 staff, 003 services, 008 categories)
014 salaries       (depends on 002 staff)
015 products       (independent)
016 audit_logs     (independent)
017 views          (depends on 009, 010, 002, 012, 015)
```

## Future Migrations (Phase 3+)

```
018_add_gst_fields_to_invoices.up.sql      → GST amount, GSTIN, HSN codes
019_create_appointments_table.up.sql        → Booking system
020_add_upi_details_to_settings.up.sql      → UPI configuration
021_create_backup_history_table.up.sql      → Backup tracking
```
