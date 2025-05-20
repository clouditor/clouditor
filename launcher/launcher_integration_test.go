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

package launcher_test

import (
	"testing"

	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/launcher"
	"clouditor.io/clouditor/v2/service/assessment"
	"clouditor.io/clouditor/v2/service/discovery"
	"clouditor.io/clouditor/v2/service/evaluation"
	"clouditor.io/clouditor/v2/service/evidence"
	"clouditor.io/clouditor/v2/service/orchestrator"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/spf13/viper"
)

// We need to test this function here, otherwise we get an import cycle with the packages discovery and launcher
func TestNewLauncher(t *testing.T) {
	type args struct {
		name  string
		specs []launcher.ServiceSpec
	}
	tests := []struct {
		name      string
		prepViper func()
		args      args
		wantL     assert.Want[*launcher.Launcher]
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
			wantL: assert.Nil[*launcher.Launcher],
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
				name:  "component",
				specs: []launcher.ServiceSpec{evidence.DefaultServiceSpec()},
			},
			wantL: assert.Nil[*launcher.Launcher],
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
			wantL: func(t *testing.T, got *launcher.Launcher) bool {
				return assert.Equal(t, &launcher.Launcher{}, got, cmpopts.IgnoreUnexported(launcher.Launcher{}))
			},
			wantErr: assert.Nil[error],
		},
		{
			name: "Happy path: discovery spec",
			prepViper: func() {
				viper.Set(config.LogLevelFlag, "info")
				viper.Set(config.DBInMemoryFlag, true)
			},
			args: args{
				name:  "component",
				specs: []launcher.ServiceSpec{discovery.DefaultServiceSpec()},
			},
			wantL: func(t *testing.T, got *launcher.Launcher) bool {
				return assert.Equal(t, &launcher.Launcher{}, got, cmpopts.IgnoreUnexported(launcher.Launcher{}))
			},
			wantErr: assert.Nil[error],
		},
		{
			name: "Happy path: assessment spec",
			prepViper: func() {
				viper.Set(config.LogLevelFlag, "info")
				viper.Set(config.DBInMemoryFlag, true)
			},
			args: args{
				name:  "component",
				specs: []launcher.ServiceSpec{assessment.DefaultServiceSpec()},
			},
			wantL: func(t *testing.T, got *launcher.Launcher) bool {
				return assert.Equal(t, &launcher.Launcher{}, got, cmpopts.IgnoreUnexported(launcher.Launcher{}))
			},
			wantErr: assert.Nil[error],
		},
		{
			name: "Happy path: orchestrator spec",
			prepViper: func() {
				viper.Set(config.LogLevelFlag, "info")
				viper.Set(config.DBInMemoryFlag, true)
			},
			args: args{
				name:  "component",
				specs: []launcher.ServiceSpec{orchestrator.DefaultServiceSpec()},
			},
			wantL: func(t *testing.T, got *launcher.Launcher) bool {
				return assert.Equal(t, &launcher.Launcher{}, got, cmpopts.IgnoreUnexported(launcher.Launcher{}))
			},
			wantErr: assert.Nil[error],
		},
		{
			name: "Happy path: evaluation spec",
			prepViper: func() {
				viper.Set(config.LogLevelFlag, "info")
				viper.Set(config.DBInMemoryFlag, true)
			},
			args: args{
				name:  "component",
				specs: []launcher.ServiceSpec{evaluation.DefaultServiceSpec()},
			},
			wantL: func(t *testing.T, got *launcher.Launcher) bool {
				return assert.Equal(t, &launcher.Launcher{}, got, cmpopts.IgnoreUnexported(launcher.Launcher{}))
			},
			wantErr: assert.Nil[error],
		},
		{
			name: "Happy path: evidence spec",
			prepViper: func() {
				viper.Set(config.LogLevelFlag, "info")
				viper.Set(config.DBInMemoryFlag, true)
			},
			args: args{
				name:  "component",
				specs: []launcher.ServiceSpec{evidence.DefaultServiceSpec()},
			},
			wantL: func(t *testing.T, got *launcher.Launcher) bool {
				return assert.Equal(t, &launcher.Launcher{}, got, cmpopts.IgnoreUnexported(launcher.Launcher{}))
			},
			wantErr: assert.Nil[error],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			tt.prepViper()

			gotL, err := launcher.NewLauncher(tt.args.name, tt.args.specs...)

			tt.wantL(t, gotL)
			tt.wantErr(t, err)
		})
	}
}
