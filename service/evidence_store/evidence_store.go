package evidenceStore

import (
	"context"
	"io"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidenceStore"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

// Service is an implementation of the Clouditor EvidenceStore service (EvidenceStoreServer)
type Service struct {
	// Currently only in-memory
	evidences map[string]*assessment.Evidence
	evidenceStore.UnimplementedEvidenceStoreServer
}

func NewService() *Service {
	return &Service{
		evidences: make(map[string]*assessment.Evidence),
	}
}

func init() {
	log = logrus.WithField("component", "evidence_store")
}

// StoreEvidence is a method implementation of the EvidenceStoreServer interface: It receives an evidence and stores it
func (s *Service) StoreEvidence(_ context.Context, evidence *assessment.Evidence) (*evidenceStore.StoreEvidenceResponse, error) {
	var (
		resp       = &evidenceStore.StoreEvidenceResponse{}
		err  error = nil
	)

	s.evidences[evidence.Id] = evidence
	resp.Status = true
	return resp, err
}

// StoreEvidences is a method implementation of the EvidenceStoreServer interface: It receives evidences and stores them
func (s *Service) StoreEvidences(stream evidenceStore.EvidenceStore_StoreEvidencesServer) (err error) {
	var receivedEvidence *assessment.Evidence
	for {
		receivedEvidence, err = stream.Recv()
		if err == io.EOF {
			log.Infof("Stopped receiving streamed evidences")
			return stream.SendAndClose(&evidenceStore.StoreEvidencesResponse{Status: true})
		}
		s.evidences[receivedEvidence.Id] = receivedEvidence
	}
}

func (s *Service) ListEvidences(_ context.Context, _ *evidenceStore.ListEvidencesRequest) (*evidenceStore.ListEvidencesResponse, error) {
	var listOfEvidences []*assessment.Evidence
	for _, v := range s.evidences {
		listOfEvidences = append(listOfEvidences, v)
	}

	return &evidenceStore.ListEvidencesResponse{Evidences: listOfEvidences}, nil
}
