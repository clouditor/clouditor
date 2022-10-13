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
var BenchmarkTypeOE int = 1
var BenchmarkTypeAssessment int = 2
var BenchmarkTypeAssessmentDetail int = 3

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
		service_disovery.WithProviders([]string{"azure", "aws", "k8s"}),
		service_disovery.WithAssessmentAddress(target),
	)
	discovery.RegisterDiscoveryServer(srv, discoveryService)

	orchestratorService := service_orchestrator.NewService(
		service_orchestrator.WithCatalogsFile("paper.json"),
	)
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
	wgDiscovery.Add(3 + 2 + 3)

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
				log.Infof("Got all assessments for evidence %s, %d in total, expecting %d, waiting: %d", e.EvidenceID, len(evidenceMap), totalResources, len(leftOvers))

				// Update our single assessment benchmark with the current one
				if benchy, ok := b["Assessment"]; ok {
					benchy.ProcessedItems = len(evidenceMap)
					benchy.Finish = e.Time
				}

				wgResources.Done()
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

	// Once all results are in the database, we can calculate the OE for each control. In the future, this would be triggered
	// directly by the orchestrator
	benchy := &benchmark{
		Key:   "Aggregation and OE",
		Type:  BenchmarkTypeOE,
		Start: time.Now(),
	}

	controls, err := orchestratorService.ListControls(context.TODO(), &orchestrator.ListControlsRequest{})
	if err != nil {
		panic(err)
	}

	for _, control := range controls.Controls {
		f, _ := orchestratorService.CalculateOEForRequirement(control)
		benchy.ProcessedItems++
		log.Infof("control %v has OE %f", control.Id, f)
	}

	benchy.Finish = time.Now()
	b[benchy.Key] = benchy

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
			benchy.TimePerItem().Microseconds(),
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
			benchy.TimePerItem().Microseconds()/1000,
		)
	}

	for _, benchy := range b {
		if benchy.Type != BenchmarkTypeOE {
			continue
		}

		fmt.Printf("%s\t\t\t& %v\t& %v\t& %v\t& %v\t& %v\\\\\n",
			benchy.Key,
			relative(startTime, benchy.Start).Milliseconds(),
			relative(startTime, benchy.Finish).Milliseconds(),
			benchy.Finish.Sub(benchy.Start).Milliseconds(),
			benchy.ProcessedItems,
			benchy.TimePerItem().Microseconds()/1000,
		)
	}

	fmt.Println("\n===== STATISTICS DETAIL FOR ASSESSMENT ====")
	fmt.Printf("min\tmax\tavg\tmedian\n")

	// TODO: min, max, median, avg for a better overview
	assessmentValues := []*benchmark{}
	for _, value := range b {
		if value.Type == BenchmarkTypeAssessmentDetail {
			// Somehow. some resources do not finish. not sure why
			if value.Finish.IsZero() {
				fmt.Printf("Skipped %s because not finished\n", value.Key)
				continue
			}
			assessmentValues = append(assessmentValues, value)
		}
	}

	fmt.Println("===== STATISTICS per Resource Type ====")
	fmt.Println("Resource Type\t& \\#policies\t& min\t& max\t& avg\t& median\\\\\\hline\\hline")

	for typ, groupedValues := range benchyPerResourceType {
		min, max, avg, median := statistics(groupedValues)
		var totalTime time.Duration
		for _, value := range groupedValues {
			totalTime += value.Finish.Sub(value.Start)
		}

		fmt.Printf("%s\t\t& %v\t\t& %.03f\t& %.03f\t& %.03f\t& %.03f\\\\\n", typ, groupedValues[0].ProcessedItems,
			float64(min.Microseconds())/1000.0,
			float64(max.Microseconds())/1000.0,
			float64(avg.Microseconds())/1000.0,
			float64(median.Microseconds())/1000.0,
		)
	}

	fmt.Println("\\hline")

	// retrieve some detailed statistics of it
	min, max, avg, median := statistics(assessmentValues)
	fmt.Printf("%s\t\t& %v\t\t& %.03f\t& %.03f\t& %.03f\t& %.03f\\\\\n", "Overall", "-",
		float64(min.Microseconds())/1000.0,
		float64(max.Microseconds())/1000.0,
		float64(avg.Microseconds())/1000.0,
		float64(median.Microseconds())/1000.0,
	)
}

func statistics(values []*benchmark) (min time.Duration, max time.Duration, avg time.Duration, median time.Duration) {
	// We need to sort the values according to their runtime, so we can easily get min/max/median
	slices.SortFunc(values, func(a, b *benchmark) bool {
		return a.RunTime() < b.RunTime()
	})

	// for _, value := range values {
	// 	fmt.Printf("%v\n", value.RunTime())
	// }

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
