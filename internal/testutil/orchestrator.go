package testutil

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
		CertificateID: "1234",
		ID:            "12345",
	}

	var mockCertificate = &orchestrator.Certificate{
		Name:        "EUCS",
		ServiceId:   "test service",
		Issuedate:   "2021-11-06",
		Standard:    "EUCS",
		Scope:       "Basic",
		Cab:         "Cab123",
		Description: "Description",
		States:      []*orchestrator.State{mockHistory},
		ID:          "1234",
	}

	return mockCertificate
}
