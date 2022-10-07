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
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
)

var BenchmarkTypeDiscovery int = 0
var BenchmarkTypeAssessment int = 1
var BenchmarkTypeAssessmentDetail int = 2

type benchmark struct {
	Type           int
	Key            string
	Start          time.Time
	Finish         time.Time
	ProcessedItems int
}

func (b *benchmark) RunTime() time.Duration {
	return b.Finish.Sub(b.Start)
}

func (b *benchmark) TimePerItem() time.Duration {
	if b.ProcessedItems == 0 {
		return 0
	}

	return b.Finish.Sub(b.Start) / time.Duration(b.ProcessedItems)
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
	wgDiscovery := sync.WaitGroup{}
	wgResources := sync.WaitGroup{}

	log.Info("Waiting for 5 discoverers to finish")

	//wgDiscovery.Add(2)
	wgDiscovery.Add(3)

	totalResources := 0
	assessmentResults := 0

	evidenceMap := map[string]bool{}
	//var mutex sync.Mutex
	go func() {
		for {
			e := <-discoveryService.Events
			if e.Type == service_disovery.DiscoveryEventTypeDiscovererFinished {
				log.Infof("Discoverer %s finished", e.Extra)
				wgDiscovery.Done()

				// Add the amount of discovered resources to the wait group
				wgResources.Add(e.ExtraInt)
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

	benchyPerResourceType := map[string][]*benchmark{}
	go func() {
		for {
			e := <-assessmentService.Events
			if e.Type == service_assessment.AssessmentEventTypeEvidenceFinished {
				// Look for the benchmark
				if benchy, ok := b[e.EvidenceID]; ok {
					benchy.ProcessedItems = e.NumAssessments
					benchy.Finish = e.Time

					benchyPerResourceType[e.ResourceType] = append(benchyPerResourceType[e.ResourceType], benchy)
				}

				leftOvers := assessmentService.LeftOvers()
				assessmentResults += e.NumAssessments
				evidenceMap[e.EvidenceID] = true
				wgResources.Done()
				log.Infof("Got all assessments for evidence %s, %d in total, expecting %d, waiting: %d", e.EvidenceID, len(evidenceMap), totalResources, len(leftOvers))

				// Update our single assessment benchmark with the current one
				if benchy, ok := b["Assessment"]; ok {
					benchy.ProcessedItems = len(evidenceMap)
					benchy.Finish = e.Time
				}

				if len(leftOvers) == 6 {
					log.Infof("why?")
				}
			} else if e.Type == service_assessment.AssessmentEventTypeEvidenceStarted {
				// Create new benchmark for detail
				benchy := &benchmark{
					Type:  BenchmarkTypeAssessmentDetail,
					Key:   e.EvidenceID,
					Start: e.Time,
				}

				b[benchy.Key] = benchy

				// Rather "simple" way of calculating the length of the assessment
				if _, ok := b["Assessment"]; !ok {
					b["Assessment"] = &benchmark{
						Type:           BenchmarkTypeAssessment,
						Key:            "Assessment",
						Start:          e.Time,
						ProcessedItems: 1,
					}
				}
			}
		}
	}()

	/*assessmentService.RegisterAssessmentResultHook(func(result *assessment.AssessmentResult, err error) {
		mutex.Lock()
		defer mutex.Unlock()

		if _, ok := evidenceMap[result.EvidenceId]; !ok {
			evidenceMap[result.EvidenceId] = true
			wgResources.Done()

			leftOvers := assessmentService.LeftOvers()

			if len(evidenceMap) == 143 {
				log.Infof("where are my 4")
			}

			log.Infof("Got assessment for evidence %s, %d in total, expecting %d, waiting: %d", result.EvidenceId, len(evidenceMap), totalResources, len(leftOvers))
		}

		ssessmentResults++

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
	})*/

	// Start collecting from our provider
	discoveryService.Start(context.Background(), &discovery.StartDiscoveryRequest{})

	wgDiscovery.Wait()
	wgResources.Wait()

	log.Infof("Received %d assessment results. Expected %d evidences", assessmentResults, totalResources)

	fmt.Println("===== STATISTICS ====")

	values := maps.Values(b)
	slices.SortFunc(values, func(a *benchmark, b *benchmark) bool {
		return a.Start.Before(b.Start)
	})
	startTime := values[0].Start

	fmt.Printf("Step\t\t\t\t\tStart [t+ms]\tFinish [t+ms]\tDuration [ms]\t#\t1/# [ms]\n")
	for _, benchy := range b {
		if benchy.Type != BenchmarkTypeDiscovery {
			continue
		}

		fmt.Printf("Discovery %s\t\t\t& %v\t& %v\t& %v\t& %v\t& %v\\\\\n",
			benchy.Key,
			relative(startTime, benchy.Start).Milliseconds(),
			relative(startTime, benchy.Finish).Milliseconds(),
			benchy.Finish.Sub(benchy.Start).Milliseconds(),
			benchy.ProcessedItems,
			benchy.TimePerItem(),
		)
	}

	for _, benchy := range b {
		if benchy.Type != BenchmarkTypeAssessment {
			continue
		}

		fmt.Printf("%s\t\t\t& %v\t& %v\t& %v\t& %v\t& %v\\\\\n",
			benchy.Key,
			relative(startTime, benchy.Start).Milliseconds(),
			relative(startTime, benchy.Finish).Milliseconds(),
			benchy.Finish.Sub(benchy.Start).Milliseconds(),
			benchy.ProcessedItems,
			benchy.TimePerItem(),
		)
	}

	fmt.Println("\n===== STATISTICS DETAIL FOR ASSESSMENT ====")
	fmt.Printf("min\tmax\tavg\tmedian\n")

	// TODO: min, max, median, avg for a better overview
	assessmentValues := []*benchmark{}
	for _, value := range b {
		if value.Type == BenchmarkTypeAssessmentDetail {
			assessmentValues = append(assessmentValues, value)
		}
	}

	// retrieve some detailed statistics of it
	min, max, avg, median := statistics(assessmentValues)
	fmt.Printf("%v\t%v\t%v\t%v\t\n\n", min, max, avg, median)

	fmt.Println("===== STATISTICS per Resource Type ====")
	fmt.Println("Resource Type\t\tTotal Time [ms]\t#\t1/# [ms]\t#policies")

	for typ, groupedValues := range benchyPerResourceType {
		var totalTime time.Duration
		for _, value := range groupedValues {
			totalTime += value.Finish.Sub(value.Start)
		}

		fmt.Printf("%s\t\t%v\t\t%v\t%v\t%v\n", typ, totalTime.Milliseconds(), len(groupedValues), (totalTime / time.Duration(len(groupedValues))), groupedValues[0].ProcessedItems)
	}
}

func statistics(values []*benchmark) (min time.Duration, max time.Duration, avg time.Duration, median time.Duration) {
	// We need to sort the values according to their runtime, so we can easily get min/max/median
	slices.SortFunc(values, func(a, b *benchmark) bool {
		return a.RunTime() < b.RunTime()
	})

	// min is the first entry, max the last
	min = values[0].RunTime()
	max = values[len(values)-1].RunTime()

	// median is the one in the middle
	if len(values)%2 == 0 {
		value1 := values[len(values)/2-1].RunTime()
		value2 := values[len(values)/2].RunTime()
		median = (value1 + value2) / 2
	} else {
		median = values[(len(values)-1)/2].RunTime()
	}

	var totalTime time.Duration
	for _, value := range values {
		totalTime += value.RunTime()
	}

	avg = totalTime / time.Duration(len(values))

	return
}

func relative(startTime time.Time, time time.Time) time.Duration {
	return time.Sub(startTime)
}

func BenchmarkAzure(b *testing.B) {
	//for n := 0; n < b.N; n++ {
	assessAzure()
	//}
}
