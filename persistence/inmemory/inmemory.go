package inmemory

import (
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/persistence/gorm"
)

// For now with GORM. Use own implementation in the future
func NewStorage() (persistence.Storage, error) {
	return gorm.NewStorage(gorm.WithInMemory())
}
