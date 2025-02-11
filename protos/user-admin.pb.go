// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.29.3
// source: user-admin.proto

package protos

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Add a new user
type AddUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Credentials *Credentials `protobuf:"bytes,1,opt,name=credentials,proto3" json:"credentials,omitempty"`
	// Granted roles
	Roles []string `protobuf:"bytes,2,rep,name=roles,proto3" json:"roles,omitempty"`
}

func (x *AddUserRequest) Reset() {
	*x = AddUserRequest{}
	mi := &file_user_admin_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddUserRequest) ProtoMessage() {}

func (x *AddUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_admin_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddUserRequest.ProtoReflect.Descriptor instead.
func (*AddUserRequest) Descriptor() ([]byte, []int) {
	return file_user_admin_proto_rawDescGZIP(), []int{0}
}

func (x *AddUserRequest) GetCredentials() *Credentials {
	if x != nil {
		return x.Credentials
	}
	return nil
}

func (x *AddUserRequest) GetRoles() []string {
	if x != nil {
		return x.Roles
	}
	return nil
}

// Update user credentials
type UpdateCredentialsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Credentials *Credentials `protobuf:"bytes,1,opt,name=credentials,proto3" json:"credentials,omitempty"`
}

func (x *UpdateCredentialsRequest) Reset() {
	*x = UpdateCredentialsRequest{}
	mi := &file_user_admin_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateCredentialsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateCredentialsRequest) ProtoMessage() {}

func (x *UpdateCredentialsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_admin_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateCredentialsRequest.ProtoReflect.Descriptor instead.
func (*UpdateCredentialsRequest) Descriptor() ([]byte, []int) {
	return file_user_admin_proto_rawDescGZIP(), []int{1}
}

func (x *UpdateCredentialsRequest) GetCredentials() *Credentials {
	if x != nil {
		return x.Credentials
	}
	return nil
}

// Update user credentials
type DropUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
}

func (x *DropUserRequest) Reset() {
	*x = DropUserRequest{}
	mi := &file_user_admin_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DropUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DropUserRequest) ProtoMessage() {}

