//
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
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: api/evidence/evidence_store.proto

package evidence

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
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
	if protoimpl.UnsafeEnabled {
		mi := &file_api_evidence_evidence_store_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreEvidenceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreEvidenceRequest) ProtoMessage() {}

func (x *StoreEvidenceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_evidence_evidence_store_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
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

type StoreEvidenceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status        bool   `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	StatusMessage string `protobuf:"bytes,2,opt,name=status_message,json=statusMessage,proto3" json:"status_message,omitempty"`
}

func (x *StoreEvidenceResponse) Reset() {
	*x = StoreEvidenceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_evidence_evidence_store_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreEvidenceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreEvidenceResponse) ProtoMessage() {}

func (x *StoreEvidenceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_evidence_evidence_store_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
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

func (x *StoreEvidenceResponse) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

func (x *StoreEvidenceResponse) GetStatusMessage() string {
	if x != nil {
		return x.StatusMessage
	}
	return ""
}

type ListEvidencesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PageSize  int32  `protobuf:"varint,1,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	PageToken string `protobuf:"bytes,2,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	OrderBy   string `protobuf:"bytes,3,opt,name=order_by,json=orderBy,proto3" json:"order_by,omitempty"`
	Asc       bool   `protobuf:"varint,4,opt,name=asc,proto3" json:"asc,omitempty"`
}

func (x *ListEvidencesRequest) Reset() {
	*x = ListEvidencesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_evidence_evidence_store_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListEvidencesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListEvidencesRequest) ProtoMessage() {}

