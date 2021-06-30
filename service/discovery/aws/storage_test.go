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
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"strconv"
	"testing"
	"time"
)

// ToDo: Needs re-writing

var cfg aws.Config

func init() {
	cfg = NewAwsDiscovery().cfg
}

func TestGetS3ServiceClient(t *testing.T) {
	// ToDo
	client := NewAwsStorageDiscovery(cfg)
	if client.client == nil {
		t.Errorf("Connection failed. Credentials are nil.")
	}
	//_, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	//if err != nil {
	//	t.Fatalf("Error: %v", err)
	//}
}

// ToDo: Begin mock stuff

//func TestListS3(t *testing.T) {
//	mockDisplayName := "MockDisplayName"
//	mockId := "MockId1234"
//	cases := []struct {
//		client func(t *testing.T) S3ListBucketsAPI
//		//expect []byte
//	}{
//		{
//
//			client: func(t *testing.T) S3ListBucketsAPI {
//				return mockS3API(func(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
//					t.Helper()
//					return &s3.ListBucketsOutput{
//							Buckets: []types.Bucket{
//								types.Bucket{
//									CreationDate: nil,
//									Name:         aws.String("FirstMockBucket"),
//								},
//								types.Bucket{
//									CreationDate: nil,
//									Name:         aws.String("SecondMockBucket"),
//								},
//							},
//							Owner: &types.Owner{
//								DisplayName: &mockDisplayName,
//								ID:          &mockId,
//							},
//						},
//						nil
//				})
//			},
//		},
//	}
//
//	for i, tt := range cases {
//		t.Run(strconv.Itoa(i), func(t *testing.T) {
//			buckets := List(tt.client(t))
//			if len(buckets.Buckets) == 0 {
//				t.Fatal("Buckets empty but shouldn't be.")
//			}
//			if o := *buckets.Owner.ID; o != mockId {
//				t.Fatalf("expected %v, but got %v", mockId, o)
//			}
//		})
//	}
//}

//func TestAreBucketsEncrypted(t *testing.T) {
//	mockDisplayName := "MockDisplayName"
//	mockId := "MockId1234"
//	cases := []struct {
//		client func(t *testing.T) S3ListBucketsAPI
//		//expect []byte
//	}{
//		{
//
//			client: func(t *testing.T) S3ListBucketsAPI {
//				return mockS3API(func(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
//					t.Helper()
//					return &s3.ListBucketsOutput{
//							Buckets: []types.Bucket{
//								types.Bucket{
//									CreationDate: nil,
//									Name:         aws.String("FirstMockBucket"),
//								},
//								types.Bucket{
//									CreationDate: nil,
//									Name:         aws.String("SecondMockBucket"),
//								},
//							},
//							Owner: &types.Owner{
//								DisplayName: &mockDisplayName,
//								ID:          &mockId,
//							},
//						},
//						nil
//				})
//			},
//		},
//	}
//
//	for i, tt := range cases {
//		t.Run(strconv.Itoa(i), func(t *testing.T) {
//			buckets := List(tt.client(t))
//			if len(buckets.Buckets) == 0 {
//				t.Fatal("Buckets empty but shouldn't be.")
//			}
//			if o := *buckets.Owner.ID; o != mockId {
//				t.Fatalf("expected %v, but got %v", mockId, o)
//			}
//		})
//	}
//}

// ToDo: End mock stuff

// Mock s3 service and methods
type mockS3API func(ctx context.Context,
	params *s3.ListBucketsInput,
	optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)

func (m mockS3API) GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
	panic("implement me")
}

func (m mockS3API) GetBucketPolicy(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
	panic("implement me")
}

func (m mockS3API) GetBucketEncryption(ctx context.Context, params *s3.GetBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error) {
	panic("implement me")
}

