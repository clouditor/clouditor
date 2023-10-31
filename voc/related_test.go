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

package voc

import (
	"reflect"
	"testing"
)

func TestVirtualMachine_Related(t *testing.T) {
	type fields struct {
		Compute           *Compute
		BlockStorage      []ResourceID
		NetworkInterfaces []ResourceID
		BootLogging       *BootLogging
		OSLogging         *OSLogging
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "Related VM resources",
			fields: fields{
				BootLogging: &BootLogging{
					Logging: &Logging{LoggingService: []ResourceID{"3"}},
				},
				OSLogging: &OSLogging{
					Logging: &Logging{LoggingService: []ResourceID{"5"}},
				},
				BlockStorage: []ResourceID{"1"},
				Compute: &Compute{
					Resource: &Resource{
						Parent: "6",
					},
					ResourceLogging: &ResourceLogging{
						&Logging{LoggingService: []ResourceID{"4"}},
					},
					NetworkInterfaces: []ResourceID{"2"},
				},
			},
			want: []string{"1", "3", "5", "2", "4", "6"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v IsCloudResource = VirtualMachine{
				Compute:      tt.fields.Compute,
				BlockStorage: tt.fields.BlockStorage,
				BootLogging:  tt.fields.BootLogging,
				OsLogging:    tt.fields.OSLogging,
			}
			if got := v.Related(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VirtualMachine.Related() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoggingService_Related(t *testing.T) {
	type fields struct {
		NetworkService *NetworkService
		Storage        []ResourceID
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "Related LoggingService resources",
			fields: fields{
				NetworkService: &NetworkService{
					Networking: &Networking{
						Resource: &Resource{},
					},
				},
				Storage: []ResourceID{"1"},
			},
			want: []string{"1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v IsCloudResource = LoggingService{
				NetworkService: tt.fields.NetworkService,
				Storage:        tt.fields.Storage,
			}
			if got := v.Related(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VirtualMachine.Related() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorageService_Related(t *testing.T) {
	type fields struct {
		NetworkService *NetworkService
		Storage        []ResourceID
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "Related LoggingService resources",
			fields: fields{
				NetworkService: &NetworkService{
					Networking: &Networking{
						Resource: &Resource{},
					},
				},
				Storage: []ResourceID{"1"},
			},
			want: []string{"1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StorageService{
				NetworkService: tt.fields.NetworkService,
				Storage:        tt.fields.Storage,
			}
			if got := s.Related(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StorageService.Related() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Related(t *testing.T) {
	type fields struct {
		Resource         *Resource
		AtRestEncryption IsAtRestEncryption
		Backups          []*Backup
		Immutability     *Immutability
		ResourceLogging  *ResourceLogging
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "Related LoggingService resources",
			fields: fields{
				Resource: &Resource{},
				Backups: []*Backup{
					{
						Storage: "1",
					},
				},
			},
			want: []string{"1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Storage{
				Resource:         tt.fields.Resource,
				AtRestEncryption: tt.fields.AtRestEncryption,
				Backups:          tt.fields.Backups,
				Immutability:     tt.fields.Immutability,
				ResourceLogging:  tt.fields.ResourceLogging,
			}
			if got := s.Related(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.Related() = %v, want %v", got, tt.want)
			}
		})
	}
}
