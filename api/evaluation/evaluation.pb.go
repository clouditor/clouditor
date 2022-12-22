//
// Copyright 2022 Fraunhofer AISEC
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
// source: api/evaluation/evaluation.proto

package evaluation

import (
	assessment "clouditor.io/clouditor/api/assessment"
	orchestrator "clouditor.io/clouditor/api/orchestrator"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "github.com/srikrsna/protoc-gen-gotag/tagger"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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

type EvaluationResult_EvaluationStatus int32

const (
	EvaluationResult_STATUS_UNSPECIFIED EvaluationResult_EvaluationStatus = 0
	EvaluationResult_COMPLIANT          EvaluationResult_EvaluationStatus = 1
	EvaluationResult_NOT_COMPLIANT      EvaluationResult_EvaluationStatus = 2
	EvaluationResult_PENDING            EvaluationResult_EvaluationStatus = 3
)

// Enum value maps for EvaluationResult_EvaluationStatus.
var (
	EvaluationResult_EvaluationStatus_name = map[int32]string{
		0: "STATUS_UNSPECIFIED",
		1: "COMPLIANT",
		2: "NOT_COMPLIANT",
		3: "PENDING",
	}
	EvaluationResult_EvaluationStatus_value = map[string]int32{
		"STATUS_UNSPECIFIED": 0,
		"COMPLIANT":          1,
		"NOT_COMPLIANT":      2,
		"PENDING":            3,
	}
)

func (x EvaluationResult_EvaluationStatus) Enum() *EvaluationResult_EvaluationStatus {
	p := new(EvaluationResult_EvaluationStatus)
	*p = x
	return p
}

func (x EvaluationResult_EvaluationStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EvaluationResult_EvaluationStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_api_evaluation_evaluation_proto_enumTypes[0].Descriptor()
}

func (EvaluationResult_EvaluationStatus) Type() protoreflect.EnumType {
	return &file_api_evaluation_evaluation_proto_enumTypes[0]
}

func (x EvaluationResult_EvaluationStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EvaluationResult_EvaluationStatus.Descriptor instead.
func (EvaluationResult_EvaluationStatus) EnumDescriptor() ([]byte, []int) {
	return file_api_evaluation_evaluation_proto_rawDescGZIP(), []int{6, 0}
}

type ListEvaluationResultsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FilteredCloudServiceId *string `protobuf:"bytes,1,opt,name=filtered_cloud_service_id,json=filteredCloudServiceId,proto3,oneof" json:"filtered_cloud_service_id,omitempty"`
	PageSize               int32   `protobuf:"varint,10,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	PageToken              string  `protobuf:"bytes,11,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	OrderBy                string  `protobuf:"bytes,12,opt,name=order_by,json=orderBy,proto3" json:"order_by,omitempty"`
	Asc                    bool    `protobuf:"varint,13,opt,name=asc,proto3" json:"asc,omitempty"`
}

func (x *ListEvaluationResultsRequest) Reset() {
	*x = ListEvaluationResultsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_evaluation_evaluation_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListEvaluationResultsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListEvaluationResultsRequest) ProtoMessage() {}

func (x *ListEvaluationResultsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_evaluation_evaluation_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListEvaluationResultsRequest.ProtoReflect.Descriptor instead.
func (*ListEvaluationResultsRequest) Descriptor() ([]byte, []int) {
	return file_api_evaluation_evaluation_proto_rawDescGZIP(), []int{0}
}

func (x *ListEvaluationResultsRequest) GetFilteredCloudServiceId() string {
	if x != nil && x.FilteredCloudServiceId != nil {
		return *x.FilteredCloudServiceId
	}
	return ""
}

func (x *ListEvaluationResultsRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListEvaluationResultsRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

func (x *ListEvaluationResultsRequest) GetOrderBy() string {
	if x != nil {
		return x.OrderBy
	}
	return ""
}

func (x *ListEvaluationResultsRequest) GetAsc() bool {
	if x != nil {
		return x.Asc
	}
	return false
}

type ListEvaluationResultsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Results       []*EvaluationResult `protobuf:"bytes,1,rep,name=results,proto3" json:"results,omitempty"`
	NextPageToken string              `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
}

func (x *ListEvaluationResultsResponse) Reset() {
	*x = ListEvaluationResultsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_evaluation_evaluation_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListEvaluationResultsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListEvaluationResultsResponse) ProtoMessage() {}

func (x *ListEvaluationResultsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_evaluation_evaluation_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListEvaluationResultsResponse.ProtoReflect.Descriptor instead.
func (*ListEvaluationResultsResponse) Descriptor() ([]byte, []int) {
	return file_api_evaluation_evaluation_proto_rawDescGZIP(), []int{1}
}

func (x *ListEvaluationResultsResponse) GetResults() []*EvaluationResult {
	if x != nil {
		return x.Results
	}
	return nil
}

func (x *ListEvaluationResultsResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

type StartEvaluationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TargetOfEvaluation *orchestrator.TargetOfEvaluation `protobuf:"bytes,1,opt,name=target_of_evaluation,json=targetOfEvaluation,proto3" json:"target_of_evaluation,omitempty"`
}

func (x *StartEvaluationRequest) Reset() {
	*x = StartEvaluationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_evaluation_evaluation_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartEvaluationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartEvaluationRequest) ProtoMessage() {}

func (x *StartEvaluationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_evaluation_evaluation_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartEvaluationRequest.ProtoReflect.Descriptor instead.
func (*StartEvaluationRequest) Descriptor() ([]byte, []int) {
	return file_api_evaluation_evaluation_proto_rawDescGZIP(), []int{2}
}

func (x *StartEvaluationRequest) GetTargetOfEvaluation() *orchestrator.TargetOfEvaluation {
	if x != nil {
		return x.TargetOfEvaluation
	}
	return nil
}

type StartEvaluationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status        bool   `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	StatusMessage string `protobuf:"bytes,2,opt,name=status_message,json=statusMessage,proto3" json:"status_message,omitempty"`
}

func (x *StartEvaluationResponse) Reset() {
	*x = StartEvaluationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_evaluation_evaluation_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartEvaluationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartEvaluationResponse) ProtoMessage() {}

func (x *StartEvaluationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_evaluation_evaluation_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartEvaluationResponse.ProtoReflect.Descriptor instead.
func (*StartEvaluationResponse) Descriptor() ([]byte, []int) {
	return file_api_evaluation_evaluation_proto_rawDescGZIP(), []int{3}
}

func (x *StartEvaluationResponse) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

func (x *StartEvaluationResponse) GetStatusMessage() string {
	if x != nil {
		return x.StatusMessage
	}
	return ""
}

type StopEvaluationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TargetOfEvaluation *orchestrator.TargetOfEvaluation `protobuf:"bytes,1,opt,name=target_of_evaluation,json=targetOfEvaluation,proto3" json:"target_of_evaluation,omitempty"`
	// The category the control ID belongs to
	CategoryName string `protobuf:"bytes,2,opt,name=category_name,json=categoryName,proto3" json:"category_name,omitempty"`
	// The control ID
	ControlId string `protobuf:"bytes,3,opt,name=control_id,json=controlId,proto3" json:"control_id,omitempty"`
}

func (x *StopEvaluationRequest) Reset() {
	*x = StopEvaluationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_evaluation_evaluation_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StopEvaluationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopEvaluationRequest) ProtoMessage() {}

func (x *StopEvaluationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_evaluation_evaluation_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopEvaluationRequest.ProtoReflect.Descriptor instead.
func (*StopEvaluationRequest) Descriptor() ([]byte, []int) {
	return file_api_evaluation_evaluation_proto_rawDescGZIP(), []int{4}
}

func (x *StopEvaluationRequest) GetTargetOfEvaluation() *orchestrator.TargetOfEvaluation {
	if x != nil {
		return x.TargetOfEvaluation
	}
	return nil
}

func (x *StopEvaluationRequest) GetCategoryName() string {
	if x != nil {
		return x.CategoryName
	}
	return ""
}

func (x *StopEvaluationRequest) GetControlId() string {
	if x != nil {
		return x.ControlId
	}
	return ""
}

type StopEvaluationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StopEvaluationResponse) Reset() {
	*x = StopEvaluationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_evaluation_evaluation_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StopEvaluationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopEvaluationResponse) ProtoMessage() {}

func (x *StopEvaluationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_evaluation_evaluation_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopEvaluationResponse.ProtoReflect.Descriptor instead.
func (*StopEvaluationResponse) Descriptor() ([]byte, []int) {
	return file_api_evaluation_evaluation_proto_rawDescGZIP(), []int{5}
}

// A evaluation result resource, representing the result after evaluating the cloud service with a specific control
type EvaluationResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Evaluation result id
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Evaluation status
	Status EvaluationResult_EvaluationStatus `protobuf:"varint,2,opt,name=status,proto3,enum=clouditor.evaluation.v1.EvaluationResult_EvaluationStatus" json:"status,omitempty"`
	// The category the control belongs to
	CategoryName string `protobuf:"bytes,3,opt,name=category_name,json=categoryName,proto3" json:"category_name,omitempty"`
	// Reference to the control id the evaluation was based on
	ControlId string `protobuf:"bytes,4,opt,name=control_id,json=controlId,proto3" json:"control_id,omitempty"`
	// Time of evaluation
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=timestamp,proto3" json:"timestamp,omitempty" gorm:"serializer:timestamppb;type:time"`
	// Reference to the target of evaluation the evaluation was based on
	TargetOfEvaluation *orchestrator.TargetOfEvaluation `protobuf:"bytes,6,opt,name=target_of_evaluation,json=targetOfEvaluation,proto3" json:"target_of_evaluation,omitempty"`
	// List of assessment results because of which the evaluation status is not 'compliant'
	FailingAssessmentResults []*assessment.AssessmentResult `protobuf:"bytes,7,rep,name=failing_assessment_results,json=failingAssessmentResults,proto3" json:"failing_assessment_results,omitempty"`
}

func (x *EvaluationResult) Reset() {
	*x = EvaluationResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_evaluation_evaluation_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EvaluationResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EvaluationResult) ProtoMessage() {}

func (x *EvaluationResult) ProtoReflect() protoreflect.Message {
	mi := &file_api_evaluation_evaluation_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EvaluationResult.ProtoReflect.Descriptor instead.
func (*EvaluationResult) Descriptor() ([]byte, []int) {
	return file_api_evaluation_evaluation_proto_rawDescGZIP(), []int{6}
}

func (x *EvaluationResult) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *EvaluationResult) GetStatus() EvaluationResult_EvaluationStatus {
	if x != nil {
		return x.Status
	}
	return EvaluationResult_STATUS_UNSPECIFIED
}

func (x *EvaluationResult) GetCategoryName() string {
	if x != nil {
		return x.CategoryName
	}
	return ""
}

func (x *EvaluationResult) GetControlId() string {
	if x != nil {
		return x.ControlId
	}
	return ""
}

func (x *EvaluationResult) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *EvaluationResult) GetTargetOfEvaluation() *orchestrator.TargetOfEvaluation {
	if x != nil {
		return x.TargetOfEvaluation
	}
	return nil
}

func (x *EvaluationResult) GetFailingAssessmentResults() []*assessment.AssessmentResult {
	if x != nil {
		return x.FailingAssessmentResults
	}
	return nil
}

var File_api_evaluation_evaluation_proto protoreflect.FileDescriptor

var file_api_evaluation_evaluation_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2f, 0x65, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x17, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x65, 0x76, 0x61,
	0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x61, 0x70, 0x69, 0x2f, 0x61,
	0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73,
	0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x23, 0x61, 0x70, 0x69, 0x2f,
	0x6f, 0x72, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2f, 0x6f, 0x72, 0x63,
	0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x13, 0x74, 0x61, 0x67, 0x67, 0x65, 0x72, 0x2f, 0x74, 0x61, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xee, 0x01,
	0x0a, 0x1c, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x47,
	0x0a, 0x19, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x65, 0x64, 0x5f, 0x63, 0x6c, 0x6f, 0x75, 0x64,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x48, 0x00, 0x52, 0x16, 0x66, 0x69,
	0x6c, 0x74, 0x65, 0x72, 0x65, 0x64, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f,
	0x73, 0x69, 0x7a, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65,
	0x53, 0x69, 0x7a, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x61, 0x67, 0x65, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x62, 0x79, 0x18,
	0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x42, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x61, 0x73, 0x63, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x08, 0x52, 0x03, 0x61, 0x73, 0x63,
	0x42, 0x1c, 0x0a, 0x1a, 0x5f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x65, 0x64, 0x5f, 0x63, 0x6c,
	0x6f, 0x75, 0x64, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x22, 0x8c,
	0x01, 0x0a, 0x1d, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x43, 0x0a, 0x07, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x29, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x65, 0x76,
	0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x61, 0x6c,
	0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x07, 0x72, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x73, 0x12, 0x26, 0x0a, 0x0f, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x61,
	0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x6e, 0x65, 0x78, 0x74, 0x50, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x83, 0x01,
	0x0a, 0x16, 0x53, 0x74, 0x61, 0x72, 0x74, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x69, 0x0a, 0x14, 0x74, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x5f, 0x6f, 0x66, 0x5f, 0x65, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74,
	0x6f, 0x72, 0x2e, 0x6f, 0x72, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e,
	0x76, 0x31, 0x2e, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x4f, 0x66, 0x45, 0x76, 0x61, 0x6c, 0x75,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52,
	0x12, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x4f, 0x66, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x22, 0x58, 0x0a, 0x17, 0x53, 0x74, 0x61, 0x72, 0x74, 0x45, 0x76, 0x61, 0x6c,
	0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0xd8, 0x01,
	0x0a, 0x15, 0x53, 0x74, 0x6f, 0x70, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x69, 0x0a, 0x14, 0x74, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x5f, 0x6f, 0x66, 0x5f, 0x65, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f,
	0x72, 0x2e, 0x6f, 0x72, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x76,
	0x31, 0x2e, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x4f, 0x66, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x12,
	0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x4f, 0x66, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x2c, 0x0a, 0x0d, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02,
	0x10, 0x01, 0x52, 0x0c, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x26, 0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x09, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x49, 0x64, 0x22, 0x18, 0x0a, 0x16, 0x53, 0x74, 0x6f, 0x70,
	0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x89, 0x05, 0x0a, 0x10, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x1b, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x0b, 0xfa, 0x42, 0x08, 0x72, 0x06, 0xd0, 0x01, 0x01, 0xb0, 0x01, 0x01,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x5c, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x3a, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72,
	0x2e, 0x65, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x45,
	0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x2e,
	0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x42, 0x08, 0xfa, 0x42, 0x05, 0x82, 0x01, 0x02, 0x10, 0x01, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x2c, 0x0a, 0x0d, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02,
	0x10, 0x01, 0x52, 0x0c, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x26, 0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x09, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x49, 0x64, 0x12, 0x66, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x2c, 0x9a, 0x84, 0x9e, 0x03, 0x27, 0x67, 0x6f,
	0x72, 0x6d, 0x3a, 0x22, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x72, 0x3a, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x70, 0x62, 0x3b, 0x74, 0x79, 0x70, 0x65, 0x3a,
	0x74, 0x69, 0x6d, 0x65, 0x22, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x12, 0x69, 0x0a, 0x14, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x6f, 0x66, 0x5f, 0x65, 0x76,
	0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2d,
	0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x6f, 0x72, 0x63, 0x68, 0x65,
	0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x4f, 0x66, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x08, 0xfa,
	0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x12, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x4f,
	0x66, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x76, 0x0a, 0x1a, 0x66,
	0x61, 0x69, 0x6c, 0x69, 0x6e, 0x67, 0x5f, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e,
	0x74, 0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x29, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x61, 0x73, 0x73, 0x65,
	0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73,
	0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x42, 0x0d, 0xfa, 0x42, 0x0a, 0x92,
	0x01, 0x07, 0x22, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x18, 0x66, 0x61, 0x69, 0x6c, 0x69,
	0x6e, 0x67, 0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x73, 0x22, 0x59, 0x0a, 0x10, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x0a, 0x12, 0x53, 0x54, 0x41, 0x54, 0x55,
	0x53, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12,
	0x0d, 0x0a, 0x09, 0x43, 0x4f, 0x4d, 0x50, 0x4c, 0x49, 0x41, 0x4e, 0x54, 0x10, 0x01, 0x12, 0x11,
	0x0a, 0x0d, 0x4e, 0x4f, 0x54, 0x5f, 0x43, 0x4f, 0x4d, 0x50, 0x4c, 0x49, 0x41, 0x4e, 0x54, 0x10,
	0x02, 0x12, 0x0b, 0x0a, 0x07, 0x50, 0x45, 0x4e, 0x44, 0x49, 0x4e, 0x47, 0x10, 0x03, 0x32, 0xa3,
	0x04, 0x0a, 0x0a, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x9e, 0x01,
	0x0a, 0x0f, 0x53, 0x74, 0x61, 0x72, 0x74, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x2f, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x65, 0x76,
	0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x72,
	0x74, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x30, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x65,
	0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61,
	0x72, 0x74, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x28, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x22, 0x3a, 0x01, 0x2a, 0x22,
	0x1d, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f,
	0x65, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x2f, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12, 0xca,
	0x01, 0x0a, 0x0e, 0x53, 0x74, 0x6f, 0x70, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x2e, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x65, 0x76,
	0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x6f, 0x70,
	0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x2f, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x65, 0x76,
	0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x6f, 0x70,
	0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x57, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x51, 0x3a, 0x14, 0x74, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x5f, 0x6f, 0x66, 0x5f, 0x65, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x22, 0x39, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2f, 0x65, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x2f, 0x7b, 0x63, 0x61, 0x74, 0x65, 0x67,
	0x6f, 0x72, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x7d, 0x2f, 0x7b, 0x63, 0x6f, 0x6e, 0x74, 0x72,
	0x6f, 0x6c, 0x5f, 0x69, 0x64, 0x7d, 0x2f, 0x73, 0x74, 0x6f, 0x70, 0x12, 0xa6, 0x01, 0x0a, 0x15,
	0x4c, 0x69, 0x73, 0x74, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x73, 0x12, 0x35, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f,
	0x72, 0x2e, 0x65, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e,
	0x4c, 0x69, 0x73, 0x74, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x36, 0x2e, 0x63,
	0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x65, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x76, 0x61, 0x6c, 0x75,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x12, 0x16, 0x2f, 0x76,
	0x31, 0x2f, 0x65, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x72, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x73, 0x42, 0x27, 0x5a, 0x25, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f,
	0x72, 0x2e, 0x69, 0x6f, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x65, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_evaluation_evaluation_proto_rawDescOnce sync.Once
	file_api_evaluation_evaluation_proto_rawDescData = file_api_evaluation_evaluation_proto_rawDesc
)

