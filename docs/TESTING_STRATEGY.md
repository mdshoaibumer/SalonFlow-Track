# SalonFlow Track — Complete Testing Strategy

## Executive Summary

| Metric | Target | Current | Gap |
|--------|--------|---------|-----|
| Backend Domain Coverage | 90%+ | ~25% | Critical |
| Backend UseCase Coverage | 85%+ | ~15% | Critical |
| Backend Repository Coverage | 80%+ | ~49% | Moderate |
| Backend Handler Coverage | 80%+ | ~4% | Critical |
| Frontend Unit Coverage | 70%+ | 0% | Critical |
| E2E Critical Paths | 100% | 0% | Critical |
| **Overall Target** | **80%+** | **~20%** | **60% gap** |

## Test Pyramid

```
         ╱╲
        ╱ E2E ╲          ← 15-20 critical user journeys (Playwright)
       ╱────────╲
      ╱Integration╲      ← API contract tests, DB integration (Go httptest)
     ╱──────────────╲
    ╱  Unit Tests     ╲   ← Domain, UseCase, Component (Go test, Vitest)
   ╱────────────────────╲
```

## Architecture

```
tests/
├── backend/
│   ├── testutil/           ← Shared factories, DB helpers
│   ├── integration/        ← HTTP-level integration tests
│   └── security/           ← Security-specific tests
├── frontend/
│   ├── __tests__/          ← Vitest component & service tests
│   ├── __mocks__/          ← MSW handlers, test doubles
│   └── e2e/                ← Playwright tests
├── fixtures/               ← Seed data, JSON fixtures
└── ci/                     ← CI/CD pipeline configs
```

---

## 1. Backend Unit Tests (Go)

### 1.1 Domain Layer Tests

Every domain entity gets: validation, business rules, edge cases.

| Entity | Tests Required |
|--------|---------------|
| Staff | Validate, commission calc, status transitions |
| Service | Validate, pricing, duration |
| Customer | Validate, phone format, loyalty |
| Invoice | Validate, line items, totals, tax calc |
| Appointment | Status transitions, overlap detection, time validation |
| Expense | Category validation, amount bounds |
| Product | Stock calc, low stock detection |
| Commission | Percentage calc, tier boundaries |
| Membership | Session tracking, expiry, renewal |
| WhatsApp | Template variable substitution |
| CloudBackup | Config validation, provider constants |
| License | Key format, expiry check |

### 1.2 UseCase Layer Tests

Pattern: Mock repositories via interfaces, test orchestration logic.

```go
// Example: appointment_usecase_test.go
func TestAppointmentUseCase_Create_Success(t *testing.T)
func TestAppointmentUseCase_Create_OverlappingSlot(t *testing.T)
func TestAppointmentUseCase_Create_InvalidCustomer(t *testing.T)
func TestAppointmentUseCase_UpdateStatus_InvalidTransition(t *testing.T)
func TestAppointmentUseCase_Delete_Completed_Fails(t *testing.T)
```

### 1.3 Handler Layer Tests

Pattern: Use `httptest.NewRecorder` + real chi router.

```go
func TestStaffHandler_Create_201(t *testing.T)
func TestStaffHandler_Create_400_InvalidBody(t *testing.T)
func TestStaffHandler_Create_400_MissingFields(t *testing.T)
func TestStaffHandler_Get_200(t *testing.T)
func TestStaffHandler_Get_404_NotFound(t *testing.T)
func TestStaffHandler_List_200_Paginated(t *testing.T)
func TestStaffHandler_Update_200(t *testing.T)
func TestStaffHandler_Delete_200(t *testing.T)
func TestStaffHandler_Delete_404(t *testing.T)
```

---

## 2. Frontend Unit Tests (Vitest)

### 2.1 Service Layer Tests

Mock `fetch` via MSW, test API client contract.