func (x *DropUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_admin_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DropUserRequest.ProtoReflect.Descriptor instead.
func (*DropUserRequest) Descriptor() ([]byte, []int) {
	return file_user_admin_proto_rawDescGZIP(), []int{2}
}

func (x *DropUserRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

// Update user credentials
type GetUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
}

func (x *GetUserRequest) Reset() {
	*x = GetUserRequest{}
	mi := &file_user_admin_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserRequest) ProtoMessage() {}

func (x *GetUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_admin_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserRequest.ProtoReflect.Descriptor instead.
func (*GetUserRequest) Descriptor() ([]byte, []int) {
	return file_user_admin_proto_rawDescGZIP(), []int{3}
}

func (x *GetUserRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

// Grant roles request
type GrantRolesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string   `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Roles    []string `protobuf:"bytes,2,rep,name=roles,proto3" json:"roles,omitempty"`
}

func (x *GrantRolesRequest) Reset() {
	*x = GrantRolesRequest{}
	mi := &file_user_admin_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GrantRolesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GrantRolesRequest) ProtoMessage() {}

func (x *GrantRolesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_admin_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GrantRolesRequest.ProtoReflect.Descriptor instead.
func (*GrantRolesRequest) Descriptor() ([]byte, []int) {
	return file_user_admin_proto_rawDescGZIP(), []int{4}
}

func (x *GrantRolesRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *GrantRolesRequest) GetRoles() []string {
	if x != nil {
		return x.Roles
	}
	return nil
}

// Revoke roles request
type RevokeRolesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string   `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Roles    []string `protobuf:"bytes,2,rep,name=roles,proto3" json:"roles,omitempty"`
}

func (x *RevokeRolesRequest) Reset() {
	*x = RevokeRolesRequest{}
	mi := &file_user_admin_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RevokeRolesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RevokeRolesRequest) ProtoMessage() {}

func (x *RevokeRolesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_admin_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RevokeRolesRequest.ProtoReflect.Descriptor instead.
func (*RevokeRolesRequest) Descriptor() ([]byte, []int) {
	return file_user_admin_proto_rawDescGZIP(), []int{5}
}

func (x *RevokeRolesRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *RevokeRolesRequest) GetRoles() []string {
	if x != nil {
		return x.Roles
	}
	return nil
}

// A list of roles.
type ListRolesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// List of roles.
	Roles []*Role `protobuf:"bytes,1,rep,name=roles,proto3" json:"roles,omitempty"`
}

func (x *ListRolesResponse) Reset() {
	*x = ListRolesResponse{}
	mi := &file_user_admin_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListRolesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRolesResponse) ProtoMessage() {}

func (x *ListRolesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_admin_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRolesResponse.ProtoReflect.Descriptor instead.
func (*ListRolesResponse) Descriptor() ([]byte, []int) {
	return file_user_admin_proto_rawDescGZIP(), []int{6}
}

func (x *ListRolesResponse) GetRoles() []*Role {
	if x != nil {
		return x.Roles
	}
	return nil
}

// A list of users.
type ListUsersResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// List of users.
	Users []*User `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`
}

func (x *ListUsersResponse) Reset() {
	*x = ListUsersResponse{}
	mi := &file_user_admin_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListUsersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListUsersResponse) ProtoMessage() {}

func (x *ListUsersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_admin_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListUsersResponse.ProtoReflect.Descriptor instead.
func (*ListUsersResponse) Descriptor() ([]byte, []int) {
	return file_user_admin_proto_rawDescGZIP(), []int{7}
}

func (x *ListUsersResponse) GetUsers() []*User {
	if x != nil {
		return x.Users
	}
	return nil
}

var File_user_admin_proto protoreflect.FileDescriptor

var file_user_admin_proto_rawDesc = []byte{
	0x0a, 0x10, 0x75, 0x73, 0x65, 0x72, 0x2d, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x10, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65,
	0x63, 0x74, 0x6f, 0x72, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x0b, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x67,
	0x0a, 0x0e, 0x41, 0x64, 0x64, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x3f, 0x0a, 0x0b, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b,
	0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x61, 0x6c, 0x73, 0x52, 0x0b, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c,
	0x73, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x05, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x22, 0x5b, 0x0a, 0x18, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x3f, 0x0a, 0x0b, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61,
	0x6c, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73,
	0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x43, 0x72, 0x65, 0x64,
	0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x52, 0x0b, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x61, 0x6c, 0x73, 0x22, 0x2d, 0x0a, 0x0f, 0x44, 0x72, 0x6f, 0x70, 0x55, 0x73, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e,
	0x61, 0x6d, 0x65, 0x22, 0x2c, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x22, 0x45, 0x0a, 0x11, 0x47, 0x72, 0x61, 0x6e, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x05, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x22, 0x46, 0x0a, 0x12, 0x52, 0x65, 0x76, 0x6f,
	0x6b, 0x65, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a,
	0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f,
	0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x72, 0x6f, 0x6c, 0x65, 0x73,
	0x22, 0x41, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a, 0x05, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65,
	0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x05, 0x72, 0x6f,
	0x6c, 0x65, 0x73, 0x22, 0x41, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70,
	0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52,
	0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x32, 0xf8, 0x04, 0x0a, 0x10, 0x55, 0x73, 0x65, 0x72, 0x41,
	0x64, 0x6d, 0x69, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x45, 0x0a, 0x07, 0x41,
	0x64, 0x64, 0x55, 0x73, 0x65, 0x72, 0x12, 0x20, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69,
	0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x41, 0x64, 0x64, 0x55, 0x73, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x00, 0x12, 0x59, 0x0a, 0x11, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x72, 0x65, 0x64,
	0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x12, 0x2a, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70,
	0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x47, 0x0a,
	0x08, 0x44, 0x72, 0x6f, 0x70, 0x55, 0x73, 0x65, 0x72, 0x12, 0x21, 0x2e, 0x61, 0x65, 0x72, 0x6f,
	0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x44, 0x72, 0x6f,
	0x70, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x45, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x12, 0x20, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65,
	0x63, 0x74, 0x6f, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e,
	0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x22, 0x00, 0x12, 0x4a, 0x0a,
	0x09, 0x4c, 0x69, 0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x1a, 0x23, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4b, 0x0a, 0x0a, 0x47, 0x72, 0x61,
	0x6e, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x12, 0x23, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70,
	0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x47, 0x72, 0x61, 0x6e, 0x74,
	0x52, 0x6f, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x4d, 0x0a, 0x0b, 0x52, 0x65, 0x76, 0x6f, 0x6b, 0x65,
	0x52, 0x6f, 0x6c, 0x65, 0x73, 0x12, 0x24, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b,
	0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x52, 0x65, 0x76, 0x6f, 0x6b, 0x65, 0x52,
	0x6f, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x4a, 0x0a, 0x09, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x6f, 0x6c,
	0x65, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x23, 0x2e, 0x61, 0x65, 0x72,
	0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x42, 0x43, 0x0a, 0x21, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69,
	0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x1c, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70,
	0x69, 0x6b, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_user_admin_proto_rawDescOnce sync.Once
	file_user_admin_proto_rawDescData = file_user_admin_proto_rawDesc
)

func file_user_admin_proto_rawDescGZIP() []byte {
	file_user_admin_proto_rawDescOnce.Do(func() {
		file_user_admin_proto_rawDescData = protoimpl.X.CompressGZIP(file_user_admin_proto_rawDescData)
	})
	return file_user_admin_proto_rawDescData
}

var file_user_admin_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_user_admin_proto_goTypes = []any{
	(*AddUserRequest)(nil),           // 0: aerospike.vector.AddUserRequest
	(*UpdateCredentialsRequest)(nil), // 1: aerospike.vector.UpdateCredentialsRequest
	(*DropUserRequest)(nil),          // 2: aerospike.vector.DropUserRequest
	(*GetUserRequest)(nil),           // 3: aerospike.vector.GetUserRequest
	(*GrantRolesRequest)(nil),        // 4: aerospike.vector.GrantRolesRequest
	(*RevokeRolesRequest)(nil),       // 5: aerospike.vector.RevokeRolesRequest
	(*ListRolesResponse)(nil),        // 6: aerospike.vector.ListRolesResponse
	(*ListUsersResponse)(nil),        // 7: aerospike.vector.ListUsersResponse
	(*Credentials)(nil),              // 8: aerospike.vector.Credentials
	(*Role)(nil),                     // 9: aerospike.vector.Role
	(*User)(nil),                     // 10: aerospike.vector.User
	(*emptypb.Empty)(nil),            // 11: google.protobuf.Empty
}
var file_user_admin_proto_depIdxs = []int32{
	8,  // 0: aerospike.vector.AddUserRequest.credentials:type_name -> aerospike.vector.Credentials
	8,  // 1: aerospike.vector.UpdateCredentialsRequest.credentials:type_name -> aerospike.vector.Credentials
	9,  // 2: aerospike.vector.ListRolesResponse.roles:type_name -> aerospike.vector.Role
	10, // 3: aerospike.vector.ListUsersResponse.users:type_name -> aerospike.vector.User
	0,  // 4: aerospike.vector.UserAdminService.AddUser:input_type -> aerospike.vector.AddUserRequest
	1,  // 5: aerospike.vector.UserAdminService.UpdateCredentials:input_type -> aerospike.vector.UpdateCredentialsRequest
	2,  // 6: aerospike.vector.UserAdminService.DropUser:input_type -> aerospike.vector.DropUserRequest
	3,  // 7: aerospike.vector.UserAdminService.GetUser:input_type -> aerospike.vector.GetUserRequest
	11, // 8: aerospike.vector.UserAdminService.ListUsers:input_type -> google.protobuf.Empty
	4,  // 9: aerospike.vector.UserAdminService.GrantRoles:input_type -> aerospike.vector.GrantRolesRequest
	5,  // 10: aerospike.vector.UserAdminService.RevokeRoles:input_type -> aerospike.vector.RevokeRolesRequest
	11, // 11: aerospike.vector.UserAdminService.ListRoles:input_type -> google.protobuf.Empty
	11, // 12: aerospike.vector.UserAdminService.AddUser:output_type -> google.protobuf.Empty
	11, // 13: aerospike.vector.UserAdminService.UpdateCredentials:output_type -> google.protobuf.Empty
	11, // 14: aerospike.vector.UserAdminService.DropUser:output_type -> google.protobuf.Empty
	10, // 15: aerospike.vector.UserAdminService.GetUser:output_type -> aerospike.vector.User
	7,  // 16: aerospike.vector.UserAdminService.ListUsers:output_type -> aerospike.vector.ListUsersResponse
	11, // 17: aerospike.vector.UserAdminService.GrantRoles:output_type -> google.protobuf.Empty
	11, // 18: aerospike.vector.UserAdminService.RevokeRoles:output_type -> google.protobuf.Empty
	6,  // 19: aerospike.vector.UserAdminService.ListRoles:output_type -> aerospike.vector.ListRolesResponse
	12, // [12:20] is the sub-list for method output_type
	4,  // [4:12] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_user_admin_proto_init() }
func file_user_admin_proto_init() {
	if File_user_admin_proto != nil {
		return
	}
	file_types_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_user_admin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_user_admin_proto_goTypes,
		DependencyIndexes: file_user_admin_proto_depIdxs,
		MessageInfos:      file_user_admin_proto_msgTypes,
	}.Build()
	File_user_admin_proto = out.File
	file_user_admin_proto_rawDesc = nil
	file_user_admin_proto_goTypes = nil
	file_user_admin_proto_depIdxs = nil
}
