package orchestratortest

import (
	"clouditor.io/clouditor/api/orchestrator"
	"time"
)

// CreateCertificateMock creates a mock certificate creation request
func CreateCertificateMock() *orchestrator.Certificate {
	mockHistory := &orchestrator.State{
		State:         "new",
		TreeId:        "12345",
		Timestamp:     time.Now().String(),
		CertificateId: "1234",
		Id:            "12345",
	}

	var mockCertificate = &orchestrator.Certificate{
		Name:           "EUCS",
		ServiceId:      "test service",
		IssueDate:      "2021-11-06",
		ExpirationDate: "2024-11-06",
		Standard:       "EUCS",
		AssuranceLevel: "Basic",
		Cab:            "Cab123",
		Description:    "Description",
		States:         []*orchestrator.State{mockHistory},
		Id:             "1234",
	}

	return mockCertificate
}
