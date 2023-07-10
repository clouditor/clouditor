// Copyright 2023 Fraunhofer AISEC
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
	"clouditor.io/clouditor/persistence/inmemory"
	"clouditor.io/clouditor/service"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	log                      *logrus.Entry
	ErrCatalogIdIsMissing    = errors.New("catalog_id is missing")
	ErrCategoryNameIsMissing = errors.New("category_name is missing")
	ErrControlIdIsMissing    = errors.New("control_id is missing")
	ErrControlNotAvailable   = errors.New("control not available")
)

const (
	// DefaultOrchestratorAddress specifies the default gRPC address of the orchestrator service.
	DefaultOrchestratorAddress = "localhost:9090"

	// defaultInterval is the default interval time for the scheduler. If no interval is set in the StartEvaluationRequest, the default value is taken.
	defaultInterval int = 5
)

// Service is an implementation of the Clouditor Evaluation service
type Service struct {
	evaluation.UnimplementedEvaluationServer

	orchestrator *api.RPCConnection[orchestrator.OrchestratorClient]

	authorizer api.Authorizer
	scheduler  map[string]*gocron.Scheduler

	// wg is used for a control that waits for its sub-controls to be evaluated
	wg map[string]*sync.WaitGroup

	// authz defines our authorization strategy, e.g., which user can access which cloud service and associated
	// resources, such as evaluation results.
	authz service.AuthorizationStrategy

	storage persistence.Storage

	// controls stores the catalog controls so that they do not always have to be retrieved from Orchestrators getControl endpoint
	// map[catalog_id][category_name-control_id]*orchestrator.Control
	catalogControls map[string]map[string]*orchestrator.Control
}

func init() {
	log = logrus.WithField("component", "evaluation")
}

// WithStorage is an option to set the storage. If not set, NewService will use inmemory storage.
func WithStorage(storage persistence.Storage) service.Option[Service] {
	return func(svc *Service) {
		svc.storage = storage
	}
}

// WithOAuth2Authorizer is an option to use an OAuth 2.0 authorizer
func WithOAuth2Authorizer(config *clientcredentials.Config) service.Option[Service] {
	return func(svc *Service) {
		svc.orchestrator.SetAuthorizer(api.NewOAuthAuthorizerFromClientCredentials(config))
	}
}

// WithAuthorizer is an option to use a pre-created authorizer
func WithAuthorizer(auth api.Authorizer) service.Option[Service] {
	return func(svc *Service) {
		svc.orchestrator.SetAuthorizer(auth)
	}
}

// WithOrchestratorAddress is an option to configure the orchestrator service gRPC address.
func WithOrchestratorAddress(target string, opts ...grpc.DialOption) service.Option[Service] {
	return func(svc *Service) {
		svc.orchestrator.Target = target
		svc.orchestrator.Opts = opts
	}
}

// NewService creates a new Evaluation service
func NewService(opts ...service.Option[Service]) *Service {
	var err error
	svc := Service{
		orchestrator:    api.NewRPCConnection(DefaultOrchestratorAddress, orchestrator.NewOrchestratorClient),
		scheduler:       make(map[string]*gocron.Scheduler),
		wg:              make(map[string]*sync.WaitGroup),
		catalogControls: make(map[string]map[string]*orchestrator.Control),
	}

	// Apply service options
	for _, o := range opts {
		o(&svc)
	}

	// Default to an in-memory storage, if nothing was explicitly set
	if svc.storage == nil {
		svc.storage, err = inmemory.NewStorage()
		if err != nil {
			log.Errorf("Could not initialize the storage: %v", err)
		}
	}

	// Default to an allow-all authorization strategy
	if svc.authz == nil {
		svc.authz = &service.AuthorizationStrategyAllowAll{}
	}

	return &svc
}

