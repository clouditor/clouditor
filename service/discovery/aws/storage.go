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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"time"
)

// awsS3Discovery handles the AWS API requests regarding the S3 service
// ToDo: Generalize from s3 to storage for including other types, like EFS -> Probably storage folder with 3 .go files
type awsS3Discovery struct {
	client        S3API
	isDiscovering bool
	buckets       []bucket
}

// bucket contains meta data about a S3 bucket
type bucket struct {
	arn             string
	name            string
	numberOfObjects int
	creationTime    time.Time
	endpoint        string
	region          string
}

// S3API describes the S3 client interface (for mock testing)
type S3API interface {
	ListBuckets(ctx context.Context,
		params *s3.ListBucketsInput,
		optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
	GetBucketEncryption(ctx context.Context,
		params *s3.GetBucketEncryptionInput,
		optFns ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error)
	GetBucketPolicy(ctx context.Context,
		params *s3.GetBucketPolicyInput,
		optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error)
	GetBucketLocation(ctx context.Context,
		params *s3.GetBucketLocationInput,
		optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error)
	GetPublicAccessBlock(ctx context.Context,
		params *s3.GetPublicAccessBlockInput,
		optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error)
	GetBucketReplication(ctx context.Context,
		params *s3.GetBucketReplicationInput,
		optFns ...func(*s3.Options)) (*s3.GetBucketReplicationOutput, error)
	GetBucketLifecycleConfiguration(ctx context.Context,
		params *s3.GetBucketLifecycleConfigurationInput,
		optFns ...func(*s3.Options)) (*s3.GetBucketLifecycleConfigurationOutput, error)
}

// BucketPolicy matches the returned bucket policy in JSON from AWS
type BucketPolicy struct {
	ID        string
	Version   string
	Statement []Statement
}
type Statement struct {
	Action   string
	Effect   string
	Resource []string
	Condition
}
type Condition struct {
	Bool
}
type Bool struct {
	AwsSecureTransport bool `json:"aws:SecureTransport"`
}

// Name is the method implementation defined in the discovery.Discoverer interface
func (d *awsS3Discovery) Name() string {
	return "AWS Blob Storage"
}

// List is the method implementation defined in the discovery.Discoverer interface
func (d *awsS3Discovery) List() (resources []voc.IsResource, err error) {
	var encryptionAtRest voc.AtRestEncryption
	var encryptionAtTransmit voc.TransportEncryption

	log.Info("Starting List() in ", d.Name())
	err = d.getBuckets()
	if err != nil {
		return
	}
	log.Println(d.buckets)
	for _, bucket := range d.buckets {
		log.Println("Getting resources for", bucket.name)
		encryptionAtRest, err = d.getEncryptionAtRest(bucket.name)
		if err != nil {
			return
		}
		encryptionAtTransmit, err = d.getTransportEncryption(bucket.name)
		if err != nil {
			return
		}
		resources = append(resources, &voc.ObjectStorageResource{
			StorageResource: voc.StorageResource{
				Resource: voc.Resource{
					ID:           bucket.arn,
					Name:         bucket.name,
					CreationTime: bucket.creationTime.Unix(),
					Type:         []string{"ObjectStorage", "Storage", "Resource"},
				},
				AtRestEncryption: &encryptionAtRest,
			},
			HttpEndpoint: &voc.HttpEndpoint{
				URL:                 bucket.endpoint,
				TransportEncryption: &encryptionAtTransmit,
			},
		})
	}
	return
}
func (b bucket) String() string {
	return fmt.Sprintf("[ARN: %v, Name: %v, Creation Time: %v, Number of objects: %v]", b.arn, b.name, b.creationTime, b.numberOfObjects)
}

// NewAwsStorageDiscovery constructs a new awsS3Discovery initializing the s3-client and isDiscovering with true
func NewAwsStorageDiscovery(cfg aws.Config) discovery.Discoverer {
	return &awsS3Discovery{
		client:        s3.NewFromConfig(cfg),
		buckets:       nil,
		isDiscovering: true,
	}
}

// S3ListBucketsAPI is the interface for the List function (used for mock testing)
type S3ListBucketsAPI interface {
	ListBuckets(ctx context.Context,
		params *s3.ListBucketsInput,
		optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
}

// getBuckets returns all buckets
func (d *awsS3Discovery) getBuckets() (err error) {
	log.Println("Getting buckets in s3...")
	var resp *s3.ListBucketsOutput
	resp, err = d.client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Error("Error occurred.")
		var oe *smithy.OperationError
		if errors.As(err, &oe) {
			log.Printf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
		}
		return
	}
	var region string
	log.Printf("Retrieved %v buckets.", len(resp.Buckets))
	for _, b := range resp.Buckets {
		//d.buckets = resp
		region, err = d.getRegion(aws.ToString(b.Name))
		if err != nil {
			return
		}
		d.buckets = append(d.buckets, bucket{
			arn:          "arn:aws:s3:::" + *b.Name,
			name:         aws.ToString(b.Name),
			creationTime: aws.ToTime(b.CreationDate),
			region:       region,
			endpoint:     "https://" + aws.ToString(b.Name) + ".s3." + region + ".amazonaws.com",
			// ToDo: Implement method for retrieving the number of objects per bucket (if needed)
			numberOfObjects: -1,
		})
	}
	return
}

// getEncryptionAtRest gets the bucket's encryption configuration
func (d *awsS3Discovery) getEncryptionAtRest(bucket string) (e voc.AtRestEncryption, err error) {
	log.Printf("Checking encryption for bucket %v.\n", bucket)
	input := s3.GetBucketEncryptionInput{
		Bucket:              aws.String(bucket),
		ExpectedBucketOwner: nil,
	}
	var resp *s3.GetBucketEncryptionOutput

	resp, err = d.client.GetBucketEncryption(context.TODO(), &input)
	if err != nil {
		log.Error("Error occurred")
		var oe *smithy.OperationError
		if errors.As(err, &oe) {
			log.Printf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
		}
		return
	}
	log.Println("Bucket is encrypted.")
	e.Algorithm = string(resp.ServerSideEncryptionConfiguration.Rules[0].ApplyServerSideEncryptionByDefault.SSEAlgorithm)
	if e.Algorithm == string(types.ServerSideEncryptionAes256) {
		e.KeyManager = "SSE-S3"
	} else {
		e.KeyManager = "SSE-KMS"
	}
	e.Enabled = true
	return
}

// "confirm that your bucket policies explicitly deny access to HTTP requests"
// https://aws.amazon.com/premiumsupport/knowledge-center/s3-bucket-policy-for-config-rule/
// getTransportEncryption loops over all statements in the bucket policy and checks if one statement denies https only == false
func (d *awsS3Discovery) getTransportEncryption(bucket string) (encryptionAtTransit voc.TransportEncryption, err error) {
	input := s3.GetBucketPolicyInput{
		Bucket:              aws.String(bucket),
		ExpectedBucketOwner: nil,
	}
	var resp *s3.GetBucketPolicyOutput

	resp, err = d.client.GetBucketPolicy(context.TODO(), &input)

	// Case 1: No bucket policy or api error -> no https only set
	if err != nil {
		log.Error("Error occurred.")
		var oe *smithy.OperationError
		if errors.As(err, &oe) {
			log.Printf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
		}
		return
	}

	// Case 2: bucket policy -> check if https only is set
	var policy BucketPolicy
	err = json.Unmarshal([]byte(aws.ToString(resp.Policy)), &policy)
	if err != nil {
		log.Error("Error occurred while unmarshalling the bucket policy:", err)
		return
	}
	// one statement has set https only -> default encryption is set
	for _, statement := range policy.Statement {
		if statement.Effect == "Deny" && !statement.Condition.AwsSecureTransport && statement.Action == "s3:*" {
			encryptionAtTransit = voc.TransportEncryption{
				Encryption: voc.Encryption{Enabled: true},
				Enforced:   true,
				TlsVersion: "TLS1.2",
			}
			return
		}
	}
	encryptionAtTransit = voc.TransportEncryption{
		// ToDo: I think you can only enforce it via a bucket policy. But not look if its enabled or not, otherwise
		Encryption: voc.Encryption{Enabled: false},
		Enforced:   false,
		TlsVersion: "",
	}
	return

}

// getRegion returns the region where the bucket resides
func (d *awsS3Discovery) getRegion(bucket string) (region string, err error) {
	input := s3.GetBucketLocationInput{
		Bucket: aws.String(bucket),
	}
	var resp *s3.GetBucketLocationOutput
	resp, err = d.client.GetBucketLocation(context.TODO(), &input)
	if err != nil {
		log.Error("Error occurred.")
		var oe *smithy.OperationError
		if errors.As(err, &oe) {
			log.Printf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
		}
		return
	}
	region = string(resp.LocationConstraint)
	return
}

// ToDo: The next checks are not defined yet (in ontology or in voc). They were checked in Clouditor 1.0

//// getPublicAccessBlockConfiguration gets the bucket's access configuration
//func (d *awsS3Discovery) getPublicAccessBlockConfiguration(bucket string) (false bool) {
//	log.Printf("Check if bucket %v has public access...", bucket)
//	input := s3.GetPublicAccessBlockInput{
//		Bucket:              aws.String(bucket),
//		ExpectedBucketOwner: nil,
//	}
//	resp, err := d.client.GetPublicAccessBlock(context.TODO(), &input)
//	if err != nil {
//		log.Errorf("Error found: %v", err)
//		return
//	}
//	log.Printf("Found: %v", resp.PublicAccessBlockConfiguration)
//
//	configs := resp.PublicAccessBlockConfiguration
//	if !configs.BlockPublicAcls || !configs.BlockPublicPolicy || !configs.IgnorePublicAcls || !configs.RestrictPublicBuckets {
//		return
//	}
//	return true
//}

// getBucketObjects returns all objects of the given bucket
// ToDo: Do we need to iterate through single bucket objects or do we only check the general bucket settings?
// ToDo: "Overload" method s.t. you can list all objects from all buckets or only from a specific set (e.g. 1) of buckets
//func (d *awsS3Discovery) getBucketObjects(myBucket string) *s3.ListObjectsV2Output {
//	Cfg, err := config.LoadDefaultConfig(context.TODO())
//	if err != nil {
//		fmt.Println("Error occurred:", err)
//		log.Fatal(err)
//	}
//
//	client := s3.NewFromConfig(Cfg)
//
//	output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
//		Bucket: aws.String(myBucket),
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	log.Println("first page results:")
//	for _, object := range output.Contents {
//		log.Printf("key=%s size=%d", aws.ToString(object.Key), object.Size)
//	}
//
//	return output
//}

//// checkBucketReplication gets the bucket's replication configuration
//func (d *awsS3Discovery) checkBucketReplication(bucket string) (false bool) {
//	log.Printf("Check if bucket '%v' is been replicated...", bucket)
//	input := s3.GetBucketReplicationInput{
//		Bucket:              aws.String(bucket),
//		ExpectedBucketOwner: nil,
//	}
//	resp, err := d.client.GetBucketReplication(context.TODO(), &input)
//	if err != nil {
//		logrus.Errorf("Error (probably no replica configuration): %v", err)
//		return
//	}
//	logrus.Println(resp.ReplicationConfiguration)
//	return true
//}
//
//// checkLifeCycleConfiguration gets the bucket's lifecycle configuration
//// ToDo
//func (d *awsS3Discovery) checkLifeCycleConfiguration(bucket string) (false bool) {
//	log.Printf("Check life cycle configuration for bucket '%v'", bucket)
//	input := s3.GetBucketLifecycleConfigurationInput{
//		Bucket:              aws.String(bucket),
//		ExpectedBucketOwner: nil,
//	}
//	resp, err := d.client.GetBucketLifecycleConfiguration(context.TODO(), &input)
//	if err != nil {
//		logrus.Errorf("Error occurred: %v", err)
//		return
//	}
//	logrus.Printf(string(resp.Rules[0].Expiration.Days))
//	return true
//}
