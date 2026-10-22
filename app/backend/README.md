# CardHub backend — minimal Go microservice scaffold

This repository is **an executable architectural example**, not the complete CardHub backend. It demonstrates **HTTP/JSON → gateway/OpenAPI → gRPC/protobuf → service → PostgreSQL**, with structured logging, centralized SQL migrations, tests, and Docker Compose. It intentionally does **not** implement authentication, Pokémon catalog, inventory/trading, or object storage.

## Start

**Prerequisites:** Docker Engine + Compose v2, internet access for the first image/dependency downloads, `curl` for the smoke test. Native Go development needs Go **1.27+**; Buf is needed only when editing `.proto` files.

```bash
# From backend/
make up               # Builds gateway, example-service and Goose; starts PostgreSQL
make smoke            # Creates a record via HTTP, fetches it, tests OpenAPI/404
make logs             # Optional: inspect JSON logs
make down             # Stop; keep the database volume
# make clean          # CAUTION: stops everything and DELETES the local DB volume
```

For the manual request sequence (once `make up` has completed):

```bash
curl -fsS http://127.0.0.1:8080/healthz
curl -fsS http://127.0.0.1:8080/openapi.json
curl -fsS -X POST http://127.0.0.1:8080/v1/examples \
  -H 'Content-Type: application/json' \
  -d '{"name":"Pikachu"}'
# Copy id from the returned JSON:
curl -fsS http://127.0.0.1:8080/v1/examples/<id>
```

Interactive documentation: **http://127.0.0.1:8080/docs**. OpenAPI JSON: **http://127.0.0.1:8080/openapi.json**, generated at runtime by Huma from the registered HTTP handlers and Go types; there is no hand-edited OpenAPI YAML file.

`docker compose up --build -d` automatically executes the **one-shot Goose migration container** after PostgreSQL is healthy. The gRPC service starts only after a successful migration. The gateway listens on **127.0.0.1:8080**. PostgreSQL and gRPC stay on the private Compose network.

## Repository map

```text
backend/
├── gateway/                # Separate Go module: Chi + Huma HTTP gateway and gRPC client
│   ├── cmd/gateway/        # Program entry point
│   └── internal/           # HTTP API + gRPC client adapter
├── service-example/        # Separate Go module: example business logic
│   ├── cmd/service-example/# Program entry point
│   └── internal/
│       ├── service/        # Domain rules and repository interface
│       ├── transport/grpc/ # gRPC handler
│       └── postgres/       # SQL repository implementation
├── protocol/               # Separate Go module: .proto source + generated client/server
│   ├── proto/example/v1/
│   ├── gen/example/v1/
│   ├── buf.yaml
│   └── buf.gen.yaml
├── foundation/             # Separate Go module: config, Zap logging, gRPC, pgxpool
├── migrations/example/     # Central location; isolated migration stream for this DB
├── scripts/smoke.sh        # End-to-end HTTP integration smoke check
├── docs/
│   ├── ARCHITECTURE.md
│   └── ADDING_A_SERVICE.md
├── compose.yaml
└── Makefile
```

## Four independent Go modules

| Directory | Module name | Imports |
| --- | --- | --- |
| `gateway` | `cardhub/backend/gateway` | `foundation`, `protocol` |
| `service-example` | `cardhub/backend/service-example` | `foundation`, `protocol` |
| `foundation` | `cardhub/backend/foundation` | Third-party infrastructure libraries only |
| `protocol` | `cardhub/backend/protocol` | Protobuf/gRPC runtimes only |

There is **no root `go.mod` or `go.work`**. The two applications declare local `replace cardhub/backend/foundation => ../foundation` and `replace cardhub/backend/protocol => ../protocol` in their own module files, following TimBRA's approach. Dockerfiles use the *repository root* as their build context so relative replacements also resolve during Docker builds. Each module may have its own generated `go.sum`; run `make deps` with an online Go installation to generate/tidy them. The Dockerfiles also tidy during initial builds.

## Developer commands

```bash
make deps        # Run `go mod tidy` in all four modules; requires Go 1.27+ / network
make test        # Run Go unit tests in all four modules (no database required)
make generate    # Regenerate protobuf and gRPC Go files (Buf + remote plugins; network)
make lint        # Buf lint
make migrate     # Re-run Goose on the example DB
make smoke       # End-to-end test against the Compose stack
```

Buf code-generation plugins are pinned in `protocol/buf.gen.yaml`. `protocol/gen/**` is checked into the scaffold so users do not need Buf merely to run Docker. Regenerate after modifying a `.proto` file and commit the resulting `.pb.go` and `_grpc.pb.go` files. HTTP/OpenAPI generation is separate from protobuf code generation.

## HTTP endpoints

| HTTP | gRPC | Outcome |
| --- | --- | --- |
| `GET /healthz` | — | `200 ok` (gateway liveness only, not a DB readiness check) |
| `GET /docs` | — | Interactive HTTP API documentation |
| `GET /openapi.json` | — | Generated OpenAPI schema |
| `POST /v1/examples` | `ExampleService.CreateExample` | HTTP `201` with `{id,name}` |
| `GET /v1/examples/{id}` | `ExampleService.GetExample` | HTTP `200` or `404` |

The gateway translates gRPC `InvalidArgument`, `NotFound`, `Unavailable`, and `DeadlineExceeded` into meaningful HTTP errors. The service applies its own validation and handles its own SQL errors; the gateway never connects to the database.

## Database and migrations

A **single PostgreSQL instance** hosts the example database `cardhub_example`. SQL files are managed at `migrations/example/`, but only `service-example/internal/postgres` queries its tables. `migrate` is an explicit one-shot service; application startup never silently mutates its schema.

When introducing real services, create `migrations/identity/`, `migrations/catalog/`, or `migrations/inventory/` **alongside separate databases and Goose migration histories** (rather than mixing ownership just because SQL files share a parent folder). See [how to add another service](docs/ADDING_A_SERVICE.md).

## Scope and limitations

- **Development only**: example routes deliberately have no JWT/authentication; do not publicly deploy them as-is.
- Internal gRPC uses plaintext inside a trusted Docker network. Add TLS/mTLS for untrusted service networks.
- Database credentials in Compose are example-only; use secret management outside local development.
- The example uses `pgx` and direct SQL to avoid unnecessary ORM/sqlc setup. Both can be revisited as query volume grows.
- `go.sum` files are normally produced by `make deps`; the provided environment lacked network module downloads. No claim of an already-verified Compose run is made in this package. Run `make up && make smoke` in a Docker-enabled environment to validate the full workflow.

For rationale and component boundaries, read [architecture](docs/ARCHITECTURE.md).
