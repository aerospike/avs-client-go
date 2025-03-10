// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.29.3
// source: index.proto

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

type StandaloneIndexMetrics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The id of the index to insert into.
	IndexId *IndexId `protobuf:"bytes,1,opt,name=indexId,proto3" json:"indexId,omitempty"`
	// Current state of the standalone index.
	State StandaloneIndexState `protobuf:"varint,2,opt,name=state,proto3,enum=aerospike.vector.StandaloneIndexState" json:"state,omitempty"`
	// The count of scanned vector records for indexing.
	ScannedVectorRecordCount uint64 `protobuf:"varint,3,opt,name=scannedVectorRecordCount,proto3" json:"scannedVectorRecordCount,omitempty"`
	// The count of indexed vector records.
	IndexedVectorRecordCount uint64 `protobuf:"varint,4,opt,name=indexedVectorRecordCount,proto3" json:"indexedVectorRecordCount,omitempty"`
}

func (x *StandaloneIndexMetrics) Reset() {
	*x = StandaloneIndexMetrics{}
	mi := &file_index_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StandaloneIndexMetrics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StandaloneIndexMetrics) ProtoMessage() {}

func (x *StandaloneIndexMetrics) ProtoReflect() protoreflect.Message {
	mi := &file_index_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StandaloneIndexMetrics.ProtoReflect.Descriptor instead.
func (*StandaloneIndexMetrics) Descriptor() ([]byte, []int) {
	return file_index_proto_rawDescGZIP(), []int{0}
}

func (x *StandaloneIndexMetrics) GetIndexId() *IndexId {
	if x != nil {
		return x.IndexId
	}
	return nil
}

func (x *StandaloneIndexMetrics) GetState() StandaloneIndexState {
	if x != nil {
		return x.State
	}
	return StandaloneIndexState_CREATING
}

func (x *StandaloneIndexMetrics) GetScannedVectorRecordCount() uint64 {
	if x != nil {
		return x.ScannedVectorRecordCount
	}
	return 0
}

func (x *StandaloneIndexMetrics) GetIndexedVectorRecordCount() uint64 {
	if x != nil {
		return x.IndexedVectorRecordCount
	}
	return 0
}

type IndexStatusResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StandaloneIndexMetrics *StandaloneIndexMetrics `protobuf:"bytes,1,opt,name=standaloneIndexMetrics,proto3,oneof" json:"standaloneIndexMetrics,omitempty"`
	// Number of unmerged index records.
	UnmergedRecordCount int64 `protobuf:"varint,2,opt,name=unmergedRecordCount,proto3" json:"unmergedRecordCount,omitempty"`
	// Number of vector records indexed (0 in case healer has not yet run).
	IndexHealerVectorRecordsIndexed int64 `protobuf:"varint,3,opt,name=indexHealerVectorRecordsIndexed,proto3" json:"indexHealerVectorRecordsIndexed,omitempty"`
	// Number of vertices in the main index (0 in case healer has not yet run).
	IndexHealerVerticesValid int64 `protobuf:"varint,4,opt,name=indexHealerVerticesValid,proto3" json:"indexHealerVerticesValid,omitempty"`
	// Status of the index.
	Status Status `protobuf:"varint,5,opt,name=status,proto3,enum=aerospike.vector.Status" json:"status,omitempty"`
}

func (x *IndexStatusResponse) Reset() {
	*x = IndexStatusResponse{}
	mi := &file_index_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IndexStatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IndexStatusResponse) ProtoMessage() {}

