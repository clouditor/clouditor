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
	"clouditor.io/clouditor/service"
	service_orchestrator "clouditor.io/clouditor/service/orchestrator"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/clientcredentials"
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

type grpcTarget struct {
	target string
	opts   []grpc.DialOption
}

// Service is an implementation of the Clouditor Assessment service. It should not be used directly,
// but rather the NewService constructor should be used. It implements the AssessmentServer interface.
type Service struct {
	// Embedded for FWD compatibility
	assessment.UnimplementedAssessmentServer

	// evidenceStoreStream sends evidences to the Evidence Store
	evidenceStoreStreams *api.StreamsOf[evidence.EvidenceStore_StoreEvidencesClient, *evidence.StoreEvidenceRequest]
	evidenceStoreAddress grpcTarget

	// orchestratorStream sends assessment results to the Orchestrator
	orchestratorStreams *api.StreamsOf[orchestrator.Orchestrator_StoreAssessmentResultsClient, *orchestrator.StoreAssessmentResultRequest]
	orchestratorClient  orchestrator.OrchestratorClient
	orchestratorAddress grpcTarget
	metricEventStream   orchestrator.Orchestrator_SubscribeMetricChangeEventsClient

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

	// pe contains the actual policy evaluation engine we use
	pe policies.PolicyEval
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
func WithEvidenceStoreAddress(address string, opts ...grpc.DialOption) ServiceOption {
	return func(s *Service) {
		if address == "" {
			address = DefaultEvidenceStoreAddress
		}

		s.evidenceStoreAddress = grpcTarget{
			target: address,
			opts:   opts,
		}
	}
}

// WithOrchestratorAddress is an option to configure the orchestrator gRPC address.
func WithOrchestratorAddress(address string, opts ...grpc.DialOption) ServiceOption {
	return func(s *Service) {
		if address == "" {
			address = DefaultOrchestratorAddress
		}

		s.orchestratorAddress = grpcTarget{
			target: address,
			opts:   opts,
		}
	}
}

// WithOAuth2Authorizer is an option to use an OAuth 2.0 authorizer
func WithOAuth2Authorizer(config *clientcredentials.Config) ServiceOption {
	return func(s *Service) {
		s.SetAuthorizer(api.NewOAuthAuthorizerFromClientCredentials(config))
	}
}

// NewService creates a new assessment service with default values.
func NewService(opts ...ServiceOption) *Service {
	s := &Service{
		results:              make(map[string]*assessment.AssessmentResult),
		evidenceStoreStreams: api.NewStreamsOf(api.WithLogger[evidence.EvidenceStore_StoreEvidencesClient, *evidence.StoreEvidenceRequest](log)),
		orchestratorStreams:  api.NewStreamsOf(api.WithLogger[orchestrator.Orchestrator_StoreAssessmentResultsClient, *orchestrator.StoreAssessmentResultRequest](log)),
		cachedConfigurations: make(map[string]cachedConfiguration),
	}

	// Apply any options
	for _, o := range opts {
		o(s)
	}

	// Set to default Evidence Store
	if s.evidenceStoreAddress.target == "" {
		s.evidenceStoreAddress = grpcTarget{
			target: DefaultEvidenceStoreAddress,
		}
	}

	// Set to default Orchestrator
	if s.orchestratorAddress.target == "" {
		s.orchestratorAddress = grpcTarget{
			target: DefaultOrchestratorAddress,
		}
	}

	// Initialize the policy evaluator after storage is set
	s.pe = policies.NewRegoEval()

	return s
}

// SetAuthorizer implements UsesAuthorizer
func (svc *Service) SetAuthorizer(auth api.Authorizer) {
	svc.authorizer = auth
}

// Authorizer implements UsesAuthorizer
func (svc *Service) Authorizer() api.Authorizer {
	return svc.authorizer
}

// AssessEvidence is a method implementation of the assessment interface: It assesses a single evidence
func (svc *Service) AssessEvidence(_ context.Context, req *assessment.AssessEvidenceRequest) (res *assessment.AssessEvidenceResponse, err error) {

	// Validate evidence
	resourceId, err := req.Evidence.Validate()
	if err != nil {
		newError := fmt.Errorf("invalid evidence: %w", err)
		log.Error(newError)
		svc.informHooks(nil, newError)

		res = &assessment.AssessEvidenceResponse{
			Status:        assessment.AssessEvidenceResponse_FAILED,
			StatusMessage: newError.Error(),
		}

		return res, status.Errorf(codes.InvalidArgument, "%v", newError)
	}

	// Assess evidence
	err = svc.handleEvidence(req.Evidence, resourceId)
	if err != nil {
		res = &assessment.AssessEvidenceResponse{
			Status:        assessment.AssessEvidenceResponse_FAILED,
			StatusMessage: err.Error(),
		}

		newError := errors.New("error while handling evidence")
		log.Errorf("%v: %v", newError, err)

		return res, status.Errorf(codes.Internal, "%v", newError)
	}

	// Create response
	res = &assessment.AssessEvidenceResponse{
		Status: assessment.AssessEvidenceResponse_ASSESSED,
	}

	return res, nil
}

// AssessEvidences is a method implementation of the assessment interface: It assesses multiple evidences (stream) and responds with a stream.
func (svc *Service) AssessEvidences(stream assessment.Assessment_AssessEvidencesServer) (err error) {
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
		res, err = svc.AssessEvidence(context.Background(), assessEvidencesReq)
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
func (svc *Service) handleEvidence(ev *evidence.Evidence, resourceId string) (err error) {
	log.Debugf("Evaluating evidence %s (%s) collected by %s at %v", ev.Id, resourceId, ev.ToolId, ev.Timestamp)
	log.Tracef("Evidence: %+v", ev)

	evaluations, err := svc.pe.Eval(ev, svc)
	if err != nil {
		newError := fmt.Errorf("could not evaluate evidence: %w", err)
		log.Error(newError)

		go svc.informHooks(nil, newError)

		return newError
	}

	// Get Evidence Store stream
	channelEvidenceStore, err := svc.evidenceStoreStreams.GetStream(svc.evidenceStoreAddress.target, "Evidence Store", svc.initEvidenceStoreStream, svc.evidenceStoreAddress.opts...)
	if err != nil {
		err = fmt.Errorf("could not get stream to evidence store (%s): %w", svc.evidenceStoreAddress.target, err)
		log.Error(err)

		go svc.informHooks(nil, err)

		return err
	}

	// Get Orchestrator stream
	channelOrchestrator, err := svc.orchestratorStreams.GetStream(svc.orchestratorAddress.target, "Orchestrator", svc.initOrchestratorStream, svc.orchestratorAddress.opts...)
	if err != nil {
		err = fmt.Errorf("could not get stream to orchestrator (%s): %w", svc.orchestratorAddress.target, err)
		log.Error(err)

		go svc.informHooks(nil, err)

		return err
	}

	// Send evidence in evidenceStoreChannel
	channelEvidenceStore.Send(&evidence.StoreEvidenceRequest{Evidence: ev})

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

		svc.resultMutex.Lock()
		// Just a little hack to quickly enable multiple results per resource
		svc.results[fmt.Sprintf("%s-%d", resourceId, i)] = result
		svc.resultMutex.Unlock()

		// Inform hooks about new assessment result
		go svc.informHooks(result, nil)

		// Send assessment result in orchestratorChannel
		channelOrchestrator.Send(&orchestrator.StoreAssessmentResultRequest{Result: result})
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
func (svc *Service) informHooks(result *assessment.AssessmentResult, err error) {
	svc.hookMutex.RLock()
	hooks := svc.resultHooks
	defer svc.hookMutex.RUnlock()

	// Inform our hook, if we have any
	if len(hooks) > 0 {
		for _, hook := range hooks {
			// We could do hook concurrent again (assuming different hooks don't interfere with each other)
			hook(result, err)
		}
	}
}

// ListAssessmentResults is a method implementation of the assessment interface
func (svc *Service) ListAssessmentResults(_ context.Context, req *assessment.ListAssessmentResultsRequest) (res *assessment.ListAssessmentResultsResponse, err error) {
	res = new(assessment.ListAssessmentResultsResponse)

	// Paginate the results according to the request
	res.Results, res.NextPageToken, err = service.PaginateMapValues(req, svc.results, func(a *assessment.AssessmentResult, b *assessment.AssessmentResult) bool {
		return a.Id < b.Id
	}, service.DefaultPaginationOpts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
	}

	return
}

func (svc *Service) RegisterAssessmentResultHook(assessmentResultsHook func(result *assessment.AssessmentResult, err error)) {
	svc.hookMutex.Lock()
	defer svc.hookMutex.Unlock()
	svc.resultHooks = append(svc.resultHooks, assessmentResultsHook)
}

// initEvidenceStoreStream initializes the stream to the Evidence Store
func (svc *Service) initEvidenceStoreStream(URL string, additionalOpts ...grpc.DialOption) (stream evidence.EvidenceStore_StoreEvidencesClient, err error) {
	log.Infof("Trying to establish a connection to evidence store service @ %v", svc.evidenceStoreAddress.target)

	// Establish connection to evidence store gRPC service
	conn, err := grpc.Dial(URL,
		api.DefaultGrpcDialOptions(URL, svc, additionalOpts...)...,
	)
	if err != nil {
		return nil, fmt.Errorf("could not connect to evidence store service: %w", err)
	}

	evidenceStoreClient := evidence.NewEvidenceStoreClient(conn)
	stream, err = evidenceStoreClient.StoreEvidences(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not set up stream to evidence store for storing evidences: %w", err)
	}

	log.Infof("Connected to Evidence Store")

	return
}

// initOrchestratorStream initializes the stream to the Orchestrator
func (svc *Service) initOrchestratorStream(URL string, additionalOpts ...grpc.DialOption) (stream orchestrator.Orchestrator_StoreAssessmentResultsClient, err error) {
	log.Infof("Trying to establish a connection to orchestrator This will disappear when I changeOrcservice @ %v", svc.orchestratorAddress.target)

	// Establish connection to orchestrator gRPC service
	err = svc.initOrchestratorClient()
	if err != nil {
		return nil, fmt.Errorf("could not set orchestrator client")
	}

	stream, err = svc.orchestratorClient.StoreAssessmentResults(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not set up stream to orchestrator for storing assessment results: %w", err)
	}

	log.Infof("Connected to Orchestrator")

	// TODO(oxisto): We should rewrite our generic StreamsOf to deal with incoming messages
	svc.metricEventStream, err = svc.orchestratorClient.SubscribeMetricChangeEvents(context.Background(), &orchestrator.SubscribeMetricChangeEventRequest{})
	if err != nil {
		return nil, fmt.Errorf("could not set up stream for listening to metric change events: %w", err)
	}

	go svc.recvEventsLoop()

	return
}

// Metrics implements MetricsSource by retrieving the metric list from the orchestrator.
func (svc *Service) Metrics() (metrics []*assessment.Metric, err error) {
	var res *orchestrator.ListMetricsResponse

	err = svc.initOrchestratorClient()
	if err != nil {
		return nil, fmt.Errorf("could not set orchestrator client")
	}

	res, err = svc.orchestratorClient.ListMetrics(context.Background(), &orchestrator.ListMetricsRequest{})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve metric list from orchestrator: %w", err)
	}

	return res.Metrics, nil
}

// MetricImplementation implements MetricsSource by retrieving the metric implementation
// from the orchestrator.
func (svc *Service) MetricImplementation(lang assessment.MetricImplementation_Language, metric string) (impl *assessment.MetricImplementation, err error) {
	// For now, the orchestrator only supports the Rego language.
	if lang != assessment.MetricImplementation_REGO {
		return nil, errors.New("unsupported language")
	}

	err = svc.initOrchestratorClient()
	if err != nil {
		return nil, fmt.Errorf("could not set orchestrator client")
	}

	// Retrieve it from the orchestrator
	impl, err = svc.orchestratorClient.GetMetricImplementation(context.Background(), &orchestrator.GetMetricImplementationRequest{
		MetricId: metric,
	})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve metric implementation for %s from orchestrator: %w", metric, err)
	}

	return
}

// MetricConfiguration implements MetricsSource by getting the corresponding metric configuration for the
// default target cloud service
func (svc *Service) MetricConfiguration(metric string) (config *assessment.MetricConfiguration, err error) {
	var (
		ok    bool
		cache cachedConfiguration
	)

	// Retrieve our cached entry
	svc.confMutex.Lock()
	cache, ok = svc.cachedConfigurations[metric]
	svc.confMutex.Unlock()

	err = svc.initOrchestratorClient()
	if err != nil {
		return nil, fmt.Errorf("could not set orchestrator client")
	}

	// Check if entry is not there or is expired
	if !ok || cache.cachedAt.After(time.Now().Add(EvictionTime)) {
		config, err = svc.orchestratorClient.GetMetricConfiguration(context.Background(), &orchestrator.GetMetricConfigurationRequest{
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

		svc.confMutex.Lock()
		// Update the metric configuration
		svc.cachedConfigurations[metric] = cache
		defer svc.confMutex.Unlock()
	}

	return cache.MetricConfiguration, nil
}

// recvEventsLoop continuously tries to receive events on the metricEventStream
func (svc *Service) recvEventsLoop() {
	for {
		var (
			event *orchestrator.MetricChangeEvent
			err   error
		)
		event, err = svc.metricEventStream.Recv()

		if errors.Is(err, io.EOF) {
			log.Debugf("no more responses from orchestrator stream: EOF")
			break
		}

		if err != nil {
			log.Errorf("error receiving response from orchestrator stream: %v", err)
			break
		}

		_ = svc.pe.HandleMetricEvent(event)
	}
}

// initOrchestratorClient set the orchestrator client
func (svc *Service) initOrchestratorClient() error {
	if svc.orchestratorClient != nil {
		log.Debug("Orchestrator client is already initialized.")
		return nil
	}

	// Establish connection to orchestrator gRPC service
	conn, err := grpc.Dial(svc.orchestratorAddress.target,
		api.DefaultGrpcDialOptions(svc.orchestratorAddress.target, svc, svc.orchestratorAddress.opts...)...,
	)
	if err != nil {
		return fmt.Errorf("could not connect to orchestrator service: %w", err)
	}

	svc.orchestratorClient = orchestrator.NewOrchestratorClient(conn)

	return nil
}
