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

	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

// awsS3Discovery handles the AWS API requests regarding the S3 service
type awsS3Discovery struct {
	storageAPI    S3API
	isDiscovering bool
	awsConfig     *Client
}

// bucket contains metadata about a S3 bucket
type bucket struct {
	arn          string
	name         string
	creationTime time.Time
	endpoint     string
	region       string
}

// S3API describes the S3 api interface which is implemented by the official AWS storageAPI and mock clients in tests
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
	ID        string      `json:"id"`
	Version   string      `json:"Version"`
	Statement []Statement `json:"Statement"`
}
type Statement struct {
	Action    interface{} `json:"Action"`
	Effect    string      `json:"Effect"`
	Resource  interface{}
	Condition `json:"Condition"`
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
	var encryptionAtRest voc.HasAtRestEncryption
	var encryptionAtTransmit *voc.TransportEncryption

	log.Infof("Collecting evidences in %s", d.Name())
	var buckets []bucket
	buckets, err = d.getBuckets()
	if err != nil {
		return
	}
	for _, b := range buckets {
		encryptionAtRest, err = d.getEncryptionAtRest(b)
		if err != nil {
			return
		}
		encryptionAtTransmit, err = d.getTransportEncryption(b.name)
		if err != nil {
			return
		}
		resources = append(resources, &voc.ObjectStorage{
			Storage: &voc.Storage{
				Resource: &voc.Resource{
					ID:           voc.ResourceID(b.arn),
					Name:         b.name,
					CreationTime: b.creationTime.Unix(),
					Type:         []string{"ObjectStorage", "Storage", "Resource"},
					GeoLocation: voc.GeoLocation{
						Region: b.region,
					},
				},
				AtRestEncryption: encryptionAtRest,
			},
			// TODO(garuppel): Update HttpEndpoint
			HttpEndpoint: &voc.HttpEndpoint{
				Url:                 b.endpoint,
				TransportEncryption: encryptionAtTransmit,
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
		storageAPI:    s3.NewFromConfig(client.cfg),
		isDiscovering: true,
		awsConfig:     client,
	}
}

// getBuckets returns all buckets
func (d *awsS3Discovery) getBuckets() (buckets []bucket, err error) {
	var resp *s3.ListBucketsOutput
	resp, err = d.storageAPI.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
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
		// Currently only buckets are retrieved that are in the region of the users specified region in the config. Since getBucketPolicy throws error if bucket region differs
		// TODO(lebogg): Retrieve all buckets (just remove if) and fix issues with other methods, e.g. getBucketPolicy
		if region == d.awsConfig.cfg.Region {
			buckets = append(buckets, bucket{
				arn:          "arn:aws:s3:::" + *b.Name,
				name:         aws.ToString(b.Name),
				creationTime: aws.ToTime(b.CreationDate),
				region:       region,
				endpoint:     "https://" + aws.ToString(b.Name) + ".s3." + region + ".amazonaws.com",
			})
		}
	}
	return
}

// getEncryptionAtRest gets the bucket's encryption configuration
func (d *awsS3Discovery) getEncryptionAtRest(bucket bucket) (e voc.HasAtRestEncryption, err error) {

	input := s3.GetBucketEncryptionInput{
		Bucket:              aws.String(bucket.name),
		ExpectedBucketOwner: nil,
	}
	var resp *s3.GetBucketEncryptionOutput

	resp, err = d.storageAPI.GetBucketEncryption(context.TODO(), &input)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "ServerSideEncryptionConfigurationNotFoundError" {
				// This error code is equivalent to "encryption not enabled": set err to nil
				e = &voc.AtRestEncryption{
					Confidentiality: nil,
					Algorithm:       "",
					Enabled:         false,
				}
				err = nil
				return
			}
			// Any other error is a connection error with AWS : Format err and return it
			err = formatError(ae)
		}
		// return any error (but according to doc: "All service API response errors implement the smithy.APIError")
		return
	}

	if alg := resp.ServerSideEncryptionConfiguration.Rules[0].ApplyServerSideEncryptionByDefault.SSEAlgorithm; alg == types.ServerSideEncryptionAes256 {
		e = voc.ManagedKeyEncryption{AtRestEncryption: &voc.AtRestEncryption{
			Algorithm: string(alg),
			Enabled:   true,
		}}
	} else {
		e = voc.CustomerKeyEncryption{
			AtRestEncryption: &voc.AtRestEncryption{
				Algorithm: "", // not available
				Enabled:   true,
			},
			// TODO(lebogg): Check in console if bucket.region is the actual region of the key arn
			KeyUrl: "arn:aws:kms:" + bucket.region + ":" + aws.ToString(d.awsConfig.accountID) + ":key/" + aws.ToString(resp.ServerSideEncryptionConfiguration.Rules[0].ApplyServerSideEncryptionByDefault.KMSMasterKeyID),
		}
	}
	return
}

