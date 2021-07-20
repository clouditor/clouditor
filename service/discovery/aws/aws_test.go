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
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

const mockRegion = "mockRegion"

// TestNewClient tests the NewClient function
func TestNewClient(t *testing.T) {
	// Mock loadDefaultConfig and store the original function back to loadDefaultConfig at the end of the test
	old := loadDefaultConfig
	defer func() { loadDefaultConfig = old }()

	// Case 1: Get config (and no error)
	loadDefaultConfig = func(ctx context.Context,
		opt ...func(options *config.LoadOptions) error) (cfg aws.Config, err error) {
		err = nil
		cfg = aws.Config{
			Region:           mockRegion,
			Credentials:      nil,
			HTTPClient:       nil,
			EndpointResolver: nil,
			Retryer:          nil,
			ConfigSources:    nil,
			APIOptions:       nil,
			Logger:           nil,
			ClientLogMode:    0,
		}
		return
	}
	client, err := NewClient()
	assert.Nil(t, err)
	assert.Equal(t, mockRegion, client.Cfg.Region)

	// Case 1: Get error (and empty config)
	loadDefaultConfig = func(ctx context.Context,
		opt ...func(options *config.LoadOptions) error) (cfg aws.Config, err error) {
		err = errors.New("error occurred while loading credentials")
		cfg = aws.Config{}
		return
	}
	client, err = NewClient()
	assert.NotNil(t, err)
	assert.Empty(t, client.Cfg.Region)

}
