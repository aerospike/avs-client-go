// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: vector-db.proto

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

// AboutClient is the client API for About service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AboutClient interface {
	Get(ctx context.Context, in *AboutRequest, opts ...grpc.CallOption) (*AboutResponse, error)
}

type aboutClient struct {
	cc grpc.ClientConnInterface
}

func NewAboutClient(cc grpc.ClientConnInterface) AboutClient {
	return &aboutClient{cc}
}

func (c *aboutClient) Get(ctx context.Context, in *AboutRequest, opts ...grpc.CallOption) (*AboutResponse, error) {
	out := new(AboutResponse)
	err := c.cc.Invoke(ctx, "/aerospike.vector.About/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AboutServer is the server API for About service.
// All implementations must embed UnimplementedAboutServer
// for forward compatibility
type AboutServer interface {
	Get(context.Context, *AboutRequest) (*AboutResponse, error)
	mustEmbedUnimplementedAboutServer()
}

// UnimplementedAboutServer must be embedded to have forward compatible implementations.
type UnimplementedAboutServer struct {
}

func (UnimplementedAboutServer) Get(context.Context, *AboutRequest) (*AboutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedAboutServer) mustEmbedUnimplementedAboutServer() {}

// UnsafeAboutServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AboutServer will
// result in compilation errors.
type UnsafeAboutServer interface {
	mustEmbedUnimplementedAboutServer()
}

func RegisterAboutServer(s grpc.ServiceRegistrar, srv AboutServer) {
	s.RegisterService(&About_ServiceDesc, srv)
}

func _About_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AboutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AboutServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aerospike.vector.About/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AboutServer).Get(ctx, req.(*AboutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// About_ServiceDesc is the grpc.ServiceDesc for About service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var About_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "aerospike.vector.About",
	HandlerType: (*AboutServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _About_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "vector-db.proto",
}

// ClusterInfoClient is the client API for ClusterInfo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClusterInfoClient interface {
	// Get the internal cluster node-Id for this server.
	GetNodeId(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*NodeId, error)
	// Get current cluster-Id for the current cluster.
	GetClusterId(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ClusterId, error)
	// Get the advertised/listening endpoints for all nodes in the cluster, given a listener name.
	GetClusterEndpoints(ctx context.Context, in *ClusterNodeEndpointsRequest, opts ...grpc.CallOption) (*ClusterNodeEndpoints, error)
	// Get per-node owned partition list for all nodes in the cluster.
	GetOwnedPartitions(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ClusterPartitions, error)
}

type clusterInfoClient struct {
	cc grpc.ClientConnInterface
}

func NewClusterInfoClient(cc grpc.ClientConnInterface) ClusterInfoClient {
	return &clusterInfoClient{cc}
}

func (c *clusterInfoClient) GetNodeId(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*NodeId, error) {
	out := new(NodeId)
	err := c.cc.Invoke(ctx, "/aerospike.vector.ClusterInfo/GetNodeId", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterInfoClient) GetClusterId(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ClusterId, error) {
	out := new(ClusterId)
	err := c.cc.Invoke(ctx, "/aerospike.vector.ClusterInfo/GetClusterId", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterInfoClient) GetClusterEndpoints(ctx context.Context, in *ClusterNodeEndpointsRequest, opts ...grpc.CallOption) (*ClusterNodeEndpoints, error) {
	out := new(ClusterNodeEndpoints)
	err := c.cc.Invoke(ctx, "/aerospike.vector.ClusterInfo/GetClusterEndpoints", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterInfoClient) GetOwnedPartitions(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ClusterPartitions, error) {
	out := new(ClusterPartitions)
	err := c.cc.Invoke(ctx, "/aerospike.vector.ClusterInfo/GetOwnedPartitions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClusterInfoServer is the server API for ClusterInfo service.
// All implementations must embed UnimplementedClusterInfoServer
// for forward compatibility
type ClusterInfoServer interface {
	// Get the internal cluster node-Id for this server.
	GetNodeId(context.Context, *emptypb.Empty) (*NodeId, error)
	// Get current cluster-Id for the current cluster.
	GetClusterId(context.Context, *emptypb.Empty) (*ClusterId, error)
	// Get the advertised/listening endpoints for all nodes in the cluster, given a listener name.
	GetClusterEndpoints(context.Context, *ClusterNodeEndpointsRequest) (*ClusterNodeEndpoints, error)
	// Get per-node owned partition list for all nodes in the cluster.
	GetOwnedPartitions(context.Context, *emptypb.Empty) (*ClusterPartitions, error)
	mustEmbedUnimplementedClusterInfoServer()
}

// UnimplementedClusterInfoServer must be embedded to have forward compatible implementations.
type UnimplementedClusterInfoServer struct {
}

func (UnimplementedClusterInfoServer) GetNodeId(context.Context, *emptypb.Empty) (*NodeId, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNodeId not implemented")
}
func (UnimplementedClusterInfoServer) GetClusterId(context.Context, *emptypb.Empty) (*ClusterId, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClusterId not implemented")
}
func (UnimplementedClusterInfoServer) GetClusterEndpoints(context.Context, *ClusterNodeEndpointsRequest) (*ClusterNodeEndpoints, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClusterEndpoints not implemented")
}
func (UnimplementedClusterInfoServer) GetOwnedPartitions(context.Context, *emptypb.Empty) (*ClusterPartitions, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOwnedPartitions not implemented")
}
func (UnimplementedClusterInfoServer) mustEmbedUnimplementedClusterInfoServer() {}

// UnsafeClusterInfoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClusterInfoServer will
// result in compilation errors.
type UnsafeClusterInfoServer interface {
	mustEmbedUnimplementedClusterInfoServer()
}

func RegisterClusterInfoServer(s grpc.ServiceRegistrar, srv ClusterInfoServer) {
	s.RegisterService(&ClusterInfo_ServiceDesc, srv)
}

func _ClusterInfo_GetNodeId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterInfoServer).GetNodeId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aerospike.vector.ClusterInfo/GetNodeId",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterInfoServer).GetNodeId(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClusterInfo_GetClusterId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterInfoServer).GetClusterId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aerospike.vector.ClusterInfo/GetClusterId",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterInfoServer).GetClusterId(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClusterInfo_GetClusterEndpoints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClusterNodeEndpointsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterInfoServer).GetClusterEndpoints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aerospike.vector.ClusterInfo/GetClusterEndpoints",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterInfoServer).GetClusterEndpoints(ctx, req.(*ClusterNodeEndpointsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClusterInfo_GetOwnedPartitions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterInfoServer).GetOwnedPartitions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aerospike.vector.ClusterInfo/GetOwnedPartitions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterInfoServer).GetOwnedPartitions(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// ClusterInfo_ServiceDesc is the grpc.ServiceDesc for ClusterInfo service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ClusterInfo_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "aerospike.vector.ClusterInfo",
	HandlerType: (*ClusterInfoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetNodeId",
			Handler:    _ClusterInfo_GetNodeId_Handler,
		},
		{
			MethodName: "GetClusterId",
			Handler:    _ClusterInfo_GetClusterId_Handler,
		},
		{
			MethodName: "GetClusterEndpoints",
			Handler:    _ClusterInfo_GetClusterEndpoints_Handler,
		},
		{
			MethodName: "GetOwnedPartitions",
			Handler:    _ClusterInfo_GetOwnedPartitions_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "vector-db.proto",
}
