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
	"encoding/json"
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
	"sync"
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
	// Embedded for FWD compatibility
	assessment.UnimplementedAssessmentServer

	Configuration

	// evidenceStoreStream send evidences to the Evidence Store
	evidenceStoreStream evidence.EvidenceStore_StoreEvidencesClient

	// resultHooks is a list of hook functions that can be used if one wants to be
	// informed about each assessment result
	resultHooks []assessment.ResultHookFunc
	// mu is used for (un)locking result hook calls
	mu sync.Mutex

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

// Configure configures assessment by setting up connections to other services. In the future, also with manual configurations in request
// TODO: Add CLI Command
func (*Service) Configure(_ context.Context, _ *assessment.ConfigureAssessmentRequest) (resp *assessment.ConfigureAssessmentResponse, err error) {
	// TODO
	// s.evidenceStoreTargetAddress = someThingFromRequest
	return &assessment.ConfigureAssessmentResponse{}, nil
}

// AssessEvidence is a method implementation of the assessment interface: It assesses a single evidence
func (s *Service) AssessEvidence(_ context.Context, req *assessment.AssessEvidenceRequest) (res *assessment.AssessEvidenceResponse, err error) {
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
			return nil, status.Errorf(codes.Internal, "could not assess evidence: %v", err)
		}
	}

	// Assess evidence
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
func (s *Service) AssessEvidences(stream assessment.Assessment_AssessEvidencesServer) (err error) { // Set up stream to Evidence Store.
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
func (s *Service) handleEvidence(evidence *evidence.Evidence, resourceId string) error {

	log.Infof("Running evidence %s (%s) collected by %s at %v", evidence.Id, resourceId, evidence.ToolId, evidence.Timestamp)
	log.Debugf("Evidence: %+v", evidence)

	evaluations, err := policies.RunEvidence(evidence)
	if err != nil {
		newError := fmt.Errorf("could not evaluate evidence: %w", err)
		log.Errorf(newError.Error())

		go s.informHooks(nil, newError)

		return newError
	}
	// The evidence is sent to the evidence store since it has been successfully evaluated
	go s.sendToEvidenceStore(evidence)

	for i, data := range evaluations {
		metricId := data["metricId"].(string)

		log.Infof("Evaluated evidence with metric '%v' as %v", metricId, data["compliant"])

		// Get output values of Rego evaluation. If they are not given, the zero value is used
		operator, _ := data["operator"].(string)
		targetValue := data["target_value"]
		compliant, _ := data["compliant"].(bool)

		convertedTargetValue, err := convertTargetValue(targetValue)
		if err != nil {
			return fmt.Errorf("could not convert target value: %w", err)
		}

		result := &assessment.AssessmentResult{
			Id:        uuid.NewString(),
			Timestamp: timestamppb.Now(),
			MetricId:  metricId,
			MetricConfiguration: &assessment.MetricConfiguration{
				TargetValue: convertedTargetValue,
				Operator:    operator,
			},
			Compliant:             compliant,
			EvidenceId:            evidence.Id,
			ResourceId:            resourceId,
			NonComplianceComments: "No comments so far",
		}

		// Just a little hack to quickly enable multiple results per resource
		s.results[fmt.Sprintf("%s-%d", resourceId, i)] = result

		// TODO(all): Currently, we send result via hook. But I think it would be cleaner to do it straight here.
		go s.informHooks(result, nil)
	}

	return nil
}

// convertTargetValue converts the given target_value in a format accepted by protobuf
// Note: Therefore, we have to be specific about the format of target values in the REGO code
func convertTargetValue(value interface{}) (convertedTargetValue *structpb.Value, err error) {
	if valueStr, ok := value.(string); ok {
		convertedTargetValue = &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: valueStr}}
		return
	}
	if valueBool, ok := value.(bool); ok {
		convertedTargetValue = &structpb.Value{Kind: &structpb.Value_BoolValue{BoolValue: valueBool}}
		return
	}
	if _, ok := value.(json.Number); ok {
		convertedTargetValue, err = assertNumber(value)
		if err != nil {
			err = fmt.Errorf("could not convert jsonNumber %v to Float: %w", value, err)
		}
		return
	}
	if _, ok := value.(int); ok {
		convertedTargetValue, err = assertNumber(value)
		if err != nil {
			err = fmt.Errorf("could not assert integer value %v: %w", value, err)
		}
		return
	}
	if _, ok := value.(float64); ok {
		convertedTargetValue, err = assertNumber(value)
		if err != nil {
			err = fmt.Errorf("could not assert float64 value %v: %w", value, err)
		}
		return
	}
	if _, ok := value.(float32); ok {
		convertedTargetValue, err = assertNumber(value)
		if err != nil {
			err = fmt.Errorf("could not assert float32 value %v: %w", value, err)
		}
		return
	}
	if listOfValues, ok := value.([]interface{}); ok {
		var listOfConvertedValues []*structpb.Value
		for _, v := range listOfValues {
			if valueStr, ok := v.(string); ok {
				listOfConvertedValues = append(listOfConvertedValues, &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: valueStr}})
			} else if valueBool, ok := v.(bool); ok {
				listOfConvertedValues = append(listOfConvertedValues, &structpb.Value{Kind: &structpb.Value_BoolValue{BoolValue: valueBool}})
			} else if valueMap, ok := v.(map[string]interface{}); ok {
				var mapOfValues = make(map[string]*structpb.Value)
				for k, mapValue := range valueMap {
					if k == "runtimeLanguage" {
						if mapValueString, ok := mapValue.(string); ok {
							mapOfValues[k] = &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: mapValueString}}
						} else {
							return nil, fmt.Errorf("runtimeLanguage assertion failed. Wanted string but got %T (%v)", mapValue, mapValue)
						}
					} else if k == "runtimeVersion" {
						mapOfValues[k], err = assertNumber(mapValue)
						if err != nil {
							return nil, fmt.Errorf("runtimeVersion assertion failed. Could not assert '%v': %w", mapValue, err)
						}
					} else {
						return nil, fmt.Errorf("no supported key. Check if it is one of the supported (e.g. runtimeLanguage or runtimeVersion)." +
							"Or consider adding it to the code")
					}
				}
				listOfConvertedValues = append(listOfConvertedValues, &structpb.Value{Kind: &structpb.Value_StructValue{StructValue: &structpb.Struct{Fields: mapOfValues}}})
			} else {
				// When v is no string or boolean, check for numbers. If assertNumber fails, throw error.
				var number *structpb.Value
				number, err = assertNumber(v)
				if err != nil {
					return nil, fmt.Errorf("after it is no string or boolean, it is also not a supported number. Got %T (%v): %w", v, v, err)
				}
				listOfConvertedValues = append(listOfConvertedValues, number)
			}
		}
		convertedTargetValue = &structpb.Value{Kind: &structpb.Value_ListValue{ListValue: &structpb.ListValue{Values: listOfConvertedValues}}}
		return
	}
	return nil, fmt.Errorf("no assertion found for %v (%T). Consider adding it", value, value)
}

