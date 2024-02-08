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

package discovery

import (
	"encoding/json"
	"fmt"
	reflect "reflect"
	"strings"

	"clouditor.io/clouditor/api/ontology"
	"google.golang.org/protobuf/proto"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

const (
	// DefaultCloudServiceID is the default service ID. Currently, our discoverers have no way to differentiate between different
	// services, but we need this feature in the future. This serves as a default to already prepare the necessary
	// structures for this feature.
	DefaultCloudServiceID   = "00000000-0000-0000-0000-000000000000"
	EvidenceCollectorToolId = "Clouditor Evidences Collection"
)

// Discoverer is a part of the discovery service that takes care of the actual discovering and translation into
// vocabulary objects.
type Discoverer interface {
	Name() string
	List() ([]ontology.IsResource, error)
	CloudServiceID() string
}

func (r *Resource) ToOntologyResource() (or ontology.IsResource, err error) {
	var (
		m  proto.Message
		ok bool
	)

	m, err = r.Properties.UnmarshalNew()
	if err != nil {
		return nil, err
	}

	if or, ok = m.(ontology.IsResource); ok {
		return or, nil
	}

	return
}

func Raw(raws ...any) string {
	var rawMap = make(map[string][]any)

	for _, raw := range raws {
		typ := reflect.TypeOf(raw).String()

		rawMap[typ] = append(rawMap[typ], raw)
	}

	b, _ := json.Marshal(rawMap)
	return string(b)
}

// ToDiscoveryResource converts a proto message that complies to the interface [ontology.IsResource] into a resource
// that can be persisted in our database ([*discovery.Resource]).
func ToDiscoveryResource(resource ontology.IsResource, csID string) (r *Resource, err error) {
	var (
		a *anypb.Any
	)

	// Convert our ontology resource into an Any proto message
	a, err = anypb.New(resource)
	if err != nil {
		return nil, fmt.Errorf("could not convert protobuf structure: %w", err)
	}

	// Build a resource struct. This will hold the latest sync state of the
	// resource for our storage layer.
	r = &Resource{
		Id:             string(resource.GetId()),
		ResourceType:   strings.Join(ontology.ResourceTypes(resource), ","),
		CloudServiceId: csID,
		Properties:     a,
	}

	return
}

// GetCloudServiceId is a shortcut to implement CloudServiceRequest. It returns
// the cloud service ID of the inner object.
func (req *UpdateResourceRequest) GetCloudServiceId() string {
	return req.Resource.GetCloudServiceId()
}
