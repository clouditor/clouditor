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
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/logging/formatter"
	"clouditor.io/clouditor/v2/server"
	"clouditor.io/clouditor/v2/server/rest"
	"clouditor.io/clouditor/v2/service"
	service_assessment "clouditor.io/clouditor/v2/service/assessment"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	srv               *server.Server
	assessmentService *service_assessment.Service

	log *logrus.Entry
)

var engineCmd = &cobra.Command{
	Use:   "engine",
	Short: "engine launches the Clouditors Discovery Engine",
	Long:  "Discovery Engine is a component of the Clouditor and starts the Discovery Service.",
	RunE:  doCmd,
}

func init() {
	log = logrus.WithField("component", "grpc")
	log.Logger.Formatter = formatter.CapitalizeFormatter{Formatter: &logrus.TextFormatter{ForceColors: true}}
	cobra.OnInitialize(initConfig)

	engineCmd = config.InitCobra(engineCmd)
}

func initConfig() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetEnvPrefix(config.EnvPrefix)
	viper.SetConfigName("clouditor")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
}

func doCmd(_ *cobra.Command, _ []string) (err error) {
	var (
		rt, _ = service.GetRuntimeInfo()
		level logrus.Level
	)

	fmt.Printf(`
           $$\                           $$\ $$\   $$\
           $$ |                          $$ |\__|  $$ |
  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
 $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
 $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
 $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
 \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
 
  Clouditor Assessment Service Version %s
  `, rt.VersionString())
	fmt.Println()

	level, err = logrus.ParseLevel(viper.GetString(config.LogLevelFlag))
	if err != nil {
		return err
	}
	logrus.SetLevel(level)
	log.Infof("Log level is set to %s", level)

	assessmentService = service_assessment.NewService(
		service_assessment.WithOAuth2Authorizer(
			// Configure the OAuth 2.0 client credentials for this service
			&clientcredentials.Config{
				ClientID:     viper.GetString(config.ServiceOAuth2ClientIDFlag),
				ClientSecret: viper.GetString(config.ServiceOAuth2ClientSecretFlag),
				TokenURL:     viper.GetString(config.ServiceOAuth2EndpointFlag),
			},
		),
		service_assessment.WithOrchestratorAddress(viper.GetString(config.OrchestratorURLFlag)),
		service_assessment.WithEvidenceStoreAddress(viper.GetString(config.EvidenceStoreURLFlag)),
	)

	// It is possible to register hook functions for the assessment service.
	//  * The hook functions in assessment are implemented in AssessEvidence(s)

	// assessmentService.RegisterAssessmentResultHook(func(result *assessment.AssessmentResult, err error) {}

	grpcPort := viper.GetUint16(config.APIgRPCPortFlag)
	httpPort := viper.GetUint16(config.APIHTTPPortFlag)

	var opts = []rest.ServerConfigOption{
		rest.WithAllowedOrigins(viper.GetStringSlice(config.APICORSAllowedOriginsFlags)),
		rest.WithAllowedHeaders(viper.GetStringSlice(config.APICORSAllowedHeadersFlags)),
		rest.WithAllowedMethods(viper.GetStringSlice(config.APICORSAllowedMethodsFlags)),
	}

	log.Infof("Starting gRPC endpoint on :%d", grpcPort)
	log.Infof("Orchestrator URL is set to %s", viper.GetString(config.OrchestratorURLFlag))
	log.Infof("Evidence Store URL is set to %s", viper.GetString(config.EvidenceStoreURLFlag))

	// Start the gRPC server
	_, srv, err = server.StartGRPCServer(
		fmt.Sprintf("0.0.0.0:%d", grpcPort),
		server.WithJWKS(viper.GetString(config.APIJWKSURLFlag)),
		server.WithAssessment(assessmentService),
		server.WithReflection(),
	)
	if err != nil {
		log.Errorf("Failed to serve gRPC endpoint: %s", err)
		return err
	}

	// Start the gRPC-HTTP gateway
	err = rest.RunServer(context.Background(),
		grpcPort,
		httpPort,
		opts...,
	)
	if err != nil && err != http.ErrServerClosed {
		log.Errorf("failed to serve gRPC-HTTP gateway: %v", err)
		return err
	}

	assessmentService.Shutdown()

	log.Infof("Stopping gRPC endpoint")
	srv.Stop()

	return nil
}

func main() {
	if err := engineCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
