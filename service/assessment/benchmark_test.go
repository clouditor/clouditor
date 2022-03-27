package assessment

import (
	"context"
	"fmt"
	"math"
	"net"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testutil"
	service_evidence "clouditor.io/clouditor/service/evidence"
	service_orchestrator "clouditor.io/clouditor/service/orchestrator"
	"clouditor.io/clouditor/voc"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var NumFunctionsMetrics = 2
var NumVMMetrics = 13
var NumLoggingServiceMetrics = 7
var NumBlockStorageMetrics = 4
var NumObjectStorageMetrics = 4
var NumIdentityMetrics = 5

func createVMWithMalwareProtection(numCloudServices int, numAccounts int, numVMs, numFunction int, b *testing.B) {
	var (
		wg   sync.WaitGroup
		err  error
		sock net.Listener
	)

	srv := grpc.NewServer()

	logrus.SetLevel(logrus.PanicLevel)

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

	addr := fmt.Sprintf("localhost:%d", sock.Addr().(*net.TCPAddr).Port)

	var count int64 = 0

	work := (NumFunctionsMetrics*NumFunctionsMetrics +
		numVMs*NumVMMetrics +
		numVMs*NumBlockStorageMetrics +
		1*NumObjectStorageMetrics +
		1*NumLoggingServiceMetrics +
		NumIdentityMetrics*numAccounts) * numCloudServices

	log.Infof("Waiting for %d results", work)

	wg.Add(work)

	svc := NewService(WithOrchestratorAddress(addr), WithEvidenceStoreAddress(addr))
	orchestratorService.RegisterAssessmentResultHook(func(result *assessment.AssessmentResult, err error) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()

		current := atomic.AddInt64(&count, 1)

		log.Debugf("Current count: %v - stats: %+v", current, svc.stats)

		// In this scenario, we expect that all results are 100 % true
		if !result.Compliant {
			b.Fatalf("Assesment result (metric %s) for resource %s is not compliant", result.MetricId, result.ResourceId)
		}

		if err == nil {
			wg.Done()
		} else {
			log.Errorf("Error in assessment hook: %s", err)
		}

	})

	// Create m parallel executions for each "cloud service"
	for j := 0; j < numCloudServices; j++ {
		go func() {
			for i := 0; i < numFunction; i++ {
				f := &voc.Function{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:   voc.ResourceID(fmt.Sprintf("%d-function", i)),
							Type: []string{"Function", "Compute", "Resource"},
						},
					},
					RuntimeVersion:  "11",
					RuntimeLanguage: "Java",
				}

				assess(svc, f, b)
			}

			for i := 0; i < numVMs; i++ {
				vm := &voc.VirtualMachine{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:   voc.ResourceID(fmt.Sprintf("%d-vm", i)),
							Type: []string{"VirtualMachine", "Compute", "Resource"},
						},
					},
					BootLogging: &voc.BootLogging{
						Logging: &voc.Logging{
							LoggingService:  []voc.ResourceID{"logging-service"},
							Enabled:         true,
							RetentionPeriod: time.Hour * 24 * 30,
						},
					},
					OSLogging: &voc.OSLogging{
						Logging: &voc.Logging{
							LoggingService:  []voc.ResourceID{"logging-service"},
							Enabled:         true,
							RetentionPeriod: time.Hour * 24 * 30,
						},
					},
					BlockStorage: []voc.ResourceID{voc.ResourceID(fmt.Sprintf("%d-storage", i))},
					AutomaticUpdates: &voc.AutomaticUpdates{
						Enabled:      true,
						SecurityOnly: true,
						Interval:     time.Hour * 24 * 2,
					},
					MalwareProtection: &voc.MalwareProtection{
						Enabled: true,
						ApplicationLogging: &voc.ApplicationLogging{
							Logging: &voc.Logging{
								LoggingService:  []voc.ResourceID{"logging-service"},
								Enabled:         true,
								RetentionPeriod: time.Hour * 24 * 30,
							},
						},
					},
				}
				s := voc.BlockStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:   voc.ResourceID(fmt.Sprintf("%d-storage", i)),
							Type: []string{"BlockStorage", "Storage"},
						},
						AtRestEncryption: &voc.CustomerKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Enabled:   true,
								Algorithm: "AES256",
							},
						},
					},
				}

				assess(svc, vm, b)
				assess(svc, s, b)
			}

			ls := voc.LoggingService{
				NetworkService: &voc.NetworkService{
					Networking: &voc.Networking{
						Resource: &voc.Resource{
							ID:   voc.ResourceID("logging-service"),
							Type: []string{"LoggingService", "NetworkService", "Networking", "Resource"},
						},
					},
					Authenticity: &voc.TokenBasedAuthentication{
						Activated: true,
					},
					TransportEncryption: &voc.TransportEncryption{
						Enforced:   true,
						Enabled:    true,
						Algorithm:  "TLS",
						TlsVersion: "TLS1.2",
					},
				},
				Storage: []voc.ResourceID{voc.ResourceID("log-storage")},
			}
			lss := voc.ObjectStorage{
				Storage: &voc.Storage{
					Resource: &voc.Resource{
						ID:   voc.ResourceID("log-storage"),
						Type: []string{"ObjectStorage", "Storage"},
					},
					AtRestEncryption: &voc.CustomerKeyEncryption{
						AtRestEncryption: &voc.AtRestEncryption{
							Enabled:   true,
							Algorithm: "AES256",
						},
					},
					Immutability: &voc.Immutability{
						Enabled: true,
					},
				},
			}

			assess(svc, ls, b)
			assess(svc, lss, b)

			for j := 0; j < numAccounts; j++ {
				user := &voc.Identity{
					Identifiable: &voc.Identifiable{
						Resource: &voc.Resource{
							ID:   voc.ResourceID(fmt.Sprintf("%d-user", j)),
							Type: []string{"Identity", "Resource"},
						},
					},
					Authenticity: []voc.HasAuthenticity{
						&voc.OTPBasedAuthentication{Activated: true},
						&voc.PasswordBasedAuthentication{Activated: true},
					},
					Privileged:            true,
					DisablePasswordPolicy: false,
					Activated:             true,
					LastActivity:          time.Now().Add(-time.Hour * 24 * 7),
				}

				assess(svc, user, b)
			}
		}()
	}

	wg.Wait()
}

