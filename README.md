# seazus-go

A minimal URL shortener API written in Go.

## Overview

`seazus-go` provides a REST API for URL shortening and user management. The project demonstrates:

- Clean package separation for `cmd`, `internal`, and `migrations`
- Structured logging with Zap
- Request middleware for logging, CORS, panic recovery, and rate limiting
- PostgreSQL data access with repository and service layers
- Domain-level error mapping and response handling

## Features

- User registration and login
- URL creation and retrieval
- Protected user routes
- Application-level rate limiting
- Structured JSON API responses

> **NOTE**
> This is still a work in progress and you are welcome to share suggestions on what could be improved on as this is my first shot at building an API with Go.

## Repository layout

- `cmd/api/` — HTTP server entrypoint, routing, and handler wiring
- `internal/config/` — configuration loading
- `internal/db/` — database connection setup
- `internal/logger/` — logging setup and helpers
- `internal/middleware/` — HTTP middleware utilities
- `internal/models/` — domain models and request payloads
- `internal/repository/` — database repositories
- `internal/service/` — business logic services
- `internal/common/` — shared response and error helpers
- `migrations/` — SQL migration files

## Prerequisites

- Go 1.22+
- PostgreSQL
- `goose` for migrations if you run database commands from `Makefile`

## Setup

1. Copy the sample environment file:

   ```bash
   cp env.sample .env
   ```

2. Fill in your database credentials and app settings.

3. Create the database and run migrations:

   ```bash
   make db-up
   ```

## Run

### Development

```bash
make dev
```

### Build and run

```bash
make build
./bin/seazus-api
```

### Direct Go run

```bash
go run ./cmd/api
```

## Testing

```bash
make test
```

## Environment variables

Use `.env` or your own environment provider to configure:

- `DB_USER`
- `DB_PASSWORD`
- `DB_HOST`
- `DB_PORT`
- `DB_NAME`
- `DB_MAX_IDLE_CONNS`
- `DB_MAX_IDLE_TIME`
- `DB_MAX_OPEN_CONNS`
- `APP_ENV`
- `PORT`
- `LOG_LEVEL`
- `TRUSTED_ORIGINS`
- `RATE_LIMIT_REQUEST_PER_TIMEFRAME`
- `RATE_LIMIT_TIMEFRAME`
- `RATE_LIMIT_ENABLED`

## Notes

- The app uses a centralized logger initialization in `internal/logger`.
- Repo and service dependencies are wired in `cmd/api/handlers/handler.go`.
- HTTP error responses are generated via shared helpers in `internal/common/response.go`.