func file_api_evaluation_evaluation_proto_rawDescGZIP() []byte {
	file_api_evaluation_evaluation_proto_rawDescOnce.Do(func() {
		file_api_evaluation_evaluation_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_evaluation_evaluation_proto_rawDescData)
	})
	return file_api_evaluation_evaluation_proto_rawDescData
}

var file_api_evaluation_evaluation_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_evaluation_evaluation_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_api_evaluation_evaluation_proto_goTypes = []interface{}{
	(EvaluationResult_EvaluationStatus)(0),  // 0: clouditor.evaluation.v1.EvaluationResult.EvaluationStatus
	(*ListEvaluationResultsRequest)(nil),    // 1: clouditor.evaluation.v1.ListEvaluationResultsRequest
	(*ListEvaluationResultsResponse)(nil),   // 2: clouditor.evaluation.v1.ListEvaluationResultsResponse
	(*StartEvaluationRequest)(nil),          // 3: clouditor.evaluation.v1.StartEvaluationRequest
	(*StartEvaluationResponse)(nil),         // 4: clouditor.evaluation.v1.StartEvaluationResponse
	(*StopEvaluationRequest)(nil),           // 5: clouditor.evaluation.v1.StopEvaluationRequest
	(*StopEvaluationResponse)(nil),          // 6: clouditor.evaluation.v1.StopEvaluationResponse
	(*EvaluationResult)(nil),                // 7: clouditor.evaluation.v1.EvaluationResult
	(*orchestrator.TargetOfEvaluation)(nil), // 8: clouditor.orchestrator.v1.TargetOfEvaluation
	(*timestamppb.Timestamp)(nil),           // 9: google.protobuf.Timestamp
	(*assessment.AssessmentResult)(nil),     // 10: clouditor.assessment.v1.AssessmentResult
}
var file_api_evaluation_evaluation_proto_depIdxs = []int32{
	7,  // 0: clouditor.evaluation.v1.ListEvaluationResultsResponse.results:type_name -> clouditor.evaluation.v1.EvaluationResult
	8,  // 1: clouditor.evaluation.v1.StartEvaluationRequest.target_of_evaluation:type_name -> clouditor.orchestrator.v1.TargetOfEvaluation
	8,  // 2: clouditor.evaluation.v1.StopEvaluationRequest.target_of_evaluation:type_name -> clouditor.orchestrator.v1.TargetOfEvaluation
	0,  // 3: clouditor.evaluation.v1.EvaluationResult.status:type_name -> clouditor.evaluation.v1.EvaluationResult.EvaluationStatus
	9,  // 4: clouditor.evaluation.v1.EvaluationResult.timestamp:type_name -> google.protobuf.Timestamp
	8,  // 5: clouditor.evaluation.v1.EvaluationResult.target_of_evaluation:type_name -> clouditor.orchestrator.v1.TargetOfEvaluation
	10, // 6: clouditor.evaluation.v1.EvaluationResult.failing_assessment_results:type_name -> clouditor.assessment.v1.AssessmentResult
	3,  // 7: clouditor.evaluation.v1.Evaluation.StartEvaluation:input_type -> clouditor.evaluation.v1.StartEvaluationRequest
	5,  // 8: clouditor.evaluation.v1.Evaluation.StopEvaluation:input_type -> clouditor.evaluation.v1.StopEvaluationRequest
	1,  // 9: clouditor.evaluation.v1.Evaluation.ListEvaluationResults:input_type -> clouditor.evaluation.v1.ListEvaluationResultsRequest
	4,  // 10: clouditor.evaluation.v1.Evaluation.StartEvaluation:output_type -> clouditor.evaluation.v1.StartEvaluationResponse
	6,  // 11: clouditor.evaluation.v1.Evaluation.StopEvaluation:output_type -> clouditor.evaluation.v1.StopEvaluationResponse
	2,  // 12: clouditor.evaluation.v1.Evaluation.ListEvaluationResults:output_type -> clouditor.evaluation.v1.ListEvaluationResultsResponse
	10, // [10:13] is the sub-list for method output_type
	7,  // [7:10] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_api_evaluation_evaluation_proto_init() }
func file_api_evaluation_evaluation_proto_init() {
	if File_api_evaluation_evaluation_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_evaluation_evaluation_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListEvaluationResultsRequest); i {
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
		file_api_evaluation_evaluation_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListEvaluationResultsResponse); i {
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
		file_api_evaluation_evaluation_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartEvaluationRequest); i {
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
		file_api_evaluation_evaluation_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartEvaluationResponse); i {
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
		file_api_evaluation_evaluation_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StopEvaluationRequest); i {
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
		file_api_evaluation_evaluation_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StopEvaluationResponse); i {
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
		file_api_evaluation_evaluation_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EvaluationResult); i {
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
	file_api_evaluation_evaluation_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_evaluation_evaluation_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_evaluation_evaluation_proto_goTypes,
		DependencyIndexes: file_api_evaluation_evaluation_proto_depIdxs,
		EnumInfos:         file_api_evaluation_evaluation_proto_enumTypes,
		MessageInfos:      file_api_evaluation_evaluation_proto_msgTypes,
	}.Build()
	File_api_evaluation_evaluation_proto = out.File
	file_api_evaluation_evaluation_proto_rawDesc = nil
	file_api_evaluation_evaluation_proto_goTypes = nil
	file_api_evaluation_evaluation_proto_depIdxs = nil
}
