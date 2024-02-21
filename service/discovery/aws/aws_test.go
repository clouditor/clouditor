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
	"testing"

	"clouditor.io/clouditor/v2/internal/testutil/assert"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go"
)

const mockRegion = "mockRegion"

// TestNewClient tests the NewClient function
func TestNewClient(t *testing.T) {
	// Mock loadDefaultConfig and store the original function back at the end of the test
	oldLoadDefaultConfig := loadDefaultConfig
	defer func() { loadDefaultConfig = oldLoadDefaultConfig }()
	// Mock newFromConfigSTS and store the original function back at the ned of the test
	oldNewFromConfigSTS := newFromConfigSTS
	defer func() { newFromConfigSTS = oldNewFromConfigSTS }()

	// Case 1: Get config (and no error)
	loadDefaultConfig = func(ctx context.Context,
		opt ...func(options *config.LoadOptions) error) (cfg aws.Config, err error) {
		err = nil
		cfg = aws.Config{
			Region: mockRegion,
		}
		return
	}
	newFromConfigSTS = func(cfg aws.Config) STSAPI {
		return mockSTSClient{}
	}

	client, err := NewClient()
	assert.NoError(t, err)
	assert.Equal(t, mockRegion, client.cfg.Region)

	// Case 2: Get error while loading credentials
	loadDefaultConfig = func(ctx context.Context,
		opt ...func(options *config.LoadOptions) error) (cfg aws.Config, err error) {
		err = errors.New("error occurred while loading credentials")
		cfg = aws.Config{}
		return
	}
	client, err = NewClient()
	assert.Error(t, err)
	assert.Nil(t, client)

	// Case 3: Get error while calling GetCallerIdentity
	newFromConfigSTS = func(cfg aws.Config) STSAPI {
		return mockSTSClientWithAPIError{}
	}
	loadDefaultConfig = func(ctx context.Context,
		opt ...func(options *config.LoadOptions) error) (cfg aws.Config, err error) {
		err = nil
		cfg = aws.Config{
			Region: mockRegion,
		}
		return
	}
	client, err = NewClient()
	assert.Error(t, err)
	assert.Nil(t, client)

}

type mockSTSClient struct{}

func (mockSTSClient) GetCallerIdentity(_ context.Context,
	_ *sts.GetCallerIdentityInput, _ ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
	return &sts.GetCallerIdentityOutput{
		Account: aws.String("12345"),
	}, nil
}

type mockSTSClientWithAPIError struct{}

func (mockSTSClientWithAPIError) GetCallerIdentity(_ context.Context,
	_ *sts.GetCallerIdentityInput, _ ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
	return nil, &smithy.OperationError{
		ServiceID:     "STS",
		OperationName: "GetCallerIdentity",
		Err:           errors.New("MaxAttemptsError"),
	}
}
