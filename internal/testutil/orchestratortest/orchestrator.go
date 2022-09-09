package orchestratortest

import (
	"time"

	"google.golang.org/protobuf/types/known/structpb"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
)

// NewCertificate creates a mock certificate
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

// NewCertificate creates a mock certificate
func NewControl() *orchestrator.Control {
	var mockControl = &orchestrator.Control{
		ShortName:         "Cont1234",
		CategoryName:      "My name",
		CategoryCatalogId: "Cat1234",
		Name:              "Mock Control",
		Description:       "This is a mock control",
		Metrics: []*assessment.Metric{{
			Id:          "MockMetric",
			Name:        "A Mock Metric",
			Description: "This Metric is a mock metric",
			Scale:       assessment.Metric_ORDINAL,
			Range: &assessment.Range{
				Range: &assessment.Range_AllowedValues{AllowedValues: &assessment.AllowedValues{
					Values: []*structpb.Value{
						structpb.NewBoolValue(false),
						structpb.NewBoolValue(true),
					}}}},
		}},
	}
	return mockControl
}

// NewCatalog creates a mock catalog
func NewCatalog() *orchestrator.Catalog {
	var mockCatalog = &orchestrator.Catalog{
		Name:        "MockCatalog",
		Id:          "Cat1234",
		Description: "This is a mock catalog",
		Categories: []*orchestrator.Category{{
			Name:        "My name",
			Description: "test",
			CatalogId:   "Cat1234",
			Controls: []*orchestrator.Control{{
				ShortName:   "Cont1234",
				Name:        "Mock Control",
				Description: "This is a mock control",
				/*Metrics: []*assessment.Metric{{
					Id:          "MockMetric",
					Name:        "A Mock Metric",
					Description: "This Metric is a mock metric",
					Scale:       assessment.Metric_ORDINAL,
					Range: &assessment.Range{
						Range: &assessment.Range_AllowedValues{AllowedValues: &assessment.AllowedValues{
							Values: []*structpb.Value{
								structpb.NewBoolValue(false),
								structpb.NewBoolValue(true),
							}}}},
				}}*/
				// create a nested control
				Controls: []*orchestrator.Control{{
					ShortName:   "Cont1234.1",
					Name:        "Mock Sub-Control",
					Description: "This is a mock sub-control",
					//Metrics:     []*assessment.Metric{},
				}},
			}},
		}}}
	return mockCatalog
}
