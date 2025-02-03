// Copyright 2023 Fraunhofer AISEC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//           $$\                           $$\ $$\   $$\
//           $$ |                          $$ |\__|  $$ |
//  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
// $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
// $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
// $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
// \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
//  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
//
// This file is part of Clouditor Community Edition.

package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	sync "sync"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

// RPCConnection is a helper struct that wraps all necessary information for a gRPC connection, which is established
// using [grpc.Dial]. It features transparent goroutine-safe lazy initialization of the connection by overloading the
// underlying [grpc.ClientConn]. The connection is established automatically once the first client call is made. If an
// [io.EOF] error is received the connection is tried to be re-established on the next client call.
type RPCConnection[T any] struct {
	// Target contains the target used in grpc.Dial. Ideally, this should not be changed after the first client call.
	Target string

	// Opts contain options used in grpc.Dial. Ideally, this should not be changed after the first client call.
	Opts []grpc.DialOption

	// Client contains a gRPC client that is used to issue the actual RPCs.
	Client T

	// authorizer is the authorizer used in grpc.Dual. Ideally, this should not be changed after the first client call.
	authorizer Authorizer
	// ClientConn is an embedded grpc.ClientConn that we hook to automatically establish the connection.
	cc *grpc.ClientConn
	// m is a mutex that synchronizes access to the client conn.
	m sync.RWMutex
}

// SetAuthorizer implements UsesAuthorizer
func (conn *RPCConnection[T]) SetAuthorizer(auth Authorizer) {
	conn.authorizer = auth
}

// Authorizer implements UsesAuthorizer
func (conn *RPCConnection[T]) Authorizer() Authorizer {
	return conn.authorizer
}

// NewRPCConnection creates a new [RPCConnection] to the target using the specified function that creates a new client.
func NewRPCConnection[T any](target string, newClientFunc func(cc grpc.ClientConnInterface) T, opts ...grpc.DialOption) *RPCConnection[T] {
	conn := &RPCConnection[T]{
		Target: target,
		Opts:   opts,
	}
	conn.Client = newClientFunc(conn)

	return conn
}

// ForceReconnect drops the established gRPC client conn and forces a re-connect at the next client call.
func (conn *RPCConnection[T]) ForceReconnect() {
	conn.cc = nil
}

// Invoke implements [grpc.ClientConnInterface].
func (conn *RPCConnection[T]) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) (err error) {
	// Make sure, this connection is established
	err = conn.init()
	if err != nil {
		return
	}

	// Then, just forward the request to the embedded client conn
	err = conn.cc.Invoke(ctx, method, args, reply, opts...)
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		log.Debugf("Caught EOF while invoking method %s, forcing connection to reconnect on next call", method)
		conn.ForceReconnect()
	}

	return
}

// NewStream implements [grpc.ClientConnInterface].
func (conn *RPCConnection[T]) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (stream grpc.ClientStream, err error) {
	// Make sure, this connection is established
	err = conn.init()
	if err != nil {
		return
	}

	// Then, just forward the request to the embedded client conn
	stream, err = conn.cc.NewStream(ctx, desc, method, opts...)
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		log.Debugf("Caught EOF while invoking method %s, forcing connection to reconnect on next call", method)
		conn.ForceReconnect()
	}

	return
}

// init takes care of actually establishing the connection to the gRPC server. If the connection is already established,
// this is a no-op. This function is go-routine safe, because potentially multiple callers could access this at the same
// time.
func (conn *RPCConnection[T]) init() (err error) {
	if conn == nil {
		return errors.New("RPC connection not configured")
	}

	// Check, if we already have a valid client connection. We use a read-only lock for a faster access.
	conn.m.RLock()
	if conn.cc != nil && conn.cc.GetState() != connectivity.Shutdown {
		defer conn.m.RUnlock()
		return nil
	}

	// If we arrive at this point we need to exchange our read-only lock for a write lock, since we are creating a new
	// connection and will write to the client conn property.
	conn.m.RUnlock()
	conn.m.Lock()
	defer conn.m.Unlock()

	// Establish a connection to the specified gRPC service
	conn.cc, err = grpc.NewClient(conn.Target,
		DefaultGrpcDialOptions(conn.Target, conn, conn.Opts...)...,
	)
	if err != nil {
		return fmt.Errorf("could not connect to gPRC target %q: %w", conn.Target, err)
	}

	return nil
}
