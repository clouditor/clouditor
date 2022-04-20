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
	"gorm.io/gorm/logger"
)

var log *logrus.Entry

type storage struct {
	db *gorm.DB
	// for options: (set default when not in opts)
	dialector gorm.Dialector
	config    gorm.Config
}

// StorageOption is a functional option type to configure the GORM storage. E.g. WithInMemory or WithPostgres
type StorageOption func(*storage)

// WithInMemory is an option to configure Storage to use an in memory DB
func WithInMemory() StorageOption {
	return func(s *storage) {
		s.dialector = sqlite.Open(":memory:")
	}
}

// WithPostgres is an option to configure Storage to use a Postgres DB
func WithPostgres(host string, port int16) StorageOption {
	return func(s *storage) {
		s.dialector = postgres.Open(fmt.Sprintf("postgres://postgres@%s:%d/postgres?sslmode=disable", host, port))
	}
}

func init() {
	log = logrus.WithField("component", "storage")
}

// NewStorage creates a new storage using GORM (which DB to use depends on the StorageOption)
func NewStorage(opts ...StorageOption) (s persistence.Storage, err error) {
	g := &storage{
		config: gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	}

	// Init storage
	log.Println("Creating storage")
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

	// After successful DB initialization, migrate the schema
	var types = []interface{}{
		&auth.User{},
		&orchestrator.CloudService{},
		&assessment.MetricImplementation{},
	}

	if err = g.db.AutoMigrate(types...); err != nil {
		err = fmt.Errorf("error during auto-migration: %w", err)
		return
	}

	s = g
	return
}

func (s *storage) Create(r interface{}) error {
	return s.db.Create(r).Error
}

func (s *storage) Get(r interface{}, conds ...interface{}) (err error) {
	err = s.db.First(r, conds...).Error
	// if record is not found, use the error message defined in the persistence package
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = persistence.ErrRecordNotFound
	}
	return
}

func (s *storage) List(r interface{}, conds ...interface{}) error {
	return s.db.Find(r, conds...).Error
}

func (s *storage) Count(r interface{}, conds ...interface{}) (count int64, err error) {
	err = s.db.Model(r).Where(conds).Count(&count).Error
	return
}

func (s *storage) Save(r interface{}, conds ...interface{}) error {
	return s.db.Where(conds).Save(r).Error
}

// Update will update the record with non-zero fields. Note that to get the entire updated record you have to call Get
func (s *storage) Update(r interface{}, query interface{}, args ...interface{}) error {
	return s.db.Model(r).Where(query, args).Updates(r).Error
}

// Delete deletes record with given id. If no record was found, returns ErrRecordNotFound
func (s *storage) Delete(r interface{}, conds ...interface{}) error {
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