func (x *ListEvidencesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_evidence_evidence_store_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
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
	return file_api_evidence_evidence_store_proto_rawDescGZIP(), []int{2}
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

type ListEvidencesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Evidences     []*Evidence `protobuf:"bytes,1,rep,name=evidences,proto3" json:"evidences,omitempty"`
	NextPageToken string      `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
}

func (x *ListEvidencesResponse) Reset() {
	*x = ListEvidencesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_evidence_evidence_store_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListEvidencesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListEvidencesResponse) ProtoMessage() {}

func (x *ListEvidencesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_evidence_evidence_store_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
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
	return file_api_evidence_evidence_store_proto_rawDescGZIP(), []int{3}
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

var File_api_evidence_evidence_store_proto protoreflect.FileDescriptor

var file_api_evidence_evidence_store_proto_rawDesc = []byte{
	0x0a, 0x21, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x2f, 0x65,
	0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x15, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x65,
	0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x61, 0x70, 0x69, 0x2f,
	0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x2f, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5d,
	0x0a, 0x14, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x45, 0x0a, 0x08, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e,
	0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64,
	0x69, 0x74, 0x6f, 0x72, 0x2e, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x8a, 0x01,
	0x02, 0x10, 0x01, 0x52, 0x08, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x22, 0x56, 0x0a,
	0x15, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x25,
	0x0a, 0x0e, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x7f, 0x0a, 0x14, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x76, 0x69,
	0x64, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a,
	0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x61,
	0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x70, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64,
	0x65, 0x72, 0x5f, 0x62, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x72, 0x64,
	0x65, 0x72, 0x42, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x73, 0x63, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x03, 0x61, 0x73, 0x63, 0x22, 0x7e, 0x0a, 0x15, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x76,
	0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x3d, 0x0a, 0x09, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x65,
	0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x69, 0x64, 0x65,
	0x6e, 0x63, 0x65, 0x52, 0x09, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x12, 0x26,
	0x0a, 0x0f, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6e, 0x65, 0x78, 0x74, 0x50, 0x61, 0x67,
	0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x32, 0xb1, 0x03, 0x0a, 0x0d, 0x45, 0x76, 0x69, 0x64, 0x65,
	0x6e, 0x63, 0x65, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x12, 0x99, 0x01, 0x0a, 0x0d, 0x53, 0x74, 0x6f,
	0x72, 0x65, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x2b, 0x2e, 0x63, 0x6c, 0x6f,
	0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69,
	0x74, 0x6f, 0x72, 0x2e, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x53, 0x74, 0x6f, 0x72, 0x65, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x2d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x27, 0x3a, 0x08, 0x65,
	0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x22, 0x1b, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x76, 0x69,
	0x64, 0x65, 0x6e, 0x63, 0x65, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x65, 0x76, 0x69, 0x64,
	0x65, 0x6e, 0x63, 0x65, 0x12, 0x71, 0x0a, 0x0e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x45, 0x76, 0x69,
	0x64, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x12, 0x2b, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74,
	0x6f, 0x72, 0x2e, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53,
	0x74, 0x6f, 0x72, 0x65, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e,
	0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x6f, 0x72,
	0x65, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x12, 0x90, 0x01, 0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74,
	0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x12, 0x2b, 0x2e, 0x63, 0x6c, 0x6f, 0x75,
	0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74,
	0x6f, 0x72, 0x2e, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4c,
	0x69, 0x73, 0x74, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x24, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1e, 0x12, 0x1c, 0x2f, 0x76,
	0x31, 0x2f, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x2f, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x42, 0x25, 0x5a, 0x23, 0x63, 0x6c,
	0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x69, 0x6f, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64,
	0x69, 0x74, 0x6f, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63,
	0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_api_evidence_evidence_store_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_api_evidence_evidence_store_proto_goTypes = []interface{}{
	(*StoreEvidenceRequest)(nil),  // 0: clouditor.evidence.v1.StoreEvidenceRequest
	(*StoreEvidenceResponse)(nil), // 1: clouditor.evidence.v1.StoreEvidenceResponse
	(*ListEvidencesRequest)(nil),  // 2: clouditor.evidence.v1.ListEvidencesRequest
	(*ListEvidencesResponse)(nil), // 3: clouditor.evidence.v1.ListEvidencesResponse
	(*Evidence)(nil),              // 4: clouditor.evidence.v1.Evidence
}
var file_api_evidence_evidence_store_proto_depIdxs = []int32{
	4, // 0: clouditor.evidence.v1.StoreEvidenceRequest.evidence:type_name -> clouditor.evidence.v1.Evidence
	4, // 1: clouditor.evidence.v1.ListEvidencesResponse.evidences:type_name -> clouditor.evidence.v1.Evidence
	0, // 2: clouditor.evidence.v1.EvidenceStore.StoreEvidence:input_type -> clouditor.evidence.v1.StoreEvidenceRequest
	0, // 3: clouditor.evidence.v1.EvidenceStore.StoreEvidences:input_type -> clouditor.evidence.v1.StoreEvidenceRequest
	2, // 4: clouditor.evidence.v1.EvidenceStore.ListEvidences:input_type -> clouditor.evidence.v1.ListEvidencesRequest
	1, // 5: clouditor.evidence.v1.EvidenceStore.StoreEvidence:output_type -> clouditor.evidence.v1.StoreEvidenceResponse
	1, // 6: clouditor.evidence.v1.EvidenceStore.StoreEvidences:output_type -> clouditor.evidence.v1.StoreEvidenceResponse
	3, // 7: clouditor.evidence.v1.EvidenceStore.ListEvidences:output_type -> clouditor.evidence.v1.ListEvidencesResponse
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_evidence_evidence_store_proto_init() }
func file_api_evidence_evidence_store_proto_init() {
	if File_api_evidence_evidence_store_proto != nil {
		return
	}
	file_api_evidence_evidence_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_api_evidence_evidence_store_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoreEvidenceRequest); i {
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
		file_api_evidence_evidence_store_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoreEvidenceResponse); i {
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
		file_api_evidence_evidence_store_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListEvidencesRequest); i {
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
		file_api_evidence_evidence_store_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListEvidencesResponse); i {
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
			RawDescriptor: file_api_evidence_evidence_store_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
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