func (x *IndexStatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_index_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IndexStatusResponse.ProtoReflect.Descriptor instead.
func (*IndexStatusResponse) Descriptor() ([]byte, []int) {
	return file_index_proto_rawDescGZIP(), []int{1}
}

func (x *IndexStatusResponse) GetStandaloneIndexMetrics() *StandaloneIndexMetrics {
	if x != nil {
		return x.StandaloneIndexMetrics
	}
	return nil
}

func (x *IndexStatusResponse) GetUnmergedRecordCount() int64 {
	if x != nil {
		return x.UnmergedRecordCount
	}
	return 0
}

func (x *IndexStatusResponse) GetIndexHealerVectorRecordsIndexed() int64 {
	if x != nil {
		return x.IndexHealerVectorRecordsIndexed
	}
	return 0
}

func (x *IndexStatusResponse) GetIndexHealerVerticesValid() int64 {
	if x != nil {
		return x.IndexHealerVerticesValid
	}
	return 0
}

func (x *IndexStatusResponse) GetStatus() Status {
	if x != nil {
		return x.Status
	}
	return Status_READY
}

type GcInvalidVerticesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IndexId *IndexId `protobuf:"bytes,1,opt,name=indexId,proto3" json:"indexId,omitempty"`
	// Vertices identified as invalid before cutoff timestamp (Unix timestamp) are garbage collected.
	CutoffTimestamp int64 `protobuf:"varint,2,opt,name=cutoffTimestamp,proto3" json:"cutoffTimestamp,omitempty"`
}

func (x *GcInvalidVerticesRequest) Reset() {
	*x = GcInvalidVerticesRequest{}
	mi := &file_index_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GcInvalidVerticesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GcInvalidVerticesRequest) ProtoMessage() {}

func (x *GcInvalidVerticesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_index_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GcInvalidVerticesRequest.ProtoReflect.Descriptor instead.
func (*GcInvalidVerticesRequest) Descriptor() ([]byte, []int) {
	return file_index_proto_rawDescGZIP(), []int{2}
}

func (x *GcInvalidVerticesRequest) GetIndexId() *IndexId {
	if x != nil {
		return x.IndexId
	}
	return nil
}

func (x *GcInvalidVerticesRequest) GetCutoffTimestamp() int64 {
	if x != nil {
		return x.CutoffTimestamp
	}
	return 0
}

type IndexCreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Definition *IndexDefinition `protobuf:"bytes,1,opt,name=definition,proto3" json:"definition,omitempty"`
}

func (x *IndexCreateRequest) Reset() {
	*x = IndexCreateRequest{}
	mi := &file_index_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IndexCreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IndexCreateRequest) ProtoMessage() {}

