package evidences

import (
	"context"
	"fmt"
	"io"

	"clouditor.io/clouditor/api/evidence"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var log *logrus.Entry

// Service is an implementation of the Clouditor req service (evidenceServer)
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
	log = logrus.WithField("component", "req")
}

// StoreEvidence is a method implementation of the evidenceServer interface: It receives an req and stores it
func (s *Service) StoreEvidence(_ context.Context, req *evidence.StoreEvidenceRequest) (*evidence.StoreEvidenceResponse, error) {
	var (
		resp = &evidence.StoreEvidenceResponse{}
		err  error
		e    = req.GetEvidence()
	)

	_, err = e.Validate()
	if err != nil {
		log.Errorf("Invalid evidence: %v", err)
		newError := fmt.Errorf("invalid evidence: %w", err)

		// Inform our hook, if we have any
		if s.EvidenceHook != nil {
			for _, hook := range s.EvidenceHook {
				go hook(nil, newError)
			}
		}

		return resp, status.Errorf(codes.InvalidArgument, "invalid req: %v", err)
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
	var (
		req *evidence.StoreEvidenceRequest
		e   *evidence.Evidence
	)
	for {
		req, err = stream.Recv()
		e = req.GetEvidence()
		if err == io.EOF {
			log.Infof("Stopped receiving streamed evidences")
			return stream.SendAndClose(&emptypb.Empty{})
		}
		s.evidences[e.Id] = e

		// Inform our hook, if we have any
		if s.EvidenceHook != nil {
			for _, hook := range s.EvidenceHook {
				go hook(e, nil)
			}
		}
	}
}

// ListEvidences is a method implementation of the evidenceServer interface: It returns the evidences lying in the req storage
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
