// Copyright 2022 Fraunhofer AISEC
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
	"clouditor.io/clouditor/api/auth"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/persistence"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var log *logrus.Entry

type Storage struct {
	db *gorm.DB
}

func init() {
	log = logrus.WithField("component", "storage")
}

// NewStorage creates a new storage using GORM
// TODO(lebogg): Maybe rename 'Storage' part in name. Have to see in usage with package name ect.
func NewStorage(inMemory bool, host string, port int16) (s *Storage, err error) {
	s = &Storage{}

	if inMemory {
		if s.db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{}); err != nil {
			return nil, err
		}
		log.Println("Using in-memory DB")
	} else {
		if s.db, err = gorm.Open(postgres.Open(fmt.Sprintf("postgres://postgres@%s:%d/postgres?sslmode=disable", host, port)), &gorm.Config{}); err != nil {
			return nil, err
		}

		log.Printf("Using postgres DB @ %s", host)
	}

	if err = s.db.AutoMigrate(&auth.User{}); err != nil {
		err = fmt.Errorf("error during auto-migration: %w", err)
	}

	if err = s.db.AutoMigrate(&orchestrator.CloudService{}); err != nil {
		err = fmt.Errorf("error during auto-migration: %w", err)
	}

	return
}

func (s *Storage) Create(r interface{}) error {
	return s.db.Create(r).Error
}

func (s *Storage) Get(r interface{}, conds ...interface{}) (err error) {
	err = s.db.First(r, conds).Error
	// if record is not found, use the error message defined in the persistence package
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = persistence.ErrRecordNotFound
	}
	return
}

func (s *Storage) List(r interface{}, conds ...interface{}) error {
	return s.db.Find(r, conds).Error
}

func (s *Storage) Count(r interface{}, conds ...interface{}) (count int64, err error) {
	// TODO(lebogg): Test if this method chain works!
	err = s.db.Model(r).Where(conds).Count(&count).Error
	return
}

func (s *Storage) Update(r interface{}, _ ...interface{}) error {
	// TODO(lebogg): Open discussion about update vs. save, i.e. only individual fields should be updates or not
	return s.db.Save(r).Error
}

// Delete deletes record with given id. If no record was found, returns ErrRecordNotFound
func (s *Storage) Delete(r interface{}, conds ...interface{}) error {
	// if id is empty remove all records -> currently used for testing.
	if len(conds) == 0 {
		return s.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(r).Error
	}
	// Remove record r with given ID
	tx := s.db.Delete(r, conds)
	if err := tx.Error; err != nil { // db error
		return err
	}
	// No record with given ID found
	if tx.RowsAffected == 0 {
		return persistence.ErrRecordNotFound
	}

	return nil
}

// GetDatabase returns the database
func (s *Storage) GetDatabase() *gorm.DB {
	return s.db
}

// Reset resets entire the database
func (s *Storage) Reset() (err error) {
	if err = s.Delete(&orchestrator.CloudService{}); err != nil {
		return
	}
	if err = s.Delete(&auth.User{}); err != nil {
		return
	}
	return
}
