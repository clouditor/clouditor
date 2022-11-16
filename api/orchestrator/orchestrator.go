// Copyright 2022 Fraunhofer AISEC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//           $$\                           $$\ $$\   $$\
//           $$ |                          $$ |\__|  $$ |
//  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
// $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
// $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
// $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
// \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
//  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
//
// This file is part of Clouditor Community Edition.

package orchestrator

import (
	"context"
	"errors"

	"clouditor.io/clouditor/api/assessment"
)

type CloudServiceHookFunc func(ctx context.Context, cld *CloudService, err error)
type ToeHookFunc func(ctx context.Context, toe *TargetOfEvaluation, status ToeStatus_Status, err error)

var (
	ErrCertificateIsNil   = errors.New("certificate is empty")
	ErrServiceIsNil       = errors.New("service is empty")
	ErrNameIsMissing      = errors.New("service name is empty")
	ErrIDIsMissing        = errors.New("service ID is empty")
	ErrCertIDIsMissing    = errors.New("certificate ID is empty")
	ErrCatalogIsNil       = errors.New("catalog is empty")
	ErrCatalogIDIsMissing = errors.New("catalog ID is empty")
	ErrToEIDIsMissing     = errors.New("toe ID is empty")
)

// Validate validates the UpdateMetricConfigurationRequest
func (req *UpdateMetricConfigurationRequest) Validate() error {
	// Check cloud service ID
	err := assessment.CheckCloudServiceID(req.CloudServiceId)
	if err != nil {
		return err
	}

	// Check metric ID
	if req.MetricId == "" {
		return assessment.ErrMetricIdMissing
	}

	return nil
}

// Validate validates the GetMetricConfigurationRequest
func (req *GetMetricConfigurationRequest) Validate() error {
	// Check cloud service ID
	err := assessment.CheckCloudServiceID(req.CloudServiceId)
	if err != nil {
		return err
	}

	// Check metric ID
	if req.MetricId == "" {
		return assessment.ErrMetricIdMissing
	}

	return nil
}

type CloudServiceRequest interface {
	GetCloudServiceId() string
}
