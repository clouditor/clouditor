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
	"context"
	"fmt"
	"io"
	"os"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/policies"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var log *logrus.Entry

var listId int32

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
	_, err = s.handleEvidence(req.Evidence)

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
	var evidence *evidence.Evidence
	var err error

	for {
		evidence, err = stream.Recv()

		// TODO: Catch error?
		_, _ = s.handleEvidence(evidence)

		if err == io.EOF {
			log.Infof("Stopped receiving streamed evidences")

			return stream.SendAndClose(&emptypb.Empty{})
		}

	}
}

func (s Service) handleEvidence(evidence *evidence.Evidence) (result *evidence.Evidence, err error) {
	log.Infof("Received evidence for resource %s", evidence.ResourceId)
	log.Debugf("Evidence: %+v", evidence)

	var file string

	listId++
	evidence.Id = fmt.Sprintf("%d", listId)

	if len(evidence.ApplicableMetrics) == 0 {
		log.Warnf("Could not find a valid metric for evidence of resource %s", evidence.ResourceId)
	}

	// TODO(oxisto): actually look up metric via orchestrator
	for _, metric := range evidence.ApplicableMetrics {
		var baseDir string = "."

		// check, if we are in the root of Clouditor
		if _, err := os.Stat("policies"); os.IsNotExist(err) {
			// in tests, we are relative to our current package
			baseDir = "../../"
		}

		file = fmt.Sprintf("%s/policies/metric%d.rego", baseDir, metric)

		// TODO(oxisto): use go embed
		data, err := policies.RunEvidence(file, evidence)
		if err != nil {
			log.Errorf("Could not evaluate evidence: %v", err)

			if s.ResultHook != nil {
				go s.ResultHook(nil, err)
			}

			return nil, err
		}

		log.Infof("Evaluated evidence as %v", data["compliant"])

		result := &assessment.Result{
			ResourceId: evidence.ResourceId,
			Compliant:  data["compliant"].(bool),
			MetricId:   int32(metric),
		}

		// just a little hack to quickly enable multiple results per resource
		s.results[fmt.Sprintf("%s-%d", evidence.ResourceId, metric)] = result

		if s.ResultHook != nil {
			go s.ResultHook(result, nil)
		}
	}

	return
}

func (s Service) ListAssessmentResults(_ context.Context, _ *assessment.ListAssessmentResultsRequest) (res *assessment.ListAssessmentResultsResponse, err error) {
	res = new(assessment.ListAssessmentResultsResponse)
	res.Results = []*assessment.Result{}

	for _, result := range s.results {
		res.Results = append(res.Results, result)
	}

	return
}
