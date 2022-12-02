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
	"errors"
	"fmt"
	"strings"
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

type Evaluator struct {
	categoryName string
	controlId    string
	// firstLevel describes if the controlId is a first level controlId (e.g., OPS13) or a second level controlId (e.g., OPS-13.1)
	firstLevel bool
}

// Service is an implementation of the Clouditor Evaluation service
type Service struct {
	evaluation.UnimplementedEvaluationServer
	orchestratorClient  orchestrator.OrchestratorClient
	orchestratorAddress grpcTarget
	authorizer          api.Authorizer
	scheduler           *gocron.Scheduler

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
		orchestratorAddress: grpcTarget{
			target: DefaultOrchestratorAddress,
		},
		results:   make(map[string]*evaluation.EvaluationResult),
		scheduler: gocron.NewScheduler(time.UTC),
	}

	// Apply service options
	for _, o := range opts {
		o(&s)
	}

	// TODO(lebogg): Implement persistence storage
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

// StartEvaluation is a method implementation of the evaluation interface: It starts the evaluation for a cloud service and a given Target of Evaluation to be evaluated (e.g., EUCS OPS-13.2)
func (s *Service) StartEvaluation(_ context.Context, req *evaluation.StartEvaluationRequest) (resp *evaluation.StartEvaluationResponse, err error) {
	var (
		schedulerTag string
		evaluator    []*Evaluator
		control      *orchestrator.Control
	)

	// Validate request
	err = req.Validate()
	if err != nil {
		resp = &evaluation.StartEvaluationResponse{
			Status:        false,
			StatusMessage: err.Error(),
		}
		return resp, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// Verify that evaluation of this service and control hasn't started already
	for _, control := range req.EvalControl {
		schedulerTag = createSchedulerTag(req.TargetOfEvaluation.CloudServiceId, control.ControlId)

		_, err := s.scheduler.FindJobsByTag(schedulerTag)
		if err == nil {
			err = status.Errorf(codes.AlreadyExists, "Evaluation for Cloud Service ID '%s' and Control ID %s started already.", req.TargetOfEvaluation.CloudServiceId, control.ControlId)
			log.Error(err)
			return nil, err
		}
	}

	// Get orchestrator client
	err = s.initOrchestratorClient()
	if err != nil {
		log.Errorf("could not set orchestrator client: %v", err)
		return
	}

	log.Info("Starting evaluation ...")

	// Get all necessary information to start the scheduler for the given control IDs
	for _, elem := range req.EvalControl {
		// Get control and check if the control has further subcontrols
		control, err = s.getControl(req.TargetOfEvaluation.CatalogId, elem.CategoryName, elem.ControlId)
		if err != nil {
			err = fmt.Errorf("could not get control for control id {%s}: %v", elem.ControlId, err)
			log.Error(err)
			return
		}

		// TODO(all): Refactor the following to ifs
		// Store current control if it is a second level control
		if len(control.Controls) == 0 {
			eval := &Evaluator{
				categoryName: control.CategoryName,
				controlId:    control.Id,
				firstLevel:   false,
			}
			evaluator = append(evaluator, eval)
		}

		// Check if the control is a first level control and has further sub-controls
		if len(control.Controls) != 0 { // Control from an upper control level, we have to get the lower controls
			// Store current control, it is a first level control
			eval := &Evaluator{
				categoryName: control.CategoryName,
				controlId:    control.Id,
				firstLevel:   true,
			}
			evaluator = append(evaluator, eval)

			// Store the sub-level controls
			for _, elem := range control.Controls {
				control, err = s.getControl(req.TargetOfEvaluation.CatalogId, elem.CategoryName, elem.Id)
				if err != nil {
					err = fmt.Errorf("could not get control for control id {%s}: %v", control.Id, err)
					log.Error(err)
					return
				}

				eval := &Evaluator{
					categoryName: control.CategoryName,
					controlId:    control.Id,
					firstLevel:   false,
				}

				evaluator = append(evaluator, eval)
			}
		}
	}

	s.scheduler.TagsUnique()

	// Start scheduler jobs
	for _, e := range evaluator {
		log.Infof("Evaluate Cloud Service '%s' for Control ID '%s' every 5 minutes...", req.TargetOfEvaluation.CloudServiceId, e.controlId)

		schedulerTag = createSchedulerTag(req.TargetOfEvaluation.CloudServiceId, e.controlId)

		// Regarding the control level the specific method is called
		if e.firstLevel {
			_, err = s.scheduler.
				Every(5).
				Minute().
				Tag(schedulerTag).
				Do(s.evaluateFirstLevelControl, req.TargetOfEvaluation, e.categoryName, e.controlId)
		} else {
			_, err = s.scheduler.
				Every(5).
				Minute().
				Tag(schedulerTag).
				Do(s.evaluateSecondLevelControl, req.TargetOfEvaluation, e.categoryName, e.controlId)
		}
		if err != nil {
			newErr := fmt.Errorf("evaluation for Cloud Service '%s' and Control ID '%s' cannot be scheduled", req.TargetOfEvaluation.CloudServiceId, e.controlId)
			log.Errorf("%s: %s", newErr, err)
			err = status.Errorf(codes.Internal, "%s", newErr)
			return
		}

		s.scheduler.StartAsync()
	}

	resp = &evaluation.StartEvaluationResponse{}

	return
}

// StopEvaluation is a method implementation of the evaluation interface: It stop the evaluation for a cloud service and its given controls (e.g., EUCS OPS-13,OPS-13.2)
func (s *Service) StopEvaluation(_ context.Context, req *evaluation.StopEvaluationRequest) (res *evaluation.StopEvaluationResponse, err error) {

	var (
		schedulerTag string
		control      *orchestrator.Control
		controlIds   = []string{}
	)

	err = req.Validate()
	if err != nil {
		err = status.Errorf(codes.InvalidArgument, "%v", err)
		return
	}

	// Check if control has further sub-controls
	control, err = s.getControl(req.TargetOfEvaluation.CatalogId, req.CategoryName, req.ControlId)
	if err != nil {
		err = fmt.Errorf("could not get control for control id {%s}: %v", req.ControlId, err)
		log.Error(err)
		return
	}

	// Get all control ids that need to be checked in the scheduler
	controlIds = append(controlIds, control.Id)
	for _, control := range control.Controls {
		controlIds = append(controlIds, control.Id)
	}

	// Check for each control id a job is running
	for _, controlId := range controlIds {
		schedulerTag = createSchedulerTag(req.TargetOfEvaluation.CloudServiceId, controlId)

		// Verify if scheduler job exists for the control id
		_, err = s.scheduler.FindJobsByTag(schedulerTag)

		if err == nil {
			// Delete the job from the scheduler
			err = s.scheduler.RemoveByTag(schedulerTag)
			if err != nil {
				err = fmt.Errorf("error when removing job from scheduler: %v", err)
				log.Error(err)
				err = status.Errorf(codes.Internal, "error when stopping scheduler")
				return
			}
		} else if strings.Contains(err.Error(), "no jobs found with given tag") {
			err = fmt.Errorf("evaluation of cloud service '%s' with '%s' not running", req.TargetOfEvaluation.CloudServiceId, controlId)
			log.Error(err)
			err = status.Errorf(codes.NotFound, err.Error())
		} else if err != nil {
			shortErr := errors.New("error when stopping scheduler")
			log.Errorf("%s: %v", shortErr, err)
			err = status.Errorf(codes.Internal, "%v", &shortErr)
			return
		}
	}

	res = &evaluation.StopEvaluationResponse{}

	return
}

// evaluateSecondLevelControl evaluates the second level controls, e.g., OPS-13.2
func (s *Service) evaluateSecondLevelControl(toe *orchestrator.TargetOfEvaluation, categoryName, controlId string) {
	log.Infof("Started evaluation for Cloud Service '%s',  Catalog ID '%s' and Controls '%s'", toe.CloudServiceId, toe.CatalogId, controlId)

	var (
		// metrics          []*assessment.Metric
		result *evaluation.EvaluationResult
		err    error
	)

	result, err = s.evaluationResultForLowerControlLevel(toe.CloudServiceId, toe.CatalogId, categoryName, controlId, toe)
	if err != nil {
		err = fmt.Errorf("error creating evaluation result: %v", err)
		log.Error(err)
		return
	}
	if result == nil {
		log.Debug("No evaluation result created")
		return
	}

	s.resultMutex.Lock()
	s.results[result.Id] = result
	s.resultMutex.Unlock()
	log.Infof("Evaluation result with ID %s stored.", result.Id)
}

// evaluateFirstLevelControl evaluates first level control, e.g., OPS-13
func (s *Service) evaluateFirstLevelControl(toe *orchestrator.TargetOfEvaluation, categoryName, controlId string) {
	// TODO(anatheka): TBD How should we evaluate the first level control(OPS-13)? It should be compliant if all sub-controls are compliant
	// Suggestion: We get all sub-control results and calculate it. But then we should start that job a bit later. Then we must start the scheduler in SingletonMode

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

// evaluationResultForLowerControlLevel return the evaluation result for the lower control level, e.g. OPS-13.7
func (s *Service) evaluationResultForLowerControlLevel(cloudServiceId, catalogId, categoryName, controlId string, toe *orchestrator.TargetOfEvaluation) (result *evaluation.EvaluationResult, err error) {

	// Get metrics from control and sub-controls
	metrics, err := s.getMetrics(catalogId, categoryName, controlId)
	if err != nil {
		err = fmt.Errorf("could not get metrics from control and sub-controls for Cloud Serivce '%s' from Orchestrator: %v", cloudServiceId, err)
		log.Error(err)
		return
	}

	// Get assessment results for the target cloud service
	// TODO(anatheka): Add FilterMetricId when PR#877 is ready
	assessmentResults, err := api.ListAllPaginated(&assessment.ListAssessmentResultsRequest{
		FilteredCloudServiceId: cloudServiceId}, s.orchestratorClient.ListAssessmentResults, func(res *assessment.ListAssessmentResultsResponse) []*assessment.AssessmentResult {
		return res.Results
	})
	if err != nil || len(assessmentResults) == 0 {
		// TODO(anatheka): Let we the scheduler running or do we want to stop it if we do not get assessment results from the orchestrator?
		// We let the scheduler running if we do not get the assessment results from the orchestrator, maybe it is only a temporary network problem
		err = fmt.Errorf("could not get assessment results for Cloud Service ID '%s' from Orchestrator: %v", cloudServiceId, err)
		log.Error(err)

		return
	}

	// Get mapping assessment results related to the metric
	mappingList := getMapping(assessmentResults, metrics)

	// Here the actual evaluation takes place. We check if all asssessment results are compliant or not. If at least one assessment result is not compliant the whole evaluation status is set to NOT_COMPLIANT. Furthermore, all non-compliant assessment results are stored in a separate list.
	// TODO(anatheka): Do we want the whole assessment result in evaluationResult.FailingAssessmentResults or only the ID?
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
	// TODO(anatheka): Or should we set it to UNSPECIFIED? What does it mean if we have no assessment results?
	if status == evaluation.EvaluationResult_PENDING {
		status = evaluation.EvaluationResult_COMPLIANT
	}

	// Create evaluation result
	// TODO(all): Store in DB
	result = &evaluation.EvaluationResult{
		Id:                       uuid.NewString(),
		Timestamp:                timestamppb.Now(),
		CategoryName:             categoryName,
		Control:                  controlId,
		TargetOfEvaluation:       toe,
		Status:                   status,
		FailingAssessmentResults: nonCompliantAssessmentResults,
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
	metrics = append(metrics, control.Metrics...)

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

// getAllControlIdsFromControl returns a list with the control id and the sub-control ids
func getAllControlIdsFromControl(control *orchestrator.Control) []string {
	controlIds := []string{}

	if control == nil {
		return []string{}
	}
	controlIds = append(controlIds, control.Id)
	for _, control := range control.Controls {
		controlIds = append(controlIds, control.Id)
	}

	return controlIds
}
