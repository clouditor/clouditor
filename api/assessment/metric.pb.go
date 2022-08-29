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
// 	protoc-gen-go v1.28.1
// 	protoc        v3.20.0
// source: api/assessment/metric.proto

package assessment

import (
	_ "github.com/srikrsna/protoc-gen-gotag/tagger"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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
	MetricImplementation_REGO                 MetricImplementation_Language = 1
)

// Enum value maps for MetricImplementation_Language.
var (
	MetricImplementation_Language_name = map[int32]string{
		0: "LANGUAGE_UNSPECIFIED",
		1: "REGO",
	}
	MetricImplementation_Language_value = map[string]int32{
		"LANGUAGE_UNSPECIFIED": 0,
		"REGO":                 1,
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
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Required. The unique identifier of the metric.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Required. The human readable name of the metric.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// The description of the metric
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	// The reference to control catalog category or domain
	Category string `protobuf:"bytes,4,opt,name=category,proto3" json:"category,omitempty"`
	// The scale of this metric, e.g. categories, ranked data or metric values.
	Scale Metric_Scale `protobuf:"varint,5,opt,name=scale,proto3,enum=clouditor.Metric_Scale" json:"scale,omitempty"`
	// The range of this metric. Depending on the scale.
	Range *Range `protobuf:"bytes,6,opt,name=range,proto3" json:"range,omitempty"`
	// The interval in seconds the evidences must be collected for the respective metric.
	// For now, we are not able to use google.protobuf.Duration because it is converted to a custom object in OpenAPI (https://github.com/google/gnostic/issues/351)
	Interval int64 `protobuf:"varint,7,opt,name=interval,proto3" json:"interval,omitempty"`
}

func (x *Metric) Reset() {
	*x = Metric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_assessment_metric_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric) ProtoMessage() {}

func (x *Metric) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_metric_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
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

func (x *Metric) GetInterval() int64 {
	if x != nil {
		return x.Interval
	}
	return 0
}

// A range resource representing the range of values
type Range struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Required.
	//
	// Types that are assignable to Range:
	//	*Range_AllowedValues
	//	*Range_Order
	//	*Range_MinMax
	Range isRange_Range `protobuf_oneof:"range"`
}

func (x *Range) Reset() {
	*x = Range{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_assessment_metric_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Range) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Range) ProtoMessage() {}

func (x *Range) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_metric_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
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

func (m *Range) GetRange() isRange_Range {
	if m != nil {
		return m.Range
	}
	return nil
}

func (x *Range) GetAllowedValues() *AllowedValues {
	if x, ok := x.GetRange().(*Range_AllowedValues); ok {
		return x.AllowedValues
	}
	return nil
}

func (x *Range) GetOrder() *Order {
	if x, ok := x.GetRange().(*Range_Order); ok {
		return x.Order
	}
	return nil
}

func (x *Range) GetMinMax() *MinMax {
	if x, ok := x.GetRange().(*Range_MinMax); ok {
		return x.MinMax
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
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Required.
	Min int64 `protobuf:"varint,1,opt,name=min,proto3" json:"min,omitempty"`
	// Required.
	Max int64 `protobuf:"varint,2,opt,name=max,proto3" json:"max,omitempty"`
}

func (x *MinMax) Reset() {
	*x = MinMax{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_assessment_metric_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MinMax) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MinMax) ProtoMessage() {}

func (x *MinMax) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_metric_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
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
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Values []*structpb.Value `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty"`
}

func (x *AllowedValues) Reset() {
	*x = AllowedValues{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_assessment_metric_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AllowedValues) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllowedValues) ProtoMessage() {}

func (x *AllowedValues) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_metric_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
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
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Values []*structpb.Value `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty"`
}

func (x *Order) Reset() {
	*x = Order{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_assessment_metric_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Order) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Order) ProtoMessage() {}

func (x *Order) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_metric_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
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
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The operator to compare the metric, such as == or >
	Operator string `protobuf:"bytes,1,opt,name=operator,proto3" json:"operator,omitempty"`
	// The target value
	TargetValue *structpb.Value `protobuf:"bytes,2,opt,name=target_value,json=targetValue,proto3" json:"target_value,omitempty"`
	// Whether this configuration is a default configuration
	IsDefault bool `protobuf:"varint,3,opt,name=is_default,json=isDefault,proto3" json:"is_default,omitempty"`
	// The last time of update
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty" gorm:"serializer:timestamppb;type:time"`
}

func (x *MetricConfiguration) Reset() {
	*x = MetricConfiguration{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_assessment_metric_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MetricConfiguration) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetricConfiguration) ProtoMessage() {}

