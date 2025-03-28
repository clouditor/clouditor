package prototest

import (
	"testing"

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testutil/assert"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// NewAny creates a new [*anypb.Any] from a [proto.Message] with an assert that no error has been thrown.
func NewAny(t *testing.T, m proto.Message) *anypb.Any {
	a, err := anypb.New(m)
	assert.NoError(t, err)

	return a
}

// NewAny creates a new [*anypb.Any] from a [proto.Message] with a panic if an error has been thrown.
func NewAnyWithPanic(m proto.Message) *anypb.Any {
	a, err := anypb.New(m)
	if err != nil {
		panic(err)
	}

	return a
}

func NewProtobufResource(t *testing.T, or ontology.IsResource) *ontology.Resource {
	r := ontology.ProtoResource(or)
	assert.NotNil(t, r)

	return r
}
