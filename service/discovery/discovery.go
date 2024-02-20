// Copyright 2016-2020 Fraunhofer AISEC
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

package discovery

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/assessment/assessmentconnect"
	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/util"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/persistence/inmemory"
	"clouditor.io/clouditor/v2/service"
	"clouditor.io/clouditor/v2/service/discovery/aws"
	"clouditor.io/clouditor/v2/service/discovery/azure"
	"clouditor.io/clouditor/v2/service/discovery/k8s"
	"connectrpc.com/connect"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	ProviderAWS   = "aws"
	ProviderK8S   = "k8s"
	ProviderAzure = "azure"
)

var log *logrus.Entry

// DiscoveryEventType defines the event types for [DiscoveryEvent].
type DiscoveryEventType int

const (
	// DiscovererStart is emmited at the start of a discovery run.
	DiscovererStart DiscoveryEventType = iota
	// DiscovererFinished is emmited at the end of a discovery run.
	DiscovererFinished
)

// DiscoveryEvent represents an event that is emitted if certain situations happen in the discoverer (defined by
// [DiscoveryEventType]). Examples would be the start or the end of the discovery. We will potentially expand this in
// the future.
type DiscoveryEvent struct {
	Type            DiscoveryEventType
	DiscovererName  string
	DiscoveredItems int
	Time            time.Time
}

type rpcOpts struct {
	URL    string
	Client connect.HTTPClient
	Opts   []connect.ClientOption
}

// Service is an implementation of the Clouditor Discovery service (plus its experimental extensions). It should not be
// used directly, but rather the NewService constructor should be used.
type Service struct {
	authorizer api.Authorizer

	assessmentStreams *api.ConnectStreamsOf[
		assessmentconnect.AssessmentClient,
		assessment.AssessEvidenceRequest,
		assessment.AssessEvidencesResponse,
	]
	assessment     assessmentconnect.AssessmentClient
	assessmentOpts rpcOpts

	storage persistence.Storage

	scheduler *gocron.Scheduler

	authz service.AuthorizationStrategy

	providers []string

	discoveryInterval time.Duration

	Events chan *DiscoveryEvent

	// csID is the cloud service ID for which we are gathering resources.
	csID string

	target string
}

func init() {
	log = logrus.WithField("component", "discovery")
}

const (
	// DefaultAssessmentAddress specifies the default gRPC address of the assessment service.
	DefaultAssessmentAddress = "localhost:9090"
)

// ServiceOption is a functional option type to configure the discovery service.
type ServiceOption func(*Service)

// WithAssessmentAddress is an option to configure the assessment service gRPC address.
func WithAssessmentAddress(target string, opts ...connect.ClientOption) ServiceOption {
	return func(s *Service) {
		s.assessmentOpts.URL = target
		s.assessmentOpts.Opts = opts
	}
}

// WithCloudServiceID is an option to configure the cloud service ID for which resources will be discovered.
func WithCloudServiceID(ID string) ServiceOption {
	return func(svc *Service) {
		svc.csID = ID
	}
}

// WithOAuth2Authorizer is an option to use an OAuth 2.0 authorizer
func WithOAuth2Authorizer(config *clientcredentials.Config) ServiceOption {
	return func(svc *Service) {
		svc.authorizer = api.NewOAuthAuthorizerFromClientCredentials(config)
		//svc.assessment.SetAuthorizer()
	}
}

// WithProviders is an option to set providers for discovering
func WithProviders(providersList []string) ServiceOption {
	if len(providersList) == 0 {
		newError := errors.New("no providers given")
		log.Error(newError)
	}

	return func(s *Service) {
		s.providers = providersList
	}
}

// WithStorage is an option to set the storage. If not set, NewService will use inmemory storage.
func WithStorage(storage persistence.Storage) ServiceOption {
	return func(s *Service) {
		s.storage = storage
	}
}

// WithDiscoveryInterval is an option to set the discovery interval. If not set, the discovery is set to 5 minutes.
func WithDiscoveryInterval(interval time.Duration) ServiceOption {
	return func(s *Service) {
		s.discoveryInterval = interval
	}
}

// WithAuthorizationStrategy is an option that configures an authorization strategy to be used with this service.
func WithAuthorizationStrategy(authz service.AuthorizationStrategy) ServiceOption {
	return func(s *Service) {
		s.authz = authz
	}
}

