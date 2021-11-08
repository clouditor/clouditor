// Copyright 2021 Fraunhofer AISEC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//           $$\                           $$\ $$\   $$\
//           $$ |                          $$ |\__|  $$ |
//  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
// $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
// $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
// $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
// \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
//  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
//
// This file is part of Clouditor Community Edition.

package assessment

import (
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/policies"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
)

var log *logrus.Entry

//var listId int32

func init() {
	log = logrus.WithField("component", "assessment")
}

type Service struct {
	ResultHook func(result *assessment.Result, err error)

	results map[string]*assessment.Result
	assessment.UnimplementedAssessmentServer
}

func NewService() assessment.AssessmentServer {
	return &Service{
		results: make(map[string]*assessment.Result),
	}
}

func (s Service) AssessEvidence(_ context.Context, req *assessment.AssessEvidenceRequest) (res *assessment.AssessEvidenceResponse, err error) {
	err = s.handleEvidence(req.Evidence)

	if err != nil {
		res = &assessment.AssessEvidenceResponse{
			Status: false,
		}

		return res, status.Errorf(codes.Internal, "Error while handling evidence: %v", err)
	}

	res = &assessment.AssessEvidenceResponse{
		Status: true,
	}

	return res, nil
}

func (s Service) AssessEvidences(stream assessment.Assessment_AssessEvidencesServer) error {
	var e *evidence.Evidence
	var err error

	for {
		e, err = stream.Recv()

		// TODO: Catch error?
		_ = s.handleEvidence(e)

		if err == io.EOF {
			log.Infof("Stopped receiving streamed evidences")

			return stream.SendAndClose(&emptypb.Empty{})
		}

	}
}

func (s Service) handleEvidence(evidence *evidence.Evidence) error {
	log.Println("Running evidence for resourceId", evidence.ResourceId, " and ID:", evidence.Id)
	log.Debugf("Evidence: %+v", evidence)

	// TODO(all): The discovery already sets up the (UU)ID?
	//listId++
	//evidence.Id = fmt.Sprintf("%d", listId)

	//if len(evidence.ApplicableMetrics) == 0 {
	//	log.Warnf("Could not find a valid metric for evidence of resource %s", evidence.ResourceId)
	//}

	// TODO(oxisto): use go embed
	evaluations, err := policies.RunEvidence(evidence)
	if err != nil {
		log.Errorf("Could not evaluate evidence: %v", err)

		if s.ResultHook != nil {
			go s.ResultHook(nil, err)
		}

		return err
	}

	for i, data := range evaluations {
		//var metricID int64
		//if metricID, err = data["metricID"].(json.Number).Int64(); err != nil {
		//	return fmt.Errorf("could not convert metricID: %v", metricID)
		//}
		log.Infof("Evaluated evidence with metric '%v' as %v", data["name"], data["compliant"])

		result := &assessment.Result{
			// TODO(lebogg): Remove metric name hack after demo
			ResourceId: evidence.ResourceId + " with metric " + data["name"].(string),
			Compliant:  data["compliant"].(bool),
			// TODO(lebogg): Currently no metric IDs are used
			MetricId: int32(0),
		}
		// just a little hack to quickly enable multiple results per resource
		s.results[fmt.Sprintf("%s-%d", evidence.ResourceId, i)] = result

		// TODO(oxisto): What is this for? (lebogg)
		if s.ResultHook != nil {
			go s.ResultHook(result, nil)
		}
	}
	fmt.Println()

	//}

	return nil
}

func (s Service) ListAssessmentResults(_ context.Context, _ *assessment.ListAssessmentResultsRequest) (res *assessment.ListAssessmentResultsResponse, err error) {
	res = new(assessment.ListAssessmentResultsResponse)
	res.Results = []*assessment.Result{}

	for _, result := range s.results {
		res.Results = append(res.Results, result)
	}

	return
}
