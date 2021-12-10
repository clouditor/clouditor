package evidences

import (
	"context"
	"io"
	"reflect"
	"testing"

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
			assert.NotNil(t, s.evidences["MockEvidenceId"])
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
