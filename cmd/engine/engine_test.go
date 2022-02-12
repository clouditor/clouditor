package main

import (
	"testing"

	"clouditor.io/clouditor/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_doCmd(t *testing.T) {
	viper.Set(DBInMemoryFlag, true)
	viper.Set(APIHTTPPortFlag, 0)
	viper.Set(APIgRPCPortFlag, 0)

	type args struct {
		in0 *cobra.Command
		in1 []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Launch with --db-in-memory",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				if err := doCmd(tt.args.in0, tt.args.in1); (err != nil) != tt.wantErr {
					t.Errorf("doCmd() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()

			<-rest.GetReadyChannel()

			assert.NotNil(t, server)
			assert.NotNil(t, authService)
			assert.NotNil(t, discoveryService)
			assert.NotNil(t, assessmentService)
			assert.NotNil(t, evidenceStoreService)
		})
	}
}
