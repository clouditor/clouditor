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
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.2
// source: api/assessment/assessment.proto

package assessment

import (
	evidence "clouditor.io/clouditor/api/evidence"
	_ "github.com/srikrsna/protoc-gen-gotag/tagger"
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

type AssessEvidenceResponse_AssessmentStatus int32

const (
	AssessEvidenceResponse_ASSESSMENT_STATUS_UNSPECIFIED AssessEvidenceResponse_AssessmentStatus = 0
	AssessEvidenceResponse_WAITING_FOR_RELATED           AssessEvidenceResponse_AssessmentStatus = 1
	AssessEvidenceResponse_ASSESSED                      AssessEvidenceResponse_AssessmentStatus = 2
	AssessEvidenceResponse_FAILED                        AssessEvidenceResponse_AssessmentStatus = 3
)

// Enum value maps for AssessEvidenceResponse_AssessmentStatus.
var (
	AssessEvidenceResponse_AssessmentStatus_name = map[int32]string{
		0: "ASSESSMENT_STATUS_UNSPECIFIED",
		1: "WAITING_FOR_RELATED",
		2: "ASSESSED",
		3: "FAILED",
	}
	AssessEvidenceResponse_AssessmentStatus_value = map[string]int32{
		"ASSESSMENT_STATUS_UNSPECIFIED": 0,
		"WAITING_FOR_RELATED":           1,
		"ASSESSED":                      2,
		"FAILED":                        3,
	}
)

func (x AssessEvidenceResponse_AssessmentStatus) Enum() *AssessEvidenceResponse_AssessmentStatus {
	p := new(AssessEvidenceResponse_AssessmentStatus)
	*p = x
	return p
}

func (x AssessEvidenceResponse_AssessmentStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AssessEvidenceResponse_AssessmentStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_api_assessment_assessment_proto_enumTypes[0].Descriptor()
}

func (AssessEvidenceResponse_AssessmentStatus) Type() protoreflect.EnumType {
	return &file_api_assessment_assessment_proto_enumTypes[0]
}

func (x AssessEvidenceResponse_AssessmentStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AssessEvidenceResponse_AssessmentStatus.Descriptor instead.
func (AssessEvidenceResponse_AssessmentStatus) EnumDescriptor() ([]byte, []int) {
	return file_api_assessment_assessment_proto_rawDescGZIP(), []int{6, 0}
}

type ListAssessmentResultsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PageSize  int32  `protobuf:"varint,1,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	PageToken string `protobuf:"bytes,2,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	OrderBy   string `protobuf:"bytes,3,opt,name=order_by,json=orderBy,proto3" json:"order_by,omitempty"`
	Asc       bool   `protobuf:"varint,4,opt,name=asc,proto3" json:"asc,omitempty"`
}

func (x *ListAssessmentResultsRequest) Reset() {
	*x = ListAssessmentResultsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_assessment_assessment_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListAssessmentResultsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAssessmentResultsRequest) ProtoMessage() {}

func (x *ListAssessmentResultsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_assessment_proto_msgTypes[0]
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
	return file_api_assessment_assessment_proto_rawDescGZIP(), []int{0}
}

func (x *ListAssessmentResultsRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListAssessmentResultsRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

func (x *ListAssessmentResultsRequest) GetOrderBy() string {
	if x != nil {
		return x.OrderBy
	}
	return ""
}

func (x *ListAssessmentResultsRequest) GetAsc() bool {
	if x != nil {
		return x.Asc
	}
	return false
}

type ListAssessmentResultsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Results       []*AssessmentResult `protobuf:"bytes,1,rep,name=results,proto3" json:"results,omitempty"`
	NextPageToken string              `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
}

func (x *ListAssessmentResultsResponse) Reset() {
	*x = ListAssessmentResultsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_assessment_assessment_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListAssessmentResultsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAssessmentResultsResponse) ProtoMessage() {}

func (x *ListAssessmentResultsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_assessment_proto_msgTypes[1]
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
	return file_api_assessment_assessment_proto_rawDescGZIP(), []int{1}
}

func (x *ListAssessmentResultsResponse) GetResults() []*AssessmentResult {
	if x != nil {
		return x.Results
	}
	return nil
}

func (x *ListAssessmentResultsResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

type ConfigureAssessmentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ConfigureAssessmentRequest) Reset() {
	*x = ConfigureAssessmentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_assessment_assessment_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfigureAssessmentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigureAssessmentRequest) ProtoMessage() {}

func (x *ConfigureAssessmentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_assessment_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigureAssessmentRequest.ProtoReflect.Descriptor instead.
func (*ConfigureAssessmentRequest) Descriptor() ([]byte, []int) {
	return file_api_assessment_assessment_proto_rawDescGZIP(), []int{2}
}

type ConfigureAssessmentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ConfigureAssessmentResponse) Reset() {
	*x = ConfigureAssessmentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_assessment_assessment_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfigureAssessmentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigureAssessmentResponse) ProtoMessage() {}

func (x *ConfigureAssessmentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_assessment_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigureAssessmentResponse.ProtoReflect.Descriptor instead.
func (*ConfigureAssessmentResponse) Descriptor() ([]byte, []int) {
	return file_api_assessment_assessment_proto_rawDescGZIP(), []int{3}
}

type TriggerAssessmentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SomeOption string `protobuf:"bytes,1,opt,name=some_option,json=someOption,proto3" json:"some_option,omitempty"`
}

func (x *TriggerAssessmentRequest) Reset() {
	*x = TriggerAssessmentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_assessment_assessment_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TriggerAssessmentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TriggerAssessmentRequest) ProtoMessage() {}

func (x *TriggerAssessmentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_assessment_proto_msgTypes[4]
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
	return file_api_assessment_assessment_proto_rawDescGZIP(), []int{4}
}

func (x *TriggerAssessmentRequest) GetSomeOption() string {
	if x != nil {
		return x.SomeOption
	}
	return ""
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
		mi := &file_api_assessment_assessment_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssessEvidenceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssessEvidenceRequest) ProtoMessage() {}

func (x *AssessEvidenceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_assessment_proto_msgTypes[5]
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
	return file_api_assessment_assessment_proto_rawDescGZIP(), []int{5}
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

	Status        AssessEvidenceResponse_AssessmentStatus `protobuf:"varint,1,opt,name=status,proto3,enum=clouditor.AssessEvidenceResponse_AssessmentStatus" json:"status,omitempty"`
	StatusMessage string                                  `protobuf:"bytes,2,opt,name=status_message,json=statusMessage,proto3" json:"status_message,omitempty"`
}

func (x *AssessEvidenceResponse) Reset() {
	*x = AssessEvidenceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_assessment_assessment_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssessEvidenceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssessEvidenceResponse) ProtoMessage() {}

func (x *AssessEvidenceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_assessment_proto_msgTypes[6]
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
	return file_api_assessment_assessment_proto_rawDescGZIP(), []int{6}
}

func (x *AssessEvidenceResponse) GetStatus() AssessEvidenceResponse_AssessmentStatus {
	if x != nil {
		return x.Status
	}
	return AssessEvidenceResponse_ASSESSMENT_STATUS_UNSPECIFIED
}

func (x *AssessEvidenceResponse) GetStatusMessage() string {
	if x != nil {
		return x.StatusMessage
	}
	return ""
}

// A result resource, representing the result after assessing the cloud resource
// with id resource_id.
type AssessmentResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Assessment result id
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Time of assessment
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty" gorm:"serializer:timestamppb;type:time"`
	// Reference to the metric the assessment was based on
	MetricId string `protobuf:"bytes,3,opt,name=metric_id,json=metricId,proto3" json:"metric_id,omitempty"`
	// Data corresponding to the metric by the given metric id
	MetricConfiguration *MetricConfiguration `protobuf:"bytes,4,opt,name=metric_configuration,json=metricConfiguration,proto3" json:"metric_configuration,omitempty"`
	// Compliant case: true or false
	Compliant bool `protobuf:"varint,5,opt,name=compliant,proto3" json:"compliant,omitempty"`
	// Reference to the assessed evidence
	EvidenceId string `protobuf:"bytes,6,opt,name=evidence_id,json=evidenceId,proto3" json:"evidence_id,omitempty"`
	// Reference to the resource of the assessed evidence
	ResourceId string `protobuf:"bytes,7,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
	// Resource types
	ResourceTypes []string `protobuf:"bytes,8,rep,name=resource_types,json=resourceTypes,proto3" json:"resource_types,omitempty"`
	// Some comments on the reason for non-compliance
	NonComplianceComments string `protobuf:"bytes,9,opt,name=non_compliance_comments,json=nonComplianceComments,proto3" json:"non_compliance_comments,omitempty"`
}

func (x *AssessmentResult) Reset() {
	*x = AssessmentResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_assessment_assessment_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssessmentResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssessmentResult) ProtoMessage() {}

func (x *AssessmentResult) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_assessment_proto_msgTypes[7]
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
	return file_api_assessment_assessment_proto_rawDescGZIP(), []int{7}
}

func (x *AssessmentResult) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *AssessmentResult) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *AssessmentResult) GetMetricId() string {
	if x != nil {
		return x.MetricId
	}
	return ""
}

func (x *AssessmentResult) GetMetricConfiguration() *MetricConfiguration {
	if x != nil {
		return x.MetricConfiguration
	}
	return nil
}

func (x *AssessmentResult) GetCompliant() bool {
	if x != nil {
		return x.Compliant
	}
	return false
}

func (x *AssessmentResult) GetEvidenceId() string {
	if x != nil {
		return x.EvidenceId
	}
	return ""
}

func (x *AssessmentResult) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

func (x *AssessmentResult) GetResourceTypes() []string {
	if x != nil {
		return x.ResourceTypes
	}
	return nil
}

func (x *AssessmentResult) GetNonComplianceComments() string {
	if x != nil {
		return x.NonComplianceComments
	}
	return ""
}

var File_api_assessment_assessment_proto protoreflect.FileDescriptor

var file_api_assessment_assessment_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74,
	0x2f, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x09, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x1a, 0x1c, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x76,
	0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x2f, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x73, 0x73, 0x65, 0x73,
	0x73, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x13, 0x74, 0x61, 0x67, 0x67, 0x65, 0x72, 0x2f, 0x74, 0x61, 0x67, 0x67, 0x65,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x87, 0x01, 0x0a, 0x1c, 0x4c, 0x69, 0x73, 0x74,
	0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65,
	0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67,
	0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x61, 0x67, 0x65, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x62, 0x79,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x42, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x61, 0x73, 0x63, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x03, 0x61, 0x73,
	0x63, 0x22, 0x7e, 0x0a, 0x1d, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d,
	0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x35, 0x0a, 0x07, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e,
	0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x52, 0x07, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x12, 0x26, 0x0a, 0x0f, 0x6e, 0x65, 0x78,
	0x74, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x6e, 0x65, 0x78, 0x74, 0x50, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x22, 0x1c, 0x0a, 0x1a, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x65, 0x41, 0x73,
	0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22,
	0x1d, 0x0a, 0x1b, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x65, 0x41, 0x73, 0x73, 0x65,
	0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x3b,
	0x0a, 0x18, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d,
	0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x6f,
	0x6d, 0x65, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x73, 0x6f, 0x6d, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x48, 0x0a, 0x15, 0x41,
	0x73, 0x73, 0x65, 0x73, 0x73, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x2f, 0x0a, 0x08, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74,
	0x6f, 0x72, 0x2e, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x08, 0x65, 0x76, 0x69,
	0x64, 0x65, 0x6e, 0x63, 0x65, 0x22, 0xf5, 0x01, 0x0a, 0x16, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73,
	0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x4a, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x32, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x41, 0x73, 0x73,
	0x65, 0x73, 0x73, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x2e, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x25, 0x0a, 0x0e,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x22, 0x68, 0x0a, 0x10, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e,
	0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x21, 0x0a, 0x1d, 0x41, 0x53, 0x53, 0x45, 0x53,
	0x53, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x4e, 0x53,
	0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x17, 0x0a, 0x13, 0x57, 0x41,
	0x49, 0x54, 0x49, 0x4e, 0x47, 0x5f, 0x46, 0x4f, 0x52, 0x5f, 0x52, 0x45, 0x4c, 0x41, 0x54, 0x45,
	0x44, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x41, 0x53, 0x53, 0x45, 0x53, 0x53, 0x45, 0x44, 0x10,
	0x02, 0x12, 0x0a, 0x0a, 0x06, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x03, 0x22, 0xb9, 0x03,
	0x0a, 0x10, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x66, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x42, 0x2c, 0x9a, 0x84, 0x9e, 0x03, 0x27, 0x67, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x73, 0x65,
	0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x72, 0x3a, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x70, 0x62, 0x3b, 0x74, 0x79, 0x70, 0x65, 0x3a, 0x74, 0x69, 0x6d, 0x65, 0x22, 0x52,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1b, 0x0a, 0x09, 0x6d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x49, 0x64, 0x12, 0x51, 0x0a, 0x14, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f,
	0x72, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x13, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x6f,
	0x6d, 0x70, 0x6c, 0x69, 0x61, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x63,
	0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x61, 0x6e, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x65, 0x76, 0x69, 0x64,
	0x65, 0x6e, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x65,
	0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x18, 0x08, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x0d, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x73, 0x12, 0x36, 0x0a, 0x17, 0x6e, 0x6f, 0x6e, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x61,
	0x6e, 0x63, 0x65, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x15, 0x6e, 0x6f, 0x6e, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x61, 0x6e, 0x63,
	0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x32, 0xd7, 0x03, 0x0a, 0x0a, 0x41, 0x73,
	0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x52, 0x0a, 0x11, 0x54, 0x72, 0x69, 0x67,
	0x67, 0x65, 0x72, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x23, 0x2e,
	0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65,
	0x72, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x89, 0x01, 0x0a,
	0x0e, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x12,
	0x20, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x41, 0x73, 0x73, 0x65,
	0x73, 0x73, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x21, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x41, 0x73,
	0x73, 0x65, 0x73, 0x73, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x32, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x2c, 0x22, 0x18, 0x2f, 0x76,
	0x31, 0x2f, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x65, 0x76, 0x69,
	0x64, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x3a, 0x08, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65,
	0x62, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x5c, 0x0a, 0x0f, 0x41, 0x73, 0x73, 0x65,
	0x73, 0x73, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x12, 0x20, 0x2e, 0x63, 0x6c,
	0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x45, 0x76,
	0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e,
	0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73,
	0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x12, 0x8a, 0x01, 0x0a, 0x15, 0x4c, 0x69, 0x73, 0x74, 0x41,
	0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73,
	0x12, 0x27, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e, 0x63, 0x6c, 0x6f, 0x75,
	0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73,
	0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x1e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x12, 0x16, 0x2f, 0x76, 0x31,
	0x2f, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x72, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x73, 0x42, 0x27, 0x5a, 0x25, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72,
	0x2e, 0x69, 0x6f, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_assessment_assessment_proto_rawDescOnce sync.Once
	file_api_assessment_assessment_proto_rawDescData = file_api_assessment_assessment_proto_rawDesc
)

func file_api_assessment_assessment_proto_rawDescGZIP() []byte {
	file_api_assessment_assessment_proto_rawDescOnce.Do(func() {
		file_api_assessment_assessment_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_assessment_assessment_proto_rawDescData)
	})
	return file_api_assessment_assessment_proto_rawDescData
}

var file_api_assessment_assessment_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_assessment_assessment_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_api_assessment_assessment_proto_goTypes = []interface{}{
	(AssessEvidenceResponse_AssessmentStatus)(0), // 0: clouditor.AssessEvidenceResponse.AssessmentStatus
	(*ListAssessmentResultsRequest)(nil),         // 1: clouditor.ListAssessmentResultsRequest
	(*ListAssessmentResultsResponse)(nil),        // 2: clouditor.ListAssessmentResultsResponse
	(*ConfigureAssessmentRequest)(nil),           // 3: clouditor.ConfigureAssessmentRequest
	(*ConfigureAssessmentResponse)(nil),          // 4: clouditor.ConfigureAssessmentResponse
	(*TriggerAssessmentRequest)(nil),             // 5: clouditor.TriggerAssessmentRequest
	(*AssessEvidenceRequest)(nil),                // 6: clouditor.AssessEvidenceRequest
	(*AssessEvidenceResponse)(nil),               // 7: clouditor.AssessEvidenceResponse
	(*AssessmentResult)(nil),                     // 8: clouditor.AssessmentResult
	(*evidence.Evidence)(nil),                    // 9: clouditor.Evidence
	(*timestamppb.Timestamp)(nil),                // 10: google.protobuf.Timestamp
	(*MetricConfiguration)(nil),                  // 11: clouditor.MetricConfiguration
	(*emptypb.Empty)(nil),                        // 12: google.protobuf.Empty
}
var file_api_assessment_assessment_proto_depIdxs = []int32{
	8,  // 0: clouditor.ListAssessmentResultsResponse.results:type_name -> clouditor.AssessmentResult
	9,  // 1: clouditor.AssessEvidenceRequest.evidence:type_name -> clouditor.Evidence
	0,  // 2: clouditor.AssessEvidenceResponse.status:type_name -> clouditor.AssessEvidenceResponse.AssessmentStatus
	10, // 3: clouditor.AssessmentResult.timestamp:type_name -> google.protobuf.Timestamp
	11, // 4: clouditor.AssessmentResult.metric_configuration:type_name -> clouditor.MetricConfiguration
	5,  // 5: clouditor.Assessment.TriggerAssessment:input_type -> clouditor.TriggerAssessmentRequest
	6,  // 6: clouditor.Assessment.AssessEvidence:input_type -> clouditor.AssessEvidenceRequest
	6,  // 7: clouditor.Assessment.AssessEvidences:input_type -> clouditor.AssessEvidenceRequest
	1,  // 8: clouditor.Assessment.ListAssessmentResults:input_type -> clouditor.ListAssessmentResultsRequest
	12, // 9: clouditor.Assessment.TriggerAssessment:output_type -> google.protobuf.Empty
	7,  // 10: clouditor.Assessment.AssessEvidence:output_type -> clouditor.AssessEvidenceResponse
	7,  // 11: clouditor.Assessment.AssessEvidences:output_type -> clouditor.AssessEvidenceResponse
	2,  // 12: clouditor.Assessment.ListAssessmentResults:output_type -> clouditor.ListAssessmentResultsResponse
	9,  // [9:13] is the sub-list for method output_type
	5,  // [5:9] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_api_assessment_assessment_proto_init() }
func file_api_assessment_assessment_proto_init() {
	if File_api_assessment_assessment_proto != nil {
		return
	}
	file_api_assessment_metric_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_api_assessment_assessment_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_api_assessment_assessment_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_api_assessment_assessment_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfigureAssessmentRequest); i {
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
		file_api_assessment_assessment_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfigureAssessmentResponse); i {
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
		file_api_assessment_assessment_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
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
		file_api_assessment_assessment_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
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
		file_api_assessment_assessment_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
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
		file_api_assessment_assessment_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
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
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_assessment_assessment_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_assessment_assessment_proto_goTypes,
		DependencyIndexes: file_api_assessment_assessment_proto_depIdxs,
		EnumInfos:         file_api_assessment_assessment_proto_enumTypes,
		MessageInfos:      file_api_assessment_assessment_proto_msgTypes,
	}.Build()
	File_api_assessment_assessment_proto = out.File
	file_api_assessment_assessment_proto_rawDesc = nil
	file_api_assessment_assessment_proto_goTypes = nil
	file_api_assessment_assessment_proto_depIdxs = nil
}
