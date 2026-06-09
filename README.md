# SalonFlow Track

A production-grade salon management desktop application built with **Go + React** for the Indian market. Ships as a single `.exe` via Wails v2 with an embedded SQLite database — zero external dependencies for end users.

![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&logoColor=white)
![React](https://img.shields.io/badge/React-19-61DAFB?logo=react&logoColor=black)
![TypeScript](https://img.shields.io/badge/TypeScript-5.6-3178C6?logo=typescript&logoColor=white)
![SQLite](https://img.shields.io/badge/SQLite-3-003B57?logo=sqlite&logoColor=white)
![Wails](https://img.shields.io/badge/Wails-v2-EB5E28)
![Tests](https://img.shields.io/badge/E2E_Tests-178_passed-brightgreen)

---

## Features

| Module | Description |
|--------|-------------|
| **Staff Management** | CRUD, designation, commission tracking, status filters |
| **Services & Categories** | Service catalog with pricing, duration, cost tracking |
| **Customer Management** | Customer profiles, visit history, spending analytics |
| **Appointments** | Booking, calendar view, walk-ins, status workflow |
| **Billing & Invoices** | Invoice generation, payment methods (UPI/Cash/Card), tax |
| **Memberships** | Plans, packages, subscriptions, session tracking |
| **Expenses** | Categorized expense tracking, recurring expenses |
| **Salary & Commissions** | Auto-calculation, advances, payroll management |
| **Inventory & Products** | Stock levels, purchases, low-stock alerts |
| **Analytics & Reports** | Revenue, staff performance, customer, P&L reports |
| **WhatsApp Integration** | Templates, automated messages, send history |
| **GST & Tax** | Configurable tax rates, GSTIN management |
| **Backup & Restore** | Local + cloud backup (Google Drive/S3) |
| **Desktop App** | Single .exe, WebView2, auto-updates, license management |

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Frontend | React 19, TypeScript 5.6, Vite 6, Tailwind CSS 3 |
| UI Components | Radix UI + Shadcn/ui, Lucide Icons |
| State | TanStack Query v5 |
| Routing | React Router v7 |
| Forms | React Hook Form + Zod validation |
| Charts | Recharts |
| Backend | Go 1.24+, Chi router, slog logging |
| Database | SQLite 3 (CGo via mattn/go-sqlite3) |
| Desktop | Wails v2 (Go + WebView2) |
| Testing | Vitest, Playwright, Go testing |

---

## Architecture

```
┌─────────────────────────────────────────────────┐
│                 Wails Desktop                     │
│  ┌───────────────┐     ┌─────────────────────┐  │
│  │   WebView2    │────▶│   Go HTTP Server    │  │
│  │  (React SPA)  │     │   (port 8080)       │  │
│  └───────────────┘     └──────────┬──────────┘  │
│                                    │             │
│                         ┌──────────▼──────────┐  │
│                         │   SQLite Database   │  │
│                         └─────────────────────┘  │
└─────────────────────────────────────────────────┘
```

**Backend Design Principles:**
- Clean Architecture (Hexagonal/Ports & Adapters)
- Domain-Driven Design
- Repository Pattern with SQLite
- Dependency Injection (manual, no framework)
- Structured error handling with typed errors

---

## Project Structure

```
SalonFlow-Track/
├── backend/                    # Go backend (module root)
│   ├── main.go                 # Wails entry point
│   ├── app.go                  # App lifecycle (HTTP server start/stop)
│   ├── container.go            # DI composition root + router
│   ├── internal/
│   │   ├── adapters/
│   │   │   ├── handler/        # HTTP handlers (26 files)
│   │   │   ├── repository/     # SQLite repository implementations
│   │   │   ├── backup/         # Local backup adapter
│   │   │   ├── cloudbackup/    # Cloud backup (GDrive/S3)
│   │   │   ├── gst/            # GST calculation engine
│   │   │   ├── license/        # License validation
│   │   │   └── printer/        # Receipt printing
│   │   ├── core/
│   │   │   ├── domain/         # Entities, value objects
│   │   │   ├── ports/          # Interfaces (repository + service)
│   │   │   └── usecase/        # Business logic
│   │   ├── config/             # YAML config loader
│   │   ├── database/           # DB connection + migration runner
│   │   └── logger/             # slog + lumberjack rotation
│   └── pkg/                    # Shared packages
│       ├── apperror/           # Typed application errors
│       ├── middleware/         # HTTP middleware (logging, recovery)
│       └── uid/                # UUID v7 generation
├── frontend/                   # React 19 SPA
│   ├── src/
│   │   ├── app/                # Router, layouts
│   │   ├── components/         # UI components (Shadcn + custom)
│   │   ├── hooks/              # React Query hooks
│   │   ├── pages/              # 34 page components
│   │   ├── services/           # API client layer (23 modules)
│   │   └── types/              # TypeScript definitions
│   ├── e2e/                    # Playwright E2E tests (18 spec files)
│   └── src/mocks/              # MSW handlers for unit tests
├── database/migrations/        # 17 SQL migration pairs (up/down)
├── scripts/                    # PowerShell automation
├── docs/                       # Architecture & strategy docs
├── tests/fixtures/             # Seed data for testing
└── build/                      # Build output (.exe)
```

---

## Getting Started

### Prerequisites

- **Go** 1.24+
- **Node.js** 20+ / npm 10+
- **MinGW-w64 GCC** (for CGo/SQLite compilation)
- **Wails CLI** v2 (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### Quick Setup

```powershell
.\scripts\setup.ps1
```

### Development (Web Mode)

```powershell
# Terminal 1: Backend
cd backend
go run .

# Terminal 2: Frontend (auto-proxies /api to :8080)
cd frontend
npm run dev
```

Frontend: `http://localhost:5173` → API: `http://localhost:8080/api/v1`

### Build Desktop App

```powershell
.\scripts\build-desktop.ps1
# Output: build/bin/SalonFlow-Track.exe (~20 MB)
```

### Run Tests

```powershell
# Frontend unit tests
cd frontend && npm test

# E2E tests (against desktop app)
.\scripts\test-desktop.ps1

# Backend tests
cd backend && go test ./...
```

---

## API Endpoints

All endpoints are prefixed with `/api/v1`:

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check + DB status |
| CRUD | `/staff` | Staff management |
| CRUD | `/services` | Service catalog |
| CRUD | `/customers` | Customer management |
| CRUD | `/appointments` | Appointment booking |
| CRUD | `/invoices` | Invoice & billing |
| CRUD | `/expenses` | Expense tracking |
| CRUD | `/products` | Product inventory |
| CRUD | `/membership-plans` | Membership plans |
| CRUD | `/membership-subscriptions` | Active subscriptions |
| CRUD | `/commissions` | Commission rules |
| CRUD | `/salary` | Salary processing |
| CRUD | `/advances` | Staff advances |
| GET | `/reports/*` | Analytics & reports |
| POST | `/whatsapp/send` | Send WhatsApp message |
| CRUD | `/whatsapp/templates` | Message templates |
| POST | `/backup/create` | Create local backup |
| POST | `/cloud-backup/*` | Cloud backup operations |

---

## Database

- **Engine:** SQLite 3 (embedded, single-file)
- **Migrations:** 17 versioned pairs in `database/migrations/`
- **Schema:** Staff, Services, Customers, Appointments, Invoices, Payments, Expenses, Products, Memberships, Commissions, Salary, Advances, WhatsApp templates, Audit logs, Performance views

Run migrations:
```powershell
.\scripts\migrate.ps1
```

---

## Testing

| Type | Tool | Count | Command |
|------|------|-------|---------|
| Unit (Frontend) | Vitest + MSW | 15+ | `npm test` |
| Unit (Backend) | Go test | 55+ | `go test ./...` |
| E2E (Desktop) | Playwright | 178 | `.\scripts\test-desktop.ps1` |

E2E tests validate the **real desktop binary** — Go backend + SQLite + React frontend running together.

---

## Development Phases

| Phase | Status | Description |
|-------|--------|-------------|
| 1. Foundation | ✅ | Architecture, config, logging, migrations, frontend setup |
| 2. Core CRUD | ✅ | Staff, Services, Customers with full API + UI |
| 3. Business Logic | ✅ | Appointments, Billing, Invoices, Payments |
| 4. Financial Modules | ✅ | Expenses, Salary, Commissions, Advances, P&L |
| 5. Advanced Features | ✅ | Memberships, WhatsApp, Analytics, Inventory |
| 6. Desktop & Ops | ✅ | Wails build, Backup, GST, Printer, License, Updates |
| 7. Testing & QA | ✅ | 178 E2E tests, unit tests, CI/CD pipeline |

---

## License

Proprietary — All rights reserved © 2024-2026 SalonFlow.
