// Copyright 2023 Fraunhofer AISEC
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
	"fmt"

	"clouditor.io/clouditor/v2/internal/util"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dataprotection/armdataprotection"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearning/armmachinelearning"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/security/armsecurity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
)

// ClientCreateFuncWithSubID is a type that describes a function to create a new Azure SDK client with a subscription ID.
type ClientCreateFuncWithSubID[T any] func(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*T, error)

// ClientCreateFuncWithoutSubID is a type that describes a function to create a new Azure SDK client without a subscription ID.
type ClientCreateFuncWithoutSubID[T any] func(credential azcore.TokenCredential, options *arm.ClientOptions) (*T, error)

// initClientWithSubID creates an Azure client if not already exists. This function will
// log the error, so calling functions should just directly return the error in
// order to avoid double-logging.
// This function is for client create functions with a subscription ID.
func initClientWithSubID[T any](existingClient *T, d *azureDiscovery, fun ClientCreateFuncWithSubID[T]) (client *T, err error) {
	if existingClient != nil {
		return existingClient, nil
	}

	var subID string
	if d.sub != nil {
		subID = util.Deref(d.sub.SubscriptionID)
	}

	client, err = fun(subID, d.cred, &d.clientOptions)
	if err != nil {
		err = fmt.Errorf("could not get %T client: %w", new(T), err)
		log.Error(err)
		return nil, err
	}

	return
}

// initClientWithoutSubID creates an Azure client if not already exists. This function will
// log the error, so calling functions should just directly return the error in
// order to avoid double-logging.
// This function is for client create functions without a subscription ID.
func initClientWithoutSubID[T any](existingClient *T, d *azureDiscovery, fun ClientCreateFuncWithoutSubID[T]) (client *T, err error) {
	if existingClient != nil {
		return existingClient, nil
	}

	client, err = fun(d.cred, &d.clientOptions)
	if err != nil {
		err = fmt.Errorf("could not get %T client: %w", new(T), err)
		log.Error(err)
		return nil, err
	}

	return
}

// initAccountsClient creates the client if not already exists
func (d *azureDiscovery) initAccountsClient() (err error) {
	d.clients.accountsClient, err = initClientWithSubID(d.clients.accountsClient, d, armstorage.NewAccountsClient)
	return
}

// initApplicationGatewayClient creates the client if not already exists
func (d *azureDiscovery) initApplicationGatewayClient() (err error) {
	d.clients.applicationGatewayClient, err = initClientWithSubID(d.clients.applicationGatewayClient, d, armnetwork.NewApplicationGatewaysClient)
	return
}

// initBackupInstancesClient creates the client if not already exists
func (d *azureDiscovery) initBackupInstancesClient() (err error) {
	d.clients.backupInstancesClient, err = initClientWithSubID(d.clients.backupInstancesClient, d, armdataprotection.NewBackupInstancesClient)

	return
}

// initBackupPoliciesClient creates the client if not already exists
func (d *azureDiscovery) initBackupPoliciesClient() (err error) {
	d.clients.backupPoliciesClient, err = initClientWithSubID(d.clients.backupPoliciesClient, d, armdataprotection.NewBackupPoliciesClient)

	return
}

// initBackupVaultsClient creates the client if not already exists
func (d *azureDiscovery) initBackupVaultsClient() (err error) {
	d.clients.backupVaultClient, err = initClientWithSubID(d.clients.backupVaultClient, d, armdataprotection.NewBackupVaultsClient)

	return
}

// initBlobContainerClient creates the client if not already exists
func (d *azureDiscovery) initBlobContainerClient() (err error) {
	d.clients.blobContainerClient, err = initClientWithSubID(d.clients.blobContainerClient, d, armstorage.NewBlobContainersClient)
	return
}

// initBlockStoragesClient creates the client if not already exists
func (d *azureDiscovery) initBlockStoragesClient() (err error) {
	d.clients.blockStorageClient, err = initClientWithSubID(d.clients.blockStorageClient, d, armcompute.NewDisksClient)
	return
}

// initCosmosDBClient creates the client if not already exists
func (d *azureDiscovery) initCosmosDBClient() (err error) {
	d.clients.cosmosDBClient, err = initClientWithSubID(d.clients.cosmosDBClient, d, armcosmos.NewDatabaseAccountsClient)

	return
}

// initMongoDResourcesBClient creates the client if not already exists
func (d *azureDiscovery) initMongoDResourcesBClient() (err error) {
	d.clients.mongoDBResourcesClient, err = initClientWithSubID(d.clients.mongoDBResourcesClient, d, armcosmos.NewMongoDBResourcesClient)

	return
}

