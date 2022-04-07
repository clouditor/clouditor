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
	"context"
	"errors"
	"io"
	"testing"

	"clouditor.io/clouditor/api/assessment"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var ErrSomeError = errors.New("some error")

func TestStreamsOf_GetStream(t *testing.T) {
	type args struct {
		target         string
		component      string
		additionalOpts []grpc.DialOption
		init           InitFuncOf[*mockClientStream]
	}
	type fields struct {
		channels map[string]*StreamChannelOf[*mockClientStream, proto.Message]
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "missing init function",
			fields: fields{},
			args: args{
				target:    "localhost",
				component: "mycomponent",
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrMissingInitFunc)
			},
		},
		{
			name:   "error in init function",
			fields: fields{},
			args: args{
				target:    "localhost",
				component: "mycomponent",
				init: func(target string, additionalOpts ...grpc.DialOption) (m *mockClientStream, err error) {
					return nil, ErrSomeError
				},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrSomeError)
			},
		},
		{
			name: "adding new stream to mock client",
			args: args{
				target:    "mock:1234",
				component: "mock component",
				init: func(target string, additionalOpts ...grpc.DialOption) (m *mockClientStream, err error) {
					return &mockClientStream{}, nil
				},
			},
			fields: fields{
				channels: map[string]*StreamChannelOf[*mockClientStream, proto.Message]{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StreamsOf[*mockClientStream, proto.Message]{
				channels: tt.fields.channels,
			}
			_, err := s.GetStream(tt.args.target, tt.args.component, tt.args.init, tt.args.additionalOpts...)

			if tt.wantErr != nil {
				tt.wantErr(t, err, tt.args)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func Test_StreamChannelOf_sendLoop(t *testing.T) {
	type args struct {
		s *StreamsOf[*mockClientStream, proto.Message]
	}
	type fields struct {
		Channel   chan proto.Message
		stream    *mockClientStream
		target    string
		component string
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		want   assert.ValueAssertionFunc
	}{
		{
			name: "send error",
			fields: fields{
				Channel: make(chan protoreflect.ProtoMessage),
				stream:  &mockClientStream{sendErr: errors.New("some error")},
				target:  "test",
			},
			args: args{
				s: &StreamsOf[*mockClientStream, proto.Message]{
					// add an existing channel with the same target name, to see if we remove it on error
					channels: map[string]*StreamChannelOf[*mockClientStream, protoreflect.ProtoMessage]{
						"test": &StreamChannelOf[*mockClientStream, proto.Message]{},
					},
				},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				s := i1.(*StreamsOf[*mockClientStream, proto.Message])

				// sendLoop should remove itself from the channels on error
				return assert.Empty(t, s.channels)
			},
		},
		{
			name: "send EOF",
			fields: fields{
				Channel: make(chan protoreflect.ProtoMessage),
				stream:  &mockClientStream{sendErr: io.EOF},
				target:  "test",
			},
			args: args{
				s: &StreamsOf[*mockClientStream, proto.Message]{
					// add an existing channel with the same target name, to see if we remove it on error
					channels: map[string]*StreamChannelOf[*mockClientStream, protoreflect.ProtoMessage]{
						"test": &StreamChannelOf[*mockClientStream, proto.Message]{},
					},
				},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				s := i1.(*StreamsOf[*mockClientStream, proto.Message])

				// sendLoop should remove itself from the channels on error
				return assert.Empty(t, s.channels)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &StreamChannelOf[*mockClientStream, proto.Message]{
				Channel:   tt.fields.Channel,
				stream:    tt.fields.stream,
				target:    tt.fields.target,
				component: tt.fields.component,
			}

			// prepare something, otherwise the sendloop will block waiting for a message
			go func() { c.Channel <- mockMessage{} }()

			c.sendLoop(tt.args.s)

			if tt.want != nil {
				tt.want(t, tt.args.s)
			}
		})
	}
}

type mockClientStream struct {
	sendErr error
}

type mockMessage struct{}

func (mockMessage) ProtoReflect() protoreflect.Message {
	return nil
}

func (m mockClientStream) Send(msg proto.Message) error {
	return m.Send(msg)
}

func (mockClientStream) CloseSend() error {
	return nil
}

func (m *mockClientStream) Recv() (req *assessment.AssessEvidenceRequest, err error) {
	return nil, nil
}

func (mockClientStream) Header() (metadata.MD, error) {
	return nil, nil
}

func (mockClientStream) Trailer() metadata.MD {
	return nil
}

func (mockClientStream) SendHeader(metadata.MD) error {
	return nil
}

func (mockClientStream) Context() context.Context {
	return nil
}

func (m mockClientStream) SendMsg(interface{}) error {
	return m.sendErr
}

func (mockClientStream) RecvMsg(interface{}) error {
	return nil
}
