package discovery

import (
	"context"
	"net"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/auth"
	service_assessment "clouditor.io/clouditor/service/assessment"
	service_auth "clouditor.io/clouditor/service/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const DefaultBufferSize = 1024 * 1024 * 2

var (
	bufAuthListener       *bufconn.Listener
	bufAssessmentListener *bufconn.Listener
)

func bufAuthDialer(context.Context, string) (net.Conn, error) {
	return bufAuthListener.Dial()
}

func bufAssessmentDialer(context.Context, string) (net.Conn, error) {
	return bufAssessmentListener.Dial()
}

// startBufAuthServer starts an auth service listening on a bufconn listener. It exposes
// real functionality of the auth server based on a internal in-memory connection for
// testing purposes.
func startBufAuthServer() (*grpc.Server, *service_auth.Service) {
	bufAuthListener = bufconn.Listen(DefaultBufferSize)

	server := grpc.NewServer()
	service := service_auth.NewService()
	auth.RegisterAuthenticationServer(server, service)

	go func() {
		if err := server.Serve(bufAuthListener); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	return server, service
}

// startBufAssessmentServer starts an assessment service listening on a bufconn listener. It exposes
// real functionality of the auth server based on a internal in-memory connection for
// testing purposes.
func startBufAssessmentServer(opts ...service_assessment.ServiceOption) (*grpc.Server, *service_assessment.Service) {
	bufAssessmentListener = bufconn.Listen(DefaultBufferSize)

	server := grpc.NewServer()
	service := service_assessment.NewService()
	assessment.RegisterAssessmentServer(server, service)

	go func() {
		if err := server.Serve(bufAssessmentListener); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	return server, service
}
