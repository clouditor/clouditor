package evaluation

import (
	"context"
	"fmt"
	"time"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/service/standalone"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Service struct {
}

var log *logrus.Entry

func init() {
	log = logrus.WithField("component", "discovery")
}

var DefaultOETime = time.Hour * 24 * 7

func (s Service) Do() (err error) {
	var client assessment.AssessmentClient

	var isStandalone bool = true

	if isStandalone {
		client = standalone.NewAssessmentClient()
	} else {
		// TODO(oxisto): support assessment on another tcp/port
		cc, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
		if err != nil {
			//return nil, status.Errorf(codes.Internal, "could not connect to assessment service: %v", err)
			return err
		}
		client = assessment.NewAssessmentClient(cc)
	}

	resp, _ := client.ListAssessmentResults(context.Background(), &assessment.ListAssessmentResultsRequest{})

	op, _ := s.calculateOEForObjective("CKM-02", resp.Results) // should contain CKM-02.2

	fmt.Printf("op: %+v", op)

	return
}

func (s Service) calculateOEForObjective(objectiveId string, results []*assessment.Result) (op float64, err error) {
}

func (s Service) calculateOEForRequirement(requirementId string, results []*assessment.Result) (op float64, err error) {
	t := time.Now().Add(-DefaultOETime)

	var n = 0

	// filter results
	for _, result := range results {
		if result.GetTimestamp().GetUtcOffset().Seconds > t.Unix() {
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
