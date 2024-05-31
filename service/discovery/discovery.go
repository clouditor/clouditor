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
	"slices"
	"strings"
	"time"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/internal/util"
	"clouditor.io/clouditor/v2/launcher"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/persistence/inmemory"
	"clouditor.io/clouditor/v2/server/rest"
	"clouditor.io/clouditor/v2/service"
	"clouditor.io/clouditor/v2/service/discovery/aws"
	"clouditor.io/clouditor/v2/service/discovery/azure"
	"clouditor.io/clouditor/v2/service/discovery/extra/csaf"
	"clouditor.io/clouditor/v2/service/discovery/k8s"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	ProviderAWS   = "aws"
	ProviderK8S   = "k8s"
	ProviderAzure = "azure"
	ProviderCSAF  = "csaf"
)

var log *logrus.Entry

// DefaultServiceSpec returns a [launcher.ServiceSpec] for this [Service] with all necessary options retrieved from the
// config system.
func DefaultServiceSpec() launcher.ServiceSpec {
	var providers []string

	// If no CSPs for discovering are given, take all implemented discoverers
	if len(viper.GetStringSlice(config.DiscoveryProviderFlag)) == 0 {
		providers = []string{ProviderAWS, ProviderAzure, ProviderK8S}
	} else {
		providers = viper.GetStringSlice(config.DiscoveryProviderFlag)
	}

	return launcher.NewServiceSpec(
		NewService,
		WithStorage,
		nil,
		WithOAuth2Authorizer(config.ClientCredentials()),
		WithCloudServiceID(viper.GetString(config.CloudServiceIDFlag)),
		WithProviders(providers),
		WithAssessmentAddress(viper.GetString(config.AssessmentURLFlag)),
	)
}

// DiscoveryEventType defines the event types for [DiscoveryEvent].
type DiscoveryEventType int

