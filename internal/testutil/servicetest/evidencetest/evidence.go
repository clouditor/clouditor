package evidencetest

import (
	"strings"

	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/util"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

<<<<<<< HEAD
// Evidence requests for testing
=======
// Mock Requests
>>>>>>> main
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

<<<<<<< HEAD
// Evidences for testing
=======
// Mock Evidence and Resource
>>>>>>> main
var (
	MockEvidence1 = &evidence.Evidence{
		Id:                   testdata.MockEvidenceID1,
		Timestamp:            timestamppb.Now(),
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
		ToolId:               testdata.MockEvidenceToolID1,
		Resource: &ontology.Resource{
			Type: &ontology.Resource_VirtualMachine{
				VirtualMachine: &ontology.VirtualMachine{
					Id:           testdata.MockVirtualMachineID1,
					Name:         testdata.MockVirtualMachineName1,
					CreationTime: timestamppb.Now(),
					Description:  testdata.MockVirtualMachineDescription1,
					BlockStorageIds: []string{
						testdata.MockBlockStorageID1,
					},
				},
			},
		},
	}
	MockEvidence2 = &evidence.Evidence{
		Id:                   testdata.MockEvidenceID2,
		Timestamp:            timestamppb.Now(),
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID2,
		ToolId:               testdata.MockEvidenceToolID2,
		Resource: &ontology.Resource{
			Type: &ontology.Resource_BlockStorage{
				BlockStorage: &ontology.BlockStorage{
					Id:           testdata.MockBlockStorageID1,
					Name:         testdata.MockBlockStorageName1,
					CreationTime: timestamppb.Now(),
					Description:  testdata.MockBlockStorageDescription1,
					ParentId:     util.Ref(testdata.MockVirtualMachineID1),
				},
			},
		},
	}
	MockVirtualMachineResource1 = &evidence.Resource{
		Id:                   testdata.MockVirtualMachineID1,
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
		ResourceType:         strings.Join(testdata.MockVirtualMachineTypes, ","),
		ToolId:               testdata.MockEvidenceToolID2,
		Properties:           &anypb.Any{},
	}
	MockVirtualMachineResource2 = &evidence.Resource{
		Id:                   testdata.MockVirtualMachineID2,
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID2,
		ResourceType:         strings.Join(testdata.MockVirtualMachineTypes, ","),
		ToolId:               testdata.MockEvidenceToolID1,
		Properties:           &anypb.Any{},
	}
	MockBlockStorageResource1 = &evidence.Resource{
		Id:                   testdata.MockBlockStorageID1,
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
		ResourceType:         strings.Join(testdata.MockBlockStorageTypes, ","),
		ToolId:               testdata.MockEvidenceToolID1,
		Properties:           &anypb.Any{},
	}
	MockBlockStorageResource2 = &evidence.Resource{
		Id:                   testdata.MockBlockStorageID2,
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID2,
		ResourceType:         strings.Join(testdata.MockBlockStorageTypes, ","),
		ToolId:               testdata.MockEvidenceToolID1,
		Properties:           &anypb.Any{},
	}
)

// Resources (evidence.Resource) for testing
var (
	MockResourceFromEvidence1, _ = evidence.ToEvidenceResource(MockEvidence1.GetOntologyResource(), testdata.MockTargetOfEvaluationID1, testdata.MockEvidenceToolID1)

	MockResourceFromEvidence2, _ = evidence.ToEvidenceResource(MockEvidence2.GetOntologyResource(), testdata.MockTargetOfEvaluationID2, testdata.MockEvidenceToolID2)
	//nolint:lll,gosec // this is a test file, so we don't care about the linter here
)
