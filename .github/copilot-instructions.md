# Copilot Instructions for rosso0815-go-crud-billing

## Project Overview

This is a Go web application for billing/invoice management built with:
- **Backend**: Go 1.25 with PostgreSQL
- **ORM**: SQLC (SQL-first code generation)
- **Templating**: Templ (type-safe HTML templates in Go)
- **Frontend**: HTMX + Bootstrap 5
- **Authentication**: OAuth2 (Google, GitLab, Gitea)
- **UI Framework**: Custom components with Templ

## Build, Test, and Lint Commands

### Setup
```bash
make setup_go     # Install Go tools (sqlc, templ, air, goose, golangci-lint)
make setup_npm    # Install frontend assets (Bootstrap, HTMX)
make sql          # Reset database, run migrations, load seed data
```

### Development
```bash
make run          # Start hot-reloading dev server with Air (port 3000)
```

### Building
```bash
make build        # Build release binary with SQLC & Templ generation
make docker_build # Build Docker image
```

### Code Quality
```bash
make audit        # Run staticcheck (configured with: -checks=all,-ST1000,-U1000)
make lint         # Run golangci-lint, go vet, and compile check
```

### Testing
```bash
go test ./...                    # Run all tests
go test -v ./web -run Test_Show  # Run specific test in web package
go test ./services -race         # Run tests with race detector
```

### Code Generation
```bash
sqlc generate     # Generate type-safe Go from SQL (output: db/generated/)
templ generate    # Generate Go from Templ files
```

## High-Level Architecture

### Directory Structure
- **`main.go`**: Entry point; initializes store, router, and HTTP server
- **`router/`**: HTTP routing with auth handlers (Gitea, GitLab, Google OAuth2)
- **`services/`**: Business logic layer with database stores
  - `customer_store.go`: Customer CRUD operations
  - `invoice_store.go`: Invoice/invoiceentry management
  - `userkv_store.go`: User key-value storage
  - `scheduler.go`: Background tasks (gocron)
- **`web/`**: Web handlers and UI components
  - `*_handler.go`: Endpoint handlers
  - `*entity*.go`: Domain models and utilities
  - `*.templ`: Templ template files
  - `ui/`: Reusable UI components (pagination, forms, tables)
- **`db/`**: Database layer
  - `schema/`: SQL migration files (goose format)
  - `query/`: Named SQL queries for SQLC
  - `generated/`: Auto-generated SQLC code (do not edit)
- **`config/`**: Configuration management
- **`static/`**: Frontend assets (CSS, JS, fonts)

### Data Flow
1. HTTP request → `router.go` → Authentication (OAuth2 if enabled)
2. Handler in `web/*_handler.go` → Calls business logic in `services/*_store.go`
3. Store → Calls auto-generated SQLC methods in `db/generated/*.go`
4. Database query result → Marshal to struct → Return to handler
5. Handler renders Templ template → Response to client

### Key Technologies
- **SQLC**: Generates type-safe Go code from SQL queries in `db/query/`
- **Templ**: Compiles `.templ` files to type-safe Go functions
- **Air**: Hot reload on file changes (configured in `.air.toml`)
- **gocron v2**: Lightweight job scheduling for background tasks

## Key Conventions

### Code Generation
1. **After adding SQL queries**: Run `make build` or `sqlc generate` to regenerate code in `db/generated/`
2. **After modifying `.templ` files**: Run `templ generate` before building
3. Both are included in `make run` via Air's pre-build command

### SQL & SQLC
- Write SQL queries in `db/query/` with `.sql` extension
- Use named queries like `-- name: GetCustomer :one` at the top of each query
- Column type overrides in `sqlc.yaml` map DB columns to Go types (e.g., `customer.customer_id` → `int`)
- Generated code uses `pgx/v5` driver for type safety

### Templ Templates
- Stored in `web/` and subdirectories (e.g., `web/ui/`, `web/customer.templ`)
- Template functions must accept `context.Context` as first parameter
- Auto-generated `*_templ.go` files are excluded from Air's file watcher (see `.air.toml`)
- Reusable components in `web/ui/` (e.g., `crud_table.templ`, `pagination.templ`)

### HTTP Handlers
- Patterns: `func handler(w http.ResponseWriter, r *http.Request)` 
- Access session via session manager from router
- Access config via dependency injection (passed to handler functions)
- Return Templ components or JSON

### Database Transactions
- Services use `db_gen.Queries` for database operations
- Connection pool managed by `services.Store` (wraps `pgxpool.Pool`)
- Store provides methods to create contexts and transactions

### Testing
- Test files follow `*_test.go` naming convention
- Tests for complex logic in `services/` package (e.g., `invoices_pgx_test.go`)
- Use `testing.T` parameter; no testing frameworks like testify in current codebase
- Test data setup typically loads seed data via `db/data/*.sql`

### Environment & Configuration
- Config file: `config.yaml` (embedded in binary)
- Environment variables override config (e.g., `PGHOST`, `PGUSER`, `PGNAME`, `PGPORT`)
- OAuth2 settings configurable via env or config file

### Hot Reload Development
- Air monitors `.go`, `.templ`, `.html`, `.css` files in project root (see `.air.toml`)
- Excludes: `assets/`, `tmp/`, `vendor/`, `db/` directories
- Pre-build command runs `templ generate` before Go build
- Binary runs with `web` argument: `./tmp/main web`

## MCP Server Integration

**PostgreSQL MCP** is configured for this project (connect via your CLI settings). This enables:
- Direct database schema introspection in Copilot
- SQL query validation without manual testing
- Quick table structure lookups during development

Configure in your Copilot CLI settings with your local PostgreSQL connection details (typically `localhost:5432`).
