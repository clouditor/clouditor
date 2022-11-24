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
	"fmt"
	"sync"
	"time"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evaluation"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/persistence"
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
	results     map[string]*evaluation.EvaluationResult
	resultMutex sync.Mutex

	storage persistence.Storage
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
	scheduler          *gocron.Scheduler
	evaluatedControlID string
}

// WithStorage is an option to set the storage. If not set, NewService will use inmemory storage.
func WithStorage(storage persistence.Storage) ServiceOption {
	return func(s *Service) {
		s.storage = storage
	}
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
		results:    make(map[string]*evaluation.EvaluationResult),
	}

	// Apply service options
	for _, o := range opts {
		o(&s)
	}

	// // Default to an in-memory storage, if nothing was explicitly set
	// if s.storage == nil {
	// 	s.storage, err = inmemory.NewStorage()
	// 	if err != nil {
	// 		log.Errorf("Could not initialize the storage: %v", err)
	// 	}
	// }

	return &s
}

// SetAuthorizer implements UsesAuthorizer
func (s *Service) SetAuthorizer(auth api.Authorizer) {
	s.authorizer = auth
}

// Authorizer implements UsesAuthorizer
func (s *Service) Authorizer() api.Authorizer {
	return s.authorizer
}

// StartEvaluation is a method implementation of the evaluation interface: It starts the evaluation for a cloud service and a given Control (e.g., EUCS OPS-13.2)
func (s *Service) StartEvaluation(_ context.Context, req *evaluation.StartEvaluationRequest) (resp *evaluation.StartEvaluationResponse, err error) {
	resp = &evaluation.StartEvaluationResponse{}
	var schedulerTag string

	err = req.Validate()
	if err != nil {
		resp = &evaluation.StartEvaluationResponse{
			Status:        false,
			StatusMessage: err.Error(),
		}
		return resp, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// Set scheduler tag
	schedulerTag = createSchedulerTag(req.TargetOfEvaluation.CloudServiceId, req.ControlId)

	// Verify that evaluating of this service and control hasn't started already
	// TODO(anatheka): Extend for one schedule per control or do we have to stop it and add with several control IDs?
	s.evaluationMutex.Lock()
	if m := s.evaluation[schedulerTag]; m != nil && m.scheduler != nil && m.scheduler.IsRunning() {
		err = status.Errorf(codes.AlreadyExists, "Cloud Service '%s' is being evaluated with Control %s already.", req.TargetOfEvaluation.CloudServiceId, req.ControlId)
		log.Error(err)
		return
	}
	s.evaluationMutex.Unlock()

	log.Info("Starting evaluation ...")
	s.scheduler.TagsUnique()

	log.Infof("Evaluate Cloud Service '%s' for Control ID '%s' every 5 minutes...", req.TargetOfEvaluation.CloudServiceId, req.ControlId)
	_, err = s.scheduler.
		Every(5).
		Minute().
		Tag(schedulerTag).
		Do(s.Evaluate, req)
	if err != nil {
		err = fmt.Errorf("evaluation for Cloud Service '%s' and Control ID '%s' cannot be scheduled", req.TargetOfEvaluation.CloudServiceId, req.ControlId)
		log.Error(err)
		err = status.Errorf(codes.Internal, "%s", err)
		log.Error(err)
		return
	}

	// Add map entry for target cloud service id and control id
	s.evaluationMutex.Lock()
	if s.evaluation[schedulerTag] == nil {
		s.evaluation[schedulerTag] = new(EvaluationScheduler)
		s.evaluation[schedulerTag].scheduler = s.scheduler
		s.evaluation[schedulerTag].evaluatedControlID = req.ControlId
	}
	s.evaluationMutex.Unlock()

	s.scheduler.StartAsync()

	return
}

func (s *Service) Evaluate(req *evaluation.StartEvaluationRequest) {
	log.Infof("Started evaluation for Cloud Service '%s',  Catalog ID '%s' and Control '%s'", req.TargetOfEvaluation.CloudServiceId, req.TargetOfEvaluation.CatalogId, req.ControlId)

	var metrics []*assessment.Metric

	// Get orchestrator client
	err := s.initOrchestratorClient()
	if err != nil {
		log.Errorf("could not set orchestrator client: %v", err)
		return
	}

	// Get metrics from control and sub-controls
	metrics, err = s.getMetrics(req.TargetOfEvaluation.CatalogId, req.CategoryName, req.ControlId)
	if err != nil {
		log.Errorf("Could not get metrics from control and sub-controlsfor Cloud Serivce '%s' from Orchestrator: %v", req.TargetOfEvaluation.CloudServiceId, err)
		return
	}

	// Get assessment results for the target cloud service
	assessmentResults, err := api.ListAllPaginated(&assessment.ListAssessmentResultsRequest{FilteredCloudServiceId: req.TargetOfEvaluation.CloudServiceId}, s.orchestratorClient.ListAssessmentResults, func(res *assessment.ListAssessmentResultsResponse) []*assessment.AssessmentResult {
		return res.Results
	})
	if err != nil {
		log.Errorf("Could not get assessment results for Cloud Serivce '%s' from Orchestrator: %v", req.TargetOfEvaluation.CloudServiceId, err)

		// TODO(anatheka): Do we need that? Or do we let it running?
		// Delete evaluation entry, it is no longer needed if we don't get the assessment results from the orchestrator
		s.evaluationMutex.Lock()
		delete(s.evaluation, createSchedulerTag(req.TargetOfEvaluation.CloudServiceId, req.ControlId))
		s.evaluationMutex.Unlock()
		return
	}

	// Get mapping assessment results related to the metric
	mappingList := getMapping(assessmentResults, metrics)

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

	// If no assessment results are available for the metric, the evaluation status is set to compliant
	if status == evaluation.EvaluationResult_PENDING {
		status = evaluation.EvaluationResult_COMPLIANT
	}

	// Create evaluation result
	// TODO(all): Store in DB
	s.resultMutex.Lock()
	result := &evaluation.EvaluationResult{
		Id:                       uuid.NewString(),
		Timestamp:                timestamppb.Now(),
		CategoryName:             req.CategoryName,
		Control:                  req.ControlId,
		TargetOfEvaluation:       req.TargetOfEvaluation,
		Status:                   status,
		FailingAssessmentResults: nonCompliantAssessmentResults,
	}
	s.results[result.Id] = result
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
	if s.evaluation[createSchedulerTag(req.TargetOfEvaluation.CloudServiceId, req.ControlId)] == nil {
		err = fmt.Errorf("evaluation of cloud service %s has not been started yet", req.TargetOfEvaluation.CloudServiceId)
		log.Error(err)
		err = status.Errorf(codes.NotFound, "%s", err)
		return
	}

	// Verify if scheduler is running
	if !s.evaluation[createSchedulerTag(req.TargetOfEvaluation.CloudServiceId, req.ControlId)].scheduler.IsRunning() {
		err = fmt.Errorf("evaluation of cloud service %s has been stopped already", req.TargetOfEvaluation.CloudServiceId)
		log.Error(err)
		err = status.Errorf(codes.NotFound, err.Error())
		return
	}

	// Stop scheduler
	s.evaluationMutex.Lock()
	err = s.evaluation[createSchedulerTag(req.TargetOfEvaluation.CloudServiceId, req.ControlId)].scheduler.RemoveByTag(createSchedulerTag(req.TargetOfEvaluation.CloudServiceId, req.ControlId))
	if err != nil {
		err = fmt.Errorf("error in removing scheduler: %v", err)
		log.Error(err)
		err = status.Errorf(codes.Internal, "error at stopping scheduler")
	}
	// Delete entry for given Cloud Service ID
	delete(s.evaluation, createSchedulerTag(req.TargetOfEvaluation.CloudServiceId, req.ControlId))
	s.evaluationMutex.Unlock()

	res = &evaluation.StopEvaluationResponse{}

	return
}

// ListEvaluationResults is a method implementation of the assessment interface
func (s *Service) ListEvaluationResults(_ context.Context, req *evaluation.ListEvaluationResultsRequest) (res *evaluation.ListEvaluationResultsResponse, err error) {
	res = new(evaluation.ListEvaluationResultsResponse)

	// Paginate the results according to the request
	res.Results, res.NextPageToken, err = service.PaginateMapValues(req, s.results, func(a *evaluation.EvaluationResult, b *evaluation.EvaluationResult) bool {
		return a.Id < b.Id
	}, service.DefaultPaginationOpts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate evaluation results: %v", err)
	}

	return
}

// getMetrics return the metrics from a given controlId
// For now a control has either sub-controls or metrics. If the control has sub-controls, get all metrics from the sub-controls.
// TODO(anatheka): Is it possible that a control has sub-controls and metrics?
func (s *Service) getMetrics(catalogId, categoryName, controlId string) (metrics []*assessment.Metric, err error) {
	var subControlMetrics []*assessment.Metric

	control, err := s.getControl(catalogId, categoryName, controlId)
	if err != nil {
		err = fmt.Errorf("could not get control for control id {%s}: %v", controlId, err)
		log.Error(err)
		return
	}

	// Add metric of control to the metrics list
	metrics = append(metrics, subControlMetrics...)

	// Add sub-control metrics to the metric list if exist
	if len(control.Controls) != 0 {
		// Get the metrics from the next sub-control level
		subControlMetrics, err = s.getMetricsFromSubControls(control)
		if err != nil {
			err = fmt.Errorf("error getting metrics from sub-controls: %v", err)
			log.Error(err)
			return
		}

		metrics = append(metrics, subControlMetrics...)
	}

	return
}

// getMetricsFromSubControls returns a list of metrics from the sub-controls
func (s *Service) getMetricsFromSubControls(control *orchestrator.Control) (metrics []*assessment.Metric, err error) {
	var subcontrol *orchestrator.Control

	for _, control := range control.Controls {
		subcontrol, err = s.getControl(control.CategoryCatalogId, control.CategoryName, control.Id)
		if err != nil {
			return
		}

		metrics = append(metrics, subcontrol.Metrics...)
	}

	return
}

// getControl return the control for the given control ID.
func (s *Service) getControl(catalogId, categoryName, controlId string) (control *orchestrator.Control, err error) {
	control, err = s.orchestratorClient.GetControl(context.Background(), &orchestrator.GetControlRequest{
		CatalogId:    catalogId,
		CategoryName: categoryName,
		ControlId:    controlId,
	})

	return
}

func createSchedulerTag(cloudServiceId, controlId string) string {
	if cloudServiceId == "" || controlId == "" {
		return ""
	}
	return fmt.Sprintf("%s-%s", cloudServiceId, controlId)
}

// TODO(all): Make a generic method for that in folder internal?
// initOrchestratorClient set the orchestrator client
func (s *Service) initOrchestratorClient() error {
	if s.orchestratorClient != nil {
		return nil
	}

	// Establish connection to orchestrator gRPC service
	conn, err := grpc.Dial(s.orchestratorAddress.target,
		api.DefaultGrpcDialOptions(s.orchestratorAddress.target, s, s.orchestratorAddress.opts...)...,
	)
	if err != nil {
		return fmt.Errorf("could not connect to orchestrator service: %w", err)
	}

	s.orchestratorClient = orchestrator.NewOrchestratorClient(conn)

	return nil
}

// getMapping returns a mapping of the assessment results to the metric
func getMapping(results []*assessment.AssessmentResult, metrics []*assessment.Metric) (mappingList []*mappingResultMetric) {
	// For each metric, find all assessment results related to it
	for _, metric := range metrics {
		var mapping *mappingResultMetric
		resultList := []*assessment.AssessmentResult{}
		// Find all assessment results for the metric ID
		for _, result := range results {
			if result.MetricId == metric.Id {
				resultList = append(resultList, result)
			}
		}

		// Add assessment results and metric ID
		mapping = &mappingResultMetric{
			metricName: metric.Id,
			results:    resultList,
		}

		// Store mapping in list
		mappingList = append(mappingList, mapping)
	}

	return
}

func (s *Service) Shutdown() {
	s.scheduler.Stop()
}
