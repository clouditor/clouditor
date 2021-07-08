package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
)

func GetS3ServiceClient() *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Error occurred:", err)
		log.Fatal(err)
	}

	return s3.NewFromConfig(cfg)
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