// StartEvaluation is a method implementation of the evaluation interface: It periodically starts the evaluation of a
// cloud service and the given controls_in_scope (e.g., EUCS OPS-13, EUCS OPS-13.2) in the target_of_evaluation. If no
// interval time is given, the default value is used.
func (svc *Service) StartEvaluation(ctx context.Context, req *evaluation.StartEvaluationRequest) (resp *evaluation.StartEvaluationResponse, err error) {
	var (
		parentJobTag string
		interval     int
		toe          *orchestrator.TargetOfEvaluation
		schedulerTag string
	)

	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the cloud service according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessCreate, req) {
		return nil, service.ErrPermissionDenied
	}

	// Set the interval to the default value if not set. If the interval is set to 0, the default interval is used.
	if req.GetInterval() == 0 {
		interval = defaultInterval
	} else {
		interval = int(req.GetInterval())
	}

	// Get all Controls from Orchestrator for the evaluation
	err = svc.cacheControls(req.GetCatalogId())
	if err != nil {
		err = fmt.Errorf("could not cache controls: %w", err)
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	// Get Target of Evaluation. The Target of Evaluation is retrieved to get the controls_in_scope, which are then
	// evaluated.
	toe, err = svc.orchestrator.Client.GetTargetOfEvaluation(context.Background(), &orchestrator.GetTargetOfEvaluationRequest{
		CloudServiceId: req.GetCloudServiceId(),
		CatalogId:      req.GetCatalogId(),
	})
	if err != nil {
		err = fmt.Errorf("could not get target of evaluation: %w", err)
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	// schedulerTag is the tag for the scheduler that contains the scheduler jobs for the specific ToE
	schedulerTag = createSchedulerTag(req.GetCloudServiceId(), req.GetCatalogId())

	// Check if scheduler for ToE has started already
	if svc.scheduler != nil && svc.scheduler[schedulerTag] != nil && svc.scheduler[schedulerTag].IsRunning() {
		err = fmt.Errorf("evaluation for Cloud Service '%s' and Catalog ID '%s' already started", toe.GetCloudServiceId(), toe.GetCatalogId())
		log.Error(err)
		return nil, status.Errorf(codes.AlreadyExists, "%s", err)
	}

	log.Info("Starting evaluation ...")

	// Create new scheduler for ToE
	svc.scheduler[schedulerTag] = gocron.NewScheduler(time.UTC)

	// Scheduler tags must be unique to find the jobs by tag name
	svc.scheduler[schedulerTag].TagsUnique()

	// Add the controls_in_scope of the target_of_evaluation including their sub-controls to the scheduler. The parent
	// control has to wait for the evaluation of the sub-controls. That's why we need to know how much sub-controls are
	// available and define the waitGroup with the number of the corresponding sub-controls. The controls_in_scope are
	// not stored in a hierarchy, so we have to get the parent control and find all related sub-controls.
	controlsInScope := createControlsInScopeHierarchy(toe.GetControlsInScope())
	for _, control := range controlsInScope {
		var (
			controls []*orchestrator.Control
		)

		// If the control does not have sub-controls don't start a scheduler.
		if control.ParentControlId != nil {
			continue
		}

		// Collect control including sub-controls from the controls_in_scope slice to start them later as scheduler jobs
		controls = append(controls, control)
		controls = append(controls, control.GetControls()...)

		// parentJobTag is the tag for the parent control (e.g., OPS-13)
		parentJobTag = createJobTag(toe.GetCloudServiceId(), toe.GetCatalogId(), control.GetId())

		// Add number of sub-controls to the WaitGroup. For the control (e.g., OPS-13) we have to wait until all the
		// sub-controls (e.g., OPS-13.1) are ready. Control is a parent control if no parentControlId exists.
		svc.wg[parentJobTag] = &sync.WaitGroup{}
		// The controls list contains also the parent control itself and must be minimized by 1 for the
		// parent_control_id.
		svc.wg[parentJobTag].Add(len(controls) - 1)

		// Add control including sub-controls to the scheduler
		for _, control := range controls {
			// If the control is a sub-control create parentJobTag
			if control.GetParentControlId() == "" {
				parentJobTag = ""
			} else {
				// parentJobTag is the tag for the parent control (e.g., OPS-13)
				parentJobTag = createJobTag(toe.GetCloudServiceId(), toe.GetCatalogId(), control.GetParentControlId())
			}
			err = svc.addJobToScheduler(control, toe, parentJobTag, interval)
			// We can return the error as it is
			if err != nil {
				return nil, err
			}
		}

		log.Infof("Scheduled to evaluate control ID '%s' for cloud service '%s' every 5 minutes...", control.GetId(), toe.GetCloudServiceId())
	}

	// Start scheduler jobs
	svc.scheduler[schedulerTag].StartAsync()
	log.Infof("Evaluation started.")

	resp = &evaluation.StartEvaluationResponse{Successful: true}

	return
}

// StopEvaluation is a method implementation of the evaluation interface: It stops the evaluation for a
// TargetOfEvaluation.
func (svc *Service) StopEvaluation(ctx context.Context, req *evaluation.StopEvaluationRequest) (resp *evaluation.StopEvaluationResponse, err error) {
	var schedulerTag string

	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the cloud service according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessCreate, req) {
		return nil, service.ErrPermissionDenied
	}

	// schedulerTag is the tag for the scheduler that contains the scheduler jobs for the specific ToE
	schedulerTag = createSchedulerTag(req.GetCloudServiceId(), req.GetCatalogId())

	// Stop scheduler for given Cloud Service
	if _, ok := svc.scheduler[schedulerTag]; ok {
		svc.scheduler[schedulerTag].Stop()
	} else {
		err = fmt.Errorf("evaluation for Cloud Service '%s' and Catalog '%s' not running", req.GetCloudServiceId(), req.GetCatalogId())
		log.Error(err)
		return nil, status.Errorf(codes.NotFound, "%v", err)
	}

	if !svc.scheduler[schedulerTag].IsRunning() {
		delete(svc.scheduler, schedulerTag)
		log.Infof("Scheduler for Cloud Service '%s' and Catalog '%s' deleted.", req.GetCloudServiceId(), req.GetCatalogId())
	} else {
		log.Errorf("Error deleting scheduler for Cloud Service '%s' and Catalog '%s' deleted.", req.GetCloudServiceId(), req.GetCatalogId())
	}

	resp = &evaluation.StopEvaluationResponse{}

	return
}

