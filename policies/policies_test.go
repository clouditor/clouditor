package policies_test

import (
	"testing"

	"clouditor.io/clouditor/policies"
)

func TestRun(t *testing.T) {
	policies.Run("tls.rego")
}
