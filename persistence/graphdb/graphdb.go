// Copyright 2024 Fraunhofer AISEC
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

package graphdb

import (
	"clouditor.io/clouditor/v2/persistence"
	"fmt"
)

import "github.com/aerospike/aerospike-client-go/v7"

type storage struct {
	db *aerospike.Client

	config config
}

type config struct {
	host string
	port int
}

func NewStorage(opts ...StorageOption) (p persistence.Storage, err error) {
	s := &storage{}

	for _, o := range opts {
		o(s)
	}
	s.db, err = aerospike.NewClient(s.config.host, s.config.port)
	if err != nil {
		err = fmt.Errorf("could not create new aerospike client: %v", err)
		return
	}

	p = s
	return
}

func WithConfig(host string, port int) StorageOption {
	return func(s *storage) {
		s.config.host = host
		s.config.port = port
	}
}

type StorageOption func(storage2 *storage)

func (s storage) Create(r any) error {
	//TODO implement me
	panic("implement me")
}

func (s storage) Save(r any, conds ...any) error {
	//TODO implement me
	panic("implement me")
}

func (s storage) Update(r any, conds ...any) error {
	//TODO implement me
	panic("implement me")
}

func (s storage) Get(r any, conds ...any) error {
	//TODO implement me
	panic("implement me")
}

func (s storage) List(r any, orderBy string, asc bool, offset int, limit int, conds ...any) error {
	//TODO implement me
	panic("implement me")
}

func (s storage) Count(r any, conds ...any) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s storage) Delete(r any, conds ...any) error {
	//TODO implement me
	panic("implement me")
}

func (s storage) Raw(r any, query string, args ...any) error {
	//TODO implement me
	panic("implement me")
}