const (
	// DiscovererStart is emitted at the start of a discovery run.
	DiscovererStart DiscoveryEventType = iota
	// DiscovererFinished is emitted at the end of a discovery run.
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

// Service is an implementation of the Clouditor Discovery service (plus its experimental extensions). It should not be
// used directly, but rather the NewService constructor should be used.
type Service struct {
	discovery.UnimplementedDiscoveryServer
	discovery.UnimplementedExperimentalDiscoveryServer

	assessmentStreams *api.StreamsOf[assessment.Assessment_AssessEvidencesClient, *assessment.AssessEvidenceRequest]
	assessment        *api.RPCConnection[assessment.AssessmentClient]

	storage persistence.Storage

	scheduler *gocron.Scheduler

	authz service.AuthorizationStrategy

	providers   []string
	discoverers []discovery.Discoverer

	discoveryInterval time.Duration

	Events chan *DiscoveryEvent

	// csID is the cloud service ID for which we are gathering resources.
	csID string

	// collectorID is the evidence collector tool ID which is gathering the resources.
	collectorID string
}

func init() {
	log = logrus.WithField("component", "discovery")
}

const (
	// DefaultAssessmentAddress specifies the default gRPC address of the assessment service.
	DefaultAssessmentAddress = "localhost:9090"
)

// WithAssessmentAddress is an option to configure the assessment service gRPC address.
func WithAssessmentAddress(target string, opts ...grpc.DialOption) service.Option[*Service] {

	return func(s *Service) {
		log.Infof("Assessment URL is set to %s", target)

		s.assessment.Target = target
		s.assessment.Opts = opts
	}
}

// WithCloudServiceID is an option to configure the cloud service ID for which resources will be discovered.
func WithCloudServiceID(ID string) service.Option[*Service] {
	return func(svc *Service) {
		log.Infof("Cloud Service ID is set to %s", ID)

		svc.csID = ID
	}
}

// WithEvidenceCollectorToolID is an option to configure the collector tool ID that is used to discover resources.
func WithEvidenceCollectorToolID(ID string) service.Option[*Service] {
	return func(svc *Service) {
		log.Infof("Evidence Collector Tool ID is set to %s", ID)

		svc.collectorID = ID
	}
}

// WithOAuth2Authorizer is an option to use an OAuth 2.0 authorizer
func WithOAuth2Authorizer(config *clientcredentials.Config) service.Option[*Service] {
	return func(svc *Service) {
		svc.assessment.SetAuthorizer(api.NewOAuthAuthorizerFromClientCredentials(config))
	}
}

// WithProviders is an option to set providers for discovering
func WithProviders(providersList []string) service.Option[*Service] {
	if len(providersList) == 0 {
		newError := errors.New("no providers given")
		log.Error(newError)
	}

	return func(s *Service) {
		s.providers = providersList
	}
}

// WithAdditionalDiscoverers is an option to add additional discoverers for discovering. Note: These are added in
// addition to the ones created by [WithProviders].
func WithAdditionalDiscoverers(discoverers []discovery.Discoverer) service.Option[*Service] {
	return func(s *Service) {
		s.discoverers = append(s.discoverers, discoverers...)
	}
}

// WithStorage is an option to set the storage. If not set, NewService will use inmemory storage.
func WithStorage(storage persistence.Storage) service.Option[*Service] {
	return func(s *Service) {
		s.storage = storage
	}
}

// WithDiscoveryInterval is an option to set the discovery interval. If not set, the discovery is set to 5 minutes.
func WithDiscoveryInterval(interval time.Duration) service.Option[*Service] {
	return func(s *Service) {
		s.discoveryInterval = interval
	}
}

// WithAuthorizationStrategy is an option that configures an authorization strategy to be used with this service.
func WithAuthorizationStrategy(authz service.AuthorizationStrategy) service.Option[*Service] {
	return func(s *Service) {
		s.authz = authz
	}
}

func NewService(opts ...service.Option[*Service]) *Service {
	var err error
	s := &Service{
		assessmentStreams: api.NewStreamsOf(api.WithLogger[assessment.Assessment_AssessEvidencesClient, *assessment.AssessEvidenceRequest](log)),
		assessment:        api.NewRPCConnection(DefaultAssessmentAddress, assessment.NewAssessmentClient),
		scheduler:         gocron.NewScheduler(time.UTC),
		Events:            make(chan *DiscoveryEvent),
		csID:              config.DefaultCloudServiceID,
		collectorID:       config.DefaultEvidenceCollectorToolID,
		authz:             &service.AuthorizationStrategyAllowAll{},
		discoveryInterval: 5 * time.Minute, // Default discovery interval is 5 minutes
	}

	// Apply any options
	for _, o := range opts {
		o(s)
	}

	// Default to an in-memory storage, if nothing was explicitly set
	if s.storage == nil {
		s.storage, err = inmemory.NewStorage()
		if err != nil {
			log.Errorf("Could not initialize the storage: %v", err)
		}
	}

	return s
}

func (svc *Service) Init() {
	var err error

	// Automatically start the discovery, if we have this flag enabled
	if viper.GetBool(config.DiscoveryAutoStartFlag) {
		go func() {
			<-rest.GetReadyChannel()
			_, err = svc.Start(context.Background(), &discovery.StartDiscoveryRequest{
				ResourceGroup: util.Ref(viper.GetString(config.DiscoveryResourceGroupFlag)),
				CsafDomain:    util.Ref(viper.GetString(config.DiscoveryCSAFDomainFlag)),
			})
			if err != nil {
				log.Errorf("Could not automatically start discovery: %v", err)
			}
		}()
	}
}

func (svc *Service) Shutdown() {
	svc.assessmentStreams.CloseAll()
	svc.scheduler.Stop()
}

// initAssessmentStream initializes the stream that is used to send evidences to the assessment service.
// If configured, it uses the Authorizer of the discovery service to authenticate requests to the assessment.
func (svc *Service) initAssessmentStream(target string, _ ...grpc.DialOption) (stream assessment.Assessment_AssessEvidencesClient, err error) {
	log.Infof("Trying to establish a connection to assessment service @ %v", target)

	// Make sure, that we re-connect
	svc.assessment.ForceReconnect()

	// Set up the stream and store it in our service struct, so we can access it later to actually
	// send the evidence data
	stream, err = svc.assessment.Client.AssessEvidences(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not set up stream for assessing evidences: %w", err)
	}

	log.Infof("Connected to Assessment")

	return
}

// Start starts discovery
func (svc *Service) Start(ctx context.Context, req *discovery.StartDiscoveryRequest) (resp *discovery.StartDiscoveryResponse, err error) {
	var (
		opts = []azure.DiscoveryOption{}
	)

	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check if cloud_service_id in the service is within allowed or one can access *all* the cloud services
	if !svc.authz.CheckAccess(ctx, service.AccessUpdate, svc) {
		return nil, service.ErrPermissionDenied
	}

	resp = &discovery.StartDiscoveryResponse{Successful: true}

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
			if req.GetResourceGroup() != "" {
				opts = append(opts, azure.WithResourceGroup(req.GetResourceGroup()))
			}
			svc.discoverers = append(svc.discoverers, azure.NewAzureDiscovery(opts...))
		case provider == ProviderK8S:
			k8sClient, err := k8s.AuthFromKubeConfig()
			if err != nil {
				log.Errorf("Could not authenticate to Kubernetes: %v", err)
				return nil, status.Errorf(codes.FailedPrecondition, "could not authenticate to Kubernetes: %v", err)
			}
			svc.discoverers = append(svc.discoverers,
				k8s.NewKubernetesComputeDiscovery(k8sClient, svc.csID),
				k8s.NewKubernetesNetworkDiscovery(k8sClient, svc.csID),
				k8s.NewKubernetesStorageDiscovery(k8sClient, svc.csID))
		case provider == ProviderAWS:
			awsClient, err := aws.NewClient()
			if err != nil {
				log.Errorf("Could not authenticate to AWS: %v", err)
				return nil, status.Errorf(codes.FailedPrecondition, "could not authenticate to AWS: %v", err)
			}
			svc.discoverers = append(svc.discoverers,
				aws.NewAwsStorageDiscovery(awsClient, svc.csID),
				aws.NewAwsComputeDiscovery(awsClient, svc.csID))
		case provider == ProviderCSAF:
			var (
				domain string
				opts   []csaf.DiscoveryOption
			)
			domain = util.Deref(req.CsafDomain)
			if domain != "" {
				opts = append(opts, csaf.WithProviderDomain(domain))
			}
			svc.discoverers = append(svc.discoverers, csaf.NewTrustedProviderDiscovery(opts...))
		default:
			newError := fmt.Errorf("provider %s not known", provider)
			log.Error(newError)
			return nil, status.Errorf(codes.InvalidArgument, "%s", newError)
		}
	}

	for _, v := range svc.discoverers {
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

	return resp, nil
}

