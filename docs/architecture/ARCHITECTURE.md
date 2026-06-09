# SalonFlow Track - Architecture Document

## Overview

SalonFlow Track is a production-grade salon management desktop application built with:
- **Frontend**: React 19 + TypeScript + Vite + Tailwind CSS + Shadcn UI
- **Desktop**: Tauri v2
- **Backend**: Go 1.24+
- **Database**: SQLite

## Architecture Principles

- Clean Architecture (Hexagonal)
- Domain Driven Design
- Repository Pattern
- Dependency Injection
- Service Layer
- SOLID Principles
- Offline-first design

## System Architecture

```
┌─────────────────────────────────────────────────────┐
│                    Tauri v2 Shell                     │
├─────────────────────────────────────────────────────┤
│                                                       │
│  ┌─────────────────┐       ┌─────────────────────┐  │
│  │   React 19 UI   │◄─────►│   Go Backend (API)  │  │
│  │   (WebView)     │  IPC  │   (Sidecar/Embed)   │  │
│  └─────────────────┘       └──────────┬──────────┘  │
│                                        │              │
│                              ┌─────────▼─────────┐   │
│                              │     SQLite DB      │   │
│                              └───────────────────┘   │
│                                                       │
└─────────────────────────────────────────────────────┘
```

## Backend Layer Architecture

```
┌───────────────────────────────────────────────────┐
│                  Interfaces Layer                    │
│          (HTTP Handlers, CLI, IPC)                   │
├───────────────────────────────────────────────────┤
│                Application Layer                     │
│          (Use Cases, DTOs, Services)                 │
├───────────────────────────────────────────────────┤
│                  Domain Layer                        │
│      (Entities, Value Objects, Domain Events)        │
├───────────────────────────────────────────────────┤
│              Infrastructure Layer                    │
│    (Database, Logger, Config, External Services)     │
└───────────────────────────────────────────────────┘
```

## Dependency Rule

Dependencies point INWARD only:
- Interfaces → Application → Domain
- Infrastructure → Domain (implements ports)
- Domain depends on NOTHING external

## Frontend Architecture

```
src/
├── app/          → Application shell, providers, router
├── pages/        → Route-level page components
├── features/     → Feature modules (colocated logic)
├── components/   → Shared UI components
├── hooks/        → Custom React hooks
├── services/     → API client, IPC bridge
├── lib/          → Utilities, helpers
└── types/        → Shared TypeScript types
```

## Data Flow

1. User interacts with React UI
2. UI calls service layer (TanStack Query)
3. Service sends request via Tauri IPC / HTTP
4. Go handler receives request
5. Handler calls use case
6. Use case orchestrates domain logic
7. Repository persists to SQLite
8. Response flows back up the chain

## Communication Pattern

- **Development**: HTTP REST (Go serves on localhost)
- **Production**: Tauri IPC commands (Go embedded as sidecar)
