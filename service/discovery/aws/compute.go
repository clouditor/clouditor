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
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	types2 "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"time"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go"
)

// computeDiscovery handles the AWS API requests regarding the EC2 service
type computeDiscovery struct {
	virtualMachineAPI EC2API
	functionAPI       LambdaAPI
	isDiscovering     bool
	awsConfig         *Client
}

// EC2API describes the EC2 api interface (for mock testing)
type EC2API interface {
	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(options *ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

// LambdaAPI describes the lambda api interface (for mock testing)
// TODO(lebogg): Is there a way to squash both, EC2 and lambda, into one interface?
type LambdaAPI interface {
	ListFunctions(ctx context.Context,
		params *lambda.ListFunctionsInput, optFns ...func(*lambda.Options)) (*lambda.ListFunctionsOutput, error)
}

// newFromConfigEC2 holds ec2.NewFromConfig(...) allowing a test function to mock it
var newFromConfigEC2 = ec2.NewFromConfig

// newFromConfigLambda holds lambda.NewFromConfig(...) allowing a test function tp mock it
var newFromConfigLambda = lambda.NewFromConfig

// NewComputeDiscovery constructs a new awsS3Discovery initializing the s3-virtualMachineAPI and isDiscovering with true
func NewComputeDiscovery(client *Client) discovery.Discoverer {
	return &computeDiscovery{
		virtualMachineAPI: newFromConfigEC2(client.Cfg),
		functionAPI:       newFromConfigLambda(client.Cfg),
		isDiscovering:     true,
		awsConfig:         client,
	}
}

func (d computeDiscovery) Name() string {
	return "AWS Compute"
}

func (d computeDiscovery) List() (resources []voc.IsResource, err error) {
	listOfVMs, err := d.discoverVirtualMachines()
	if err != nil {
		return
	}
	for _, machine := range listOfVMs {
		resources = append(resources, &machine)
	}

	listOfFunctions, err := d.discoverFunctions()
	if err != nil {
		return
	}
	for _, function := range listOfFunctions {
		resources = append(resources, &function)
	}

	return
}

// TODO(all): Do we want to cover all VMs or only VMs in current region?
func (d *computeDiscovery) discoverVirtualMachines() ([]voc.VirtualMachineResource, error) {
	resp, err := d.virtualMachineAPI.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			err = formatError(ae)
		}
		return nil, err
	}
	var resources []voc.VirtualMachineResource
	for _, reservation := range resp.Reservations {
		for _, vm := range reservation.Instances {
			computeResource := voc.ComputeResource{
				Resource: voc.Resource{
					ID:   d.getARN(vm),
					Name: d.getNameOfVM(vm),
					// TODO(all): Currently only the launch time can be derived directly. We could derive the creation
					// time of the attached volume. But this 1st) requires an additional API Call. It is 2nd) a rather
					// ugly solution since although it is likely to be not detached, it is NOT guaranteed.
					CreationTime: vm.LaunchTime.Unix(),
					Type:         []string{"VirtualMachine", "Compute", "Resource"},
				},
				NetworkInterfaces: d.getNetworkInterfacesOfVM(vm),
			}

			resources = append(resources, voc.VirtualMachineResource{
				ComputeResource: computeResource,
				BlockStorage:    d.getBlockStorageIDsOfVM(vm),
				// TODO(lebogg): How to derive logs
				Log: getLogsOfVM(vm),
			})
		}
	}
	return resources, nil
}

// TODO(all): lastModified for creation Time?
// TODO(all): lambda can have "elastic network interfaces" if it is connected to a VPC. But you only get IDs of SecGroup, Subnet and VPC
// TODO(lebogg): FunctionVersion in input to ALL?
// TODO(lebogg): "Lambda returns up to 50 functions per call" -> Whats when there are >50? I think "NextMarker" (string)
func (d *computeDiscovery) discoverFunctions() ([]voc.FunctionResource, error) {
	input := &lambda.ListFunctionsInput{
		FunctionVersion: types2.FunctionVersionAll,
	}

	resp, err := d.functionAPI.ListFunctions(context.TODO(), input)

	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			err = formatError(ae)
		}
		return nil, err
	}

	var resources []voc.FunctionResource
	for _, function := range resp.Functions {
		lastModified, err := parseTime(function.LastModified)
		if err != nil {
			return resources, err
		}
		resources = append(resources, voc.FunctionResource{
			ComputeResource: voc.ComputeResource{
				Resource: voc.Resource{
					ID:           voc.ResourceID(aws.ToString(function.FunctionArn)),
					Name:         aws.ToString(function.FunctionName),
					CreationTime: lastModified,
					Type:         []string{"Function", "Compute", "Resource"},
				},
				// TODO(lebogg): Can I retrieve network interface IDs? (VPC ID possible)
				NetworkInterfaces: nil,
			},
		})
	}
	return resources, nil
}

func parseTime(t *string) (int64, error) {
	parsedT, err := time.Parse(time.RFC3339, *t)
	if err != nil {
		return 0, err
	}
	return parsedT.Unix(), nil
}

// TODO(lebogg): Try in other discoverer (e.g. storage) and maybe put it in aws.go
func formatError(ae smithy.APIError) error {
	return fmt.Errorf("code: %v, fault: %v, message: %v", ae.ErrorCode(), ae.ErrorFault(), ae.ErrorMessage())
}

// TODO(all): Currently there is no option to find out if logs are enabled -> Default value false?
// getLogsOfVM checks if logging is enabled
func getLogsOfVM(_ types.Instance) (l *voc.Log) {
	l = new(voc.Log)
	l.Enabled = false
	return
}

func (d *computeDiscovery) getBlockStorageIDsOfVM(vm types.Instance) (blockStorageIDs []voc.ResourceID) {
	for _, mapping := range vm.BlockDeviceMappings {
		blockStorageIDs = append(blockStorageIDs, voc.ResourceID(aws.ToString(mapping.Ebs.VolumeId)))
	}
	return
}

// getNetworkInterfacesOfVM returns the network interface IDs by iterating through the VMs network interfaces
func (d *computeDiscovery) getNetworkInterfacesOfVM(vm types.Instance) (networkInterfaceIDs []voc.ResourceID) {
	for _, networkInterface := range vm.NetworkInterfaces {
		networkInterfaceIDs = append(networkInterfaceIDs, voc.ResourceID(aws.ToString(networkInterface.NetworkInterfaceId)))
	}
	return
}

// getNameOfVM returns the name if exists (= a tag with key 'name' exists), otherwise instance ID is used
func (d *computeDiscovery) getNameOfVM(vm types.Instance) string {
	for _, tag := range vm.Tags {
		if aws.ToString(tag.Key) == "Name" {
			return aws.ToString(tag.Value)
		}
	}
	// If no tag with 'name' was found, return instanceId instead
	return aws.ToString(vm.InstanceId)
}

func (d computeDiscovery) getARN(vm types.Instance) voc.ResourceID {
	// TODO(lebogg): Get Account ID
	return voc.ResourceID("arn:aws:ec2:" +
		d.awsConfig.Cfg.Region + ":" +
		aws.ToString(d.awsConfig.accountID) +
		":instance/" +
		aws.ToString(vm.InstanceId))
}
