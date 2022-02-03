// Copyright 2016-2022 Fraunhofer AISEC
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

package rest

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"clouditor.io/clouditor/api/evidence"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/auth"
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/api/orchestrator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	log *logrus.Entry

	// srv holds the global http.Server of our REST API.
	srv *http.Server

	// sock holds the listener socket of our REST API.
	sock net.Listener

	ready = make(chan bool)

	// cors holds the global CORS configuration.
	cors *corsConfig
)

// corsConfig holds all necessary configuration options for Cross-Origin Resource Sharing of our REST API.
type corsConfig struct {
	// allowedOrigins contains a list of the allowed origins
	allowedOrigins []string

	// allowedHeaders contains a list of the allowed headers
	allowedHeaders []string

	// allowedMethods contains a list of the allowed methods
	allowedMethods []string
}

// CORSConfigOption represents functional-style options to modify the CORS configuration in RunServer.
type CORSConfigOption func(*corsConfig)

var (
	// DefaultAllowedOrigins contains a nil slice, as per default, no origins are allowed.
	DefaultAllowedOrigins []string = nil

	// DefaultAllowedHeaders contains sensible defaults for the Access-Control-Allow-Headers header.
	// Please adjust accordingly in production using WithAllowedHeaders.
	DefaultAllowedHeaders = []string{"Content-Type", "Accept", "Authorization"}

	// DefaultAllowedMethods contains sensible defaults for the Access-Control-Allow-Methods header.
	// Please adjust accordingly in production using WithAllowedMethods.
	DefaultAllowedMethods = []string{"GET", "POST", "PUT", "DELETE"}
)

// WithAllowedOrigins is an option to supply allowed origins in CORS.
func WithAllowedOrigins(origins []string) CORSConfigOption {
	return func(cc *corsConfig) {
		cc.allowedOrigins = origins
	}
}

// WithAllowedHeaders is an option to supply allowed headers in CORS.
func WithAllowedHeaders(headers []string) CORSConfigOption {
	return func(cc *corsConfig) {
		cc.allowedHeaders = headers
	}
}

// WithAllowedMethods is an option to supply allowed methods in CORS.
func WithAllowedMethods(methods []string) CORSConfigOption {
	return func(cc *corsConfig) {
		cc.allowedMethods = methods
	}
}

func init() {
	log = logrus.WithField("component", "rest")

	// initialize the CORS config with restrictive default values, e.g. no origin allowed
	cors = &corsConfig{
		allowedOrigins: DefaultAllowedOrigins,
		allowedHeaders: DefaultAllowedHeaders,
		allowedMethods: DefaultAllowedMethods,
	}
}

// RunServer starts our REST API. The REST API is a reverse proxy using grpc-gateway that
// exposes certain gRPC calls as RESTful HTTP methods.
func RunServer(ctx context.Context, grpcPort int, httpPort int, corsOpts ...CORSConfigOption) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, o := range corsOpts {
		o(cors)
	}

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := auth.RegisterAuthenticationHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), opts); err != nil {
		return fmt.Errorf("failed to connect to authentication gRPC service %w", err)
	}

	if err := discovery.RegisterDiscoveryHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), opts); err != nil {
		return fmt.Errorf("failed to connect to discovery gRPC service %w", err)
	}

	if err := assessment.RegisterAssessmentHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), opts); err != nil {
		return fmt.Errorf("failed to connect to assessment gRPC service %w", err)
	}

	if err := orchestrator.RegisterOrchestratorHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), opts); err != nil {
		return fmt.Errorf("failed to connect to orchestrator gRPC service %w", err)
	}

	if err := evidence.RegisterEvidenceStoreHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), opts); err != nil {
		return fmt.Errorf("failed to connect to evidence gRPC service %w", err)
	}

	srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: handleCORS(mux),
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			break
		}

		ready <- false

		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		log.Printf("Shutting down REST gateway on :%d", httpPort)

		_ = srv.Shutdown(ctx)
	}()

	log.Printf("Starting REST gateway on :%d", httpPort)

	sock, err = net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}

	ready <- true

	return srv.Serve(sock)
}

// GetServerPort returns the actual port used by the REST API
func GetServerPort() (int, error) {
	if sock == nil {
		return 0, errors.New("server socket is empty")
	}

	tcpAddr, ok := sock.Addr().(*net.TCPAddr)

	if !ok {
		return 0, errors.New("server socket address is not a TCP address")
	}

	return tcpAddr.Port, nil
}

// handleCORS adds an appropriate http.HandlerFunc to an existing http.Handler to configure
// Cross-Origin Resource Sharing (CORS) according to our global configuration.
func handleCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check, if we allow this specific origin
		origin := r.Header.Get("Origin")
		if originAllowed(origin) {
			// Set the appropriate access control header
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Vary", "Origin")

			// Additionally, we need to handle preflight (OPTIONS) requests to specify allowed headers and methods
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(cors.allowedHeaders, ","))
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(cors.allowedMethods, ","))
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}

// originAllowed checks if the supplied origin is allowed according to our global CORS configuration.
func originAllowed(origin string) bool {
	for _, v := range cors.allowedOrigins {
		if origin == v {
			return true
		}
	}

	return false
}
