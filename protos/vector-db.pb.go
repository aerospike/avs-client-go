// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.29.3
// source: vector-db.proto

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

// Node roles determine the tasks a node will undertake.
type NodeRole int32

const (
	// Node will handle ANN query related tasks.
	NodeRole_INDEX_QUERY NodeRole = 0
	// Node will handle vector index update/maintenance tasks including
	// index, insert, delete and update of vector records and index healing.
	// This value has been DEPRECATED. Use `NodeRole.INDEXER` instead.
	//
	// Deprecated: Marked as deprecated in vector-db.proto.
	NodeRole_INDEX_UPDATE NodeRole = 1
	// Node will handle key value read operations.
	NodeRole_KV_READ NodeRole = 2
	// Node will handle standalone index building.
	NodeRole_STANDALONE_INDEXER NodeRole = 3
	// Node will handle key value write operations.
	NodeRole_KV_WRITE NodeRole = 4
	// Node will handle vector index update/maintenance tasks including
	// index, insert, delete and update of vector records and index healing.
	NodeRole_INDEXER NodeRole = 5
)

// Enum value maps for NodeRole.
var (
	NodeRole_name = map[int32]string{
		0: "INDEX_QUERY",
		1: "INDEX_UPDATE",
		2: "KV_READ",
		3: "STANDALONE_INDEXER",
		4: "KV_WRITE",
		5: "INDEXER",
	}
	NodeRole_value = map[string]int32{
		"INDEX_QUERY":        0,
		"INDEX_UPDATE":       1,
		"KV_READ":            2,
		"STANDALONE_INDEXER": 3,
		"KV_WRITE":           4,
		"INDEXER":            5,
	}
)

func (x NodeRole) Enum() *NodeRole {
	p := new(NodeRole)
	*p = x
	return p
}

func (x NodeRole) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (NodeRole) Descriptor() protoreflect.EnumDescriptor {
	return file_vector_db_proto_enumTypes[0].Descriptor()
}

func (NodeRole) Type() protoreflect.EnumType {
	return &file_vector_db_proto_enumTypes[0]
}

func (x NodeRole) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use NodeRole.Descriptor instead.
func (NodeRole) EnumDescriptor() ([]byte, []int) {
	return file_vector_db_proto_rawDescGZIP(), []int{0}
}

// The about request message.
type AboutRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AboutRequest) Reset() {
	*x = AboutRequest{}
	mi := &file_vector_db_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AboutRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AboutRequest) ProtoMessage() {}

func (x *AboutRequest) ProtoReflect() protoreflect.Message {
	mi := &file_vector_db_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AboutRequest.ProtoReflect.Descriptor instead.
func (*AboutRequest) Descriptor() ([]byte, []int) {
	return file_vector_db_proto_rawDescGZIP(), []int{0}
}

// The about response message.
type AboutResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// AVS server version.
	Version string `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	// Assigned node roles.
	Roles []NodeRole `protobuf:"varint,2,rep,packed,name=roles,proto3,enum=aerospike.vector.NodeRole" json:"roles,omitempty"`
	// Self node id.
	SelfNodeId *NodeId `protobuf:"bytes,3,opt,name=selfNodeId,proto3" json:"selfNodeId,omitempty"`
}

func (x *AboutResponse) Reset() {
	*x = AboutResponse{}
	mi := &file_vector_db_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AboutResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AboutResponse) ProtoMessage() {}

func (x *AboutResponse) ProtoReflect() protoreflect.Message {
	mi := &file_vector_db_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AboutResponse.ProtoReflect.Descriptor instead.
func (*AboutResponse) Descriptor() ([]byte, []int) {
	return file_vector_db_proto_rawDescGZIP(), []int{1}
}

func (x *AboutResponse) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *AboutResponse) GetRoles() []NodeRole {
	if x != nil {
		return x.Roles
	}
	return nil
}

func (x *AboutResponse) GetSelfNodeId() *NodeId {
	if x != nil {
		return x.SelfNodeId
	}
	return nil
}

// Cluster Id
type ClusterId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ClusterId) Reset() {
	*x = ClusterId{}
	mi := &file_vector_db_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ClusterId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClusterId) ProtoMessage() {}

func (x *ClusterId) ProtoReflect() protoreflect.Message {
	mi := &file_vector_db_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClusterId.ProtoReflect.Descriptor instead.
func (*ClusterId) Descriptor() ([]byte, []int) {
	return file_vector_db_proto_rawDescGZIP(), []int{2}
}

func (x *ClusterId) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

// A server node Id
type NodeId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *NodeId) Reset() {
	*x = NodeId{}
	mi := &file_vector_db_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NodeId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeId) ProtoMessage() {}

func (x *NodeId) ProtoReflect() protoreflect.Message {
	mi := &file_vector_db_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeId.ProtoReflect.Descriptor instead.
func (*NodeId) Descriptor() ([]byte, []int) {
	return file_vector_db_proto_rawDescGZIP(), []int{3}
}

func (x *NodeId) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

// Server endpoint.
type ServerEndpoint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// IP address or DNS name.
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	// Listening port.
	Port uint32 `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	// Indicates if this is a TLS enabled port.
	IsTls bool `protobuf:"varint,3,opt,name=isTls,proto3" json:"isTls,omitempty"`
}

