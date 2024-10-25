// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.2
// source: transact.proto

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

// The type of write operation.
type WriteType int32

const (
	// Insert the record if it does not exist.
	// Update the record by replacing specified fields.
	WriteType_UPSERT WriteType = 0
	// If the record exists update the record by replacing specified fields.
	// Fails if the record does not exist.
	WriteType_UPDATE_ONLY WriteType = 1
	// Insert / create the record if it does not exists.
	// Fails if the record already exist.
	WriteType_INSERT_ONLY WriteType = 2
	// Replace all fields in the record if it exists, else create the
	// record if it does not exists.
	WriteType_REPLACE WriteType = 3
	// Replace all fields in the record if it exists.
	// Fails if the record does not exist.
	WriteType_REPLACE_ONLY WriteType = 4
)

// Enum value maps for WriteType.
var (
	WriteType_name = map[int32]string{
		0: "UPSERT",
		1: "UPDATE_ONLY",
		2: "INSERT_ONLY",
		3: "REPLACE",
		4: "REPLACE_ONLY",
	}
	WriteType_value = map[string]int32{
		"UPSERT":       0,
		"UPDATE_ONLY":  1,
		"INSERT_ONLY":  2,
		"REPLACE":      3,
		"REPLACE_ONLY": 4,
	}
)

func (x WriteType) Enum() *WriteType {
	p := new(WriteType)
	*p = x
	return p
}

