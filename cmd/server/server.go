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

package main

import (
	"os"

	"clouditor.io/clouditor/v2/cli/commands/evidence"
	"clouditor.io/clouditor/v2/server/commands"
	"clouditor.io/clouditor/v2/server/commands/assessment"
	"clouditor.io/clouditor/v2/server/commands/discovery"
	"clouditor.io/clouditor/v2/server/commands/evaluation"
	"clouditor.io/clouditor/v2/server/commands/orchestrator"
	"clouditor.io/clouditor/v2/server/commands/standalone"
	"github.com/spf13/cobra"
)

func newRootCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "server",
		Short: "server launches a Clouditor server",
		Long:  "It can be used to launch an all-in-one solution of several microservices or each service can be started individually.",
	}

	// AddCommands adds all subcommands
	cmd.AddCommand(
		assessment.NewAssessmentCommand(),
		discovery.NewDiscoveryCommand(),
		evaluation.NewEvaluationCommand(),
		evidence.NewEvidenceCommand(),
		orchestrator.NewOrchestratorCommand(),
		standalone.NewStandaloneCommand(),
	)

	commands.BindPersistentFlags(cmd)

	return cmd
}

func main() {
	var cmd = newRootCommand()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
