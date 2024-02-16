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

package service

import (
	"runtime/debug"
	"time"

	"clouditor.io/clouditor/v2/api/runtime"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	// version can contain a version that is set externally using ldflags. Otherwise it will be empty
	version string

	// rt is a struct for all necessary Clouditor runtime information
	rt runtime.Runtime

	// populated is set to true once the runtime information has been succesfully populated
	populated = false
)

// populateRuntimeInfo populates the runtime info using the build info and other sources.
func populateRuntimeInfo() {
	if populated {
		return
	}

	// Set build info
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		log.Errorf("Could not read build info. Runtime information will not be complete.")
	} else {
		if version != "" {
			rt.ReleaseVersion = &version
		}
		rt.GolangVersion = buildInfo.GoVersion
		// Set dependencies
		for _, d := range buildInfo.Deps {
			rt.Dependencies = append(rt.Dependencies, &runtime.Dependency{
				Path:    d.Path,
				Version: d.Version,
			})
		}

		// Set version control system info
		for _, setting := range buildInfo.Settings {
			switch setting.Key {
			case "vcs.time":
				var t, err = time.Parse(time.RFC3339, setting.Value)
				if err == nil {
					rt.CommitTime = timestamppb.New(t)
				}
			case "vcs.revision":
				rt.CommitHash = setting.Value
			case "vcs":
				rt.Vcs = setting.Value
			}
		}

		populated = true
	}
}

// GetRuntimeInfo implements method to get Clouditors runtime information
func GetRuntimeInfo() (*runtime.Runtime, error) {
	// Make sure, runtime info is populated
	populateRuntimeInfo()

	return &rt, nil
}
