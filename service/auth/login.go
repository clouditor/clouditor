package auth

import (
	"context"

	"clouditor.io/clouditor"
)

//go:generate protoc -I ../../proto -I ../../third_party auth.proto --go_out=../.. --go-grpc_out=../..

type Service struct {
	clouditor.UnimplementedAuthenticationServer
}

func (s Service) Login(ctx context.Context, in *clouditor.LoginRequest) (response *clouditor.LoginResponse, err error) {
	response = &clouditor.LoginResponse{Token: "my-token"}

	return response, nil
}
