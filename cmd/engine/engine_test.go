package main

import (
	"os"
	"testing"

	"clouditor.io/clouditor/internal/testutil/clitest"
	"clouditor.io/clouditor/server/rest"
	service_discovery "clouditor.io/clouditor/service/discovery"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
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
		want      assert.ValueAssertionFunc
		wantErr   bool
	}{
		{
			name: "Launch with --db-in-memory",
			prepViper: func() {
				viper.Set(DBInMemoryFlag, true)
				viper.Set(APIStartEmbeddedOAuth2ServerFlag, true)
				viper.Set(APIHTTPPortFlag, 0)
				viper.Set(APIgRPCPortFlag, 0)
				viper.Set(LogLevelFlag, DefaultLogLevel)
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				discoveryService := i1.(*service_discovery.Service)
				return assert.NotNil(t, discoveryService)
			},
		},
		{
			name: "Launch with invalid postgres port",
			prepViper: func() {
				viper.Set(DBPortFlag, 0)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			tt.prepViper()

			go func() {
				err := doCmd(tt.args.in0, tt.args.in1)
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
