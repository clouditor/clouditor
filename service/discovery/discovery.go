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
	"strings"
	"time"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/persistence/inmemory"
	"clouditor.io/clouditor/service"
	"clouditor.io/clouditor/service/discovery/aws"
	"clouditor.io/clouditor/service/discovery/azure"
	"clouditor.io/clouditor/service/discovery/k8s"
	"clouditor.io/clouditor/voc"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
)

const (
	ProviderAWS   = "aws"
	ProviderK8S   = "k8s"
	ProviderAzure = "azure"
)

var log *logrus.Entry

type grpcTarget struct {
	target string
	opts   []grpc.DialOption
}

// DiscoveryEventType defines the event types for [DiscoveryEvent].
type DiscoveryEventType int

const (
	// DiscovererStart is emmited at the start of a discovery run.
	DiscovererStart DiscoveryEventType = iota
	// DiscovererFinished is emmited at the end of a discovery run.
	DiscovererFinished
)

// DiscoveryEvent represents an event that is ommited if certain situations happen in the discoverer (defined by
// [DiscoveryEventType]). Examples would be the start or the end of the discovery. We will potentially expand this in
// the future.
type DiscoveryEvent struct {
	Type            DiscoveryEventType
	DiscovererName  string
	DiscoveredItems int
	Time            time.Time
}

// Service is an implementation of the Clouditor Discovery service.
// It should not be used directly, but rather the NewService constructor
// should be used.
type Service struct {
	discovery.UnimplementedDiscoveryServer

	configurations map[discovery.Discoverer]*Configuration

	assessmentStreams *api.StreamsOf[assessment.Assessment_AssessEvidencesClient, *assessment.AssessEvidenceRequest]
	assessmentAddress grpcTarget

	storage persistence.Storage

	scheduler *gocron.Scheduler

	authorizer api.Authorizer

	providers []string

	Events chan *DiscoveryEvent

	// csID is the cloud service ID for which we are gathering resources.
	csID string
}

