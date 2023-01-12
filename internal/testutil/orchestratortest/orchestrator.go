package orchestratortest

import (
	"time"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	// Cloud Service
	MockCloudServiceID   = "11111111-1111-1111-1111-111111111111"
	MockCloudServiceName = "Mock Cloud Service"

	// Catalog
	MockCatalogID                 = "Cat1234"
	MockCatalogName               = "Mock Catalog"
	MockCatalogDescription        = "This is a mock catalog"
	MockCategoryName              = "Mock Category Name"
	MockCategoryDescription       = "This is a mock category"
	MockControlID                 = "Cont1234"
	MockControlName               = "Mock Control Name"
	MockControlDescription        = "This is a mock control"
	MockSubControlID              = "Cont1234.1"
	MockSubControlName            = "Mock Sub-Control Name"
	MockSubControlDescription     = "This is a mock sub-control"
	MockAnotherControlID          = "Cont4567"
	MockAnotherControlName        = "Mock Another Control Name"
	MockAnotherControlDescription = "This is a another mock control"
	AssuranceLevelBasic           = "basic"
	AssuranceLevelSubstantial     = "substantial"
	AssuranceLevelHigh            = "high"

	// Metric
	MockMetricID          = "Mock Metric"
	MockMetricName        = "Mock Metric Name"
	MockMetricDescription = "This is a mock metric"

	// Certificate
	MockCertificateID          = "1234"
	MockCertificateName        = "EUCS"
	MockCertificateDescription = "This is a mock certificate"
	MockCertificateCab         = "Cab123"
	MockCertificateStandard    = "EUCS"

	// State
	MockStateId     = "12345"
	MockStateState  = "new"
	MockStateTreeID = "12345"
)

// var (
// 	MockCloudServiceID   = "11111111-1111-1111-1111-111111111111"
// 	MockCloudServiceName = "Mock Cloud Service"
// 	MockCertificateID    = "1234"
// 	MockCatalogID        = "Cat1234"
// 	MockCategoryName     = "My name"
// 	MockControlID        = "Cont1234"
// 	MockSubControlID     = "Cont1234.1"
// 	MockAnotherControlID = "Cont4567"
// 	AssuranceLevelHigh   = "high"
// )

// NewCertificate creates a mock certificate
func NewCertificate() *orchestrator.Certificate {
	var mockCertificate = &orchestrator.Certificate{
		Id:             MockCertificateID,
		Name:           MockCertificateName,
		CloudServiceId: MockCloudServiceID,
		IssueDate:      "2021-11-06",
		ExpirationDate: "2024-11-06",
		Standard:       MockCertificateName,
		AssuranceLevel: AssuranceLevelHigh,
		Cab:            MockCertificateCab,
		Description:    MockCertificateDescription,
		States: []*orchestrator.State{{
			State:         MockStateState,
			TreeId:        MockStateTreeID,
			Timestamp:     time.Now().String(),
			CertificateId: MockCertificateID,
			Id:            MockStateId,
		}},
	}

	return mockCertificate
}

// NewCatalog creates a mock catalog
func NewCatalog() *orchestrator.Catalog {
	var mockCatalog = &orchestrator.Catalog{
		Name:        MockCatalogName,
		Id:          MockCatalogID,
		Description: MockCatalogDescription,
		AllInScope:  true,
		Categories: []*orchestrator.Category{{
			Name:        MockCategoryName,
			Description: MockCategoryDescription,
			CatalogId:   MockCatalogID,
			Controls: []*orchestrator.Control{{
				Id:                MockControlID,
				Name:              MockControlName,
				CategoryName:      MockCategoryName,
				CategoryCatalogId: MockCatalogID,
				Description:       MockControlDescription,
				Metrics: []*assessment.Metric{{
					Id:          MockMetricID,
					Name:        MockMetricName,
					Description: MockMetricDescription,
					Scale:       assessment.Metric_ORDINAL,
					Range: &assessment.Range{
						Range: &assessment.Range_AllowedValues{AllowedValues: &assessment.AllowedValues{
							Values: []*structpb.Value{
								structpb.NewBoolValue(false),
								structpb.NewBoolValue(true),
							}}}},
				}},
				Controls: []*orchestrator.Control{{
					Id:                MockSubControlID,
					Name:              MockSubControlName,
					Description:       MockSubControlDescription,
					Metrics:           []*assessment.Metric{},
					CategoryName:      MockCategoryName,
					CategoryCatalogId: MockCatalogID,
				}},
			},
				{
					Id:                MockAnotherControlID,
					Name:              MockAnotherControlName,
					CategoryName:      MockCategoryName,
					CategoryCatalogId: MockCatalogID,
				}},
		}}}

	return mockCatalog
}

