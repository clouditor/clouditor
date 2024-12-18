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
// 	protoc-gen-go v1.36.0
// 	protoc        (unknown)
// source: api/assessment/assessment.proto

package assessment

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	evidence "clouditor.io/clouditor/v2/api/evidence"
	_ "github.com/srikrsna/protoc-gen-gotag/tagger"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	structpb "google.golang.org/protobuf/types/known/structpb"
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

type AssessmentStatus int32

const (
	AssessmentStatus_ASSESSMENT_STATUS_UNSPECIFIED         AssessmentStatus = 0
	AssessmentStatus_ASSESSMENT_STATUS_WAITING_FOR_RELATED AssessmentStatus = 1
	AssessmentStatus_ASSESSMENT_STATUS_ASSESSED            AssessmentStatus = 2
	AssessmentStatus_ASSESSMENT_STATUS_FAILED              AssessmentStatus = 3
)

// Enum value maps for AssessmentStatus.
var (
	AssessmentStatus_name = map[int32]string{
		0: "ASSESSMENT_STATUS_UNSPECIFIED",
		1: "ASSESSMENT_STATUS_WAITING_FOR_RELATED",
		2: "ASSESSMENT_STATUS_ASSESSED",
		3: "ASSESSMENT_STATUS_FAILED",
	}
	AssessmentStatus_value = map[string]int32{
		"ASSESSMENT_STATUS_UNSPECIFIED":         0,
		"ASSESSMENT_STATUS_WAITING_FOR_RELATED": 1,
		"ASSESSMENT_STATUS_ASSESSED":            2,
		"ASSESSMENT_STATUS_FAILED":              3,
	}
)

func (x AssessmentStatus) Enum() *AssessmentStatus {
	p := new(AssessmentStatus)
	*p = x
	return p
}

func (x AssessmentStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AssessmentStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_api_assessment_assessment_proto_enumTypes[0].Descriptor()
}

func (AssessmentStatus) Type() protoreflect.EnumType {
	return &file_api_assessment_assessment_proto_enumTypes[0]
}

func (x AssessmentStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AssessmentStatus.Descriptor instead.
func (AssessmentStatus) EnumDescriptor() ([]byte, []int) {
	return file_api_assessment_assessment_proto_rawDescGZIP(), []int{0}
}

type ConfigureAssessmentRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ConfigureAssessmentRequest) Reset() {
	*x = ConfigureAssessmentRequest{}
	mi := &file_api_assessment_assessment_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ConfigureAssessmentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigureAssessmentRequest) ProtoMessage() {}

func (x *ConfigureAssessmentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_assessment_proto_msgTypes[0]
	if x != nil {
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
	return file_api_assessment_assessment_proto_rawDescGZIP(), []int{0}
}

type ConfigureAssessmentResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ConfigureAssessmentResponse) Reset() {
	*x = ConfigureAssessmentResponse{}
	mi := &file_api_assessment_assessment_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ConfigureAssessmentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigureAssessmentResponse) ProtoMessage() {}

func (x *ConfigureAssessmentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_assessment_proto_msgTypes[1]
	if x != nil {
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
	return file_api_assessment_assessment_proto_rawDescGZIP(), []int{1}
}

type CalculateComplianceRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ControlId     string                 `protobuf:"bytes,1,opt,name=control_id,json=controlId,proto3" json:"control_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CalculateComplianceRequest) Reset() {
	*x = CalculateComplianceRequest{}
	mi := &file_api_assessment_assessment_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CalculateComplianceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CalculateComplianceRequest) ProtoMessage() {}

func (x *CalculateComplianceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_assessment_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CalculateComplianceRequest.ProtoReflect.Descriptor instead.
func (*CalculateComplianceRequest) Descriptor() ([]byte, []int) {
	return file_api_assessment_assessment_proto_rawDescGZIP(), []int{2}
}

func (x *CalculateComplianceRequest) GetControlId() string {
	if x != nil {
		return x.ControlId
	}
	return ""
}

type AssessEvidenceRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Evidence      *evidence.Evidence     `protobuf:"bytes,1,opt,name=evidence,proto3" json:"evidence,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AssessEvidenceRequest) Reset() {
	*x = AssessEvidenceRequest{}
	mi := &file_api_assessment_assessment_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AssessEvidenceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssessEvidenceRequest) ProtoMessage() {}

