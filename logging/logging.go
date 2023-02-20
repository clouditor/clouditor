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

	"clouditor.io/clouditor/api/orchestrator"
	"github.com/sirupsen/logrus"
)

const (
	Create   = "created"
	Remove   = "removed"
	Update   = "updated"
	Register = "registered"
	Load     = "loaded"
)

// LogMessage creates a logging message with the given parameters
//   - log *logrus.Entry
//   - loglevel the message must have (LoglevelDebug, LoglevelInfo, LoglevelError)
//   - operation that is performed (Create, Remove, Update, Register, Load)
//   - Optional. params must contain the catalog ID.
//
// The message looks like one of the following depending on the given information
//   - "*orchestrator.Catalog created with ID 'Cat1234'."
//   - "*orchestrator.Certificate created with ID 'Cert1234' for Cloud Service '00000000-0000-0000-0000-000000000000'."
//   - "*orchestrator.TargetOfEvaluation created with ID 'ToE1234' for Cloud Service '00000000-0000-0000-0000-000000000000' and Catalog 'EUCS'."
func LogMessage(log *logrus.Entry, loglevel logrus.Level, operation string, req orchestrator.LogRequest, params ...string) {
	var (
		message string
	)

	if req.GetPayloadID() != "" {
		message = fmt.Sprintf("%s %s with ID '%s'", req.GetType(), operation, req.GetPayloadID())
	} else {
		message = fmt.Sprintf("%s %s", req.GetType(), operation)
	}

	if req.GetCloudServiceId() != "" && len(params) > 0 {
		message = fmt.Sprintf("%s for Cloud Service '%s' and Catalog '%s'", message, req.GetCloudServiceId(), params[0])
	} else if req.GetCloudServiceId() != "" {
		message = fmt.Sprintf("%s for Cloud Service '%s'", message, req.GetCloudServiceId())
	}

	switch loglevel {
	case logrus.DebugLevel:
		log.Debugf("%s.", message)
	case logrus.InfoLevel:
		log.Infof("%s.", message)
	case logrus.ErrorLevel:
		log.Errorf("%s.", message)
	}
}
