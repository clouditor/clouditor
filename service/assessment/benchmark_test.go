package assessment

import (
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/logging/formatter"
	service_evidence "clouditor.io/clouditor/service/evidence"
	service_orchestrator "clouditor.io/clouditor/service/orchestrator"
	"clouditor.io/clouditor/voc"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
)

func createEvidences(n int, m int, b *testing.B) int {
	var (
		wg   sync.WaitGroup
		err  error
		sock net.Listener
	)

	logrus.SetLevel(logrus.InfoLevel)

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

	addr := fmt.Sprintf("localhost:%d", sock.Addr().(*net.TCPAddr).Port)

	svc := NewService(WithOrchestratorAddress(addr), WithEvidenceStoreAddress(addr))

	svc.RegisterAssessmentResultHook(func(result *assessment.AssessmentResult, err error) {
		wg.Done()
		current := atomic.AddInt64(&count, 1)
		log.Debugf("Current count: %v", current)
	})

	// Create m parallel executions of our evidence creation
	for j := 0; j < m; j++ {
		go func() {
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

				if i%100 == 0 {
					log.Infof("Currently @ %v", i)
				}

				_, err := svc.AssessEvidence(context.Background(), &assessment.AssessEvidenceRequest{
					Evidence: &e,
				})
				if err != nil {
					b.Errorf("Error while calling AssessEvidence: %v", err)
				}
			}
		}()
	}

	wg.Wait()

	return 0
}

func createEvidencesForTest(n int, m int, t *testing.T) int {
	var (
		wg   sync.WaitGroup
		err  error
		sock net.Listener
	)

	log.Logger.Formatter = formatter.CapitalizeFormatter{Formatter: &logrus.TextFormatter{ForceColors: true}}
	logrus.SetLevel(logrus.DebugLevel)

	srv := grpc.NewServer()

	orchestratorService := service_orchestrator.NewService()
	orchestrator.RegisterOrchestratorServer(srv, orchestratorService)

	evidenceService := service_evidence.NewService()
	evidence.RegisterEvidenceStoreServer(srv, evidenceService)

	sock, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Errorf("could not listen: %v", err)
	}
	assert.Nil(t, err)

	go func() {
		err := srv.Serve(sock)
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error while creating gRPC server: %v", err)
		}
	}()

	wg.Add(n * m * 7)

	var count int64 = 0

	addr := fmt.Sprintf("localhost:%d", sock.Addr().(*net.TCPAddr).Port)

	svc := NewService(WithOrchestratorAddress(addr), WithEvidenceStoreAddress(addr))

	orchestratorService.RegisterAssessmentResultHook(func(result *assessment.AssessmentResult, err error) {
		wg.Done()
		current := atomic.AddInt64(&count, 1)
		log.Debugf("Current count: %v", current)
	})

	// Create m parallel executions of our evidence creation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer srv.Stop()

	for j := 0; j < m; j++ {
		go func() {
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

				if i%100 == 0 {
					log.Infof("Currently @ %v", i)
				}

				_, err := svc.AssessEvidence(ctx, &assessment.AssessEvidenceRequest{
					Evidence: &e,
				})

				assert.Nil(t, err)
			}
		}()
	}

	wg.Wait()

	return 0
}

func benchmarkAssessEvidence(i int, m int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		createEvidences(i, m, b)
	}
}

//func Test_benchmarkAssessEvidence(t *testing.T) {
//	for n := 0; n < 1; n++ {
//		createEvidencesForTest(100, 2, t)
//	}
//}

func BenchmarkAssessEvidence1(b *testing.B) {
	benchmarkAssessEvidence(1, 1, b)
}

func BenchmarkAssessEvidence2(b *testing.B) {
	benchmarkAssessEvidence(2, 1, b)
}

func BenchmarkAssessEvidence10(b *testing.B) {
	benchmarkAssessEvidence(10, 1, b)
}

func BenchmarkAssessEvidence10x2(b *testing.B) {
	benchmarkAssessEvidence(10, 2, b)
}

func BenchmarkAssessEvidence100(b *testing.B) {
	benchmarkAssessEvidence(100, 1, b)
}

func BenchmarkAssessEvidence1000(b *testing.B) {
	benchmarkAssessEvidence(1000, 1, b)
}

func BenchmarkAssessEvidence1000x2(b *testing.B) {
	benchmarkAssessEvidence(1000, 2, b)
}

func BenchmarkAssessEvidence1000x10(b *testing.B) {
	benchmarkAssessEvidence(1000, 10, b)
}

func BenchmarkAssessEvidence3000(b *testing.B) {
	benchmarkAssessEvidence(3000, 2, b)
}

func BenchmarkAssessEvidence10000(b *testing.B) {
	benchmarkAssessEvidence(10000, 1, b)
}
