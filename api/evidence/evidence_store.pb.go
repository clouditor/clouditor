// Copyright 2016-2022 Fraunhofer AISEC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//           $$\                           $$\ $$\   $$\
//           $$ |                          $$ |\__|  $$ |
//  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
// $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
// $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
// $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
// \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
//  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
//
// This file is part of Clouditor Community Edition.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: api/evidence/evidence_store.proto

package evidence

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type StoreEvidenceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Evidence *Evidence `protobuf:"bytes,1,opt,name=evidence,proto3" json:"evidence,omitempty"`
}

func (x *StoreEvidenceRequest) Reset() {
	*x = StoreEvidenceRequest{}
	mi := &file_api_evidence_evidence_store_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StoreEvidenceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreEvidenceRequest) ProtoMessage() {}

func (x *StoreEvidenceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_evidence_evidence_store_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreEvidenceRequest.ProtoReflect.Descriptor instead.
func (*StoreEvidenceRequest) Descriptor() ([]byte, []int) {
	return file_api_evidence_evidence_store_proto_rawDescGZIP(), []int{0}
}

func (x *StoreEvidenceRequest) GetEvidence() *Evidence {
	if x != nil {
		return x.Evidence
	}
	return nil
}

// StoreEvidenceResponse belongs to StoreEvidence, which uses a custom unary RPC and therefore requires a response message according to the style convention. Since no return values are required, this is empty.
type StoreEvidenceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StoreEvidenceResponse) Reset() {
	*x = StoreEvidenceResponse{}
	mi := &file_api_evidence_evidence_store_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StoreEvidenceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreEvidenceResponse) ProtoMessage() {}

func (x *StoreEvidenceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_evidence_evidence_store_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreEvidenceResponse.ProtoReflect.Descriptor instead.
func (*StoreEvidenceResponse) Descriptor() ([]byte, []int) {
	return file_api_evidence_evidence_store_proto_rawDescGZIP(), []int{1}
}

// StoreEvidencesResponse belongs to StoreEvidences, which uses a custom bidirectional streaming RPC and therefore requires a response message according to the style convention. The bidirectional streaming needs the status and its message in the response for error handling.
type StoreEvidencesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status        bool   `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	StatusMessage string `protobuf:"bytes,2,opt,name=status_message,json=statusMessage,proto3" json:"status_message,omitempty"`
}

func (x *StoreEvidencesResponse) Reset() {
	*x = StoreEvidencesResponse{}
	mi := &file_api_evidence_evidence_store_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StoreEvidencesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreEvidencesResponse) ProtoMessage() {}

func (x *StoreEvidencesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_evidence_evidence_store_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreEvidencesResponse.ProtoReflect.Descriptor instead.
func (*StoreEvidencesResponse) Descriptor() ([]byte, []int) {
	return file_api_evidence_evidence_store_proto_rawDescGZIP(), []int{2}
}

func (x *StoreEvidencesResponse) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

func (x *StoreEvidencesResponse) GetStatusMessage() string {
	if x != nil {
		return x.StatusMessage
	}
	return ""
}

type ListEvidencesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Filter *Filter `protobuf:"bytes,1,opt,name=filter,proto3,oneof" json:"filter,omitempty"`
	// page_size: 0 = default (50 is default value), > 0 = set value (i.e. page_size = 5 -> SQL-Limit = 5)
	PageSize  int32  `protobuf:"varint,10,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	PageToken string `protobuf:"bytes,11,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	OrderBy   string `protobuf:"bytes,12,opt,name=order_by,json=orderBy,proto3" json:"order_by,omitempty"`
	Asc       bool   `protobuf:"varint,13,opt,name=asc,proto3" json:"asc,omitempty"`
}

func (x *ListEvidencesRequest) Reset() {
	*x = ListEvidencesRequest{}
	mi := &file_api_evidence_evidence_store_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListEvidencesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListEvidencesRequest) ProtoMessage() {}