```typescript
describe('staffService', () => {
  it('listStaff returns paginated results')
  it('createStaff sends correct payload')
  it('handles 400 validation errors')
  it('handles network failures gracefully')
  it('handles 500 server errors')
})
```

### 2.2 Hook Tests

Use `@testing-library/react-hooks` with QueryClientProvider.

```typescript
describe('useStaff', () => {
  it('fetches staff on mount')
  it('returns loading state initially')
  it('returns error on API failure')
  it('invalidates cache after mutation')
})
```

### 2.3 Component Tests

```typescript
describe('StaffPage', () => {
  it('renders loading skeleton initially')
  it('renders staff table with data')
  it('opens create dialog on button click')
  it('submits form and shows success toast')
  it('shows error message on validation failure')
  it('handles pagination correctly')
})
```

---

## 3. Integration Tests

### 3.1 API Contract Tests (Go)

Full HTTP roundtrip: Request → Router → Handler → UseCase → Repository → SQLite → Response.

```go
func TestAPI_Staff_CRUD_Integration(t *testing.T) {
    // Setup: in-memory DB with migrations
    // Create → verify 201, response shape
    // List → verify pagination, filtering
    // Update → verify 200, changed fields
    // Delete → verify 200, subsequent 404
}
```

### 3.2 Cross-Module Integration

```go
func TestAPI_Invoice_WithServices_Integration(t *testing.T)
func TestAPI_Appointment_WithStaffAndCustomer(t *testing.T)
func TestAPI_Membership_SellAndUseSession(t *testing.T)
func TestAPI_Commission_CalculatedFromInvoice(t *testing.T)
```

---

## 4. E2E Tests (Playwright)

### 4.1 Critical User Journeys

| # | Journey | Priority |
|---|---------|----------|
| 1 | Staff: Create → Verify in list → Edit → Delete | P0 |
| 2 | Service: Create with pricing → Appears in billing | P0 |
| 3 | Customer: Create → Book appointment → Complete | P0 |
| 4 | Billing: Create invoice → Auto-calculate → Print | P0 |
| 5 | Appointment: Book → Confirm → Complete → History | P0 |
| 6 | Membership: Create plan → Sell → Use session | P1 |
| 7 | Expense: Create → Appears in P&L report | P1 |
| 8 | Backup: Create → Verify in list → Restore | P1 |
| 9 | GST: Configure → Invoice shows tax breakdown | P1 |
| 10 | Reports: Revenue report shows correct totals | P2 |

### 4.2 Negative Scenarios

- Submit empty forms → validation messages appear
- Navigate to non-existent route → 404 page
- API timeout → error boundary shows message
- Submit duplicate phone number → business error shown

---

## 5. Security Tests

### 5.1 Input Validation

```go
func TestSQLInjection_StaffName(t *testing.T)       // '; DROP TABLE--
func TestXSS_CustomerNotes(t *testing.T)            // <script>alert(1)</script>
func TestPathTraversal_BackupRestore(t *testing.T)  // ../../etc/passwd
func TestOversizedPayload_Rejected(t *testing.T)    // 10MB JSON body
func TestInvalidUUID_Returns400(t *testing.T)       // not-a-uuid
```

### 5.2 Business Logic Security

```go
func TestCannotDeleteActiveSubscription(t *testing.T)
func TestCannotOverrideCompletedInvoice(t *testing.T)
func TestLicenseKeyTampering_Detected(t *testing.T)
func TestBackupEncryption_ContentNotReadable(t *testing.T)
```

---

## 6. Performance Tests

### 6.1 Benchmarks (Go)

```go
func BenchmarkInvoiceCreate(b *testing.B)
func BenchmarkStaffList_1000Records(b *testing.B)
func BenchmarkAppointmentFilter_DateRange(b *testing.B)
func BenchmarkAnalytics_DashboardStats(b *testing.B)
```

### 6.2 Load Thresholds