func NewService(opts ...ServiceOption) *Service {
	var err error

	s := &Service{
		assessmentStreams: api.NewConnectStreamsOf[
			assessmentconnect.AssessmentClient,
			assessment.AssessEvidenceRequest,
			assessment.AssessEvidencesResponse,
		](assessmentconnect.AssessmentClient.AssessEvidences),
		assessmentOpts: rpcOpts{
			URL:    DefaultAssessmentAddress,
			Client: http.DefaultClient,
		},
		scheduler:         gocron.NewScheduler(time.UTC),
		Events:            make(chan *DiscoveryEvent),
		csID:              discovery.DefaultCloudServiceID,
		authz:             &service.AuthorizationStrategyAllowAll{},
		discoveryInterval: 5 * time.Minute, // Default discovery interval is 5 minutes
	}

	// Apply any options
	for _, o := range opts {
		o(s)
	}

	// Create RPC client(s)
	// TODO: interceptors
	s.assessment = assessmentconnect.NewAssessmentClient(
		http.DefaultClient,
		s.assessmentOpts.URL,
		s.assessmentOpts.Opts...,
	)

	// Default to an in-memory storage, if nothing was explicitly set
	if s.storage == nil {
		s.storage, err = inmemory.NewStorage()
		if err != nil {
			log.Errorf("Could not initialize the storage: %v", err)
		}
	}

	return s
}

// Start starts discovery
func (svc *Service) Start(ctx context.Context, req *connect.Request[discovery.StartDiscoveryRequest]) (res *connect.Response[discovery.StartDiscoveryResponse], err error) {
	var (
		opts       = []azure.DiscoveryOption{}
		discoverer []discovery.Discoverer
	)

	// Validate request
	err = api.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	// Check if cloud_service_id in the service is within allowed or one can access *all* the cloud services
	if !svc.authz.CheckAccess(ctx, service.AccessUpdate, svc) {
		return nil, service.ErrPermissionDenied
	}

	res = connect.NewResponse(&discovery.StartDiscoveryResponse{Successful: true})

	log.Infof("Starting discovery...")
	svc.scheduler.TagsUnique()

	// Configure discoverers for given providers
	for _, provider := range svc.providers {
		switch {
		case provider == ProviderAzure:
			authorizer, err := azure.NewAuthorizer()
			if err != nil {
				log.Errorf("Could not authenticate to Azure: %v", err)
				return nil, status.Errorf(codes.FailedPrecondition, "could not authenticate to Azure: %v", err)
			}
			// Add authorizer and cloudServiceID
			opts = append(opts, azure.WithAuthorizer(authorizer), azure.WithCloudServiceID(svc.csID))
			// Check if resource group is given and append to discoverer
			if req.Msg.GetResourceGroup() != "" {
				opts = append(opts, azure.WithResourceGroup(req.Msg.GetResourceGroup()))
			}
			discoverer = append(discoverer, azure.NewAzureDiscovery(opts...))
		case provider == ProviderK8S:
			k8sClient, err := k8s.AuthFromKubeConfig()
			if err != nil {
				log.Errorf("Could not authenticate to Kubernetes: %v", err)
				return nil, status.Errorf(codes.FailedPrecondition, "could not authenticate to Kubernetes: %v", err)
			}
			discoverer = append(discoverer,
				k8s.NewKubernetesComputeDiscovery(k8sClient, svc.csID),
				k8s.NewKubernetesNetworkDiscovery(k8sClient, svc.csID),
				k8s.NewKubernetesStorageDiscovery(k8sClient, svc.csID))
		case provider == ProviderAWS:
			awsClient, err := aws.NewClient()
			if err != nil {
				log.Errorf("Could not authenticate to AWS: %v", err)
				return nil, status.Errorf(codes.FailedPrecondition, "could not authenticate to AWS: %v", err)
			}
			discoverer = append(discoverer,
				aws.NewAwsStorageDiscovery(awsClient, svc.csID),
				aws.NewAwsComputeDiscovery(awsClient, svc.csID))
		default:
			newError := fmt.Errorf("provider %s not known", provider)
			log.Error(newError)
			return nil, status.Errorf(codes.InvalidArgument, "%s", newError)
		}
	}

	for _, v := range discoverer {
		log.Infof("Scheduling {%s} to execute every {%v} minutes...", v.Name(), svc.discoveryInterval.Minutes())

		_, err = svc.scheduler.
			Every(svc.discoveryInterval).
			Tag(v.Name()).
			Do(svc.StartDiscovery, v)
		if err != nil {
			newError := fmt.Errorf("could not schedule job for {%s}: %v", v.Name(), err)
			log.Error(newError)
			return nil, status.Errorf(codes.Aborted, "%s", newError)
		}
	}

	svc.scheduler.StartAsync()

	return res, nil
}

