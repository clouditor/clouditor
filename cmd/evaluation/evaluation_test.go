// Copyright 2016-2020 Fraunhofer AISEC
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

package main

import (
	"os"
	"testing"

	"clouditor.io/clouditor/v2/api/evaluation"
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/clitest"
	"clouditor.io/clouditor/v2/server/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	clitest.AutoChdir()

	os.Exit(m.Run())
}

func Test_doCmd(t *testing.T) {
	type args struct {
		in0 *cobra.Command
		in1 []string
	}
	tests := []struct {
		name      string
		prepViper func()
		args      args
		want      assert.Want[evaluation.EvaluationServer]
		wantErr   assert.WantErr
	}{
		{
			name: "Launch with invalid postgres port",
			prepViper: func() {
				viper.Set(config.LogLevelFlag, config.DefaultLogLevel)
				viper.Set(config.DBPortFlag, 0)
			},
			want: assert.Nil[evaluation.EvaluationServer],
			wantErr: func(t *testing.T, err error) bool {
				return assert.ErrorContains(t, err, "could not create storage:")
			},
		},
		{
			name: "Happy path: launch with --db-in-memory",
			prepViper: func() {
				viper.Set(config.DBInMemoryFlag, true)
				viper.Set(config.APIgRPCPortFlag, "9081")
				viper.Set(config.APIHTTPPortFlag, "8081")
				viper.Set(config.LogLevelFlag, config.DefaultLogLevel)
			},
			want: func(t *testing.T, got evaluation.EvaluationServer) bool {
				return assert.NotNil(t, got)
			},
			wantErr: assert.Nil[error],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			tt.prepViper()

			go func() {
				err := doCmd(tt.args.in0, tt.args.in1)
				tt.wantErr(t, err)

				if err != nil {
					// Signal that we are ready anyway, so that we fail properly
					rest.GetReadyChannel() <- false
				}
			}()

			success := <-rest.GetReadyChannel()
			if success {
				assert.NotNil(t, srv)
				assert.NotNil(t, evaluationService)
			}

			tt.want(t, evaluationService)
		})
	}
}
