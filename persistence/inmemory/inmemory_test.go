package inmemory

import (
	"testing"

	"clouditor.io/clouditor/api/auth"
	"github.com/stretchr/testify/assert"
)

// TestNewStorage is a simple test for NewStorage. If we implement in-memory our own, add more (table) tests
func TestNewStorage(t *testing.T) {
	s, err := NewStorage()
	assert.NoError(t, err)

	// Test to create new user
	userInput := &auth.User{
		Username: "SomeName",
		Password: "SomePassword",
		Email:    "SomeMail",
		FullName: "SomeFullName",
	}
	err = s.Create(userInput)
	assert.NoError(t, err)

	// Test if we get same user via its name
	userOutput := &auth.User{}
	err = s.Get(userOutput, "Username = ?", "SomeName")
	assert.NoError(t, err)
	assert.Equal(t, userInput, userOutput)

	list := make([]auth.User, 0)
	err = s.List(&list)
	assert.ErrorIs(t, err, nil)
	assert.NotEmpty(t, list)
}
