//
// Copyright 2021 Fraunhofer AISEC
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
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.8
// source: orchestrator.proto

package orchestrator

import (
	assessment "clouditor.io/clouditor/api/assessment"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type AssessmentResult_ComplianceStatus int32

const (
	AssessmentResult_COMPLIANT     AssessmentResult_ComplianceStatus = 0
	AssessmentResult_NON_COMPLIANT AssessmentResult_ComplianceStatus = 1
)

// Enum value maps for AssessmentResult_ComplianceStatus.
var (
	AssessmentResult_ComplianceStatus_name = map[int32]string{
		0: "COMPLIANT",
		1: "NON_COMPLIANT",
	}
	AssessmentResult_ComplianceStatus_value = map[string]int32{
		"COMPLIANT":     0,
		"NON_COMPLIANT": 1,
	}
)

func (x AssessmentResult_ComplianceStatus) Enum() *AssessmentResult_ComplianceStatus {
	p := new(AssessmentResult_ComplianceStatus)
	*p = x
	return p
}

func (x AssessmentResult_ComplianceStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AssessmentResult_ComplianceStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_orchestrator_proto_enumTypes[0].Descriptor()
}

func (AssessmentResult_ComplianceStatus) Type() protoreflect.EnumType {
	return &file_orchestrator_proto_enumTypes[0]
}

func (x AssessmentResult_ComplianceStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AssessmentResult_ComplianceStatus.Descriptor instead.
func (AssessmentResult_ComplianceStatus) EnumDescriptor() ([]byte, []int) {
	return file_orchestrator_proto_rawDescGZIP(), []int{0, 0}
}

type AssessmentResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// the ID in a uuid format
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// the ID of the metric it refers to
	MetricId    string                            `protobuf:"bytes,2,opt,name=metric_id,json=metricId,proto3" json:"metric_id,omitempty"`
	Result      AssessmentResult_ComplianceStatus `protobuf:"varint,3,opt,name=result,proto3,enum=clouditor.AssessmentResult_ComplianceStatus" json:"result,omitempty"`
	TargetValue string                            `protobuf:"bytes,4,opt,name=target_value,json=targetValue,proto3" json:"target_value,omitempty"`
	Evidence    *assessment.Evidence              `protobuf:"bytes,5,opt,name=evidence,proto3" json:"evidence,omitempty"`
}

func (x *AssessmentResult) Reset() {
	*x = AssessmentResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orchestrator_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssessmentResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssessmentResult) ProtoMessage() {}

func (x *AssessmentResult) ProtoReflect() protoreflect.Message {
	mi := &file_orchestrator_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssessmentResult.ProtoReflect.Descriptor instead.
func (*AssessmentResult) Descriptor() ([]byte, []int) {
	return file_orchestrator_proto_rawDescGZIP(), []int{0}
}

func (x *AssessmentResult) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *AssessmentResult) GetMetricId() string {
	if x != nil {
		return x.MetricId
	}
	return ""
}

func (x *AssessmentResult) GetResult() AssessmentResult_ComplianceStatus {
	if x != nil {
		return x.Result
	}
	return AssessmentResult_COMPLIANT
}

func (x *AssessmentResult) GetTargetValue() string {
	if x != nil {
		return x.TargetValue
	}
	return ""
}

func (x *AssessmentResult) GetEvidence() *assessment.Evidence {
	if x != nil {
		return x.Evidence
	}
	return nil
}

// Represents an external tool or service that offers assessments according to
// certain metrics
type AssessmentTool struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// TODO: we might change this to just an identifier, if we only allow tools to
	// refer to existing IDs
	AvailableMetrics []*assessment.Metric `protobuf:"bytes,3,rep,name=available_metrics,json=availableMetrics,proto3" json:"available_metrics,omitempty"`
}

func (x *AssessmentTool) Reset() {
	*x = AssessmentTool{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orchestrator_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssessmentTool) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssessmentTool) ProtoMessage() {}

func (x *AssessmentTool) ProtoReflect() protoreflect.Message {
	mi := &file_orchestrator_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssessmentTool.ProtoReflect.Descriptor instead.
func (*AssessmentTool) Descriptor() ([]byte, []int) {
	return file_orchestrator_proto_rawDescGZIP(), []int{1}
}

func (x *AssessmentTool) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *AssessmentTool) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *AssessmentTool) GetAvailableMetrics() []*assessment.Metric {
	if x != nil {
		return x.AvailableMetrics
	}
	return nil
}

type RegisterAssessmentToolRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tool *AssessmentTool `protobuf:"bytes,1,opt,name=tool,proto3" json:"tool,omitempty"`
}

func (x *RegisterAssessmentToolRequest) Reset() {
	*x = RegisterAssessmentToolRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orchestrator_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterAssessmentToolRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterAssessmentToolRequest) ProtoMessage() {}

func (x *RegisterAssessmentToolRequest) ProtoReflect() protoreflect.Message {
	mi := &file_orchestrator_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterAssessmentToolRequest.ProtoReflect.Descriptor instead.
func (*RegisterAssessmentToolRequest) Descriptor() ([]byte, []int) {
	return file_orchestrator_proto_rawDescGZIP(), []int{2}
}

func (x *RegisterAssessmentToolRequest) GetTool() *AssessmentTool {
	if x != nil {
		return x.Tool
	}
	return nil
}

type ListAssessmentToolsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// filter tools by metric id
	MetricId string `protobuf:"bytes,1,opt,name=metric_id,json=metricId,proto3" json:"metric_id,omitempty"`
}

func (x *ListAssessmentToolsRequest) Reset() {
	*x = ListAssessmentToolsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orchestrator_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListAssessmentToolsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAssessmentToolsRequest) ProtoMessage() {}

func (x *ListAssessmentToolsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_orchestrator_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAssessmentToolsRequest.ProtoReflect.Descriptor instead.
func (*ListAssessmentToolsRequest) Descriptor() ([]byte, []int) {
	return file_orchestrator_proto_rawDescGZIP(), []int{3}
}

func (x *ListAssessmentToolsRequest) GetMetricId() string {
	if x != nil {
		return x.MetricId
	}
	return ""
}

type ListAssessmentToolsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tools []*AssessmentTool `protobuf:"bytes,1,rep,name=tools,proto3" json:"tools,omitempty"`
}

func (x *ListAssessmentToolsResponse) Reset() {
	*x = ListAssessmentToolsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orchestrator_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListAssessmentToolsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAssessmentToolsResponse) ProtoMessage() {}

func (x *ListAssessmentToolsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_orchestrator_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAssessmentToolsResponse.ProtoReflect.Descriptor instead.
func (*ListAssessmentToolsResponse) Descriptor() ([]byte, []int) {
	return file_orchestrator_proto_rawDescGZIP(), []int{4}
}

func (x *ListAssessmentToolsResponse) GetTools() []*AssessmentTool {
	if x != nil {
		return x.Tools
	}
	return nil
}

type GetAssessmentToolRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetAssessmentToolRequest) Reset() {
	*x = GetAssessmentToolRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orchestrator_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAssessmentToolRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAssessmentToolRequest) ProtoMessage() {}

