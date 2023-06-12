package testutil

import (
	"clouditor.io/clouditor/voc"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
)

// ToStruct transforms r to a struct and asserts if it was successful
func ToStruct(r voc.IsCloudResource, t assert.TestingT) (s *structpb.Value) {
	s, err := voc.ToStruct(r)
	if err != nil {
		assert.Error(t, err)
	}

	return
}
