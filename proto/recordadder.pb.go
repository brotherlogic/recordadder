// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0-devel
// 	protoc        (unknown)
// source: recordadder.proto

package proto

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Queue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Requests         []*AddRecordRequest `protobuf:"bytes,1,rep,name=requests,proto3" json:"requests,omitempty"`
	ProcessedRecords int32               `protobuf:"varint,2,opt,name=processed_records,json=processedRecords,proto3" json:"processed_records,omitempty"`
	LastAdditionDate int64               `protobuf:"varint,3,opt,name=last_addition_date,json=lastAdditionDate,proto3" json:"last_addition_date,omitempty"`
}

func (x *Queue) Reset() {
	*x = Queue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_recordadder_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Queue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Queue) ProtoMessage() {}

func (x *Queue) ProtoReflect() protoreflect.Message {
	mi := &file_recordadder_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Queue.ProtoReflect.Descriptor instead.
func (*Queue) Descriptor() ([]byte, []int) {
	return file_recordadder_proto_rawDescGZIP(), []int{0}
}

func (x *Queue) GetRequests() []*AddRecordRequest {
	if x != nil {
		return x.Requests
	}
	return nil
}

func (x *Queue) GetProcessedRecords() int32 {
	if x != nil {
		return x.ProcessedRecords
	}
	return 0
}

func (x *Queue) GetLastAdditionDate() int64 {
	if x != nil {
		return x.LastAdditionDate
	}
	return 0
}

type AddRecordRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             int32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Cost           int32 `protobuf:"varint,2,opt,name=cost,proto3" json:"cost,omitempty"`
	Folder         int32 `protobuf:"varint,3,opt,name=folder,proto3" json:"folder,omitempty"`
	AccountingYear int32 `protobuf:"varint,5,opt,name=accounting_year,json=accountingYear,proto3" json:"accounting_year,omitempty"`
	ResetFolder    int32 `protobuf:"varint,4,opt,name=reset_folder,json=resetFolder,proto3" json:"reset_folder,omitempty"`
	Arrived        bool  `protobuf:"varint,6,opt,name=arrived,proto3" json:"arrived,omitempty"`
}

func (x *AddRecordRequest) Reset() {
	*x = AddRecordRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_recordadder_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddRecordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddRecordRequest) ProtoMessage() {}

func (x *AddRecordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_recordadder_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddRecordRequest.ProtoReflect.Descriptor instead.
func (*AddRecordRequest) Descriptor() ([]byte, []int) {
	return file_recordadder_proto_rawDescGZIP(), []int{1}
}

func (x *AddRecordRequest) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *AddRecordRequest) GetCost() int32 {
	if x != nil {
		return x.Cost
	}
	return 0
}

func (x *AddRecordRequest) GetFolder() int32 {
	if x != nil {
		return x.Folder
	}
	return 0
}

func (x *AddRecordRequest) GetAccountingYear() int32 {
	if x != nil {
		return x.AccountingYear
	}
	return 0
}

func (x *AddRecordRequest) GetResetFolder() int32 {
	if x != nil {
		return x.ResetFolder
	}
	return 0
}

func (x *AddRecordRequest) GetArrived() bool {
	if x != nil {
		return x.Arrived
	}
	return false
}

type AddRecordResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ExpectedAdditionDate int64 `protobuf:"varint,1,opt,name=expected_addition_date,json=expectedAdditionDate,proto3" json:"expected_addition_date,omitempty"`
}

func (x *AddRecordResponse) Reset() {
	*x = AddRecordResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_recordadder_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddRecordResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddRecordResponse) ProtoMessage() {}

func (x *AddRecordResponse) ProtoReflect() protoreflect.Message {
	mi := &file_recordadder_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddRecordResponse.ProtoReflect.Descriptor instead.
func (*AddRecordResponse) Descriptor() ([]byte, []int) {
	return file_recordadder_proto_rawDescGZIP(), []int{2}
}

func (x *AddRecordResponse) GetExpectedAdditionDate() int64 {
	if x != nil {
		return x.ExpectedAdditionDate
	}
	return 0
}

type ListQueueRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListQueueRequest) Reset() {
	*x = ListQueueRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_recordadder_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListQueueRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListQueueRequest) ProtoMessage() {}

