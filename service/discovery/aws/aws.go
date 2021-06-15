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

// awsDiscovery holds configurations across all services within AWS
type awsDiscovery struct {
	cfg aws.Config
}

// NewAwsDiscovery constructs a new awsDiscovery
// ToDo: "Overload" (switch) with staticCredentialsProvider
func NewAwsDiscovery() *awsDiscovery {
	d := &awsDiscovery{}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logrus.Errorf("Could not load default config: %v", err)
	}
	// ToDo: Test if proper region was loaded and maybe remove line
	logrus.Printf("Loaded credentials in region: %v", cfg.Region)
	d.cfg = cfg
	return d
}

// ToDo: I should make the services mor OO like
// DiscoverAll ToDo: Accumulate all service responses into, e.g., one JSON
func (d *awsDiscovery) discoverAll(*awsDiscovery) {
	logrus.Println("Discovering all services (s3,ec2).")
	//rawBuckets := List(GetS3Client(d.cfg))
	//for i, e := range rawBuckets.Buckets {
	//	bucket := GetObjectsOfBucket()
	//}

}