func (m mockS3API) GetPublicAccessBlock(ctx context.Context, params *s3.GetPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {
	panic("implement me")
}

func (m mockS3API) GetBucketReplication(ctx context.Context, params *s3.GetBucketReplicationInput, optFns ...func(*s3.Options)) (*s3.GetBucketReplicationOutput, error) {
	panic("implement me")
}

func (m mockS3API) GetBucketLifecycleConfiguration(ctx context.Context, params *s3.GetBucketLifecycleConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLifecycleConfigurationOutput, error) {
	panic("implement me")
}

func (m mockS3API) ListBuckets(ctx context.Context,
	params *s3.ListBucketsInput,
	optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {

	// ToDo what does "m(xxx)" mean? What does it return from mockS3API?
	return m(ctx, params, optFns...)
}

type mockS3APINew struct{}

func (m mockS3APINew) ListBuckets(ctx context.Context,
	params *s3.ListBucketsInput,
	optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {

	output := &s3.ListBucketsOutput{Buckets: []types.Bucket{
		{
			Name:         aws.String("mockbucket1"),
			CreationDate: aws.Time(time.Now()),
		},
		{
			Name:         aws.String("mockbucket2"),
			CreationDate: aws.Time(time.Now().Add(-time.Hour * 24)),
		},
	}}
	return output, nil
}

func (m mockS3APINew) GetBucketEncryption(ctx context.Context, params *s3.GetBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error) {
	// ToDo: Switch between different buckets (params) -> different SSEAlgorithm and KeyManager
	switch aws.ToString(params.Bucket) {
	case "mockbucket1":
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
	case "mockbucket2":
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
		oe := &smithy.OperationError{
			ServiceID:     "MockS3API",
			OperationName: "GetBucketEncryption",
			Err:           errors.New("failed to resolve service endpoint"),
		}
		return nil, oe
	}
}

func (m mockS3APINew) GetBucketPolicy(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (output *s3.GetBucketPolicyOutput, err error) {
	switch aws.ToString(params.Bucket) {
	case "mockbucket1": // statement has the right format/properties
		policy := Policy{
			ID:      "Mock Policy ID 1234",
			Version: "2012-10-17",
			Statement: []Statement{
				{
					Action:    "s3:*",
					Effect:    "Deny",
					Resource:  []string{"*"},
					Condition: Condition{Bool{awsSecureTransport: false}},
				}},
		}
		policyJson, err := json.Marshal(policy)
		if err != nil {
			log.Error(err)
		}
		output = &s3.GetBucketPolicyOutput{
			Policy: aws.String(string(policyJson)),
		}
		err = nil
	case "mockbucket2": // JSON failure
		output = &s3.GetBucketPolicyOutput{
			Policy: aws.String(""),
		}
		err = nil
	case "mockbucket3": // Effect audit instead of deny
		policy := Policy{
			ID:      "Mock Policy ID 1234",
			Version: "2012-10-17",
			Statement: []Statement{
				{
					Action:    "s3:*",
					Effect:    "Audit",
					Resource:  []string{"*"},
					Condition: Condition{Bool{awsSecureTransport: false}},
				}},
		}
		policyJson, err := json.Marshal(policy)
		if err != nil {
			log.Error(err)
		}
		output = &s3.GetBucketPolicyOutput{
			Policy: aws.String(string(policyJson)),
		}
		err = nil
	}
	return
}

func (m mockS3APINew) GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (output *s3.GetBucketLocationOutput, err error) {
	switch aws.ToString(params.Bucket) {
	case "Mock Bucket 1":
		output = &s3.GetBucketLocationOutput{
			LocationConstraint: "eu-central-1",
		}
		err = nil
	default:
		output = nil
		err = errors.New("MockS3API: No bucket policy found for given bucket or bucket does not exist")
	}
	return
}

func (m mockS3APINew) GetPublicAccessBlock(ctx context.Context, params *s3.GetPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {
	panic("implement me")
}

func (m mockS3APINew) GetBucketReplication(ctx context.Context, params *s3.GetBucketReplicationInput, optFns ...func(*s3.Options)) (*s3.GetBucketReplicationOutput, error) {
	panic("implement me")
}

func (m mockS3APINew) GetBucketLifecycleConfiguration(ctx context.Context, params *s3.GetBucketLifecycleConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLifecycleConfigurationOutput, error) {
	panic("implement me")
}

type mockS3APIWitHErrors struct{}

func (m mockS3APIWitHErrors) ListBuckets(ctx context.Context,
	params *s3.ListBucketsInput,
	optFns ...func(*s3.Options)) (o *s3.ListBucketsOutput, e error) {
	oe := &smithy.OperationError{
		ServiceID:     "MockS3API",
		OperationName: "GetBucketEncryption",
		Err:           errors.New("failed to resolve service endpoint"),
	}
	return nil, oe
}

func (m mockS3APIWitHErrors) GetBucketEncryption(ctx context.Context, params *s3.GetBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error) {
	oe := &smithy.OperationError{
		ServiceID:     "MockS3API",
		OperationName: "GetBucketEncryption",
		Err:           errors.New("failed to resolve service endpoint"),
	}
	return nil, oe
}

func (m mockS3APIWitHErrors) GetBucketPolicy(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
	oe := &smithy.OperationError{
		ServiceID:     "MockS3API",
		OperationName: "GetBucketEncryption",
		Err:           errors.New("failed to resolve service endpoint"),
	}
	return nil, oe
}

func (m mockS3APIWitHErrors) GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
	oe := &smithy.OperationError{
		ServiceID:     "MockS3API",
		OperationName: "GetBucketEncryption",
		Err:           errors.New("failed to resolve service endpoint"),
	}
	return nil, oe
}

func (m mockS3APIWitHErrors) GetPublicAccessBlock(ctx context.Context, params *s3.GetPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {
	oe := &smithy.OperationError{
		ServiceID:     "MockS3API",
		OperationName: "GetBucketEncryption",
		Err:           errors.New("failed to resolve service endpoint"),
	}
	return nil, oe
}

func (m mockS3APIWitHErrors) GetBucketReplication(ctx context.Context, params *s3.GetBucketReplicationInput, optFns ...func(*s3.Options)) (*s3.GetBucketReplicationOutput, error) {
	oe := &smithy.OperationError{
		ServiceID:     "MockS3API",
		OperationName: "GetBucketEncryption",
		Err:           errors.New("failed to resolve service endpoint"),
	}
	return nil, oe
}

func (m mockS3APIWitHErrors) GetBucketLifecycleConfiguration(ctx context.Context, params *s3.GetBucketLifecycleConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLifecycleConfigurationOutput, error) {
	oe := &smithy.OperationError{
		ServiceID:     "MockS3API",
		OperationName: "GetBucketEncryption",
		Err:           errors.New("failed to resolve service endpoint"),
	}
	return nil, oe
}

// TestGetBucketsNew tests the getBuckets method (with other form of mocking implementation)
// ToDo: Its simpler and shorter but I would like the other one more (with "cases")
func TestGetBucketsNew(t *testing.T) {
	// ToDo: It is not better to use a pointer (&awsS3Discovery), is it?
	// ToDo: I should mock the initialization of d as well? Meaning creating a init function in storage.go and mock it
	d := awsS3Discovery{
		client:        mockS3APINew{},
		buckets:       nil,
		isDiscovering: false,
	}
	d.getBuckets()
	log.Print("Testing number of buckets")
	if e, a := 2, len(d.buckets); e != a {
		t.Error("EXPECTED:", e, "GOT:", a)
	}
	log.Print("Testing name of first bucket")
	if e, a := "mockbucket1", d.buckets[0].name; e != a {
		t.Error("EXPECTED:", e, "GOT:", a)
	}

	d = awsS3Discovery{
		client:        mockS3APIWitHErrors{},
		buckets:       nil,
		isDiscovering: false,
	}

	d.getBuckets()
	if d.buckets != nil {
		t.Error("EXPECTED empty list of buckets")
	}
}

// TestGetEncryptionAtRest tests the getEncryptionAtRest method
func TestGetEncryptionAtRest(t *testing.T) {
	d := awsS3Discovery{
		client:        mockS3APINew{},
		buckets:       nil,
		isDiscovering: false,
	}
	// First case
	isEncrypted, algorithm, manager := d.getEncryptionAtRest("mockbucket1")
	if isEncrypted == false {
		t.Error("Expected:", true, ".Got:", isEncrypted)
	}
	if e := "AES256"; algorithm != e {
		t.Error("Expected:", e, ".Got:", algorithm)
	}
	if e := "SSE-S3"; manager != e {
		t.Error("Expected:", e, ".Got:", manager)
	}
	// Second case
	isEncrypted, algorithm, manager = d.getEncryptionAtRest("mockbucket2")
	if isEncrypted == false {
		t.Error("Expected:", true, ".Got:", isEncrypted)
	}
	if e := "aws:kms"; algorithm != e {
		t.Error("Expected:", e, ".Got:", algorithm)
	}
	if e := "SSE-KMS"; manager != e {
		t.Error("Expected:", e, ".Got:", manager)
	}

	// Third case
	isEncrypted, algorithm, manager = d.getEncryptionAtRest("Mock Bucket with no encryption")
	if isEncrypted == true {
		t.Error("Expected:", false, ".Got:", isEncrypted)
	}
}

// TestGetTransportEncryption tests the getTransportEncryption method
// ToDo: Check Test again
func TestGetTransportEncryption(t *testing.T) {
	// Case 1 (error case):
	d := awsS3Discovery{
		client:        mockS3APIWitHErrors{},
		buckets:       nil,
		isDiscovering: false,
	}
	if isEncrypted, algorithm, enforced, version := d.getTransportEncryption(""); isEncrypted == true || algorithm != "" || enforced == true || version != "" {
		t.Errorf("Expected isEncrypted: %v, algorithm: %v, enforced: %v, version: %v."+
			"Got: %v, %v, %v, %v", false, "", false, "", isEncrypted, algorithm, enforced, version)
	}

	d = awsS3Discovery{
		client:        mockS3APINew{},
		buckets:       nil,
		isDiscovering: false,
	}

	// Case 2
	if isEncrypted, algorithm, enforced, version := d.getTransportEncryption("mockbucket1"); isEncrypted == false || algorithm != "TLS" || enforced == false || version != "1.2" {
		t.Errorf("Expected isEncrypted: %v, algorithm: %v, enforced: %v, version: %v."+
			"Got: %v, %v, %v, %v", true, "TLS", true, "1.2", isEncrypted, algorithm, enforced, version)
	}

	// Case 3: JSON failure
	if isEncrypted, algorithm, enforced, version := d.getTransportEncryption("mockbucket2"); isEncrypted == true || algorithm != "" || enforced == true || version != "" {
		t.Errorf("Expected isEncrypted: %v, algorithm: %v, enforced: %v, version: %v."+
			"Got: %v, %v, %v, %v", false, "", false, "", isEncrypted, algorithm, enforced, version)
	}

	// Case 4:
	if isEncrypted, algorithm, enforced, version := d.getTransportEncryption("mockbucket3"); isEncrypted == true || algorithm != "" || enforced == true || version != "" {
		t.Errorf("Expected isEncrypted: %v, algorithm: %v, enforced: %v, version: %v."+
			"Got: %v, %v, %v, %v", false, "", false, "", isEncrypted, algorithm, enforced, version)
	}

}

// TestGetGeoLocation tests the getGeoLocation method
func TestGetGeoLocation(t *testing.T) {
	d := awsS3Discovery{
		client:        mockS3APINew{},
		buckets:       nil,
		isDiscovering: false,
	}
	if e, a := "eu-central-1", d.getGeoLocation("Mock Bucket 1"); e != a {
		t.Error("Expected: ", e, ". Got:", a)
	}
	if e, a := "", d.getGeoLocation("Mock Bucket 2"); e != a {
		t.Error("Expected: ", "(empty string)", ". Got:", a)
	}

}

// TestGetPublicAccessBlockConfiguration tests the getPublicAccessBlockConfiguration method
// ToDo: When needed or I have the time (since this check is not necessary now)
func TestGetPublicAccessBlockConfiguration(t *testing.T) {
	d := awsS3Discovery{
		client:        mockS3APINew{},
		buckets:       nil,
		isDiscovering: false,
	}
	log.Print(d)
	//d.getPublicAccessBlockConfiguration("Mock Bucket 1")
}

// TestGetBuckets tests the getBuckets method
// ToDo: Remove test when deciding to use the new mock implementation variant
func TestGetBuckets(t *testing.T) {
	cases := []struct {
		client func(t *testing.T) S3API
	}{
		{
			// ToDo: Can i put multiple functions into client here? (If I test every function separately, maybe I dont need it)
			client: func(t *testing.T) S3API {
				return mockS3API(func(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
					t.Helper()
					return &s3.ListBucketsOutput{
						Buckets: []types.Bucket{
							{
								CreationDate: nil,
								Name:         aws.String("AWS Mock Bucket 1"),
							},
							{
								CreationDate: nil,
								Name:         aws.String("AWS Mock Bucket 2"),
							},
						},
						Owner: &types.Owner{
							DisplayName: aws.String("Mock Display Name"),
							ID:          aws.String("MockId1234"),
						},
					}, nil
				})
			},
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			d := &awsS3Discovery{
				client:        tt.client(t),
				buckets:       nil,
				isDiscovering: false,
			}
			d.getBuckets()
			log.Println("Here are the buckets: ", d.buckets)
		})
	}
}

// TestGetEncryptionAtRest tests the getEncryptionAtRest method
//func TestGetEncryptionAtRest(t *testing.T) {
//	cases := []struct {
//		client func(t *testing.T) S3API
//	}{
//		{
//			client: func(t *testing.T) S3API {
//				return mockS3API(func(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
//
//				})
//			}}}
//}

// ToDo: Works with my credentials -> Mock it
func TestCheckEncryption_withCredentials(t *testing.T) {
	d := NewAwsStorageDiscovery(cfg)
	d.getBuckets()
	for i, bucket := range d.buckets {
		isEncrypted, _, _ := d.getEncryptionAtRest(bucket.name)
		if i == 0 && isEncrypted {
			fmt.Printf("Expected that bucket %v is not encrypted, but it is", bucket)
		} else if i == 1 && !isEncrypted {
			t.Errorf("Expected that bucket %v is encrypted, but it is not", bucket)
		}
	}
}

func TestCheckPublicAccessBlockConfiguration(t *testing.T) {
	d := NewAwsStorageDiscovery(NewAwsDiscovery().cfg)
	d.getBuckets()
	for _, bucket := range d.buckets {
		if d.getPublicAccessBlockConfiguration(bucket.name) == false {
			t.Fatalf("Expected no public access of bucket. But public access is enabled for %v.", bucket)
		}
	}
}

func TestCheckBucketReplication(t *testing.T) {
	d := NewAwsStorageDiscovery(NewAwsDiscovery().cfg)
	d.getBuckets()
	for _, bucket := range d.buckets {
		if d.checkBucketReplication(bucket.name) == true {
			t.Fatalf("Expected no replication setting for bucket. But replication is set for bucket '%v'.", bucket)
		}
	}
}

func TestCheckLifeCycleConfiguration(t *testing.T) {
	d := NewAwsStorageDiscovery(NewAwsDiscovery().cfg)
	d.getBuckets()
	for _, bucket := range d.buckets {
		if d.checkLifeCycleConfiguration(bucket.name) == true {
			t.Fatalf("Expected no life cycle configuration setting for bucket. But it is set for bucket '%v'.", bucket)
		}
	}
}

//func TestGetObjectsOfBucket_whenNotEmpty(t *testing.T) {
//	if bucketObjects := GetObjectsOfBucket(os.Getenv("TESTBUCKET")); len(bucketObjects.Contents) == 0 {
//		t.Errorf("No buckets found")
//	}
//}

// TestName tests the Name method
func TestName(t *testing.T) {
	d := awsS3Discovery{
		client:        mockS3APINew{},
		buckets:       nil,
		isDiscovering: false,
	}
	log.Println("Testing the name of the AWS Blob Storage Discovery")
	if e, a := "Aws Blob Storage", d.Name(); e != a {
		t.Error("EXPECTED:", e, "GOT", a)
	}
}

// TestList tests the List method
func TestList(t *testing.T) {
	d := awsS3Discovery{
		client:        mockS3APINew{},
		buckets:       nil,
		isDiscovering: false,
	}
	resources, err := d.List()
	if err != nil {
		t.Error("Error occurred:", err)
	}
	log.Println("Testing number of resources (buckets)")
	if a := len(resources); a != 2 {
		t.Error("EXPECTED: 2", "GOT:", a)
	}

	expectedResourceNames := []string{"mockbucket1", "mockbucket2", "mockbucket3"}
	expectedResourceAtRestEncryptions := []bool{true, true, false}
	expectedResourceTransportEncryptions := []bool{true, false, false}
	for i, r := range resources {
		log.Println("Testing name for resource (bucket)", i+1)
		if e, a := expectedResourceNames[i], r.GetName(); e != a {
			t.Error("EXPECTED:", e, "GOT:", a)
		}
		log.Println("Testing type of resource", i+1)
		if e := "ObjectStorage"; !r.HasType(e) {
			t.Error(e, "not found as type")
		}
		log.Println("Testing at rest encryption of resource", i+1)
		if e, a := expectedResourceAtRestEncryptions[i], r.AtRestEncryption.Enabled; e != a {
			log.Println(r.AtRestEncryption.Enabled)
			t.Error("EXPECTED", e, "GOT", a)
		}
		log.Println("Testing transport encryption of resource", i+1)
		if e, a := expectedResourceTransportEncryptions[i], r.HttpEndpoint.TransportEncryption.Enabled; e != a {
			t.Error("EXPECTED", e, "GOT", a)
		}
	}

}
