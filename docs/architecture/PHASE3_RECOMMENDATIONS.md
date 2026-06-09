# Phase 3 Recommendations

## Before Starting Phase 3 (Implementation)

### 1. UUIDv7 Migration

Phase 1 used `google/uuid` (v4). For Phase 2+ entities, switch to **UUIDv7** for:
- Time-ordered keys (better B-tree locality in SQLite)
- No index fragmentation

**Action**: Add `github.com/google/uuid` v1.6+ which supports `uuid.NewV7()` or use a dedicated UUIDv7 library.

### 2. Database Connection Update

Ensure the SQLite connection enables:
```
_journal_mode=WAL         ← Already done
_busy_timeout=5000        ← Already done
_foreign_keys=ON          ← Already done (CRITICAL for cascades)
_cache_size=-64000        ← Increase to 64MB
_synchronous=NORMAL       ← Already done
```

### 3. Transaction Helper

Create a generic transaction wrapper for operations spanning multiple tables:

```go
func (db *DB) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil { return err }
    defer tx.Rollback()
    
    if err := fn(tx); err != nil { return err }
    return tx.Commit()
}
```

### 4. Recommended Implementation Order

```
Phase 3a: Customer + Service CRUD (simplest, enables testing)
Phase 3b: Invoice + Payment (core billing flow)
Phase 3c: Staff Performance queries (derives from billing)
Phase 3d: Incentive rules + calculation
Phase 3e: Advance + Salary generation
Phase 3f: Products + Inventory
Phase 3g: Expenses
Phase 3h: Audit logging middleware
```

### 5. API Design Principles

| Principle | Approach |
|-----------|----------|
| Versioning | `/api/v1/` prefix (already done) |
| Pagination | `?page=1&per_page=20` with total count in meta |
| Filtering | Query parameters: `?status=active&date_from=2026-06-01` |
| Sorting | `?sort_by=created_at&sort_dir=desc` |
| Soft Delete | `DELETE /api/v1/customers/:id` sets deleted_at |
| Bulk Ops | POST with array body where needed |

### 6. Frontend State Architecture

```
TanStack Query Keys Convention:
- ['customers']                  → list
- ['customers', id]              → detail
- ['customers', { filters }]     → filtered list
- ['invoices', 'daily', date]    → daily invoices
- ['staff-performance', staffId, month] → performance
```

### 7. Data Integrity Checks

Before production, implement:
- [ ] Invoice total = Σ(items) - discount (invariant check)
- [ ] Payment total ≤ Invoice total (constraint)
- [ ] Stock never negative (CHECK constraint in DB)
- [ ] Salary unique per staff/month (UNIQUE index done)
- [ ] Advance balance = amount - recovered (application check)

### 8. Performance Considerations

| Query | Expected Load | Optimization |
|-------|--------------|-------------|
| Daily performance | 4 staff × ~10 invoices/day | View is sufficient |
| Monthly report | ~300 invoices aggregated | SQLite handles fine |
| Low stock | ~50 products | Partial index done |
| Invoice search | Up to 10K records/year | Date + status index |

SQLite benchmarks for this scale:
- Simple SELECT: <1ms
- Aggregate query: <5ms
- INSERT with indexes: <2ms
- Full table scan (10K rows): <10ms

No performance concern for single-salon single-branch usage.

### 9. Backup Strategy

Implement SQLite backup using `.backup` command:
- Daily automatic backup to a separate folder
- Keep last 30 days of backups
- Backup before any migration

### 10. Error Recovery

Design for these failure scenarios:
- Power loss during invoice creation → Draft status, incomplete is OK
- App crash during salary generation → Regenerate (idempotent by staff+month)
- Corrupted DB → Restore from backup + replay audit log

### 11. Testing Priority

| Module | Test Priority | Why |
|--------|--------------|-----|
| Salary Calculation | P0 | Money must be accurate |
| Incentive Engine | P0 | Commission errors = trust issues |
| Invoice Totals | P0 | Billing must be correct |
| Stock Adjustment | P1 | Inventory accuracy |
| Advance Recovery | P1 | Deduction correctness |
| Soft Delete | P2 | Data integrity |
| Audit Logging | P2 | Compliance |

### 12. Migration Naming Convention

```
NNN_verb_noun[_qualifier].up.sql
NNN_verb_noun[_qualifier].down.sql

Examples:
006_add_soft_delete_to_phase1.up.sql
009_create_invoices_table.up.sql
018_add_gst_fields_to_invoices.up.sql
```

---

## Architecture Decision: No CQRS Needed

For a single-branch salon with 4 staff and ~30 transactions/day:
- CQRS is overengineering
- Single read/write model through repository pattern is sufficient
- SQLite views handle read-side optimizations
- If performance becomes a concern (unlikely), add materialized tables later

## Architecture Decision: Offline-First Sync (Future)

When cloud sync is needed:
- Use `updated_at` timestamps for conflict detection
- Last-write-wins with manual conflict resolution for invoices
- Append-only audit log serves as event source for replay
- Consider SQLite `ROWID` for ordering guarantees
