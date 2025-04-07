package evidencetest

import (
	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/util"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	MockListEvidenceRequest1 = &evidence.ListEvidencesRequest{
		Filter: &evidence.Filter{
			TargetOfEvaluationId: util.Ref(testdata.MockTargetOfEvaluationID1),
			ToolId:               util.Ref(testdata.MockEvidenceToolID1),
		},
	}
	MockListEvidenceRequest2 = &evidence.ListEvidencesRequest{
		Filter: &evidence.Filter{
			TargetOfEvaluationId: util.Ref(testdata.MockTargetOfEvaluationID2),
			ToolId:               util.Ref(testdata.MockEvidenceToolID2),
		},
	}
)

var (
	MockEvidence1 = &evidence.Evidence{
		Id:                   testdata.MockEvidenceID1,
		Timestamp:            timestamppb.Now(),
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
		ToolId:               testdata.MockEvidenceToolID1,
		Resource:             nil,
	}
	MockEvidence2 = &evidence.Evidence{
		Id:                   testdata.MockEvidenceID2,
		Timestamp:            timestamppb.Now(),
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID2,
		ToolId:               testdata.MockEvidenceToolID2,
		Resource:             nil,
	}
)
