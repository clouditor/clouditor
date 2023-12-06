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
	"errors"
	"reflect"
	"strings"
	"time"

	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Entry
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
	List() ([]voc.IsCloudResource, error)
	CloudServiceID() string
}

// Authorizer authorizes a Cloud service
type Authorizer interface {
	Authorize() (err error)
}

// typeRegistry
var typeRegistry = make(map[string]reflect.Type)

// TODO(oxisto): auto-generate them as part of the voc?
func init() {
	types := []any{
		voc.Application{},
		voc.Library{},
		voc.TranslationUnitDeclaration{},
		voc.CodeRepository{},
		voc.Function{},
		voc.VirtualMachine{},
		voc.ObjectStorageService{},
		voc.ObjectStorage{},
		voc.NetworkInterface{},
		voc.ResourceGroup{},
		voc.BlockStorage{},
		voc.DatabaseService{},
		voc.DatabaseStorage{},
		voc.FileStorageService{},
		voc.FileStorage{},
		voc.WebApp{},
	}
	for _, v := range types {
		t := reflect.TypeOf(v)
		typeRegistry[t.String()] = t
	}
}

// NewResource creates a new voc resource.
func NewResource(d Discoverer, ID voc.ResourceID, name string, creationTime *time.Time, location voc.GeoLocation, labels map[string]string, parent voc.ResourceID, typ []string, raw ...interface{}) *voc.Resource {
	rawString, err := voc.ToStringInterface(raw)
	if err != nil {
		log.Errorf("%v: %v", voc.ErrConvertingStructToString, err)
	}

	return &voc.Resource{
		ID:           ID,
		ServiceID:    d.CloudServiceID(),
		CreationTime: util.SafeTimestamp(creationTime),
		Name:         name,
		GeoLocation:  location,
		Type:         typ,
		Labels:       labels,
		Parent:       parent,
		Raw:          rawString,
	}
}

func (r *Resource) ToVocResource() (voc.IsCloudResource, error) {
	var (
		b   []byte
		err error
	)

	typ := strings.Split(r.ResourceType, ",")[0]

	var t, ok = typeRegistry["voc."+typ]
	if !ok {
		return nil, errors.New("invalid type")
	}
	var v = reflect.New(t).Interface().(voc.IsCloudResource)

	b, err = r.Properties.GetStructValue().MarshalJSON()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, v)
	return v, err
}

// GetCloudServiceId is a shortcut to implement CloudServiceRequest. It returns
// the cloud service ID of the inner object.
func (req *UpdateResourceRequest) GetCloudServiceId() string {
	return req.Resource.GetCloudServiceId()
}
