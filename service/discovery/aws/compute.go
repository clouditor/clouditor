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

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/api/ontology"
	"clouditor.io/clouditor/internal/util"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	typesEC2 "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	typesLambda "github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

// computeDiscovery handles the AWS API requests regarding the computing services (EC2 and Lambda)
type computeDiscovery struct {
	virtualMachineAPI EC2API
	functionAPI       LambdaAPI
	isDiscovering     bool
	awsConfig         *Client
	csID              string
}

// EC2API describes the EC2 api interface which is implemented by the official AWS client and mock clients in tests
type EC2API interface {
	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(options *ec2.Options)) (*ec2.DescribeInstancesOutput, error)

	DescribeVolumes(ctx context.Context,
		params *ec2.DescribeVolumesInput,
		optFns ...func(options *ec2.Options)) (*ec2.DescribeVolumesOutput, error)

	DescribeNetworkInterfaces(ctx context.Context,
		params *ec2.DescribeNetworkInterfacesInput,
		optFns ...func(options *ec2.Options)) (*ec2.DescribeNetworkInterfacesOutput, error)
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
func NewAwsComputeDiscovery(client *Client, cloudServiceID string) discovery.Discoverer {
	return &computeDiscovery{
		virtualMachineAPI: newFromConfigEC2(client.cfg),
		functionAPI:       newFromConfigLambda(client.cfg),
		isDiscovering:     true,
		awsConfig:         client,
		csID:              cloudServiceID,
	}
}

// Name is the method implementation defined in the discovery.Discoverer interface
func (*computeDiscovery) Name() string {
	return "AWS Compute"
}

// List is the method implementation defined in the discovery.Discoverer interface
func (d *computeDiscovery) List() (resources []ontology.IsResource, err error) {
	log.Infof("Collecting evidences in %s", d.Name())

	// Even though technically volumes are "storage", they are part of the EC2 API and therefore discovered here
	volumes, err := d.discoverVolumes()
	if err != nil {
		return nil, fmt.Errorf("could not discover volumes: %w", err)
	}
	for _, volume := range volumes {
		resources = append(resources, volume)
	}

	// Even though technically network interfaces are "network", they are part of the EC2 API and therefore discovered here
	ifcs, err := d.discoverNetworkInterfaces()
	if err != nil {
		return nil, fmt.Errorf("could not discover volumes: %w", err)
	}
	for _, ifc := range ifcs {
		resources = append(resources, ifc)
	}

	listOfVMs, err := d.discoverVirtualMachines()
	if err != nil {
		return nil, fmt.Errorf("could not discover virtual machines: %w", err)
	}
	for _, machine := range listOfVMs {
		resources = append(resources, machine)
	}

	listOfFunctions, err := d.discoverFunctions()
	if err != nil {
		return nil, fmt.Errorf("could not discover functions: %w", err)
	}
	for _, function := range listOfFunctions {
		resources = append(resources, function)
	}

	return
}

func (d *computeDiscovery) CloudServiceID() string {
	return d.csID
}

// discoverVolumes discovers all volumes (in the current region)
func (d *computeDiscovery) discoverVolumes() ([]*ontology.BlockStorage, error) {
	res, err := d.virtualMachineAPI.DescribeVolumes(context.TODO(), &ec2.DescribeVolumesInput{})
	if err != nil {
		return nil, prettyError(err)
	}

	var blocks []*ontology.BlockStorage
	for i := range res.Volumes {
		volume := &res.Volumes[i]

		atRest := &ontology.ManagedKeyEncryption{
			Enabled: util.Deref(volume.Encrypted),
		}

		// AWS uses a fixed algorithm, if enabled
		if atRest.Enabled {
			atRest.Algorithm = "AES-256"
		}

		blocks = append(blocks, &ontology.BlockStorage{
			Id:           d.arnify("volume", volume.VolumeId),
			Name:         d.nameOrID(volume.Tags, volume.VolumeId),
			CreationTime: timestamppb.New(util.Deref(volume.CreateTime)),
			GeoLocation: &ontology.GeoLocation{
				Region: d.awsConfig.cfg.Region,
			},
			Labels: d.labels(volume.Tags),
			AtRestEncryption: &ontology.AtRestEncryption{
				Type: &ontology.AtRestEncryption_ManagedKeyEncryption{
					ManagedKeyEncryption: atRest,
				},
			},
			Raw: discovery.Raw(&res.Volumes[i]),
		})
	}

	return blocks, nil
}

// discoverNetworkInterfaces discovers all network interfaces (in the current region)
func (d *computeDiscovery) discoverNetworkInterfaces() ([]*ontology.NetworkInterface, error) {
	res, err := d.virtualMachineAPI.DescribeNetworkInterfaces(context.TODO(), &ec2.DescribeNetworkInterfacesInput{})
	if err != nil {
		return nil, prettyError(err)
	}

	var ifcs []*ontology.NetworkInterface
	for i := range res.NetworkInterfaces {
		ifc := &res.NetworkInterfaces[i]

		ifcs = append(ifcs, &ontology.NetworkInterface{
			Id:   d.arnify("network-interface", ifc.NetworkInterfaceId),
			Name: d.nameOrID(ifc.TagSet, ifc.NetworkInterfaceId),
			GeoLocation: &ontology.GeoLocation{
				Region: d.awsConfig.cfg.Region,
			},
			Labels: d.labels(ifc.TagSet),
			Raw:    discovery.Raw(&res.NetworkInterfaces[i]),
		})
	}

	return ifcs, nil
}

