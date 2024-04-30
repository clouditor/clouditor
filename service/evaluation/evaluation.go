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
	"slices"
	"strings"
	"time"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/evaluation"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/internal/util"
	"clouditor.io/clouditor/v2/launcher"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/persistence/inmemory"
	"clouditor.io/clouditor/v2/service"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/sync/errgroup"
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

// DefaultServiceSpec returns a [launcher.ServiceSpec] for this [Service] with all necessary options retrieved from the
// config system.
func DefaultServiceSpec() launcher.ServiceSpec {
	return launcher.NewServiceSpec(
		NewService,
		WithStorage,
		nil,
		WithOAuth2Authorizer(config.ClientCredentials()),
		WithOrchestratorAddress(viper.GetString(config.OrchestratorURLFlag)),
	)
}

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

	scheduler *gocron.Scheduler

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
func WithStorage(storage persistence.Storage) service.Option[*Service] {
	return func(svc *Service) {
		svc.storage = storage
	}
}

// WithOAuth2Authorizer is an option to use an OAuth 2.0 authorizer
func WithOAuth2Authorizer(config *clientcredentials.Config) service.Option[*Service] {
	return func(svc *Service) {
		svc.orchestrator.SetAuthorizer(api.NewOAuthAuthorizerFromClientCredentials(config))
	}
}

// WithAuthorizer is an option to use a pre-created authorizer
func WithAuthorizer(auth api.Authorizer) service.Option[*Service] {
	return func(svc *Service) {
		svc.orchestrator.SetAuthorizer(auth)
	}
}

// WithOrchestratorAddress is an option to configure the orchestrator service gRPC address.
func WithOrchestratorAddress(target string, opts ...grpc.DialOption) service.Option[*Service] {
	return func(svc *Service) {
		log.Infof("Orchestrator URL is set to %s", target)

		svc.orchestrator.Target = target
		svc.orchestrator.Opts = opts
	}
}

