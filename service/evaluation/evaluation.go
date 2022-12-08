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

type WaitGroup struct {
	wg      *sync.WaitGroup
	wgMutex sync.Mutex
	count   int32
}

// type Evaluator struct {
// 	categoryName string
// 	controlId    string
// 	// firstLevel describes if the controlId is a first level controlId (e.g., OPS13) or a second level controlId (e.g., OPS-13.1)
// 	firstLevel bool
// 	// subControlList is a list of all sub-controls included in the first level control
// 	subControlList []string
// }

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

// StartEvaluation is a method implementation of the evaluation interface: It starts the evaluation for a cloud service and a given Target of Evaluation to be evaluated (e.g., EUCS OPS-13.2)
func (s *Service) StartEvaluation(_ context.Context, req *evaluation.StartEvaluationRequest) (resp *evaluation.StartEvaluationResponse, err error) {
	var (
		schedulerTag string
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

	// Get orchestrator client
	err = s.initOrchestratorClient()
	if err != nil {
		err = fmt.Errorf("could not set orchestrator client: %v", err)
		log.Error(err)
		resp = &evaluation.StartEvaluationResponse{
			Status:        false,
			StatusMessage: err.Error(),
		}
		return resp, status.Errorf(codes.Internal, "%s", err)
	}

	log.Info("Starting evaluation ...")

	// Scheduler tags must be unique to find the jobs by tag name
	s.scheduler.TagsUnique()

	// Add the controlsInScope of the TargetOfEvaluation including their sub-controls to the scheduler
	// TODO(anatheka): We should get the control from the orchestrator, as we do not know if the request is correct
	for _, c := range req.TargetOfEvaluation.GetControlsInScope() {
		controls := []*orchestrator.Control{}
		parentSchedulerTag := ""

		// The schedulerTag of the current control
		schedulerTag = createSchedulerTag(req.TargetOfEvaluation.GetCloudServiceId(), c.GetId())

		// Check if scheduler job for current controlId is already running
		// TODO(anatheka): Do we want to append the errors for already running jobs or do we ignore the error value?
		jobs, err := s.scheduler.FindJobsByTag(schedulerTag)
		if len(jobs) > 0 && err == nil {
			err = fmt.Errorf("evaluation for cloud service id '%s' and control id '%s' already started", req.TargetOfEvaluation.GetCloudServiceId(), c.GetId())
			log.Error(err)
			resp = &evaluation.StartEvaluationResponse{
				Status:        false,
				StatusMessage: err.Error(),
			}
			return resp, status.Errorf(codes.AlreadyExists, "%s", err)
		}

		// Add control including sub-controls as jobs to the scheduler
		controls = append(controls, c)
		controls = append(controls, c.GetControls()...)

		// Add number of sub-controls to the WaitGroup. For the first level control (e.g., OPS-13) we have to wait until all the sub-controls (e.g., OPS-13.1) are ready.
		// Control is a first level control if no parentControlId exists
		if c.ParentControlId == nil {
			s.wg[schedulerTag] = &WaitGroup{
				wg:      &sync.WaitGroup{},
				wgMutex: sync.Mutex{},
				count:   0,
			}
			s.wg[schedulerTag].wgMutex.Lock()
			s.wg[schedulerTag].wg.Add(len(c.GetControls()))
			s.wg[schedulerTag].wgMutex.Unlock()

			//TODO(anatheka): Assumption: a single sublevel control has no parentControlId
			// parentSchedulerTag is the tag for the controls parent
			parentSchedulerTag = createSchedulerTag(req.TargetOfEvaluation.GetCloudServiceId(), c.GetId())
		}

		// Add control including sub-controls to the scheduler
		for _, control := range controls {
			resp, err = s.addJobToScheduler(control, req.GetTargetOfEvaluation(), parentSchedulerTag)
			// We can return the error as it is
			if err != nil {
				return resp, err
			}
		}

		log.Infof("Evaluate Cloud Service '%s' for Control ID '%s' every 5 minutes...", req.GetTargetOfEvaluation().GetCloudServiceId(), c.GetId())
	}

	// Start scheduler jobs
	s.scheduler.StartAsync()
	log.Infof("Scheduler started...")

	resp = &evaluation.StartEvaluationResponse{Status: true}

	return
}

// addJobToScheduler adds a job for the given control to the scheduler
func (s *Service) addJobToScheduler(c *orchestrator.Control, toe *orchestrator.TargetOfEvaluation, parentSchedulerTag string) (resp *evaluation.StartEvaluationResponse, err error) {
	// schedulerTag is the tag for the given control
	schedulerTag := createSchedulerTag(toe.GetCloudServiceId(), c.GetId())

	// Regarding the control level the specific method is called. We have to decide if the control is sub-control that can be evaluated direclty or a first level control that has to wait for the results of the sub-level controls.
	if c.ParentControlId == nil { // first level control
		_, err = s.scheduler.
			Every(1).
			Minute().
			Tag(schedulerTag).
			Do(s.evaluateFirstLevelControl, toe, c.GetCategoryName(), c.GetId(), schedulerTag, c.GetControls())
	} else { // second level control
		_, err = s.scheduler.
			Every(1).
			Minute().
			Tag(schedulerTag).
			Do(s.evaluateSecondLevelControl, toe, c.GetCategoryName(), c.GetId(), parentSchedulerTag)
	}
	if err != nil {
		shortErr := fmt.Sprintf("evaluation for Cloud Service '%s' and Control ID '%s' cannot be scheduled", toe.GetCloudServiceId(), c.GetId())
		err = fmt.Errorf("%v: %v", shortErr, err)
		log.Error(err)
		resp = &evaluation.StartEvaluationResponse{
			Status:        false,
			StatusMessage: err.Error(),
		}
		return resp, status.Errorf(codes.Internal, "%s", shortErr)
	}

	log.Debugf("Cloud Service '%s' with Control ID '%s' added to scheduler", toe.GetCloudServiceId(), c.GetId())
	resp = &evaluation.StartEvaluationResponse{}

	return
}

// StopEvaluation is a method implementation of the evaluation interface: It stop the evaluation for a cloud service and its given control (e.g., EUCS OPS-13,OPS-13.2). Only first level controls (e.g., OPS-13) or individually started sub-controls (e.g., OPS-13.1) can be stopped. A sub-control with a running evaluation for the parent control cannot be stopped as it is needed for the parent control.
func (s *Service) StopEvaluation(_ context.Context, req *evaluation.StopEvaluationRequest) (res *evaluation.StopEvaluationResponse, err error) {

	var (
		schedulerTag string
		control      *orchestrator.Control
	)

	err = req.Validate()
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

	// We can only stop the evaluation for first level controls or sub-controls that are started individually
	// Check if the control is a first level control
	if control.ParentControlId == nil {
		// schedulerTag = createSchedulerTag(req.TargetOfEvaluation.CloudServiceId, control.GetId())

		// Control is a first level control, stop the scheduler job for the control and all sub-controls
		err = s.stopSchedulerJobs(getSchedulerTagsForControlIds(getAllControlIdsFromControl(control), req.GetTargetOfEvaluation().GetCloudServiceId())) //createSchedulerTag(req.GetTargetOfEvaluation().GetCloudServiceId(), control.GetId()))
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

		err = s.handleJobError(err, req.TargetOfEvaluation.GetCloudServiceId(), req.GetControlId())
		// No need of further error handling
		if err != nil {
			return
		}
	}

	res = &evaluation.StopEvaluationResponse{}

	return
}

// handleJobError handles the error from scheduler.FindJobsByTag()
// TODO(anatheka): Refactor error handling
func (s *Service) handleJobError(err error, cloudServiceId, controlId string) error {
	var schedulerTag string

	if err == nil {
		// Scheduler job for control id cannot be removed because parent control is currently evaluated
		err = fmt.Errorf("evaluation of control id '%s' for cloud service '%s' can not be stopped because the control is a sub-control of the evaluated control '%s'", controlId, cloudServiceId, controlId)
		log.Error(err)
		err = status.Errorf(codes.NotFound, err.Error())
	} else if strings.Contains(err.Error(), "no jobs found with given tag") {
		// Scheduler job can be removed because the parent control is currently not evaluated
		schedulerTag = createSchedulerTag(cloudServiceId, controlId)

		// Verify if scheduler job exists for the control id
		_, err = s.scheduler.FindJobsByTag(schedulerTag)
		if err == nil {
			err = s.stopSchedulerJob(schedulerTag)
			if err != nil {
				err = fmt.Errorf("error when stopping scheduler job for cloud service id '%s' with control id '%s'", cloudServiceId, controlId)
				log.Error(err)
				err = status.Errorf(codes.Internal, err.Error())
				return err
			}
			log.Debugf("Scheduler job for cloud service id '%s' with control id '%s'.", cloudServiceId, controlId)
		} else if strings.Contains(err.Error(), "no jobs found with given tag") {
			err = fmt.Errorf("evaluation for cloud service id '%s' with '%s' not running", cloudServiceId, controlId)
			log.Error(err)
			err = status.Errorf(codes.NotFound, err.Error())
			return err
		} else if err != nil {
			shortErr := fmt.Errorf("error when stopping scheduler job for cloud service id '%s' with '%s'", cloudServiceId, controlId)
			log.Errorf("%s: %v", shortErr, err)
			err = status.Errorf(codes.Internal, "%v", &shortErr)
			return err
		}
	}

	return err
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

// stopSchedulerJobs stops a scheduler job for the given scheduler tag
func (s *Service) stopSchedulerJobs(controlIds []string) (err error) {
	// Delete the job from the scheduler
	for _, schedulerTag := range controlIds {
		err = s.stopSchedulerJob(schedulerTag)
		if err != nil {
			return
		}
	}

	return
}

// evaluateSecondLevelControl evaluates the second level controls, e.g., OPS-13.2
func (s *Service) evaluateSecondLevelControl(toe *orchestrator.TargetOfEvaluation, categoryName, controlId, parentSchedulerTag string) {
	log.Infof("Start evaluation for Cloud Service '%s', Catalog ID '%s' and Control '%s'", toe.CloudServiceId, toe.CatalogId, controlId)
	log.Infof("Debug 10: %s", controlId)

	var (
		schedulerTag string
		result       *evaluation.EvaluationResult
		err          error
	)

	schedulerTag = createSchedulerTag(toe.GetCloudServiceId(), controlId)

	result, err = s.evaluationResultForLowerControlLevel(toe.GetCloudServiceId(), toe.GetCatalogId(), categoryName, controlId, schedulerTag, toe)
	if err != nil {
		err = fmt.Errorf("error creating evaluation result: %v", err)
		log.Error(err)
		log.Infof("Debug Done 11: %s", controlId)
	} else if result == nil {
		log.Debug("No evaluation result created")
		log.Infof("Debug Done 13: %s", controlId)
	} else {
		log.Infof("Debug Done 15: %s", controlId)

		s.resultMutex.Lock()
		s.results[result.Id] = result
		s.resultMutex.Unlock()

		log.Infof("Evaluation result for Cloud Service '%s' and ControlID '%s' stored with ID '%s'.", toe.GetCloudServiceId(), controlId, result.Id)
	}

	if parentSchedulerTag != "" {
		s.wg[parentSchedulerTag].wgMutex.Lock()
		s.wg[parentSchedulerTag].count = +1
		s.wg[parentSchedulerTag].wg.Done()
		s.wg[parentSchedulerTag].wgMutex.Unlock()
	}
	log.Infof("Debug Done 16: %s", controlId)
}

// evaluateFirstLevelControl evaluates a first level control, e.g., OPS-13
// TODO(all): Note: That is a first try. I'm not convinced, but I can't think of anything better at the moment.
func (s *Service) evaluateFirstLevelControl(toe *orchestrator.TargetOfEvaluation, categoryName, controlId, schedulerTag string, subControls []*orchestrator.Control) {
	var (
		status                        = evaluation.EvaluationResult_STATUS_UNSPECIFIED
		nonCompliantAssessmentResults []*assessment.AssessmentResult
		result                        *evaluation.EvaluationResult
	)

	// TODO(anatheka): The OPS-13 eval result is still missing! WEITERMACHEN!!!
	// TODO(anatheka): TBD How should we evaluate the first level control(OPS-13)? It should be compliant if all sub-controls are compliant
	// Suggestion: We get all sub-control results and calculate it. But then we should start that job a bit later. Then we must start the scheduler in SingletonMode
	log.Infof("Debug 1: %s", controlId)

	// Wait till all sub-controls are evaluated
	s.wg[schedulerTag].wgMutex.Lock()
	s.wg[schedulerTag].wg.Wait()
	s.wg[schedulerTag].wgMutex.Unlock()

	evaluations, err := s.ListEvaluationResults(context.Background(), &evaluation.ListEvaluationResultsRequest{
		FilteredCloudServiceId: toe.CloudServiceId,
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

	log.Infof("Debug 2: %s", controlId)
	s.resultMutex.Lock()
	s.results[result.Id] = result
	s.resultMutex.Unlock()
	log.Infof("Evaluation result for ControlID '%s' with ID '%s' stored.", controlId, result.Id)
	log.Infof("Debug 3: %s", controlId)
}

// TODO(anatheka): Move to internals
// controlContains return true if the given control id is a sub-control of the controls slice
func controlContains(controls []*orchestrator.Control, controlId string) bool {
	for _, c := range controls {
		if c.GetId() == controlId {
			return true
		}
	}

	return false
}

// ListEvaluationResults is a method implementation of the assessment interface
func (s *Service) ListEvaluationResults(_ context.Context, req *evaluation.ListEvaluationResultsRequest) (res *evaluation.ListEvaluationResultsResponse, err error) {
	res = new(evaluation.ListEvaluationResultsResponse)

	// Paginate the results according to the request
	res.Results, res.NextPageToken, err = service.PaginateMapValues(req, s.results, func(a *evaluation.EvaluationResult, b *evaluation.EvaluationResult) bool {
		return a.Id < b.Id
	}, service.DefaultPaginationOpts)
	if err != nil {
		err = fmt.Errorf("could not paginate evaluation results: %v", err)
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return
}

// evaluationResultForLowerControlLevel return the evaluation result for the lower control level, e.g. OPS-13.7
// TODO(anatheka): Refactor cloudServiceId, etc. with toe.*
func (s *Service) evaluationResultForLowerControlLevel(cloudServiceId, catalogId, categoryName, controlId, schedulerTag string, toe *orchestrator.TargetOfEvaluation) (result *evaluation.EvaluationResult, err error) {
	// Get metrics from control and sub-controls
	metrics, err := s.getMetrics(catalogId, categoryName, controlId)

	if err != nil {
		err = fmt.Errorf("could not get metrics from control and sub-controls for Cloud Serivce '%s' from Orchestrator: %v", cloudServiceId, err)
		log.Error(err)
		// s.wgMutex.Lock()
		// log.Infof("Debug Done: %s", controlId)
		// s.wg[schedulerTag].Done()
		// s.wgMutex.Unlock()

		return
	}

	// Get assessment results for the target cloud service
	// TODO(anatheka): Add FilterMetricId when PR#877 is ready
	assessmentResults, err := api.ListAllPaginated(&assessment.ListAssessmentResultsRequest{
		FilteredCloudServiceId: &cloudServiceId}, s.orchestratorClient.ListAssessmentResults, func(res *assessment.ListAssessmentResultsResponse) []*assessment.AssessmentResult {
		return res.Results
	})

	if err != nil || len(assessmentResults) == 0 {
		// TODO(anatheka): Let we the scheduler running or do we want to stop it if we do not get assessment results from the orchestrator?
		// We let the scheduler running if we do not get the assessment results from the orchestrator, maybe it is only a temporary network problem
		err = fmt.Errorf("could not get assessment results for Cloud Service ID '%s' from Orchestrator: %v", cloudServiceId, err)
		// log.Error(err)
		// s.wgMutex.Lock()
		// log.Infof("Debug Done: %s", controlId)
		// s.wg[schedulerTag].Done()
		// s.wgMutex.Unlock()

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
		ControlId:                controlId,
		TargetOfEvaluation:       toe,
		Status:                   status,
		FailingAssessmentResults: nonCompliantAssessmentResults,
	}

	// s.wgMutex.Lock()
	// log.Infof("Debug Done: %s", controlId)
	// s.wg[schedulerTag].Done()
	// s.wgMutex.Unlock()

	return
}

// getMetrics returns all  metrics from a given controlId
// For now a control has either sub-controls or metrics. If the control has sub-controls, get all metrics from the sub-controls.
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

// getSchedulerTagsForControlIds return for a given list of control ids the corresponding scheduler tags
func getSchedulerTagsForControlIds(controlIds []string, cloudServiceId string) (schedulerTags []string) {
	for _, controlId := range controlIds {
		schedulerTags = append(schedulerTags, createSchedulerTag(cloudServiceId, controlId))
	}

	return
}
