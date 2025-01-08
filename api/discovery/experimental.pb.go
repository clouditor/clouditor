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
// 	protoc-gen-go v1.36.2
// 	protoc        (unknown)
// source: api/discovery/experimental.proto

package discovery

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

type UpdateResourceRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Resource      *Resource              `protobuf:"bytes,1,opt,name=resource,proto3" json:"resource,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateResourceRequest) Reset() {
	*x = UpdateResourceRequest{}
	mi := &file_api_discovery_experimental_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateResourceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateResourceRequest) ProtoMessage() {}

func (x *UpdateResourceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_discovery_experimental_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateResourceRequest.ProtoReflect.Descriptor instead.
func (*UpdateResourceRequest) Descriptor() ([]byte, []int) {
	return file_api_discovery_experimental_proto_rawDescGZIP(), []int{0}
}

func (x *UpdateResourceRequest) GetResource() *Resource {
	if x != nil {
		return x.Resource
	}
	return nil
}

type ListGraphEdgesRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	PageSize      int32                  `protobuf:"varint,10,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	PageToken     string                 `protobuf:"bytes,11,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	OrderBy       string                 `protobuf:"bytes,12,opt,name=order_by,json=orderBy,proto3" json:"order_by,omitempty"`
	Asc           bool                   `protobuf:"varint,13,opt,name=asc,proto3" json:"asc,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListGraphEdgesRequest) Reset() {
	*x = ListGraphEdgesRequest{}
	mi := &file_api_discovery_experimental_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListGraphEdgesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListGraphEdgesRequest) ProtoMessage() {}

func (x *ListGraphEdgesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_discovery_experimental_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListGraphEdgesRequest.ProtoReflect.Descriptor instead.
func (*ListGraphEdgesRequest) Descriptor() ([]byte, []int) {
	return file_api_discovery_experimental_proto_rawDescGZIP(), []int{1}
}

func (x *ListGraphEdgesRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListGraphEdgesRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

func (x *ListGraphEdgesRequest) GetOrderBy() string {
	if x != nil {
		return x.OrderBy
	}
	return ""
}

func (x *ListGraphEdgesRequest) GetAsc() bool {
	if x != nil {
		return x.Asc
	}
	return false
}

type ListGraphEdgesResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Edges         []*GraphEdge           `protobuf:"bytes,1,rep,name=edges,proto3" json:"edges,omitempty"`
	NextPageToken string                 `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListGraphEdgesResponse) Reset() {
	*x = ListGraphEdgesResponse{}
	mi := &file_api_discovery_experimental_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListGraphEdgesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListGraphEdgesResponse) ProtoMessage() {}

