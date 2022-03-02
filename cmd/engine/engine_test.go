package main

import (
	"clouditor.io/clouditor/logging/formatter"
	"github.com/sirupsen/logrus"
	"testing"

	"clouditor.io/clouditor/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func init() {
	log = logrus.WithField("component", "engine-tests")
	log.Logger.Formatter = formatter.CapitalizeFormatter{Formatter: &logrus.TextFormatter{ForceColors: true, FullTimestamp: true}}
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
		wantErr   bool
	}{
		{
			name: "Launch with --db-in-memory",
			prepViper: func() {
				viper.Set(DBInMemoryFlag, true)
				viper.Set(APIHTTPPortFlag, 0)
				viper.Set(APIgRPCPortFlag, 0)
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
				assert.NotNil(t, server)
				assert.NotNil(t, authService)
				assert.NotNil(t, discoveryService)
				assert.NotNil(t, assessmentService)
				assert.NotNil(t, evidenceStoreService)
			}
		})
	}
}
