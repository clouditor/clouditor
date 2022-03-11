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

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
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

// NewAwsComputeDiscovery constructs a new awsS3Discovery initializing the s3-virtualMachineAPI and isDiscovering with true
func NewAwsComputeDiscovery(client *Client) discovery.Discoverer {
	return &computeDiscovery{
		virtualMachineAPI: newFromConfigEC2(client.cfg),
		functionAPI:       newFromConfigLambda(client.cfg),
		isDiscovering:     true,
		awsConfig:         client,
	}
}

// Name is the method implementation defined in the discovery.Discoverer interface
func (*computeDiscovery) Name() string {
	return "AWS Compute"
}

// List is the method implementation defined in the discovery.Discoverer interface
func (d computeDiscovery) List() (resources []voc.IsCloudResource, err error) {
	log.Infof("Collecting evidences in %s", d.Name())
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
				Resource: &voc.Resource{
					ID:           d.addARNToVM(vm),
					Name:         d.getNameOfVM(vm),
					CreationTime: 0,
					Type:         []string{"VirtualMachine", "Compute", "Resource"},
					GeoLocation: voc.GeoLocation{
						Region: d.awsConfig.cfg.Region,
					},
				},
				NetworkInterface: d.getNetworkInterfacesOfVM(vm),
			}

			resources = append(resources, voc.VirtualMachine{
				Compute:      computeResource,
				BlockStorage: d.mapBlockStorageIDsOfVM(vm),
				BootLogging:  d.getBootLog(vm),
				OSLogging:    d.getOSLog(vm),
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
				Resource: &voc.Resource{
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
func (*computeDiscovery) getBootLog(_ *typesEC2.Instance) (l *voc.BootLogging) {
	l = &voc.BootLogging{
		Logging: &voc.Logging{
			Auditing:        nil,
			LoggingService:  nil,
			Enabled:         false,
			RetentionPeriod: 0,
		},
	}
	return
}

// getOSLog checks if OS logging is enabled
// Currently there is no option to find out if any logs are enabled -> Assign default zero values
func (*computeDiscovery) getOSLog(_ *typesEC2.Instance) (l *voc.OSLogging) {
	l = &voc.OSLogging{
		Logging: &voc.Logging{
			Auditing:        nil,
			LoggingService:  nil,
			Enabled:         false,
			RetentionPeriod: 0,
		},
	}
	return
}

// mapBlockStorageIDsOfVM returns block storages IDs by iterating the VMs block storages
func (*computeDiscovery) mapBlockStorageIDsOfVM(vm *typesEC2.Instance) (blockStorageIDs []voc.ResourceID) {
	// Loop through mappings using an index, since BlockDeviceMappings is an array of a struct
	// and not of a pointer; otherwise we would copy a lot of data
	for i := range vm.BlockDeviceMappings {
		mapping := &vm.BlockDeviceMappings[i]
		blockStorageIDs = append(blockStorageIDs, voc.ResourceID(aws.ToString(mapping.Ebs.VolumeId)))
	}
	return
}

// getNetworkInterfacesOfVM returns the network interface IDs by iterating the VMs network interfaces
func (*computeDiscovery) getNetworkInterfacesOfVM(vm *typesEC2.Instance) (networkInterfaceIDs []voc.ResourceID) {
	// Loop through mappings using an index, since is NetworkInterfaces an array of a struct
	// and not of a pointer; otherwise we would copy a lot of data
	for i := range vm.NetworkInterfaces {
		ifc := &vm.NetworkInterfaces[i]
		networkInterfaceIDs = append(networkInterfaceIDs, voc.ResourceID(aws.ToString(ifc.NetworkInterfaceId)))
	}
	return
}

// getNameOfVM returns the name if exists (i.e. a tag with key 'name' exists), otherwise instance ID is used
func (*computeDiscovery) getNameOfVM(vm *typesEC2.Instance) string {
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
