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
	"crypto/ecdsa"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/auth"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	oauth2 "github.com/oxisto/oauth2go"
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

// ServerConfigOption represents functional-style options to modify the server configuration in RunServer.
type ServerConfigOption func(*corsConfig, *runtime.ServeMux)

var (
	// DefaultAllowedOrigins contains a nil slice, as per default, no origins are allowed.
	DefaultAllowedOrigins []string = nil

	// DefaultAllowedHeaders contains sensible defaults for the Access-Control-Allow-Headers header.
	// Please adjust accordingly in production using WithAllowedHeaders.
	DefaultAllowedHeaders = []string{"Content-Type", "Accept", "Authorization"}

	// DefaultAllowedMethods contains sensible defaults for the Access-Control-Allow-Methods header.
	// Please adjust accordingly in production using WithAllowedMethods.
	DefaultAllowedMethods = []string{"GET", "POST", "PUT", "DELETE"}

	// DefaultAPIHTTPPort specifies the default port for the REST API.
	DefaultAPIHTTPPort uint16 = 8080
)

func init() {
	log = logrus.WithField("component", "rest")

	// initialize the CORS config with restrictive default values, e.g. no origin allowed
	cors = &corsConfig{
		allowedOrigins: DefaultAllowedOrigins,
		allowedHeaders: DefaultAllowedHeaders,
		allowedMethods: DefaultAllowedMethods,
	}
}

// WithAllowedOrigins is an option to supply allowed origins in CORS.
func WithAllowedOrigins(origins []string) ServerConfigOption {
	return func(cc *corsConfig, _ *runtime.ServeMux) {
		cc.allowedOrigins = origins
	}
}

// WithAllowedHeaders is an option to supply allowed headers in CORS.
func WithAllowedHeaders(headers []string) ServerConfigOption {
	return func(cc *corsConfig, _ *runtime.ServeMux) {
		cc.allowedHeaders = headers
	}
}

// WithAllowedMethods is an option to supply allowed methods in CORS.
func WithAllowedMethods(methods []string) ServerConfigOption {
	return func(cc *corsConfig, _ *runtime.ServeMux) {
		cc.allowedMethods = methods
	}
}

// WithAdditionalHandler is an option to add an additional handler func in the REST server.
func WithAdditionalHandler(method string, path string, h runtime.HandlerFunc) ServerConfigOption {
	return func(_ *corsConfig, sm *runtime.ServeMux) {
		_ = sm.HandlePath(method, path, h)
	}
}

// WithEmbeddedOAuth2Server configures our REST gateway to include an embedded OAuth 2.0
// authorization server with the given parameters. This server can be used to authenticate
// and authorize API clients, such as our own micro-services as well as users logging in
// through the UI.
//
// For the various options to configure the OAuth 2.0 server, please refer to
// https://pkg.go.dev/github.com/oxisto/oauth2go#AuthorizationServerOption.
//
// In production scenarios, the usage of a dedicated authentication and authorization server is
// recommended.
func WithEmbeddedOAuth2Server(keyPath string, keyPassword string, saveOnCreate bool, opts ...oauth2.AuthorizationServerOption) ServerConfigOption {
	return func(cc *corsConfig, sm *runtime.ServeMux) {
		opts = append(opts, oauth2.WithSigningKeysFunc(func() map[int]*ecdsa.PrivateKey {
			return auth.LoadSigningKeys(keyPath, keyPassword, saveOnCreate)
		}))

		authSrv := oauth2.NewServer("", opts...)
		authHandler := func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			http.StripPrefix("/v1/auth", authSrv.Handler).ServeHTTP(w, r)
		}

		WithAdditionalHandler("GET", "/.well-known/jwks.json", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			authSrv.Handler.(*http.ServeMux).ServeHTTP(w, r)
		})(cc, sm)
		WithAdditionalHandler("GET", "/v1/auth/login", authHandler)(cc, sm)
		WithAdditionalHandler("GET", "/v1/auth/authorize", authHandler)(cc, sm)
		WithAdditionalHandler("POST", "/v1/auth/login", authHandler)(cc, sm)
		WithAdditionalHandler("POST", "/v1/auth/token", authHandler)(cc, sm)
	}
}

// RunServer starts our REST API. The REST API is a reverse proxy using grpc-gateway that
// exposes certain gRPC calls as RESTful HTTP methods.
func RunServer(ctx context.Context, grpcPort uint16, httpPort uint16, serverOpts ...ServerConfigOption) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	for _, o := range serverOpts {
		o(cors, mux)
	}

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

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
		Addr:              fmt.Sprintf(":%d", httpPort),
		Handler:           handleCORS(mux),
		ReadHeaderTimeout: 2 * time.Second,
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			break
		}

		StopServer(ctx)
	}()

	log.Printf("Starting REST gateway on :%d", httpPort)

	sock, err = net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}

	go func() {
		ready <- true
	}()

	return srv.Serve(sock)
}

// StopServer is in charge of stopping the REST API.
func StopServer(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	log.Printf("Shutting down REST gateway")

	// Clear our ready channel, otherwise, this will block the exit
	select {
	case <-ready:
	default:
	}

	_ = srv.Shutdown(ctx)
}

// GetReadyChannel returns a channel which will notify when the server is ready
func GetReadyChannel() chan bool {
	return ready
}

// GetServerPort returns the actual port used by the REST API
func GetServerPort() (uint16, error) {
	if sock == nil {
		return 0, errors.New("server socket is empty")
	}

	tcpAddr, ok := sock.Addr().(*net.TCPAddr)

	if !ok {
		return 0, errors.New("server socket address is not a TCP address")
	}

	return tcpAddr.AddrPort().Port(), nil
}

// handleCORS adds an appropriate http.HandlerFunc to an existing http.Handler to configure
// Cross-Origin Resource Sharing (CORS) according to our global configuration.
func handleCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check, if we allow this specific origin
		origin := r.Header.Get("Origin")
		if cors.OriginAllowed(origin) {
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

// OriginAllowed checks if the supplied origin is allowed according to our global CORS configuration.
func (cors *corsConfig) OriginAllowed(origin string) bool {
	// If no origin is specified, we are running in a non-browser environment and
	// this means, that all origins are allowed
	if origin == "" {
		return true
	}

	for _, v := range cors.allowedOrigins {
		if origin == v {
			return true
		}
	}

	return false
}
