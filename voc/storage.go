// Copyright 2021 Fraunhofer AISEC
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

package voc

import (
	"time"
)

type IsResource interface {
	GetID() string
	GetName() string
	GetCreationTime() *time.Time
}

type Resource struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	CreationTime int64  `json:"creationTime"`
}

func (r *Resource) GetID() string {
	return r.ID
}

func (r *Resource) GetName() string {
	return r.Name
}

func (r *Resource) GetCreationTime() *time.Time {
	t := time.Unix(r.CreationTime, 0)
	return &t
}

type HasAtRestEncryption interface {
	GetAtRestEncryption() *AtRestEncryption
}

type HasHttpEndpoint interface {
	GetHttpEndpoint() *HttpEndpoint
}

type IsStorage interface {
	IsResource

	HasAtRestEncryption
}

type StorageResource struct {
	Resource

	AtRestEncryption *AtRestEncryption `json:"atRestEncryption"`
}

func (s *StorageResource) GetAtRestEncryption() *AtRestEncryption {
	return s.AtRestEncryption
}

type IsObjectStorage interface {
	IsStorage
	HasHttpEndpoint
}

type ObjectStorageResource struct {
	StorageResource

	HttpEndpoint *HttpEndpoint `json:"httpEndpoint"`
}

type BlockStorageResource struct {
	StorageResource
}
