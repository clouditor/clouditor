package assessment

import (
	"context"
	"testing"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/prototest"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestService_AssessEvidenceWaitFor(t *testing.T) {
	s := NewService()
	s.evidenceStore = api.NewRPCConnection("bufnet", evidence.NewEvidenceStoreClient, grpc.WithContextDialer(bufConnDialer))
	s.orchestrator = api.NewRPCConnection("bufnet", orchestrator.NewOrchestratorClient, grpc.WithContextDialer(bufConnDialer))

	// Add our first evidence
	resp, err := s.AssessEvidence(context.Background(), &assessment.AssessEvidenceRequest{
		Evidence: &evidence.Evidence{
			Id: "11111111-1111-1111-1111-111111111111",
			Resource: prototest.NewAny(t, &ontology.VirtualMachine{
				Id:              "my-resource",
				Name:            "my resource",
				BlockStorageIds: []string{"my-third-resource"},
			}),
			CloudServiceId:                 testdata.MockCloudServiceID1,
			ToolId:                         "my-tool",
			Timestamp:                      timestamppb.Now(),
			ExperimentalRelatedResourceIds: []string{"my-third-resource"},
		},
	})

	assert.ErrorIs(t, err, nil)
	assert.Equal(t, assessment.AssessmentStatus_ASSESSMENT_STATUS_WAITING_FOR_RELATED, resp.Status)

	// For more fun, we add a second evidence also waiting for the same resource
	resp, err = s.AssessEvidence(context.Background(), &assessment.AssessEvidenceRequest{
		Evidence: &evidence.Evidence{
			Id: "22222222-2222-2222-2222-222222222222",
			Resource: prototest.NewAny(t, &ontology.VirtualMachine{
				Id:              "my-other-resource",
				Name:            "my other resource",
				BlockStorageIds: []string{"my-third-resource"},
			}),
			CloudServiceId:                 testdata.MockCloudServiceID1,
			ToolId:                         "my-tool",
			Timestamp:                      timestamppb.Now(),
			ExperimentalRelatedResourceIds: []string{"my-third-resource"},
		},
	})

	assert.ErrorIs(t, err, nil)
	assert.Equal(t, assessment.AssessmentStatus_ASSESSMENT_STATUS_WAITING_FOR_RELATED, resp.Status)

	// Finally, an evidence for the resource we are waiting for. For even more fun, this evidence
	// also depends on the first resource, creating a mutual dependency.
	resp, err = s.AssessEvidence(context.Background(), &assessment.AssessEvidenceRequest{
		Evidence: &evidence.Evidence{
			Id: "33333333-3333-3333-3333-333333333333",
			Resource: prototest.NewAny(t, &ontology.BlockStorage{
				Id:   "my-third-resource",
				Name: "my third resource",
			}),
			CloudServiceId:                 testdata.MockCloudServiceID1,
			ToolId:                         "my-tool",
			Timestamp:                      timestamppb.Now(),
			ExperimentalRelatedResourceIds: []string{"my-resource"},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, assessment.AssessmentStatus_ASSESSMENT_STATUS_ASSESSED, resp.Status)

	s.wg.Wait()

	// No requests should be left over
	assert.Empty(t, s.requests)
}
