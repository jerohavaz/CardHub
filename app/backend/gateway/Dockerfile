FROM golang:1.27.3-alpine AS builder
WORKDIR /workspace
COPY foundation ./foundation
COPY protocol ./protocol
COPY gateway ./gateway
WORKDIR /workspace/gateway
# Resolve sums during build; run `make deps` locally to write go.sum files first.
RUN go mod tidy && CGO_ENABLED=0 go build -trimpath -o /out/gateway ./cmd/gateway

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /out/gateway /gateway
EXPOSE 8080
ENTRYPOINT ["/gateway"]