// NewService creates a new Evaluation service
func NewService(opts ...service.Option[*Service]) *Service {
	var err error
	svc := Service{
		orchestrator:    api.NewRPCConnection(DefaultOrchestratorAddress, orchestrator.NewOrchestratorClient),
		scheduler:       gocron.NewScheduler(time.Local),
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

func (svc *Service) Init() {}

func (svc *Service) Shutdown() {
	svc.scheduler.Stop()
}

// StartEvaluation is a method implementation of the evaluation interface: It periodically starts the evaluation of a
// cloud service and the given catalog in the target_of_evaluation. If no interval time is given, the default value is
// used.
func (svc *Service) StartEvaluation(ctx context.Context, req *evaluation.StartEvaluationRequest) (resp *evaluation.StartEvaluationResponse, err error) {
	var (
		interval int
		toe      *orchestrator.TargetOfEvaluation
		catalog  *orchestrator.Catalog
		jobs     []*gocron.Job
	)

	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the cloud service according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessCreate, req) {
		return nil, service.ErrPermissionDenied
	}

	// Make sure that the scheduler is already running
	svc.scheduler.StartAsync()

	// Set the interval to the default value if not set. If the interval is set to 0, the default interval is used.
	if req.GetInterval() == 0 {
		interval = defaultInterval
	} else {
		interval = int(req.GetInterval())
	}

	// Get all Controls from Orchestrator for the evaluation
	err = svc.cacheControls(req.CatalogId)
	if err != nil {
		err = fmt.Errorf("could not cache controls: %w", err)
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	// Get Target of Evaluation
	toe, err = svc.orchestrator.Client.GetTargetOfEvaluation(context.Background(), &orchestrator.GetTargetOfEvaluationRequest{
		CloudServiceId: req.CloudServiceId,
		CatalogId:      req.CatalogId,
	})
	if err != nil {
		err = fmt.Errorf("could not get target of evaluation: %w", err)
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	// Retrieve the catalog
	catalog, err = svc.orchestrator.Client.GetCatalog(context.Background(), &orchestrator.GetCatalogRequest{
		CatalogId: req.CatalogId,
	})
	if err != nil {
		err = fmt.Errorf("could not get catalog: %w", err)
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	// Check, if a previous job exists and/or is running
	jobs, err = svc.scheduler.FindJobsByTag(req.CloudServiceId, req.CatalogId)
	if err != nil && !errors.Is(err, gocron.ErrJobNotFoundWithTag) {
		err = fmt.Errorf("error while retrieving existing scheduler job: %w", err)
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "%s", err)
	} else if len(jobs) > 0 {
		err = fmt.Errorf("evaluation for Cloud Service '%s' and Catalog ID '%s' already started", toe.GetCloudServiceId(), toe.GetCatalogId())
		log.Error(err)
		return nil, status.Errorf(codes.AlreadyExists, "%s", err)
	}

	log.Info("Starting evaluation ...")

	// Add job to scheduler
	err = svc.addJobToScheduler(ctx, toe, catalog, interval)
	// We can return the error as it is
	if err != nil {
		return nil, err
	}

	log.Infof("Scheduled to evaluate catalog ID '%s' for cloud service '%s' every %d minutes...",
		catalog.GetId(),
		toe.GetCloudServiceId(),
		interval,
	)

	resp = &evaluation.StartEvaluationResponse{Successful: true}

	return
}

// StopEvaluation is a method implementation of the evaluation interface: It stops the evaluation for a
// TargetOfEvaluation.
func (svc *Service) StopEvaluation(ctx context.Context, req *evaluation.StopEvaluationRequest) (resp *evaluation.StopEvaluationResponse, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the cloud service according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessCreate, req) {
		return nil, service.ErrPermissionDenied
	}

	// Stop jobs(s) for given cloud service and catalog
	err = svc.scheduler.RemoveByTags(req.CloudServiceId, req.CatalogId)
	if err != nil && errors.Is(err, gocron.ErrJobNotFoundWithTag) {
		return nil, status.Errorf(codes.FailedPrecondition, "job for cloud service '%s' and catalog '%s' not running", req.CloudServiceId, req.CatalogId)
	} else if err != nil {
		err = fmt.Errorf("error while removing jobs for cloud service '%s' and catalog '%s': %w", req.CloudServiceId, req.CatalogId, err)
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "%s", err)
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
	err = api.Validate(req)
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

		if util.Deref(req.Filter.ValidManualOnly) {
			query = append(query, "status IN ?")
			args = append(args, []any{
				evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT_MANUALLY,
				evaluation.EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT_MANUALLY,
			})

			query = append(query, "valid_until IS NULL OR valid_until >= CURRENT_TIMESTAMP")
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

// CreateEvaluationResult is a method implementation of the assessment interface
func (svc *Service) CreateEvaluationResult(ctx context.Context, req *evaluation.CreateEvaluationResultRequest) (res *evaluation.EvaluationResult, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// We only allow manually created statuses
	if req.Result.Status != evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT_MANUALLY &&
		req.Result.Status != evaluation.EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT_MANUALLY {
		return nil, status.Errorf(codes.InvalidArgument, "only manually set statuses are allowed")
	}

	// The ValidUntil field must be checked separately as it is an optional field and not checked by the request
	// validation. It is only mandatory when manually creating a result.
	if req.Result.ValidUntil == nil {
		return nil, status.Errorf(codes.InvalidArgument, "validity must be set")
	}

	// Check, if this request has access to the cloud service according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessCreate, req) {
		return nil, service.ErrPermissionDenied
	}

	res = req.Result
	res.Id = uuid.NewString()
	err = svc.storage.Create(res)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	return res, nil
}

// addJobToScheduler adds a job for the given control to the scheduler and sets the scheduler interval to the given interval
func (svc *Service) addJobToScheduler(ctx context.Context, toe *orchestrator.TargetOfEvaluation, catalog *orchestrator.Catalog, interval int) (err error) {
	// Check inputs and log error
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

	_, err = svc.scheduler.
		Every(interval).
		Minute().
		Tag(toe.GetCloudServiceId(), toe.GetCatalogId()).
		Do(svc.evaluateCatalog, ctx, toe, catalog, interval)
	if err != nil {
		err = fmt.Errorf("evaluation for cloud service '%s' and catalog ID '%s' cannot be scheduled: %w", toe.GetCloudServiceId(), catalog.GetId(), err)
		log.Error(err)
		return status.Errorf(codes.Internal, "%s", err)
	}

	log.Debugf("cloud service '%s' with catalog ID '%s' added to scheduler", toe.GetCloudServiceId(), catalog.GetId())

	return
}

// evaluateCatalog evaluates all [orchestrator.Control] items in the catalog whether their associated metrics are
// fulfilled or not.
func (svc *Service) evaluateCatalog(ctx context.Context, toe *orchestrator.TargetOfEvaluation, catalog *orchestrator.Catalog, interval int) error {
	var (
		controls []*orchestrator.Control
		relevant []*orchestrator.Control
		ignored  []string
		manual   map[string][]*evaluation.EvaluationResult
		err      error
		g        *errgroup.Group
		cancel   context.CancelFunc
	)

	// Retrieve all controls that match our assurance level, sorted by the control ID for easier debugging
	controls = values(svc.catalogControls[toe.CatalogId])
	slices.SortFunc(controls, func(a *orchestrator.Control, b *orchestrator.Control) int {
		return strings.Compare(a.Id, b.Id)
	})

	// First, look for any manual evaluation results that are still within their validity period, to see whether we need
	// to ignore some of the automated ones
	//
	// TODO(oxisto): Its problematic to use the context from the original StartEvaluation request, since this token
	// might time out at some point
	results, err := api.ListAllPaginated(&evaluation.ListEvaluationResultsRequest{
		Filter: &evaluation.ListEvaluationResultsRequest_Filter{
			CloudServiceId:  &toe.CloudServiceId,
			CatalogId:       &toe.CatalogId,
			ValidManualOnly: util.Ref(true),
		},
		LatestByControlId: util.Ref(true),
	},
		func(ctx context.Context, req *evaluation.ListEvaluationResultsRequest, opts ...grpc.CallOption) (*evaluation.ListEvaluationResultsResponse, error) {
			return svc.ListEvaluationResults(ctx, req)
		}, func(res *evaluation.ListEvaluationResultsResponse) []*evaluation.EvaluationResult {
			return res.Results
		})
	if err != nil {
		log.Error(err)
		return err
	}

	manual = make(map[string][]*evaluation.EvaluationResult)

	// Gather a list of controls, we are ignoring
	ignored = make([]string, 0, len(results))
	for _, result := range results {
		if result.ParentControlId != nil {
			manual[*result.ParentControlId] = append(manual[*result.ParentControlId], result)
		} else {
			ignored = append(ignored, result.ControlId)
		}
	}

	// Filter relevant controls
	for _, c := range controls {
		// Only parent controls
		if c.ParentControlId != nil {
			continue
		}

		// If we ignore the control, we can skip it
		if slices.Contains(ignored, c.Id) {
			continue
		}

		if c.IsRelevantFor(toe, catalog) {
			relevant = append(relevant, c)
		}
	}

	log.Infof("Starting catalog evaluation for Cloud Service '%s', Catalog ID '%s'. Waiting for the evaluation of %d control(s)",
		toe.CloudServiceId,
		toe.CatalogId,
		len(relevant),
	)

	// We are using a timeout of half the interval, so that we are not running into overlapping executions
	ctx, cancel = context.WithTimeout(ctx, time.Duration(interval)*time.Minute/2)
	defer cancel()

	g, ctx = errgroup.WithContext(ctx)
	for _, control := range relevant {
		control := control // https://golang.org/doc/faq#closures_and_goroutines needed until Go 1.22 (loopvar)
		g.Go(func() error {
			err := svc.evaluateControl(ctx, toe, catalog, control, manual[control.Id])
			if err != nil {
				return err
			}

			return nil
		})
	}

	// Wait until all sub-controls are evaluated
	err = g.Wait()
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

// evaluateControl evaluates a control, e.g., OPS-13. Therefore, the method needs to wait till all sub-controls (e.g.,
// OPS-13.1) are evaluated.
func (svc *Service) evaluateControl(ctx context.Context, toe *orchestrator.TargetOfEvaluation, catalog *orchestrator.Catalog, control *orchestrator.Control, manual []*evaluation.EvaluationResult) (err error) {
	var (
		status   = evaluation.EvaluationStatus_EVALUATION_STATUS_PENDING
		result   *evaluation.EvaluationResult
		results  []*evaluation.EvaluationResult
		relevant []*orchestrator.Control
		ignored  []string
		g        *errgroup.Group
	)

	// Gather a list of sub control IDs that we have manual results for and thus we are ignoring
	ignored = make([]string, 0, len(manual))
	for _, result := range manual {
		ignored = append(ignored, result.ControlId)
	}

	// Filter relevant controls
	for _, sub := range control.Controls {
		// If we ignore the control, we can skip it
		if slices.Contains(ignored, sub.Id) {
			continue
		}

		if sub.IsRelevantFor(toe, catalog) {
			relevant = append(relevant, sub)
		}
	}

	log.Infof("Starting control evaluation for Cloud Service '%s', Catalog ID '%s' and Control '%s'. Waiting for the evaluation of %d sub-control(s)",
		toe.CloudServiceId,
		toe.CatalogId,
		control.Id,
		len(relevant),
	)

	// Prepare the results slice
	results = make([]*evaluation.EvaluationResult, len(relevant)+len(manual))

	g, ctx = errgroup.WithContext(ctx)
	for i, sub := range relevant {
		i, sub := i, sub // https://golang.org/doc/faq#closures_and_goroutines needed until Go 1.22 (loopvar)
		g.Go(func() error {
			result, err := svc.evaluateSubcontrol(ctx, toe, sub)
			if err != nil {
				return err
			}

			results[i] = result
			return nil
		})
	}

	// Wait until all sub-controls are evaluated
	err = g.Wait()
	if err != nil {
		log.Error(err)
		return
	}

	// Copy the manual results
	copy(results[len(relevant):], manual)

	var nonCompliantAssessmentResults = []string{}

	for _, r := range results {
		if (r.Status == evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT ||
			r.Status == evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT_MANUALLY) &&
			(status != evaluation.EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT &&
				status != evaluation.EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT_MANUALLY) {
			status = evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT
		} else if r.Status == evaluation.EvaluationStatus_EVALUATION_STATUS_PENDING &&
			status == evaluation.EvaluationStatus_EVALUATION_STATUS_PENDING {
			status = evaluation.EvaluationStatus_EVALUATION_STATUS_PENDING
		} else {
			status = evaluation.EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT
			nonCompliantAssessmentResults = append(nonCompliantAssessmentResults, r.GetFailingAssessmentResultIds()...)
		}
	}

	// Create evaluation result
	result = &evaluation.EvaluationResult{
		Id:                         uuid.NewString(),
		Timestamp:                  timestamppb.Now(),
		ControlCategoryName:        control.CategoryName,
		ControlId:                  control.Id,
		CloudServiceId:             toe.CloudServiceId,
		ControlCatalogId:           toe.CatalogId,
		Status:                     status,
		FailingAssessmentResultIds: nonCompliantAssessmentResults,
	}

	err = svc.storage.Create(result)
	if err != nil {
		log.Errorf("error storing evaluation result for control ID '%s' (in cloud service %s) in database: %v",
			control.Id,
			toe.CloudServiceId,
			err)
		return
	}

	log.Infof("Evaluation result for control ID '%s' (in cloud service %s) was %s", control.Id, toe.CloudServiceId, result.Status.String())

	return
}

// evaluateSubcontrol evaluates the sub-controls, e.g., OPS-13.2
func (svc *Service) evaluateSubcontrol(_ context.Context, toe *orchestrator.TargetOfEvaluation, control *orchestrator.Control) (eval *evaluation.EvaluationResult, err error) {
	var (
		assessments                   []*assessment.AssessmentResult
		status                        evaluation.EvaluationStatus
		nonCompliantAssessmentResults []string
	)

	if toe == nil || control == nil {
		log.Errorf("input is missing")
		return
	}

	// Get metrics from control and sub-controls
	metrics, err := svc.getAllMetricsFromControl(toe.GetCatalogId(), control.CategoryName, control.Id)
	if err != nil {
		log.Errorf("could not get metrics for controlID '%s' and Cloud Service '%s' from Orchestrator: %v", control.Id, toe.GetCloudServiceId(), err)
		return
	}

	if len(metrics) != 0 {
		// Get latest assessment_results by resource_id filtered by
		// * cloud service id
		// * metric ids
		assessments, err = api.ListAllPaginated(&orchestrator.ListAssessmentResultsRequest{
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
		} else if len(assessments) == 0 {
			// We let the scheduler running if we do not get the assessment results from the orchestrator, maybe it is
			// only a temporary network problem
			log.Debugf("no assessment results for Cloud Service ID '%s' and MetricIds '%s' available", toe.GetCloudServiceId(), getMetricIds(metrics))
		}
	} else {
		log.Debugf("no metrics are available for the given control")
	}

	// If no assessment_results are available we are stuck at pending
	if len(assessments) == 0 {
		status = evaluation.EvaluationStatus_EVALUATION_STATUS_PENDING
	} else {
		// Otherwise, there are some results and first we assume that everything is compliant, unless someone proves it
		// otherwise
		status = evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT
	}

	// Here the actual evaluation takes place. We check if the assessment results are compliant.
	for _, results := range assessments {
		if !results.Compliant {
			nonCompliantAssessmentResults = append(nonCompliantAssessmentResults, results.GetId())
			status = evaluation.EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT
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
		return nil, err
	}

	log.Infof("Evaluation result for %s (in cloud service %s) was %s", control.Id, toe.GetCloudServiceId(), eval.Status.String())

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
		err      error
		tag      string
		controls []*orchestrator.Control
	)

	if catalogId == "" {
		return ErrCatalogIdIsMissing
	}

	// Get controls for given catalog
	controls, err = api.ListAllPaginated[*orchestrator.ListControlsResponse](&orchestrator.ListControlsRequest{
		CatalogId: catalogId,
	}, svc.orchestrator.Client.ListControls, func(res *orchestrator.ListControlsResponse) []*orchestrator.Control {
		return res.Controls
	})
	if err != nil {
		return err
	}

	if len(controls) == 0 {
		return fmt.Errorf("no controls for catalog '%s' available", catalogId)
	}

	// Store controls in map
	svc.catalogControls[catalogId] = make(map[string]*orchestrator.Control)
	for _, control := range controls {
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

// TODO(oxisto): We can remove it with maps.Values in Go 1.22+
func values[M ~map[K]V, K comparable, V any](m M) []V {
	rr := make([]V, 0, len(m))

	for _, v := range m {
		rr = append(rr, v)
	}

	return rr
}
