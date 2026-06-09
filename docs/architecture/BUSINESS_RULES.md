# Business Rules

## 1. Invoice Rules

| Rule ID | Rule | Enforcement |
|---------|------|-------------|
| INV-001 | Invoice must have at least one item before completion | Domain entity validation |
| INV-002 | Invoice number is auto-generated: `INV-YYYYMM-XXXX` | Application layer |
| INV-003 | Cannot cancel a completed invoice with payments | Domain entity logic |
| INV-004 | Each invoice item must be assigned to a staff member | Database NOT NULL |
| INV-005 | Invoice total = Σ(item totals) - discount + tax | Domain recalculation |
| INV-006 | Payment status auto-updates: unpaid → partial → paid | Domain entity logic |
| INV-007 | Cancelled invoices do NOT count toward staff performance | Query filter |
| INV-008 | Item name is denormalized (snapshot at time of billing) | Application layer |
| INV-009 | Invoice date defaults to today but can be backdated | UI/API validation |
| INV-010 | Discount cannot exceed subtotal | Domain validation |

## 2. Payment Rules

| Rule ID | Rule | Enforcement |
|---------|------|-------------|
| PAY-001 | Payment amount must be positive | Domain validation |
| PAY-002 | Total payments cannot exceed invoice total | Application layer |
| PAY-003 | Payment method must be: cash, card, or upi | Domain validation |
| PAY-004 | Refund is recorded as separate payment (is_refund=true) | Domain flag |
| PAY-005 | UPI/card payments should store reference number | UI suggestion (not mandatory) |

## 3. Staff Rules

| Rule ID | Rule | Enforcement |
|---------|------|-------------|
| STF-001 | Staff name is required | Domain validation |
| STF-002 | Staff must have a role: stylist, assistant, receptionist | Domain validation |
| STF-003 | Inactive staff cannot be assigned new invoice items | Application layer |
| STF-004 | Deleting staff is soft-delete only | Repository layer |
| STF-005 | Base salary is stored on staff record for salary generation | Schema design |
| STF-006 | Staff phone should be unique (warn, not block) | Application layer |

## 4. Customer Rules

| Rule ID | Rule | Enforcement |
|---------|------|-------------|
| CUS-001 | Customer name is required | Domain validation |
| CUS-002 | Customer phone is required | Domain validation |
| CUS-003 | Phone should be unique (deduplicate walk-ins) | Database UNIQUE |
| CUS-004 | Visit count auto-increments on invoice completion | Application layer |
| CUS-005 | Gender: male, female, other | Domain validation |
| CUS-006 | Walk-in can be recorded with minimal data (name + phone) | Business requirement |

## 5. Service Rules

| Rule ID | Rule | Enforcement |
|---------|------|-------------|
| SVC-001 | Service name is required | Domain validation |
| SVC-002 | Price cannot be negative | Domain validation |
| SVC-003 | Duration must be positive (in minutes) | Domain validation |
| SVC-004 | Service belongs to exactly one category | Schema FK |
| SVC-005 | Inactive services can still appear in historical invoices | Denormalized name |
| SVC-006 | Service price can change without affecting past invoices | Snapshot pattern |

## 6. Incentive Rules

| Rule ID | Rule | Enforcement |
|---------|------|-------------|
| INC-001 | Revenue slab rules require at least one slab | Application validation |
| INC-002 | Slabs must not overlap (min/max ranges) | Application validation |
| INC-003 | Service-based rules require a service_id | Domain validation |
| INC-004 | If staff_id is NULL, rule applies to ALL staff | Query logic |
| INC-005 | Only active rules in effective date range are evaluated | Query filter |
| INC-006 | Commission calculated on completed invoice items only | Query filter |
| INC-007 | Monthly period: 1st to last day of month | Application layer |
| INC-008 | Incentive must be approved before inclusion in salary | Status workflow |

### Revenue Slab Calculation Example

```
Rule: Monthly Revenue Commission
Staff: Nazim

Slabs:
  0 - 10000    → 0%
  10001 - 20000 → 5%
  20001 - ∞     → 10%

If Nazim's monthly revenue = ₹25,000
  → Falls in slab 3 (20001-∞)
  → Commission = ₹25,000 × 10% = ₹2,500
```

### Service-Based Calculation Example

```
Rule: Hair Color Commission (15%)
Service: Hair Color (ID: xxx)

If Nazim did 5 Hair Color services @ ₹2000 each
  → Revenue from Hair Color = ₹10,000
  → Commission = ₹10,000 × 15% = ₹1,500
```

## 7. Salary Rules

| Rule ID | Rule | Enforcement |
|---------|------|-------------|
| SAL-001 | One salary record per staff per month | Database UNIQUE(staff_id, month, year) |
| SAL-002 | Net = Base + Commissions + Bonus - Advance Deductions - Other Deductions | Domain calculation |
| SAL-003 | Salary flows: Draft → Approved → Paid | Domain state machine |
| SAL-004 | Only draft salaries can be modified | Domain validation |
| SAL-005 | Advance deductions auto-update advance balance_amount | Application layer |
| SAL-006 | Paid salary records the payment method | Domain field |
| SAL-007 | Salary generation pulls: base from staff, incentives from staff_incentives, advances from advances | Application orchestration |

