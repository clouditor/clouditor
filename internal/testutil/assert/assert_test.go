// Package assert contains helpful assertion helpers. Under the hood it currently uses
// github.com/stretchr/testify/assert, but this might change in the future. In order to keep this transparent to the
// tests, unit tests should exclusively use this package. This also helps us keep track how often the individual assert
// functions are used and whether we can reduce the API surface of this package.
package assert

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MyStruct struct {
	A string

	b int
}

// fakeT is an implementation of [assert.AssertT] that does nothing, so that we can assert that an assert fails without
// failing ;)
type fakeT struct{}

func (fakeT) Errorf(format string, args ...any) {
	// do nothing
}

func TestEqual(t *testing.T) {
	type args struct {
		t    TestingT
		want any
		got  any
		opts []cmp.Option
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "compare literals",
			args: args{
				t:    t,
				want: "5",
				got:  "5",
			},
			want: true,
		},
		{
			name: "compare ordinary structs with unexported fields",
			args: args{
				t:    t,
				want: &MyStruct{A: "test", b: 1},
				got:  &MyStruct{A: "test", b: 1},
				opts: []cmp.Option{CompareAllUnexported()},
			},
			want: true,
		},
		{
			name: "compare protobuf",
			args: args{
				t:    &fakeT{},
				want: timestamppb.New(time.Unix(0, 0)),
				got:  timestamppb.Now(),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Equal(tt.args.t, tt.args.want, tt.args.got, tt.args.opts...); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotEqual(t *testing.T) {
	type args struct {
		t    TestingT
		want any
		got  any
		opts []cmp.Option
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "compare literals fail",
			args: args{
				t:    &fakeT{},
				want: "5",
				got:  "5",
			},
			want: false,
		},
		{
			name: "compare literals",
			args: args{
				t:    t,
				want: "6",
				got:  "5",
			},
			want: true,
		},
		{
			name: "compare protobuf",
			args: args{
				t:    t,
				want: timestamppb.New(time.Unix(0, 0)),
				got:  timestamppb.Now(),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NotEqual(tt.args.t, tt.args.want, tt.args.got, tt.args.opts...); got != tt.want {
				t.Errorf("NotEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
