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
	sync "sync"
	"testing"

	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"

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
		init           InitFuncOf[*recordedClientStream]
	}
	type fields struct {
		channels map[string]*StreamChannelOf[*recordedClientStream, proto.Message]
	}
	tests := []struct {
		name     string
		args     args
		fields   fields
		wantErr  assert.ErrorAssertionFunc
		wantRcvd int
	}{
		{
			name:   "missing init function",
			fields: fields{},
			args: args{
				target:    testdata.MockGRPCTarget,
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
				target:    testdata.MockGRPCTarget,
				component: "mycomponent",
				init: func(target string, additionalOpts ...grpc.DialOption) (m *recordedClientStream, err error) {
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
				init: func(target string, additionalOpts ...grpc.DialOption) (m *recordedClientStream, err error) {
					return &recordedClientStream{}, nil
				},
			},
			fields: fields{
				channels: map[string]*StreamChannelOf[*recordedClientStream, proto.Message]{},
			},
		},
		{
			name: "restarting stream",
			fields: fields{
				channels: map[string]*StreamChannelOf[*recordedClientStream, proto.Message]{
					"mock:1234": {
						dead: true,
						channel: func() chan proto.Message {
							// put 2 left over messages into the channel
							var c = make(chan proto.Message)
							go func() {
								c <- &assessment.AssessEvidenceRequest{Evidence: &evidence.Evidence{Id: testdata.MockEvidenceID1}}
							}()
							go func() {
								c <- &assessment.AssessEvidenceRequest{Evidence: &evidence.Evidence{Id: testdata.MockEvidenceID2}}
							}()
							return c
						}(),
						target:    "mock:1234",
						component: "mock",
					},
				},
			},
			args: args{
				target:    "mock:1234",
				component: "mock component",
				init: func(target string, additionalOpts ...grpc.DialOption) (m *recordedClientStream, err error) {
					return &recordedClientStream{}, nil
				},
			},
			wantRcvd: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StreamsOf[*recordedClientStream, proto.Message]{
				channels: tt.fields.channels,
				log:      defaultLog(),
			}
			c, err := s.GetStream(tt.args.target, tt.args.component, tt.args.init, tt.args.additionalOpts...)

			if tt.wantErr != nil {
				tt.wantErr(t, err, tt.args)
			} else {
				assert.Nil(t, err)

				// wait until our stream has received the wanted messages
				c.stream.wg.Add(tt.wantRcvd)
				c.stream.wg.Wait()
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
		want   assert.Want[*StreamsOf[*mockClientStream, proto.Message]]
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
					log: defaultLog(),
				},
			},
			want: func(t *testing.T, got *StreamsOf[*mockClientStream, protoreflect.ProtoMessage]) bool {
				// sendLoop should declare the channel dead
				return assert.True(t, got.channels["test"].dead)
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
					log: defaultLog(),
				},
			},
			want: func(t *testing.T, got *StreamsOf[*mockClientStream, protoreflect.ProtoMessage]) bool {
				// sendLoop should declare the channel dead
				return assert.True(t, got.channels["test"].dead)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &StreamChannelOf[*mockClientStream, proto.Message]{
				channel:   tt.fields.Channel,
				stream:    tt.fields.stream,
				target:    tt.fields.target,
				component: tt.fields.component,
			}

			// overwrite the stream channel to make sure we are dealing with the same stream object
			tt.args.s.channels["test"] = c

			// prepare something, otherwise the sendloop will block waiting for a message
			go func() { c.Send(&mockMessage{}) }()

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
	return m.SendMsg(msg)
}

func (mockClientStream) CloseSend() error {
	return nil
}

func (mockClientStream) Recv() (req *assessment.AssessEvidenceRequest, err error) {
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

type recordedClientStream struct {
	mockClientStream
	recvd []proto.Message
	wg    sync.WaitGroup
}

func (r *recordedClientStream) SendMsg(msg interface{}) error {
	r.recvd = append(r.recvd, msg.(proto.Message))
	r.wg.Done()
	return nil
}
