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
// 	protoc-gen-go v1.36.4
// 	protoc        (unknown)
// source: api/evidence/evidence.proto

package evidence

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	_ "github.com/srikrsna/protoc-gen-gotag/tagger"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
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

// An evidence resource
type Evidence struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// the ID in a uuid format
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// time of evidence creation
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty" gorm:"serializer:timestamppb;type:timestamp"`
	// Reference to a certification target (e.g., service, organization) this evidence was gathered from
	CertificationTargetId string `protobuf:"bytes,3,opt,name=certification_target_id,json=certificationTargetId,proto3" json:"certification_target_id,omitempty"`
	// Reference to the tool which provided the evidence
	ToolId string `protobuf:"bytes,4,opt,name=tool_id,json=toolId,proto3" json:"tool_id,omitempty"`
	// Optional. Contains the evidence in its original form without following a
	// defined schema, e.g. the raw JSON
	Raw *string `protobuf:"bytes,5,opt,name=raw,proto3,oneof" json:"raw,omitempty"`
	// Semantic representation of the Cloud resource according to our defined
	// ontology
	Resource *anypb.Any `protobuf:"bytes,6,opt,name=resource,proto3" json:"resource,omitempty" gorm:"serializer:anypb;type:json"`
	// Very experimental property. Use at own risk. This property will be deleted again.
	//
	// Related resource IDs. The assessment will wait until all evidences for related resource have arrived in the
	// assessment and are recent enough. In the future, this will be replaced with information in the "related" edges in
	// the resource. For now, this needs to be set manually in the evidence.
	ExperimentalRelatedResourceIds []string `protobuf:"bytes,999,rep,name=experimental_related_resource_ids,json=experimentalRelatedResourceIds,proto3" json:"experimental_related_resource_ids,omitempty" gorm:"serializer:json"`
	unknownFields                  protoimpl.UnknownFields
	sizeCache                      protoimpl.SizeCache
}

func (x *Evidence) Reset() {
	*x = Evidence{}
	mi := &file_api_evidence_evidence_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Evidence) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Evidence) ProtoMessage() {}

func (x *Evidence) ProtoReflect() protoreflect.Message {
	mi := &file_api_evidence_evidence_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Evidence.ProtoReflect.Descriptor instead.
func (*Evidence) Descriptor() ([]byte, []int) {
	return file_api_evidence_evidence_proto_rawDescGZIP(), []int{0}
}

func (x *Evidence) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Evidence) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *Evidence) GetCertificationTargetId() string {
	if x != nil {
		return x.CertificationTargetId
	}
	return ""
}

func (x *Evidence) GetToolId() string {
	if x != nil {
		return x.ToolId
	}
	return ""
}

func (x *Evidence) GetRaw() string {
	if x != nil && x.Raw != nil {
		return *x.Raw
	}
	return ""
}

func (x *Evidence) GetResource() *anypb.Any {
	if x != nil {
		return x.Resource
	}
	return nil
}

func (x *Evidence) GetExperimentalRelatedResourceIds() []string {
	if x != nil {
		return x.ExperimentalRelatedResourceIds
	}
	return nil
}

var File_api_evidence_evidence_proto protoreflect.FileDescriptor

