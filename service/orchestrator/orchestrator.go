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
var DefaultRequirementsFile = "requirements.json"

// Service is an implementation of the Clouditor Orchestrator service
type Service struct {
	orchestrator.UnimplementedOrchestratorServer

	// metricConfigurations holds a double-map of metric configurations associated first by service ID and then metric ID
	metricConfigurations map[string]map[string]*assessment.MetricConfiguration

	// metrics contains map of our metric definitions
	metrics map[string]*assessment.Metric

	// mm is a mutex for metric related maps
	mm sync.Mutex

	// cloudServiceHooks is a list of hook functions that can be used to inform
	// about updated CloudServices
	cloudServiceHooks []orchestrator.CloudServiceHookFunc
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

	requirements []*orchestrator.Requirement

	events chan *orchestrator.MetricChangeEvent
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

func WithRequirements(r []*orchestrator.Requirement) ServiceOption {
	return func(s *Service) {
		s.requirements = r
	}
}

// WithStorage is an option to set the storage. If not set, NewService will use inmemory storage.
func WithStorage(storage persistence.Storage) ServiceOption {
	return func(s *Service) {
		s.storage = storage
	}
}

// NewService creates a new Orchestrator service
func NewService(opts ...ServiceOption) *Service {
	var err error
	s := Service{
		results:              make(map[string]*assessment.AssessmentResult),
		metricConfigurations: make(map[string]map[string]*assessment.MetricConfiguration),
		metricsFile:          DefaultMetricsFile,
		events:               make(chan *orchestrator.MetricChangeEvent, 1000),
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

	// Load requirements if nothing was specified
	if s.requirements == nil {
		if s.requirements, err = LoadRequirements(DefaultRequirementsFile); err != nil {
			log.Errorf("Could not load embedded requirements. Will continue with empty requirements list: %v", err)
		}
	}

	if err = s.loadMetrics(); err != nil {
		log.Errorf("Could not load embedded metrics. Will continue with empty metric list: %v", err)
	}

	return &s
}

// informHooks informs the registered hook functions
func (s *Service) informHooks(cld *orchestrator.CloudService, err error) {
	s.hookMutex.RLock()
	hooks := s.cloudServiceHooks
	defer s.hookMutex.RUnlock()

	// Inform our hook, if we have any
	if len(hooks) > 0 {
		for _, hook := range hooks {
			// We could do hook concurrent again (assuming different hooks don't interfere with each other)
			hook(cld, err)
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
	_, err = req.Result.Validate()

	if err != nil {
		newError := fmt.Errorf("invalid assessment result: %w", err)
		log.Error(newError)

		go s.informHook(nil, newError)

		resp = &orchestrator.StoreAssessmentResultResponse{
			Status:        false,
			StatusMessage: newError.Error(),
		}

		return resp, status.Errorf(codes.InvalidArgument, "%v", newError)
	}

	s.results[req.Result.Id] = req.Result

	go s.informHook(req.Result, nil)

	resp = &orchestrator.StoreAssessmentResultResponse{
		Status: true,
	}

	return resp, nil
}

func (s *Service) StoreAssessmentResults(stream orchestrator.Orchestrator_StoreAssessmentResultsServer) (err error) {
	var (
		result *orchestrator.StoreAssessmentResultRequest
		res    *orchestrator.StoreAssessmentResultResponse
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
		res, err = s.StoreAssessmentResult(context.Background(), storeAssessmentResultReq)
		if err != nil {
			log.Errorf("Error storing assessment result: %v", err)
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

// ListAssessmentResults is a method implementation of the orchestrator interface
func (svc *Service) ListAssessmentResults(_ context.Context, req *assessment.ListAssessmentResultsRequest) (res *assessment.ListAssessmentResultsResponse, err error) {
	res = new(assessment.ListAssessmentResultsResponse)

	// Paginate the results according to the request
	res.Results, res.NextPageToken, err = service.PaginateMapValues(req, svc.results, func(a *assessment.AssessmentResult, b *assessment.AssessmentResult) bool {
		return a.Id < b.Id
	}, service.DefaultPaginationOpts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
	}

	return
}

func (s *Service) RegisterAssessmentResultHook(hook func(result *assessment.AssessmentResult, err error)) {
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
	if req == nil {
		return nil,
			status.Errorf(codes.InvalidArgument, orchestrator.ErrRequestIsNil.Error())
	}
	if req.Certificate == nil {
		return nil,
			status.Errorf(codes.InvalidArgument, orchestrator.ErrCertificateIsNil.Error())
	}
	if req.Certificate.Id == "" {
		return nil,
			status.Errorf(codes.InvalidArgument, orchestrator.ErrCertIDIsMissing.Error())
	}

	// Persist the new certificate in our database
	err := svc.storage.Create(req.Certificate)
	if err != nil {
		return nil,
			status.Errorf(codes.Internal, "could not add certificate to the database: %v", err)
	}

	// Return certificate
	return req.Certificate, nil
}

// GetCertificate implements method for getting a certificate, e.g. to show its state in the UI
func (svc *Service) GetCertificate(_ context.Context, req *orchestrator.GetCertificateRequest) (response *orchestrator.Certificate, err error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, orchestrator.ErrRequestIsNil.Error())
	}
	if req.CertificateId == "" {
		return nil, status.Errorf(codes.NotFound, orchestrator.ErrCertIDIsMissing.Error())
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
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, orchestrator.ErrRequestIsNil.Error())
	}

	res = new(orchestrator.ListCertificatesResponse)

	res.Certificates, res.NextPageToken, err = service.PaginateStorage[*orchestrator.Certificate](req, svc.storage, service.DefaultPaginationOpts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
	}
	return
}

// UpdateCertificate implements method for updating an existing certificate
func (svc *Service) UpdateCertificate(_ context.Context, req *orchestrator.UpdateCertificateRequest) (response *orchestrator.Certificate, err error) {
	if req.CertificateId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "certificate id is empty")
	}

	if req.Certificate == nil {
		return nil, status.Errorf(codes.InvalidArgument, "certificate is empty")
	}

	count, err := svc.storage.Count(req.Certificate, "Certificate_id = ?", req.CertificateId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	if count == 0 {
		return nil, status.Error(codes.NotFound, "certificate not found")
	}

	response = req.Certificate
	response.Id = req.CertificateId

	err = svc.storage.Save(response, "Id = ?", response.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}
	return
}

// RemoveCertificate implements method for removing a certificate
func (svc *Service) RemoveCertificate(_ context.Context, req *orchestrator.RemoveCertificateRequest) (response *emptypb.Empty, err error) {
	if req.CertificateId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "certificate id is empty")
	}

	err = svc.storage.Delete(&orchestrator.Certificate{}, "Id = ?", req.CertificateId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "certificate not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	return &emptypb.Empty{}, nil
}
