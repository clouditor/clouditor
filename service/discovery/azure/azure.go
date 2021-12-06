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

	// for now, do not re-authorize. in the future, we would probably need to check, if
	// the token is still valid. or maybe Azure does this for us?
	if a.isAuthorized {
		return
	}

	subClient := subscriptions.NewClient()
	a.apply(&subClient.Client)

	// get first subcription
	page, _ := subClient.List(context.Background())
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

func ResourceGroupName(id string) string {
	return strings.Split(id, "/")[4]
}
