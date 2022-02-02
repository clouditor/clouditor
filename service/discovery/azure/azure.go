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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/resources/mgmt/subscriptions"
	"github.com/Azure/go-autorest/autorest"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

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
	// TODO: How do we check, if the token is still valid?
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

func NewAuthorizer() (autorest.Authorizer, error) {
	// create an authorizer from file or as fallback from the CLI
	// if authorizer is from CLI, the access token expires after 75 minutes
	authorizer, err := auth.NewAuthorizerFromFile(autorest_azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		log.Errorf("Could not authenticate to Azure with authorizer from file: %v", err)
		log.Infof("Fallback to Azure with authorizer from CLI.")
		authorizer, err = auth.NewAuthorizerFromCLI()
		if err != nil {
			log.Errorf("Could not authenticate to Azure with authorizer from CLI: %v", err)
			return nil, status.Errorf(codes.FailedPrecondition, "could not authenticate to Azure: %v", err)
		}
		log.Info("Using Azure authorizer from CLI. The discovery times out after 1 hour.")
	} else {
		log.Info("Using Azure authorizer from file.")
	}

	return authorizer, nil
}

func getResourceGroupName(id string) string {
	return strings.Split(id, "/")[4]
}
