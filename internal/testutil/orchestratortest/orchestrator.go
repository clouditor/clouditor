package orchestratortest

import (
	"time"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testdata"
	"google.golang.org/protobuf/types/known/structpb"
)

// NewCertificate creates a mock certificate
func NewCertificate() *orchestrator.Certificate {
	var mockCertificate = &orchestrator.Certificate{
		Id:             testdata.MockCertificateID,
		Name:           testdata.MockCertificateName,
		CloudServiceId: testdata.MockCloudServiceID,
		IssueDate:      "2021-11-06",
		ExpirationDate: "2024-11-06",
		Standard:       testdata.MockCertificateName,
		AssuranceLevel: testdata.AssuranceLevelHigh,
		Cab:            testdata.MockCertificateCab,
		Description:    testdata.MockCertificateDescription,
		States: []*orchestrator.State{{
			State:         testdata.MockStateState,
			TreeId:        testdata.MockStateTreeID,
			Timestamp:     time.Now().String(),
			CertificateId: testdata.MockCertificateID,
			Id:            testdata.MockStateId,
		}},
	}

	return mockCertificate
}

// NewCatalog creates a mock catalog
func NewCatalog() *orchestrator.Catalog {
	var mockCatalog = &orchestrator.Catalog{
		Name:        testdata.MockCatalogName,
		Id:          testdata.MockCatalogID,
		Description: testdata.MockCatalogDescription,
		AllInScope:  true,
		Categories: []*orchestrator.Category{{
			Name:        testdata.MockCategoryName,
			Description: testdata.MockCategoryDescription,
			CatalogId:   testdata.MockCatalogID,
			Controls: []*orchestrator.Control{{
				Id:                testdata.MockControlID,
				Name:              testdata.MockControlName,
				CategoryName:      testdata.MockCategoryName,
				CategoryCatalogId: testdata.MockCatalogID,
				Description:       testdata.MockControlDescription,
				Metrics: []*assessment.Metric{{
					Id:          testdata.MockMetricID,
					Name:        testdata.MockMetricName,
					Description: testdata.MockMetricDescription,
					Scale:       assessment.Metric_ORDINAL,
					Range: &assessment.Range{
						Range: &assessment.Range_AllowedValues{AllowedValues: &assessment.AllowedValues{
							Values: []*structpb.Value{
								structpb.NewBoolValue(false),
								structpb.NewBoolValue(true),
							}}}},
				}},
				Controls: []*orchestrator.Control{{
					Id:                testdata.MockSubControlID,
					Name:              testdata.MockSubControlName,
					Description:       testdata.MockSubControlDescription,
					Metrics:           []*assessment.Metric{},
					CategoryName:      testdata.MockCategoryName,
					CategoryCatalogId: testdata.MockCatalogID,
				}},
			},
				{
					Id:                testdata.MockAnotherControlID,
					Name:              testdata.MockAnotherControlName,
					CategoryName:      testdata.MockCategoryName,
					CategoryCatalogId: testdata.MockCatalogID,
				}},
		}}}

	return mockCatalog
}

func NewTargetOfEvaluation() *orchestrator.TargetOfEvaluation {
	var toe = &orchestrator.TargetOfEvaluation{
		CloudServiceId: testdata.MockCloudServiceID,
		CatalogId:      testdata.MockCatalogID,
		AssuranceLevel: &testdata.AssuranceLevelHigh,
	}

	// Our test catalog does not allow scoping, so we need to emulate what we do in CreateTargetOfEvaluation
	toe.ControlsInScope = []*orchestrator.Control{{
		Id:                testdata.MockControlID,
		CategoryName:      testdata.MockCategoryName,
		CategoryCatalogId: testdata.MockCatalogID,
		Name:              testdata.MockControlName,
	}, {
		Id:                testdata.MockSubControlID,
		CategoryName:      testdata.MockCategoryName,
		CategoryCatalogId: testdata.MockCatalogID,
		Name:              testdata.MockSubControlName,
	}}

	return toe
}

func NewCloudService() *orchestrator.CloudService {
	return &orchestrator.CloudService{
		Id:                testdata.MockCloudServiceID,
		Name:              testdata.MockCloudServiceName,
		Description:       testdata.MockCloudServiceDescription,
		CatalogsInScope:   []*orchestrator.Catalog{},
		ConfiguredMetrics: []*assessment.Metric{},
	}
}

func NewAnotherCloudService() *orchestrator.CloudService {
	return &orchestrator.CloudService{
		Id:                testdata.MockAnotherCloudServiceID,
		Name:              testdata.MockAnotherCloudServiceName,
		Description:       testdata.MockAnotherCloudServiceDescription,
		CatalogsInScope:   []*orchestrator.Catalog{},
		ConfiguredMetrics: []*assessment.Metric{},
	}
}
