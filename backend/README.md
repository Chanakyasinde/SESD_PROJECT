# Inventory Management System Backend (Go + Gin + MongoDB)

## Run

1. Copy `.env.example` to `.env`
2. Install dependencies and run:

```bash
go mod tidy
go run ./cmd/api
```

## API Base URL

`/api/v1`

## Core Modules

- Domain: entities and interface contracts
- Usecase: business rules independent from frameworks
- Repository: MongoDB implementation of interfaces
- Handler: Gin HTTP adapters
- Middleware: auth, role checks, logging, recovery
- Infrastructure: config, database wiring, security
