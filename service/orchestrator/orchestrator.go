// Copyright 2016-2022 Fraunhofer AISEC
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

package orchestrator

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"sync"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/persistence/inmemory"
	"clouditor.io/clouditor/service"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

//go:embed *.json
var f embed.FS

var defaultMetricConfigurations map[string]*assessment.MetricConfiguration
var log *logrus.Entry

var DefaultMetricsFile = "metrics.json"

var DefaultCatalogsFile = "catalogs.json"

// Service is an implementation of the Clouditor Orchestrator service
type Service struct {
	orchestrator.UnimplementedOrchestratorServer

	// cloudServiceHooks is a list of hook functions that can be used to inform
	// about updated CloudServices
	cloudServiceHooks []orchestrator.CloudServiceHookFunc

	// toeHooks is a list of hook functions that can be used to inform about updated Target of Evaluations
	toeHooks []orchestrator.TargetOfEvaluationHookFunc

	// hookMutex is used for (un)locking hook calls
	hookMutex sync.RWMutex

	// Currently only in-memory
	results map[string]*assessment.AssessmentResult

	// Hook
	AssessmentResultHooks []func(result *assessment.AssessmentResult, err error)
	// mu is used for (un)locking result hook calls
	mu sync.Mutex

	storage persistence.Storage

	metricsFile string

	// loadMetricsFunc is a function that is used to initially load metrics at the start of the orchestrator
	loadMetricsFunc func() ([]*assessment.Metric, error)

	catalogsFile string

	// loadCatalogsFunc is a function that is used to initially load catalogs at the start of the orchestrator
	loadCatalogsFunc func() ([]*orchestrator.Catalog, error)

	events chan *orchestrator.MetricChangeEvent

	// authz defines our authorization strategy, e.g., which user can access which cloud service and associated
	// resources, such as evidences and assessment results.
	authz service.AuthorizationStrategy
}

func init() {
	log = logrus.WithField("component", "orchestrator")
}

// ServiceOption is a function-style option to configure the Orchestrator Service
type ServiceOption func(*Service)

// WithMetricsFile can be used to load a different metrics file
func WithMetricsFile(file string) ServiceOption {
	return func(s *Service) {
		s.metricsFile = file
	}
}

// WithExternalMetrics can be used to load metric definitions from an external source
func WithExternalMetrics(f func() ([]*assessment.Metric, error)) ServiceOption {
	return func(s *Service) {
		s.loadMetricsFunc = f
	}
}

// WithCatalogsFile can be used to load a different catalogs file
func WithCatalogsFile(file string) ServiceOption {
	return func(s *Service) {
		s.catalogsFile = file
	}
}

// WithExternalCatalogs can be used to load catalog definitions from an external source
func WithExternalCatalogs(f func() ([]*orchestrator.Catalog, error)) ServiceOption {
	return func(s *Service) {
		s.loadCatalogsFunc = f
	}
}

// WithStorage is an option to set the storage. If not set, NewService will use inmemory storage.
func WithStorage(storage persistence.Storage) ServiceOption {
	return func(s *Service) {
		s.storage = storage
	}
}

// WithAuthorizationStrategyJWT is an option that configures an JWT-based authorization strategy using a specific claim key.
func WithAuthorizationStrategyJWT(key string) ServiceOption {
	return func(s *Service) {
		s.authz = &service.AuthorizationStrategyJWT{Key: key}
	}
}

// NewService creates a new Orchestrator service
func NewService(opts ...ServiceOption) *Service {
	var err error
	s := Service{
		results:      make(map[string]*assessment.AssessmentResult),
		metricsFile:  DefaultMetricsFile,
		catalogsFile: DefaultCatalogsFile,
		events:       make(chan *orchestrator.MetricChangeEvent, 1000),
	}

	// Apply service options
	for _, o := range opts {
		o(&s)
	}

	// Default to an in-memory storage, if nothing was explicitly set
	if s.storage == nil {
		s.storage, err = inmemory.NewStorage()
		if err != nil {
			log.Errorf("Could not initialize the storage: %v", err)
		}
	}

	// Default to an allow-all authorization strategy
	if s.authz == nil {
		s.authz = &service.AuthorizationStrategyAllowAll{}
	}

	if err = s.loadMetrics(); err != nil {
		log.Errorf("Could not load embedded metrics. Will continue with empty metric list: %v", err)
	}

	if err = s.loadCatalogs(); err != nil {
		log.Errorf("Could not load embedded catalogs: %v", err)
	}

	return &s
}

// informHooks informs the registered hook functions
func (s *Service) informHooks(ctx context.Context, cld *orchestrator.CloudService, err error) {
	s.hookMutex.RLock()
	hooks := s.cloudServiceHooks
	defer s.hookMutex.RUnlock()

	// Inform our hook, if we have any
	if len(hooks) > 0 {
		for _, hook := range hooks {
			// We could do hook concurrent again (assuming different hooks don't interfere with each other)
			hook(ctx, cld, err)
		}
	}
}

