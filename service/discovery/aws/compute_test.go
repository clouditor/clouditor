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
	"clouditor.io/clouditor/internal/util"
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	types2 "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"reflect"
	"testing"
	"time"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/voc"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	typesEC2 "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	typesLambda "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/middleware"
	"github.com/stretchr/testify/assert"
)

const (
	mockVM1            = "mockVM1"
	mockVM1ID          = "i-0b0c58ade95f269f7"
	blockVolumeId      = "blockVolumeID"
	networkInterfaceId = "networkInterfaceId"
	mockVMCreationTime = "2012-11-01T22:08:41+00:00"

	mockFunction1ID           = "arn:aws:lambda:eu-central-1:123456789:function:mock-function:1"
	mockFunction1             = "MockFunction1"
	mockFunction1Region       = "eu-central-1"
	mockFunction1CreationTime = "2012-11-01T22:08:41.0+00:00"

	mockDefaultBaselineID = "pb-09bcfbcf275c8c953"
	// mockDefaultBaselineARN includes mockDefaultBaselineID at the end
	mockDefaultBaselineARN = "arn:aws:ssm:eu-central-1:416089608788:patchbaseline/pb-09bcfbcf275c8c953"
)

// mockEC2API implements the EC2API interface for mock testing
type mockEC2API struct {
}

// mockEC2APIWithErrors implements the EC2API interface (API call returning error) for mock testing
type mockEC2APIWithErrors struct {
}

// mockLambdaAPI implements the LambdaAPI interface for mock testing
type mockLambdaAPI struct {
}

// mockLambdaAPI implements the LambdaAPI interface for mock testing if >50 Lambda functions are discovered (not only 50)
type mockLambdaAPI51LambdaFunctions struct {
}

// mockLambdaAPIWithErrors implements the LambdaAPI interface (API call returning error) for mock testing
type mockLambdaAPIWithErrors struct {
}

// mockSSMAPI implements the SystemsManagerAPI interface for mock testing
type mockSSMAPI struct{}

func (m mockSSMAPI) DescribeInstancePatchStates(ctx context.Context, params *ssm.DescribeInstancePatchStatesInput, optFns ...func(*ssm.Options)) (*ssm.DescribeInstancePatchStatesOutput, error) {
	switch params.InstanceIds[0] {
	case mockVM1ID:
		return &ssm.DescribeInstancePatchStatesOutput{
			InstancePatchStates: []types2.InstancePatchState{
				{
					BaselineId: util.Ref(mockDefaultBaselineID),
					InstanceId: util.Ref(params.InstanceIds[0]),
					PatchGroup: nil,
				},
			},
		}, nil
	}
	panic("implement me")
}

func (m mockSSMAPI) GetPatchBaseline(ctx context.Context, params *ssm.GetPatchBaselineInput, optFns ...func(*ssm.Options)) (*ssm.GetPatchBaselineOutput, error) {
	//TODO(lebogg): Add other cases
	switch util.Deref(params.BaselineId) {
	case mockVM1ID:
		return &ssm.GetPatchBaselineOutput{
			ApprovedPatchesEnableNonSecurity: util.Ref(false),
			BaselineId:                       params.BaselineId,
			ResultMetadata:                   middleware.Metadata{},
		}, nil
	default:
		return &ssm.GetPatchBaselineOutput{
			ApprovedPatchesEnableNonSecurity: util.Ref(false),
			BaselineId:                       params.BaselineId,
			ResultMetadata:                   middleware.Metadata{},
		}, nil
	}

}

// ListFunctions is the method implementation of the LambdaAPI interface
func (mockLambdaAPI) ListFunctions(_ context.Context, _ *lambda.ListFunctionsInput, _ ...func(*lambda.Options)) (*lambda.ListFunctionsOutput, error) {
	return &lambda.ListFunctionsOutput{
		Functions: []typesLambda.FunctionConfiguration{
			{
				FunctionArn:  aws.String(mockFunction1ID),
				FunctionName: aws.String(mockFunction1),
				LastModified: aws.String(mockFunction1CreationTime),
				Runtime:      "Java11",
			},
		},
		NextMarker:     nil,
		ResultMetadata: middleware.Metadata{},
	}, nil
}

