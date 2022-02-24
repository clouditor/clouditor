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
	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/policies"
	service_orchestrator "clouditor.io/clouditor/service/orchestrator"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var log *logrus.Entry

func init() {
	log = logrus.WithField("component", "assessment")
}

const (
	// EvictionTime is the time after which an entry in the metric configuration is invalid
	EvictionTime = time.Hour * 1
)

type cachedConfiguration struct {
	cachedAt time.Time
	*assessment.MetricConfiguration
}

// Service is an implementation of the Clouditor Assessment service. It should not be used directly,
// but rather the NewService constructor should be used. It implements the AssessmentServer interface.
type Service struct {
	// Embedded for FWD compatibility
	assessment.UnimplementedAssessmentServer

	// evidenceStoreStream sends evidences to the Evidence Store
	evidenceStoreStream  evidence.EvidenceStore_StoreEvidencesClient
	evidenceStoreAddress string

	// orchestratorStream sends ARs to the Orchestrator
	orchestratorStream  orchestrator.Orchestrator_StoreAssessmentResultsClient
	orchestratorClient  orchestrator.OrchestratorClient
	orchestratorAddress string

	// resultHooks is a list of hook functions that can be used if one wants to be
	// informed about each assessment result
	resultHooks []assessment.ResultHookFunc
	// hookMutex is used for (un)locking result hook calls
	hookMutex sync.Mutex

	// Currently, results are just stored as a map (=in-memory). In the future, we will use a DB.
	results map[string]*assessment.AssessmentResult

	// cachedConfigurations holds cached metric configurations for faster access with key being the corresponding
	// metric name
	cachedConfigurations map[string]cachedConfiguration

	authorizer api.Authorizer
}

const (
	// DefaultEvidenceStoreAddress specifies the default gRPC address of the evidence store.
	DefaultEvidenceStoreAddress = "localhost:9090"

	// DefaultOrchestratorAddress specifies the default gRPC address of the orchestrator.
	DefaultOrchestratorAddress = "localhost:9090"
)

// ServiceOption is a functional option type to configure the assessment service.
type ServiceOption func(*Service)

// WithEvidenceStoreAddress is an option to configure the evidence store gRPC address.
func WithEvidenceStoreAddress(address string) ServiceOption {
	return func(s *Service) {
		s.evidenceStoreAddress = address
	}
}

// WithOrchestratorAddress is an option to configure the orchestrator gRPC address.
func WithOrchestratorAddress(address string) ServiceOption {
	return func(s *Service) {
		s.orchestratorAddress = address
	}
}

// WithInternalAuthorizer is an option to use an authorizer to the internal Clouditor auth service.
func WithInternalAuthorizer(address string, username string, password string, opts ...grpc.DialOption) ServiceOption {
	return func(s *Service) {
		s.SetAuthorizer(api.NewInternalAuthorizerFromPassword(address, username, password, opts...))
	}
}

// NewService creates a new assessment service with default values.
func NewService(opts ...ServiceOption) *Service {
	s := &Service{
		results:              make(map[string]*assessment.AssessmentResult),
		evidenceStoreAddress: DefaultEvidenceStoreAddress,
		orchestratorAddress:  DefaultOrchestratorAddress,
		cachedConfigurations: make(map[string]cachedConfiguration),
	}

	// Apply any options
	for _, o := range opts {
		o(s)
	}

	return s
}

// SetAuthorizer implements UsesAuthorizer
func (s *Service) SetAuthorizer(auth api.Authorizer) {
	s.authorizer = auth
}

// Authorizer implements UsesAuthorizer
func (s *Service) Authorizer() api.Authorizer {
	return s.authorizer
}

// AssessEvidence is a method implementation of the assessment interface: It assesses a single evidence
func (s *Service) AssessEvidence(_ context.Context, req *assessment.AssessEvidenceRequest) (res *assessment.AssessEvidenceResponse, err error) {
	resourceId, err := req.Evidence.Validate()
	if err != nil {
		newError := fmt.Errorf("invalid evidence: %w", err)
		log.Error(newError)
		s.informHooks(nil, newError)

		res = &assessment.AssessEvidenceResponse{
			Status:        assessment.AssessEvidenceResponse_FAILED,
			StatusMessage: newError.Error(),
		}

		return res, status.Errorf(codes.InvalidArgument, "%v", newError)
	}

	// Assess evidence
	err = s.handleEvidence(req.Evidence, resourceId)

	if err != nil {
		res = &assessment.AssessEvidenceResponse{
			Status:        assessment.AssessEvidenceResponse_FAILED,
			StatusMessage: err.Error(),
		}

		newError := errors.New("error while handling evidence")
		log.Error(newError)

		return res, status.Errorf(codes.Internal, "%v", newError)
	}

	res = &assessment.AssessEvidenceResponse{
		Status: assessment.AssessEvidenceResponse_ASSESSED,
	}

	return res, nil
}

