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
	"reflect"
	"testing"

	"clouditor.io/clouditor/api/orchestrator"
)

func TestLoadRequirements(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name             string
		args             args
		wantRequirements []*orchestrator.Requirement
		wantErr          bool
	}{
		{
			name: "load",
			args: args{
				file: "requirements.json",
			},
			wantRequirements: []*orchestrator.Requirement{
				{
					Id:          "Req-1",
					Name:        "Make-it-Secure",
					Description: "You should make everything secure",
					Metrics: &orchestrator.Requirement_Metrics{
						MetricIds: []string{
							"TransportEncryptionEnabled",
							"TransportEncryptionAlgorithm",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRequirements, err := LoadRequirements(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadRequirements() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRequirements, tt.wantRequirements) {
				t.Errorf("LoadRequirements() = %v, want %v", gotRequirements, tt.wantRequirements)
			}
		})
	}
}