func (svc *Service) Shutdown() {
	log.Info("Shutting down discovery service")

	svc.assessmentStreams.CloseAll()
	svc.scheduler.Stop()
}

func (svc *Service) StartDiscovery(discoverer discovery.Discoverer) {
	var (
		err  error
		list []ontology.IsResource
	)

	go func() {
		svc.Events <- &DiscoveryEvent{
			Type:           DiscovererStart,
			DiscovererName: discoverer.Name(),
			Time:           time.Now(),
		}
	}()

	list, err = discoverer.List()

	if err != nil {
		log.Errorf("Could not retrieve resources from discoverer '%s': %v", discoverer.Name(), err)
		return
	}

	// Notify event listeners that the discoverer is finished
	go func() {
		svc.Events <- &DiscoveryEvent{
			Type:            DiscovererFinished,
			DiscovererName:  discoverer.Name(),
			DiscoveredItems: len(list),
			Time:            time.Now(),
		}
	}()

	for _, resource := range list {
		// Build a resource struct. This will hold the latest sync state of the
		// resource for our storage layer.
		r, err := discovery.ToDiscoveryResource(resource, svc.GetCloudServiceId())
		if err != nil {
			log.Errorf("Could not convert resource: %v", err)
			continue
		}

		// Persist the latest state of the resource
		err = svc.storage.Save(&r, "id = ?", r.Id)
		if err != nil {
			log.Errorf("Could not save resource with ID '%s' to storage: %v", r.Id, err)
		}

		a, err := anypb.New(resource)
		if err != nil {
			log.Errorf("Could not wrap resource message into Any protobuf object: %v", err)
			continue
		}

		e := &evidence.Evidence{
			Id:             uuid.New().String(),
			CloudServiceId: svc.GetCloudServiceId(),
			Timestamp:      timestamppb.Now(),
			Raw:            util.Ref(resource.GetRaw()),
			ToolId:         discovery.EvidenceCollectorToolId,
			Resource:       a,
		}

		// Get Evidence Store stream
		channel, err := svc.assessmentStreams.GetStream(svc.assessment, "Assessment")
		if err != nil {
			err = fmt.Errorf("could not get stream to assessment service (%s): %w", svc.target, err)
			log.Error(err)
			continue
		}

		channel.Send(&assessment.AssessEvidenceRequest{Evidence: e})
	}
}

func (svc *Service) ListResources(ctx context.Context, req *connect.Request[discovery.ListResourcesRequest]) (res *connect.Response[discovery.ListResourcesResponse], err error) {
	var (
		query   []string
		args    []any
		all     bool
		allowed []string
	)

	// Validate request
	err = api.Validate(req.Msg)
	if err != nil {
		return nil, err
	}

	// Filtering the resources by
	// * cloud service ID
	// * resource type
	if req.Msg.Filter != nil {
		// Check if cloud_service_id in filter is within allowed or one can access *all* the cloud services
		if !svc.authz.CheckAccess(ctx, service.AccessRead, req.Msg.Filter) {
			return nil, service.ErrPermissionDenied
		}

		if req.Msg.Filter.CloudServiceId != nil {
			query = append(query, "cloud_service_id = ?")
			args = append(args, req.Msg.Filter.GetCloudServiceId())
		}
		if req.Msg.Filter.Type != nil {
			query = append(query, "(resource_type LIKE ? OR resource_type LIKE ? OR resource_type LIKE ?)")
			args = append(args, req.Msg.Filter.GetType()+",%", "%,"+req.Msg.Filter.GetType()+",%", "%,"+req.Msg.Filter.GetType())
		}
	}

	// We need to further restrict our query according to the cloud service we are allowed to "see".
	//
	// TODO(oxisto): This is suboptimal, since we are now calling AllowedCloudServices twice. Once here
	//  and once above in CheckAccess.
	all, allowed = svc.authz.AllowedCloudServices(ctx)
	if !all {
		query = append(query, "cloud_service_id IN ?")
		args = append(args, allowed)
	}

	res = connect.NewResponse(&discovery.ListResourcesResponse{})

	// Join query with AND and prepend the query
	args = append([]any{strings.Join(query, " AND ")}, args...)

	res.Msg.Results, res.Msg.NextPageToken, err = service.PaginateStorage[*discovery.Resource](req.Msg, svc.storage, service.DefaultPaginationOpts, args...)

	return
}

// GetCloudServiceId implements CloudServiceRequest for this service. This is a little trick, so that we can call
// CheckAccess directly on the service. This is necessary because the discovery service itself is tied to a specific
// cloud service ID, instead of the individual requests that are made against the service.
func (svc *Service) GetCloudServiceId() string {
	return svc.csID
}
