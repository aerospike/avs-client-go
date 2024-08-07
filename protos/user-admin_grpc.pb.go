// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: user-admin.proto

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

// UserAdminServiceClient is the client API for UserAdminService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserAdminServiceClient interface {
	// Add a new user.
	AddUser(ctx context.Context, in *AddUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Update user credentials.
	UpdateCredentials(ctx context.Context, in *UpdateCredentialsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Drop a user.
	DropUser(ctx context.Context, in *DropUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Get details for a user.
	GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*User, error)
	// List users.
	ListUsers(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListUsersResponse, error)
	// Grant roles to a user.
	GrantRoles(ctx context.Context, in *GrantRolesRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Revoke roles from a user.
	RevokeRoles(ctx context.Context, in *RevokeRolesRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// List roles.
	ListRoles(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListRolesResponse, error)
}

type userAdminServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserAdminServiceClient(cc grpc.ClientConnInterface) UserAdminServiceClient {
	return &userAdminServiceClient{cc}
}

func (c *userAdminServiceClient) AddUser(ctx context.Context, in *AddUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/aerospike.vector.UserAdminService/AddUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userAdminServiceClient) UpdateCredentials(ctx context.Context, in *UpdateCredentialsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/aerospike.vector.UserAdminService/UpdateCredentials", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userAdminServiceClient) DropUser(ctx context.Context, in *DropUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/aerospike.vector.UserAdminService/DropUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userAdminServiceClient) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/aerospike.vector.UserAdminService/GetUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userAdminServiceClient) ListUsers(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListUsersResponse, error) {
	out := new(ListUsersResponse)
	err := c.cc.Invoke(ctx, "/aerospike.vector.UserAdminService/ListUsers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userAdminServiceClient) GrantRoles(ctx context.Context, in *GrantRolesRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/aerospike.vector.UserAdminService/GrantRoles", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userAdminServiceClient) RevokeRoles(ctx context.Context, in *RevokeRolesRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/aerospike.vector.UserAdminService/RevokeRoles", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userAdminServiceClient) ListRoles(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListRolesResponse, error) {
	out := new(ListRolesResponse)
	err := c.cc.Invoke(ctx, "/aerospike.vector.UserAdminService/ListRoles", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserAdminServiceServer is the server API for UserAdminService service.
// All implementations must embed UnimplementedUserAdminServiceServer
// for forward compatibility
type UserAdminServiceServer interface {
	// Add a new user.
	AddUser(context.Context, *AddUserRequest) (*emptypb.Empty, error)
	// Update user credentials.
	UpdateCredentials(context.Context, *UpdateCredentialsRequest) (*emptypb.Empty, error)
	// Drop a user.
	DropUser(context.Context, *DropUserRequest) (*emptypb.Empty, error)
	// Get details for a user.
	GetUser(context.Context, *GetUserRequest) (*User, error)
	// List users.
	ListUsers(context.Context, *emptypb.Empty) (*ListUsersResponse, error)
	// Grant roles to a user.
	GrantRoles(context.Context, *GrantRolesRequest) (*emptypb.Empty, error)
	// Revoke roles from a user.
	RevokeRoles(context.Context, *RevokeRolesRequest) (*emptypb.Empty, error)
	// List roles.
	ListRoles(context.Context, *emptypb.Empty) (*ListRolesResponse, error)
	mustEmbedUnimplementedUserAdminServiceServer()
}

// UnimplementedUserAdminServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserAdminServiceServer struct {
}

func (UnimplementedUserAdminServiceServer) AddUser(context.Context, *AddUserRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUser not implemented")
}
func (UnimplementedUserAdminServiceServer) UpdateCredentials(context.Context, *UpdateCredentialsRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCredentials not implemented")
}
func (UnimplementedUserAdminServiceServer) DropUser(context.Context, *DropUserRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DropUser not implemented")
}
func (UnimplementedUserAdminServiceServer) GetUser(context.Context, *GetUserRequest) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (UnimplementedUserAdminServiceServer) ListUsers(context.Context, *emptypb.Empty) (*ListUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUsers not implemented")
}
func (UnimplementedUserAdminServiceServer) GrantRoles(context.Context, *GrantRolesRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GrantRoles not implemented")
}
func (UnimplementedUserAdminServiceServer) RevokeRoles(context.Context, *RevokeRolesRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RevokeRoles not implemented")
}
func (UnimplementedUserAdminServiceServer) ListRoles(context.Context, *emptypb.Empty) (*ListRolesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRoles not implemented")
}
func (UnimplementedUserAdminServiceServer) mustEmbedUnimplementedUserAdminServiceServer() {}

// UnsafeUserAdminServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserAdminServiceServer will
// result in compilation errors.
type UnsafeUserAdminServiceServer interface {
	mustEmbedUnimplementedUserAdminServiceServer()
}

func RegisterUserAdminServiceServer(s grpc.ServiceRegistrar, srv UserAdminServiceServer) {
	s.RegisterService(&UserAdminService_ServiceDesc, srv)
}

func _UserAdminService_AddUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserAdminServiceServer).AddUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aerospike.vector.UserAdminService/AddUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserAdminServiceServer).AddUser(ctx, req.(*AddUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserAdminService_UpdateCredentials_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCredentialsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserAdminServiceServer).UpdateCredentials(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aerospike.vector.UserAdminService/UpdateCredentials",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserAdminServiceServer).UpdateCredentials(ctx, req.(*UpdateCredentialsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserAdminService_DropUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DropUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserAdminServiceServer).DropUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aerospike.vector.UserAdminService/DropUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserAdminServiceServer).DropUser(ctx, req.(*DropUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserAdminService_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserAdminServiceServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aerospike.vector.UserAdminService/GetUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserAdminServiceServer).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserAdminService_ListUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserAdminServiceServer).ListUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aerospike.vector.UserAdminService/ListUsers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserAdminServiceServer).ListUsers(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserAdminService_GrantRoles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GrantRolesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserAdminServiceServer).GrantRoles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aerospike.vector.UserAdminService/GrantRoles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserAdminServiceServer).GrantRoles(ctx, req.(*GrantRolesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserAdminService_RevokeRoles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RevokeRolesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserAdminServiceServer).RevokeRoles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aerospike.vector.UserAdminService/RevokeRoles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserAdminServiceServer).RevokeRoles(ctx, req.(*RevokeRolesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserAdminService_ListRoles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserAdminServiceServer).ListRoles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aerospike.vector.UserAdminService/ListRoles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserAdminServiceServer).ListRoles(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// UserAdminService_ServiceDesc is the grpc.ServiceDesc for UserAdminService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserAdminService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "aerospike.vector.UserAdminService",
	HandlerType: (*UserAdminServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddUser",
			Handler:    _UserAdminService_AddUser_Handler,
		},
		{
			MethodName: "UpdateCredentials",
			Handler:    _UserAdminService_UpdateCredentials_Handler,
		},
		{
			MethodName: "DropUser",
			Handler:    _UserAdminService_DropUser_Handler,
		},
		{
			MethodName: "GetUser",
			Handler:    _UserAdminService_GetUser_Handler,
		},
		{
			MethodName: "ListUsers",
			Handler:    _UserAdminService_ListUsers_Handler,
		},
		{
			MethodName: "GrantRoles",
			Handler:    _UserAdminService_GrantRoles_Handler,
		},
		{
			MethodName: "RevokeRoles",
			Handler:    _UserAdminService_RevokeRoles_Handler,
		},
		{
			MethodName: "ListRoles",
			Handler:    _UserAdminService_ListRoles_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user-admin.proto",
}