func assess(svc assessment.AssessmentServer, r voc.IsCloudResource, t assert.TestingT) {
	_, err := svc.AssessEvidence(context.Background(), &assessment.AssessEvidenceRequest{
		Evidence: &evidence.Evidence{
			Id:                 uuid.NewString(),
			Timestamp:          timestamppb.Now(),
			ToolId:             "mytool",
			Resource:           testutil.ToStruct(r, t),
			RelatedResourceIds: r.Related(),
		},
	})
	if err != nil {
		t.Errorf("Error while calling AssessEvidence: %v", err)
	}
}

var numVMs = []int{10 /*, 100, 1000, 5000*/}

func BenchmarkComplex(b *testing.B) {
	for _, k := range numVMs {
		for l := 0.; l <= 2; l++ {
			parallel := int(math.Pow(2, l))
			b.Run(fmt.Sprintf("%d/%d", k, parallel), func(b *testing.B) {
				//for n := 0; n < b.N; n++ {
				createVMWithMalwareProtection(parallel, k, k/10, k/10, b)
				//}
			})
		}
	}
}

func createVMEvidences(n int, m int, b *testing.B) {
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

	wg.Add(n * m * 9)

	var count int64 = 0

	addr := fmt.Sprintf("localhost:%d", sock.Addr().(*net.TCPAddr).Port)

	svc := NewService(WithOrchestratorAddress(addr), WithEvidenceStoreAddress(addr))

	orchestratorService.RegisterAssessmentResultHook(func(result *assessment.AssessmentResult, err error) {
		current := atomic.AddInt64(&count, 1)

		log.Debugf("Current count: %v - stats: %+v", current, svc.stats)

		wg.Done()
	})

	// Create m parallel executions of our evidence creation
	for j := 0; j < m; j++ {
		go func() {
			// Create evidences for n (1 resource per )
			for i := 0; i < n; i++ {
				if i%100 == 0 {
					log.Infof("Currently @ %v - stats: %+v", i, svc.stats)
				}

				vm := voc.VirtualMachine{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:   voc.ResourceID(fmt.Sprintf("%d-%d-vm", j, i)),
							Type: []string{"VirtualMachine", "Compute", "Resource"},
						},
					},
					BlockStorage: nil,
				}

				assess(svc, vm, b)
			}
		}()
	}

	wg.Wait()
}

