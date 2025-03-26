package evaluationtest

import (
	"time"

	"clouditor.io/clouditor/v2/api/evaluation"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/util"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	MockEvaluationResult1 = &evaluation.EvaluationResult{
		Id:                   testdata.MockEvaluationResult1ID,
		Timestamp:            timestamppb.New(time.Unix(5, 0)),
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
		AuditScopeId:         testdata.MockAuditScopeID1,
		ControlCategoryName:  testdata.MockCategoryName,
		ControlCatalogId:     testdata.MockCatalogID1,
		Status:               evaluation.EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT,
		ControlId:            testdata.MockControlID1,
		AssessmentResultIds:  []string{testdata.MockAssessmentResult1ID},
	}
	MockEvaluationResult2 = &evaluation.EvaluationResult{
		Id:                   testdata.MockEvaluationResult2ID,
		Timestamp:            timestamppb.New(time.Unix(3, 0)),
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
		AuditScopeId:         testdata.MockAuditScopeID1,
		ControlCategoryName:  testdata.MockCategoryName,
		ControlCatalogId:     testdata.MockCatalogID1,
		Status:               evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT,
		ControlId:            testdata.MockSubControlID11,
		ParentControlId:      util.Ref(testdata.MockControlID1),
		AssessmentResultIds:  []string{},
	}
	MockEvaluationResult22 = &evaluation.EvaluationResult{
		Id:                   testdata.MockEvaluationResult9ID,
		Timestamp:            timestamppb.New(time.Unix(5, 0)),
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
		AuditScopeId:         testdata.MockAuditScopeID1,
		ControlCategoryName:  testdata.MockCategoryName,
		ControlCatalogId:     testdata.MockCatalogID1,
		Status:               evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT,
		ControlId:            testdata.MockSubControlID11,
		ParentControlId:      util.Ref(testdata.MockControlID1),
		AssessmentResultIds:  []string{},
	}
	MockEvaluationResult3 = &evaluation.EvaluationResult{
		Id:                   testdata.MockEvaluationResult3ID,
		Timestamp:            timestamppb.New(time.Unix(1, 0)),
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
		AuditScopeId:         testdata.MockAuditScopeID1,
		ControlCategoryName:  testdata.MockCategoryName,
		ControlCatalogId:     testdata.MockCatalogID1,
		Status:               evaluation.EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT,
		ControlId:            testdata.MockSubControlID12,
		ParentControlId:      util.Ref(testdata.MockControlID1),
		AssessmentResultIds:  []string{},
	}
	MockEvaluationResult4 = &evaluation.EvaluationResult{
		Id:                   testdata.MockEvaluationResult4ID,
		Timestamp:            timestamppb.New(time.Unix(1, 0)),
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
		AuditScopeId:         testdata.MockAuditScopeID1,
		ControlCategoryName:  testdata.MockCategoryName,
		ControlCatalogId:     testdata.MockCatalogID1,
		Status:               evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT,
		ControlId:            testdata.MockControlID2,
		AssessmentResultIds:  []string{},
	}
	MockEvaluationResult5 = &evaluation.EvaluationResult{
		Id:                   testdata.MockEvaluationResult5ID,
		Timestamp:            timestamppb.New(time.Unix(3, 0)),
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
		AuditScopeId:         testdata.MockAuditScopeID1,
		ControlCategoryName:  testdata.MockCategoryName,
		ControlCatalogId:     testdata.MockCatalogID1,
		Status:               evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT,
		ControlId:            testdata.MockSubControlID21,
		ParentControlId:      util.Ref(testdata.MockControlID2),
		AssessmentResultIds:  []string{},
	}
	MockEvaluationResult6 = &evaluation.EvaluationResult{
		Id:                   testdata.MockEvaluationResult6ID,
		Timestamp:            timestamppb.New(time.Unix(3, 0)),
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
		AuditScopeId:         testdata.MockAuditScopeID1,
		ControlCategoryName:  testdata.MockCategoryName,
		ControlCatalogId:     testdata.MockCatalogID1,
		Status:               evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT,
		ControlId:            testdata.MockSubControlID22,
		ParentControlId:      util.Ref(testdata.MockControlID2),
		AssessmentResultIds:  []string{},
	}

	MockEvaluationResults = []*evaluation.EvaluationResult{
		MockEvaluationResult1,
		MockEvaluationResult2,
		MockEvaluationResult22,
		MockEvaluationResult3,
		MockEvaluationResult4,
		MockEvaluationResult5,
		MockEvaluationResult6,
	}

	MockEvaluationResultsWithoutResultsForParentControl = []*evaluation.EvaluationResult{
		MockEvaluationResult2,
		MockEvaluationResult22,
		MockEvaluationResult3,
		MockEvaluationResult5,
		MockEvaluationResult6,
	}
)
