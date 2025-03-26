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
	"slices"

	"clouditor.io/clouditor/v2/api/assessment"

	"google.golang.org/protobuf/proto"
)

type CertificationTargetHookFunc func(ctx context.Context, cld *TargetOfEvaluation, err error)
type AuditScopeHookFunc func(ctx context.Context, event *AuditScopeChangeEvent, err error)

// GetCertificationTargetId is a shortcut to implement CertificationTargetRequest. It returns
// the certification target ID of the inner object.
func (req *StoreAssessmentResultRequest) GetCertificationTargetId() string {
	return req.GetResult().GetCertificationTargetId()
}

// GetCertificationTargetId is a shortcut to implement CertificationTargetRequest. It returns
// the certification target ID of the inner object.
func (req *CreateCertificateRequest) GetCertificationTargetId() string {
	return req.GetCertificate().GetCertificationTargetId()
}

// GetCertificationTargetId is a shortcut to implement CertificationTargetRequest. It returns
// the certification target ID of the inner object.
func (req *UpdateCertificateRequest) GetCertificationTargetId() string {
	return req.GetCertificate().GetCertificationTargetId()
}

// GetCertificationTargetId is a shortcut to implement CertificationTargetRequest. It returns
// the certification target ID of the inner object.
func (req *CreateCertificationTargetRequest) GetCertificationTargetId() string {
	return req.GetCertificationTarget().GetId()
}

// GetCertificationTargetId is a shortcut to implement CertificationTargetRequest. It returns
// the certification target ID of the inner object.
func (req *UpdateCertificationTargetRequest) GetCertificationTargetId() string {
	return req.TargetOfEvaluation.GetId()
}

// GetCertificationTargetId is a shortcut to implement CertificationTargetRequest. It returns
// the certification target ID of the inner object.
func (req *CreateAuditScopeRequest) GetCertificationTargetId() string {
	return req.GetAuditScope().GetCertificationTargetId()
}

// GetCertificationTargetId is a shortcut to implement CertificationTargetRequest. It returns
// the certification target ID of the inner object.
func (req *UpdateAuditScopeRequest) GetCertificationTargetId() string {
	return req.GetAuditScope().GetCertificationTargetId()
}

func (req *StoreAssessmentResultRequest) GetPayload() proto.Message {
	return req.Result
}

func (req *RegisterAssessmentToolRequest) GetPayload() proto.Message {
	return req.Tool
}

func (req *UpdateAssessmentToolRequest) GetPayload() proto.Message {
	return req.Tool
}

func (req *DeregisterAssessmentToolRequest) GetPayload() proto.Message {
	return &AssessmentTool{Id: req.ToolId}
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

func (req *CreateCertificateRequest) GetPayload() proto.Message {
	return req.Certificate
}

func (req *UpdateCertificateRequest) GetPayload() proto.Message {
	return req.Certificate
}

func (req *RemoveCertificateRequest) GetPayload() proto.Message {
	return &Certificate{Id: req.CertificateId}
}

func (req *CreateCertificationTargetRequest) GetPayload() proto.Message {
	return req.TargetOfEvaluation
}

func (req *UpdateCertificationTargetRequest) GetPayload() proto.Message {
	return req.TargetOfEvaluation
}

func (req *RemoveCertificationTargetRequest) GetPayload() proto.Message {
	return &TargetOfEvaluation{Id: req.CertificationTargetId}
}

func (req *CreateMetricRequest) GetPayload() proto.Message {
	return req.Metric
}

func (req *UpdateMetricRequest) GetPayload() proto.Message {
	return req.Metric
}

func (req *RemoveMetricRequest) GetPayload() proto.Message {
	return &assessment.Metric{Id: req.MetricId}
}

func (req *UpdateMetricConfigurationRequest) GetPayload() proto.Message {
	return req.Configuration
}

func (req *UpdateMetricImplementationRequest) GetPayload() proto.Message {
	return req.Implementation
}

func (req *CreateAuditScopeRequest) GetPayload() proto.Message {
	return req.AuditScope
}

func (req *UpdateAuditScopeRequest) GetPayload() proto.Message {
	return req.AuditScope
}

func (req *RemoveAuditScopeRequest) GetPayload() proto.Message {
	return &AuditScope{Id: req.AuditScopeId}
}

// IsRelevantFor checks, whether this control is relevant for the given audit scope. For now this mainly
// checks, whether the assurance level matches, if the Audit Scope has one. In the future, this could also include checks, if
// the control is somehow out of scope.
func (c *Control) IsRelevantFor(auditScope *AuditScope, catalog *Catalog) bool {
	// If the catalog does not have an assurance level, we are good to go
	if len(catalog.AssuranceLevels) == 0 {
		return true
	}

	// If the control does not explicitly specify an assurance level, we are also ok
	if c.AssuranceLevel == nil || auditScope.AssuranceLevel == nil {
		return true
	}

	// Otherwise, we need to retrieve the possible assurance levels (in order) from the catalogs and compare the
	// indices
	idxControl := slices.Index(catalog.AssuranceLevels, *c.AssuranceLevel)
	idxToe := slices.Index(catalog.AssuranceLevels, *auditScope.AssuranceLevel)

	return idxControl <= idxToe
}
