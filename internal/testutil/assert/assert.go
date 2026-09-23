// Package assert contains helpful assertion helpers. Under the hood it currently uses
// github.com/stretchr/testify/assert, but this might change in the future. In order to keep this transparent to the
// tests, unit tests should exclusively use this package. This also helps us keep track how often the individual assert
// functions are used and whether we can reduce the API surface of this package.
package assert

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/testing/protocmp"
)

var True = assert.True
var False = assert.False
var NotEmpty = assert.NotEmpty
var Contains = assert.Contains
var ErrorContains = assert.ErrorContains
var NoError = assert.NoError
var Error = assert.Error
var ErrorIs = assert.ErrorIs
var Fail = assert.Fail
var Same = assert.Same

type TestingT = assert.TestingT
type ErrorAssertionFunc = assert.ErrorAssertionFunc

// Want is a function type that can hold asserts in order to check the validity of "got".
type Want[T any] func(t *testing.T, got T) bool

var _ Want[any] = AnyValue[any]
var _ Want[any] = Nil[any]

// WantErr is a function type that can hold asserts in order to check the error of "err".
type WantErr func(t *testing.T, err error) bool

var _ WantErr = AnyValue[error]
var _ WantErr = Nil[error]

// CompareAllUnexported is a [cmp.Option] that allows the introspection of all un-exported fields in order to use them in
// [Equal] or [NotEqual].
func CompareAllUnexported() cmp.Option {
	return cmp.Exporter(func(reflect.Type) bool { return true })
}

// stripZeroValues is a cmp.Transformer that, after protocmp.Transform(),
// removes entries with zero values from protocmp.Message maps. This makes
// nil optional fields (absent key) compare equal to optional fields set
// to their zero value (e.g. *bool→false vs nil *bool).
var stripZeroValues = cmp.Transformer("stripZero", func(m protocmp.Message) protocmp.Message {
	result := protocmp.Message{}
	for k, v := range m {
		if !isZeroValue(v) {
			result[k] = v
		}
	}
	return result
})

func isZeroValue(v interface{}) bool {
	switch x := v.(type) {
	case bool:
		return !x
	case string:
		return x == ""
	case int:
		return x == 0
	case int32:
		return x == 0
	case int64:
		return x == 0
	case float32:
		return x == 0
	case float64:
		return x == 0
	case protocmp.Message:
		// Skip the "@type" field which is always present in protocmp.Message
		for k, sv := range x {
			if k == "@type" {
				continue
			}
			if !isZeroValue(sv) {
				return false
			}
		}
		return true
	case []interface{}:
		return len(x) == 0
	case nil:
		return true
	default:
		return false
	}
}

// Equal asserts that [got] and [want] are Equal. Under the hood, this uses the go-cmp package in combination with
// protocmp and also supplies a diff, in case the messages to do not match.
//
// Note: By default the option protocmp.Transform() will be used. This can cause problems with structs that are NOT
// protobuf messages and contain un-exported fields. In this case, [CompareAllUnexported] can be used instead.
func Equal[T any](t TestingT, want T, got T, opts ...cmp.Option) bool {
	tt, ok := t.(*testing.T)
	if ok {
		tt.Helper()
	}

	opts = append(opts, protocmp.Transform(), stripZeroValues)

	if cmp.Equal(got, want, opts...) {
		return true
	}

	return assert.Fail(t, "Not equal, but expected to be equal", cmp.Diff(got, want, opts...))
}

// NotEqual is similar to [Equal], but inverse.
func NotEqual[T any](t TestingT, want T, got T, opts ...cmp.Option) bool {
	tt, ok := t.(*testing.T)
	if ok {
		tt.Helper()
	}

	opts = append(opts, protocmp.Transform(), stripZeroValues)

	if !cmp.Equal(got, want, opts...) {
		return true
	}

	return assert.Fail(t, "Equal, but excepted to be not equal", cmp.Diff(got, want, opts...))
}

func Nil[T any](t *testing.T, obj T) bool {
	t.Helper()

	return assert.Nil(t, obj)
}

func NotNil[T any](t *testing.T, obj T) bool {
	t.Helper()

	return assert.NotNil(t, obj)
}

func Empty[T any](t *testing.T, obj T) bool {
	t.Helper()

	return assert.Empty(t, obj)
}

// Is asserts that a certain incoming object a (of type [any]) is of type T. It will return a type casted variant of
// that object in the return value obj, if it succeeded.
func Is[T any](t TestingT, a any) (obj T) {
	var ok bool

	obj, ok = a.(T)
	assert.True(t, ok)

	return
}

// AnyValue is a [Want] that accepts any value of T.
func AnyValue[T any](*testing.T, T) bool {
	return true
}

// Optional asserts the [want] function, if it is not nil. Otherwise, the assertion is ignored. This is helpful if an extra [Want] func is specified only for a select sub-set of table tests.
func Optional[T any](t *testing.T, want Want[T], got T) bool {
	if want != nil {
		return want(t, got)
	}

	return true
}
