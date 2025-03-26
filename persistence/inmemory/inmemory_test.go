package inmemory

import (
	"testing"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest/orchestratortest"
)

// TestNewStorage is a simple test for NewStorage. If we implement in-memory our own, add more (table) tests
func TestNewStorage(t *testing.T) {
	s, err := NewStorage()
	assert.NoError(t, err)

	// Test to create new target of evaluation
	userInput := orchestratortest.NewCertificationTarget()
	err = s.Create(userInput)
	assert.NoError(t, err)

	// Test if we get same user via its name
	userOutput := &orchestrator.TargetOfEvaluation{}
	err = s.Get(userOutput, "name = ?", testdata.MockCertificationTargetName1)
	assert.NoError(t, err)
	assert.NoError(t, api.Validate(userOutput))
	assert.Equal(t, userInput, userOutput)
}
