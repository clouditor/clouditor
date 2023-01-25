package inmemory

import (
	"testing"

	"clouditor.io/clouditor/api/auth"
	"clouditor.io/clouditor/internal/testdata"
	"github.com/stretchr/testify/assert"
)

// TestNewStorage is a simple test for NewStorage. If we implement in-memory our own, add more (table) tests
func TestNewStorage(t *testing.T) {
	s, err := NewStorage()
	assert.NoError(t, err)

	// Test to create new user
	userInput := testdata.NewUser1()
	err = s.Create(userInput)
	assert.NoError(t, err)

	// Test if we get same user via its name
	userOutput := &auth.User{}
	err = s.Get(userOutput, "Username = ?", "SomeName")
	assert.NoError(t, err)
	assert.NoError(t, userOutput.Validate())
	assert.Equal(t, userInput, userOutput)
}
