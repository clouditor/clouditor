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

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

const (
	mockBucket1             = "mockbucket1"
	mockBucket1Region       = "eu-central-1"
	mockBucket1Endpoint     = "https://" + mockBucket1 + ".s3." + mockBucket1Region + ".amazonaws.com"
	mockBucket1CreationTime = "2012-11-01T22:08:41+00:00"
	mockBucket2             = "mockbucket2"
	mockBucket2Region       = mockBucket1Region
	mockBucket2Endpoint     = "https://" + mockBucket2 + ".s3." + mockBucket1Region + ".amazonaws.com"
	mockBucket2CreationTime = "2013-12-02T22:08:41+00:00"
	mockBucket2KeyId        = "1234abcd-12ab-34cd-56ed-1234567890ab"
	mockBucket3             = "mockbucket3"
)

type mockS3APINew struct{}

func (mockS3APINew) ListBuckets(_ context.Context,
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

func (mockS3APINew) GetBucketEncryption(_ context.Context,
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
						BucketKeyEnabled: util.Ref(false),
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
							SSEAlgorithm:   "aws:kms",
							KMSMasterKeyID: aws.String(mockBucket2KeyId),
						},
						BucketKeyEnabled: util.Ref(false),
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

func (mockS3APINew) GetBucketPolicy(_ context.Context,
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
					Resource:  "*",
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
					Resource:  "*",
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

func (mockS3APINew) GetBucketLocation(_ context.Context,
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

func (mockS3APINew) GetPublicAccessBlock(_ context.Context,
	_ *s3.GetPublicAccessBlockInput,
	_ ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {
	panic("implement me")
}

func (mockS3APINew) GetBucketReplication(_ context.Context,
	_ *s3.GetBucketReplicationInput,
	_ ...func(*s3.Options)) (*s3.GetBucketReplicationOutput, error) {
	panic("implement me")
}

func (mockS3APINew) GetBucketLifecycleConfiguration(_ context.Context,
	_ *s3.GetBucketLifecycleConfigurationInput,
	_ ...func(*s3.Options)) (*s3.GetBucketLifecycleConfigurationOutput, error) {
	panic("implement me")
}

type mockS3APIWitHErrors struct{}

func (mockS3APIWitHErrors) ListBuckets(_ context.Context,
	_ *s3.ListBucketsInput,
	_ ...func(*s3.Options)) (o *s3.ListBucketsOutput, e error) {
	oe := &smithy.OperationError{
		ServiceID:     "MockS3API",
		OperationName: "ListBuckets",
		Err:           errors.New("failed to resolve service endpoint"),
	}
	return nil, oe
}

func (mockS3APIWitHErrors) GetBucketEncryption(_ context.Context, _ *s3.GetBucketEncryptionInput, _ ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error) {
	ae := &smithy.GenericAPIError{
		Message: "failed to resolve service endpoint",
	}
	return nil, ae
}

func (mockS3APIWitHErrors) GetBucketPolicy(_ context.Context, _ *s3.GetBucketPolicyInput, _ ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
	ae := &smithy.GenericAPIError{
		Message: "failed to resolve service endpoint",
	}
	return nil, ae
}

func (mockS3APIWitHErrors) GetBucketLocation(_ context.Context, _ *s3.GetBucketLocationInput, _ ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
	oe := &smithy.OperationError{
		ServiceID:     "MockS3API",
		OperationName: "GetBucketLocation",
		Err:           errors.New("failed to resolve service endpoint"),
	}
	return nil, oe
}

func (mockS3APIWitHErrors) GetPublicAccessBlock(_ context.Context, _ *s3.GetPublicAccessBlockInput, _ ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {
	oe := &smithy.OperationError{
		ServiceID:     "MockS3API",
		OperationName: "GetPublicAccessBlock",
		Err:           errors.New("failed to resolve service endpoint"),
	}
	return nil, oe
}

func (mockS3APIWitHErrors) GetBucketReplication(_ context.Context, _ *s3.GetBucketReplicationInput, _ ...func(*s3.Options)) (*s3.GetBucketReplicationOutput, error) {
	oe := &smithy.OperationError{
		ServiceID:     "MockS3API",
		OperationName: "GetBucketReplication",
		Err:           errors.New("failed to resolve service endpoint"),
	}
	return nil, oe
}

func (mockS3APIWitHErrors) GetBucketLifecycleConfiguration(_ context.Context, _ *s3.GetBucketLifecycleConfigurationInput, _ ...func(*s3.Options)) (*s3.GetBucketLifecycleConfigurationOutput, error) {
	oe := &smithy.OperationError{
		ServiceID:     "MockS3API",
		OperationName: "GetBucketLifecycleConfiguration",
		Err:           errors.New("failed to resolve service endpoint"),
	}
	return nil, oe
}

// TestGetBuckets tests the getBuckets method (with other form of mocking implementation)
func TestAwsS3Discovery_getBuckets(t *testing.T) {
	d := awsS3Discovery{
		storageAPI:    mockS3APINew{},
		isDiscovering: false,
		awsConfig: &Client{
			cfg:       aws.Config{Region: mockBucket1Region},
			accountID: nil,
		},
	}
	buckets, err := d.getBuckets()
	assert.NoError(t, err)

	log.Print("Testing number of buckets")
	// We fetch buckets currently only of users region
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
		storageAPI:    mockS3APIWitHErrors{},
		isDiscovering: false,
	}

	_, err = d.getBuckets()
	assert.Error(t, err)
}

// TestGetEncryptionAtRest tests the getEncryptionAtRest method
func TestAwsS3Discovery_getEncryptionAtRest(t *testing.T) {
	var (
		encryptionAtRest    *ontology.AtRestEncryption
		managedEncryption   *ontology.ManagedKeyEncryption
		customerEncryption  *ontology.CustomerKeyEncryption
		rawEncryptionAtRest *s3.GetBucketEncryptionOutput

		err           error
		mockAccountID = "123456789"
	)

	d := awsS3Discovery{
		storageAPI:    mockS3APINew{},
		isDiscovering: false,
		awsConfig: &Client{
			cfg:       aws.Config{},
			accountID: aws.String(mockAccountID),
		},
	}

	// First case: SSE-S3 encryption
	encryptionAtRest, rawEncryptionAtRest, err = d.getEncryptionAtRest(&bucket{name: mockBucket1})
	assert.NoError(t, err)
	managedEncryption = encryptionAtRest.GetManagedKeyEncryption()
	assert.True(t, managedEncryption.Enabled)
	assert.Equal(t, "AES256", managedEncryption.Algorithm)
	assert.NotEmpty(t, rawEncryptionAtRest)

	// Second case: SSE-KMS encryption
	encryptionAtRest, rawEncryptionAtRest, err = d.getEncryptionAtRest(&bucket{name: mockBucket2, region: mockBucket2Region})
	customerEncryption = encryptionAtRest.GetCustomerKeyEncryption()
	assert.NoError(t, err)
	assert.True(t, customerEncryption.Enabled)
	assert.Equal(t, "", customerEncryption.Algorithm)
	assert.Equal(t, "arn:aws:kms:"+mockBucket2Region+":"+mockAccountID+":key/"+mockBucket2KeyId, customerEncryption.KeyUrl)
	assert.NotEmpty(t, rawEncryptionAtRest)

	// Third case: No encryption
	encryptionAtRest, rawEncryptionAtRest, err = d.getEncryptionAtRest(&bucket{name: "mockbucket3"})
	assert.NoError(t, err)
	assert.Nil(t, encryptionAtRest)
	assert.Empty(t, rawEncryptionAtRest)

	// 4th case: Connection error
	d = awsS3Discovery{
		storageAPI:    mockS3APIWitHErrors{},
		isDiscovering: false,
	}
	_, _, err = d.getEncryptionAtRest(&bucket{name: "mockbucket4"})
	assert.Error(t, err)
}

// TestGetTransportEncryption tests the getTransportEncryption method
func TestAwsS3Discovery_getTransportEncryption(t *testing.T) {
	var rawBucketPolicy *s3.GetBucketPolicyOutput

	// Case 1: Connection error
	d := awsS3Discovery{
		storageAPI:    mockS3APIWitHErrors{},
		isDiscovering: false,
	}
	_, rawBucketPolicy, err := d.getTransportEncryption("")
	assert.Error(t, err)
	assert.Empty(t, rawBucketPolicy)

	d = awsS3Discovery{
		storageAPI:    mockS3APINew{},
		isDiscovering: false,
	}

	// Case 2: Enforced
	encryptionAtTransit, rawBucketPolicy, err := d.getTransportEncryption(mockBucket1)
	assert.NoError(t, err)
	assert.True(t, encryptionAtTransit.Enabled)
	assert.Equal(t, float32(1.2), encryptionAtTransit.ProtocolVersion)
	assert.True(t, encryptionAtTransit.Enforced)
	assert.NotEmpty(t, rawBucketPolicy)

	// Case 3: JSON failure
	encryptionAtTransit, rawBucketPolicy, err = d.getTransportEncryption(mockBucket2)
	assert.Error(t, err)
	assert.Nil(t, encryptionAtTransit)
	assert.NotEmpty(t, rawBucketPolicy)

	// Case 4: Not enforced
	encryptionAtTransit, rawBucketPolicy, err = d.getTransportEncryption(mockBucket3)
	assert.NoError(t, err)
	assert.True(t, encryptionAtTransit.Enabled)
	assert.Equal(t, float32(1.2), encryptionAtTransit.ProtocolVersion)
	assert.False(t, encryptionAtTransit.Enforced)
	assert.NotEmpty(t, rawBucketPolicy)

	// Case 5: No bucket policy == not enforced
	encryptionAtTransit, rawBucketPolicy, err = d.getTransportEncryption("")
	assert.NoError(t, err)
	assert.True(t, encryptionAtTransit.Enabled)
	assert.Equal(t, float32(1.2), encryptionAtTransit.ProtocolVersion)
	assert.False(t, encryptionAtTransit.Enforced)
	assert.Empty(t, rawBucketPolicy)
}

// TestGetRegion tests the getRegion method
func TestAwsS3Discovery_getRegion(t *testing.T) {
	d := awsS3Discovery{
		storageAPI:    mockS3APINew{},
		isDiscovering: false,
	}
	actualRegion, rawRegion, err := d.getRegion(mockBucket1)
	assert.NoError(t, err)
	assert.NotEmpty(t, rawRegion)
	assert.Equal(t, mockBucket1Region, actualRegion)

	actualRegion, rawRegion, err = d.getRegion(mockBucket2)
	assert.NoError(t, err)
	assert.NotEmpty(t, rawRegion)
	assert.Equal(t, mockBucket2Region, actualRegion)

	// Error case
	_, rawRegion, err = d.getRegion("mockbucketNotAvailable")
	assert.Empty(t, rawRegion)
	assert.Error(t, err)

}

// TestName tests the Name method
func TestAwsS3Discovery_Name(t *testing.T) {
	d := awsS3Discovery{
		storageAPI:    mockS3APINew{},
		isDiscovering: false,
	}

	assert.Equal(t, "AWS Blob Storage", d.Name())
}

// TestList tests the List method
func TestAwsS3Discovery_List(t *testing.T) {
	d := awsS3Discovery{
		storageAPI:    mockS3APINew{},
		isDiscovering: false,
		awsConfig: &Client{
			cfg: aws.Config{
				Region:        mockBucket1Region,
				Credentials:   nil,
				HTTPClient:    nil,
				Retryer:       nil,
				ConfigSources: nil,
				APIOptions:    nil,
				Logger:        nil,
				ClientLogMode: 0,
			},
			accountID: aws.String("123456789"),
		},
	}
	resources, err := d.List()
	assert.NotNil(t, err)

	log.Println("Testing number of resources (buckets)")
	assert.Equal(t, 2, len(resources))

	expectedResourceNames := []string{mockBucket1, "mockbucket2", "mockbucket3"}

	// Check first element: voc.ObjectStorage
	log.Println("Testing name for resource (bucket)", 1)
	assert.Equal(t, expectedResourceNames[0], resources[0].GetName())
	log.Println("Testing type of resource", 1)
	assert.True(t, ontology.HasType(resources[0], "ObjectStorage"))
	expectedRaw := "{\"**s3.GetBucketEncryptionOutput\":[{\"ServerSideEncryptionConfiguration\":{\"Rules\":[{\"ApplyServerSideEncryptionByDefault\":{\"SSEAlgorithm\":\"AES256\",\"KMSMasterKeyID\":null},\"BucketKeyEnabled\":false}]},\"ResultMetadata\":{}}],\"**s3.GetBucketPolicyOutput\":[{\"Policy\":\"{\\\"id\\\":\\\"Mock BucketPolicy ID 1234\\\",\\\"Version\\\":\\\"2012-10-17\\\",\\\"Statement\\\":[{\\\"Action\\\":\\\"s3:*\\\",\\\"Effect\\\":\\\"Deny\\\",\\\"Resource\\\":\\\"*\\\",\\\"Condition\\\":{\\\"aws:SecureTransport\\\":false}}]}\",\"ResultMetadata\":{}}],\"*[]interface {}\":[[{\"BucketArn\":null,\"BucketRegion\":null,\"CreationDate\":\"2012-11-01T22:08:41Z\",\"Name\":\"mockbucket1\"},{\"LocationConstraint\":\"eu-central-1\",\"ResultMetadata\":{}}]],\"*aws.bucket\":[{}]}"
	assert.Equal(t, expectedRaw, resources[0].GetRaw())

	// Check second element: voc.ObjectStorageService
	log.Println("Testing name for resource (bucket)", 2)
	assert.Equal(t, expectedResourceNames[0], resources[1].GetName())
	log.Println("Testing type of resource", 2)
	assert.True(t, ontology.HasType(resources[1], "ObjectStorageService"))
}

func Test_awsS3Discovery_TargetOfEvaluationID(t *testing.T) {
	type fields struct {
		storageAPI    S3API
		isDiscovering bool
		awsConfig     *Client
		ctID          string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy path",
			fields: fields{
				ctID: testdata.MockTargetOfEvaluationID1,
			},
			want: testdata.MockTargetOfEvaluationID1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &awsS3Discovery{
				storageAPI:    tt.fields.storageAPI,
				isDiscovering: tt.fields.isDiscovering,
				awsConfig:     tt.fields.awsConfig,
				ctID:          tt.fields.ctID,
			}
			if got := d.TargetOfEvaluationID(); got != tt.want {
				t.Errorf("awsS3Discovery.TargetOfEvaluationID() = %v, want %v", got, tt.want)
			}
		})
	}
}
