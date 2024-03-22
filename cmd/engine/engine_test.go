package main

import (
	"os"
	"testing"

	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/clitest"
	"clouditor.io/clouditor/v2/server/rest"
	service_discovery "clouditor.io/clouditor/v2/service/discovery"

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
		want      assert.Want[*service_discovery.Service]
		wantErr   assert.WantErr
	}{
		{
			name: "Launch without log level",
			prepViper: func() {
				viper.Set(config.LogLevelFlag, "")
			},
			want: assert.Nil[*service_discovery.Service],
			wantErr: func(t *testing.T, err error) bool {
				return assert.ErrorContains(t, err, "not a valid logrus Level:")
			},
		},
		{
			name: "Launch with --db-in-memory",
			prepViper: func() {
				viper.Set(config.DBInMemoryFlag, true)
				viper.Set(config.APIStartEmbeddedOAuth2ServerFlag, true)
				viper.Set(config.APIHTTPPortFlag, 0)
				viper.Set(config.APIgRPCPortFlag, 0)
				viper.Set(config.LogLevelFlag, config.DefaultLogLevel)
			},
			want: func(t *testing.T, got *service_discovery.Service) bool {
				return assert.NotNil(t, got)
			},
			wantErr: assert.Nil[error],
		},
		{
			name: "Launch with invalid postgres port",
			prepViper: func() {
				viper.Set(config.LogLevelFlag, config.DefaultLogLevel)
				viper.Set(config.DBPortFlag, 0)
			},
			wantErr: func(t *testing.T, err error) bool {
				return assert.ErrorContains(t, err, "could not create storage:")
			},
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
				assert.NotNil(t, discoveryService)
				assert.NotNil(t, assessmentService)
				assert.NotNil(t, evidenceStoreService)
			}

			if tt.want != nil {
				tt.want(t, discoveryService)
			}
		})
	}
}
