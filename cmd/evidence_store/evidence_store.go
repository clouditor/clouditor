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
	"clouditor.io/clouditor/v2/internal/launcher"
	"clouditor.io/clouditor/v2/server"
	service_evidenceStore "clouditor.io/clouditor/v2/service/evidence"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var engineCmd = &cobra.Command{
	Use:   "evidence-store",
	Short: "evidence-store launches the Clouditor Evidence Store Service",
	Long:  "Evidence Store is a component of the Clouditor and starts the Evidence Store Service.",
	RunE:  doCmd,
}

func init() {
	config.InitCobra(engineCmd)
}

func doCmd(cmd *cobra.Command, _ []string) (err error) {
	l, err := launcher.NewLauncher[service_evidenceStore.Service](
		cmd.Use,
		service_evidenceStore.NewService,
		service_evidenceStore.WithStorage,
		func(svc *service_evidenceStore.Service) ([]server.StartGRPCServerOption, error) {
			// It is possible to register hook functions for the evidenceStore.
			//  * The hook functions in evidenceStore are implemented in StoreEvidence(s)
			// evidenceStoreService.RegisterEvidenceHook(func(result *evidence.Evidence, err error) {})

			return []server.StartGRPCServerOption{
				server.WithJWKS(viper.GetString(config.APIJWKSURLFlag)),
				server.WithEvidenceStore(svc),
				server.WithReflection(),
			}, nil
		},
	)
	if err != nil {
		return err
	}

	// Start the gRPC server and the corresponding gRPC-HTTP gateway
	return l.Launch()
}

func main() {
	if err := engineCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
