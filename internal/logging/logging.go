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

package logging

import (
	"fmt"
	"strings"

	"clouditor.io/clouditor/internal/api"
	"clouditor.io/clouditor/internal/util"
	"github.com/sirupsen/logrus"
)

// RequestType specifies the type of request
type RequestType int

const (
	Assess RequestType = iota
	Add
	Create
	Register
	Remove
	Store
	Send
	Update
)

// String returns the RequestType as string.
func (r RequestType) String() string {
	switch r {
	case Assess:
		return "assessed"
	case Add:
		return "added"
	case Create:
		return "created"
	case Register:
		return "registered"
	case Remove:
		return "removed"
	case Store:
		return "stored"
	case Send:
		return "sent"
	case Update:
		return "updated"
	default:
		return "unspecified"
	}
}

// LogRequest creates a logging message with the given parameters
//   - log *logrus.Entry
//   - level the message must have
//   - reqType is the request type
//   - Optional. params for self-created string messages to be appended to the created log message. The elements do not need a space at the beginning of the message.
//
// The message looks like one of the following depending on the given information
//
//	"*orchestrator.Catalog created with ID 'Cat1234'."
//
// or
//
//	"*orchestrator.Certificate created with ID 'Cert1234' for Cloud Service '00000000-0000-0000-0000-000000000000'."
//
// or
//
//	"*orchestrator.TargetOfEvaluation created with ID 'ToE1234' for Cloud Service '00000000-0000-0000-0000-000000000000' and Catalog 'EUCS'."
func LogRequest(log *logrus.Entry, level logrus.Level, reqType RequestType, req api.PayloadRequest, params ...string) {
	var (
		message string
	)

	// Check if inputs are available
	if log == nil || util.IsNil(req) {
		return
	}

	// Retrieve the payload from the request. The request itself is usually
	// a wrapper around the sent object.
	payload := req.GetPayload()
	if util.IsNil(payload) {
		return
	}

	// We can retrieve the name via the proto descriptor. This should be
	// sufficiently fast and also gives us the non-pointer type in comparison to
	// the %T printf directive.
	name := req.GetPayload().ProtoReflect().Descriptor().Name()

	// Check, if our payload has an ID field
	idreq, ok := payload.(interface{ GetId() string })
	if ok && idreq.GetId() != "" {
		message = fmt.Sprintf("%s with ID '%s' %s", name, idreq.GetId(), reqType.String())
	} else {
		message = fmt.Sprintf("%s %s", name, reqType.String())
	}

	// Check, if it is a cloud service request. In this case we can append the
	// information about the target cloud service. However, we only want to do
	// that, if the payload type is not a cloud service itself.
	csreq, ok := req.(api.CloudServiceRequest)
	// If params is not empty, the elements are joined and added to the message
	if name != "CloudService" && ok && len(params) > 0 {
		message = fmt.Sprintf("%s for Cloud Service '%s' %s", message, csreq.GetCloudServiceId(), strings.Join(params, " "))
	} else if name != "CloudService" && ok {
		message = fmt.Sprintf("%s for Cloud Service '%s'", message, csreq.GetCloudServiceId())
	}

	log.Logf(level, "%s.", message)
}
