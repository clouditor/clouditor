package inmemory

import (
	"testing"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/testutil/servicetest/orchestratortest"
	"github.com/stretchr/testify/assert"
)

// TestNewStorage is a simple test for NewStorage. If we implement in-memory our own, add more (table) tests
func TestNewStorage(t *testing.T) {
	s, err := NewStorage()
	assert.NoError(t, err)

	// Test to create new cloud service
	userInput := orchestratortest.NewCloudService()
	err = s.Create(userInput)
	assert.NoError(t, err)

	// Test if we get same user via its name
	userOutput := &orchestrator.CloudService{}
	err = s.Get(userOutput, "name = ?", testdata.MockCloudServiceName1)
	assert.NoError(t, err)
	assert.NoError(t, api.ValidateRequest(userOutput))
	assert.Equal(t, userInput, userOutput)
}