var file_api_evidence_evidence_proto_rawDesc = string([]byte{
	0x0a, 0x1b, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x2f, 0x65,
	0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x63,
	0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63,
	0x65, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13, 0x74,
	0x61, 0x67, 0x67, 0x65, 0x72, 0x2f, 0x74, 0x61, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xec, 0x03, 0x0a, 0x08, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x12,
	0x18, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05,
	0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12, 0x71, 0x0a, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x37, 0xba, 0x48, 0x03, 0xc8, 0x01, 0x01,
	0x9a, 0x84, 0x9e, 0x03, 0x2c, 0x67, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x73, 0x65, 0x72, 0x69, 0x61,
	0x6c, 0x69, 0x7a, 0x65, 0x72, 0x3a, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x70,
	0x62, 0x3b, 0x74, 0x79, 0x70, 0x65, 0x3a, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x22, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x40, 0x0a, 0x17,
	0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba,
	0x48, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x15, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x49, 0x64, 0x12, 0x20,
	0x0a, 0x07, 0x74, 0x6f, 0x6f, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x06, 0x74, 0x6f, 0x6f, 0x6c, 0x49, 0x64,
	0x12, 0x1e, 0x0a, 0x03, 0x72, 0x61, 0x77, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba,
	0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x48, 0x00, 0x52, 0x03, 0x72, 0x61, 0x77, 0x88, 0x01, 0x01,
	0x12, 0x5e, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x42, 0x2c, 0xba, 0x48, 0x03, 0xc8, 0x01, 0x01,
	0x9a, 0x84, 0x9e, 0x03, 0x21, 0x67, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x73, 0x65, 0x72, 0x69, 0x61,
	0x6c, 0x69, 0x7a, 0x65, 0x72, 0x3a, 0x61, 0x6e, 0x79, 0x70, 0x62, 0x3b, 0x74, 0x79, 0x70, 0x65,
	0x3a, 0x6a, 0x73, 0x6f, 0x6e, 0x22, 0x52, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x12, 0x67, 0x0a, 0x21, 0x65, 0x78, 0x70, 0x65, 0x72, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x61, 0x6c,
	0x5f, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x5f, 0x69, 0x64, 0x73, 0x18, 0xe7, 0x07, 0x20, 0x03, 0x28, 0x09, 0x42, 0x1b, 0x9a, 0x84,
	0x9e, 0x03, 0x16, 0x67, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69,
	0x7a, 0x65, 0x72, 0x3a, 0x6a, 0x73, 0x6f, 0x6e, 0x22, 0x52, 0x1e, 0x65, 0x78, 0x70, 0x65, 0x72,
	0x69, 0x6d, 0x65, 0x6e, 0x74, 0x61, 0x6c, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x65, 0x64, 0x52, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x73, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x72, 0x61,
	0x77, 0x42, 0x28, 0x5a, 0x26, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x69,
	0x6f, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2f, 0x76, 0x32, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
})

var (
	file_api_evidence_evidence_proto_rawDescOnce sync.Once
	file_api_evidence_evidence_proto_rawDescData []byte
)

func file_api_evidence_evidence_proto_rawDescGZIP() []byte {
	file_api_evidence_evidence_proto_rawDescOnce.Do(func() {
		file_api_evidence_evidence_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_evidence_evidence_proto_rawDesc), len(file_api_evidence_evidence_proto_rawDesc)))
	})
	return file_api_evidence_evidence_proto_rawDescData
}

var file_api_evidence_evidence_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_api_evidence_evidence_proto_goTypes = []any{
	(*Evidence)(nil),              // 0: clouditor.evidence.v1.Evidence
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
	(*anypb.Any)(nil),             // 2: google.protobuf.Any
}
var file_api_evidence_evidence_proto_depIdxs = []int32{
	1, // 0: clouditor.evidence.v1.Evidence.timestamp:type_name -> google.protobuf.Timestamp
	2, // 1: clouditor.evidence.v1.Evidence.resource:type_name -> google.protobuf.Any
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_evidence_evidence_proto_init() }
func file_api_evidence_evidence_proto_init() {
	if File_api_evidence_evidence_proto != nil {
		return
	}
	file_api_evidence_evidence_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_evidence_evidence_proto_rawDesc), len(file_api_evidence_evidence_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_evidence_evidence_proto_goTypes,
		DependencyIndexes: file_api_evidence_evidence_proto_depIdxs,
		MessageInfos:      file_api_evidence_evidence_proto_msgTypes,
	}.Build()
	File_api_evidence_evidence_proto = out.File
	file_api_evidence_evidence_proto_goTypes = nil
	file_api_evidence_evidence_proto_depIdxs = nil
}
