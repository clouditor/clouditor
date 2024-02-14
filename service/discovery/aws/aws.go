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

package aws

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("component", "aws-discovery")

// loadDefaultConfig holds config.LoadDefaultConfig() so that NewClient() can use it and test function can mock it
var loadDefaultConfig = config.LoadDefaultConfig

// newFromConfigSTS holds sts.NewFromConfig() so that NewClient() can use it and test function can mock it
var newFromConfigSTS = loadSTSClient

// Client holds configurations across all services within AWS
type Client struct {
	cfg aws.Config
	// accountID is needed for ARN creation
	accountID *string
}

// STSAPI describes the STS api interface which is implemented by the official AWS client and mock clients in tests
type STSAPI interface {
	GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
}

// NewClient constructs a new AwsClient
// TODO(lebogg): "Overload" (switch) with staticCredentialsProvider
func NewClient() (*Client, error) {
	c := &Client{}

	// load configuration
	cfg, err := loadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("could not load default config: %w", err)
	}
	c.cfg = cfg

	// load accountID
	stsClient := newFromConfigSTS(cfg)
	resp, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, prettyError(err)
	}
	c.accountID = resp.Account

	return c, err
}

// formatError returns AWS API specific error code transformed into the default error type
func formatError(ae smithy.APIError) error {
	return fmt.Errorf("code: %v, fault: %v, message: %v", ae.ErrorCode(), ae.ErrorFault(), ae.ErrorMessage())
}

// prettyError returns an AWS API specific error code if it is an AWS error (using [formatError]), otherwise, just the error itself.
func prettyError(err error) error {
	var ae smithy.APIError
	if errors.As(err, &ae) {
		err = formatError(ae)
	}
	return err
}

// loadSTSClient creates the STS client using the STS api interface (for mock testing)
func loadSTSClient(cfg aws.Config) STSAPI {
	client := sts.NewFromConfig(cfg)
	return client
}
