package orchestratortest

import (
	"time"

	"clouditor.io/clouditor/api/orchestrator"
)

// NewCertificate creates a mock certificate creation request
func NewCertificate() *orchestrator.Certificate {
	var mockCertificate = &orchestrator.Certificate{
		Name:           "EUCS",
		ServiceId:      "test service",
		IssueDate:      "2021-11-06",
		ExpirationDate: "2024-11-06",
		Standard:       "EUCS",
		AssuranceLevel: "Basic",
		Cab:            "Cab123",
		Description:    "Description",
		States: []*orchestrator.State{{
			State:         "new",
			TreeId:        "12345",
			Timestamp:     time.Now().String(),
			CertificateId: "1234",
			Id:            "12345",
		}},
		Id: "1234",
	}

	return mockCertificate
}
