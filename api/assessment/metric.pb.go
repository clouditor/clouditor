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
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: api/assessment/metric.proto

package assessment

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	_ "github.com/srikrsna/protoc-gen-gotag/tagger"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
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

// The values a Scale accepts
type Metric_Scale int32

const (
	Metric_SCALE_UNSPECIFIED Metric_Scale = 0
	Metric_NOMINAL           Metric_Scale = 1
	Metric_ORDINAL           Metric_Scale = 2
	Metric_METRIC            Metric_Scale = 3
)

// Enum value maps for Metric_Scale.
var (
	Metric_Scale_name = map[int32]string{
		0: "SCALE_UNSPECIFIED",
		1: "NOMINAL",
		2: "ORDINAL",
		3: "METRIC",
	}
	Metric_Scale_value = map[string]int32{
		"SCALE_UNSPECIFIED": 0,
		"NOMINAL":           1,
		"ORDINAL":           2,
		"METRIC":            3,
	}
)

func (x Metric_Scale) Enum() *Metric_Scale {
	p := new(Metric_Scale)
	*p = x
	return p
}

func (x Metric_Scale) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Metric_Scale) Descriptor() protoreflect.EnumDescriptor {
	return file_api_assessment_metric_proto_enumTypes[0].Descriptor()
}

func (Metric_Scale) Type() protoreflect.EnumType {
	return &file_api_assessment_metric_proto_enumTypes[0]
}

func (x Metric_Scale) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Metric_Scale.Descriptor instead.
func (Metric_Scale) EnumDescriptor() ([]byte, []int) {
	return file_api_assessment_metric_proto_rawDescGZIP(), []int{0, 0}
}

type MetricImplementation_Language int32

const (
	MetricImplementation_LANGUAGE_UNSPECIFIED MetricImplementation_Language = 0
	MetricImplementation_LANGUAGE_REGO        MetricImplementation_Language = 1
)

// Enum value maps for MetricImplementation_Language.
var (
	MetricImplementation_Language_name = map[int32]string{
		0: "LANGUAGE_UNSPECIFIED",
		1: "LANGUAGE_REGO",
	}
	MetricImplementation_Language_value = map[string]int32{
		"LANGUAGE_UNSPECIFIED": 0,
		"LANGUAGE_REGO":        1,
	}
)

func (x MetricImplementation_Language) Enum() *MetricImplementation_Language {
	p := new(MetricImplementation_Language)
	*p = x
	return p
}

func (x MetricImplementation_Language) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MetricImplementation_Language) Descriptor() protoreflect.EnumDescriptor {
	return file_api_assessment_metric_proto_enumTypes[1].Descriptor()
}

func (MetricImplementation_Language) Type() protoreflect.EnumType {
	return &file_api_assessment_metric_proto_enumTypes[1]
}

func (x MetricImplementation_Language) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MetricImplementation_Language.Descriptor instead.
func (MetricImplementation_Language) EnumDescriptor() ([]byte, []int) {
	return file_api_assessment_metric_proto_rawDescGZIP(), []int{6, 0}
}

