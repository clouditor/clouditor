/*
 * Copyright 2021 Fraunhofer AISEC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *           $$\                           $$\ $$\   $$\
 *           $$ |                          $$ |\__|  $$ |
 *  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
 * $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
 * $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
 * $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
 * \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
 *  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
 *
 * This file is part of Clouditor Community Edition.
 */

package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/resources/mgmt/subscriptions"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

type AzureAuthorizer struct {
	authorizer autorest.Authorizer
	sub        subscriptions.Subscription
	ctx        context.Context
}

var azureAuthorizer = &AzureAuthorizer{}

func init() {
	log = logrus.WithField("component", "azure-discovery")

	// Initialise Azure paramters
	AuthorizeAzure()

}

func AuthorizeAzure() (err error) {
	// create an authorizer from env vars or Azure Managed Service Identity
	azureAuthorizer.authorizer, err = auth.NewAuthorizerFromCLI()
	if err != nil {
		log.Errorf("Could not authenticate to Azure: %s", err)
		return err
	}

	subClient := subscriptions.NewClient()
	subClient.Authorizer = azureAuthorizer.authorizer

	// get first subcription
	page, _ := subClient.List(context.Background())
	azureAuthorizer.sub = page.Values()[0]

	azureAuthorizer.ctx = context.Background()

	log.Infof("Using %s as subscription", *azureAuthorizer.sub.SubscriptionID)

	return nil
}
