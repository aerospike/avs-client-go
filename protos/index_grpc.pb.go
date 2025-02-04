// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: index.proto

package protos

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	IndexService_Create_FullMethodName            = "/aerospike.vector.IndexService/Create"
	IndexService_Update_FullMethodName            = "/aerospike.vector.IndexService/Update"
	IndexService_Drop_FullMethodName              = "/aerospike.vector.IndexService/Drop"
	IndexService_List_FullMethodName              = "/aerospike.vector.IndexService/List"
	IndexService_Get_FullMethodName               = "/aerospike.vector.IndexService/Get"
	IndexService_GetStatus_FullMethodName         = "/aerospike.vector.IndexService/GetStatus"
	IndexService_GcInvalidVertices_FullMethodName = "/aerospike.vector.IndexService/GcInvalidVertices"
)

// IndexServiceClient is the client API for IndexService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Service to manage indices.
type IndexServiceClient interface {
	// Create an index.
	Create(ctx context.Context, in *IndexCreateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Create an index.
	Update(ctx context.Context, in *IndexUpdateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Drop an index.
	Drop(ctx context.Context, in *IndexDropRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// List available indices.
	List(ctx context.Context, in *IndexListRequest, opts ...grpc.CallOption) (*IndexDefinitionList, error)
	// Get the index definition.
	Get(ctx context.Context, in *IndexGetRequest, opts ...grpc.CallOption) (*IndexDefinition, error)
	// Query status of an index.
	// NOTE: API is subject to change.
	GetStatus(ctx context.Context, in *IndexStatusRequest, opts ...grpc.CallOption) (*IndexStatusResponse, error)
	// Garbage collect vertices identified as invalid before cutoff timestamp.
	GcInvalidVertices(ctx context.Context, in *GcInvalidVerticesRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type indexServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewIndexServiceClient(cc grpc.ClientConnInterface) IndexServiceClient {
	return &indexServiceClient{cc}
}

func (c *indexServiceClient) Create(ctx context.Context, in *IndexCreateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, IndexService_Create_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexServiceClient) Update(ctx context.Context, in *IndexUpdateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, IndexService_Update_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexServiceClient) Drop(ctx context.Context, in *IndexDropRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, IndexService_Drop_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexServiceClient) List(ctx context.Context, in *IndexListRequest, opts ...grpc.CallOption) (*IndexDefinitionList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(IndexDefinitionList)
	err := c.cc.Invoke(ctx, IndexService_List_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexServiceClient) Get(ctx context.Context, in *IndexGetRequest, opts ...grpc.CallOption) (*IndexDefinition, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(IndexDefinition)
	err := c.cc.Invoke(ctx, IndexService_Get_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexServiceClient) GetStatus(ctx context.Context, in *IndexStatusRequest, opts ...grpc.CallOption) (*IndexStatusResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(IndexStatusResponse)
	err := c.cc.Invoke(ctx, IndexService_GetStatus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexServiceClient) GcInvalidVertices(ctx context.Context, in *GcInvalidVerticesRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, IndexService_GcInvalidVertices_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IndexServiceServer is the server API for IndexService service.
// All implementations must embed UnimplementedIndexServiceServer
// for forward compatibility.
//
// Service to manage indices.
type IndexServiceServer interface {
	// Create an index.
	Create(context.Context, *IndexCreateRequest) (*emptypb.Empty, error)
	// Create an index.
	Update(context.Context, *IndexUpdateRequest) (*emptypb.Empty, error)
	// Drop an index.
	Drop(context.Context, *IndexDropRequest) (*emptypb.Empty, error)
	// List available indices.
	List(context.Context, *IndexListRequest) (*IndexDefinitionList, error)
	// Get the index definition.
	Get(context.Context, *IndexGetRequest) (*IndexDefinition, error)
	// Query status of an index.
	// NOTE: API is subject to change.
	GetStatus(context.Context, *IndexStatusRequest) (*IndexStatusResponse, error)
	// Garbage collect vertices identified as invalid before cutoff timestamp.
	GcInvalidVertices(context.Context, *GcInvalidVerticesRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedIndexServiceServer()
}

// UnimplementedIndexServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedIndexServiceServer struct{}

func (UnimplementedIndexServiceServer) Create(context.Context, *IndexCreateRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedIndexServiceServer) Update(context.Context, *IndexUpdateRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedIndexServiceServer) Drop(context.Context, *IndexDropRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Drop not implemented")
}
func (UnimplementedIndexServiceServer) List(context.Context, *IndexListRequest) (*IndexDefinitionList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedIndexServiceServer) Get(context.Context, *IndexGetRequest) (*IndexDefinition, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedIndexServiceServer) GetStatus(context.Context, *IndexStatusRequest) (*IndexStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStatus not implemented")
}
func (UnimplementedIndexServiceServer) GcInvalidVertices(context.Context, *GcInvalidVerticesRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GcInvalidVertices not implemented")
}
func (UnimplementedIndexServiceServer) mustEmbedUnimplementedIndexServiceServer() {}
func (UnimplementedIndexServiceServer) testEmbeddedByValue()                      {}

// UnsafeIndexServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IndexServiceServer will
// result in compilation errors.
type UnsafeIndexServiceServer interface {
	mustEmbedUnimplementedIndexServiceServer()
}

func RegisterIndexServiceServer(s grpc.ServiceRegistrar, srv IndexServiceServer) {
	// If the following call pancis, it indicates UnimplementedIndexServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&IndexService_ServiceDesc, srv)
}

func _IndexService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IndexCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IndexService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexServiceServer).Create(ctx, req.(*IndexCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IndexService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IndexUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IndexService_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexServiceServer).Update(ctx, req.(*IndexUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IndexService_Drop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IndexDropRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServiceServer).Drop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IndexService_Drop_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexServiceServer).Drop(ctx, req.(*IndexDropRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IndexService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IndexListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IndexService_List_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexServiceServer).List(ctx, req.(*IndexListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IndexService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IndexGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IndexService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexServiceServer).Get(ctx, req.(*IndexGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IndexService_GetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IndexStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServiceServer).GetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IndexService_GetStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexServiceServer).GetStatus(ctx, req.(*IndexStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IndexService_GcInvalidVertices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GcInvalidVerticesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServiceServer).GcInvalidVertices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IndexService_GcInvalidVertices_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexServiceServer).GcInvalidVertices(ctx, req.(*GcInvalidVerticesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// IndexService_ServiceDesc is the grpc.ServiceDesc for IndexService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var IndexService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "aerospike.vector.IndexService",
	HandlerType: (*IndexServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _IndexService_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _IndexService_Update_Handler,
		},
		{
			MethodName: "Drop",
			Handler:    _IndexService_Drop_Handler,
		},
		{
			MethodName: "List",
			Handler:    _IndexService_List_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _IndexService_Get_Handler,
		},
		{
			MethodName: "GetStatus",
			Handler:    _IndexService_GetStatus_Handler,
		},
		{
			MethodName: "GcInvalidVertices",
			Handler:    _IndexService_GcInvalidVertices_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "index.proto",
}
