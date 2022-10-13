package orchestrator

import (
	"time"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"golang.org/x/exp/slices"
)

type Evaluator interface {
	Evaluate(res *assessment.AssessmentResult)
}

var DefaultOETime = time.Hour * 24 * 7

/*func (svc *Service) Evaluate(res *assessment.AssessmentResult) (err error) {
var ok bool
// Build a map of requirements and their results

/*metric, ok := svc.metric(res.MetricId)
if !ok {
	return errors.New("could not evaluate: invalid metric")
}
// new average = old average * (n-1)/n + new value /n

//reqID := metric.Category

eval, ok := svc.EvalMetrics[reqID]
if !ok {
	svc.EvalMetrics[reqID] = &EvalMetric{
		requirementID: reqID,
		fullfilled:    true,
		op:            1,
		n:             1,
	}
}*/

// Build a map of requirements and their results

/*metric, err := svc.GetMetric(context.TODO(), &orchestrator.GetMetricRequest{MetricId: res.MetricId})
	if err != nil {
		return errors.New("could not evaluate: invalid metric")
	}

	reqID := metric.Category

	var list []*assessment.AssessmentResult
	if list, ok = svc.evaluations[reqID]; !ok {
		list = []*assessment.AssessmentResult{}
	}

	svc.evaluations[reqID] = append(list, res)

	log.Debugf("We have %d assessment results for requirement %s", len(list)+1, reqID)

	svc.calculateMetrics()

	return nil
}

func (svc *Service) calculateMetrics() {
	// do several things

	// first, we want the current state of each requirement
	m := EvalMetric{
		fullfilled: true,
	}

	// TODO(oxisto): optimization: only calculate if a requirement has new results, for now we calculate all
	for reqID, list := range svc.evaluations {
		for _, res := range list {
			if !res.Compliant {
				m.fullfilled = false
			}
		}

		m.op, _ = svc.calculateOEForRequirement(reqID, list)

		log.Infof("Requirement %s is now %+v", reqID, m)
	}
}*/

func (svc *Service) CalculateOEForRequirement(control *orchestrator.Control) (op float64, err error) {
	t := time.Now().Add(-DefaultOETime)

	var n = 0

	var metricIds []string
	for _, metric := range control.Metrics {
		metricIds = append(metricIds, metric.Id)
	}

	// filter results
	for _, result := range svc.results {
		if !slices.Contains(metricIds, result.MetricId) {
			continue
		}

		if result.GetTimestamp().Seconds > t.Unix() {
			if result.Compliant {
				op += 1
			}
			n += 1
		}
	}

	if n > 0 {
		op /= float64(n)
	} else {
		op = 0
	}

	return
}
