package gorm

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/persistence"
	"github.com/stretchr/testify/assert"
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

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TimestampSerializer.Value() = %v, want %v", got, tt.want)
			}
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
		want    assert.ValueAssertionFunc
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
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				// output of protojson is randomized (see
				// https://github.com/protocolbuffers/protobuf-go/commit/582ab3de426ef0758666e018b422dd20390f7f26),
				// so we need to unmarshal it to compare the contents in a
				// stable way
				b, ok := i1.([]byte)
				if assert.True(t, ok) {
					return false
				}

				var m map[string]interface{}
				err := json.Unmarshal(b, &m)
				assert.NoError(t, err)

				return assert.Equal(t, m, map[string]interface{}{
					"@type": "type.googleapis.com/clouditor.CloudService",
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
			want:    assert.Empty,
			wantErr: assert.NoError,
		},
		{
			name: "field wrong type",
			args: args{
				field:      &schema.Field{Name: "config"},
				dst:        reflect.Value{},
				fieldValue: "string",
			},
			want: assert.Empty,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, persistence.ErrUnsupportedType)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := AnySerializer{}

			got, err := tr.Value(tt.args.ctx, tt.args.field, tt.args.dst, tt.args.fieldValue)

			// Validate the error via the ErrorAssertionFunc function
			tt.wantErr(t, err, tt.args)

			// Validate the response via the ValueAssertionFunc function
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
