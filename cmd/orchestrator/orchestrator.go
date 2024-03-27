// Copyright 2016-2020 Fraunhofer AISEC
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

	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/logging/formatter"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/server"
	service_orchestrator "clouditor.io/clouditor/v2/service/orchestrator"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	srv                 *server.Server
	orchestratorService *service_orchestrator.Service
	db                  persistence.Storage

	log *logrus.Entry
)

var engineCmd = &cobra.Command{
	Use:   "orchestrator",
	Short: "orchestrator launches the Clouditor Orchestrator Service",
	Long:  "Orchestrator is a component of the Clouditor and starts the Orchestrator Service.",
	RunE:  doCmd,
}

func init() {
	log = logrus.WithField("component", "orchestrator-grpc")
	log.Logger.Formatter = formatter.CapitalizeFormatter{Formatter: &logrus.TextFormatter{ForceColors: true}}
	cobra.OnInitialize(config.InitConfig)

	engineCmd = config.InitCobra(engineCmd)
}

func doCmd(_ *cobra.Command, _ []string) (err error) {
	var (
		grpcOpts []server.StartGRPCServerOption
	)

	// Print Clouditor header with the used Clouditor version
	config.PrintClouditorHeader("Clouditor Orchestrator Service")

	// Set log level
	log, err = config.SetLogLevel(log)
	if err != nil {
		return err
	}

	// Set storage
	db, err = config.SetStorage()
	if err != nil {
		return err
	}

	// Create new Orchestrator Service
	orchestratorService = service_orchestrator.NewService(service_orchestrator.WithStorage(db))

	// It is possible to register hook functions for the orchestrator.
	//  * The hook functions in orchestrator are implemented in StoreAssessmentResult(s)

	// orchestratorService.RegisterAssessmentResultHook(func(result *assessment.AssessmentResult, err error) {})

	// Create default target Cloud Service
	if viper.GetBool(config.CreateDefaultTarget) {
		_, err := orchestratorService.CreateDefaultTargetCloudService()
		if err != nil {
			log.Errorf("could not register default target cloud service: %v", err)
		}
	}

	// Start the gRPC server and the corresponding gRPC-HTTP gateway
	grpcOpts = []server.StartGRPCServerOption{
		server.WithJWKS(viper.GetString(config.APIJWKSURLFlag)),
		server.WithOrchestrator(orchestratorService),
		server.WithReflection()}

	srv, err = config.StartServer(log, grpcOpts...)
	if err != nil {
		log.Errorf("could not register default target cloud service: %v", err)
	}

	return nil
}

func main() {
	if err := engineCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
