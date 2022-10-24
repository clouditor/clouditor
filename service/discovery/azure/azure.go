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
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"

	"github.com/sirupsen/logrus"

	"clouditor.io/clouditor/internal/util"
)

const (
	StorageComponent = "storage"
	ComputeComponent = "compute"
	NetworkComponent = "network"
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

func init() {
	log = logrus.WithField("component", "azure-discovery")
}

type azureDiscovery struct {
	isAuthorized bool

	sub                 armsubscription.Subscription
	cred                azcore.TokenCredential
	clientOptions       arm.ClientOptions
	discovererComponent string
	clients             clients
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
		log.Fatalf("%s: %+v", ErrCouldNotAuthenticate, err)
	}

	return cred, nil
}

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