func createIdentityEvidences(n int, m int, b *testing.B) {
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

	wg.Add(n * m * NumIdentityMetrics)

	var count int64 = 0

	addr := fmt.Sprintf("localhost:%d", sock.Addr().(*net.TCPAddr).Port)

	svc := NewService(WithOrchestratorAddress(addr), WithEvidenceStoreAddress(addr))

	orchestratorService.RegisterAssessmentResultHook(func(result *assessment.AssessmentResult, err error) {
		current := atomic.AddInt64(&count, 1)

		log.Debugf("Current count: %v - stats: %+v", current, svc.stats)

		wg.Done()
	})

	// Create m parallel executions of our evidence creation
	for j := 0; j < m; j++ {
		go func() {
			// Create evidences for n (1 resource per )
			for i := 0; i < n; i++ {
				if i%100 == 0 {
					log.Infof("Currently @ %v - stats: %+v", i, svc.stats)
				}

				i := voc.Identity{
					Identifiable: &voc.Identifiable{
						Resource: &voc.Resource{
							ID:   voc.ResourceID(fmt.Sprintf("%d-%d-identity", j, i)),
							Type: []string{"Identity", "Identifiable", "Resource"},
						},
					},
					Privileged: true,
					Activated:  true,
				}

				assess(svc, i, b)
			}
		}()
	}

	wg.Wait()
}

/*func createRoleEvidences(n int, m int, b *testing.B) {
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

	wg.Add(n * m * 2)

	var count int64 = 0

	addr := fmt.Sprintf("localhost:%d", sock.Addr().(*net.TCPAddr).Port)

	svc := NewService(WithOrchestratorAddress(addr), WithEvidenceStoreAddress(addr))

	orchestratorService.RegisterAssessmentResultHook(func(result *assessment.AssessmentResult, err error) {
		current := atomic.AddInt64(&count, 1)

		log.Debugf("Current count: %v - stats: %+v", current, svc.stats)

		wg.Done()
	})

	// Create m parallel executions of our evidence creation
	for j := 0; j < m; j++ {
		go func() {
			// Create evidences for n (1 resource per )
			for i := 0; i < n; i++ {
				if i%100 == 0 {
					log.Infof("Currently @ %v - stats: %+v", i, svc.stats)
				}

				i := voc.RoleAssignment{
					Identifiable: &voc.Identifiable{
						Resource: &voc.Resource{
							ID:   voc.ResourceID(fmt.Sprintf("%d-%d-identity", j, i)),
							Type: []string{"Identity", "Identifiable", "Resource"},
						},
					},
					MixedDuties: 0.5,
				}

				assess(svc, i, b)
			}
		}()
	}

	wg.Wait()
}*/

func createServiceEvidences(n int, m int, b *testing.B) {
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

	wg.Add(n * m * 4)

	var count int64 = 0

	addr := fmt.Sprintf("localhost:%d", sock.Addr().(*net.TCPAddr).Port)

	svc := NewService(WithOrchestratorAddress(addr), WithEvidenceStoreAddress(addr))

	orchestratorService.RegisterAssessmentResultHook(func(result *assessment.AssessmentResult, err error) {
		current := atomic.AddInt64(&count, 1)

		log.Debugf("Current count: %v - stats: %+v", current, svc.stats)

		wg.Done()
	})

	// Create m parallel executions of our evidence creation
	for j := 0; j < m; j++ {
		go func() {
			// Create evidences for n (1 resource per )
			for i := 0; i < n; i++ {
				if i%100 == 0 {
					log.Infof("Currently @ %v - stats: %+v", i, svc.stats)
				}

				i := voc.NetworkService{
					Networking: &voc.Networking{
						Resource: &voc.Resource{
							ID:   voc.ResourceID(fmt.Sprintf("%d-%d-network-service", j, i)),
							Type: []string{"NetworkService", "Networking", "Resource"},
						},
					},
					Authenticity: &voc.TokenBasedAuthentication{},
					TransportEncryption: &voc.TransportEncryption{
						Enabled: true,
					},
				}

				assess(svc, i, b)
			}
		}()
	}

	wg.Wait()
}

