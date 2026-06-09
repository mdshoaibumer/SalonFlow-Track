# Testing Strategy

## Overview

SalonFlow Track employs a multi-layer testing strategy aligned with the testing pyramid.

## Backend Testing (Go)

### Unit Tests
- **Location**: `*_test.go` files alongside source code
- **Scope**: Domain entities, use cases, utility functions
- **Framework**: Standard `testing` package
- **Mocking**: Interfaces (ports) enable easy test doubles

```
backend/
├── internal/
│   ├── core/
│   │   ├── domain/entity/
│   │   │   └── entities_test.go      ← entity validation tests
│   │   └── usecases/
│   │       └── health_test.go         ← use case unit tests
│   ├── adapters/
│   │   └── repository/sqlite/
│   │       └── settings_test.go       ← repository integration tests
│   └── infrastructure/
│       └── config/
│           └── config_test.go         ← config loading tests
└── pkg/
    ├── apperror/
    │   └── errors_test.go
    └── response/
        └── response_test.go
```

### Integration Tests
- **Scope**: Repository layer with real SQLite (in-memory)
- **Pattern**: Use `testing.T` with setup/teardown
- **Database**: `:memory:` SQLite for isolated tests

### Running Backend Tests
```powershell
cd backend
go test ./...                    # All tests
go test ./internal/core/...     # Domain & use case tests only
go test -race ./...             # Race condition detection
go test -cover ./...            # Coverage report
```

## Frontend Testing

### Unit Tests
- **Framework**: Vitest
- **Component Testing**: React Testing Library
- **Location**: `*.test.tsx` alongside components

### Test Structure
```
frontend/src/
├── components/
│   └── layout/
│       └── Sidebar.test.tsx
├── pages/
│   └── DashboardPage.test.tsx
├── services/
│   └── api-client.test.ts
└── lib/
    └── utils.test.ts
```

### Running Frontend Tests
```powershell
cd frontend
npm test                # Run all tests
npm run test:coverage   # With coverage
npm run test:watch      # Watch mode
```

## Test Conventions

### Naming
- Go: `TestFunctionName_Scenario_ExpectedBehavior`
- TS: `describe('Component') → it('should do X when Y')`

### Test Structure (AAA Pattern)
```
Arrange → Set up test data and dependencies
Act     → Execute the function/component under test
Assert  → Verify expected outcomes
```

### What to Test
| Layer | What to Test |
|-------|-------------|
| Domain entities | Validation, business rules |
| Use cases | Orchestration logic, error paths |
| Repositories | CRUD operations, edge cases |
| Handlers | Request parsing, response format |
| Components | Rendering, user interactions |
| Services | API call construction, error handling |

### What NOT to Test
- Third-party libraries (Shadcn, TanStack Query internals)
- Simple getters/setters with no logic
- Configuration constants
- Auto-generated code

## Coverage Targets
- Domain layer: 90%+
- Use cases: 85%+
- Repositories: 80%+
- Frontend components: 70%+
- Overall: 80%+
