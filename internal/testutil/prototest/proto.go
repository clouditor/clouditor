package prototest

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/anypb"
)

// NewAny creates a new [*anypb.Any] from a [proto.Message] with an assert that no error has been thrown.
func NewAny(t *testing.T, m proto.Message) *anypb.Any {
	a, err := anypb.New(m)
	assert.NoError(t, err)

	return a
}

// NewAny creates a new [*anypb.Any] from a [proto.Message] with a panic that no error has been thrown.
func NewAnyWithPanic(m proto.Message) *anypb.Any {
	a, err := anypb.New(m)
	if err != nil {
		panic(err)
	}

	return a
}

// Equal asserts that [got] and [want] are two equal protobuf messages. Under the hood, this uses the go-cmp package in
// combination with protocmp and also supplies a diff, in case the messages to do not match.
func Equal[T proto.Message](t *testing.T, want T, got T) bool {
	return assert.Empty(t, cmp.Diff(got, want, protocmp.Transform()))
}

// EqualSlice asserts that [got] and [want] are two equal protobuf message slices. Under the hood, this uses the go-cmp
// package in combination with protocmp and also supplies a diff, in case the messages to do not match.
func EqualSlice[T proto.Message](t *testing.T, want []T, got []T) bool {
	return assert.Empty(t, cmp.Diff(got, want, protocmp.Transform()))
}
