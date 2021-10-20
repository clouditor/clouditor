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
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	typesEC2 "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	typesLambda "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/smithy-go"
)

// computeDiscovery handles the AWS API requests regarding the computing services (EC2 and Lambda)
type computeDiscovery struct {
	virtualMachineAPI EC2API
	functionAPI       LambdaAPI
	isDiscovering     bool
	awsConfig         *Client
}

// EC2API describes the EC2 api interface which is implemented by the official AWS client and mock clients in tests
type EC2API interface {
	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(options *ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

// LambdaAPI describes the lambda api interface which is implemented by the official AWS client and mock clients in tests
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
		virtualMachineAPI: newFromConfigEC2(client.cfg),
		functionAPI:       newFromConfigLambda(client.cfg),
		isDiscovering:     true,
		awsConfig:         client,
	}
}

// Name is the method implementation defined in the discovery.Discoverer interface
func (d computeDiscovery) Name() string {
	return "AWS Compute"
}

// List is the method implementation defined in the discovery.Discoverer interface
func (d computeDiscovery) List() (resources []voc.IsCloudResource, err error) {
	listOfVMs, err := d.discoverVirtualMachines()
	if err != nil {
		return nil, fmt.Errorf("could not discover virtual machines: %w", err)
	}
	for _, machine := range listOfVMs {
		resources = append(resources, &machine)
	}

	listOfFunctions, err := d.discoverFunctions()
	if err != nil {
		return nil, fmt.Errorf("could not discover functions: %w", err)
	}
	for _, function := range listOfFunctions {
		resources = append(resources, &function)
	}

	return
}

// discoverVirtualMachines discovers all VMs (in the current region)
func (d *computeDiscovery) discoverVirtualMachines() ([]voc.VirtualMachine, error) {
	resp, err := d.virtualMachineAPI.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			err = formatError(ae)
		}
		return nil, err
	}
	var resources []voc.VirtualMachine
	for _, reservation := range resp.Reservations {
		for i := range reservation.Instances {
			vm := &reservation.Instances[i]
			computeResource := &voc.Compute{
				CloudResource: &voc.CloudResource{
					ID:           d.addARNToVM(vm),
					Name:         d.getNameOfVM(vm),
					CreationTime: 0,
					Type:         []string{"VirtualMachine", "Compute", "Resource"},
					GeoLocation: voc.GeoLocation{
						Region: d.awsConfig.cfg.Region,
					},
				},
			}

			resources = append(resources, voc.VirtualMachine{
				Compute:          computeResource,
				BlockStorage:     d.mapBlockStorageIDsOfVM(vm),
				NetworkInterface: d.getNetworkInterfacesOfVM(vm),
				BootLog:          getBootLog(vm),
				OSLog:            getOSLog(vm),
			})
		}
	}
	return resources, nil
}

// discoverFunctions discovers all lambda functions
func (d *computeDiscovery) discoverFunctions() (resources []voc.Function, err error) {
	// 'listFunctions' discovers up to 50 Lambda functions per execution -> loop through when response has nextMarker set
	var resp *lambda.ListFunctionsOutput
	var nextMarker *string
	for {
		resp, err = d.functionAPI.ListFunctions(context.TODO(), &lambda.ListFunctionsInput{
			Marker: nextMarker,
		})
		if err != nil {
			var ae smithy.APIError
			if errors.As(err, &ae) {
				err = formatError(ae)
			}
			return nil, err
		}
		resources = append(resources, d.mapFunctionResources(resp.Functions)...)

		if nextMarker = resp.NextMarker; nextMarker == nil {
			break
		}
	}

	return
}

// mapFunctionResources iterates functionConfigurations and returns a list of corresponding FunctionResources
func (d *computeDiscovery) mapFunctionResources(functions []typesLambda.FunctionConfiguration) (resources []voc.Function) {
	for i := range functions {
		function := &functions[i]
		resources = append(resources, voc.Function{
			Compute: &voc.Compute{
				CloudResource: &voc.CloudResource{
					ID:           voc.ResourceID(aws.ToString(function.FunctionArn)),
					Name:         aws.ToString(function.FunctionName),
					CreationTime: 0,
					Type:         []string{"Function", "Compute", "Resource"},
					GeoLocation: voc.GeoLocation{
						Region: d.awsConfig.cfg.Region,
					},
				},
			}})
	}
	return
}

// getBootLog checks if boot logging is enabled
// Currently there is no option to find out if any logs are enabled -> Assign default zero values
func getBootLog(_ *typesEC2.Instance) (l *voc.BootLog) {
	l = &voc.BootLog{
		Log: &voc.Log{
			Auditing:        nil,
			Output:          nil,
			Enabled:         false,
			RetentionPeriod: 0,
		},
	}
	return
}

// getOSLog checks if OS logging is enabled
// Currently there is no option to find out if any logs are enabled -> Assign default zero values
func getOSLog(_ *typesEC2.Instance) (l *voc.OSLog) {
	l = &voc.OSLog{
		Log: &voc.Log{
			Auditing:        nil,
			Output:          nil,
			Enabled:         false,
			RetentionPeriod: 0,
		},
	}
	return
}

// mapBlockStorageIDsOfVM returns block storages IDs by iterating the VMs block storages
func (d *computeDiscovery) mapBlockStorageIDsOfVM(vm *typesEC2.Instance) (blockStorageIDs []voc.ResourceID) {
	for _, mapping := range vm.BlockDeviceMappings {
		blockStorageIDs = append(blockStorageIDs, voc.ResourceID(aws.ToString(mapping.Ebs.VolumeId)))
	}
	return
}

// getNetworkInterfacesOfVM returns the network interface IDs by iterating the VMs network interfaces
func (d *computeDiscovery) getNetworkInterfacesOfVM(vm *typesEC2.Instance) (networkInterfaceIDs []voc.ResourceID) {
	for _, networkInterface := range vm.NetworkInterfaces {
		networkInterfaceIDs = append(networkInterfaceIDs, voc.ResourceID(aws.ToString(networkInterface.NetworkInterfaceId)))
	}
	return
}

// getNameOfVM returns the name if exists (i.e. a tag with key 'name' exists), otherwise instance ID is used
func (d *computeDiscovery) getNameOfVM(vm *typesEC2.Instance) string {
	for _, tag := range vm.Tags {
		if aws.ToString(tag.Key) == "Name" {
			return aws.ToString(tag.Value)
		}
	}
	// If no tag with 'name' was found, return instanceId instead
	return aws.ToString(vm.InstanceId)
}

// addARNToVM generates the ARN of a VM instance
func (d computeDiscovery) addARNToVM(vm *typesEC2.Instance) voc.ResourceID {
	return voc.ResourceID("arn:aws:ec2:" +
		d.awsConfig.cfg.Region + ":" +
		aws.ToString(d.awsConfig.accountID) +
		":instance/" +
		aws.ToString(vm.InstanceId))
}
