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
	"clouditor.io/clouditor/internal/constants"
	"context"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	typesEC2 "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	typesLambda "github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

// LatestLambdaGoVersion is the latest go version supported by Lambda functions:
// According to doc it is always the latest version of the official go releases
const LatestLambdaGoVersion = "20"

// LatestLambdaNodeJSVersion is the latest node.js version supported by Lambda functions:
// See aws-sdk-go-v2/service/lambda@v1.37.0/types/enums.go:338
const LatestLambdaNodeJSVersion = "16"

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
func (d *computeDiscovery) List() (resources []voc.IsCloudResource, err error) {
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

// discoverVolumes discoveres all volumes (in the current region)
func (d *computeDiscovery) discoverVolumes() ([]*voc.BlockStorage, error) {
	res, err := d.virtualMachineAPI.DescribeVolumes(context.TODO(), &ec2.DescribeVolumesInput{})
	if err != nil {
		return nil, prettyError(err)
	}

	var blocks []*voc.BlockStorage
	for i := range res.Volumes {
		volume := &res.Volumes[i]

		atRest := &voc.AtRestEncryption{
			Enabled: util.Deref(volume.Encrypted),
		}

		// AWS uses a fixed algorithm, if enabled
		if atRest.Enabled {
			atRest.Algorithm = "AES-256"
		}

		blocks = append(blocks, &voc.BlockStorage{
			Storage: &voc.Storage{
				Resource: discovery.NewResource(d,
					d.arnify("volume", volume.VolumeId),
					d.nameOrID(volume.Tags, volume.VolumeId),
					volume.CreateTime,
					voc.GeoLocation{
						Region: d.awsConfig.cfg.Region,
					},
					d.labels(volume.Tags),
					voc.BlockStorageType,
					&res.Volumes[i],
				),
				AtRestEncryption: atRest,
			},
		})
	}

	return blocks, nil
}

// discoverNetworkInterfaces discovers all network interfaces (in the current region)
func (d *computeDiscovery) discoverNetworkInterfaces() ([]voc.NetworkInterface, error) {
	res, err := d.virtualMachineAPI.DescribeNetworkInterfaces(context.TODO(), &ec2.DescribeNetworkInterfacesInput{})
	if err != nil {
		return nil, prettyError(err)
	}

	var ifcs []voc.NetworkInterface
	for i := range res.NetworkInterfaces {
		ifc := &res.NetworkInterfaces[i]

		ifcs = append(ifcs, voc.NetworkInterface{
			Networking: &voc.Networking{
				Resource: discovery.NewResource(
					d,
					d.arnify("network-interface", ifc.NetworkInterfaceId),
					d.nameOrID(ifc.TagSet, ifc.NetworkInterfaceId),
					nil,
					voc.GeoLocation{
						Region: d.awsConfig.cfg.Region,
					},
					d.labels(ifc.TagSet),
					[]string{"NetworkInterface", "Networking", "Resource"},
					&res.NetworkInterfaces[i],
				),
			},
		})
	}

	return ifcs, nil
}

// discoverVirtualMachines discovers all VMs (in the current region)
func (d *computeDiscovery) discoverVirtualMachines() ([]*voc.VirtualMachine, error) {
	resp, err := d.virtualMachineAPI.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, prettyError(err)
	}
	var resources []*voc.VirtualMachine
	for _, reservation := range resp.Reservations {
		for i := range reservation.Instances {
			vm := &reservation.Instances[i]
			computeResource := &voc.Compute{
				Resource: discovery.NewResource(
					d,
					d.arnify("instance", vm.InstanceId),
					d.getNameOfVM(vm),
					nil,
					voc.GeoLocation{
						Region: d.awsConfig.cfg.Region,
					},
					d.labels(vm.Tags),
					[]string{"VirtualMachine", "Compute", "Resource"},
					&reservation,
				),

				NetworkInterfaces: d.getNetworkInterfacesOfVM(vm),
			}

			resources = append(resources, &voc.VirtualMachine{
				Compute:      computeResource,
				BlockStorage: d.mapBlockStorageIDsOfVM(vm),
				BootLogging:  d.getBootLog(vm),
				OsLogging:    d.getOSLog(vm),
			})
		}
	}
	return resources, nil
}

