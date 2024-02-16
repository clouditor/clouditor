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

package gorm

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm/schema"

	"clouditor.io/clouditor/v2/internal/util"
	"clouditor.io/clouditor/v2/persistence"
)

// TimestampSerializer is a GORM serializer that allows the serialization and deserialization of the
// google.protobuf.Timestamp protobuf message type.
type TimestampSerializer struct{}

// Value implements https://pkg.go.dev/gorm.io/gorm/schema#SerializerValuerInterface to indicate
// how this struct will be saved into an SQL database field.
func (TimestampSerializer) Value(_ context.Context, _ *schema.Field, _ reflect.Value, fieldValue interface{}) (interface{}, error) {
	var (
		t  *timestamppb.Timestamp
		ok bool
	)

	if util.IsNil(fieldValue) {
		return nil, nil
	}

	if t, ok = fieldValue.(*timestamppb.Timestamp); !ok {
		return nil, persistence.ErrUnsupportedType
	}

	return t.AsTime(), nil
}

// Scan implements https://pkg.go.dev/gorm.io/gorm/schema#SerializerInterface to indicate how
// this struct can be loaded from an SQL database field.
func (TimestampSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
	var t *timestamppb.Timestamp

	if dbValue != nil {
		switch v := dbValue.(type) {
		case time.Time:
			t = timestamppb.New(v)
		default:
			return persistence.ErrUnsupportedType
		}

		field.ReflectValueOf(ctx, dst).Set(reflect.ValueOf(t))
	}

	return
}

// AnySerializer is a GORM serializer that allows the serialization and deserialization of the
// google.protobuf.Any protobuf message type using a JSONB field.
type AnySerializer struct{}

// Value implements https://pkg.go.dev/gorm.io/gorm/schema#SerializerValuerInterface to indicate
// how this struct will be saved into an SQL database field.
func (AnySerializer) Value(_ context.Context, _ *schema.Field, _ reflect.Value, fieldValue interface{}) (interface{}, error) {
	var (
		a  *anypb.Any
		ok bool
	)

	if util.IsNil(fieldValue) {
		return nil, nil
	}

	if a, ok = fieldValue.(*anypb.Any); !ok {
		return nil, persistence.ErrUnsupportedType
	}

	return protojson.Marshal(a)
}

// Scan implements https://pkg.go.dev/gorm.io/gorm/schema#SerializerInterface to indicate how
// this struct can be loaded from an SQL database field.
func (AnySerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
	var (
		a anypb.Any
	)

	if dbValue != nil {
		var bytes []byte
		switch v := dbValue.(type) {
		case []byte:
			bytes = v
		case string:
			bytes = []byte(v)
		default:
			return persistence.ErrUnsupportedType
		}

		err = protojson.Unmarshal(bytes, &a)
		if err != nil {
			return fmt.Errorf("could not unmarshal JSONB value into protobuf message: %w", err)
		}
	}

	field.ReflectValueOf(ctx, dst).Set(reflect.ValueOf(&a))
	return
}

// ValueSerializer is a GORM serializer that allows the serialization and deserialization of the
// google.protobuf.Value protobuf message type.
type ValueSerializer struct{}

// Value implements https://pkg.go.dev/gorm.io/gorm/schema#SerializerValuerInterface to indicate
// how this struct will be saved into an SQL database field.
func (ValueSerializer) Value(_ context.Context, _ *schema.Field, _ reflect.Value, fieldValue interface{}) (interface{}, error) {
	var (
		v  *structpb.Value
		ok bool
	)

	if util.IsNil(fieldValue) {
		return nil, nil
	}

	if v, ok = fieldValue.(*structpb.Value); !ok {
		return nil, persistence.ErrUnsupportedType
	}

	return v.MarshalJSON()
}

// Scan implements https://pkg.go.dev/gorm.io/gorm/schema#SerializerInterface to indicate how
// this struct can be loaded from an SQL database field.
func (ValueSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
	v := new(structpb.Value)

	if dbValue != nil {
		switch d := dbValue.(type) {
		case []byte:
			err = v.UnmarshalJSON(d)
			if err != nil {
				return err
			}
		default:
			return persistence.ErrUnsupportedType
		}

		field.ReflectValueOf(ctx, dst).Set(reflect.ValueOf(v))
	}

	return
}
