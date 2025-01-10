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

// package openstack contains a Clouditor discoverer for OpenStack-based cloud environments.
package openstack

import (
	"context"
	"errors"
	"fmt"
	"os"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/pagination"
	"github.com/sirupsen/logrus"
)

const (
	RegionName = "OS_REGION_NAME"
)

var (
	log *logrus.Entry

	ErrGettingAuthOptionsFromEnv = errors.New("error getting auth options from environment")
)

type openstackDiscovery struct {
	ctID     string
	clients  clients
	authOpts *gophercloud.AuthOptions
}

type clients struct {
	provider       *gophercloud.ProviderClient
	identityClient *gophercloud.ServiceClient
	computeClient  *gophercloud.ServiceClient
	networkClient  *gophercloud.ServiceClient
	storageClient  *gophercloud.ServiceClient
}

func (*openstackDiscovery) Name() string {
	return "OpenStack"
}

func (*openstackDiscovery) Description() string {
	return "Discovery OpenStack."
}

func (d *openstackDiscovery) CertificationTargetID() string {
	return d.ctID
}

type DiscoveryOption func(d *openstackDiscovery)

func WithCertificationTargetID(ctID string) DiscoveryOption {
	return func(d *openstackDiscovery) {
		d.ctID = ctID
	}
}

// WithAuthorizer is an option to set the authentication options
func WithAuthorizer(o gophercloud.AuthOptions) DiscoveryOption {
	return func(d *openstackDiscovery) {
		d.authOpts = util.Ref(o)
	}
}

func init() {
	log = logrus.WithField("component", "openstack-discovery")
}

func NewOpenstackDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &openstackDiscovery{
		ctID: config.DefaultCertificationTargetID,
	}

	// Apply options
	for _, opt := range opts {
		opt(d)
	}

	// WithAuthorizer is mandatory, since it cannot be checked directly whether WithAuthorizer was passed, we check if authOpts is set before returning the discoverer
	if d.authOpts == nil {
		return nil
	}

	return d
}

// authorize authorizes to OpenStack and asserts the following clients
// * compute client
// * network client
// * block storage client
// * identity client
func (d *openstackDiscovery) authorize() (err error) {
	if d.clients.provider == nil {
		d.clients.provider, err = openstack.AuthenticatedClient(context.Background(), util.Deref(d.authOpts))
		if err != nil {
			return fmt.Errorf("error while authenticating: %w", err)
		}
	}

	// TODO(all): Move to compute and storage files?
	// Compute client
	if d.clients.computeClient == nil {
		d.clients.computeClient, err = openstack.NewComputeV2(d.clients.provider, gophercloud.EndpointOpts{
			Region: os.Getenv(RegionName),
		})
		if err != nil {
			return fmt.Errorf("could not create compute client: %w", err)
		}
	}

	// Network client
	if d.clients.networkClient == nil {
		d.clients.networkClient, err = openstack.NewNetworkV2(d.clients.provider, gophercloud.EndpointOpts{
			Region: os.Getenv(RegionName),
		})
		if err != nil {
			return fmt.Errorf("could not create network client: %w", err)
		}
	}

	// Storage client
	if d.clients.storageClient == nil {
		d.clients.storageClient, err = openstack.NewBlockStorageV2(d.clients.provider, gophercloud.EndpointOpts{
			Region: os.Getenv(RegionName),
			Type:   "block-storage", // We have to use block-storage here, otherwise volumev3 is used as type and that does not work. volumev3 is not available in the service catalog for now. We have to wait until it is fixed, see: https://github.com/gophercloud/gophercloud/issues/3207
		})
		if err != nil {
			return fmt.Errorf("could not create block storage client: %w", err)
		}
	}

	// Identity client
	if d.clients.identityClient == nil {
		d.clients.identityClient, err = openstack.NewIdentityV3(d.clients.provider, gophercloud.EndpointOpts{
			Region: os.Getenv(RegionName),
		})
		if err != nil {
			return fmt.Errorf("could not create identity client: %w", err)
		}
	}

	return
}

func NewAuthorizer() (gophercloud.AuthOptions, error) {
	ao, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		log.Error("%w: %w", ErrGettingAuthOptionsFromEnv, err)
	}
	return ao, err

}

// List discovers the following OpenStack resource types and translates them into the Clouditor ontology:
// * Domains
// * Projects
// * Network interfaces
// * Servers
// * Block storages
func (d *openstackDiscovery) List() (list []ontology.IsResource, err error) {
	// Discover domains resource
	domains, err := d.discoverDomains()
	if err != nil {
		return nil, fmt.Errorf("could not discover domains: %w", err)
	}
	list = append(list, domains...)

	// Discover project resources
	projects, err := d.discoverProjects()
	if err != nil {
		return nil, fmt.Errorf("could not discover projects: %w", err)
	}
	list = append(list, projects...)

	// Discover networks interfaces
	networks, err := d.discoverNetworkInterfaces()
	if err != nil {
		return nil, fmt.Errorf("could not discover network interfaces: %w", err)
	}
	list = append(list, networks...)

	// Discover servers
	servers, err := d.discoverServer()
	if err != nil {
		return nil, fmt.Errorf("could not discover servers: %w", err)
	}
	list = append(list, servers...)

	// Discover block storage
	storage, err := d.discoverBlockStorage()
	if err != nil {
		return nil, fmt.Errorf("could not discover block storage: %w", err)
	}
	list = append(list, storage...)

	return
}

type ClientFunc func() (*gophercloud.ServiceClient, error)
type ListFunc[O any] func(client *gophercloud.ServiceClient, opts O) pagination.Pager
type HandlerFunc[T any, R ontology.IsResource] func(in *T) (r R, err error)
type ExtractorFunc[T any] func(r pagination.Page) ([]T, error)

// genericList is a function leveraging type parameters that takes care of listing OpenStack
// resources using a
// - ClientFunc, which returns the needed client,
// - a ListFunc l, which returns paginated results,
// - a handler which converts them into an appropriate Clouditor ontology object,
// - an extractor that extracts the results into gophercloud specific objects and
// - optional options
func genericList[T any, O any, R ontology.IsResource](d *openstackDiscovery,
	clientGetter ClientFunc,
	l ListFunc[O],
	handler HandlerFunc[T, R],
	extractor ExtractorFunc[T],
	opts O,
) (list []ontology.IsResource, err error) {
	if err = d.authorize(); err != nil {
		return nil, fmt.Errorf("could not authorize openstack: %w", err)
	}

	client, err := clientGetter()
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	err = l(client, opts).EachPage(context.Background(), func(_ context.Context, p pagination.Page) (bool, error) {
		x, err := extractor(p)

		if err != nil {
			return false, fmt.Errorf("could not extract items from paginated result: %w", err)
		}

		for _, s := range x {
			r, err := handler(&s)
			if err != nil {
				return false, fmt.Errorf("could not convert into Clouditor ontology: %w", err)
			}

			log.Debugf("Adding resource %+v", s)

			list = append(list, r)
		}

		return true, nil
	})

	if err != nil {
		return nil, fmt.Errorf("could not list resources: %w", err)
	}

	return
}
