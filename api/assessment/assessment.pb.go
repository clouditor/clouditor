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
	unsafe "unsafe"
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
	// ComplianceComment contains a human readable description on the reason for (non)-compliance.
	ComplianceComment string `protobuf:"bytes,9,opt,name=compliance_comment,json=complianceComment,proto3" json:"compliance_comment,omitempty"`
	// ComplianceDetails contains machine-readable details about which comparisons lead to a (non)-compliance.
	ComplianceDetails []*ComparisonResult `protobuf:"bytes,10,rep,name=compliance_details,json=complianceDetails,proto3" json:"compliance_details,omitempty" gorm:"serializer:json"`
	// The target of evaluation which this assessment result belongs to
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

func (x *AssessmentResult) GetComplianceComment() string {
	if x != nil {
		return x.ComplianceComment
	}
	return ""
}

func (x *AssessmentResult) GetComplianceDetails() []*ComparisonResult {
	if x != nil {
		return x.ComplianceDetails
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

// An optional structure containing more details how a comparison inside an assessment result was done and if it was successful.
type ComparisonResult struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Property is the property that was compared
	Property string `protobuf:"bytes,1,opt,name=property,proto3" json:"property,omitempty"`
	// Value is the value in the property
	Value *structpb.Value `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
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

func (x *ComparisonResult) GetValue() *structpb.Value {
	if x != nil {
		return x.Value
	}
	return nil
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

const file_api_assessment_assessment_proto_rawDesc = "" +
	"\n" +
	"\x1fapi/assessment/assessment.proto\x12\x17clouditor.assessment.v1\x1a\x1bapi/assessment/metric.proto\x1a\x1bapi/evidence/evidence.proto\x1a\x1bbuf/validate/validate.proto\x1a\x1cgoogle/api/annotations.proto\x1a\x1fgoogle/api/field_behavior.proto\x1a\x1bgoogle/protobuf/empty.proto\x1a\x1cgoogle/protobuf/struct.proto\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x13tagger/tagger.proto\"\x1c\n" +
	"\x1aConfigureAssessmentRequest\"\x1d\n" +
	"\x1bConfigureAssessmentResponse\";\n" +
	"\x1aCalculateComplianceRequest\x12\x1d\n" +
	"\n" +
	"control_id\x18\x01 \x01(\tR\tcontrolId\"\\\n" +
	"\x15AssessEvidenceRequest\x12C\n" +
	"\bevidence\x18\x01 \x01(\v2\x1f.clouditor.evidence.v1.EvidenceB\x06\xbaH\x03\xc8\x01\x01R\bevidence\"[\n" +
	"\x16AssessEvidenceResponse\x12A\n" +
	"\x06status\x18\x01 \x01(\x0e2).clouditor.assessment.v1.AssessmentStatusR\x06status\"\x83\x01\n" +
	"\x17AssessEvidencesResponse\x12A\n" +
	"\x06status\x18\x01 \x01(\x0e2).clouditor.assessment.v1.AssessmentStatusR\x06status\x12%\n" +
	"\x0estatus_message\x18\x02 \x01(\tR\rstatusMessage\"\xcb\x06\n" +
	"\x10AssessmentResult\x12\x1b\n" +
	"\x02id\x18\x01 \x01(\tB\v\xe0A\x02\xbaH\x05r\x03\xb0\x01\x01R\x02id\x12t\n" +
	"\ttimestamp\x18\x02 \x01(\v2\x1a.google.protobuf.TimestampB:\xe0A\x02\xbaH\x03\xc8\x01\x01\x9a\x84\x9e\x03,gorm:\"serializer:timestamppb;type:timestamp\"R\ttimestamp\x12'\n" +
	"\tmetric_id\x18\x03 \x01(\tB\n" +
	"\xe0A\x02\xbaH\x04r\x02\x10\x01R\bmetricId\x12\x85\x01\n" +
	"\x14metric_configuration\x18\x04 \x01(\v2,.clouditor.assessment.v1.MetricConfigurationB$\xe0A\x02\xbaH\x03\xc8\x01\x01\x9a\x84\x9e\x03\x16gorm:\"serializer:json\"R\x13metricConfiguration\x12\x1c\n" +
	"\tcompliant\x18\x05 \x01(\bR\tcompliant\x12,\n" +
	"\vevidence_id\x18\x06 \x01(\tB\v\xe0A\x02\xbaH\x05r\x03\xb0\x01\x01R\n" +
	"evidenceId\x12+\n" +
	"\vresource_id\x18\a \x01(\tB\n" +
	"\xe0A\x02\xbaH\x04r\x02\x10\x01R\n" +
	"resourceId\x12M\n" +
	"\x0eresource_types\x18\b \x03(\tB&\xe0A\x02\xbaH\x05\x92\x01\x02\b\x01\x9a\x84\x9e\x03\x16gorm:\"serializer:json\"R\rresourceTypes\x129\n" +
	"\x12compliance_comment\x18\t \x01(\tB\n" +
	"\xe0A\x02\xbaH\x04r\x02\x10\x01R\x11complianceComment\x12u\n" +
	"\x12compliance_details\x18\n" +
	" \x03(\v2).clouditor.assessment.v1.ComparisonResultB\x1b\x9a\x84\x9e\x03\x16gorm:\"serializer:json\"R\x11complianceDetails\x12C\n" +
	"\x17certification_target_id\x18\x14 \x01(\tB\v\xe0A\x02\xbaH\x05r\x03\xb0\x01\x01R\x15certificationTargetId\x12(\n" +
	"\atool_id\x18\x15 \x01(\tB\n" +
	"\xe0A\x02\xbaH\x04r\x02\x10\x01H\x00R\x06toolId\x88\x01\x01B\n" +
	"\n" +
	"\b_tool_id\"\xb6\x02\n" +
	"\x10ComparisonResult\x12&\n" +
	"\bproperty\x18\x01 \x01(\tB\n" +
	"\xe0A\x02\xbaH\x04r\x02\x10\x01R\bproperty\x127\n" +
	"\x05value\x18\x02 \x01(\v2\x16.google.protobuf.ValueB\t\xe0A\x02\xbaH\x03\xc8\x01\x01R\x05value\x12A\n" +
	"\boperator\x18\x03 \x01(\tB%\xe0A\x02\xbaH\x1fr\x1d2\x1b^(<|>|<=|>=|==|isIn|allIn)$R\boperator\x12_\n" +
	"\ftarget_value\x18\x04 \x01(\v2\x16.google.protobuf.ValueB$\xe0A\x02\xbaH\x03\xc8\x01\x01\x9a\x84\x9e\x03\x16gorm:\"serializer:json\"R\vtargetValue\x12\x1d\n" +
	"\asuccess\x18\x05 \x01(\bB\x03\xe0A\x02R\asuccess*\x9e\x01\n" +
	"\x10AssessmentStatus\x12!\n" +
	"\x1dASSESSMENT_STATUS_UNSPECIFIED\x10\x00\x12)\n" +
	"%ASSESSMENT_STATUS_WAITING_FOR_RELATED\x10\x01\x12\x1e\n" +
	"\x1aASSESSMENT_STATUS_ASSESSED\x10\x02\x12\x1c\n" +
	"\x18ASSESSMENT_STATUS_FAILED\x10\x032\x8d\x03\n" +
	"\n" +
	"Assessment\x12d\n" +
	"\x13CalculateCompliance\x123.clouditor.assessment.v1.CalculateComplianceRequest\x1a\x16.google.protobuf.Empty\"\x00\x12\x9d\x01\n" +
	"\x0eAssessEvidence\x12..clouditor.assessment.v1.AssessEvidenceRequest\x1a/.clouditor.assessment.v1.AssessEvidenceResponse\"*\x82\xd3\xe4\x93\x02$:\bevidence\"\x18/v1/assessment/evidences\x12y\n" +
	"\x0fAssessEvidences\x12..clouditor.assessment.v1.AssessEvidenceRequest\x1a0.clouditor.assessment.v1.AssessEvidencesResponse\"\x00(\x010\x01B*Z(clouditor.io/clouditor/v2/api/assessmentb\x06proto3"

var (
	file_api_assessment_assessment_proto_rawDescOnce sync.Once
	file_api_assessment_assessment_proto_rawDescData []byte
)

func file_api_assessment_assessment_proto_rawDescGZIP() []byte {
	file_api_assessment_assessment_proto_rawDescOnce.Do(func() {
		file_api_assessment_assessment_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_assessment_assessment_proto_rawDesc), len(file_api_assessment_assessment_proto_rawDesc)))
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
	8,  // 5: clouditor.assessment.v1.AssessmentResult.compliance_details:type_name -> clouditor.assessment.v1.ComparisonResult
	12, // 6: clouditor.assessment.v1.ComparisonResult.value:type_name -> google.protobuf.Value
	12, // 7: clouditor.assessment.v1.ComparisonResult.target_value:type_name -> google.protobuf.Value
	3,  // 8: clouditor.assessment.v1.Assessment.CalculateCompliance:input_type -> clouditor.assessment.v1.CalculateComplianceRequest
	4,  // 9: clouditor.assessment.v1.Assessment.AssessEvidence:input_type -> clouditor.assessment.v1.AssessEvidenceRequest
	4,  // 10: clouditor.assessment.v1.Assessment.AssessEvidences:input_type -> clouditor.assessment.v1.AssessEvidenceRequest
	13, // 11: clouditor.assessment.v1.Assessment.CalculateCompliance:output_type -> google.protobuf.Empty
	5,  // 12: clouditor.assessment.v1.Assessment.AssessEvidence:output_type -> clouditor.assessment.v1.AssessEvidenceResponse
	6,  // 13: clouditor.assessment.v1.Assessment.AssessEvidences:output_type -> clouditor.assessment.v1.AssessEvidencesResponse
	11, // [11:14] is the sub-list for method output_type
	8,  // [8:11] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
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
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_assessment_assessment_proto_rawDesc), len(file_api_assessment_assessment_proto_rawDesc)),
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
	file_api_assessment_assessment_proto_goTypes = nil
	file_api_assessment_assessment_proto_depIdxs = nil
}
