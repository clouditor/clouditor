// Maybe rename file to make more anambigous to the services
package aws

import (
	pb "clouditor.io/clouditor/api/discovery"
	"context"
	"flag"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	tls  = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	port = flag.Int("port", 10000, "The server port")
)

type discoveryServer struct {
	pb.UnimplementedDiscoveryServer
	isActive      pb.StartDiscoveryResponse
	queryResponse pb.QueryResponse
}

func (s *discoveryServer) Query(ctx context.Context) (*pb.QueryResponse, error) {
	listOfValues := []*structpb.Value{structpb.NewStringValue("Value 1"), structpb.NewStringValue("Value 2")}
	result := structpb.ListValue{Values: listOfValues}
	return &pb.QueryResponse{Result: &result}, nil
}

//func startServer() {
//	// flag.Parse()
//	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
//	if err != nil {
//		log.Fatalf("failed to listen: %v", err)
//	}
//	//var opts []grpc.ServerOption
//	// if *tls...
//	grpcServer := grpc.NewServer()
//	pb.RegisterDiscoveryServer(grpcServer, newServer())
//
//}
//
//func newServer() *discoveryServer {
//	s := &discoveryServer{}
//	return s
//}
