// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
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
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// IndexServiceClient is the client API for IndexService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type IndexServiceClient interface {
	// Create an index.
	Create(ctx context.Context, in *IndexDefinition, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Drop an index.
	Drop(ctx context.Context, in *IndexId, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// List available indices.
	List(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*IndexDefinitionList, error)
	// Get the index definition.
	Get(ctx context.Context, in *IndexId, opts ...grpc.CallOption) (*IndexDefinition, error)
	// Query status of an index.
	// NOTE: API is subject to change.
	GetStatus(ctx context.Context, in *IndexId, opts ...grpc.CallOption) (*IndexStatusResponse, error)
}

type indexServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewIndexServiceClient(cc grpc.ClientConnInterface) IndexServiceClient {
	return &indexServiceClient{cc}
}

func (c *indexServiceClient) Create(ctx context.Context, in *IndexDefinition, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/aerospike.vector.IndexService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexServiceClient) Drop(ctx context.Context, in *IndexId, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/aerospike.vector.IndexService/Drop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexServiceClient) List(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*IndexDefinitionList, error) {
	out := new(IndexDefinitionList)
	err := c.cc.Invoke(ctx, "/aerospike.vector.IndexService/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexServiceClient) Get(ctx context.Context, in *IndexId, opts ...grpc.CallOption) (*IndexDefinition, error) {
	out := new(IndexDefinition)
	err := c.cc.Invoke(ctx, "/aerospike.vector.IndexService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexServiceClient) GetStatus(ctx context.Context, in *IndexId, opts ...grpc.CallOption) (*IndexStatusResponse, error) {
	out := new(IndexStatusResponse)
	err := c.cc.Invoke(ctx, "/aerospike.vector.IndexService/GetStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IndexServiceServer is the server API for IndexService service.
// All implementations must embed UnimplementedIndexServiceServer
// for forward compatibility
type IndexServiceServer interface {
	// Create an index.
	Create(context.Context, *IndexDefinition) (*emptypb.Empty, error)
	// Drop an index.
	Drop(context.Context, *IndexId) (*emptypb.Empty, error)
	// List available indices.
	List(context.Context, *emptypb.Empty) (*IndexDefinitionList, error)
	// Get the index definition.
	Get(context.Context, *IndexId) (*IndexDefinition, error)
	// Query status of an index.
	// NOTE: API is subject to change.
	GetStatus(context.Context, *IndexId) (*IndexStatusResponse, error)
	mustEmbedUnimplementedIndexServiceServer()
}

// UnimplementedIndexServiceServer must be embedded to have forward compatible implementations.
type UnimplementedIndexServiceServer struct {
}

func (UnimplementedIndexServiceServer) Create(context.Context, *IndexDefinition) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedIndexServiceServer) Drop(context.Context, *IndexId) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Drop not implemented")
}
func (UnimplementedIndexServiceServer) List(context.Context, *emptypb.Empty) (*IndexDefinitionList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedIndexServiceServer) Get(context.Context, *IndexId) (*IndexDefinition, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedIndexServiceServer) GetStatus(context.Context, *IndexId) (*IndexStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStatus not implemented")
}
func (UnimplementedIndexServiceServer) mustEmbedUnimplementedIndexServiceServer() {}

// UnsafeIndexServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IndexServiceServer will
// result in compilation errors.
type UnsafeIndexServiceServer interface {
	mustEmbedUnimplementedIndexServiceServer()
}

func RegisterIndexServiceServer(s grpc.ServiceRegistrar, srv IndexServiceServer) {
	s.RegisterService(&IndexService_ServiceDesc, srv)
}

func _IndexService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IndexDefinition)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aerospike.vector.IndexService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexServiceServer).Create(ctx, req.(*IndexDefinition))
	}
	return interceptor(ctx, in, info, handler)
}

func _IndexService_Drop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IndexId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServiceServer).Drop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aerospike.vector.IndexService/Drop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexServiceServer).Drop(ctx, req.(*IndexId))
	}
	return interceptor(ctx, in, info, handler)
}

func _IndexService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aerospike.vector.IndexService/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexServiceServer).List(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _IndexService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IndexId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aerospike.vector.IndexService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexServiceServer).Get(ctx, req.(*IndexId))
	}
	return interceptor(ctx, in, info, handler)
}

func _IndexService_GetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IndexId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServiceServer).GetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aerospike.vector.IndexService/GetStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexServiceServer).GetStatus(ctx, req.(*IndexId))
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
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "index.proto",
}
