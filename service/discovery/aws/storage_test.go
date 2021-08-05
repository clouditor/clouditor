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
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/stretchr/testify/assert"
)

const (
	mockBucket1             = "mockbucket1"
	mockBucket1Endpoint     = "https://mockbucket1.s3.eu-central-1.amazonaws.com"
	mockBucket1Region       = "eu-central-1"
	mockBucket1CreationTime = "2012-11-01T22:08:41+00:00"
	mockBucket2             = "mockbucket2"
	mockBucket2Endpoint     = "https://mockbucket2.s3.eu-west-1.amazonaws.com"
	mockBucket2Region       = "eu-west-1"
	mockBucket2CreationTime = "2013-12-02T22:08:41+00:00"
	mockBucket3             = "mockbucket3"
)

type mockS3APINew struct{}

func (m mockS3APINew) ListBuckets(_ context.Context,
	_ *s3.ListBucketsInput,
	_ ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	creationDate1, err := time.Parse(time.RFC3339, mockBucket1CreationTime)
	if err != nil {
		log.Error(err)
	}
	creationDate2 := creationDate1.AddDate(1, 1, 1)

	output := &s3.ListBucketsOutput{Buckets: []types.Bucket{
		{
			Name:         aws.String(mockBucket1),
			CreationDate: aws.Time(creationDate1),
		},
		{
			Name:         aws.String(mockBucket2),
			CreationDate: aws.Time(creationDate2),
		},
	}}
	return output, nil
}

func (m mockS3APINew) GetBucketEncryption(_ context.Context,
	params *s3.GetBucketEncryptionInput,
	_ ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error) {
	switch aws.ToString(params.Bucket) {
	case mockBucket1:
		output := &s3.GetBucketEncryptionOutput{
			ServerSideEncryptionConfiguration: &types.ServerSideEncryptionConfiguration{
				Rules: []types.ServerSideEncryptionRule{
					{
						ApplyServerSideEncryptionByDefault: &types.ServerSideEncryptionByDefault{
							SSEAlgorithm: "AES256",
						},
						BucketKeyEnabled: false,
					},
				},
			},
		}
		return output, nil
	case mockBucket2:
		output := &s3.GetBucketEncryptionOutput{
			ServerSideEncryptionConfiguration: &types.ServerSideEncryptionConfiguration{
				Rules: []types.ServerSideEncryptionRule{
					{
						ApplyServerSideEncryptionByDefault: &types.ServerSideEncryptionByDefault{
							SSEAlgorithm: "aws:kms",
						},
						BucketKeyEnabled: false,
					},
				},
			},
		}
		return output, nil
	default:
		ae := smithy.GenericAPIError{
			Code:    "ServerSideEncryptionConfigurationNotFoundError",
			Message: "No encryption set",
			Fault:   0,
		}
		return nil, &ae
	}
}

func (m mockS3APINew) GetBucketPolicy(_ context.Context,
	params *s3.GetBucketPolicyInput,
	_ ...func(*s3.Options)) (output *s3.GetBucketPolicyOutput, err error) {
	switch aws.ToString(params.Bucket) {
	case mockBucket1: // statement has the right format/properties
		policy := BucketPolicy{
			ID:      "Mock BucketPolicy ID 1234",
			Version: "2012-10-17",
			Statement: []Statement{
				{
					Action:    "s3:*",
					Effect:    "Deny",
					Resource:  []string{"*"},
					Condition: Condition{Bool{AwsSecureTransport: false}},
				}},
		}
		policyJson, err := json.Marshal(policy)
		if err != nil {
			log.Error(err)
		}
		output = &s3.GetBucketPolicyOutput{
			Policy: aws.String(string(policyJson)),
		}
	case mockBucket2: // JSON failure
		output = &s3.GetBucketPolicyOutput{
			Policy: aws.String(""),
		}
		err = nil
	case mockBucket3: // Effect audit instead of deny
		policy := BucketPolicy{
			ID:      "Mock BucketPolicy ID 1234",
			Version: "2012-10-17",
			Statement: []Statement{
				{
					Action:    "s3:*",
					Effect:    "Audit",
					Resource:  []string{"*"},
					Condition: Condition{Bool{AwsSecureTransport: false}},
				}},
		}
		policyJson, err := json.Marshal(policy)
		if err != nil {
			log.Error(err)
		}
		output = &s3.GetBucketPolicyOutput{
			Policy: aws.String(string(policyJson)),
		}
	default:
		output = nil
		err = &smithy.GenericAPIError{
			Code: "NoSuchBucketPolicy",
		}
	}
	return
}