func (x *IndexCreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_index_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IndexCreateRequest.ProtoReflect.Descriptor instead.
func (*IndexCreateRequest) Descriptor() ([]byte, []int) {
	return file_index_proto_rawDescGZIP(), []int{3}
}

func (x *IndexCreateRequest) GetDefinition() *IndexDefinition {
	if x != nil {
		return x.Definition
	}
	return nil
}

type IndexUpdateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IndexId *IndexId `protobuf:"bytes,1,opt,name=indexId,proto3" json:"indexId,omitempty"`
	// Optional labels associated with the index.
	Labels map[string]string `protobuf:"bytes,2,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Types that are assignable to Update:
	//
	//	*IndexUpdateRequest_HnswIndexUpdate
	Update isIndexUpdateRequest_Update `protobuf_oneof:"update"`
	// Mode of the index.
	Mode *IndexMode `protobuf:"varint,4,opt,name=mode,proto3,enum=aerospike.vector.IndexMode,oneof" json:"mode,omitempty"`
}

func (x *IndexUpdateRequest) Reset() {
	*x = IndexUpdateRequest{}
	mi := &file_index_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IndexUpdateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IndexUpdateRequest) ProtoMessage() {}

func (x *IndexUpdateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_index_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IndexUpdateRequest.ProtoReflect.Descriptor instead.
func (*IndexUpdateRequest) Descriptor() ([]byte, []int) {
	return file_index_proto_rawDescGZIP(), []int{4}
}

func (x *IndexUpdateRequest) GetIndexId() *IndexId {
	if x != nil {
		return x.IndexId
	}
	return nil
}

func (x *IndexUpdateRequest) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (m *IndexUpdateRequest) GetUpdate() isIndexUpdateRequest_Update {
	if m != nil {
		return m.Update
	}
	return nil
}

func (x *IndexUpdateRequest) GetHnswIndexUpdate() *HnswIndexUpdate {
	if x, ok := x.GetUpdate().(*IndexUpdateRequest_HnswIndexUpdate); ok {
		return x.HnswIndexUpdate
	}
	return nil
}

func (x *IndexUpdateRequest) GetMode() IndexMode {
	if x != nil && x.Mode != nil {
		return *x.Mode
	}
	return IndexMode_DISTRIBUTED
}

type isIndexUpdateRequest_Update interface {
	isIndexUpdateRequest_Update()
}

type IndexUpdateRequest_HnswIndexUpdate struct {
	HnswIndexUpdate *HnswIndexUpdate `protobuf:"bytes,3,opt,name=hnswIndexUpdate,proto3,oneof"`
}

func (*IndexUpdateRequest_HnswIndexUpdate) isIndexUpdateRequest_Update() {}

type IndexDropRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IndexId *IndexId `protobuf:"bytes,1,opt,name=indexId,proto3" json:"indexId,omitempty"`
}

func (x *IndexDropRequest) Reset() {
	*x = IndexDropRequest{}
	mi := &file_index_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IndexDropRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IndexDropRequest) ProtoMessage() {}

func (x *IndexDropRequest) ProtoReflect() protoreflect.Message {
	mi := &file_index_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IndexDropRequest.ProtoReflect.Descriptor instead.
func (*IndexDropRequest) Descriptor() ([]byte, []int) {
	return file_index_proto_rawDescGZIP(), []int{5}
}

func (x *IndexDropRequest) GetIndexId() *IndexId {
	if x != nil {
		return x.IndexId
	}
	return nil
}

type IndexListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Apply default values to parameters which are not set by user.
	ApplyDefaults *bool `protobuf:"varint,1,opt,name=applyDefaults,proto3,oneof" json:"applyDefaults,omitempty"`
}

func (x *IndexListRequest) Reset() {
	*x = IndexListRequest{}
	mi := &file_index_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IndexListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IndexListRequest) ProtoMessage() {}

func (x *IndexListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_index_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IndexListRequest.ProtoReflect.Descriptor instead.
func (*IndexListRequest) Descriptor() ([]byte, []int) {
	return file_index_proto_rawDescGZIP(), []int{6}
}

func (x *IndexListRequest) GetApplyDefaults() bool {
	if x != nil && x.ApplyDefaults != nil {
		return *x.ApplyDefaults
	}
	return false
}

type IndexGetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IndexId *IndexId `protobuf:"bytes,1,opt,name=indexId,proto3" json:"indexId,omitempty"`
	// Apply default values to parameters which are not set by user.
	ApplyDefaults *bool `protobuf:"varint,2,opt,name=applyDefaults,proto3,oneof" json:"applyDefaults,omitempty"`
}

func (x *IndexGetRequest) Reset() {
	*x = IndexGetRequest{}
	mi := &file_index_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IndexGetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IndexGetRequest) ProtoMessage() {}

func (x *IndexGetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_index_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IndexGetRequest.ProtoReflect.Descriptor instead.
func (*IndexGetRequest) Descriptor() ([]byte, []int) {
	return file_index_proto_rawDescGZIP(), []int{7}
}

func (x *IndexGetRequest) GetIndexId() *IndexId {
	if x != nil {
		return x.IndexId
	}
	return nil
}

func (x *IndexGetRequest) GetApplyDefaults() bool {
	if x != nil && x.ApplyDefaults != nil {
		return *x.ApplyDefaults
	}
	return false
}

type IndexStatusRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IndexId *IndexId `protobuf:"bytes,1,opt,name=indexId,proto3" json:"indexId,omitempty"`
}

func (x *IndexStatusRequest) Reset() {
	*x = IndexStatusRequest{}
	mi := &file_index_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IndexStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IndexStatusRequest) ProtoMessage() {}

func (x *IndexStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_index_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IndexStatusRequest.ProtoReflect.Descriptor instead.
func (*IndexStatusRequest) Descriptor() ([]byte, []int) {
	return file_index_proto_rawDescGZIP(), []int{8}
}

func (x *IndexStatusRequest) GetIndexId() *IndexId {
	if x != nil {
		return x.IndexId
	}
	return nil
}

var File_index_proto protoreflect.FileDescriptor

var file_index_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x61,
	0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x1a,
	0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0b, 0x74, 0x79,
	0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x83, 0x02, 0x0a, 0x16, 0x53, 0x74,
	0x61, 0x6e, 0x64, 0x61, 0x6c, 0x6f, 0x6e, 0x65, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x73, 0x12, 0x33, 0x0a, 0x07, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b,
	0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x49, 0x64,
	0x52, 0x07, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x49, 0x64, 0x12, 0x3c, 0x0a, 0x05, 0x73, 0x74, 0x61,
	0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x26, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73,
	0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x53, 0x74, 0x61, 0x6e,
	0x64, 0x61, 0x6c, 0x6f, 0x6e, 0x65, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x3a, 0x0a, 0x18, 0x73, 0x63, 0x61, 0x6e, 0x6e,
	0x65, 0x64, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x18, 0x73, 0x63, 0x61, 0x6e, 0x6e,
	0x65, 0x64, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x12, 0x3a, 0x0a, 0x18, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x64, 0x56, 0x65,
	0x63, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x18, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x64, 0x56, 0x65,
	0x63, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22,
	0x81, 0x03, 0x0a, 0x13, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x65, 0x0a, 0x16, 0x73, 0x74, 0x61, 0x6e, 0x64,
	0x61, 0x6c, 0x6f, 0x6e, 0x65, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70,
	0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x53, 0x74, 0x61, 0x6e, 0x64,
	0x61, 0x6c, 0x6f, 0x6e, 0x65, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x73, 0x48, 0x00, 0x52, 0x16, 0x73, 0x74, 0x61, 0x6e, 0x64, 0x61, 0x6c, 0x6f, 0x6e, 0x65, 0x49,
	0x6e, 0x64, 0x65, 0x78, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x88, 0x01, 0x01, 0x12, 0x30,
	0x0a, 0x13, 0x75, 0x6e, 0x6d, 0x65, 0x72, 0x67, 0x65, 0x64, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x13, 0x75, 0x6e, 0x6d,
	0x65, 0x72, 0x67, 0x65, 0x64, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x12, 0x48, 0x0a, 0x1f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x48, 0x65, 0x61, 0x6c, 0x65, 0x72, 0x56,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x49, 0x6e, 0x64, 0x65,
	0x78, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x1f, 0x69, 0x6e, 0x64, 0x65, 0x78,
	0x48, 0x65, 0x61, 0x6c, 0x65, 0x72, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x63, 0x6f,
	0x72, 0x64, 0x73, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x64, 0x12, 0x3a, 0x0a, 0x18, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x48, 0x65, 0x61, 0x6c, 0x65, 0x72, 0x56, 0x65, 0x72, 0x74, 0x69, 0x63, 0x65,
	0x73, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x18, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x48, 0x65, 0x61, 0x6c, 0x65, 0x72, 0x56, 0x65, 0x72, 0x74, 0x69, 0x63, 0x65,
	0x73, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x12, 0x30, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x18, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69,
	0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x19, 0x0a, 0x17, 0x5f, 0x73, 0x74, 0x61,
	0x6e, 0x64, 0x61, 0x6c, 0x6f, 0x6e, 0x65, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x73, 0x22, 0x79, 0x0a, 0x18, 0x47, 0x63, 0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x56, 0x65, 0x72, 0x74, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x33, 0x0a, 0x07, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x49, 0x64, 0x52, 0x07, 0x69, 0x6e, 0x64,
	0x65, 0x78, 0x49, 0x64, 0x12, 0x28, 0x0a, 0x0f, 0x63, 0x75, 0x74, 0x6f, 0x66, 0x66, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0f, 0x63,
	0x75, 0x74, 0x6f, 0x66, 0x66, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x57,
	0x0a, 0x12, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x41, 0x0a, 0x0a, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73,
	0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x49, 0x6e, 0x64, 0x65,
	0x78, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x64, 0x65, 0x66,
	0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xe6, 0x02, 0x0a, 0x12, 0x49, 0x6e, 0x64, 0x65,
	0x78, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x33,
	0x0a, 0x07, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x19, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x49, 0x64, 0x52, 0x07, 0x69, 0x6e, 0x64, 0x65,
	0x78, 0x49, 0x64, 0x12, 0x48, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x30, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e,
	0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x4d, 0x0a,
	0x0f, 0x68, 0x6e, 0x73, 0x77, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69,
	0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x48, 0x6e, 0x73, 0x77, 0x49, 0x6e,
	0x64, 0x65, 0x78, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x48, 0x00, 0x52, 0x0f, 0x68, 0x6e, 0x73,
	0x77, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x34, 0x0a, 0x04,
	0x6d, 0x6f, 0x64, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x61, 0x65, 0x72,
	0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x49, 0x6e,
	0x64, 0x65, 0x78, 0x4d, 0x6f, 0x64, 0x65, 0x48, 0x01, 0x52, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x88,
	0x01, 0x01, 0x1a, 0x39, 0x0a, 0x0b, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x08, 0x0a,
	0x06, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x6d, 0x6f, 0x64, 0x65,
	0x22, 0x47, 0x0a, 0x10, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x44, 0x72, 0x6f, 0x70, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x33, 0x0a, 0x07, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b,
	0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x49, 0x64,
	0x52, 0x07, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x49, 0x64, 0x22, 0x4f, 0x0a, 0x10, 0x49, 0x6e, 0x64,
	0x65, 0x78, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x29, 0x0a,
	0x0d, 0x61, 0x70, 0x70, 0x6c, 0x79, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x0d, 0x61, 0x70, 0x70, 0x6c, 0x79, 0x44, 0x65, 0x66,
	0x61, 0x75, 0x6c, 0x74, 0x73, 0x88, 0x01, 0x01, 0x42, 0x10, 0x0a, 0x0e, 0x5f, 0x61, 0x70, 0x70,
	0x6c, 0x79, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x73, 0x22, 0x83, 0x01, 0x0a, 0x0f, 0x49,
	0x6e, 0x64, 0x65, 0x78, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x33,
	0x0a, 0x07, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x19, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x49, 0x64, 0x52, 0x07, 0x69, 0x6e, 0x64, 0x65,
	0x78, 0x49, 0x64, 0x12, 0x29, 0x0a, 0x0d, 0x61, 0x70, 0x70, 0x6c, 0x79, 0x44, 0x65, 0x66, 0x61,
	0x75, 0x6c, 0x74, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x0d, 0x61, 0x70,
	0x70, 0x6c, 0x79, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x73, 0x88, 0x01, 0x01, 0x42, 0x10,
	0x0a, 0x0e, 0x5f, 0x61, 0x70, 0x70, 0x6c, 0x79, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x73,
	0x22, 0x49, 0x0a, 0x12, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x33, 0x0a, 0x07, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x49,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70,
	0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78,
	0x49, 0x64, 0x52, 0x07, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x49, 0x64, 0x32, 0x8c, 0x05, 0x0a, 0x0c,
	0x49, 0x6e, 0x64, 0x65, 0x78, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x48, 0x0a, 0x06,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x24, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69,
	0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x48, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x12, 0x24, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00,
	0x12, 0x44, 0x0a, 0x04, 0x44, 0x72, 0x6f, 0x70, 0x12, 0x22, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73,
	0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x49, 0x6e, 0x64, 0x65,
	0x78, 0x44, 0x72, 0x6f, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x53, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x22,
	0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x25, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x44, 0x65, 0x66, 0x69, 0x6e,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x00, 0x12, 0x4d, 0x0a, 0x03, 0x47,
	0x65, 0x74, 0x12, 0x21, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x47, 0x65, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b,
	0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x44, 0x65,
	0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x00, 0x12, 0x5a, 0x0a, 0x09, 0x47, 0x65,
	0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x24, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70,
	0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e,
	0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x59, 0x0a, 0x11, 0x47, 0x63, 0x49, 0x6e, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x56, 0x65, 0x72, 0x74, 0x69, 0x63, 0x65, 0x73, 0x12, 0x2a, 0x2e, 0x61, 0x65,
	0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x47,
	0x63, 0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x56, 0x65, 0x72, 0x74, 0x69, 0x63, 0x65, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22,
	0x00, 0x12, 0x47, 0x0a, 0x10, 0x41, 0x72, 0x65, 0x49, 0x6e, 0x64, 0x69, 0x63, 0x65, 0x73, 0x49,
	0x6e, 0x53, 0x79, 0x6e, 0x63, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x19, 0x2e,
	0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x65, 0x61, 0x6e, 0x22, 0x00, 0x42, 0x43, 0x0a, 0x21, 0x63, 0x6f,
	0x6d, 0x2e, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x76, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x50,
	0x01, 0x5a, 0x1c, 0x61, 0x65, 0x72, 0x6f, 0x73, 0x70, 0x69, 0x6b, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_index_proto_rawDescOnce sync.Once
	file_index_proto_rawDescData = file_index_proto_rawDesc
)

func file_index_proto_rawDescGZIP() []byte {
	file_index_proto_rawDescOnce.Do(func() {
		file_index_proto_rawDescData = protoimpl.X.CompressGZIP(file_index_proto_rawDescData)
	})
	return file_index_proto_rawDescData
}

var file_index_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_index_proto_goTypes = []any{
	(*StandaloneIndexMetrics)(nil),   // 0: aerospike.vector.StandaloneIndexMetrics
	(*IndexStatusResponse)(nil),      // 1: aerospike.vector.IndexStatusResponse
	(*GcInvalidVerticesRequest)(nil), // 2: aerospike.vector.GcInvalidVerticesRequest
	(*IndexCreateRequest)(nil),       // 3: aerospike.vector.IndexCreateRequest
	(*IndexUpdateRequest)(nil),       // 4: aerospike.vector.IndexUpdateRequest
	(*IndexDropRequest)(nil),         // 5: aerospike.vector.IndexDropRequest
	(*IndexListRequest)(nil),         // 6: aerospike.vector.IndexListRequest
	(*IndexGetRequest)(nil),          // 7: aerospike.vector.IndexGetRequest
	(*IndexStatusRequest)(nil),       // 8: aerospike.vector.IndexStatusRequest
	nil,                              // 9: aerospike.vector.IndexUpdateRequest.LabelsEntry
	(*IndexId)(nil),                  // 10: aerospike.vector.IndexId
	(StandaloneIndexState)(0),        // 11: aerospike.vector.StandaloneIndexState
	(Status)(0),                      // 12: aerospike.vector.Status
	(*IndexDefinition)(nil),          // 13: aerospike.vector.IndexDefinition
	(*HnswIndexUpdate)(nil),          // 14: aerospike.vector.HnswIndexUpdate
	(IndexMode)(0),                   // 15: aerospike.vector.IndexMode
	(*emptypb.Empty)(nil),            // 16: google.protobuf.Empty
	(*IndexDefinitionList)(nil),      // 17: aerospike.vector.IndexDefinitionList
	(*Boolean)(nil),                  // 18: aerospike.vector.Boolean
}
var file_index_proto_depIdxs = []int32{
	10, // 0: aerospike.vector.StandaloneIndexMetrics.indexId:type_name -> aerospike.vector.IndexId
	11, // 1: aerospike.vector.StandaloneIndexMetrics.state:type_name -> aerospike.vector.StandaloneIndexState
	0,  // 2: aerospike.vector.IndexStatusResponse.standaloneIndexMetrics:type_name -> aerospike.vector.StandaloneIndexMetrics
	12, // 3: aerospike.vector.IndexStatusResponse.status:type_name -> aerospike.vector.Status
	10, // 4: aerospike.vector.GcInvalidVerticesRequest.indexId:type_name -> aerospike.vector.IndexId
	13, // 5: aerospike.vector.IndexCreateRequest.definition:type_name -> aerospike.vector.IndexDefinition
	10, // 6: aerospike.vector.IndexUpdateRequest.indexId:type_name -> aerospike.vector.IndexId
	9,  // 7: aerospike.vector.IndexUpdateRequest.labels:type_name -> aerospike.vector.IndexUpdateRequest.LabelsEntry
	14, // 8: aerospike.vector.IndexUpdateRequest.hnswIndexUpdate:type_name -> aerospike.vector.HnswIndexUpdate
	15, // 9: aerospike.vector.IndexUpdateRequest.mode:type_name -> aerospike.vector.IndexMode
	10, // 10: aerospike.vector.IndexDropRequest.indexId:type_name -> aerospike.vector.IndexId
	10, // 11: aerospike.vector.IndexGetRequest.indexId:type_name -> aerospike.vector.IndexId
	10, // 12: aerospike.vector.IndexStatusRequest.indexId:type_name -> aerospike.vector.IndexId
	3,  // 13: aerospike.vector.IndexService.Create:input_type -> aerospike.vector.IndexCreateRequest
	4,  // 14: aerospike.vector.IndexService.Update:input_type -> aerospike.vector.IndexUpdateRequest
	5,  // 15: aerospike.vector.IndexService.Drop:input_type -> aerospike.vector.IndexDropRequest
	6,  // 16: aerospike.vector.IndexService.List:input_type -> aerospike.vector.IndexListRequest
	7,  // 17: aerospike.vector.IndexService.Get:input_type -> aerospike.vector.IndexGetRequest
	8,  // 18: aerospike.vector.IndexService.GetStatus:input_type -> aerospike.vector.IndexStatusRequest
	2,  // 19: aerospike.vector.IndexService.GcInvalidVertices:input_type -> aerospike.vector.GcInvalidVerticesRequest
	16, // 20: aerospike.vector.IndexService.AreIndicesInSync:input_type -> google.protobuf.Empty
	16, // 21: aerospike.vector.IndexService.Create:output_type -> google.protobuf.Empty
	16, // 22: aerospike.vector.IndexService.Update:output_type -> google.protobuf.Empty
	16, // 23: aerospike.vector.IndexService.Drop:output_type -> google.protobuf.Empty
	17, // 24: aerospike.vector.IndexService.List:output_type -> aerospike.vector.IndexDefinitionList
	13, // 25: aerospike.vector.IndexService.Get:output_type -> aerospike.vector.IndexDefinition
	1,  // 26: aerospike.vector.IndexService.GetStatus:output_type -> aerospike.vector.IndexStatusResponse
	16, // 27: aerospike.vector.IndexService.GcInvalidVertices:output_type -> google.protobuf.Empty
	18, // 28: aerospike.vector.IndexService.AreIndicesInSync:output_type -> aerospike.vector.Boolean
	21, // [21:29] is the sub-list for method output_type
	13, // [13:21] is the sub-list for method input_type
	13, // [13:13] is the sub-list for extension type_name
	13, // [13:13] is the sub-list for extension extendee
	0,  // [0:13] is the sub-list for field type_name
}

func init() { file_index_proto_init() }
func file_index_proto_init() {
	if File_index_proto != nil {
		return
	}
	file_types_proto_init()
	file_index_proto_msgTypes[1].OneofWrappers = []any{}
	file_index_proto_msgTypes[4].OneofWrappers = []any{
		(*IndexUpdateRequest_HnswIndexUpdate)(nil),
	}
	file_index_proto_msgTypes[6].OneofWrappers = []any{}
	file_index_proto_msgTypes[7].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_index_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_index_proto_goTypes,
		DependencyIndexes: file_index_proto_depIdxs,
		MessageInfos:      file_index_proto_msgTypes,
	}.Build()
	File_index_proto = out.File
	file_index_proto_rawDesc = nil
	file_index_proto_goTypes = nil
	file_index_proto_depIdxs = nil
}
