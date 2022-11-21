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

	"github.com/go-co-op/gocron"
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

// Evaluate is a method implementation of the evaluation interface: It starts the evaluation for a cloud service
func (s *Service) Evaluate(_ context.Context, req *evaluation.EvaluateRequest) (resp *evaluation.EvaluateResponse, err error) {
	resp = &evaluation.EvaluateResponse{}

	err = req.Validate()
	if err != nil {
		resp = &evaluation.EvaluateResponse{
			Status:        false,
			StatusMessage: err.Error(),
		}
		return resp, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	log.Info("Starting evaluation ...")
	s.scheduler.TagsUnique()

	log.Infof("Evaluate Cloud Service {%s} every 5 minutes...", req.Toe.CloudServiceId)
	_, err = s.scheduler.
		Every(5).
		Minute().
		Tag(req.Toe.CloudServiceId).
		Do(s.StartEvaluation, req)

	s.scheduler.StartAsync()

	return
}

func (s *Service) StartEvaluation(req *evaluation.EvaluateRequest) {
	log.Infof("Started evaluation for Cloud Service '%s',  Catalog ID '%s' and Control '%s'", req.Toe.CloudServiceId, req.Toe.CatalogId, req.ControlId)

	// Verify that monitoring of this service and controls hasn't started already
	// TODO(anatheka)

	// Establish connection to orchestrator gRPC service
	conn, err := grpc.Dial(s.orchestratorAddress.target,
		api.DefaultGrpcDialOptions(s.orchestratorAddress.target, s, s.orchestratorAddress.opts...)...)
	if err != nil {
		log.Errorf("could not connect to orchestrator service: %v", err)
	}

	// Find associated metrics to control
	// For each metric, find all assessment results for the target cloud service
	// If at least one assessment result is non-compliant --> evaluation = fail; otherwise --> evaluation = ok
	// Optional: Add failing assessment results to list of failing assessment results

	// Create orchestrator client
	s.orchestratorClient = orchestrator.NewOrchestratorClient(conn)

	// Get control and the associated metrics
	control, err := s.orchestratorClient.GetControl(context.Background(), &orchestrator.GetControlRequest{
		CatalogId:    req.Toe.CatalogId,
		CategoryName: req.CategoryName,
		ControlId:    req.ControlId,
	})
	if err != nil {
		log.Errorf("Could not get control '%s' from Orchestrator: %v", req.ControlId, err)
		// TODO(anatheka): Do we need that?
		s.Shutdown()
		return
	}
	metrics := control.Metrics

	// Get assessment results for the target cloud service
	// TODO(anatheka): The filtered_cloud_service_id option does not work: access denied
	// TODO(anatheka): Get all results, not only the first page.
	results, err := s.orchestratorClient.ListAssessmentResults(context.Background(), &assessment.ListAssessmentResultsRequest{
		// FilteredCloudServiceId: req.Toe.CloudServiceId,
	})
	if err != nil {
		log.Errorf("Could not get assessment results for Cloud Serivce '%s' from Orchestrator", req.Toe.CloudServiceId)
		// TODO(anatheka): Do we need that?
		s.Shutdown()
		return
	}

	// For each metric, find all assessment results for the target cloud service
	var mappingList []mappingResultMetric
	for _, metric := range metrics {
		var mapping mappingResultMetric
		var resultList []*assessment.AssessmentResult
		for _, result := range results.Results {
			if result.MetricId == metric.Id {
				resultList = append(resultList, result)
			}
		}
		mapping.metricName = metric.Name
		mapping.results = resultList
		mappingList = append(mappingList, mapping)
	}

	log.Debugf("Found %d metrics.", len(mappingList))
	// Do evaluation and find all non-comüpliant assessment results

}

func (s *Service) Shutdown() {
	s.scheduler.Stop()
}