func (x *MetricConfiguration) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_metric_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
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

// MetricImplementation defines the implementation of an individual metric.
type MetricImplementation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The metric which is implemented
	MetricId string `protobuf:"bytes,1,opt,name=metric_id,json=metricId,proto3" json:"metric_id,omitempty"`
	// The language this metric is implemented in
	Lang MetricImplementation_Language `protobuf:"varint,2,opt,name=lang,proto3,enum=clouditor.MetricImplementation_Language" json:"lang,omitempty"`
	// The actual implementation
	Code string `protobuf:"bytes,3,opt,name=code,proto3" json:"code,omitempty"`
	// The last time of update
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty" gorm:"serializer:timestamppb;type:time"`
}

func (x *MetricImplementation) Reset() {
	*x = MetricImplementation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_assessment_metric_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MetricImplementation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetricImplementation) ProtoMessage() {}

func (x *MetricImplementation) ProtoReflect() protoreflect.Message {
	mi := &file_api_assessment_metric_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
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

var file_api_assessment_metric_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74,
	0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x63,
	0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13, 0x74, 0x61, 0x67, 0x67, 0x65, 0x72, 0x2f,
	0x74, 0x61, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa3, 0x02, 0x0a,
	0x06, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a,
	0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x2d, 0x0a, 0x05, 0x73, 0x63, 0x61,
	0x6c, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64,
	0x69, 0x74, 0x6f, 0x72, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x53, 0x63, 0x61, 0x6c,
	0x65, 0x52, 0x05, 0x73, 0x63, 0x61, 0x6c, 0x65, 0x12, 0x26, 0x0a, 0x05, 0x72, 0x61, 0x6e, 0x67,
	0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69,
	0x74, 0x6f, 0x72, 0x2e, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x05, 0x72, 0x61, 0x6e, 0x67, 0x65,
	0x12, 0x1a, 0x0a, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x22, 0x44, 0x0a, 0x05,
	0x53, 0x63, 0x61, 0x6c, 0x65, 0x12, 0x15, 0x0a, 0x11, 0x53, 0x43, 0x41, 0x4c, 0x45, 0x5f, 0x55,
	0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07,
	0x4e, 0x4f, 0x4d, 0x49, 0x4e, 0x41, 0x4c, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x4f, 0x52, 0x44,
	0x49, 0x4e, 0x41, 0x4c, 0x10, 0x02, 0x12, 0x0a, 0x0a, 0x06, 0x4d, 0x45, 0x54, 0x52, 0x49, 0x43,
	0x10, 0x03, 0x22, 0xab, 0x01, 0x0a, 0x05, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x41, 0x0a, 0x0e,
	0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72,
	0x2e, 0x41, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x48, 0x00,
	0x52, 0x0d, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x12,
	0x28, 0x0a, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10,
	0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72,
	0x48, 0x00, 0x52, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x2c, 0x0a, 0x07, 0x6d, 0x69, 0x6e,
	0x5f, 0x6d, 0x61, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x63, 0x6c, 0x6f,
	0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x4d, 0x69, 0x6e, 0x4d, 0x61, 0x78, 0x48, 0x00, 0x52,
	0x06, 0x6d, 0x69, 0x6e, 0x4d, 0x61, 0x78, 0x42, 0x07, 0x0a, 0x05, 0x72, 0x61, 0x6e, 0x67, 0x65,
	0x22, 0x2c, 0x0a, 0x06, 0x4d, 0x69, 0x6e, 0x4d, 0x61, 0x78, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x69,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x6d, 0x69, 0x6e, 0x12, 0x10, 0x0a, 0x03,
	0x6d, 0x61, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x6d, 0x61, 0x78, 0x22, 0x3f,
	0x0a, 0x0d, 0x41, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x12,
	0x2e, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x22,
	0x37, 0x0a, 0x05, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x2e, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x22, 0xf4, 0x01, 0x0a, 0x13, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x1a, 0x0a, 0x08, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x39, 0x0a, 0x0c,
	0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0b, 0x74, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x73, 0x5f, 0x64, 0x65,
	0x66, 0x61, 0x75, 0x6c, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x69, 0x73, 0x44,
	0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x12, 0x67, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x64, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x2c, 0x9a, 0x84, 0x9e, 0x03, 0x27, 0x67, 0x6f, 0x72,
	0x6d, 0x3a, 0x22, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x72, 0x3a, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x70, 0x62, 0x3b, 0x74, 0x79, 0x70, 0x65, 0x3a, 0x74,
	0x69, 0x6d, 0x65, 0x22, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22,
	0x9e, 0x02, 0x0a, 0x14, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x49, 0x6d, 0x70, 0x6c, 0x65, 0x6d,
	0x65, 0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x6d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x49, 0x64, 0x12, 0x3c, 0x0a, 0x04, 0x6c, 0x61, 0x6e, 0x67, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x28, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x49, 0x6d, 0x70, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x52, 0x04, 0x6c,
	0x61, 0x6e, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x67, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x2c, 0x9a, 0x84, 0x9e, 0x03, 0x27, 0x67, 0x6f,
	0x72, 0x6d, 0x3a, 0x22, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x72, 0x3a, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x70, 0x62, 0x3b, 0x74, 0x79, 0x70, 0x65, 0x3a,
	0x74, 0x69, 0x6d, 0x65, 0x22, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x22, 0x2e, 0x0a, 0x08, 0x4c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x14,
	0x4c, 0x41, 0x4e, 0x47, 0x55, 0x41, 0x47, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49,
	0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x52, 0x45, 0x47, 0x4f, 0x10, 0x01,
	0x42, 0x32, 0x5a, 0x30, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x69, 0x6f,
	0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61,
	0x73, 0x73, 0x65, 0x73, 0x73, 0x6d, 0x65, 0x6e, 0x74, 0x3b, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73,
	0x6d, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_assessment_metric_proto_rawDescOnce sync.Once
	file_api_assessment_metric_proto_rawDescData = file_api_assessment_metric_proto_rawDesc
)