func (x *ServerEndpoint) Reset() {
	*x = ServerEndpoint{}
	mi := &file_vector_db_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ServerEndpoint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerEndpoint) ProtoMessage() {}

func (x *ServerEndpoint) ProtoReflect() protoreflect.Message {
	mi := &file_vector_db_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerEndpoint.ProtoReflect.Descriptor instead.
func (*ServerEndpoint) Descriptor() ([]byte, []int) {
	return file_vector_db_proto_rawDescGZIP(), []int{4}
}

func (x *ServerEndpoint) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *ServerEndpoint) GetPort() uint32 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *ServerEndpoint) GetIsTls() bool {
	if x != nil {
		return x.IsTls
	}
	return false
}

// Server endpoint.
type ServerEndpointList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// List of server endpoints.
	Endpoints []*ServerEndpoint `protobuf:"bytes,1,rep,name=endpoints,proto3" json:"endpoints,omitempty"`
}

func (x *ServerEndpointList) Reset() {
	*x = ServerEndpointList{}
	mi := &file_vector_db_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ServerEndpointList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerEndpointList) ProtoMessage() {}

func (x *ServerEndpointList) ProtoReflect() protoreflect.Message {
	mi := &file_vector_db_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerEndpointList.ProtoReflect.Descriptor instead.
func (*ServerEndpointList) Descriptor() ([]byte, []int) {
	return file_vector_db_proto_rawDescGZIP(), []int{5}
}

func (x *ServerEndpointList) GetEndpoints() []*ServerEndpoint {
	if x != nil {
		return x.Endpoints
	}
	return nil
}

// Cluster endpoint.
type ClusterNodeEndpoints struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Endpoints map[uint64]*ServerEndpointList `protobuf:"bytes,1,rep,name=endpoints,proto3" json:"endpoints,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *ClusterNodeEndpoints) Reset() {
	*x = ClusterNodeEndpoints{}
	mi := &file_vector_db_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ClusterNodeEndpoints) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClusterNodeEndpoints) ProtoMessage() {}

func (x *ClusterNodeEndpoints) ProtoReflect() protoreflect.Message {
	mi := &file_vector_db_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClusterNodeEndpoints.ProtoReflect.Descriptor instead.
func (*ClusterNodeEndpoints) Descriptor() ([]byte, []int) {
	return file_vector_db_proto_rawDescGZIP(), []int{6}
}

func (x *ClusterNodeEndpoints) GetEndpoints() map[uint64]*ServerEndpointList {
	if x != nil {
		return x.Endpoints
	}
	return nil
}

// Cluster endpoint.
type ClusterNodeEndpointsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Optional name of the listener.
	// If not specified the "default" listener endpoints are returned.
	ListenerName *string `protobuf:"bytes,1,opt,name=listenerName,proto3,oneof" json:"listenerName,omitempty"`
}

func (x *ClusterNodeEndpointsRequest) Reset() {
	*x = ClusterNodeEndpointsRequest{}
	mi := &file_vector_db_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ClusterNodeEndpointsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClusterNodeEndpointsRequest) ProtoMessage() {}

func (x *ClusterNodeEndpointsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_vector_db_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClusterNodeEndpointsRequest.ProtoReflect.Descriptor instead.
func (*ClusterNodeEndpointsRequest) Descriptor() ([]byte, []int) {
	return file_vector_db_proto_rawDescGZIP(), []int{7}
}

func (x *ClusterNodeEndpointsRequest) GetListenerName() string {
	if x != nil && x.ListenerName != nil {
		return *x.ListenerName
	}
	return ""
}

type ClusteringState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Indicates if this node is in a cluster.
	IsInCluster bool `protobuf:"varint,1,opt,name=isInCluster,proto3" json:"isInCluster,omitempty"`
	// Unique identifier for the current cluster epoch.
	// Valid only if the node is in cluster.
	ClusterId *ClusterId `protobuf:"bytes,2,opt,name=clusterId,proto3" json:"clusterId,omitempty"`
	// Current cluster members.
	Members []*NodeId `protobuf:"bytes,3,rep,name=members,proto3" json:"members,omitempty"`
}

func (x *ClusteringState) Reset() {
	*x = ClusteringState{}
	mi := &file_vector_db_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ClusteringState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClusteringState) ProtoMessage() {}

func (x *ClusteringState) ProtoReflect() protoreflect.Message {
	mi := &file_vector_db_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClusteringState.ProtoReflect.Descriptor instead.
func (*ClusteringState) Descriptor() ([]byte, []int) {
	return file_vector_db_proto_rawDescGZIP(), []int{8}
}

func (x *ClusteringState) GetIsInCluster() bool {
	if x != nil {
		return x.IsInCluster
	}
	return false
}

func (x *ClusteringState) GetClusterId() *ClusterId {
	if x != nil {
		return x.ClusterId
	}
	return nil
}

func (x *ClusteringState) GetMembers() []*NodeId {
	if x != nil {
		return x.Members
	}
	return nil
}

var File_vector_db_proto protoreflect.FileDescriptor

var file_vector_db_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2d, 0x64, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x10, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x0e, 0x0a, 0x0c, 0x41, 0x62, 0x6f, 0x75, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x22, 0x95, 0x01, 0x0a, 0x0d, 0x41, 0x62, 0x6f, 0x75, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x30, 0x0a, 0x05,
	0x72, 0x6f, 0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x61, 0x65,
	0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x4e,
	0x6f, 0x64, 0x65, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x05, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x12, 0x38,
	0x0a, 0x0a, 0x73, 0x65, 0x6c, 0x66, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x18, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x52, 0x0a, 0x73, 0x65,
	0x6c, 0x66, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x22, 0x1b, 0x0a, 0x09, 0x43, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x49, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x02, 0x69, 0x64, 0x22, 0x18, 0x0a, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x22,
	0x54, 0x0a, 0x0e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e,
	0x74, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x70,
	0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x69, 0x73, 0x54, 0x6c, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05,
	0x69, 0x73, 0x54, 0x6c, 0x73, 0x22, 0x54, 0x0a, 0x12, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x45,
	0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x3e, 0x0a, 0x09, 0x65,
	0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20,
	0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x52, 0x09, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x22, 0xcf, 0x01, 0x0a, 0x14,
	0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x45, 0x6e, 0x64, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x73, 0x12, 0x53, 0x0a, 0x09, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x35, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70,
	0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x2e,
	0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x09,
	0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x1a, 0x62, 0x0a, 0x0e, 0x45, 0x6e, 0x64,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x3a, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x61,
	0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e,
	0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x4c, 0x69,
	0x73, 0x74, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x57, 0x0a,
	0x1b, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x45, 0x6e, 0x64, 0x70,
	0x6f, 0x69, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x0c,
	0x6c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x00, 0x52, 0x0c, 0x6c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x4e, 0x61,
	0x6d, 0x65, 0x88, 0x01, 0x01, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x65, 0x6e,
	0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0xa2, 0x01, 0x0a, 0x0f, 0x43, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x69, 0x73,
	0x49, 0x6e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0b, 0x69, 0x73, 0x49, 0x6e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12, 0x39, 0x0a, 0x09,
	0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1b, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x64, 0x52, 0x09, 0x63, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x64, 0x12, 0x32, 0x0a, 0x07, 0x6d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73,
	0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x4e, 0x6f, 0x64, 0x65,
	0x49, 0x64, 0x52, 0x07, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x2a, 0x71, 0x0a, 0x08, 0x4e,
	0x6f, 0x64, 0x65, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x0f, 0x0a, 0x0b, 0x49, 0x4e, 0x44, 0x45, 0x58,
	0x5f, 0x51, 0x55, 0x45, 0x52, 0x59, 0x10, 0x00, 0x12, 0x14, 0x0a, 0x0c, 0x49, 0x4e, 0x44, 0x45,
	0x58, 0x5f, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x10, 0x01, 0x1a, 0x02, 0x08, 0x01, 0x12, 0x0b,
	0x0a, 0x07, 0x4b, 0x56, 0x5f, 0x52, 0x45, 0x41, 0x44, 0x10, 0x02, 0x12, 0x16, 0x0a, 0x12, 0x53,
	0x54, 0x41, 0x4e, 0x44, 0x41, 0x4c, 0x4f, 0x4e, 0x45, 0x5f, 0x49, 0x4e, 0x44, 0x45, 0x58, 0x45,
	0x52, 0x10, 0x03, 0x12, 0x0c, 0x0a, 0x08, 0x4b, 0x56, 0x5f, 0x57, 0x52, 0x49, 0x54, 0x45, 0x10,
	0x04, 0x12, 0x0b, 0x0a, 0x07, 0x49, 0x4e, 0x44, 0x45, 0x58, 0x45, 0x52, 0x10, 0x05, 0x32, 0x58,
	0x0a, 0x0c, 0x41, 0x62, 0x6f, 0x75, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x48,
	0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x1e, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b,
	0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x41, 0x62, 0x6f, 0x75, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b,
	0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x41, 0x62, 0x6f, 0x75, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x32, 0xdf, 0x02, 0x0a, 0x12, 0x43, 0x6c, 0x75,
	0x73, 0x74, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x3f, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x1a, 0x18, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65,
	0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x22, 0x00,
	0x12, 0x45, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1b, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73,
	0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x43, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x49, 0x64, 0x22, 0x00, 0x12, 0x51, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x43, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x21, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b,
	0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x69, 0x6e, 0x67, 0x53, 0x74, 0x61, 0x74, 0x65, 0x22, 0x00, 0x12, 0x6e, 0x0a, 0x13, 0x47, 0x65,
	0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x73, 0x12, 0x2d, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65,
	0x63, 0x74, 0x6f, 0x72, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65,
	0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x26, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x45,
	0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x22, 0x00, 0x42, 0x43, 0x0a, 0x21, 0x63, 0x6f,
	0x6d, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x50,
	0x01, 0x5a, 0x1c, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_vector_db_proto_rawDescOnce sync.Once
	file_vector_db_proto_rawDescData = file_vector_db_proto_rawDesc
)

func file_vector_db_proto_rawDescGZIP() []byte {
	file_vector_db_proto_rawDescOnce.Do(func() {
		file_vector_db_proto_rawDescData = protoimpl.X.CompressGZIP(file_vector_db_proto_rawDescData)
	})
	return file_vector_db_proto_rawDescData
}

var file_vector_db_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_vector_db_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_vector_db_proto_goTypes = []any{
	(NodeRole)(0),                       // 0: aerospike.vector.NodeRole
	(*AboutRequest)(nil),                // 1: aerospike.vector.AboutRequest
	(*AboutResponse)(nil),               // 2: aerospike.vector.AboutResponse
	(*ClusterId)(nil),                   // 3: aerospike.vector.ClusterId
	(*NodeId)(nil),                      // 4: aerospike.vector.NodeId
	(*ServerEndpoint)(nil),              // 5: aerospike.vector.ServerEndpoint
	(*ServerEndpointList)(nil),          // 6: aerospike.vector.ServerEndpointList
	(*ClusterNodeEndpoints)(nil),        // 7: aerospike.vector.ClusterNodeEndpoints
	(*ClusterNodeEndpointsRequest)(nil), // 8: aerospike.vector.ClusterNodeEndpointsRequest
	(*ClusteringState)(nil),             // 9: aerospike.vector.ClusteringState
	nil,                                 // 10: aerospike.vector.ClusterNodeEndpoints.EndpointsEntry
	(*emptypb.Empty)(nil),               // 11: google.protobuf.Empty
}
var file_vector_db_proto_depIdxs = []int32{
	0,  // 0: aerospike.vector.AboutResponse.roles:type_name -> aerospike.vector.NodeRole
	4,  // 1: aerospike.vector.AboutResponse.selfNodeId:type_name -> aerospike.vector.NodeId
	5,  // 2: aerospike.vector.ServerEndpointList.endpoints:type_name -> aerospike.vector.ServerEndpoint
	10, // 3: aerospike.vector.ClusterNodeEndpoints.endpoints:type_name -> aerospike.vector.ClusterNodeEndpoints.EndpointsEntry
	3,  // 4: aerospike.vector.ClusteringState.clusterId:type_name -> aerospike.vector.ClusterId
	4,  // 5: aerospike.vector.ClusteringState.members:type_name -> aerospike.vector.NodeId
	6,  // 6: aerospike.vector.ClusterNodeEndpoints.EndpointsEntry.value:type_name -> aerospike.vector.ServerEndpointList
	1,  // 7: aerospike.vector.AboutService.Get:input_type -> aerospike.vector.AboutRequest
	11, // 8: aerospike.vector.ClusterInfoService.GetNodeId:input_type -> google.protobuf.Empty
	11, // 9: aerospike.vector.ClusterInfoService.GetClusterId:input_type -> google.protobuf.Empty
	11, // 10: aerospike.vector.ClusterInfoService.GetClusteringState:input_type -> google.protobuf.Empty
	8,  // 11: aerospike.vector.ClusterInfoService.GetClusterEndpoints:input_type -> aerospike.vector.ClusterNodeEndpointsRequest
	2,  // 12: aerospike.vector.AboutService.Get:output_type -> aerospike.vector.AboutResponse
	4,  // 13: aerospike.vector.ClusterInfoService.GetNodeId:output_type -> aerospike.vector.NodeId
	3,  // 14: aerospike.vector.ClusterInfoService.GetClusterId:output_type -> aerospike.vector.ClusterId
	9,  // 15: aerospike.vector.ClusterInfoService.GetClusteringState:output_type -> aerospike.vector.ClusteringState
	7,  // 16: aerospike.vector.ClusterInfoService.GetClusterEndpoints:output_type -> aerospike.vector.ClusterNodeEndpoints
	12, // [12:17] is the sub-list for method output_type
	7,  // [7:12] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_vector_db_proto_init() }
func file_vector_db_proto_init() {
	if File_vector_db_proto != nil {
		return
	}
	file_vector_db_proto_msgTypes[7].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_vector_db_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_vector_db_proto_goTypes,
		DependencyIndexes: file_vector_db_proto_depIdxs,
		EnumInfos:         file_vector_db_proto_enumTypes,
		MessageInfos:      file_vector_db_proto_msgTypes,
	}.Build()
	File_vector_db_proto = out.File
	file_vector_db_proto_rawDesc = nil
	file_vector_db_proto_goTypes = nil
	file_vector_db_proto_depIdxs = nil
}
