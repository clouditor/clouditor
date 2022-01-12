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
	// resultHooks is a list of hook functions that can be used if one wants to be
	assessment.UnimplementedAssessmentServer

	Configuration

	// evidenceStoreStream send evidences to the Evidence Store
	evidenceStoreStream evidence.EvidenceStore_StoreEvidencesClient

	// ResultHook is a hook function that can be used if one wants to be
	// informed about each assessment result
	resultHooks []assessment.ResultHookFunc

	// Currently, results are just stored as a map (=in-memory). In the future, we will use a DB.
	results map[string]*assessment.AssessmentResult
}

// Configuration contains an assessment service's information about connection to other services, e.g. Evidence Store
type Configuration struct {
	evidenceStoreTargetAddress string
}

// NewService creates a new assessment service with default values
func NewService() *Service {
	return &Service{
		results:       make(map[string]*assessment.AssessmentResult),
		Configuration: Configuration{evidenceStoreTargetAddress: "localhost:9090"},
	}
}

// Start starts assessment by setting up connections to other services. In the future, also with manual configurations in request
// TODO: Maybe rename (proto) to ConfigureConnections. For now, we do it in assess evidence
// TODO: Add CLI Command
func (s *Service) Start(_ context.Context, _ *assessment.StartAssessmentRequest) (resp *assessment.StartAssessmentResponse, err error) {
	// Establish connection to evidenceStore component
	conn, err := grpc.Dial(s.Configuration.evidenceStoreTargetAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not connect to evidence store service: %v", err)
	}
	// defer conn.Close()
	evidenceStoreClient := evidence.NewEvidenceStoreClient(conn)
	s.evidenceStoreStream, err = evidenceStoreClient.StoreEvidences(context.Background())
	if err != nil {
		// TODO(all): We dont have to print it, since the caller (e.g. CLI) can/should do it, right?
		return nil, status.Errorf(codes.Internal, "could not set up stream for storing evidences: %v", err)
	}

	log.Infof("Assessment Started")
	return &assessment.StartAssessmentResponse{}, nil
}

// AssessEvidence is a method implementation of the assessment interface: It assesses a single evidence
func (s Service) AssessEvidence(_ context.Context, req *assessment.AssessEvidenceRequest) (res *assessment.AssessEvidenceResponse, err error) {
	resourceId, err := req.Evidence.Validate()
	if err != nil {
		log.Errorf("Invalid evidence: %v", err)
		newError := fmt.Errorf("invalid evidence: %w", err)

		s.informHooks(nil, newError)

		res = &assessment.AssessEvidenceResponse{
			Status: false,
		}

		return res, status.Errorf(codes.InvalidArgument, "invalid req: %v", newError)
	}
	// Set up stream to Evidence Store.
	// TODO(lebogg): If assessment is used as standalone service (without fwd evidences) maybe adapt code, i.e. no error when ConfigureConnections sets no evidence store
	if s.evidenceStoreStream == nil {
		if err = s.setEvidenceStoreStream(); err != nil {
			return nil, err
		}
	}

	// Assess evidence
	err = s.handleEvidence(req.Evidence, resourceId)

	err = s.handleEvidence(req.Evidence, resourceId)
	if err != nil {
		res = &assessment.AssessEvidenceResponse{
			Status: false,
		}

		return res, status.Errorf(codes.Internal, "error while handling evidence: %v", err)
	}

	res = &assessment.AssessEvidenceResponse{
		Status: true,
	}

	return res, nil
}

// AssessEvidences is a method implementation of the assessment interface: It assesses multiple evidences (stream)
func (s Service) AssessEvidences(stream assessment.Assessment_AssessEvidencesServer) (err error) { // Set up stream to Evidence Store.
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

		// Call AssessEvidence for assessing a single evidence
		assessEvidencesReq := &assessment.AssessEvidenceRequest{
			Evidence: req.Evidence,
		}
		_, err = s.AssessEvidence(context.Background(), assessEvidencesReq)
		if err != nil {
			return err
		}
	}
}

