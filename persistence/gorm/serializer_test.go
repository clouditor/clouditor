// Copyright 2016-2022 Fraunhofer AISEC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//	         $$\                           $$\ $$\   $$\
//	         $$ |                          $$ |\__|  $$ |
//	$$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
//
// $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
// $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
// $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
// \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
//
//	\_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
//
// This file is part of Clouditor Community Edition.
package gorm

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/persistence"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm/schema"
)

func TestTimestampSerializer_Value(t *testing.T) {
	type args struct {
		ctx        context.Context
		field      *schema.Field
		dst        reflect.Value
		fieldValue interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "ok field",
			args: args{
				field:      &schema.Field{Name: "timestamp"},
				dst:        reflect.Value{},
				fieldValue: timestamppb.New(time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)),
			},
			want:    time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC),
			wantErr: assert.NoError,
		},
		{
			name: "nil field",
			args: args{
				field:      &schema.Field{Name: "timestamp"},
				dst:        reflect.Value{},
				fieldValue: nil,
			},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "field wrong type",
			args: args{
				field:      &schema.Field{Name: "timestamp"},
				dst:        reflect.Value{},
				fieldValue: "string",
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, persistence.ErrUnsupportedType)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := TimestampSerializer{}

			got, err := tr.Value(tt.args.ctx, tt.args.field, tt.args.dst, tt.args.fieldValue)
			tt.wantErr(t, err, tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTimestampSerializer_Scan(t *testing.T) {
	type args struct {
		ctx     context.Context
		field   *schema.Field
		dst     reflect.Value
		dbValue interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "db wrong type",
			args: args{
				field:   &schema.Field{},
				dbValue: "string",
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, persistence.ErrUnsupportedType)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := TimestampSerializer{}
			err := tr.Scan(tt.args.ctx, tt.args.field, tt.args.dst, tt.args.dbValue)

			tt.wantErr(t, err, tt.args)
		})
	}
}

func TestAnySerializer_Value(t *testing.T) {
	type args struct {
		ctx        context.Context
		field      *schema.Field
		dst        reflect.Value
		fieldValue interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    assert.Want[any]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "ok field",
			args: args{
				field: &schema.Field{Name: "config"},
				dst:   reflect.Value{},
				fieldValue: func() *anypb.Any {
					a, _ := anypb.New(&orchestrator.CloudService{
						Id: "my-service",
					})
					return a
				}(),
			},
			want: func(t *testing.T, got any) bool {
				// output of protojson is randomized (see
				// https://github.com/protocolbuffers/protobuf-go/commit/582ab3de426ef0758666e018b422dd20390f7f26),
				// so we need to unmarshal it to compare the contents in a
				// stable way
				b := assert.Is[[]byte](t, got)
				if !assert.NotNil(t, b) {
					return false
				}

				var m map[string]interface{}
				err := json.Unmarshal(b, &m)
				assert.NoError(t, err)

				return assert.Equal(t, m, map[string]interface{}{
					"@type": "type.googleapis.com/clouditor.orchestrator.v1.CloudService",
					"id":    "my-service",
				})
			},
			wantErr: assert.NoError,
		},
		{
			name: "nil field",
			args: args{
				field:      &schema.Field{Name: "config"},
				dst:        reflect.Value{},
				fieldValue: nil,
			},
			want:    assert.Nil[any],
			wantErr: assert.NoError,
		},
		{
			name: "field wrong type",
			args: args{
				field:      &schema.Field{Name: "config"},
				dst:        reflect.Value{},
				fieldValue: "string",
			},
			want: assert.Nil[any],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, persistence.ErrUnsupportedType)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := AnySerializer{}

			got, err := tr.Value(tt.args.ctx, tt.args.field, tt.args.dst, tt.args.fieldValue)
			tt.wantErr(t, err, tt.args)
			tt.want(t, got)
		})
	}
}

func TestAnySerializer_Scan(t *testing.T) {
	type args struct {
		ctx     context.Context
		field   *schema.Field
		dst     reflect.Value
		dbValue interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "db wrong type",
			args: args{
				field:   &schema.Field{},
				dbValue: "string",
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not unmarshal JSONB value into protobuf message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := AnySerializer{}
			err := tr.Scan(tt.args.ctx, tt.args.field, tt.args.dst, tt.args.dbValue)

			tt.wantErr(t, err, tt.args)

		})
	}
}
