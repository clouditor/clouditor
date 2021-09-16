package evidenceStore

import (
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidenceStore"
	"context"
	"github.com/sirupsen/logrus"
	"io"
)

var log *logrus.Entry

// Service is an implementation of the Clouditor EvidenceStore service (EvidenceStoreServer)
type Service struct {
	// Currently only in-memory. ToDo(lebogg): Add connection to DB
	evidences map[string]*assessment.Evidence
	evidenceStore.UnimplementedEvidenceStoreServer
}

func NewService() evidenceStore.EvidenceStoreServer {
	return &Service{
		evidences: make(map[string]*assessment.Evidence),
	}
}

func init() {
	log = logrus.WithField("component", "evidenceStore")
}

func (s *Service) StoreEvidence(_ context.Context, evidence *assessment.Evidence) (resp *evidenceStore.StoreEvidenceResponse, err error) {
	log.Warnf("Storing evidence in-memory. But there is no other functionality here!")
	s.evidences[evidence.Id] = evidence
	resp.Status = true
	return
}

func (s *Service) StoreEvidences(stream evidenceStore.EvidenceStore_StoreEvidencesServer) (err error) {
	var receivedEvidence *assessment.Evidence
	log.Warnf("Storing evidences in-memory. But there is no other functionality here!")
	for {
		receivedEvidence, err = stream.Recv()
		s.evidences[receivedEvidence.Id] = receivedEvidence
		if err == io.EOF {
			log.Infof("Stopped receiving streamed evidences")
			return stream.SendAndClose(&evidenceStore.StoreEvidencesResponse{Status: true})
		}
	}
}
