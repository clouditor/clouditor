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

func TestGetS3ServiceClient(t *testing.T) {
	if client := GetS3Client(); client == nil {
		t.Errorf("Connection failed. Credentials are nil.")
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
			buckets, err := List(tt.client(t))
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
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
			buckets, err := List(tt.client(t))
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
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
