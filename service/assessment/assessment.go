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
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/logging"
	"clouditor.io/clouditor/policies"
	"clouditor.io/clouditor/service"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	// isEvidenceStoreDisabled specifies if evidences shall be discarded (when true).
	isEvidenceStoreDisabled bool
	// evidenceStoreStream sends evidences to the Evidence Store
	evidenceStoreStreams *api.StreamsOf[evidence.EvidenceStore_StoreEvidencesClient, *evidence.StoreEvidenceRequest]
	evidenceStore        *api.RPCConnection[evidence.EvidenceStoreClient]

	// orchestratorStream sends assessment results to the Orchestrator
	orchestratorStreams *api.StreamsOf[orchestrator.Orchestrator_StoreAssessmentResultsClient, *orchestrator.StoreAssessmentResultRequest]
	orchestrator        *api.RPCConnection[orchestrator.OrchestratorClient]
	metricEventStream   orchestrator.Orchestrator_SubscribeMetricChangeEventsClient

	// resultHooks is a list of hook functions that can be used if one wants to be
	// informed about each assessment result
	resultHooks []assessment.ResultHookFunc
	// hookMutex is used for (un)locking result hook calls
	hookMutex sync.RWMutex

	// cachedConfigurations holds cached metric configurations for faster access with key being the corresponding
	// metric name
	cachedConfigurations map[string]cachedConfiguration
	// TODO(oxisto): combine with hookMutex and replace with a generic version of a mutex'd map
	confMutex sync.Mutex

	authz service.AuthorizationStrategy

	// pe contains the actual policy evaluation engine we use
	pe policies.PolicyEval

	// evalPkg specifies the package used for the evaluation engine
	evalPkg string
}

const (
	// DefaultEvidenceStoreAddress specifies the default gRPC address of the evidence store.
	DefaultEvidenceStoreAddress = "localhost:9090"

	// DefaultOrchestratorAddress specifies the default gRPC address of the orchestrator.
	DefaultOrchestratorAddress = "localhost:9090"
)

// WithoutEvidenceStore is a service option to discard evidences and don't send them to an evidence store
func WithoutEvidenceStore() service.Option[Service] {
	return func(svc *Service) {
		svc.isEvidenceStoreDisabled = true
	}
}

// WithEvidenceStoreAddress is an option to configure the evidence store gRPC address.
func WithEvidenceStoreAddress(address string, opts ...grpc.DialOption) service.Option[Service] {
	return func(svc *Service) {
		svc.evidenceStore.Target = address
		svc.evidenceStore.Opts = opts
	}
}

// WithOrchestratorAddress is an option to configure the orchestrator gRPC address.
func WithOrchestratorAddress(target string, opts ...grpc.DialOption) service.Option[Service] {
	return func(svc *Service) {
		svc.orchestrator.Target = target
		svc.orchestrator.Opts = opts
	}
}

// WithOAuth2Authorizer is an option to use an OAuth 2.0 authorizer
func WithOAuth2Authorizer(config *clientcredentials.Config) service.Option[Service] {
	return func(s *Service) {
		auth := api.NewOAuthAuthorizerFromClientCredentials(config)
		s.evidenceStore.SetAuthorizer(auth)
		s.orchestrator.SetAuthorizer(auth)
	}
}

// WithAuthorizer is an option to use a pre-created authorizer
func WithAuthorizer(auth api.Authorizer) service.Option[Service] {
	return func(s *Service) {
		s.evidenceStore.SetAuthorizer(auth)
		s.orchestrator.SetAuthorizer(auth)
	}
}

// WithRegoPackageName is an option to configure the Rego package name
func WithRegoPackageName(pkg string) service.Option[Service] {
	return func(s *Service) {
		s.evalPkg = pkg
	}
}

// WithAuthorizationStrategy is an option that configures an authorization strategy.
func WithAuthorizationStrategy(authz service.AuthorizationStrategy) service.Option[Service] {
	return func(svc *Service) {
		svc.authz = authz
	}
}

