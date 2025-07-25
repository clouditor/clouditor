package orchestratortest

import (
	"time"

	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/util"
	"github.com/google/uuid"
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
		Id:                   testdata.MockCertificateID,
		Name:                 testdata.MockCertificateName,
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
		IssueDate:            timeStamp.AddDate(-5, 0, 0).String(),
		ExpirationDate:       timeStamp.AddDate(5, 0, 0).String(),
		Standard:             testdata.MockCertificateName,
		AssuranceLevel:       testdata.AssuranceLevelHigh,
		Cab:                  testdata.MockCertificateCab,
		Description:          testdata.MockCertificateDescription,
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
		Id:                   testdata.MockCertificateID2,
		Name:                 testdata.MockCertificateName2,
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID2,
		IssueDate:            timeStamp.AddDate(-5, 0, 0).String(),
		ExpirationDate:       timeStamp.AddDate(5, 0, 0).String(),
		Standard:             testdata.MockCertificateName2,
		AssuranceLevel:       testdata.AssuranceLevelHigh,
		Cab:                  testdata.MockCertificateCab2,
		Description:          testdata.MockCertificateDescription2,
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
		Id:              testdata.MockCatalogID1,
		Description:     testdata.MockCatalogDescription,
		AllInScope:      true,
		AssuranceLevels: []string{testdata.AssuranceLevelBasic, testdata.AssuranceLevelSubstantial, testdata.AssuranceLevelHigh},
		Categories: []*orchestrator.Category{{
			Name:        testdata.MockCategoryName,
			Description: testdata.MockCategoryDescription,
			CatalogId:   testdata.MockCatalogID1,
			Controls: []*orchestrator.Control{
				MockControl1,
				MockControl2,
			},
		}}}

	return mockCatalog
}

// NewAuditScope creates a new Audit Scope. The assurance level is set if available. A different audit scope id is set, if available.
func NewAuditScope(assuranceLevel, auditScopeId, targetOfEvaluationID, auditScopeName string) *orchestrator.AuditScope {
	var auditScope = &orchestrator.AuditScope{
		Id:                   auditScopeId,
		Name:                 auditScopeName,
		TargetOfEvaluationId: targetOfEvaluationID,
		CatalogId:            testdata.MockCatalogID1,
	}

	if assuranceLevel != "" {
		auditScope.AssuranceLevel = &assuranceLevel
	}

	if auditScopeId == "" {
		auditScope.Id = uuid.NewString()
	}

	if targetOfEvaluationID == "" {
		auditScope.TargetOfEvaluationId = testdata.MockTargetOfEvaluationID1
	}

	return auditScope
}

func NewTargetOfEvaluation() *orchestrator.TargetOfEvaluation {
	return &orchestrator.TargetOfEvaluation{
		Id:                testdata.MockTargetOfEvaluationID1,
		Name:              testdata.MockTargetOfEvaluationName1,
		Description:       testdata.MockTargetOfEvaluationDescription1,
		ConfiguredMetrics: []*assessment.Metric{},
	}
}

func NewAnotherTargetOfEvaluation() *orchestrator.TargetOfEvaluation {
	return &orchestrator.TargetOfEvaluation{
		Id:                testdata.MockTargetOfEvaluationID2,
		Name:              testdata.MockTargetOfEvaluationName2,
		Description:       testdata.MockTargetOfEvaluationDescription2,
		ConfiguredMetrics: []*assessment.Metric{},
	}
}

