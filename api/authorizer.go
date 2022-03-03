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
	"fmt"

	"clouditor.io/clouditor/api/auth"
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

	// AuthURL contains the base URL of the authentication server this authorizer uses. This
	// can be a gRPC address (for example if the Clouditor internal authentication service is used) or
	// a regular HTTP(S) URL.
	AuthURL() string
}

// UsesAuthorizer is an interface to denote that a struct is willing to accept and use
// an Authorizer
type UsesAuthorizer interface {
	SetAuthorizer(auth Authorizer)
	Authorizer() Authorizer
}

// DefaultInternalAuthorizerAddress specifies the default gRPC address of the internal Clouditor auth service.
const DefaultInternalAuthorizerAddress = "localhost:9090"

// internalAuthorizer is an authorizer that uses OAuth 2.0 client credentials and does a OAuth client
// credentials flow
type oauthAuthorizer struct {
	// TODO(oxisto): check, if this is really needed
	authURL string

	*protectedToken
}

// internalAuthorizer is an authorizer that uses the Clouditor internal auth server (using gRPC) and
// does a login flow using username and password
type internalAuthorizer struct {
	// TODO(oxisto): check, if this is really needed
	authURL string

	*protectedToken
}

// baseAuthorizer contains fields that are shared by all authorizers
//
// TODO(oxisto): Rename to something that is more tied to the grpc per-request interface
type protectedToken struct {
	oauth2.TokenSource
}

// NewInternalAuthorizerFromPassword creates a new authorizer based on a (gRPC) URL of the
// authentication server and a username / password combination
func NewInternalAuthorizerFromPassword(url string, username string, password string, grpcOptions ...grpc.DialOption) Authorizer {
	var authorizer = &internalAuthorizer{
		authURL: url,
		protectedToken: &protectedToken{
			TokenSource: oauth2.ReuseTokenSource(nil, &internalTokenSource{
				authURL:     url,
				username:    username,
				password:    password,
				grpcOptions: grpcOptions,
			}),
		},
	}

	return authorizer
}

// NewInternalAuthorizerFromToken creates a new authorizer based on a (gRPC) URL of the
// authentication server and an oauth2.Token. It will attempt to refresh an expired access token,
// if a refresh token is supplied. Note, that this does a similar flow as OAuth 2.0, but it uses
// our internal direct gRPC connection to the authentication server, whereas OAuth 2.0 would use
// a POST request with application/x-www-form-urlencoded data.
func NewInternalAuthorizerFromToken(url string, token *oauth2.Token, grpcOptions ...grpc.DialOption) Authorizer {
	var refreshToken string

	if token != nil {
		refreshToken = token.RefreshToken
	}

	var authorizer = &internalAuthorizer{
		authURL: url,
		protectedToken: &protectedToken{
			TokenSource: oauth2.ReuseTokenSource(token, &internalTokenSource{
				authURL:      url,
				refreshToken: refreshToken,
				grpcOptions:  grpcOptions,
			}),
		},
	}

	return authorizer
}

// NewOAuthAuthorizerFromClientCredentials creates a new authorizer based on an OAuth 2.0 client credentials.
// It will attempt to refresh an expired access token, if a refresh token is supplied. Note, that this does a
// similar flow as OAuth 2.0, but it uses our internal direct gRPC connection to the authentication server,
// whereas OAuth 2.0 would use a POST request with application/x-www-form-urlencoded data.
func NewOAuthAuthorizerFromClientCredentials(config *clientcredentials.Config) Authorizer {
	var authorizer = &oauthAuthorizer{
		authURL: config.TokenURL,
		protectedToken: &protectedToken{
			TokenSource: oauth2.ReuseTokenSource(nil, config.TokenSource(context.Background())),
		},
	}

	return authorizer
}

// NewOAuthAuthorizerFromConfig creates a new authorizer based on an OAuth 2.0 config
func NewOAuthAuthorizerFromConfig(config *oauth2.Config, token *oauth2.Token) Authorizer {
	var authorizer = &oauthAuthorizer{
		authURL: config.Endpoint.AuthURL,
		protectedToken: &protectedToken{
			TokenSource: config.TokenSource(context.Background(), token),
		},
	}

	return authorizer
}

