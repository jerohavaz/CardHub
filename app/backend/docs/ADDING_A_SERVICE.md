# How to add another microservice

This is a concrete example using **Catalog**, whose actual implementation is **not** part of the scaffold. Start by reading the working `service-example/` and the actual `protocol/proto/example/v1/example.proto`.

## 1. Add a new module, not a placeholder package

```text
backend/catalog/
├── go.mod
├── Dockerfile
├── cmd/catalog/main.go
└── internal/
    ├── service/service.go
    ├── transport/grpc/handler.go
    └── postgres/repository.go
```

The module is `cardhub/backend/catalog` with a separate `go.mod` and the same pinned dependency versions as the other services. Use TimBRA-style local dependencies:

```go
require (
    cardhub/backend/foundation v0.0.0
    cardhub/backend/protocol v0.0.0
)
replace cardhub/backend/foundation => ../foundation
replace cardhub/backend/protocol => ../protocol
```

## 2. Define its versioned protobuf contract

Add `protocol/proto/catalog/v1/catalog.proto`:

```proto
syntax = "proto3";
package cardhub.catalog.v1;
option go_package = "cardhub/backend/protocol/gen/catalog/v1;catalogv1";

service CatalogService {
  rpc GetCard(GetCardRequest) returns (GetCardResponse);
}
message GetCardRequest { string id = 1; }
message GetCardResponse { string id = 1; string name = 2; }
```

From the root, run `make lint && make generate`. Commit `.proto` source and generated `protocol/gen/catalog/v1/` Go code together. Preserve field numbers when evolving existing contracts; changing them breaks wire compatibility.

## 3. Implement the three internal layers

`internal/service` owns domain rules and a narrow `Repository` interface, **without** pgx or gRPC imports. `internal/transport/grpc` implements the generated `CatalogServiceServer`, translates protobuf types and maps errors to gRPC status codes. `internal/postgres` owns SQL and translates `pgx.ErrNoRows` to domain-level not-found errors. Keep SQL out of `gateway` and `foundation`.

## 4. Add an independent database with centrally stored migrations

Create `migrations/catalog/00001_create_cards.sql`, with Goose `Up`/`Down` comments. Provision `catalog_db` in the existing PostgreSQL instance with separate credentials, add a `migrate-catalog` Compose one-shot job with `GOOSE_MIGRATION_DIR=/migrations/catalog`, and point Catalog only at `catalog_db`. Each database has its own `goose_db_version` history.

For local development, provisioning scripts under `postgres-init/` can create additional databases/roles through `/docker-entrypoint-initdb.d/` **on a fresh data volume**. For deployed environments, provision databases separately from application schema migrations. Do not give Catalog SQL access to Identity or Inventory tables.

## 5. Wire the new service and public gateway

- In `catalog/cmd/catalog/main.go`, initialize `foundation/logging`, `foundation/postgres` and `foundation/grpcserver`, then register `catalogv1.RegisterCatalogServiceServer`.
- Copy `service-example/Dockerfile`, changing the copied service directory/binary. Keep `context: .` and private gRPC networking.
- In `compose.yaml`, declare `catalog` with `CATALOG_GRPC_ADDR=catalog:9001`, its database URL and dependency on a successful `migrate-catalog`.
- Add a generated gRPC client adapter in `gateway/internal/client/catalog.go`.
- Register a Huma HTTP route like `GET /v1/cards/{id}` in `gateway/internal/api/`; this updates `/docs` and `/openapi.json` automatically. **Do not write an additional OpenAPI YAML file.**
- Handle status mapping at the gateway boundary, but leave domain authorization and writes within Catalog.

## 6. Verify the vertical slice

1. Unit-test Catalog logic using a fake Repository (example: `service-example/internal/service/service_test.go`).
2. Test gRPC responses and codes using an in-memory service (example: `service-example/internal/transport/grpc/handler_test.go`).
3. Mock the new generated gRPC client in the HTTP gateway tests and assert both route behavior and an OpenAPI path entry.
4. Extend `scripts/smoke.sh` with a request going through **HTTP → Catalog gRPC → Catalog PostgreSQL**.
5. Run `make deps && make test && make up && make smoke` and verify the old example still works.

Use the same process for Identity and Inventory. **Keep Collection and Trading inside Inventory** for atomic trade execution, rather than introducing a separate Trading service. Add real authentication/authorization before exposing non-demo application endpoints.
