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
	"strings"

	"clouditor.io/clouditor/api/request"
	"github.com/sirupsen/logrus"
)

// RequestType specifies the type of request
type RequestType int

const (
	Create RequestType = iota
	Remove
	Update
	Register
	Store
	Assess
	Send
)

// String returns the RequestType as string.
func (r RequestType) String() string {
	switch r {
	case Create:
		return "created"
	case Remove:
		return "removed"
	case Update:
		return "updated"
	case Register:
		return "registered"
	case Store:
		return "stored"
	case Assess:
		return "assessed"
	case Send:
		return "sent"
	default:
		return "unspecified"
	}
}

// LogRequest creates a logging message with the given parameters
//   - log *logrus.Entry
//   - loglevel the message must have
//   - reqType is the request type
//   - Optional. params for self-created string messages to be appended to the created log message. The elements do not need a space at the beginning of the message.
//
// The message looks like one of the following depending on the given information
//   - "*orchestrator.Catalog created with ID 'Cat1234'."
//   - "*orchestrator.Certificate created with ID 'Cert1234' for Cloud Service '00000000-0000-0000-0000-000000000000'."
//   - "*orchestrator.TargetOfEvaluation created with ID 'ToE1234' for Cloud Service '00000000-0000-0000-0000-000000000000' and Catalog 'EUCS'."
func LogRequest(log *logrus.Entry, level logrus.Level, reqType RequestType, req request.PayloadRequest, params ...string) {
	var (
		message string
	)

	// Retrieve the payload from the request. The request itself is usually
	// a wrapper around the sent object.
	payload := req.GetPayload()
	if payload == nil {
		return
	}

	// Check, if our payload has an ID field
	idreq, ok := payload.(interface{ GetId() string })
	// If the request type is "Send" then the message must look different
	if reqType == Send {
		if ok && len(params) > 0 {
			message = fmt.Sprintf("%T with ID '%s' %s to %s", req.GetPayload(), idreq.GetId(), reqType.String(), strings.Join(params, " "))
		} else if !ok && len(params) > 0 {
			message = fmt.Sprintf("%T %s to %s", req.GetPayload(), reqType.String(), strings.Join(params, " "))
		} else {
			message = fmt.Sprintf("%T %s", req.GetPayload(), reqType.String())
		}
	} else {
		if ok {
			message = fmt.Sprintf("%T %s with ID '%s'", req.GetPayload(), reqType.String(), idreq.GetId())
		} else {
			message = fmt.Sprintf("%T %s", req.GetPayload(), reqType.String())
		}
	}

	// Check, if it is a cloud service request. In this case we can append the
	// information about the target cloud service
	csreq, ok := req.(request.CloudServiceRequest)
	// If params is not empty, the elements are joined and added to the message
	if ok && len(params) > 0 {
		message = fmt.Sprintf("%s for Cloud Service '%s' %s", message, csreq.GetCloudServiceId(), strings.Join(params, " "))
	} else if ok {
		message = fmt.Sprintf("%s for Cloud Service '%s'", message, csreq.GetCloudServiceId())
	}

	log.Logf(level, "%s.", message)
}
