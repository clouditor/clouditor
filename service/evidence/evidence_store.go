package evidences

import (
	"context"
	"io"

	"clouditor.io/clouditor/api/evidence"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var log *logrus.Entry

// Service is an implementation of the Clouditor evidence service (evidenceServer)
type Service struct {
	// Currently only in-memory
	evidences map[string]*evidence.Evidence

	// Hook
	EvidenceHook []func(result *evidence.Evidence, err error)

	evidence.UnimplementedEvidenceStoreServer
}

func NewService() *Service {
	return &Service{
		evidences: make(map[string]*evidence.Evidence),
	}
}

func init() {
	log = logrus.WithField("component", "evidence")
}

// StoreEvidence is a method implementation of the evidenceServer interface: It receives an evidence and stores it
func (s *Service) StoreEvidence(_ context.Context, e *evidence.Evidence) (*evidence.StoreEvidenceResponse, error) {
	var (
		resp = &evidence.StoreEvidenceResponse{}
		err  error
	)

	_, err = e.Validate()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid evidence: %v", err)
	}

	s.evidences[e.Id] = e
	resp.Status = true

	// Inform our hook, if we have any
	if s.EvidenceHook != nil {
		for _, hook := range s.EvidenceHook {
			go hook(e, nil)
		}
	}

	return resp, err
}

// StoreEvidences is a method implementation of the evidenceServer interface: It receives evidences and stores them
func (s *Service) StoreEvidences(stream evidence.EvidenceStore_StoreEvidencesServer) (err error) {
	var receivedEvidence *evidence.Evidence

	for {
		receivedEvidence, err = stream.Recv()
		if err == io.EOF {
			log.Infof("Stopped receiving streamed evidences")
			return stream.SendAndClose(&emptypb.Empty{})
		}
		s.evidences[receivedEvidence.Id] = receivedEvidence

		// Inform our hook, if we have any
		if s.EvidenceHook != nil {
			for _, hook := range s.EvidenceHook {
				go hook(receivedEvidence, nil)
			}
		}
	}
}

// ListEvidences is a method implementation of the evidenceServer interface: It returns the evidences lying in the evidence storage
func (s *Service) ListEvidences(_ context.Context, _ *evidence.ListEvidencesRequest) (*evidence.ListEvidencesResponse, error) {
	var listOfEvidences []*evidence.Evidence
	for _, v := range s.evidences {
		listOfEvidences = append(listOfEvidences, v)
	}

	return &evidence.ListEvidencesResponse{Evidences: listOfEvidences}, nil
}

func (s *Service) RegisterEvidenceHook(evidenceHook func(result *evidence.Evidence, err error)) {
	s.EvidenceHook = append(s.EvidenceHook, evidenceHook)
}
