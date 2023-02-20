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

	"google.golang.org/protobuf/proto"
)

type CloudServiceHookFunc func(ctx context.Context, cld *CloudService, err error)
type TargetOfEvaluationHookFunc func(ctx context.Context, event *TargetOfEvaluationChangeEvent, err error)

// CloudServiceRequest represents any kind of RPC request, that contains a
// reference to a cloud service.
//
// Note: GetCloudServiceId() is already implemented by the generated protobuf
// code for the following messages because they directly have a cloud_service id
// field:
//   - RemoveControlFromScopeRequest
//   - ListControlsInScopeRequest
//   - GetCloudServiceRequest
//   - RemoveCloudServiceRequest
//   - UpdateMetricConfigurationRequest
//   - GetMetricConfigurationRequest
//   - ListMetricConfigurationRequest
//   - MetricChangeEvent
//   - TargetOfEvaluation
//   - RemoveTargetOfEvaluationRequest
//   - GetTargetOfEvaluationRequest
//   - ListTargetsOfEvaluationRequest
//   - Certificate
//
// All other requests, especially in cases where the cloud service ID is
// embedded in a sub-field need to explicitly implement this interface in order.
// This interface is for example used by authorization checks.
type CloudServiceRequest interface {
	GetCloudServiceId() string
	proto.Message
}

type PayloadRequest interface {
	GetPayload() proto.Message
}

// GetCloudServiceId is a shortcut to implement CloudServiceRequest. It returns
// the cloud service ID of the inner object.
func (req *AddControlToScopeRequest) GetCloudServiceId() string {
	return req.Scope.GetTargetOfEvaluationCloudServiceId()
}

// GetCloudServiceId is a shortcut to implement CloudServiceRequest. It returns
// the cloud service ID of the inner object.
func (req *UpdateControlInScopeRequest) GetCloudServiceId() string {
	return req.Scope.GetTargetOfEvaluationCloudServiceId()
}

// GetCloudServiceId is a shortcut to implement CloudServiceRequest. It returns
// the cloud service ID of the inner object.
func (req *UpdateCloudServiceRequest) GetCloudServiceId() string {
	return req.CloudService.GetId()
}

// GetCloudServiceId is a shortcut to implement CloudServiceRequest. It returns
// the cloud service ID of the inner object.
func (req *StoreAssessmentResultRequest) GetCloudServiceId() string {
	return req.GetResult().GetCloudServiceId()
}

// GetCloudServiceId is a shortcut to implement CloudServiceRequest. It returns
// the cloud service ID of the inner object.
func (req *CreateTargetOfEvaluationRequest) GetCloudServiceId() string {
	return req.GetTargetOfEvaluation().GetCloudServiceId()
}

// GetCloudServiceId is a shortcut to implement CloudServiceRequest. It returns
// the cloud service ID of the inner object.
func (req *UpdateTargetOfEvaluationRequest) GetCloudServiceId() string {
	return req.GetTargetOfEvaluation().GetCloudServiceId()
}

// GetCloudServiceId is a shortcut to implement CloudServiceRequest. It returns
// the cloud service ID of the inner object.
func (req *CreateCertificateRequest) GetCloudServiceId() string {
	return req.GetCertificate().GetCloudServiceId()
}

// GetCloudServiceId is a shortcut to implement CloudServiceRequest. It returns
// the cloud service ID of the inner object.
func (req *UpdateCertificateRequest) GetCloudServiceId() string {
	return req.GetCertificate().GetCloudServiceId()
}

// GetCloudServiceId is a shortcut to implement CloudServiceRequest. It returns
// the cloud service ID of the inner object.
func (req *RegisterCloudServiceRequest) GetCloudServiceId() string {
	return req.GetCloudService().GetId()
}

func (req *CreateCatalogRequest) GetPayload() proto.Message {
	return req.Catalog
}

func (req *UpdateCatalogRequest) GetPayload() proto.Message {
	return req.Catalog
}

func (req *RemoveCatalogRequest) GetPayload() proto.Message {
	return &Catalog{Id: req.CatalogId}
}

func (req *CreateTargetOfEvaluationRequest) GetPayload() proto.Message {
	return req.TargetOfEvaluation
}

func (req *UpdateTargetOfEvaluationRequest) GetPayload() proto.Message {
	return req.TargetOfEvaluation
}

func (req *RemoveTargetOfEvaluationRequest) GetPayload() proto.Message {
	return &TargetOfEvaluation{CloudServiceId: req.CloudServiceId, CatalogId: req.CatalogId}
}

func (req *CreateCertificateRequest) GetPayload() proto.Message {
	return req.Certificate
}

func (req *UpdateCertificateRequest) GetPayload() proto.Message {
	return req.Certificate
}

func (req *RemoveCertificateRequest) GetPayload() proto.Message {
	return &Certificate{Id: req.CertificateId}
}

func (req *CreateMetricRequest) GetPayload() proto.Message {
	return req.Metric
}

func (req *UpdateMetricRequest) GetPayload() proto.Message {
	return req.Metric
}

func (req *UpdateMetricImplementationRequest) GetPayload() proto.Message {
	return req.Implementation
}

func (req *UpdateMetricConfigurationRequest) GetPayload() proto.Message {
	return req.Configuration
}

func (req *RegisterCloudServiceRequest) GetPayload() proto.Message {
	return req.CloudService
}

func (req *UpdateCloudServiceRequest) GetPayload() proto.Message {
	return req.CloudService
}

func (req *RemoveCloudServiceRequest) GetPayload() proto.Message {
	return &CloudService{Id: req.CloudServiceId}
}

// TableName overrides the table name used by ControlInScope to `controls_in_scope`
func (*ControlInScope) TableName() string {
	return "controls_in_scope"
}