// initDatabasesClient creates the client if not already exists
func (d *azureDiscovery) initDatabasesClient() (err error) {
	d.clients.databasesClient, err = initClientWithSubID(d.clients.databasesClient, d, armsql.NewDatabasesClient)

	return
}

// initDefenderClient creates the client if not already exists
func (d *azureDiscovery) initDefenderClient() (err error) {
	d.clients.defenderClient, err = initClientWithoutSubID(d.clients.defenderClient, d, armsecurity.NewPricingsClient)

	return
}

// initDiskEncryptonSetClient creates the client if not already exists
func (d *azureDiscovery) initDiskEncryptonSetClient() (err error) {
	d.clients.diskEncSetClient, err = initClientWithSubID(d.clients.diskEncSetClient, d, armcompute.NewDiskEncryptionSetsClient)

	return
}

// initFileStorageClient creates the client if not already exists
func (d *azureDiscovery) initFileStorageClient() (err error) {
	d.clients.fileStorageClient, err = initClientWithSubID(d.clients.fileStorageClient, d, armstorage.NewFileSharesClient)

	return
}

// initLoadBalancersClient creates the client if not already exists
func (d *azureDiscovery) initLoadBalancersClient() (err error) {
	d.clients.loadBalancerClient, err = initClientWithSubID(d.clients.loadBalancerClient, d, armnetwork.NewLoadBalancersClient)

	return
}

// initNetworkInterfacesClient creates the client if not already exists
func (d *azureDiscovery) initNetworkInterfacesClient() (err error) {
	d.clients.networkInterfacesClient, err = initClientWithSubID(d.clients.networkInterfacesClient, d, armnetwork.NewInterfacesClient)

	return
}

// initNetworkSecurityGroupClient creates the client if not already exists
func (d *azureDiscovery) initNetworkSecurityGroupClient() (err error) {
	d.clients.networkSecurityGroupsClient, err = initClientWithSubID(d.clients.networkSecurityGroupsClient, d, armnetwork.NewSecurityGroupsClient)

	return
}

// azureDiscovery creates the client if not already exists
func (d *azureDiscovery) initResourceGroupsClient() (err error) {
	d.clients.rgClient, err = initClientWithSubID(d.clients.rgClient, d, armresources.NewResourceGroupsClient)

	return
}

// initSQLServersClient creates the client if not already exists
func (d *azureDiscovery) initSQLServersClient() (err error) {
	d.clients.sqlServersClient, err = initClientWithSubID(d.clients.sqlServersClient, d, armsql.NewServersClient)

	return
}

// initThreatProtectionClient creates the client if not already exists
func (d *azureDiscovery) initThreatProtectionClient() (err error) {
	d.clients.threatProtectionClient, err = initClientWithSubID(d.clients.threatProtectionClient, d, armsql.NewDatabaseAdvancedThreatProtectionSettingsClient)

	return
}

// initVirtualMachinesClient creates the client if not already exists
func (d *azureDiscovery) initVirtualMachinesClient() (err error) {
	d.clients.virtualMachinesClient, err = initClientWithSubID(d.clients.virtualMachinesClient, d, armcompute.NewVirtualMachinesClient)

	return
}

// initWebAppsClient creates the client if not already exists
func (d *azureDiscovery) initWebAppsClient() (err error) {
	d.clients.webAppsClient, err = initClientWithSubID(d.clients.webAppsClient, d, armappservice.NewWebAppsClient)

	return
}

// initMachineLearningComputeClient creates the client if not already exists
func (d *azureDiscovery) initMachineLearningComputeClient() (err error) {
	d.clients.mlComputeClient, err = initClientWithSubID(d.clients.mlComputeClient, d, armmachinelearning.NewComputeClient)

	return
}

// initMLWorkspaceClient creates the client if not already exists
func (d *azureDiscovery) initMLWorkspaceClient() (err error) {
	d.clients.mlWorkspaceClient, err = initClientWithSubID(d.clients.mlWorkspaceClient, d, armmachinelearning.NewWorkspacesClient)

	return
}

// initDiagnosticsSettingsClient creates the client if not already exists
func (d *azureDiscovery) initDiagnosticsSettingsClient() (err error) {
	d.clients.diagnosticSettingsClient, err = initClientWithoutSubID(d.clients.diagnosticSettingsClient, d, armmonitor.NewDiagnosticSettingsClient)

	return
}
