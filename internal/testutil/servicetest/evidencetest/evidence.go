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
			CloudServiceId: util.Ref(testdata.MockCloudServiceID1),
			ToolId:         util.Ref(testdata.MockEvidenceToolID1),
		},
	}
	MockListEvidenceRequest2 = &evidence.ListEvidencesRequest{
		Filter: &evidence.Filter{
			CloudServiceId: util.Ref(testdata.MockCloudServiceID2),
			ToolId:         util.Ref(testdata.MockEvidenceToolID2),
		},
	}
)

var (
	MockEvidence1 = &evidence.Evidence{
		Id:             testdata.MockEvidenceID1,
		Timestamp:      timestamppb.Now(),
		CloudServiceId: testdata.MockCloudServiceID1,
		ToolId:         testdata.MockEvidenceToolID1,
		Raw:            util.Ref("This Raw field must be of length >1"),
		Resource:       nil,
	}
	MockEvidence2 = &evidence.Evidence{
		Id:             testdata.MockEvidenceID2,
		Timestamp:      timestamppb.Now(),
		CloudServiceId: testdata.MockCloudServiceID2,
		ToolId:         testdata.MockEvidenceToolID2,
		Raw:            util.Ref(""),
		Resource:       nil,
	}
)
