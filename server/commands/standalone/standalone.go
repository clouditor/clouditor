// Copyright 2024 Fraunhofer AISEC
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

package standalone

import (
	"clouditor.io/clouditor/v2/internal/launcher"
	command_assessment "clouditor.io/clouditor/v2/server/commands/assessment"
	command_discovery "clouditor.io/clouditor/v2/server/commands/discovery"
	command_evaluation "clouditor.io/clouditor/v2/server/commands/evaluation"
	command_orchestrator "clouditor.io/clouditor/v2/server/commands/orchestrator"
	"clouditor.io/clouditor/v2/service/assessment"
	"clouditor.io/clouditor/v2/service/discovery"
	"clouditor.io/clouditor/v2/service/evaluation"
	evidences "clouditor.io/clouditor/v2/service/evidence"
	"clouditor.io/clouditor/v2/service/orchestrator"

	"github.com/spf13/cobra"
)

func NewStandaloneCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "standalone",
		Short: "Starts a server which contains all Clouditor service",
		Long:  "This command starts all Clouditor services in standalone mode",
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := launcher.NewLauncher(cmd.Use,
				assessment.DefaultServiceSpec(),
				discovery.DefaultServiceSpec(),
				evidences.DefaultServiceSpec(),
				evaluation.DefaultServiceSpec(),
				orchestrator.DefaultServiceSpec(),
			)
			if err != nil {
				return err
			}

			return l.Launch()
		},
	}

	command_assessment.BindFlags(cmd)
	command_discovery.BindFlags(cmd)
	command_evaluation.BindFlags(cmd)
	command_orchestrator.BindFlags(cmd)

	return cmd
}
