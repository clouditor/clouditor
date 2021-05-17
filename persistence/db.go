/*
 * Copyright 2016-2020 Fraunhofer AISEC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *           $$\                           $$\ $$\   $$\
 *           $$ |                          $$ |\__|  $$ |
 *  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
 * $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
 * $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
 * $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
 * \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
 *  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
 *
 * This file is part of Clouditor Community Edition.
 */

package persistence

import (
	"fmt"
	"math/rand"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/schema"

	"github.com/cayleygraph/quad"

	"clouditor.io/clouditor/api/auth"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/cayleygraph/quad/voc"
	_ "github.com/cayleygraph/quad/voc/core"
)

var log *logrus.Entry
var db *gorm.DB
var store *cayley.Handle
var schemaConfig *schema.Config

func init() {
	log = logrus.WithField("component", "db")
}

func InitDB(inMemory bool, host string, port int16) (err error) {
	if inMemory {
		if db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{}); err != nil {
			return err
		}

		log.Println("Using in-memory DB")
	} else {
		if db, err = gorm.Open(postgres.Open(fmt.Sprintf("postgres://postgres@%s:%d/postgres?sslmode=disable", host, port)), &gorm.Config{}); err != nil {
			return err
		}

		log.Printf("Using postgres DB @ %s", host)
	}

	db.AutoMigrate(&auth.User{})

	// experimental cayley stuff

	// Create a brand new graph
	store, err = cayley.NewMemoryGraph()
	if err != nil {
		log.Fatalln(err)
	}

	voc.RegisterPrefix("cloud:", "https://clouditor.io")

	schemaConfig = schema.NewConfig()
	// Override a function to generate IDs. Can be changed to generate UUIDs, for example.
	schemaConfig.GenerateID = func(_ interface{}) quad.Value {
		return quad.BNode(fmt.Sprintf("node%d", rand.Intn(1000)))
	}

	return nil
}

// GetDatabase returns the database
func GetDatabase() *gorm.DB {
	return db
}

func GetStore() *cayley.Handle {
	return store
}

func GetSchemaConfig() *schema.Config {
	return schemaConfig
}
