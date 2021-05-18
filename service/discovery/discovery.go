/*
 * Copyright 2016-2020 Fraunhofer AISEC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *           $$\                           $$\ $$\   $$\
 *           $$ |                          $$ |\__|  $$ |
 *  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
 * $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
 * $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
 * $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
 * \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
 *  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
 *
 * This file is part of Clouditor Community Edition.
 */

package discovery

import (
	"context"

	"clouditor.io/clouditor/api/discovery"
	"github.com/sirupsen/logrus"
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

// Start starts discovery
func (s Service) Start(ctx context.Context, request *discovery.StartDiscoveryRequest) (response *discovery.StartDiscoveryResponse, err error) {
	response = &discovery.StartDiscoveryResponse{Successful: true}

	var discovery StorageDiscoverer = &azureStorageDiscovery{}
	discovery.List()

	return response, nil
}
