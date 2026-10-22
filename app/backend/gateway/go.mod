module cardhub/backend/gateway

go 1.27

require (
	cardhub/backend/foundation v0.0.0
	cardhub/backend/protocol v0.0.0
	github.com/danielgtaylor/huma/v2 v2.39.1
	github.com/go-chi/chi/v5 v5.5.0
	go.uber.org/zap v1.28.0
	google.golang.org/grpc v1.84.0
)

require (
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260706201446-f0a921348800 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

replace cardhub/backend/foundation => ../foundation

replace cardhub/backend/protocol => ../protocol