func (x WriteType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (WriteType) Descriptor() protoreflect.EnumDescriptor {
	return file_transact_proto_enumTypes[0].Descriptor()
}

func (WriteType) Type() protoreflect.EnumType {
	return &file_transact_proto_enumTypes[0]
}

func (x WriteType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use WriteType.Descriptor instead.
func (WriteType) EnumDescriptor() ([]byte, []int) {
	return file_transact_proto_rawDescGZIP(), []int{0}
}

// The type of projection.
type ProjectionType int32

const (
	ProjectionType_ALL       ProjectionType = 0
	ProjectionType_NONE      ProjectionType = 1
	ProjectionType_SPECIFIED ProjectionType = 2
)

// Enum value maps for ProjectionType.
var (
	ProjectionType_name = map[int32]string{
		0: "ALL",
		1: "NONE",
		2: "SPECIFIED",
	}
	ProjectionType_value = map[string]int32{
		"ALL":       0,
		"NONE":      1,
		"SPECIFIED": 2,
	}
)

func (x ProjectionType) Enum() *ProjectionType {
	p := new(ProjectionType)
	*p = x
	return p
}

func (x ProjectionType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ProjectionType) Descriptor() protoreflect.EnumDescriptor {
	return file_transact_proto_enumTypes[1].Descriptor()
}

func (ProjectionType) Type() protoreflect.EnumType {
	return &file_transact_proto_enumTypes[1]
}

func (x ProjectionType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ProjectionType.Descriptor instead.
func (ProjectionType) EnumDescriptor() ([]byte, []int) {
	return file_transact_proto_rawDescGZIP(), []int{1}
}

// Put request to insert/update a record.
type PutRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The key for the record to insert/update
	Key *Key `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// The type of the put/write request. Defaults to UPSERT.
	WriteType *WriteType `protobuf:"varint,2,opt,name=writeType,proto3,enum=aerospike.vector.WriteType,oneof" json:"writeType,omitempty"`
	// The record fields.
	Fields []*Field `protobuf:"bytes,3,rep,name=fields,proto3" json:"fields,omitempty"`
	// Ignore the in-memory queue full error. These records would be written to
	// storage and later, the index healer would pick for indexing.
	IgnoreMemQueueFull bool `protobuf:"varint,4,opt,name=ignoreMemQueueFull,proto3" json:"ignoreMemQueueFull,omitempty"`
}

func (x *PutRequest) Reset() {
	*x = PutRequest{}
	mi := &file_transact_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PutRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutRequest) ProtoMessage() {}

func (x *PutRequest) ProtoReflect() protoreflect.Message {
	mi := &file_transact_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutRequest.ProtoReflect.Descriptor instead.
func (*PutRequest) Descriptor() ([]byte, []int) {
	return file_transact_proto_rawDescGZIP(), []int{0}
}

func (x *PutRequest) GetKey() *Key {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *PutRequest) GetWriteType() WriteType {
	if x != nil && x.WriteType != nil {
		return *x.WriteType
	}
	return WriteType_UPSERT
}

func (x *PutRequest) GetFields() []*Field {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *PutRequest) GetIgnoreMemQueueFull() bool {
	if x != nil {
		return x.IgnoreMemQueueFull
	}
	return false
}

// Get request to insert/update a record.
type GetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The key for the record to insert/update
	Key *Key `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// The field selector.
	Projection *ProjectionSpec `protobuf:"bytes,2,opt,name=projection,proto3" json:"projection,omitempty"`
}

func (x *GetRequest) Reset() {
	*x = GetRequest{}
	mi := &file_transact_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRequest) ProtoMessage() {}

func (x *GetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_transact_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRequest.ProtoReflect.Descriptor instead.
func (*GetRequest) Descriptor() ([]byte, []int) {
	return file_transact_proto_rawDescGZIP(), []int{1}
}

func (x *GetRequest) GetKey() *Key {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *GetRequest) GetProjection() *ProjectionSpec {
	if x != nil {
		return x.Projection
	}
	return nil
}

// Check if a record exists.
type ExistsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The key for the record check
	Key *Key `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *ExistsRequest) Reset() {
	*x = ExistsRequest{}
	mi := &file_transact_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ExistsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExistsRequest) ProtoMessage() {}

func (x *ExistsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_transact_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExistsRequest.ProtoReflect.Descriptor instead.
func (*ExistsRequest) Descriptor() ([]byte, []int) {
	return file_transact_proto_rawDescGZIP(), []int{2}
}

func (x *ExistsRequest) GetKey() *Key {
	if x != nil {
		return x.Key
	}
	return nil
}

// Delete request to delete a record.
type DeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The key for the record to delete
	Key *Key `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *DeleteRequest) Reset() {
	*x = DeleteRequest{}
	mi := &file_transact_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRequest) ProtoMessage() {}

func (x *DeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_transact_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRequest.ProtoReflect.Descriptor instead.
func (*DeleteRequest) Descriptor() ([]byte, []int) {
	return file_transact_proto_rawDescGZIP(), []int{3}
}

func (x *DeleteRequest) GetKey() *Key {
	if x != nil {
		return x.Key
	}
	return nil
}

// Request to check whether the given record is indexed for the specified
// index.
type IsIndexedRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The key of the aerospike record holding Vector.
	Key *Key `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// Index in which the vector is indexed.
	IndexId *IndexId `protobuf:"bytes,2,opt,name=indexId,proto3" json:"indexId,omitempty"`
}

func (x *IsIndexedRequest) Reset() {
	*x = IsIndexedRequest{}
	mi := &file_transact_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IsIndexedRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IsIndexedRequest) ProtoMessage() {}

func (x *IsIndexedRequest) ProtoReflect() protoreflect.Message {
	mi := &file_transact_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IsIndexedRequest.ProtoReflect.Descriptor instead.
func (*IsIndexedRequest) Descriptor() ([]byte, []int) {
	return file_transact_proto_rawDescGZIP(), []int{4}
}

func (x *IsIndexedRequest) GetKey() *Key {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *IsIndexedRequest) GetIndexId() *IndexId {
	if x != nil {
		return x.IndexId
	}
	return nil
}

// A projection filter.
type ProjectionFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The type of the selector. Defaults to ALL.
	Type *ProjectionType `protobuf:"varint,1,opt,name=type,proto3,enum=aerospike.vector.ProjectionType,oneof" json:"type,omitempty"`
	// Names of desired fields / selectors.
	Fields []string `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`
}

func (x *ProjectionFilter) Reset() {
	*x = ProjectionFilter{}
	mi := &file_transact_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProjectionFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProjectionFilter) ProtoMessage() {}

func (x *ProjectionFilter) ProtoReflect() protoreflect.Message {
	mi := &file_transact_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProjectionFilter.ProtoReflect.Descriptor instead.
func (*ProjectionFilter) Descriptor() ([]byte, []int) {
	return file_transact_proto_rawDescGZIP(), []int{5}
}

func (x *ProjectionFilter) GetType() ProjectionType {
	if x != nil && x.Type != nil {
		return *x.Type
	}
	return ProjectionType_ALL
}

func (x *ProjectionFilter) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

// Projection to select which fields are returned.
// A field is returned if it passes the include filter
// and does not pass the exclude filter.
type ProjectionSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The fields to include.
	Include *ProjectionFilter `protobuf:"bytes,1,opt,name=include,proto3" json:"include,omitempty"`
	// The fields to exclude.
	Exclude *ProjectionFilter `protobuf:"bytes,2,opt,name=exclude,proto3" json:"exclude,omitempty"`
}

func (x *ProjectionSpec) Reset() {
	*x = ProjectionSpec{}
	mi := &file_transact_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProjectionSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProjectionSpec) ProtoMessage() {}

func (x *ProjectionSpec) ProtoReflect() protoreflect.Message {
	mi := &file_transact_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProjectionSpec.ProtoReflect.Descriptor instead.
func (*ProjectionSpec) Descriptor() ([]byte, []int) {
	return file_transact_proto_rawDescGZIP(), []int{6}
}

func (x *ProjectionSpec) GetInclude() *ProjectionFilter {
	if x != nil {
		return x.Include
	}
	return nil
}

func (x *ProjectionSpec) GetExclude() *ProjectionFilter {
	if x != nil {
		return x.Exclude
	}
	return nil
}

type VectorSearchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The index identifier.
	Index *IndexId `protobuf:"bytes,1,opt,name=index,proto3" json:"index,omitempty"`
	// The query vector.
	QueryVector *Vector `protobuf:"bytes,2,opt,name=queryVector,proto3" json:"queryVector,omitempty"`
	// Maximum number of results to return.
	Limit uint32 `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
	// Projection to select fields.
	Projection *ProjectionSpec `protobuf:"bytes,4,opt,name=projection,proto3" json:"projection,omitempty"`
	// Optional parameters to tune the search.
	//
	// Types that are assignable to SearchParams:
	//
	//	*VectorSearchRequest_HnswSearchParams
	SearchParams isVectorSearchRequest_SearchParams `protobuf_oneof:"searchParams"`
}

func (x *VectorSearchRequest) Reset() {
	*x = VectorSearchRequest{}
	mi := &file_transact_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VectorSearchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VectorSearchRequest) ProtoMessage() {}

func (x *VectorSearchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_transact_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VectorSearchRequest.ProtoReflect.Descriptor instead.
func (*VectorSearchRequest) Descriptor() ([]byte, []int) {
	return file_transact_proto_rawDescGZIP(), []int{7}
}

func (x *VectorSearchRequest) GetIndex() *IndexId {
	if x != nil {
		return x.Index
	}
	return nil
}

func (x *VectorSearchRequest) GetQueryVector() *Vector {
	if x != nil {
		return x.QueryVector
	}
	return nil
}

func (x *VectorSearchRequest) GetLimit() uint32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *VectorSearchRequest) GetProjection() *ProjectionSpec {
	if x != nil {
		return x.Projection
	}
	return nil
}

func (m *VectorSearchRequest) GetSearchParams() isVectorSearchRequest_SearchParams {
	if m != nil {
		return m.SearchParams
	}
	return nil
}

func (x *VectorSearchRequest) GetHnswSearchParams() *HnswSearchParams {
	if x, ok := x.GetSearchParams().(*VectorSearchRequest_HnswSearchParams); ok {
		return x.HnswSearchParams
	}
	return nil
}

type isVectorSearchRequest_SearchParams interface {
	isVectorSearchRequest_SearchParams()
}

type VectorSearchRequest_HnswSearchParams struct {
	HnswSearchParams *HnswSearchParams `protobuf:"bytes,5,opt,name=hnswSearchParams,proto3,oneof"`
}

func (*VectorSearchRequest_HnswSearchParams) isVectorSearchRequest_SearchParams() {}

var File_transact_proto protoreflect.FileDescriptor

var file_transact_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x10, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x0b, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe4, 0x01, 0x0a,
	0x0a, 0x50, 0x75, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73,
	0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x4b, 0x65, 0x79, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x3e, 0x0a, 0x09, 0x77, 0x72, 0x69, 0x74, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70,
	0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x57, 0x72, 0x69, 0x74, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x48, 0x00, 0x52, 0x09, 0x77, 0x72, 0x69, 0x74, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x88, 0x01, 0x01, 0x12, 0x2f, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65,
	0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x52, 0x06, 0x66,
	0x69, 0x65, 0x6c, 0x64, 0x73, 0x12, 0x2e, 0x0a, 0x12, 0x69, 0x67, 0x6e, 0x6f, 0x72, 0x65, 0x4d,
	0x65, 0x6d, 0x51, 0x75, 0x65, 0x75, 0x65, 0x46, 0x75, 0x6c, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x12, 0x69, 0x67, 0x6e, 0x6f, 0x72, 0x65, 0x4d, 0x65, 0x6d, 0x51, 0x75, 0x65, 0x75,
	0x65, 0x46, 0x75, 0x6c, 0x6c, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x77, 0x72, 0x69, 0x74, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x22, 0x77, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x27, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15,
	0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x2e, 0x4b, 0x65, 0x79, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x40, 0x0a, 0x0a, 0x70, 0x72,
	0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20,
	0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x70, 0x65, 0x63,
	0x52, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x38, 0x0a, 0x0d,
	0x45, 0x78, 0x69, 0x73, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x61, 0x65, 0x72,
	0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x4b, 0x65,
	0x79, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22, 0x38, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65,
	0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x4b, 0x65, 0x79, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x22, 0x70, 0x0a, 0x10, 0x49, 0x73, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x64, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x15, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65,
	0x63, 0x74, 0x6f, 0x72, 0x2e, 0x4b, 0x65, 0x79, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x33, 0x0a,
	0x07, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19,
	0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x49, 0x64, 0x52, 0x07, 0x69, 0x6e, 0x64, 0x65, 0x78,
	0x49, 0x64, 0x22, 0x6e, 0x0a, 0x10, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x39, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x20, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65,
	0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x48, 0x00, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x88, 0x01,
	0x01, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x74, 0x79,
	0x70, 0x65, 0x22, 0x8c, 0x01, 0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x53, 0x70, 0x65, 0x63, 0x12, 0x3c, 0x0a, 0x07, 0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69,
	0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x07, 0x69, 0x6e, 0x63, 0x6c,
	0x75, 0x64, 0x65, 0x12, 0x3c, 0x0a, 0x07, 0x65, 0x78, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65,
	0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x07, 0x65, 0x78, 0x63, 0x6c, 0x75, 0x64,
	0x65, 0x22, 0xbc, 0x02, 0x0a, 0x13, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x53, 0x65, 0x61, 0x72,
	0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2f, 0x0a, 0x05, 0x69, 0x6e, 0x64,
	0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73,
	0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x49, 0x6e, 0x64, 0x65,
	0x78, 0x49, 0x64, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x3a, 0x0a, 0x0b, 0x71, 0x75,
	0x65, 0x72, 0x79, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x18, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x2e, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x0b, 0x71, 0x75, 0x65, 0x72, 0x79,
	0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x40, 0x0a, 0x0a,
	0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x20, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x70,
	0x65, 0x63, 0x52, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x50,
	0x0a, 0x10, 0x68, 0x6e, 0x73, 0x77, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x50, 0x61, 0x72, 0x61,
	0x6d, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73,
	0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x48, 0x6e, 0x73, 0x77,
	0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x48, 0x00, 0x52, 0x10,
	0x68, 0x6e, 0x73, 0x77, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73,
	0x42, 0x0e, 0x0a, 0x0c, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73,
	0x2a, 0x58, 0x0a, 0x09, 0x57, 0x72, 0x69, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0a, 0x0a,
	0x06, 0x55, 0x50, 0x53, 0x45, 0x52, 0x54, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x50, 0x44,
	0x41, 0x54, 0x45, 0x5f, 0x4f, 0x4e, 0x4c, 0x59, 0x10, 0x01, 0x12, 0x0f, 0x0a, 0x0b, 0x49, 0x4e,
	0x53, 0x45, 0x52, 0x54, 0x5f, 0x4f, 0x4e, 0x4c, 0x59, 0x10, 0x02, 0x12, 0x0b, 0x0a, 0x07, 0x52,
	0x45, 0x50, 0x4c, 0x41, 0x43, 0x45, 0x10, 0x03, 0x12, 0x10, 0x0a, 0x0c, 0x52, 0x45, 0x50, 0x4c,
	0x41, 0x43, 0x45, 0x5f, 0x4f, 0x4e, 0x4c, 0x59, 0x10, 0x04, 0x2a, 0x32, 0x0a, 0x0e, 0x50, 0x72,
	0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x07, 0x0a, 0x03,
	0x41, 0x4c, 0x4c, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x01, 0x12,
	0x0d, 0x0a, 0x09, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x02, 0x32, 0xc3,
	0x03, 0x0a, 0x0f, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x3d, 0x0a, 0x03, 0x50, 0x75, 0x74, 0x12, 0x1c, 0x2e, 0x61, 0x65, 0x72, 0x6f,
	0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x50, 0x75, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22,
	0x00, 0x12, 0x3f, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x1c, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73,
	0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69,
	0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x22, 0x00, 0x12, 0x43, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x1f, 0x2e, 0x61,
	0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x46, 0x0a, 0x06, 0x45, 0x78, 0x69, 0x73, 0x74,
	0x73, 0x12, 0x1f, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65,
	0x63, 0x74, 0x6f, 0x72, 0x2e, 0x45, 0x78, 0x69, 0x73, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x19, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x65, 0x61, 0x6e, 0x22, 0x00, 0x12,
	0x4c, 0x0a, 0x09, 0x49, 0x73, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x64, 0x12, 0x22, 0x2e, 0x61,
	0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e,
	0x49, 0x73, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x19, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x65, 0x61, 0x6e, 0x22, 0x00, 0x12, 0x55, 0x0a,
	0x0c, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x12, 0x25, 0x2e,
	0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x2e, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65,
	0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x4e, 0x65, 0x69, 0x67, 0x68, 0x62, 0x6f, 0x72,
	0x22, 0x00, 0x30, 0x01, 0x42, 0x43, 0x0a, 0x21, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x65, 0x72, 0x6f,
	0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x1c, 0x61, 0x65, 0x72,
	0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x76, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_transact_proto_rawDescOnce sync.Once
	file_transact_proto_rawDescData = file_transact_proto_rawDesc
)

func file_transact_proto_rawDescGZIP() []byte {
	file_transact_proto_rawDescOnce.Do(func() {
		file_transact_proto_rawDescData = protoimpl.X.CompressGZIP(file_transact_proto_rawDescData)
	})
	return file_transact_proto_rawDescData
}

var file_transact_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_transact_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_transact_proto_goTypes = []any{
	(WriteType)(0),              // 0: aerospike.vector.WriteType
	(ProjectionType)(0),         // 1: aerospike.vector.ProjectionType
	(*PutRequest)(nil),          // 2: aerospike.vector.PutRequest
	(*GetRequest)(nil),          // 3: aerospike.vector.GetRequest
	(*ExistsRequest)(nil),       // 4: aerospike.vector.ExistsRequest
	(*DeleteRequest)(nil),       // 5: aerospike.vector.DeleteRequest
	(*IsIndexedRequest)(nil),    // 6: aerospike.vector.IsIndexedRequest
	(*ProjectionFilter)(nil),    // 7: aerospike.vector.ProjectionFilter
	(*ProjectionSpec)(nil),      // 8: aerospike.vector.ProjectionSpec
	(*VectorSearchRequest)(nil), // 9: aerospike.vector.VectorSearchRequest
	(*Key)(nil),                 // 10: aerospike.vector.Key
	(*Field)(nil),               // 11: aerospike.vector.Field
	(*IndexId)(nil),             // 12: aerospike.vector.IndexId
	(*Vector)(nil),              // 13: aerospike.vector.Vector
	(*HnswSearchParams)(nil),    // 14: aerospike.vector.HnswSearchParams
	(*emptypb.Empty)(nil),       // 15: google.protobuf.Empty
	(*Record)(nil),              // 16: aerospike.vector.Record
	(*Boolean)(nil),             // 17: aerospike.vector.Boolean
	(*Neighbor)(nil),            // 18: aerospike.vector.Neighbor
}
var file_transact_proto_depIdxs = []int32{
	10, // 0: aerospike.vector.PutRequest.key:type_name -> aerospike.vector.Key
	0,  // 1: aerospike.vector.PutRequest.writeType:type_name -> aerospike.vector.WriteType
	11, // 2: aerospike.vector.PutRequest.fields:type_name -> aerospike.vector.Field
	10, // 3: aerospike.vector.GetRequest.key:type_name -> aerospike.vector.Key
	8,  // 4: aerospike.vector.GetRequest.projection:type_name -> aerospike.vector.ProjectionSpec
	10, // 5: aerospike.vector.ExistsRequest.key:type_name -> aerospike.vector.Key
	10, // 6: aerospike.vector.DeleteRequest.key:type_name -> aerospike.vector.Key
	10, // 7: aerospike.vector.IsIndexedRequest.key:type_name -> aerospike.vector.Key
	12, // 8: aerospike.vector.IsIndexedRequest.indexId:type_name -> aerospike.vector.IndexId
	1,  // 9: aerospike.vector.ProjectionFilter.type:type_name -> aerospike.vector.ProjectionType
	7,  // 10: aerospike.vector.ProjectionSpec.include:type_name -> aerospike.vector.ProjectionFilter
	7,  // 11: aerospike.vector.ProjectionSpec.exclude:type_name -> aerospike.vector.ProjectionFilter
	12, // 12: aerospike.vector.VectorSearchRequest.index:type_name -> aerospike.vector.IndexId
	13, // 13: aerospike.vector.VectorSearchRequest.queryVector:type_name -> aerospike.vector.Vector
	8,  // 14: aerospike.vector.VectorSearchRequest.projection:type_name -> aerospike.vector.ProjectionSpec
	14, // 15: aerospike.vector.VectorSearchRequest.hnswSearchParams:type_name -> aerospike.vector.HnswSearchParams
	2,  // 16: aerospike.vector.TransactService.Put:input_type -> aerospike.vector.PutRequest
	3,  // 17: aerospike.vector.TransactService.Get:input_type -> aerospike.vector.GetRequest
	5,  // 18: aerospike.vector.TransactService.Delete:input_type -> aerospike.vector.DeleteRequest
	4,  // 19: aerospike.vector.TransactService.Exists:input_type -> aerospike.vector.ExistsRequest
	6,  // 20: aerospike.vector.TransactService.IsIndexed:input_type -> aerospike.vector.IsIndexedRequest
	9,  // 21: aerospike.vector.TransactService.VectorSearch:input_type -> aerospike.vector.VectorSearchRequest
	15, // 22: aerospike.vector.TransactService.Put:output_type -> google.protobuf.Empty
	16, // 23: aerospike.vector.TransactService.Get:output_type -> aerospike.vector.Record
	15, // 24: aerospike.vector.TransactService.Delete:output_type -> google.protobuf.Empty
	17, // 25: aerospike.vector.TransactService.Exists:output_type -> aerospike.vector.Boolean
	17, // 26: aerospike.vector.TransactService.IsIndexed:output_type -> aerospike.vector.Boolean
	18, // 27: aerospike.vector.TransactService.VectorSearch:output_type -> aerospike.vector.Neighbor
	22, // [22:28] is the sub-list for method output_type
	16, // [16:22] is the sub-list for method input_type
	16, // [16:16] is the sub-list for extension type_name
	16, // [16:16] is the sub-list for extension extendee
	0,  // [0:16] is the sub-list for field type_name
}

func init() { file_transact_proto_init() }
func file_transact_proto_init() {
	if File_transact_proto != nil {
		return
	}
	file_types_proto_init()
	file_transact_proto_msgTypes[0].OneofWrappers = []any{}
	file_transact_proto_msgTypes[5].OneofWrappers = []any{}
	file_transact_proto_msgTypes[7].OneofWrappers = []any{
		(*VectorSearchRequest_HnswSearchParams)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_transact_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_transact_proto_goTypes,
		DependencyIndexes: file_transact_proto_depIdxs,
		EnumInfos:         file_transact_proto_enumTypes,
		MessageInfos:      file_transact_proto_msgTypes,
	}.Build()
	File_transact_proto = out.File
	file_transact_proto_rawDesc = nil
	file_transact_proto_goTypes = nil
	file_transact_proto_depIdxs = nil
}
