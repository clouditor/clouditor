package api

import (
	sync "sync"

	"connectrpc.com/connect"
)

type ConnectConnection[T any] struct {
	// BaseURL contains the URL used in establishing the connection. This cannot be changed after the first client call.
	BaseURL string

	// Opts contain options used in the client creation. Ideally, this should not be changed after the first client
	// call.
	Opts []connect.ClientOption

	// Client contains a connect client that is used to issue the actual RPCs.
	Client T

	// authorizer is the authorizer used in grpc.Dual. Ideally, this should not be changed after the first client call.
	authorizer Authorizer

	// m is a mutex that synchronizes access to the client conn.
	m sync.RWMutex
}

// SetAuthorizer implements UsesAuthorizer
func (conn *ConnectConnection[T]) SetAuthorizer(auth Authorizer) {
	conn.authorizer = auth
}

// Authorizer implements UsesAuthorizer
func (conn *ConnectConnection[T]) Authorizer() Authorizer {
	return conn.authorizer
}

type NewClientFuncType[T any] func(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) T

// NewConnectConnection creates a new [ConnectConnection] to the target using the specified function that creates a new
// client.
func NewConnectConnection[T any](newClientFunc NewClientFuncType[T], httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) *ConnectConnection[T] {
	conn := &ConnectConnection[T]{
		BaseURL: baseURL,
		Opts:    opts,
	}
	conn.Client = newClientFunc(httpClient, baseURL, opts...)

	return conn
}
