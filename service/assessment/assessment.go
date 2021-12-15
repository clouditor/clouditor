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

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/policies"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var log *logrus.Entry

func init() {
	log = logrus.WithField("component", "assessment")
}

type Service struct {
	// ResultHook is a hook function that can be used if one wants to be
	// informed about each assessment result
	ResultHook func(result *assessment.Result, err error)

	AssessmentResults map[string]*assessment.Result
	assessment.UnimplementedAssessmentServer
}

func NewService() *Service {
	return &Service{
		AssessmentResults: make(map[string]*assessment.Result),
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

		if err == io.EOF {
			log.Infof("Stopped receiving streamed evidences")

			return stream.SendAndClose(&emptypb.Empty{})
		}

		// TODO: Catch error?
		_ = s.handleEvidence(e)
	}
}

func (s Service) handleEvidence(evidence *evidence.Evidence) error {
	resourceId, err := evidence.Validate()
	if err != nil {
		return fmt.Errorf("invalid evidence: %w", err)
	}

	log.Infof("Running evidence %s (%s) collected by %s at %v", evidence.Id, resourceId, evidence.ToolId, evidence.Timestamp)
	log.Debugf("Evidence: %+v", evidence)

	evaluations, err := policies.RunEvidence(evidence)
	if err != nil {
		log.Errorf("Could not evaluate evidence: %v", err)

		// Inform our hook, if we have any
		if s.ResultHook != nil {
			go s.ResultHook(nil, err)
		}

		return err
	}

	for i, data := range evaluations {
		metricId := data["metricId"].(string)

		log.Infof("Evaluated evidence with metric '%v' as %v", metricId, data["compliant"])

		// Get output values of Rego evaluation. If they are not given, the zero value is used
		operator, _ := data["operator"].(string)
		targetValue, _ := data["target_value"].(*structpb.Value)
		compliant, _ := data["compliant"].(bool)

		result := &assessment.Result{
			Id:        uuid.NewString(),
			Timestamp: timestamppb.Now(),
			MetricId:  metricId,
			MetricData: &assessment.MetricConfiguration{
				TargetValue: targetValue,
				Operator:    operator,
			},
			Compliant:             compliant,
			EvidenceId:            evidence.Id,
			ResourceId:            resourceId,
			NonComplianceComments: "No comments so far",
		}

		// Just a little hack to quickly enable multiple AssessmentResults per resource
		s.AssessmentResults[fmt.Sprintf("%s-%d", resourceId, i)] = result

		// Inform our hook, if we have any
		if s.ResultHook != nil {
			go s.ResultHook(result, nil)
		}
	}

	return nil
}

func (s Service) ListAssessmentResults(_ context.Context, _ *assessment.ListAssessmentResultsRequest) (res *assessment.ListAssessmentResultsResponse, err error) {
	res = new(assessment.ListAssessmentResultsResponse)
	res.Results = []*assessment.Result{}

	for _, result := range s.AssessmentResults {
		res.Results = append(res.Results, result)
	}

	return
}

func (s *Service) RegisterAssessmentResultHook(assessmentResultsHook func(result *assessment.Result, err error)) {
	s.ResultHook = assessmentResultsHook
}
