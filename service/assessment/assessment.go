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
	"clouditor.io/clouditor/policies"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var log *logrus.Entry

var listId int32

//go:generate protoc -I ../../proto -I ../../third_party assessment.proto evidence.proto --go_out=../.. --go-grpc_out=../..  --openapi_out=../../openapi/assessment

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

func (s Service) StoreEvidence(ctx context.Context, req *assessment.StoreEvidenceRequest) (res *assessment.Evidence, err error) {
	res, err = s.handleEvidence(req.Evidence)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error while handling evidence: %v", err)
	}

	return
}

func (s Service) StreamEvidences(stream assessment.Assessment_StreamEvidencesServer) error {
	var evidence *assessment.Evidence
	var err error

	for {
		evidence, err = stream.Recv()

		_, _ = s.handleEvidence(evidence)

		if err == io.EOF {
			log.Infof("Stopped receiving streamed evidence")

			return stream.SendAndClose(&emptypb.Empty{})
		}

	}
}

func (s Service) handleEvidence(evidence *assessment.Evidence) (result *assessment.Evidence, err error) {
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

	result = evidence

	return
}

func (s Service) ListAssessmentResults(ctx context.Context, req *assessment.ListAssessmentResultsRequest) (res *assessment.ListAssessmentResultsResponse, err error) {
	res = new(assessment.ListAssessmentResultsResponse)
	res.Results = []*assessment.Result{}

	for _, result := range s.results {
		res.Results = append(res.Results, result)
	}

	return
}