func (x *ListEvidencesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_evidence_evidence_store_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListEvidencesRequest.ProtoReflect.Descriptor instead.
func (*ListEvidencesRequest) Descriptor() ([]byte, []int) {
	return file_api_evidence_evidence_store_proto_rawDescGZIP(), []int{3}
}

func (x *ListEvidencesRequest) GetFilter() *Filter {
	if x != nil {
		return x.Filter
	}
	return nil
}

func (x *ListEvidencesRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListEvidencesRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

func (x *ListEvidencesRequest) GetOrderBy() string {
	if x != nil {
		return x.OrderBy
	}
	return ""
}

func (x *ListEvidencesRequest) GetAsc() bool {
	if x != nil {
		return x.Asc
	}
	return false
}

// Allows specifying Filters for ListEvidencesRequest
type Filter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CertificationTargetId *string `protobuf:"bytes,1,opt,name=certification_target_id,json=certificationTargetId,proto3,oneof" json:"certification_target_id,omitempty"`
	ToolId                *string `protobuf:"bytes,2,opt,name=tool_id,json=toolId,proto3,oneof" json:"tool_id,omitempty"`
}

func (x *Filter) Reset() {
	*x = Filter{}
	mi := &file_api_evidence_evidence_store_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Filter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Filter) ProtoMessage() {}

func (x *Filter) ProtoReflect() protoreflect.Message {
	mi := &file_api_evidence_evidence_store_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Filter.ProtoReflect.Descriptor instead.
func (*Filter) Descriptor() ([]byte, []int) {
	return file_api_evidence_evidence_store_proto_rawDescGZIP(), []int{4}
}

func (x *Filter) GetCertificationTargetId() string {
	if x != nil && x.CertificationTargetId != nil {
		return *x.CertificationTargetId
	}
	return ""
}

func (x *Filter) GetToolId() string {
	if x != nil && x.ToolId != nil {
		return *x.ToolId
	}
	return ""
}

type ListEvidencesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Evidences     []*Evidence `protobuf:"bytes,1,rep,name=evidences,proto3" json:"evidences,omitempty"`
	NextPageToken string      `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
}

func (x *ListEvidencesResponse) Reset() {
	*x = ListEvidencesResponse{}
	mi := &file_api_evidence_evidence_store_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListEvidencesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListEvidencesResponse) ProtoMessage() {}

func (x *ListEvidencesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_evidence_evidence_store_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListEvidencesResponse.ProtoReflect.Descriptor instead.
func (*ListEvidencesResponse) Descriptor() ([]byte, []int) {
	return file_api_evidence_evidence_store_proto_rawDescGZIP(), []int{5}
}

func (x *ListEvidencesResponse) GetEvidences() []*Evidence {
	if x != nil {
		return x.Evidences
	}
	return nil
}

func (x *ListEvidencesResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

type GetEvidenceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EvidenceId string `protobuf:"bytes,1,opt,name=evidence_id,json=evidenceId,proto3" json:"evidence_id,omitempty"`
}

func (x *GetEvidenceRequest) Reset() {
	*x = GetEvidenceRequest{}
	mi := &file_api_evidence_evidence_store_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetEvidenceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetEvidenceRequest) ProtoMessage() {}

func (x *GetEvidenceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_evidence_evidence_store_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetEvidenceRequest.ProtoReflect.Descriptor instead.
func (*GetEvidenceRequest) Descriptor() ([]byte, []int) {
	return file_api_evidence_evidence_store_proto_rawDescGZIP(), []int{6}
}

func (x *GetEvidenceRequest) GetEvidenceId() string {
	if x != nil {
		return x.EvidenceId
	}
	return ""
}

var File_api_evidence_evidence_store_proto protoreflect.FileDescriptor

var file_api_evidence_evidence_store_proto_rawDesc = []byte{
	0x0a, 0x21, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x2f, 0x65,
	0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x15, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x65,
	0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x61, 0x70, 0x69, 0x2f,
	0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x2f, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x5b, 0x0a, 0x14, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x45, 0x76, 0x69, 0x64, 0x65,
	0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x43, 0x0a, 0x08, 0x65, 0x76,
	0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x63,
	0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x42, 0x06, 0xba,
	0x48, 0x03, 0xc8, 0x01, 0x01, 0x52, 0x08, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x22,
	0x17, 0x0a, 0x15, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x57, 0x0a, 0x16, 0x53, 0x74, 0x6f, 0x72,
	0x65, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x22, 0xc6, 0x01, 0x0a, 0x14, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e,
	0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3a, 0x0a, 0x06, 0x66, 0x69,
	0x6c, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x63, 0x6c, 0x6f,
	0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x48, 0x00, 0x52, 0x06, 0x66, 0x69, 0x6c,
	0x74, 0x65, 0x72, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73,
	0x69, 0x7a, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53,
	0x69, 0x7a, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x62, 0x79, 0x18, 0x0c,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x42, 0x79, 0x12, 0x10, 0x0a,
	0x03, 0x61, 0x73, 0x63, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x08, 0x52, 0x03, 0x61, 0x73, 0x63, 0x42,
	0x09, 0x0a, 0x07, 0x5f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x22, 0x9f, 0x01, 0x0a, 0x06, 0x46,
	0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x45, 0x0a, 0x17, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01,
	0x48, 0x00, 0x52, 0x15, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12, 0x26, 0x0a, 0x07,
	0x74, 0x6f, 0x6f, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba,
	0x48, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0x48, 0x01, 0x52, 0x06, 0x74, 0x6f, 0x6f, 0x6c, 0x49,
	0x64, 0x88, 0x01, 0x01, 0x42, 0x1a, 0x0a, 0x18, 0x5f, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x69, 0x64,
	0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x74, 0x6f, 0x6f, 0x6c, 0x5f, 0x69, 0x64, 0x22, 0x7e, 0x0a, 0x15,
	0x4c, 0x69, 0x73, 0x74, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3d, 0x0a, 0x09, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63,
	0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64,
	0x69, 0x74, 0x6f, 0x72, 0x2e, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x09, 0x65, 0x76, 0x69, 0x64, 0x65,
	0x6e, 0x63, 0x65, 0x73, 0x12, 0x26, 0x0a, 0x0f, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x61, 0x67,
	0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6e,
	0x65, 0x78, 0x74, 0x50, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x3f, 0x0a, 0x12,
	0x47, 0x65, 0x74, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x29, 0x0a, 0x0b, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05, 0x72, 0x03, 0xb0, 0x01,
	0x01, 0x52, 0x0a, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x49, 0x64, 0x32, 0xc2, 0x04,
	0x0a, 0x0d, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x12,
	0x99, 0x01, 0x0a, 0x0d, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63,
	0x65, 0x12, 0x2b, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x65, 0x76,
	0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x45,
	0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c,
	0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x65, 0x76, 0x69, 0x64, 0x65,
	0x6e, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x45, 0x76, 0x69, 0x64,
	0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x2d, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x27, 0x3a, 0x08, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x22, 0x1b,
	0x2f, 0x76, 0x31, 0x2f, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x5f, 0x73, 0x74, 0x6f,
	0x72, 0x65, 0x2f, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x72, 0x0a, 0x0e, 0x53,
	0x74, 0x6f, 0x72, 0x65, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x12, 0x2b, 0x2e,
	0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e,
	0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x45, 0x76, 0x69, 0x64, 0x65,
	0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2d, 0x2e, 0x63, 0x6c, 0x6f,
	0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x12,
	0x90, 0x01, 0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65,
	0x73, 0x12, 0x2b, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x65, 0x76,
	0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x76,
	0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c,
	0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x65, 0x76, 0x69, 0x64, 0x65,
	0x6e, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x76, 0x69, 0x64, 0x65,
	0x6e, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x24, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x1e, 0x12, 0x1c, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e,
	0x63, 0x65, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63,
	0x65, 0x73, 0x12, 0x8d, 0x01, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e,
	0x63, 0x65, 0x12, 0x29, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x65,
	0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x45, 0x76,
	0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e,
	0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e,
	0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x22, 0x32,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x2c, 0x12, 0x2a, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x76, 0x69, 0x64,
	0x65, 0x6e, 0x63, 0x65, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x65, 0x76, 0x69, 0x64, 0x65,
	0x6e, 0x63, 0x65, 0x73, 0x2f, 0x7b, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x5f, 0x69,
	0x64, 0x7d, 0x42, 0x28, 0x5a, 0x26, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e,
	0x69, 0x6f, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2f, 0x76, 0x32, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_evidence_evidence_store_proto_rawDescOnce sync.Once
	file_api_evidence_evidence_store_proto_rawDescData = file_api_evidence_evidence_store_proto_rawDesc
)

func file_api_evidence_evidence_store_proto_rawDescGZIP() []byte {
	file_api_evidence_evidence_store_proto_rawDescOnce.Do(func() {
		file_api_evidence_evidence_store_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_evidence_evidence_store_proto_rawDescData)
	})
	return file_api_evidence_evidence_store_proto_rawDescData
}

var file_api_evidence_evidence_store_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_api_evidence_evidence_store_proto_goTypes = []any{
	(*StoreEvidenceRequest)(nil),   // 0: clouditor.evidence.v1.StoreEvidenceRequest
	(*StoreEvidenceResponse)(nil),  // 1: clouditor.evidence.v1.StoreEvidenceResponse
	(*StoreEvidencesResponse)(nil), // 2: clouditor.evidence.v1.StoreEvidencesResponse
	(*ListEvidencesRequest)(nil),   // 3: clouditor.evidence.v1.ListEvidencesRequest
	(*Filter)(nil),                 // 4: clouditor.evidence.v1.Filter
	(*ListEvidencesResponse)(nil),  // 5: clouditor.evidence.v1.ListEvidencesResponse
	(*GetEvidenceRequest)(nil),     // 6: clouditor.evidence.v1.GetEvidenceRequest
	(*Evidence)(nil),               // 7: clouditor.evidence.v1.Evidence
}
var file_api_evidence_evidence_store_proto_depIdxs = []int32{
	7, // 0: clouditor.evidence.v1.StoreEvidenceRequest.evidence:type_name -> clouditor.evidence.v1.Evidence
	4, // 1: clouditor.evidence.v1.ListEvidencesRequest.filter:type_name -> clouditor.evidence.v1.Filter
	7, // 2: clouditor.evidence.v1.ListEvidencesResponse.evidences:type_name -> clouditor.evidence.v1.Evidence
	0, // 3: clouditor.evidence.v1.EvidenceStore.StoreEvidence:input_type -> clouditor.evidence.v1.StoreEvidenceRequest
	0, // 4: clouditor.evidence.v1.EvidenceStore.StoreEvidences:input_type -> clouditor.evidence.v1.StoreEvidenceRequest
	3, // 5: clouditor.evidence.v1.EvidenceStore.ListEvidences:input_type -> clouditor.evidence.v1.ListEvidencesRequest
	6, // 6: clouditor.evidence.v1.EvidenceStore.GetEvidence:input_type -> clouditor.evidence.v1.GetEvidenceRequest
	1, // 7: clouditor.evidence.v1.EvidenceStore.StoreEvidence:output_type -> clouditor.evidence.v1.StoreEvidenceResponse
	2, // 8: clouditor.evidence.v1.EvidenceStore.StoreEvidences:output_type -> clouditor.evidence.v1.StoreEvidencesResponse
	5, // 9: clouditor.evidence.v1.EvidenceStore.ListEvidences:output_type -> clouditor.evidence.v1.ListEvidencesResponse
	7, // 10: clouditor.evidence.v1.EvidenceStore.GetEvidence:output_type -> clouditor.evidence.v1.Evidence
	7, // [7:11] is the sub-list for method output_type
	3, // [3:7] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_api_evidence_evidence_store_proto_init() }
func file_api_evidence_evidence_store_proto_init() {
	if File_api_evidence_evidence_store_proto != nil {
		return
	}
	file_api_evidence_evidence_proto_init()
	file_api_evidence_evidence_store_proto_msgTypes[3].OneofWrappers = []any{}
	file_api_evidence_evidence_store_proto_msgTypes[4].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_evidence_evidence_store_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_evidence_evidence_store_proto_goTypes,
		DependencyIndexes: file_api_evidence_evidence_store_proto_depIdxs,
		MessageInfos:      file_api_evidence_evidence_store_proto_msgTypes,
	}.Build()
	File_api_evidence_evidence_store_proto = out.File
	file_api_evidence_evidence_store_proto_rawDesc = nil
	file_api_evidence_evidence_store_proto_goTypes = nil
	file_api_evidence_evidence_store_proto_depIdxs = nil
}
