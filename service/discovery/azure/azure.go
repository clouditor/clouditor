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
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dataprotection/armdataprotection"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/security/armsecurity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"

	"github.com/sirupsen/logrus"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/internal/constants"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
)

const (
	StorageComponent = "storage"
	ComputeComponent = "compute"
	NetworkComponent = "network"

	DefenderStorageType        = "StorageAccounts"
	DefenderVirtualMachineType = "VirtualMachines"

	DataSourceTypeDisc                 = "Microsoft.Compute/disks"
	DataSourceTypeStorageAccountObject = "Microsoft.Storage/storageAccounts/blobServices"

	Duration30Days = time.Duration(30 * time.Hour * 24)
	Duration7Days  = time.Duration(7 * time.Hour * 24)

	AES256 = "AES256"

	RetentionPeriod90Days = 90
)

var (
	log *logrus.Entry

	ErrCouldNotAuthenticate     = errors.New("could not authenticate to Azure")
	ErrCouldNotGetSubscriptions = errors.New("could not get azure subscription")
	ErrNoCredentialsConfigured  = errors.New("no credentials were configured")
	ErrGettingNextPage          = errors.New("error getting next page")
	ErrVaultInstanceIsEmpty     = errors.New("vault and/or instance is nil")
)

// TODO (anatheka): Merge discoverers:
// - fix backupStorage
// - fix tests

// TODO (anatheka): Discoveres:
// NewAzureResourceGroupDiscovery
// NewAzureComputeDiscovery
// NewAzureStorageDiscovery
// NewAzureNetworkDiscovery

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

	sub  *armsubscription.Subscription
	cred azcore.TokenCredential
	// rg optionally contains the name of a resource group. If this is not nil, all discovery calls will be scoped to the particular resource group.
	rg                  *string
	clientOptions       arm.ClientOptions
	discovererComponent string
	clients             clients
	csID                string
	backupMap           map[string]*backup
	defenderProperties  map[string]*defenderProperties
}

type backup struct {
	backup         map[string][]*voc.Backup
	backupStorages []voc.IsCloudResource
}

type clients struct {
	// Storage
	blobContainerClient *armstorage.BlobContainersClient
	fileStorageClient   *armstorage.FileSharesClient
	accountsClient      *armstorage.AccountsClient

	// DB
	databasesClient        *armsql.DatabasesClient
	sqlServersClient       *armsql.ServersClient
	threatProtectionClient *armsql.DatabaseAdvancedThreatProtectionSettingsClient
	cosmosDBClient         *armcosmos.DatabaseAccountsClient

	// Network
	networkInterfacesClient     *armnetwork.InterfacesClient
	loadBalancerClient          *armnetwork.LoadBalancersClient
	applicationGatewayClient    *armnetwork.ApplicationGatewaysClient
	networkSecurityGroupsClient *armnetwork.SecurityGroupsClient

	// AppService
	sitesClient *armappservice.WebAppsClient

	// Compute
	virtualMachinesClient *armcompute.VirtualMachinesClient
	blockStorageClient    *armcompute.DisksClient
	diskEncSetClient      *armcompute.DiskEncryptionSetsClient

	// Security
	defenderClient *armsecurity.PricingsClient

	// Data protection
	backupPoliciesClient  *armdataprotection.BackupPoliciesClient
	backupVaultClient     *armdataprotection.BackupVaultsClient
	backupInstancesClient *armdataprotection.BackupInstancesClient

	// Resource groups
	rgClient *armresources.ResourceGroupsClient
}

func NewAzureDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &azureDiscovery{
		csID:               discovery.DefaultCloudServiceID,
		backupMap:          make(map[string]*backup),
		defenderProperties: make(map[string]*defenderProperties),
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
// - ResourceGroup resource
func (d *azureDiscovery) List() (list []voc.IsCloudResource, err error) {
	if err = d.authorize(); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrCouldNotAuthenticate, err)
	}

	// Discover resource group resources
	log.Info("Discover Azure resource group resources...")
	rg, err := d.discoverResourceGroups()
	if err != nil {
		return nil, fmt.Errorf("could not discover resource groups: %w", err)
	}
	list = append(list, rg...)

	// Discover storage resources
	log.Info("Discover Azure storage resources...")

	// Discover Defender for X properties to add it to the required resource properties
	d.defenderProperties, err = d.discoverDefender()
	if err != nil {
		return nil, fmt.Errorf("could not discover Defender for X: %w", err)
	}

	// Discover storage accounts
	storageAccounts, err := d.discoverStorageAccounts()
	if err != nil {
		return nil, fmt.Errorf("could not discover storage accounts: %w", err)
	}
	list = append(list, storageAccounts...)

	// Discover sql databases
	dbs, err := d.discoverSqlServers()
	if err != nil {
		return nil, fmt.Errorf("could not discover sql databases: %w", err)
	}
	list = append(list, dbs...)

	// Discover Cosmos DB
	cosmosDB, err := d.discoverCosmosDB()
	if err != nil {
		return nil, fmt.Errorf("could not discover cosmos db accounts: %w", err)
	}
	list = append(list, cosmosDB...)

	// Discover compute resources
	log.Info("Discover Azure compute resources...")

	// Discover backup vaults
	err = d.discoverBackupVaults()
	if err != nil {
		log.Errorf("could not discover backup vaults: %v", err)
	}

	// Discover block storage
	storage, err := d.discoverBlockStorages()
	if err != nil {
		return nil, fmt.Errorf("could not discover block storage: %w", err)
	}
	list = append(list, storage...)

	// Add backup block storages
	if d.backupMap[DataSourceTypeDisc] != nil && d.backupMap[DataSourceTypeDisc].backupStorages != nil {
		list = append(list, d.backupMap[DataSourceTypeDisc].backupStorages...)
	}

	// Discover virtual machines
	virtualMachines, err := d.discoverVirtualMachines()
	if err != nil {
		return nil, fmt.Errorf("could not discover virtual machines: %w", err)
	}
	list = append(list, virtualMachines...)

	// Discover functions and web apps
	resources, err := d.discoverFunctionsWebApps()
	if err != nil {
		return nil, fmt.Errorf("could not discover functions: %w", err)
	}

	list = append(list, resources...)

	// Discover network resources
	log.Info("Discover Azure network resources...")

	// Discover network interfaces
	networkInterfaces, err := d.discoverNetworkInterfaces()
	if err != nil {
		return nil, fmt.Errorf("could not discover network interfaces: %w", err)
	}
	list = append(list, networkInterfaces...)

	// Discover Load Balancer
	loadBalancer, err := d.discoverLoadBalancer()
	if err != nil {
		return list, fmt.Errorf("could not discover load balancer: %w", err)
	}
	list = append(list, loadBalancer...)

	// Discover Application Gateway
	ag, err := d.discoverApplicationGateway()
	if err != nil {
		return list, fmt.Errorf("could not discover application gateways: %w", err)
	}
	list = append(list, ag...)

	return list, nil
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
	a.sub = subList[0]

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

	// initialize defender client
	if err := d.initDefenderClient(); err != nil {
		return nil, fmt.Errorf("could not initialize defender client: %w", err)
	}

	if d.clients.defenderClient == nil {
		return nil, errors.New("defenderClient not set")
	}

	// List all pricings to get the enabled Defender for X
	pricingsList, err := d.clients.defenderClient.List(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not discover pricings")
	}

	for _, elem := range pricingsList.Value {
		if util.Deref(elem.Properties.PricingTier) == armsecurity.PricingTierFree {
			pricings[util.Deref(elem.Name)] = &defenderProperties{
				monitoringLogDataEnabled: false,
				securityAlertsEnabled:    false,
			}
		} else {
			pricings[util.Deref(elem.Name)] = &defenderProperties{
				monitoringLogDataEnabled: true,
				securityAlertsEnabled:    true,
			}
		}
	}

	return pricings, nil
}