// handleEvidence is the helper method for the actual assessment used by AssessEvidence and AssessEvidences
func (s Service) handleEvidence(evidence *evidence.Evidence, resourceId string) error {

	log.Infof("Running evidence %s (%s) collected by %s at %v", evidence.Id, resourceId, evidence.ToolId, evidence.Timestamp)
	log.Debugf("Evidence: %+v", evidence)

	evaluations, err := policies.RunEvidence(evidence)
	if err != nil {
		newError := fmt.Errorf("could not evaluate evidence: %w", err)
		log.Errorf(newError.Error())

		s.informHooks(nil, newError)

		return newError
	}
	// The evidence is sent to the evidence store since it has been successfully evaluated
	go s.sendToEvidenceStore(evidence)

	for i, data := range evaluations {
		metricId := data["metricId"].(string)

		log.Infof("Evaluated evidence with metric '%v' as %v", metricId, data["compliant"])

		// Get output values of Rego evaluation. If they are not given, the zero value is used
		operator, _ := data["operator"].(string)
		targetValue, _ := data["target_value"]
		compliant, _ := data["compliant"].(bool)

		result := &assessment.AssessmentResult{
			Id:        uuid.NewString(),
			Timestamp: timestamppb.Now(),
			MetricId:  metricId,
			MetricConfiguration: &assessment.MetricConfiguration{
				TargetValue: convertTargetValue(targetValue),
				Operator:    operator,
			},
			Compliant:             compliant,
			EvidenceId:            evidence.Id,
			ResourceId:            resourceId,
			NonComplianceComments: "No comments so far",
		}

		// Just a little hack to quickly enable multiple results per resource
		s.results[fmt.Sprintf("%s-%d", resourceId, i)] = result

		// TODO(all): Currently, we sent result via hook. But I think it would be cleaner to do it straight here.
		s.informHooks(result, nil)
	}

	return nil
}

func convertTargetValue(value interface{}) (convertedTargetValue *structpb.Value) {
	if str, ok := value.(string); ok {
		convertedTargetValue = &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: str}}
	} else {
		// TODO(lebogg): FOr now, for Testing!!! Change!
		convertedTargetValue = &structpb.Value{Kind: &structpb.Value_BoolValue{BoolValue: true}}
	}
	return
}

// informHooks informs the registered hook functions
func (s Service) informHooks(result *assessment.AssessmentResult, err error) {
	// Inform our hook, if we have any
	if s.resultHooks != nil {
		for _, hook := range s.resultHooks {
			go hook(result, err)
		}
	}
}

// sendToEvidenceStore sends the provided evidence to the evidence store
func (s Service) sendToEvidenceStore(e *evidence.Evidence) {
	// Check if evidenceStoreStream is already established (via Start method)
	if s.evidenceStoreStream != nil {
		err := s.evidenceStoreStream.Send(&evidence.StoreEvidenceRequest{Evidence: e})
		if err != nil {
			log.WithError(err)
		}

	} else {
		log.Errorf("Evidence couldn't be sent to Evidence Store. Did you establish connection by starting the assessment?")
	}
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

func (s *Service) RegisterAssessmentResultHook(assessmentResultsHook func(result *assessment.AssessmentResult, err error)) {
	s.resultHooks = append(s.resultHooks, assessmentResultsHook)
}

func (s *Service) setEvidenceStoreStream() error {
	log.Infof("Establishing connection to Evidence Store")
	// Establish connection to evidenceStore component
	conn, err := grpc.Dial(s.Configuration.evidenceStoreTargetAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return status.Errorf(codes.Internal, "could not connect to evidence store service: %v", err)
	}
	evidenceStoreClient := evidence.NewEvidenceStoreClient(conn)
	s.evidenceStoreStream, err = evidenceStoreClient.StoreEvidences(context.Background())
	if err != nil {
		return status.Errorf(codes.Internal, "could not set up stream for storing evidences: %v", err)
	}
	log.Infof("Connected to Evidence Store")
	return nil
}
