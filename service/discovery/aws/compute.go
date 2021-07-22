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
	"clouditor.io/clouditor/voc"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go"
)

// computeDiscovery handles the AWS API requests regarding the EC2 service
type computeDiscovery struct {
	client        EC2API
	isDiscovering bool
}

type EC2API interface {
	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(options *ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

func (d computeDiscovery) Name() string {
	return "AWS Compute"
}

func (d computeDiscovery) List() ([]voc.IsResource, error) {
	panic("implement me")
}

func (d *computeDiscovery) discoverVirtualMachines() ([]voc.VirtualMachineResource, error) {
	resp, err := d.client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			err = fmt.Errorf("code: %v, fault: %v, message: %v", ae.ErrorCode(), ae.ErrorFault(), ae.ErrorMessage())
		}
		return nil, err
	}
	var resources []voc.VirtualMachineResource
	for _, reservation := range resp.Reservations {
		for _, vm := range reservation.Instances {
			computeResource := voc.ComputeResource{
				Resource: voc.Resource{
					// TODO(all): What do we accept as the ID: ARN or EC2-Instance-ID?
					ID:   voc.ResourceID(aws.ToString(vm.InstanceId)),
					Name: getNameOfVM(vm),
					// TODO(all): Is launchTime sufficient?
					CreationTime: vm.LaunchTime.Unix(),
					Type:         []string{"VirtualMachine", "Compute", "Resource"},
				},
				NetworkInterfaces: getNetworkInterfacesOfVM(vm),
			}

			resources = append(resources, voc.VirtualMachineResource{
				ComputeResource: computeResource,
				BlockStorage:    getBlockStorageIDsOfVM(vm),
				// TODO(lebogg): How to derive logs
				Log: getLogsOfVM(vm),
			})
		}
	}
	return resources, nil
}

// TODO(lebogg) Find out how to get logs
// getLogsOfVM checks if logging is enabled
func getLogsOfVM(_ types.Instance) (l *voc.Log) {
	l = new(voc.Log)
	return
}

func getBlockStorageIDsOfVM(vm types.Instance) (blockStorageIDs []voc.ResourceID) {
	for _, mapping := range vm.BlockDeviceMappings {
		blockStorageIDs = append(blockStorageIDs, voc.ResourceID(aws.ToString(mapping.Ebs.VolumeId)))
	}
	return
}

// getNetworkInterfacesOfVM returns the network interface IDs by iterating through the VMs network interfaces
func getNetworkInterfacesOfVM(vm types.Instance) (networkInterfaceIDs []voc.ResourceID) {
	for _, networkInterface := range vm.NetworkInterfaces {
		networkInterfaceIDs = append(networkInterfaceIDs, voc.ResourceID(aws.ToString(networkInterface.NetworkInterfaceId)))
	}
	return
}

// getNameOfVM returns the name if exists (= a tag with key 'name' exists), otherwise empty string
func getNameOfVM(vm types.Instance) (name string) {
	for _, tag := range vm.Tags {
		if aws.ToString(tag.Key) == "name" {
			name = aws.ToString(tag.Value)
		}
	}
	return
}
