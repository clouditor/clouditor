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

package voc

type IsNetwork interface {
	IsResource
}

type NetworkResource struct {
	Resource
}

type HttpEndpoint struct {
	Resource // TODO(oxisto): should actually be a functionality, not a resource

	URL string `json:"url"`

	TransportEncryption *TransportEncryption `json:"transportEncryption"`
}

type NetworkService struct {
	NetworkResource

	IPs   []string
	Ports []int16
}

// LoadBalancer
type LoadBalancerResource struct {
	NetworkService

	AccessRestriction *AccessRestriction `json:"accessRestriction"`
	HttpEndpoints     []*HttpEndpoint    `json:"httpEndpoint"`
}

// Network Interface
type NetworkInterface struct {
	NetworkResource

	AccessRestriction *AccessRestriction `json:"accessRestriction"`
	AttachedTo        ResourceID         `json:"attachedTo"`
}

func (n *NetworkInterface) GetAccessRestriction() *AccessRestriction {
	return n.AccessRestriction
}
