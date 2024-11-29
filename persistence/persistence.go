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

import (
	"errors"
	"strings"
)

var (
	ErrRecordNotFound         = errors.New("record not in the database")
	ErrConstraintFailed       = errors.New("constraint failed")
	ErrUniqueConstraintFailed = errors.New("unique constraint failed")
	ErrUnsupportedType        = errors.New("unsupported type")
	ErrDatabase               = errors.New("database error")
)

// Storage comprises a database interface
// TODO(all): I think with generics we could get sth. like `[T]Get(r T, id string) T` which is maybe more intuitive
// then changing the parameter? Maybe its only my impression
type Storage interface {

	// Create creates a new object and put it into the DB
	Create(r any) error

	// Save updates the record r (workaround with conds needed that e.g. user has no specific ID, and we can not touch
	// the generated (gRPC) code s.t. user.username has primary key annotation)
	Save(r any, conds ...any) error

	// Update updates the record with given id of the DB with non-zero values in r
	Update(r any, conds ...any) error

	// Get gets the record which meet the given conditions
	Get(r any, conds ...any) error

	// List lists all records in database which meet the (optionally) given conditions with a certain limit after an
	// offset. If no limit is desired, the value -1 can be specified. Optionally set orderBy (column) and asc (true =
	// ascending, false = descending) for ordering the results.
	// Whitelist the set of possible column names to avoid injections.
	List(r any, orderBy string, asc bool, offset int, limit int, conds ...any) error

	// Count counts the number of records which meet the (optionally) given conditions
	Count(r any, conds ...any) (int64, error)

	// Delete deletes the record with given id of the DB
	Delete(r any, conds ...any) error

	// Raw executes a raw SQL statement and stores the result in r
	Raw(r any, query string, args ...any) error
}

// BuildConds prepares the conds used in [Storage.List] out of arrays of query and args.
func BuildConds(query []string, args []any) (conds []any) {
	conds = append([]any{strings.Join(query, " AND ")}, args...)
	return
}
