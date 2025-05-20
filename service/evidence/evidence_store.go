// Copyright 2016-2023 Fraunhofer AISEC
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

package evidence

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"sync"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/internal/logging"
	"clouditor.io/clouditor/v2/launcher"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/persistence/inmemory"
	"clouditor.io/clouditor/v2/server"
	"clouditor.io/clouditor/v2/service"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var log *logrus.Entry

// DefaultServiceSpec returns a [launcher.ServiceSpec] for this [Service] with all necessary options retrieved from the
// config system.
func DefaultServiceSpec() launcher.ServiceSpec {
	return launcher.NewServiceSpec(
		NewService,
		WithStorage,
		func(svc *Service) ([]server.StartGRPCServerOption, error) {
			// It is possible to register hook functions for the evidenceStore.
			//  * The hook functions in evidenceStore are implemented in StoreEvidence(s)

			// evidenceStoreService.RegisterEvidenceHook(func(result *evidence.Evidence, err error) {})

			return nil, nil
		},
		WithOAuth2Authorizer(config.ClientCredentials()),
	)
}

// Service is an implementation of the Clouditor req service (evidenceServer)
type Service struct {
	storage persistence.Storage

	assessmentStreams *api.StreamsOf[assessment.Assessment_AssessEvidencesClient, *assessment.AssessEvidenceRequest]
	assessment        *api.RPCConnection[assessment.AssessmentClient]

	// channel that is used to send evidences from the StoreEvidence method to the worker threat to process the evidence
	channelEvidence chan *evidence.Evidence

	// evidenceHooks is a list of hook functions that can be used if one wants to be
	// informed about each evidence
	evidenceHooks []evidence.EvidenceHookFunc
	// mu is used for (un)locking result hook calls
	mu sync.Mutex

	// authz defines our authorization strategy, e.g., which user can access which target of evaluation and associated
	// resources, such as evidences and assessment results.
	authz service.AuthorizationStrategy

	evidence.UnimplementedEvidenceStoreServer
	evidence.UnimplementedExperimentalResourcesServer
}

func init() {
	log = logrus.WithField("component", "Evidence Store")
}

func WithStorage(storage persistence.Storage) service.Option[*Service] {
	return func(svc *Service) {
		svc.storage = storage
	}
}

// WithOAuth2Authorizer is an option to use an OAuth 2.0 authorizer
func WithOAuth2Authorizer(config *clientcredentials.Config) service.Option[*Service] {
	return func(s *Service) {
		auth := api.NewOAuthAuthorizerFromClientCredentials(config)
		s.assessment.SetAuthorizer(auth)
	}
}

// WithAssessmentAddress is an option to configure the assessment service gRPC address.
func WithAssessmentAddress(target string, opts ...grpc.DialOption) service.Option[*Service] {

	return func(s *Service) {
		log.Infof("Assessment URL is set to %s", target)

		s.assessment.Target = target
		s.assessment.Opts = opts
	}
}

func NewService(opts ...service.Option[*Service]) (svc *Service) {
	var (
		err error
	)
	svc = &Service{
		assessmentStreams: api.NewStreamsOf(api.WithLogger[assessment.Assessment_AssessEvidencesClient, *assessment.AssessEvidenceRequest](log)),
		assessment:        api.NewRPCConnection(config.DefaultAssessmentURL, assessment.NewAssessmentClient),
		channelEvidence:   make(chan *evidence.Evidence, 1000),
	}

	for _, o := range opts {
		o(svc)
	}

	// Default to an allow-all authorization strategy
	if svc.authz == nil {
		svc.authz = &service.AuthorizationStrategyAllowAll{}
	}

	if svc.storage == nil {
		svc.storage, err = inmemory.NewStorage()
		if err != nil {
			log.Errorf("Could not initialize the storage: %v", err)
		}
	}

	return
}

func (svc *Service) Init() {

	// Start a worker thread to process the evidence that is being passed to the StoreEvidence function in order to utilize the fire-and-forget strategy.
	// To do this, we want an channel, that contains the evidences and call another function that processes the evidence.
	go func() {
		for {
			// Wait for a new evidence to be passed to the channel
			evidence := <-svc.channelEvidence

			// Process the evidence
			err := svc.handleEvidence(evidence)
			if err != nil {
				log.Errorf("Error while processing evidence: %v", err)
			}
		}
	}()
}

func (svc *Service) Shutdown() {
	svc.assessmentStreams.CloseAll()
}