func createStorageEvidences(n int, m int, b *testing.B) {
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

	wg.Add(n * m * 4)

	var count int64 = 0

	addr := fmt.Sprintf("localhost:%d", sock.Addr().(*net.TCPAddr).Port)

	svc := NewService(WithOrchestratorAddress(addr), WithEvidenceStoreAddress(addr))

	orchestratorService.RegisterAssessmentResultHook(func(result *assessment.AssessmentResult, err error) {
		current := atomic.AddInt64(&count, 1)

		log.Debugf("Current count: %v - stats: %+v", current, svc.stats)

		wg.Done()
	})

	// Create m parallel executions of our evidence creation
	for j := 0; j < m; j++ {
		go func() {
			// Create evidences for n (1 resource per )
			for i := 0; i < n; i++ {
				if i%100 == 0 {
					log.Infof("Currently @ %v - stats: %+v", i, svc.stats)
				}

				vm := voc.BlockStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:   voc.ResourceID(fmt.Sprintf("%d-%d-storage", j, i)),
							Type: []string{"BlockStorage", "Storage", "Resource"},
						},
					},
				}

				assess(svc, vm, b)
			}
		}()
	}

	wg.Wait()
}

var numEvidences = []int{10, 100, 1000, 5000, 10000, 20000, 30000, 40000, 50000}

var create = []createFuncType{createVMEvidences, createStorageEvidences, createIdentityEvidences, createServiceEvidences}

type createFuncType func(int, int, *testing.B)

func BenchmarkAssessVMEvidence(b *testing.B) {
	for _, k := range numEvidences {
		for l := 1; l <= 3; l++ {
			b.Run(fmt.Sprintf("%d/%d", k, l), func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					createVMEvidences(k, l, b)
				}
			})
		}
	}
}

func BenchmarkAssessStorageEvidence(b *testing.B) {
	for _, k := range numEvidences {
		for l := 1; l <= 3; l++ {
			b.Run(fmt.Sprintf("%d/%d", k, l), func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					createStorageEvidences(k, l, b)
				}
			})
		}
	}
}

func BenchmarkAssessIdentityEvidence(b *testing.B) {
	for _, k := range numEvidences {
		for l := 1; l <= 3; l++ {
			b.Run(fmt.Sprintf("%d/%d", k, l), func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					createIdentityEvidences(k, l, b)
				}
			})
		}
	}
}

func BenchmarkEvidenceTypes(b *testing.B) {
	numEvidences := 10000

	for _, k := range create {
		_, name, _ := strings.Cut(runtime.FuncForPC(reflect.ValueOf(k).Pointer()).Name(), "assessment.")

		b.Run(fmt.Sprintf("%s/%d", name, numEvidences), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				k(numEvidences, 1, b)
			}
		})
	}
}

func BenchmarkAssessStorageEvidence2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		createStorageEvidences(2, 1, b)
	}
}

func BenchmarkAssessVMEvidence2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		createVMEvidences(2, 1, b)
	}
}

func BenchmarkAssessVMEvidence4(b *testing.B) {
	for n := 0; n < b.N; n++ {
		createVMEvidences(4, 1, b)
	}
}

func BenchmarkAssessVMEvidence10(b *testing.B) {
	for n := 0; n < b.N; n++ {
		createVMEvidences(10, 1, b)
	}
}

func BenchmarkAssessVMEvidence10x2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		createVMEvidences(10, 2, b)
	}
}

func BenchmarkAssessVMEvidence100(b *testing.B) {
	for n := 0; n < b.N; n++ {
		createVMEvidences(100, 1, b)
	}
}

func BenchmarkAssessVMEvidence1000(b *testing.B) {
	for n := 0; n < b.N; n++ {
		createVMEvidences(1000, 1, b)
	}
}

func BenchmarkAssessVMEvidence1000x2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		createVMEvidences(1000, 2, b)
	}
}

func BenchmarkAssessVMEvidence1000x10(b *testing.B) {
	for n := 0; n < b.N; n++ {
		createVMEvidences(1000, 10, b)
	}
}

func BenchmarkAssessVMEvidence3000(b *testing.B) {
	for n := 0; n < b.N; n++ {
		createVMEvidences(3000, 1, b)
	}
}

func BenchmarkAssessVMEvidence10000(b *testing.B) {
	for n := 0; n < b.N; n++ {
		createVMEvidences(10000, 1, b)
	}
}

func BenchmarkAssessVMEvidence10000x2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		createVMEvidences(10000, 2, b)
	}
}

func BenchmarkAssessVMEvidence30000x4(b *testing.B) {
	for n := 0; n < b.N; n++ {
		createVMEvidences(30000, 4, b)
	}
}
