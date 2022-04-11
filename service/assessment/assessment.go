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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/policies"
	service_orchestrator "clouditor.io/clouditor/service/orchestrator"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	log *logrus.Entry
)

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
	evidenceStoreStreams *api.StreamsOf[evidence.EvidenceStore_StoreEvidencesClient, *evidence.StoreEvidenceRequest, *evidence.StoreEvidenceResponse]
	evidenceStoreAddress string

	// orchestratorStream sends assessment results to the Orchestrator
	orchestratorStream  orchestrator.Orchestrator_StoreAssessmentResultsClient
	orchestratorClient  orchestrator.OrchestratorClient
	orchestratorAddress string
	// orchestratorChannel stores assessment results for the Orchestrator
	orchestratorChannel chan *assessment.AssessmentResult

	// resultHooks is a list of hook functions that can be used if one wants to be
	// informed about each assessment result
	resultHooks []assessment.ResultHookFunc
	// hookMutex is used for (un)locking result hook calls
	hookMutex sync.RWMutex

	// Currently, results are just stored as a map (=in-memory). In the future, we will use a DB.
	results     map[string]*assessment.AssessmentResult
	resultMutex sync.Mutex

	// cachedConfigurations holds cached metric configurations for faster access with key being the corresponding
	// metric name
	cachedConfigurations map[string]cachedConfiguration
	// TODO(oxisto): combine with hookMutex and replace with a generic version of a mutex'd map
	confMutex sync.Mutex

	authorizer api.Authorizer

	// grpcOpts contains additional gRPC options that will be appended to all gRPC dial calls. Useful for testing.
	grpcOpts []grpc.DialOption
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

// WithOAuth2Authorizer is an option to use an OAuth 2.0 authorizer
func WithOAuth2Authorizer(config *clientcredentials.Config) ServiceOption {
	return func(s *Service) {
		s.SetAuthorizer(api.NewOAuthAuthorizerFromClientCredentials(config))
	}
}

// WithAdditionalGRPCOpts is an option to configure additional gRPC options.
func WithAdditionalGRPCOpts(opts ...grpc.DialOption) ServiceOption {
	return func(s *Service) {
		s.grpcOpts = opts
	}
}

// NewService creates a new assessment service with default values.
func NewService(opts ...ServiceOption) *Service {
	s := &Service{
		results:              make(map[string]*assessment.AssessmentResult),
		evidenceStoreAddress: DefaultEvidenceStoreAddress,
		evidenceStoreStreams: api.NewStreamsOf(api.WithLogger[evidence.EvidenceStore_StoreEvidencesClient, *evidence.StoreEvidenceRequest, *evidence.StoreEvidenceResponse](log)),
		orchestratorAddress:  DefaultOrchestratorAddress,
		orchestratorChannel:  make(chan *assessment.AssessmentResult, 1000),
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

	// Validate evidence
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

	// Check if orchestrator stream exists
	if s.orchestratorStream == nil {
		err = s.initOrchestratorStream()
		if err != nil {
			log.Errorf("error initialising orchestrator stream: %v", err)
		}
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

	// Create response
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
		// Receive requests from client
		req, err = stream.Recv()

		// If no more input of the stream is available, return
		if errors.Is(err, io.EOF) {
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

		// Check for send errors
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			newError := fmt.Errorf("cannot send response to the client: %w", err)
			log.Error(newError)
			return status.Errorf(codes.Unknown, "%v", newError)
		}
	}
}

// handleEvidence is the helper method for the actual assessment used by AssessEvidence and AssessEvidences
func (s *Service) handleEvidence(ev *evidence.Evidence, resourceId string) (err error) {
	log.Debugf("Evaluating evidence %s (%s) collected by %s at %v", ev.Id, resourceId, ev.ToolId, ev.Timestamp)
	log.Tracef("Evidence: %+v", ev)

	evaluations, err := policies.RunEvidence(ev, s)
	if err != nil {
		newError := fmt.Errorf("could not evaluate evidence: %w", err)
		log.Error(newError)

		go s.informHooks(nil, newError)

		return newError
	}

	channel, err := s.evidenceStoreStreams.GetStream(s.evidenceStoreAddress, "Evidence Store", s.initEvidenceStoreStream, s.grpcOpts...)
	if err != nil {
		err = fmt.Errorf("could not get stream to evidence store (%s): %w", s.evidenceStoreAddress, err)
		log.Error(err)

		go s.informHooks(nil, err)

		return err
	}

	// Store evidence in evidenceStoreChannel
	channel.Send(&evidence.StoreEvidenceRequest{Evidence: ev})

	for i, data := range evaluations {
		metricId := data.MetricId

		log.Debugf("Evaluated evidence %v with metric '%v' as %v", ev.Id, metricId, data.Compliant)

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
			EvidenceId:            ev.Id,
			ResourceId:            resourceId,
			NonComplianceComments: "No comments so far",
		}

		s.resultMutex.Lock()
		// Just a little hack to quickly enable multiple results per resource
		s.results[fmt.Sprintf("%s-%d", resourceId, i)] = result
		s.resultMutex.Unlock()

		// Inform hooks about new assessment result
		go s.informHooks(result, nil)

		// Store assessment result in orchestratorChannel
		s.orchestratorChannel <- result
	}

	return nil
}

