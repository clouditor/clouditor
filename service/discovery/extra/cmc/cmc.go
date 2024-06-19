// Copyright 2024 Fraunhofer AISEC
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

package cmc

import (
	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/config"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

func init() {
	log = logrus.WithField("component", "cmc-discovery")
}

type cmcDiscovery struct {
	CmcAddr string // TODO: What is the CmcAddr we need here?

	// TODO(all): What is the domain for?
	domain string

	// serviceId is the id we are trying to get the attestation for
	serviceId string

	// CloudServiceID
	csID string
}

type DiscoveryOption func(a *cmcDiscovery)

func (*cmcDiscovery) Name() string {
	return "CMC Discovery"
}

func (*cmcDiscovery) Description() string {
	return "Discovery attestation reports from CMC"
}

func NewCMCDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &cmcDiscovery{
		csID: config.DefaultCloudServiceID,
	}

	// Apply options
	for _, opt := range opts {
		opt(d)
	}

	return d
}

func (a *cmcDiscovery) CloudServiceID() string {
	return a.csID
}

func (d *cmcDiscovery) List() (list []ontology.IsResource, err error) {
	log.Infof("Fetching attestation reports from CMC %s for resource %s", d.CmcAddr) //TODO: How do we know the resource for which we get the attestation report?

	return d.discoverReports()
}