// initAssessmentStream initializes the stream that is used to send evidences to the assessment service.
// If configured, it uses the Authorizer of the evidence store service to authenticate requests to the assessment.
func (svc *Service) initAssessmentStream(target string, _ ...grpc.DialOption) (stream assessment.Assessment_AssessEvidencesClient, err error) {
	log.Infof("Trying to establish a connection to assessment service @ %v", target)

	// Make sure, that we re-connect
	svc.assessment.ForceReconnect()

	// Set up the stream and store it in our service struct, so we can access it later to actually
	// send the evidence data
	stream, err = svc.assessment.Client.AssessEvidences(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not set up stream to assessment for assessing evidence: %w", err)
	}

	log.Infof("Stream to AssessEvidences established")

	return
}

// StoreEvidence is a method implementation of the evidenceServer interface: It receives a req and stores it
func (svc *Service) StoreEvidence(ctx context.Context, req *evidence.StoreEvidenceRequest) (res *evidence.StoreEvidenceResponse, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the target of evaluation according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessUpdate, req) {
		return nil, service.ErrPermissionDenied
	}

	// Store evidence
	err = svc.storage.Create(req.Evidence)
	if err != nil && errors.Is(err, persistence.ErrUniqueConstraintFailed) {
		return nil, status.Error(codes.AlreadyExists, persistence.ErrEntryAlreadyExists.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "%v: %v", persistence.ErrDatabase, err)
	}

	// Store Resource
	// Build a resource struct. This will hold the latest sync state of the
	// resource for our storage layer. This is needed to store the resource in our DB.s
	r, err := evidence.ToEvidenceResource(req.Evidence.GetOntologyResource(), req.GetTargetOfEvaluationId(), req.Evidence.GetToolId())
	if err != nil {
		log.Errorf("Could not convert resource: %v", err)
	}

	// Persist the latest state of the resource
	err = svc.storage.Save(&r, "id = ?", r.Id)
	if err != nil {
		log.Errorf("Could not save resource with ID '%s' to storage: %v", r.Id, err)
	}

	go svc.informHooks(ctx, req.Evidence, nil)

	// Send evidence to the channel for processing (fire and forget)
	svc.channelEvidence <- req.Evidence

	res = &evidence.StoreEvidenceResponse{}

	logging.LogRequest(log, logrus.DebugLevel, logging.Store, req)

	return res, nil
}

func (svc *Service) handleEvidence(evidence *evidence.Evidence) error {
	// TODO(anatheka): It must be checked if the evidence changed since the last time and then send to the assessment service. Add in separate PR

	// Get Assessment stream
	channelAssessment, err := svc.assessmentStreams.GetStream(svc.assessment.Target, "Assessment", svc.initAssessmentStream, svc.assessment.Opts...)
	if err != nil {
		err = fmt.Errorf("could not get stream to assessment service (%s): %w", svc.assessment.Target, err)
		log.Error(err)
		return err
	}

	// Send evidence to assessment service
	channelAssessment.Send(&assessment.AssessEvidenceRequest{Evidence: evidence})

	return nil
}