// convertTargetValue converts v in a format accepted by protobuf (structpb.Value)
func convertTargetValue(v interface{}) (s *structpb.Value, err error) {
	var b []byte

	// json.Marshal and json.Unmarshal is used instead of structpb.NewValue() which cannot handle json numbers
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
	s.hookMutex.RLock()
	hooks := s.resultHooks
	defer s.hookMutex.RUnlock()

	// Inform our hook, if we have any
	if len(hooks) > 0 {
		for _, hook := range hooks {
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
	s.hookMutex.Lock()
	defer s.hookMutex.Unlock()
	s.resultHooks = append(s.resultHooks, assessmentResultsHook)
}

// initEvidenceStoreStream initializes the stream to the Evidence Store
func (s *Service) initEvidenceStoreStream(URL string, additionalOpts ...grpc.DialOption) (stream evidence.EvidenceStore_StoreEvidencesClient, err error) {
	log.Infof("Trying to establish a connection to evidence store service @ %v", s.evidenceStoreAddress)

	// Establish connection to evidence store gRPC service
	conn, err := grpc.Dial(URL,
		api.DefaultGrpcDialOptions(URL, s, additionalOpts...)...,
	)
	if err != nil {
		return nil, fmt.Errorf("could not connect to evidence store service: %w", err)
	}

	evidenceStoreClient := evidence.NewEvidenceStoreClient(conn)
	stream, err = evidenceStoreClient.StoreEvidences(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not set up stream for storing evidences: %w", err)
	}

	return
}

// initOrchestratorStream initializes the stream to the Orchestrator
func (s *Service) initOrchestratorStream(additionalOpts ...grpc.DialOption) error {
	log.Infof("Trying to establish a connection to orchestrator service @ %v", s.orchestratorAddress)

	// Establish connection to orchestrator gRPC service
	conn, err := grpc.Dial(s.orchestratorAddress,
		api.DefaultGrpcDialOptions(s.orchestratorAddress, s, additionalOpts...)...,
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

	// Receive responses from Orchestrator
	// Currently we do not process the responses
	go func() {
		i := 1

		for {
			_, err := s.orchestratorStream.Recv()

			if errors.Is(err, io.EOF) {
				log.Debugf("no more responses from orchestrator stream: EOF")
				break
			}

			if err != nil {
				log.Errorf("error receiving response from orchestrator stream: %v", err)
				break
			}

			if i%100 == 0 {
				log.Debugf("orchestratorStream recv responses currently @ %v", i)
			}

			i++
		}
	}()

	// Send assessment results in orchestratorChannel to the Orchestrator
	go func() {
		i := 1
		for result := range s.orchestratorChannel {
			log.Debugf("Sending assessment result (%v) to Orchestrator", result.Id)

			req := &orchestrator.StoreAssessmentResultRequest{
				Result: result,
			}

			err := s.orchestratorStream.Send(req)
			if errors.Is(err, io.EOF) {
				log.Debugf("EOF")
				break
			}
			if err != nil {
				log.Errorf("Error when sending assessment result to Orchestrator: %v", err)
				break
			}

			log.Debugf("Assessment result (%v) sent to Orchestrator", result.Id)

			if i%100 == 0 {
				log.Debugf("orchestratorStream send assessment results currently @ %v", i)
			}
			i++
		}
	}()

	return nil
}

// MetricConfiguration implements MetricConfigurationSource by getting the corresponding metric configuration for the
// default target cloud service
func (s *Service) MetricConfiguration(metric string) (config *assessment.MetricConfiguration, err error) {
	var (
		ok    bool
		cache cachedConfiguration
	)

	// Retrieve our cached entry
	s.confMutex.Lock()
	cache, ok = s.cachedConfigurations[metric]
	s.confMutex.Unlock()

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

		s.confMutex.Lock()
		// Update the metric configuration
		s.cachedConfigurations[metric] = cache
		defer s.confMutex.Unlock()
	}

	return cache.MetricConfiguration, nil
}
