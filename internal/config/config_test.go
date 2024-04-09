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

package config

import (
	"testing"

	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/spf13/viper"
	"golang.org/x/oauth2/clientcredentials"
)

func TestClientCredentials(t *testing.T) {
	tests := []struct {
		name      string
		prepViper func()
		want      assert.Want[*clientcredentials.Config]
	}{
		{
			name: "Happy path",
			prepViper: func() {
				viper.Set(ServiceOAuth2ClientIDFlag, "clientID")
				viper.Set(ServiceOAuth2ClientSecretFlag, "clientSecret")
				viper.Set(ServiceOAuth2EndpointFlag, "1.1.1.1")
			},
			want: func(t *testing.T, got *clientcredentials.Config) bool {
				want := &clientcredentials.Config{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					TokenURL:     "1.1.1.1",
				}
				return assert.Equal(t, want, got, cmpopts.IgnoreUnexported(clientcredentials.Config{}))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			tt.prepViper()

			got := ClientCredentials()

			tt.want(t, got)
		})
	}
}
