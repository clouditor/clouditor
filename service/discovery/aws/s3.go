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
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
)

// S3ListBucketsAPI is the interface for the List function (used for mock testing)
type S3ListBucketsAPI interface {
	// ListBuckets ToDo: What are and do the (general) context interface, the ListBucketsInput and s3.Options parameters
	// ListBuckets ToDo: Name has to be "ListBuckets" like in s3.client, right?
	ListBuckets(ctx context.Context,
		params *s3.ListBucketsInput,
		optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
}

func GetS3Client() *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	return s3.NewFromConfig(cfg)
}

func List(clientApi S3ListBucketsAPI) (resp *s3.ListBucketsOutput, err error) {
	input := &s3.ListBucketsInput{}
	// ToDo: Why doesnt "err" throw error due to the second declaration with the same name "err"
	// ToDo: One return line
	resp, err = clientApi.ListBuckets(context.TODO(), input)
	if err != nil {
		log.Fatalf("Error occured while retrieving buckets: %v", err)
	}
	log.Println("Retrieved buckets.")
	return

}

func GetObjectsOfBucket(myBucket string) *s3.ListObjectsV2Output {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Error occurred:", err)
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg)

	output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(myBucket),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("first page results:")
	for _, object := range output.Contents {
		log.Printf("key=%s size=%d", aws.ToString(object.Key), object.Size)
	}

	return output
}
