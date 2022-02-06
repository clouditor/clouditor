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

type Authorizer interface {
	credentials.PerRPCCredentials
	oauth2.TokenSource
}

// InternalAuthorizer is an authorizer that uses the Clouditor internal auth server (using gRPC) and
// does a login flow using username and password
type InternalAuthorizer struct {
	Url string

	// GrpcOptions contains additional grpc dial options
	GrpcOptions []grpc.DialOption
	Username    string
	Password    string

	client auth.AuthenticationClient
	conn   grpc.ClientConnInterface

	// token contains the current token. It will be checked for expiry in the Token() method.
	token *oauth2.Token

	// tokenMutex is a Mutex that allowes concurrent access to our token
	tokenMutex sync.RWMutex
}

func (i *InternalAuthorizer) init() (err error) {
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
	opts = append(opts, i.GrpcOptions...)

	if i.conn, err = grpc.Dial(i.Url, opts...); err != nil {
		return fmt.Errorf("could not connect: %w", err)
	}

	// Store the client
	i.client = auth.NewAuthenticationClient(i.conn)

	return nil
}

// Token is an implementation for the interface oauth2.TokenSource so we can use this authorizer
// in (almost) the same way as any other OAuth 2.0 token endpoint. It will fetch an access token
// using the stored username / password credentials and refresh the access token, if it is expired.
func (i *InternalAuthorizer) Token() (token *oauth2.Token, err error) {
	var resp *auth.LoginResponse

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
	if token != nil && token.Expiry.After(time.Now()) {
		// Defer the unlock, if we return here
		defer i.tokenMutex.RUnlock()
		return token, nil
	}

	// Unlock a previous read, if we did not return. Otherwise we will block ourselves.
	i.tokenMutex.RUnlock()

	// Lock the token for modification
	i.tokenMutex.Lock()
	// Defer the Unlock so we definitely unlock once we exit this function
	defer i.tokenMutex.Unlock()

	// Otherwise, we need to re-authenticate
	resp, err = i.client.Login(context.TODO(), &auth.LoginRequest{
		Username: i.Username,
		Password: i.Password,
	})
	if err != nil {
		// Return without refreshing the token. At this point, the token will still be invalid
		// and the next call to Token() will try again. This way we can mitigate temporary errors.
		return nil, fmt.Errorf("error while logging in: %w", err)
	}

	// Store the current token, if login was successful
	token = &oauth2.Token{
		AccessToken: resp.AccessToken,
		Expiry:      resp.Expiry.AsTime(),
	}

	return token, nil
}

// GetRequestMetadata is an implementation for credentials.PerRPCCredentials. It is called before
// each RPC request and is used to inject our client credentials into the context of the RPC call.
func (i *InternalAuthorizer) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	// Fetch a token from our token source. This will also refresh an access token, if it has expired
	token, err := i.Token()
	if err != nil {
		return nil, err
	}

	ri, _ := credentials.RequestInfoFromContext(ctx)
	if err = credentials.CheckSecurityLevel(ri.AuthInfo, credentials.PrivacyAndIntegrity); err != nil {
		return nil, fmt.Errorf("unable to transfer InternalAuthorizer PerRPCCredentials: %v", err)
	}

	return map[string]string{
		"authorization": token.Type() + " " + token.AccessToken,
	}, nil
}

func (i *InternalAuthorizer) RequireTransportSecurity() bool {
	// TODO(oxisto): This should be set to true because we transmit credentials (except localhost)
	return false
}
