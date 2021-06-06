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

// Package voc contains the vocabulary for Cloud resources and their properties
// that can be discovered using Clouditor
package voc

import (
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/structpb"
)

type IsResource interface {
	GetID() string
	GetName() string
	GetCreationTime() *time.Time
}

type Resource struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	CreationTime int64  `json:"creationTime"` // is set to 0 if no creation time is available
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

func ToStruct(r IsResource) (s *structpb.Value, err error) {
	s = new(structpb.Value)

	// this is probably not the fastest approach, but this
	// way, no extra libraries are needed and no extra struct tags
	// except `json` are required. there is also no significant
	// speed increase in marshaling the whole resource list, because
	// we first need to build it out of the map anyway
	b, _ := json.Marshal(r)
	if err = json.Unmarshal(b, &s); err != nil {
		return nil, fmt.Errorf("JSON unmarshal failed: %v", err)
	}

	return
}
