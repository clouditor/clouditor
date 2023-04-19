// Copyright 2021 Fraunhofer AISEC
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

package azure

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dataprotection/armdataprotection"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/security/armsecurity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"

	"github.com/sirupsen/logrus"

	"clouditor.io/clouditor/internal/util"
)

const (
	StorageComponent           = "storage"
	ComputeComponent           = "compute"
	NetworkComponent           = "network"
	DefenderStorageType        = "AccountStorages"
	DefenderVirtualMachineType = "VirtualMachines"
)

var (
	log *logrus.Entry

	ErrCouldNotAuthenticate     = errors.New("could not authenticate to Azure")
	ErrCouldNotGetSubscriptions = errors.New("could not get azure subscription")
	ErrNoCredentialsConfigured  = errors.New("no credentials were configured")
	ErrGettingNextPage          = errors.New("error getting next page")
)

type DiscoveryOption func(a *azureDiscovery)

func WithSender(sender policy.Transporter) DiscoveryOption {
	return func(a *azureDiscovery) {
		a.clientOptions.Transport = sender
	}
}

func WithAuthorizer(authorizer azcore.TokenCredential) DiscoveryOption {
	return func(a *azureDiscovery) {
		a.cred = authorizer
	}
}

func WithCloudServiceID(csID string) DiscoveryOption {
	return func(a *azureDiscovery) {
		a.csID = csID
	}
}

// WithResourceGroup is a [DiscoveryOption] that scopes the discovery to a specific resource group.
func WithResourceGroup(rg string) DiscoveryOption {
	return func(a *azureDiscovery) {
		a.rg = &rg
	}
}

func init() {
	log = logrus.WithField("component", "azure-discovery")
}

type azureDiscovery struct {
	isAuthorized bool

	sub  armsubscription.Subscription
	cred azcore.TokenCredential
	// rg optionally contains the name of a resource group. If this is not nil, all discovery calls will be scoped to the particular resource group.
	rg                  *string
	clientOptions       arm.ClientOptions
	discovererComponent string
	clients             clients
	csID                string
}

type clients struct {
	blobContainerClient     *armstorage.BlobContainersClient
	fileStorageClient       *armstorage.FileSharesClient
	accountsClient          *armstorage.AccountsClient
	networkInterfacesClient *armnetwork.InterfacesClient
	loadBalancerClient      *armnetwork.LoadBalancersClient
	functionsClient         *armappservice.WebAppsClient
	virtualMachinesClient   *armcompute.VirtualMachinesClient
	blockStorageClient      *armcompute.DisksClient
	diskEncSetClient        *armcompute.DiskEncryptionSetsClient
	dataProtectionClient    *armdataprotection.BackupPoliciesClient
	backupVaultClient       *armdataprotection.BackupVaultsClient
	backupInstancesClient   *armdataprotection.BackupInstancesClient
}

func (a *azureDiscovery) CloudServiceID() string {
	return a.csID
}

func (a *azureDiscovery) authorize() (err error) {
	if a.isAuthorized {
		return
	}

	if a.cred == nil {
		return ErrNoCredentialsConfigured
	}

	// Create new subscriptions client
	subClient, err := armsubscription.NewSubscriptionsClient(a.cred, &a.clientOptions)
	if err != nil {
		err = fmt.Errorf("could not get new subscription client: %w", err)
		return err
	}

	// Get subscriptions
	subPager := subClient.NewListPager(nil)
	subList := make([]*armsubscription.Subscription, 0)
	for subPager.More() {
		pageResponse, err := subPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("%s: %w", ErrCouldNotGetSubscriptions, err)
			log.Error(err)
			return err
		}
		subList = append(subList, pageResponse.ListResult.Value...)
	}

	// check if list of subscriptions is empty
	if len(subList) == 0 {
		err = errors.New("list of subscriptions is empty")
		return
	}

	// get first subscription
	a.sub = *subList[0]

	log.Infof("Azure %s discoverer uses %s as subscription", a.discovererComponent, *a.sub.SubscriptionID)

	a.isAuthorized = true

	return nil
}

// NewAuthorizer returns the Azure credential using one of the following authentication types (in the following order):
//
//	EnvironmentCredential
//	ManagedIdentityCredential
//	AzureCLICredential
func NewAuthorizer() (*azidentity.DefaultAzureCredential, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Errorf("%s: %+v", ErrCouldNotAuthenticate, err)
		return nil, fmt.Errorf("%s: %w", ErrCouldNotAuthenticate, err)
	}

	return cred, nil
}

type defenderProperties struct {
	monitoringLogDataEnabled bool
	securityAlertsEnabled    bool
}