func (x *AssessEvidenceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_assessment_proto_msgTypes[3]
	if x != nil {
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
	return file_api_assessment_assessment_proto_rawDescGZIP(), []int{3}
}

func (x *AssessEvidenceRequest) GetEvidence() *evidence.Evidence {
	if x != nil {
		return x.Evidence
	}
	return nil
}

// AssessEvidenceResponse belongs to AssessEvidence, which uses a custom unary
// RPC and therefore requires a response message according to the style
// convention. Since no return values are required, this is empty.
type AssessEvidenceResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Status        AssessmentStatus       `protobuf:"varint,1,opt,name=status,proto3,enum=clouditor.assessment.v1.AssessmentStatus" json:"status,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AssessEvidenceResponse) Reset() {
	*x = AssessEvidenceResponse{}
	mi := &file_api_assessment_assessment_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AssessEvidenceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssessEvidenceResponse) ProtoMessage() {}

func (x *AssessEvidenceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_assessment_proto_msgTypes[4]
	if x != nil {
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
	return file_api_assessment_assessment_proto_rawDescGZIP(), []int{4}
}

func (x *AssessEvidenceResponse) GetStatus() AssessmentStatus {
	if x != nil {
		return x.Status
	}
	return AssessmentStatus_ASSESSMENT_STATUS_UNSPECIFIED
}

// AssessEvidencesResponse belongs to AssessEvidences, which uses a custom
// bidirectional streaming RPC and therefore requires a response message
// according to the style convention. The bidirectional streaming needs the
// status and its message in the response for error handling.
type AssessEvidencesResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Status        AssessmentStatus       `protobuf:"varint,1,opt,name=status,proto3,enum=clouditor.assessment.v1.AssessmentStatus" json:"status,omitempty"`
	StatusMessage string                 `protobuf:"bytes,2,opt,name=status_message,json=statusMessage,proto3" json:"status_message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AssessEvidencesResponse) Reset() {
	*x = AssessEvidencesResponse{}
	mi := &file_api_assessment_assessment_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AssessEvidencesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssessEvidencesResponse) ProtoMessage() {}

func (x *AssessEvidencesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_assessment_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssessEvidencesResponse.ProtoReflect.Descriptor instead.
func (*AssessEvidencesResponse) Descriptor() ([]byte, []int) {
	return file_api_assessment_assessment_proto_rawDescGZIP(), []int{5}
}

func (x *AssessEvidencesResponse) GetStatus() AssessmentStatus {
	if x != nil {
		return x.Status
	}
	return AssessmentStatus_ASSESSMENT_STATUS_UNSPECIFIED
}

func (x *AssessEvidencesResponse) GetStatusMessage() string {
	if x != nil {
		return x.StatusMessage
	}
	return ""
}

// A result resource, representing the result after assessing the cloud resource
// with id resource_id.
type AssessmentResult struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Assessment result id
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Time of assessment
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty" gorm:"serializer:timestamppb;type:timestamp"`
	// Reference to the metric the assessment was based on
	MetricId string `protobuf:"bytes,3,opt,name=metric_id,json=metricId,proto3" json:"metric_id,omitempty"`
	// Data corresponding to the metric by the given metric id
	MetricConfiguration *MetricConfiguration `protobuf:"bytes,4,opt,name=metric_configuration,json=metricConfiguration,proto3" json:"metric_configuration,omitempty" gorm:"serializer:json"`
	// Compliant case: true or false
	Compliant bool `protobuf:"varint,5,opt,name=compliant,proto3" json:"compliant,omitempty"`
	// Reference to the assessed evidence
	EvidenceId string `protobuf:"bytes,6,opt,name=evidence_id,json=evidenceId,proto3" json:"evidence_id,omitempty"`
	// Reference to the resource of the assessed evidence
	ResourceId string `protobuf:"bytes,7,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
	// Resource types
	ResourceTypes []string `protobuf:"bytes,8,rep,name=resource_types,json=resourceTypes,proto3" json:"resource_types,omitempty" gorm:"serializer:json"`
	// Some comments on the reason for non-compliance
	NonComplianceComments string              `protobuf:"bytes,9,opt,name=non_compliance_comments,json=nonComplianceComments,proto3" json:"non_compliance_comments,omitempty"`
	NonComplianceDetails  []*ComparisonResult `protobuf:"bytes,10,rep,name=non_compliance_details,json=nonComplianceDetails,proto3" json:"non_compliance_details,omitempty" gorm:"serializer:json"`
	// The certification target which this assessment result belongs to
	CertificationTargetId string `protobuf:"bytes,20,opt,name=certification_target_id,json=certificationTargetId,proto3" json:"certification_target_id,omitempty"`
	// Reference to the tool which provided the assessment result
	ToolId        *string `protobuf:"bytes,21,opt,name=tool_id,json=toolId,proto3,oneof" json:"tool_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AssessmentResult) Reset() {
	*x = AssessmentResult{}
	mi := &file_api_assessment_assessment_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AssessmentResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssessmentResult) ProtoMessage() {}