// NewService creates a new assessment service with default values.
func NewService(opts ...service.Option[Service]) *Service {
	svc := &Service{
		evidenceStoreStreams: api.NewStreamsOf(api.WithLogger[evidence.EvidenceStore_StoreEvidencesClient, *evidence.StoreEvidenceRequest](log)),
		orchestratorStreams:  api.NewStreamsOf(api.WithLogger[orchestrator.Orchestrator_StoreAssessmentResultsClient, *orchestrator.StoreAssessmentResultRequest](log)),
		cachedConfigurations: make(map[string]cachedConfiguration),
		evidenceStore:        api.NewRPCConnection(DefaultEvidenceStoreAddress, evidence.NewEvidenceStoreClient),
		orchestrator:         api.NewRPCConnection(DefaultOrchestratorAddress, orchestrator.NewOrchestratorClient),
	}

	// Apply any options
	for _, o := range opts {
		o(svc)
	}

	// Set to default Rego package
	if svc.evalPkg == "" {
		svc.evalPkg = policies.DefaultRegoPackage
	}

	// Initialize the policy evaluator after options are set
	svc.pe = policies.NewRegoEval(policies.WithPackageName(svc.evalPkg))

	// Default to an allow-all authorization strategy
	if svc.authz == nil {
		svc.authz = &service.AuthorizationStrategyAllowAll{}
	}

	return svc
}

