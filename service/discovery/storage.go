// Copyright 2021 Fraunhofer AISEC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package discovery

import "time"

type Resource interface {
	ID() string
	Name() string
	CreationTime() *time.Time
}

type resource struct {
	id           string
	name         string
	creationTime *time.Time
}

func (r *resource) ID() string {
	return r.id
}

func (r *resource) Name() string {
	return r.name
}

func (r *resource) CreationTime() *time.Time {
	return r.creationTime
}

type Storage interface {
	Resource

	AtRestEncryption() *AtRestEncryption
}

type storage struct {
	resource

	atRestEncryption *AtRestEncryption
}

func (s *storage) AtRestEncryption() *AtRestEncryption {
	return s.atRestEncryption
}

type ObjectStorage interface {
	Storage

	HttpEndpoint() *HttpEndpoint
}

type objectStorage struct {
	storage

	httpEndpoint *HttpEndpoint
}

func (s *objectStorage) HttpEndpoint() *HttpEndpoint {
	return s.httpEndpoint
}

type BlockStorage interface {
	Storage
}

/*type blockStorage struct {
	storage
}*/

type StorageDiscoverer interface {
	List() ([]Storage, error)
}
