package policies_test

import (
	"testing"

	"clouditor.io/clouditor/policies"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	err := policies.Run("tls.rego")

	assert.Nil(t, err)
}