// A metric resource
type Metric struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Required. The unique identifier of the metric.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Required. The human readable name of the metric.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// The description of the metric
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	// The reference to control catalog category or domain
	Category string `protobuf:"bytes,4,opt,name=category,proto3" json:"category,omitempty"`
	// The scale of this metric, e.g. categories, ranked data or metric values.
	Scale Metric_Scale `protobuf:"varint,5,opt,name=scale,proto3,enum=clouditor.assessment.v1.Metric_Scale" json:"scale,omitempty"`
	// The range of this metric. Depending on the scale.
	Range *Range `protobuf:"bytes,6,opt,name=range,proto3" json:"range,omitempty"`
	// The interval in seconds the evidences must be collected for the respective
	// metric.
	Interval *durationpb.Duration `protobuf:"bytes,7,opt,name=interval,proto3" json:"interval,omitempty" gorm:"serializer:durationpb;type:interval"`
	// The implementation of this metric. This ensures that we are modelling an
	// association between a Metric and its MetricImplementation.
	Implementation *MetricImplementation `protobuf:"bytes,8,opt,name=implementation,proto3,oneof" json:"implementation,omitempty"`
	// Optional, but required if the metric is removed. The metric is not deleted
	// for backward compatibility and the timestamp is set to the time of removal.
	DeprecatedSince *timestamppb.Timestamp `protobuf:"bytes,10,opt,name=deprecated_since,json=deprecatedSince,proto3,oneof" json:"deprecated_since,omitempty" gorm:"serializer:timestamppb;type:timestamp"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *Metric) Reset() {
	*x = Metric{}
	mi := &file_api_assessment_metric_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Metric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric) ProtoMessage() {}

func (x *Metric) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_metric_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metric.ProtoReflect.Descriptor instead.
func (*Metric) Descriptor() ([]byte, []int) {
	return file_api_assessment_metric_proto_rawDescGZIP(), []int{0}
}

func (x *Metric) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Metric) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Metric) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Metric) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *Metric) GetScale() Metric_Scale {
	if x != nil {
		return x.Scale
	}
	return Metric_SCALE_UNSPECIFIED
}

func (x *Metric) GetRange() *Range {
	if x != nil {
		return x.Range
	}
	return nil
}

func (x *Metric) GetInterval() *durationpb.Duration {
	if x != nil {
		return x.Interval
	}
	return nil
}

func (x *Metric) GetImplementation() *MetricImplementation {
	if x != nil {
		return x.Implementation
	}
	return nil
}

func (x *Metric) GetDeprecatedSince() *timestamppb.Timestamp {
	if x != nil {
		return x.DeprecatedSince
	}
	return nil
}

// A range resource representing the range of values
type Range struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Required.
	//
	// Types that are valid to be assigned to Range:
	//
	//	*Range_AllowedValues
	//	*Range_Order
	//	*Range_MinMax
	Range         isRange_Range `protobuf_oneof:"range"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Range) Reset() {
	*x = Range{}
	mi := &file_api_assessment_metric_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Range) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Range) ProtoMessage() {}

func (x *Range) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_metric_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Range.ProtoReflect.Descriptor instead.
func (*Range) Descriptor() ([]byte, []int) {
	return file_api_assessment_metric_proto_rawDescGZIP(), []int{1}
}

func (x *Range) GetRange() isRange_Range {
	if x != nil {
		return x.Range
	}
	return nil
}

func (x *Range) GetAllowedValues() *AllowedValues {
	if x != nil {
		if x, ok := x.Range.(*Range_AllowedValues); ok {
			return x.AllowedValues
		}
	}
	return nil
}

func (x *Range) GetOrder() *Order {
	if x != nil {
		if x, ok := x.Range.(*Range_Order); ok {
			return x.Order
		}
	}
	return nil
}

func (x *Range) GetMinMax() *MinMax {
	if x != nil {
		if x, ok := x.Range.(*Range_MinMax); ok {
			return x.MinMax
		}
	}
	return nil
}

type isRange_Range interface {
	isRange_Range()
}

type Range_AllowedValues struct {
	// used for nominal scale
	AllowedValues *AllowedValues `protobuf:"bytes,1,opt,name=allowed_values,json=allowedValues,proto3,oneof"`
}

type Range_Order struct {
	// used for ordinal scale
	Order *Order `protobuf:"bytes,2,opt,name=order,proto3,oneof"`
}

type Range_MinMax struct {
	// used for metric scale
	MinMax *MinMax `protobuf:"bytes,3,opt,name=min_max,json=minMax,proto3,oneof"`
}

func (*Range_AllowedValues) isRange_Range() {}

func (*Range_Order) isRange_Range() {}

func (*Range_MinMax) isRange_Range() {}