var (
	// MockAuditScopeCertTargetID1 has the Id 'MockAuditScopeID1' and TargetOfEvaluationID 'TargetOfEvaluationID1'
	MockAuditScopeCertTargetID1 = &orchestrator.AuditScope{
		Id:                   testdata.MockAuditScopeID1,
		Name:                 testdata.MockAuditScopeName1,
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
		CatalogId:            testdata.MockCatalogID1,
		AssuranceLevel:       &testdata.AssuranceLevelBasic,
	}

	// MockAuditScopeCertTargetID2 has the Id 'MockAuditScopeID2' and TargetOfEvaluationID 'TargetOfEvaluationID1'
	MockAuditScopeCertTargetID2 = &orchestrator.AuditScope{
		Id:                   testdata.MockAuditScopeID2,
		Name:                 testdata.MockAuditScopeName2,
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
		CatalogId:            testdata.MockCatalogID2,
		AssuranceLevel:       &testdata.AssuranceLevelBasic,
	}

	// MockAuditScopeCertTargetID3 has the Id 'MockAuditScopeID3' and TargetOfEvaluationID 'TargetOfEvaluationID2'
	MockAuditScopeCertTargetID3 = &orchestrator.AuditScope{
		Id:                   testdata.MockAuditScopeID3,
		Name:                 testdata.MockAuditScopeName3,
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID2,
		CatalogId:            testdata.MockCatalogID1,
		AssuranceLevel:       &testdata.AssuranceLevelBasic,
	}

	MockAssessmentResultRequest1 = &orchestrator.GetAssessmentResultRequest{
		Id: testdata.MockAssessmentResult1ID,
	}
	MockAssessmentResultRequest2 = &orchestrator.GetAssessmentResultRequest{
		Id: testdata.MockAssessmentResult2ID,
	}
	MockAssessmentResult1 = &assessment.AssessmentResult{
		Id:                   testdata.MockAssessmentResult1ID,
		CreatedAt:            timestamppb.New(time.Unix(1, 0)),
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
		MetricId:             testdata.MockMetricID1,
		Compliant:            true,
		EvidenceId:           testdata.MockEvidenceID1,
		ResourceId:           testdata.MockVirtualMachineID1,
		ResourceTypes:        testdata.MockVirtualMachineTypes,
		ComplianceComment:    assessment.DefaultCompliantMessage,
		MetricConfiguration: &assessment.MetricConfiguration{
			Operator:             "==",
			TargetValue:          structpb.NewBoolValue(true),
			IsDefault:            true,
			MetricId:             testdata.MockMetricID1,
			TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
		},
		ToolId:           util.Ref(assessment.AssessmentToolId),
		HistoryUpdatedAt: timestamppb.New(time.Unix(1, 0)),
		History: []*assessment.Record{
			{
				EvidenceRecordedAt: timestamppb.New(time.Unix(1, 0)),
				EvidenceId:         testdata.MockEvidenceID1,
			},
		},
	}
	MockAssessmentResult2 = &assessment.AssessmentResult{
		Id:                   testdata.MockAssessmentResult2ID,
		CreatedAt:            timestamppb.New(time.Unix(1, 0)),
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID2,
		MetricId:             testdata.MockMetricID1,
		Compliant:            true,
		EvidenceId:           testdata.MockEvidenceID1,
		ResourceId:           testdata.MockVirtualMachineID1,
		ResourceTypes:        testdata.MockVirtualMachineTypes,
		ComplianceComment:    assessment.DefaultCompliantMessage,
		MetricConfiguration: &assessment.MetricConfiguration{
			Operator:             "==",
			TargetValue:          structpb.NewBoolValue(true),
			IsDefault:            true,
			MetricId:             testdata.MockMetricID1,
			TargetOfEvaluationId: testdata.MockTargetOfEvaluationID2,
		},
		ToolId:           util.Ref(assessment.AssessmentToolId),
		HistoryUpdatedAt: timestamppb.New(time.Unix(1, 0)),
		History: []*assessment.Record{
			{
				EvidenceRecordedAt: timestamppb.New(time.Unix(1, 0)),
				EvidenceId:         testdata.MockEvidenceID1,
			},
		},
	}
	MockAssessmentResult3 = &assessment.AssessmentResult{
		Id:                   testdata.MockAssessmentResult3ID,
		CreatedAt:            timestamppb.New(time.Unix(1, 0)),
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
		MetricId:             testdata.MockMetricID2,
		Compliant:            false,
		EvidenceId:           testdata.MockEvidenceID1,
		ResourceId:           testdata.MockVirtualMachineID1,
		ResourceTypes:        testdata.MockVirtualMachineTypes,
		ComplianceComment:    assessment.DefaultNonCompliantMessage,
		MetricConfiguration: &assessment.MetricConfiguration{
			Operator:             "==",
			TargetValue:          structpb.NewBoolValue(true),
			IsDefault:            true,
			MetricId:             testdata.MockMetricID2,
			TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
		},
		ToolId:           util.Ref(assessment.AssessmentToolId),
		HistoryUpdatedAt: timestamppb.New(time.Unix(1, 0)),
		History: []*assessment.Record{
			{
				EvidenceRecordedAt: timestamppb.New(time.Unix(1, 0)),
				EvidenceId:         testdata.MockEvidenceID1,
			},
		},
	}
	MockAssessmentResult4 = &assessment.AssessmentResult{
		Id:                   testdata.MockAssessmentResult4ID,
		CreatedAt:            timestamppb.New(time.Unix(1, 0)),
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID2,
		MetricId:             testdata.MockMetricID2,
		Compliant:            false,
		EvidenceId:           testdata.MockEvidenceID1,
		ResourceId:           testdata.MockVirtualMachineID2,
		ResourceTypes:        testdata.MockVirtualMachineTypes,
		ComplianceComment:    assessment.DefaultNonCompliantMessage,
		MetricConfiguration: &assessment.MetricConfiguration{
			Operator:             "==",
			TargetValue:          structpb.NewBoolValue(true),
			IsDefault:            true,
			MetricId:             testdata.MockMetricID2,
			TargetOfEvaluationId: testdata.MockTargetOfEvaluationID2,
		},
		ToolId:           util.Ref(testdata.MockAssessmentResultToolID),
		HistoryUpdatedAt: timestamppb.New(time.Unix(1, 0)),
		History: []*assessment.Record{
			{
				EvidenceRecordedAt: timestamppb.New(time.Unix(1, 0)),
				EvidenceId:         testdata.MockEvidenceID1,
			},
		},
	}
	MockAssessmentResults = []*assessment.AssessmentResult{MockAssessmentResult1, MockAssessmentResult2, MockAssessmentResult3, MockAssessmentResult4}

	MockControl1 = &orchestrator.Control{
		Id:                testdata.MockControlID1,
		Name:              testdata.MockControlName,
		CategoryName:      testdata.MockCategoryName,
		CategoryCatalogId: testdata.MockCatalogID1,
		Description:       testdata.MockControlDescription,
		Controls: []*orchestrator.Control{
			{
				Id:                             testdata.MockSubControlID11,
				Name:                           testdata.MockSubControlName,
				CategoryName:                   testdata.MockCategoryName,
				CategoryCatalogId:              testdata.MockCatalogID1,
				Description:                    testdata.MockSubControlDescription,
				AssuranceLevel:                 &testdata.AssuranceLevelBasic,
				ParentControlId:                util.Ref(testdata.MockControlID1),
				ParentControlCategoryName:      util.Ref(testdata.MockCategoryName),
				ParentControlCategoryCatalogId: util.Ref(testdata.MockCatalogID1),
				Metrics: []*assessment.Metric{{
					Id:          testdata.MockMetricID1,
					Description: testdata.MockMetricDescription1,
					Category:    testdata.MockMetricCategory1,
					Version:     "1.0",
					Comments:    testdata.MockMetricComments1,
				}}},
		}}
	MockControl11 = &orchestrator.Control{
		Id:                             testdata.MockSubControlID11,
		Name:                           testdata.MockSubControlName,
		CategoryName:                   testdata.MockCategoryName,
		CategoryCatalogId:              testdata.MockCatalogID1,
		Description:                    testdata.MockSubControlDescription,
		AssuranceLevel:                 &testdata.AssuranceLevelBasic,
		ParentControlId:                util.Ref(testdata.MockControlID1),
		ParentControlCategoryName:      util.Ref(testdata.MockCategoryName),
		ParentControlCategoryCatalogId: util.Ref(testdata.MockCatalogID1),
		Metrics: []*assessment.Metric{{
			Id:          testdata.MockMetricID1,
			Description: testdata.MockMetricDescription1,
			Category:    testdata.MockMetricCategory1,
			Version:     "1.0",
			Comments:    testdata.MockMetricComments1,
		},
		}}
	MockControl2 = &orchestrator.Control{
		Id:                testdata.MockControlID2,
		Name:              testdata.MockControlName,
		CategoryName:      testdata.MockCategoryName,
		CategoryCatalogId: testdata.MockCatalogID1,
		Description:       testdata.MockControlDescription,
		Controls: []*orchestrator.Control{
			{
				Id:                             testdata.MockSubControlID21,
				Name:                           testdata.MockControlName,
				CategoryName:                   testdata.MockCategoryName,
				CategoryCatalogId:              testdata.MockCatalogID1,
				Description:                    testdata.MockControlDescription,
				AssuranceLevel:                 &testdata.AssuranceLevelBasic,
				ParentControlId:                util.Ref(testdata.MockControlID2),
				ParentControlCategoryName:      util.Ref(testdata.MockCategoryName),
				ParentControlCategoryCatalogId: util.Ref(testdata.MockCatalogID1),
				Metrics: []*assessment.Metric{{
					Id:          testdata.MockMetricID1,
					Description: testdata.MockMetricDescription1,
					Category:    testdata.MockMetricCategory1,
					Version:     "1.0",
					Comments:    "This is a comment",
				}},
			},
		},
	}
	MockControl3 = &orchestrator.Control{
		Id:                testdata.MockControlID3,
		Name:              testdata.MockControlName,
		CategoryName:      testdata.MockCategoryName,
		CategoryCatalogId: testdata.MockCatalogID1,
		Description:       testdata.MockControlDescription,
		Controls: []*orchestrator.Control{
			{
				Id:                             testdata.MockSubControlID31,
				Name:                           testdata.MockControlName,
				CategoryName:                   testdata.MockCategoryName,
				CategoryCatalogId:              testdata.MockCatalogID1,
				Description:                    testdata.MockControlDescription,
				AssuranceLevel:                 &testdata.AssuranceLevelSubstantial,
				ParentControlId:                util.Ref(testdata.MockControlID3),
				ParentControlCategoryName:      util.Ref(testdata.MockCategoryName),
				ParentControlCategoryCatalogId: util.Ref(testdata.MockCatalogID1),
			}},
	}
	MockControl4 = &orchestrator.Control{
		Id:                testdata.MockControlID4,
		Name:              testdata.MockControlName,
		CategoryName:      testdata.MockCategoryName,
		CategoryCatalogId: testdata.MockCatalogID1,
		Description:       testdata.MockControlDescription,
		AssuranceLevel:    &testdata.AssuranceLevelHigh,
		Metrics: []*assessment.Metric{{
			Id:          testdata.MockMetricID1,
			Description: testdata.MockMetricDescription1,
			Category:    testdata.MockMetricCategory1,
			Version:     "1.0",
			Comments:    "This is a comment",
		}},
	}
	MockControl5 = &orchestrator.Control{
		Id:                testdata.MockControlID5,
		Name:              testdata.MockControlName,
		CategoryName:      testdata.MockCategoryName,
		CategoryCatalogId: testdata.MockCatalogID1,
		Description:       testdata.MockControlDescription,
		AssuranceLevel:    nil,
		Metrics: []*assessment.Metric{{
			Id:          testdata.MockMetricID1,
			Description: testdata.MockMetricDescription1,
			Category:    testdata.MockMetricCategory1,
			Version:     "1.0",
			Comments:    "This is a comment",
		}},
	}
	// Control without sub-control
	MockControl6 = &orchestrator.Control{
		Id:                testdata.MockControlID1,
		Name:              testdata.MockControlName,
		CategoryName:      testdata.MockCategoryName,
		CategoryCatalogId: testdata.MockCatalogID1,
		Description:       testdata.MockControlDescription,
	}
	MockControls = []*orchestrator.Control{MockControl1, MockControl2, MockControl3, MockControl4, MockControl5}
)

func NewMetric() *assessment.Metric {
	return &assessment.Metric{
		Id:          testdata.MockMetricID1,
		Description: testdata.MockMetricDescription1,
		Version:     "1.0",
		Comments:    "This is a comment",
	}

}
