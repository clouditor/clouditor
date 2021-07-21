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
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/middleware"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	mockVM1            = "mockVM1"
	mockVM1ID          = "mockVM1ID"
	mockVM1Log         = "Mock Log for mockVM1ID"
	mockVM1Time        = 0
	blockVolumeId      = "blockVolumeID"
	networkInterfaceId = "networkInterfaceId"
)

type mockEC2API struct {
}

type mockEC2APIWithErrors struct {
}

func (m mockEC2API) DescribeInstances(_ context.Context, _ *ec2.DescribeInstancesInput, _ ...func(options *ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	// block device mappings for output struct
	blockDeviceMappings := []types.InstanceBlockDeviceMapping{
		{
			DeviceName: aws.String("/dev/sdh"),
			Ebs: &types.EbsInstanceBlockDevice{
				AttachTime:          nil,
				DeleteOnTermination: nil,
				Status:              "",
				VolumeId:            aws.String(blockVolumeId),
			},
		},
	}
	// tags for output struct
	tags := []types.Tag{
		{
			Key:   aws.String("name"),
			Value: aws.String(mockVM1ID),
		},
	}
	networkInterfaces := []types.InstanceNetworkInterface{
		{
			NetworkInterfaceId: aws.String(networkInterfaceId),
		},
	}

	// output struct containing all necessary information
	output := &ec2.DescribeInstancesOutput{
		NextToken: nil,
		Reservations: []types.Reservation{{
			Groups: nil,
			Instances: []types.Instance{{
				BlockDeviceMappings: blockDeviceMappings,
				InstanceId:          aws.String(mockVM1ID),
				NetworkInterfaces:   networkInterfaces,
				Tags:                tags,
			}},
			OwnerId:       nil,
			RequesterId:   nil,
			ReservationId: nil,
		}},
		ResultMetadata: middleware.Metadata{},
	}
	return output, nil
}

func (m mockEC2APIWithErrors) DescribeInstances(_ context.Context, _ *ec2.DescribeInstancesInput, _ ...func(options *ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	err := &smithy.GenericAPIError{
		Code:    "ConnectionError",
		Message: "Couldn't resolve host. Bad connection?",
	}
	return nil, err
}

func TestListCompute(t *testing.T) {
	d := computeDiscovery{
		client:        mockEC2API{},
		isDiscovering: true,
	}
	list, err := d.List()
	assert.Nil(t, err)
	assert.NotEmpty(t, list)
}

func TestDiscoverVirtualMachines(t *testing.T) {
	d := computeDiscovery{
		client: mockEC2API{},
	}
	machines, err := d.discoverVirtualMachines()
	assert.Nil(t, err)
	// Possible
	testMachine := machines[0]
	assert.Equal(t, mockVM1, testMachine.Name)
	// Possible
	assert.Equal(t, mockVM1ID, testMachine.ID)
	// Possible
	assert.NotEmpty(t, testMachine.BlockStorage)
	// TODO(lebogg): Possible to fetch? Via CloudWatch?
	//assert.Equal(t, mockVM1Log, machines.Log)
	// Possible
	assert.NotEmpty(t, testMachine.NetworkInterfaces)
	// Possible
	assert.Equal(t, mockVM1Time, testMachine.CreationTime)

	d = computeDiscovery{
		client: mockEC2APIWithErrors{},
	}
	_, err = d.discoverVirtualMachines()
	assert.NotNil(t, err)

}

// TODO(lebogg): Testing logs
func TestLogging(t *testing.T) {
	client, _ := NewClient()
	d := cloudwatchlogs.NewFromConfig(client.Cfg)
	input1 := &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String("testGroup"),
		LogStreamName: aws.String("testGroupTestStream"),
		EndTime:       nil,
		Limit:         nil,
		NextToken:     nil,
		StartFromHead: nil,
		StartTime:     nil,
	}
	resp, _ := d.GetLogEvents(context.TODO(), input1)
	fmt.Println(resp.Events)
	//input := &cloudwatchlogs.GetLogRecordInput{}
	//fmt.Println(d.GetLogRecord(context.TODO(), input))
}
