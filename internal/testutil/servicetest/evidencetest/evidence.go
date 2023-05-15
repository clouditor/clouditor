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
		Filter: &evidence.Filter{
			CloudServiceId: util.Ref(testdata.MockCloudServiceID),
			ToolId:         util.Ref(testdata.MockEvidenceToolID),
		},
	}
	MockListEvidenceRequest2 = &evidence.ListEvidencesRequest{
		Filter: &evidence.Filter{
			CloudServiceId: util.Ref(testdata.MockAnotherCloudServiceID),
			ToolId:         util.Ref(testdata.MockAnotherEvidenceToolID),
		},
	}
)

var (
	MockEvidence1 = &evidence.Evidence{
		Id:             testdata.MockEvidenceID,
		Timestamp:      timestamppb.Now(),
		CloudServiceId: testdata.MockCloudServiceID,
		ToolId:         testdata.MockEvidenceToolID,
		Raw:            util.Ref("This Raw field must be of length >1"),
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
