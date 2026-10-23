# Banking Auth

A Go authentication service for a banking application. It issues JWT access tokens, validates token claims for protected routes, and supports refresh-token flows for session renewal.

## Overview

This project exposes a small REST API for authentication and authorization. It is designed to work with the banking sample API and validates user credentials against a backing database before generating signed tokens.

Key capabilities:

- User login with username/password
- Access token generation and validation
- Refresh token generation and renewal
- Role-based authorization checks
- SQL Server or MySQL-compatible repository backends
- Environment-based configuration

## Project Structure

- `app/` — HTTP handlers and server bootstrap
- `constants/` — route constants
- `domain/` — token claims, repositories, and auth domain logic
- `dto/` — request and response models
- `service/` — business logic for login, verification, and refresh
- `main.go` — application entry point

## API Endpoints

### POST `/auth/login`
Authenticates a user and returns a pair of tokens.

Request body:

```json
{
  "username": "alice",
  "password": "secret"
}
```

Response:

```json
{
  "accessToken": "<jwt>",
  "refreshToken": "<jwt>"
}
```

### GET `/auth/verify`
Validates a provided JWT against the configured secret and checks whether the user is authorized for the requested route.

Query parameters commonly include:

- `token` — JWT to validate
- `routeName` — route being accessed
- `customer_id` — expected customer identifier for user-scoped requests
- `account_id` — optional account identifier verification

### POST `/auth/refresh`
Refreshes an access token using a valid refresh token.

Request body:

```json
{
  "accessToken": "<jwt>",
  "refreshToken": "<jwt>"
}
```

## Configuration

The application reads configuration from environment variables. A sample `.env` file is included:

```env
SERVER_ADDRESS="localhost"
SERVER_PORT="8081"

DB_USER="banking-master"
DB_PASSWORD="your-password"
DB_ADDRESS="your-db-host"
DB_PORT="1433"
DB_NAME="banking"
```

These values are used to bind the HTTP server and connect to the configured database.

## Running the Service

1. Install dependencies:

```bash
go mod download
```

2. Start the service:

```bash
go run .
```

3. The service will listen on the configured address and port, defaulting to `http://localhost:8081`.

## Notes

- The application uses JWT signing with the secret defined in the domain layer.
- Authentication and authorization rules are enforced by service and repository layers.
- The project is intended to be paired with a banking API that consumes the generated tokens.

## License

This project is provided as a sample codebase for educational and demonstration purposes.
