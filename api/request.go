// Copyright 2023 Fraunhofer AISEC
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

package api

import (
	"clouditor.io/clouditor/v2/internal/api"
)

// PayloadRequest describes any kind of requests that carries a certain payload.
// This is for example a Create/Update request carrying an embedded message,
// which should be updated or created.
type PayloadRequest = api.PayloadRequest

// CertificationTargetRequest represents any kind of RPC request, that contains a
// reference to a certification target.
//
// Note: GetCertificationTargetId() is already implemented by the generated protobuf
// code for the following messages because they directly have a certification_target id
// field:
//   - orchestrator.RemoveControlFromScopeRequest
//   - orchestrator.ListControlsInScopeRequest
//   - orchestrator.GetCertificationTargetRequest
//   - orchestrator.RemoveCertificationTargetRequest
//   - orchestrator.UpdateMetricConfigurationRequest
//   - orchestrator.GetMetricConfigurationRequest
//   - orchestrator.ListMetricConfigurationRequest
//   - orchestrator.MetricChangeEvent
//   - orchestrator.AuditScope
//   - orchestrator.RemoveAuditScopeRequest
//   - orchestrator.GetAuditScopeRequest
//   - orchestrator.ListAuditScopesRequest
//   - orchestrator.Certificate
//
// All other requests, especially in cases where the certification target ID is
// embedded in a sub-field need to explicitly implement this interface in order.
// This interface is for example used by authorization checks.
type CertificationTargetRequest = api.CertificationTargetRequest
