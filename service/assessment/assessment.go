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
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
)

var log *logrus.Entry

func init() {
	log = logrus.WithField("component", "assessment")
}

/*
Service is an implementation of the Clouditor Assessment service. It should not be used directly, but rather the
NewService constructor should be used. It implements the AssessmentServer interface
*/
type Service struct {
	assessment.UnimplementedAssessmentServer

	// TODO(lebogg): comment
	evidenceStoreStream evidence.EvidenceStore_StoreEvidencesClient

	// ResultHook is a hook function that can be used if one wants to be
	// informed about each assessment result
	ResultHook func(result *assessment.AssessmentResult, err error)

	results map[string]*assessment.AssessmentResult
}

func NewService() assessment.AssessmentServer {
	return &Service{
		results: make(map[string]*assessment.AssessmentResult),
	}
}

func (s Service) Start(_ context.Context, _ *assessment.StartAssessmentRequest) (*assessment.StartAssessmentResponse, error) {
	var (
		cc                  *grpc.ClientConn
		evidenceStoreClient evidence.EvidenceStoreClient
		err                 error
	)

	// Establish connection to evidenceStore component
	// TODO(lebogg): port hardcoded for now -> add configuration struct or/and field in request for setting the port
	cc, err = grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not connect to evidence store service: %v", err)
	}
	evidenceStoreClient = evidence.NewEvidenceStoreClient(cc)
	// TODO(lebogg): Find out which errors can occur and add tests accordingly
	s.evidenceStoreStream, err = evidenceStoreClient.StoreEvidences(context.Background())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not set up stream for storing evidences: %v", err)
	}
	return &assessment.StartAssessmentResponse{}, nil
}

// AssessEvidence is a method implementation of the assessment interface: It assesses a single evidence
func (s Service) AssessEvidence(_ context.Context, req *assessment.AssessEvidenceRequest) (res *assessment.AssessEvidenceResponse, err error) {
	err = s.handleEvidence(req.Evidence)

	if err != nil {
		res = &assessment.AssessEvidenceResponse{
			Status: false,
		}

		return res, status.Errorf(codes.Internal, "Error while handling evidence: %v", err)
	}

	// TODO(lebogg): Send evidence to evidence store (as stream)

	res = &assessment.AssessEvidenceResponse{
		Status: true,
	}

	return res, nil
}

// AssessEvidences is a method implementation of the assessment interface: It assesses multiple evidences (stream)
// TODO(lebogg): Send evidence to evidence store (as stream)
func (s Service) AssessEvidences(stream assessment.Assessment_AssessEvidencesServer) (err error) {
	var (
		req *assessment.AssessEvidenceRequest
	)

	for {
		req, err = stream.Recv()

		if err != nil {
			// If no more input of the stream is available, return SendAndClose `error`
			if err == io.EOF {
				log.Infof("Stopped receiving streamed evidences")
				return stream.SendAndClose(&emptypb.Empty{})
			}
			return err
		}

		err = s.handleEvidence(req.Evidence)
		if err != nil {
			return status.Errorf(codes.Internal, "Error while handling evidence: %v", err)
		}

	}
}

// handleEvidence is the common evidence assessment of AssessEvidence and AssessEvidences
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

		result := &assessment.AssessmentResult{
			Id:        uuid.NewString(),
			Timestamp: timestamppb.Now(),
			MetricId:  metricId,
			MetricConfiguration: &assessment.MetricConfiguration{
				TargetValue: targetValue,
				Operator:    operator,
			},
			Compliant:             compliant,
			EvidenceId:            evidence.Id,
			ResourceId:            resourceId,
			NonComplianceComments: "No comments so far",
		}

		// Just a little hack to quickly enable multiple results per resource
		s.results[fmt.Sprintf("%s-%d", resourceId, i)] = result

		// Inform our hook, if we have any
		if s.ResultHook != nil {
			go s.ResultHook(result, nil)
		}
	}

	return nil
}

// ListAssessmentResults is a method implementation of the assessment interface
func (s Service) ListAssessmentResults(_ context.Context, _ *assessment.ListAssessmentResultsRequest) (res *assessment.ListAssessmentResultsResponse, err error) {
	res = new(assessment.ListAssessmentResultsResponse)
	res.Results = []*assessment.AssessmentResult{}

	for _, result := range s.results {
		res.Results = append(res.Results, result)
	}

	return
}
