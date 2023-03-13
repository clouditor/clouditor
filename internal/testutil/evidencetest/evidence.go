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
		PageSize:  testdata.MockPageSize,
		PageToken: testdata.MockPageToken,
		OrderBy:   testdata.MockOrderBy,
		Asc:       testdata.MockAsc,
		Filter: &evidence.Filter{
			CloudServiceId: testdata.MockCloudServiceID,
			ToolId:         testdata.MockEvidenceToolID,
		},
	}
	MockListEvidenceRequest2 = &evidence.ListEvidencesRequest{
		PageSize:  testdata.MockAnotherPageSize,
		PageToken: testdata.MockAnotherPageToken,
		OrderBy:   testdata.MockAnotherOrderBy,
		Asc:       testdata.MockAnotherAsc,
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
)
