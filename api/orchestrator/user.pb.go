// Copyright 2025 Fraunhofer AISEC
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
// source: api/orchestrator/user.proto

package orchestrator

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	_ "github.com/srikrsna/protoc-gen-gotag/tagger"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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

// A user resource
// TODO(lebogg): Think about adding more fields to the user resource (creation time, expiration time, etc.)
type User struct {
	state     protoimpl.MessageState `protogen:"open.v1"`
	Id        string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	FirstName string                 `protobuf:"bytes,2,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName  string                 `protobuf:"bytes,3,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	Email     string                 `protobuf:"bytes,4,opt,name=email,proto3" json:"email,omitempty"`
	// target_of_certification_ids defines the scope of the user
	TargetOfEvaluationIds []*TargetOfEvaluation `protobuf:"bytes,5,rep,name=target_of_evaluation_ids,json=targetOfEvaluationIds,proto3" json:"target_of_evaluation_ids,omitempty" gorm:"many2many:user_certification_target;"`
	unknownFields         protoimpl.UnknownFields
	sizeCache             protoimpl.SizeCache
}

func (x *User) Reset() {
	*x = User{}
	mi := &file_api_orchestrator_user_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_api_orchestrator_user_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_api_orchestrator_user_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *User) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *User) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

func (x *User) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *User) GetTargetOfEvaluationIds() []*TargetOfEvaluation {
	if x != nil {
		return x.TargetOfEvaluationIds
	}
	return nil
}

var File_api_orchestrator_user_proto protoreflect.FileDescriptor

const file_api_orchestrator_user_proto_rawDesc = "" +
	"\n" +
	"\x1bapi/orchestrator/user.proto\x12\x19clouditor.orchestrator.v1\x1a#api/orchestrator/orchestrator.proto\x1a\x1bbuf/validate/validate.proto\x1a\x1cgoogle/api/annotations.proto\x1a\x1fgoogle/api/field_behavior.proto\x1a\x13tagger/tagger.proto\"\x8f\x02\n" +
	"\x04User\x12\x1a\n" +
	"\x02id\x18\x01 \x01(\tB\n" +
	"\xe0A\x02\xbaH\x04r\x02\x10\x01R\x02id\x12\x1d\n" +
	"\n" +
	"first_name\x18\x02 \x01(\tR\tfirstName\x12\x1b\n" +
	"\tlast_name\x18\x03 \x01(\tR\blastName\x12\x14\n" +
	"\x05email\x18\x04 \x01(\tR\x05email\x12\x98\x01\n" +
	"\x18target_of_evaluation_ids\x18\x05 \x03(\v2-.clouditor.orchestrator.v1.TargetOfEvaluationB0\x9a\x84\x9e\x03+gorm:\"many2many:user_certification_target;\"R\x15targetOfEvaluationIdsB,Z*clouditor.io/clouditor/v2/api/orchestratorb\x06proto3"

var (
	file_api_orchestrator_user_proto_rawDescOnce sync.Once
	file_api_orchestrator_user_proto_rawDescData []byte
)

func file_api_orchestrator_user_proto_rawDescGZIP() []byte {
	file_api_orchestrator_user_proto_rawDescOnce.Do(func() {
		file_api_orchestrator_user_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_orchestrator_user_proto_rawDesc), len(file_api_orchestrator_user_proto_rawDesc)))
	})
	return file_api_orchestrator_user_proto_rawDescData
}

var file_api_orchestrator_user_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_api_orchestrator_user_proto_goTypes = []any{
	(*User)(nil),               // 0: clouditor.orchestrator.v1.User
	(*TargetOfEvaluation)(nil), // 1: clouditor.orchestrator.v1.TargetOfEvaluation
}
var file_api_orchestrator_user_proto_depIdxs = []int32{
	1, // 0: clouditor.orchestrator.v1.User.target_of_evaluation_ids:type_name -> clouditor.orchestrator.v1.TargetOfEvaluation
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_orchestrator_user_proto_init() }
func file_api_orchestrator_user_proto_init() {
	if File_api_orchestrator_user_proto != nil {
		return
	}
	file_api_orchestrator_orchestrator_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_orchestrator_user_proto_rawDesc), len(file_api_orchestrator_user_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_orchestrator_user_proto_goTypes,
		DependencyIndexes: file_api_orchestrator_user_proto_depIdxs,
		MessageInfos:      file_api_orchestrator_user_proto_msgTypes,
	}.Build()
	File_api_orchestrator_user_proto = out.File
	file_api_orchestrator_user_proto_goTypes = nil
	file_api_orchestrator_user_proto_depIdxs = nil
}
