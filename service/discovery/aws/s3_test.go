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
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"testing"
)

// ToDo: Needs re-writing

var cfg aws.Config

func init() {
	cfg = NewAwsDiscovery().cfg
}

func TestGetS3ServiceClient(t *testing.T) {
	// ToDo
	client := NewS3Discovery(cfg)
	if client == nil {
		t.Errorf("Connection failed. Credentials are nil.")
	}
	//_, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	//if err != nil {
	//	t.Fatalf("Error: %v", err)
	//}
}

// ToDo: Begin mock stuff

type mockListBucketsAPI func(ctx context.Context,
	params *s3.ListBucketsInput,
	optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)

func (m mockListBucketsAPI) ListBuckets(ctx context.Context,
	params *s3.ListBucketsInput,
	optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {

	return m(ctx, params, optFns...)
}

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
//				return mockListBucketsAPI(func(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
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
//				return mockListBucketsAPI(func(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
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

// ToDo: Works with my credentials -> Mock it
func TestCheckEncryption(t *testing.T) {
	d := NewS3Discovery(cfg)
	d.getBuckets(d.client)
	for i, bucket := range d.bucketNames {
		isEncrypted := d.checkEncryption(bucket)
		if i == 0 && isEncrypted {
			t.Errorf("Expected that bucket %v is not encrypted, but it is", bucket)
		} else if i == 1 && !isEncrypted {
			t.Errorf("Expected that bucket %v is encrypted, but it is not", bucket)
		}
	}
}

func TestCheckPublicAccessBlockConfiguration(t *testing.T) {
	d := NewS3Discovery(NewAwsDiscovery().cfg)
	d.getBuckets(d.client)
	for _, bucket := range d.bucketNames {
		if d.checkPublicAccessBlockConfiguration(bucket) == false {
			t.Fatalf("Expected no public access of bucket. But public access is enabled for %v.", bucket)
		}
	}
}

func TestCheckBucketReplication(t *testing.T) {
	d := NewS3Discovery(NewAwsDiscovery().cfg)
	d.getBuckets(d.client)
	for _, bucket := range d.bucketNames {
		if d.checkBucketReplication(bucket) == true {
			t.Fatalf("Expected no replication setting for bucket. But replication is set for bucket '%v'.", bucket)
		}
	}
}

func TestCheckLifeCycleConfiguration(t *testing.T) {
	d := NewS3Discovery(NewAwsDiscovery().cfg)
	d.getBuckets(d.client)
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