// discoverBackupVaults receives all backup vaults in the subscription.
// Since the backups for storage and compute are discovered together, the discovery is performed here and results are stored in the azureDiscovery receiver.
func (d *azureDiscovery) discoverBackupVaults() error {

	if d.backupMap != nil && len(d.backupMap) > 0 {
		log.Debug("Backup Vaults already discovered.")
		return nil
	}

	// initialize backup vaults client
	if err := d.initBackupVaultsClient(); err != nil {
		return err
	}

	// initialize backup instances client
	if err := d.initBackupInstancesClient(); err != nil {
		return err
	}

	// initialize backup policies client
	if err := d.initBackupPoliciesClient(); err != nil {
		return err
	}

	if d.clients.backupVaultClient == nil || d.clients.backupInstancesClient == nil {
		return errors.New("backupVaultClient and/or backupInstancesClient missing")
	}

	// List all backup vaults
	err := listPager(d,
		d.clients.backupVaultClient.NewGetInSubscriptionPager,
		d.clients.backupVaultClient.NewGetInResourceGroupPager,
		func(res armdataprotection.BackupVaultsClientGetInSubscriptionResponse) []*armdataprotection.BackupVaultResource {
			return res.Value
		},
		func(res armdataprotection.BackupVaultsClientGetInResourceGroupResponse) []*armdataprotection.BackupVaultResource {
			return res.Value
		},
		func(vault *armdataprotection.BackupVaultResource) error {
			instances, err := d.discoverBackupInstances(resourceGroupName(util.Deref(vault.ID)), util.Deref(vault.Name))
			if err != nil {
				err := fmt.Errorf("could not discover backup instances: %v", err)
				return err
			}

			for _, instance := range instances {
				dataSourceType := util.Deref(instance.Properties.DataSourceInfo.DatasourceType)

				// Get retention from backup policy
				policy, err := d.clients.backupPoliciesClient.Get(context.Background(), resourceGroupName(*vault.ID), *vault.Name, backupPolicyName(*instance.Properties.PolicyInfo.PolicyID), &armdataprotection.BackupPoliciesClientGetOptions{})
				if err != nil {
					err := fmt.Errorf("could not get backup policy '%s': %w", *instance.Properties.PolicyInfo.PolicyID, err)
					log.Error(err)
					continue
				}

				// TODO(all):Maybe we should differentiate the backup retention period for different resources, e.g., disk vs blobs (Metrics)
				retention := policy.BaseBackupPolicyResource.Properties.(*armdataprotection.BackupPolicy).PolicyRules[0].(*armdataprotection.AzureRetentionRule).Lifecycles[0].DeleteAfter.(*armdataprotection.AbsoluteDeleteOption).GetDeleteOption().Duration

				resp, err := d.handleInstances(vault, instance)
				if err != nil {
					err := fmt.Errorf("could not handle instance")
					return err
				}

				// Check if map entry already exists
				_, ok := d.backupMap[dataSourceType]
				if !ok {
					d.backupMap[dataSourceType] = &backup{
						backup: make(map[string][]*voc.Backup),
					}
				}

				// Store voc.Backup in backupMap
				d.backupMap[dataSourceType].backup[util.Deref(instance.Properties.DataSourceInfo.ResourceID)] = []*voc.Backup{
					{
						Enabled:         true,
						RetentionPeriod: retentionDuration(util.Deref(retention)),
						Storage:         voc.ResourceID(util.Deref(instance.ID)),
						TransportEncryption: &voc.TransportEncryption{
							Enabled:    true,
							Enforced:   true,
							Algorithm:  constants.TLS,
							TlsVersion: constants.TLS1_2, // https://learn.microsoft.com/en-us/azure/backup/transport-layer-security#why-enable-tls-12 (Last access: 04/27/2023)
						},
					},
				}

				d.backupMap[dataSourceType].backupStorages = append(d.backupMap[dataSourceType].backupStorages, resp)
			}
			return nil
		})

	if err != nil {
		return err
	}

	return nil
}

// backupsEmptyCheck checks if the backups list is empty and returns voc.Backup with enabled = false.
func backupsEmptyCheck(backups []*voc.Backup) []*voc.Backup {
	if len(backups) == 0 {
		return []*voc.Backup{
			{
				Enabled:         false,
				RetentionPeriod: -1,
				Interval:        -1,
			},
		}
	}

	return backups
}

