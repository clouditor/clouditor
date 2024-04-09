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

package launcher

import (
	"testing"

	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/server"
	"clouditor.io/clouditor/v2/server/rest"
	"clouditor.io/clouditor/v2/service"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func TestLauncher_initStorage(t *testing.T) {
	type fields struct {
		name     string
		srv      *server.Server
		db       persistence.Storage
		log      *logrus.Entry
		grpcOpts []server.StartGRPCServerOption
		services []service.Service
	}
	tests := []struct {
		name      string
		prepViper func()
		fields    fields
		wantErr   assert.WantErr
	}{
		{
			name: "error setting storage config",
			prepViper: func() {
				viper.Set(config.DBHostFlag, "localhost")
				viper.Set(config.DBPortFlag, "8888")
				viper.Set(config.DBUserNameFlag, "postgres")
				viper.Set(config.DBPasswordFlag, "postgres")
				viper.Set(config.DBNameFlag, "testDB")
				viper.Set(config.DBSSLModeFlag, "disable")
			},
			fields: fields{
				name: "component",
			},
			wantErr: func(t *testing.T, err error) bool {
				return assert.ErrorContains(t, err, "could not create storage")
			},
		},
		{
			name: "Happy path: in-memory storage",
			prepViper: func() {
				viper.Set(config.DBInMemoryFlag, true)
			},
			fields: fields{
				name: "component",
			},
			wantErr: assert.Nil[error],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Launcher{
				name:     tt.fields.name,
				srv:      tt.fields.srv,
				db:       tt.fields.db,
				log:      tt.fields.log,
				grpcOpts: tt.fields.grpcOpts,
				services: tt.fields.services,
			}

			viper.Reset()
			tt.prepViper()

			err := l.initStorage()

			tt.wantErr(t, err)
		})
	}
}

func Test_printClouditorHeader(t *testing.T) {
	type args struct {
		component string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Happy path",
			args: args{
				component: "Component",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printClouditorHeader(tt.args.component)
		})
	}
}

func TestLauncher_initLogging(t *testing.T) {
	type fields struct {
		name     string
		srv      *server.Server
		db       persistence.Storage
		log      *logrus.Entry
		grpcOpts []server.StartGRPCServerOption
		services []service.Service
	}
	tests := []struct {
		name      string
		prepViper func()
		fields    fields
		want      assert.Want[*Launcher]
		wantErr   assert.WantErr
	}{
		{
			name: "error setting log level",
			prepViper: func() {
				viper.Set(config.LogLevelFlag, "wrongLogLevel")
			},
			fields: fields{
				name: "component",
			},
			want: assert.Nil[*Launcher],
			wantErr: func(t *testing.T, err error) bool {
				return assert.ErrorContains(t, err, "could not set log level")
			},
		},
		{
			name: "Happy path",
			prepViper: func() {
				viper.Set(config.LogLevelFlag, "info")
			},
			fields: fields{
				name: "component",
			},
			want: func(t *testing.T, got *Launcher) bool {
				return assert.NotNil(t, got.log)
			},
			wantErr: assert.Nil[error],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Launcher{
				name:     tt.fields.name,
				srv:      tt.fields.srv,
				db:       tt.fields.db,
				log:      tt.fields.log,
				grpcOpts: tt.fields.grpcOpts,
				services: tt.fields.services,
			}

			viper.Reset()
			tt.prepViper()

			err := l.initLogging()

			tt.wantErr(t, err)
		})
	}
}

func TestNewLauncher(t *testing.T) {
	type args struct {
		name  string
		specs []ServiceSpec
	}
	tests := []struct {
		name      string
		prepViper func()
		args      args
		wantL     assert.Want[*Launcher]
		wantErr   assert.WantErr
	}{
		{
			name: "error setting log level",
			prepViper: func() {
				viper.Set(config.LogLevelFlag, "wrongLogLevel")
			},
			args: args{
				name: "component",
			},
			wantL: assert.Nil[*Launcher],
			wantErr: func(t *testing.T, err error) bool {
				return assert.ErrorContains(t, err, "could not set log level")
			},
		},
		{
			name: "error setting storage",
			prepViper: func() {
				viper.Set(config.LogLevelFlag, "info")
			},
			args: args{
				name: "component",
			},
			wantL: assert.Nil[*Launcher],
			wantErr: func(t *testing.T, err error) bool {
				return assert.ErrorContains(t, err, "could not create storage")
			},
		},
		{
			// We are not able to check the individual fields in the Launcher struct, because they are all unexported. We need that check only to test for a valid execution of the method.
			name: "Happy path: without specs",
			prepViper: func() {
				viper.Set(config.LogLevelFlag, "info")
				viper.Set(config.DBInMemoryFlag, true)
			},
			args: args{
				name: "component",
			},
			wantL: func(t *testing.T, got *Launcher) bool {
				return assert.Equal(t, &Launcher{}, got, cmpopts.IgnoreUnexported(Launcher{}))
			},
			wantErr: assert.Nil[error],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			tt.prepViper()

			gotL, err := NewLauncher(tt.args.name, tt.args.specs...)

			tt.wantL(t, gotL)
			tt.wantErr(t, err)
		})
	}
}

func TestLauncher_Launch(t *testing.T) {
	type fields struct {
		name     string
		srv      *server.Server
		db       persistence.Storage
		log      *logrus.Entry
		grpcOpts []server.StartGRPCServerOption
		services []service.Service
	}
	tests := []struct {
		name      string
		prepViper func()
		fields    fields
		wantErr   bool
	}{
		{
			name: "Happy path: without embedded OAuth 2.0 server",
			prepViper: func() {
			},
			fields: fields{
				log: logrus.NewEntry(logrus.New()),
			},
			wantErr: true,
		},
		{
			name: "Happy path: with embedded OAuth 2.0 server",
			prepViper: func() {
				viper.Set(config.APIStartEmbeddedOAuth2ServerFlag, true)
				viper.Set(config.APIKeyPathFlag, "keyPath")
				viper.Set(config.APIKeyPasswordFlag, "passwd")
				viper.Set(config.APIKeySaveOnCreateFlag, true)
				viper.Set(config.DashboardURLFlag, "1.2.3.4")
				viper.Set(config.ServiceOAuth2ClientIDFlag, "clientID")
				viper.Set(config.APIDefaultUserFlag, "defaultUser")
				viper.Set(config.APIDefaultPasswordFlag, "defaultPasswd")
				viper.Set(config.APIHTTPPortFlag, 0)
				viper.Set(config.APIgRPCPortFlag, 0)
				viper.Set(config.LogLevelFlag, config.DefaultLogLevel)
			},
			fields: fields{
				log: logrus.NewEntry(logrus.New()),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			tt.prepViper()

			l := &Launcher{
				name:     tt.fields.name,
				srv:      tt.fields.srv,
				db:       tt.fields.db,
				log:      tt.fields.log,
				grpcOpts: tt.fields.grpcOpts,
				services: tt.fields.services,
			}

			go func() {
				err := l.Launch()
				if (err != nil) != tt.wantErr {
					t.Errorf("doCmd() error = %v, wantErr %v", err, tt.wantErr)
				}

				if err != nil {
					// Signal that we are ready anyway, so that we fail properly
					rest.GetReadyChannel() <- false
				}
			}()

			success := <-rest.GetReadyChannel()
			if success {
				assert.NotNil(t, l.srv)
			}
		})
	}
}
