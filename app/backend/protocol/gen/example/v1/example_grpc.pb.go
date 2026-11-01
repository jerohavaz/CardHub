// Code generated to match protoc-gen-go-grpc conventions. DO NOT EDIT.
// Run `make generate` to regenerate with Buf remote plugins.
// source: proto/example/v1/example.proto

package examplev1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const _ = grpc.SupportPackageIsVersion9
const (
	ExampleService_CreateExample_FullMethodName = "/cardhub.example.v1.ExampleService/CreateExample"
	ExampleService_GetExample_FullMethodName    = "/cardhub.example.v1.ExampleService/GetExample"
)

type ExampleServiceClient interface {
	CreateExample(ctx context.Context, in *CreateExampleRequest, opts ...grpc.CallOption) (*CreateExampleResponse, error)
	GetExample(ctx context.Context, in *GetExampleRequest, opts ...grpc.CallOption) (*GetExampleResponse, error)
}
type exampleServiceClient struct{ cc grpc.ClientConnInterface }

func NewExampleServiceClient(cc grpc.ClientConnInterface) ExampleServiceClient {
	return &exampleServiceClient{cc}
}
func (c *exampleServiceClient) CreateExample(ctx context.Context, in *CreateExampleRequest, opts ...grpc.CallOption) (*CreateExampleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateExampleResponse)
	err := c.cc.Invoke(ctx, ExampleService_CreateExample_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *exampleServiceClient) GetExample(ctx context.Context, in *GetExampleRequest, opts ...grpc.CallOption) (*GetExampleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetExampleResponse)
	err := c.cc.Invoke(ctx, ExampleService_GetExample_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type ExampleServiceServer interface {
	CreateExample(context.Context, *CreateExampleRequest) (*CreateExampleResponse, error)
	GetExample(context.Context, *GetExampleRequest) (*GetExampleResponse, error)
	mustEmbedUnimplementedExampleServiceServer()
}
type UnimplementedExampleServiceServer struct{}

func (UnimplementedExampleServiceServer) CreateExample(context.Context, *CreateExampleRequest) (*CreateExampleResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method CreateExample not implemented")
}
func (UnimplementedExampleServiceServer) GetExample(context.Context, *GetExampleRequest) (*GetExampleResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetExample not implemented")
}
func (UnimplementedExampleServiceServer) mustEmbedUnimplementedExampleServiceServer() {}
func (UnimplementedExampleServiceServer) testEmbeddedByValue()                        {}

type UnsafeExampleServiceServer interface{ mustEmbedUnimplementedExampleServiceServer() }

func RegisterExampleServiceServer(s grpc.ServiceRegistrar, srv ExampleServiceServer) {
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ExampleService_ServiceDesc, srv)
}
func _ExampleService_CreateExample_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateExampleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExampleServiceServer).CreateExample(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: ExampleService_CreateExample_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExampleServiceServer).CreateExample(ctx, req.(*CreateExampleRequest))
	}
	return interceptor(ctx, in, info, handler)
}
func _ExampleService_GetExample_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetExampleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExampleServiceServer).GetExample(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: ExampleService_GetExample_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExampleServiceServer).GetExample(ctx, req.(*GetExampleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var ExampleService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cardhub.example.v1.ExampleService",
	HandlerType: (*ExampleServiceServer)(nil),
	Methods:     []grpc.MethodDesc{{MethodName: "CreateExample", Handler: _ExampleService_CreateExample_Handler}, {MethodName: "GetExample", Handler: _ExampleService_GetExample_Handler}},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "proto/example/v1/example.proto",
}