func (x *GetAssessmentToolRequest) ProtoReflect() protoreflect.Message {
	mi := &file_orchestrator_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAssessmentToolRequest.ProtoReflect.Descriptor instead.
func (*GetAssessmentToolRequest) Descriptor() ([]byte, []int) {
	return file_orchestrator_proto_rawDescGZIP(), []int{5}
}

func (x *GetAssessmentToolRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type UpdateAssessmentToolRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string          `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Tool *AssessmentTool `protobuf:"bytes,2,opt,name=tool,proto3" json:"tool,omitempty"`
}

func (x *UpdateAssessmentToolRequest) Reset() {
	*x = UpdateAssessmentToolRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orchestrator_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateAssessmentToolRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateAssessmentToolRequest) ProtoMessage() {}

func (x *UpdateAssessmentToolRequest) ProtoReflect() protoreflect.Message {
	mi := &file_orchestrator_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateAssessmentToolRequest.ProtoReflect.Descriptor instead.
func (*UpdateAssessmentToolRequest) Descriptor() ([]byte, []int) {
	return file_orchestrator_proto_rawDescGZIP(), []int{6}
}

func (x *UpdateAssessmentToolRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateAssessmentToolRequest) GetTool() *AssessmentTool {
	if x != nil {
		return x.Tool
	}
	return nil
}

type DeregisterAssessmentToolRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeregisterAssessmentToolRequest) Reset() {
	*x = DeregisterAssessmentToolRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orchestrator_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeregisterAssessmentToolRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeregisterAssessmentToolRequest) ProtoMessage() {}

func (x *DeregisterAssessmentToolRequest) ProtoReflect() protoreflect.Message {
	mi := &file_orchestrator_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeregisterAssessmentToolRequest.ProtoReflect.Descriptor instead.
func (*DeregisterAssessmentToolRequest) Descriptor() ([]byte, []int) {
	return file_orchestrator_proto_rawDescGZIP(), []int{7}
}

func (x *DeregisterAssessmentToolRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

var File_orchestrator_proto protoreflect.FileDescriptor

var file_orchestrator_proto_rawDesc = []byte{
	0x0a, 0x12, 0x6f, 0x72, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x1a,
	0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65,
	0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0e, 0x65, 0x76, 0x69, 0x64,
	0x65, 0x6e, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8f, 0x02, 0x0a, 0x10, 0x41,
	0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x1b, 0x0a, 0x09, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x49, 0x64, 0x12, 0x44, 0x0a, 0x06,
	0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2c, 0x2e, 0x63,
	0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d,
	0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x69,
	0x61, 0x6e, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x2f, 0x0a, 0x08, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69,
	0x74, 0x6f, 0x72, 0x2e, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x08, 0x65, 0x76,
	0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x22, 0x34, 0x0a, 0x10, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x69,
	0x61, 0x6e, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0d, 0x0a, 0x09, 0x43, 0x4f,
	0x4d, 0x50, 0x4c, 0x49, 0x41, 0x4e, 0x54, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x4e, 0x4f, 0x4e,
	0x5f, 0x43, 0x4f, 0x4d, 0x50, 0x4c, 0x49, 0x41, 0x4e, 0x54, 0x10, 0x01, 0x22, 0x74, 0x0a, 0x0e,
	0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x6f, 0x6f, 0x6c, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x3e, 0x0a, 0x11, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x5f,
	0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e,
	0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x52, 0x10, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x22, 0x4e, 0x0a, 0x1d, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x41, 0x73,
	0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x6f, 0x6f, 0x6c, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x2d, 0x0a, 0x04, 0x74, 0x6f, 0x6f, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x19, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x41, 0x73,
	0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x6f, 0x6f, 0x6c, 0x52, 0x04, 0x74, 0x6f,
	0x6f, 0x6c, 0x22, 0x39, 0x0a, 0x1a, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73,
	0x6d, 0x65, 0x6e, 0x74, 0x54, 0x6f, 0x6f, 0x6c, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1b, 0x0a, 0x09, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x49, 0x64, 0x22, 0x4e, 0x0a,
	0x1b, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x54,
	0x6f, 0x6f, 0x6c, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2f, 0x0a, 0x05,
	0x74, 0x6f, 0x6f, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x63, 0x6c,
	0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65,
	0x6e, 0x74, 0x54, 0x6f, 0x6f, 0x6c, 0x52, 0x05, 0x74, 0x6f, 0x6f, 0x6c, 0x73, 0x22, 0x2a, 0x0a,
	0x18, 0x47, 0x65, 0x74, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x6f,
	0x6f, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x5c, 0x0a, 0x1b, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x6f, 0x6f,
	0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2d, 0x0a, 0x04, 0x74, 0x6f, 0x6f, 0x6c,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74,
	0x6f, 0x72, 0x2e, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x6f, 0x6f,
	0x6c, 0x52, 0x04, 0x74, 0x6f, 0x6f, 0x6c, 0x22, 0x31, 0x0a, 0x1f, 0x44, 0x65, 0x72, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x54,
	0x6f, 0x6f, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x32, 0x86, 0x06, 0x0a, 0x0c, 0x4f,
	0x72, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x8d, 0x01, 0x0a, 0x16,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65,
	0x6e, 0x74, 0x54, 0x6f, 0x6f, 0x6c, 0x12, 0x28, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74,
	0x6f, 0x72, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x41, 0x73, 0x73, 0x65, 0x73,
	0x73, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x6f, 0x6f, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x19, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x41, 0x73, 0x73,
	0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x6f, 0x6f, 0x6c, 0x22, 0x2e, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x28, 0x22, 0x20, 0x2f, 0x76, 0x31, 0x2f, 0x6f, 0x72, 0x63, 0x68, 0x65, 0x73, 0x74,
	0x61, 0x74, 0x6f, 0x72, 0x2f, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x5f,
	0x74, 0x6f, 0x6f, 0x6c, 0x73, 0x3a, 0x04, 0x74, 0x6f, 0x6f, 0x6c, 0x12, 0xbf, 0x01, 0x0a, 0x13,
	0x4c, 0x69, 0x73, 0x74, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x6f,
	0x6f, 0x6c, 0x73, 0x12, 0x25, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e,
	0x4c, 0x69, 0x73, 0x74, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x6f,
	0x6f, 0x6c, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x63, 0x6c, 0x6f,
	0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x73, 0x73, 0x65, 0x73,
	0x73, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x6f, 0x6f, 0x6c, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x59, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x53, 0x12, 0x21, 0x2f, 0x76, 0x31, 0x2f,
	0x6f, 0x72, 0x63, 0x68, 0x65, 0x73, 0x74, 0x61, 0x74, 0x6f, 0x72, 0x2f, 0x61, 0x73, 0x73, 0x65,
	0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x6f, 0x6f, 0x6c, 0x73, 0x2f, 0x5a, 0x2e, 0x12,
	0x2c, 0x2f, 0x76, 0x31, 0x2f, 0x6f, 0x72, 0x63, 0x68, 0x65, 0x73, 0x74, 0x61, 0x74, 0x6f, 0x72,
	0x2f, 0x7b, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x5f, 0x69, 0x64, 0x7d, 0x2f, 0x61, 0x73, 0x73,
	0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x6f, 0x6f, 0x6c, 0x73, 0x12, 0x82, 0x01,
	0x0a, 0x11, 0x47, 0x65, 0x74, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x54,
	0x6f, 0x6f, 0x6c, 0x12, 0x23, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e,
	0x47, 0x65, 0x74, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x6f, 0x6f,
	0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64,
	0x69, 0x74, 0x6f, 0x72, 0x2e, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x54,
	0x6f, 0x6f, 0x6c, 0x22, 0x2d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x27, 0x12, 0x25, 0x2f, 0x76, 0x31,
	0x2f, 0x6f, 0x72, 0x63, 0x68, 0x65, 0x73, 0x74, 0x61, 0x74, 0x6f, 0x72, 0x2f, 0x61, 0x73, 0x73,
	0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x6f, 0x6f, 0x6c, 0x73, 0x2f, 0x7b, 0x69,
	0x64, 0x7d, 0x12, 0x8e, 0x01, 0x0a, 0x14, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x73, 0x73,
	0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x6f, 0x6f, 0x6c, 0x12, 0x26, 0x2e, 0x63, 0x6c,
	0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x73,
	0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x6f, 0x6f, 0x6c, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e,
	0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x6f, 0x6f, 0x6c, 0x22, 0x33,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x2d, 0x1a, 0x25, 0x2f, 0x76, 0x31, 0x2f, 0x6f, 0x72, 0x63, 0x68,
	0x65, 0x73, 0x74, 0x61, 0x74, 0x6f, 0x72, 0x2f, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65,
	0x6e, 0x74, 0x5f, 0x74, 0x6f, 0x6f, 0x6c, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x3a, 0x04, 0x74,
	0x6f, 0x6f, 0x6c, 0x12, 0x8d, 0x01, 0x0a, 0x18, 0x44, 0x65, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x65, 0x72, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x6f, 0x6f, 0x6c,
	0x12, 0x2a, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x44, 0x65, 0x72,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e,
	0x74, 0x54, 0x6f, 0x6f, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x2d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x27, 0x2a, 0x25, 0x2f, 0x76,
	0x31, 0x2f, 0x6f, 0x72, 0x63, 0x68, 0x65, 0x73, 0x74, 0x61, 0x74, 0x6f, 0x72, 0x2f, 0x61, 0x73,
	0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x6f, 0x6f, 0x6c, 0x73, 0x2f, 0x7b,
	0x69, 0x64, 0x7d, 0x42, 0x12, 0x5a, 0x10, 0x61, 0x70, 0x69, 0x2f, 0x6f, 0x72, 0x63, 0x68, 0x65,
	0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_orchestrator_proto_rawDescOnce sync.Once
	file_orchestrator_proto_rawDescData = file_orchestrator_proto_rawDesc
)

func file_orchestrator_proto_rawDescGZIP() []byte {
	file_orchestrator_proto_rawDescOnce.Do(func() {
		file_orchestrator_proto_rawDescData = protoimpl.X.CompressGZIP(file_orchestrator_proto_rawDescData)
	})
	return file_orchestrator_proto_rawDescData
}

var file_orchestrator_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_orchestrator_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_orchestrator_proto_goTypes = []interface{}{
	(AssessmentResult_ComplianceStatus)(0),  // 0: clouditor.AssessmentResult.ComplianceStatus
	(*AssessmentResult)(nil),                // 1: clouditor.AssessmentResult
	(*AssessmentTool)(nil),                  // 2: clouditor.AssessmentTool
	(*RegisterAssessmentToolRequest)(nil),   // 3: clouditor.RegisterAssessmentToolRequest
	(*ListAssessmentToolsRequest)(nil),      // 4: clouditor.ListAssessmentToolsRequest
	(*ListAssessmentToolsResponse)(nil),     // 5: clouditor.ListAssessmentToolsResponse
	(*GetAssessmentToolRequest)(nil),        // 6: clouditor.GetAssessmentToolRequest
	(*UpdateAssessmentToolRequest)(nil),     // 7: clouditor.UpdateAssessmentToolRequest
	(*DeregisterAssessmentToolRequest)(nil), // 8: clouditor.DeregisterAssessmentToolRequest
	(*assessment.Evidence)(nil),             // 9: clouditor.Evidence
	(*assessment.Metric)(nil),               // 10: clouditor.Metric
	(*emptypb.Empty)(nil),                   // 11: google.protobuf.Empty
}
var file_orchestrator_proto_depIdxs = []int32{
	0,  // 0: clouditor.AssessmentResult.result:type_name -> clouditor.AssessmentResult.ComplianceStatus
	9,  // 1: clouditor.AssessmentResult.evidence:type_name -> clouditor.Evidence
	10, // 2: clouditor.AssessmentTool.available_metrics:type_name -> clouditor.Metric
	2,  // 3: clouditor.RegisterAssessmentToolRequest.tool:type_name -> clouditor.AssessmentTool
	2,  // 4: clouditor.ListAssessmentToolsResponse.tools:type_name -> clouditor.AssessmentTool
	2,  // 5: clouditor.UpdateAssessmentToolRequest.tool:type_name -> clouditor.AssessmentTool
	3,  // 6: clouditor.Orchestrator.RegisterAssessmentTool:input_type -> clouditor.RegisterAssessmentToolRequest
	4,  // 7: clouditor.Orchestrator.ListAssessmentTools:input_type -> clouditor.ListAssessmentToolsRequest
	6,  // 8: clouditor.Orchestrator.GetAssessmentTool:input_type -> clouditor.GetAssessmentToolRequest
	7,  // 9: clouditor.Orchestrator.UpdateAssessmentTool:input_type -> clouditor.UpdateAssessmentToolRequest
	8,  // 10: clouditor.Orchestrator.DeregisterAssessmentTool:input_type -> clouditor.DeregisterAssessmentToolRequest
	2,  // 11: clouditor.Orchestrator.RegisterAssessmentTool:output_type -> clouditor.AssessmentTool
	5,  // 12: clouditor.Orchestrator.ListAssessmentTools:output_type -> clouditor.ListAssessmentToolsResponse
	2,  // 13: clouditor.Orchestrator.GetAssessmentTool:output_type -> clouditor.AssessmentTool
	2,  // 14: clouditor.Orchestrator.UpdateAssessmentTool:output_type -> clouditor.AssessmentTool
	11, // 15: clouditor.Orchestrator.DeregisterAssessmentTool:output_type -> google.protobuf.Empty
	11, // [11:16] is the sub-list for method output_type
	6,  // [6:11] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_orchestrator_proto_init() }
func file_orchestrator_proto_init() {
	if File_orchestrator_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_orchestrator_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssessmentResult); i {
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
		file_orchestrator_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssessmentTool); i {
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
		file_orchestrator_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterAssessmentToolRequest); i {
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
		file_orchestrator_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListAssessmentToolsRequest); i {
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
		file_orchestrator_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListAssessmentToolsResponse); i {
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
		file_orchestrator_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAssessmentToolRequest); i {
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
		file_orchestrator_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateAssessmentToolRequest); i {
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
		file_orchestrator_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeregisterAssessmentToolRequest); i {
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
			RawDescriptor: file_orchestrator_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_orchestrator_proto_goTypes,
		DependencyIndexes: file_orchestrator_proto_depIdxs,
		EnumInfos:         file_orchestrator_proto_enumTypes,
		MessageInfos:      file_orchestrator_proto_msgTypes,
	}.Build()
	File_orchestrator_proto = out.File
	file_orchestrator_proto_rawDesc = nil
	file_orchestrator_proto_goTypes = nil
	file_orchestrator_proto_depIdxs = nil
}
