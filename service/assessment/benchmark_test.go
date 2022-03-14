package assessment

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/voc"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func createSome(n int, b *testing.B) int {
	var wg sync.WaitGroup

	wg.Add(n * 7)

	svc := NewService(WithOrchestratorAddress("bufcon"))
	svc.initOrchestratorStream(grpc.WithContextDialer(bufConnDialer))
	svc.initEvidenceStoreStream(grpc.WithContextDialer(bufConnDialer))
	svc.RegisterAssessmentResultHook(func(result *assessment.AssessmentResult, err error) {
		wg.Done()
	})

	// Create evidences for n resources (1 per resource)
	for i := 0; i < n; i++ {
		r, _ := voc.ToStruct(&voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{
			ID:   voc.ResourceID(fmt.Sprintf("%d", i)),
			Type: []string{"VirtualMachine", "Compute", "Resource"},
		}}})

		e := evidence.Evidence{
			Id:        uuid.NewString(),
			Timestamp: timestamppb.Now(),
			ToolId:    "mytool",
			Resource:  r,
		}

		if i%100 == 0 || i > 2700 {
			log.Infof("Currently @ %v", i)
		}
		if i == 2782 {
			fmt.Printf("last call")
		}

		_, err := svc.AssessEvidence(context.Background(), &assessment.AssessEvidenceRequest{
			Evidence: &e,
		})
		if err != nil {
			b.Errorf("Error while calling AssessEvidence: %v", err)
		}
	}

	wg.Wait()
	svc.orchestratorStream.CloseSend()
	svc.evidenceStoreStream.CloseSend()

	return 0
}

func benchmarkAssessEvidence(i int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		createSome(i, b)
	}
}

func BenchmarkAssessEvidence1(b *testing.B) {
	benchmarkAssessEvidence(1, b)
}

func BenchmarkAssessEvidence2(b *testing.B) {
	benchmarkAssessEvidence(2, b)
}

func BenchmarkAssessEvidence10(b *testing.B) {
	benchmarkAssessEvidence(10, b)
}

func BenchmarkAssessEvidence100(b *testing.B) {
	benchmarkAssessEvidence(100, b)
}

func BenchmarkAssessEvidence1000(b *testing.B) {
	benchmarkAssessEvidence(1000, b)
}

func BenchmarkAssessEvidence3000(b *testing.B) {
	benchmarkAssessEvidence(3000, b)
}

func BenchmarkAssessEvidence10000(b *testing.B) {
	benchmarkAssessEvidence(10000, b)
}