// AssessEvidences is a method implementation of the assessment interface: It assesses multiple evidences (stream) and responds with a stream.
func (s *Service) AssessEvidences(stream assessment.Assessment_AssessEvidencesServer) (err error) {
	var (
		req *assessment.AssessEvidenceRequest
		res *assessment.AssessEvidenceResponse
	)

	for {

		// TODO(all): Check context?
		req, err = stream.Recv()

		// If no more input of the stream is available, return
		if err == io.EOF {
			return nil
		}
		if err != nil {
			newError := fmt.Errorf("cannot receive stream request: %w", err)
			log.Error(newError)
			return status.Errorf(codes.Unknown, "%v", newError)
		}

		// Call AssessEvidence for assessing a single evidence
		assessEvidencesReq := &assessment.AssessEvidenceRequest{
			Evidence: req.Evidence,
		}
		res, err = s.AssessEvidence(context.Background(), assessEvidencesReq)
		if err != nil {
			log.Errorf("Error assessing evidence: %v", err)
		}

		// Send response back to the client
		err = stream.Send(res)
		if err != nil {
			newError := fmt.Errorf("cannot send response to the client: %w", err)
			log.Error(newError)
			return status.Errorf(codes.Unknown, "%v", newError)
		}
	}
}

// handleEvidence is the helper method for the actual assessment used by AssessEvidence and AssessEvidences
func (s *Service) handleEvidence(evidence *evidence.Evidence, resourceId string) (err error) {
	log.Infof("Evaluating evidence %s (%s) collected by %s at %v", evidence.Id, resourceId, evidence.ToolId, evidence.Timestamp)
	log.Debugf("Evidence: %+v", evidence)

	evaluations, err := policies.RunEvidence(evidence, s)
	if err != nil {
		newError := fmt.Errorf("could not evaluate evidence: %w", err)
		log.Error(newError)

		go s.informHooks(nil, newError)

		return newError
	}
	// The evidence is sent to the evidence store since it has been successfully evaluated
	// We could also just log, but an AR could be sent without an accompanying evidence in the Evidence Store
	err = s.sendToEvidenceStore(evidence)
	if err != nil {
		return fmt.Errorf("could not send evidence to the evidence store: %v", err)
	}

	for i, data := range evaluations {
		metricId := data.MetricId

		log.Infof("Evaluated evidence with metric '%v' as %v", metricId, data.Compliant)

		targetValue := data.TargetValue

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
				Operator:    data.Operator,
			},
			Compliant:             data.Compliant,
			EvidenceId:            evidence.Id,
			ResourceId:            resourceId,
			NonComplianceComments: "No comments so far",
		}

		// Just a little hack to quickly enable multiple results per resource
		s.results[fmt.Sprintf("%s-%d", resourceId, i)] = result

		// Inform hooks about new assessment result
		go s.informHooks(result, nil)

		// Send assessment result to the Orchestrator
		err = s.sendToOrchestrator(result)
		if err != nil {
			return fmt.Errorf("could not send assessment result to orchestrator: %v", err)
		}
	}
	return nil
}

// convertTargetValue converts v in a format accepted by protobuf (structpb.Value)
func convertTargetValue(v interface{}) (s *structpb.Value, err error) {
	var b []byte

	// json.Marshal and json.Unmarshaling is used instead of structpb.NewValue() which cannot handle json numbers
	if b, err = json.Marshal(v); err != nil {
		return nil, fmt.Errorf("JSON marshal failed: %w", err)
	}
	if err = json.Unmarshal(b, &s); err != nil {
		return nil, fmt.Errorf("JSON unmarshal failed: %w", err)
	}
	return

}

