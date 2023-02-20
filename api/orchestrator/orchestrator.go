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

type T any

type LogRequest interface {
	GetPayloadID() string
	GetPayload() T
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

func (req *CreateCatalogRequest) GetPayloadID() string {
	return req.GetCatalog().GetId()
}

func (req *CreateCatalogRequest) GetPayload() T {
	return req.Catalog
}

func (req *UpdateCatalogRequest) GetPayloadID() string {
	return req.GetCatalog().GetId()
}

func (req *UpdateCatalogRequest) GetPayload() T {
	return req.Catalog
}

func (req *RemoveCatalogRequest) GetPayloadID() string {
	return req.GetCatalogId()
}

func (req *RemoveCatalogRequest) GetPayload() T {
	return Catalog{}
}

func (req *CreateTargetOfEvaluationRequest) GetPayloadID() string {
	return ""
}

func (req *CreateTargetOfEvaluationRequest) GetPayload() T {
	return req.TargetOfEvaluation
}

func (req *CreateTargetOfEvaluationRequest) GetCloudServiceId() string {
	return req.GetTargetOfEvaluation().GetCloudServiceId()
}

func (req *UpdateTargetOfEvaluationRequest) GetPayloadID() string {
	return ""
}

func (req *UpdateTargetOfEvaluationRequest) GetPayload() T {
	return req.TargetOfEvaluation
}

func (req *UpdateTargetOfEvaluationRequest) GetCloudServiceId() string {
	return req.GetTargetOfEvaluation().GetCloudServiceId()
}

func (req *RemoveTargetOfEvaluationRequest) GetPayloadID() string {
	return ""
}

func (req *RemoveTargetOfEvaluationRequest) GetPayload() T {
	return TargetOfEvaluation{}
}

func (req *CreateCertificateRequest) GetPayloadID() string {
	return req.GetCertificate().GetId()
}

func (req *CreateCertificateRequest) GetPayload() T {
	return req.Certificate
}

func (req *CreateCertificateRequest) GetCloudServiceId() string {
	return req.GetCertificate().GetCloudServiceId()
}

func (req *UpdateCertificateRequest) GetPayloadID() string {
	return req.GetCertificate().GetId()
}

func (req *UpdateCertificateRequest) GetPayload() T {
	return req.Certificate
}

func (req *UpdateCertificateRequest) GetCloudServiceId() string {
	return req.GetCertificate().GetCloudServiceId()
}

func (req *RemoveCertificateRequest) GetPayloadID() string {
	return req.GetCertificateId()
}

func (req *RemoveCertificateRequest) GetPayload() T {
	return Certificate{}
}

func (req *CreateMetricRequest) GetPayloadID() string {
	return req.GetMetric().GetId()
}

func (req *CreateMetricRequest) GetPayload() T {
	return req.Metric
}

func (req *UpdateMetricRequest) GetPayloadID() string {
	return req.GetMetric().GetId()
}

func (req *UpdateMetricRequest) GetPayload() T {
	return req.Metric
}

func (req *UpdateMetricImplementationRequest) GetPayloadID() string {
	return req.GetImplementation().GetMetricId()
}

func (req *UpdateMetricImplementationRequest) GetPayload() T {
	return req.Implementation
}

func (req *UpdateMetricConfigurationRequest) GetPayloadID() string {
	return req.GetConfiguration().GetMetricId()
}

func (req *UpdateMetricConfigurationRequest) GetPayload() T {
	return req.Configuration
}

func (req *RegisterCloudServiceRequest) GetPayloadID() string {
	return req.GetCloudService().GetId()
}

func (req *RegisterCloudServiceRequest) GetPayload() T {
	return req.CloudService
}

func (req *RegisterCloudServiceRequest) GetCloudServiceId() string {
	return req.GetCloudService().GetId()
}

func (req *UpdateCloudServiceRequest) GetPayloadID() string {
	return req.GetCloudService().GetId()
}

func (req *UpdateCloudServiceRequest) GetPayload() T {
	return req.CloudService
}

func (req *RemoveCloudServiceRequest) GetPayloadID() string {
	return req.GetCloudServiceId()
}

func (req *RemoveCloudServiceRequest) GetPayload() T {
	return CloudService{}
}

// TableName overrides the table name used by ControlInScope to `controls_in_scope`
func (*ControlInScope) TableName() string {
	return "controls_in_scope"
}
