package service

import (
	"errors"
	"io"
	"strings"

	"google.golang.org/grpc"
)

// StreamRequester is an interface that can be implemented by a struct
// to wrap it into a suitable message to be transported into stream
type StreamRequester[T any] interface {
	StreamRequest() *T
	GetId() string
}

// StreamSendLoop takes care of reading from a channel and sends the content of the
// channel into a grpc.ClientStream. The main use case is to consolidate the sending
// of messages into a gRPC into one goroutine, because sending from different goroutines
// is not allowed.
func StreamSendLoop[T StreamRequester[S], S any](channel chan T, stream grpc.ClientStream, typ string, target string) {
	for msg := range channel {
		err := stream.SendMsg(msg.StreamRequest())
		if errors.Is(err, io.EOF) {
			log.Infof("Stream to %s was closed", target)
			break
		}

		if err != nil {
			log.Errorf("Error when sending %s to %s: %v", err, typ, target)
			break
		}

		log.Debugf("%s (%v) sent to %s", strings.ToUpper(typ), msg.GetId(), target)
	}
}
