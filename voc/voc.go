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
	"errors"
	"fmt"
	"reflect"
	"time"

	"google.golang.org/protobuf/types/known/structpb"
)

var (
	ErrConvertingStructToString = errors.New("error converting struct to string")
)

type IsCloudResource interface {
	GetID() ResourceID
	GetServiceID() string
	SetServiceID(ID string)
	GetName() string
	GetType() []string
	HasType(string) bool
	GetCreationTime() *time.Time
	GetRaw() string
	Related() []string
}

type IsSecurityFeature interface {
	Type() string
}

type ResourceID string

// Resource file from Ontology currently not used. How do we merge this file with the 'Resource Ontology file'
type Resource struct {
	ID ResourceID `json:"id"`
	// ServiceID contains the ID of the cloud service to which this resource belongs. When creating new resources using
	// the NewResource function of the discovery API, this gets filled automatically.
	ServiceID    string `json:"serviceId"`
	Name         string `json:"name"`
	CreationTime int64  `json:"creationTime"` // is set to 0 if no creation time is available
	// The resource type. It is an array, because a type can be derived from another
	Type        []string          `json:"type"`
	GeoLocation GeoLocation       `json:"geoLocation"`
	Labels      map[string]string `json:"labels"`
	Raw         string            `json:"raw"`
	Parent      ResourceID        `json:"parent"`
}

func (r *Resource) GetID() ResourceID {
	return r.ID
}

func (r *Resource) GetServiceID() string {
	return r.ServiceID
}

func (r *Resource) SetServiceID(ID string) {
	r.ServiceID = ID
}

func (r *Resource) GetName() string {
	return r.Name
}

func (r *Resource) GetType() []string {
	return r.Type
}

// HasType checks whether the resource has the particular resourceType
func (r *Resource) HasType(resourceType string) (ok bool) {
	for _, value := range r.Type {
		if value == resourceType {
			ok = true
			break
		}
	}

	return
}

func (r *Resource) GetRaw() string {
	return r.Raw
}

func (r *Resource) GetCreationTime() *time.Time {
	t := time.Unix(r.CreationTime, 0)
	return &t
}

// ToStringInterface returns a string representation of the input
func ToStringInterface(r []interface{}) (s string, err error) {
	var (
		b          []byte
		rawInfoMap = make(map[string][]interface{})
	)

	if r == nil {
		return "", nil
	}

	for i := range r {
		typ := reflect.TypeOf(r[i]).String()

		rawInfoMap[typ] = append(rawInfoMap[typ], r[i])
	}

	if b, err = json.Marshal(rawInfoMap); err != nil {
		return "", fmt.Errorf("JSON marshal failed: %w", err)
	}

	return string(b), nil
}

func ToStruct(r IsCloudResource) (s *structpb.Value, err error) {
	var b []byte

	s = new(structpb.Value)

	// this is probably not the fastest approach, but this
	// way, no extra libraries are needed and no extra struct tags
	// except `json` are required. there is also no significant
	// speed increase in marshaling the whole resource list, because
	// we first need to build it out of the map anyway
	if b, err = json.Marshal(r); err != nil {
		return nil, fmt.Errorf("JSON marshal failed: %w", err)
	}
	if err = json.Unmarshal(b, &s); err != nil {
		return nil, fmt.Errorf("JSON unmarshal failed: %w", err)
	}

	return
}

type IsStorage interface {
	IsCloudResource
	IsAtRestEncryption
}

type IsAtRestEncryption interface {
	IsSecurityFeature
	atRestEncryption()
	IsEnabled() bool
}

func (*AtRestEncryption) atRestEncryption() {}
func (a *AtRestEncryption) IsEnabled() bool {
	return a.Enabled
}

type IsTransportEncryption interface {
	IsSecurityFeature
	transportEncryption()
	IsEnabled() bool
}

func (*TransportEncryption) transportEncryption() {}
func (a *TransportEncryption) IsEnabled() bool {
	return a.Enabled
}

type IsAuthorization interface {
	IsSecurityFeature
	authorization()
}

func (*Authorization) authorization() {}

type IsAuthenticity interface {
	IsSecurityFeature
	authenticity()
}

func (*Authenticity) authenticity() {}

type IsAccessRestriction interface {
	IsSecurityFeature
	accessRestriction()
}

func (*AccessRestriction) accessRestriction() {}

type HasHttpEndpoint interface {
	GetHttpEndpoint() *HttpEndpoint
}

type IsCompute interface {
	IsCloudResource
}

type IsNetwork interface {
	IsCloudResource
}
