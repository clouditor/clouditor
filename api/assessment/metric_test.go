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

package assessment

import (
	"database/sql/driver"
	"testing"

	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/persistence"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestMetricConfiguration_Validate(t *testing.T) {
	type fields struct {
		MetricConfiguration *MetricConfiguration
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "MetricConfiguration Operator is empty",
			fields: fields{
				MetricConfiguration: &MetricConfiguration{
					TargetValue: testdata.MockMetricConfigurationTargetValueString,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "MetricId: value length must be at least 1 runes")
			},
		},
		{
			name: "MetricConfiguration TargetValue is empty",
			fields: fields{
				MetricConfiguration: &MetricConfiguration{
					Operator: "==",
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "TargetValue: value is required")
			},
		},
		{
			name: "Successful Validation",
			fields: fields{
				MetricConfiguration: &MetricConfiguration{
					TargetValue:    testdata.MockMetricConfigurationTargetValueString,
					Operator:       "==",
					MetricId:       testdata.MockMetricID1,
					CloudServiceId: testdata.MockCloudServiceID1,
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields.MetricConfiguration
			tt.wantErr(t, c.Validate())
		})
	}
}

func TestRange_GormDataType(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Works correctly",
			want: "jsonb",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Range{}
			if got := r.GormDataType(); got != tt.want {
				t.Errorf("Range.GormDataType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRange_Value(t *testing.T) {
	type fields struct {
		Range *Range
	}
	tests := []struct {
		name    string
		fields  fields
		want    driver.Value
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Error in Range",
			fields: fields{
				Range: &Range{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not marshal JSON")
			},
		},
		{
			name: "Range is empty",
			fields: fields{
				Range: nil,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.fields.Range
			got, err := r.Value()
			tt.wantErr(t, err)
			if err == nil {
				assert.Equal(t, got, tt.want)
			}
		})
	}
}

func TestRange_Scan(t *testing.T) {
	type fields struct {
		Range *Range
	}
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Value of wrong type",
			fields: fields{
				Range: &Range{},
			},
			args: args{
				value: "wrongType",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, persistence.ErrUnsupportedType.Error())
			},
		},
		{
			name: "Error at unmarshalling",
			fields: fields{
				Range: &Range{},
			},
			args: args{
				value: []byte{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not unmarshal JSON")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.fields.Range
			err := r.Scan(tt.args.value)
			tt.wantErr(t, err)
		})
	}
}

func TestRange_MarshalJSON(t *testing.T) {
	type fields struct {
		Range *Range
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Unknown range type",
			fields: fields{
				Range: &Range{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, persistence.ErrUnsupportedType.Error())
			},
		},
		{
			name: "Correct range type Range_AllowedValues",
			fields: fields{
				Range: &Range{
					Range: &Range_AllowedValues{},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: "Correct range type Range_MinMax",
			fields: fields{
				Range: &Range{
					Range: &Range_MinMax{},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: "Correct range type Range_Order",
			fields: fields{
				Range: &Range{
					Range: &Range_Order{},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.fields.Range
			_, err := r.MarshalJSON()
			tt.wantErr(t, err)
		})
	}
}

func TestRange_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Range *Range
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Empty input",
			fields: fields{
				Range: &Range{},
			},
			args: args{
				b: []byte{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "unexpected end of JSON input")
			},
		},
		{
			name: "Invalid input",
			fields: fields{
				Range: &Range{
					Range: &Range_AllowedValues{},
				},
			},
			args: args{
				b: []byte("Error"),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid character")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.fields.Range
			tt.wantErr(t, r.UnmarshalJSON(tt.args.b))
		})
	}
}

func TestCheckCloudServiceID(t *testing.T) {
	type args struct {
		serviceID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Missing serviceID",
			args: args{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrCloudServiceIDIsMissing.Error())
			},
		},
		{
			name: "Invalid serviceID",
			args: args{
				serviceID: "00000000-0000-0000-000000000000",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrCloudServiceIDIsInvalid.Error())
			},
		},
		{
			name: "Happy path",
			args: args{
				serviceID: testdata.MockCloudServiceID1,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckCloudServiceID(tt.args.serviceID)
			tt.wantErr(t, err)
		})
	}
}

func TestMetricConfiguration_Hash(t *testing.T) {
	type fields struct {
		sizeCache      protoimpl.SizeCache
		unknownFields  protoimpl.UnknownFields
		Operator       string
		TargetValue    *structpb.Value
		IsDefault      bool
		UpdatedAt      *timestamppb.Timestamp
		MetricId       string
		CloudServiceId string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy path",
			fields: fields{
				Operator:    "<",
				TargetValue: &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: 5}},
			},
			want: "PC1udW1iZXJfdmFsdWU6NQ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &MetricConfiguration{
				sizeCache:      tt.fields.sizeCache,
				unknownFields:  tt.fields.unknownFields,
				Operator:       tt.fields.Operator,
				TargetValue:    tt.fields.TargetValue,
				IsDefault:      tt.fields.IsDefault,
				UpdatedAt:      tt.fields.UpdatedAt,
				MetricId:       tt.fields.MetricId,
				CloudServiceId: tt.fields.CloudServiceId,
			}
			if got := x.Hash(); got != tt.want {
				t.Errorf("MetricConfiguration.Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}
