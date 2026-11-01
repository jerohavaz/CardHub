// Package client connects the public API to internal gRPC services.
package client

import (
	examplev1 "cardhub/backend/protocol/gen/example/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func DialExample(target string) (examplev1.ExampleServiceClient, *grpc.ClientConn, error) {
	// Insecure credentials are ONLY for internal trusted Compose networking.
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	return examplev1.NewExampleServiceClient(conn), conn, nil
}