func (mockLambdaAPI51LambdaFunctions) ListFunctions(_ context.Context, input *lambda.ListFunctionsInput, _ ...func(*lambda.Options)) (output *lambda.ListFunctionsOutput, err error) {
	var lambdaFunctions []typesLambda.FunctionConfiguration
	nextMarker := "ShowNext"
	if input.Marker == nil {
		for i := 0; i < 50; i++ {
			lambdaFunctions = append(lambdaFunctions, typesLambda.FunctionConfiguration{
				// We have to set a time in a right format, otherwise the discoverer fails (parse error)
				LastModified: aws.String(mockFunction1CreationTime),
			})
		}
		output = &lambda.ListFunctionsOutput{
			Functions:  lambdaFunctions,
			NextMarker: aws.String(nextMarker),
		}
	} else if *input.Marker == nextMarker {
		for i := 0; i < 5; i++ {
			lambdaFunctions = append(lambdaFunctions, typesLambda.FunctionConfiguration{
				// We have to set a time in a right format, otherwise the discoverer fails (parse error)
				LastModified: aws.String(mockFunction1CreationTime),
			})
		}
		output = &lambda.ListFunctionsOutput{
			Functions:  lambdaFunctions,
			NextMarker: nil,
		}
	}
	return
}

// ListFunctions is the method implementation of the LambdaAPI interface
func (mockLambdaAPIWithErrors) ListFunctions(_ context.Context, _ *lambda.ListFunctionsInput, _ ...func(*lambda.Options)) (*lambda.ListFunctionsOutput, error) {
	err := &smithy.GenericAPIError{
		Code:    "500",
		Message: "Internal Server Error",
	}
	return nil, err
}

// DescribeInstances is the method implementation of the EC2API interface
func (mockEC2API) DescribeInstances(_ context.Context, _ *ec2.DescribeInstancesInput, _ ...func(options *ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
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
			Key:   aws.String("Name"),
			Value: aws.String(mockVM1),
		},
	}
	networkInterfaces := []types.InstanceNetworkInterface{
		{
			NetworkInterfaceId: aws.String(networkInterfaceId),
		},
	}
	// launch time
	launchTime, err := time.Parse(time.RFC3339, mockVMCreationTime)
	if err != nil {
		log.Error(err)
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
				LaunchTime:          aws.Time(launchTime),
			}},
			OwnerId:       nil,
			RequesterId:   nil,
			ReservationId: nil,
		}},
		ResultMetadata: middleware.Metadata{},
	}
	return output, nil
}

// DescribeVolumes is the method implementation of the EC2API interface
func (mockEC2API) DescribeVolumes(_ context.Context, _ *ec2.DescribeVolumesInput, _ ...func(options *ec2.Options)) (*ec2.DescribeVolumesOutput, error) {
	output := &ec2.DescribeVolumesOutput{
		NextToken: nil,
		Volumes: []types.Volume{
			{
				VolumeId:   aws.String(blockVolumeId),
				CreateTime: aws.Time(time.Now()),
				Tags: []types.Tag{
					{Key: aws.String("Name"), Value: aws.String("My Volume")},
				},
			},
			{
				VolumeId:   aws.String("othervolume"),
				CreateTime: aws.Time(time.Now()),
			},
		},
		ResultMetadata: middleware.Metadata{},
	}

	return output, nil
}

// DescribeNetworkInterfaces is the method implementation of the EC2API interface
func (mockEC2API) DescribeNetworkInterfaces(_ context.Context, _ *ec2.DescribeNetworkInterfacesInput, _ ...func(options *ec2.Options)) (*ec2.DescribeNetworkInterfacesOutput, error) {
	output := &ec2.DescribeNetworkInterfacesOutput{
		NextToken: nil,
		NetworkInterfaces: []types.NetworkInterface{
			{
				NetworkInterfaceId: aws.String(networkInterfaceId),
				TagSet: []types.Tag{
					{Key: aws.String("Name"), Value: aws.String("My Network Interface")},
				},
			},
		},
		ResultMetadata: middleware.Metadata{},
	}

	return output, nil
}