// ListEvaluationResults is a method implementation of the assessment interface
func (svc *Service) ListEvaluationResults(ctx context.Context, req *evaluation.ListEvaluationResultsRequest) (res *evaluation.ListEvaluationResultsResponse, err error) {
	var (
		// filtered_values []*evaluation.EvaluationResult
		allowed   []string
		all       bool
		query     []string
		partition []string
		args      []any
	)

	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	// Retrieve list of allowed cloud service according to our authorization strategy. No need to specify any conditions
	// to our storage request, if we are allowed to see all cloud services.
	all, allowed = svc.authz.AllowedCloudServices(ctx)

	// Filtering evaluation results by
	// * cloud service ID
	// * control ID
	// * sub-controls
	if req.Filter != nil {
		// Check if cloud_service_id in filter is within allowed or one can access *all* the cloud services
		if !svc.authz.CheckAccess(ctx, service.AccessRead, req.Filter) {
			return nil, service.ErrPermissionDenied
		}

		if req.Filter.CloudServiceId != nil {
			query = append(query, "cloud_service_id = ?")
			args = append(args, req.Filter.GetCloudServiceId())
		}

		if req.Filter.CatalogId != nil {
			query = append(query, "control_catalog_id = ?")
			args = append(args, req.Filter.GetCatalogId())
		}

		if req.Filter.ControlId != nil {
			query = append(query, "control_id = ?")
			args = append(args, req.Filter.GetControlId())
		}

		// TODO(anatheka): change that, in other catalogs maybe it's not that easy to get the sub-control by name
		if req.Filter.SubControls != nil {
			partition = append(partition, "control_id")
			query = append(query, "control_id LIKE ?")
			args = append(args, fmt.Sprintf("%s%%", req.Filter.GetSubControls()))
		}

		if util.Deref(req.Filter.ParentsOnly) {
			query = append(query, "parent_control_id IS NULL")
		}
	}

	// In any case, we need to make sure that we only select evaluation results of cloud services that we have access to
	// (if we do not have access to all)
	if !all {
		query = append(query, "cloud_service_id IN ?")
		args = append(args, allowed)
	}

	res = new(evaluation.ListEvaluationResultsResponse)

	// If we want to have it grouped by resource ID, we need to do a raw query
	if req.GetLatestByControlId() {
		// In the raw SQL, we need to build the whole WHERE statement
		var where string
		var p = ""

		if len(query) > 0 {
			where = "WHERE " + strings.Join(query, " AND ")
		}

		if len(partition) > 0 {
			p = ", " + strings.Join(partition, ",")
		}

		// Execute the raw SQL statement
		err = svc.storage.Raw(&res.Results,
			fmt.Sprintf(`WITH sorted_results AS (
				SELECT *, ROW_NUMBER() OVER (PARTITION BY control_id %s ORDER BY timestamp DESC) AS row_number
				FROM evaluation_results
				%s
		  	)
		  	SELECT * FROM sorted_results WHERE row_number = 1 ORDER BY control_catalog_id, control_id;`, p, where), args...)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "database error: %v", err)
		}
	} else {
		// join query with AND and prepend the query
		args = append([]any{strings.Join(query, " AND ")}, args...)

		// Paginate the results according to the request
		res.Results, res.NextPageToken, err = service.PaginateStorage[*evaluation.EvaluationResult](req, svc.storage, service.DefaultPaginationOpts, args...)
		if err != nil {
			err = fmt.Errorf("could not paginate evaluation results: %w", err)
			log.Error(err)
			return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
		}
	}

	return
}

