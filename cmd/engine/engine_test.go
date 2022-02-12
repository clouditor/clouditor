package main

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Test_doCmd(t *testing.T) {
	viper.Set(DBInMemoryFlag, true)
	viper.Set(APIgRPCPortFlag, 0)
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
		})
	}
}
