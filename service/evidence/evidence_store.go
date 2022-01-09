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

	// evidenceHooks is a list of hook functions that can be used if one wants to be
	// informed about each evidence
	evidenceHooks []evidence.EvidenceHookFunc

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
func (s *Service) StoreEvidence(_ context.Context, req *evidence.StoreEvidenceRequest) (resp *evidence.StoreEvidenceResponse, err error) {

	_, err = req.Evidence.Validate()
	if err != nil {
		log.Errorf("Invalid evidence: %v", err)
		newError := fmt.Errorf("invalid evidence: %w", err)

		s.informHooks(nil, newError)

		resp = &evidence.StoreEvidenceResponse{
			Status: false,
		}

		return resp, status.Errorf(codes.InvalidArgument, "invalid req: %v", err)
	}

	s.evidences[req.Evidence.Id] = req.Evidence
	s.informHooks(req.Evidence, nil)

	resp = &evidence.StoreEvidenceResponse{
		Status: true,
	}

	return resp, nil
}

// StoreEvidences is a method implementation of the evidenceServer interface: It receives evidences and stores them
func (s *Service) StoreEvidences(stream evidence.EvidenceStore_StoreEvidencesServer) (err error) {
	var req *evidence.StoreEvidenceRequest

	for {
		req, err = stream.Recv()

		if err != nil {
			// If no more input of the stream is available, return SendAndClose `error`
			if err == io.EOF {
				log.Infof("Stopped receiving streamed evidences")
				return stream.SendAndClose(&emptypb.Empty{})
			}

			return err
		}

		// Call StoreEvidence() for storing a single evidence
		evidenceRequest := &evidence.StoreEvidenceRequest{
			Evidence: req.Evidence,
		}
		_, err = s.StoreEvidence(context.Background(), evidenceRequest)
		if err != nil {
			return err
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

func (s *Service) RegisterEvidenceHook(evidenceHook evidence.EvidenceHookFunc) {
	s.evidenceHooks = append(s.evidenceHooks, evidenceHook)
}

func (s Service) informHooks(result *evidence.Evidence, err error) {
	// Inform our hook, if we have any
	if s.evidenceHooks != nil {
		for _, hook := range s.evidenceHooks {
			go hook(result, err)
		}
	}
}
