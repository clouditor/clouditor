package assessment

import (
	"context"
	"fmt"
	"math"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"

	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/api/orchestrator"
	service_evidence "clouditor.io/clouditor/v2/service/evidence"
	service_orchestrator "clouditor.io/clouditor/v2/service/orchestrator"
	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func createEvidences(n int, m int, b *testing.B) int {
	var (
		wg   sync.WaitGroup
		err  error
		sock net.Listener
		port uint16
	)

	logrus.SetLevel(logrus.PanicLevel)

	srv := grpc.NewServer()

	orchestratorService := service_orchestrator.NewService()
	orchestrator.RegisterOrchestratorServer(srv, orchestratorService)

	evidenceService := service_evidence.NewService()
	evidence.RegisterEvidenceStoreServer(srv, evidenceService)

	sock, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		b.Errorf("could not listen: %v", err)
	}

	go func() {
		err := srv.Serve(sock)
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error while creating gRPC server: %v", err)
		}
	}()
	defer srv.Stop()

	wg.Add(n * m * 7)

	var count int64 = 0

	port = sock.Addr().(*net.TCPAddr).AddrPort().Port()

	addr := fmt.Sprintf("localhost:%d", port)

	svc := NewService(WithOrchestratorAddress(addr), WithEvidenceStoreAddress(addr))

	orchestratorService.RegisterAssessmentResultHook(func(ctx context.Context, result *assessment.AssessmentResult, err error) {
		wg.Done()
		current := atomic.AddInt64(&count, 1)
		log.Debugf("Current count: %v", current)
	})

	// Create m parallel executions of our evidence creation
	for j := 0; j < m; j++ {
		go func() {
			// Create evidences for n resources (1 per resource)
			for i := 0; i < n; i++ {
				a, _ := anypb.New(&ontology.VirtualMachine{
					Id: fmt.Sprintf("%d", i),
				})

				e := evidence.Evidence{
					Id:        uuid.NewString(),
					Timestamp: timestamppb.Now(),
					ToolId:    "mytool",
					Resource:  a,
				}

				if i%100 == 0 {
					log.Infof("Currently @ %v", i)
				}

				_, err := svc.AssessEvidence(context.Background(), connect.NewRequest(&assessment.AssessEvidenceRequest{
					Evidence: &e,
				}))
				if err != nil {
					b.Errorf("Error while calling AssessEvidence: %v", err)
				}
			}
		}()
	}

	wg.Wait()

	return 0
}

var numEvidences = []int{10, 100, 1000, 5000, 10000, 20000, 30000, 40000, 50000}

func BenchmarkAssessEvidence(b *testing.B) {
	for _, k := range numEvidences {
		for l := 0.; l <= 2; l++ {
			parallel := int(math.Pow(2, l))
			b.Run(fmt.Sprintf("%d/%d", k, parallel), func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					createEvidences(k, parallel, b)
				}
			})
		}
	}
}

func benchmarkAssessEvidenceInternal(i int, m int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		createEvidences(i, m, b)
	}
}

func BenchmarkAssessEvidence1(b *testing.B) {
	benchmarkAssessEvidenceInternal(1, 1, b)
}

func BenchmarkAssessEvidence2(b *testing.B) {
	benchmarkAssessEvidenceInternal(2, 1, b)
}

func BenchmarkAssessEvidence10(b *testing.B) {
	benchmarkAssessEvidenceInternal(10, 1, b)
}

func BenchmarkAssessEvidence10x2(b *testing.B) {
	benchmarkAssessEvidenceInternal(10, 2, b)
}

func BenchmarkAssessEvidence100(b *testing.B) {
	benchmarkAssessEvidenceInternal(100, 1, b)
}

func BenchmarkAssessEvidence1000(b *testing.B) {
	benchmarkAssessEvidenceInternal(1000, 1, b)
}

func BenchmarkAssessEvidence1000x2(b *testing.B) {
	benchmarkAssessEvidenceInternal(1000, 2, b)
}

func BenchmarkAssessEvidence1000x10(b *testing.B) {
	benchmarkAssessEvidenceInternal(1000, 10, b)
}

func BenchmarkAssessEvidence3000(b *testing.B) {
	benchmarkAssessEvidenceInternal(3000, 1, b)
}

func BenchmarkAssessEvidence10000(b *testing.B) {
	benchmarkAssessEvidenceInternal(10000, 1, b)
}