// "confirm that your bucket policies explicitly deny access to HTTP requests"
// https://aws.amazon.com/premiumsupport/knowledge-center/s3-bucket-policy-for-config-rule/
// getTransportEncryption loops over all statements in the bucket policy and checks if one statement denies https only == false
func (d *awsS3Discovery) getTransportEncryption(bucket string) (*voc.TransportEncryption, error) {
	input := s3.GetBucketPolicyInput{
		Bucket:              aws.String(bucket),
		ExpectedBucketOwner: nil,
	}
	var resp *s3.GetBucketPolicyOutput
	var err error

	resp, err = d.storageAPI.GetBucketPolicy(context.TODO(), &input)

	// encryption at transit (https) is always enabled and TLS version fixed

	// Case 1: No bucket policy in place or api error -> 'https only' is not set
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "NoSuchBucketPolicy" {
				// This error code is equivalent to "encryption not enforced": set err to nil
				return &voc.TransportEncryption{
					Enforced:   false,
					Enabled:    true,
					TlsVersion: "TLS1.2",
					Algorithm:  "TLS",
				}, nil
			}
			// Any other error is a connection error with AWS : Format err and return it
			err = formatError(ae)
		}
		// return any error (but according to doc: "All service API response errors implement the smithy.APIError")
		return nil, err
	}

	// Case 2: bucket policy -> check if https only is set
	// TODO(lebogg): bucket policy json fail still means that https is enabled (it always is). Still return error?
	var policy BucketPolicy
	err = json.Unmarshal([]byte(aws.ToString(resp.Policy)), &policy)
	if err != nil {
		return nil, fmt.Errorf("error occurred while unmarshalling the bucket policy: %v", err)
	}
	// one statement has set https only -> default encryption is set
	for _, statement := range policy.Statement {
		if a, ok := statement.Action.(string); ok {
			if statement.Effect == "Deny" && !statement.Condition.AwsSecureTransport && a == "s3:*" {
				return &voc.TransportEncryption{
					Enforced:   true,
					Enabled:    true,
					TlsVersion: "TLS1.2",
					Algorithm:  "TLS",
				}, nil
			}
		}
		if actions, ok := statement.Action.([]string); ok {
			for _, a := range actions {
				if statement.Effect == "Deny" && !statement.Condition.AwsSecureTransport && a == "s3:*" {
					return &voc.TransportEncryption{
						Enforced:   true,
						Enabled:    true,
						TlsVersion: "TLS1.2",
						Algorithm:  "TLS",
					}, nil
				}
			}
		}
	}
	return &voc.TransportEncryption{
		Enforced:   false,
		Enabled:    true,
		TlsVersion: "TLS1.2",
		Algorithm:  "TLS",
	}, nil

}

// getRegion returns the region where the bucket resides
func (d *awsS3Discovery) getRegion(bucket string) (region string, err error) {
	input := s3.GetBucketLocationInput{
		Bucket: aws.String(bucket),
	}
	var resp *s3.GetBucketLocationOutput
	resp, err = d.storageAPI.GetBucketLocation(context.TODO(), &input)
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
