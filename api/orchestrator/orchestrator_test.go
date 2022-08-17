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

package orchestrator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/assessment"
)

func TestListCertificatesRequest_Validate(t *testing.T) {
	type fields struct {
		PageSize  int32
		PageToken string
		OrderBy   string
		Asc       bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Valid with id",
			fields: fields{
				OrderBy: "Id",
				Asc:     true,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(t, err)
			},
		},
		{
			name: "Valid with empty string",
			fields: fields{
				OrderBy: "",
				Asc:     true,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(t, err)
			},
		},
		{
			name: "Invalid",
			fields: fields{
				OrderBy: "notAField",
				Asc:     true,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Contains(t, err.Error(), api.ErrInvalidColumnName.Error())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ListCertificatesRequest{
				PageSize:  tt.fields.PageSize,
				PageToken: tt.fields.PageToken,
				OrderBy:   tt.fields.OrderBy,
				Asc:       tt.fields.Asc,
			}
			err := api.ValidateListRequest[*Certificate](req)
			tt.wantErr(t, err)
		})
	}
}

func TestListCloudServicesRequest_Validate(t *testing.T) {
	type fields struct {
		PageSize  int32
		PageToken string
		OrderBy   string
		Asc       bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Valid with id",
			fields: fields{
				OrderBy: "Id",
				Asc:     true,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(t, err)
			},
		},
		{
			name: "Valid with empty string",
			fields: fields{
				OrderBy: "",
				Asc:     true,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(t, err)
			},
		},
		{
			name: "Invalid",
			fields: fields{
				OrderBy: "notAField",
				Asc:     true,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Contains(t, err.Error(), api.ErrInvalidColumnName.Error())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ListCloudServicesRequest{
				PageSize:  tt.fields.PageSize,
				PageToken: tt.fields.PageToken,
				OrderBy:   tt.fields.OrderBy,
				Asc:       tt.fields.Asc,
			}
			tt.wantErr(t, api.ValidateListRequest[*CloudService](req), "Validate()")
		})
	}
}

func TestListMetricsRequest_Validate(t *testing.T) {
	type fields struct {
		PageSize  int32
		PageToken string
		OrderBy   string
		Asc       bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Valid with id",
			fields: fields{
				OrderBy: "Category",
				Asc:     true,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(t, err)
			},
		},
		{
			name: "Valid with empty string",
			fields: fields{
				OrderBy: "",
				Asc:     true,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(t, err)
			},
		},
		{
			name: "Invalid",
			fields: fields{
				OrderBy: "notAField",
				Asc:     true,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Contains(t, err.Error(), api.ErrInvalidColumnName.Error())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ListMetricsRequest{
				PageSize:  tt.fields.PageSize,
				PageToken: tt.fields.PageToken,
				OrderBy:   tt.fields.OrderBy,
				Asc:       tt.fields.Asc,
			}
			tt.wantErr(t, api.ValidateListRequest[*assessment.Metric](req), "Validate()")
		})
	}
}

func TestListControlsRequest_Validate(t *testing.T) {
	type fields struct {
		PageSize  int32
		PageToken string
		OrderBy   string
		Asc       bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Valid with id",
			fields: fields{
				OrderBy: "Name",
				Asc:     true,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(t, err)
			},
		},
		{
			name: "Valid with empty string",
			fields: fields{
				OrderBy: "",
				Asc:     true,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(t, err)
			},
		},
		{
			name: "Invalid column",
			fields: fields{
				OrderBy: "notAField",
				Asc:     true,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Contains(t, err.Error(), api.ErrInvalidColumnName.Error())
			},
		},
		{
			name: "Req is nil",
			fields: fields{
				OrderBy: "notAField",
				Asc:     true,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Contains(t, err.Error(), api.ErrInvalidColumnName.Error())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ListControlsRequest{
				PageSize:  tt.fields.PageSize,
				PageToken: tt.fields.PageToken,
				OrderBy:   tt.fields.OrderBy,
				Asc:       tt.fields.Asc,
			}
			tt.wantErr(t, api.ValidateListRequest[*Control](req), "Validate()")
		})
	}
}
