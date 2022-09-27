package clouditor_test

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"

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

func assessAzure() {
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

			log.Infof("Got assessment for evidenec %s, %d in total, expecting %d, waiting: %d", result.EvidenceId, len(evidenceMap), totalResources, len(leftOvers))

			if len(leftOvers) == 7 {
				log.Infof("the remaining 7")
				wg.Add(-7)
			}
		}

		assessmentResults++
	})

	// Start collecting from our provider
	discoveryService.Start(context.Background(), &discovery.StartDiscoveryRequest{})

	wg.Wait()

	log.Infof("Received %d assessment results. Expected %d evidences", assessmentResults, totalResources)
}

func BenchmarkAzure(b *testing.B) {
	//for n := 0; n < b.N; n++ {
	assessAzure()
	//}
}