// addJobToScheduler adds a job for the given control to the scheduler and sets the scheduler interval to the given interval
func (svc *Service) addJobToScheduler(c *orchestrator.Control, toe *orchestrator.TargetOfEvaluation, parentJobTag string, interval int) (err error) {
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

	// jobTag is the tag for the given control
	jobTag := createJobTag(toe.GetCloudServiceId(), toe.GetCatalogId(), c.GetId())
	// schedulerTag is the tag for the scheduler that contains the scheduler jobs for the specific ToE
	schedulerTag := createSchedulerTag(toe.GetCloudServiceId(), toe.GetCatalogId())

	// Regarding the control level the specific method is called every X minutes based on the given interval. We have to decide if a sub-control is started individually or a parent control that has to wait for the results of the sub-controls.
	// If a parent control with its sub-controls is started, the parentJobTag is empty and the ParentControlId is not set.
	// If a sub-control is started individually the parentJobTag is empty and ParentControlId is set.
	if parentJobTag == "" && c.ParentControlId == nil { // parent control
		_, err = svc.scheduler[schedulerTag].
			Every(interval).
			Minute().
			Tag(jobTag).
			Do(svc.evaluateControl, toe, c.GetCategoryName(), c.GetId(), jobTag, c.GetControls())
	} else { // sub-control
		_, err = svc.scheduler[schedulerTag].
			Every(interval).
			Minute().
			Tag(jobTag).
			Do(svc.evaluateSubcontrol, toe, c, parentJobTag)
	}
	if err != nil {
		err = fmt.Errorf("evaluation for Cloud Service '%s' and Control ID '%s' cannot be scheduled: %w", toe.GetCloudServiceId(), c.GetId(), err)
		log.Error(err)
		return status.Errorf(codes.Internal, "%s", err)
	}

	log.Debugf("Cloud Service '%s' with Control ID '%s' added to scheduler", toe.GetCloudServiceId(), c.GetId())

	return
}

