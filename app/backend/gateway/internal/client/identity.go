package client

import (
	identityv1 "cardhub/backend/protocol/gen/identity/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func DialIdentity(
	target string,
) (identityv1.IdentityServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}

	return identityv1.NewIdentityServiceClient(conn), conn, nil
}
