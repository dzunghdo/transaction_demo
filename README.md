# Transaction Demo

A Go-based transaction management demonstration application built with Gin web framework, GORM ORM, and PostgreSQL database. This project demonstrates transaction handling patterns using the go-transaction-manager library.

## Prerequisites

Before setting up the project, ensure you have the following installed:

- **Go**: Version 1.24 or higher ([Download Go](https://golang.org/downloads/))
- **Docker**: For running PostgreSQL database ([Install Docker](https://docs.docker.com/get-docker/))
- **Docker Compose**: For orchestrating services ([Install Docker Compose](https://docs.docker.com/compose/install/))
- **Make**: For running Makefile commands (usually pre-installed on macOS/Linux)

## Quick Start

### 1. Clone and Setup

```bash
# Clone the repository
git clone git@github.com:dzunghdo/transaction_demo.git
cd transaction_demo

# Download and vendor dependencies
make mod
```

### 2. Start Database

```bash
# Start Docker Compose
docker-compose up -d 

```

### 3. Setup Development Tools (Optional)

Install additional development tools for migrations and documentation:

```bash
make dev-tools

# Generate mock implementations for testing
make mock
```

### 4. Database Migration

This project uses [Goose](https://github.com/pressly/goose) for database migrations. 
Migration files are stored in the `db/migrations/` directory.

```bash
# Run database migrations
make migrate-up

# To rollback migrations (if needed)
make migrate-down
```

### 5. Build and Run Application

```bash
# Build the application
make build

# Run the application
./srv
```

Alternatively, you can run directly with Go:

```bash
go run cmd/srv/main.go
```

The application will start on `http://localhost:10000`

## Configuration

The application uses environment-based configuration files located in `app/config/env/`. 
Different configuration files can be used for different environments (e.g., `local.yaml`, `staging.yaml`, `production.yaml`).
You can set the `APP_ENV` environment variable to specify which configuration file to use.

## Database Setup

The project uses PostgreSQL 15 with the following default configuration:

- **Host**: localhost
- **Port**: 15432 (mapped from container port 5432)
- **Database**: example_db
- **Username**: postgres
- **Password**: root123

## Available Make Commands

```bash
# Install and vendor Go dependencies
make mod

# Build the application binary
make build

# Install development tools (Swagger, Goose migration tool, etc.)
make dev-tools

# Run database migrations up
make migrate-up

# Run database migrations down
make migrate-down

# Generate Swagger documentation
make swag

# Generate mock implementations for interfaces
make mock
```

## Project Structure

```
transaction_demo/
├── app/
│   ├── config/           # Configuration management
│   ├── constant/         # Application constants
│   ├── domain/          # Domain layer (entities, repositories, services)
│   ├── external/        # External integrations (database implementations)
│   ├── interface/       # Interface layer (API handlers, routes)
│   ├── registry/        # Dependency injection setup
│   └── usecase/         # Business logic layer
├── cmd/
│   └── srv/             # Application entry point
├── db/
│   └── migrations/      # Database migration files
├── docker/              # Docker configuration
├── docs/                # API documentation
└── vendor/              # Vendored dependencies
```

## API Documentation

To generate API documentation:

```bash
make swag
```

This will generate Swagger documentation in the `docs/` directory.

## Development

### Database Migrations

This project uses Goose for database schema management. Create new migration files in `db/migrations/` using Goose:

```bash
# Create a new migration file
goose -dir db/migrations create migration_name sql

# Run all pending migrations
make migrate-up

# Rollback the last migration
make migrate-down

# Check migration status
goose -dir db/migrations postgres "postgresql://postgres:root123@localhost:15432/example_db?sslmode=disable" status
```

### Mock Generation

This project uses [GoMock](https://github.com/golang/mock) for generating mock implementations of interfaces for testing purposes. Mock files are automatically generated from interfaces with `//go:generate` directives.

```bash
# Generate all mock implementations
make mock

# Alternatively, run go generate directly
go generate ./...

# Generate mocks for specific package
go generate ./app/domain/repository/...
```

Mock files are generated in the following locations:
- `app/domain/repository/mock/` - Repository interface mocks
- `cmd/shared/db/mock/` - Database transaction manager mocks

**Note**: Make sure you have installed the development tools first (`make dev-tools`) to ensure GoMock is available.