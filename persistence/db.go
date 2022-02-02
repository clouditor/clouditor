// Copyright 2016-2020 Fraunhofer AISEC
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

package persistence

import (
	"clouditor.io/clouditor/api/auth"
	"clouditor.io/clouditor/api/orchestrator"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var log *logrus.Entry

// TODO(lebogg): Lack for better term
type GormX struct {
	db *gorm.DB
}

func init() {
	log = logrus.WithField("component", "db")
}

func (g *GormX) Init(inMemory bool, host string, port int16) (err error) {
	if inMemory {
		if g.db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{}); err != nil {
			return err
		}

		log.Println("Using in-memory DB")
	} else {
		if g.db, err = gorm.Open(postgres.Open(fmt.Sprintf("postgres://postgres@%s:%d/postgres?sslmode=disable", host, port)), &gorm.Config{}); err != nil {
			return err
		}

		log.Printf("Using postgres DB @ %s", host)
	}

	if err = g.db.AutoMigrate(&auth.User{}); err != nil {
		return fmt.Errorf("error during auto-migration: %w", err)
	}

	if err = g.db.AutoMigrate(&orchestrator.CloudService{}); err != nil {
		return fmt.Errorf("error during auto-migration: %w", err)
	}

	return nil
}

func (g *GormX) Create(r interface{}) error {
	return g.db.Create(r).Error
}

func (g *GormX) Read(r interface{}, conds ...interface{}) (err error) {
	switch r.(type) {
	case *[]*orchestrator.CloudService:
		err = g.db.Find(r, conds).Error
		break
	case *orchestrator.CloudService:
		err = g.db.First(r, conds).Error
		break
	default:
		err = fmt.Errorf("unsupported type: %v (%T)", r, r)
	}
	return
}

func (g *GormX) Update(r interface{}) error {
	// g.db.Model(r).Count()
	return g.db.Save(r).Error
}

// Delete deletes record with given id. If no record was found, returns ErrRecordNotFound
func (g *GormX) Delete(r interface{}, id string) error {
	tx := g.db.Delete(r, "Id = ?", id)
	if err := tx.Error; err != nil {
		return err
	}
	// No record with given ID found
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// GetDatabase returns the database
func (g *GormX) GetDatabase() *gorm.DB {
	return g.db
}