func assertNumber(rawTargetValue interface{}) (*structpb.Value, error) {
	if valueJSONNumber, ok := rawTargetValue.(json.Number); ok {
		asFloat, err := valueJSONNumber.Float64()
		if err != nil {
			return nil, fmt.Errorf("could not assert JSON number: %w", err)
		}
		return &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: asFloat}}, nil
	}
	if valueFloat64, ok := rawTargetValue.(float64); ok {
		return &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: valueFloat64}}, nil
	}
	if valueFloat32, ok := rawTargetValue.(float32); ok {
		return &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: float64(valueFloat32)}}, nil
	}
	if valueInt, ok := rawTargetValue.(int); ok {
		return &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: float64(valueInt)}}, nil
	}
	return &structpb.Value{}, fmt.Errorf("number assertions not possible for '%v' (%T)", rawTargetValue, rawTargetValue)

}

// informHooks informs the registered hook functions
func (s *Service) informHooks(result *assessment.AssessmentResult, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Inform our hook, if we have any
	if s.resultHooks != nil {
		for _, hook := range s.resultHooks {
			// TODO(all): We could do hook concurrent again (assuming different hooks don't interfere with each other)
			hook(result, err)
		}
	}
}

// sendToEvidenceStore sends the provided evidence to the evidence store
func (s *Service) sendToEvidenceStore(e *evidence.Evidence) {
	// Check if evidenceStoreStream is already established (via Configure method)
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
func (s *Service) ListAssessmentResults(_ context.Context, _ *assessment.ListAssessmentResultsRequest) (res *assessment.ListAssessmentResultsResponse, err error) {
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
		return fmt.Errorf("could not connect to evidence store service: %v", err)
	}
	evidenceStoreClient := evidence.NewEvidenceStoreClient(conn)
	s.evidenceStoreStream, err = evidenceStoreClient.StoreEvidences(context.Background())
	if err != nil {
		return fmt.Errorf("could not set up stream for storing evidences: %v", err)
	}
	log.Infof("Connected to Evidence Store")
	return nil
}
