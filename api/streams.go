// Copyright 2022 Fraunhofer AISEC
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
	"errors"
	"fmt"
	"io"
	"sync"

	"clouditor.io/clouditor/v2/internal/logging"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	ErrMissingInitFunc = errors.New("missing stream initializer function")
)

// IDField is the name of the ID field used in protobuf messages
const IDField = "id"

// StreamChannelOf provides a channel around a connection to a grpc.ClientStream to send messages of type MsgType to
// that particular stream, using an internal go routine. This is necessary, because gRPC does not allow sending to a
// stream from multiple goroutines directly.
type StreamChannelOf[StreamType grpc.ClientStream, MsgType proto.Message] struct {
	// channel can be used to send a message to the stream
	channel chan MsgType

	// stream to the component
	stream StreamType

	// target of the component (host and port usually)
	target string

	// component name
	component string

	// dead specifies that this channel lost connection and needs to be re-started
	dead bool
}

// InitFuncOf describes a function with type parameters that creates any kind of stream towards a gRPC server specified
// in target and returns the stream or an error. Additional gRPC dial options can be specified in additionalOpts.
type InitFuncOf[StreamType grpc.ClientStream] func(target string, additionalOpts ...grpc.DialOption) (stream StreamType, err error)

// StreamsOf handles stream channels to multiple gRPC servers, identified by a unique target (usually host and port).
// Since gRPC does only allow to send to a stream using one goroutine, each stream provides a go channel that can be
// used to send messages to the particular stream.
//
// A stream for a given target can be retrieved with the GetStream function, which automatically initializes the stream
// if it does not exist.
type StreamsOf[StreamType grpc.ClientStream, MsgType proto.Message] struct {
	mutex    sync.RWMutex
	channels map[string]*StreamChannelOf[StreamType, MsgType]
	log      *logrus.Entry
}

// StreamsOfOption is a functional option type to configure the StreamOf type.
type StreamsOfOption[StreamType grpc.ClientStream, MsgType proto.Message] func(*StreamsOf[StreamType, MsgType])

// WithLogger can be used to specify a dedicated logger entry which is used for logging. Otherwise, the default logging
// entry of logrus is used.
func WithLogger[StreamType grpc.ClientStream, MsgType proto.Message](log *logrus.Entry) StreamsOfOption[StreamType, MsgType] {
	return func(s *StreamsOf[StreamType, MsgType]) {
		s.log = log
	}
}

// NewStreamsOf creates a new StreamsOf object and initializes all the necessary objects for it.
func NewStreamsOf[StreamType grpc.ClientStream, MsgType proto.Message](opts ...StreamsOfOption[StreamType, MsgType]) (s *StreamsOf[StreamType, MsgType]) {
	s = &StreamsOf[StreamType, MsgType]{
		channels: map[string]*StreamChannelOf[StreamType, MsgType]{},
	}

	// Apply options
	for _, o := range opts {
		o(s)
	}

	// Default to a default logger
	if s.log == nil {
		s.log = defaultLog()
	}

	return s
}

// GetStream tries to retrieve a stream for the given target and component. If no stream exists, it tries to
// create a new stream using the supplied init function. An error is returned if the initialization is not
// successful.
func (s *StreamsOf[StreamType, MsgType]) GetStream(target string, component string, init InitFuncOf[StreamType], opts ...grpc.DialOption) (c *StreamChannelOf[StreamType, MsgType], err error) {
	var (
		ok bool
	)

	// Try to retrieve the stream, given the target. We can RLock here, because we only need read access.
	s.mutex.RLock()
	c, ok = s.channels[target]
	s.mutex.RUnlock()

	// No stream found, let's try to add one
	if !ok {
		c, err = s.addStream(target, component, init, opts...)
		if err != nil {
			return nil, fmt.Errorf("could not add stream for %s with target '%s': %w", component, target, err)
		}
	} else if c.dead {
		// We could have a dead stream that we need to restart. in this case, we can recycle a few things, e.g. the channel
		c, err = s.restartStream(c, init, opts...)
		if err != nil {
			return nil, fmt.Errorf("could not restart stream for %s with target '%s': %w", component, target, err)
		}
	}

	return c, nil
}

// CloseAll closes all streams
func (s *StreamsOf[StreamType, MsgType]) CloseAll() {
	for _, channel := range s.channels {
		_ = channel.stream.CloseSend()
	}
}