// Defines a range of values through a (inclusive) minimum and a maximum
type MinMax struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Required.
	Min int64 `protobuf:"varint,1,opt,name=min,proto3" json:"min,omitempty"`
	// Required.
	Max           int64 `protobuf:"varint,2,opt,name=max,proto3" json:"max,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MinMax) Reset() {
	*x = MinMax{}
	mi := &file_api_assessment_metric_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MinMax) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MinMax) ProtoMessage() {}

func (x *MinMax) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_metric_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MinMax.ProtoReflect.Descriptor instead.
func (*MinMax) Descriptor() ([]byte, []int) {
	return file_api_assessment_metric_proto_rawDescGZIP(), []int{2}
}

func (x *MinMax) GetMin() int64 {
	if x != nil {
		return x.Min
	}
	return 0
}

func (x *MinMax) GetMax() int64 {
	if x != nil {
		return x.Max
	}
	return 0
}

// Defines a range
type AllowedValues struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Values        []*structpb.Value      `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AllowedValues) Reset() {
	*x = AllowedValues{}
	mi := &file_api_assessment_metric_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AllowedValues) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllowedValues) ProtoMessage() {}

func (x *AllowedValues) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_metric_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AllowedValues.ProtoReflect.Descriptor instead.
func (*AllowedValues) Descriptor() ([]byte, []int) {
	return file_api_assessment_metric_proto_rawDescGZIP(), []int{3}
}

func (x *AllowedValues) GetValues() []*structpb.Value {
	if x != nil {
		return x.Values
	}
	return nil
}

// Defines a range of values in a pre-defined order from the lowest to the
// highest.
type Order struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Values        []*structpb.Value      `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Order) Reset() {
	*x = Order{}
	mi := &file_api_assessment_metric_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Order) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Order) ProtoMessage() {}

func (x *Order) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_metric_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Order.ProtoReflect.Descriptor instead.
func (*Order) Descriptor() ([]byte, []int) {
	return file_api_assessment_metric_proto_rawDescGZIP(), []int{4}
}

func (x *Order) GetValues() []*structpb.Value {
	if x != nil {
		return x.Values
	}
	return nil
}

// Defines the operator and a target value for an individual metric
type MetricConfiguration struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The operator to compare the metric, such as == or >
	Operator string `protobuf:"bytes,1,opt,name=operator,proto3" json:"operator,omitempty"`
	// The target value
	TargetValue *structpb.Value `protobuf:"bytes,2,opt,name=target_value,json=targetValue,proto3" json:"target_value,omitempty" gorm:"serializer:json"`
	// Whether this configuration is a default configuration
	IsDefault bool `protobuf:"varint,3,opt,name=is_default,json=isDefault,proto3" json:"is_default,omitempty"`
	// The last time of update
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty" gorm:"serializer:timestamppb;type:timestamp"`
	// The metric this configuration belongs to
	MetricId string `protobuf:"bytes,5,opt,name=metric_id,json=metricId,proto3" json:"metric_id,omitempty" gorm:"primaryKey"`
	// The certification target this configuration belongs to
	CertificationTargetId string `protobuf:"bytes,6,opt,name=certification_target_id,json=certificationTargetId,proto3" json:"certification_target_id,omitempty" gorm:"primaryKey"`
	unknownFields         protoimpl.UnknownFields
	sizeCache             protoimpl.SizeCache
}

func (x *MetricConfiguration) Reset() {
	*x = MetricConfiguration{}
	mi := &file_api_assessment_metric_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MetricConfiguration) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetricConfiguration) ProtoMessage() {}

func (x *MetricConfiguration) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_metric_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetricConfiguration.ProtoReflect.Descriptor instead.
func (*MetricConfiguration) Descriptor() ([]byte, []int) {
	return file_api_assessment_metric_proto_rawDescGZIP(), []int{5}
}

func (x *MetricConfiguration) GetOperator() string {
	if x != nil {
		return x.Operator
	}
	return ""
}

func (x *MetricConfiguration) GetTargetValue() *structpb.Value {
	if x != nil {
		return x.TargetValue
	}
	return nil
}

func (x *MetricConfiguration) GetIsDefault() bool {
	if x != nil {
		return x.IsDefault
	}
	return false
}

func (x *MetricConfiguration) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *MetricConfiguration) GetMetricId() string {
	if x != nil {
		return x.MetricId
	}
	return ""
}

