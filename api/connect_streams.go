package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	sync "sync"

	"clouditor.io/clouditor/v2/internal/logging"
	"connectrpc.com/connect"
	"github.com/sirupsen/logrus"
)

// ConnectInitFuncOf describes a function with type parameters that creates any kind of stream towards a gRPC server specified
// in target and returns the stream or an error. Additional gRPC dial options can be specified in additionalOpts.
type ConnectInitFuncOf[Client any, Req any, Res any] func(client Client, ctx context.Context) (stream *connect.BidiStreamForClient[Req, Res])

// ConnectStreamChannelOf provides a channel around a connection to a grpc.ClientStream to send messages of type MsgType to
// that particular stream, using an internal go routine. This is necessary, because gRPC does not allow sending to a
// stream from multiple goroutines directly.
type ConnectStreamChannelOf[Req any, Res any] struct {
	// channel can be used to send a message to the stream
	channel chan *Req

	// stream to the component
	stream *connect.BidiStreamForClient[Req, Res]

	// target of the component (host and port usually)
	target string

	// component name
	component string

	// dead specifies that this channel lost connection and needs to be re-started
	dead bool
}

// ConnectStreamsOf handles stream channels to multiple gRPC servers, identified by a unique target (usually host and port).
// Since gRPC does only allow to send to a stream using one goroutine, each stream provides a go channel that can be
// used to send messages to the particular stream.
//
// A stream for a given target can be retrieved with the GetStream function, which automatically initializes the stream
// if it does not exist.
type ConnectStreamsOf[Client comparable, Req any, Res any] struct {
	mutex    sync.RWMutex
	channels map[Client]*ConnectStreamChannelOf[Req, Res]
	initFunc ConnectInitFuncOf[Client, Req, Res]
	log      *logrus.Entry
}

// NewConnectStreamsOf creates a new ConnectStreamsOf object and initializes all the necessary objects for it.
func NewConnectStreamsOf[Client comparable, Req any, Res any](initFunc ConnectInitFuncOf[Client, Req, Res]) (s *ConnectStreamsOf[Client, Req, Res]) {
	s = &ConnectStreamsOf[Client, Req, Res]{
		channels: map[Client]*ConnectStreamChannelOf[Req, Res]{},
		initFunc: initFunc,
	}

	// Apply options
	/*for _, o := range opts {
		o(s)
	}*/

	// Default to a default logger
	if s.log == nil {
		s.log = defaultLog()
	}

	return s
}

// GetStream tries to retrieve a stream for the given target and component. If no stream exists, it tries to create a
// new stream using the supplied init function. An error is returned if the initialization is not successful.
func (s *ConnectStreamsOf[Client, Req, Res]) GetStream(client Client, component string /*, opts ...grpc.DialOption*/) (c *ConnectStreamChannelOf[Req, Res], err error) {
	var (
		ok bool
	)

	// Try to retrieve the stream, given the target. We can RLock here, because we only need read access.
	s.mutex.RLock()
	c, ok = s.channels[client]
	s.mutex.RUnlock()

	// TODO
	target := ""

	// No stream found, let's try to add one
	if !ok {
		c, err = s.addStream(client, target, component)
		if err != nil {
			return nil, fmt.Errorf("could not add stream for %s with target '%s': %w", component, target, err)
		}
	} else if c.dead {
		// We could have a dead stream that we need to restart. in this case, we can recycle a few things, e.g. the
		// channel
		c, err = s.restartStream(client, c)
		if err != nil {
			return nil, fmt.Errorf("could not restart stream for %s with target '%s': %w", c.component, c.target, err)
		}
	}

	return c, nil
}

// CloseAll closes all streams
func (s *ConnectStreamsOf[Client, Req, Res]) CloseAll() {
	for _, channel := range s.channels {
		_ = channel.stream.CloseRequest()
		_ = channel.stream.CloseResponse()
	}
}

