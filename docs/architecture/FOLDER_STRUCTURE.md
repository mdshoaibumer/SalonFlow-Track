# Folder Structure

```
salonflow-track/
├── frontend/                    # React 19 + TypeScript frontend
│   ├── public/                  # Static assets
│   ├── src/
│   │   ├── app/                 # Application shell & providers
│   │   │   ├── providers/       # Context providers (theme, auth, query)
│   │   │   ├── router/          # Route definitions & guards
│   │   │   └── layouts/         # Page layout components
│   │   ├── pages/               # Route-level page components
│   │   ├── features/            # Feature-specific modules
│   │   ├── components/          # Shared reusable UI components
│   │   │   └── ui/              # Shadcn UI components
│   │   ├── hooks/               # Custom React hooks
│   │   ├── services/            # API layer & IPC bridge
│   │   ├── lib/                 # Utilities, constants, helpers
│   │   └── types/               # Shared TypeScript type definitions
│   ├── index.html
│   ├── package.json
│   ├── tsconfig.json
│   ├── tailwind.config.ts
│   ├── vite.config.ts
│   └── components.json          # Shadcn UI config
│
├── backend/                     # Go backend service
│   ├── cmd/
│   │   └── server/              # Application entry point
│   │       └── main.go
│   ├── internal/
│   │   ├── core/
│   │   │   ├── domain/          # Entities, value objects, domain errors
│   │   │   │   ├── entity/      # Domain entities
│   │   │   │   ├── valueobject/ # Value objects
│   │   │   │   └── event/       # Domain events
│   │   │   ├── ports/           # Interfaces (repository, service contracts)
│   │   │   │   ├── repository/  # Repository interfaces
│   │   │   │   └── service/     # Service interfaces
│   │   │   └── usecases/        # Application use cases
│   │   ├── adapters/
│   │   │   ├── repository/      # Repository implementations (SQLite)
│   │   │   └── handlers/        # HTTP/IPC handlers
│   │   │       └── http/        # REST API handlers
│   │   └── infrastructure/
│   │       ├── database/        # DB connection, migrations
│   │       ├── logger/          # Structured logging (slog)
│   │       ├── config/          # Configuration loading
│   │       └── server/          # HTTP server setup
│   ├── pkg/                     # Shared packages (errors, middleware)
│   │   ├── apperror/            # Application error types
│   │   ├── middleware/          # HTTP middleware
│   │   ├── response/            # Standard API responses
│   │   └── validator/           # Request validation
│   ├── go.mod
│   ├── go.sum
│   └── config.yaml              # Application configuration
│
├── database/                    # Database assets
│   └── migrations/              # SQL migration files
│
├── docs/                        # Documentation
│   ├── architecture/            # Architecture decisions & diagrams
│   └── api/                     # API documentation
│
├── scripts/                     # Development & build scripts
│   ├── setup.ps1                # Windows setup script
│   ├── migrate.ps1              # Run migrations
│   └── dev.ps1                  # Start dev environment
│
├── src-tauri/                   # Tauri v2 configuration (Phase 2)
│
└── README.md
```

## Folder Responsibilities

### Frontend

| Folder | Purpose |
|--------|---------|
| `app/providers` | React context providers for global state (theme, query client) |
| `app/router` | Route configuration, lazy loading, route guards |
| `app/layouts` | Reusable page layouts (sidebar + header + content) |
| `pages/` | One component per route, thin orchestration layer |
| `features/` | Feature-specific components, hooks, and logic colocated |
| `components/` | Shared atomic/molecule components |
| `components/ui/` | Shadcn UI generated components |
| `hooks/` | Reusable custom hooks (useLocalStorage, useDebounce) |
| `services/` | HTTP client, Tauri IPC bridge, API functions |
| `lib/` | Pure utilities (formatters, constants, cn helper) |
| `types/` | Shared TypeScript interfaces and type aliases |

### Backend

| Folder | Purpose |
|--------|---------|
| `cmd/server/` | Application entry point, bootstrap |
| `core/domain/entity/` | Business entities with behavior |
| `core/domain/valueobject/` | Immutable value objects |
| `core/ports/repository/` | Repository interface definitions |
| `core/ports/service/` | Service interface definitions |
| `core/usecases/` | Application business logic orchestration |
| `adapters/repository/` | SQLite repository implementations |
| `adapters/handlers/http/` | HTTP request handlers |
| `infrastructure/database/` | SQLite connection pool, migration runner |
| `infrastructure/logger/` | slog configuration, log rotation |
| `infrastructure/config/` | YAML config loading & validation |
| `infrastructure/server/` | HTTP server lifecycle management |
| `pkg/apperror/` | Standardized error types |
| `pkg/middleware/` | HTTP middleware (logging, recovery, CORS) |
| `pkg/response/` | Standard JSON response helpers |
| `pkg/validator/` | Input validation utilities |