// DescribeInstances is the method implementation of the EC2API interface
func (mockEC2APIWithErrors) DescribeInstances(_ context.Context, _ *ec2.DescribeInstancesInput, _ ...func(options *ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	err := &smithy.GenericAPIError{
		Code:    "ConnectionError",
		Message: "Couldn't resolve host. Bad connection?",
	}
	return nil, err
}

// DescribeVolumes is the method implementation of the EC2API interface
func (mockEC2APIWithErrors) DescribeVolumes(_ context.Context, _ *ec2.DescribeVolumesInput, _ ...func(options *ec2.Options)) (*ec2.DescribeVolumesOutput, error) {
	err := &smithy.GenericAPIError{
		Code:    "ConnectionError",
		Message: "Couldn't resolve host. Bad connection?",
	}
	return nil, err
}

// DescribeNetworkInterfaces is the method implementation of the EC2API interface
func (mockEC2APIWithErrors) DescribeNetworkInterfaces(_ context.Context, _ *ec2.DescribeNetworkInterfacesInput, _ ...func(options *ec2.Options)) (*ec2.DescribeNetworkInterfacesOutput, error) {
	err := &smithy.GenericAPIError{
		Code:    "ConnectionError",
		Message: "Couldn't resolve host. Bad connection?",
	}
	return nil, err
}

func TestComputeDiscovery_List(t *testing.T) {
	d := computeDiscovery{
		virtualMachineAPI: mockEC2API{},
		functionAPI:       mockLambdaAPI{},
		isDiscovering:     true,
		awsConfig: &Client{
			cfg: aws.Config{
				Region: "eu-central-1",
			},
			accountID: aws.String("MockAccountID1234"),
		},
	}
	list, err := d.List()
	assert.NoError(t, err)
	assert.NotEmpty(t, list)

	d = computeDiscovery{
		virtualMachineAPI: mockEC2APIWithErrors{},
	}
	_, err = d.List()
	assert.Error(t, err)

	d = computeDiscovery{
		virtualMachineAPI: mockEC2API{},
		functionAPI:       mockLambdaAPIWithErrors{},
		isDiscovering:     true,
		awsConfig: &Client{
			cfg: aws.Config{
				Region: "eu-central-1",
			},
			accountID: aws.String("MockAccountID1234"),
		},
	}
	_, err = d.List()
	assert.Error(t, err)
}

func TestComputeDiscovery_discoverVirtualMachines(t *testing.T) {
	d := computeDiscovery{
		virtualMachineAPI: mockEC2API{},
		isDiscovering:     true,
		awsConfig: &Client{
			cfg: aws.Config{
				Region: "eu-central-1",
			},
			accountID: aws.String("MockAccountID1234"),
		},
	}
	machines, err := d.discoverVirtualMachines()
	assert.NoError(t, err)
	testMachine := machines[0]
	assert.Equal(t, mockVM1, testMachine.Name)
	assert.Equal(t, voc.ResourceID("arn:aws:ec2:eu-central-1:MockAccountID1234:instance/mockVM1ID"), testMachine.ID)
	assert.NotEmpty(t, testMachine.BlockStorage)
	assert.False(t, testMachine.BootLogging.Enabled)
	assert.False(t, testMachine.OsLogging.Enabled)
	assert.Equal(t, int64(0), testMachine.CreationTime)
	assert.Equal(t, mockFunction1Region, testMachine.GeoLocation.Region)

	d = computeDiscovery{
		virtualMachineAPI: mockEC2APIWithErrors{},
	}
	_, err = d.discoverVirtualMachines()
	assert.Error(t, err)

}

func TestComputeDiscover_Name(t *testing.T) {
	d := computeDiscovery{
		virtualMachineAPI: mockEC2API{},
		isDiscovering:     true,
		awsConfig:         &Client{},
	}
	assert.Equal(t, "AWS Compute", d.Name())
}

func TestComputeDiscovery_getNameOfVM(t *testing.T) {
	type fields struct {
		api           EC2API
		isDiscovering bool
		awsConfig     *Client
	}
	type args struct {
		vm types.Instance
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"First Test without tag",
			fields{},
			args{vm: types.Instance{InstanceId: aws.String(mockVM1ID)}},
			mockVM1ID,
		},
		{
			"Second test with tag",
			fields{},
			args{vm: types.Instance{InstanceId: aws.String(mockVM1ID),
				Tags: []types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(mockVM1),
					},
				}}},
			mockVM1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &computeDiscovery{
				virtualMachineAPI: tt.fields.api,
				isDiscovering:     tt.fields.isDiscovering,
				awsConfig:         tt.fields.awsConfig,
			}
			if got := d.getNameOfVM(&tt.args.vm); got != tt.want {
				t.Errorf("getNameOfVM() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComputeDiscovery_discoverFunctions(t *testing.T) {
	type fields struct {
		virtualMachineAPI EC2API
		functionAPI       LambdaAPI
		isDiscovering     bool
		awsConfig         *Client
		csID              string
	}
	mockClient := &Client{
		cfg: aws.Config{
			Region: "eu-central-1",
		},
	}
	// creationTime, _ := time.Parse(time.RFC3339, mockFunction1CreationTime)
	tests := []struct {
		name    string
		fields  fields
		want    assert.ValueAssertionFunc
		wantErr bool
	}{
		{
			name: "Happy path",
			fields: fields{
				functionAPI: mockLambdaAPI{},
				awsConfig:   mockClient,
				csID:        testdata.MockCloudServiceID1,
			},
			want: func(t assert.TestingT, i1 interface{}, i ...interface{}) bool {
				functions, ok := i1.([]*voc.Function)
				assert.True(t, ok)
				f := functions[0]
				assert.Equal(t, mockClient.cfg.Region, f.GeoLocation.Region)
				assert.Equal(t, "Java", f.RuntimeLanguage)
				assert.Equal(t, "11", f.RuntimeVersion)
				return assert.Equal(t, voc.ResourceID(mockFunction1ID), f.ID)
			},
			wantErr: false,
		},
		{
			name: "Error - ",
			fields: fields{
				functionAPI: mockLambdaAPIWithErrors{},
			},
			want:    assert.Nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &computeDiscovery{
				virtualMachineAPI: tt.fields.virtualMachineAPI,
				functionAPI:       tt.fields.functionAPI,
				isDiscovering:     tt.fields.isDiscovering,
				awsConfig:         tt.fields.awsConfig,
				csID:              tt.fields.csID,
			}
			got, err := d.discoverFunctions()
			if (err != nil) != tt.wantErr {
				t.Errorf("discoverFunctions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.want(t, got)
		})
	}

	// Testing the case where two API Calls have to be made due to limit of returned functions
	d := computeDiscovery{
		functionAPI: mockLambdaAPI51LambdaFunctions{},
		awsConfig:   mockClient,
	}
	functions, err := d.discoverFunctions()
	assert.NoError(t, err)
	assert.Less(t, 50, len(functions))

}

func TestComputeDiscovery_NewComputeDiscovery(t *testing.T) {
	// Mock newFromConfigs and store the original functions back at the end of the test
	oldEC2 := newFromConfigEC2
	defer func() { newFromConfigEC2 = oldEC2 }()
	oldLambda := newFromConfigLambda
	defer func() { newFromConfigLambda = oldLambda }()

	newFromConfigEC2 = func(cfg aws.Config, optFns ...func(*ec2.Options)) *ec2.Client {
		return &ec2.Client{}
	}
	newFromConfigLambda = func(cfg aws.Config, optFns ...func(*lambda.Options)) *lambda.Client {
		return &lambda.Client{}
	}

	type args struct {
		client *Client
		csID   string
	}
	mockClient := &Client{
		cfg: aws.Config{
			Region: "eu-central-1",
		},
		accountID: aws.String("1234"),
	}
	tests := []struct {
		name string
		args args
		want discovery.Discoverer
	}{
		{
			args: args{client: mockClient, csID: testdata.MockCloudServiceID1},
			want: &computeDiscovery{
				virtualMachineAPI: &ec2.Client{},
				functionAPI:       &lambda.Client{},
				isDiscovering:     true,
				awsConfig:         mockClient,
				csID:              testdata.MockCloudServiceID1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAwsComputeDiscovery(tt.args.client, tt.args.csID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAwsComputeDiscovery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitRuntime(t *testing.T) {
	type args struct {
		runtime typesLambda.Runtime
	}
	tests := []struct {
		name         string
		args         args
		wantLanguage string
		wantVersion  string
	}{
		{
			name:         "Nodejs without version",
			args:         args{runtime: typesLambda.RuntimeNodejs},
			wantLanguage: string(typesLambda.RuntimeNodejs),
			wantVersion:  LatestLambdaNodeJSVersion,
		},
		{
			name:         "Nodejs with version",
			args:         args{runtime: typesLambda.RuntimeNodejs12x},
			wantLanguage: string(typesLambda.RuntimeNodejs),
			wantVersion:  "12.x",
		},
		{
			name:         "Java with version",
			args:         args{runtime: typesLambda.RuntimeJava11},
			wantLanguage: "java",
			wantVersion:  "11",
		},
		{
			name:         "Go (always latest official version)",
			args:         args{runtime: typesLambda.RuntimeGo1x},
			wantLanguage: "go",
			wantVersion:  LatestLambdaGoVersion,
		},
		{
			name:         "Some new language not considered yet",
			args:         args{runtime: "SomeNewSupportedLanguage"},
			wantLanguage: "SomeNewSupportedLanguage",
			wantVersion:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLanguage, gotVersion := splitRuntime(tt.args.runtime)
			assert.Equalf(t, tt.wantLanguage, gotLanguage, "splitRuntime(%v)", tt.args.runtime)
			assert.Equalf(t, tt.wantVersion, gotVersion, "splitRuntime(%v)", tt.args.runtime)
		})
	}
}

func Test_useOfficialLanguageName(t *testing.T) {
	type args struct {
		l string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Go",
			args: args{"go"},
			want: constants.Go,
		},
		{
			name: "Java",
			args: args{"java"},
			want: constants.Java,
		},
		{
			name: "NodeJS",
			args: args{"nodejs"},
			want: constants.NodeJS,
		},
		{
			name: "Version not supported yet",
			args: args{"NotYetSupportedVersion"},
			want: "NotYetSupportedVersion",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, useOfficialLanguageName(tt.args.l), "useOfficialLanguageName(%v)", tt.args.l)
		})
	}
}

func Test_toRuntimeLanguage(t *testing.T) {
	type args struct {
		runtime typesLambda.Runtime
	}
	tests := []struct {
		name         string
		args         args
		wantLanguage string
	}{
		{
			name:         "Java from Java11",
			args:         args{"java11"},
			wantLanguage: "Java",
		},
		{
			name:         "New Language 42",
			args:         args{"New42"},
			wantLanguage: "New",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantLanguage, toRuntimeLanguage(tt.args.runtime), "toRuntimeLanguage(%v)", tt.args.runtime)
		})
	}
}

func Test_toRuntimeVersion(t *testing.T) {
	type args struct {
		runtime typesLambda.Runtime
	}
	tests := []struct {
		name        string
		args        args
		wantVersion string
	}{
		{
			name:        "Java from Java11",
			args:        args{"java11"},
			wantVersion: "11",
		},
		{
			name:        "New Language 42",
			args:        args{"New42"},
			wantVersion: "42",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantVersion, toRuntimeVersion(tt.args.runtime), "toRuntimeVersion(%v)", tt.args.runtime)
		})
	}
}

