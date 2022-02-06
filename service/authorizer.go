package service

import (
	"context"
	"fmt"

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

	client auth.AuthenticationClient
	conn   grpc.ClientConnInterface

	// GrpcOptions contains additional grpc dial options
	GrpcOptions []grpc.DialOption
	Username    string
	Password    string
}

func (i *InternalAuthorizer) init() (err error) {
	// TODO(oxisto): set flag depending on target url, insecure only for localhost
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(i),
	}

	if i.conn, err = grpc.Dial(i.Url, opts...); err != nil {
		return fmt.Errorf("could not connect: %w", err)
	}

	return nil
}

func (i *InternalAuthorizer) Token() (token *oauth2.Token, err error) {
	var resp *auth.LoginResponse

	if i.conn == nil {
		err = i.init()
		if err != nil {
			return nil, fmt.Errorf("could not initialize connection to auth service: %w", err)
		}
	}

	i.client = auth.NewAuthenticationClient(i.conn)
	resp, err = i.client.Login(context.TODO(), &auth.LoginRequest{
		Username: i.Username,
		Password: i.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("error while logging in: %w", err)
	}

	return &oauth2.Token{
		AccessToken: resp.Token,
	}, nil
}

func (i *InternalAuthorizer) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
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
