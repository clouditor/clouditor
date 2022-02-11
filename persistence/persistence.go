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

package persistence

import "errors"

var (
	ErrRecordNotFound = errors.New("record not in the database")
)

// Storage comprises a database interface
// TODO(all): I think with generics we could get sth. like `[T]Get(r T, id string) T` which is maybe more intuitive
// then changing the parameter? Maybe its only my impression
type Storage interface {

	// Create creates a new object and put it into the DB
	Create(r interface{}) error

	// Get gets the record which meet the given conditions
	Get(r interface{}, conds ...interface{}) error

	// List lists all records in database which meet the (optionally) given conditions
	List(r interface{}, conds ...interface{}) error

	// Count counts the number of records which meet the (optionally) given conditions
	Count(r interface{}, conds ...interface{}) (int64, error)

	// Update updates the record with given id of the DB (only change non-zero values of r)
	Update(r interface{}, conds ...interface{}) error

	// Delete deletes the record with given id of the DB
	Delete(r interface{}, conds ...interface{}) error
}