type Configuration struct {
	Interval time.Duration
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
func WithAssessmentAddress(address string, opts ...grpc.DialOption) ServiceOption {
	return func(s *Service) {
		s.assessmentAddress = grpcTarget{
			target: address,
			opts:   opts,
		}
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
	return func(s *Service) {
		s.SetAuthorizer(api.NewOAuthAuthorizerFromClientCredentials(config))
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

func NewService(opts ...ServiceOption) *Service {
	var err error
	s := &Service{
		assessmentStreams: api.NewStreamsOf(api.WithLogger[assessment.Assessment_AssessEvidencesClient, *assessment.AssessEvidenceRequest](log)),
		assessmentAddress: grpcTarget{
			target: DefaultAssessmentAddress,
		},
		scheduler:      gocron.NewScheduler(time.UTC),
		configurations: make(map[discovery.Discoverer]*Configuration),
		Events:         make(chan *DiscoveryEvent),
		csID:           discovery.DefaultCloudServiceID,
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

// SetAuthorizer implements UsesAuthorizer.
func (svc *Service) SetAuthorizer(auth api.Authorizer) {
	svc.authorizer = auth
}

// Authorizer implements UsesAuthorizer.
func (svc *Service) Authorizer() api.Authorizer {
	return svc.authorizer
}

// initAssessmentStream initializes the stream that is used to send evidences to the assessment service.
// If configured, it uses the Authorizer of the discovery service to authenticate requests to the assessment.
func (svc *Service) initAssessmentStream(target string, additionalOpts ...grpc.DialOption) (stream assessment.Assessment_AssessEvidencesClient, err error) {
	log.Infof("Trying to establish a connection to assessment service @ %v", target)

	// Establish connection to assessment gRPC service
	conn, err := grpc.Dial(target,
		api.DefaultGrpcDialOptions(target, svc, additionalOpts...)...,
	)
	if err != nil {
		return nil, fmt.Errorf("could not connect to assessment service: %w", err)
	}

	client := assessment.NewAssessmentClient(conn)

	// Set up the stream and store it in our service struct, so we can access it later to actually
	// send the evidence data
	stream, err = client.AssessEvidences(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not set up stream for assessing evidences: %w", err)
	}

	log.Infof("Connected to Assessment")

	return
}

// Start starts discovery
func (svc *Service) Start(_ context.Context, req *discovery.StartDiscoveryRequest) (resp *discovery.StartDiscoveryResponse, err error) {
	var opts = []azure.DiscoveryOption{}
	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	resp = &discovery.StartDiscoveryResponse{Successful: true}

	log.Infof("Starting discovery...")
	svc.scheduler.TagsUnique()

	var discoverer []discovery.Discoverer

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

			discoverer = append(discoverer,
				azure.NewAzureComputeDiscovery(opts...),
				azure.NewAzureStorageDiscovery(opts...),
				azure.NewAzureNetworkDiscovery(opts...))
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
		svc.configurations[v] = &Configuration{
			Interval: 5 * time.Minute,
		}

		log.Infof("Scheduling {%s} to execute every 5 minutes...", v.Name())

		_, err = svc.scheduler.
			Every(5).
			Minute().
			Tag(v.Name()).
			Do(svc.StartDiscovery, v)
		if err != nil {
			log.Errorf("Could not schedule job for {%s}: %v", v.Name(), err)
		}
	}

	svc.scheduler.StartAsync()

	return resp, nil
}

func (svc *Service) Shutdown() {
	svc.scheduler.Stop()
}

func (svc *Service) StartDiscovery(discoverer discovery.Discoverer) {
	var (
		err  error
		list []voc.IsCloudResource
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
		// Set the cloud service ID to the one of the discoverer
		resource.SetServiceID(svc.csID)

		var (
			v *structpb.Value
		)

		v, err = voc.ToStruct(resource)
		if err != nil {
			log.Errorf("Could not convert resource to protobuf struct: %v", err)
		}

		// Build a resource struct. This will hold the latest sync state of the
		// resource for our storage layer.
		r := &discovery.Resource{
			Id:             string(resource.GetID()),
			ResourceType:   strings.Join(resource.GetType(), ","),
			CloudServiceId: resource.GetServiceID(),
			Properties:     v,
		}

		// Persist the latest state of the resource
		err = svc.storage.Save(&r, "id = ?", r.Id)
		if err != nil {
			log.Errorf("Could not save resource with ID '%s' to storage: %v", r.Id, err)
		}

		// TODO(all): What is the raw type in our case?
		e := &evidence.Evidence{
			Id:             uuid.New().String(),
			CloudServiceId: resource.GetServiceID(),
			Timestamp:      timestamppb.Now(),
			ToolId:         "Clouditor Evidences Collection",
			Raw:            nil,
			Resource:       v,
		}

		// Get Evidence Store stream
		channel, err := svc.assessmentStreams.GetStream(svc.assessmentAddress.target, "Assessment", svc.initAssessmentStream, svc.assessmentAddress.opts...)
		if err != nil {
			err = fmt.Errorf("could not get stream to assessment service (%s): %w", svc.assessmentAddress.target, err)
			log.Error(err)
			continue
		}

		channel.Send(&assessment.AssessEvidenceRequest{Evidence: e})
	}
}

func (svc *Service) Query(_ context.Context, req *discovery.QueryRequest) (res *discovery.QueryResponse, err error) {
	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	var query []string
	var args []any

	// Filtering the resources by
	// * cloud service ID
	// * resource type
	if req.Filter != nil {
		if req.Filter.CloudServiceId != nil {
			query = append(query, "cloud_service_id = ?")
			args = append(args, req.Filter.GetCloudServiceId())
		}
		if req.Filter.Type != nil {
			query = append(query, "(resource_type LIKE ? OR resource_type LIKE ? OR resource_type LIKE ?)")
			args = append(args, req.Filter.GetType()+",%", "%,"+req.Filter.GetType()+",%", "%,"+req.Filter.GetType())
		}
	}

	res = new(discovery.QueryResponse)

	// Join query with AND and prepend the query
	args = append([]any{strings.Join(query, " AND ")}, args...)

	res.Results, res.NextPageToken, err = service.PaginateStorage[*discovery.Resource](req, svc.storage, service.DefaultPaginationOpts, args...)

	return
}
