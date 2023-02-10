// Copyright 2016-2022 Fraunhofer AISEC
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

package clouditor

import (
	"context"
	"runtime/debug"

	"clouditor.io/clouditor/api/orchestrator"
	"github.com/sirupsen/logrus"
)

var (
	Version string
	log     *logrus.Entry
)

type Service struct {
	// runtime is a struct for all necessary Clouditor runtime information
	runtime *orchestrator.Runtime
}

// NewService creates a new Orchestrator service
func NewService() *Service {
	var err error
	s := Service{
		runtime: &orchestrator.Runtime{
			Dependencies: []*orchestrator.Dependency{},
		},
	}

	// Set build info
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		log.Errorf("error reading build info: %v", err)
	} else {
		s.runtime.GolangVersion = buildInfo.GoVersion
		// Set dependencies
		for _, d := range buildInfo.Deps {
			s.runtime.Dependencies = append(s.runtime.Dependencies, &orchestrator.Dependency{
				Path:    d.Path,
				Version: d.Version,
			})
		}

		// Set version control system info
		for _, setting := range buildInfo.Settings {
			switch setting.Key {
			case "vcs.revision":
				s.runtime.CommitHash = setting.Value
			case "vcs":
				s.runtime.Vcs = setting.Value
			}
		}
	}
	return &s
}

// GetRuntimeInfo implements method to get Clouditors runtime information
func (svc *Service) GetRuntimeInfo(_ context.Context, _ *orchestrator.RuntimeRequest) (runtime *orchestrator.RuntimeResponse, err error) {
	return &orchestrator.RuntimeResponse{Runtime: svc.runtime}, nil
}