func (x *ListQueueRequest) ProtoReflect() protoreflect.Message {
	mi := &file_recordadder_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListQueueRequest.ProtoReflect.Descriptor instead.
func (*ListQueueRequest) Descriptor() ([]byte, []int) {
	return file_recordadder_proto_rawDescGZIP(), []int{3}
}

type ListQueueResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Requests []*AddRecordRequest `protobuf:"bytes,1,rep,name=requests,proto3" json:"requests,omitempty"`
}

func (x *ListQueueResponse) Reset() {
	*x = ListQueueResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_recordadder_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListQueueResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListQueueResponse) ProtoMessage() {}

func (x *ListQueueResponse) ProtoReflect() protoreflect.Message {
	mi := &file_recordadder_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListQueueResponse.ProtoReflect.Descriptor instead.
func (*ListQueueResponse) Descriptor() ([]byte, []int) {
	return file_recordadder_proto_rawDescGZIP(), []int{4}
}

func (x *ListQueueResponse) GetRequests() []*AddRecordRequest {
	if x != nil {
		return x.Requests
	}
	return nil
}

type UpdateRecordRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        int32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Available bool  `protobuf:"varint,2,opt,name=available,proto3" json:"available,omitempty"`
}

func (x *UpdateRecordRequest) Reset() {
	*x = UpdateRecordRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_recordadder_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateRecordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateRecordRequest) ProtoMessage() {}

func (x *UpdateRecordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_recordadder_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateRecordRequest.ProtoReflect.Descriptor instead.
func (*UpdateRecordRequest) Descriptor() ([]byte, []int) {
	return file_recordadder_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateRecordRequest) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateRecordRequest) GetAvailable() bool {
	if x != nil {
		return x.Available
	}
	return false
}

type UpdateRecordResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UpdateRecordResponse) Reset() {
	*x = UpdateRecordResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_recordadder_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateRecordResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateRecordResponse) ProtoMessage() {}

func (x *UpdateRecordResponse) ProtoReflect() protoreflect.Message {
	mi := &file_recordadder_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateRecordResponse.ProtoReflect.Descriptor instead.
func (*UpdateRecordResponse) Descriptor() ([]byte, []int) {
	return file_recordadder_proto_rawDescGZIP(), []int{6}
}

type DeleteRecordRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteRecordRequest) Reset() {
	*x = DeleteRecordRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_recordadder_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteRecordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRecordRequest) ProtoMessage() {}

func (x *DeleteRecordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_recordadder_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRecordRequest.ProtoReflect.Descriptor instead.
func (*DeleteRecordRequest) Descriptor() ([]byte, []int) {
	return file_recordadder_proto_rawDescGZIP(), []int{7}
}

func (x *DeleteRecordRequest) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type DeleteRecordResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteRecordResponse) Reset() {
	*x = DeleteRecordResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_recordadder_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteRecordResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRecordResponse) ProtoMessage() {}

