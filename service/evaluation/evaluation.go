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

// interval is the default interval time for the scheduler. If no interval is set in the StartEvaluationRequest, the default value is taken.
const defaultInterval int = 5

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

	// WaitGroup for each started evaluation
	wg map[string]*WaitGroup

	// Currently, evaluation results are just stored as a map (=in-memory). In the future, we will use a DB.
	results     map[string]*evaluation.EvaluationResult
	resultMutex sync.Mutex

	storage persistence.Storage
}

func init() {
	log = logrus.WithField("component", "evaluation")
}

const (
	// DefaultOrchestratorAddress specifies the default gRPC address of the orchestrator service.
	DefaultOrchestratorAddress = "localhost:9090"
	DefaultAssessmentAddress   = "localhost:9090"
)

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

// StartEvaluation is a method implementation of the evaluation interface: It starts the evaluation for a cloud service and the given controls_in_scope (e.g., EUCS OPS-13.2) in the target_of_evaluation periodically. If no inteval time is given, the default value of 5 minutes is used.
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

	// Get orchestrator client
	err = s.initOrchestratorClient()
	if err != nil {
		err = fmt.Errorf("could not set orchestrator client: %v", err)
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	// Get Target of Evaluation
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

	// Add the controlsInScope of the TargetOfEvaluation including their sub-controls to the scheduler
	for _, c := range toe.GetControlsInScope() {
		var controls []*orchestrator.Control

		// We must get the control from the orchestrator, as we do not know if the request is correct and the parent_ids are probably missing
		control, err := s.getControl(c.GetCategoryCatalogId(), c.GetCategoryName(), c.GetId())
		if err != nil {
			err = fmt.Errorf("could not get control for control id {%s}: %v", c.GetId(), err)
			log.Error(err)
			return nil, status.Errorf(codes.Internal, "%s", err)
		}

		// The parent control scheduler tag
		parentSchedulerTag := ""

		// The schedulerTag of the current control
		schedulerTag = createSchedulerTag(toe.GetCloudServiceId(), c.GetId())

		// Check if scheduler job for current controlId is already running
		// TODO(anatheka): Do we want to append the errors for already running jobs or do we ignore the error value?
		jobs, err := s.scheduler.FindJobsByTag(schedulerTag)
		if len(jobs) > 0 && err == nil {
			err = fmt.Errorf("evaluation for cloud service id '%s' and control id '%s' already started", toe.GetCloudServiceId(), c.GetId())
			log.Error(err)
			return nil, status.Errorf(codes.AlreadyExists, "%s", err)
		}

		// Collect control including sub-controls to start them later as scheduler jobs
		controls = append(controls, control)
		controls = append(controls, control.GetControls()...)

		// Add number of sub-controls to the WaitGroup. For the first level control (e.g., OPS-13) we have to wait until all the sub-controls (e.g., OPS-13.1) are ready.
		// Control is a first level control if no parentControlId exists. That is the reason why we have to get the control from the orchestrator above, if the parent_control_id is missing this fails
		if control.ParentControlId == nil {
			s.wg[schedulerTag] = &WaitGroup{
				wg:      &sync.WaitGroup{},
				wgMutex: sync.Mutex{},
			}
			s.wg[schedulerTag].wgMutex.Lock()
			s.wg[schedulerTag].wg.Add(len(c.GetControls()))
			s.wg[schedulerTag].wgMutex.Unlock()

			// parentSchedulerTag is the tag for the parent control, that is only needed if a first level control is schedules
			parentSchedulerTag = createSchedulerTag(toe.GetCloudServiceId(), control.GetId())
		}

		// Add control including sub-controls to the scheduler
		for _, control := range controls {
			err = s.addJobToScheduler(control, toe, parentSchedulerTag, interval)
			// We can return the error as it is
			if err != nil {
				return nil, err
			}
		}

		log.Infof("Evaluate Cloud Service '%s' for Control ID '%s' every 5 minutes...", toe.GetCloudServiceId(), control.GetId())
	}

	// Start scheduler jobs
	s.scheduler.StartAsync()
	log.Infof("Scheduler started...")

	resp = &evaluation.StartEvaluationResponse{}

	return
}

