# Wallet API

A mini wallet and expense tracker REST API built with Go, Gin, PostgreSQL, GORM, and JWT authentication.

## Features
- User authentication with JWT
- Wallet management (deposit, withdraw, transfer)
- Transaction history with filtering and pagination
- Monthly expense summary by category
- Role-based access control (user/admin)
- Race condition prevention with database row locking

## Tech Stack
- Go + Gin (REST API)
- PostgreSQL + GORM (database)
- JWT (authentication)
- Swagger (API docs)
- Docker (containerization)

## Prerequisites
- Go 1.24+
- PostgreSQL
- Docker (optional)

## Running with Docker
```bash
docker compose up --build
```
The API will be available at `http://localhost:8080`.

## Running without Docker

1. Create a PostgreSQL database
2. Copy `.env.example` to `.env` and fill in your values
3. Run the server:
```bash
go run .
```

## Environment Variables
| Variable | Description | Example |
|----------|-------------|---------|
| DB_HOST | Database host | localhost |
| DB_PORT | Database port | 5432 |
| DB_USER | Database user | postgres |
| DB_PASSWORD | Database password | postgres |
| DB_NAME | Database name | wallet |
| TEST_DB_NAME | Test database name | wallet_test |
| JWT_SECRET | JWT signing secret | supersecret |

## Running Tests
```bash
go test ./tests/... -v
```

## API Documentation
Start the server and visit `http://localhost:8080/swagger/index.html`

## Endpoints
| Method | Path | Description | Auth |
|--------|------|-------------|------|
| POST | /signup | Create account | No |
| POST | /login | Login and get token | No |
| GET | /wallet | Get current wallet | Yes |
| POST | /wallet/deposit | Deposit money | Yes |
| POST | /wallet/withdraw | Withdraw money | Yes |
| POST | /wallet/transfer | Transfer to another user | Yes |
| GET | /transactions | List transactions | Yes |
| GET | /transactions/summary | Monthly summary by category | Yes |