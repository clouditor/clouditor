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
	"fmt"
	"sync"

	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/api/runtime"
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/launcher"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/persistence/inmemory"
	"clouditor.io/clouditor/v2/server"
	"clouditor.io/clouditor/v2/service"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//go:embed *.json
var f embed.FS

var DefaultMetricsFile = "metrics.json"
var DefaultCatalogsFolder = "catalogs"

var (
	defaultMetricConfigurations map[string]*assessment.MetricConfiguration
	log                         *logrus.Entry
)

func DefaultServiceSpec() launcher.ServiceSpec {
	return launcher.NewServiceSpec(
		NewService,
		WithStorage,
		func(svc *Service) ([]server.StartGRPCServerOption, error) {
			// It is possible to register hook functions for the orchestrator.
			//  * The hook functions in orchestrator are implemented in StoreAssessmentResult(s)

			// svc.RegisterAssessmentResultHook(func(result *assessment.AssessmentResult, err error) {})

			// Create default target Cloud Service
			if viper.GetBool(config.CreateDefaultTargetFlag) {
				_, err := svc.CreateDefaultTargetCloudService()
				if err != nil {
					return nil, fmt.Errorf("could not register default target cloud service: %v", err)
				}
			}

			return []server.StartGRPCServerOption{
				server.WithOrchestrator(svc),
			}, nil
		},
	)
}

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

	// Hook
	AssessmentResultHooks []assessment.ResultHookFunc
	// mu is used for (un)locking result hook calls
	mu sync.Mutex

	storage persistence.Storage

	metricsFile string

	// loadMetricsFunc is a function that is used to initially load metrics at the start of the orchestrator
	loadMetricsFunc func() ([]*assessment.Metric, error)

	catalogsFolder string

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

// WithMetricsFile can be used to load a different metrics file
func WithMetricsFile(file string) service.Option[*Service] {
	return func(s *Service) {
		s.metricsFile = file
	}
}

// WithExternalMetrics can be used to load metric definitions from an external source
func WithExternalMetrics(f func() ([]*assessment.Metric, error)) service.Option[*Service] {
	return func(s *Service) {
		s.loadMetricsFunc = f
	}
}

// WithCatalogsFolder can be used to load catalog files from a different catalogs folder
func WithCatalogsFolder(folder string) service.Option[*Service] {
	return func(s *Service) {
		s.catalogsFolder = folder
	}
}

// WithExternalCatalogs can be used to load catalog definitions from an external source
func WithExternalCatalogs(f func() ([]*orchestrator.Catalog, error)) service.Option[*Service] {
	return func(s *Service) {
		s.loadCatalogsFunc = f
	}
}

// WithStorage is an option to set the storage. If not set, NewService will use inmemory storage.
func WithStorage(storage persistence.Storage) service.Option[*Service] {
	return func(s *Service) {
		s.storage = storage
	}
}

// WithAuthorizationStrategyJWT is an option that configures an JWT-based authorization strategy using a specific claim key.
func WithAuthorizationStrategyJWT(key string, allowAllKey string) service.Option[*Service] {
	return func(s *Service) {
		s.authz = &service.AuthorizationStrategyJWT{CloudServicesKey: key, AllowAllKey: allowAllKey}
	}
}

func WithAuthorizationStrategy(authz service.AuthorizationStrategy) service.Option[*Service] {
	return func(s *Service) {
		s.authz = authz
	}
}

// NewService creates a new Orchestrator service
func NewService(opts ...service.Option[*Service]) *Service {
	var err error
	s := Service{
		metricsFile:    DefaultMetricsFile,
		catalogsFolder: DefaultCatalogsFolder,
		events:         make(chan *orchestrator.MetricChangeEvent, 1000),
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

func (svc *Service) Init() {}

func (svc *Service) Shutdown() {}

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

// GetRuntimeInfo implements a method to retrieve runtime information
func (*Service) GetRuntimeInfo(_ context.Context, _ *runtime.GetRuntimeInfoRequest) (res *runtime.Runtime, err error) {
	return service.GetRuntimeInfo()
}