func (x *MetricConfiguration) GetCertificationTargetId() string {
	if x != nil {
		return x.CertificationTargetId
	}
	return ""
}

// MetricImplementation defines the implementation of an individual metric.
type MetricImplementation struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The metric which is implemented
	MetricId string `protobuf:"bytes,1,opt,name=metric_id,json=metricId,proto3" json:"metric_id,omitempty" gorm:"primaryKey"`
	// The language this metric is implemented in
	Lang MetricImplementation_Language `protobuf:"varint,2,opt,name=lang,proto3,enum=clouditor.assessment.v1.MetricImplementation_Language" json:"lang,omitempty"`
	// The actual implementation
	Code string `protobuf:"bytes,3,opt,name=code,proto3" json:"code,omitempty"`
	// The last time of update
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty" gorm:"serializer:timestamppb;type:timestamp"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MetricImplementation) Reset() {
	*x = MetricImplementation{}
	mi := &file_api_assessment_metric_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MetricImplementation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetricImplementation) ProtoMessage() {}

func (x *MetricImplementation) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_metric_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetricImplementation.ProtoReflect.Descriptor instead.
func (*MetricImplementation) Descriptor() ([]byte, []int) {
	return file_api_assessment_metric_proto_rawDescGZIP(), []int{6}
}

func (x *MetricImplementation) GetMetricId() string {
	if x != nil {
		return x.MetricId
	}
	return ""
}

func (x *MetricImplementation) GetLang() MetricImplementation_Language {
	if x != nil {
		return x.Lang
	}
	return MetricImplementation_LANGUAGE_UNSPECIFIED
}

func (x *MetricImplementation) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *MetricImplementation) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

var File_api_assessment_metric_proto protoreflect.FileDescriptor

