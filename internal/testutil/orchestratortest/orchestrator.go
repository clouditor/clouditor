package orchestratortest

import (
	"time"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testdata"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		Name:            testdata.MockCatalogName,
		Id:              testdata.MockCatalogID,
		Description:     testdata.MockCatalogDescription,
		AllInScope:      true,
		AssuranceLevels: []string{testdata.AssuranceLevelBasic, testdata.AssuranceLevelSubstantial, testdata.AssuranceLevelHigh},
		Categories: []*orchestrator.Category{{
			Name:        testdata.MockCategoryName,
			Description: testdata.MockCategoryDescription,
			CatalogId:   testdata.MockCatalogID,
			Controls: []*orchestrator.Control{{
				Id:                testdata.MockControlID1,
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

// NewTargetOfEvaluation creates a new Target of Evaluation. The assurance level is set if available.
func NewTargetOfEvaluation(assuranceLevel string) *orchestrator.TargetOfEvaluation {
	var toe = &orchestrator.TargetOfEvaluation{
		CloudServiceId: testdata.MockCloudServiceID,
		CatalogId:      testdata.MockCatalogID,
	}

	if assuranceLevel != "" {
		toe.AssuranceLevel = &assuranceLevel
	}

	// Our test catalog does not allow scoping, so we need to emulate what we do in CreateTargetOfEvaluation
	toe.ControlsInScope = []*orchestrator.Control{{
		Id:                testdata.MockControlID1,
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

var (
	MockAssessmentResult1 = &assessment.AssessmentResult{
		Id:             testdata.MockAssessmentResult1ID,
		Timestamp:      timestamppb.New(time.Unix(1, 0)),
		CloudServiceId: testdata.MockCloudServiceID,
		MetricId:       testdata.MockMetricID,
		Compliant:      true,
		EvidenceId:     testdata.MockEvidenceID,
		ResourceId:     testdata.MockResourceID,
		ResourceTypes:  []string{"Resource"},
		MetricConfiguration: &assessment.MetricConfiguration{
			Operator:       "==",
			TargetValue:    structpb.NewBoolValue(true),
			IsDefault:      true,
			MetricId:       testdata.MockMetricID,
			CloudServiceId: testdata.MockCloudServiceID,
		},
	}
	MockAssessmentResult2 = &assessment.AssessmentResult{
		Id:             testdata.MockAssessmentResult2ID,
		Timestamp:      timestamppb.New(time.Unix(1, 0)),
		CloudServiceId: testdata.MockAnotherCloudServiceID,
		MetricId:       testdata.MockMetricID,
		Compliant:      true,
		EvidenceId:     testdata.MockEvidenceID,
		ResourceId:     testdata.MockResourceID,
		ResourceTypes:  []string{"Resource"},
		MetricConfiguration: &assessment.MetricConfiguration{
			Operator:       "==",
			TargetValue:    structpb.NewBoolValue(true),
			IsDefault:      true,
			MetricId:       testdata.MockMetricID,
			CloudServiceId: testdata.MockAnotherCloudServiceID,
		},
	}
	MockAssessmentResult3 = &assessment.AssessmentResult{
		Id:             testdata.MockAssessmentResult3ID,
		Timestamp:      timestamppb.New(time.Unix(1, 0)),
		CloudServiceId: testdata.MockCloudServiceID,
		MetricId:       testdata.MockAnotherMetricID,
		Compliant:      false,
		EvidenceId:     testdata.MockEvidenceID,
		ResourceId:     testdata.MockResourceID,
		ResourceTypes:  []string{"Resource"},
		MetricConfiguration: &assessment.MetricConfiguration{
			Operator:       "==",
			TargetValue:    structpb.NewBoolValue(true),
			IsDefault:      true,
			MetricId:       testdata.MockMetricID,
			CloudServiceId: testdata.MockCloudServiceID,
		},
	}
	MockAssessmentResult4 = &assessment.AssessmentResult{
		Id:             testdata.MockAssessmentResult4ID,
		Timestamp:      timestamppb.New(time.Unix(1, 0)),
		CloudServiceId: testdata.MockAnotherCloudServiceID,
		MetricId:       testdata.MockAnotherMetricID,
		Compliant:      false,
		EvidenceId:     testdata.MockEvidenceID,
		ResourceId:     testdata.MockResourceID,
		ResourceTypes:  []string{"Resource"},
		MetricConfiguration: &assessment.MetricConfiguration{
			Operator:       "==",
			TargetValue:    structpb.NewBoolValue(true),
			IsDefault:      true,
			MetricId:       testdata.MockAnotherMetricID,
			CloudServiceId: testdata.MockAnotherCloudServiceID,
		},
	}
	MockAssessmentResults = []*assessment.AssessmentResult{MockAssessmentResult1, MockAssessmentResult2, MockAssessmentResult3, MockAssessmentResult4}

	MockControl1 = &orchestrator.Control{
		Id:                testdata.MockControlID1,
		Name:              testdata.MockControlName,
		CategoryName:      testdata.MockCategoryName,
		CategoryCatalogId: testdata.MockCatalogID,
		Description:       testdata.MockControlDescription,
		AssuranceLevel:    nil,
		Metrics: []*assessment.Metric{{
			Id:          testdata.MockMetricID,
			Name:        testdata.MockMetricName,
			Description: testdata.MockMetricDescription,
			Scale:       assessment.Metric_ORDINAL,
			Range: &assessment.Range{
				Range: &assessment.Range_AllowedValues{
					AllowedValues: &assessment.AllowedValues{
						Values: []*structpb.Value{
							structpb.NewBoolValue(false),
							structpb.NewBoolValue(true),
						},
					},
				},
			},
		}},
	}
	MockControl2 = &orchestrator.Control{
		Id:                testdata.MockControlID2,
		Name:              testdata.MockControlName,
		CategoryName:      testdata.MockCategoryName,
		CategoryCatalogId: testdata.MockCatalogID,
		Description:       testdata.MockControlDescription,
		AssuranceLevel:    &testdata.AssuranceLevelBasic,
		Metrics: []*assessment.Metric{{
			Id:          testdata.MockMetricID,
			Name:        testdata.MockMetricName,
			Description: testdata.MockMetricDescription,
			Scale:       assessment.Metric_ORDINAL,
			Range: &assessment.Range{
				Range: &assessment.Range_AllowedValues{
					AllowedValues: &assessment.AllowedValues{
						Values: []*structpb.Value{
							structpb.NewBoolValue(false),
							structpb.NewBoolValue(true),
						},
					},
				},
			},
		}},
	}
	MockControl3 = &orchestrator.Control{
		Id:                testdata.MockControlID3,
		Name:              testdata.MockControlName,
		CategoryName:      testdata.MockCategoryName,
		CategoryCatalogId: testdata.MockCatalogID,
		Description:       testdata.MockControlDescription,
		AssuranceLevel:    &testdata.AssuranceLevelSubstantial,
	}
	MockControl4 = &orchestrator.Control{
		Id:                testdata.MockControlID4,
		Name:              testdata.MockControlName,
		CategoryName:      testdata.MockCategoryName,
		CategoryCatalogId: testdata.MockCatalogID,
		Description:       testdata.MockControlDescription,
		AssuranceLevel:    &testdata.AssuranceLevelHigh,
		Metrics: []*assessment.Metric{{
			Id:          testdata.MockMetricID,
			Name:        testdata.MockMetricName,
			Description: testdata.MockMetricDescription,
			Scale:       assessment.Metric_ORDINAL,
			Range: &assessment.Range{
				Range: &assessment.Range_AllowedValues{
					AllowedValues: &assessment.AllowedValues{
						Values: []*structpb.Value{
							structpb.NewBoolValue(false),
							structpb.NewBoolValue(true),
						},
					},
				},
			},
		}},
	}
	MockControl5 = &orchestrator.Control{
		Id:                testdata.MockControlID5,
		Name:              testdata.MockControlName,
		CategoryName:      testdata.MockCategoryName,
		CategoryCatalogId: testdata.MockCatalogID,
		Description:       testdata.MockControlDescription,
		AssuranceLevel:    nil,
		Metrics: []*assessment.Metric{{
			Id:          testdata.MockMetricID,
			Name:        testdata.MockMetricName,
			Description: testdata.MockMetricDescription,
			Scale:       assessment.Metric_ORDINAL,
			Range: &assessment.Range{
				Range: &assessment.Range_AllowedValues{
					AllowedValues: &assessment.AllowedValues{
						Values: []*structpb.Value{
							structpb.NewBoolValue(false),
							structpb.NewBoolValue(true),
						},
					},
				},
			},
		}},
	}
	MockControls = []*orchestrator.Control{MockControl1, MockControl2, MockControl3, MockControl4, MockControl5}
)
