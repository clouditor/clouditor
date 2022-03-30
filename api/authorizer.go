// Copyright 2022 Fraunhofer AISEC
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

package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// Authorizer represents an interface which provides a token used for authenticating a client in server-client communication.
// More specifically, this interfaces requires credentials.PerRPCCredentials, which enables this to be used by a gRPC client
// to communicate with a gRPC server that requires per-RPC credentials.
type Authorizer interface {
	credentials.PerRPCCredentials
	oauth2.TokenSource
}

// UsesAuthorizer is an interface to denote that a struct is willing to accept and use
// an Authorizer
type UsesAuthorizer interface {
	SetAuthorizer(auth Authorizer)
	Authorizer() Authorizer
}

// oauthAuthorizer is an authorizer that uses OAuth 2.0 client credentials and does a OAuth client
// credentials flow
type oauthAuthorizer struct {
	oauth2.TokenSource
}

// NewOAuthAuthorizerFromClientCredentials creates a new authorizer based on an OAuth 2.0 client credentials.
func NewOAuthAuthorizerFromClientCredentials(config *clientcredentials.Config) Authorizer {
	var authorizer = &oauthAuthorizer{
		TokenSource: oauth2.ReuseTokenSource(nil, config.TokenSource(context.Background())),
	}

	return authorizer
}

// NewOAuthAuthorizerFromConfig creates a new authorizer based on an OAuth 2.0 config.
func NewOAuthAuthorizerFromConfig(config *oauth2.Config, token *oauth2.Token) Authorizer {
	var authorizer = &oauthAuthorizer{
		TokenSource: config.TokenSource(context.Background(), token),
	}

	return authorizer
}

// GetRequestMetadata is an implementation for credentials.PerRPCCredentials. It is called before
// each RPC request and is used to inject our client credentials into the context of the RPC call.
func (p *oauthAuthorizer) GetRequestMetadata(ctx context.Context, _ ...string) (map[string]string, error) {
	// Fetch a token from our token source. This will also refresh an access token, if it has expired
	token, err := p.Token()
	if err != nil {
		return nil, err
	}

	if p.RequireTransportSecurity() {
		ri, _ := credentials.RequestInfoFromContext(ctx)
		if err = credentials.CheckSecurityLevel(ri.AuthInfo, credentials.PrivacyAndIntegrity); err != nil {
			return nil, fmt.Errorf("unable to transfer OAuthAuthorizer PerRPCCredentials: %v", err)
		}
	}

	return map[string]string{
		"authorization": token.Type() + " " + token.AccessToken,
	}, nil
}

func (*oauthAuthorizer) RequireTransportSecurity() bool {
	// TODO(oxisto): This should be set to true because we transmit credentials (except localhost)
	return false
}

// DefaultGrpcDialOptions returns a set of sensible default list of grpc.DialOption values. It includes
// transport credentials and configures per-RPC credentials using an authorizer, if one is configured.
func DefaultGrpcDialOptions(hostport string, s UsesAuthorizer, additionalOpts ...grpc.DialOption) (opts []grpc.DialOption) {
	var (
		port string
		host string
		err  error
	)

	host, port, err = net.SplitHostPort(hostport)

	// TODO(oxisto): make a better distinction, for now this is ok
	if err == nil && port == "443" {
		// Use default TLS configuration using the system cert store
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			ServerName:         host,
			MinVersion:         12,
			InsecureSkipVerify: false,
		})))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// In practice, we should always have an authorizer, so we could fail early here. However,
	// if the server-side has not enabled the auth middleware (for example in testing), it is perfectly
	// fine to at least attempt to run it without one. If the server-side has enabled auth middleware
	// and does not receive any client credentials, the actual RPC call will then fail later.
	authorizer := s.Authorizer()

	if authorizer != nil {
		opts = append(opts, grpc.WithPerRPCCredentials(authorizer))
	}

	// Appply any additional options that we might have
	opts = append(opts, additionalOpts...)

	return opts
}