func (m mockS3APINew) GetBucketLocation(_ context.Context,
	params *s3.GetBucketLocationInput,
	_ ...func(*s3.Options)) (output *s3.GetBucketLocationOutput, err error) {
	switch aws.ToString(params.Bucket) {
	case mockBucket1:
		output = &s3.GetBucketLocationOutput{
			LocationConstraint: mockBucket1Region,
		}
		err = nil
	case mockBucket2:
		output = &s3.GetBucketLocationOutput{
			LocationConstraint: mockBucket2Region,
		}
		err = nil
	default:
		output = nil
		err = &smithy.OperationError{
			ServiceID:     "MockS3API",
			OperationName: "GetBucketLocation",
			Err:           errors.New("no bucket location found for given bucket. Bucket does not exist"),
		}
	}
	return
}

func (m mockS3APINew) GetPublicAccessBlock(_ context.Context,
	_ *s3.GetPublicAccessBlockInput,
	_ ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {
	panic("implement me")
}

func (m mockS3APINew) GetBucketReplication(_ context.Context,
	_ *s3.GetBucketReplicationInput,
	_ ...func(*s3.Options)) (*s3.GetBucketReplicationOutput, error) {
	panic("implement me")
}

func (m mockS3APINew) GetBucketLifecycleConfiguration(_ context.Context,
	_ *s3.GetBucketLifecycleConfigurationInput,
	_ ...func(*s3.Options)) (*s3.GetBucketLifecycleConfigurationOutput, error) {
	panic("implement me")
}

type mockS3APIWitHErrors struct{}

func (m mockS3APIWitHErrors) ListBuckets(_ context.Context,
	_ *s3.ListBucketsInput,
	_ ...func(*s3.Options)) (o *s3.ListBucketsOutput, e error) {
	oe := &smithy.OperationError{
		ServiceID:     "MockS3API",
		OperationName: "ListBuckets",
		Err:           errors.New("failed to resolve service endpoint"),
	}
	return nil, oe
}

func (m mockS3APIWitHErrors) GetBucketEncryption(_ context.Context, _ *s3.GetBucketEncryptionInput, _ ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error) {
	ae := &smithy.GenericAPIError{
		Message: "failed to resolve service endpoint",
	}
	return nil, ae
}

func (m mockS3APIWitHErrors) GetBucketPolicy(_ context.Context, _ *s3.GetBucketPolicyInput, _ ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
	ae := &smithy.GenericAPIError{
		Message: "failed to resolve service endpoint",
	}
	return nil, ae
}

func (m mockS3APIWitHErrors) GetBucketLocation(_ context.Context, _ *s3.GetBucketLocationInput, _ ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
	oe := &smithy.OperationError{
		ServiceID:     "MockS3API",
		OperationName: "GetBucketLocation",
		Err:           errors.New("failed to resolve service endpoint"),
	}
	return nil, oe
}

func (m mockS3APIWitHErrors) GetPublicAccessBlock(_ context.Context, _ *s3.GetPublicAccessBlockInput, _ ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {
	oe := &smithy.OperationError{
		ServiceID:     "MockS3API",
		OperationName: "GetPublicAccessBlock",
		Err:           errors.New("failed to resolve service endpoint"),
	}
	return nil, oe
}

func (m mockS3APIWitHErrors) GetBucketReplication(_ context.Context, _ *s3.GetBucketReplicationInput, _ ...func(*s3.Options)) (*s3.GetBucketReplicationOutput, error) {
	oe := &smithy.OperationError{
		ServiceID:     "MockS3API",
		OperationName: "GetBucketReplication",
		Err:           errors.New("failed to resolve service endpoint"),
	}
	return nil, oe
}

func (m mockS3APIWitHErrors) GetBucketLifecycleConfiguration(_ context.Context, _ *s3.GetBucketLifecycleConfigurationInput, _ ...func(*s3.Options)) (*s3.GetBucketLifecycleConfigurationOutput, error) {
	oe := &smithy.OperationError{
		ServiceID:     "MockS3API",
		OperationName: "GetBucketLifecycleConfiguration",
		Err:           errors.New("failed to resolve service endpoint"),
	}
	return nil, oe
}

// TestGetBuckets tests the getBuckets method (with other form of mocking implementation)
func TestGetBuckets(t *testing.T) {
	d := awsS3Discovery{
		client:        mockS3APINew{},
		isDiscovering: false,
	}
	buckets, err := d.getBuckets()
	assert.Nil(t, err)

	log.Print("Testing number of buckets")
	assert.Equal(t, 2, len(buckets))

	log.Print("Testing name of first bucket")
	assert.Equal(t, mockBucket1, buckets[0].name)
	log.Print("Testing region of first bucket")
	assert.Equal(t, mockBucket1Region, buckets[0].region)
	log.Print("Testing endpoint of first bucket")
	assert.Equal(t, mockBucket1Endpoint, buckets[0].endpoint)
	log.Print("Testing creation time of first bucket")
	expectedCreationTime1, _ := time.Parse(time.RFC3339, mockBucket1CreationTime)
	assert.Equal(t, expectedCreationTime1.String(), buckets[0].creationTime.String())

	log.Print("Testing name of second bucket")
	assert.Equal(t, mockBucket2, buckets[1].name)
	log.Print("Testing region of second bucket")
	assert.Equal(t, mockBucket2Region, buckets[1].region)
	log.Print("Testing endpoint of second bucket")
	assert.Equal(t, mockBucket2Endpoint, buckets[1].endpoint)
	log.Print("Testing creation time of second bucket")
	expectedCreationTime2, _ := time.Parse(time.RFC3339, mockBucket2CreationTime)
	assert.Equal(t, expectedCreationTime2.String(), buckets[1].creationTime.String())

	// API error case
	d = awsS3Discovery{
		client:        mockS3APIWitHErrors{},
		isDiscovering: false,
	}

	_, err = d.getBuckets()
	assert.NotNil(t, err)
}

// TestGetEncryptionAtRest tests the getEncryptionAtRest method
func TestGetEncryptionAtRest(t *testing.T) {
	d := awsS3Discovery{
		client:        mockS3APINew{},
		isDiscovering: false,
	}

	// First case: SSE-S3 encryption
	encryptionAtRest, err := d.getEncryptionAtRest(mockBucket1)
	assert.Nil(t, err)
	assert.True(t, encryptionAtRest.Enabled)
	assert.Equal(t, "AES256", encryptionAtRest.Algorithm)
	assert.Equal(t, "SSE-S3", encryptionAtRest.KeyManager)

	// Second case: SSE-KMS encryption
	encryptionAtRest, err = d.getEncryptionAtRest(mockBucket2)
	assert.Nil(t, err)
	assert.True(t, encryptionAtRest.Enabled)
	assert.Equal(t, "aws:kms", encryptionAtRest.Algorithm)
	assert.Equal(t, "SSE-KMS", encryptionAtRest.KeyManager)

	// Third case: No encryption
	encryptionAtRest, err = d.getEncryptionAtRest("mockbucket3")
	assert.Nil(t, err)
	assert.False(t, encryptionAtRest.Enabled)

	// 4th case: Connection error
	d = awsS3Discovery{
		client:        mockS3APIWitHErrors{},
		isDiscovering: false,
	}
	_, err = d.getEncryptionAtRest("mockbucket4")
	assert.NotNil(t, err)
}

// TestGetTransportEncryption tests the getTransportEncryption method
func TestGetTransportEncryption(t *testing.T) {
	// Case 1: Connection error
	d := awsS3Discovery{
		client:        mockS3APIWitHErrors{},
		isDiscovering: false,
	}
	_, err := d.getTransportEncryption("")
	assert.NotNil(t, err)

	d = awsS3Discovery{
		client:        mockS3APINew{},
		isDiscovering: false,
	}

	// Case 2: Enforced
	encryptionAtTransit, err := d.getTransportEncryption(mockBucket1)
	assert.Nil(t, err)
	assert.True(t, encryptionAtTransit.Enabled)
	assert.Equal(t, "TLS1.2", encryptionAtTransit.TlsVersion)
	assert.True(t, encryptionAtTransit.Enforced)

	// Case 3: JSON failure
	encryptionAtTransit, err = d.getTransportEncryption(mockBucket2)
	assert.NotNil(t, err)
	assert.True(t, encryptionAtTransit.Enabled)
	assert.Equal(t, "TLS1.2", encryptionAtTransit.TlsVersion)
	assert.False(t, encryptionAtTransit.Enforced)

	// Case 4: Not enforced
	encryptionAtTransit, err = d.getTransportEncryption(mockBucket3)
	assert.Nil(t, err)
	assert.True(t, encryptionAtTransit.Enabled)
	assert.Equal(t, "TLS1.2", encryptionAtTransit.TlsVersion)
	assert.False(t, encryptionAtTransit.Enforced)

	// Case 5: No bucket policy == not enforced
	encryptionAtTransit, err = d.getTransportEncryption("")
	assert.Nil(t, err)
	assert.True(t, encryptionAtTransit.Enabled)
	assert.Equal(t, "TLS1.2", encryptionAtTransit.TlsVersion)
	assert.False(t, encryptionAtTransit.Enforced)

}

// TestGetRegion tests the getRegion method
func TestGetRegion(t *testing.T) {
	d := awsS3Discovery{
		client:        mockS3APINew{},
		isDiscovering: false,
	}
	actualRegion, err := d.getRegion(mockBucket1)
	assert.Nil(t, err)
	assert.Equal(t, mockBucket1Region, actualRegion)

	actualRegion, err = d.getRegion(mockBucket2)
	assert.Nil(t, err)
	assert.Equal(t, mockBucket2Region, actualRegion)

	// Error case
	_, err = d.getRegion("mockbucketNotAvailable")
	assert.NotNil(t, err)

}

// TestName tests the Name method
func TestName(t *testing.T) {
	d := awsS3Discovery{
		client:        mockS3APINew{},
		isDiscovering: false,
	}

	assert.Equal(t, "AWS Blob Storage", d.Name())
}

// TestList tests the List method
func TestList(t *testing.T) {
	d := awsS3Discovery{
		client:        mockS3APINew{},
		isDiscovering: false,
	}
	resources, err := d.List()
	assert.NotNil(t, err, "EXPECTED error because MockBucket2 should throw JSON error. But GOT no error")

	log.Println("Testing number of resources (buckets)")
	assert.Equal(t, 1, len(resources))

	expectedResourceNames := []string{mockBucket1, "mockbucket2", "mockbucket3"}
	//expectedResourceAtRestEncryptions := []bool{true, true, false}
	//expectedResourceTransportEncryptions := []bool{true, false, false}
	for i, r := range resources {
		log.Println("Testing name for resource (bucket)", i+1)
		assert.Equal(t, expectedResourceNames[i], r.GetName())
		log.Println("Testing type of resource", i+1)
		assert.True(t, r.HasType("ObjectStorage"))
	}
}
