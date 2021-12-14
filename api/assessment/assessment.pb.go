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
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.1
// source: assessment.proto

package assessment

import (
	evidence "clouditor.io/clouditor/api/evidence"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TriggerAssessmentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SomeOption string `protobuf:"bytes,1,opt,name=someOption,proto3" json:"someOption,omitempty"`
}

func (x *TriggerAssessmentRequest) Reset() {
	*x = TriggerAssessmentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_assessment_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TriggerAssessmentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TriggerAssessmentRequest) ProtoMessage() {}

func (x *TriggerAssessmentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_assessment_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TriggerAssessmentRequest.ProtoReflect.Descriptor instead.
func (*TriggerAssessmentRequest) Descriptor() ([]byte, []int) {
	return file_assessment_proto_rawDescGZIP(), []int{0}
}

func (x *TriggerAssessmentRequest) GetSomeOption() string {
	if x != nil {
		return x.SomeOption
	}
	return ""
}

type ListAssessmentResultsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListAssessmentResultsRequest) Reset() {
	*x = ListAssessmentResultsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_assessment_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListAssessmentResultsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAssessmentResultsRequest) ProtoMessage() {}

func (x *ListAssessmentResultsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_assessment_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAssessmentResultsRequest.ProtoReflect.Descriptor instead.
func (*ListAssessmentResultsRequest) Descriptor() ([]byte, []int) {
	return file_assessment_proto_rawDescGZIP(), []int{1}
}

type ListAssessmentResultsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Results []*Result `protobuf:"bytes,1,rep,name=results,proto3" json:"results,omitempty"`
}

func (x *ListAssessmentResultsResponse) Reset() {
	*x = ListAssessmentResultsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_assessment_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListAssessmentResultsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAssessmentResultsResponse) ProtoMessage() {}

func (x *ListAssessmentResultsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_assessment_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAssessmentResultsResponse.ProtoReflect.Descriptor instead.
func (*ListAssessmentResultsResponse) Descriptor() ([]byte, []int) {
	return file_assessment_proto_rawDescGZIP(), []int{2}
}

func (x *ListAssessmentResultsResponse) GetResults() []*Result {
	if x != nil {
		return x.Results
	}
	return nil
}

type AssessEvidenceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Evidence *evidence.Evidence `protobuf:"bytes,1,opt,name=evidence,proto3" json:"evidence,omitempty"`
}

func (x *AssessEvidenceRequest) Reset() {
	*x = AssessEvidenceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_assessment_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssessEvidenceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssessEvidenceRequest) ProtoMessage() {}

func (x *AssessEvidenceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_assessment_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssessEvidenceRequest.ProtoReflect.Descriptor instead.
func (*AssessEvidenceRequest) Descriptor() ([]byte, []int) {
	return file_assessment_proto_rawDescGZIP(), []int{3}
}

func (x *AssessEvidenceRequest) GetEvidence() *evidence.Evidence {
	if x != nil {
		return x.Evidence
	}
	return nil
}

type AssessEvidenceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status bool `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *AssessEvidenceResponse) Reset() {
	*x = AssessEvidenceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_assessment_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssessEvidenceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssessEvidenceResponse) ProtoMessage() {}

func (x *AssessEvidenceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_assessment_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssessEvidenceResponse.ProtoReflect.Descriptor instead.
func (*AssessEvidenceResponse) Descriptor() ([]byte, []int) {
	return file_assessment_proto_rawDescGZIP(), []int{4}
}

func (x *AssessEvidenceResponse) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

// A result resource, representing the result after assessing the cloud resource
// with id resource_id.
type Result struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Assessment result id
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Time of assessment
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// Reference to the metric the assessment was based on
	MetricId string `protobuf:"bytes,3,opt,name=metric_id,json=metricId,proto3" json:"metric_id,omitempty"`
	// Data corresponding to the metric by the given metric id
	MetricData *MetricConfiguration `protobuf:"bytes,4,opt,name=metric_data,json=metricData,proto3" json:"metric_data,omitempty"`
	// Compliant case: true or false
	Compliant bool `protobuf:"varint,5,opt,name=compliant,proto3" json:"compliant,omitempty"`
	// Reference to the assessed evidence
	EvidenceId string `protobuf:"bytes,6,opt,name=evidence_id,json=evidenceId,proto3" json:"evidence_id,omitempty"`
	// Reference to the resource of the assessed evidence
	ResourceId string `protobuf:"bytes,7,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
	// Some comments on the reason for non-compliance
	NonComplianceComments string `protobuf:"bytes,8,opt,name=non_compliance_comments,json=nonComplianceComments,proto3" json:"non_compliance_comments,omitempty"`
}

