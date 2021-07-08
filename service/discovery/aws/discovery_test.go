package aws

import (
	"os"
	"testing"
)

// ToDo: Mock an AWS (S3) service client
func TestGetS3ServiceClient(t *testing.T) {
	if client := GetS3ServiceClient(); client == nil {
		t.Errorf("Connection failed. Credentials are nil.")
	}
}

func TestGetObjectsOfBucket_whenNotEmpty(t *testing.T) {
	if bucketObjects := GetObjectsOfBucket(os.Getenv("TESTBUCKET")); len(bucketObjects.Contents) == 0 {
		t.Errorf("No buckets found")
	}
}
