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

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testutil"
	service_evidence "clouditor.io/clouditor/service/evidence"
	service_orchestrator "clouditor.io/clouditor/service/orchestrator"
	"clouditor.io/clouditor/voc"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func createEvidences(n int, m int, b *testing.B) int {
	var (
		wg   sync.WaitGroup
		err  error
		sock net.Listener
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

	wg.Add(n / 2 * m * (3 + 8))

	var count int64 = 0

	addr := fmt.Sprintf("localhost:%d", sock.Addr().(*net.TCPAddr).Port)

	svc := NewService(WithOrchestratorAddress(addr), WithEvidenceStoreAddress(addr))

	orchestratorService.RegisterAssessmentResultHook(func(result *assessment.AssessmentResult, err error) {
		wg.Done()
		current := atomic.AddInt64(&count, 1)
		log.Debugf("Current count: %v", current)
	})

	// Create m parallel executions of our evidence creation
	for j := 0; j < m; j++ {
		go func() {
			// Create evidences for n/2 resources (2 per resource)
			for i := 0; i < n/2; i++ {
				if i%100 == 0 {
					log.Infof("Currently @ %v", i)
					log.Infof("Current stats: %+v", svc.stats)
				}

				vm := voc.VirtualMachine{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:   voc.ResourceID(fmt.Sprintf("%d-vm", i)),
							Type: []string{"VirtualMachine", "Compute", "Resource"},
						},
					},
					BlockStorage: []voc.ResourceID{voc.ResourceID(fmt.Sprintf("%d-storage", i))},
				}

				_, err := svc.AssessEvidence(context.Background(), &assessment.AssessEvidenceRequest{
					Evidence: &evidence.Evidence{
						Id:                 uuid.NewString(),
						Timestamp:          timestamppb.Now(),
						ToolId:             "mytool",
						Resource:           testutil.ToStruct(vm, b),
						RelatedResourceIds: vm.Related(),
					},
				})
				if err != nil {
					b.Errorf("Error while calling AssessEvidence: %v", err)
				}

				s := voc.BlockStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:   voc.ResourceID(fmt.Sprintf("%d-storage", i)),
							Type: []string{"BlockStorage", "Storage"},
						},
						AtRestEncryption: &voc.AtRestEncryption{
							Enabled:   true,
							Algorithm: "AES-256",
						},
					},
				}

				_, err = svc.AssessEvidence(context.Background(), &assessment.AssessEvidenceRequest{
					Evidence: &evidence.Evidence{
						Id:                 uuid.NewString(),
						Timestamp:          timestamppb.Now(),
						ToolId:             "mytool",
						Resource:           testutil.ToStruct(s, b),
						RelatedResourceIds: s.Related(),
					},
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

func benchmarkAssessEvidence(i int, m int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		createEvidences(i, m, b)
	}
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
	benchmarkAssessEvidence(3000, 1, b)
}

func BenchmarkAssessEvidence10000(b *testing.B) {
	benchmarkAssessEvidence(10000, 1, b)
}