// discoverDefender discovers Defender for X services and returns a map with the following properties for each defender type
// * monitoringLogDataEnabled
// * securityAlertsEnabled
// The property will be set to the individual resources, e.g., compute, storage in the corresponding discoverers
func (d *azureDiscovery) discoverDefender() (map[string]*defenderProperties, error) {
	var pricings = make(map[string]*defenderProperties)

	// Create new defender client
	defenderClient, err := armsecurity.NewPricingsClient(util.Deref(d.sub.SubscriptionID), d.cred, &d.clientOptions)
	if err != nil {
		err = fmt.Errorf("could not get new defender client: %w", err)
		return nil, err
	}

	// List all pricings to get the enabled Defender for X
	pricingsList, err := defenderClient.List(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not discover pricing")
	}

	for _, elem := range pricingsList.Value {
		if *elem.Properties.PricingTier == armsecurity.PricingTierFree {
			pricings[*elem.Name] = &defenderProperties{
				monitoringLogDataEnabled: false,
				securityAlertsEnabled:    false,
			}
		} else {
			pricings[*elem.Name] = &defenderProperties{
				monitoringLogDataEnabled: true,
				securityAlertsEnabled:    true,
			}
		}
	}

	return pricings, nil
}

// resourceGroupName returns the resource group name of a given Azure ID
func resourceGroupName(id string) string {
	return strings.Split(id, "/")[4]
}

// labels converts the resource tags to the vocabulary label
func labels(tags map[string]*string) map[string]string {
	l := make(map[string]string)

	for tag, i := range tags {
		l[tag] = util.Deref(i)
	}

	return l
}

// ClientCreateFunc is a type that describes a function to create a new Azure SDK client.
type ClientCreateFunc[T any] func(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*T, error)

// initClient creates an Azure client if not already exists
func initClient[T any](existingClient *T, d *azureDiscovery, fun ClientCreateFunc[T]) (client *T, err error) {
	if existingClient != nil {
		return existingClient, nil
	}

	client, err = fun(util.Deref(d.sub.SubscriptionID), d.cred, &d.clientOptions)
	if err != nil {
		err = fmt.Errorf("could not get %T client: %w", new(T), err)
		log.Debug(err)
		return nil, err
	}

	return
}

// listPager loops all values from a [runtime.Pager] object from the Azure SDK and issues a callback for each item. It
// takes the following arguments:
//   - d, an [azureDiscovery] struct,
//   - newListAllPager, a function that supplies a [runtime.Pager] listing all resources of a specific Azure client,
//   - newListByResourceGroupPager, a function that supplies a [runtime.Pager] listing all resources of a specific resource group,
//   - resToValues1, a function that takes the response from a single page of newListAllPager and returns its values,
//   - resToValues2, a function that takes the response from a single page of newListByResourceGroupPager and returns its values,
//   - callback, a function that is called for each item in every page.
//
// This function will then decide to use newListAllPager or newListByResourceGroupPager depending on whether a resource
// group scope is set in the [azureDiscovery] object.
//
// This function makes heavy use of the following type constraints (generics):
//   - O1, a type that represents an option argument to the newListAllPager function, e.g. *[armcompute.VirtualMachinesClientListAllOptions],
//   - R1, a type that represents the return type of the newListAllPager function, e.g. [armcompute.VirtualMachinesClientListAllResponse],
//   - O2, a type that represents an option argument to the newListByResourceGroupPager function, e.g. *[armcompute.VirtualMachinesClientListOptions],
//   - R1, a type that represents the return type of the newListAllPager function, e.g. [armcompute.VirtualMachinesClientListResponse],
//   - T, a type that represents the final resource that is supplied to the callback, e.g. *[armcompute.VirtualMachine].
func listPager[O1 any, R1 any, O2 any, R2 any, T any](
	d *azureDiscovery,
	newListAllPager func(options O1) *runtime.Pager[R1],
	newListByResourceGroupPager func(resourceGroupName string, options O2) *runtime.Pager[R2],
	allPagerResponseToValues func(res R1) []*T,
	allByResourceGroupPagerResponseToValues func(res R2) []*T,
	callback func(disk *T) error,
) error {
	// If the resource group is empty, we invoke the all-pager
	if d.rg == nil {
		pager := newListAllPager(*new(O1))
		// Invoke a callback for each page
		return allPages(pager, func(page R1) error {
			// Retrieve all resources of every page
			values := allPagerResponseToValues(page)
			for _, resource := range values {
				// Invoke the outer callback for each resource
				err := callback(resource)
				// We abort with the supplied error, if the callback issued an error
				if err != nil {
					return err
				}
			}

			return nil
		})
	} else {
		// Otherwise, we ivnoke the by-resource-group-pager
		pager := newListByResourceGroupPager(*d.rg, *new(O2))
		// Invoke a callback for each page
		return allPages(pager, func(page R2) error {
			// Retrieve all resources of every page
			values := allByResourceGroupPagerResponseToValues(page)
			for _, resource := range values {
				// Invoke the outer callback for each resource
				err := callback(resource)
				// We abort with the supplied error, if the callback issued an error
				if err != nil {
					return err
				}
			}

			return nil
		})
	}
}

// allPages loops through all pages of a [runtime.Pager] and issues a callback to each page.
func allPages[T any](pager *runtime.Pager[T], callback func(page T) error) error {
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			return fmt.Errorf("%s: %w", ErrGettingNextPage, err)
		}

		err = callback(page)
		if err != nil {
			return err
		}
	}

	return nil
}
