// Copyright 2016-2022 Fraunhofer AISEC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//           $$\                           $$\ $$\   $$\
//           $$ |                          $$ |\__|  $$ |
//  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
// $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
// $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
// $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
// \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
//  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
//
// This file is part of Clouditor Community Edition.

package gorm

import (
	"errors"
	"fmt"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/auth"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/persistence"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var log *logrus.Entry

type storage struct {
	db *gorm.DB
	// for options: (set default when not in opts)
	dialector gorm.Dialector
	config    gorm.Config

	// types contains all types that we need to auto-migrate into database tables
	types []any

	// maxConn is the maximum number of connections. 0 means unlimited.
	maxConn int
}

// DefaultTypes contains a list of internal types that need to be migrated by default
var DefaultTypes = []any{
	&auth.User{},
	&orchestrator.CloudService{},
	&assessment.MetricImplementation{},
	&assessment.Metric{},
	&orchestrator.Certificate{},
	&orchestrator.State{},
	&orchestrator.Requirement{},
}

// StorageOption is a functional option type to configure the GORM storage. E.g. WithInMemory or WithPostgres
type StorageOption func(*storage)

// WithInMemory is an option to configure Storage to use an in memory DB
func WithInMemory() StorageOption {
	return func(s *storage) {
		s.dialector = sqlite.Open(":memory:")
	}
}

// WithMaxOpenConns is an option to configure the maximum number of open connections
func WithMaxOpenConns(max int) StorageOption {
	return func(s *storage) {
		s.maxConn = max
	}
}

// WithPostgres is an option to configure Storage to use a Postgres DB
func WithPostgres(host string, port uint16, user string, pw string, db string, sslmode string) StorageOption {
	return func(s *storage) {
		s.dialector = postgres.Open(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", user, pw, host, port, db, sslmode))
	}
}

// WithLogger is an option to configure Storage to use a Logger
func WithLogger(logger logger.Interface) StorageOption {
	return func(s *storage) {
		s.config.Logger = logger
	}
}

// WithAdditionalAutoMigration is an option to add additional types to GORM's auto-migration.
func WithAdditionalAutoMigration(types ...any) StorageOption {
	return func(s *storage) {
		s.types = append(s.types, types...)
	}
}

func init() {
	log = logrus.WithField("component", "storage")
}

// NewStorage creates a new storage using GORM (which DB to use depends on the StorageOption)
func NewStorage(opts ...StorageOption) (s persistence.Storage, err error) {
	log.Println("Creating storage")
	// Create storage with default gorm config
	g := &storage{
		config: gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		},
		types: DefaultTypes,
	}

	// Add options and/or override default ones
	for _, o := range opts {
		o(g)
	}

	if g.dialector == nil {
		WithInMemory()(g)
	}

	g.db, err = gorm.Open(g.dialector, &g.config)
	if err != nil {
		return nil, err
	}

	if g.maxConn > 0 {
		sql, err := g.db.DB()
		if err != nil {
			return nil, fmt.Errorf("could not retrieve sql.DB: %v", err)
		}

		sql.SetMaxOpenConns(g.maxConn)
	}

	schema.RegisterSerializer("timestamppb", &TimestampSerializer{})
	schema.RegisterSerializer("anypb", &AnySerializer{})

	// After successful DB initialization, migrate the schema
	if err = g.db.AutoMigrate(g.types...); err != nil {
		err = fmt.Errorf("error during auto-migration: %w", err)
		return
	}

	s = g
	return
}

func (s *storage) Create(r any) error {
	return s.db.Create(r).Error
}

func (s *storage) Get(r any, conds ...any) (err error) {
	// Preload all associations for r being filled with all items (including relationships)
	err = s.db.Preload(clause.Associations).First(r, conds...).Error
	// if record is not found, use the error message defined in the persistence package
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = persistence.ErrRecordNotFound
	}
	return
}

func (s *storage) List(r any, orderBy string, asc bool, offset int, limit int, conds ...any) error {
	var query = s.db
	// Set default direction to "ascending"
	var orderDirection = "asc"

	if limit != -1 {
		query = s.db.Limit(limit)
	}
	// Set direction to "descending"
	if !asc {
		orderDirection = "desc"
	}
	orderStmt := orderBy + " " + orderDirection
	// No explicit ordering
	if orderBy == "" {
		orderStmt = ""
	}

	return query.Offset(offset).Preload(clause.Associations).Order(orderStmt).Find(r, conds...).Error
}

func (s *storage) Count(r any, conds ...any) (count int64, err error) {
	err = s.db.Model(r).Where(conds).Count(&count).Error
	return
}

func (s *storage) Save(r any, conds ...any) error {
	return s.db.Where(conds).Save(r).Error
}

// Update will update the record with non-zero fields. Note that to get the entire updated record you have to call Get
func (s *storage) Update(r any, query any, args ...any) error {
	return s.db.Model(r).Where(query, args).Updates(r).Error
}

// Delete deletes record with given id. If no record was found, returns ErrRecordNotFound
func (s *storage) Delete(r any, conds ...any) error {
	// Remove record r with given ID
	tx := s.db.Delete(r, conds...)
	if err := tx.Error; err != nil { // db error
		return err
	}
	// No record with given ID found
	if tx.RowsAffected == 0 {
		return persistence.ErrRecordNotFound
	}

	return nil
}