// discoverVirtualMachines discovers all VMs (in the current region)
func (d *computeDiscovery) discoverVirtualMachines() ([]*ontology.VirtualMachine, error) {
	resp, err := d.virtualMachineAPI.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, prettyError(err)
	}
	var resources []*ontology.VirtualMachine
	for _, reservation := range resp.Reservations {
		for i := range reservation.Instances {
			vm := &reservation.Instances[i]

			resources = append(resources, &ontology.VirtualMachine{
				Id:   d.arnify("instance", vm.InstanceId),
				Name: d.getNameOfVM(vm),
				GeoLocation: &ontology.GeoLocation{
					Region: d.awsConfig.cfg.Region,
				},
				Labels:              d.labels(vm.Tags),
				NetworkInterfaceIds: d.getNetworkInterfacesOfVM(vm),
				BlockStorageIds:     d.mapBlockStorageIDsOfVM(vm),
				BootLogging:         d.getBootLog(vm),
				OsLogging:           d.getOSLog(vm),
				Raw:                 discovery.Raw(&reservation),
			})
		}
	}

	return resources, nil
}

// discoverFunctions discovers all lambda functions
func (d *computeDiscovery) discoverFunctions() (resources []*ontology.Function, err error) {
	// 'listFunctions' discovers up to 50 Lambda functions per execution -> loop through when response has nextMarker set
	var resp *lambda.ListFunctionsOutput
	var nextMarker *string
	for {
		resp, err = d.functionAPI.ListFunctions(context.TODO(), &lambda.ListFunctionsInput{
			Marker: nextMarker,
		})
		if err != nil {
			return nil, prettyError(err)
		}
		resources = append(resources, d.mapFunctionResources(resp.Functions)...)

		if nextMarker = resp.NextMarker; nextMarker == nil {
			break
		}
	}

	return
}

// mapFunctionResources iterates functionConfigurations and returns a list of corresponding FunctionResources
func (d *computeDiscovery) mapFunctionResources(functions []typesLambda.FunctionConfiguration) (resources []*ontology.Function) {
	// TODO(all): Labels are missing
	for i := range functions {
		function := &functions[i]

		resources = append(resources, &ontology.Function{
			Id:   aws.ToString(function.FunctionArn),
			Name: aws.ToString(function.FunctionName),
			GeoLocation: &ontology.GeoLocation{
				Region: d.awsConfig.cfg.Region,
			},
			Raw: discovery.Raw(&functions[i]),
		})
	}
	return
}

// getBootLog checks if boot logging is enabled
// Currently there is no option to find out if any logs are enabled -> Assign default zero values
func (*computeDiscovery) getBootLog(_ *typesEC2.Instance) (l *ontology.BootLogging) {
	l = &ontology.BootLogging{
		Enabled: false,
	}
	return
}

// getOSLog checks if OS logging is enabled
// Currently there is no option to find out if any logs are enabled -> Assign default zero values
func (*computeDiscovery) getOSLog(_ *typesEC2.Instance) (l *ontology.OSLogging) {
	l = &ontology.OSLogging{
		Enabled: false,
	}
	return
}

// mapBlockStorageIDsOfVM returns block storages IDs by iterating the VMs block storages
func (d *computeDiscovery) mapBlockStorageIDsOfVM(vm *typesEC2.Instance) (blockStorageIDs []string) {
	// Loop through mappings using an index, since BlockDeviceMappings is an array of a struct
	// and not of a pointer; otherwise we would copy a lot of data
	for i := range vm.BlockDeviceMappings {
		mapping := &vm.BlockDeviceMappings[i]
		blockStorageIDs = append(blockStorageIDs, d.arnify("volume", mapping.Ebs.VolumeId))
	}
	return
}

// getNetworkInterfacesOfVM returns the network interface IDs by iterating the VMs network interfaces
func (d *computeDiscovery) getNetworkInterfacesOfVM(vm *typesEC2.Instance) (networkInterfaceIDs []string) {
	// Loop through mappings using an index, since is NetworkInterfaces an array of a struct
	// and not of a pointer; otherwise we would copy a lot of data
	for i := range vm.NetworkInterfaces {
		ifc := &vm.NetworkInterfaces[i]
		networkInterfaceIDs = append(networkInterfaceIDs, d.arnify("network-interface", ifc.NetworkInterfaceId))
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

// nameOrID returns the name if exists (i.e. a tag with key 'name' exists), otherwise instance ID is used
func (*computeDiscovery) nameOrID(tags []typesEC2.Tag, ID *string) string {
	for _, tag := range tags {
		if aws.ToString(tag.Key) == "Name" {
			return aws.ToString(tag.Value)
		}
	}

	// If no tag with 'name' was found, return ID instead
	return aws.ToString(ID)
}

func (*computeDiscovery) labels(tags []typesEC2.Tag) (labels map[string]string) {
	labels = map[string]string{}

	for _, tag := range tags {
		labels[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
	}

	return
}

// addARNToVolume generates the ARN of a volume instance
func (d *computeDiscovery) arnify(typ string, ID *string) string {
	return "arn:aws:ec2:" +
		d.awsConfig.cfg.Region + ":" +
		aws.ToString(d.awsConfig.accountID) +
		":" + typ + "/" +
		aws.ToString(ID)
}
