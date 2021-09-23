package evidence_store

import (
	"context"
	"io"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidence_store"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
)

var log *logrus.Entry

// Service is an implementation of the Clouditor evidence_store service (evidence_storeServer)
type Service struct {
	// Currently only in-memory
	evidences map[string]*assessment.Evidence
	evidence_store.UnimplementedEvidenceStoreServer
}

func NewService() *Service {
	return &Service{
		evidences: make(map[string]*assessment.Evidence),
	}
}

func init() {
	log = logrus.WithField("component", "evidence_store")
}

// StoreEvidence is a method implementation of the evidence_storeServer interface: It receives an evidence and stores it
func (s *Service) StoreEvidence(_ context.Context, evidence *assessment.Evidence) (*evidence_store.StoreEvidenceResponse, error) {
	var (
		resp       = &evidence_store.StoreEvidenceResponse{}
		err  error = nil
	)

	s.evidences[evidence.Id] = evidence
	resp.Status = true
	return resp, err
}

// StoreEvidences is a method implementation of the evidence_storeServer interface: It receives evidences and stores them
func (s *Service) StoreEvidences(stream evidence_store.EvidenceStore_StoreEvidencesServer) (err error) {
	var receivedEvidence *assessment.Evidence
	for {
		receivedEvidence, err = stream.Recv()
		if err == io.EOF {
			log.Infof("Stopped receiving streamed evidences")
			return stream.SendAndClose(&emptypb.Empty{})
		}
		s.evidences[receivedEvidence.Id] = receivedEvidence
	}
}

// ListEvidences is a method implementation of the evidence_storeServer interface: It returns the evidences lying in the evidence storage
func (s *Service) ListEvidences(_ context.Context, _ *evidence_store.ListEvidencesRequest) (*evidence_store.ListEvidencesResponse, error) {
	var listOfEvidences []*assessment.Evidence
	for _, v := range s.evidences {
		listOfEvidences = append(listOfEvidences, v)
	}

	return &evidence_store.ListEvidencesResponse{Evidences: listOfEvidences}, nil
}
