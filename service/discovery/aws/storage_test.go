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
	"testing"
	"time"
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
		oe := &smithy.OperationError{
			ServiceID:     "MockS3API",
			OperationName: "GetBucketEncryption",
			Err:           errors.New("failed to resolve service endpoint"),
		}
		return nil, oe
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
	case "mockbucket2": // JSON failure
		output = &s3.GetBucketPolicyOutput{
			Policy: aws.String(""),
		}
		err = nil
	case "mockbucket3": // Effect audit instead of deny
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
		OperationName: "GetBucketEncryption",
		Err:           errors.New("failed to resolve service endpoint"),
	}
	return nil, oe
}

func (m mockS3APIWitHErrors) GetBucketEncryption(_ context.Context, _ *s3.GetBucketEncryptionInput, _ ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error) {
	oe := &smithy.OperationError{
		ServiceID:     "MockS3API",
		OperationName: "GetBucketEncryption",
		Err:           errors.New("failed to resolve service endpoint"),
	}
	return nil, oe
}

func (m mockS3APIWitHErrors) GetBucketPolicy(_ context.Context, _ *s3.GetBucketPolicyInput, _ ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
	oe := &smithy.OperationError{
		ServiceID:     "MockS3API",
		OperationName: "GetBucketEncryption",
		Err:           errors.New("failed to resolve service endpoint"),
	}
	return nil, oe
}

func (m mockS3APIWitHErrors) GetBucketLocation(_ context.Context, _ *s3.GetBucketLocationInput, _ ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
	oe := &smithy.OperationError{
		ServiceID:     "MockS3API",
		OperationName: "GetBucketEncryption",
		Err:           errors.New("failed to resolve service endpoint"),
	}
	return nil, oe
}