func (x *Result) Reset() {
	*x = Result{}
	if protoimpl.UnsafeEnabled {
		mi := &file_assessment_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Result) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Result) ProtoMessage() {}

func (x *Result) ProtoReflect() protoreflect.Message {
	mi := &file_assessment_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Result.ProtoReflect.Descriptor instead.
func (*Result) Descriptor() ([]byte, []int) {
	return file_assessment_proto_rawDescGZIP(), []int{5}
}

func (x *Result) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Result) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *Result) GetMetricId() string {
	if x != nil {
		return x.MetricId
	}
	return ""
}

func (x *Result) GetMetricData() *MetricConfiguration {
	if x != nil {
		return x.MetricData
	}
	return nil
}

func (x *Result) GetCompliant() bool {
	if x != nil {
		return x.Compliant
	}
	return false
}

func (x *Result) GetEvidenceId() string {
	if x != nil {
		return x.EvidenceId
	}
	return ""
}

func (x *Result) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

func (x *Result) GetNonComplianceComments() string {
	if x != nil {
		return x.NonComplianceComments
	}
	return ""
}

var File_assessment_proto protoreflect.FileDescriptor

var file_assessment_proto_rawDesc = []byte{
	0x0a, 0x10, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x09, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x1a, 0x1c, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70,
	0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0e, 0x65, 0x76, 0x69, 0x64, 0x65,
	0x6e, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3a, 0x0a, 0x18, 0x54, 0x72, 0x69, 0x67, 0x67,
	0x65, 0x72, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x6f, 0x6d, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x6f, 0x6d, 0x65, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x22, 0x1e, 0x0a, 0x1c, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x73, 0x73, 0x65, 0x73,
	0x73, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x22, 0x4c, 0x0a, 0x1d, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x73, 0x73, 0x65, 0x73,
	0x73, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a, 0x07, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f,
	0x72, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x07, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x73, 0x22, 0x48, 0x0a, 0x15, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x45, 0x76, 0x69, 0x64, 0x65,
	0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2f, 0x0a, 0x08, 0x65, 0x76,
	0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x63,
	0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63,
	0x65, 0x52, 0x08, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x22, 0x30, 0x0a, 0x16, 0x41,
	0x73, 0x73, 0x65, 0x73, 0x73, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0xc8, 0x02,
	0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x12, 0x1b, 0x0a, 0x09, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x5f, 0x69, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x49, 0x64, 0x12,
	0x3f, 0x0a, 0x0b, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72,
	0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x1c, 0x0a, 0x09, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x61, 0x6e, 0x74, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x09, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x61, 0x6e, 0x74, 0x12, 0x1f,
	0x0a, 0x0b, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x49, 0x64, 0x12,
	0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64,
	0x12, 0x36, 0x0a, 0x17, 0x6e, 0x6f, 0x6e, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x61, 0x6e,
	0x63, 0x65, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x15, 0x6e, 0x6f, 0x6e, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65,
	0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x32, 0xae, 0x03, 0x0a, 0x0a, 0x41, 0x73, 0x73,
	0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x52, 0x0a, 0x11, 0x54, 0x72, 0x69, 0x67, 0x67,
	0x65, 0x72, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x23, 0x2e, 0x63,
	0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72,
	0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x8d, 0x01, 0x0a, 0x15,
	0x4c, 0x69, 0x73, 0x74, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x73, 0x12, 0x27, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f,
	0x72, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28,
	0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41,
	0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x21, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1b,
	0x22, 0x16, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74,
	0x2f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x62, 0x01, 0x2a, 0x12, 0x7a, 0x0a, 0x0e, 0x41,
	0x73, 0x73, 0x65, 0x73, 0x73, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x20, 0x2e,
	0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73,
	0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x21, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x41, 0x73, 0x73, 0x65,
	0x73, 0x73, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x23, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1d, 0x22, 0x18, 0x2f, 0x76, 0x31, 0x2f,
	0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x65, 0x76, 0x69, 0x64, 0x65,
	0x6e, 0x63, 0x65, 0x73, 0x62, 0x01, 0x2a, 0x12, 0x40, 0x0a, 0x0f, 0x41, 0x73, 0x73, 0x65, 0x73,
	0x73, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x12, 0x13, 0x2e, 0x63, 0x6c, 0x6f,
	0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x1a,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x28, 0x01, 0x42, 0x10, 0x5a, 0x0e, 0x61, 0x70, 0x69,
	0x2f, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_assessment_proto_rawDescOnce sync.Once
	file_assessment_proto_rawDescData = file_assessment_proto_rawDesc
)

