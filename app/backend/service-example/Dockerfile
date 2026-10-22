FROM golang:1.27.3-alpine AS builder
WORKDIR /workspace
COPY foundation ./foundation
COPY protocol ./protocol
COPY service-example ./service-example
WORKDIR /workspace/service-example
RUN go mod tidy && CGO_ENABLED=0 go build -trimpath -o /out/service-example ./cmd/service-example

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /out/service-example /service-example
EXPOSE 9000
ENTRYPOINT ["/service-example"]
