package orchestratortest

import (
	"time"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testvariables"
	"google.golang.org/protobuf/types/known/structpb"
)

// NewCertificate creates a mock certificate
func NewCertificate() *orchestrator.Certificate {
	var mockCertificate = &orchestrator.Certificate{
		Id:             testvariables.MockCertificateID,
		Name:           testvariables.MockCertificateName,
		CloudServiceId: testvariables.MockCloudServiceID,
		IssueDate:      "2021-11-06",
		ExpirationDate: "2024-11-06",
		Standard:       testvariables.MockCertificateName,
		AssuranceLevel: testvariables.AssuranceLevelHigh,
		Cab:            testvariables.MockCertificateCab,
		Description:    testvariables.MockCertificateDescription,
		States: []*orchestrator.State{{
			State:         testvariables.MockStateState,
			TreeId:        testvariables.MockStateTreeID,
			Timestamp:     time.Now().String(),
			CertificateId: testvariables.MockCertificateID,
			Id:            testvariables.MockStateId,
		}},
	}

	return mockCertificate
}

// NewCatalog creates a mock catalog
func NewCatalog() *orchestrator.Catalog {
	var mockCatalog = &orchestrator.Catalog{
		Name:        testvariables.MockCatalogName,
		Id:          testvariables.MockCatalogID,
		Description: testvariables.MockCatalogDescription,
		AllInScope:  true,
		Categories: []*orchestrator.Category{{
			Name:        testvariables.MockCategoryName,
			Description: testvariables.MockCategoryDescription,
			CatalogId:   testvariables.MockCatalogID,
			Controls: []*orchestrator.Control{{
				Id:                testvariables.MockControlID,
				Name:              testvariables.MockControlName,
				CategoryName:      testvariables.MockCategoryName,
				CategoryCatalogId: testvariables.MockCatalogID,
				Description:       testvariables.MockControlDescription,
				Metrics: []*assessment.Metric{{
					Id:          testvariables.MockMetricID,
					Name:        testvariables.MockMetricName,
					Description: testvariables.MockMetricDescription,
					Scale:       assessment.Metric_ORDINAL,
					Range: &assessment.Range{
						Range: &assessment.Range_AllowedValues{AllowedValues: &assessment.AllowedValues{
							Values: []*structpb.Value{
								structpb.NewBoolValue(false),
								structpb.NewBoolValue(true),
							}}}},
				}},
				Controls: []*orchestrator.Control{{
					Id:                testvariables.MockSubControlID,
					Name:              testvariables.MockSubControlName,
					Description:       testvariables.MockSubControlDescription,
					Metrics:           []*assessment.Metric{},
					CategoryName:      testvariables.MockCategoryName,
					CategoryCatalogId: testvariables.MockCatalogID,
				}},
			},
				{
					Id:                testvariables.MockAnotherControlID,
					Name:              testvariables.MockAnotherControlName,
					CategoryName:      testvariables.MockCategoryName,
					CategoryCatalogId: testvariables.MockCatalogID,
				}},
		}}}

	return mockCatalog
}

func NewTargetOfEvaluation() *orchestrator.TargetOfEvaluation {
	var toe = &orchestrator.TargetOfEvaluation{
		CloudServiceId: testvariables.MockCloudServiceID,
		CatalogId:      testvariables.MockCatalogID,
		AssuranceLevel: &testvariables.AssuranceLevelHigh,
	}

	// Our test catalog does not allow scoping, so we need to emulate what we do in CreateTargetOfEvaluation
	toe.ControlsInScope = []*orchestrator.Control{{
		Id:                testvariables.MockControlID,
		CategoryName:      testvariables.MockCategoryName,
		CategoryCatalogId: testvariables.MockCatalogID,
		Name:              testvariables.MockControlName,
	}, {
		Id:                testvariables.MockSubControlID,
		CategoryName:      testvariables.MockCategoryName,
		CategoryCatalogId: testvariables.MockCatalogID,
		Name:              testvariables.MockSubControlName,
	}}

	return toe
}
