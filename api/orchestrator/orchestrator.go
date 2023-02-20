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
	reflect "reflect"

	"google.golang.org/protobuf/proto"
)

type CloudServiceHookFunc func(ctx context.Context, cld *CloudService, err error)
type TargetOfEvaluationHookFunc func(ctx context.Context, event *TargetOfEvaluationChangeEvent, err error)

// CloudServiceRequest represents any kind of RPC request, that contains a
// reference to a cloud service.
type CloudServiceRequest interface {
	GetCloudServiceId() string
	proto.Message
}

type LogRequest interface {
	GetPayloadID() string
	GetType() string
	CloudServiceRequest
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

func (req *CreateCatalogRequest) GetType() string {
	return reflect.TypeOf(req.Catalog).String()
}

func (req *CreateCatalogRequest) GetCloudServiceId() string {
	return ""
}

func (req *UpdateCatalogRequest) GetPayloadID() string {
	return req.GetCatalog().GetId()
}

func (req *UpdateCatalogRequest) GetType() string {
	return reflect.TypeOf(req.Catalog).String()
}

func (req *UpdateCatalogRequest) GetCloudServiceId() string {
	return ""
}

func (req *RemoveCatalogRequest) GetPayloadID() string {
	return req.GetCatalogId()
}

func (req *RemoveCatalogRequest) GetType() string {
	return reflect.TypeOf(Catalog{}).String()
}

func (req *RemoveCatalogRequest) GetCloudServiceId() string {
	return ""
}

func (req *CreateTargetOfEvaluationRequest) GetPayloadID() string {
	return ""
}

func (req *CreateTargetOfEvaluationRequest) GetType() string {
	return reflect.TypeOf(req.TargetOfEvaluation).String()
}

func (req *CreateTargetOfEvaluationRequest) GetCloudServiceId() string {
	return req.GetTargetOfEvaluation().GetCloudServiceId()
}

func (req *UpdateTargetOfEvaluationRequest) GetPayloadID() string {
	return ""
}

func (req *UpdateTargetOfEvaluationRequest) GetType() string {
	return reflect.TypeOf(req.TargetOfEvaluation).String()
}

func (req *UpdateTargetOfEvaluationRequest) GetCloudServiceId() string {
	return req.GetTargetOfEvaluation().GetCloudServiceId()
}

func (req *RemoveTargetOfEvaluationRequest) GetPayloadID() string {
	return ""
}

func (req *RemoveTargetOfEvaluationRequest) GetType() string {
	return reflect.TypeOf(TargetOfEvaluation{}).String()
}

func (req *CreateCertificateRequest) GetPayloadID() string {
	return req.GetCertificate().GetId()
}

func (req *CreateCertificateRequest) GetType() string {
	return reflect.TypeOf(req.Certificate).String()
}

func (req *CreateCertificateRequest) GetCloudServiceId() string {
	return req.GetCertificate().GetCloudServiceId()
}

func (req *UpdateCertificateRequest) GetPayloadID() string {
	return req.GetCertificate().GetId()
}

func (req *UpdateCertificateRequest) GetType() string {
	return reflect.TypeOf(req.Certificate).String()
}

func (req *UpdateCertificateRequest) GetCloudServiceId() string {
	return req.GetCertificate().GetCloudServiceId()
}

func (req *RemoveCertificateRequest) GetPayloadID() string {
	return req.GetCertificateId()
}

func (req *RemoveCertificateRequest) GetType() string {
	return reflect.TypeOf(Certificate{}).String()
}

func (req *RemoveCertificateRequest) GetCloudServiceId() string {
	return ""
}

func (req *CreateMetricRequest) GetPayloadID() string {
	return req.GetMetric().GetId()
}

func (req *CreateMetricRequest) GetType() string {
	return reflect.TypeOf(req.Metric).String()
}

func (req *CreateMetricRequest) GetCloudServiceId() string {
	return ""
}

func (req *UpdateMetricRequest) GetPayloadID() string {
	return req.GetMetric().GetId()
}

func (req *UpdateMetricRequest) GetType() string {
	return reflect.TypeOf(req.Metric).String()
}

func (req *UpdateMetricRequest) GetCloudServiceId() string {
	return ""
}

func (req *UpdateMetricImplementationRequest) GetPayloadID() string {
	return req.GetImplementation().GetMetricId()
}

func (req *UpdateMetricImplementationRequest) GetType() string {
	return reflect.TypeOf(req.Implementation).String()
}

func (req *UpdateMetricImplementationRequest) GetCloudServiceId() string {
	return ""
}

func (req *UpdateMetricConfigurationRequest) GetPayloadID() string {
	return req.GetConfiguration().GetMetricId()
}

func (req *UpdateMetricConfigurationRequest) GetType() string {
	return reflect.TypeOf(req.Configuration).String()
}

func (req *RegisterCloudServiceRequest) GetPayloadID() string {
	return req.GetCloudService().GetId()
}

func (req *RegisterCloudServiceRequest) GetType() string {
	return reflect.TypeOf(req.CloudService).String()
}

func (req *RegisterCloudServiceRequest) GetCloudServiceId() string {
	return req.GetCloudService().GetId()
}

func (req *UpdateCloudServiceRequest) GetPayloadID() string {
	return req.GetCloudService().GetId()
}

func (req *UpdateCloudServiceRequest) GetType() string {
	return reflect.TypeOf(req.CloudService).String()
}

func (req *RemoveCloudServiceRequest) GetPayloadID() string {
	return req.GetCloudServiceId()
}

func (req *RemoveCloudServiceRequest) GetType() string {
	return reflect.TypeOf(CloudService{}).String()
}

// TableName overrides the table name used by ControlInScope to `controls_in_scope`
func (*ControlInScope) TableName() string {
	return "controls_in_scope"
}