// StopEvaluation is a method implementation of the evaluation interface: It stop the evaluation for a cloud service and its given control (e.g., EUCS OPS-13,OPS-13.2). Only first level controls (e.g., OPS-13) or individually started sub-controls (e.g., OPS-13.1) can be stopped. A sub-control with a running evaluation for the parent control cannot be stopped as it is needed for the parent control.
func (s *Service) StopEvaluation(_ context.Context, req *evaluation.StopEvaluationRequest) (resp *evaluation.StopEvaluationResponse, err error) {

	var (
		schedulerTag string
		control      *orchestrator.Control
	)

	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		err = status.Errorf(codes.InvalidArgument, "%v", err)
		return
	}

	// Get control
	control, err = s.getControl(req.TargetOfEvaluation.CatalogId, req.CategoryName, req.ControlId)
	if err != nil {
		err = fmt.Errorf("could not get control for control id {%s}: %v", req.ControlId, err)
		log.Error(err)
		return
	}

	// Currently, we can only stop the evaluation for first level controls or sub-controls that are started individually. We are not able to stop a single sub-level control for a started first level control, e.g., OPS-13 is started with it's second level controls OPS-13.1, OPS-13.2 and OPS-13.3 than it is not possible only to stop OPS-13.2.
	// Check if the control is a first level control
	if control.ParentControlId == nil {
		// Control is a first level control, stop the scheduler job for the control and all sub-controls
		err = s.stopSchedulerJobs(getSchedulerTagsForControlIds(getAllControlIdsFromControl(control), req.GetTargetOfEvaluation().GetCloudServiceId()))
		if err != nil && strings.Contains(err.Error(), gocron.ErrJobNotFoundWithTag.Error()) {
			err = fmt.Errorf("evaluation for cloud service id '%s' with '%s' not running", req.GetTargetOfEvaluation().GetCloudServiceId(), req.GetControlId())
			log.Error(err)
			err = status.Errorf(codes.FailedPrecondition, "%v", err)
			return
		} else if err != nil {
			err = fmt.Errorf("error when stopping scheduler job for control id '%s'", control.GetId())
			log.Error(err)
			err = status.Errorf(codes.Internal, "%v", err)
			return
		}

		// TODO(anatheka): WEITERMACHEN!!! How do we have to delete the WaitGroup? We have to wait until the sub-controls evaluation is finished

		log.Infof("Evaluation stopped for Cloud Service ID '%s' with Control ID '%s'", req.TargetOfEvaluation.GetCloudServiceId(), control.GetId())
	} else {
		// Control is a second level control, check if the parent control is scheduled
		schedulerTag = createSchedulerTag(req.TargetOfEvaluation.CloudServiceId, *control.ParentControlId)
		_, err = s.scheduler.FindJobsByTag(schedulerTag)

		err = s.handleFindParentControlJobError(err, req.TargetOfEvaluation.GetCloudServiceId(), req.GetControlId())
		// No need of further error handling
		if err != nil {
			return
		}
	}

	resp = &evaluation.StopEvaluationResponse{}

	return
}

