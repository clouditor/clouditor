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

// NewCertificate creates a mock certificate
func NewCertificate() *orchestrator.Certificate {
	var mockCertificate = &orchestrator.Certificate{
		Id:             testdata.MockCertificateID,
		Name:           testdata.MockCertificateName,
		CloudServiceId: testdata.MockCloudServiceID1,
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

	// Our test catalog does not allow scoping, so we need to emulate what we do in CreateTargetOfEvaluation
	toe.ControlsInScope = []*orchestrator.Control{
		{
			Id:                testdata.MockControlID1,
			CategoryName:      testdata.MockCategoryName,
			CategoryCatalogId: testdata.MockCatalogID,
			Name:              testdata.MockControlName,
			Description:       testdata.MockControlDescription,
		},
		{
			Id:                             testdata.MockSubControlID11,
			CategoryName:                   testdata.MockCategoryName,
			CategoryCatalogId:              testdata.MockCatalogID,
			Name:                           testdata.MockSubControlName,
			Description:                    testdata.MockSubControlDescription,
			ParentControlId:                util.Ref(testdata.MockControlID1),
			ParentControlCategoryName:      util.Ref(testdata.MockCategoryName),
			ParentControlCategoryCatalogId: util.Ref(testdata.MockCatalogID),
			AssuranceLevel:                 &testdata.AssuranceLevelBasic,
		},
	}

	return toe
}

// NewTargetOfEvaluationWithoutControlsInScope creates a new Target of Evaluation without controls_in_scope. The assurance level is set if available.
func NewTargetOfEvaluationWithoutControlsInScope(assuranceLevel string) *orchestrator.TargetOfEvaluation {
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
		ToolId: assessment.AssessmentToolId,
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
		ToolId: assessment.AssessmentToolId,
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
		ToolId: assessment.AssessmentToolId,
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
		ToolId: assessment.AssessmentToolId,
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

	MockControlsInScope1 = &orchestrator.Control{
		Id:                testdata.MockControlID1,
		Name:              testdata.MockControlName,
		CategoryName:      testdata.MockCategoryName,
		CategoryCatalogId: testdata.MockCatalogID,
		Description:       testdata.MockControlDescription,
	}
	MockControlsInScopeSubControl11 = &orchestrator.Control{
		Id:                             testdata.MockSubControlID11,
		Name:                           testdata.MockControlName,
		CategoryName:                   testdata.MockCategoryName,
		CategoryCatalogId:              testdata.MockCatalogID,
		Description:                    testdata.MockControlDescription,
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
		}}}
	MockControlsInScope2 = &orchestrator.Control{
		Id:                testdata.MockControlID2,
		Name:              testdata.MockControlName,
		CategoryName:      testdata.MockCategoryName,
		CategoryCatalogId: testdata.MockCatalogID,
		Description:       testdata.MockControlDescription,
	}
	MockControlsInScopeSubControl21 = &orchestrator.Control{
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
		}},
	}
	MockControlsInScope3 = &orchestrator.Control{
		Id:                testdata.MockControlID3,
		Name:              testdata.MockControlName,
		CategoryName:      testdata.MockCategoryName,
		CategoryCatalogId: testdata.MockCatalogID,
		Description:       testdata.MockControlDescription,
	}
	MockControlsInScopeSubControl31 = &orchestrator.Control{
		Id:                             testdata.MockSubControlID31,
		Name:                           testdata.MockControlName,
		CategoryName:                   testdata.MockCategoryName,
		CategoryCatalogId:              testdata.MockCatalogID,
		Description:                    testdata.MockControlDescription,
		AssuranceLevel:                 &testdata.AssuranceLevelSubstantial,
		ParentControlId:                util.Ref(testdata.MockControlID3),
		ParentControlCategoryName:      util.Ref(testdata.MockCategoryName),
		ParentControlCategoryCatalogId: util.Ref(testdata.MockCatalogID),
	}
	MockControlsInScopeSubControl32 = &orchestrator.Control{
		Id:                             testdata.MockSubControlID32,
		Name:                           testdata.MockControlName,
		CategoryName:                   testdata.MockCategoryName,
		CategoryCatalogId:              testdata.MockCatalogID,
		Description:                    testdata.MockControlDescription,
		AssuranceLevel:                 &testdata.AssuranceLevelHigh,
		ParentControlId:                util.Ref(testdata.MockControlID3),
		ParentControlCategoryName:      util.Ref(testdata.MockCategoryName),
		ParentControlCategoryCatalogId: util.Ref(testdata.MockCatalogID),
	}
	MockControlsInScope4 = &orchestrator.Control{
		Id:                testdata.MockControlID4,
		Name:              testdata.MockControlName,
		CategoryName:      testdata.MockCategoryName,
		CategoryCatalogId: testdata.MockCatalogID,
		Description:       testdata.MockControlDescription,
		Metrics:           []*assessment.Metric{},
	}
	MockControlsInScopeSubControl4 = &orchestrator.Control{
		Id:                             testdata.MockSubControlID32,
		Name:                           testdata.MockControlName,
		CategoryName:                   testdata.MockCategoryName,
		CategoryCatalogId:              testdata.MockCatalogID,
		Description:                    testdata.MockControlDescription,
		AssuranceLevel:                 &testdata.AssuranceLevelSubstantial,
		ParentControlId:                util.Ref(testdata.MockControlID3),
		ParentControlCategoryName:      util.Ref(testdata.MockCategoryName),
		ParentControlCategoryCatalogId: util.Ref(testdata.MockCatalogID),
	}
	MockControlsInScope5 = &orchestrator.Control{
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
	MockControlsInScope = []*orchestrator.Control{MockControlsInScope1, MockControlsInScopeSubControl11, MockControlsInScope2, MockControlsInScopeSubControl21, MockControlsInScope3, MockControlsInScopeSubControl31, MockControlsInScope4, MockControlsInScope5, MockControlsInScopeSubControl32}
)
