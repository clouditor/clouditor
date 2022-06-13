package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFieldNames(t *testing.T) {
	type someStruct struct {
		Name   string
		secret int
	}
	type args struct {
		aStruct any
	}
	tests := []struct {
		name           string
		args           args
		wantFieldNames []string
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name: "nil struct",
			args: args{aStruct: nil},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrStructIsNil)
			},
		},
		{
			name: "not of type struct",
			args: args{aStruct: ""},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrNoStruct)
			},
		},
		{
			name:           "successful",
			args:           args{aStruct: someStruct{}},
			wantFieldNames: []string{"Name", "secret"},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFieldNames, err := GetFieldNames(tt.args.aStruct)
			if !tt.wantErr(t, err, fmt.Sprintf("GetFieldNames(%v)", tt.args.aStruct)) {
				return
			}
			// Type assertion avoids linter err 'field `secret` is unused (unused)'
			_ = tt.args.aStruct.(someStruct)
			assert.Equalf(t, tt.wantFieldNames, gotFieldNames, "GetFieldNames(%v)", tt.args.aStruct)
		})
	}
}
