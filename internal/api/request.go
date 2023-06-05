package api

import (
	"google.golang.org/protobuf/proto"
)

// PayloadRequest describes any kind of requests that carries a certain payload.
// Defined in internal/api to avoid cyclic dependencies.
type PayloadRequest interface {
	GetPayload() proto.Message
	proto.Message
}

// CloudServiceRequest represents any kind of RPC request, that contains a
// reference to a cloud service. Defined in internal/api to avoid cyclic
// dependencies.
type CloudServiceRequest interface {
	GetCloudServiceId() string
}
