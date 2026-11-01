FROM golang:1.27.3-alpine AS builder
# Goose is a build-time tool, not a runtime dependency of application services.
RUN CGO_ENABLED=0 go install -tags='no_sqlite3 no_mysql no_mssql no_clickhouse no_azuresql no_ydb no_vertica no_libsql' github.com/pressly/goose/v3/cmd/goose@v3.29.0

FROM alpine:3.23
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY migrations /migrations
ENTRYPOINT ["goose"]
CMD ["up"]
