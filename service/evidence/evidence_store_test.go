package evidences

import (
	"context"
	"io"
	"reflect"
	"runtime"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/voc"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

// TestNewService is a simply test for NewService
func TestNewService(t *testing.T) {
	tests := []struct {
		name string
		want evidence.EvidenceStoreServer
	}{
		{
			name: "EvidenceStoreServer created with empty req map",
			want: &Service{
				evidences:                        make(map[string]*evidence.Evidence),
				UnimplementedEvidenceStoreServer: evidence.UnimplementedEvidenceStoreServer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewService(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestStoreEvidence tests StoreEvidence
func TestStoreEvidence(t *testing.T) {
	type args struct {
		in0 context.Context
		req *evidence.StoreEvidenceRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *evidence.StoreEvidenceResponse
		wantErr  bool
	}{
		{
			name: "Store req to the map",
			args: args{
				in0: context.TODO(),
				req: &evidence.StoreEvidenceRequest{Evidence: &evidence.Evidence{
					Id:        "MockEvidenceId",
					ServiceId: "MockServiceId",
					ToolId:    "MockTool",
					Timestamp: timestamppb.Now(),
					Raw:       "",
					Resource: toStruct(voc.VirtualMachine{
						Compute: &voc.Compute{CloudResource: &voc.CloudResource{
							ID: "mock-id",
						}},
					}, t),
				}},
			},
			wantErr:  false,
			wantResp: &evidence.StoreEvidenceResponse{Status: true},
		},
		{
			name: "Store an evidence without toolId to the map",
			args: args{
				in0: context.TODO(),
				req: &evidence.StoreEvidenceRequest{Evidence: &evidence.Evidence{
					Id:        "MockEvidenceId-1",
					ServiceId: "MockServiceId-1",
					Timestamp: timestamppb.Now(),
					Raw:       "",
					Resource: toStruct(voc.VirtualMachine{
						Compute: &voc.Compute{
							CloudResource: &voc.CloudResource{
								ID: "mock-id-1",
							},
						},
					}, t),
				},
				},
			},
			wantErr: true,
			wantResp: &evidence.StoreEvidenceResponse{
				Status: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			gotResp, err := s.StoreEvidence(tt.args.in0, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreEvidence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("StoreEvidence() gotResp = %v, want %v", gotResp, tt.wantResp)
			}

			if gotResp.Status {
				assert.NotNil(t, s.evidences["MockEvidenceId"])
			} else {
				assert.Empty(t, s.evidences)
			}
		})
	}
}

// TestStoreEvidences tests StoreEvidences
func TestStoreEvidences(t *testing.T) {
	type args struct {
		stream evidence.EvidenceStore_StoreEvidencesServer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Store 2 evidences to the map",
			args:    args{stream: &mockStreamer{counter: 0}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			if err := s.StoreEvidences(tt.args.stream); (err != nil) != tt.wantErr {
				t.Errorf("StoreEvidences() error = %v, wantErr %v", err, tt.wantErr)
				assert.Equal(t, 2, len(s.evidences))
			}
		})
	}
}

// TestListEvidences tests List req
func TestListEvidences(t *testing.T) {
	s := NewService()
	s.evidences["MockEvidenceId-1"] = &evidence.Evidence{
		Id:        "MockEvidenceId-1",
		ServiceId: "MockServiceId-1",
		Timestamp: timestamppb.Now(),
		Raw:       "",
		Resource:  nil,
	}
	s.evidences["MockEvidenceId-2"] = &evidence.Evidence{
		Id:        "MockEvidenceId-2",
		ServiceId: "MockServiceId-2",
		Timestamp: timestamppb.Now(),
		Raw:       "",
		Resource:  nil,
	}

	resp, err := s.ListEvidences(context.TODO(), &evidence.ListEvidencesRequest{})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(resp.Evidences))
}

func TestEvidenceHook(t *testing.T) {
	var ready1 = make(chan bool)
	var ready2 = make(chan bool)
	hookCallCounter := 0

	firstHookFunction := func(evidence *evidence.Evidence, err error) {
		hookCallCounter++
		log.Println("Hello from inside the firstHookFunction")

		ready1 <- true
	}

	secondHookFunction := func(evidence *evidence.Evidence, err error) {
		hookCallCounter++
		log.Println("Hello from inside the secondHookFunction")

		ready2 <- true
	}

	service := NewService()
	service.RegisterEvidenceHook(firstHookFunction)
	service.RegisterEvidenceHook(secondHookFunction)

	// Check if first hook is registered
	funcName1 := runtime.FuncForPC(reflect.ValueOf(service.EvidenceHook[0]).Pointer()).Name()
	funcName2 := runtime.FuncForPC(reflect.ValueOf(firstHookFunction).Pointer()).Name()
	assert.Equal(t, funcName1, funcName2)

	// Check if second hook is registered
	funcName1 = runtime.FuncForPC(reflect.ValueOf(service.EvidenceHook[1]).Pointer()).Name()
	funcName2 = runtime.FuncForPC(reflect.ValueOf(secondHookFunction).Pointer()).Name()
	assert.Equal(t, funcName1, funcName2)

	// Check GRPC call
	type args struct {
		in0      context.Context
		evidence *evidence.StoreEvidenceRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *evidence.StoreEvidenceResponse
		wantErr  bool
	}{
		{
			name: "Store an evidence to the map",
			args: args{
				in0: context.TODO(),
				evidence: &evidence.StoreEvidenceRequest{Evidence: &evidence.Evidence{
					Id:        "MockEvidenceId-1",
					ServiceId: "MockServiceId-1",
					Timestamp: timestamppb.Now(),
					Raw:       "",
					ToolId:    "mockToolId-1",
					Resource: toStruct(voc.VirtualMachine{
						Compute: &voc.Compute{
							CloudResource: &voc.CloudResource{
								ID: "mock-id-1",
							},
						},
					}, t),
				},
				},
			},
			wantErr: false,
			wantResp: &evidence.StoreEvidenceResponse{
				Status: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hookCallCounter = 0
			s := service
			gotResp, err := s.StoreEvidence(tt.args.in0, tt.args.evidence)

			//make the test wait
			select {
			case <-ready1:
				break
			case <-time.After(10 * time.Second):
				log.Println("Timeout while waiting for first StoreEvidence to be ready")
			}

			select {
			case <-ready2:
				break
			case <-time.After(10 * time.Second):
				log.Println("Timeout while waiting for second StoreEvidence to be ready")
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("StoreEvidence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("StoreEvidence() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
			assert.NotEmpty(t, s.evidences)
			assert.Equal(t, 2, hookCallCounter)
		})
	}

}

type mockStreamer struct {
	counter int
}

func (mockStreamer) SendAndClose(_ *emptypb.Empty) error {
	return nil
}

func (m *mockStreamer) Recv() (req *evidence.StoreEvidenceRequest, err error) {
	if m.counter == 0 {
		m.counter++
		return &evidence.StoreEvidenceRequest{Evidence: &evidence.Evidence{
			Id:        "MockEvidenceId-1",
			ServiceId: "MockServiceId-1",
			Timestamp: timestamppb.Now(),
			Raw:       "",
			Resource:  nil,
		}}, nil
	} else if m.counter == 1 {
		m.counter++
		return &evidence.StoreEvidenceRequest{Evidence: &evidence.Evidence{
			Id:        "MockEvidenceId-2",
			ServiceId: "MockServiceId-2",
			Timestamp: timestamppb.Now(),
			Raw:       "",
			Resource:  nil,
		}}, nil
	} else {
		return nil, io.EOF
	}
}

func (mockStreamer) SetHeader(_ metadata.MD) error {
	panic("implement me")
}

func (mockStreamer) SendHeader(_ metadata.MD) error {
	panic("implement me")
}

func (mockStreamer) SetTrailer(_ metadata.MD) {
	panic("implement me")
}

func (mockStreamer) Context() context.Context {
	panic("implement me")
}

func (mockStreamer) SendMsg(_ interface{}) error {
	panic("implement me")
}

func (mockStreamer) RecvMsg(_ interface{}) error {
	panic("implement me")
}

func toStruct(r voc.IsCloudResource, t *testing.T) (s *structpb.Value) {
	s, err := voc.ToStruct(r)
	if err != nil {
		assert.NotNil(t, err)
	}

	return
}
