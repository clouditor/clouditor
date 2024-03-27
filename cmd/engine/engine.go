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

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/evaluation"
	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/internal/util"
	"clouditor.io/clouditor/v2/logging/formatter"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/persistence/gorm"
	"clouditor.io/clouditor/v2/persistence/inmemory"
	"clouditor.io/clouditor/v2/server"
	"clouditor.io/clouditor/v2/server/rest"
	"clouditor.io/clouditor/v2/service"
	service_assessment "clouditor.io/clouditor/v2/service/assessment"
	service_discovery "clouditor.io/clouditor/v2/service/discovery"
	service_evaluation "clouditor.io/clouditor/v2/service/evaluation"
	service_evidenceStore "clouditor.io/clouditor/v2/service/evidence"
	service_orchestrator "clouditor.io/clouditor/v2/service/orchestrator"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	srv                  *server.Server
	discoveryService     *service_discovery.Service
	orchestratorService  *service_orchestrator.Service
	assessmentService    *service_assessment.Service
	evidenceStoreService evidence.EvidenceStoreServer
	evaluationService    evaluation.EvaluationServer
	db                   persistence.Storage
	providers            []string

	log *logrus.Entry
)

var engineCmd = &cobra.Command{
	Use:   "engine",
	Short: "engine launches the Clouditor Engine",
	Long:  "Clouditor Engine is the main component of Clouditor. It is an all-in-one solution of several microservices, which also can be started individually.",
	RunE:  doCmd,
}

func init() {
	log = logrus.WithField("component", "grpc")
	log.Logger.Formatter = formatter.CapitalizeFormatter{Formatter: &logrus.TextFormatter{ForceColors: true}}
	config.InitCobra(engineCmd)
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
 
  Version %s
  `, rt.VersionString())
	fmt.Println()

	level, err = logrus.ParseLevel(viper.GetString(config.LogLevelFlag))
	if err != nil {
		return err
	}
	logrus.SetLevel(level)

	if viper.GetBool(config.DBInMemoryFlag) {
		db, err = inmemory.NewStorage()
	} else {
		db, err = gorm.NewStorage(gorm.WithPostgres(
			viper.GetString(config.DBHostFlag),
			viper.GetUint16(config.DBPortFlag),
			viper.GetString(config.DBUserNameFlag),
			viper.GetString(config.DBPasswordFlag),
			viper.GetString(config.DBNameFlag),
			viper.GetString(config.DBSSLModeFlag),
		))
	}
	if err != nil {
		// We could also just log the error and forward db = nil which will result in inmemory storages for each service
		// below
		return fmt.Errorf("could not create storage: %w", err)
	}

	// If no CSPs for discovering is given, take all implemented discoverers
	if len(viper.GetStringSlice(config.DiscoveryProviderFlag)) == 0 {
		providers = []string{service_discovery.ProviderAWS, service_discovery.ProviderAzure, service_discovery.ProviderK8S}
	} else {
		providers = viper.GetStringSlice(config.DiscoveryProviderFlag)
	}

	discoveryService = service_discovery.NewService(
		service_discovery.WithProviders(providers),
		service_discovery.WithStorage(db),
		service_discovery.WithOAuth2Authorizer(
			// Configure the OAuth 2.0 client credentials for this service
			&clientcredentials.Config{
				ClientID:     viper.GetString(config.ServiceOAuth2ClientIDFlag),
				ClientSecret: viper.GetString(config.ServiceOAuth2ClientSecretFlag),
				TokenURL:     viper.GetString(config.ServiceOAuth2EndpointFlag),
			}),
	)

	orchestratorService = service_orchestrator.NewService(service_orchestrator.WithStorage(db))

	assessmentService = service_assessment.NewService(
		service_assessment.WithOAuth2Authorizer(
			// Configure the OAuth 2.0 client credentials for this service
			&clientcredentials.Config{
				ClientID:     viper.GetString(config.ServiceOAuth2ClientIDFlag),
				ClientSecret: viper.GetString(config.ServiceOAuth2ClientSecretFlag),
				TokenURL:     viper.GetString(config.ServiceOAuth2EndpointFlag),
			},
		),
	)

	evidenceStoreService = service_evidenceStore.NewService(service_evidenceStore.WithStorage(db))

	evaluationService = service_evaluation.NewService(
		service_evaluation.WithOAuth2Authorizer(
			// Configure the OAuth 2.0 client credentials for this service
			&clientcredentials.Config{
				ClientID:     viper.GetString(config.ServiceOAuth2ClientIDFlag),
				ClientSecret: viper.GetString(config.ServiceOAuth2ClientSecretFlag),
				TokenURL:     viper.GetString(config.ServiceOAuth2EndpointFlag),
			},
		),
		service_evaluation.WithStorage(db),
	)

	// It is possible to register hook functions for the orchestrator, evidenceStore and assessment service.
	//  * The hook functions in orchestrator are implemented in StoreAssessmentResult(s)
	//  * The hook functions in evidenceStore are implemented in StoreEvidence(s)
	//  * The hook functions in assessment are implemented in AssessEvidence(s)

	// orchestratorService.RegisterAssessmentResultHook(func(result *assessment.AssessmentResult, err error) {})
	// evidenceStoreService.RegisterEvidenceHook(func(result *evidence.Evidence, err error) {})
	// assessmentService.RegisterAssessmentResultHook(func(result *assessment.AssessmentResult, err error) {}

	if viper.GetBool(config.CreateDefaultTargetFlag) {
		_, err := orchestratorService.CreateDefaultTargetCloudService()
		if err != nil {
			log.Errorf("could not register default target cloud service: %v", err)
		}
	}

	grpcPort := viper.GetUint16(config.APIgRPCPortFlag)
	httpPort := viper.GetUint16(config.APIHTTPPortFlag)

	var opts = []rest.ServerConfigOption{
		rest.WithAllowedOrigins(viper.GetStringSlice(config.APICORSAllowedOriginsFlags)),
		rest.WithAllowedHeaders(viper.GetStringSlice(config.APICORSAllowedHeadersFlags)),
		rest.WithAllowedMethods(viper.GetStringSlice(config.APICORSAllowedMethodsFlags)),
	}

	// Automatically start the discovery, if we have this flag enabled
	if viper.GetBool(config.DiscoveryAutoStartFlag) {
		go func() {
			<-rest.GetReadyChannel()
			_, err = discoveryService.Start(context.Background(), &discovery.StartDiscoveryRequest{
				ResourceGroup: util.Ref(viper.GetString(config.DiscoveryResourceGroupFlag)),
			})
			if err != nil {
				log.Errorf("Could not automatically start discovery: %v", err)
			}
		}()
	}

	log.Infof("Starting gRPC endpoint on :%d", grpcPort)

	// Start the gRPC server
	_, srv, err = server.StartGRPCServer(
		fmt.Sprintf("0.0.0.0:%d", grpcPort),
		server.WithJWKS(viper.GetString(config.APIJWKSURLFlag)),
		server.WithDiscovery(discoveryService),
		server.WithExperimentalDiscovery(discoveryService),
		server.WithOrchestrator(orchestratorService),
		server.WithAssessment(assessmentService),
		server.WithEvidenceStore(evidenceStoreService),
		server.WithEvaluation(evaluationService),
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
	discoveryService.Shutdown()

	log.Infof("Stopping gRPC endpoint")
	srv.Stop()

	return nil
}

func main() {
	if err := engineCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
