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
	"clouditor.io/clouditor/internal/util"
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

const (
	// DefaultOrchestratorAddress specifies the default gRPC address of the orchestrator service.
	DefaultOrchestratorAddress = "localhost:9090"

	// defaultinterval is the default interval time for the scheduler. If no interval is set in the StartEvaluationRequest, the default value is taken.
	defaultInterval int = 5
)

type grpcTarget struct {
	target string
	opts   []grpc.DialOption
}

type WaitGroup struct {
	wg      *sync.WaitGroup
	wgMutex sync.Mutex
}

// Service is an implementation of the Clouditor Evaluation service
type Service struct {
	evaluation.UnimplementedEvaluationServer
	orchestratorClient  orchestrator.OrchestratorClient
	orchestratorAddress grpcTarget
	authorizer          api.Authorizer
	scheduler           *gocron.Scheduler

	// wg is used for a control that waits for its sub-controls to be evaluated
	wg map[string]*WaitGroup

	// Currently, evaluation results are just stored as a map (=in-memory). In the future, we will use a DB.
	results     map[string]*evaluation.EvaluationResult
	resultMutex sync.Mutex

	storage persistence.Storage

	// TODO(all): Comment in once the evaluation results are stored in storage
	// // authz defines our authorization strategy, e.g., which user can access which cloud service and associated
	// // resources, such as evidences and assessment results.
	// authz service.AuthorizationStrategy
}

func init() {
	log = logrus.WithField("component", "evaluation")
}

// ServiceOption is a function-style option to configure the Evaluation Service
type ServiceOption func(*Service)