// AssessEvidence is a method implementation of the assessment interface: It assesses a single evidence
func (svc *Service) AssessEvidence(ctx context.Context, req *assessment.AssessEvidenceRequest) (resp *assessment.AssessEvidenceResponse, err error) {
	resp = &assessment.AssessEvidenceResponse{}

	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	// Validate evidence
	resourceId, err := req.Evidence.ValidateWithResource()
	if err != nil {
		err = fmt.Errorf("invalid evidence: %w", err)
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// Check if cloud_service_id in the service is within allowed or one can access *all* the cloud services
	if !svc.authz.CheckAccess(ctx, service.AccessUpdate, req) {
		return nil, service.ErrPermissionDenied
	}

	// Assess evidence
	_, err = svc.handleEvidence(req.Evidence, resourceId)
	if err != nil {
		err = fmt.Errorf("error while handling evidence: %v", err)
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	logging.LogRequest(log, logrus.DebugLevel, logging.Assess, req)

	return resp, nil
}

// AssessEvidences is a method implementation of the assessment interface: It assesses multiple evidences (stream) and responds with a stream.
func (svc *Service) AssessEvidences(stream assessment.Assessment_AssessEvidencesServer) (err error) {
	var (
		req *assessment.AssessEvidenceRequest
		res *assessment.AssessEvidencesResponse
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
		_, err = svc.AssessEvidence(stream.Context(), assessEvidencesReq)
		if err != nil {
			// Create response message. The AssessEvidence method does not need that message, so we have to create it here for the stream response.
			res = &assessment.AssessEvidencesResponse{
				Status:        assessment.AssessEvidencesResponse_FAILED,
				StatusMessage: err.Error(),
			}
		} else {
			res = &assessment.AssessEvidencesResponse{
				Status: assessment.AssessEvidencesResponse_ASSESSED,
			}
		}

		// Send response back to the client
		err = stream.Send(res)

		// Check for send errors
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			err = fmt.Errorf("cannot send response to the client: %w", err)
			log.Error(err)
			return status.Errorf(codes.Unknown, "%v", err)
		}
	}
}

// handleEvidence is the helper method for the actual assessment used by AssessEvidence and AssessEvidences
func (svc *Service) handleEvidence(ev *evidence.Evidence, resourceId string) (results []*assessment.AssessmentResult, err error) {
	var types []string

	log.Debugf("Evaluating evidence %s (%s) collected by %s at %s", ev.Id, resourceId, ev.ToolId, ev.Timestamp.AsTime())
	log.Tracef("Evidence: %+v", ev)

	evaluations, err := svc.pe.Eval(ev, svc)
	if err != nil {
		newError := fmt.Errorf("could not evaluate evidence: %w", err)

		go svc.informHooks(nil, newError)

		return nil, newError
	}

	// Send evidence via Evidence Store stream if sending evidences is not disabled
	if !svc.isEvidenceStoreDisabled {
		// Get Evidence Store stream
		channelEvidenceStore, err := svc.evidenceStoreStreams.GetStream(svc.evidenceStore.Target, "Evidence Store", svc.initEvidenceStoreStream, svc.evidenceStore.Opts...)
		if err != nil {
			err = fmt.Errorf("could not get stream to evidence store (%s): %w", svc.evidenceStore.Target, err)

			go svc.informHooks(nil, err)

			return nil, err
		}
		channelEvidenceStore.Send(&evidence.StoreEvidenceRequest{Evidence: ev})
	}

	// Get Orchestrator stream
	channelOrchestrator, err := svc.orchestratorStreams.GetStream(svc.orchestrator.Target, "Orchestrator", svc.initOrchestratorStream, svc.orchestrator.Opts...)
	if err != nil {
		err = fmt.Errorf("could not get stream to orchestrator (%s): %w", svc.orchestrator.Target, err)

		go svc.informHooks(nil, err)

		return nil, err
	}

	for _, data := range evaluations {
		// That there is an empty (nil) evaluation should be caught beforehand, but you never know.
		if data == nil {
			log.Errorf("One empty policy evaluation detected for evidence '%s'. That should not happen.",
				ev.GetId())
			continue
		}
		metricID := data.MetricID

		log.Debugf("Evaluated evidence %v with metric '%v' as %v", ev.Id, metricID, data.Compliant)

		types, err = ev.ResourceTypes()
		if err != nil {
			return nil, fmt.Errorf("could not extract resource types from evidence: %w", err)
		}

		result := &assessment.AssessmentResult{
			Id:                    uuid.NewString(),
			Timestamp:             timestamppb.Now(),
			CloudServiceId:        ev.GetCloudServiceId(),
			MetricId:              metricID,
			MetricConfiguration:   data.Config,
			Compliant:             data.Compliant,
			EvidenceId:            ev.GetId(),
			ResourceId:            resourceId,
			ResourceTypes:         types,
			NonComplianceComments: "No comments so far",
		}

		// Inform hooks about new assessment result
		go svc.informHooks(result, nil)

		// Send assessment result in orchestratorChannel
		channelOrchestrator.Send(&orchestrator.StoreAssessmentResultRequest{Result: result})

		results = append(results, result)
	}

	return results, nil
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

func (svc *Service) RegisterAssessmentResultHook(assessmentResultsHook func(result *assessment.AssessmentResult, err error)) {
	svc.hookMutex.Lock()
	defer svc.hookMutex.Unlock()
	svc.resultHooks = append(svc.resultHooks, assessmentResultsHook)
}

// initEvidenceStoreStream initializes the stream to the Evidence Store
func (svc *Service) initEvidenceStoreStream(target string, _ ...grpc.DialOption) (stream evidence.EvidenceStore_StoreEvidencesClient, err error) {
	log.Infof("Trying to establish a stream to evidence store service @ %v", target)

	// Make sure, that we re-connect
	svc.evidenceStore.ForceReconnect()

	stream, err = svc.evidenceStore.Client.StoreEvidences(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not set up stream to evidence store for storing evidences: %w", err)
	}

	log.Infof("Connected to Evidence Store")

	return
}

// initOrchestratorStream initializes the stream to the Orchestrator
func (svc *Service) initOrchestratorStream(target string, _ ...grpc.DialOption) (stream orchestrator.Orchestrator_StoreAssessmentResultsClient, err error) {
	log.Infof("Trying to establish a stream to orchestrator service @ %v", target)

	// Make sure, that we re-connect
	svc.orchestrator.ForceReconnect()

	stream, err = svc.orchestrator.Client.StoreAssessmentResults(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not set up stream to orchestrator for storing assessment results: %w", err)
	}

	log.Infof("Stream to StoreAssessmentResults established")

	// TODO(oxisto): We should rewrite our generic StreamsOf to deal with incoming messages
	svc.metricEventStream, err = svc.orchestrator.Client.SubscribeMetricChangeEvents(context.Background(), &orchestrator.SubscribeMetricChangeEventRequest{})
	if err != nil {
		return nil, fmt.Errorf("could not set up stream for listening to metric change events: %w", err)
	}

	log.Infof("Stream to SubscribeMetricChangeEvents established")

	go svc.recvEventsLoop()

	return
}

// Metrics implements MetricsSource by retrieving the metric list from the orchestrator.
func (svc *Service) Metrics() (metrics []*assessment.Metric, err error) {
	var res *orchestrator.ListMetricsResponse

	res, err = svc.orchestrator.Client.ListMetrics(context.Background(), &orchestrator.ListMetricsRequest{})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve metric list from orchestrator: %w", err)
	}

	return res.Metrics, nil
}

// MetricImplementation implements MetricsSource by retrieving the metric implementation
// from the orchestrator.
func (svc *Service) MetricImplementation(lang assessment.MetricImplementation_Language, metric string) (impl *assessment.MetricImplementation, err error) {
	// For now, the orchestrator only supports the Rego language.
	if lang != assessment.MetricImplementation_LANGUAGE_REGO {
		return nil, errors.New("unsupported language")
	}

	// Retrieve it from the orchestrator
	impl, err = svc.orchestrator.Client.GetMetricImplementation(context.Background(), &orchestrator.GetMetricImplementationRequest{
		MetricId: metric,
	})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve metric implementation for %s from orchestrator: %w", metric, err)
	}

	return
}