// addStream stores a stream to the given component and starts a goroutine for sending messages from the channel to the given component
func (s *StreamsOf[StreamType, MsgType]) addStream(target string, component string, init InitFuncOf[StreamType], opts ...grpc.DialOption) (c *StreamChannelOf[StreamType, MsgType], err error) {
	// We need an init func
	if init == nil {
		return nil, ErrMissingInitFunc
	}

	// Initialize the stream using our init function
	stream, err := init(target, opts...)
	if err != nil {
		return nil, fmt.Errorf("could not init stream: %w", err)
	}

	// Create our stream channel struct
	c = &StreamChannelOf[StreamType, MsgType]{
		stream:    stream,
		component: component,
		target:    target,
		channel:   make(chan MsgType, 1000),
	}

	// Update the stream map. This time we need a real lock for an update
	s.mutex.Lock()
	s.channels[target] = c
	s.mutex.Unlock()

	s.log.Infof("Established stream to %s (%s)", component, target)

	// Start go routine for receiving messages from the stream (especially relevant for bi-directional streams).
	go c.recvLoop(s)

	// Start go routine for sending messages from the channel to the stream
	go c.sendLoop(s)

	return c, nil
}

// restartStream restarts a stream to the given component and starts a goroutine for sending messages from the channel to the given component
func (s *StreamsOf[StreamType, MsgType]) restartStream(c *StreamChannelOf[StreamType, MsgType], init InitFuncOf[StreamType], opts ...grpc.DialOption) (*StreamChannelOf[StreamType, MsgType], error) {
	var err error

	// We need an init func
	if init == nil {
		return nil, ErrMissingInitFunc
	}

	// Initialize the stream using our init function
	c.stream, err = init(c.target, opts...)
	if err != nil {
		return nil, fmt.Errorf("could not init stream: %w", err)
	}

	// Revive the stream
	c.dead = false

	s.log.Infof("Re-Established stream to %s (%s)", c.component, c.target)

	// Start go routine for receiving messages from the stream (especially relevant for bi-directional streams).
	go c.recvLoop(s)

	// Start go routine for sending messages from the channel to the stream
	go c.sendLoop(s)

	return c, nil
}

// sendLoop continuously fetches new messages from the channel inside c and sends them to the appropriate stream.
func (c *StreamChannelOf[StreamType, MsgType]) sendLoop(s *StreamsOf[StreamType, MsgType]) {
	var err error

	// Fetch new (or queued old) messages from the channel. This will block.
	for m := range c.channel {
		// We want to log some additional information about this stream and its
		// payload. The logging functions are safe to call with a nil request,
		// so we can avoid checking, if this succeeds
		preq, _ := any(m).(PayloadRequest)

		// Try to send the message in our stream
		err = c.stream.SendMsg(m)
		if err != nil {
			if errors.Is(err, io.EOF) {
				s.log.Infof("Stream to %s (%s) closed with EOF", c.component, c.target)
			} else {
				// Some other error than EOF occurred
				s.log.Errorf("Error when sending message to %s (%s): %v", c.component, c.target, err)

				// Close the stream gracefully. We can ignore any error resulting from the close here
				_ = c.stream.CloseSend()
			}

			// Declare the stream as dead
			c.dead = true

			// Put the message back on the channel, so that it does not get lost
			go func() {
				logging.LogRequest(s.log, logrus.DebugLevel, logging.Store, preq, fmt.Sprintf("back into queue for %s (%s)", c.component, c.target))
				c.channel <- m
			}()
			return
		}

		logging.LogRequest(s.log, logrus.DebugLevel, logging.Send, preq, fmt.Sprintf("to %s (%s)", c.component, c.target))
	}
}

// recvLoop continuously receives message from the stream. Currently, they are just discarded. In the future, we might
// want to send them back to the caller. But we need to receive them, otherwise the buffer of the stream gets congested.
func (c *StreamChannelOf[StreamType, MsgType]) recvLoop(s *StreamsOf[StreamType, MsgType]) {
	for {
		// TODO(oxisto): Check, if this also works for uni-directional streams
		// emptypb.Empty is used for now to give a correctly typed message to RecvMsg. In the future, use
		// types of response message of respective RPCs.

		msg := new(emptypb.Empty)
		err := c.stream.RecvMsg(msg)

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			s.log.Errorf("Error receiving response from stream: %v", err)
			break
		}
	}
}

// Send sends the message into the stream via the channel. Since this uses the receive operator on the channel,
// this function may block until the message is received on the sendLoop of this StreamChannelOf or if
// the buffer of the channel is full.
func (c *StreamChannelOf[StreamType, MsgType]) Send(msg MsgType) {
	c.channel <- msg
}

// defaultLog returns the default logger, if none is specified.
func defaultLog() *logrus.Entry {
	return logrus.NewEntry(logrus.StandardLogger())
}
