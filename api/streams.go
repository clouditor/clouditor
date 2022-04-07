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

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

var (
	ErrMissingInitFunc = errors.New("missing stream initializer function")
)

var (
	log = logrus.WithField("component", "api-streams")
)

// StreamChannelOf provides a channel around a connection to a grpc.ClientStream to send messages of type MsgType to
// that particular stream, using an internal go routine. This is necessary, because gRPC does not allow to send to a
// stream from multiple goroutines directly.
type StreamChannelOf[StreamType grpc.ClientStream, MsgType proto.Message] struct {
	// Channel can be used to send message to the stream
	Channel chan MsgType

	// stream to the component
	stream StreamType

	// url of the component
	url string

	// component name
	component string
}

// InitFuncOf describes a function with type paramters that creates any kind of stream towards a gRPC server specified
// in URL and returns the stream or an error. Additional gRPC dial options can be specified in additionalOpts.
type InitFuncOf[StreamType grpc.ClientStream] func(URL string, additionalOpts ...grpc.DialOption) (stream StreamType, err error)

// StreamsOf handles stream channels to multiple gRPC servers, idenfitied by a unique URL. Since gRPC does only allow to
// send to a stream using one goroutine, each stream provides a go channel that can be used to send messages to the
// particular stream.
//
// A stream for a given URL can be retrieved with the GetStream function, which automatically initializes the stream if
// it does not exist.
type StreamsOf[StreamType grpc.ClientStream, MsgType proto.Message] struct {
	mutex    sync.RWMutex
	channels map[string]*StreamChannelOf[StreamType, MsgType]
}

// NewStreamsOf creates a new StreamsOf object and initializes all the necessary objects for it.
func NewStreamsOf[StreamType grpc.ClientStream, MsgType proto.Message]() *StreamsOf[StreamType, MsgType] {
	return &StreamsOf[StreamType, MsgType]{
		channels: map[string]*StreamChannelOf[StreamType, MsgType]{},
	}
}

// GetStream tries to retrieve a stream for the given URL and component. If no stream exists, it tries to
// create a new stream using the supplied init function. An error is returned if the initialization is not
// successful.
func (s *StreamsOf[StreamType, MsgType]) GetStream(URL string, component string, init InitFuncOf[StreamType], opts ...grpc.DialOption) (c *StreamChannelOf[StreamType, MsgType], err error) {
	var (
		ok bool
	)

	// Try to retrieve the stream, given the URL. We can RLock here, because we only need read access.
	s.mutex.RLock()
	c, ok = s.channels[URL]
	s.mutex.RUnlock()

	// No stream found, let's try to add one
	if !ok {
		c, err = s.addStream(URL, component, init, opts...)
		if err != nil {
			return nil, fmt.Errorf("could not add stream for %s with URL '%s': %w", component, URL, err)
		}
	}

	return c, nil
}

// addStream stores a stream to the given component and starts a goroutine for sending messages from the channel to the given component
func (s *StreamsOf[StreamType, MsgType]) addStream(URL string, component string, init InitFuncOf[StreamType], opts ...grpc.DialOption) (c *StreamChannelOf[StreamType, MsgType], err error) {
	// We need an init func
	if init == nil {
		return nil, ErrMissingInitFunc
	}

	// Initialize the stream using our init function
	stream, err := init(URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("could not init stream: %w", err)
	}

	// Create our stream channel struct
	c = &StreamChannelOf[StreamType, MsgType]{
		stream:    stream,
		component: component,
		url:       URL,
		Channel:   make(chan MsgType, 1000),
	}

	// Update the stream map. This time we need a real lock for an update
	s.mutex.Lock()
	s.channels[URL] = c
	s.mutex.Unlock()

	log.Infof("Established stream to %s (%s)", component, URL)

	// Start go routine for receiving messages from the stream (especially relevant for bi-directional streams).
	go c.recvLoop()

	// Start go routine for sending messages from the channel to the stream
	go c.sendLoop(s)

	return c, nil
}

// removeStream deletes the channel from the stream map.
func (s *StreamsOf[StreamType, MsgType]) removeStream(URL string) {
	log.Debugf("Removing stream channel for URL %s", URL)

	s.mutex.Lock()
	delete(s.channels, URL)
	s.mutex.Unlock()
}

// sendLoop continuously fetches new messages from the channel inside c and sends them to the appropriate stream.
func (c *StreamChannelOf[StreamType, MsgType]) sendLoop(s *StreamsOf[StreamType, MsgType]) {
	var err error

	// Fetch new messages from channel (this will block)
	for e := range c.Channel {
		// Try to send the message in our stream
		err = c.stream.SendMsg(e)
		if errors.Is(err, io.EOF) {
			log.Infof("Stream to %s (%s) closed with EOF", c.component, c.url)

			// Remove the stream from our map and end this goroutine
			s.removeStream(c.url)
			return
		}

		// Some other error than EOF occured
		if err != nil {
			log.Errorf("Error when sending message to %s (%s): %v", c.component, c.url, err)

			// Close the stream gracefully. We can ignore any error resulting from the close here
			_ = c.stream.CloseSend()

			// Remove the stream from our map and end this goroutine
			s.removeStream(c.url)
			return
		}

		log.Debugf("Message sent to %s (%s)", c.component, c.url)
	}
}

// recvLoop continously receives message from the stream. Currently they are just discarded. In the future, we might
// want to send them back to the caller. But we need to receive them, otherwise the buffer of the stream gets congested.
func (c *StreamChannelOf[StreamType, MsgType]) recvLoop() {
	for {
		// TODO(oxisto): Check, if this also works for uni-directional streams
		var msg proto.Message
		err := c.stream.RecvMsg(&msg)

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			log.Errorf("Error receiving response from stream: %v", err)
			break
		}
	}
}
