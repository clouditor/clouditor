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
	"encoding/json"
	"fmt"

	"clouditor.io/clouditor/api/orchestrator"
)

// LoadRequirements loads requirements definitions from a JSON file.
func LoadRequirements(file string) (requirements []*orchestrator.Requirement, err error) {
	var (
		b []byte
	)

	b, err = f.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error while loading %s: %w", file, err)
	}

	err = json.Unmarshal(b, &requirements)
	if err != nil {
		return nil, fmt.Errorf("error in JSON marshal: %w", err)
	}

	return requirements, nil
}

func (svc *Service) ListRequirements(_ context.Context, _ *orchestrator.ListRequirementsRequest) (response *orchestrator.ListRequirementsResponse, err error) {
	return &orchestrator.ListRequirementsResponse{
		Requirements: svc.requirements,
	}, nil
}
