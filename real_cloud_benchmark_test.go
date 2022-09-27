package clouditor_test

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/api/orchestrator"
	service_assessment "clouditor.io/clouditor/service/assessment"
	service_disovery "clouditor.io/clouditor/service/discovery"
	service_evidence "clouditor.io/clouditor/service/evidence"
	service_orchestrator "clouditor.io/clouditor/service/orchestrator"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var BenchmarkTypeDiscovery int = 0
var BenchmarkTypeAssessment int = 1

type benchmark struct {
	Type           int
	Key            string
	Start          time.Time
	Finish         time.Time
	ProcessedItems int
}

func assessAzure() {
	b := map[string]*benchmark{}

	srv := grpc.NewServer()

	logrus.SetLevel(logrus.DebugLevel)

	log := logrus.WithField("test", true)

	lis, _ := net.Listen("tcp", ":0")
	port := lis.Addr().(*net.TCPAddr).AddrPort().Port()
	target := fmt.Sprintf("localhost:%d", port)

	discoveryService := service_disovery.NewService(
		service_disovery.WithProviders([]string{"azure"}),
		service_disovery.WithAssessmentAddress(target),
	)
	discovery.RegisterDiscoveryServer(srv, discoveryService)

	orchestratorService := service_orchestrator.NewService()
	orchestrator.RegisterOrchestratorServer(srv, orchestratorService)

	assessmentService := service_assessment.NewService(
		service_assessment.WithEvidenceStoreAddress(target),
		service_assessment.WithOrchestratorAddress(target),
	)
	assessment.RegisterAssessmentServer(srv, assessmentService)

	evidenceService := service_evidence.NewService()
	evidence.RegisterEvidenceStoreServer(srv, evidenceService)

	go srv.Serve(lis)
	wg := sync.WaitGroup{}

	log.Info("Waiting for 3 discoverers to finish")

	wg.Add(3)

	totalResources := 0
	assessmentResults := 0

	evidenceMap := map[string]bool{}
	var mutex sync.Mutex
	go func() {
		for {
			e := <-discoveryService.Events
			if e.Type == service_disovery.DiscoveryEventTypeDiscovererFinished {
				log.Infof("Discoverer %s finished", e.Extra)
				wg.Done()

				// Add the amount of discovered resources to the wait group
				wg.Add(e.ExtraInt)
				totalResources += e.ExtraInt
				log.Infof("Waiting for %d resources of the discoverer. %d in total", e.ExtraInt, totalResources)

				// Look for the benchmark
				if benchy, ok := b[e.Extra]; ok {
					benchy.ProcessedItems = e.ExtraInt
					benchy.Finish = e.Time
				}
			} else if e.Type == service_disovery.DiscoveryEventTypeDiscovererStart {
				// Create new benchmark
				benchy := benchmark{
					Type:  BenchmarkTypeDiscovery,
					Key:   e.Extra,
					Start: e.Time,
				}

				b[benchy.Key] = &benchy
			}
		}
	}()

	assessmentService.RegisterAssessmentResultHook(func(result *assessment.AssessmentResult, err error) {
		mutex.Lock()
		defer mutex.Unlock()

		if _, ok := evidenceMap[result.EvidenceId]; !ok {
			evidenceMap[result.EvidenceId] = true
			wg.Done()

			leftOvers := assessmentService.LeftOvers()

			log.Infof("Got assessment for evidence %s, %d in total, expecting %d, waiting: %d", result.EvidenceId, len(evidenceMap), totalResources, len(leftOvers))
		}

		assessmentResults++

		// Rather "simple" way of calculating the length of the assessment
		if benchy, ok := b["Assessment"]; !ok {
			b["Assessment"] = &benchmark{
				Type:           BenchmarkTypeAssessment,
				Key:            "Assessment",
				Start:          result.Timestamp.AsTime().Local(),
				ProcessedItems: 1,
			}
		} else {
			benchy.ProcessedItems = assessmentResults
			benchy.Finish = result.Timestamp.AsTime().Local()
		}
	})

	// Start collecting from our provider
	discoveryService.Start(context.Background(), &discovery.StartDiscoveryRequest{})

	wg.Wait()

	log.Infof("Received %d assessment results. Expected %d evidences", assessmentResults, totalResources)

	fmt.Println("===== STATISTICS ====")

	fmt.Printf("Step\t\t\t\t\tStart\tFinish\tDuration\t#\tsec/#\n")
	for _, benchy := range b {
		if benchy.Type != BenchmarkTypeDiscovery {
			continue
		}

		fmt.Printf("Discovery %s\t\t\t%v\t%v\t%v\t%v\t%v\n",
			benchy.Key,
			benchy.Start.Format("15:04:05.00"),
			benchy.Finish.Format("15:04:05.00"),
			benchy.Finish.Sub(benchy.Start),
			benchy.ProcessedItems,
			benchy.Finish.Sub(benchy.Start)/time.Duration(benchy.ProcessedItems),
		)
	}

	for _, benchy := range b {
		if benchy.Type != BenchmarkTypeAssessment {
			continue
		}

		fmt.Printf("%s\t\t\t\t%v\t%v\t%v\t%v\t%v\n",
			benchy.Key, benchy.Start.Format("15:04:05.00"),
			benchy.Finish.Format("15:04:05.00"),
			benchy.Finish.Sub(benchy.Start),
			benchy.ProcessedItems,
			benchy.Finish.Sub(benchy.Start)/time.Duration(benchy.ProcessedItems),
		)
	}
}

func BenchmarkAzure(b *testing.B) {
	//for n := 0; n < b.N; n++ {
	assessAzure()
	//}
}