// WithStorage is an option to set the storage. If not set, NewService will use inmemory storage.
func WithStorage(storage persistence.Storage) service.Option[Service] {
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
func WithOrchestratorAddress(address string, opts ...grpc.DialOption) service.Option[Service] {
	return func(s *Service) {
		s.orchestratorAddress = grpcTarget{
			target: address,
			opts:   opts,
		}
	}
}

// NewService creates a new Evaluation service
func NewService(opts ...service.Option[Service]) *Service {
	s := Service{
		orchestratorAddress: grpcTarget{
			target: DefaultOrchestratorAddress,
		},
		results:   make(map[string]*evaluation.EvaluationResult),
		scheduler: gocron.NewScheduler(time.UTC),
		wg:        make(map[string]*WaitGroup),
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

// StartEvaluation is a method implementation of the evaluation interface: It periodically starts the evaluation of a cloud service and the given controls_in_scope (e.g., EUCS OPS-13, EUCS OPS-13.2) in the target_of_evaluation. If no inteval time is given, the default value is used.
func (s *Service) StartEvaluation(_ context.Context, req *evaluation.StartEvaluationRequest) (resp *evaluation.StartEvaluationResponse, err error) {
	var (
		schedulerTag string
		interval     int
		toe          *orchestrator.TargetOfEvaluation
	)

	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// Set the interval to the default value if not set. If the interval is set to 0, the default interval is used.
	if req.GetInterval() == 0 {
		interval = defaultInterval
	} else {
		interval = int(req.GetInterval())
	}

	// Get orchestrator client. The orchestrator client is used to retrieve necessary information from the Orchestrator, such as assessment_results, controls or targets_of_evaluation.
	err = s.initOrchestratorClient()
	if err != nil {
		err = fmt.Errorf("could not set orchestrator client: %v", err)
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	// Get Target of Evaluation. The Target of Evaluation is retrieved to get the controls_in_scope, which are then evaluated.
	toe, err = s.orchestratorClient.GetTargetOfEvaluation(context.Background(), &orchestrator.GetTargetOfEvaluationRequest{
		CloudServiceId: req.GetCloudServiceId(),
		CatalogId:      req.GetCatalogId(),
	})
	if err != nil {
		err = fmt.Errorf("could not get target of evaluation: %v", err)
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	log.Info("Starting evaluation ...")

	// Scheduler tags must be unique to find the jobs by tag name
	s.scheduler.TagsUnique()

	// Add the controls_in_scope of the target_of_evaluation including their sub-controls to the scheduler. The parent control has to wait for the evaluation of the sub-controls. That's why we need to know how much sub-controls are available and define the waitGroup with the number of the corresponding sub-controls. The controls_in_scope are not stored in a hierarchy, so we have to get the parent control and find all related sub-controls.
	controlsInScope := getControlsInScopeHierarchy(toe.GetControlsInScope())
	for _, control := range controlsInScope {
		var (
			controls []*orchestrator.Control

			// The parent control scheduler tag
			parentSchedulerTag = ""
		)

		// If the control does not have sub-controls don't start a scheduler.
		if control.ParentControlId != nil {
			continue
		}

		// Collect control including sub-controls from the controls_in_scope slice to start them later as scheduler jobs
		controls = append(controls, control)
		controls = append(controls, control.GetControls()...)

		// parentSchedulerTag is the tag for the parent control (e.g., OPS-13)
		parentSchedulerTag = getSchedulerTag(toe.GetCloudServiceId(), control.GetId())

		// Check if scheduler job for current control_id is already running
		jobs, err := s.scheduler.FindJobsByTag(schedulerTag)
		if len(jobs) > 0 && err == nil {
			err = fmt.Errorf("evaluation for Cloud Service '%s' and Control ID '%s' already started", toe.GetCloudServiceId(), control.GetId())
			log.Error(err)
			return nil, status.Errorf(codes.AlreadyExists, "%s", err)
		}

		// Add number of sub-controls to the WaitGroup. For the control (e.g., OPS-13) we have to wait until all the sub-controls (e.g., OPS-13.1) are ready.
		// Control is a parent control if no parentControlId exists.
		s.wg[schedulerTag] = &WaitGroup{
			wg:      &sync.WaitGroup{},
			wgMutex: sync.Mutex{},
		}
		s.wg[schedulerTag].wgMutex.Lock()
		// The controls list contains also the parent control itself and must be minimized by 1.
		s.wg[schedulerTag].wg.Add(len(controls) - 1)
		s.wg[schedulerTag].wgMutex.Unlock()

		// Add control including sub-controls to the scheduler
		for _, control := range controls {
			err = s.addJobToScheduler(control, toe, parentSchedulerTag, interval)
			// We can return the error as it is
			if err != nil {
				return nil, err
			}
		}

		log.Infof("Evaluate Control ID '%s' for Cloud Service '%s' every 5 minutes...", control.GetId(), toe.GetCloudServiceId())
	}

	// Start scheduler jobs
	s.scheduler.StartAsync()
	log.Infof("Evaluation started.")

	resp = &evaluation.StartEvaluationResponse{}

	return
}

// StopEvaluation is a method implementation of the evaluation interface: It stops the evaluation for a TargetOfEvaluation.
func (s *Service) StopEvaluation(_ context.Context, req *evaluation.StopEvaluationRequest) (resp *evaluation.StopEvaluationResponse, err error) {

	var (
		schedulerTag string
		// control      *orchestrator.Control
		toe *orchestrator.TargetOfEvaluation
	)

	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		err = status.Errorf(codes.InvalidArgument, "%v", err)
		return
	}

	// Get Target of Evaluation. The Target of Evaluation is retrieved to get the controls_in_scope, which are then stopped from the evaluation.
	toe, err = s.orchestratorClient.GetTargetOfEvaluation(context.Background(), &orchestrator.GetTargetOfEvaluationRequest{
		CloudServiceId: req.GetCloudServiceId(),
		CatalogId:      req.GetCatalogId(),
	})
	if err != nil {
		err = fmt.Errorf("could not get target of evaluation: %v", err)
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	// Stop all controls in the list controls_in_scope
	errorControlIds := []string{}
	for _, control := range toe.ControlsInScope {
		// Check if the control is a parent control
		if control.ParentControlId == nil {
			// Control is a parent control, stop the scheduler job for the control and all sub-controls
			err = s.stopSchedulerJobs(getSchedulerTagsForControlIds(getAllControlIdsFromControl(control), req.GetCloudServiceId()))
			if err != nil && !strings.Contains(err.Error(), gocron.ErrJobNotFoundWithTag.Error()) {
				log.Error(err)
				errorControlIds = append(errorControlIds, control.GetId())
			}

			log.Debugf("Evaluation stopped for Cloud Service ID '%s' with Control ID '%s'", req.GetCloudServiceId(), control.GetId())
		} else {
			// Control is a sub-control
			schedulerTag = getSchedulerTag(req.GetCloudServiceId(), control.GetId())
			err = s.stopSchedulerJob(schedulerTag)
			// No need of further error handling
			if err != nil && !strings.Contains(err.Error(), gocron.ErrJobNotFoundWithTag.Error()) {
				log.Error(err)
				errorControlIds = append(errorControlIds, control.GetId())
			}

			log.Debugf("Evaluation stopped for Cloud Service ID '%s' with Control ID '%s'", req.GetCloudServiceId(), control.GetId())
		}
	}

	if len(errorControlIds) > 0 {
		log.Infof("Error stopping scheduler for Controls '%v' for Cloud Service '%s'.", strings.Join(errorControlIds, ", "), req.GetCloudServiceId())
	} else {
		log.Infof("Evaluation for Cloud Service '%s' and Catalog '%s' stopped.", req.GetCloudServiceId(), req.GetCatalogId())
	}

	resp = &evaluation.StopEvaluationResponse{}

	return
}

// ListEvaluationResults is a method implementation of the assessment interface
func (s *Service) ListEvaluationResults(_ context.Context, req *evaluation.ListEvaluationResultsRequest) (res *evaluation.ListEvaluationResultsResponse, err error) {
	var (
		filtered_values []*evaluation.EvaluationResult
		// TODO(all): Comment in once the evaluation results are stored in storage
		// allowed         []string
		// all             bool
		// query           []string
		// args            []any
	)

	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	// TODO(all): Comment in once the evaluation results are stored in storage
	// // Retrieve list of allowed cloud service according to our authorization strategy. No need to specify any conditions to our storage request, if we are allowed to see all cloud services.
	// all, allowed = s.authz.AllowedCloudServices(ctx)

	// // The content of the filtered cloud service ID must be in the list of allowed cloud service IDs,
	// // unless one can access *all* the cloud services.
	// if !all && req.FilteredCloudServiceId != nil && !slices.Contains(allowed, req.GetFilteredCloudServiceId()) {
	// 	return nil, service.ErrPermissionDenied
	// }

	// // Filtering evaluation results by
	// // * cloud service ID
	// // * control ID
	// if req.GetFilteredCloudServiceId() != "" {
	// 	query = append(query, "cloud_service_id = ?")
	// 	args = append(args, req.GetFilteredCloudServiceId())
	// }
	// if req.GetFilteredControlId() != "" {
	// 	query = append(query, "control_id = ?")
	// 	args = append(args, req.GetFilteredControlId())
	// }

	// // In any case, we need to make sure that we only select assessment results of cloud services that we have access to (if we do not have access to all)
	// if !all {
	// 	query = append(query, "cloud_service_id IN ?")
	// 	args = append(args, allowed)
	// }

	// // join query with AND and prepend the query
	// args = append([]any{strings.Join(query, " AND ")}, args...)

	// res = new(evaluation.ListEvaluationResultsResponse)

	// // Paginate the results according to the request
	// res.Results, res.NextPageToken, err = service.PaginateStorage[*evaluation.EvaluationResult](req, s.storage, service.DefaultPaginationOpts, args...)
	// if err != nil {
	// 	err = fmt.Errorf("could not paginate evaluation results: %v", err)
	// 	log.Error(err)
	// 	return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
	// }

	for _, v := range s.results {
		if req.FilteredCloudServiceId != nil && v.GetCloudServiceId() != req.GetFilteredCloudServiceId() {
			continue
		}

		if req.FilteredControlId != nil && v.ControlId != req.GetFilteredControlId() {
			continue
		}

		filtered_values = append(filtered_values, v)
	}

	res = new(evaluation.ListEvaluationResultsResponse)

	// Paginate the results according to the request
	res.Results, res.NextPageToken, err = service.PaginateSlice(req, filtered_values, func(a *evaluation.EvaluationResult, b *evaluation.EvaluationResult) bool {
		return a.Id < b.Id
	}, service.DefaultPaginationOpts)
	if err != nil {
		err = fmt.Errorf("could not paginate evaluation results: %v", err)
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return
}

// addJobToScheduler adds a job for the given control to the scheduler and sets the scheduler interval to the given interval
func (s *Service) addJobToScheduler(c *orchestrator.Control, toe *orchestrator.TargetOfEvaluation, parentSchedulerTag string, interval int) (err error) {
	// Check inputs and log error
	if c == nil {
		err = errors.New("control is invalid")
	}
	if toe == nil {
		err = errors.New("target of evaluation is invalid")
	}
	if interval == 0 {
		err = errors.New("interval is invalid")
	}
	if err != nil {
		log.Error(err)
		return status.Errorf(codes.Internal, "evaluation cannot be scheduled: %v", err)
	}

	// schedulerTag is the tag for the given control
	schedulerTag := getSchedulerTag(toe.GetCloudServiceId(), c.GetId())

	// Regarding the control level the specific method is called every X minutes based on the given interval. We have to decide if a sub-control is started individually or a parent control that has to wait for the results of the sub-controls.
	// If a parent control with its sub-controls is started, the parentSchedulerTag is empty and the ParentControlId is not set.
	// If a sub-control is started individually the parentSchedulerTag is empty and ParentControlId is set.
	if parentSchedulerTag == "" && c.ParentControlId == nil { // parent control
		_, err = s.scheduler.
			Every(interval).
			Minute().
			Tag(schedulerTag).
			Do(s.evaluateControl, toe, c.GetCategoryName(), c.GetId(), schedulerTag, c.GetControls())
	} else { // sub-control
		_, err = s.scheduler.
			Every(interval).
			Minute().
			Tag(schedulerTag).
			Do(s.evaluateSubcontrol, toe, c.GetCategoryName(), c.GetId(), parentSchedulerTag)
	}
	if err != nil {
		err = fmt.Errorf("evaluation for Cloud Service '%s' and Control ID '%s' cannot be scheduled: %v", toe.GetCloudServiceId(), c.GetId(), err)
		log.Error(err)
		return status.Errorf(codes.Internal, "%s", err)
	}

	log.Debugf("Cloud Service '%s' with Control ID '%s' added to scheduler", toe.GetCloudServiceId(), c.GetId())

	return
}

// stopSchedulerJob stops a scheduler job for the given scheduler tag
func (s *Service) stopSchedulerJob(schedulerTag string) (err error) {
	// Remove job for tag from the scheduler
	err = s.scheduler.RemoveByTag(schedulerTag)
	if err != nil {
		err = fmt.Errorf("error while removing job for tag '%s' from scheduler: %v", schedulerTag, err)
		return
	}

	return
}

// stopSchedulerJobs stops all scheduler jobs for the given scheduler tags
func (s *Service) stopSchedulerJobs(schedulerTags []string) (err error) {
	// Remove all job for given tags from the scheduler
	for _, schedulerTag := range schedulerTags {
		err = s.stopSchedulerJob(schedulerTag)
		if err != nil {
			return
		}
	}

	return
}

// evaluateControl evaluates a control, e.g., OPS-13. Therefere, the method needs to wait till all sub-controls (e.g., OPS-13.1) are evaluated.
// TODO(all): Note: That is a first try. I'm not convinced, but I can't think of anything better at the moment.
func (s *Service) evaluateControl(toe *orchestrator.TargetOfEvaluation, categoryName, controlId, schedulerTag string, subControls []*orchestrator.Control) {
	var (
		status                        = evaluation.EvaluationResult_STATUS_UNSPECIFIED
		nonCompliantAssessmentResults []string
		result                        *evaluation.EvaluationResult
	)

	log.Infof("Start control evaluation for Cloud Service '%s', Catalog ID '%s' and Control '%s'", toe.CloudServiceId, toe.CatalogId, controlId)

	// Wait till all sub-controls are evaluated
	s.wg[schedulerTag].wg.Wait()

	evaluations, err := s.ListEvaluationResults(context.Background(), &evaluation.ListEvaluationResultsRequest{
		FilteredCloudServiceId: &toe.CloudServiceId,
		FilteredControlId:      &controlId,
		// TODO(all): If evaluation results are stored in DB: getLastedByControlID or similar.
	})
	if err != nil {
		err = fmt.Errorf("error list evaluation results: %v", err)
		log.Error(err)
		return
	}

	// For the next iteration set wg again to the number of sub-controls
	s.wg[schedulerTag].wgMutex.Lock()
	s.wg[schedulerTag].wg.Add(len(subControls))
	s.wg[schedulerTag].wgMutex.Unlock()

	for _, r := range evaluations.Results {
		if r.Status == evaluation.EvaluationResult_COMPLIANT && status != evaluation.EvaluationResult_NOT_COMPLIANT {
			status = evaluation.EvaluationResult_COMPLIANT
		} else {
			status = evaluation.EvaluationResult_NOT_COMPLIANT
			nonCompliantAssessmentResults = append(nonCompliantAssessmentResults, r.GetFailingAssessmentResultsId()...)
		}
	}

	// Create evaluation result
	result = &evaluation.EvaluationResult{
		Id:                         uuid.NewString(),
		Timestamp:                  timestamppb.Now(),
		CategoryName:               categoryName,
		ControlId:                  controlId,
		CloudServiceId:             toe.GetCloudServiceId(),
		Status:                     status,
		FailingAssessmentResultsId: nonCompliantAssessmentResults,
	}

	// TODO(all): Store in DB
	s.resultMutex.Lock()
	s.results[result.Id] = result
	s.resultMutex.Unlock()

	log.Debugf("Evaluation result for ControlID '%s' with ID '%s' stored.", controlId, result.Id)
}

// evaluateSubcontrol evaluates the sub-controls, e.g., OPS-13.2
func (s *Service) evaluateSubcontrol(toe *orchestrator.TargetOfEvaluation, categoryName, controlId, parentSchedulerTag string) {
	var (
		result *evaluation.EvaluationResult
		err    error
	)

	result, err = s.evaluationResultForSubcontrol(toe.GetCloudServiceId(), toe.GetCatalogId(), categoryName, controlId, toe)
	if err != nil {
		err = fmt.Errorf("error creating evaluation result: %v", err)
		log.Error(err)
	} else if result == nil {
		log.Debug("No evaluation result created")
	} else {

		//TODO(all): Store in DB
		s.resultMutex.Lock()
		s.results[result.Id] = result
		s.resultMutex.Unlock()

		log.Infof("Evaluation result for Cloud Service '%s' and ControlID '%s' stored with ID '%s'.", toe.GetCloudServiceId(), controlId, result.Id)
	}

	// If the parentSchedulerTag is not empty, we have do decrement the WaitGroup so that the parent control can also be evaluated when all sub-controls are evaluated.
	if parentSchedulerTag != "" && s.wg[parentSchedulerTag] != nil {
		s.wg[parentSchedulerTag].wgMutex.Lock()
		s.wg[parentSchedulerTag].wg.Done()
		s.wg[parentSchedulerTag].wgMutex.Unlock()
	}
}

// evaluationResultForSubcontrol return the evaluation result for the sub-control, e.g. OPS-13.7
// TODO(anatheka): Refactor cloudServiceId, etc. with toe.*
func (s *Service) evaluationResultForSubcontrol(cloudServiceId, catalogId, categoryName, controlId string, toe *orchestrator.TargetOfEvaluation) (eval *evaluation.EvaluationResult, err error) {
	var (
		assessmentResults []*assessment.AssessmentResult
	)

	// Get metrics from control and sub-controls
	metrics, err := s.getAllMetricsFromControl(catalogId, categoryName, controlId)
	if err != nil {
		err = fmt.Errorf("could not get metrics for controlID '%s' for Cloud Service '%s' from Orchestrator: %v", controlId, cloudServiceId, err)
		return
	}

	// TODO(anatheka): Is the evaluation result COMPLIANT or should it be UNSPECIFIED if no metrics exist for the control?
	// Currently, the evaluation continues when no metrics exist and the evaluation result becomes COMPLIANT.
	if len(metrics) != 0 {
		// Get latest assessment_results by resource_id filtered by
		// * cloud service id
		// * metric ids
		assessmentResults, err = api.ListAllPaginated(&assessment.ListAssessmentResultsRequest{
			FilteredCloudServiceId: &cloudServiceId,
			FilteredMetricId:       getMetricIds(metrics),
			LatestByResourceId:     util.Ref(true),
		}, s.orchestratorClient.ListAssessmentResults, func(res *assessment.ListAssessmentResultsResponse) []*assessment.AssessmentResult {
			return res.Results
		})

		if err != nil || len(assessmentResults) == 0 {
			// We let the scheduler running if we do not get the assessment results from the orchestrator, maybe it is only a temporary network problem
			err = fmt.Errorf("could not get assessment results for Cloud Service ID '%s' and MetricIds '%s' from Orchestrator: %v", cloudServiceId, getMetricIds(metrics), err)
			return nil, err
		} else if len(assessmentResults) == 0 {
			// We let the scheduler running if we do not get the assessment results from the orchestrator, maybe it is only a temporary network problem
			err = fmt.Errorf("no assessment results for Cloud Service ID '%s' and MetricIds '%s' available", cloudServiceId, getMetricIds(metrics))
			return nil, err
		}
	} else {
		log.Debugf("no metrics are available for the given control")
	}

	// Here the actual evaluation takes place. For every resource_id we check if the asssessment results are compliant. If at least one assessment result per resource_id is not compliant the whole evaluation status is set to NOT_COMPLIANT. Furthermore, all non-compliant assessment_result_ids are stored in a separate list.

	// Get a map of the assessment results, so that we have all assessment results for a specific resource_id together for evaluation
	assessmentResultsMap := getAssessmentResultMap(assessmentResults)
	for _, result := range assessmentResultsMap {
		var nonCompliantAssessmentResults []string
		status := evaluation.EvaluationResult_PENDING

		// If no assessment_results are available continue
		if len(result) == 0 {
			continue
		}

		for i := range result {
			if !result[i].Compliant {
				nonCompliantAssessmentResults = append(nonCompliantAssessmentResults, result[i].GetId()) // It does not matter which element is used here, since they all have the same resource_id.
				status = evaluation.EvaluationResult_NOT_COMPLIANT
			}
		}

		// If no assessment results are available for the metric, no evaluation result is created
		if status == evaluation.EvaluationResult_PENDING {
			log.Debugf("No assessment results for resource '%s' available.", result[0].GetResourceId())
			continue
		}

		// Create evaluation result
		// TODO(all): Store in DB
		eval = &evaluation.EvaluationResult{
			Id:                         uuid.NewString(),
			Timestamp:                  timestamppb.Now(),
			CategoryName:               categoryName,
			ControlId:                  controlId,
			CloudServiceId:             toe.GetCloudServiceId(),
			CatalogId:                  toe.GetCatalogId(),
			ResourceId:                 result[0].GetResourceId(), // It does not matter which element is used here, since they all have the same resource_id.
			Status:                     status,
			FailingAssessmentResultsId: nonCompliantAssessmentResults,
		}
	}

	return
}

// getMetricIds returns the metric Ids for the given metrics
func getMetricIds(metrics []*assessment.Metric) []string {
	var metricIds []string

	for _, m := range metrics {
		metricIds = append(metricIds, m.GetId())
	}

	return metricIds
}

// getAllMetricsFromControl returns all metrics from a given controlId
// For now a control has either sub-controls or metrics. If the control has sub-controls, get also all metrics from the sub-controls.
func (s *Service) getAllMetricsFromControl(catalogId, categoryName, controlId string) (metrics []*assessment.Metric, err error) {
	var subControlMetrics []*assessment.Metric

	control, err := s.getControl(catalogId, categoryName, controlId)
	if err != nil {
		err = fmt.Errorf("could not get control for control id {%s}: %v", controlId, err)
		return
	}

	// Add metric of control to the metrics list
	metrics = append(metrics, control.Metrics...)

	// Add sub-control metrics to the metric list if exist
	if len(control.Controls) != 0 {
		// Get the metrics from the next sub-control
		subControlMetrics, err = s.getMetricsFromSubcontrols(control)
		if err != nil {
			err = fmt.Errorf("error getting metrics from sub-controls: %v", err)
			return
		}

		metrics = append(metrics, subControlMetrics...)
	}

	return
}

// getMetricsFromSubcontrols returns a list of metrics from the sub-controls.
func (s *Service) getMetricsFromSubcontrols(control *orchestrator.Control) (metrics []*assessment.Metric, err error) {
	var subcontrol *orchestrator.Control

	if control == nil {
		return nil, errors.New("control is missing")
	}

	for _, control := range control.Controls {
		subcontrol, err = s.getControl(control.CategoryCatalogId, control.CategoryName, control.Id)
		if err != nil {
			return
		}

		metrics = append(metrics, subcontrol.Metrics...)
	}

	return
}

// getAllControlIdsFromControl returns a list with the control and the sub-control ids.
func getAllControlIdsFromControl(control *orchestrator.Control) []string {
	// controlIds := make([]string, 0)
	var controlIds []string

	if control == nil {
		return []string{}
	}
	controlIds = append(controlIds, control.Id)
	for _, control := range control.Controls {
		controlIds = append(controlIds, control.Id)
	}

	return controlIds
}

// getControl returns the control for the given control_id.
func (s *Service) getControl(catalogId, categoryName, controlId string) (control *orchestrator.Control, err error) {

	if s.orchestratorClient == nil {
		err := s.initOrchestratorClient()
		if err != nil {
			return nil, err
		}
	}

	control, err = s.orchestratorClient.GetControl(context.Background(), &orchestrator.GetControlRequest{
		CatalogId:    catalogId,
		CategoryName: categoryName,
		ControlId:    controlId,
	})

	return
}

// TODO(all): Make a generic method for that in folder internal?
// initOrchestratorClient sets the orchestrator client.
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

// getSchedulerTag creates the scheduler tag for a given cloud_service_id and control_id.
func getSchedulerTag(cloudServiceId, controlId string) string {
	if cloudServiceId == "" || controlId == "" {
		return ""
	}
	return fmt.Sprintf("%s-%s", cloudServiceId, controlId)
}

// getSchedulerTagsForControlIds return for a given list of control_ids the corresponding scheduler_tags.
func getSchedulerTagsForControlIds(controlIds []string, cloudServiceId string) (schedulerTags []string) {
	schedulerTags = []string{}

	if len(controlIds) == 0 || cloudServiceId == "" {
		return
	}

	for _, controlId := range controlIds {
		schedulerTags = append(schedulerTags, getSchedulerTag(cloudServiceId, controlId))
	}

	return
}

// getControlsInScopeHierarchy return a controls list as hierarchy regarding the parent and sub-controls.
func getControlsInScopeHierarchy(controls []*orchestrator.Control) (controlsHierarchy []*orchestrator.Control) {
	var temp = make(map[string]*orchestrator.Control)

	for i := range controls {
		if controls[i].ParentControlId == nil {
			temp[controls[i].GetId()] = controls[i]
		}
	}

	for i := range controls {
		if controls[i].ParentControlId != nil {
			temp[controls[i].GetParentControlId()].Controls = append(temp[controls[i].GetParentControlId()].Controls, controls[i])
		}
	}

	for i := range temp {
		controlsHierarchy = append(controlsHierarchy, temp[i])
	}

	return
}

func getAssessmentResultMap(results []*assessment.AssessmentResult) map[string][]*assessment.AssessmentResult {
	var hierarchyResults = make(map[string][]*assessment.AssessmentResult)

	for i := range results {
		hierarchyResults[results[i].GetId()] = append(hierarchyResults[results[i].GetId()], results[i])
	}

	return hierarchyResults
}
