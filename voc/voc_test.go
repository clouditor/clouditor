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

package voc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
)

type test struct {
	A int    `json:"a"`
	B string `json:"b"`
	C bool   `json:"c"`
}

func TestAuthenticityInterface(t *testing.T) {
	tests := []struct {
		isAuthenticity IsAuthenticity
		typ            string
	}{
		{
			isAuthenticity: &NoAuthentication{},
			typ:            "NoAuthentication",
		},
		{
			isAuthenticity: &SingleSignOn{},
			typ:            "SingleSignOn",
		},
		{
			isAuthenticity: &OTPBasedAuthentication{},
			typ:            "OTPBasedAuthentication",
		},
		{
			isAuthenticity: &PasswordBasedAuthentication{},
			typ:            "PasswordBasedAuthentication",
		},
		{
			isAuthenticity: &TokenBasedAuthentication{},
			typ:            "TokenBasedAuthentication",
		},
	}

	for _, tt := range tests {
		t.Run(tt.typ, func(t *testing.T) {
			var isa = tt.isAuthenticity
			assert.Equal(t, tt.typ, isa.Type())
		})
	}
}

func TestAuthorizationInterface(t *testing.T) {
	tests := []struct {
		isAuthorization IsAuthorization
		typ             string
	}{
		{
			isAuthorization: &AccessRestriction{},
			typ:             "AccessRestriction",
		},
		{
			isAuthorization: &ABAC{},
			typ:             "ABAC",
		},
		{
			isAuthorization: &RBAC{},
			typ:             "RBAC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.typ, func(t *testing.T) {
			var isa = tt.isAuthorization
			assert.Equal(t, tt.typ, isa.Type())
		})
	}
}

func TestToStringInterface(t *testing.T) {
	var tmp = []interface{}{}

	s, err := ToStringInterface(tmp)
	assert.NoError(t, err)
	assert.Equal(t, "{}", s)

	tmp = append(tmp, test{
		A: 200,
		B: "TestStruct",
		C: false,
	})
	s, err = ToStringInterface(tmp)
	assert.NoError(t, err)
	assert.Equal(t, "{\"voc.test\":[{\"a\":200,\"b\":\"TestStruct\",\"c\":false}]}", s)
}

func TestToStruct(t *testing.T) {
	type args struct {
		r IsCloudResource
	}
	tests := []struct {
		name    string
		args    args
		wantS   *structpb.Value
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "Empty input",
			args:    args{},
			wantS:   nil,
			wantErr: assert.NoError,
		},
		{
			name: "Happy path",
			args: args{
				r: &Resource{
					ID:   "my-resource-id",
					Name: "my-resource-name",
					Type: ObjectStorageType,
				},
			},
			wantS: &structpb.Value{
				Kind: &structpb.Value_StructValue{
					StructValue: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"creationTime": structpb.NewNumberValue(0),
							"geoLocation": structpb.NewStructValue(&structpb.Struct{
								Fields: map[string]*structpb.Value{
									"region": structpb.NewStringValue(""),
								}}),
							"id":        structpb.NewStringValue("my-resource-id"),
							"labels":    structpb.NewNullValue(),
							"name":      structpb.NewStringValue("my-resource-name"),
							"parent":    structpb.NewStringValue(""),
							"raw":       structpb.NewStringValue(""),
							"serviceId": structpb.NewStringValue(""),
							"type": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{
								structpb.NewStringValue("ObjectStorage"),
								structpb.NewStringValue("Storage"),
								structpb.NewStringValue("Resource"),
							}}),
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotS, err := ToStruct(tt.args.r)

			tt.wantErr(t, err)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantS.String(), gotS.String())
			}
		})
	}
}
