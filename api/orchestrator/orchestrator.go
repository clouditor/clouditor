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
	"errors"
)

var (
	ErrRequestIsNil  = errors.New("request is empty")
	ErrServiceIsNil  = errors.New("service is empty")
	ErrNameIsMissing = errors.New("service name is empty")
	ErrIDIsMissing   = errors.New("service ID is empty")
)

// TODO(lebogg for oxisto): I kept that req was checked to be not nil - but it is necessary?

// Validate validates the request for RPC RegisterCloudService
func (s *RegisterCloudServiceRequest) Validate() (err error) {
	if s == nil {
		return ErrRequestIsNil
	}
	if s.Service == nil {
		return ErrServiceIsNil
	}
	if s.Service.Name == "" {
		return ErrNameIsMissing
	}
	return
}

// Validate validates the request for RPC GetCloudService
func (s *GetCloudServiceRequest) Validate() (err error) {
	if s == nil {
		return ErrRequestIsNil
	}
	if s.ServiceId == "" {
		return ErrIDIsMissing
	}
	return
}

// Validate validates the request for RPC UpdateCloudService
func (s *UpdateCloudServiceRequest) Validate() (err error) {
	if s == nil {
		return ErrRequestIsNil
	}
	// TODO(all): See TODO in cloud_service.go -> Remove ServiceID from req?
	// If not I will differentiate both Error messages
	if s.Service.Id == "" {
		return ErrIDIsMissing
	}
	// TODO(all): Otherwise, name will be overwritten with empty string -> See comment in Update in db.go: Save vs Update
	if s.Service.Name == "" {
		return ErrNameIsMissing
	}
	if s.ServiceId == "" {
		return ErrIDIsMissing
	}
	return
}

// Validate validates the request for RPC GetCloudService
func (s *RemoveCloudServiceRequest) Validate() (err error) {
	if s == nil {
		return ErrRequestIsNil
	}
	if s.ServiceId == "" {
		return ErrIDIsMissing
	}
	return
}
