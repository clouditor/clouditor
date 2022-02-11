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

package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"clouditor.io/clouditor/api/auth"
	"golang.org/x/oauth2"
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

// InternalAuthorizer is an authorizer that uses the Clouditor internal auth server (using gRPC) and
// does a login flow using username and password
type internalAuthorizer struct {
	authURL string

	// GrpcOptions contains additional grpc dial options
	grpcOptions []grpc.DialOption

	username string
	password string

	client auth.AuthenticationClient
	conn   grpc.ClientConnInterface

	// token contains the current token. It will be checked for expiry in the Token() method.
	token *oauth2.Token

	// tokenMutex is a Mutex that allowes concurrent access to our token
	tokenMutex sync.RWMutex
}

// NewInternalAuthorizerFromPassword creates a new authorizedbased on a (gRPC) URL of the
// authentication server and a username / password combination
func NewInternalAuthorizerFromPassword(url string, username string, password string, grpcOptions ...grpc.DialOption) Authorizer {
	return &internalAuthorizer{
		authURL:     url,
		username:    username,
		password:    password,
		grpcOptions: grpcOptions,
	}
}

// NewInternalAuthorizerFromToken creates a new authorizer based on a (gRPC) URL of the
// authentication server and an oauth2.Token. It will attempt to refresh an expired access token,
// if a refresh token is supplied. Note, that this does a similar flow as OAuth 2.0, but it uses
// our internal direct gRPC connection to the authentication server, whereas OAuth 2.0 would use
// a POST request with application/x-www-form-urlencoded data.
func NewInternalAuthorizerFromToken(url string, token *oauth2.Token, grpcOptions ...grpc.DialOption) Authorizer {
	return &internalAuthorizer{
		authURL:     url,
		token:       token,
		grpcOptions: grpcOptions,
	}
}

// init initializes the authorizer. This is called when fetching a token, if this authorizer has not been initialized.
func (i *internalAuthorizer) init() (err error) {
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
// using the stored username / password credentials and refresh the access token, if it is expired.
func (i *internalAuthorizer) Token() (*oauth2.Token, error) {
	var (
		resp *auth.TokenResponse
		err  error
	)

	// We do a lazy initialization here, so the first request might take a little bit longer.
	// This might not be entirely thread-safe.
	if i.conn == nil {
		err = i.init()
		if err != nil {
			return nil, fmt.Errorf("could not initialize connection to auth service: %w", err)
		}
	}

	// Lock the token for reading, so that we are sure to get the recent token,
	// if another call is currently fetching a new one and thus modifying it.
	i.tokenMutex.RLock()

	// Check if we already have a token or if it is still ok
	if i.token != nil && i.token.Expiry.After(time.Now()) {
		// Defer the unlock, if we return here
		defer i.tokenMutex.RUnlock()
		return i.token, nil
	}

	// Unlock a previous read, if we did not return. Otherwise we will block ourselves.
	i.tokenMutex.RUnlock()

	// Lock the token for modification
	i.tokenMutex.Lock()
	// Defer the Unlock so we definitely unlock once we exit this function
	defer i.tokenMutex.Unlock()

	// Otherwise, we need to re-authenticate
	// TODO(oxisto): In the future we should split this into two authorizers
	if i.token != nil && i.token.RefreshToken != "" {
		resp, err = i.client.Token(context.TODO(), &auth.TokenRequest{
			GrantType:    "refresh_token",
			RefreshToken: i.token.RefreshToken,
		})
	} else {
		resp, err = i.client.Login(context.TODO(), &auth.LoginRequest{
			Username: i.username,
			Password: i.password,
		})
	}
	if err != nil {
		// Return without refreshing the token. At this point, the token will still be invalid
		// and the next call to Token() will try again. This way we can mitigate temporary errors.
		return nil, fmt.Errorf("error while logging in: %w", err)
	}

	// Store the current token, if login was successful
	i.token = &oauth2.Token{
		AccessToken:  resp.AccessToken,
		Expiry:       resp.Expiry.AsTime(),
		TokenType:    resp.TokenType,
		RefreshToken: resp.RefreshToken,
	}

	return i.token, nil
}

// AuthURL is an implementation needed for Authorizer.
func (i *internalAuthorizer) AuthURL() string {
	return i.authURL
}

// GetRequestMetadata is an implementation for credentials.PerRPCCredentials. It is called before
// each RPC request and is used to inject our client credentials into the context of the RPC call.
func (i *internalAuthorizer) GetRequestMetadata(ctx context.Context, _ ...string) (map[string]string, error) {
	// Fetch a token from our token source. This will also refresh an access token, if it has expired
	token, err := i.Token()
	if err != nil {
		return nil, err
	}

	if i.RequireTransportSecurity() {
		ri, _ := credentials.RequestInfoFromContext(ctx)
		if err = credentials.CheckSecurityLevel(ri.AuthInfo, credentials.PrivacyAndIntegrity); err != nil {
			return nil, fmt.Errorf("unable to transfer InternalAuthorizer PerRPCCredentials: %v", err)
		}
	}

	return map[string]string{
		"authorization": token.Type() + " " + token.AccessToken,
	}, nil
}

func (*internalAuthorizer) RequireTransportSecurity() bool {
	// TODO(oxisto): This should be set to true because we transmit credentials (except localhost)
	return false
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