func (x *ListGraphEdgesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_discovery_experimental_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListGraphEdgesResponse.ProtoReflect.Descriptor instead.
func (*ListGraphEdgesResponse) Descriptor() ([]byte, []int) {
	return file_api_discovery_experimental_proto_rawDescGZIP(), []int{2}
}

func (x *ListGraphEdgesResponse) GetEdges() []*GraphEdge {
	if x != nil {
		return x.Edges
	}
	return nil
}

func (x *ListGraphEdgesResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

type GraphEdge struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Source        string                 `protobuf:"bytes,2,opt,name=source,proto3" json:"source,omitempty"`
	Target        string                 `protobuf:"bytes,3,opt,name=target,proto3" json:"target,omitempty"`
	Type          string                 `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GraphEdge) Reset() {
	*x = GraphEdge{}
	mi := &file_api_discovery_experimental_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GraphEdge) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GraphEdge) ProtoMessage() {}

func (x *GraphEdge) ProtoReflect() protoreflect.Message {
	mi := &file_api_discovery_experimental_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GraphEdge.ProtoReflect.Descriptor instead.
func (*GraphEdge) Descriptor() ([]byte, []int) {
	return file_api_discovery_experimental_proto_rawDescGZIP(), []int{3}
}

func (x *GraphEdge) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *GraphEdge) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *GraphEdge) GetTarget() string {
	if x != nil {
		return x.Target
	}
	return ""
}

func (x *GraphEdge) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

var File_api_discovery_experimental_proto protoreflect.FileDescriptor

var file_api_discovery_experimental_proto_rawDesc = []byte{
	0x0a, 0x20, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2f,
	0x65, 0x78, 0x70, 0x65, 0x72, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x22, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x64, 0x69,
	0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x65, 0x78, 0x70, 0x65, 0x72, 0x69,
	0x6d, 0x65, 0x6e, 0x74, 0x61, 0x6c, 0x1a, 0x1d, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x69, 0x73, 0x63,
	0x6f, 0x76, 0x65, 0x72, 0x79, 0x2f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61,
	0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65,
	0x6c, 0x64, 0x5f, 0x62, 0x65, 0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x5a, 0x0a, 0x15, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x41, 0x0a, 0x08, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x63,
	0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65,
	0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x42, 0x03,
	0xe0, 0x41, 0x02, 0x52, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x22, 0x80, 0x01,
	0x0a, 0x15, 0x4c, 0x69, 0x73, 0x74, 0x47, 0x72, 0x61, 0x70, 0x68, 0x45, 0x64, 0x67, 0x65, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f,
	0x73, 0x69, 0x7a, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65,
	0x53, 0x69, 0x7a, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x61, 0x67, 0x65, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x62, 0x79, 0x18,
	0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x42, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x61, 0x73, 0x63, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x08, 0x52, 0x03, 0x61, 0x73, 0x63,
	0x22, 0x8a, 0x01, 0x0a, 0x16, 0x4c, 0x69, 0x73, 0x74, 0x47, 0x72, 0x61, 0x70, 0x68, 0x45, 0x64,
	0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x48, 0x0a, 0x05, 0x65,
	0x64, 0x67, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x63, 0x6c, 0x6f,
	0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79,
	0x2e, 0x76, 0x31, 0x65, 0x78, 0x70, 0x65, 0x72, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x61, 0x6c, 0x2e,
	0x47, 0x72, 0x61, 0x70, 0x68, 0x45, 0x64, 0x67, 0x65, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x05,
	0x65, 0x64, 0x67, 0x65, 0x73, 0x12, 0x26, 0x0a, 0x0f, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x61,
	0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x6e, 0x65, 0x78, 0x74, 0x50, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x88, 0x01,
	0x0a, 0x09, 0x47, 0x72, 0x61, 0x70, 0x68, 0x45, 0x64, 0x67, 0x65, 0x12, 0x1a, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xe0, 0x41, 0x02, 0xba, 0x48, 0x04, 0x72,
	0x02, 0x10, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12, 0x22, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xe0, 0x41, 0x02, 0xba, 0x48, 0x04, 0x72,
	0x02, 0x10, 0x01, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x22, 0x0a, 0x06, 0x74,
	0x61, 0x72, 0x67, 0x65, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xe0, 0x41, 0x02,
	0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x12,
	0x17, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0xe0,
	0x41, 0x02, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x32, 0xfe, 0x02, 0x0a, 0x15, 0x45, 0x78, 0x70,
	0x65, 0x72, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x61, 0x6c, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65,
	0x72, 0x79, 0x12, 0xab, 0x01, 0x0a, 0x0e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x39, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f,
	0x72, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x65, 0x78,
	0x70, 0x65, 0x72, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x61, 0x6c, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x20, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x64, 0x69, 0x73,
	0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x22, 0x3c, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x36, 0x3a, 0x01, 0x2a, 0x22, 0x31, 0x2f,
	0x76, 0x31, 0x65, 0x78, 0x70, 0x65, 0x72, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x61, 0x6c, 0x2f, 0x64,
	0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2f, 0x7b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x69, 0x64, 0x7d,
	0x12, 0xb6, 0x01, 0x0a, 0x0e, 0x4c, 0x69, 0x73, 0x74, 0x47, 0x72, 0x61, 0x70, 0x68, 0x45, 0x64,
	0x67, 0x65, 0x73, 0x12, 0x39, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e,
	0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x65, 0x78, 0x70, 0x65,
	0x72, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x61, 0x6c, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x47, 0x72, 0x61,
	0x70, 0x68, 0x45, 0x64, 0x67, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x3a,
	0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x65, 0x78, 0x70, 0x65, 0x72, 0x69, 0x6d, 0x65, 0x6e,
	0x74, 0x61, 0x6c, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x47, 0x72, 0x61, 0x70, 0x68, 0x45, 0x64, 0x67,
	0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x2d, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x27, 0x12, 0x25, 0x2f, 0x76, 0x31, 0x65, 0x78, 0x70, 0x65, 0x72, 0x69, 0x6d, 0x65, 0x6e,
	0x74, 0x61, 0x6c, 0x2f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2f, 0x67, 0x72,
	0x61, 0x70, 0x68, 0x2f, 0x65, 0x64, 0x67, 0x65, 0x73, 0x42, 0x29, 0x5a, 0x27, 0x63, 0x6c, 0x6f,
	0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x69, 0x6f, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69,
	0x74, 0x6f, 0x72, 0x2f, 0x76, 0x32, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x69, 0x73, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_discovery_experimental_proto_rawDescOnce sync.Once
	file_api_discovery_experimental_proto_rawDescData = file_api_discovery_experimental_proto_rawDesc
)

func file_api_discovery_experimental_proto_rawDescGZIP() []byte {
	file_api_discovery_experimental_proto_rawDescOnce.Do(func() {
		file_api_discovery_experimental_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_discovery_experimental_proto_rawDescData)
	})
	return file_api_discovery_experimental_proto_rawDescData
}

var file_api_discovery_experimental_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_api_discovery_experimental_proto_goTypes = []any{
	(*UpdateResourceRequest)(nil),  // 0: clouditor.discovery.v1experimental.UpdateResourceRequest
	(*ListGraphEdgesRequest)(nil),  // 1: clouditor.discovery.v1experimental.ListGraphEdgesRequest
	(*ListGraphEdgesResponse)(nil), // 2: clouditor.discovery.v1experimental.ListGraphEdgesResponse
	(*GraphEdge)(nil),              // 3: clouditor.discovery.v1experimental.GraphEdge
	(*Resource)(nil),               // 4: clouditor.discovery.v1.Resource
}
var file_api_discovery_experimental_proto_depIdxs = []int32{
	4, // 0: clouditor.discovery.v1experimental.UpdateResourceRequest.resource:type_name -> clouditor.discovery.v1.Resource
	3, // 1: clouditor.discovery.v1experimental.ListGraphEdgesResponse.edges:type_name -> clouditor.discovery.v1experimental.GraphEdge
	0, // 2: clouditor.discovery.v1experimental.ExperimentalDiscovery.UpdateResource:input_type -> clouditor.discovery.v1experimental.UpdateResourceRequest
	1, // 3: clouditor.discovery.v1experimental.ExperimentalDiscovery.ListGraphEdges:input_type -> clouditor.discovery.v1experimental.ListGraphEdgesRequest
	4, // 4: clouditor.discovery.v1experimental.ExperimentalDiscovery.UpdateResource:output_type -> clouditor.discovery.v1.Resource
	2, // 5: clouditor.discovery.v1experimental.ExperimentalDiscovery.ListGraphEdges:output_type -> clouditor.discovery.v1experimental.ListGraphEdgesResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_discovery_experimental_proto_init() }
func file_api_discovery_experimental_proto_init() {
	if File_api_discovery_experimental_proto != nil {
		return
	}
	file_api_discovery_discovery_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_discovery_experimental_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_discovery_experimental_proto_goTypes,
		DependencyIndexes: file_api_discovery_experimental_proto_depIdxs,
		MessageInfos:      file_api_discovery_experimental_proto_msgTypes,
	}.Build()
	File_api_discovery_experimental_proto = out.File
	file_api_discovery_experimental_proto_rawDesc = nil
	file_api_discovery_experimental_proto_goTypes = nil
	file_api_discovery_experimental_proto_depIdxs = nil
}
