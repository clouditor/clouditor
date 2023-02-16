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

package logging

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

const (
	LoglevelDebug = "debug"
	LoglevelInfo  = "info"
	LoglevelError = "error"
)

// LogCreateMessage creates a 'create' logging message with the given parameters
//   - log *logrus.Entry
//   - typ is the type of the created object
//   - id is the ID of the created object
//   - loglevel the message must have (LoglevelDebug, LoglevelInfo, LoglevelError)
//   - Optional. cloudServiceId is the ID of the Cloud Service.
//   - Optional. catalogId is the ID of the used Catalog.
//
// The message looks like one of the following
//   - "*orchestrator.Catalog created with ID 'Cat1234'."
//   - "*orchestrator.Certificate created with ID 'Cert1234' for Cloud Service '00000000-0000-0000-0000-000000000000'."
//   - "*orchestrator.TargetOfEvaluation created with ID 'ToE1234' for Cloud Service '00000000-0000-0000-0000-000000000000' and Catalog 'EUCS'."
func LogCreateMessage(log *logrus.Entry, typ, id, loglevel, cloudServiceId, catalogId string) {
	var (
		message string
	)

	message = fmt.Sprintf("%s created with ID '%s'", typ, id)

	if cloudServiceId != "" && catalogId != "" {
		message = fmt.Sprintf("%s for Cloud Service '%s' and Catalog '%s'.", message, cloudServiceId, catalogId)
	} else if cloudServiceId != "" {
		message = fmt.Sprintf("%s for Cloud Service '%s'", message, cloudServiceId)
	}

	switch loglevel {
	case LoglevelDebug:
		log.Debugf("%s.", message)
	case LoglevelInfo:
		log.Infof("%s.", message)
	case LoglevelError:
		log.Errorf("%s.", message)
	}
}

// DebugUpdate creates a debug logging message for update
func DebugUpdate(t, id, name string) string {

	return ""
}

// DebugRemove creates a debug logging message for remove
func DebugRemove(t, id, name string) string {

	return ""
}

// DebugUpdate creates a debug logging message for load
func DebugLoad(t, id, name string) string {

	return ""
}

// log.Debugf("Catalog created with name '%s'.", req.Catalog.GetName())
// log.Debugf("Metric created with name '%s'.", req.Metric.GetName())
// log.Debugf("Metric created with name '%s'.", req.Metric.GetName())
// log.Debugf("Certificate created for Cloud Service ID '%s'.", req.Certificate.GetCloudServiceId())
// log.Debugf("ToE created for Cloud Service ID '%s' and Catalog ID '%s'.", req.TargetOfEvaluation.GetCloudServiceId(), req.TargetOfEvaluation.GetCatalogId())

// log.Debugf("Catalog updated with name '%s' and id '%s'.", req.Catalog.GetName(), req.Catalog.GetId())
// log.Debugf("Cloud Service updated with name '%s' and id '%s'.", req.CloudService.GetName(), req.CloudService.GetId())
// log.Debugf("Metric updated with id '%s'.", req.Metric.GetId())
// log.Debugf("Metric implemenatation updated for metric id '%s'.", req.Implementation.GetMetricId())
// log.Debugf("Metric conifguration updated for metric id '%s' and Cloud Service ID '%s'.", req.GetMetricId(), req.GetCloudServiceId())
// log.Debugf("Certificate updated with id '%s' for Cloud Service ID '%s'.", req.Certificate.GetId(), req.Certificate.GetCloudServiceId())
// log.Debugf("ToE updated for Cloud Service ID '%s' and Catalog ID '%s'.", req.TargetOfEvaluation.GetCloudServiceId(), req.TargetOfEvaluation.GetCatalogId())

// log.Debugf("Catalog removed with id '%s'.", req.GetCatalogId())
// log.Debugf("Cloud Service removed with id '%s'.", req.GetCloudServiceId())
// log.Debugf("Certificate removed with id '%s'.", req.GetCertificateId())
// log.Debugf("ToE removed for Cloud Service ID '%s' and Catalog ID '%s'.", req.GetCloudServiceId(), req.GetCatalogId())

// log.Debugf("Catalog loaded with name '%s' and id '%s'.", catalogs[0].GetName(), catalogs[0].GetId())

// log.Debugf("Cloud Service registered with name '%s'.", req.CloudService.GetName())
