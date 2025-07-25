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
	DomainID   = "OS_PROJECT_DOMAIN_ID"
	DomainName = "OS_USER_DOMAIN_NAME"
)

var (
	log *logrus.Entry

	ErrGettingAuthOptionsFromEnv = errors.New("error getting auth options from environment")
)

type openstackDiscovery struct {
	ctID     string
	clients  clients
	authOpts *gophercloud.AuthOptions
	region   string
	domain   *domain
	project  *project
}

type domain struct {
	domainID   string
	domainName string
}

type project struct {
	// It is not possible to add the OS_TENANT_ID or OS_TENANT_NAME. It results in an error: "Error authenticating with application credential: Application credentials cannot request a scope."
	projectID   string
	projectName string
}

type clients struct {
	provider       *gophercloud.ProviderClient
	identityClient *gophercloud.ServiceClient
	computeClient  *gophercloud.ServiceClient
	networkClient  *gophercloud.ServiceClient
	storageClient  *gophercloud.ServiceClient
	clusterClient  *gophercloud.ServiceClient
}

func (*openstackDiscovery) Name() string {
	return "OpenStack"
}

func (*openstackDiscovery) Description() string {
	return "Discovery OpenStack."
}

func (d *openstackDiscovery) TargetOfEvaluationID() string {
	return d.ctID
}

type DiscoveryOption func(d *openstackDiscovery)

func WithTargetOfEvaluationID(ctID string) DiscoveryOption {
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
	region := os.Getenv(RegionName)
	if region == "" {
		region = "unknown"
	}

	d := &openstackDiscovery{
		ctID:   config.DefaultTargetOfEvaluationID,
		region: os.Getenv(RegionName),
		domain: &domain{
			domainID:   os.Getenv(DomainID),
			domainName: os.Getenv(DomainName),
		},
		// Currently, the project ID cannot be specified as an environment variable in conjunction with application credentials.
		project: &project{},
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

	// Compute client
	if d.clients.computeClient == nil {
		d.clients.computeClient, err = openstack.NewComputeV2(d.clients.provider, gophercloud.EndpointOpts{
			Region: d.region,
		})
		if err != nil {
			return fmt.Errorf("could not create compute client: %w", err)
		}
	}

	// Network client
	if d.clients.networkClient == nil {
		d.clients.networkClient, err = openstack.NewNetworkV2(d.clients.provider, gophercloud.EndpointOpts{
			Region: d.region,
		})
		if err != nil {
			return fmt.Errorf("could not create network client: %w", err)
		}
	}

	// Storage client
	if d.clients.storageClient == nil {
		d.clients.storageClient, err = openstack.NewBlockStorageV3(d.clients.provider, gophercloud.EndpointOpts{
			Region: d.region,
		})
		if err != nil {
			return fmt.Errorf("could not create block storage client: %w", err)
		}
	}

	// Identity client
	if d.clients.identityClient == nil {
		d.clients.identityClient, err = openstack.NewIdentityV3(d.clients.provider, gophercloud.EndpointOpts{
			Region: d.region,
		})
		if err != nil {
			return fmt.Errorf("could not create identity client: %w", err)
		}
	}

	// Cluster client
	if d.clients.clusterClient == nil {
		d.clients.clusterClient, err = openstack.NewContainerInfraV1(d.clients.provider, gophercloud.EndpointOpts{
			Region: d.region,
		})
		if err != nil {
			return fmt.Errorf("could not create cluster client: %w", err)
		}
	}

	return
}

func NewAuthorizer() (gophercloud.AuthOptions, error) {
	ao, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		log.Error("%w: %w", ErrGettingAuthOptionsFromEnv, err)
	}

	ao.AllowReauth = true // Allow re-authentication if the token expires
	return ao, err

}

// List discovers the following OpenStack resource types and translates them into the Clouditor ontology:
// * Servers
// * Network interfaces
// * Block storages
// * Domains
// * Projects
func (d *openstackDiscovery) List() (list []ontology.IsResource, err error) {
	var (
		servers  []ontology.IsResource
		networks []ontology.IsResource
		storages []ontology.IsResource
		projects []ontology.IsResource
		domains  []ontology.IsResource
		clusters []ontology.IsResource
	)

	if err = d.authorize(); err != nil {
		return nil, fmt.Errorf("could not authorize openstack: %w", err)
	}

	// First, we need to discover the resources to obtain the domain and project ID. Domains and projects are discovered last, or they are set manually if discovery is not possible due to insufficient permissions. Currently, application credentials in OpenStack are always created for a specific project within a specific domain, making discovery essentially unnecessary. The code will be retained in case this changes in the future.

	// Discover servers
	servers, err = d.discoverServer()
	if err != nil {
		log.Errorf("could not discover servers: %v", err)
	}
	list = append(list, servers...)

	// Discover networks interfaces
	networks, err = d.discoverNetworkInterfaces()
	if err != nil {
		log.Errorf("could not discover network interfaces: %v", err)
	}
	list = append(list, networks...)

	// Discover block storage
	storages, err = d.discoverBlockStorage()
	if err != nil {
		log.Errorf("could not discover block storage: %v", err)
	}
	list = append(list, storages...)

	// Discover clusters
	clusters, err = d.discoverCluster()
	if err != nil {
		log.Errorf("could not discover clusters: %v", err)
	}
	list = append(list, clusters...)

	// Discover project resources
	projects, err = d.discoverProjects()
	if err != nil {
		log.Errorf("could not discover projects/tenants: %v", err)
	}
	list = append(list, projects...)

	// Discover domains resource
	domains, err = d.discoverDomains()
	if err != nil {
		log.Errorf("could not discover domains: %v", err)
	}
	list = append(list, domains...)

	return list, nil
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
	client, err := clientGetter()
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	err = l(client, opts).EachPage(context.Background(), func(_ context.Context, p pagination.Page) (bool, error) {
		x, err := extractor(p)

		if err != nil {
			return false, fmt.Errorf("could not extract items from paginated result: %w", err)
		}

		// Check if project/tenant ID is already stored
		if d.project.projectID == "" {
			d.setProjectInfo(x)
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
