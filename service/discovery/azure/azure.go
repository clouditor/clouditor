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
	autorest_azure "github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/resources/mgmt/subscriptions"
	"github.com/Azure/go-autorest/autorest"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Entry

	ErrCouldNotAuthenticate = errors.New("could not authenticate to Azure with any credentials of the chain")
)

type DiscoveryOption interface {
	apply(*autorest.Client)
}

type senderOption struct {
	sender autorest.Sender
}

func (o senderOption) apply(client *autorest.Client) {
	client.Sender = o.sender
}

func WithSender(sender autorest.Sender) DiscoveryOption {
	return &senderOption{sender}
}

type authorizerOption struct {
	authorizer autorest.Authorizer
}

func WithAuthorizer(authorizer autorest.Authorizer) DiscoveryOption {
	return &authorizerOption{authorizer: authorizer}
}

func init() {
	log = logrus.WithField("component", "azure-discovery")
}

func (a authorizerOption) apply(client *autorest.Client) {
	client.Authorizer = a.authorizer
}

type azureDiscovery struct {
	authOption *authorizerOption
	sub        subscriptions.Subscription

	isAuthorized bool

	options []DiscoveryOption
}

func (a *azureDiscovery) authorize() (err error) {
	if a.authOption == nil {
		return errors.New("no authorized was available")
	}

	// If using NewAuthorizerFromFile() in discovery file, we do not need to re-authorize.
	// If using NewAuthorizerFromCLI() in discovery file, the token expires after 75 minutes.
	if a.isAuthorized {
		return
	}

	subClient := subscriptions.NewClient()
	a.apply(&subClient.Client)

	// get subscriptions
	page, err := subClient.List(context.Background())
	if err != nil {
		err = fmt.Errorf("could not get azure subscription: %v", err)
		return
	}

	// check if list of subscriptions is empty
	if len(page.Values()) == 0 {
		err = errors.New("list of subscriptions is empty")
		return
	}

	// get first subscription
	a.sub = page.Values()[0]

	log.Infof("Using %s as subscription", *a.sub.SubscriptionID)

	a.isAuthorized = true

	return nil
}

func (a azureDiscovery) apply(client *autorest.Client) {
	if a.authOption != nil {
		a.authOption.apply(client)
	}

	for _, v := range a.options {
		v.apply(client)
	}
}

// NewAuthorizer creates authorizer for connecting to Azure using a custom credential chain (ENV, from file and from CLI)
func NewAuthorizer() (authorizer autorest.Authorizer, err error) {
	// First, try to create authorizer via credentials from the environment
	authorizer, err = auth.NewAuthorizerFromEnvironment()
	if err == nil {
		log.Infof("Using authorizer from environment")
		return
	}
	log.Infof("Could not authenticate to Azure with authorizer from environment: %v", err)
	log.Infof("Fallback to authorizer from file")

	// Create authorizer from file
	authorizer, err = auth.NewAuthorizerFromFile(autorest_azure.PublicCloud.ResourceManagerEndpoint)
	if err == nil {
		log.Infof("Using authorizer from file")
		return
	}
	log.Infof("Could not authenticate to Azure with authorizer from file: %v", err)
	log.Infof("Fallback to authorizer from CLI.")

	// Create authorizer from CLI
	authorizer, err = auth.NewAuthorizerFromCLI()
	if err == nil {
		// if authorizer is from CLI, the access token expires after 75 minutes
		log.Info("Using authorizer from CLI. The discovery times out after 1 hour.")
		return
	}
	log.Infof("Could not authenticate to Azure with authorizer from CLI: %v", err)

	// Authorizer couldn't be created
	log.Error(ErrCouldNotAuthenticate)
	return nil, ErrCouldNotAuthenticate
}

func getResourceGroupName(id string) string {
	return strings.Split(id, "/")[4]
}