var file_api_assessment_metric_proto_rawDesc = string([]byte{
	0x0a, 0x1b, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74,
	0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x63,
	0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d,
	0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x62, 0x65, 0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x13, 0x74, 0x61, 0x67, 0x67, 0x65, 0x72, 0x2f, 0x74, 0x61, 0x67, 0x67,
	0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc1, 0x05, 0x0a, 0x06, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x12, 0x1a, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x0a, 0xe0, 0x41, 0x02, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x1e, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xe0,
	0x41, 0x02, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x23, 0x0a, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x08, 0x63, 0x61,
	0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x45, 0x0a, 0x05, 0x73, 0x63, 0x61, 0x6c, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x25, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f,
	0x72, 0x2e, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x53, 0x63, 0x61, 0x6c, 0x65, 0x42, 0x08, 0xba, 0x48,
	0x05, 0x82, 0x01, 0x02, 0x10, 0x01, 0x52, 0x05, 0x73, 0x63, 0x61, 0x6c, 0x65, 0x12, 0x3c, 0x0a,
	0x05, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x63,
	0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d,
	0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x42, 0x06, 0xba, 0x48,
	0x03, 0xc8, 0x01, 0x01, 0x52, 0x05, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x66, 0x0a, 0x08, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x2f, 0x9a, 0x84, 0x9e, 0x03, 0x2a, 0x67,
	0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x72, 0x3a,
	0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x3b, 0x74, 0x79, 0x70, 0x65, 0x3a,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x22, 0x52, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x76, 0x61, 0x6c, 0x12, 0x5a, 0x0a, 0x0e, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x63, 0x6c,
	0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65,
	0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x49, 0x6d, 0x70, 0x6c,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52, 0x0e, 0x69, 0x6d,
	0x70, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12,
	0x7d, 0x0a, 0x10, 0x64, 0x65, 0x70, 0x72, 0x65, 0x63, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x73, 0x69,
	0x6e, 0x63, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x31, 0x9a, 0x84, 0x9e, 0x03, 0x2c, 0x67, 0x6f, 0x72, 0x6d,
	0x3a, 0x22, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x72, 0x3a, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x70, 0x62, 0x3b, 0x74, 0x79, 0x70, 0x65, 0x3a, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x48, 0x01, 0x52, 0x0f, 0x64, 0x65, 0x70, 0x72,
	0x65, 0x63, 0x61, 0x74, 0x65, 0x64, 0x53, 0x69, 0x6e, 0x63, 0x65, 0x88, 0x01, 0x01, 0x22, 0x44,
	0x0a, 0x05, 0x53, 0x63, 0x61, 0x6c, 0x65, 0x12, 0x15, 0x0a, 0x11, 0x53, 0x43, 0x41, 0x4c, 0x45,
	0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0b,
	0x0a, 0x07, 0x4e, 0x4f, 0x4d, 0x49, 0x4e, 0x41, 0x4c, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x4f,
	0x52, 0x44, 0x49, 0x4e, 0x41, 0x4c, 0x10, 0x02, 0x12, 0x0a, 0x0a, 0x06, 0x4d, 0x45, 0x54, 0x52,
	0x49, 0x43, 0x10, 0x03, 0x42, 0x11, 0x0a, 0x0f, 0x5f, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x6d, 0x65,
	0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x13, 0x0a, 0x11, 0x5f, 0x64, 0x65, 0x70, 0x72,
	0x65, 0x63, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x73, 0x69, 0x6e, 0x63, 0x65, 0x22, 0xd5, 0x01, 0x0a,
	0x05, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x4f, 0x0a, 0x0e, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65,
	0x64, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26,
	0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x61, 0x73, 0x73, 0x65, 0x73,
	0x73, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x48, 0x00, 0x52, 0x0d, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65,
	0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x12, 0x36, 0x0a, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74,
	0x6f, 0x72, 0x2e, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31,
	0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x48, 0x00, 0x52, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x12,
	0x3a, 0x0a, 0x07, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x61, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1f, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x61, 0x73, 0x73,
	0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x69, 0x6e, 0x4d, 0x61,
	0x78, 0x48, 0x00, 0x52, 0x06, 0x6d, 0x69, 0x6e, 0x4d, 0x61, 0x78, 0x42, 0x07, 0x0a, 0x05, 0x72,
	0x61, 0x6e, 0x67, 0x65, 0x22, 0x2c, 0x0a, 0x06, 0x4d, 0x69, 0x6e, 0x4d, 0x61, 0x78, 0x12, 0x10,
	0x0a, 0x03, 0x6d, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x6d, 0x69, 0x6e,
	0x12, 0x10, 0x0a, 0x03, 0x6d, 0x61, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x6d,
	0x61, 0x78, 0x22, 0x3f, 0x0a, 0x0d, 0x41, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x73, 0x12, 0x2e, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x73, 0x22, 0x37, 0x0a, 0x05, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x2e, 0x0a, 0x06,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x22, 0xe5, 0x03, 0x0a,
	0x13, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x41, 0x0a, 0x08, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x25, 0xe0, 0x41, 0x02, 0xba, 0x48, 0x1f, 0x72, 0x1d,
	0x32, 0x1b, 0x5e, 0x28, 0x3c, 0x7c, 0x3e, 0x7c, 0x3c, 0x3d, 0x7c, 0x3e, 0x3d, 0x7c, 0x3d, 0x3d,
	0x7c, 0x69, 0x73, 0x49, 0x6e, 0x7c, 0x61, 0x6c, 0x6c, 0x49, 0x6e, 0x29, 0x24, 0x52, 0x08, 0x6f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x5f, 0x0a, 0x0c, 0x74, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x24, 0xe0, 0x41, 0x02, 0xba, 0x48, 0x03, 0xc8, 0x01, 0x01,
	0x9a, 0x84, 0x9e, 0x03, 0x16, 0x67, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x73, 0x65, 0x72, 0x69, 0x61,
	0x6c, 0x69, 0x7a, 0x65, 0x72, 0x3a, 0x6a, 0x73, 0x6f, 0x6e, 0x22, 0x52, 0x0b, 0x74, 0x61, 0x72,
	0x67, 0x65, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x22, 0x0a, 0x0a, 0x69, 0x73, 0x5f, 0x64,
	0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x42, 0x03, 0xe0, 0x41,
	0x02, 0x52, 0x09, 0x69, 0x73, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x12, 0x6c, 0x0a, 0x0a,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x31, 0x9a, 0x84,
	0x9e, 0x03, 0x2c, 0x67, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69,
	0x7a, 0x65, 0x72, 0x3a, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x70, 0x62, 0x3b,
	0x74, 0x79, 0x70, 0x65, 0x3a, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x52,
	0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x3d, 0x0a, 0x09, 0x6d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x42, 0x20, 0xe0,
	0x41, 0x02, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x9a, 0x84, 0x9e, 0x03, 0x11, 0x67, 0x6f,
	0x72, 0x6d, 0x3a, 0x22, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x4b, 0x65, 0x79, 0x22, 0x52,
	0x08, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x49, 0x64, 0x12, 0x59, 0x0a, 0x17, 0x63, 0x65, 0x72,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x42, 0x21, 0xe0, 0x41, 0x02, 0xba,
	0x48, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0x9a, 0x84, 0x9e, 0x03, 0x11, 0x67, 0x6f, 0x72, 0x6d,
	0x3a, 0x22, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x4b, 0x65, 0x79, 0x22, 0x52, 0x15, 0x63,
	0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x49, 0x64, 0x22, 0xf2, 0x02, 0x0a, 0x14, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x49,
	0x6d, 0x70, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x3d, 0x0a,
	0x09, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x20, 0xe0, 0x41, 0x02, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x9a, 0x84, 0x9e, 0x03,
	0x11, 0x67, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x4b, 0x65,
	0x79, 0x22, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x49, 0x64, 0x12, 0x54, 0x0a, 0x04,
	0x6c, 0x61, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x36, 0x2e, 0x63, 0x6c, 0x6f,
	0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e,
	0x74, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x49, 0x6d, 0x70, 0x6c, 0x65,
	0x6d, 0x65, 0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4c, 0x61, 0x6e, 0x67, 0x75, 0x61,
	0x67, 0x65, 0x42, 0x08, 0xba, 0x48, 0x05, 0x82, 0x01, 0x02, 0x10, 0x01, 0x52, 0x04, 0x6c, 0x61,
	0x6e, 0x67, 0x12, 0x1e, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x0a, 0xe0, 0x41, 0x02, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x12, 0x6c, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x42, 0x31, 0x9a, 0x84, 0x9e, 0x03, 0x2c, 0x67, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x73,
	0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x72, 0x3a, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x70, 0x62, 0x3b, 0x74, 0x79, 0x70, 0x65, 0x3a, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x22, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x22, 0x37, 0x0a, 0x08, 0x4c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x14,
	0x4c, 0x41, 0x4e, 0x47, 0x55, 0x41, 0x47, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49,
	0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x4c, 0x41, 0x4e, 0x47, 0x55, 0x41,
	0x47, 0x45, 0x5f, 0x52, 0x45, 0x47, 0x4f, 0x10, 0x01, 0x42, 0x2a, 0x5a, 0x28, 0x63, 0x6c, 0x6f,
	0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x69, 0x6f, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69,
	0x74, 0x6f, 0x72, 0x2f, 0x76, 0x32, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x73, 0x73, 0x65, 0x73,
	0x73, 0x6d, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_api_assessment_metric_proto_rawDescOnce sync.Once
	file_api_assessment_metric_proto_rawDescData []byte
)

func file_api_assessment_metric_proto_rawDescGZIP() []byte {
	file_api_assessment_metric_proto_rawDescOnce.Do(func() {
		file_api_assessment_metric_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_assessment_metric_proto_rawDesc), len(file_api_assessment_metric_proto_rawDesc)))
	})
	return file_api_assessment_metric_proto_rawDescData
}