func (d *azureDiscovery) handleInstances(vault *armdataprotection.BackupVaultResource, instance *armdataprotection.BackupInstanceResource) (resource voc.IsCloudResource, err error) {
	if vault == nil || instance == nil {
		return nil, ErrVaultInstanceIsEmpty
	}

	raw, err := voc.ToStringInterface([]interface{}{instance, vault})
	if err != nil {
		log.Error(err)
	}

	if *instance.Properties.DataSourceInfo.DatasourceType == "Microsoft.Storage/storageAccounts/blobServices" {
		resource = &voc.ObjectStorage{
			Storage: &voc.Storage{
				Resource: &voc.Resource{
					ID:           voc.ResourceID(*instance.ID),
					Name:         *instance.Name,
					CreationTime: 0,
					GeoLocation: voc.GeoLocation{
						Region: *vault.Location,
					},
					Labels:    nil,
					ServiceID: d.csID,
					Type:      voc.ObjectStorageType,
					Parent:    resourceGroupID(instance.ID),
					Raw:       raw,
				},
			},
		}
	} else if *instance.Properties.DataSourceInfo.DatasourceType == "Microsoft.Compute/disks" {
		resource = &voc.BlockStorage{
			Storage: &voc.Storage{
				Resource: &voc.Resource{
					ID:           voc.ResourceID(*instance.ID),
					Name:         *instance.Name,
					ServiceID:    d.csID,
					CreationTime: 0,
					Type:         voc.BlockStorageType,
					GeoLocation: voc.GeoLocation{
						Region: *vault.Location,
					},
					Labels: nil,
					Parent: resourceGroupID(instance.ID),
					Raw:    raw,
				},
			},
		}
	}

	return
}

// retentionDuration returns the retention string as time.Duration
func retentionDuration(retention string) time.Duration {
	if retention == "" {
		return time.Duration(0)
	}

	// Delete first and last character
	r := retention[1 : len(retention)-1]

	// string to int
	d, err := strconv.Atoi(r)
	if err != nil {
		log.Errorf("could not convert string to int")
		return time.Duration(0)
	}

	// Create duration in hours
	duration := time.Duration(time.Duration(d) * time.Hour * 24)

	return duration
}

// discoverBackupInstances retrieves the instances in a given backup vault.
// Note: It is only possible to backup a storage account with all containers in it.
func (d *azureDiscovery) discoverBackupInstances(resourceGroup, vaultName string) ([]*armdataprotection.BackupInstanceResource, error) {
	var (
		list armdataprotection.BackupInstancesClientListResponse
		err  error
	)

	if resourceGroup == "" || vaultName == "" {
		return nil, errors.New("missing resource group and/or vault name")
	}

	// List all instances in the given backup vault
	listPager := d.clients.backupInstancesClient.NewListPager(resourceGroup, vaultName, &armdataprotection.BackupInstancesClientListOptions{})
	for listPager.More() {
		list, err = listPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("%s: %v", ErrGettingNextPage, err)
			return nil, err
		}
	}

	return list.Value, nil
}

// resourceGroupName returns the resource group name of a given Azure ID
func resourceGroupName(id string) string {
	return strings.Split(id, "/")[4]
}

func resourceGroupID(ID *string) voc.ResourceID {
	// split according to "/"
	s := strings.Split(util.Deref(ID), "/")

	// We cannot really return an error here, so we just return an empty string
	if len(s) < 5 {
		return ""
	}

	id := strings.Join(s[:5], "/")

	return voc.ResourceID(id)
}

// backupPolicyName returns the backup policy name of a given Azure ID
func backupPolicyName(id string) string {
	// split according to "/"
	s := strings.Split(id, "/")

	// We cannot really return an error here, so we just return an empty string
	if len(s) < 10 {
		return ""
	}
	return s[10]
}

// labels converts the resource tags to the vocabulary label
func labels(tags map[string]*string) map[string]string {
	l := make(map[string]string)

	for tag, i := range tags {
		l[tag] = util.Deref(i)
	}

	return l
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
