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
		t.Fatalf("Could not initialize in-memory db: %v", err)
	}

	for _, f := range funcs {
		f(s)
	}

	return
}

// StorageWithError can be used to introduce various errors in a storage operation during unit testing.
type StorageWithError struct {
	CreateErr error
	SaveErr   error
	UpdateErr error
	GetErr    error
	ListErr   error
	RawErr    error
	CountRes  int64
	CountErr  error
	DeleteErr error
}

func (s *StorageWithError) Create(_ any) error         { return s.CreateErr }
func (s *StorageWithError) Save(_ any, _ ...any) error { return s.SaveErr }
func (*StorageWithError) Update(_ any, _ ...any) error {
	return nil
}
func (s *StorageWithError) Get(_ any, _ ...any) error { return s.GetErr }
func (s *StorageWithError) List(_ any, _ string, _ bool, _ int, _ int, _ ...any) error {
	return s.ListErr
}
func (s *StorageWithError) Raw(_ any, _ string, _ ...any) error {
	return s.RawErr
}
func (s *StorageWithError) Count(_ any, _ ...any) (int64, error) {
	return s.CountRes, s.CountErr
}
func (s *StorageWithError) Delete(_ any, _ ...any) error { return s.DeleteErr }
