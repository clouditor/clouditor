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

	// Check if evidence is valid and inform hook. If not, return status error with code InvalidArgument
	resp, err = s.validateEvidence(req.Evidence)
	if err != nil {
		return resp, err
	}

	err = s.handleEvidence(req.GetEvidence())
	if err != nil {
		resp = &evidence.StoreEvidenceResponse{
			Status: false,
		}
		return resp, status.Errorf(codes.Internal, "error while handling evidence: %v", err)
	}

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

		// Check if evidence is valid and inform hook. If not, return status error with code InvalidArgument
		_, err = s.validateEvidence(req.Evidence)
		if err != nil {
			return err
		}

		err = s.handleEvidence(req.GetEvidence())
		if err != nil {
			return status.Errorf(codes.Internal, "Error while handling evidence: %v", err)
		}
	}
}

//validateEvidence checks if evidence is valid, informs the hook function and created the status error.
func (s *Service) validateEvidence(e *evidence.Evidence) (resp *evidence.StoreEvidenceResponse, err error) {
	_, err = e.Validate()
	if err != nil {
		log.Errorf("Invalid evidence: %v", err)
		newError := fmt.Errorf("invalid evidence: %w", err)

		s.informHooks(nil, newError)

		resp = &evidence.StoreEvidenceResponse{
			Status: false,
		}

		return resp, status.Errorf(codes.InvalidArgument, "invalid req: %v", err)
	} else {
		resp = &evidence.StoreEvidenceResponse{}

		return resp, nil
	}

}

func (s *Service) handleEvidence(e *evidence.Evidence) (err error) {
	s.evidences[e.Id] = e
	s.informHooks(e, nil)

	return
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
