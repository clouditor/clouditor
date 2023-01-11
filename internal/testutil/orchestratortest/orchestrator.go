package orchestratortest

import (
	"time"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	MockCertificateID    = "1234"
	MockCatalogID        = "Cat1234"
	MockCategoryName     = "My name"
	MockControlID        = "Cont1234"
	MockSubControlID     = "Cont1234.1"
	MockAnotherControlID = "Cont4567"
	AssuranceLevelHigh   = "high"
	MockServiceID        = "11111111-1111-1111-1111-111111111111"
)

// NewCertificate creates a mock certificate
func NewCertificate() *orchestrator.Certificate {
	var mockCertificate = &orchestrator.Certificate{
		Id:             MockCertificateID,
		Name:           "EUCS",
		CloudServiceId: MockServiceID,
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
			CertificateId: MockCertificateID,
			Id:            "12345",
		}},
	}

	return mockCertificate
}

// NewCatalog creates a mock catalog
func NewCatalog() *orchestrator.Catalog {
	var mockCatalog = &orchestrator.Catalog{
		Name:        "MockCatalog",
		Id:          "Cat1234",
		Description: "This is a mock catalog",
		AllInScope:  true,
		Categories: []*orchestrator.Category{{
			Name:        "My name",
			Description: "test",
			CatalogId:   "Cat1234",
			Controls: []*orchestrator.Control{{
				Id:                "Cont1234",
				Name:              "Mock Control",
				CategoryName:      MockCategoryName,
				CategoryCatalogId: MockCatalogID,
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
				Controls: []*orchestrator.Control{{
					Id:                "Cont1234.1",
					Name:              "Mock Sub-Control",
					Description:       "This is a mock sub-control",
					Metrics:           []*assessment.Metric{},
					CategoryName:      MockCategoryName,
					CategoryCatalogId: MockCatalogID,
				}},
			},
				{
					Id:                "Cont4567",
					Name:              "Another Mock Control",
					CategoryName:      MockCategoryName,
					CategoryCatalogId: MockCatalogID,
				}},
		}}}
	return mockCatalog
}

func NewTargetOfEvaluation() *orchestrator.TargetOfEvaluation {
	var toe = &orchestrator.TargetOfEvaluation{
		CloudServiceId: MockServiceID,
		CatalogId:      MockCatalogID,
		AssuranceLevel: &AssuranceLevelHigh,
	}

	// Our test catalog does not allow scoping, so we need to emulate what we do in CreateTargetOfEvaluation
	toe.ControlsInScope = []*orchestrator.Control{{
		Id:                "Cont1234",
		CategoryName:      "My name",
		CategoryCatalogId: "Cat1234",
		Name:              "My name",
	}, {
		Id:                "Cont1234.1",
		CategoryName:      "My name",
		CategoryCatalogId: "Cat1234",
		Name:              "My name",
	}}

	return toe
}
