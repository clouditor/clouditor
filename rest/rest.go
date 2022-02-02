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

package rest

import (
	"context"
	"fmt"
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
	log  *logrus.Entry
	cors *corsConfig
)

// corsConfig holds all necessary configuration options for Cross-Origin Resource Sharing of our REST API.
type corsConfig struct {
	allowedOrigins []string
	allowedHeaders []string
	allowedMethods []string
}

type CORSConfigOption func(*corsConfig)

var (
	// DefaultAllowedOrigins contains an nil slice, as per default, no origins are allowed.
	DefaultAllowedOrigins []string = nil

	// DefaultAllowedHeaders contains sensible defaults for the Access-Control-Allow-Headers header.
	// Please adjust accordingly in production using WithAllowedHeaders.
	DefaultAllowedHeaders = []string{"Content-Type", "Accept", "Authorization"}

	// DefaultAllowedMethods contains sensible defaults for the Access-Control-Allow-Methods header.
	// Please adjust accordingly in production using WithAllowedMethods.
	DefaultAllowedMethods = []string{"GET", "POST", "PUT", "DELETE"}
)

func WithAllowedOrigins(origins []string) CORSConfigOption {
	return func(cc *corsConfig) {
		cc.allowedOrigins = origins
	}
}

func WithAllowedHeaders(headers []string) CORSConfigOption {
	return func(cc *corsConfig) {
		cc.allowedHeaders = headers
	}
}

func WithAllowedMethods(methods []string) CORSConfigOption {
	return func(cc *corsConfig) {
		cc.allowedMethods = methods
	}
}

func init() {
	log = logrus.WithField("component", "rest")

	cors = &corsConfig{
		allowedOrigins: DefaultAllowedOrigins,
		allowedHeaders: DefaultAllowedHeaders,
		allowedMethods: DefaultAllowedMethods,
	}
}

func RunServer(ctx context.Context, grpcPort int, httpPort int, corsOpts ...CORSConfigOption) error {
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

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: allowCORS(mux),
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			break
		}

		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		log.Printf("Shutting down REST gateway on :%d", httpPort)

		_ = srv.Shutdown(ctx)
	}()

	log.Printf("Starting REST gateway on :%d", httpPort)

	return srv.ListenAndServe()
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check, if we allow this specific origin
		origin := r.Header.Get("Origin")
		if originAllowed(origin) {
			// set the appropriate access control header
			w.Header().Set("Access-Control-Allow-Origin", origin)

			// additionally, we need to handle preflight (OPTIONS) requests differently
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				handlePreflight(w, r)
				return
			}
		}
	})
}

func handlePreflight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(cors.allowedHeaders, ","))
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(cors.allowedMethods, ","))
}

func originAllowed(origin string) bool {
	if cors.allowedOrigins == nil || len(cors.allowedOrigins) == 0 {
		return false
	}

	for _, v := range cors.allowedOrigins {
		if origin == v {
			return true
		}
	}

	return false
}