func (x *DeleteRecordResponse) ProtoReflect() protoreflect.Message {
	mi := &file_recordadder_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRecordResponse.ProtoReflect.Descriptor instead.
func (*DeleteRecordResponse) Descriptor() ([]byte, []int) {
	return file_recordadder_proto_rawDescGZIP(), []int{8}
}

var File_recordadder_proto protoreflect.FileDescriptor

var file_recordadder_proto_rawDesc = []byte{
	0x0a, 0x11, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x61, 0x64, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x61, 0x64, 0x64, 0x65, 0x72,
	0x22, 0x9d, 0x01, 0x0a, 0x05, 0x51, 0x75, 0x65, 0x75, 0x65, 0x12, 0x39, 0x0a, 0x08, 0x72, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x72,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x61, 0x64, 0x64, 0x65, 0x72, 0x2e, 0x41, 0x64, 0x64, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x08, 0x72, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x73, 0x12, 0x2b, 0x0a, 0x11, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73,
	0x65, 0x64, 0x5f, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x10, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x64, 0x52, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x73, 0x12, 0x2c, 0x0a, 0x12, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x61, 0x64, 0x64, 0x69, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x10,
	0x6c, 0x61, 0x73, 0x74, 0x41, 0x64, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x65,
	0x22, 0xb4, 0x01, 0x0a, 0x10, 0x41, 0x64, 0x64, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x73, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x6f, 0x6c,
	0x64, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x66, 0x6f, 0x6c, 0x64, 0x65,
	0x72, 0x12, 0x27, 0x0a, 0x0f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x5f,
	0x79, 0x65, 0x61, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x61, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x59, 0x65, 0x61, 0x72, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65,
	0x73, 0x65, 0x74, 0x5f, 0x66, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0b, 0x72, 0x65, 0x73, 0x65, 0x74, 0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x12, 0x18, 0x0a,
	0x07, 0x61, 0x72, 0x72, 0x69, 0x76, 0x65, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07,
	0x61, 0x72, 0x72, 0x69, 0x76, 0x65, 0x64, 0x22, 0x49, 0x0a, 0x11, 0x41, 0x64, 0x64, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x34, 0x0a, 0x16,
	0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x64, 0x64, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x14, 0x65, 0x78,
	0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x41, 0x64, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x61,
	0x74, 0x65, 0x22, 0x12, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x51, 0x75, 0x65, 0x75, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x4e, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x51, 0x75,
	0x65, 0x75, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x39, 0x0a, 0x08, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e,
	0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x61, 0x64, 0x64, 0x65, 0x72, 0x2e, 0x41, 0x64, 0x64, 0x52,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x08, 0x72, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x22, 0x43, 0x0a, 0x13, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1c, 0x0a,
	0x09, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x09, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x22, 0x16, 0x0a, 0x14, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x25, 0x0a, 0x13, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x22, 0x16, 0x0a, 0x14, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x32, 0xdc, 0x02, 0x0a, 0x10, 0x41, 0x64, 0x64, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4c, 0x0a, 0x09, 0x41, 0x64, 0x64, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x12, 0x1d, 0x2e, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x61, 0x64, 0x64,
	0x65, 0x72, 0x2e, 0x41, 0x64, 0x64, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x61, 0x64, 0x64, 0x65,
	0x72, 0x2e, 0x41, 0x64, 0x64, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4c, 0x0a, 0x09, 0x4c, 0x69, 0x73, 0x74, 0x51, 0x75, 0x65,
	0x75, 0x65, 0x12, 0x1d, 0x2e, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x61, 0x64, 0x64, 0x65, 0x72,
	0x2e, 0x4c, 0x69, 0x73, 0x74, 0x51, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1e, 0x2e, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x61, 0x64, 0x64, 0x65, 0x72, 0x2e,
	0x4c, 0x69, 0x73, 0x74, 0x51, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x55, 0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x12, 0x20, 0x2e, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x61, 0x64, 0x64, 0x65,
	0x72, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x61, 0x64,
	0x64, 0x65, 0x72, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x55, 0x0a, 0x0c, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x20, 0x2e, 0x72, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x61, 0x64, 0x64, 0x65, 0x72, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x72,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x61, 0x64, 0x64, 0x65, 0x72, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x42, 0x2b, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x62, 0x72, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x2f, 0x72, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x61, 0x64, 0x64, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_recordadder_proto_rawDescOnce sync.Once
	file_recordadder_proto_rawDescData = file_recordadder_proto_rawDesc
)

func file_recordadder_proto_rawDescGZIP() []byte {
	file_recordadder_proto_rawDescOnce.Do(func() {
		file_recordadder_proto_rawDescData = protoimpl.X.CompressGZIP(file_recordadder_proto_rawDescData)
	})
	return file_recordadder_proto_rawDescData
}

var file_recordadder_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_recordadder_proto_goTypes = []interface{}{
	(*Queue)(nil),                // 0: recordadder.Queue
	(*AddRecordRequest)(nil),     // 1: recordadder.AddRecordRequest
	(*AddRecordResponse)(nil),    // 2: recordadder.AddRecordResponse
	(*ListQueueRequest)(nil),     // 3: recordadder.ListQueueRequest
	(*ListQueueResponse)(nil),    // 4: recordadder.ListQueueResponse
	(*UpdateRecordRequest)(nil),  // 5: recordadder.UpdateRecordRequest
	(*UpdateRecordResponse)(nil), // 6: recordadder.UpdateRecordResponse
	(*DeleteRecordRequest)(nil),  // 7: recordadder.DeleteRecordRequest
	(*DeleteRecordResponse)(nil), // 8: recordadder.DeleteRecordResponse
}
var file_recordadder_proto_depIdxs = []int32{
	1, // 0: recordadder.Queue.requests:type_name -> recordadder.AddRecordRequest
	1, // 1: recordadder.ListQueueResponse.requests:type_name -> recordadder.AddRecordRequest
	1, // 2: recordadder.AddRecordService.AddRecord:input_type -> recordadder.AddRecordRequest
	3, // 3: recordadder.AddRecordService.ListQueue:input_type -> recordadder.ListQueueRequest
	5, // 4: recordadder.AddRecordService.UpdateRecord:input_type -> recordadder.UpdateRecordRequest
	7, // 5: recordadder.AddRecordService.DeleteRecord:input_type -> recordadder.DeleteRecordRequest
	2, // 6: recordadder.AddRecordService.AddRecord:output_type -> recordadder.AddRecordResponse
	4, // 7: recordadder.AddRecordService.ListQueue:output_type -> recordadder.ListQueueResponse
	6, // 8: recordadder.AddRecordService.UpdateRecord:output_type -> recordadder.UpdateRecordResponse
	8, // 9: recordadder.AddRecordService.DeleteRecord:output_type -> recordadder.DeleteRecordResponse
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_recordadder_proto_init() }
func file_recordadder_proto_init() {
	if File_recordadder_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_recordadder_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Queue); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_recordadder_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddRecordRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_recordadder_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddRecordResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_recordadder_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListQueueRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_recordadder_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListQueueResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_recordadder_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateRecordRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_recordadder_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateRecordResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_recordadder_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteRecordRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_recordadder_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteRecordResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_recordadder_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_recordadder_proto_goTypes,
		DependencyIndexes: file_recordadder_proto_depIdxs,
		MessageInfos:      file_recordadder_proto_msgTypes,
	}.Build()
	File_recordadder_proto = out.File
	file_recordadder_proto_rawDesc = nil
	file_recordadder_proto_goTypes = nil
	file_recordadder_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AddRecordServiceClient is the client API for AddRecordService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AddRecordServiceClient interface {
	AddRecord(ctx context.Context, in *AddRecordRequest, opts ...grpc.CallOption) (*AddRecordResponse, error)
	ListQueue(ctx context.Context, in *ListQueueRequest, opts ...grpc.CallOption) (*ListQueueResponse, error)
	UpdateRecord(ctx context.Context, in *UpdateRecordRequest, opts ...grpc.CallOption) (*UpdateRecordResponse, error)
	DeleteRecord(ctx context.Context, in *DeleteRecordRequest, opts ...grpc.CallOption) (*DeleteRecordResponse, error)
}

type addRecordServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAddRecordServiceClient(cc grpc.ClientConnInterface) AddRecordServiceClient {
	return &addRecordServiceClient{cc}
}

func (c *addRecordServiceClient) AddRecord(ctx context.Context, in *AddRecordRequest, opts ...grpc.CallOption) (*AddRecordResponse, error) {
	out := new(AddRecordResponse)
	err := c.cc.Invoke(ctx, "/recordadder.AddRecordService/AddRecord", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *addRecordServiceClient) ListQueue(ctx context.Context, in *ListQueueRequest, opts ...grpc.CallOption) (*ListQueueResponse, error) {
	out := new(ListQueueResponse)
	err := c.cc.Invoke(ctx, "/recordadder.AddRecordService/ListQueue", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *addRecordServiceClient) UpdateRecord(ctx context.Context, in *UpdateRecordRequest, opts ...grpc.CallOption) (*UpdateRecordResponse, error) {
	out := new(UpdateRecordResponse)
	err := c.cc.Invoke(ctx, "/recordadder.AddRecordService/UpdateRecord", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *addRecordServiceClient) DeleteRecord(ctx context.Context, in *DeleteRecordRequest, opts ...grpc.CallOption) (*DeleteRecordResponse, error) {
	out := new(DeleteRecordResponse)
	err := c.cc.Invoke(ctx, "/recordadder.AddRecordService/DeleteRecord", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AddRecordServiceServer is the server API for AddRecordService service.
type AddRecordServiceServer interface {
	AddRecord(context.Context, *AddRecordRequest) (*AddRecordResponse, error)
	ListQueue(context.Context, *ListQueueRequest) (*ListQueueResponse, error)
	UpdateRecord(context.Context, *UpdateRecordRequest) (*UpdateRecordResponse, error)
	DeleteRecord(context.Context, *DeleteRecordRequest) (*DeleteRecordResponse, error)
}

// UnimplementedAddRecordServiceServer can be embedded to have forward compatible implementations.
type UnimplementedAddRecordServiceServer struct {
}

func (*UnimplementedAddRecordServiceServer) AddRecord(context.Context, *AddRecordRequest) (*AddRecordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddRecord not implemented")
}
func (*UnimplementedAddRecordServiceServer) ListQueue(context.Context, *ListQueueRequest) (*ListQueueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListQueue not implemented")
}
func (*UnimplementedAddRecordServiceServer) UpdateRecord(context.Context, *UpdateRecordRequest) (*UpdateRecordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRecord not implemented")
}
func (*UnimplementedAddRecordServiceServer) DeleteRecord(context.Context, *DeleteRecordRequest) (*DeleteRecordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRecord not implemented")
}

func RegisterAddRecordServiceServer(s *grpc.Server, srv AddRecordServiceServer) {
	s.RegisterService(&_AddRecordService_serviceDesc, srv)
}

func _AddRecordService_AddRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddRecordServiceServer).AddRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/recordadder.AddRecordService/AddRecord",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddRecordServiceServer).AddRecord(ctx, req.(*AddRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AddRecordService_ListQueue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListQueueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddRecordServiceServer).ListQueue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/recordadder.AddRecordService/ListQueue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddRecordServiceServer).ListQueue(ctx, req.(*ListQueueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AddRecordService_UpdateRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddRecordServiceServer).UpdateRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/recordadder.AddRecordService/UpdateRecord",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddRecordServiceServer).UpdateRecord(ctx, req.(*UpdateRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AddRecordService_DeleteRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddRecordServiceServer).DeleteRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/recordadder.AddRecordService/DeleteRecord",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddRecordServiceServer).DeleteRecord(ctx, req.(*DeleteRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AddRecordService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "recordadder.AddRecordService",
	HandlerType: (*AddRecordServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddRecord",
			Handler:    _AddRecordService_AddRecord_Handler,
		},
		{
			MethodName: "ListQueue",
			Handler:    _AddRecordService_ListQueue_Handler,
		},
		{
			MethodName: "UpdateRecord",
			Handler:    _AddRecordService_UpdateRecord_Handler,
		},
		{
			MethodName: "DeleteRecord",
			Handler:    _AddRecordService_DeleteRecord_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "recordadder.proto",
}
