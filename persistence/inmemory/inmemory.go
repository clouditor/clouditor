package inmemory

import (
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/persistence/gorm"
)

// NewStorage creates a new in memory storage
// For now with GORM. Use own implementation in the future
func NewStorage() (persistence.Storage, error) {
	return gorm.NewStorage(gorm.WithInMemory())
}