| Operation | Max Latency | Target |
|-----------|------------|--------|
| GET list (50 items) | 200ms | <100ms |
| POST create | 100ms | <50ms |
| GET dashboard stats | 500ms | <200ms |
| Backup create (50MB DB) | 5s | <3s |
| Report generation | 2s | <1s |

---

## 7. Backup & Restore Tests

```go
func TestBackup_CreateFile_ValidSQLite(t *testing.T)
func TestBackup_RestorePreservesData(t *testing.T)
func TestBackup_ConcurrentAccess_NoCorruption(t *testing.T)
func TestCloudBackup_Upload_SimulatedProvider(t *testing.T)
func TestCloudBackup_Download_IntegrityCheck(t *testing.T)
func TestBackup_MaxVersionsRespected(t *testing.T)
func TestBackup_AutoSchedule_Triggers(t *testing.T)
```

---

## 8. License Tests

```go
func TestLicense_ValidKey_Activates(t *testing.T)
func TestLicense_ExpiredKey_Rejected(t *testing.T)
func TestLicense_WrongDevice_Rejected(t *testing.T)
func TestLicense_TamperedSignature_Detected(t *testing.T)
func TestLicense_GracePeriod_AllowsAccess(t *testing.T)
func TestLicense_RevokedKey_BlocksAccess(t *testing.T)
```

---

## 9. Auto Update Tests

```go
func TestUpdate_CheckAvailable_NewVersion(t *testing.T)
func TestUpdate_CheckAvailable_AlreadyLatest(t *testing.T)
func TestUpdate_Download_VerifiesChecksum(t *testing.T)
func TestUpdate_Apply_RollbackOnFailure(t *testing.T)
func TestUpdate_History_Recorded(t *testing.T)
```

---

## 10. Data Migration Tests

```go
func TestMigrations_UpDown_AllVersions(t *testing.T)
func TestMigrations_DataPreserved_AfterUp(t *testing.T)
func TestMigrations_NoDataLoss_OnRollback(t *testing.T)
func TestMigrations_ConcurrentMigration_Safe(t *testing.T)
func TestMigrations_CorruptDB_GracefulError(t *testing.T)
```

---

## Running Tests

### Backend
```powershell
# All tests
$env:CGO_ENABLED="1"; $env:GOARCH="386"; go test ./... -v

# With coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Race detection
go test -race ./...

# Benchmarks
go test -bench=. -benchmem ./...

# Specific module
go test ./internal/adapters/repository/sqlite/ -run TestAppointment -v
```

### Frontend
```powershell
# Unit tests
npm test

# Coverage
npm run test:coverage

# Watch mode
npm run test:watch

# E2E
npx playwright test

# E2E with UI
npx playwright test --ui
```

### CI/CD
```powershell
# Full pipeline locally
npm run ci:test  # Runs lint + type-check + unit + e2e
```

---

## Test Data Strategy

### Factories
- Go: `testutil/factory.go` — generates valid domain objects
- TS: `__tests__/factories.ts` — generates valid API responses

### Seed Data
- `fixtures/seed.sql` — minimal dataset for E2E tests
- `fixtures/large-dataset.sql` — 10K records for performance tests

### Isolation
- Each Go test uses `:memory:` SQLite — zero shared state
- Each Vitest test has fresh QueryClient — no cache leakage
- Each Playwright test resets API state via seed endpoint

---

## Quality Gates (CI/CD)

| Gate | Threshold | Blocks PR |
|------|-----------|-----------|
| Go unit tests | 100% pass | Yes |
| Go coverage | ≥80% | Yes |
| TypeScript type check | 0 errors | Yes |
| Vitest unit tests | 100% pass | Yes |
| Frontend coverage | ≥70% | Yes |
| Playwright E2E | 100% P0 pass | Yes |
| Security scan | 0 critical | Yes |
| Bundle size | <1.5MB gzip | Warning |
| Build time | <60s | Warning |
