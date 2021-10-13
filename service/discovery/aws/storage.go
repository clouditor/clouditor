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
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

// awsS3Discovery handles the AWS API requests regarding the S3 service
type awsS3Discovery struct {
	client        S3API
	isDiscovering bool
}

// bucket contains meta data about a S3 bucket
type bucket struct {
	arn          string
	name         string
	creationTime time.Time
	endpoint     string
	region       string
}

// S3API describes the S3 api interface which is implemented by the official AWS client and mock clients in tests
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
func (d *awsS3Discovery) List() (resources []voc.IsCloudResource, err error) {
	var encryptionAtRest *voc.AtRestEncryption
	var encryptionAtTransmit voc.TransportEncryption

	log.Info("Starting List() in ", d.Name())
	var buckets []bucket
	buckets, err = d.getBuckets()
	if err != nil {
		return
	}
	log.Println("Found", len(buckets), "buckets.")
	for _, bucket := range buckets {
		log.Println("Getting resources for", bucket.name)
		encryptionAtRest, err = d.getEncryptionAtRest(bucket.name)
		if err != nil {
			return
		}
		encryptionAtTransmit, err = d.getTransportEncryption(bucket.name)
		if err != nil {
			return
		}
		resources = append(resources, &voc.ObjectStorage{
			Storage: &voc.Storage{
				CloudResource: &voc.CloudResource{
					ID:           voc.ResourceID(bucket.arn),
					Name:         bucket.name,
					CreationTime: bucket.creationTime.Unix(),
					Type:         []string{"ObjectStorage", "Storage", "Resource"},
					GeoLocation: voc.GeoLocation{
						Region: bucket.region,
					},
				},
				AtRestEncryption: &voc.CustomerKeyEncryption{AtRestEncryption: encryptionAtRest},
			},
			HttpEndpoint: &voc.HttpEndpoint{
				Url:                 bucket.endpoint,
				TransportEncryption: &encryptionAtTransmit,
			},
		})
	}
	return
}
func (b bucket) String() string {
	return fmt.Sprintf("[ARN: %v, Name: %v, Creation Time: %v]", b.arn, b.name, b.creationTime)
}

// NewAwsStorageDiscovery constructs a new awsS3Discovery initializing the s3-api and isDiscovering with true
func NewAwsStorageDiscovery(client *Client) discovery.Discoverer {
	return &awsS3Discovery{
		client:        s3.NewFromConfig(client.cfg),
		isDiscovering: true,
	}
}

// getBuckets returns all buckets
func (d *awsS3Discovery) getBuckets() (buckets []bucket, err error) {
	log.Println("Getting buckets in s3...")
	var resp *s3.ListBucketsOutput
	resp, err = d.client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			err = formatError(ae)
		}
		return
	}
	var region string
	for _, b := range resp.Buckets {
		region, err = d.getRegion(aws.ToString(b.Name))
		if err != nil {
			return
		}
		buckets = append(buckets, bucket{
			arn:          "arn:aws:s3:::" + *b.Name,
			name:         aws.ToString(b.Name),
			creationTime: aws.ToTime(b.CreationDate),
			region:       region,
			endpoint:     "https://" + aws.ToString(b.Name) + ".s3." + region + ".amazonaws.com",
		})
	}
	return
}

// getEncryptionAtRest gets the bucket's encryption configuration
func (d *awsS3Discovery) getEncryptionAtRest(bucket string) (e *voc.AtRestEncryption, err error) {
	log.Printf("Checking encryption for bucket %v.", bucket)
	e = new(voc.AtRestEncryption)
	input := s3.GetBucketEncryptionInput{
		Bucket:              aws.String(bucket),
		ExpectedBucketOwner: nil,
	}
	var resp *s3.GetBucketEncryptionOutput

	resp, err = d.client.GetBucketEncryption(context.TODO(), &input)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "ServerSideEncryptionConfigurationNotFoundError" {
				// This error code is equivalent to "encryption not enabled": set err to nil
				e.Enabled = false
				err = nil
				return
			}
			// Any other error is a connection error with AWS : Format err and return it
			err = formatError(ae)
		}
		// return any error (but according to doc: "All service API response errors implement the smithy.APIError")
		return
	}
	e.Algorithm = string(resp.ServerSideEncryptionConfiguration.Rules[0].ApplyServerSideEncryptionByDefault.SSEAlgorithm)
	//if e.Algorithm == string(types.ServerSideEncryptionAes256) {
	//	e.KeyManager = "SSE-S3"
	//} else {
	//	e.KeyManager = "SSE-KMS"
	//}
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

	// encryption at transit (https) is always enabled and TLS version fixed
	encryptionAtTransit.Enabled = true
	encryptionAtTransit.TlsVersion = "TLS1.2"

	// Case 1: No bucket policy in place or api error -> 'https only' is not set
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "NoSuchBucketPolicy" {
				// This error code is equivalent to "encryption not enforced": set err to nil
				encryptionAtTransit.Enforced = false
				err = nil
				return
			}
			// Any other error is a connection error with AWS : Format err and return it
			err = formatError(ae)
		}
		// return any error (but according to doc: "All service API response errors implement the smithy.APIError")
		return
	}

	// Case 2: bucket policy -> check if https only is set
	// TODO(lebogg): bucket policy json fail still means that https is enabled (it always is). Still return error?
	var policy BucketPolicy
	err = json.Unmarshal([]byte(aws.ToString(resp.Policy)), &policy)
	if err != nil {
		err = fmt.Errorf("error occurred while unmarshalling the bucket policy:%v", err)
		return
	}
	// one statement has set https only -> default encryption is set
	for _, statement := range policy.Statement {
		if statement.Effect == "Deny" && !statement.Condition.AwsSecureTransport && statement.Action == "s3:*" {
			encryptionAtTransit.Enforced = true
			return
		}
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
		var oe *smithy.OperationError
		if errors.As(err, &oe) {
			err = fmt.Errorf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
		}
		return
	}
	region = string(resp.LocationConstraint)
	return
}
