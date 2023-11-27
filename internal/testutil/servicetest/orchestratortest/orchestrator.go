package orchestratortest

import (
	"time"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/util"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NewCertificateOption is an option for NewCertificate to modify single properties for testing , e.g., Update endpoint
type NewCertificateOption func(*orchestrator.Certificate)

func WithDescription(description string) NewCertificateOption {
	return func(certificate *orchestrator.Certificate) {
		certificate.Description = description
	}
}

// NewCertificate creates a mock certificate.
func NewCertificate(opts ...NewCertificateOption) *orchestrator.Certificate {
	timeStamp := time.Date(2011, 7, 1, 0, 0, 0, 0, time.UTC)
	var mockCertificate = &orchestrator.Certificate{
		Id:             testdata.MockCertificateID,
		Name:           testdata.MockCertificateName,
		CloudServiceId: testdata.MockCloudServiceID1,
		IssueDate:      timeStamp.AddDate(-5, 0, 0).String(),
		ExpirationDate: timeStamp.AddDate(5, 0, 0).String(),
		Standard:       testdata.MockCertificateName,
		AssuranceLevel: testdata.AssuranceLevelHigh,
		Cab:            testdata.MockCertificateCab,
		Description:    testdata.MockCertificateDescription,
		States: []*orchestrator.State{{
			State:         testdata.MockStateState,
			TreeId:        testdata.MockStateTreeID,
			Timestamp:     timeStamp.String(),
			CertificateId: testdata.MockCertificateID,
			Id:            testdata.MockStateId,
		}},
	}

	for _, o := range opts {
		o(mockCertificate)
	}

	return mockCertificate
}

// NewCertificate2 creates a mock certificate with other properties than NewCertificate
func NewCertificate2() *orchestrator.Certificate {
	timeStamp := time.Date(2014, 12, 1, 0, 0, 0, 0, time.UTC)
	var mockCertificate = &orchestrator.Certificate{
		Id:             testdata.MockCertificateID2,
		Name:           testdata.MockCertificateName2,
		CloudServiceId: testdata.MockCloudServiceID2,
		IssueDate:      timeStamp.AddDate(-5, 0, 0).String(),
		ExpirationDate: timeStamp.AddDate(5, 0, 0).String(),
		Standard:       testdata.MockCertificateName2,
		AssuranceLevel: testdata.AssuranceLevelHigh,
		Cab:            testdata.MockCertificateCab2,
		Description:    testdata.MockCertificateDescription2,
		States: []*orchestrator.State{{
			State:         testdata.MockStateState2,
			TreeId:        testdata.MockStateTreeID2,
			Timestamp:     timeStamp.String(),
			CertificateId: testdata.MockCertificateID2,
			Id:            testdata.MockStateId2,
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
			Controls: []*orchestrator.Control{
				MockControl1,
				MockControl2,
			},
		}}}

	return mockCatalog
}

// NewTargetOfEvaluation creates a new Target of Evaluation. The assurance level is set if available.
func NewTargetOfEvaluation(assuranceLevel string) *orchestrator.TargetOfEvaluation {
	var toe = &orchestrator.TargetOfEvaluation{
		CloudServiceId: testdata.MockCloudServiceID1,
		CatalogId:      testdata.MockCatalogID,
	}

	if assuranceLevel != "" {
		toe.AssuranceLevel = &assuranceLevel
	}

	return toe
}

func NewCloudService() *orchestrator.CloudService {
	return &orchestrator.CloudService{
		Id:                testdata.MockCloudServiceID1,
		Name:              testdata.MockCloudServiceName1,
		Description:       testdata.MockCloudServiceDescription1,
		CatalogsInScope:   []*orchestrator.Catalog{},
		ConfiguredMetrics: []*assessment.Metric{},
	}
}

func NewAnotherCloudService() *orchestrator.CloudService {
	return &orchestrator.CloudService{
		Id:                testdata.MockCloudServiceID2,
		Name:              testdata.MockCloudServiceName2,
		Description:       testdata.MockCloudServiceDescription2,
		CatalogsInScope:   []*orchestrator.Catalog{},
		ConfiguredMetrics: []*assessment.Metric{},
	}
}

var (
	MockAssessmentResultRequest1 = &orchestrator.GetAssessmentResultRequest{
		Id: testdata.MockAssessmentResult1ID,
	}
	MockAssessmentResultRequest2 = &orchestrator.GetAssessmentResultRequest{
		Id: testdata.MockAssessmentResult2ID,
	}
	MockAssessmentResult1 = &assessment.AssessmentResult{
		Id:             testdata.MockAssessmentResult1ID,
		Timestamp:      timestamppb.New(time.Unix(1, 0)),
		CloudServiceId: testdata.MockCloudServiceID1,
		MetricId:       testdata.MockMetricID1,
		Compliant:      true,
		EvidenceId:     testdata.MockEvidenceID1,
		ResourceId:     testdata.MockResourceID1,
		ResourceTypes:  []string{"Resource"},
		MetricConfiguration: &assessment.MetricConfiguration{
			Operator:       "==",
			TargetValue:    structpb.NewBoolValue(true),
			IsDefault:      true,
			MetricId:       testdata.MockMetricID1,
			CloudServiceId: testdata.MockCloudServiceID1,
		},
		ToolId: util.Ref(assessment.AssessmentToolId),
	}
	MockAssessmentResult2 = &assessment.AssessmentResult{
		Id:             testdata.MockAssessmentResult2ID,
		Timestamp:      timestamppb.New(time.Unix(1, 0)),
		CloudServiceId: testdata.MockCloudServiceID2,
		MetricId:       testdata.MockMetricID1,
		Compliant:      true,
		EvidenceId:     testdata.MockEvidenceID1,
		ResourceId:     testdata.MockResourceID1,
		ResourceTypes:  []string{"Resource"},
		MetricConfiguration: &assessment.MetricConfiguration{
			Operator:       "==",
			TargetValue:    structpb.NewBoolValue(true),
			IsDefault:      true,
			MetricId:       testdata.MockMetricID1,
			CloudServiceId: testdata.MockCloudServiceID2,
		},
		ToolId: util.Ref(assessment.AssessmentToolId),
	}
	MockAssessmentResult3 = &assessment.AssessmentResult{
		Id:             testdata.MockAssessmentResult3ID,
		Timestamp:      timestamppb.New(time.Unix(1, 0)),
		CloudServiceId: testdata.MockCloudServiceID1,
		MetricId:       testdata.MockMetricID2,
		Compliant:      false,
		EvidenceId:     testdata.MockEvidenceID1,
		ResourceId:     testdata.MockResourceID1,
		ResourceTypes:  []string{"Resource"},
		MetricConfiguration: &assessment.MetricConfiguration{
			Operator:       "==",
			TargetValue:    structpb.NewBoolValue(true),
			IsDefault:      true,
			MetricId:       testdata.MockMetricID2,
			CloudServiceId: testdata.MockCloudServiceID1,
		},
		ToolId: util.Ref(assessment.AssessmentToolId),
	}
	MockAssessmentResult4 = &assessment.AssessmentResult{
		Id:             testdata.MockAssessmentResult4ID,
		Timestamp:      timestamppb.New(time.Unix(1, 0)),
		CloudServiceId: testdata.MockCloudServiceID2,
		MetricId:       testdata.MockMetricID2,
		Compliant:      false,
		EvidenceId:     testdata.MockEvidenceID1,
		ResourceId:     testdata.MockResourceID2,
		ResourceTypes:  []string{"Resource"},
		MetricConfiguration: &assessment.MetricConfiguration{
			Operator:       "==",
			TargetValue:    structpb.NewBoolValue(true),
			IsDefault:      true,
			MetricId:       testdata.MockMetricID2,
			CloudServiceId: testdata.MockCloudServiceID2,
		},
		ToolId: util.Ref(testdata.MockAssessmentResultToolID),
	}
	MockAssessmentResults = []*assessment.AssessmentResult{MockAssessmentResult1, MockAssessmentResult2, MockAssessmentResult3, MockAssessmentResult4}

	MockControl1 = &orchestrator.Control{
		Id:                testdata.MockControlID1,
		Name:              testdata.MockControlName,
		CategoryName:      testdata.MockCategoryName,
		CategoryCatalogId: testdata.MockCatalogID,
		Description:       testdata.MockControlDescription,
		Controls: []*orchestrator.Control{
			{
				Id:                             testdata.MockSubControlID11,
				Name:                           testdata.MockSubControlName,
				CategoryName:                   testdata.MockCategoryName,
				CategoryCatalogId:              testdata.MockCatalogID,
				Description:                    testdata.MockSubControlDescription,
				AssuranceLevel:                 &testdata.AssuranceLevelBasic,
				ParentControlId:                util.Ref(testdata.MockControlID1),
				ParentControlCategoryName:      util.Ref(testdata.MockCategoryName),
				ParentControlCategoryCatalogId: util.Ref(testdata.MockCatalogID),
				Metrics: []*assessment.Metric{{
					Id:          testdata.MockMetricID1,
					Name:        testdata.MockMetricName1,
					Description: testdata.MockMetricDescription1,
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
				}}},
		}}
	MockControl11 = &orchestrator.Control{
		Id:                             testdata.MockSubControlID11,
		Name:                           testdata.MockSubControlName,
		CategoryName:                   testdata.MockCategoryName,
		CategoryCatalogId:              testdata.MockCatalogID,
		Description:                    testdata.MockSubControlDescription,
		AssuranceLevel:                 &testdata.AssuranceLevelBasic,
		ParentControlId:                util.Ref(testdata.MockControlID1),
		ParentControlCategoryName:      util.Ref(testdata.MockCategoryName),
		ParentControlCategoryCatalogId: util.Ref(testdata.MockCatalogID),
		Metrics: []*assessment.Metric{{
			Id:          testdata.MockMetricID1,
			Name:        testdata.MockMetricName1,
			Description: testdata.MockMetricDescription1,
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
		},
		}}
	MockControl2 = &orchestrator.Control{
		Id:                testdata.MockControlID2,
		Name:              testdata.MockControlName,
		CategoryName:      testdata.MockCategoryName,
		CategoryCatalogId: testdata.MockCatalogID,
		Description:       testdata.MockControlDescription,
		Controls: []*orchestrator.Control{
			{
				Id:                             testdata.MockSubControlID21,
				Name:                           testdata.MockControlName,
				CategoryName:                   testdata.MockCategoryName,
				CategoryCatalogId:              testdata.MockCatalogID,
				Description:                    testdata.MockControlDescription,
				AssuranceLevel:                 &testdata.AssuranceLevelBasic,
				ParentControlId:                util.Ref(testdata.MockControlID2),
				ParentControlCategoryName:      util.Ref(testdata.MockCategoryName),
				ParentControlCategoryCatalogId: util.Ref(testdata.MockCatalogID),
				Metrics: []*assessment.Metric{{
					Id:          testdata.MockMetricID1,
					Name:        testdata.MockMetricName1,
					Description: testdata.MockMetricDescription1,
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
				}}},
		}}
	MockControl3 = &orchestrator.Control{
		Id:                testdata.MockControlID3,
		Name:              testdata.MockControlName,
		CategoryName:      testdata.MockCategoryName,
		CategoryCatalogId: testdata.MockCatalogID,
		Description:       testdata.MockControlDescription,
		Controls: []*orchestrator.Control{
			{
				Id:                             testdata.MockSubControlID31,
				Name:                           testdata.MockControlName,
				CategoryName:                   testdata.MockCategoryName,
				CategoryCatalogId:              testdata.MockCatalogID,
				Description:                    testdata.MockControlDescription,
				AssuranceLevel:                 &testdata.AssuranceLevelSubstantial,
				ParentControlId:                util.Ref(testdata.MockControlID3),
				ParentControlCategoryName:      util.Ref(testdata.MockCategoryName),
				ParentControlCategoryCatalogId: util.Ref(testdata.MockCatalogID),
			}},
	}
	MockControl4 = &orchestrator.Control{
		Id:                testdata.MockControlID4,
		Name:              testdata.MockControlName,
		CategoryName:      testdata.MockCategoryName,
		CategoryCatalogId: testdata.MockCatalogID,
		Description:       testdata.MockControlDescription,
		AssuranceLevel:    &testdata.AssuranceLevelHigh,
		Metrics: []*assessment.Metric{{
			Id:          testdata.MockMetricID1,
			Name:        testdata.MockMetricName1,
			Description: testdata.MockMetricDescription1,
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
			Id:          testdata.MockMetricID1,
			Name:        testdata.MockMetricName1,
			Description: testdata.MockMetricDescription1,
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
	// Control without sub-control
	MockControl6 = &orchestrator.Control{
		Id:                testdata.MockControlID1,
		Name:              testdata.MockControlName,
		CategoryName:      testdata.MockCategoryName,
		CategoryCatalogId: testdata.MockCatalogID,
		Description:       testdata.MockControlDescription,
	}
	MockControls = []*orchestrator.Control{MockControl1, MockControl2, MockControl3, MockControl4, MockControl5}
)

func NewMetric() *assessment.Metric {
	return &assessment.Metric{
		Id:          testdata.MockMetricID1,
		Name:        testdata.MockMetricName1,
		Description: testdata.MockMetricDescription1,
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
		Deprecated: false,
	}

}
