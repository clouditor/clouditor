// Copyright 2025 Fraunhofer AISEC
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

// package ionos contains a Clouditor discoverer for IONOS Cloud environments.
package ionos

import (
	"errors"
	"fmt"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/config"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Entry

	ErrCouldNotAuthenticate     = errors.New("could not authenticate to Azure")
	ErrCouldNotGetSubscriptions = errors.New("could not get azure subscription")
	ErrGettingNextPage          = errors.New("could not get next page")
	ErrNoCredentialsConfigured  = errors.New("no credentials were configured")
	ErrSubscriptionNotFound     = errors.New("SubscriptionNotFound")
	ErrVaultInstanceIsEmpty     = errors.New("vault and/or instance is nil")
)

func (*ionosDiscovery) Name() string {
	return "IONOS Cloud"
}

func (*ionosDiscovery) Description() string {
	return "Discovery IONOS Cloud."
}

type DiscoveryOption func(d *ionosDiscovery)

func WithAuthorizer(config *ionoscloud.Configuration) DiscoveryOption {
	return func(d *ionosDiscovery) {
		d.authConfig = config
	}
}

func WithTargetOfEvaluationID(ctID string) DiscoveryOption {
	return func(a *ionosDiscovery) {
		a.ctID = ctID
	}
}

// WithResourceGroup is a [DiscoveryOption] that scopes the discovery to a specific resource group.
func WithResourceGroup(rg string) DiscoveryOption {
	return func(d *ionosDiscovery) {
		d.rg = &rg
	}
}

func init() {
	log = logrus.WithField("component", "azure-discovery")
}

type ionosDiscovery struct {
	authConfig *ionoscloud.Configuration // authConfig contains the IONOS Cloud configuration, which is used to authenticate against the IONOS Cloud API
	// rg optionally contains the name of a resource group. If this is not nil, all discovery calls will be scoped to the particular resource group.
	rg *string
	// discovererComponent string
	clients clients
	ctID    string
}

type clients struct {
	computeClient *ionoscloud.APIClient
}

func NewIonosDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &ionosDiscovery{
		ctID: config.DefaultTargetOfEvaluationID,
	}

	// Apply options
	for _, opt := range opts {
		opt(d)
	}

	return d
}

// List discovers the following Azure resources types:
// - Storage resource
// - Compute resource
// - Network resource
func (d *ionosDiscovery) List() (list []ontology.IsResource, err error) {
	if err = d.authorize(); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrCouldNotAuthenticate, err)
	}

	// Discover storage resources
	log.Info("Discover IONOS Cloud storage resources...")

	// Discover storage accounts
	// storageAccounts, err := d.discoverStorage()
	// if err != nil {
	// 	return nil, fmt.Errorf("could not discover storage accounts: %w", err)
	// }
	// list = append(list, storageAccounts...)

	// Discover compute resources
	log.Info("Discover IONOS Cloud compute resources...")

	// Discover block storage
	// storage, err := d.discoverBlockStorages()
	// if err != nil {
	// 	return nil, fmt.Errorf("could not discover block storage: %w", err)
	// }
	// list = append(list, storage...)

	// Discover virtual machines
	virtualMachines, err := d.discoverServer()
	if err != nil {
		return nil, fmt.Errorf("could not discover virtual machines: %w", err)
	}
	list = append(list, virtualMachines...)

	// Discover network resources
	log.Info("Discover IONOS Cloud network resources...")

	// Discover network interfaces
	// networkInterfaces, err := d.discoverNetworkInterfaces()
	// if err != nil {
	// 	return nil, fmt.Errorf("could not discover network interfaces: %w", err)
	// }
	// list = append(list, networkInterfaces...)

	return list, nil
}

func (d *ionosDiscovery) TargetOfEvaluationID() string {
	return d.ctID
}

func (d *ionosDiscovery) authorize() (err error) {
	if d.clients.computeClient == nil {
		d.clients.computeClient = ionoscloud.NewAPIClient(d.authConfig)
	}

	return nil
}

// NewAuthorizer returns the IONOS Cloud configuration
func NewAuthorizer() (*ionoscloud.Configuration, error) {
	// authClient := auth.NewAPIClient(auth.NewConfigurationFromEnv())
	// jwt, _, err := authClient.TokensApi.TokensGenerate(context.Background()).Execute()
	// if err != nil {
	// 	return nil, nil, fmt.Errorf("error occurred while generating token (%w)", err)
	// }
	// if !jwt.HasToken() {
	// 	return nil, nil, errors.New("could not generate token")
	// }

	config := ionoscloud.NewConfigurationFromEnv()
	// sharedConfiguration := shared.NewConfigurationFromEnv()
	// if sharedConfiguration == nil {
	// 	return nil, fmt.Errorf("%w: %s", ErrNoCredentialsConfigured, "IONOS Cloud credentials are not configured in the environment")
	// }

	return config, nil
}