// StoreEvidences is a method implementation of the evidenceServer interface: It receives evidences and stores them
func (svc *Service) StoreEvidences(stream evidence.EvidenceStore_StoreEvidencesServer) (err error) {
	var (
		req *evidence.StoreEvidenceRequest
		res *evidence.StoreEvidencesResponse
	)

	for {
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

		// Call StoreEvidence() for storing a single evidence
		evidenceRequest := &evidence.StoreEvidenceRequest{
			Evidence: req.Evidence,
		}
		_, err = svc.StoreEvidence(stream.Context(), evidenceRequest)
		if err != nil {
			log.Errorf("Error storing evidence: %v", err)
			// Create response message. The StoreEvidence method does not need that message, so we have to create it here for the stream response.
			res = &evidence.StoreEvidencesResponse{
				Status:        evidence.EvidenceStatus_EVIDENCE_STATUS_ERROR,
				StatusMessage: err.Error(),
			}
		} else {
			res = &evidence.StoreEvidencesResponse{
				Status: evidence.EvidenceStatus_EVIDENCE_STATUS_OK,
			}
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

// ListEvidences is a method implementation of the evidenceServer interface: It returns the evidences lying in the storage
func (svc *Service) ListEvidences(ctx context.Context, req *evidence.ListEvidencesRequest) (res *evidence.ListEvidencesResponse, err error) {
	var (
		all     bool
		allowed []string
		query   []string
		args    []any
	)
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Retrieve list of allowed target of evaluation according to our authorization strategy. No need to specify any additional
	// conditions to our storage request, if we are allowed to see all target of evaluations.
	all, allowed = svc.authz.AllowedTargetOfEvaluations(ctx)
	if !all && req.GetFilter().GetTargetOfEvaluationId() != "" && !slices.Contains(allowed, req.GetFilter().GetTargetOfEvaluationId()) {
		return nil, service.ErrPermissionDenied
	}

	res = new(evidence.ListEvidencesResponse)

	// Apply filter options
	if filter := req.GetFilter(); filter != nil {
		if TargetOfEvaluationId := filter.GetTargetOfEvaluationId(); TargetOfEvaluationId != "" {
			query = append(query, "target_of_evaluation_id = ?")
			args = append(args, TargetOfEvaluationId)
		}
		if toolId := filter.GetToolId(); toolId != "" {
			query = append(query, "tool_id = ?")
			args = append(args, toolId)
		}
	}

	// In any case, we need to make sure that we only select evidences of target of evaluations that we have access to
	if !all {
		query = append(query, "target_of_evaluation_id IN ?")
		args = append(args, allowed)
	}

	// Paginate the evidences according to the request
	res.Evidences, res.NextPageToken, err = service.PaginateStorage[*evidence.Evidence](req, svc.storage,
		service.DefaultPaginationOpts, persistence.BuildConds(query, args)...)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
	}

	return
}

// GetEvidence is a method implementation of the evidenceServer interface: It returns a particular evidence in the storage
func (svc *Service) GetEvidence(ctx context.Context, req *evidence.GetEvidenceRequest) (res *evidence.Evidence, err error) {
	var (
		all     bool
		allowed []string
		conds   []any
	)

	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Retrieve list of allowed target of evaluation according to our authorization strategy. No need to specify any additional
	// conditions to our storage request, if we are allowed to see all target of evaluations.
	all, allowed = svc.authz.AllowedTargetOfEvaluations(ctx)
	if !all {
		conds = []any{"id = ? AND target_of_evaluation_id IN ?", req.EvidenceId, allowed}
	} else {
		conds = []any{"id = ?", req.EvidenceId}
	}

	res = new(evidence.Evidence)

	err = svc.storage.Get(res, conds...)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "evidence not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	return
}

func (svc *Service) ListResources(ctx context.Context, req *evidence.ListResourcesRequest) (res *evidence.ListResourcesResponse, err error) {
	var (
		query   []string
		args    []any
		all     bool
		allowed []string
	)

	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Filtering the resources by
	// * target of evaluation ID
	// * resource type
	// * tool ID
	if req.Filter != nil {
		// Check if target_of_evaluation_id in filter is within allowed or one can access *all* the target of evaluations
		if !svc.authz.CheckAccess(ctx, service.AccessRead, req.Filter) {
			return nil, service.ErrPermissionDenied
		}

		if req.Filter.TargetOfEvaluationId != nil {
			query = append(query, "target_of_evaluation_id = ?")
			args = append(args, req.Filter.GetTargetOfEvaluationId())
		}
		if req.Filter.Type != nil {
			query = append(query, "(resource_type LIKE ? OR resource_type LIKE ? OR resource_type LIKE ?)")
			args = append(args, req.Filter.GetType()+",%", "%,"+req.Filter.GetType()+",%", "%,"+req.Filter.GetType())
		}
		if req.Filter.ToolId != nil {
			query = append(query, "tool_id = ?")
			args = append(args, req.Filter.GetToolId())
		}
	}

	// We need to further restrict our query according to the target of evaluation we are allowed to "see".
	//
	// TODO(oxisto): This is suboptimal, since we are now calling AllowedTargetOfEvaluations twice. Once here
	//  and once above in CheckAccess.
	all, allowed = svc.authz.AllowedTargetOfEvaluations(ctx)
	if !all {
		query = append(query, "target_of_evaluation_id IN ?")
		args = append(args, allowed)
	}

	res = new(evidence.ListResourcesResponse)

	// Join query with AND and prepend the query
	args = append([]any{strings.Join(query, " AND ")}, args...)

	res.Results, res.NextPageToken, err = service.PaginateStorage[*evidence.Resource](req, svc.storage, service.DefaultPaginationOpts, args...)

	return
}

func (svc *Service) RegisterEvidenceHook(evidenceHook evidence.EvidenceHookFunc) {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	svc.evidenceHooks = append(svc.evidenceHooks, evidenceHook)
}

func (svc *Service) informHooks(ctx context.Context, result *evidence.Evidence, err error) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	// Inform our hook, if we have any
	if svc.evidenceHooks != nil {
		for _, hook := range svc.evidenceHooks {
			// TODO(all): We could do hook concurrent again (assuming different hooks don't interfere with each other)
			hook(ctx, result, err)
		}
	}
}
