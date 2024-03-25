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

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/internal/util"
	"clouditor.io/clouditor/v2/logging/formatter"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/persistence/gorm"
	"clouditor.io/clouditor/v2/persistence/inmemory"
	"clouditor.io/clouditor/v2/server"
	"clouditor.io/clouditor/v2/server/rest"
	"clouditor.io/clouditor/v2/service"
	service_discovery "clouditor.io/clouditor/v2/service/discovery"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	srv              *server.Server
	discoveryService *service_discovery.Service
	db               persistence.Storage
	providers        []string

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
	cobra.OnInitialize(config.InitConfig)

	engineCmd = config.InitCobra(engineCmd)
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
 
  Clouditor Discovery Service Version %s
  `, rt.VersionString())
	fmt.Println()

	level, err = logrus.ParseLevel(viper.GetString(config.LogLevelFlag))
	if err != nil {
		return err
	}
	logrus.SetLevel(level)
	log.Infof("Log level is set to %s", level)

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

	// If no CSPs for discovering are given, take all implemented discoverers
	if len(viper.GetStringSlice(config.DiscoveryProviderFlag)) == 0 {
		providers = []string{service_discovery.ProviderAWS, service_discovery.ProviderAzure, service_discovery.ProviderK8S}
	} else {
		providers = viper.GetStringSlice(config.DiscoveryProviderFlag)
	}

	discoveryService = service_discovery.NewService(
		service_discovery.WithCloudServiceID(viper.GetString(config.CloudServiceIDFlag)),
		service_discovery.WithProviders(providers),
		service_discovery.WithStorage(db),
		service_discovery.WithAssessmentAddress(viper.GetString(config.AssessmentURLFlag)),
		service_discovery.WithOAuth2Authorizer(
			// Configure the OAuth 2.0 client credentials for this service
			&clientcredentials.Config{
				ClientID:     viper.GetString(config.ServiceOAuth2ClientIDFlag),
				ClientSecret: viper.GetString(config.ServiceOAuth2ClientSecretFlag),
				TokenURL:     viper.GetString(config.ServiceOAuth2EndpointFlag),
			}),
	)

	grpcPort := viper.GetUint16(config.APIgRPCPortDiscoveryFlag)
	httpPort := viper.GetUint16(config.APIHTTPPortDiscoveryFlag)

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
	log.Infof("Assessment URL is set to %s", viper.GetString(config.AssessmentURLFlag))

	// Start the gRPC server
	_, srv, err = server.StartGRPCServer(
		fmt.Sprintf("0.0.0.0:%d", grpcPort),
		server.WithDiscovery(discoveryService),
		server.WithExperimentalDiscovery(discoveryService),
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
