package api

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestStatusFromWrappedError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		args   args
		wantS  *status.Status
		wantOk bool
	}{
		{
			name: "err is status.Status",
			args: args{
				err: status.Errorf(codes.NotFound, "not found"),
			},
			wantS:  status.New(codes.NotFound, "not found"),
			wantOk: true,
		},
		{
			name: "wrapped in fmt.Errorf",
			args: args{
				err: fmt.Errorf("some error: %w", status.Errorf(codes.NotFound, "not found")),
			},
			wantS:  status.New(codes.NotFound, "not found"),
			wantOk: true,
		},
		{
			name: "no status",
			args: args{
				err: fmt.Errorf("some error: %w", errors.New("some other error")),
			},
			wantS:  nil,
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotS, gotOk := StatusFromWrappedError(tt.args.err)
			if !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("StatusFromWrappedError() gotS = %v, want %v", gotS, tt.wantS)
			}
			if gotOk != tt.wantOk {
				t.Errorf("StatusFromWrappedError() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
