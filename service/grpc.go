package service

import (
	"fmt"
	"net"

	"clouditor.io/clouditor/api/orchestrator"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
)

type serverOption func(srv *grpc.Server)

func WithOrchestrator(service orchestrator.OrchestratorServer) serverOption {
	return func(srv *grpc.Server) {
		orchestrator.RegisterOrchestratorServer(srv, service)
	}
}

func StartGRPCServer(jwksURL string, opts ...serverOption) (sock net.Listener, srv *grpc.Server, err error) {
	var addr = ":0"

	// create a new socket for gRPC communication
	sock, err = net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, fmt.Errorf("could not listen: %w", err)
	}

	authConfig := ConfigureAuth(WithJWKSURL(jwksURL))

	// We also add our authentication middleware, because we usually add additional service later
	srv = grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_auth.UnaryServerInterceptor(authConfig.AuthFunc),
		),
	)

	for _, o := range opts {
		o(srv)
	}

	go func() {
		// serve the gRPC socket
		_ = srv.Serve(sock)
	}()

	return sock, srv, nil
}
