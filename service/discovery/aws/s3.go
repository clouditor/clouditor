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
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
	"log"
)

// awsS3Discovery handles the AWS API requests regarding the S3 service
type awsS3Discovery struct {
	client *s3.Client
	// ToDo: Change type to bucketsType to list of bucket structs
	buckets       interface{}
	isDiscovering bool
}

// NewS3Discovery constructs a new awsS3Discovery initializing the s3-client and isDiscovering with true
// ToDo: cfg as copy (instead of pointer) since we do not want to change it?
func NewS3Discovery(cfg aws.Config) *awsS3Discovery {
	return &awsS3Discovery{
		client:        s3.NewFromConfig(cfg),
		buckets:       nil,
		isDiscovering: true,
	}
}

// S3ListBucketsAPI is the interface for the List function (used for mock testing)
// ToDo: Is it a good idea to do so, i.e. integrating test stuff here? It is recommended by the AWS SDK documentation
type S3ListBucketsAPI interface {
	ListBuckets(ctx context.Context,
		params *s3.ListBucketsInput,
		optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
}

// discoverBuckets retrieves all buckets and its configurations needed for assessing
func (d *awsS3Discovery) discoverBuckets() {
	// ToDo: Double check that d.client can't be zero due to "private" access with "public" constructor?
	d.getBuckets(d.client)
	d.checkEncryption()
	d.checkPublicAccessBlockConfiguration()
	d.checkBucketReplication()
	d.checkLifeCycleConfiguration()
}

// ToDo: Decide if functions are attached to awsS3Discovery or not. It is, e.g., weird that discoverBuckets
// calls "d.getBuckets(d.client)"

// getBuckets returns all buckets
func (d *awsS3Discovery) getBuckets(clientApi S3ListBucketsAPI) {
	logrus.Println("Discovering s3:")
	// ToDo: One return line
	resp, err := clientApi.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatalf("Error occured while retrieving buckets: %v", err)
	}
	logrus.Printf("Retrieved %v buckets.", len(resp.Buckets))
	d.buckets = resp

	// ToDo: Call getBucketObjects and associate the objects with the buckets
}

// getBucketObjects returns all objects of the given bucket
// ToDo: Do we need to iterate through single bucket objects or do we only check the general bucket settings?
// ToDo: "Overload" method s.t. you can list all objects from all buckets or only from a specific set (e.g. 1) of buckets
func (d *awsS3Discovery) getBucketObjects(myBucket string) *s3.ListObjectsV2Output {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Error occurred:", err)
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg)

	output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(myBucket),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("first page results:")
	for _, object := range output.Contents {
		log.Printf("key=%s size=%d", aws.ToString(object.Key), object.Size)
	}

	return output
}

// checkEncryption gets the bucket's encryption configuration
// ToDo
func (d *awsS3Discovery) checkEncryption() {
	log.Println("ToDo: Implement checkEncryption")
}

// checkPublicAccessBlockConfiguration gets the bucket's access configuration
// ToDo
func (d *awsS3Discovery) checkPublicAccessBlockConfiguration() {
	log.Println("ToDo: Implement checkPublicAccessBlockConfiguration")

}

// checkBucketReplication gets the bucket's replication configuration
// ToDo
func (d *awsS3Discovery) checkBucketReplication() {
	log.Println("ToDo: Implement checkBucketReplication")

}

// checkLifeCycleConfiguration gets the bucket's lifecycle configuration
// ToDo
func (d *awsS3Discovery) checkLifeCycleConfiguration() {
	log.Println("ToDo: Implement checkLifeCycleConfiguration")
}
