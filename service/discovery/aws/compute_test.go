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
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go/middleware"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockEC2 struct {
}

func (m mockEC2) DescribeInstances(_ context.Context, _ *ec2.DescribeInstancesInput, _ ...func(options *ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	output := &ec2.DescribeInstancesOutput{
		NextToken: nil,
		Reservations: []types.Reservation{{
			Groups: nil,
			Instances: []types.Instance{{
				AmiLaunchIndex:                          nil,
				Architecture:                            "",
				BlockDeviceMappings:                     nil,
				BootMode:                                "",
				CapacityReservationId:                   nil,
				CapacityReservationSpecification:        nil,
				ClientToken:                             nil,
				CpuOptions:                              nil,
				EbsOptimized:                            nil,
				ElasticGpuAssociations:                  nil,
				ElasticInferenceAcceleratorAssociations: nil,
				EnaSupport:                              nil,
				EnclaveOptions:                          nil,
				HibernationOptions:                      nil,
				Hypervisor:                              "",
				IamInstanceProfile:                      nil,
				ImageId:                                 nil,
				InstanceId:                              nil,
				InstanceLifecycle:                       "",
				InstanceType:                            "",
				KernelId:                                nil,
				KeyName:                                 nil,
				LaunchTime:                              nil,
				Licenses:                                nil,
				MetadataOptions:                         nil,
				Monitoring:                              nil,
				NetworkInterfaces:                       nil,
				OutpostArn:                              nil,
				Placement:                               nil,
				Platform:                                "",
				PrivateDnsName:                          nil,
				PrivateIpAddress:                        nil,
				ProductCodes:                            nil,
				PublicDnsName:                           nil,
				PublicIpAddress:                         nil,
				RamdiskId:                               nil,
				RootDeviceName:                          nil,
				RootDeviceType:                          "",
				SecurityGroups:                          nil,
				SourceDestCheck:                         nil,
				SpotInstanceRequestId:                   nil,
				SriovNetSupport:                         nil,
				State:                                   nil,
				StateReason:                             nil,
				StateTransitionReason:                   nil,
				SubnetId:                                nil,
				Tags:                                    nil,
				VirtualizationType:                      "",
				VpcId:                                   nil,
			}},
			OwnerId:       nil,
			RequesterId:   nil,
			ReservationId: nil,
		}},
		ResultMetadata: middleware.Metadata{},
	}
	return output, nil

}

func TestListCompute(t *testing.T) {
	d := computeDiscovery{
		client:        mockEC2{},
		isDiscovering: true,
	}
	list, err := d.List()
	assert.NotNil(t, err)
	assert.NotEmpty(t, list)
}

const (
	mockVM1     = "mockVM1"
	mockVM1ID   = "mockVM1ID"
	mockVM1Log  = "Mock Log for mockVM1ID"
	mockVM1Time = 0
)

func TestGetVMs(t *testing.T) {
	d := computeDiscovery{
		client:        mockEC2{},
		isDiscovering: true,
	}
	machines, err := d.getVMs()
	assert.NotNil(t, err)
	// Possible
	assert.Equal(t, mockVM1, machines.Name)
	// Possible
	assert.Equal(t, mockVM1ID, machines.ID)
	// Possible
	assert.NotEmpty(t, machines.BlockStorage)
	// TODO(lebogg): Possible to fetch? Via CloudWatch?
	assert.Equal(t, mockVM1Log, machines.Log)
	// Possible
	assert.NotEmpty(t, machines.NetworkInterfaces)
	// Possible
	assert.Equal(t, mockVM1Time, machines.CreationTime)

}
