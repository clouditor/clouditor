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

package util

import (
	"reflect"
	"testing"

	"clouditor.io/clouditor/internal/testutil/prototest"
)

// func TestRemoveIndexFromSlice(t *testing.T) {
// 	inputValue := []*prototest.TestStruct{
// 		{
// 			TestName:        "testName_1",
// 			TestId:          "testId_1",
// 			TestDescription: "testDescription_1",
// 		},
// 		{
// 			TestName:        "testName_2",
// 			TestId:          "testId_2",
// 			TestDescription: "testDescription_2",
// 		},
// 		{
// 			TestName:        "testName_3",
// 			TestId:          "testId_3",
// 			TestDescription: "testDescription_3",
// 		},
// 	}
// 	outputValueRemoveFirstElement := []*prototest.TestStruct{
// 		{
// 			TestName:        "testName_2",
// 			TestId:          "testId_2",
// 			TestDescription: "testDescription_2",
// 		},
// 		{
// 			TestName:        "testName_3",
// 			TestId:          "testId_3",
// 			TestDescription: "testDescription_3",
// 		},
// 	}

// 	outputValueRemoveSecondElement := []*prototest.TestStruct{
// 		{
// 			TestName:        "testName_1",
// 			TestId:          "testId_1",
// 			TestDescription: "testDescription_1",
// 		},
// 		{
// 			TestName:        "testName_3",
// 			TestId:          "testId_3",
// 			TestDescription: "testDescription_3",
// 		},
// 	}

// 	outputValueRemoveLastElement := []*prototest.TestStruct{
// 		{
// 			TestName:        "testName_1",
// 			TestId:          "testId_1",
// 			TestDescription: "testDescription_1",
// 		},
// 		{
// 			TestName:        "testName_2",
// 			TestId:          "testId_2",
// 			TestDescription: "testDescription_2",
// 		},
// 	}

// 	emptyValue := []*prototest.TestStruct{}

// 	// Remove first element from inputValue
// 	assert.Equal(t, outputValueRemoveFirstElement, RemoveIndexFromSlice(inputValue, 0))

// 	// Remove second element from inputValue
// 	result := RemoveIndexFromSlice(inputValue, 1)
// 	assert.Equal(t, outputValueRemoveSecondElement, result)

// 	// Remove last element from inputValue
// 	assert.Equal(t, outputValueRemoveLastElement, RemoveIndexFromSlice(inputValue, 2))

// 	// Remove n+1 element from inputValue
// 	assert.Equal(t, inputValue, RemoveIndexFromSlice(inputValue, 3))

// 	// Remove second element from empty inputValue
// 	assert.Equal(t, emptyValue, RemoveIndexFromSlice(emptyValue, 1))

// }

func TestRemoveIndexFromSlice(t *testing.T) {
	type args struct {
		slice []*prototest.TestStruct
		index int
	}
	tests := []struct {
		name string
		args args
		want []*prototest.TestStruct
	}{
		{
			name: "Remove first element",
			args: args{
				slice: []*prototest.TestStruct{
					{
						TestName:        "testName_1",
						TestId:          "testId_1",
						TestDescription: "testDescription_1",
					},
					{
						TestName:        "testName_2",
						TestId:          "testId_2",
						TestDescription: "testDescription_2",
					},
					{
						TestName:        "testName_3",
						TestId:          "testId_3",
						TestDescription: "testDescription_3",
					},
				},
				index: 0,
			},
			want: []*prototest.TestStruct{
				{
					TestName:        "testName_2",
					TestId:          "testId_2",
					TestDescription: "testDescription_2",
				},
				{
					TestName:        "testName_3",
					TestId:          "testId_3",
					TestDescription: "testDescription_3",
				},
			},
		},
		{
			name: "Remove second element",
			args: args{
				slice: []*prototest.TestStruct{
					{
						TestName:        "testName_1",
						TestId:          "testId_1",
						TestDescription: "testDescription_1",
					},
					{
						TestName:        "testName_2",
						TestId:          "testId_2",
						TestDescription: "testDescription_2",
					},
					{
						TestName:        "testName_3",
						TestId:          "testId_3",
						TestDescription: "testDescription_3",
					},
				},
				index: 1,
			},
			want: []*prototest.TestStruct{
				{
					TestName:        "testName_1",
					TestId:          "testId_1",
					TestDescription: "testDescription_1",
				},
				{
					TestName:        "testName_3",
					TestId:          "testId_3",
					TestDescription: "testDescription_3",
				},
			},
		},
		{
			name: "Remove last element",
			args: args{
				slice: []*prototest.TestStruct{
					{
						TestName:        "testName_1",
						TestId:          "testId_1",
						TestDescription: "testDescription_1",
					},
					{
						TestName:        "testName_2",
						TestId:          "testId_2",
						TestDescription: "testDescription_2",
					},
					{
						TestName:        "testName_3",
						TestId:          "testId_3",
						TestDescription: "testDescription_3",
					},
				},
				index: 2,
			},
			want: []*prototest.TestStruct{
				{
					TestName:        "testName_1",
					TestId:          "testId_1",
					TestDescription: "testDescription_1",
				},
				{
					TestName:        "testName_2",
					TestId:          "testId_2",
					TestDescription: "testDescription_2",
				},
			},
		},
		{
			name: "Remove n+1 element",
			args: args{
				slice: []*prototest.TestStruct{
					{
						TestName:        "testName_1",
						TestId:          "testId_1",
						TestDescription: "testDescription_1",
					},
					{
						TestName:        "testName_2",
						TestId:          "testId_2",
						TestDescription: "testDescription_2",
					},
					{
						TestName:        "testName_3",
						TestId:          "testId_3",
						TestDescription: "testDescription_3",
					},
				},
				index: 3,
			},
			want: []*prototest.TestStruct{
				{
					TestName:        "testName_1",
					TestId:          "testId_1",
					TestDescription: "testDescription_1",
				},
				{
					TestName:        "testName_2",
					TestId:          "testId_2",
					TestDescription: "testDescription_2",
				},
				{
					TestName:        "testName_3",
					TestId:          "testId_3",
					TestDescription: "testDescription_3",
				},
			},
		},
		{
			name: "Empty input",
			args: args{
				slice: []*prototest.TestStruct{},
				index: 3,
			},
			want: []*prototest.TestStruct{},
		},
		{
			name: "Missing input",
			args: args{
				slice: nil,
				index: 3,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveIndexFromSlice(tt.args.slice, tt.args.index); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveIndexFromSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