func NewTargetOfEvaluation() *orchestrator.TargetOfEvaluation {
	var toe = &orchestrator.TargetOfEvaluation{
		CloudServiceId: MockCloudServiceID,
		CatalogId:      MockCatalogID,
		AssuranceLevel: &AssuranceLevelHigh,
	}

	// Our test catalog does not allow scoping, so we need to emulate what we do in CreateTargetOfEvaluation
	toe.ControlsInScope = []*orchestrator.Control{{
		Id:                MockControlID,
		CategoryName:      MockCategoryName,
		CategoryCatalogId: MockCatalogID,
		Name:              MockControlName,
	}, {
		Id:                MockSubControlID,
		CategoryName:      MockCategoryName,
		CategoryCatalogId: MockCatalogID,
		Name:              MockSubControlName,
	}}

	return toe
}

// // NewCertificate creates a mock certificate
// func NewCertificate() *orchestrator.Certificate {
// 	var mockCertificate = &orchestrator.Certificate{
// 		Id:             MockCertificateID,
// 		Name:           MockCertificateName,
// 		CloudServiceId: MockCloudServiceID,
// 		IssueDate:      "2021-11-06",
// 		ExpirationDate: "2024-11-06",
// 		Standard:       MockCertificateName,
// 		AssuranceLevel: AssuranceLevelHigh,
// 		Cab:            MockCertificateCab,
// 		Description:    MockCertificateDescription,
// 		States: []*orchestrator.State{{
// 			State:         MockStateState,
// 			TreeId:        MockStateTreeID,
// 			Timestamp:     time.Now().String(),
// 			CertificateId: MockCertificateID,
// 			Id:            MockStateId,
// 		}},
// 	}

// 	return mockCertificate
// }

// // NewCatalog creates a mock catalog
// func NewCatalog() *orchestrator.Catalog {
// 	var mockCatalog = &orchestrator.Catalog{
// 		Name:        MockCatalogName,
// 		Id:          MockCatalogID,
// 		Description: MockCatalogDescription,
// 		AllInScope:  true,
// 		Categories: []*orchestrator.Category{{
// 			Name:        MockCategoryName,
// 			Description: MockCategoryDescription,
// 			CatalogId:   MockCatalogID,
// 			Controls: []*orchestrator.Control{{
// 				Id:                MockControlID,
// 				Name:              MockControlName,
// 				CategoryName:      MockCategoryName,
// 				CategoryCatalogId: MockCatalogID,
// 				Description:       MockCategoryDescription,
// 				Metrics: []*assessment.Metric{{
// 					Id:          MockMetricID,
// 					Name:        MockMetricName,
// 					Description: MockMetricDescription,
// 					Scale:       assessment.Metric_ORDINAL,
// 					Range: &assessment.Range{
// 						Range: &assessment.Range_AllowedValues{AllowedValues: &assessment.AllowedValues{
// 							Values: []*structpb.Value{
// 								structpb.NewBoolValue(false),
// 								structpb.NewBoolValue(true),
// 							}}}},
// 				}},
// 				Controls: []*orchestrator.Control{{
// 					Id:                MockSubControlID,
// 					Name:              MockSubControlName,
// 					Description:       MockSubControlDescription,
// 					Metrics:           []*assessment.Metric{},
// 					CategoryName:      MockCategoryName,
// 					CategoryCatalogId: MockCatalogID,
// 				}},
// 			},
// 				{
// 					Id:                MockAnotherControlID,
// 					Name:              MockAnotherControlName,
// 					CategoryName:      MockCategoryName,
// 					CategoryCatalogId: MockCatalogID,
// 				}},
// 		}}}

// 	return mockCatalog
// }

// func NewTargetOfEvaluation() *orchestrator.TargetOfEvaluation {
// 	var toe = &orchestrator.TargetOfEvaluation{
// 		CloudServiceId: MockCloudServiceID,
// 		CatalogId:      MockCatalogID,
// 		AssuranceLevel: &AssuranceLevelHigh,
// 	}

// 	// Our test catalog does not allow scoping, so we need to emulate what we do in CreateTargetOfEvaluation
// 	toe.ControlsInScope = []*orchestrator.Control{{
// 		Id:                MockControlID,
// 		CategoryName:      MockCategoryName,
// 		CategoryCatalogId: MockCatalogID,
// 		Name:              MockControlName,
// 	}, {
// 		Id:                MockSubControlID,
// 		CategoryName:      MockCategoryName,
// 		CategoryCatalogId: MockCatalogID,
// 		Name:              MockSubControlName,
// 	}}

// 	return toe
// }
