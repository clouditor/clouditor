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

// CertificationTargetRequest represents any kind of RPC request, that contains a
// reference to a target of evaluation. Defined in internal/api to avoid cyclic
// dependencies.
type CertificationTargetRequest interface {
	GetCertificationTargetId() string
}
