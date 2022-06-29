package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFieldNames(t *testing.T) {
	var (
		fieldnames []string
		err        error
	)
	type someStruct struct {
		Name string
		// TODO(lebogg): When unexported, linter throws err 'field `secret` is unused (unused)'
		Secret int
	}

	// Type has to be struct, not string
	_, err = GetFieldNames[string]()
	assert.ErrorIs(t, err, ErrNoStruct)

	// Type has to be struct, not int
	_, err = GetFieldNames[int]()
	assert.ErrorIs(t, err, ErrNoStruct)

	// Successful
	fieldnames, err = GetFieldNames[someStruct]()
	assert.Equal(t, fieldnames, []string{"Name", "Secret"})

}