// GetRequestMetadata is an implementation for credentials.PerRPCCredentials. It is called before
// each RPC request and is used to inject our client credentials into the context of the RPC call.
func (p *protectedToken) GetRequestMetadata(ctx context.Context, _ ...string) (map[string]string, error) {
	// Fetch a token from our token source. This will also refresh an access token, if it has expired
	token, err := p.Token()
	if err != nil {
		return nil, err
	}

	if p.RequireTransportSecurity() {
		ri, _ := credentials.RequestInfoFromContext(ctx)
		if err = credentials.CheckSecurityLevel(ri.AuthInfo, credentials.PrivacyAndIntegrity); err != nil {
			return nil, fmt.Errorf("unable to transfer InternalAuthorizer PerRPCCredentials: %v", err)
		}
	}

	return map[string]string{
		"authorization": token.Type() + " " + token.AccessToken,
	}, nil
}

func (*protectedToken) RequireTransportSecurity() bool {
	// TODO(oxisto): This should be set to true because we transmit credentials (except localhost)
	return false
}

// AuthURL is an implementation needed for Authorizer. It returns the gRPC address of the authorization server.
func (i *internalAuthorizer) AuthURL() string {
	return i.authURL
}

type internalTokenSource struct {
	authURL string

	refreshToken string
	username     string
	password     string

	// grpcOptions contains additional grpc dial options
	grpcOptions []grpc.DialOption

	client auth.AuthenticationClient
	conn   grpc.ClientConnInterface
}

// init initializes the authorizer. This is called when fetching a token, if this authorizer has not been initialized.
func (i *internalTokenSource) init() (err error) {
	// Note, that we do NOT want any credentials.PerRPCCredentials dial option on this connection because
	// the API of the auth service is available without token authentication. Otherwise, a first login
	// would be impossible.
	//
	// TODO(oxisto): set transport credentials depending on target url, insecure only for localhost
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// Apply any extra gRPC options that we might have configured. The main use case for this is a unit test,
	// but clients might want to use this to tweak some values as well.
	opts = append(opts, i.grpcOptions...)

	if i.conn, err = grpc.Dial(i.authURL, opts...); err != nil {
		return fmt.Errorf("could not connect: %w", err)
	}

	// Store the client
	i.client = auth.NewAuthenticationClient(i.conn)

	return nil
}

// Token is an implementation for the interface oauth2.TokenSource so we can use this authorizer
// in (almost) the same way as any other OAuth 2.0 token endpoint. It will fetch an access token
// (and possibly a refresh token), if no token is stored yet in the authorizer or if the token is expired.
//
// The "refresh" of a token will either be done by a refresh token (if it exists) or by the stored combination
// of username / password.
func (i *internalTokenSource) Token() (token *oauth2.Token, err error) {
	var resp *auth.TokenResponse

	// We do a lazy initialization here, so the first request might take a little bit longer.
	// This might not be entirely thread-safe.
	if i.conn == nil {
		err = i.init()
		if err != nil {
			return nil, fmt.Errorf("could not initialize connection to auth service: %w", err)
		}
	}

	// Otherwise, we need to re-authenticate
	if i.refreshToken != "" {
		resp, err = i.client.Token(context.TODO(), &auth.TokenRequest{
			GrantType:    "refresh_token",
			RefreshToken: i.refreshToken,
		})
	} else {
		resp, err = i.client.Login(context.TODO(), &auth.LoginRequest{
			Username: i.username,
			Password: i.password,
		})
	}

	if err != nil {
		return nil, fmt.Errorf("could not fetch token via authentication client: %w", err)
	}

	token = &oauth2.Token{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		Expiry:       resp.GetExpiry().AsTime(),
		TokenType:    resp.TokenType,
	}
	i.refreshToken = resp.RefreshToken

	return
}

// AuthURL is an implementation needed for Authorizer. It returns the OAuth 2.0 token endpoint.
func (o *oauthAuthorizer) AuthURL() string {
	return o.authURL
}

// DefaultGrpcDialOptions returns a set of sensible default list of grpc.DialOption values. It includes
// transport credentials and configures per-RPC credentials using an authorizer, if one is configured.
func DefaultGrpcDialOptions(s UsesAuthorizer, additionalOpts ...grpc.DialOption) (opts []grpc.DialOption) {
	// TODO(oxisto): Enable TLS to external based on the URL (scheme)
	opts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

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
