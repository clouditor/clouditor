// Copyright 2024 Fraunhofer AISEC
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

package openstack

import (
	"testing"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
)

func TestNewOpenstackDiscovery(t *testing.T) {
	type args struct {
		opts []DiscoveryOption
	}
	tests := []struct {
		name string
		args args
		want assert.Want[discovery.Discoverer]
	}{
		{
			name: "error: oauthOpts not set",
			args: args{},
			want: assert.Nil[discovery.Discoverer],
		},
		{
			name: "Happy path: with certification target id",
			args: args{
				opts: []DiscoveryOption{
					WithAuthorizer(&AuthOptions{
						IdentityEndpoint: testdata.MockOpenstackIdentitiyEndpoint,
						Username:         testdata.MockOpenstackUsername,
						Password:         testdata.MockOpenstackPassword,
						TenantName:       testdata.MockOpenstackTenantName,
						AllowReauth:      true,
					}),
					WithCertificationTargetID(testdata.MockCertificationTargetID2),
				},
			},
			want: func(t *testing.T, got discovery.Discoverer) bool {
				assert.Equal(t, testdata.MockCertificationTargetID2, got.CertificationTargetID())
				return assert.NotNil(t, got)
			},
		},
		{
			name: "Happy path: with authorizer",
			args: args{
				opts: []DiscoveryOption{
					WithAuthorizer(&AuthOptions{
						IdentityEndpoint: testdata.MockOpenstackIdentitiyEndpoint,
						Username:         testdata.MockOpenstackUsername,
						Password:         testdata.MockOpenstackPassword,
						TenantName:       testdata.MockOpenstackTenantName,
						AllowReauth:      true,
					}),
				},
			},
			want: func(t *testing.T, got discovery.Discoverer) bool {
				assert.Equal(t, config.DefaultCertificationTargetID, got.CertificationTargetID())
				return assert.NotNil(t, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewOpenstackDiscovery(tt.args.opts...)
			tt.want(t, got)
		})
	}
}
