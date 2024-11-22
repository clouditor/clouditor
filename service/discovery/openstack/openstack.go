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
	StorageComponent = "storage"
	ComputeComponent = "compute"
	NetworkComponent = "network"
	RegionName       = "OS_REGION_NAME"
)

var (
	log *logrus.Entry

	ErrConversionProtobufToAuthOptions = errors.New("could not convert protobuf value to openstack.authOptions")
	ErrGettingAuthOptionsFromEnv       = errors.New("error getting auth options from environment")
	ErrCouldNotAuthenticate            = errors.New("could not authenticate to Azure")
	ErrGettingNextPage                 = errors.New("error getting next page")
	ErrNoCredentialsConfigured         = errors.New("no credentials were configured")
)

type openstackDiscovery struct {
	isAuthorized bool

	// sub  *armsubscription.Subscription
	// cred azcore.TokenCredential
	// clientOptions       arm.ClientOptions
	// clients  clients //TODO(all): Implement
	csID     string
	clients  clients
	provider *gophercloud.ProviderClient
	authOpts *gophercloud.AuthOptions
}

type AuthOptions struct {
	IdentityEndpoint string `json:"identityEndpoint"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	TenantName       string `json:"tenantName"`
	AllowReauth      bool   `json:"allowReauth"`
}

type clients struct {
	provider      *gophercloud.ProviderClient
	computeClient *gophercloud.ServiceClient
	networkClient *gophercloud.ServiceClient
	storageClient *gophercloud.ServiceClient
	authOpts      *gophercloud.AuthOptions
}

func (*openstackDiscovery) Name() string {
	return "OpenStack"
}

func (*openstackDiscovery) Description() string {
	return "Discovery OpenStack."
}

func (a *openstackDiscovery) CertificationTargetID() string {
	return a.csID
}

type DiscoveryOption func(a *openstackDiscovery)

func WithCertificationTargetID(csID string) DiscoveryOption {
	return func(a *openstackDiscovery) {
		a.csID = csID
	}
}

// WithAuthorizer is an option to set the authentication options
func WithAuthorizer(o gophercloud.AuthOptions) DiscoveryOption {
	return func(d *openstackDiscovery) {
		d.authOpts = util.Ref(o)
		// d.authOpts = &gophercloud.AuthOptions{
		// 	IdentityEndpoint: o.IdentityEndpoint, // "https://identityHost:portNumber/v2.0"
		// 	Username:         o.Username,
		// 	Password:         o.Password,
		// 	TenantName:       o.TenantName,
		// 	AllowReauth:      o.AllowReauth,
		// }
	}
}

func WithProvider(p *gophercloud.ProviderClient) DiscoveryOption {
	return func(d *openstackDiscovery) {
		d.provider = p
	}
}

func init() {
	log = logrus.WithField("component", "openstack-discovery")
}

func NewOpenstackDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &openstackDiscovery{
		csID: config.DefaultCertificationTargetID,
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

// authorize authorizes to Openstack and asserts the following clients
// * compute client
// * block storage client
func (d *openstackDiscovery) authorize() (err error) {
	if d.provider == nil {
		d.provider, err = openstack.AuthenticatedClient(context.Background(), *d.authOpts)
		if err != nil {
			return fmt.Errorf("error while authenticating: %w", err)
		}
	}

	// TODO(all): Move to compute and storage files?
	if d.clients.computeClient == nil {
		d.clients.computeClient, err = openstack.NewComputeV2(d.provider, gophercloud.EndpointOpts{
			Region: os.Getenv(RegionName),
		})
		if err != nil {
			return fmt.Errorf("could not create compute client: %w", err)
		}
	}

	if d.clients.networkClient == nil {
		d.clients.networkClient, err = openstack.NewNetworkV2(d.provider, gophercloud.EndpointOpts{
			Region: os.Getenv(RegionName),
		})
		if err != nil {
			return fmt.Errorf("could not create network client: %w", err)
		}
	}

	if d.clients.storageClient == nil {
		d.clients.storageClient, err = openstack.NewBlockStorageV2(d.provider, gophercloud.EndpointOpts{
			Region: os.Getenv(RegionName),
			Type:   "block-storage", // We have to use block-storage here, otherwise volumev3 is used as type and that does not work. volumev3 is not available in the service catalog for now. We have to wait until it is fixed, see: https://github.com/gophercloud/gophercloud/issues/3207
		})
		if err != nil {
			return fmt.Errorf("could not create block storage client: %w", err)
		}
	}

	return
}

func NewAuthorizer( /*value *structpb.Value*/ ) (gophercloud.AuthOptions, error) {
	// TODO(anatheka): Das ist gophercloud options vs. eigens definierte options
	//  Get auth options from environment
	ao, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		log.Error("error getting auth options from environment: %w", err)
	}
	return ao, nil

	// // Get AuthOpts from protobuf value
	// authOpts, err = toAuthOptions(value)
	// if err != nil {
	// 	return nil, ErrConversionProtobufToAuthOptions
	// }
	// return ao, nil
}

// TODO(anatheka): Do we need that anymore?
// // toAuthOptions converts the protobuf value to Openstack AuthOptions
// func toAuthOptions(v *structpb.Value) (authOpts *AuthOptions, err error) {
// 	// Get openstack auth opts from configuration
// 	value := v.GetStructValue().AsMap()

// 	if len(value) == 0 {
// 		return nil, fmt.Errorf("converting raw configuration to map is empty")
// 	}

// 	// First, we have to marshal the configuration map
// 	jsonbody, err := json.Marshal(value)
// 	if err != nil {
// 		err = fmt.Errorf("could not marshal configuration")
// 		return
// 	}

// 	// Then, we can store it back to the gophercloud.AuthOptions
// 	if err = json.Unmarshal(jsonbody, &authOpts); err != nil {
// 		err = fmt.Errorf("could not parse configuration: %w", err)
// 		return
// 	}

// 	return
// }

// List lists OpenStack servers (compute resources) and translates them into the Clouditor ontology
func (d *openstackDiscovery) List() (list []ontology.IsResource, err error) {

	// Discover network interfaces
	network, err := d.discoverNetworkInterfaces()
	if err != nil {
		return nil, fmt.Errorf("could not discover network interfaces: %w", err)
	}
	list = append(list, network...)

	// Discover servers
	servers, err := d.discoverServers()
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
