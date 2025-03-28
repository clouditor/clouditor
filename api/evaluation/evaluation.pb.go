// Copyright 2023 Fraunhofer AISEC
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
// source: api/evaluation/evaluation.proto

package evaluation

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	_ "github.com/srikrsna/protoc-gen-gotag/tagger"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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

type EvaluationStatus int32

const (
	EvaluationStatus_EVALUATION_STATUS_UNSPECIFIED            EvaluationStatus = 0
	EvaluationStatus_EVALUATION_STATUS_COMPLIANT              EvaluationStatus = 1
	EvaluationStatus_EVALUATION_STATUS_COMPLIANT_MANUALLY     EvaluationStatus = 2
	EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT          EvaluationStatus = 3
	EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT_MANUALLY EvaluationStatus = 4
	EvaluationStatus_EVALUATION_STATUS_PENDING                EvaluationStatus = 10
)

// Enum value maps for EvaluationStatus.
var (
	EvaluationStatus_name = map[int32]string{
		0:  "EVALUATION_STATUS_UNSPECIFIED",
		1:  "EVALUATION_STATUS_COMPLIANT",
		2:  "EVALUATION_STATUS_COMPLIANT_MANUALLY",
		3:  "EVALUATION_STATUS_NOT_COMPLIANT",
		4:  "EVALUATION_STATUS_NOT_COMPLIANT_MANUALLY",
		10: "EVALUATION_STATUS_PENDING",
	}
	EvaluationStatus_value = map[string]int32{
		"EVALUATION_STATUS_UNSPECIFIED":            0,
		"EVALUATION_STATUS_COMPLIANT":              1,
		"EVALUATION_STATUS_COMPLIANT_MANUALLY":     2,
		"EVALUATION_STATUS_NOT_COMPLIANT":          3,
		"EVALUATION_STATUS_NOT_COMPLIANT_MANUALLY": 4,
		"EVALUATION_STATUS_PENDING":                10,
	}
)

func (x EvaluationStatus) Enum() *EvaluationStatus {
	p := new(EvaluationStatus)
	*p = x
	return p
}

func (x EvaluationStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EvaluationStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_api_evaluation_evaluation_proto_enumTypes[0].Descriptor()
}

func (EvaluationStatus) Type() protoreflect.EnumType {
	return &file_api_evaluation_evaluation_proto_enumTypes[0]
}