// ListEvaluationResults is a method implementation of the assessment interface
func (s *Service) ListEvaluationResults(_ context.Context, req *evaluation.ListEvaluationResultsRequest) (res *evaluation.ListEvaluationResultsResponse, err error) {
	var filtered_values []*evaluation.EvaluationResult

	// Filtering evaluation results by
	// * cloud service ID
	// * control ID
	for _, v := range s.results {
		if req.FilteredCloudServiceId != nil && v.TargetOfEvaluation.GetCloudServiceId() != req.GetFilteredCloudServiceId() {
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
	if parentSchedulerTag == "" {
		err = errors.New("parent scheduler tag is invalid")
	}
	if interval == 0 {
		err = errors.New("interval is invalid")
	}
	if err != nil {
		log.Error(err)
		return status.Errorf(codes.Internal, "%s", "evaluation cannot be scheduled")
	}

	// schedulerTag is the tag for the given control
	schedulerTag := createSchedulerTag(toe.GetCloudServiceId(), c.GetId())

	// Regarding the control level the specific method is called every X minutes based on the given interval. We have to decide if the control is sub-control that can be evaluated direclty or a first level control that has to wait for the results of the sub-level controls.
	if c.ParentControlId == nil { // first level control
		_, err = s.scheduler.
			Every(interval).
			Minute().
			Tag(schedulerTag).
			Do(s.evaluateFirstLevelControl, toe, c.GetCategoryName(), c.GetId(), schedulerTag, c.GetControls())
	} else { // second level control
		_, err = s.scheduler.
			Every(interval).
			Minute().
			Tag(schedulerTag).
			Do(s.evaluateSecondLevelControl, toe, c.GetCategoryName(), c.GetId(), parentSchedulerTag)
	}
	if err != nil {
		err = fmt.Errorf("evaluation for Cloud Service '%s' and Control ID '%s' cannot be scheduled: %v", toe.GetCloudServiceId(), c.GetId(), err)
		log.Error(err)
		return status.Errorf(codes.Internal, "%s", err)
	}

	log.Debugf("Cloud Service '%s' with Control ID '%s' added to scheduler", toe.GetCloudServiceId(), c.GetId())

	return
}

// handleFindParentControlJobError handles the scheduler.FindJobsByTag() error for the parent control id
func (s *Service) handleFindParentControlJobError(err error, cloudServiceId, controlId string) error {
	var schedulerTag string

	if err == nil {
		// Scheduler job for parent control cannot be removed because parent control is currently evaluated
		err = fmt.Errorf("evaluation of control id '%s' for cloud service '%s' can not be stopped because the control is a sub-control of the evaluated control '%s'", controlId, cloudServiceId, controlId)
		log.Error(err)
		return status.Errorf(codes.NotFound, err.Error())
	} else if strings.Contains(err.Error(), "no jobs found with given tag") {
		// Scheduler job can be removed because the parent control is currently not evaluated
		schedulerTag = createSchedulerTag(cloudServiceId, controlId)

		// Verify if scheduler job exists for the control id (sub-level control)
		_, err = s.scheduler.FindJobsByTag(schedulerTag)
		if err == nil {
			// Scheduler job is running and could be stopped
			err = s.stopSchedulerJob(schedulerTag)
			if err == nil {
				log.Debugf("Scheduler job for cloud service id '%s' with control id '%s' stopped.", cloudServiceId, controlId)
				return nil
			}
		} else if strings.Contains(err.Error(), "no jobs found with given tag") {
			// Scheduler job does not exist
			err = fmt.Errorf("evaluation for cloud service id '%s' with '%s' not running", cloudServiceId, controlId)
			log.Error(err)
			return status.Errorf(codes.NotFound, err.Error())
		}
	}

	// This is the error message in case a job can be stopped but an error has occurred.
	err = fmt.Errorf("error when stopping scheduler job for cloud service id '%s' with control id '%s'", cloudServiceId, controlId)
	log.Error(err)
	return status.Errorf(codes.Internal, err.Error())
}

// stopSchedulerJob stops a scheduler job for the given scheduler tag
func (s *Service) stopSchedulerJob(schedulerTag string) (err error) {
	// Delete the job from the scheduler
	err = s.scheduler.RemoveByTag(schedulerTag)
	if err != nil {
		err = fmt.Errorf("error when removing job for tag '%s' from scheduler: %v", schedulerTag, err)
		log.Error(err)
		return
	}

	return
}

// stopSchedulerJobs stops all scheduler jobs for the given scheduler tags
func (s *Service) stopSchedulerJobs(schedulerTags []string) (err error) {
	// Delete the job from the scheduler
	for _, schedulerTag := range schedulerTags {
		err = s.stopSchedulerJob(schedulerTag)
		if err != nil {
			return
		}
	}

	return
}

// evaluateFirstLevelControl evaluates a first level control, e.g., OPS-13. Therefere, the method needs to wait till all sub-controls are evaluated.
// TODO(all): Note: That is a first try. I'm not convinced, but I can't think of anything better at the moment.
func (s *Service) evaluateFirstLevelControl(toe *orchestrator.TargetOfEvaluation, categoryName, controlId, schedulerTag string, subControls []*orchestrator.Control) {
	var (
		status                        = evaluation.EvaluationResult_STATUS_UNSPECIFIED
		nonCompliantAssessmentResults []*assessment.AssessmentResult
		result                        *evaluation.EvaluationResult
	)

	log.Infof("Start first level control evaluation for Cloud Service '%s', Catalog ID '%s' and Control '%s'", toe.CloudServiceId, toe.CatalogId, controlId)

	// Wait till all sub-controls are evaluated
	s.wg[schedulerTag].wg.Wait()

	evaluations, err := s.ListEvaluationResults(context.Background(), &evaluation.ListEvaluationResultsRequest{
		FilteredCloudServiceId: &toe.CloudServiceId,
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

	// Find all needed evaluation results and calculate compliant status and add non-compliant assessment results
	for _, r := range evaluations.Results {
		if controlContains(toe.GetControlsInScope(), r.GetControlId()) && r.Status == evaluation.EvaluationResult_COMPLIANT && status != evaluation.EvaluationResult_NOT_COMPLIANT {
			status = evaluation.EvaluationResult_COMPLIANT
		} else {
			status = evaluation.EvaluationResult_NOT_COMPLIANT
			nonCompliantAssessmentResults = append(nonCompliantAssessmentResults, r.FailingAssessmentResults...)
		}
	}

	// Create evaluation result
	// TODO(all): Store in DB
	result = &evaluation.EvaluationResult{
		Id:                       uuid.NewString(),
		Timestamp:                timestamppb.Now(),
		CategoryName:             categoryName,
		ControlId:                controlId,
		TargetOfEvaluation:       toe,
		Status:                   status,
		FailingAssessmentResults: nonCompliantAssessmentResults,
	}

	s.resultMutex.Lock()
	s.results[result.Id] = result
	s.resultMutex.Unlock()

	log.Infof("Evaluation result for ControlID '%s' with ID '%s' stored.", controlId, result.Id)
}

// evaluateSecondLevelControl evaluates the second level controls, e.g., OPS-13.2
func (s *Service) evaluateSecondLevelControl(toe *orchestrator.TargetOfEvaluation, categoryName, controlId, parentSchedulerTag string) {
	log.Infof("Start second level control evaluation for Cloud Service '%s', Catalog ID '%s' and Control '%s'", toe.CloudServiceId, toe.CatalogId, controlId)

	var (
		result *evaluation.EvaluationResult
		err    error
	)

	result, err = s.evaluationResultForSecondControlLevel(toe.GetCloudServiceId(), toe.GetCatalogId(), categoryName, controlId, toe)
	if err != nil {
		err = fmt.Errorf("error creating evaluation result: %v", err)
		log.Error(err)
	} else if result == nil {
		log.Debug("No evaluation result created")
	} else {

		s.resultMutex.Lock()
		s.results[result.Id] = result
		s.resultMutex.Unlock()

		log.Infof("Evaluation result for Cloud Service '%s' and ControlID '%s' stored with ID '%s'.", toe.GetCloudServiceId(), controlId, result.Id)
	}

	// If the parentSchedulerTag is not empty, we have do decrement the WaitGroup so that the parent control can also be evaluated when all sub-controls are evaluated.
	if parentSchedulerTag != "" {
		s.wg[parentSchedulerTag].wgMutex.Lock()
		s.wg[parentSchedulerTag].wg.Done()
		s.wg[parentSchedulerTag].wgMutex.Unlock()
	}
}

// evaluationResultForSecondControlLevel return the evaluation result for the lower control level, e.g. OPS-13.7
// TODO(anatheka): Refactor cloudServiceId, etc. with toe.*
func (s *Service) evaluationResultForSecondControlLevel(cloudServiceId, catalogId, categoryName, controlId string, toe *orchestrator.TargetOfEvaluation) (result *evaluation.EvaluationResult, err error) {
	var (
		assessmentResults []*assessment.AssessmentResult
	)

	// Get metrics from control and sub-controls
	metrics, err := s.getMetrics(catalogId, categoryName, controlId)
	if err != nil {
		err = fmt.Errorf("could not get metrics from control and sub-controls for Cloud Serivce '%s' from Orchestrator: %v", cloudServiceId, err)
		return
	}

	// TODO(anatheka): Is the evaluation result COMPLIANT or should it be UNSPECIFIED if no metrics exist for the control?
	// Currently, the evaluation continues when no metrics exist and the evaluation result becomes COMPLIANT.
	if len(metrics) != 0 {
		// Get assessment results filtered for
		// * cloud service id
		// * metric ids
		assessmentResults, err = api.ListAllPaginated(&assessment.ListAssessmentResultsRequest{
			FilteredCloudServiceId: &cloudServiceId,
			FilteredMetricId:       getMetricIds(metrics),
		}, s.orchestratorClient.ListAssessmentResults, func(res *assessment.ListAssessmentResultsResponse) []*assessment.AssessmentResult {
			return res.Results
		})

		if err != nil || len(assessmentResults) == 0 {
			// TODO(anatheka): Let we the scheduler running or do we want to stop it if we do not get assessment results from the orchestrator?
			// We let the scheduler running if we do not get the assessment results from the orchestrator, maybe it is only a temporary network problem
			err = fmt.Errorf("could not get assessment results for Cloud Service ID '%s' from Orchestrator: %v", cloudServiceId, err)
			return nil, err
		}
	} else {
		log.Infof("no metrics are evailable for the given control")
	}

	// Here the actual evaluation takes place. We check if all asssessment results are compliant or not. If at least one assessment result is not compliant the whole evaluation status is set to NOT_COMPLIANT. Furthermore, all non-compliant assessment results are stored in a separate list.
	// TODO(anatheka): Do we want the whole assessment result in evaluationResult.FailingAssessmentResults or only the IDs?
	var nonCompliantAssessmentResults []*assessment.AssessmentResult
	status := evaluation.EvaluationResult_PENDING
	for _, elem := range assessmentResults {
		if !elem.Compliant {
			nonCompliantAssessmentResults = append(nonCompliantAssessmentResults, elem)
			status = evaluation.EvaluationResult_NOT_COMPLIANT
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
		ControlId:                controlId,
		TargetOfEvaluation:       toe,
		Status:                   status,
		FailingAssessmentResults: nonCompliantAssessmentResults,
	}

	return
}

// TODO(anatheka): Move to internals?
// controlContains return true if the given control id is a sub-control of the controls slice
func controlContains(controls []*orchestrator.Control, controlId string) bool {
	for _, c := range controls {
		if c.GetId() == controlId {
			return true
		}
	}

	return false
}

// getMetricIds returns the metric Ids for the given metrics
func getMetricIds(metrics []*assessment.Metric) []string {
	var metricIds []string

	for _, m := range metrics {
		metricIds = append(metricIds, m.GetId())
	}

	return metricIds
}

// getMetrics returns all metrics from a given controlId
// For now a control has either sub-controls or metrics. If the control has sub-controls, get also all metrics from the sub-controls.
// TODO(anatheka): Is it possible that a control has sub-controls and metrics?
// TODO(anatheka): Refactor catalogId, etc.?
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

// getMetricsFromSubControls returns a list of metrics from the sub-controls.
func (s *Service) getMetricsFromSubControls(control *orchestrator.Control) (metrics []*assessment.Metric, err error) {
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

// getAllControlIdsFromControl returns a list with the control id and the sub-control ids.
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

// createSchedulerTag creates the scheduler tag for a given cloud_service_id and control_id.
func createSchedulerTag(cloudServiceId, controlId string) string {
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
		schedulerTags = append(schedulerTags, createSchedulerTag(cloudServiceId, controlId))
	}

	return
}