// addStream stores a stream to the given component and starts a goroutine for sending messages from the channel to the given component
func (s *ConnectStreamsOf[Client, Req, Res]) addStream(client Client, target, component string) (c *ConnectStreamChannelOf[Req, Res], err error) {
	// Create a new stream using our init func
	stream := s.initFunc(client, context.TODO())

	// Create our stream channel struct
	c = &ConnectStreamChannelOf[Req, Res]{
		stream:    stream,
		target:    target,
		component: component,
		channel:   make(chan *Req, 1000),
	}

	// Update the stream map. This time we need a real lock for an update
	s.mutex.Lock()
	s.channels[client] = c
	s.mutex.Unlock()

	s.log.Infof("Established stream to %s (%s)", component, target)

	// Start go routine for receiving messages from the stream (especially relevant for bi-directional streams).
	go c.recvLoop(s.log)

	// Start go routine for sending messages from the channel to the stream
	go c.sendLoop(s.log)

	return c, nil
}

// restartStream restarts a stream to the given component and starts a goroutine for sending messages from the channel
// to the given component
func (s *ConnectStreamsOf[Client, Req, Res]) restartStream(client Client, c *ConnectStreamChannelOf[Req, Res]) (*ConnectStreamChannelOf[Req, Res], error) {
	// Create a new stream using our init func
	c.stream = s.initFunc(client, context.TODO())

	// Revive the stream
	c.dead = false

	s.log.Infof("Re-Established stream to %s (%s)", c.component, c.target)

	// Start go routine for receiving messages from the stream (especially relevant for bi-directional streams).
	go c.recvLoop(s.log)

	// Start go routine for sending messages from the channel to the stream
	go c.sendLoop(s.log)

	return c, nil
}

// sendLoop continuously fetches new messages from the channel inside c and sends them to the appropriate stream.
func (c *ConnectStreamChannelOf[Req, Res]) sendLoop(log *logrus.Entry) {
	var err error

	// Fetch new (or queued old) messages from the channel. This will block.
	for m := range c.channel {
		// We want to log some additional information about this stream and its
		// payload. The logging functions are safe to call with a nil request,
		// so we can avoid checking, if this succeeds
		preq, _ := any(m).(PayloadRequest)

		// Try to send the message in our stream
		err = c.stream.Send(m)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Infof("Stream to %s (%s) closed with EOF", c.component, c.target)
			} else {
				// Some other error than EOF occurred
				log.Errorf("Error when sending message to %s (%s): %v", c.component, c.target, err)

				// Close the stream gracefully. We can ignore any error resulting from the close here
				_ = c.stream.CloseRequest()
			}

			// Declare the stream as dead
			c.dead = true

			// Put the message back on the channel, so that it does not get lost
			go func() {
				logging.LogRequest(log, logrus.DebugLevel, logging.Store, preq, fmt.Sprintf("back into queue for %s (%s)", c.component, c.target))
				c.channel <- m
			}()
			return
		}

		logging.LogRequest(log, logrus.DebugLevel, logging.Send, preq, fmt.Sprintf("to %s (%s)", c.component, c.target))
	}
}

// recvLoop continuously receives message from the stream. Currently, they are just discarded. In the future, we might
// want to send them back to the caller. But we need to receive them, otherwise the buffer of the stream gets congested.
func (c *ConnectStreamChannelOf[Req, Res]) recvLoop(log *logrus.Entry) {
	for {
		// TODO(oxisto): Check, if this also works for uni-directional streams
		// emptypb.Empty is used for now to give a correctly typed message to RecvMsg. In the future, use
		// types of response message of respective RPCs.
		msg, err := c.stream.Receive()
		if errors.Is(err, io.EOF) {
			break
		}

		// For now we discard the response
		_ = msg

		if err != nil {
			log.Errorf("Error receiving response from stream: %v", err)
			break
		}
	}
}

// Send sends the message into the stream via the channel. Since this uses the receive operator on the channel,
// this function may block until the message is received on the sendLoop of this StreamChannelOf or if
// the buffer of the channel is full.
func (c *ConnectStreamChannelOf[Req, Res]) Send(msg *Req) {
	c.channel <- msg
}
