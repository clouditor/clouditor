//go:build exclude

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
	"clouditor.io/clouditor/internal/constants"
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
	csID          string
}

// bucket contains metadata about a S3 bucket
type bucket struct {
	arn          string
	name         string
	creationTime time.Time
	endpoint     string
	region       string
	raw          []interface{}
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
func (*awsS3Discovery) Name() string {
	return "AWS Blob Storage"
}

// List is the method implementation defined in the discovery.Discoverer interface
func (d *awsS3Discovery) List() (resources []voc.IsCloudResource, err error) {
	var (
		encryptionAtRest    voc.IsAtRestEncryption
		encryptionAtTransit *voc.TransportEncryption
		rawBucketEncOutput  *s3.GetBucketEncryptionOutput
		rawBucketTranspEnc  *s3.GetBucketPolicyOutput
	)

	log.Infof("Collecting evidences in %s", d.Name())
	var buckets []bucket
	buckets, err = d.getBuckets()
	if err != nil {
		return
	}

	for _, b := range buckets {
		encryptionAtRest, rawBucketEncOutput, err = d.getEncryptionAtRest(&b)
		if err != nil {
			return
		}
		encryptionAtTransit, rawBucketTranspEnc, err = d.getTransportEncryption(b.name)
		if err != nil {
			return
		}

		resources = append(resources,
			// Add ObjectStorage
			&voc.ObjectStorage{
				Storage: &voc.Storage{
					Resource: discovery.NewResource(d,
						voc.ResourceID(b.arn),
						b.name,
						&b.creationTime,
						voc.GeoLocation{
							Region: b.region,
						},
						nil,
						"",
						voc.ObjectStorageType,
						&b, &rawBucketEncOutput, &rawBucketTranspEnc, &b.raw),
					AtRestEncryption: encryptionAtRest,
				},
			},
			// Add ObjectStorageService
			&voc.ObjectStorageService{
				StorageService: &voc.StorageService{
					Storage: []voc.ResourceID{voc.ResourceID(b.arn)},
					NetworkService: &voc.NetworkService{
						Networking: &voc.Networking{
							Resource: discovery.NewResource(d,
								voc.ResourceID(b.arn),
								b.name,
								&b.creationTime,
								voc.GeoLocation{Region: b.region},
								nil,
								"",
								voc.ObjectStorageServiceType,
								&b, &rawBucketEncOutput, &rawBucketTranspEnc, &b.raw,
							),
						},
						TransportEncryption: encryptionAtTransit,
					},
				},
				HttpEndpoint: &voc.HttpEndpoint{
					Url:                 b.endpoint,
					TransportEncryption: encryptionAtTransit,
				},
			})
	}
	return
}

func (d *awsS3Discovery) CloudServiceID() string {
	return d.csID
}

func (b *bucket) String() string {
	return fmt.Sprintf("[ARN: %v, Name: %v, Creation Time: %v]", b.arn, b.name, b.creationTime)
}

// NewAwsStorageDiscovery constructs a new awsS3Discovery initializing the s3-api and isDiscovering with true
func NewAwsStorageDiscovery(client *Client, cloudServiceID string) discovery.Discoverer {
	return &awsS3Discovery{
		storageAPI:    s3.NewFromConfig(client.cfg),
		isDiscovering: true,
		awsConfig:     client,
		csID:          cloudServiceID,
	}
}

// getBuckets returns all buckets
func (d *awsS3Discovery) getBuckets() (buckets []bucket, err error) {
	var resp *s3.ListBucketsOutput
	resp, err = d.storageAPI.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return nil, prettyError(err)
	}
	var region string
	for _, b := range resp.Buckets {
		var (
			rawRegion *s3.GetBucketLocationOutput
		)

		region, rawRegion, err = d.getRegion(aws.ToString(b.Name))
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
				raw:          []interface{}{b, rawRegion},
			})
		}
	}
	return
}

// getEncryptionAtRest gets the bucket's encryption configuration
func (d *awsS3Discovery) getEncryptionAtRest(bucket *bucket) (e voc.IsAtRestEncryption, resp *s3.GetBucketEncryptionOutput, err error) {
	input := s3.GetBucketEncryptionInput{
		Bucket:              aws.String(bucket.name),
		ExpectedBucketOwner: nil,
	}

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
		e = &voc.ManagedKeyEncryption{AtRestEncryption: &voc.AtRestEncryption{
			Algorithm: string(alg),
			Enabled:   true,
		}}
	} else {
		e = &voc.CustomerKeyEncryption{
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
func (d *awsS3Discovery) getTransportEncryption(bucket string) (*voc.TransportEncryption, *s3.GetBucketPolicyOutput, error) {
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
					TlsVersion: constants.TLS1_2,
					Algorithm:  constants.TLS,
				}, resp, nil
			}
			// Any other error is a connection error with AWS : Format err and return it
			err = formatError(ae)
		}
		// return any error (but according to doc: "All service API response errors implement the smithy.APIError")
		return nil, resp, err
	}

	// Case 2: bucket policy -> check if https only is set
	// TODO(lebogg): bucket policy json fail still means that https is enabled (it always is). Still return error?
	var policy BucketPolicy
	err = json.Unmarshal([]byte(aws.ToString(resp.Policy)), &policy)
	if err != nil {
		return nil, resp, fmt.Errorf("error occurred while unmarshalling the bucket policy: %v", err)
	}
	// one statement has set https only -> default encryption is set
	for _, statement := range policy.Statement {
		if a, ok := statement.Action.(string); ok {
			if statement.Effect == "Deny" && !statement.Condition.AwsSecureTransport && a == "s3:*" {
				return &voc.TransportEncryption{
					Enforced:   true,
					Enabled:    true,
					TlsVersion: constants.TLS1_2,
					Algorithm:  constants.TLS,
				}, resp, nil
			}
		}
		if actions, ok := statement.Action.([]string); ok {
			for _, a := range actions {
				if statement.Effect == "Deny" && !statement.Condition.AwsSecureTransport && a == "s3:*" {
					return &voc.TransportEncryption{
						Enforced:   true,
						Enabled:    true,
						TlsVersion: constants.TLS1_2,
						Algorithm:  constants.TLS,
					}, resp, nil
				}
			}
		}
	}
	return &voc.TransportEncryption{
		Enforced:   false,
		Enabled:    true,
		TlsVersion: constants.TLS1_2,
		Algorithm:  constants.TLS,
	}, resp, nil

}

// getRegion returns the region where the bucket resides
func (d *awsS3Discovery) getRegion(bucket string) (region string, resp *s3.GetBucketLocationOutput, err error) {
	input := s3.GetBucketLocationInput{
		Bucket: aws.String(bucket),
	}
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
