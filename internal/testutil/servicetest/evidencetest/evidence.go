package evidencetest

import (
	"strings"

	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testdata"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Mock Requests
var (
	MockListEvidenceRequest1 = &evidence.ListEvidencesRequest{
		Filter: &evidence.Filter{
			TargetOfEvaluationId: new(testdata.MockTargetOfEvaluationID1),
			ToolId:               new(testdata.MockEvidenceToolID1),
		},
	}
	MockListEvidenceRequest2 = &evidence.ListEvidencesRequest{
		Filter: &evidence.Filter{
			TargetOfEvaluationId: new(testdata.MockTargetOfEvaluationID2),
			ToolId:               new(testdata.MockEvidenceToolID2),
		},
	}
)

// Mock Evidence and Resource
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
	MockEvidence3 = &evidence.Evidence{
		Id:                   testdata.MockEvidenceID2,
		Timestamp:            timestamppb.Now(),
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID2,
		ToolId:               testdata.MockEvidenceToolID2,
		Resource: &ontology.Resource{
			Type: &ontology.Resource_VirtualMachine{
				VirtualMachine: &ontology.VirtualMachine{
					Id:   new(testdata.MockVirtualMachineID1),
					Name: new(testdata.MockVirtualMachineName1),
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