// MetricConfiguration implements MetricsSource by getting the corresponding metric configuration for the
// default target cloud service
func (svc *Service) MetricConfiguration(cloudServiceID, metricID string) (config *assessment.MetricConfiguration, err error) {
	var (
		ok    bool
		cache cachedConfiguration
		key   string
	)

	// Calculate the cache key
	key = fmt.Sprintf("%s-%s", cloudServiceID, metricID)

	// Retrieve our cached entry
	svc.confMutex.Lock()
	cache, ok = svc.cachedConfigurations[key]
	svc.confMutex.Unlock()

	// Check if entry is not there or is expired
	if !ok || cache.cachedAt.After(time.Now().Add(EvictionTime)) {
		config, err = svc.orchestrator.Client.GetMetricConfiguration(context.Background(), &orchestrator.GetMetricConfigurationRequest{
			CloudServiceId: cloudServiceID,
			MetricId:       metricID,
		})

		if err != nil {
			return nil, fmt.Errorf("could not retrieve metric configuration for %s: %w", metricID, err)
		}

		cache = cachedConfiguration{
			cachedAt:            time.Now(),
			MetricConfiguration: config,
		}

		svc.confMutex.Lock()
		// Update the metric configuration
		svc.cachedConfigurations[key] = cache
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

		svc.handleMetricEvent(event)
	}
}

func (svc *Service) handleMetricEvent(event *orchestrator.MetricChangeEvent) {
	var key string

	// In case the configuration has changed, we need to clear our configuration cache. Otherwise the policy evaluation
	// will clear the Rego cache, but still refer to the old metric configuration (until it expires). Handle metric event in our policy
	// evaluation
	if event.GetType() == orchestrator.MetricChangeEvent_TYPE_CONFIG_CHANGED {
		// Evict the metric configuration from cache
		svc.confMutex.Lock()

		// Calculate the cache key
		key = fmt.Sprintf("%s-%s", event.CloudServiceId, event.MetricId)

		delete(svc.cachedConfigurations, key)
		svc.confMutex.Unlock()
	}

	// Forward the event to the policy evaluator
	_ = svc.pe.HandleMetricEvent(event)
}
