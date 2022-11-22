// Copyright 2022 Fraunhofer AISEC
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

package evaluation

import (
	"context"
	"sync"
	"time"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evaluation"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/service"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

type grpcTarget struct {
	target string
	opts   []grpc.DialOption
}

type mappingResultMetric struct {
	metricName string
	results    []*assessment.AssessmentResult
}

// Service is an implementation of the Clouditor Evaluation service
type Service struct {
	evaluation.UnimplementedEvaluationServer
	scheduler           *gocron.Scheduler
	orchestratorClient  orchestrator.OrchestratorClient
	orchestratorAddress grpcTarget
	authorizer          api.Authorizer

	// evaluation contains lists of controls (value) per target cloud service (key) that are currently evaluated
	evaluation      map[string]*EvaluationScheduler
	evaluationMutex sync.RWMutex

	// Currently, evaluation results are just stored as a map (=in-memory). In the future, we will use a DB.
	results     []*evaluation.EvaluationResult
	resultMutex sync.Mutex
}

func init() {
	log = logrus.WithField("component", "orchestrator")
}

const (
	// DefaultOrchestratorAddress specifies the default gRPC address of the orchestrator service.
	DefaultOrchestratorAddress = "localhost:9090"
	DefaultAssessmentAddress   = "localhost:9090"
)

// ServiceOption is a function-style option to configure the Evaluation Service
type ServiceOption func(*Service)

type EvaluationScheduler struct {
	scheduler           *gocron.Scheduler
	evaluatedControlIDs []string
}

// WithOAuth2Authorizer is an option to use an OAuth 2.0 authorizer
func WithOAuth2Authorizer(config *clientcredentials.Config) service.Option[Service] {
	return func(s *Service) {
		s.SetAuthorizer(api.NewOAuthAuthorizerFromClientCredentials(config))
	}
}

// WithAuthorizer is an option to use a pre-created authorizer
func WithAuthorizer(auth api.Authorizer) service.Option[Service] {
	return func(s *Service) {
		s.SetAuthorizer(auth)
	}
}

// WithOrchestratorAddress is an option to configure the orchestrator service gRPC address.
func WithOrchestratorAddress(address string, opts ...grpc.DialOption) ServiceOption {
	return func(s *Service) {
		s.orchestratorAddress = grpcTarget{
			target: address,
			opts:   opts,
		}
	}
}

// NewService creates a new Evaluation service
func NewService(opts ...ServiceOption) *Service {
	s := Service{
		scheduler: gocron.NewScheduler(time.UTC),
		orchestratorAddress: grpcTarget{
			target: DefaultOrchestratorAddress,
		},
		evaluation: make(map[string]*EvaluationScheduler),
	}

	// Apply service options
	for _, o := range opts {
		o(&s)
	}

	return &s
}

// SetAuthorizer implements UsesAuthorizer
func (svc *Service) SetAuthorizer(auth api.Authorizer) {
	svc.authorizer = auth
}

// Authorizer implements UsesAuthorizer
func (svc *Service) Authorizer() api.Authorizer {
	return svc.authorizer
}

// StartEvaluation is a method implementation of the evaluation interface: It starts the evaluation for a cloud service and a given Control (e.g., EUCS OPS-13.2)
func (s *Service) StartEvaluation(_ context.Context, req *evaluation.StartEvaluationRequest) (resp *evaluation.StartEvaluationResponse, err error) {
	resp = &evaluation.StartEvaluationResponse{}

	err = req.Validate()
	if err != nil {
		resp = &evaluation.StartEvaluationResponse{
			Status:        false,
			StatusMessage: err.Error(),
		}
		return resp, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// Verify that evaluating of this service and control hasn't started already
	// TODO(anatheka): Extend for one schedule per control or do we have to stop it and add with several control IDs?
	s.evaluationMutex.Lock()
	if m := s.evaluation[req.Toe.CloudServiceId]; m != nil && /*slices.Contains(s.evaluation[req.Toe.CloudServiceId].evaluatedControlIDs, req.ControlId) &&*/ m.scheduler != nil && m.scheduler.IsRunning() {
		err = status.Errorf(codes.AlreadyExists, "Service %s is being evaluated with Control %s already.", req.Toe.CloudServiceId, req.ControlId)
		return
	}
	s.evaluationMutex.Unlock()

	log.Info("Starting evaluation ...")
	s.scheduler.TagsUnique()

	log.Infof("Evaluate Cloud Service {%s} every 5 minutes...", req.Toe.CloudServiceId)
	_, err = s.scheduler.
		Every(5).
		Minute().
		Tag(req.Toe.CloudServiceId).
		Do(s.Evaluate, req)

	// Add map entry for target cloud service id or if already exists add new control ID
	s.evaluationMutex.Lock()
	if s.evaluation[req.Toe.CloudServiceId] == nil {
		s.evaluation[req.Toe.CloudServiceId] = new(EvaluationScheduler)
		s.evaluation[req.Toe.CloudServiceId].scheduler = s.scheduler
		s.evaluation[req.Toe.CloudServiceId].evaluatedControlIDs = []string{req.ControlId}
	} else {
		s.evaluation[req.Toe.CloudServiceId].evaluatedControlIDs = append(s.evaluation[req.Toe.CloudServiceId].evaluatedControlIDs, req.ControlId)
	}
	s.evaluationMutex.Unlock()

	s.scheduler.StartAsync()

	return
}

func (s *Service) Evaluate(req *evaluation.StartEvaluationRequest) {
	log.Infof("Started evaluation for Cloud Service '%s',  Catalog ID '%s' and Control '%s'", req.Toe.CloudServiceId, req.Toe.CatalogId, req.ControlId)

	// Establish connection to orchestrator gRPC service
	// TODO (anatheka): Use getOrchestratorClient method or something similiar. Do we have already a method for that?
	conn, err := grpc.Dial(s.orchestratorAddress.target,
		api.DefaultGrpcDialOptions(s.orchestratorAddress.target, s, s.orchestratorAddress.opts...)...)
	if err != nil {
		log.Errorf("could not connect to orchestrator service: %v", err)
	}

	// Create orchestrator client
	s.orchestratorClient = orchestrator.NewOrchestratorClient(conn)

	// Get control and the associated metrics
	metrics, err := s.getMetricFromControl(req, req.CategoryName)
	if err != nil {
		log.Errorf("Could not get metrics for control ID '%s' from Orchestrator: %v", req.ControlId, err)

		// TODO(anatheka): Do we need that?
		// Delete evaluation entry, it is not longer needed
		s.evaluationMutex.Lock()
		delete(s.evaluation, req.Toe.CloudServiceId)
		s.evaluationMutex.Unlock()

		return
	}

	// Get assessment results for the target cloud service
	// TODO(anatheka): The filtered_cloud_service_id option does not work: access denied
	assessmentResults, err := api.ListAllPaginated(&assessment.ListAssessmentResultsRequest{ /*FilteredCloudServiceId: req.Toe.CloudServiceId*/ }, s.orchestratorClient.ListAssessmentResults, func(res *assessment.ListAssessmentResultsResponse) []*assessment.AssessmentResult {
		return res.Results
	})
	if err != nil {
		log.Errorf("Could not get assessment results for Cloud Serivce '%s' from Orchestrator", req.Toe.CloudServiceId)

		// TODO(anatheka): Do we need that?
		// Delete evaluation entry, it is no longer needed if we don't get the assessment results from the orchestrator
		s.evaluationMutex.Lock()
		delete(s.evaluation, req.Toe.CloudServiceId)
		s.evaluationMutex.Unlock()
		return
	}

	// For each metric, find all assessment results for the target cloud service
	var mappingList []mappingResultMetric
	for _, metric := range metrics {
		log.Debugf("Metric with ID '%s'", metric.Id)
		var mapping mappingResultMetric
		var resultList []*assessment.AssessmentResult
		for _, result := range assessmentResults {
			log.Debugf("Assessment result with metric ID '%s'", result.MetricId)

			if result.MetricId == metric.Id {
				resultList = append(resultList, result)
			}
		}
		mapping.metricName = metric.Id
		mapping.results = resultList
		mappingList = append(mappingList, mapping)
	}

	// Do evaluation and find all non-compliant assessment results
	var nonCompliantAssessmentResults []*assessment.AssessmentResult
	status := evaluation.EvaluationResult_PENDING
	for _, item := range mappingList {
		for _, elem := range item.results {
			if !elem.Compliant {
				nonCompliantAssessmentResults = append(nonCompliantAssessmentResults, elem)
				status = evaluation.EvaluationResult_NOT_COMPLIANT
			}
		}
	}

	// Create evaluation result
	// TODO(all): Store in DB
	s.resultMutex.Lock()
	result := &evaluation.EvaluationResult{
		Id:                       uuid.NewString(),
		Timestamp:                timestamppb.Now(),
		Control:                  req.ControlId,
		CloudServiceId:           req.Toe.CloudServiceId,
		TargetOfEvaluation:       req.Toe.String(),
		Status:                   status,
		FailingAssessmentResults: nonCompliantAssessmentResults,
	}
	s.results = append(s.results, result)
	s.resultMutex.Unlock()

	log.Infof("Evaluation result with ID %s stored.", result.Id)
}

// StopEvaluation is a method implementation of the evaluation interface: It starts the evaluation for a cloud service and a given Control (e.g., EUCS OPS-13.2)
func (s *Service) StopEvaluation(_ context.Context, req *evaluation.StopEvaluationRequest) (res *evaluation.StopEvaluationResponse, err error) {
	err = req.Validate()
	if err != nil {
		err = status.Errorf(codes.InvalidArgument, "%v", err)
		return
	}

	// Verify that the service is evaluated currently
	if s.evaluation[req.Toe.CloudServiceId] == nil {
		err = status.Errorf(codes.NotFound, "Evaluation of cloud service %s has not been started yet.", req.Toe.CloudServiceId)
		return
	}

	// Verify if scheduler is running
	if !s.evaluation[req.Toe.CloudServiceId].scheduler.IsRunning() {
		err = status.Errorf(codes.NotFound, "Evaluation of cloud service %s has been stopped already", req.Toe.CloudServiceId)
		return
	}

	// Stop scheduler
	s.evaluationMutex.Lock()
	s.evaluation[req.Toe.CloudServiceId].scheduler.RemoveByTag(req.Toe.CloudServiceId)
	// Delete entry for given Cloud Service ID
	delete(s.evaluation, req.Toe.CloudServiceId)
	s.evaluationMutex.Unlock()

	res = &evaluation.StopEvaluationResponse{}

	return
}

// getMetricFromControl return a list of metrics for the given control ID. For now it is only possible to get the metrics for the lowest control level.
func (s *Service) getMetricFromControl(req *evaluation.StartEvaluationRequest, categoryName string) (metrics []*assessment.Metric, err error) {
	control, err := s.orchestratorClient.GetControl(context.Background(), &orchestrator.GetControlRequest{
		CatalogId:    req.Toe.CatalogId,
		CategoryName: req.CategoryName,
		ControlId:    req.ControlId,
	})

	metrics = control.Metrics

	return
}

func (s *Service) Shutdown() {
	s.scheduler.Stop()
}