func (x EvaluationStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EvaluationStatus.Descriptor instead.
func (EvaluationStatus) EnumDescriptor() ([]byte, []int) {
	return file_api_evaluation_evaluation_proto_rawDescGZIP(), []int{0}
}

type StartEvaluationRequest struct {
	state        protoimpl.MessageState `protogen:"open.v1"`
	AuditScopeId string                 `protobuf:"bytes,1,opt,name=audit_scope_id,json=auditScopeId,proto3" json:"audit_scope_id,omitempty"`
	// The interval time in minutes the evaluation executes periodically. The
	// default interval is set to 5 minutes.
	Interval      *int32 `protobuf:"varint,3,opt,name=interval,proto3,oneof" json:"interval,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StartEvaluationRequest) Reset() {
	*x = StartEvaluationRequest{}
	mi := &file_api_evaluation_evaluation_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StartEvaluationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartEvaluationRequest) ProtoMessage() {}

func (x *StartEvaluationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_evaluation_evaluation_proto_msgTypes[0]
	if x != nil {
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
	return file_api_evaluation_evaluation_proto_rawDescGZIP(), []int{0}
}

func (x *StartEvaluationRequest) GetAuditScopeId() string {
	if x != nil {
		return x.AuditScopeId
	}
	return ""
}

func (x *StartEvaluationRequest) GetInterval() int32 {
	if x != nil && x.Interval != nil {
		return *x.Interval
	}
	return 0
}

type StartEvaluationResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Successful    bool                   `protobuf:"varint,1,opt,name=successful,proto3" json:"successful,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StartEvaluationResponse) Reset() {
	*x = StartEvaluationResponse{}
	mi := &file_api_evaluation_evaluation_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StartEvaluationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartEvaluationResponse) ProtoMessage() {}

func (x *StartEvaluationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_evaluation_evaluation_proto_msgTypes[1]
	if x != nil {
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
	return file_api_evaluation_evaluation_proto_rawDescGZIP(), []int{1}
}

func (x *StartEvaluationResponse) GetSuccessful() bool {
	if x != nil {
		return x.Successful
	}
	return false
}

type CreateEvaluationResultRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Result        *EvaluationResult      `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateEvaluationResultRequest) Reset() {
	*x = CreateEvaluationResultRequest{}
	mi := &file_api_evaluation_evaluation_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateEvaluationResultRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateEvaluationResultRequest) ProtoMessage() {}

func (x *CreateEvaluationResultRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_evaluation_evaluation_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateEvaluationResultRequest.ProtoReflect.Descriptor instead.
func (*CreateEvaluationResultRequest) Descriptor() ([]byte, []int) {
	return file_api_evaluation_evaluation_proto_rawDescGZIP(), []int{2}
}

func (x *CreateEvaluationResultRequest) GetResult() *EvaluationResult {
	if x != nil {
		return x.Result
	}
	return nil
}

type StopEvaluationRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	AuditScopeId  string                 `protobuf:"bytes,1,opt,name=audit_scope_id,json=auditScopeId,proto3" json:"audit_scope_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StopEvaluationRequest) Reset() {
	*x = StopEvaluationRequest{}
	mi := &file_api_evaluation_evaluation_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StopEvaluationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopEvaluationRequest) ProtoMessage() {}

func (x *StopEvaluationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_evaluation_evaluation_proto_msgTypes[3]
	if x != nil {
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
	return file_api_evaluation_evaluation_proto_rawDescGZIP(), []int{3}
}

func (x *StopEvaluationRequest) GetAuditScopeId() string {
	if x != nil {
		return x.AuditScopeId
	}
	return ""
}

type StopEvaluationResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StopEvaluationResponse) Reset() {
	*x = StopEvaluationResponse{}
	mi := &file_api_evaluation_evaluation_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StopEvaluationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopEvaluationResponse) ProtoMessage() {}

func (x *StopEvaluationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_evaluation_evaluation_proto_msgTypes[4]
	if x != nil {
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
	return file_api_evaluation_evaluation_proto_rawDescGZIP(), []int{4}
}

type ListEvaluationResultsRequest struct {
	state  protoimpl.MessageState               `protogen:"open.v1"`
	Filter *ListEvaluationResultsRequest_Filter `protobuf:"bytes,1,opt,name=filter,proto3,oneof" json:"filter,omitempty"`
	// Optional. Latest results grouped by control_id.
	LatestByControlId *bool  `protobuf:"varint,2,opt,name=latest_by_control_id,json=latestByControlId,proto3,oneof" json:"latest_by_control_id,omitempty"`
	PageSize          int32  `protobuf:"varint,10,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	PageToken         string `protobuf:"bytes,11,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	OrderBy           string `protobuf:"bytes,12,opt,name=order_by,json=orderBy,proto3" json:"order_by,omitempty"`
	Asc               bool   `protobuf:"varint,13,opt,name=asc,proto3" json:"asc,omitempty"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *ListEvaluationResultsRequest) Reset() {
	*x = ListEvaluationResultsRequest{}
	mi := &file_api_evaluation_evaluation_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListEvaluationResultsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListEvaluationResultsRequest) ProtoMessage() {}

func (x *ListEvaluationResultsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_evaluation_evaluation_proto_msgTypes[5]
	if x != nil {
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
	return file_api_evaluation_evaluation_proto_rawDescGZIP(), []int{5}
}

func (x *ListEvaluationResultsRequest) GetFilter() *ListEvaluationResultsRequest_Filter {
	if x != nil {
		return x.Filter
	}
	return nil
}

func (x *ListEvaluationResultsRequest) GetLatestByControlId() bool {
	if x != nil && x.LatestByControlId != nil {
		return *x.LatestByControlId
	}
	return false
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
	state         protoimpl.MessageState `protogen:"open.v1"`
	Results       []*EvaluationResult    `protobuf:"bytes,1,rep,name=results,proto3" json:"results,omitempty"`
	NextPageToken string                 `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListEvaluationResultsResponse) Reset() {
	*x = ListEvaluationResultsResponse{}
	mi := &file_api_evaluation_evaluation_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListEvaluationResultsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListEvaluationResultsResponse) ProtoMessage() {}

func (x *ListEvaluationResultsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_evaluation_evaluation_proto_msgTypes[6]
	if x != nil {
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
	return file_api_evaluation_evaluation_proto_rawDescGZIP(), []int{6}
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

// A evaluation result resource, representing the result after evaluating the
// target of evaluation with a specific control target_of_evaluation_id, category_name and
// catalog_id are necessary to get the corresponding AuditScope
type EvaluationResult struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Evaluation result id
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// The Target of Evaluation ID the evaluation belongs to
	TargetOfEvaluationId string `protobuf:"bytes,2,opt,name=target_of_evaluation_id,json=targetOfEvaluationId,proto3" json:"target_of_evaluation_id,omitempty"`
	// The Audit Scope ID the evaluation belongs to
	AuditScopeId string `protobuf:"bytes,3,opt,name=audit_scope_id,json=auditScopeId,proto3" json:"audit_scope_id,omitempty"`
	// The control id the evaluation was based on
	ControlId string `protobuf:"bytes,4,opt,name=control_id,json=controlId,proto3" json:"control_id,omitempty"`
	// The category the evaluated control belongs to
	ControlCategoryName string `protobuf:"bytes,5,opt,name=control_category_name,json=controlCategoryName,proto3" json:"control_category_name,omitempty"`
	// The catalog the evaluated control belongs to
	ControlCatalogId string `protobuf:"bytes,6,opt,name=control_catalog_id,json=controlCatalogId,proto3" json:"control_catalog_id,omitempty"`
	// Optionally, specifies the parent control ID, if this is a sub-control
	ParentControlId *string `protobuf:"bytes,7,opt,name=parent_control_id,json=parentControlId,proto3,oneof" json:"parent_control_id,omitempty"`
	// Evaluation status
	Status EvaluationStatus `protobuf:"varint,8,opt,name=status,proto3,enum=clouditor.evaluation.v1.EvaluationStatus" json:"status,omitempty"`
	// Time of evaluation
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=timestamp,proto3" json:"timestamp,omitempty" gorm:"serializer:timestamppb;type:timestamp"`
	// List of assessment results because of which the evaluation status is compliant or not compliant
	AssessmentResultIds []string `protobuf:"bytes,10,rep,name=assessment_result_ids,json=assessmentResultIds,proto3" json:"assessment_result_ids,omitempty" gorm:"serializer:json"`
	Comment             *string  `protobuf:"bytes,11,opt,name=comment,proto3,oneof" json:"comment,omitempty"`
	// Optional, but required if the status is one of the "manually" ones. This
	// denotes how long the (manual) created evaluation result is valid. During
	// this time, no automatic results are generated for the specific control.
	ValidUntil    *timestamppb.Timestamp `protobuf:"bytes,20,opt,name=valid_until,json=validUntil,proto3,oneof" json:"valid_until,omitempty" gorm:"serializer:timestamppb;type:timestamp"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EvaluationResult) Reset() {
	*x = EvaluationResult{}
	mi := &file_api_evaluation_evaluation_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EvaluationResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EvaluationResult) ProtoMessage() {}

func (x *EvaluationResult) ProtoReflect() protoreflect.Message {
	mi := &file_api_evaluation_evaluation_proto_msgTypes[7]
	if x != nil {
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
	return file_api_evaluation_evaluation_proto_rawDescGZIP(), []int{7}
}

func (x *EvaluationResult) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *EvaluationResult) GetTargetOfEvaluationId() string {
	if x != nil {
		return x.TargetOfEvaluationId
	}
	return ""
}

func (x *EvaluationResult) GetAuditScopeId() string {
	if x != nil {
		return x.AuditScopeId
	}
	return ""
}

func (x *EvaluationResult) GetControlId() string {
	if x != nil {
		return x.ControlId
	}
	return ""
}

func (x *EvaluationResult) GetControlCategoryName() string {
	if x != nil {
		return x.ControlCategoryName
	}
	return ""
}

func (x *EvaluationResult) GetControlCatalogId() string {
	if x != nil {
		return x.ControlCatalogId
	}
	return ""
}

func (x *EvaluationResult) GetParentControlId() string {
	if x != nil && x.ParentControlId != nil {
		return *x.ParentControlId
	}
	return ""
}

func (x *EvaluationResult) GetStatus() EvaluationStatus {
	if x != nil {
		return x.Status
	}
	return EvaluationStatus_EVALUATION_STATUS_UNSPECIFIED
}

func (x *EvaluationResult) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *EvaluationResult) GetAssessmentResultIds() []string {
	if x != nil {
		return x.AssessmentResultIds
	}
	return nil
}

func (x *EvaluationResult) GetComment() string {
	if x != nil && x.Comment != nil {
		return *x.Comment
	}
	return ""
}

func (x *EvaluationResult) GetValidUntil() *timestamppb.Timestamp {
	if x != nil {
		return x.ValidUntil
	}
	return nil
}

type ListEvaluationResultsRequest_Filter struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Optional. Lists only evaluation results for a specific target of evaluation.
	TargetOfEvaluationId *string `protobuf:"bytes,1,opt,name=target_of_evaluation_id,json=targetOfEvaluationId,proto3,oneof" json:"target_of_evaluation_id,omitempty"`
	// Optional. Lists only evaluation results for a specific catalog.
	CatalogId *string `protobuf:"bytes,2,opt,name=catalog_id,json=catalogId,proto3,oneof" json:"catalog_id,omitempty"`
	// Optional. Lists only evaluation results for a specific control id.
	ControlId *string `protobuf:"bytes,3,opt,name=control_id,json=controlId,proto3,oneof" json:"control_id,omitempty"`
	// Optional. Lists all evaluation results for the given initial control id
	// substring, e.g., if the substring 'CMK-01.' is given it returns the
	// controls CMK-01.1B, CMK-01.1S, CMK-01.1H.
	SubControls *string `protobuf:"bytes,4,opt,name=sub_controls,json=subControls,proto3,oneof" json:"sub_controls,omitempty"`
	// Optional. Lists only results for parent controls
	ParentsOnly *bool `protobuf:"varint,5,opt,name=parents_only,json=parentsOnly,proto3,oneof" json:"parents_only,omitempty"`
	// Optional. Lists only manual results in their validity period
	ValidManualOnly *bool `protobuf:"varint,6,opt,name=valid_manual_only,json=validManualOnly,proto3,oneof" json:"valid_manual_only,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *ListEvaluationResultsRequest_Filter) Reset() {
	*x = ListEvaluationResultsRequest_Filter{}
	mi := &file_api_evaluation_evaluation_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListEvaluationResultsRequest_Filter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListEvaluationResultsRequest_Filter) ProtoMessage() {}

func (x *ListEvaluationResultsRequest_Filter) ProtoReflect() protoreflect.Message {
	mi := &file_api_evaluation_evaluation_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListEvaluationResultsRequest_Filter.ProtoReflect.Descriptor instead.
func (*ListEvaluationResultsRequest_Filter) Descriptor() ([]byte, []int) {
	return file_api_evaluation_evaluation_proto_rawDescGZIP(), []int{5, 0}
}

func (x *ListEvaluationResultsRequest_Filter) GetTargetOfEvaluationId() string {
	if x != nil && x.TargetOfEvaluationId != nil {
		return *x.TargetOfEvaluationId
	}
	return ""
}

func (x *ListEvaluationResultsRequest_Filter) GetCatalogId() string {
	if x != nil && x.CatalogId != nil {
		return *x.CatalogId
	}
	return ""
}

func (x *ListEvaluationResultsRequest_Filter) GetControlId() string {
	if x != nil && x.ControlId != nil {
		return *x.ControlId
	}
	return ""
}

func (x *ListEvaluationResultsRequest_Filter) GetSubControls() string {
	if x != nil && x.SubControls != nil {
		return *x.SubControls
	}
	return ""
}

func (x *ListEvaluationResultsRequest_Filter) GetParentsOnly() bool {
	if x != nil && x.ParentsOnly != nil {
		return *x.ParentsOnly
	}
	return false
}

func (x *ListEvaluationResultsRequest_Filter) GetValidManualOnly() bool {
	if x != nil && x.ValidManualOnly != nil {
		return *x.ValidManualOnly
	}
	return false
}

var File_api_evaluation_evaluation_proto protoreflect.FileDescriptor

const file_api_evaluation_evaluation_proto_rawDesc = "" +
	"\n" +
	"\x1fapi/evaluation/evaluation.proto\x12\x17clouditor.evaluation.v1\x1a\x1bbuf/validate/validate.proto\x1a\x1cgoogle/api/annotations.proto\x1a\x1fgoogle/api/field_behavior.proto\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x13tagger/tagger.proto\"\x82\x01\n" +
	"\x16StartEvaluationRequest\x121\n" +
	"\x0eaudit_scope_id\x18\x01 \x01(\tB\v\xe0A\x02\xbaH\x05r\x03\xb0\x01\x01R\fauditScopeId\x12(\n" +
	"\binterval\x18\x03 \x01(\x05B\a\xbaH\x04\x1a\x02 \x00H\x00R\binterval\x88\x01\x01B\v\n" +
	"\t_interval\"9\n" +
	"\x17StartEvaluationResponse\x12\x1e\n" +
	"\n" +
	"successful\x18\x01 \x01(\bR\n" +
	"successful\"j\n" +
	"\x1dCreateEvaluationResultRequest\x12I\n" +
	"\x06result\x18\x01 \x01(\v2).clouditor.evaluation.v1.EvaluationResultB\x06\xbaH\x03\xc8\x01\x01R\x06result\"J\n" +
	"\x15StopEvaluationRequest\x121\n" +
	"\x0eaudit_scope_id\x18\x01 \x01(\tB\v\xe0A\x02\xbaH\x05r\x03\xb0\x01\x01R\fauditScopeId\"\x18\n" +
<<<<<<< HEAD
	"\x16StopEvaluationResponse\"\xe3\x05\n" +
=======
	"\x16StopEvaluationResponse\"\xe4\x05\n" +
>>>>>>> b9e0eb0f (Add `User` message to Orchestrator folder)
	"\x1cListEvaluationResultsRequest\x12Y\n" +
	"\x06filter\x18\x01 \x01(\v2<.clouditor.evaluation.v1.ListEvaluationResultsRequest.FilterH\x00R\x06filter\x88\x01\x01\x124\n" +
	"\x14latest_by_control_id\x18\x02 \x01(\bH\x01R\x11latestByControlId\x88\x01\x01\x12\x1b\n" +
	"\tpage_size\x18\n" +
	" \x01(\x05R\bpageSize\x12\x1d\n" +
	"\n" +
	"page_token\x18\v \x01(\tR\tpageToken\x12\x19\n" +
	"\border_by\x18\f \x01(\tR\aorderBy\x12\x10\n" +
<<<<<<< HEAD
	"\x03asc\x18\r \x01(\bR\x03asc\x1a\xa4\x03\n" +
	"\x06Filter\x12D\n" +
	"\x17target_of_evaluation_id\x18\x01 \x01(\tB\b\xbaH\x05r\x03\xb0\x01\x01H\x00R\x14targetOfEvaluationId\x88\x01\x01\x12+\n" +
=======
	"\x03asc\x18\r \x01(\bR\x03asc\x1a\xa5\x03\n" +
	"\x06Filter\x12E\n" +
	"\x17certification_target_id\x18\x01 \x01(\tB\b\xbaH\x05r\x03\xb0\x01\x01H\x00R\x15certificationTargetId\x88\x01\x01\x12+\n" +
>>>>>>> b9e0eb0f (Add `User` message to Orchestrator folder)
	"\n" +
	"catalog_id\x18\x02 \x01(\tB\a\xbaH\x04r\x02\x10\x01H\x01R\tcatalogId\x88\x01\x01\x12+\n" +
	"\n" +
	"control_id\x18\x03 \x01(\tB\a\xbaH\x04r\x02\x10\x01H\x02R\tcontrolId\x88\x01\x01\x12/\n" +
	"\fsub_controls\x18\x04 \x01(\tB\a\xbaH\x04r\x02\x10\x01H\x03R\vsubControls\x88\x01\x01\x12&\n" +
	"\fparents_only\x18\x05 \x01(\bH\x04R\vparentsOnly\x88\x01\x01\x12/\n" +
	"\x11valid_manual_only\x18\x06 \x01(\bH\x05R\x0fvalidManualOnly\x88\x01\x01B\x1a\n" +
<<<<<<< HEAD
	"\x18_target_of_evaluation_idB\r\n" +
=======
	"\x18_certification_target_idB\r\n" +
>>>>>>> b9e0eb0f (Add `User` message to Orchestrator folder)
	"\v_catalog_idB\r\n" +
	"\v_control_idB\x0f\n" +
	"\r_sub_controlsB\x0f\n" +
	"\r_parents_onlyB\x14\n" +
	"\x12_valid_manual_onlyB\t\n" +
	"\a_filterB\x17\n" +
	"\x15_latest_by_control_id\"\x8c\x01\n" +
	"\x1dListEvaluationResultsResponse\x12C\n" +
	"\aresults\x18\x01 \x03(\v2).clouditor.evaluation.v1.EvaluationResultR\aresults\x12&\n" +
<<<<<<< HEAD
	"\x0fnext_page_token\x18\x02 \x01(\tR\rnextPageToken\"\xd3\x06\n" +
	"\x10EvaluationResult\x12\x1b\n" +
	"\x02id\x18\x01 \x01(\tB\v\xe0A\x02\xbaH\x05r\x03\xb0\x01\x01R\x02id\x12?\n" +
	"\x17target_of_evaluation_id\x18\x02 \x01(\tB\b\xbaH\x05r\x03\xb0\x01\x01R\x14targetOfEvaluationId\x12.\n" +
=======
	"\x0fnext_page_token\x18\x02 \x01(\tR\rnextPageToken\"\xd4\x06\n" +
	"\x10EvaluationResult\x12\x1b\n" +
	"\x02id\x18\x01 \x01(\tB\v\xe0A\x02\xbaH\x05r\x03\xb0\x01\x01R\x02id\x12@\n" +
	"\x17certification_target_id\x18\x02 \x01(\tB\b\xbaH\x05r\x03\xb0\x01\x01R\x15certificationTargetId\x12.\n" +
>>>>>>> b9e0eb0f (Add `User` message to Orchestrator folder)
	"\x0eaudit_scope_id\x18\x03 \x01(\tB\b\xbaH\x05r\x03\xb0\x01\x01R\fauditScopeId\x12&\n" +
	"\n" +
	"control_id\x18\x04 \x01(\tB\a\xbaH\x04r\x02\x10\x01R\tcontrolId\x12;\n" +
	"\x15control_category_name\x18\x05 \x01(\tB\a\xbaH\x04r\x02\x10\x01R\x13controlCategoryName\x125\n" +
	"\x12control_catalog_id\x18\x06 \x01(\tB\a\xbaH\x04r\x02\x10\x01R\x10controlCatalogId\x12/\n" +
	"\x11parent_control_id\x18\a \x01(\tH\x00R\x0fparentControlId\x88\x01\x01\x12N\n" +
	"\x06status\x18\b \x01(\x0e2).clouditor.evaluation.v1.EvaluationStatusB\v\xe0A\x02\xbaH\x05\x82\x01\x02\x10\x01R\x06status\x12n\n" +
	"\ttimestamp\x18\t \x01(\v2\x1a.google.protobuf.TimestampB4\xe0A\x02\x9a\x84\x9e\x03,gorm:\"serializer:timestamppb;type:timestamp\"R\ttimestamp\x12^\n" +
	"\x15assessment_result_ids\x18\n" +
	" \x03(\tB*\xe0A\x02\xbaH\t\x92\x01\x06\"\x04r\x02\x10\x01\x9a\x84\x9e\x03\x16gorm:\"serializer:json\"R\x13assessmentResultIds\x12\x1d\n" +
	"\acomment\x18\v \x01(\tH\x01R\acomment\x88\x01\x01\x12s\n" +
	"\vvalid_until\x18\x14 \x01(\v2\x1a.google.protobuf.TimestampB1\x9a\x84\x9e\x03,gorm:\"serializer:timestamppb;type:timestamp\"H\x02R\n" +
	"validUntil\x88\x01\x01B\x14\n" +
	"\x12_parent_control_idB\n" +
	"\n" +
	"\b_commentB\x0e\n" +
	"\f_valid_until*\xf2\x01\n" +
	"\x10EvaluationStatus\x12!\n" +
	"\x1dEVALUATION_STATUS_UNSPECIFIED\x10\x00\x12\x1f\n" +
	"\x1bEVALUATION_STATUS_COMPLIANT\x10\x01\x12(\n" +
	"$EVALUATION_STATUS_COMPLIANT_MANUALLY\x10\x02\x12#\n" +
	"\x1fEVALUATION_STATUS_NOT_COMPLIANT\x10\x03\x12,\n" +
	"(EVALUATION_STATUS_NOT_COMPLIANT_MANUALLY\x10\x04\x12\x1d\n" +
	"\x19EVALUATION_STATUS_PENDING\x10\n" +
	"2\xb5\x05\n" +
	"\n" +
	"Evaluation\x12\xac\x01\n" +
	"\x0fStartEvaluation\x12/.clouditor.evaluation.v1.StartEvaluationRequest\x1a0.clouditor.evaluation.v1.StartEvaluationResponse\"6\x82\xd3\xe4\x93\x020\"./v1/evaluation/evaluate/{audit_scope_id}/start\x12\xa8\x01\n" +
	"\x0eStopEvaluation\x12..clouditor.evaluation.v1.StopEvaluationRequest\x1a/.clouditor.evaluation.v1.StopEvaluationResponse\"5\x82\xd3\xe4\x93\x02/\"-/v1/evaluation/evaluate/{audit_scope_id}/stop\x12\xa6\x01\n" +
	"\x15ListEvaluationResults\x125.clouditor.evaluation.v1.ListEvaluationResultsRequest\x1a6.clouditor.evaluation.v1.ListEvaluationResultsResponse\"\x1e\x82\xd3\xe4\x93\x02\x18\x12\x16/v1/evaluation/results\x12\xa3\x01\n" +
	"\x16CreateEvaluationResult\x126.clouditor.evaluation.v1.CreateEvaluationResultRequest\x1a).clouditor.evaluation.v1.EvaluationResult\"&\x82\xd3\xe4\x93\x02 :\x06result\"\x16/v1/evaluation/resultsB*Z(clouditor.io/clouditor/v2/api/evaluationb\x06proto3"

var (
	file_api_evaluation_evaluation_proto_rawDescOnce sync.Once
	file_api_evaluation_evaluation_proto_rawDescData []byte
)

func file_api_evaluation_evaluation_proto_rawDescGZIP() []byte {
	file_api_evaluation_evaluation_proto_rawDescOnce.Do(func() {
		file_api_evaluation_evaluation_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_evaluation_evaluation_proto_rawDesc), len(file_api_evaluation_evaluation_proto_rawDesc)))
	})
	return file_api_evaluation_evaluation_proto_rawDescData
}

var file_api_evaluation_evaluation_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_evaluation_evaluation_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_api_evaluation_evaluation_proto_goTypes = []any{
	(EvaluationStatus)(0),                       // 0: clouditor.evaluation.v1.EvaluationStatus
	(*StartEvaluationRequest)(nil),              // 1: clouditor.evaluation.v1.StartEvaluationRequest
	(*StartEvaluationResponse)(nil),             // 2: clouditor.evaluation.v1.StartEvaluationResponse
	(*CreateEvaluationResultRequest)(nil),       // 3: clouditor.evaluation.v1.CreateEvaluationResultRequest
	(*StopEvaluationRequest)(nil),               // 4: clouditor.evaluation.v1.StopEvaluationRequest
	(*StopEvaluationResponse)(nil),              // 5: clouditor.evaluation.v1.StopEvaluationResponse
	(*ListEvaluationResultsRequest)(nil),        // 6: clouditor.evaluation.v1.ListEvaluationResultsRequest
	(*ListEvaluationResultsResponse)(nil),       // 7: clouditor.evaluation.v1.ListEvaluationResultsResponse
	(*EvaluationResult)(nil),                    // 8: clouditor.evaluation.v1.EvaluationResult
	(*ListEvaluationResultsRequest_Filter)(nil), // 9: clouditor.evaluation.v1.ListEvaluationResultsRequest.Filter
	(*timestamppb.Timestamp)(nil),               // 10: google.protobuf.Timestamp
}
var file_api_evaluation_evaluation_proto_depIdxs = []int32{
	8,  // 0: clouditor.evaluation.v1.CreateEvaluationResultRequest.result:type_name -> clouditor.evaluation.v1.EvaluationResult
	9,  // 1: clouditor.evaluation.v1.ListEvaluationResultsRequest.filter:type_name -> clouditor.evaluation.v1.ListEvaluationResultsRequest.Filter
	8,  // 2: clouditor.evaluation.v1.ListEvaluationResultsResponse.results:type_name -> clouditor.evaluation.v1.EvaluationResult
	0,  // 3: clouditor.evaluation.v1.EvaluationResult.status:type_name -> clouditor.evaluation.v1.EvaluationStatus
	10, // 4: clouditor.evaluation.v1.EvaluationResult.timestamp:type_name -> google.protobuf.Timestamp
	10, // 5: clouditor.evaluation.v1.EvaluationResult.valid_until:type_name -> google.protobuf.Timestamp
	1,  // 6: clouditor.evaluation.v1.Evaluation.StartEvaluation:input_type -> clouditor.evaluation.v1.StartEvaluationRequest
	4,  // 7: clouditor.evaluation.v1.Evaluation.StopEvaluation:input_type -> clouditor.evaluation.v1.StopEvaluationRequest
	6,  // 8: clouditor.evaluation.v1.Evaluation.ListEvaluationResults:input_type -> clouditor.evaluation.v1.ListEvaluationResultsRequest
	3,  // 9: clouditor.evaluation.v1.Evaluation.CreateEvaluationResult:input_type -> clouditor.evaluation.v1.CreateEvaluationResultRequest
	2,  // 10: clouditor.evaluation.v1.Evaluation.StartEvaluation:output_type -> clouditor.evaluation.v1.StartEvaluationResponse
	5,  // 11: clouditor.evaluation.v1.Evaluation.StopEvaluation:output_type -> clouditor.evaluation.v1.StopEvaluationResponse
	7,  // 12: clouditor.evaluation.v1.Evaluation.ListEvaluationResults:output_type -> clouditor.evaluation.v1.ListEvaluationResultsResponse
	8,  // 13: clouditor.evaluation.v1.Evaluation.CreateEvaluationResult:output_type -> clouditor.evaluation.v1.EvaluationResult
	10, // [10:14] is the sub-list for method output_type
	6,  // [6:10] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_api_evaluation_evaluation_proto_init() }
func file_api_evaluation_evaluation_proto_init() {
	if File_api_evaluation_evaluation_proto != nil {
		return
	}
	file_api_evaluation_evaluation_proto_msgTypes[0].OneofWrappers = []any{}
	file_api_evaluation_evaluation_proto_msgTypes[5].OneofWrappers = []any{}
	file_api_evaluation_evaluation_proto_msgTypes[7].OneofWrappers = []any{}
	file_api_evaluation_evaluation_proto_msgTypes[8].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_evaluation_evaluation_proto_rawDesc), len(file_api_evaluation_evaluation_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_evaluation_evaluation_proto_goTypes,
		DependencyIndexes: file_api_evaluation_evaluation_proto_depIdxs,
		EnumInfos:         file_api_evaluation_evaluation_proto_enumTypes,
		MessageInfos:      file_api_evaluation_evaluation_proto_msgTypes,
	}.Build()
	File_api_evaluation_evaluation_proto = out.File
	file_api_evaluation_evaluation_proto_goTypes = nil
	file_api_evaluation_evaluation_proto_depIdxs = nil
}