func (s *Service) RegisterCloudServiceHook(hook orchestrator.CloudServiceHookFunc) {
	s.hookMutex.Lock()
	defer s.hookMutex.Unlock()
	s.cloudServiceHooks = append(s.cloudServiceHooks, hook)
}

// StoreAssessmentResult is a method implementation of the orchestrator interface: It receives an assessment result and stores it
func (s *Service) StoreAssessmentResult(_ context.Context, req *orchestrator.StoreAssessmentResultRequest) (resp *orchestrator.StoreAssessmentResultResponse, err error) {
	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	s.results[req.Result.Id] = req.Result

	go s.informHook(req.Result, nil)

	return nil, nil
}

func (s *Service) StoreAssessmentResults(stream orchestrator.Orchestrator_StoreAssessmentResultsServer) (err error) {
	var (
		result *orchestrator.StoreAssessmentResultRequest
		res    *orchestrator.StoreAssessmentResultsResponse
	)

	for {
		result, err = stream.Recv()

		// If no more input of the stream is available, return
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			newError := fmt.Errorf("cannot receive stream request: %w", err)
			log.Error(newError)
			return status.Errorf(codes.Unknown, "%v", newError)
		}

		// Call StoreAssessmentResult() for storing a single assessment
		storeAssessmentResultReq := &orchestrator.StoreAssessmentResultRequest{
			Result: result.Result,
		}
		_, err = s.StoreAssessmentResult(context.Background(), storeAssessmentResultReq)
		if err != nil {
			// Create response message. The StoreAssessmentResult method does not need that message, so we have to create it here for the stream response.
			res = &orchestrator.StoreAssessmentResultsResponse{
				Status:        false,
				StatusMessage: err.Error(),
			}
		} else {
			res = &orchestrator.StoreAssessmentResultsResponse{
				Status: true,
			}
		}

		log.Debugf("Assessment result received (%v)", result.Result.Id)

		err = stream.Send(res)

		// Check for send errors
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			newError := fmt.Errorf("cannot stream response to the client: %w", err)
			log.Error(newError)
			return status.Errorf(codes.Unknown, "%v", newError.Error())
		}
	}
}

func (s *Service) RegisterAssessmentResultHook(hook func(result *assessment.AssessmentResult, err error)) {
	s.hookMutex.Lock()
	defer s.hookMutex.Unlock()
	s.AssessmentResultHooks = append(s.AssessmentResultHooks, hook)
}

func (s *Service) informHook(result *assessment.AssessmentResult, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Inform our hook, if we have any
	if s.AssessmentResultHooks != nil {
		for _, hook := range s.AssessmentResultHooks {
			// TODO(all): We could do hook concurrent again (assuming different hooks don't interfere with each other)
			hook(result, err)
		}
	}
}

// CreateCertificate implements method for creating a new certificate
func (svc *Service) CreateCertificate(_ context.Context, req *orchestrator.CreateCertificateRequest) (
	*orchestrator.Certificate, error) {
	// Validate request
	err := service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	// Persist the new certificate in our database
	err = svc.storage.Create(req.Certificate)
	if err != nil {
		return nil,
			status.Errorf(codes.Internal, "could not add certificate to the database: %v", err)
	}

	// Return certificate
	return req.Certificate, nil
}

// GetCertificate implements method for getting a certificate, e.g. to show its state in the UI
func (svc *Service) GetCertificate(_ context.Context, req *orchestrator.GetCertificateRequest) (response *orchestrator.Certificate, err error) {
	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	response = new(orchestrator.Certificate)
	err = svc.storage.Get(response, "Id = ?", req.CertificateId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "certificate not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}
	return response, nil
}

// ListCertificates implements method for getting a certificate, e.g. to show its state in the UI
func (svc *Service) ListCertificates(_ context.Context, req *orchestrator.ListCertificatesRequest) (res *orchestrator.ListCertificatesResponse, err error) {
	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	res = new(orchestrator.ListCertificatesResponse)

	res.Certificates, res.NextPageToken, err = service.PaginateStorage[*orchestrator.Certificate](req, svc.storage,
		service.DefaultPaginationOpts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
	}
	return
}

// UpdateCertificate implements method for updating an existing certificate
func (svc *Service) UpdateCertificate(_ context.Context, req *orchestrator.UpdateCertificateRequest) (response *orchestrator.Certificate, err error) {
	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	count, err := svc.storage.Count(req.Certificate, req.Certificate.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	if count == 0 {
		return nil, status.Error(codes.NotFound, "certificate not found")
	}

	response = req.Certificate

	err = svc.storage.Save(response, "Id = ?", response.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}
	return
}

// RemoveCertificate implements method for removing a certificate
func (svc *Service) RemoveCertificate(_ context.Context, req *orchestrator.RemoveCertificateRequest) (response *emptypb.Empty, err error) {
	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	err = svc.storage.Delete(&orchestrator.Certificate{}, "Id = ?", req.CertificateId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "certificate not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	return &emptypb.Empty{}, nil
}
