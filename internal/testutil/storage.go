package testutil

import (
	"testing"

	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/persistence/inmemory"
)

// NewInMemoryStorage uses the inmemory package to create a new in-memory storage that can be used
// for unit testing. The funcs varargs can be used to immediately execute storage operations on it.
func NewInMemoryStorage(t *testing.T, funcs ...func(s persistence.Storage)) (s persistence.Storage) {
	var err error

	s, err = inmemory.NewStorage()
	if err != nil {
		t.Errorf("Could not initialize in-memory db: %v", err)
	}

	for _, f := range funcs {
		f(s)
	}

	return
}

// StorageWithError can be used to introduce various errors in a storage operation during unit testing.
type StorageWithError struct {
	SaveErr error
	GetErr  error
}

func (*StorageWithError) Create(r interface{}) error                       { return nil }
func (s *StorageWithError) Save(r interface{}, conds ...interface{}) error { return s.SaveErr }
func (*StorageWithError) Update(r interface{}, query interface{}, args ...interface{}) error {
	return nil
}
func (s *StorageWithError) Get(r interface{}, conds ...interface{}) error          { return s.GetErr }
func (*StorageWithError) List(r interface{}, conds ...interface{}) error           { return nil }
func (*StorageWithError) Count(r interface{}, conds ...interface{}) (int64, error) { return 0, nil }
func (*StorageWithError) Delete(r interface{}, conds ...interface{}) error         { return nil }