// evaluateControl evaluates a control, e.g., OPS-13. Therefore, the method needs to wait till all sub-controls (e.g.,
// OPS-13.1) are evaluated.
//
// TODO(all): Note: That is a first try. I'm not convinced, but I can't think of anything better at the moment.
func (svc *Service) evaluateControl(toe *orchestrator.TargetOfEvaluation, categoryName, controlId, schedulerTag string, subControls []*orchestrator.Control) {
	var (
		status     = evaluation.EvaluationStatus_EVALUATION_STATUS_PENDING
		evalResult *evaluation.EvaluationResult
	)

	log.Infof("Start control evaluation for Cloud Service '%s', Catalog ID '%s' and Control '%s'", toe.CloudServiceId, toe.CatalogId, controlId)

	// Wait till all sub-controls are evaluated
	svc.wg[schedulerTag].Wait()

	// For the next iteration set wg again to the number of sub-controls
	svc.wg[schedulerTag].Add(len(subControls))

	evaluations, err := svc.ListEvaluationResults(context.Background(), &evaluation.ListEvaluationResultsRequest{
		Filter: &evaluation.ListEvaluationResultsRequest_Filter{
			CloudServiceId: &toe.CloudServiceId,
			SubControls:    util.Ref(fmt.Sprintf("%s.", controlId)), // We only need the sub-controls for the evaluation.
		},
		LatestByControlId: util.Ref(true),
	})
	if err != nil {
		err = fmt.Errorf("error list evaluation results: %w", err)
		log.Error(err)
		return
	}

	// If no evaluation results for the sub-controls are available return and do not create a new evaluation result
	if len(evaluations.Results) == 0 {
		return
	}

	var nonCompliantAssessmentResults = []string{}

	for _, r := range evaluations.Results {
		if r.Status == evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT && status != evaluation.EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT {
			status = evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT
		} else if r.Status == evaluation.EvaluationStatus_EVALUATION_STATUS_PENDING && status == evaluation.EvaluationStatus_EVALUATION_STATUS_PENDING {
			status = evaluation.EvaluationStatus_EVALUATION_STATUS_PENDING
		} else {
			status = evaluation.EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT
			nonCompliantAssessmentResults = append(nonCompliantAssessmentResults, r.GetFailingAssessmentResultIds()...)
		}
	}

	// Create evaluation result
	evalResult = &evaluation.EvaluationResult{
		Id:                         uuid.NewString(),
		Timestamp:                  timestamppb.Now(),
		ControlCategoryName:        categoryName,
		ControlId:                  controlId,
		CloudServiceId:             toe.GetCloudServiceId(),
		ControlCatalogId:           toe.GetCatalogId(),
		Status:                     status,
		FailingAssessmentResultIds: nonCompliantAssessmentResults,
	}

	err = svc.storage.Create(evalResult)
	if err != nil {
		log.Errorf("error storing evaluation result for control ID '%s' in database: %v", controlId, err)
		return
	}

	log.Debugf("Evaluation result stored for ControlID '%s' and Cloud Service ID '%s' with ID '%s'.", controlId, toe.GetCloudServiceId(), evalResult.Id)

}

