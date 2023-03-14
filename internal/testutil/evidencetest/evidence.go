package evidencetest

import (
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/util"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	MockListEvidenceRequest1 = &evidence.ListEvidencesRequest{
		PageSize:  testdata.MockEvidencePageSize,
		PageToken: testdata.MockEvidencePageToken,
		OrderBy:   testdata.MockEvidenceOrderBy,
		Asc:       testdata.MockEvidenceAsc,
		Filter: &evidence.Filter{
			CloudServiceId: testdata.MockCloudServiceID,
			ToolId:         testdata.MockEvidenceToolID,
		},
	}
	MockListEvidenceRequest2 = &evidence.ListEvidencesRequest{
		PageSize:  testdata.MockAnotherEvidencePageSize,
		PageToken: testdata.MockAnotherEvidencePageToken,
		OrderBy:   testdata.MockAnotherEvidenceOrderBy,
		Asc:       testdata.MockAnotherEvidenceAsc,
		Filter: &evidence.Filter{
			CloudServiceId: testdata.MockAnotherCloudServiceID,
			ToolId:         testdata.MockAnotherEvidenceToolID,
		},
	}
)

var (
	MockEvidence1 = &evidence.Evidence{
		Id:             testdata.MockEvidenceID,
		Timestamp:      timestamppb.Now(),
		CloudServiceId: testdata.MockCloudServiceID,
		ToolId:         testdata.MockEvidenceToolID,
		Raw:            util.Ref(""),
		Resource:       structpb.NewNullValue(),
	}
	MockEvidence2 = &evidence.Evidence{
		Id:             testdata.MockAnotherEvidenceID,
		Timestamp:      timestamppb.Now(),
		CloudServiceId: testdata.MockAnotherCloudServiceID,
		ToolId:         testdata.MockAnotherEvidenceToolID,
		Raw:            util.Ref(""),
		Resource:       structpb.NewNullValue(),
	}
)
