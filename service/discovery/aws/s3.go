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
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/sirupsen/logrus"
)

// awsS3Discovery handles the AWS API requests regarding the S3 service
// ToDo: Generalize from s3 to other storage types like EFS
type awsS3Discovery struct {
	client *s3.Client
	// ToDo: Change type to bucketsType to list of bucket structs
	buckets       interface{}
	bucketNames   []string
	isDiscovering bool
}

func (d *awsS3Discovery) Name() string {
	return "Aws Storage Account"
}

func (d *awsS3Discovery) List() (resources []voc.IsResource, err error) {
	log.Info("Getting buckets")
	resp, err := d.client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Errorf("Could not retrieve buckets: %v", err)
	}
	for _, bucket := range resp.Buckets {
		isEncrypted, algorithm, keyManager := d.checkEncryption(*bucket.Name)
		resources = append(resources, &voc.ObjectStorageResource{
			StorageResource: voc.StorageResource{
				Resource: voc.Resource{
					// ToDo: "An Amazon S3 bucket name is globally unique". But I bucket.Name isn't the full URL
					ID: *bucket.Name,
					// ToDo: Maybe get bucket name without URL prefix?
					Name:         *bucket.Name,
					CreationTime: bucket.CreationDate.Unix(),
					Type:         []string{"ObjectStorage", "Storage", "Resource"},
				},
				AtRestEncryption: voc.NewAtRestEncryption(isEncrypted, algorithm, keyManager),
			},
			// ToDo: Why AtRestEncryption with constructor and with direct access?
			HttpEndpoint: &voc.HttpEndpoint{
				// What is with voc.Resource here? Is it even a resource?
				// ToDo: I don't know how (yet)?
				URL:                 "",
				TransportEncryption: nil,
			},
		})
	}
	return
}

// NewS3Discovery constructs a new awsS3Discovery initializing the s3-client and isDiscovering with true
// ToDo: Discard method
func NewS3Discovery(cfg aws.Config) *awsS3Discovery {
	return &awsS3Discovery{
		client:        s3.NewFromConfig(cfg),
		buckets:       nil,
		isDiscovering: true,
	}
}

// NewAwsStorageDiscovery constructs a new awsS3Discovery initializing the s3-client and isDiscovering with true
// ToDo: cfg as copy (instead of pointer) since we do not want to change it?
func NewAwsStorageDiscovery(cfg aws.Config) discovery.Discoverer {
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
	d.checkEncryption("")
	d.checkPublicAccessBlockConfiguration("")
	d.checkBucketReplication("")
	d.checkLifeCycleConfiguration("")
}

// ToDo: Decide if functions are attached to awsS3Discovery or not. It is, e.g., weird that discoverBuckets
// calls "d.getBuckets(d.client)"

// getBuckets returns all buckets
func (d *awsS3Discovery) getBuckets(clientApi S3ListBucketsAPI) {
	logrus.Println("Discovering s3:")
	resp, err := clientApi.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatalf("Error occured while retrieving buckets: %v", err)
	}
	logrus.Printf("Retrieved %v buckets.", len(resp.Buckets))
	d.buckets = resp
	for _, bucket := range resp.Buckets {
		d.bucketNames = append(d.bucketNames, *bucket.Name)
	}

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
// ToDo: algorithm "", "AES256" or "aws:kms" as algorithm. Is this the right format?
// ToDo: keyManager "", "SSE-S3" or "SSE-KMS". Is this the right format?
func (d *awsS3Discovery) checkEncryption(bucket string) (bool, string, string) {
	log.Printf("Checking encryption for bucket %v.\n", bucket)
	input := s3.GetBucketEncryptionInput{
		Bucket:              aws.String(bucket),
		ExpectedBucketOwner: nil,
	}
	resp, err := d.client.GetBucketEncryption(context.TODO(), &input)
	if err != nil {
		logrus.Errorf("Probably no encryption enabled. Error: %v", err)
		// ToDo: Simply Return good?
		return false, "", ""
	}
	log.Println("Bucket is encrypted.")
	// ToDo: Why there are multiple rules? For now I check only the first
	var algorithm string
	var keyManager string
	algorithm = string(resp.ServerSideEncryptionConfiguration.Rules[0].ApplyServerSideEncryptionByDefault.SSEAlgorithm)
	if algorithm == string(types.ServerSideEncryptionAes256) {
		keyManager = "SSE-S3"
	} else {
		keyManager = "SSE-KMS"
	}
	//for i, rule := range resp.ServerSideEncryptionConfiguration.Rules {
	//	algorithm = string(rule.ApplyServerSideEncryptionByDefault.SSEAlgorithm)
	//	keyManager = string(rule.)
	//	log.Printf("Rule %v: %v", i, algorithm)
	//}

	return true, algorithm, keyManager
}

// checkPublicAccessBlockConfiguration gets the bucket's access configuration
func (d *awsS3Discovery) checkPublicAccessBlockConfiguration(bucket string) (false bool) {
	log.Printf("Check if bucket %v has public access...", bucket)
	input := s3.GetPublicAccessBlockInput{
		Bucket:              aws.String(bucket),
		ExpectedBucketOwner: nil,
	}
	resp, err := d.client.GetPublicAccessBlock(context.TODO(), &input)
	if err != nil {
		log.Errorf("Error found: %v", err)
		return
	}
	log.Printf("Found: %v", resp.PublicAccessBlockConfiguration)

	configs := resp.PublicAccessBlockConfiguration
	// ToDo: Currently every configuration has to be unset. Anyway in future, the configs are saved to struct since the
	// discovery does no assessment
	if !configs.BlockPublicAcls || !configs.BlockPublicPolicy || !configs.IgnorePublicAcls || !configs.RestrictPublicBuckets {
		return
	}
	return true
}

// checkBucketReplication gets the bucket's replication configuration
func (d *awsS3Discovery) checkBucketReplication(bucket string) (false bool) {
	log.Printf("Check if bucket '%v' is been replicated...", bucket)
	input := s3.GetBucketReplicationInput{
		Bucket:              aws.String(bucket),
		ExpectedBucketOwner: nil,
	}
	resp, err := d.client.GetBucketReplication(context.TODO(), &input)
	if err != nil {
		logrus.Errorf("Error (probably no replica configuration): %v", err)
		return
	}
	logrus.Println(resp.ReplicationConfiguration)
	return true
}

// checkLifeCycleConfiguration gets the bucket's lifecycle configuration
// ToDo
func (d *awsS3Discovery) checkLifeCycleConfiguration(bucket string) (false bool) {
	log.Printf("Check life cycle configuration for bucket '%v'", bucket)
	input := s3.GetBucketLifecycleConfigurationInput{
		Bucket:              aws.String(bucket),
		ExpectedBucketOwner: nil,
	}
	resp, err := d.client.GetBucketLifecycleConfiguration(context.TODO(), &input)
	if err != nil {
		logrus.Errorf("Error occurred: %v", err)
		return
	}
	logrus.Printf(string(resp.Rules[0].Expiration.Days))
	return true
}
