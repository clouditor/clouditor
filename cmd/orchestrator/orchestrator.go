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
	"fmt"
	"os"

	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/internal/launcher"
	"clouditor.io/clouditor/v2/server"
	service_orchestrator "clouditor.io/clouditor/v2/service/orchestrator"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var engineCmd = &cobra.Command{
	Use:   "orchestrator",
	Short: "orchestrator launches the Clouditor Orchestrator Service",
	Long:  "Orchestrator is a component of the Clouditor and starts the Orchestrator Service.",
	RunE:  doCmd,
}

var OrchestratorSpec = launcher.NewServiceSpec(
	service_orchestrator.NewService,
	service_orchestrator.WithStorage,
	func(svc *service_orchestrator.Service) ([]server.StartGRPCServerOption, error) {
		// It is possible to register hook functions for the orchestrator.
		//  * The hook functions in orchestrator are implemented in StoreAssessmentResult(s)

		// svc.RegisterAssessmentResultHook(func(result *assessment.AssessmentResult, err error) {})

		// Create default target Cloud Service
		if viper.GetBool(config.CreateDefaultTargetFlag) {
			_, err := svc.CreateDefaultTargetCloudService()
			if err != nil {
				return nil, fmt.Errorf("could not register default target cloud service: %v", err)
			}
		}

		return []server.StartGRPCServerOption{
			server.WithOrchestrator(svc),
		}, nil
	},
)

func init() {
	config.InitCobra(engineCmd)
}

func doCmd(cmd *cobra.Command, _ []string) (err error) {
	l, err := launcher.NewLauncher(
		cmd.Use,
		OrchestratorSpec,
	)
	if err != nil {
		return err
	}

	return l.Launch()
}

func main() {
	if err := engineCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
