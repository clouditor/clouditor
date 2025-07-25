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
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: api/discovery/discovery.proto

package discovery

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	_ "github.com/srikrsna/protoc-gen-gotag/tagger"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type StartDiscoveryRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ResourceGroup *string                `protobuf:"bytes,1,opt,name=resource_group,json=resourceGroup,proto3,oneof" json:"resource_group,omitempty"`
	CsafDomain    *string                `protobuf:"bytes,2,opt,name=csaf_domain,json=csafDomain,proto3,oneof" json:"csaf_domain,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StartDiscoveryRequest) Reset() {
	*x = StartDiscoveryRequest{}
	mi := &file_api_discovery_discovery_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StartDiscoveryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartDiscoveryRequest) ProtoMessage() {}

func (x *StartDiscoveryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_discovery_discovery_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartDiscoveryRequest.ProtoReflect.Descriptor instead.
func (*StartDiscoveryRequest) Descriptor() ([]byte, []int) {
	return file_api_discovery_discovery_proto_rawDescGZIP(), []int{0}
}

func (x *StartDiscoveryRequest) GetResourceGroup() string {
	if x != nil && x.ResourceGroup != nil {
		return *x.ResourceGroup
	}
	return ""
}

func (x *StartDiscoveryRequest) GetCsafDomain() string {
	if x != nil && x.CsafDomain != nil {
		return *x.CsafDomain
	}
	return ""
}

type StartDiscoveryResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Successful    bool                   `protobuf:"varint,1,opt,name=successful,proto3" json:"successful,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StartDiscoveryResponse) Reset() {
	*x = StartDiscoveryResponse{}
	mi := &file_api_discovery_discovery_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StartDiscoveryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartDiscoveryResponse) ProtoMessage() {}

func (x *StartDiscoveryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_discovery_discovery_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartDiscoveryResponse.ProtoReflect.Descriptor instead.
func (*StartDiscoveryResponse) Descriptor() ([]byte, []int) {
	return file_api_discovery_discovery_proto_rawDescGZIP(), []int{1}
}

func (x *StartDiscoveryResponse) GetSuccessful() bool {
	if x != nil {
		return x.Successful
	}
	return false
}

var File_api_discovery_discovery_proto protoreflect.FileDescriptor

const file_api_discovery_discovery_proto_rawDesc = "" +
	"\n" +
	"\x1dapi/discovery/discovery.proto\x12\x16clouditor.discovery.v1\x1a\x1bbuf/validate/validate.proto\x1a\x1cgoogle/api/annotations.proto\x1a\x19google/protobuf/any.proto\x1a\x13tagger/tagger.proto\"\x8c\x01\n" +
	"\x15StartDiscoveryRequest\x12*\n" +
	"\x0eresource_group\x18\x01 \x01(\tH\x00R\rresourceGroup\x88\x01\x01\x12$\n" +
	"\vcsaf_domain\x18\x02 \x01(\tH\x01R\n" +
	"csafDomain\x88\x01\x01B\x11\n" +
	"\x0f_resource_groupB\x0e\n" +
	"\f_csaf_domain\"8\n" +
	"\x16StartDiscoveryResponse\x12\x1e\n" +
	"\n" +
	"successful\x18\x01 \x01(\bR\n" +
	"successful2\x97\x01\n" +
	"\tDiscovery\x12\x89\x01\n" +
	"\x05Start\x12-.clouditor.discovery.v1.StartDiscoveryRequest\x1a..clouditor.discovery.v1.StartDiscoveryResponse\"!\x82\xd3\xe4\x93\x02\x1b:\x01*b\x01*\"\x13/v1/discovery/startB)Z'clouditor.io/clouditor/v2/api/discoveryb\x06proto3"

var (
	file_api_discovery_discovery_proto_rawDescOnce sync.Once
	file_api_discovery_discovery_proto_rawDescData []byte
)

func file_api_discovery_discovery_proto_rawDescGZIP() []byte {
	file_api_discovery_discovery_proto_rawDescOnce.Do(func() {
		file_api_discovery_discovery_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_discovery_discovery_proto_rawDesc), len(file_api_discovery_discovery_proto_rawDesc)))
	})
	return file_api_discovery_discovery_proto_rawDescData
}

var file_api_discovery_discovery_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_discovery_discovery_proto_goTypes = []any{
	(*StartDiscoveryRequest)(nil),  // 0: clouditor.discovery.v1.StartDiscoveryRequest
	(*StartDiscoveryResponse)(nil), // 1: clouditor.discovery.v1.StartDiscoveryResponse
}
var file_api_discovery_discovery_proto_depIdxs = []int32{
	0, // 0: clouditor.discovery.v1.Discovery.Start:input_type -> clouditor.discovery.v1.StartDiscoveryRequest
	1, // 1: clouditor.discovery.v1.Discovery.Start:output_type -> clouditor.discovery.v1.StartDiscoveryResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_discovery_discovery_proto_init() }
func file_api_discovery_discovery_proto_init() {
	if File_api_discovery_discovery_proto != nil {
		return
	}
	file_api_discovery_discovery_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_discovery_discovery_proto_rawDesc), len(file_api_discovery_discovery_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_discovery_discovery_proto_goTypes,
		DependencyIndexes: file_api_discovery_discovery_proto_depIdxs,
		MessageInfos:      file_api_discovery_discovery_proto_msgTypes,
	}.Build()
	File_api_discovery_discovery_proto = out.File
	file_api_discovery_discovery_proto_goTypes = nil
	file_api_discovery_discovery_proto_depIdxs = nil
}
