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

package csaf

import (
	"fmt"
	"net/http"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/csaf-poc/csaf_distribution/v3/csaf"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

func init() {
	log = logrus.WithField("component", "csaf-discovery")
}

type csafDiscovery struct {
	domain string
	csID   string
	client *http.Client
}

type DiscoveryOption func(a *csafDiscovery)

func WithProviderDomain(domain string) DiscoveryOption {
	return func(a *csafDiscovery) {
		a.domain = domain
	}
}

func WithCloudServiceID(csID string) DiscoveryOption {
	return func(a *csafDiscovery) {
		a.csID = csID
	}
}

func NewTrustedProviderDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &csafDiscovery{
		csID:   discovery.DefaultCloudServiceID,
		domain: "wid.cert-bund.de",
		client: http.DefaultClient,
	}

	// Apply options
	for _, opt := range opts {
		opt(d)
	}

	return d
}

func (*csafDiscovery) Name() string {
	return "CSAF Trusted Provider Discovery"
}

func (*csafDiscovery) Description() string {
	return "Discovery CSAF documents from a CSAF trusted provider"
}

func (d *csafDiscovery) CloudServiceID() string {
	return d.csID
}

func (d *csafDiscovery) List() (list []ontology.IsResource, err error) {
	log.Info("Fetching CSAF documents from provider")

	loader := csaf.NewProviderMetadataLoader(d.client)
	metadataFiles := loader.Enumerate(d.domain)
	_ = metadataFiles

	// TODO: actually discover evidences in future PR
	for _, pmd := range metadata {
		if !pmd.Valid() {
			return nil, fmt.Errorf("could not load provider-metadata.json from %s", d.domain)
		}
		log.Info("Found valid CSAF PMD file")

		// create resource
		r := &ontology.ServiceMetadataDocument{
			Id:           string(pmd.URL),
			CreationTime: timestamppb.Now(),
		}
		list = append(list, r)
	}

	if err != nil {
		return nil, err
	}

	return list, nil
}