// informHooks informs the registered hook functions
func (s *Service) informHooks(result *assessment.AssessmentResult, err error) {
	s.hookMutex.Lock()
	defer s.hookMutex.Unlock()

	// Inform our hook, if we have any
	if s.resultHooks != nil {
		for _, hook := range s.resultHooks {
			// We could do hook concurrent again (assuming different hooks don't interfere with each other)
			hook(result, err)
		}
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

// initEvidenceStoreStream initializes the stream to the Evidence Store
func (s *Service) initEvidenceStoreStream(additionalOpts ...grpc.DialOption) error {
	log.Infof("Trying to establish a connection to evidence store service @ %v", s.evidenceStoreAddress)

	// Establish connection to evidence store gRPC service
	conn, err := grpc.Dial(s.evidenceStoreAddress,
		api.DefaultGrpcDialOptions(s, additionalOpts...)...,
	)
	if err != nil {
		return fmt.Errorf("could not connect to evidence store service: %w", err)
	}

	evidenceStoreClient := evidence.NewEvidenceStoreClient(conn)
	s.evidenceStoreStream, err = evidenceStoreClient.StoreEvidences(context.Background())
	if err != nil {
		return fmt.Errorf("could not set up stream for storing evidences: %w", err)
	}

	log.Infof("Connected to Evidence Store")

	return nil
}

// sendToEvidenceStore sends evidence e to the Evidence Store
func (s *Service) sendToEvidenceStore(e *evidence.Evidence) error {
	if s.evidenceStoreStream == nil {
		err := s.initEvidenceStoreStream()
		if err != nil {
			return fmt.Errorf("could not initialize stream to Evidence Store: %v", err)
		}
	}

	log.Infof("Sending evidence (%v) to Evidence Store", e.Id)

	err := s.evidenceStoreStream.Send(&evidence.StoreEvidenceRequest{Evidence: e})
	if err != nil {
		return err
	}
	return nil
}

// initOrchestratorStream initializes the stream to the Orchestrator
func (s *Service) initOrchestratorStream(additionalOpts ...grpc.DialOption) error {
	log.Infof("Trying to establish a connection to orchestrator service @ %v", s.orchestratorAddress)

	// Establish connection to orchestrator gRPC service
	conn, err := grpc.Dial(s.orchestratorAddress,
		api.DefaultGrpcDialOptions(s, additionalOpts...)...,
	)
	if err != nil {
		return fmt.Errorf("could not connect to orchestrator service: %w", err)
	}

	s.orchestratorClient = orchestrator.NewOrchestratorClient(conn)
	s.orchestratorStream, err = s.orchestratorClient.StoreAssessmentResults(context.Background())
	if err != nil {
		return fmt.Errorf("could not set up stream for storing assessment results: %w", err)
	}

	log.Infof("Connected to Orchestrator")

	return nil
}

// sendToOrchestrator sends the assessment result to the Orchestrator
func (s *Service) sendToOrchestrator(result *assessment.AssessmentResult) error {
	if s.orchestratorStream == nil || s.orchestratorClient == nil {
		err := s.initOrchestratorStream()
		if err != nil {
			return fmt.Errorf("could not initialize stream to Orchestrator: %v", err)
		}
	}

	log.Infof("Sending assessment result (%v) to Orchestrator", result.Id)

	req := &orchestrator.StoreAssessmentResultRequest{
		Result: result,
	}
	err := s.orchestratorStream.Send(req)
	if err != nil {
		log.Errorf("Error when sending assessment result to Orchestrator: %v", err)
		return err
	}

	return nil
}

// MetricConfiguration implements MetricConfigurationSource by getting the corresponding metric configuration for the
// default target cloud service
func (s *Service) MetricConfiguration(metric string) (config *assessment.MetricConfiguration, err error) {
	var (
		ok    bool
		cache cachedConfiguration
	)

	// Lazy init of connection
	if s.orchestratorStream == nil || s.orchestratorClient == nil {
		err = s.initOrchestratorStream()
		if err != nil {
			return nil, fmt.Errorf("could not initialize connection to orchestrator: %v", err)
		}
	}

	// Retrieve our cached entry
	cache, ok = s.cachedConfigurations[metric]

	// Check if entry is not there or is expired
	if !ok || cache.cachedAt.After(time.Now().Add(EvictionTime)) {
		config, err = s.orchestratorClient.GetMetricConfiguration(context.Background(), &orchestrator.GetMetricConfigurationRequest{
			ServiceId: service_orchestrator.DefaultTargetCloudServiceId,
			MetricId:  metric,
		})

		if err != nil {
			return nil, fmt.Errorf("could not retrieve metric configuration for %s: %w", metric, err)
		}

		cache = cachedConfiguration{
			cachedAt:            time.Now(),
			MetricConfiguration: config,
		}

		// Update the metric configuration
		s.cachedConfigurations[metric] = cache
	}

	return cache.MetricConfiguration, nil
}