func (svc *Service) StartDiscovery(discoverer discovery.Discoverer) {
	var (
		err  error
		list []ontology.IsResource
	)

	fmt.Println("Collector id: ", svc.collectorID)

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
		r, err := discovery.ToDiscoveryResource(resource, svc.GetCloudServiceId(), svc.collectorID)
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
			ToolId:         svc.collectorID,
			Resource:       a,
		}

		// Only enabled related evidences for some specific resources for now
		if slices.Contains(ontology.ResourceTypes(resource), "SecurityAdvisoryService") {
			edges := ontology.Related(resource)
			for _, edge := range edges {
				e.ExperimentalRelatedResourceIds = append(e.ExperimentalRelatedResourceIds, edge.Value)
			}
		}

		// Get Evidence Store stream
		channel, err := svc.assessmentStreams.GetStream(svc.assessment.Target, "Assessment", svc.initAssessmentStream, svc.assessment.Opts...)
		if err != nil {
			err = fmt.Errorf("could not get stream to assessment service (%s): %w", svc.assessment.Target, err)
			log.Error(err)
			continue
		}

		channel.Send(&assessment.AssessEvidenceRequest{Evidence: e})
	}
}

func (svc *Service) ListResources(ctx context.Context, req *discovery.ListResourcesRequest) (res *discovery.ListResourcesResponse, err error) {
	var (
		query   []string
		args    []any
		all     bool
		allowed []string
	)

	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Filtering the resources by
	// * cloud service ID
	// * resource type
	if req.Filter != nil {
		// Check if cloud_service_id in filter is within allowed or one can access *all* the cloud services
		if !svc.authz.CheckAccess(ctx, service.AccessRead, req.Filter) {
			return nil, service.ErrPermissionDenied
		}

		if req.Filter.CloudServiceId != nil {
			query = append(query, "cloud_service_id = ?")
			args = append(args, req.Filter.GetCloudServiceId())
		}
		if req.Filter.Type != nil {
			query = append(query, "(resource_type LIKE ? OR resource_type LIKE ? OR resource_type LIKE ?)")
			args = append(args, req.Filter.GetType()+",%", "%,"+req.Filter.GetType()+",%", "%,"+req.Filter.GetType())
		}
		if req.Filter.ToolId != nil {
			query = append(query, "(tool_id = ?)")
			args = append(args, req.Filter.GetToolId())
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

	res = new(discovery.ListResourcesResponse)

	// Join query with AND and prepend the query
	args = append([]any{strings.Join(query, " AND ")}, args...)

	res.Results, res.NextPageToken, err = service.PaginateStorage[*discovery.Resource](req, svc.storage, service.DefaultPaginationOpts, args...)

	return
}

// GetCloudServiceId implements CloudServiceRequest for this service. This is a little trick, so that we can call
// CheckAccess directly on the service. This is necessary because the discovery service itself is tied to a specific
// cloud service ID, instead of the individual requests that are made against the service.
func (svc *Service) GetCloudServiceId() string {
	return svc.csID
}
