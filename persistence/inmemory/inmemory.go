package inmemory

import (
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/persistence/gorm"
	"gorm.io/gorm/logger"
)

// NewStorage creates a new in-memory storage. For now this uses the gorm provider with the gorm.WithInMemory. In the
// future we want to supply our own independent implementation. It automatically sets the maximum concurrent connections
// to 1 because the in-memory sqlite driver has problems with more than concurrent connection (see
// https://github.com/mattn/go-sqlite3/issues/511).
func NewStorage() (persistence.Storage, error) {
	return gorm.NewStorage(
		gorm.WithInMemory(),
		gorm.WithMaxOpenConns(1),
		gorm.WithLogger(logger.Default.LogMode(logger.Silent)))
}
