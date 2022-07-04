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
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	"github.com/Azure/go-autorest/autorest/to"

	"github.com/sirupsen/logrus"
	// "github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/resources/mgmt/subscriptions"
	// "github.com/Azure/go-autorest/autorest"
)

var (
	log *logrus.Entry

	ErrCouldNotAuthenticate = errors.New("could not authenticate to Azure with any authorizer " +
		"(from environment, file, or CLI)")
)

type DiscoveryOption interface {
	apply(azcore.TokenCredential)
}

// type senderOption struct {
// 	sender autorest.Sender
// }

// func (o senderOption) apply(client *autorest.Client) {
// 	client.Sender = o.sender
// }

// func WithSender(sender autorest.Sender) DiscoveryOption {
// 	return &senderOption{sender}
// }

// credentialOption contains the client secret credential
type credentialOption struct {
	credential azcore.TokenCredential
}

func WithAuthorizer(credential azcore.TokenCredential) DiscoveryOption {
	return &credentialOption{credential: credential}
}

func init() {
	log = logrus.WithField("component", "azure-discovery")
}

// apply sets the credential
func (a credentialOption) apply(credential azcore.TokenCredential) {
	credential = a.credential
}

// azureDiscovery contains the necessary
type azureDiscovery struct {
	authCredentials *credentialOption // Contains the credentials
	sub             armsubscription.Subscription

	isAuthorized bool

	options []DiscoveryOption
}

func (a *azureDiscovery) authorize() (err error) {
	if a.authCredentials == nil {
		return errors.New("no authorized was available")
	}

	// TODO(anatheka): Still after Azure sdk update
	// If using NewAuthorizerFromFile() in discovery file, we do not need to re-authorize.
	// If using NewAuthorizerFromCLI() in discovery file, the token expires after 75 minutes.
	if a.isAuthorized {
		return
	}

	cred, err := NewAuthorizer()
	if err != nil {
		err = fmt.Errorf("could not get azure credentials: %w", err)
		log.Error(err)
		return err
	}
	a.apply(cred)

	// Create new subscriptions client
	subClient, err := armsubscription.NewSubscriptionsClient(cred, &arm.ClientOptions{})
	if err != nil {
		err = fmt.Errorf("could not get new subscription client: %w", err)
		return err
	}

	// get subscriptions
	subPager := subClient.NewListPager(nil)
	subList := make([]*armsubscription.Subscription, 0)
	for subPager.More() {
		pageResponse, err := subPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("could not get azure subscription: %w", err)
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

	log.Infof("Using %s as subscription", *a.sub.SubscriptionID)

	a.isAuthorized = true

	return nil
}

func (a *azureDiscovery) apply(cred azcore.TokenCredential) {
	if a.authCredentials != nil {
		a.authCredentials.apply(cred)
	}

	for _, v := range a.options {
		v.apply(cred)
	}
}

// NewAuthorizer returns the Azure credential using one of the following authentication types (in the following order):
//  EnvironmentCredential
//  ManagedIdentityCredential
//  AzureCLICredential
func NewAuthorizer() (*azidentity.DefaultAzureCredential, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("Authentication failure: %+v", err)
	}

	return cred, nil
}

func resourceGroupName(id string) string {
	return strings.Split(id, "/")[4]
}

// labels return the tags from resources in the format map[string]string
func labels(tags map[string]*string) map[string]string {
	labels := make(map[string]string)

	for tag, i := range tags {
		labels[tag] = to.String(i)
	}

	return labels
}
