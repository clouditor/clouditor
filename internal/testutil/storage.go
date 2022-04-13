package testutil

import (
	"testing"

	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/persistence/inmemory"
)

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

type StorageWithError struct {
	SaveErr error
}

func (*StorageWithError) Create(r interface{}) error                       { return nil }
func (s *StorageWithError) Save(r interface{}, conds ...interface{}) error { return s.SaveErr }
func (*StorageWithError) Update(r interface{}, query interface{}, args ...interface{}) error {
	return nil
}
func (*StorageWithError) Get(r interface{}, conds ...interface{}) error            { return nil }
func (*StorageWithError) List(r interface{}, conds ...interface{}) error           { return nil }
func (*StorageWithError) Count(r interface{}, conds ...interface{}) (int64, error) { return 0, nil }
func (*StorageWithError) Delete(r interface{}, conds ...interface{}) error         { return nil }
