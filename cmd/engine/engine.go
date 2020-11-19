/*
 * Copyright 2016-2020 Fraunhofer AISEC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *           $$\                           $$\ $$\   $$\
 *           $$ |                          $$ |\__|  $$ |
 *  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
 * $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
 * $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
 * $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
 * \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
 *  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
 *
 * This file is part of Clouditor Community Edition.
 */

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"clouditor.io/clouditor"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/rest"
	"clouditor.io/clouditor/service/auth"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	APIDefaultUserFlag     = "api-default-user"
	APIDefaultPasswordFlag = "api-default-password"
	APISecretFlag          = "api-secret"
	APIgRPCPortFlag        = "api-grpc-port"
	APIHTTPPortFlag        = "api-http-port"
	DBUserNameFlag         = "db-user-name"
	DBPasswordFlag         = "db-password"
	DBHostFlag             = "db-host"
	DBNameFlag             = "db-name"
	DBPortFlag             = "db-port"
	DBInMemoryFlag         = "db-in-memory"

	DefaultAPIDefaultUser     = "clouditor"
	DefaultAPIDefaultPassword = "clouditor"
	DefaultAPISecret          = "changeme"
	DefaultAPIgRPCPort        = 9090
	DefaultAPIHTTPPort        = 8080
	DefaultDBUserName         = "postgres"
	DefaultDBPassword         = "postgres"
	DefaultDBHost             = "localhost"
	DefaultDBName             = "postgres"
	DefaultDBPort             = 5432
	DefaultDBInMemory         = false

	EnvPrefix = "CLOUDITOR"
)

var server *grpc.Server
var authService *auth.Service

var engineCmd = &cobra.Command{
	Use:   "engine",
	Short: "engine launches the Clouditor Engine",
	Long:  "Clouditor Engine is the main component of Clouditor",
	RunE:  doCmd,
}

func init() {
	cobra.OnInitialize(initConfig)

	engineCmd.Flags().String(APIDefaultUserFlag, DefaultAPIDefaultUser, "Specifies the default API username")
	engineCmd.Flags().String(APIDefaultPasswordFlag, DefaultAPIDefaultPassword, "Specifies the default API password")
	engineCmd.Flags().String(APISecretFlag, DefaultAPISecret, "Specifies the secret used by API tokens")
	engineCmd.Flags().Int16(APIgRPCPortFlag, DefaultAPIgRPCPort, "Specifies the port used for the gRPC API")
	engineCmd.Flags().Int16(APIHTTPPortFlag, DefaultAPIHTTPPort, "Specifies the port used for the HTTP API")
	engineCmd.Flags().String(DBUserNameFlag, DefaultDBUserName, "Provides user name of database")
	engineCmd.Flags().String(DBPasswordFlag, DefaultDBPassword, "Provides password of database")
	engineCmd.Flags().String(DBHostFlag, DefaultDBHost, "Provides address of database")
	engineCmd.Flags().String(DBNameFlag, DefaultDBName, "Provides name of database")
	engineCmd.Flags().Int16(DBPortFlag, DefaultDBPort, "Provides port for database")
	engineCmd.Flags().Bool(DBInMemoryFlag, DefaultDBInMemory, "Uses an in-memory database which is not persisted at all")

	viper.BindPFlag(APIDefaultUserFlag, engineCmd.Flags().Lookup(APIDefaultUserFlag))
	viper.BindPFlag(APIDefaultPasswordFlag, engineCmd.Flags().Lookup(APIDefaultPasswordFlag))
	viper.BindPFlag(APISecretFlag, engineCmd.Flags().Lookup(APISecretFlag))
	viper.BindPFlag(APIgRPCPortFlag, engineCmd.Flags().Lookup(APIgRPCPortFlag))
	viper.BindPFlag(APIHTTPPortFlag, engineCmd.Flags().Lookup(APIHTTPPortFlag))
	viper.BindPFlag(DBUserNameFlag, engineCmd.Flags().Lookup(DBUserNameFlag))
	viper.BindPFlag(DBPasswordFlag, engineCmd.Flags().Lookup(DBPasswordFlag))
	viper.BindPFlag(DBHostFlag, engineCmd.Flags().Lookup(DBHostFlag))
	viper.BindPFlag(DBNameFlag, engineCmd.Flags().Lookup(DBNameFlag))
	viper.BindPFlag(DBPortFlag, engineCmd.Flags().Lookup(DBPortFlag))
	viper.BindPFlag(DBInMemoryFlag, engineCmd.Flags().Lookup(DBInMemoryFlag))
}

func initConfig() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetEnvPrefix(EnvPrefix)
	viper.SetConfigName("clouditor")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
}

func doCmd(cmd *cobra.Command, args []string) (err error) {
	log.Println("Welcome to new Clouditor 2.0")

	fmt.Println(`
           $$\                           $$\ $$\   $$\
           $$ |                          $$ |\__|  $$ |
  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
 $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
 $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
 $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
 \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
 `)

	persistence.InitDB(viper.GetBool(DBInMemoryFlag),
		viper.GetString(DBHostFlag),
		int16(viper.GetInt(DBPortFlag)))

	authService = &auth.Service{
		TokenSecret: viper.GetString(APISecretFlag),
	}

	createDefaultUser()

	grpcPort := viper.GetInt(APIgRPCPortFlag)
	httpPort := viper.GetInt(APIHTTPPortFlag)

	// create a new socket for gRPC communication
	sock, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("could not listen: %v", err)
	}

	server = grpc.NewServer()
	clouditor.RegisterAuthenticationServer(server, authService)

	// enable reflection, primary for testing in early stages
	reflection.Register(server)

	// start the gRPC-HTTP gateway
	go func() {
		err = rest.RunServer(context.Background(), grpcPort, httpPort)
		if errors.Is(err, http.ErrServerClosed) {
			os.Exit(0)
			return
		}

		if err != nil {
			log.Fatalf("failed to serve gRPC-HTTP gateway: %v", err)
		}
	}()

	log.Printf("Starting gRPC endpoint on :%d", grpcPort)

	// serve the gRPC socket
	if err := server.Serve(sock); err != nil {
		log.Printf("failed to serve gRPC endpoint: %v", err)
		return err
	}

	return nil
}

func main() {
	if err := engineCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

// createDefaultUser creates a default user in the database
func createDefaultUser() {
	db := persistence.GetDatabase()

	var count int64
	db.Model(&clouditor.User{}).Count(&count)

	if count == 0 {
		password, _ := authService.HashPassword(viper.GetString(APIDefaultPasswordFlag))

		user := clouditor.User{
			Username: viper.GetString(APIDefaultUserFlag),
			FullName: viper.GetString(APIDefaultUserFlag),
			Password: string(password),
		}

		log.Printf("Creating default user %s\n", user.Username)

		db.Create(&user)
	}
}