// evaluateSubcontrol evaluates the sub-controls, e.g., OPS-13.2
func (svc *Service) evaluateSubcontrol(toe *orchestrator.TargetOfEvaluation, control *orchestrator.Control, parentJobTag string) {
	var (
		eval                          *evaluation.EvaluationResult
		err                           error
		assessmentResults             []*assessment.AssessmentResult
		status                        evaluation.EvaluationStatus
		nonCompliantAssessmentResults []string
	)

	if toe == nil || control == nil || parentJobTag == "" {
		log.Errorf("input is missing")
		return
	}

	// Get metrics from control and sub-controls
	metrics, err := svc.getAllMetricsFromControl(toe.GetCatalogId(), control.CategoryName, control.Id)
	if err != nil {
		log.Errorf("could not get metrics for controlID '%s' and Cloud Service '%s' from Orchestrator: %v", control.Id, toe.GetCloudServiceId(), err)
		// If the parentJobTag is not empty, we have do decrement the WaitGroup so that the parent control can also be
		// evaluated when all sub-controls are evaluated.
		if parentJobTag != "" && svc.wg[parentJobTag] != nil {
			svc.wg[parentJobTag].Done()
		}
		return
	}

	if len(metrics) != 0 {
		// Get latest assessment_results by resource_id filtered by
		// * cloud service id
		// * metric ids
		assessmentResults, err = api.ListAllPaginated(&orchestrator.ListAssessmentResultsRequest{
			Filter: &orchestrator.Filter{
				CloudServiceId: &toe.CloudServiceId,
				MetricIds:      getMetricIds(metrics),
			},
			LatestByResourceId: util.Ref(true),
		}, svc.orchestrator.Client.ListAssessmentResults, func(res *orchestrator.ListAssessmentResultsResponse) []*assessment.AssessmentResult {
			return res.Results
		})

		if err != nil {
			// We let the scheduler running if we do not get the assessment results from the orchestrator, maybe it is
			// only a temporary network problem
			log.Errorf("could not get assessment results for Cloud Service ID '%s' and MetricIds '%s' from Orchestrator: %v", toe.GetCloudServiceId(), getMetricIds(metrics), err)
		} else if len(assessmentResults) == 0 {
			// We let the scheduler running if we do not get the assessment results from the orchestrator, maybe it is
			// only a temporary network problem
			log.Debugf("no assessment results for Cloud Service ID '%s' and MetricIds '%s' available", toe.GetCloudServiceId(), getMetricIds(metrics))
		}
	} else {
		log.Debugf("no metrics are available for the given control")
	}

	status = evaluation.EvaluationStatus_EVALUATION_STATUS_PENDING

	// Here the actual evaluation takes place. We check if the assessment results are compliant.
	// If the latest assessment result per resource_id is not compliant the whole evaluation status is set to
	// NOT_COMPLIANT. Furthermore, all non-compliant assessment_result_ids are stored in a separate list.

	// Get a map of the assessment results, so that we have all assessment results for a specific resource_id and
	// metric_id together for evaluation
	assessmentResultsMap := createAssessmentResultMap(assessmentResults)
	for _, results := range assessmentResultsMap {
		// If no assessment_results are available continue
		if len(results) == 0 {
			continue
		}

		// Otherwise, there are some results and first we assume that everything is compliant, unless someone proves it
		// otherwise
		status = evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT

		for i := range results {
			if !results[i].Compliant {
				nonCompliantAssessmentResults = append(nonCompliantAssessmentResults, results[i].GetId())
				status = evaluation.EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT
			}
		}
	}

	// Create evaluation result
	eval = &evaluation.EvaluationResult{
		Id:                         uuid.NewString(),
		Timestamp:                  timestamppb.Now(),
		ControlCategoryName:        control.CategoryName,
		ControlId:                  control.Id,
		ParentControlId:            control.ParentControlId,
		CloudServiceId:             toe.GetCloudServiceId(),
		ControlCatalogId:           toe.GetCatalogId(),
		Status:                     status,
		FailingAssessmentResultIds: nonCompliantAssessmentResults,
	}

	err = svc.storage.Create(eval)
	if err != nil {
		log.Errorf("error storing evaluation result for control ID '%s' in database: %v", control.Id, err)
	}

	log.Debugf("Evaluation result stored for ControlID '%s' and Cloud Service ID '%s' with ID '%s'.", control.Id, toe.GetCloudServiceId(), eval.Id)

	// If the parentJobTag is not empty, we have do decrement the WaitGroup so that the parent control can also be
	// evaluated when all sub-controls are evaluated.
	if parentJobTag != "" && svc.wg[parentJobTag] != nil {
		svc.wg[parentJobTag].Done()
	}
}

// getMetricIds returns the metric Ids for the given metrics
func getMetricIds(metrics []*assessment.Metric) []string {
	var metricIds []string

	for _, m := range metrics {
		metricIds = append(metricIds, m.GetId())
	}

	return metricIds
}

// getAllMetricsFromControl returns all metrics from a given controlId.
//
// For now a control has either sub-controls or metrics. If the control has sub-controls, get also all metrics from the
// sub-controls.
func (svc *Service) getAllMetricsFromControl(catalogId, categoryName, controlId string) (metrics []*assessment.Metric, err error) {
	var subControlMetrics []*assessment.Metric

	control, err := svc.getControl(catalogId, categoryName, controlId)
	if err != nil {
		err = fmt.Errorf("could not get control for control id {%s}: %w", controlId, err)
		return
	}

	// Add metric of control to the metrics list
	metrics = append(metrics, control.Metrics...)

	// Add sub-control metrics to the metric list if exist
	if len(control.Controls) != 0 {
		// Get the metrics from the next sub-control
		subControlMetrics, err = svc.getMetricsFromSubcontrols(control)
		if err != nil {
			err = fmt.Errorf("error getting metrics from sub-controls: %w", err)
			return
		}

		metrics = append(metrics, subControlMetrics...)
	}

	return
}

