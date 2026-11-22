# CardHub Backend

Go microservices backend using HTTP/JSON, gRPC, PostgreSQL, and Docker Compose.

## Stack

- **Go** — backend services
- **Chi + Huma** — HTTP gateway and OpenAPI
- **gRPC + Protobuf** — internal communication
- **PostgreSQL + Goose** — database and migrations
- **JWT + Argon2id** — authentication

## Getting Started

Requires Docker and Docker Compose. Native development requires Go 1.27+.

From `app/backend/`:

```bash
cp .env.example .env
make up
```

API: `http://localhost:8080`

Swagger UI: `http://localhost:8080/docs`

## Services

| Directory | Purpose |
|---|---|
| `gateway/` | HTTP API, authentication middleware, gRPC clients |
| `identity/` | User accounts, login, JWTs, refresh sessions |
| `service-example/` | Example gRPC service |
| `protocol/` | Protobuf definitions and generated code |
| `foundation/` | Shared logging, configuration, database and gRPC utilities |
| `migrations/` | Service-specific SQL migrations |

## Authentication

| Method | Endpoint | Access |
|---|---|---|
| POST | `/v1/auth/register` | Public |
| POST | `/v1/auth/login` | Public |
| POST | `/v1/auth/refresh` | Public |
| POST | `/v1/auth/logout` | Public |
| GET | `/v1/users/me` | Bearer JWT |

Other application endpoints require a valid access token.

## Development

Run commands from `app/backend/`:

```bash
make up          # Build and start services
make down        # Stop services
make logs        # View logs
make test        # Run tests
make migrate     # Apply database migrations
make generate    # Generate gRPC code
make lint        # Lint protobuf definitions
make fmt         # Format Go code
```

Configuration is managed through `.env` and `compose.yaml`.

**Note:** `make clean` deletes Docker volumes, including local database data and JWT keys.

## Documentation

- [Architecture](docs/ARCHITECTURE.md)
- [Adding a Service](docs/ADDING_A_SERVICE.md)