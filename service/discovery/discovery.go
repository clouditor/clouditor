// Copyright 2016-2020 Fraunhofer AISEC
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
	"context"
	"encoding/json"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/service/discovery/azure"
	"clouditor.io/clouditor/voc"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

var log *logrus.Entry

//go:generate protoc -I ../../proto -I ../../third_party discovery.proto --go_out=../.. --go-grpc_out=../.. --openapi_out=../../openapi/discovery

// Service is an implementation of the Clouditor Discovery service
type Service struct {
	discovery.UnimplementedDiscoveryServer
}

func init() {
	log = logrus.WithField("component", "discovery")
}

var resources map[string]voc.IsResource = make(map[string]voc.IsResource)

// Start starts discovery
func (s Service) Start(ctx context.Context, request *discovery.StartDiscoveryRequest) (response *discovery.StartDiscoveryResponse, err error) {
	response = &discovery.StartDiscoveryResponse{Successful: true}

	log.Infof("Starting discovery...")

	var discoverer discovery.Discoverer = azure.NewAzureStorageDiscovery()

	list, _ := discoverer.List()

	for _, v := range list {
		resources[string(v.GetID())] = v
	}

	return response, nil
}

func (s Service) Query(ctx context.Context, request *emptypb.Empty) (response *discovery.QueryResponse, err error) {
	var r []*structpb.Value

	for _, v := range resources {
		var s structpb.Value

		// this is probably not the fastest approach, but this
		// way, no extra libraries are needed and no extra struct tags
		// except `json` are required. there is also no significant
		// speed increase in marshaling the whole resource list, because
		// we first need to build it out of the map anyway
		b, _ := json.Marshal(v)
		json.Unmarshal(b, &s)
		r = append(r, &s)
	}

	return &discovery.QueryResponse{
		Result: &structpb.ListValue{Values: r},
	}, nil
}