func file_assessment_proto_rawDescGZIP() []byte {
	file_assessment_proto_rawDescOnce.Do(func() {
		file_assessment_proto_rawDescData = protoimpl.X.CompressGZIP(file_assessment_proto_rawDescData)
	})
	return file_assessment_proto_rawDescData
}

var file_assessment_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_assessment_proto_goTypes = []interface{}{
	(*TriggerAssessmentRequest)(nil),      // 0: clouditor.TriggerAssessmentRequest
	(*ListAssessmentResultsRequest)(nil),  // 1: clouditor.ListAssessmentResultsRequest
	(*ListAssessmentResultsResponse)(nil), // 2: clouditor.ListAssessmentResultsResponse
	(*AssessEvidenceRequest)(nil),         // 3: clouditor.AssessEvidenceRequest
	(*AssessEvidenceResponse)(nil),        // 4: clouditor.AssessEvidenceResponse
	(*Result)(nil),                        // 5: clouditor.Result
	(*evidence.Evidence)(nil),             // 6: clouditor.Evidence
	(*timestamppb.Timestamp)(nil),         // 7: google.protobuf.Timestamp
	(*MetricConfiguration)(nil),           // 8: clouditor.MetricConfiguration
	(*emptypb.Empty)(nil),                 // 9: google.protobuf.Empty
}
var file_assessment_proto_depIdxs = []int32{
	5, // 0: clouditor.ListAssessmentResultsResponse.results:type_name -> clouditor.Result
	6, // 1: clouditor.AssessEvidenceRequest.evidence:type_name -> clouditor.Evidence
	7, // 2: clouditor.Result.timestamp:type_name -> google.protobuf.Timestamp
	8, // 3: clouditor.Result.metric_data:type_name -> clouditor.MetricConfiguration
	0, // 4: clouditor.Assessment.TriggerAssessment:input_type -> clouditor.TriggerAssessmentRequest
	1, // 5: clouditor.Assessment.ListAssessmentResults:input_type -> clouditor.ListAssessmentResultsRequest
	3, // 6: clouditor.Assessment.AssessEvidence:input_type -> clouditor.AssessEvidenceRequest
	6, // 7: clouditor.Assessment.AssessEvidences:input_type -> clouditor.Evidence
	9, // 8: clouditor.Assessment.TriggerAssessment:output_type -> google.protobuf.Empty
	2, // 9: clouditor.Assessment.ListAssessmentResults:output_type -> clouditor.ListAssessmentResultsResponse
	4, // 10: clouditor.Assessment.AssessEvidence:output_type -> clouditor.AssessEvidenceResponse
	9, // 11: clouditor.Assessment.AssessEvidences:output_type -> google.protobuf.Empty
	8, // [8:12] is the sub-list for method output_type
	4, // [4:8] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_assessment_proto_init() }
func file_assessment_proto_init() {
	if File_assessment_proto != nil {
		return
	}
	file_metric_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_assessment_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TriggerAssessmentRequest); i {
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
		file_assessment_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListAssessmentResultsRequest); i {
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
		file_assessment_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListAssessmentResultsResponse); i {
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
		file_assessment_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssessEvidenceRequest); i {
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
		file_assessment_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssessEvidenceResponse); i {
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
		file_assessment_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Result); i {
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
			RawDescriptor: file_assessment_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_assessment_proto_goTypes,
		DependencyIndexes: file_assessment_proto_depIdxs,
		MessageInfos:      file_assessment_proto_msgTypes,
	}.Build()
	File_assessment_proto = out.File
	file_assessment_proto_rawDesc = nil
	file_assessment_proto_goTypes = nil
	file_assessment_proto_depIdxs = nil
}
