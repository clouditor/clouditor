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

package util

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/voc"
)

func TestDeepCopyOfMap(t *testing.T) {
	type args struct {
		originalMap map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "Empty input",
			args: args{
				originalMap: map[string]interface{}{},
			},
			want: map[string]interface{}{},
		},
		{
			name: "map[string]string input",
			args: args{
				originalMap: map[string]interface{}{
					"testKey1": "testValue1",
					"testKey2": "testValue2",
				},
			},
			want: map[string]interface{}{
				"testKey1": "testValue1",
				"testKey2": "testValue2",
			},
		},
		{
			name: "map[string]interface{} input",
			args: args{
				originalMap: map[string]interface{}{
					"testKey1": map[string]interface{}{
						"testKey11": "testValue11",
					},
					"testKey2": map[string]interface{}{
						"testKey21": "testValue21",
					},
				},
			},
			want: map[string]interface{}{
				"testKey1": map[string]interface{}{
					"testKey11": "testValue11",
				},
				"testKey2": map[string]interface{}{
					"testKey21": "testValue21",
				},
			},
		},
		{
			name: "map[string]map[string]string input",
			args: args{
				originalMap: map[string]interface{}{
					"testKey1": map[string]string{
						"test2ndKey1": "test2ndValue1",
						"test1ndKey2": "test2ndValue2",
					},
					"testKey2": map[string]string{
						"test2ndKey3": "test2ndValue3",
						"test2ndKey4": "test2ndValue4",
					},
				},
			},
			want: map[string]interface{}{
				"testKey1": map[string]string{
					"test2ndKey1": "test2ndValue1",
					"test1ndKey2": "test2ndValue2",
				},
				"testKey2": map[string]string{
					"test2ndKey3": "test2ndValue3",
					"test2ndKey4": "test2ndValue4",
				},
			},
		},
		{
			name: "map[string]Resource input",
			args: args{
				originalMap: map[string]interface{}{
					"testKey1": &voc.Resource{
						ID:   testdata.MockResourceID,
						Name: testdata.MockResourceName,
						Type: []string{"Resource"},
					},
					"testKey2": &voc.Resource{
						ID:   "00000000-0000-0000-0000-000000000001",
						Name: "testResource2",
						Type: []string{"Resource"},
					},
				},
			},
			want: map[string]interface{}{
				"testKey1": &voc.Resource{
					ID:   testdata.MockResourceID,
					Name: testdata.MockResourceName,
					Type: []string{"Resource"},
				},
				"testKey2": &voc.Resource{
					ID:   "00000000-0000-0000-0000-000000000001",
					Name: "testResource2",
					Type: []string{"Resource"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, DeepCopyOfMap(tt.args.originalMap), "DeepCopyOfMap(%v)", tt.args.originalMap)
		})
	}
}

func TestDeepCopy(t *testing.T) {
	type args struct {
		original []interface{}
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: "Empty input",
			args: args{
				original: []interface{}{},
			},
			want: []interface{}(nil),
		},
		{
			name: "[]string input",
			args: args{
				original: []interface{}{
					"testValue1",
					"testValue2",
				},
			},
			want: []interface{}{
				"testValue1",
				"testValue2",
			},
		},
		{
			name: "[]Resource input",
			args: args{
				original: []interface{}{
					&voc.Resource{
						ID:   testdata.MockResourceID,
						Name: testdata.MockResourceName,
						Type: []string{"Resource"},
					},
					&voc.Resource{
						ID:   "00000000-0000-0000-000000000002",
						Name: "testResource2",
						Type: []string{"Resource"},
					},
				},
			},
			want: []interface{}{
				&voc.Resource{
					ID:   testdata.MockResourceID,
					Name: testdata.MockResourceName,
					Type: []string{"Resource"},
				},
				&voc.Resource{
					ID:   "00000000-0000-0000-000000000002",
					Name: "testResource2",
					Type: []string{"Resource"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, DeepCopy(tt.args.original), "DeepCopy(%v)", tt.args.original)
		})
	}
}