### Salary Calculation Flow

```
1. Get staff.base_salary → ₹15,000
2. Get approved staff_incentives for the month
   - Revenue commission: ₹2,500
   - Service commission: ₹1,500
   → Total commission: ₹4,000
3. Get pending advances
   - Advance taken on 5th: ₹3,000
4. Calculate:
   Base Pay:      +₹15,000
   Commission:    +₹4,000
   Advance:       -₹3,000
   ─────────────────────────
   Net Salary:     ₹16,000
```

## 8. Advance Rules

| Rule ID | Rule | Enforcement |
|---------|------|-------------|
| ADV-001 | Amount must be positive | Domain validation |
| ADV-002 | Reason is required | Domain validation |
| ADV-003 | Recovery amount cannot exceed balance | Domain validation |
| ADV-004 | Status auto-updates: pending → partial → recovered | Domain logic |
| ADV-005 | Pending/partial advances are auto-included in next salary | Application layer |
| ADV-006 | Partial recovery allowed (can deduct less than full amount) | Business rule |

## 9. Expense Rules

| Rule ID | Rule | Enforcement |
|---------|------|-------------|
| EXP-001 | Amount must be positive | Domain validation |
| EXP-002 | Category is required | Domain validation |
| EXP-003 | Description is required | Domain validation |
| EXP-004 | Default categories: Rent, Electricity, Supplies, Marketing, Maintenance, Miscellaneous | Seed data |
| EXP-005 | Custom categories can be added | CRUD |
| EXP-006 | System categories (is_default=true) cannot be deleted | Application layer |

## 10. Inventory Rules

| Rule ID | Rule | Enforcement |
|---------|------|-------------|
| INV-001 | Stock cannot go negative | Domain validation |
| INV-002 | Every stock change creates a StockTransaction | Application layer |
| INV-003 | Purchase increases stock, consumption/sale decreases | Domain logic |
| INV-004 | Low stock = current_stock ≤ min_stock_level | Domain method |
| INV-005 | Product sold to customer creates invoice item + stock out | Application orchestration |
| INV-006 | Inventory valuation = Σ(current_stock × cost_price) | Query aggregation |
| INV-007 | SKU should be unique if provided | Database constraint |

## 11. Staff Performance (DSR) Rules

| Rule ID | Rule | Enforcement |
|---------|------|-------------|
| DSR-001 | Revenue = Σ(invoice_items.total_price) WHERE staff_id AND completed | Query |
| DSR-002 | Only completed invoices count | Filter: status='completed' |
| DSR-003 | Customer count = distinct customers served that day | Query DISTINCT |
| DSR-004 | Cancelled invoices excluded from all calculations | Filter |
| DSR-005 | Date range aggregations: daily, weekly (Mon-Sun), monthly | Query logic |
| DSR-006 | Product sales revenue counted separately from service revenue | item_type filter |

### DSR Query Pattern

```sql
-- Daily Staff Performance
SELECT 
    ii.staff_id,
    s.name AS staff_name,
    DATE(i.date) AS work_date,
    COUNT(DISTINCT i.customer_id) AS customer_count,
    SUM(CASE WHEN ii.item_type = 'service' THEN ii.total_price ELSE 0 END) AS service_revenue,
    SUM(CASE WHEN ii.item_type = 'product' THEN ii.total_price ELSE 0 END) AS product_revenue,
    SUM(ii.total_price) AS total_revenue,
    COUNT(ii.id) AS item_count
FROM invoice_items ii
JOIN invoices i ON i.id = ii.invoice_id
JOIN staff s ON s.id = ii.staff_id
WHERE i.status = 'completed'
  AND i.deleted_at IS NULL
  AND DATE(i.date) = ?
GROUP BY ii.staff_id, s.name, DATE(i.date)
ORDER BY total_revenue DESC;
```

## 12. Audit Rules

| Rule ID | Rule | Enforcement |
|---------|------|-------------|
| AUD-001 | Invoice create/update/cancel must be audited | Application layer |
| AUD-002 | Salary approve/pay must be audited | Application layer |
| AUD-003 | Expense create/delete must be audited | Application layer |
| AUD-004 | Advance create must be audited | Application layer |
| AUD-005 | Old/new values stored as JSON for diff capability | Schema design |
| AUD-006 | Audit logs are append-only (never update/delete) | No UPDATE/DELETE |
| AUD-007 | Audit log retention: 2 years minimum | Housekeeping policy |

## 13. Soft Delete Rules

| Rule | Description |
|------|-------------|
| All customer-facing entities support soft delete | `deleted_at` column |
| Soft-deleted records excluded from list queries | WHERE deleted_at IS NULL |
| Soft-deleted records still visible in historical reports | Historical queries include deleted |
| Audit logs are NEVER soft-deleted | Immutable record |
| Settings are NEVER soft-deleted | Hard delete only |
