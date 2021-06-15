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
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"os"
	"strconv"
	"testing"
)

// ToDo: Needs re-writing

func TestGetS3ServiceClient(t *testing.T) {
	cfg := NewAwsDiscovery().cfg
	client := GetS3Client(cfg)
	if client == nil {
		t.Errorf("Connection failed. Credentials are nil.")
	}
	_, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
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

func TestListS3(t *testing.T) {
	mockDisplayName := "MockDisplayName"
	mockId := "MockId1234"
	cases := []struct {
		client func(t *testing.T) S3ListBucketsAPI
		//expect []byte
	}{
		{

			client: func(t *testing.T) S3ListBucketsAPI {
				return mockListBucketsAPI(func(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
					t.Helper()
					return &s3.ListBucketsOutput{
							Buckets: []types.Bucket{
								types.Bucket{
									CreationDate: nil,
									Name:         aws.String("FirstMockBucket"),
								},
								types.Bucket{
									CreationDate: nil,
									Name:         aws.String("SecondMockBucket"),
								},
							},
							Owner: &types.Owner{
								DisplayName: &mockDisplayName,
								ID:          &mockId,
							},
						},
						nil
				})
			},
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			buckets := List(tt.client(t))
			if len(buckets.Buckets) == 0 {
				t.Fatal("Buckets empty but shouldn't be.")
			}
			if o := *buckets.Owner.ID; o != mockId {
				t.Fatalf("expected %v, but got %v", mockId, o)
			}
		})
	}
}

func TestAreBucketsEncrypted(t *testing.T) {
	mockDisplayName := "MockDisplayName"
	mockId := "MockId1234"
	cases := []struct {
		client func(t *testing.T) S3ListBucketsAPI
		//expect []byte
	}{
		{

			client: func(t *testing.T) S3ListBucketsAPI {
				return mockListBucketsAPI(func(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
					t.Helper()
					return &s3.ListBucketsOutput{
							Buckets: []types.Bucket{
								types.Bucket{
									CreationDate: nil,
									Name:         aws.String("FirstMockBucket"),
								},
								types.Bucket{
									CreationDate: nil,
									Name:         aws.String("SecondMockBucket"),
								},
							},
							Owner: &types.Owner{
								DisplayName: &mockDisplayName,
								ID:          &mockId,
							},
						},
						nil
				})
			},
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			buckets := List(tt.client(t))
			if len(buckets.Buckets) == 0 {
				t.Fatal("Buckets empty but shouldn't be.")
			}
			if o := *buckets.Owner.ID; o != mockId {
				t.Fatalf("expected %v, but got %v", mockId, o)
			}
		})
	}
}

// ToDo: End mock stuff

func TestGetObjectsOfBucket_whenNotEmpty(t *testing.T) {
	if bucketObjects := GetObjectsOfBucket(os.Getenv("TESTBUCKET")); len(bucketObjects.Contents) == 0 {
		t.Errorf("No buckets found")
	}
}
