# SalonFlow Track - Development Roadmap

## Phase 1: Foundation (CURRENT) ✅

**Goal**: Rock-solid architecture, zero business features.

- [x] Project structure & folder conventions
- [x] Backend Clean Architecture (Go)
- [x] Domain entities & ports
- [x] Repository pattern with SQLite
- [x] Dependency injection container
- [x] Configuration system (YAML, env-aware)
- [x] Structured logging (slog + rotation)
- [x] Error handling framework
- [x] Database migration system
- [x] Foundation tables (users, staff, services, settings, license)
- [x] HTTP server with graceful shutdown
- [x] Health check endpoint
- [x] Frontend architecture (React 19 + TypeScript + Vite)
- [x] Routing with React Router
- [x] Layout system (Sidebar + Header)
- [x] Theme support (light/dark/system)
- [x] Tailwind + Shadcn UI setup
- [x] API client service layer
- [x] TanStack Query integration
- [x] TypeScript type definitions
- [x] Development scripts (setup, dev, migrate)
- [x] Testing strategy documentation

---

## Phase 2: Core CRUD Modules

**Goal**: Staff, Services, and Settings management.

### Staff Module
- [ ] Staff CRUD use cases
- [ ] Staff HTTP handlers (GET, POST, PUT, DELETE)
- [ ] Staff repository implementation
- [ ] Staff list page with table
- [ ] Staff add/edit form with validation
- [ ] Staff detail view

### Services Module
- [ ] Service CRUD use cases
- [ ] Service HTTP handlers
- [ ] Service repository implementation
- [ ] Service list page with categories
- [ ] Service add/edit form
- [ ] Category management

### Settings Module
- [ ] Settings CRUD handlers
- [ ] Salon profile settings (name, address, phone)
- [ ] Operating hours configuration
- [ ] Settings page with sections

### Technical
- [ ] Request validation middleware
- [ ] Pagination support
- [ ] Search/filter utilities
- [ ] Form validation (frontend - Zod)
- [ ] Toast notifications
- [ ] Loading states & skeletons

---

## Phase 3: Appointments & Calendar

**Goal**: Core booking functionality.

- [ ] Appointment entity & domain logic
- [ ] Appointment CRUD
- [ ] Calendar view (daily/weekly)
- [ ] Time slot management
- [ ] Staff availability
- [ ] Walk-in vs scheduled appointments
- [ ] Appointment status flow (booked → in-progress → completed)
- [ ] Conflict detection (double booking)

---

## Phase 4: Billing & Invoicing

**Goal**: Invoice generation and payment tracking.

- [ ] Invoice entity & domain logic
- [ ] Invoice generation from appointments
- [ ] GST calculation engine
- [ ] Invoice PDF generation
- [ ] Payment recording
- [ ] UPI integration (architecture)
- [ ] Daily cash register
- [ ] Payment methods (cash, card, UPI)

---

## Phase 5: Reporting & Analytics

**Goal**: Business intelligence.

- [ ] Revenue reports (daily, weekly, monthly)
- [ ] Staff performance metrics
- [ ] Service popularity analytics
- [ ] Customer visit frequency
- [ ] Export to Excel/PDF
- [ ] Dashboard widgets with real data

---

## Phase 6: Licensing & Updates

**Goal**: Commercial readiness.

- [ ] License key validation logic
- [ ] License activation flow
- [ ] Expiry handling & grace period
- [ ] Auto-update mechanism
- [ ] Update download & installation
- [ ] Version management

---

## Phase 7: Tauri Desktop Integration

**Goal**: Package as native desktop app.

- [ ] Tauri v2 configuration
- [ ] IPC command bridge (replace HTTP in prod)
- [ ] System tray integration
- [ ] Native file dialogs
- [ ] Auto-start on boot option
- [ ] Windows installer (MSI/NSIS)
- [ ] Go sidecar bundling

---

## Phase 8: Cloud Sync (Future)

**Goal**: Optional cloud connectivity.

- [ ] Sync architecture design
- [ ] Conflict resolution strategy
- [ ] Offline-first queue
- [ ] Multi-device support
- [ ] Cloud backup
- [ ] Remote access API

---

## Architecture Decisions Record (ADR)

| Decision | Choice | Reason |
|----------|--------|--------|
| Database | SQLite | Offline-first, zero config, portable |
| Backend | Go | Fast, single binary, excellent concurrency |
| Frontend | React 19 | Ecosystem, developer experience |
| Desktop | Tauri v2 | Small binary, native performance |
| State | TanStack Query | Server state management, caching |
| Styling | Tailwind + Shadcn | Rapid UI development, consistent design |
| Architecture | Clean Architecture | Testability, maintainability, flexibility |
| Logging | slog | Standard library, structured, performant |
| Config | YAML | Human-readable, hierarchical |
| UUID | v4 | No sequential exposure, globally unique |
