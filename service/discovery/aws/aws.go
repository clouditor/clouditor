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
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("component", "aws-discovery")

// loadDefaultConfig holds config.LoadDefaultConfig() so the test function can mock it
var loadDefaultConfig = config.LoadDefaultConfig

// Client holds configurations across all services within AWS
// TODO(lebogg): deepsource.io wants the struct to exported since NewAwsStorageDiscovery is exported. Encapsulation?
type Client struct {
	Cfg aws.Config
}

// NewClient constructs a new AwsClient
// TODO(lebogg): "Overload" (switch) with staticCredentialsProvider
func NewClient() (*Client, error) {
	c := &Client{}
	cfg, err := loadDefaultConfig(context.TODO())
	if err != nil {
		log.Errorf("Could not load default config: %v", err)
	}
	c.Cfg = cfg
	return c, err
}
