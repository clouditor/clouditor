package testutil

import (
	"context"
	"io"
	"slices"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

// ListRecvStreamerOf implements a grpc.ClientStream that mocks the reception of a list of messages
// and then EOFs.
type ListRecvStreamerOf[MsgType proto.Message] struct {
	Messages []MsgType
}

func (ListRecvStreamerOf[MsgType]) CloseSend() error {
	return nil
}

func (l *ListRecvStreamerOf[MsgType]) Recv() (req MsgType, err error) {
	if len(l.Messages) == 0 {
		// For some reason we cannot return nil here, so we have to hack this
		empty := new(MsgType)
		return *empty, io.EOF
	}

	msg := l.Messages[0]
	l.Messages = slices.Delete(l.Messages, 0, 1)
	return msg, nil
}

func (ListRecvStreamerOf[MsgType]) Header() (metadata.MD, error) {
	return nil, nil
}

func (ListRecvStreamerOf[MsgType]) Trailer() metadata.MD {
	return nil
}

func (ListRecvStreamerOf[MsgType]) Context() context.Context {
	return nil
}

func (ListRecvStreamerOf[MsgType]) SendMsg(_ interface{}) error {
	return nil
}

func (ListRecvStreamerOf[MsgType]) RecvMsg(_ interface{}) error {
	// TODO(oxisto): We took a shortcut here and implemented Recv directly.
	// However, some users might call RecvMsg directly and we should implement this
	// and call RecvMsg in Recv instead.
	return nil
}