func (x *AssessmentResult) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_assessment_proto_msgTypes[6]
	if x != nil {
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
	return file_api_assessment_assessment_proto_rawDescGZIP(), []int{6}
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

func (x *AssessmentResult) GetNonComplianceDetails() []*ComparisonResult {
	if x != nil {
		return x.NonComplianceDetails
	}
	return nil
}

func (x *AssessmentResult) GetCertificationTargetId() string {
	if x != nil {
		return x.CertificationTargetId
	}
	return ""
}

func (x *AssessmentResult) GetToolId() string {
	if x != nil && x.ToolId != nil {
		return *x.ToolId
	}
	return ""
}

// An optional structure containing more details how a comparison inside an assessment result was done and if it was succesful.
type ComparisonResult struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Property is the property that was compared
	Property string `protobuf:"bytes,1,opt,name=property,proto3" json:"property,omitempty"`
	// Value is the value in the property
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	// Operator is the operator used in the comparison
	Operator string `protobuf:"bytes,3,opt,name=operator,proto3" json:"operator,omitempty"`
	// TargetValue is the target value used in the comparison
	TargetValue *structpb.Value `protobuf:"bytes,4,opt,name=target_value,json=targetValue,proto3" json:"target_value,omitempty" gorm:"serializer:json"`
	// Success is true, if the comparison was sucessful
	Success       bool `protobuf:"varint,5,opt,name=success,proto3" json:"success,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ComparisonResult) Reset() {
	*x = ComparisonResult{}
	mi := &file_api_assessment_assessment_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ComparisonResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ComparisonResult) ProtoMessage() {}

func (x *ComparisonResult) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_assessment_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ComparisonResult.ProtoReflect.Descriptor instead.
func (*ComparisonResult) Descriptor() ([]byte, []int) {
	return file_api_assessment_assessment_proto_rawDescGZIP(), []int{7}
}

func (x *ComparisonResult) GetProperty() string {
	if x != nil {
		return x.Property
	}
	return ""
}

func (x *ComparisonResult) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *ComparisonResult) GetOperator() string {
	if x != nil {
		return x.Operator
	}
	return ""
}

func (x *ComparisonResult) GetTargetValue() *structpb.Value {
	if x != nil {
		return x.TargetValue
	}
	return nil
}

func (x *ComparisonResult) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

var File_api_assessment_assessment_proto protoreflect.FileDescriptor

var file_api_assessment_assessment_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74,
	0x2f, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x17, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x61, 0x73, 0x73,
	0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x61, 0x70, 0x69, 0x2f,
	0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x76, 0x69,
	0x64, 0x65, 0x6e, 0x63, 0x65, 0x2f, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e,
	0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74,
	0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13, 0x74, 0x61, 0x67,
	0x67, 0x65, 0x72, 0x2f, 0x74, 0x61, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x1c, 0x0a, 0x1a, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x65, 0x41, 0x73, 0x73,
	0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x1d,
	0x0a, 0x1b, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x65, 0x41, 0x73, 0x73, 0x65, 0x73,
	0x73, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x3b, 0x0a,
	0x1a, 0x43, 0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x69,
	0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x49, 0x64, 0x22, 0x5c, 0x0a, 0x15, 0x41, 0x73,
	0x73, 0x65, 0x73, 0x73, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x43, 0x0a, 0x08, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f,
	0x72, 0x2e, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76,
	0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x42, 0x06, 0xba, 0x48, 0x03, 0xc8, 0x01, 0x01, 0x52, 0x08,
	0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x22, 0x5b, 0x0a, 0x16, 0x41, 0x73, 0x73, 0x65,
	0x73, 0x73, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x41, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x29, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x61,
	0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x73, 0x73,
	0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x83, 0x01, 0x0a, 0x17, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73,
	0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x41, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x29, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x61, 0x73,
	0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x73, 0x73, 0x65,
	0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0xb4, 0x06, 0x0a, 0x10,
	0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x12, 0x18, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48,
	0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12, 0x71, 0x0a, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x37, 0xba, 0x48, 0x03, 0xc8, 0x01,
	0x01, 0x9a, 0x84, 0x9e, 0x03, 0x2c, 0x67, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x73, 0x65, 0x72, 0x69,
	0x61, 0x6c, 0x69, 0x7a, 0x65, 0x72, 0x3a, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x70, 0x62, 0x3b, 0x74, 0x79, 0x70, 0x65, 0x3a, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x22, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x24, 0x0a,
	0x09, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x49, 0x64, 0x12, 0x82, 0x01, 0x0a, 0x14, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x5f, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x61,
	0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x42, 0x21, 0xba, 0x48, 0x03, 0xc8, 0x01, 0x01, 0x9a, 0x84, 0x9e, 0x03, 0x16, 0x67, 0x6f, 0x72,
	0x6d, 0x3a, 0x22, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x72, 0x3a, 0x6a, 0x73,
	0x6f, 0x6e, 0x22, 0x52, 0x13, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x6f, 0x6d, 0x70,
	0x6c, 0x69, 0x61, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x63, 0x6f, 0x6d,
	0x70, 0x6c, 0x69, 0x61, 0x6e, 0x74, 0x12, 0x29, 0x0a, 0x0b, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e,
	0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05,
	0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x0a, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x49,
	0x64, 0x12, 0x28, 0x0a, 0x0b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52,
	0x0a, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x4a, 0x0a, 0x0e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x18, 0x08, 0x20,
	0x03, 0x28, 0x09, 0x42, 0x23, 0xba, 0x48, 0x05, 0x92, 0x01, 0x02, 0x08, 0x01, 0x9a, 0x84, 0x9e,
	0x03, 0x16, 0x67, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a,
	0x65, 0x72, 0x3a, 0x6a, 0x73, 0x6f, 0x6e, 0x22, 0x52, 0x0d, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x73, 0x12, 0x36, 0x0a, 0x17, 0x6e, 0x6f, 0x6e, 0x5f, 0x63,
	0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e,
	0x74, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x15, 0x6e, 0x6f, 0x6e, 0x43, 0x6f, 0x6d,
	0x70, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12,
	0x7c, 0x0a, 0x16, 0x6e, 0x6f, 0x6e, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x61, 0x6e, 0x63,
	0x65, 0x5f, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x18, 0x0a, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x29, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x61, 0x73, 0x73, 0x65,
	0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x61, 0x72,
	0x69, 0x73, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x42, 0x1b, 0x9a, 0x84, 0x9e, 0x03,
	0x16, 0x67, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65,
	0x72, 0x3a, 0x6a, 0x73, 0x6f, 0x6e, 0x22, 0x52, 0x14, 0x6e, 0x6f, 0x6e, 0x43, 0x6f, 0x6d, 0x70,
	0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x12, 0x40, 0x0a,
	0x17, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74,
	0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x14, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08,
	0xba, 0x48, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x15, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x49, 0x64, 0x12,
	0x25, 0x0a, 0x07, 0x74, 0x6f, 0x6f, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x15, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x48, 0x00, 0x52, 0x06, 0x74, 0x6f, 0x6f,
	0x6c, 0x49, 0x64, 0x88, 0x01, 0x01, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x74, 0x6f, 0x6f, 0x6c, 0x5f,
	0x69, 0x64, 0x22, 0xf6, 0x01, 0x0a, 0x10, 0x43, 0x6f, 0x6d, 0x70, 0x61, 0x72, 0x69, 0x73, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x70, 0x65,
	0x72, 0x74, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x70, 0x65,
	0x72, 0x74, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x3e, 0x0a, 0x08, 0x6f, 0x70, 0x65,
	0x72, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x22, 0xba, 0x48, 0x1f,
	0x72, 0x1d, 0x32, 0x1b, 0x5e, 0x28, 0x3c, 0x7c, 0x3e, 0x7c, 0x3c, 0x3d, 0x7c, 0x3e, 0x3d, 0x7c,
	0x3d, 0x3d, 0x7c, 0x69, 0x73, 0x49, 0x6e, 0x7c, 0x61, 0x6c, 0x6c, 0x49, 0x6e, 0x29, 0x24, 0x52,
	0x08, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x56, 0x0a, 0x0c, 0x74, 0x61, 0x72,
	0x67, 0x65, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x1b, 0x9a, 0x84, 0x9e, 0x03, 0x16, 0x67, 0x6f,
	0x72, 0x6d, 0x3a, 0x22, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x72, 0x3a, 0x6a,
	0x73, 0x6f, 0x6e, 0x22, 0x52, 0x0b, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x2a, 0x9e, 0x01, 0x0a, 0x10,
	0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x21, 0x0a, 0x1d, 0x41, 0x53, 0x53, 0x45, 0x53, 0x53, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x53,
	0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45,
	0x44, 0x10, 0x00, 0x12, 0x29, 0x0a, 0x25, 0x41, 0x53, 0x53, 0x45, 0x53, 0x53, 0x4d, 0x45, 0x4e,
	0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x57, 0x41, 0x49, 0x54, 0x49, 0x4e, 0x47,
	0x5f, 0x46, 0x4f, 0x52, 0x5f, 0x52, 0x45, 0x4c, 0x41, 0x54, 0x45, 0x44, 0x10, 0x01, 0x12, 0x1e,
	0x0a, 0x1a, 0x41, 0x53, 0x53, 0x45, 0x53, 0x53, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x53, 0x54, 0x41,
	0x54, 0x55, 0x53, 0x5f, 0x41, 0x53, 0x53, 0x45, 0x53, 0x53, 0x45, 0x44, 0x10, 0x02, 0x12, 0x1c,
	0x0a, 0x18, 0x41, 0x53, 0x53, 0x45, 0x53, 0x53, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x53, 0x54, 0x41,
	0x54, 0x55, 0x53, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x03, 0x32, 0x8d, 0x03, 0x0a,
	0x0a, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x64, 0x0a, 0x13, 0x43,
	0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x61, 0x6e,
	0x63, 0x65, 0x12, 0x33, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x61,
	0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x61, 0x6c,
	0x63, 0x75, 0x6c, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22,
	0x00, 0x12, 0x9d, 0x01, 0x0a, 0x0e, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x45, 0x76, 0x69, 0x64,
	0x65, 0x6e, 0x63, 0x65, 0x12, 0x2e, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72,
	0x2e, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x41,
	0x73, 0x73, 0x65, 0x73, 0x73, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x2f, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72,
	0x2e, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x41,
	0x73, 0x73, 0x65, 0x73, 0x73, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x2a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x24, 0x3a, 0x08, 0x65,
	0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x22, 0x18, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x73, 0x73,
	0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65,
	0x73, 0x12, 0x79, 0x0a, 0x0f, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x45, 0x76, 0x69, 0x64, 0x65,
	0x6e, 0x63, 0x65, 0x73, 0x12, 0x2e, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72,
	0x2e, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x41,
	0x73, 0x73, 0x65, 0x73, 0x73, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x30, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72,
	0x2e, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x41,
	0x73, 0x73, 0x65, 0x73, 0x73, 0x45, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x2a, 0x5a, 0x28,
	0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x69, 0x6f, 0x2f, 0x63, 0x6c, 0x6f,
	0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2f, 0x76, 0x32, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x73,
	0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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
var file_api_assessment_assessment_proto_goTypes = []any{
	(AssessmentStatus)(0),               // 0: clouditor.assessment.v1.AssessmentStatus
	(*ConfigureAssessmentRequest)(nil),  // 1: clouditor.assessment.v1.ConfigureAssessmentRequest
	(*ConfigureAssessmentResponse)(nil), // 2: clouditor.assessment.v1.ConfigureAssessmentResponse
	(*CalculateComplianceRequest)(nil),  // 3: clouditor.assessment.v1.CalculateComplianceRequest
	(*AssessEvidenceRequest)(nil),       // 4: clouditor.assessment.v1.AssessEvidenceRequest
	(*AssessEvidenceResponse)(nil),      // 5: clouditor.assessment.v1.AssessEvidenceResponse
	(*AssessEvidencesResponse)(nil),     // 6: clouditor.assessment.v1.AssessEvidencesResponse
	(*AssessmentResult)(nil),            // 7: clouditor.assessment.v1.AssessmentResult
	(*ComparisonResult)(nil),            // 8: clouditor.assessment.v1.ComparisonResult
	(*evidence.Evidence)(nil),           // 9: clouditor.evidence.v1.Evidence
	(*timestamppb.Timestamp)(nil),       // 10: google.protobuf.Timestamp
	(*MetricConfiguration)(nil),         // 11: clouditor.assessment.v1.MetricConfiguration
	(*structpb.Value)(nil),              // 12: google.protobuf.Value
	(*emptypb.Empty)(nil),               // 13: google.protobuf.Empty
}
var file_api_assessment_assessment_proto_depIdxs = []int32{
	9,  // 0: clouditor.assessment.v1.AssessEvidenceRequest.evidence:type_name -> clouditor.evidence.v1.Evidence
	0,  // 1: clouditor.assessment.v1.AssessEvidenceResponse.status:type_name -> clouditor.assessment.v1.AssessmentStatus
	0,  // 2: clouditor.assessment.v1.AssessEvidencesResponse.status:type_name -> clouditor.assessment.v1.AssessmentStatus
	10, // 3: clouditor.assessment.v1.AssessmentResult.timestamp:type_name -> google.protobuf.Timestamp
	11, // 4: clouditor.assessment.v1.AssessmentResult.metric_configuration:type_name -> clouditor.assessment.v1.MetricConfiguration
	8,  // 5: clouditor.assessment.v1.AssessmentResult.non_compliance_details:type_name -> clouditor.assessment.v1.ComparisonResult
	12, // 6: clouditor.assessment.v1.ComparisonResult.target_value:type_name -> google.protobuf.Value
	3,  // 7: clouditor.assessment.v1.Assessment.CalculateCompliance:input_type -> clouditor.assessment.v1.CalculateComplianceRequest
	4,  // 8: clouditor.assessment.v1.Assessment.AssessEvidence:input_type -> clouditor.assessment.v1.AssessEvidenceRequest
	4,  // 9: clouditor.assessment.v1.Assessment.AssessEvidences:input_type -> clouditor.assessment.v1.AssessEvidenceRequest
	13, // 10: clouditor.assessment.v1.Assessment.CalculateCompliance:output_type -> google.protobuf.Empty
	5,  // 11: clouditor.assessment.v1.Assessment.AssessEvidence:output_type -> clouditor.assessment.v1.AssessEvidenceResponse
	6,  // 12: clouditor.assessment.v1.Assessment.AssessEvidences:output_type -> clouditor.assessment.v1.AssessEvidencesResponse
	10, // [10:13] is the sub-list for method output_type
	7,  // [7:10] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_api_assessment_assessment_proto_init() }
func file_api_assessment_assessment_proto_init() {
	if File_api_assessment_assessment_proto != nil {
		return
	}
	file_api_assessment_metric_proto_init()
	file_api_assessment_assessment_proto_msgTypes[6].OneofWrappers = []any{}
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
