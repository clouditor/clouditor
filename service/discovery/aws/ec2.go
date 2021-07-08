// ToDo: Divide package in services? Such that list, e.g., does not have to be renamed?
package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"log"
)

// EC2DescribeInstancesAPI is the interface for ListInstances function (used for mock testing)
type EC2DescribeInstancesAPI interface {
	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

// ToDo: Currently lists reservations instead of instances only
func ListInstances(api EC2DescribeInstancesAPI) (resp *ec2.DescribeInstancesOutput, err error) {
	input := &ec2.DescribeInstancesInput{}
	resp, err = api.DescribeInstances(context.TODO(), input)
	if err != nil {
		log.Fatalf("Error occured while retrieving instances: %v", err)
	}
	log.Println("Retrieved instances.")
	return
}
