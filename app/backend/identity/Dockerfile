FROM golang:1.27.3-alpine AS builder
WORKDIR /workspace
COPY foundation ./foundation
COPY protocol ./protocol
COPY identity ./identity
WORKDIR /workspace/identity
RUN go mod tidy && CGO_ENABLED=0 go build -trimpath -o /out/identity ./cmd/identity

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /out/identity /identity
EXPOSE 9000
ENTRYPOINT ["/identity"]
