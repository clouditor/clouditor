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
	"github.com/spf13/cobra"
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

func TestInitConfig(t *testing.T) {
	tests := []struct {
		name      string
		prepViper func()
	}{
		{
			name:      "Happy path",
			prepViper: func() {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			tt.prepViper()

			InitConfig()

			assert.Equal(t, EnvPrefix, viper.GetEnvPrefix())
		})
	}
}

func TestInitCobra_EmbeddedOAuth2ServerEnabled(t *testing.T) {
	type args struct {
		engineCmd *cobra.Command
	}
	tests := []struct {
		name      string
		prepViper func()
		args      args
		want      assert.Want[bool]
	}{
		{
			name: "Happy path: EmbeddedOAuth2ServerEnabled set to false",
			prepViper: func() {
				viper.Set(EmbeddedOAuth2ServerEnabledFlag, false)
			},
			args: args{engineCmd: &cobra.Command{}},
			want: func(t *testing.T, got bool) bool {
				return assert.Equal(t, false, got)
			},
		},
		{
			name: "Happy path: EmbeddedOAuth2ServerEnabled set to true",
			prepViper: func() {
				viper.Set(EmbeddedOAuth2ServerEnabledFlag, true)
			},
			args: args{engineCmd: &cobra.Command{}},
			want: func(t *testing.T, got bool) bool {
				return assert.Equal(t, true, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			tt.prepViper()

			InitCobra(tt.args.engineCmd)

			tt.want(t, viper.GetBool(EmbeddedOAuth2ServerEnabledFlag))
		})
	}
}

func TestInitCobra_EmbeddedOAuth2ServerPublicURL(t *testing.T) {
	type args struct {
		engineCmd *cobra.Command
	}
	tests := []struct {
		name      string
		prepViper func()
		args      args
		want      assert.Want[string]
	}{
		{
			name: "Happy path: EmbeddedOAuth2ServerPublicURLFlag set",
			prepViper: func() {
				viper.Set(EmbeddedOAuth2ServerPublicURLFlag, "http://localhost")
			},
			args: args{engineCmd: &cobra.Command{}},
			want: func(t *testing.T, got string) bool {
				return assert.Equal(t, "http://localhost", got)
			},
		},
		{
			name: "Happy path: EmbeddedOAuth2ServerPublicURLFlag not set",
			prepViper: func() {
				viper.Set(EmbeddedOAuth2ServerPublicURLFlag, "")
			},
			args: args{engineCmd: &cobra.Command{}},
			want: func(t *testing.T, got string) bool {
				return assert.Equal(t, "", got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			tt.prepViper()

			InitCobra(tt.args.engineCmd)

			tt.want(t, viper.GetString(EmbeddedOAuth2ServerPublicURLFlag))
		})
	}
}