func Test_computeDiscovery_getAutomaticUpdates(t *testing.T) {
	type fields struct {
		virtualMachineAPI EC2API
		functionAPI       LambdaAPI
		systemManagerAPI  SSMAPI
		isDiscovering     bool
		awsConfig         *Client
		csID              string
	}
	type args struct {
		vm *typesEC2.Instance
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantAu *voc.AutomaticUpdates
	}{
		{
			name: "Happy path - all compliant",
			fields: fields{
				systemManagerAPI: mockSSMAPI{},
			},
			args: args{vm: &typesEC2.Instance{InstanceId: util.Ref(mockVM1ID)}},
			wantAu: &voc.AutomaticUpdates{
				Enabled:      true,
				SecurityOnly: true,
				Interval:     1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &computeDiscovery{
				virtualMachineAPI: tt.fields.virtualMachineAPI,
				functionAPI:       tt.fields.functionAPI,
				systemManagerAPI:  tt.fields.systemManagerAPI,
				isDiscovering:     tt.fields.isDiscovering,
				awsConfig:         tt.fields.awsConfig,
				csID:              tt.fields.csID,
			}
			assert.Equalf(t, tt.wantAu, d.getAutomaticUpdates(tt.args.vm), "getAutomaticUpdates(%v)", tt.args.vm)
		})
	}
}
