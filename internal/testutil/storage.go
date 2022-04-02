package testutil

import (
	"testing"

	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/persistence/inmemory"
)

func NewInMemoryStorage(t *testing.T) (s persistence.Storage) {
	var err error

	s, err = inmemory.NewStorage()
	if err != nil {
		t.Errorf("Could not initialize in-memory db: %v", err)
	}

	return
}
