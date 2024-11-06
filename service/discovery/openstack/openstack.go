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
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/config"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/structpb"
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
	provider *gophercloud.ProviderClient
	compute  *gophercloud.ServiceClient
	storage  *gophercloud.ServiceClient
	authOpts *gophercloud.AuthOptions
}

type AuthOptions struct {
	IdentityEndpoint string `json:"identityEndpoint"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	TenantName       string `json:"tenantName"`
	AllowReauth      bool   `json:"allowReauth"`
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
func WithAuthorizer(o *AuthOptions) DiscoveryOption {
	return func(d *openstackDiscovery) {
		d.authOpts = &gophercloud.AuthOptions{
			IdentityEndpoint: o.IdentityEndpoint, // "https://identityHost:portNumber/v2.0"
			Username:         o.Username,
			Password:         o.Password,
			TenantName:       o.TenantName,
			AllowReauth:      o.AllowReauth,
		}
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
		d.provider, err = openstack.AuthenticatedClient(*d.authOpts)
		if err != nil {
			return fmt.Errorf("error while authenticating: %w", err)
		}
	}

	if d.compute == nil {
		d.compute, err = openstack.NewComputeV2(d.provider, gophercloud.EndpointOpts{
			Region: os.Getenv(RegionName),
		})

		if err != nil {
			return fmt.Errorf("could not create compute client: %w", err)
		}
	}

	if d.storage == nil {
		d.storage, err = openstack.NewBlockStorageV3(d.provider, gophercloud.EndpointOpts{
			Region: os.Getenv(RegionName),
		})

		if err != nil {
			return fmt.Errorf("could not create block storage client: %w", err)
		}
	}

	return
}

func NewAuthorizer(value *structpb.Value) (*AuthOptions, error) {
	// Get AuthOpts from protobuf value
	authOpts, err := toAuthOptions(value)
	if err != nil {
		return nil, ErrConversionProtobufToAuthOptions
	}
	return authOpts, nil
}

// toAuthOptions converts the protobuf value to Openstack AuthOptions
func toAuthOptions(v *structpb.Value) (authOpts *AuthOptions, err error) {
	// Get openstack auth opts from configuration
	value := v.GetStructValue().AsMap()

	if value == nil {
		return nil, fmt.Errorf("converting raw configuration to map is nil")
	} else if len(value) == 0 {
		return nil, fmt.Errorf("converting raw configuration to map is empty")
	}

	// First, we have to marshal the configuration map
	jsonbody, err := json.Marshal(value)
	if err != nil {
		err = fmt.Errorf("could not marshal configuration")
		return
	}

	// Then, we can store it back to the gophercloud.AuthOptions
	if err = json.Unmarshal(jsonbody, &authOpts); err != nil {
		err = fmt.Errorf("could not parse configuration: %w", err)
		return
	}

	return
}

type ClientFunc func() (*gophercloud.ServiceClient, error)
type ListFunc[O any] func(client *gophercloud.ServiceClient, opts O) pagination.Pager
type HandlerFunc[T any, R ontology.IsResource] func(in *T) (r R, err error)
type ExtractorFunc[T any] func(r pagination.Page) ([]T, error)

// genericList is a function leveraging type parameters that takes care of listing OpenStack
// resources using a ClientFunc, which returns the needed client, a ListFunc l, which returns paginated results,
// an extractor that extracts the results into gophercloud specific objects and a handler which converts them
// into an appropriate Clouditor vocabulary object.
func genericList[T any, O any, R ontology.IsResource](d *openstackDiscovery, clientGetter ClientFunc,
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

	err = l(client, opts).EachPage(func(p pagination.Page) (bool, error) {
		x, err := extractor(p)

		if err != nil {
			return false, fmt.Errorf("could not extract items from paginated result: %w", err)
		}

		for _, s := range x {
			r, err := handler(&s)
			if err != nil {
				return false, fmt.Errorf("could not convert into Clouditor vocabulary: %w", err)
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
