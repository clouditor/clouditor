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

package evidence

import (
	"context"
	"reflect"

	"clouditor.io/clouditor/v2/api/ontology"

	"google.golang.org/protobuf/proto"
)

type EvidenceHookFunc func(ctx context.Context, evidence *Evidence, err error)

func (req *StoreEvidenceRequest) GetPayload() proto.Message {
	return req.Evidence
}

func (ev *Evidence) GetOntologyResource() ontology.IsResource {
	var (
		resource ontology.IsResource
		ok       bool
	)

	if ev.Resource == nil || ev.Resource.Type == nil {
		return nil
	}

	// A little bit of dark Go magic
	typ := reflect.ValueOf(ev.Resource.Type).Elem()
	resource, ok = typ.Field(0).Interface().(ontology.IsResource)
	if !ok {
		return nil
	}

	return resource
}
