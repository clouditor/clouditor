package ccl_test

import (
	"encoding/json"
	"testing"

	"clouditor.io/clouditor/ccl"
	"github.com/stretchr/testify/assert"
)

func TestRule(t *testing.T) {
	var j = `{"field": "value"}`
	var o map[string]interface{}

	json.Unmarshal([]byte(j), &o)

	success, err := ccl.RunRule("test_rule.ccl", o)

	assert.Nil(t, err)
	assert.True(t, success)
}