// discoverFunctions discovers all lambda functions
func (d *computeDiscovery) discoverFunctions() (resources []*voc.Function, err error) {
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
func (d *computeDiscovery) mapFunctionResources(functions []typesLambda.FunctionConfiguration) (resources []*voc.Function) {
	// TODO(all): Labels are missing
	for i := range functions {
		f := &functions[i]

		resources = append(resources, &voc.Function{
			// General Compute fields
			Compute: &voc.Compute{
				Resource: discovery.NewResource(
					d,
					voc.ResourceID(aws.ToString(f.FunctionArn)),
					aws.ToString(f.FunctionName),
					nil,
					voc.GeoLocation{
						Region: d.awsConfig.cfg.Region,
					},
					nil,
					voc.FunctionType,
					&functions[i],
				),
			},
			// Function-specific fields
			// TODO(lebogg): Test
			RuntimeLanguage: toRuntimeLanguage(f.Runtime),
			RuntimeVersion:  toRuntimeVersion(f.Runtime),
		})
	}
	return
}

// toRuntimeLanguage returns the runtime language of runtime (runtime contains both language und version in one string)
func toRuntimeLanguage(runtime typesLambda.Runtime) (language string) {
	language, _ = splitRuntime(runtime)
	return
}

// toRuntimeVersion returns the runtime version of runtime (runtime contains both language und version in one string)
func toRuntimeVersion(runtime typesLambda.Runtime) (version string) {
	_, version = splitRuntime(runtime)
	return
}

// splitRuntime splits runtime into the runtime language and version. It goes through the string, character by
// character, and divides the string when the first digit is reached. If there is no version number, the most recent
// version supported by lambda functions for this language is assumed.
func splitRuntime(runtime typesLambda.Runtime) (language, version string) {
	input := string(runtime)
	var separator int
	for i := 0; i < len(input); i++ {
		if unicode.IsDigit(rune(input[i])) {
			separator, _ = strconv.Atoi(string(input[i]))
			break
		}
	}
	strArr := strings.SplitN(input, strconv.Itoa(separator), 2)
	// Go lambda functions always use the latest go version according to documentation. See
	// https://github.com/aws-samples/sessions-with-aws-sam/tree/master/go-al2#golang-installation)
	if l := strArr[0]; l == "go" {
		return useOfficialLanguageName(l), LatestLambdaGoVersion
	}
	// For node.js, there is the possibility that only nodejs is returned w/o a version attached. We assume the latest
	// version supported by lambda functions
	if l := strArr[0]; l == "nodejs" && len(strArr) == 1 {
		return useOfficialLanguageName(l), LatestLambdaNodeJSVersion
	}
	// Currently not reachable but to avoid "index out of range" error in the future when a new language might be
	// supported which we won't have considered yet.
	if l := strArr[0]; len(strArr) == 1 {
		log.Warnf("Runtime '%s' is not considered yet. Maybe it got newly introduced into AWS lambdas.",
			l)
		return l, ""
	}
	return useOfficialLanguageName(strArr[0]), strconv.Itoa(separator) + strArr[1]
}

// useOfficialLanguageName converts the given language names specified by AWS to the names needed for the Rego
// expressions
func useOfficialLanguageName(l string) string {
	switch l {
	case "go":
		return constants.Go
	case "java":
		return constants.Java
	case "nodejs":
		return constants.NodeJS
	default:
		return l
	}
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
func (d *computeDiscovery) mapBlockStorageIDsOfVM(vm *typesEC2.Instance) (blockStorageIDs []voc.ResourceID) {
	// Loop through mappings using an index, since BlockDeviceMappings is an array of a struct
	// and not of a pointer; otherwise we would copy a lot of data
	for i := range vm.BlockDeviceMappings {
		mapping := &vm.BlockDeviceMappings[i]
		blockStorageIDs = append(blockStorageIDs, d.arnify("volume", mapping.Ebs.VolumeId))
	}
	return
}

// getNetworkInterfacesOfVM returns the network interface IDs by iterating the VMs network interfaces
func (d *computeDiscovery) getNetworkInterfacesOfVM(vm *typesEC2.Instance) (networkInterfaceIDs []voc.ResourceID) {
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

// addARNToVolume generates the ARN of a volumne instance
func (d *computeDiscovery) arnify(typ string, ID *string) voc.ResourceID {
	return voc.ResourceID("arn:aws:ec2:" +
		d.awsConfig.cfg.Region + ":" +
		aws.ToString(d.awsConfig.accountID) +
		":" + typ + "/" +
		aws.ToString(ID))
}
