package evaluationtest

import (
	"time"

	"clouditor.io/clouditor/api/evaluation"
	"clouditor.io/clouditor/internal/testdata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	MockEvaluationResult1 = &evaluation.EvaluationResult{
		Id:                         testdata.MockEvaluationResult1ID,
		Timestamp:                  timestamppb.New(time.Unix(5, 0)),
		CloudServiceId:             testdata.MockCloudServiceID,
		ControlCategoryName:        testdata.MockCategoryName,
		ControlCatalogId:           testdata.MockCatalogID,
		ResourceId:                 testdata.MockResourceID,
		Status:                     evaluation.EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT,
		ControlId:                  testdata.MockControlID1,
		FailingAssessmentResultIds: []string{},
	}
	MockEvaluationResult2 = &evaluation.EvaluationResult{
		Id:                         testdata.MockEvaluationResult2ID,
		Timestamp:                  timestamppb.New(time.Unix(3, 0)),
		CloudServiceId:             testdata.MockCloudServiceID,
		ControlCategoryName:        testdata.MockCategoryName,
		ControlCatalogId:           testdata.MockCatalogID,
		ResourceId:                 testdata.MockResourceID,
		Status:                     evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT,
		ControlId:                  testdata.MockSubControlID11,
		FailingAssessmentResultIds: []string{},
	}
	MockEvaluationResult22 = &evaluation.EvaluationResult{
		Id:                         testdata.MockEvaluationResult9ID,
		Timestamp:                  timestamppb.New(time.Unix(5, 0)),
		CloudServiceId:             testdata.MockCloudServiceID,
		ControlCategoryName:        testdata.MockCategoryName,
		ControlCatalogId:           testdata.MockCatalogID,
		ResourceId:                 testdata.MockResourceID,
		Status:                     evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT,
		ControlId:                  testdata.MockSubControlID11,
		FailingAssessmentResultIds: []string{},
	}
	MockEvaluationResult3 = &evaluation.EvaluationResult{
		Id:                         testdata.MockEvaluationResult3ID,
		Timestamp:                  timestamppb.New(time.Unix(1, 0)),
		CloudServiceId:             testdata.MockCloudServiceID,
		ControlCategoryName:        testdata.MockCategoryName,
		ControlCatalogId:           testdata.MockCatalogID,
		ResourceId:                 testdata.MockResourceID,
		Status:                     evaluation.EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT,
		ControlId:                  testdata.MockSubControlID12,
		FailingAssessmentResultIds: []string{},
	}
	MockEvaluationResult4 = &evaluation.EvaluationResult{
		Id:                         testdata.MockEvaluationResult4ID,
		Timestamp:                  timestamppb.New(time.Unix(1, 0)),
		CloudServiceId:             testdata.MockCloudServiceID,
		ControlCategoryName:        testdata.MockCategoryName,
		ControlCatalogId:           testdata.MockCatalogID,
		ResourceId:                 testdata.MockResourceID,
		Status:                     evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT,
		ControlId:                  testdata.MockControlID2,
		FailingAssessmentResultIds: []string{},
	}
	MockEvaluationResult5 = &evaluation.EvaluationResult{
		Id:                         testdata.MockEvaluationResult5ID,
		Timestamp:                  timestamppb.New(time.Unix(3, 0)),
		CloudServiceId:             testdata.MockCloudServiceID,
		ControlCategoryName:        testdata.MockCategoryName,
		ControlCatalogId:           testdata.MockCatalogID,
		ResourceId:                 testdata.MockResourceID,
		Status:                     evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT,
		ControlId:                  testdata.MockSubControlID21,
		FailingAssessmentResultIds: []string{},
	}
	MockEvaluationResult6 = &evaluation.EvaluationResult{
		Id:                         testdata.MockEvaluationResult6ID,
		Timestamp:                  timestamppb.New(time.Unix(3, 0)),
		CloudServiceId:             testdata.MockCloudServiceID,
		ControlCategoryName:        testdata.MockCategoryName,
		ControlCatalogId:           testdata.MockCatalogID,
		ResourceId:                 testdata.MockResourceID,
		Status:                     evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT,
		ControlId:                  testdata.MockSubControlID22,
		FailingAssessmentResultIds: []string{},
	}
	MockEvaluationResult7 = &evaluation.EvaluationResult{
		Id:                         testdata.MockEvaluationResult7ID,
		Timestamp:                  timestamppb.New(time.Unix(3, 0)),
		CloudServiceId:             testdata.MockAnotherCloudServiceID,
		ControlCategoryName:        testdata.MockCategoryName,
		ControlCatalogId:           testdata.MockCatalogID,
		ResourceId:                 testdata.MockAnotherResourceID,
		Status:                     evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT,
		ControlId:                  testdata.MockControlID1,
		FailingAssessmentResultIds: []string{},
	}
	MockEvaluationResult8 = &evaluation.EvaluationResult{
		Id:                         testdata.MockEvaluationResult8ID,
		Timestamp:                  timestamppb.New(time.Unix(2, 0)),
		CloudServiceId:             testdata.MockAnotherCloudServiceID,
		ControlCategoryName:        testdata.MockCategoryName,
		ControlCatalogId:           testdata.MockCatalogID,
		ResourceId:                 testdata.MockAnotherResourceID,
		Status:                     evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT,
		ControlId:                  testdata.MockSubControlID11,
		FailingAssessmentResultIds: []string{},
	}
	MockEvaluationResults = []*evaluation.EvaluationResult{MockEvaluationResult1, MockEvaluationResult2, MockEvaluationResult22, MockEvaluationResult3, MockEvaluationResult4, MockEvaluationResult5, MockEvaluationResult6, MockEvaluationResult7, MockEvaluationResult8}

	MockEvaluationResultsWithoutResultsForParentControl = []*evaluation.EvaluationResult{MockEvaluationResult2, MockEvaluationResult22, MockEvaluationResult3, MockEvaluationResult5, MockEvaluationResult6, MockEvaluationResult8}
)
