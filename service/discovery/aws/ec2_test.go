package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go/middleware"
	"testing"
)

type mockDescribeInstancesAPI func(ctx context.Context,
	params *ec2.DescribeInstancesInput,
	optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)

func (m mockDescribeInstancesAPI) DescribeInstances(ctx context.Context,
	params *ec2.DescribeInstancesInput,
	optFns ...func(options *ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	return m(ctx, params, optFns...)
}

func TestListEC2(t *testing.T) {
	mockInstanceId1 := "MockInstanceId1"
	client := func(t *testing.T) EC2DescribeInstancesAPI {
		return mockDescribeInstancesAPI(func(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
			t.Helper()
			return &ec2.DescribeInstancesOutput{
					NextToken: nil,
					Reservations: []types.Reservation{
						{
							Instances: []types.Instance{
								{
									InstanceId: aws.String(mockInstanceId1),
								},
								{
									InstanceId: aws.String("MockInstanceId2"),
								},
							},
						},
					},
					ResultMetadata: middleware.Metadata{},
				},
				nil
		})
	}
	reservations, err := ListInstances(client(t))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if i := *reservations.Reservations[0].Instances[0].InstanceId; i != mockInstanceId1 {
		t.Fatalf("First instanceId wrong. Expected: %v, got %v.", i, mockInstanceId1)
	}
	if len(reservations.Reservations[0].Instances) == 0 {
		t.Fatal("Amount of reservations is 0 but shouldn't be")
	}

}
