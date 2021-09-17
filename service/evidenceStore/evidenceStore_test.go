package evidenceStore

import (
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidenceStore"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"io"
	"reflect"
	"testing"
)

// TestNewService is a simply test for NewService
func TestNewService(t *testing.T) {
	tests := []struct {
		name string
		want evidenceStore.EvidenceStoreServer
	}{
		{
			name: "EvidenceStoreServer created with empty evidence map",
			want: &Service{
				evidences:                        make(map[string]*assessment.Evidence),
				UnimplementedEvidenceStoreServer: evidenceStore.UnimplementedEvidenceStoreServer{},
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
		in0      context.Context
		evidence *assessment.Evidence
	}
	tests := []struct {
		name     string
		args     args
		wantResp *evidenceStore.StoreEvidenceResponse
		wantErr  bool
	}{
		{
			name: "Store evidence to the map",
			args: args{
				in0: context.TODO(),
				evidence: &assessment.Evidence{
					Id:                "MockEvidenceId",
					ServiceId:         "MockServiceId",
					ResourceId:        "MockResourceId",
					Timestamp:         "TimeXY",
					ApplicableMetrics: []int32{1, 2},
					Raw:               "",
					Resource:          nil,
				},
			},
			wantErr:  false,
			wantResp: &evidenceStore.StoreEvidenceResponse{Status: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			gotResp, err := s.StoreEvidence(tt.args.in0, tt.args.evidence)
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

// TestStoreEvidence tests StoreEvidence
func TestStoreEvidences(t *testing.T) {
	type args struct {
		stream evidenceStore.EvidenceStore_StoreEvidencesServer
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
			}
		})
	}
}

type mockStreamer struct {
	counter int
}

func (mockStreamer) SendAndClose(r *evidenceStore.StoreEvidencesResponse) error {
	if r.Status {
		return nil
	}
	return fmt.Errorf("error occured while sending and closing")
}

func (m *mockStreamer) Recv() (*assessment.Evidence, error) {
	if m.counter == 0 {
		m.counter++
		return &assessment.Evidence{
			Id:                "MockEvidenceId-1",
			ServiceId:         "MockServiceId-1",
			ResourceId:        "MockResourceId-1",
			Timestamp:         "TimeXY",
			ApplicableMetrics: []int32{1, 2},
			Raw:               "",
			Resource:          nil,
		}, nil
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
