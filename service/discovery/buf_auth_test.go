package discovery

import (
	"context"
	"net"

	"clouditor.io/clouditor/api/auth"
	service_auth "clouditor.io/clouditor/service/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var bufAuthListener *bufconn.Listener

func bufAuthDialer(context.Context, string) (net.Conn, error) {
	return bufAuthListener.Dial()
}

func startBufAuthServer() (*grpc.Server, *service_auth.Service) {
	const bufSize = 1024 * 1024 * 2
	bufAuthListener = bufconn.Listen(bufSize)

	s := grpc.NewServer()
	authService := service_auth.NewService()
	auth.RegisterAuthenticationServer(s, authService)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	return s, authService
}
