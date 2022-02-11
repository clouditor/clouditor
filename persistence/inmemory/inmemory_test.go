package inmemory

import (
	"clouditor.io/clouditor/api/auth"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNewStorage is a simple test for NewStorage. If we implement inmemory our own, add more (table) tests
func TestNewStorage(t *testing.T) {
	s, err := NewStorage()
	assert.ErrorIs(t, err, nil)

	// Test to create new user
	userInput := &auth.User{
		Username: "SomeName",
		Password: "SomePassword",
		Email:    "SomeMail",
		FullName: "SomeFullName",
	}
	err = s.Create(userInput)
	assert.ErrorIs(t, err, nil)

	// Test if we get same user via its ID
	userOutput := &auth.User{}
	err = s.Get(userOutput, "Id = ?", "SomeName")
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, userInput, userOutput)
}