func file_api_assessment_metric_proto_rawDescGZIP() []byte {
	file_api_assessment_metric_proto_rawDescOnce.Do(func() {
		file_api_assessment_metric_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_assessment_metric_proto_rawDescData)
	})
	return file_api_assessment_metric_proto_rawDescData
}

var file_api_assessment_metric_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_api_assessment_metric_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_api_assessment_metric_proto_goTypes = []interface{}{
	(Metric_Scale)(0),                  // 0: clouditor.Metric.Scale
	(MetricImplementation_Language)(0), // 1: clouditor.MetricImplementation.Language
	(*Metric)(nil),                     // 2: clouditor.Metric
	(*Range)(nil),                      // 3: clouditor.Range
	(*MinMax)(nil),                     // 4: clouditor.MinMax
	(*AllowedValues)(nil),              // 5: clouditor.AllowedValues
	(*Order)(nil),                      // 6: clouditor.Order
	(*MetricConfiguration)(nil),        // 7: clouditor.MetricConfiguration
	(*MetricImplementation)(nil),       // 8: clouditor.MetricImplementation
	(*structpb.Value)(nil),             // 9: google.protobuf.Value
	(*timestamppb.Timestamp)(nil),      // 10: google.protobuf.Timestamp
}
var file_api_assessment_metric_proto_depIdxs = []int32{
	0,  // 0: clouditor.Metric.scale:type_name -> clouditor.Metric.Scale
	3,  // 1: clouditor.Metric.range:type_name -> clouditor.Range
	5,  // 2: clouditor.Range.allowed_values:type_name -> clouditor.AllowedValues
	6,  // 3: clouditor.Range.order:type_name -> clouditor.Order
	4,  // 4: clouditor.Range.min_max:type_name -> clouditor.MinMax
	9,  // 5: clouditor.AllowedValues.values:type_name -> google.protobuf.Value
	9,  // 6: clouditor.Order.values:type_name -> google.protobuf.Value
	9,  // 7: clouditor.MetricConfiguration.target_value:type_name -> google.protobuf.Value
	10, // 8: clouditor.MetricConfiguration.updated_at:type_name -> google.protobuf.Timestamp
	1,  // 9: clouditor.MetricImplementation.lang:type_name -> clouditor.MetricImplementation.Language
	10, // 10: clouditor.MetricImplementation.updated_at:type_name -> google.protobuf.Timestamp
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_api_assessment_metric_proto_init() }
func file_api_assessment_metric_proto_init() {
	if File_api_assessment_metric_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_assessment_metric_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metric); i {
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
		file_api_assessment_metric_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Range); i {
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
		file_api_assessment_metric_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MinMax); i {
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
		file_api_assessment_metric_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AllowedValues); i {
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
		file_api_assessment_metric_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Order); i {
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
		file_api_assessment_metric_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MetricConfiguration); i {
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
		file_api_assessment_metric_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MetricImplementation); i {
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
	file_api_assessment_metric_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*Range_AllowedValues)(nil),
		(*Range_Order)(nil),
		(*Range_MinMax)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_assessment_metric_proto_rawDesc,
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
	file_api_assessment_metric_proto_rawDesc = nil
	file_api_assessment_metric_proto_goTypes = nil
	file_api_assessment_metric_proto_depIdxs = nil
}