// getMetricsFromSubcontrols returns a list of metrics from the sub-controls.
func (svc *Service) getMetricsFromSubcontrols(control *orchestrator.Control) (metrics []*assessment.Metric, err error) {
	var subcontrol *orchestrator.Control

	if control == nil {
		return nil, errors.New("control is missing")
	}

	for _, control := range control.Controls {
		subcontrol, err = svc.getControl(control.CategoryCatalogId, control.CategoryName, control.Id)
		if err != nil {
			return
		}

		metrics = append(metrics, subcontrol.Metrics...)
	}

	return
}

// cacheControls caches the catalog controls for the given catalog.
func (svc *Service) cacheControls(catalogId string) error {
	var (
		err  error
		tag  string
		resp *orchestrator.ListControlsResponse
	)

	if catalogId == "" {
		return ErrCatalogIdIsMissing
	}

	// Get controls for given catalog
	resp, err = svc.orchestrator.Client.ListControls(context.Background(), &orchestrator.ListControlsRequest{
		CatalogId: catalogId,
	})
	if err != nil {
		return err
	}

	if len(resp.Controls) == 0 {
		return fmt.Errorf("no controls for catalog '%s' available", catalogId)
	}

	// Store controls in map
	svc.catalogControls[catalogId] = make(map[string]*orchestrator.Control)
	for _, control := range resp.Controls {
		tag = fmt.Sprintf("%s-%s", control.GetCategoryName(), control.GetId())
		svc.catalogControls[catalogId][tag] = control
	}

	return nil
}

// getControl returns the control for the given catalogID, CategoryName and controlID.
func (svc *Service) getControl(catalogId, categoryName, controlId string) (control *orchestrator.Control, err error) {
	if catalogId == "" {
		return nil, ErrCatalogIdIsMissing
	} else if categoryName == "" {
		return nil, ErrCategoryNameIsMissing
	} else if controlId == "" {
		return nil, ErrControlIdIsMissing
	}

	tag := fmt.Sprintf("%s-%s", categoryName, controlId)

	control, ok := svc.catalogControls[catalogId][tag]
	if !ok {
		return nil, ErrControlNotAvailable
	}

	return
}

// createSchedulerTag creates a tag for the schedulers job in the format 'cloud_service_id-catalog_id' ('00000000-0000-0000-0000-000000000000-EUCS')
func createSchedulerTag(cloudServiceId, catalogId string) string {
	if cloudServiceId == "" || catalogId == "" {
		return ""
	}
	return fmt.Sprintf("%s-%s", cloudServiceId, catalogId)
}

// createJobTag creates the schedulers job tag for a given catalog_id, category_name and control_id.
func createJobTag(catalogId, categoryName, controlId string) string {
	if catalogId == "" || categoryName == "" || controlId == "" {
		return ""
	}
	return fmt.Sprintf("%s-%s-%s", catalogId, categoryName, controlId)
}

// createControlsInScopeHierarchy return a controls list as hierarchy regarding the parent and sub-controls.
func createControlsInScopeHierarchy(controls []*orchestrator.Control) (controlsHierarchy []*orchestrator.Control) {
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

// createAssessmentResultMap returns a map with the resource_id as key and the assessment results as a value slice. We
// need that map if we have more than one assessment_result for evaluation, e.g., if we have two assessment_results for
// 2 different metrics.
func createAssessmentResultMap(results []*assessment.AssessmentResult) map[string][]*assessment.AssessmentResult {
	var hierarchyResults = make(map[string][]*assessment.AssessmentResult)

	for _, result := range results {
		hierarchyResults[result.GetResourceId()] = append(hierarchyResults[result.GetResourceId()], result)
	}

	return hierarchyResults
}
