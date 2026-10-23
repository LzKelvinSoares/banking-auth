# Banking Auth

A Go HTTP service that authenticates banking users and issues JWT access and refresh tokens. It also verifies access tokens and checks whether a token's role may access a requested route.

## Requirements

- Go 1.27.1 or later
- SQL Server
- A database containing the authentication tables described below

The service currently configures the SQL Server repository at startup. Although the shared database library supports other database types, this application does not select them.

## Configuration

The service loads environment variables from `.env` in the working directory, or from the process environment. Configure:

| Variable | Description |
| --- | --- |
| `SERVER_ADDRESS` | Address to bind the HTTP server to |
| `SERVER_PORT` | HTTP server port |
| `DB_USER` | SQL Server username |
| `DB_PASSWORD` | SQL Server password |
| `DB_ADDRESS` | SQL Server hostname |
| `DB_PORT` | SQL Server port (commonly `1433`) |
| `DB_NAME` | Database name |
| `HMAC_SECRET` | Secret used to sign and validate tokens; must contain at least 32 bytes |

For local development, create a `.env` file with these variables. Keep real credentials and signing secrets out of source control. When running in a container, set `SERVER_ADDRESS=0.0.0.0` and use the same `SERVER_PORT` that you publish.

The database must provide:

- `users` with `username`, `password`, `customer_id`, and `role` columns
- `accounts` with `customer_id` and `account_id` columns
- `refresh_token_store` with a `refresh_token` column

The login query uses these tables to load a user's role and account identifiers. Refresh tokens are registered in `refresh_token_store`.

## Run

```sh
go mod download
go run .
```

The server listens at `SERVER_ADDRESS:SERVER_PORT`.

To build the Docker image:

```sh
docker build -f dockerfile -t banking-auth .
```

Provide the configuration variables to the container at runtime and publish its configured port.

## API

All request and response bodies use JSON.

### `POST /auth/login`

Authenticates a user and returns an access token and a registered refresh token.

Request:

```json
{
  "username": "alice",
  "password": "secret"
}
```

Response:

```json
{
  "access_token": "<jwt>",
  "refresh_token": "<jwt>"
}
```

Access tokens expire after one hour. Refresh tokens expire after 30 days.

### `GET /auth/verify`

Verifies the access token and checks its role against the requested route. Supply these query parameters:

| Parameter | Description |
| --- | --- |
| `token` | Access token to verify |
| `routeName` | Protected route name, such as `GetCustomer` or `MakeTransaction` |
| `customer_id` | Customer identifier; checked against the token for user-role tokens |
| `account_id` | Optional account identifier; if supplied for a user-role token, it must belong to that token |

Configured route permissions:

| Role | Allowed route names |
| --- | --- |
| `admin` | `GetAllCustomers`, `GetCustomer`, `NewAccount`, `MakeTransaction` |
| `user` | `GetCustomer`, `MakeTransaction` |

Successful response:

```json
{
  "isAuthorized": true
}
```

### `POST /auth/refresh`

Validates the supplied access token and registered refresh token, then returns a new access token.

Request:

```json
{
  "access_token": "<current-access-jwt>",
  "refresh_token": "<registered-refresh-jwt>"
}
```

Response:

```json
{
  "access_token": "<new-access-jwt>"
}
```

## Project structure

- `app/` — HTTP handlers and server setup
- `constants/` — API route paths
- `domain/` — token claims, repository interfaces, and database repositories
- `dto/` — API request and response models
- `service/` — authentication, token verification, and refresh logic
- `main.go` — application entry point
