package examplev1

import (
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestGeneratedProtobufRoundTrip(t *testing.T) {
	before := &CreateExampleResponse{Id: "123e4567-e89b-12d3-a456-426614174000", Name: "Pikachu"}
	data, err := proto.Marshal(before)
	if err != nil {
		t.Fatal(err)
	}
	var after CreateExampleResponse
	if err := proto.Unmarshal(data, &after); err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(before, &after) {
		t.Fatalf("round trip: %#v -> %#v", before, &after)
	}
	if ExampleService_ServiceDesc.ServiceName != "cardhub.example.v1.ExampleService" {
		t.Fatal("gRPC name mismatch")
	}
}