var file_api_assessment_metric_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_api_assessment_metric_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_api_assessment_metric_proto_goTypes = []any{
	(Metric_Scale)(0),                  // 0: clouditor.assessment.v1.Metric.Scale
	(MetricImplementation_Language)(0), // 1: clouditor.assessment.v1.MetricImplementation.Language
	(*Metric)(nil),                     // 2: clouditor.assessment.v1.Metric
	(*Range)(nil),                      // 3: clouditor.assessment.v1.Range
	(*MinMax)(nil),                     // 4: clouditor.assessment.v1.MinMax
	(*AllowedValues)(nil),              // 5: clouditor.assessment.v1.AllowedValues
	(*Order)(nil),                      // 6: clouditor.assessment.v1.Order
	(*MetricConfiguration)(nil),        // 7: clouditor.assessment.v1.MetricConfiguration
	(*MetricImplementation)(nil),       // 8: clouditor.assessment.v1.MetricImplementation
	(*durationpb.Duration)(nil),        // 9: google.protobuf.Duration
	(*timestamppb.Timestamp)(nil),      // 10: google.protobuf.Timestamp
	(*structpb.Value)(nil),             // 11: google.protobuf.Value
}
var file_api_assessment_metric_proto_depIdxs = []int32{
	0,  // 0: clouditor.assessment.v1.Metric.scale:type_name -> clouditor.assessment.v1.Metric.Scale
	3,  // 1: clouditor.assessment.v1.Metric.range:type_name -> clouditor.assessment.v1.Range
	9,  // 2: clouditor.assessment.v1.Metric.interval:type_name -> google.protobuf.Duration
	8,  // 3: clouditor.assessment.v1.Metric.implementation:type_name -> clouditor.assessment.v1.MetricImplementation
	10, // 4: clouditor.assessment.v1.Metric.deprecated_since:type_name -> google.protobuf.Timestamp
	5,  // 5: clouditor.assessment.v1.Range.allowed_values:type_name -> clouditor.assessment.v1.AllowedValues
	6,  // 6: clouditor.assessment.v1.Range.order:type_name -> clouditor.assessment.v1.Order
	4,  // 7: clouditor.assessment.v1.Range.min_max:type_name -> clouditor.assessment.v1.MinMax
	11, // 8: clouditor.assessment.v1.AllowedValues.values:type_name -> google.protobuf.Value
	11, // 9: clouditor.assessment.v1.Order.values:type_name -> google.protobuf.Value
	11, // 10: clouditor.assessment.v1.MetricConfiguration.target_value:type_name -> google.protobuf.Value
	10, // 11: clouditor.assessment.v1.MetricConfiguration.updated_at:type_name -> google.protobuf.Timestamp
	1,  // 12: clouditor.assessment.v1.MetricImplementation.lang:type_name -> clouditor.assessment.v1.MetricImplementation.Language
	10, // 13: clouditor.assessment.v1.MetricImplementation.updated_at:type_name -> google.protobuf.Timestamp
	14, // [14:14] is the sub-list for method output_type
	14, // [14:14] is the sub-list for method input_type
	14, // [14:14] is the sub-list for extension type_name
	14, // [14:14] is the sub-list for extension extendee
	0,  // [0:14] is the sub-list for field type_name
}

func init() { file_api_assessment_metric_proto_init() }
func file_api_assessment_metric_proto_init() {
	if File_api_assessment_metric_proto != nil {
		return
	}
	file_api_assessment_metric_proto_msgTypes[0].OneofWrappers = []any{}
	file_api_assessment_metric_proto_msgTypes[1].OneofWrappers = []any{
		(*Range_AllowedValues)(nil),
		(*Range_Order)(nil),
		(*Range_MinMax)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_assessment_metric_proto_rawDesc), len(file_api_assessment_metric_proto_rawDesc)),
			NumEnums:      2,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_assessment_metric_proto_goTypes,
		DependencyIndexes: file_api_assessment_metric_proto_depIdxs,
		EnumInfos:         file_api_assessment_metric_proto_enumTypes,
		MessageInfos:      file_api_assessment_metric_proto_msgTypes,
	}.Build()
	File_api_assessment_metric_proto = out.File
	file_api_assessment_metric_proto_goTypes = nil
	file_api_assessment_metric_proto_depIdxs = nil
}
