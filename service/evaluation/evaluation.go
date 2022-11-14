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

// Service is an implementation of the Clouditor Evaluation service
type Service struct {
	evaluation.UnimplementedEvaluationServer
	scheduler           *gocron.Scheduler
	orchestratorClient  orchestrator.OrchestratorClient
	orchestratorAddress grpcTarget
	authorizer          api.Authorizer
	// csID is the cloud service ID for which we start the evaluation
	csID string
}

func init() {
	log = logrus.WithField("component", "orchestrator")
}

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

// NewService creates a new Evaluation service
func NewService(opts ...ServiceOption) *Service {
	s := Service{
		scheduler: gocron.NewScheduler(time.UTC),
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

	return
}

func (s *Service) StartEvaluation(req *evaluation.EvaluateRequest) {
	log.Infof("Started evaluation for Cloud Service '%s' and Catalog ID '%s'", req.Toe.CloudServiceId, req.Toe.CatalogId)

	// Find associated metrics for control
	// ControlId
	// ToE
	// * catalog_id
	// * cloud_service_id
	// * assurance_level

	// Establish connection to orchestrator gRPC service
	conn, err := grpc.Dial(s.orchestratorAddress.target,
		api.DefaultGrpcDialOptions(s.orchestratorAddress.target, s, s.orchestratorAddress.opts...)...)
	if err != nil {
		log.Errorf("could not connect to orchestrator service: %v", err)
	}

	s.orchestratorClient = orchestrator.NewOrchestratorClient(conn)
	catalog, err := s.orchestratorClient.GetCatalog(context.Background(), &orchestrator.GetCatalogRequest{CatalogId: req.Toe.CatalogId})
	if err != nil {
		log.Errorf("Could not get catalog '%s' from Orchestrator", req.Toe.CatalogId)
	}

	log.Infof("Catalog: %v", catalog)

}

func (s *Service) Shutdown() {
	log.Debug(s.)
	s.scheduler.Stop()
}