func (m mockS3APIWitHErrors) GetPublicAccessBlock(_ context.Context, _ *s3.GetPublicAccessBlockInput, _ ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {
	oe := &smithy.OperationError{
		ServiceID:     "MockS3API",
		OperationName: "GetBucketEncryption",
		Err:           errors.New("failed to resolve service endpoint"),
	}
	return nil, oe
}

func (m mockS3APIWitHErrors) GetBucketReplication(_ context.Context, _ *s3.GetBucketReplicationInput, _ ...func(*s3.Options)) (*s3.GetBucketReplicationOutput, error) {
	oe := &smithy.OperationError{
		ServiceID:     "MockS3API",
		OperationName: "GetBucketEncryption",
		Err:           errors.New("failed to resolve service endpoint"),
	}
	return nil, oe
}

func (m mockS3APIWitHErrors) GetBucketLifecycleConfiguration(_ context.Context, _ *s3.GetBucketLifecycleConfigurationInput, _ ...func(*s3.Options)) (*s3.GetBucketLifecycleConfigurationOutput, error) {
	oe := &smithy.OperationError{
		ServiceID:     "MockS3API",
		OperationName: "GetBucketEncryption",
		Err:           errors.New("failed to resolve service endpoint"),
	}
	return nil, oe
}

// TestGetBuckets tests the getBuckets method (with other form of mocking implementation)
func TestGetBuckets(t *testing.T) {
	d := awsS3Discovery{
		client:        mockS3APINew{},
		buckets:       nil,
		isDiscovering: false,
	}
	err := d.getBuckets()
	if err != nil {
		t.Error("EXPECTED no errors. GOT one:", err)
	}
	log.Print("Testing number of buckets")
	if e, a := 2, len(d.buckets); e != a {
		t.Error("EXPECTED:", e, "GOT:", a)
	}

	log.Print("Testing name of first bucket")
	if e, a := mockBucket1, d.buckets[0].name; e != a {
		t.Error("EXPECTED:", e, "GOT:", a)
	}
	log.Print("Testing region of first bucket")
	if e, a := mockBucket1Region, d.buckets[0].region; e != a {
		t.Error("EXPECTED:", e, "GOT:", a)
	}
	log.Print("Testing endpoint of first bucket")
	if e, a := mockBucket1Endpoint, d.buckets[0].endpoint; e != a {
		t.Error("EXPECTED:", e, "GOT:", a)
	}
	log.Print("Testing creation time of first bucket")
	expectedCreationTime1, _ := time.Parse(time.RFC3339, mockBucket1CreationTime)
	if e, a := expectedCreationTime1, d.buckets[0].creationTime; e.String() != a.String() {
		t.Error("EXPECTED:", e, "GOT:", a)
	}

	log.Print("Testing name of second bucket")
	if e, a := mockBucket2, d.buckets[1].name; e != a {
		t.Error("EXPECTED:", e, "GOT:", a)
	}
	log.Print("Testing region of second bucket")
	if e, a := mockBucket2Region, d.buckets[1].region; e != a {
		t.Error("EXPECTED:", e, "GOT:", a)
	}
	log.Print("Testing endpoint of second bucket")
	if e, a := mockBucket2Endpoint, d.buckets[1].endpoint; e != a {
		t.Error("EXPECTED:", e, "GOT:", a)
	}
	log.Print("Testing creation time of second bucket")
	expectedCreationTime2, _ := time.Parse(time.RFC3339, mockBucket2CreationTime)
	if e, a := expectedCreationTime2, d.buckets[1].creationTime; e.String() != a.String() {
		t.Error("EXPECTED:", e, "GOT:", a)
	}

	fmt.Println(d.buckets[1].name)
	fmt.Println(d.buckets[1].creationTime)

	d = awsS3Discovery{
		client:        mockS3APIWitHErrors{},
		buckets:       nil,
		isDiscovering: false,
	}

	if err = d.getBuckets(); err == nil {
		t.Error("EXPECTED error. GOT none")
	}
}

// TestGetEncryptionAtRest tests the getEncryptionAtRest method
// ToDo: Check errors explicitly
func TestGetEncryptionAtRest(t *testing.T) {
	d := awsS3Discovery{
		client:        mockS3APINew{},
		buckets:       nil,
		isDiscovering: false,
	}

	// First case
	encryptionAtRest, _ := d.getEncryptionAtRest(mockBucket1)
	if isEncrypted := encryptionAtRest.Encryption.Enabled; !isEncrypted {
		t.Error("Expected:", true, ".Got:", isEncrypted)
	}
	if e, a := "AES256", encryptionAtRest.Algorithm; e != a {
		t.Error("Expected:", e, ".Got:", a)
	}
	if e, a := "SSE-S3", encryptionAtRest.KeyManager; e != a {
		t.Error("Expected:", e, ".Got:", a)
	}

	// Second case
	encryptionAtRest, _ = d.getEncryptionAtRest(mockBucket2)
	if isEncrypted := encryptionAtRest.Encryption.Enabled; !isEncrypted {
		t.Error("Expected:", true, ".Got:", isEncrypted)
	}
	if e, a := "aws:kms", encryptionAtRest.Algorithm; e != a {
		t.Error("Expected:", e, ".Got:", a)
	}
	if e, a := "SSE-KMS", encryptionAtRest.KeyManager; e != a {
		t.Error("Expected:", e, ".Got:", a)
	}

	// Third case (no encryption at all)
	encryptionAtRest, _ = d.getEncryptionAtRest("Mock Bucket with no encryption")
	if isEncrypted := encryptionAtRest.Encryption.Enabled; isEncrypted {
		t.Error("Expected:", false, ".Got:", isEncrypted)
	}
}

// TestGetTransportEncryption tests the getTransportEncryption method
func TestGetTransportEncryption(t *testing.T) {
	// Case 1 (error case):
	d := awsS3Discovery{
		client:        mockS3APIWitHErrors{},
		buckets:       nil,
		isDiscovering: false,
	}
	encryptionAtTransit, err := d.getTransportEncryption("")
	if err == nil || encryptionAtTransit.Encryption.Enabled || encryptionAtTransit.TlsVersion != "" || encryptionAtTransit.Enforced {
		t.Errorf("Expected error but got no one. encryptionEnabled: %v, TLS version: %v, isEnforced: %v."+
			"Got: %v, %v, %v", false, "", false, encryptionAtTransit.Encryption.Enabled, encryptionAtTransit.TlsVersion, encryptionAtTransit.Enforced)
	}

	d = awsS3Discovery{
		client:        mockS3APINew{},
		buckets:       nil,
		isDiscovering: false,
	}

	// Case 2
	encryptionAtTransit, _ = d.getTransportEncryption(mockBucket1)
	if !encryptionAtTransit.Encryption.Enabled || encryptionAtTransit.TlsVersion != "TLS1.2" || !encryptionAtTransit.Enforced {
		t.Errorf("Expected encryptionEnabled: %v, TLS version: %v, isEnforced: %v."+
			"Got: %v, %v, %v", true, "TLS1.2", true, encryptionAtTransit.Encryption.Enabled, encryptionAtTransit.TlsVersion, encryptionAtTransit.Enforced)
	}

	// Case 3: JSON failure
	encryptionAtTransit, _ = d.getTransportEncryption(mockBucket2)
	if encryptionAtTransit.Encryption.Enabled || encryptionAtTransit.TlsVersion != "" || encryptionAtTransit.Enforced {
		t.Errorf("Expected isEncrypted: %v, TLS version: %v, enforced: %v. "+
			"Got: %v, %v, %v.", false, "", false, encryptionAtTransit.Encryption.Enabled, encryptionAtTransit.TlsVersion, encryptionAtTransit.Enforced)
	}

	// Case 4
	encryptionAtTransit, _ = d.getTransportEncryption("mockbucket3")
	if encryptionAtTransit.Encryption.Enabled || encryptionAtTransit.TlsVersion != "" || encryptionAtTransit.Enforced {
		t.Errorf("Expected isEncrypted: %v, TLS version: %v, enforced: %v."+
			"Got: %v, %v, %v.", false, "", false, encryptionAtTransit.Encryption.Enabled, encryptionAtTransit.TlsVersion, encryptionAtTransit.Enforced)
	}

}

// TestGetRegion tests the getRegion method
// ToDo: Test errors explicitly
func TestGetRegion(t *testing.T) {
	d := awsS3Discovery{
		client:        mockS3APINew{},
		buckets:       nil,
		isDiscovering: false,
	}
	if a, _ := d.getRegion(mockBucket1); mockBucket1Region != a {
		t.Error("Expected: ", mockBucket1Region, ". Got:", a)
	}
	if a, _ := d.getRegion(mockBucket2); mockBucket2Region != a {
		t.Error("Expected: ", mockBucket2Region, ". Got:", a)
	}

	// Error case
	if _, err := d.getRegion("mockbucketNotAvailable"); err == nil {
		t.Error("EXPECTED error. GOT none")
	}

}

// TestName tests the Name method
func TestName(t *testing.T) {
	d := awsS3Discovery{
		client:        mockS3APINew{},
		buckets:       nil,
		isDiscovering: false,
	}
	log.Println("Testing the name of the AWS Blob Storage Discovery")
	if e, a := "AWS Blob Storage", d.Name(); e != a {
		t.Error("EXPECTED:", e, "GOT", a)
	}
}

// TestList tests the List method
// ToDo: Currently no full coverage since testing all errors in single tests yields early termination here
// ToDo: One possibility: Create new MockAPIs for each case...
func TestList(t *testing.T) {
	d := awsS3Discovery{
		client:        mockS3APINew{},
		buckets:       nil,
		isDiscovering: false,
	}
	resources, err := d.List()
	if err == nil {
		t.Error("EXPECTED error because MockBucket2 should throw JSON error. But GOT no error")
	}
	log.Println("Testing number of resources (buckets)")
	if e, a := 1, len(resources); e != a {
		t.Error("EXPECTED: 2", "GOT:", a)
	}

	expectedResourceNames := []string{mockBucket1, "mockbucket2", "mockbucket3"}
	//expectedResourceAtRestEncryptions := []bool{true, true, false}
	//expectedResourceTransportEncryptions := []bool{true, false, false}
	for i, r := range resources {
		log.Println("Testing name for resource (bucket)", i+1)
		if e, a := expectedResourceNames[i], r.GetName(); e != a {
			t.Error("EXPECTED:", e, "GOT:", a)
		}
		log.Println("Testing type of resource", i+1)
		if e := "ObjectStorage"; !r.HasType(e) {
			t.Error(e, "not found as type")
		}
		// ToDo: How to convert to ObjectStorageResource s.t. we can access atRestEncryption?
		//r = voc.ObjectStorageResource(r)
		//log.Println("Testing at rest encryption of resource", i+1)
		//if e, a := expectedResourceAtRestEncryptions[i], r.AtRestEncryption.Enabled; e != a {
		//	log.Println(r.AtRestEncryption.Enabled)
		//	t.Error("EXPECTED", e, "GOT", a)
		//}
		//log.Println("Testing transport encryption of resource", i+1)
		//if e, a := expectedResourceTransportEncryptions[i], r.HttpEndpoint.TransportEncryption.Enabled; e != a {
		//	t.Error("EXPECTED", e, "GOT", a)
		//}
	}

}

// ToDO: Old API implementation. Delete in the future
//// Mock s3 service and methods
//type mockS3API func(ctx context.Context,
//	params *s3.ListBucketsInput,
//	optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
//
//func (m mockS3API) GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
//	panic("implement me")
//}
//
//func (m mockS3API) GetBucketPolicy(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
//	panic("implement me")
//}
//
//func (m mockS3API) GetBucketEncryption(ctx context.Context, params *s3.GetBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error) {
//	panic("implement me")
//}
//
//func (m mockS3API) GetPublicAccessBlock(ctx context.Context, params *s3.GetPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {
//	panic("implement me")
//}
//
//func (m mockS3API) GetBucketReplication(ctx context.Context, params *s3.GetBucketReplicationInput, optFns ...func(*s3.Options)) (*s3.GetBucketReplicationOutput, error) {
//	panic("implement me")
//}
//
//func (m mockS3API) GetBucketLifecycleConfiguration(ctx context.Context, params *s3.GetBucketLifecycleConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLifecycleConfigurationOutput, error) {
//	panic("implement me")
//}
//
//func (m mockS3API) ListBuckets(ctx context.Context,
//	params *s3.ListBucketsInput,
//	optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
//
//	// ToDo what does "m(xxx)" mean? What does it return from mockS3API?
//	return m(ctx, params, optFns...)
//}

// ToDo: Potential future implementation or implementation via credentials

//var Cfg aws.Config
//
//func init() {
//	Cfg = NewClient().Cfg
//}

//// TestGetPublicAccessBlockConfiguration tests the getPublicAccessBlockConfiguration method
//func TestGetPublicAccessBlockConfiguration(t *testing.T) {
//	d := awsS3Discovery{
//		client:        mockS3APINew{},
//		buckets:       nil,
//		isDiscovering: false,
//	}
//	log.Print(d)
//	//d.getPublicAccessBlockConfiguration("Mock Bucket 1")
//}

//func TestCheckPublicAccessBlockConfiguration(t *testing.T) {
//	d := NewAwsStorageDiscovery(NewClient().Cfg)
//	d.getBuckets()
//	for _, bucket := range d.buckets {
//		if d.getPublicAccessBlockConfiguration(bucket.name) == false {
//			t.Fatalf("Expected no public access of bucket. But public access is enabled for %v.", bucket)
//		}
//	}
//}
//
//func TestCheckBucketReplication(t *testing.T) {
//	d := NewAwsStorageDiscovery(NewClient().Cfg)
//	d.getBuckets()
//	for _, bucket := range d.buckets {
//		if d.checkBucketReplication(bucket.name) == true {
//			t.Fatalf("Expected no replication setting for bucket. But replication is set for bucket '%v'.", bucket)
//		}
//	}
//}
//
//func TestCheckLifeCycleConfiguration(t *testing.T) {
//	d := NewAwsStorageDiscovery(NewClient().Cfg)
//	d.getBuckets()
//	for _, bucket := range d.buckets {
//		if d.checkLifeCycleConfiguration(bucket.name) == true {
//			t.Fatalf("Expected no life cycle configuration setting for bucket. But it is set for bucket '%v'.", bucket)
//		}
//	}
//}

//func TestGetObjectsOfBucket_whenNotEmpty(t *testing.T) {
//	if bucketObjects := GetObjectsOfBucket(os.Getenv("TESTBUCKET")); len(bucketObjects.Contents) == 0 {
//		t.Errorf("No buckets found")
//	}
//}
