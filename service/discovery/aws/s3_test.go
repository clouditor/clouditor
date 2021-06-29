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
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"strconv"
	"testing"
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
			Name: aws.String("AWS Mock Bucket 1"),
		},
		{
			Name: aws.String("AWS Mock Bucket 2"),
		},
	}}
	return output, nil
}

func (m mockS3APINew) GetBucketEncryption(ctx context.Context, params *s3.GetBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error) {
	// ToDo: Switch between different buckets (params) -> different SSEAlgorithm and KeyManager
	switch aws.ToString(params.Bucket) {
	case "Mock Bucket 1":
		output := &s3.GetBucketEncryptionOutput{
			ServerSideEncryptionConfiguration: &types.ServerSideEncryptionConfiguration{
				Rules: []types.ServerSideEncryptionRule{
					{
						ApplyServerSideEncryptionByDefault: &types.ServerSideEncryptionByDefault{
							SSEAlgorithm:   "AES256",
							KMSMasterKeyID: aws.String("SSE-S3"),
						},
						BucketKeyEnabled: false,
					},
				},
			},
		}
		return output, nil
	default:
		return nil, errors.New("there was an error")
	}
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

// TestGetBucketsNew tests the getBuckets method (with other form of mocking implementation)
// ToDo: Its simpler and shorter but I would like the other one more (with "cases")
func TestGetBucketsNew(t *testing.T) {
	// ToDo: It is not better to use a pointer (&awsS3Discovery), is it?
	d := awsS3Discovery{
		client:        mockS3APINew{},
		buckets:       nil,
		bucketNames:   nil,
		isDiscovering: false,
	}
	d.getBuckets()
	log.Println("Here are the buckets: ", d.bucketNames)
}

func TestGetEncryptionConf(t *testing.T) {
	d := awsS3Discovery{
		client:        mockS3APINew{},
		buckets:       nil,
		bucketNames:   nil,
		isDiscovering: false,
	}
	isEncrypted, algorithm, manager := d.getEncryptionConf("Mock Bucket 1")
	if isEncrypted == false {
		t.Error("Expected:", true, ".Got:", isEncrypted)
	}
	if e := "AES256"; algorithm != e {
		t.Error("Expected:", e, ".Got:", algorithm)
	}
	if e := "SSE-S3"; manager != e {
		t.Error("Expected:", e, ".Got:", manager)
	}
	isEncrypted, algorithm, manager = d.getEncryptionConf("Mock Bucket with no encryption")
	if isEncrypted == true {
		t.Error("Expected:", false, ".Got:", isEncrypted)
	}
}

// TestGetBuckets tests the getBuckets method
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
				bucketNames:   nil,
				isDiscovering: false,
			}
			d.getBuckets()
			log.Println("Here are the buckets: ", d.bucketNames)
		})
	}
}

// TestGetEncryptionConf tests the getEncryptionConf method
//func TestGetEncryptionConf(t *testing.T) {
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
	for i, bucket := range d.bucketNames {
		isEncrypted, _, _ := d.getEncryptionConf(bucket)
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
	for _, bucket := range d.bucketNames {
		if d.checkPublicAccessBlockConfiguration(bucket) == false {
			t.Fatalf("Expected no public access of bucket. But public access is enabled for %v.", bucket)
		}
	}
}

func TestCheckBucketReplication(t *testing.T) {
	d := NewAwsStorageDiscovery(NewAwsDiscovery().cfg)
	d.getBuckets()
	for _, bucket := range d.bucketNames {
		if d.checkBucketReplication(bucket) == true {
			t.Fatalf("Expected no replication setting for bucket. But replication is set for bucket '%v'.", bucket)
		}
	}
}

func TestCheckLifeCycleConfiguration(t *testing.T) {
	d := NewAwsStorageDiscovery(NewAwsDiscovery().cfg)
	d.getBuckets()
	for _, bucket := range d.bucketNames {
		if d.checkLifeCycleConfiguration(bucket) == true {
			t.Fatalf("Expected no life cycle configuration setting for bucket. But it is set for bucket '%v'.", bucket)
		}
	}
}

//func TestGetObjectsOfBucket_whenNotEmpty(t *testing.T) {
//	if bucketObjects := GetObjectsOfBucket(os.Getenv("TESTBUCKET")); len(bucketObjects.Contents) == 0 {
//		t.Errorf("No buckets found")
//	}
//}
