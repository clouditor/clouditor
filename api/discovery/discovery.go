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
	"time"

	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
)

// DefaultCloudServiceID is the default service ID. Currently, our discoverers have no way to differentiate between different
// services, but we need this feature in the future. This serves as a default to already prepare the necessary
// structures for this feature.
const DefaultCloudServiceID = "00000000-0000-0000-0000-000000000000"

// Discoverer is a part of the discovery service that takes care of the actual discovering and translation into
// vocabulary objects.
type Discoverer interface {
	Name() string
	List() ([]voc.IsCloudResource, error)
	CloudServiceID() string
}

// Authorizer authorizes a Cloud service
type Authorizer interface {
	Authorize() (err error)
}

// NewResource creates a new resource.
func NewResource(d Discoverer, ID voc.ResourceID, name string, creationTime *time.Time, location voc.GeoLocation, labels map[string]string, typ []string) *voc.Resource {
	return &voc.Resource{
		ID:           ID,
		ServiceID:    d.CloudServiceID(),
		CreationTime: util.SafeTimestamp(creationTime),
		Name:         name,
		GeoLocation:  location,
		Type:         typ,
		Labels:       labels,
	}
}
