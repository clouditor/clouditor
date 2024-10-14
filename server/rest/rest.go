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

	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/evaluation"
	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	oauth2 "github.com/oxisto/oauth2go"
	"github.com/oxisto/oauth2go/storage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	log *logrus.Entry

	// srv holds the global http.Server of our REST API.
	srv *http.Server

	// httpPort holds the used HTTP port of the http.Server
	httpPort uint16

	// sock holds the listener socket of our REST API.
	sock net.Listener

	ready = make(chan bool)

	cnf restConfig
)

// restConfig holds different configuration options for the REST gateway.
type restConfig struct {
	// cors holds the global CORS configuration.
	cors *corsConfig

	// opts contains the gRPC options used for the gateway-to-backend calls.
	opts []grpc.DialOption
}

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
type ServerConfigOption func(*restConfig, *runtime.ServeMux)

func init() {
	log = logrus.WithField("component", "rest")

	// initialize the CORS config with restrictive default values, e.g. no origin allowed
	cnf.cors = &corsConfig{
		allowedOrigins: config.DefaultAllowedOrigins,
		allowedHeaders: config.DefaultAllowedHeaders,
		allowedMethods: config.DefaultAllowedMethods,
	}
}

// WithAllowedOrigins is an option to supply allowed origins in CORS.
func WithAllowedOrigins(origins []string) ServerConfigOption {
	return func(c *restConfig, _ *runtime.ServeMux) {
		c.cors.allowedOrigins = origins
	}
}

// WithAllowedHeaders is an option to supply allowed headers in CORS.
func WithAllowedHeaders(headers []string) ServerConfigOption {
	return func(c *restConfig, _ *runtime.ServeMux) {
		c.cors.allowedHeaders = headers
	}
}

// WithAllowedMethods is an option to supply allowed methods in CORS.
func WithAllowedMethods(methods []string) ServerConfigOption {
	return func(c *restConfig, _ *runtime.ServeMux) {
		c.cors.allowedMethods = methods
	}
}

// WithAdditionalHandler is an option to add an additional handler func in the REST server.
func WithAdditionalHandler(method string, path string, h runtime.HandlerFunc) ServerConfigOption {
	return func(_ *restConfig, sm *runtime.ServeMux) {
		_ = sm.HandlePath(method, path, h)
	}
}

// WithAdditionalGRPCOpts is an option to add an additional gRPC dial options in the REST server communication to the
// backend.
func WithAdditionalGRPCOpts(opts []grpc.DialOption) ServerConfigOption {
	return func(c *restConfig, sm *runtime.ServeMux) {
		c.opts = append(c.opts, opts...)
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
	return func(c *restConfig, sm *runtime.ServeMux) {
		var publicURL string

		// Check if an embedded OAuth2 server public URL is given, otherwise take localhost.
		if config.EmbeddedOAuth2ServerPublicURLFlag == "" {
			publicURL = fmt.Sprintf("http://localhost:%d/v1/auth", httpPort)
		} else {
			publicURL = fmt.Sprintf("%s:%d/v1/auth", viper.GetString(config.EmbeddedOAuth2ServerPublicURLFlag), httpPort)
		}

		log.Infof("Using embedded OAuth2.0 server on %s", publicURL)

		// Configure the options for the embedded auth server
		opts = append(opts,
			oauth2.WithSigningKeysFunc(func() map[int]*ecdsa.PrivateKey {
				// Expand path, because this could contain ~
				path, err := util.ExpandPath(keyPath)
				if err != nil {
					// Just use the current working dir if it fails
					path = "."
				}

				return storage.LoadSigningKeys(path, keyPassword, saveOnCreate)
			}),
			oauth2.WithPublicURL(publicURL),
		)

		// Create a new embedded OAuth 2.0 server to serve as our auth server
		authSrv := oauth2.NewServer("", opts...)
		authHandler := func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			http.StripPrefix("/v1/auth", authSrv.Handler).ServeHTTP(w, r)
		}

		// Map specific paths in our REST server to our auth server
		WithAdditionalHandler("GET", "/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			authSrv.Handler.ServeHTTP(w, r)
		})(c, sm)
		WithAdditionalHandler("GET", "/v1/auth/certs", authHandler)(c, sm)
		WithAdditionalHandler("GET", "/v1/auth/login", authHandler)(c, sm)
		WithAdditionalHandler("GET", "/v1/auth/authorize", authHandler)(c, sm)
		WithAdditionalHandler("POST", "/v1/auth/login", authHandler)(c, sm)
		WithAdditionalHandler("POST", "/v1/auth/token", authHandler)(c, sm)
	}
}

// RunServer starts our REST API. The REST API is a reverse proxy using grpc-gateway that
// exposes certain gRPC calls as RESTful HTTP methods.
func RunServer(ctx context.Context, grpcPort uint16, port uint16, serverOpts ...ServerConfigOption) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	httpPort = port

	mux := runtime.NewServeMux()

	cnf.opts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	for _, o := range serverOpts {
		o(&cnf, mux)
	}

	if err := discovery.RegisterDiscoveryHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), cnf.opts); err != nil {
		return fmt.Errorf("failed to connect to discovery gRPC service %w", err)
	}

	if err := discovery.RegisterExperimentalDiscoveryHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), cnf.opts); err != nil {
		return fmt.Errorf("failed to connect to discovery gRPC service %w", err)
	}

	if err := assessment.RegisterAssessmentHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), cnf.opts); err != nil {
		return fmt.Errorf("failed to connect to assessment gRPC service %w", err)
	}

	if err := orchestrator.RegisterOrchestratorHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), cnf.opts); err != nil {
		return fmt.Errorf("failed to connect to orchestrator gRPC service %w", err)
	}

	if err := evaluation.RegisterEvaluationHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), cnf.opts); err != nil {
		return fmt.Errorf("failed to connect to evaluation gRPC service %w", err)
	}

	if err := evidence.RegisterEvidenceStoreHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), cnf.opts); err != nil {
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

	log.WithFields(logrus.Fields{
		"allowed-origins": cnf.cors.allowedOrigins,
		"allowed-methods": cnf.cors.allowedMethods,
		"allowed-headers": cnf.cors.allowedHeaders,
	}).Info("Applying CORS configuration...")

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
		if cnf.cors.OriginAllowed(origin) {
			// Set the appropriate access control header
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Vary", "Origin")

			// Additionally, we need to handle preflight (OPTIONS) requests to specify allowed headers and methods
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(cnf.cors.allowedHeaders, ","))
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(cnf.cors.allowedMethods, ","))
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
