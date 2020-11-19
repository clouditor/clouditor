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
	"fmt"
	"log"
	"net"

	"clouditor.io/clouditor"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/rest"
	"clouditor.io/clouditor/service/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var server *grpc.Server
var authService auth.Service

func main() {
	persistence.InitPostgreSQL("localhost")

	createDefaultUser()

	fmt.Printf("Welcome to new Clouditor 2.0\n\n")

	grpcPort := 9090
	httpPort := 8080

	// create a new socket for gRPC communication
	sock, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("could not listen: %v", err)
	}

	server = grpc.NewServer()
	clouditor.RegisterAuthenticationServer(server, &authService)

	// enable reflection, primary for testing in early stages
	reflection.Register(server)

	// start the gRPC-HTTP gateway
	go func() {
		if err = rest.RunServer(context.Background(), grpcPort, httpPort); err != nil {
			log.Fatalf("failed to serve gRPC-HTTP gateway: %v", err)
		}
	}()

	// serve the gRPC socket
	if err := server.Serve(sock); err != nil {
		log.Fatalf("failed to serve gRPC endpoint: %v", err)
	}
}

func createDefaultUser() {
	db := persistence.GetDatabase()

	var count int64
	db.Model(&clouditor.User{}).Count(&count)

	if count == 0 {
		password, _ := authService.HashPassword("clouditor")

		user := clouditor.User{
			Username: "clouditor",
			FullName: "clouditor",
			Password: string(password),
		}

		log.Printf("Creating default user %s\n", user.Username)

		db.Create(&user)
	}
}
