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
	"time"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/internal/util"
	"clouditor.io/clouditor/v2/launcher"
	"clouditor.io/clouditor/v2/server/rest"
	"clouditor.io/clouditor/v2/service"
	"clouditor.io/clouditor/v2/service/discovery/aws"
	"clouditor.io/clouditor/v2/service/discovery/azure"
	"clouditor.io/clouditor/v2/service/discovery/extra/csaf"
	"clouditor.io/clouditor/v2/service/discovery/k8s"
	"clouditor.io/clouditor/v2/service/discovery/openstack"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	ProviderAWS       = "aws"
	ProviderK8S       = "k8s"
	ProviderAzure     = "azure"
	ProviderOpenstack = "openstack"
	ProviderCSAF      = "csaf"

	// DiscovererStart is emitted at the start of a discovery run.
	DiscovererStart DiscoveryEventType = iota
	// DiscovererFinished is emitted at the end of a discovery run.
	DiscovererFinished
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
		nil,
		nil,
		WithOAuth2Authorizer(config.ClientCredentials()),
		WithTargetOfEvaluationID(viper.GetString(config.TargetOfEvaluationIDFlag)),
		WithProviders(providers),
		WithEvidenceStoreAddress(viper.GetString(config.EvidenceStoreURLFlag)),
	)
}

// DiscoveryEventType defines the event types for [DiscoveryEvent].
type DiscoveryEventType int

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

	evidenceStoreStreams *api.StreamsOf[evidence.EvidenceStore_StoreEvidencesClient, *evidence.StoreEvidenceRequest]
	evidenceStore        *api.RPCConnection[evidence.EvidenceStoreClient]

	scheduler *gocron.Scheduler

	authz service.AuthorizationStrategy

	providers   []string
	discoverers []discovery.Discoverer

	discoveryInterval time.Duration

	Events chan *DiscoveryEvent

	// ctID is the target of evaluation ID for which we are gathering resources.
	ctID string

	// collectorID is the evidence collector tool ID which is gathering the resources.
	collectorID string
}

func init() {
	log = logrus.WithField("component", "discovery")
}

// WithEvidenceStoreAddress is an option to configure the evidence store service gRPC address.
func WithEvidenceStoreAddress(target string, opts ...grpc.DialOption) service.Option[*Service] {

	return func(s *Service) {
		log.Infof("Evidence Store URL is set to %s", target)

		s.evidenceStore.Target = target
		s.evidenceStore.Opts = opts
	}
}

// WithTargetOfEvaluationID is an option to configure the target of evaluation ID for which resources will be discovered.
func WithTargetOfEvaluationID(ID string) service.Option[*Service] {
	return func(svc *Service) {
		log.Infof("Target of Evaluation ID is set to %s", ID)

		svc.ctID = ID
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
		svc.evidenceStore.SetAuthorizer(api.NewOAuthAuthorizerFromClientCredentials(config))
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
	s := &Service{
		evidenceStoreStreams: api.NewStreamsOf(api.WithLogger[evidence.EvidenceStore_StoreEvidencesClient, *evidence.StoreEvidenceRequest](log)),
		evidenceStore:        api.NewRPCConnection(string(config.DefaultEvidenceStoreURL), evidence.NewEvidenceStoreClient),
		scheduler:            gocron.NewScheduler(time.UTC),
		Events:               make(chan *DiscoveryEvent),
		ctID:                 config.DefaultTargetOfEvaluationID,
		collectorID:          config.DefaultEvidenceCollectorToolID,
		authz:                &service.AuthorizationStrategyAllowAll{},
		discoveryInterval:    5 * time.Minute, // Default discovery interval is 5 minutes
	}

	// Apply any options
	for _, o := range opts {
		o(s)
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
	svc.evidenceStoreStreams.CloseAll()
	svc.scheduler.Stop()
}

// initEvidenceStoreStream initializes the stream that is used to send evidences to the evidence store service.
// If configured, it uses the Authorizer of the discovery service to authenticate requests to the evidence store.
func (svc *Service) initEvidenceStoreStream(target string, _ ...grpc.DialOption) (stream evidence.EvidenceStore_StoreEvidencesClient, err error) {
	log.Infof("Trying to establish a connection to evidence store service @ %v", target)

	// Make sure, that we re-connect
	svc.evidenceStore.ForceReconnect()

	// Set up the stream and store it in our service struct, so we can access it later to actually
	// send the evidence data
	stream, err = svc.evidenceStore.Client.StoreEvidences(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not set up stream for storing evidences: %w", err)
	}

	log.Infof("Connected to Evidence Store″")

	return
}

// Start starts discovery
func (svc *Service) Start(ctx context.Context, req *discovery.StartDiscoveryRequest) (resp *discovery.StartDiscoveryResponse, err error) {
	var (
		optsAzure     = []azure.DiscoveryOption{}
		optsOpenstack = []openstack.DiscoveryOption{}
	)

	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check if target_of_evaluation_id in the service is within allowed or one can access *all* the target of evaluations
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
			// Add authorizer and TargetOfEvaluationID
			optsAzure = append(optsAzure, azure.WithAuthorizer(authorizer), azure.WithTargetOfEvaluationID(svc.ctID))
			// Check if resource group is given and append to discoverer
			if req.GetResourceGroup() != "" {
				optsAzure = append(optsAzure, azure.WithResourceGroup(req.GetResourceGroup()))
			}
			svc.discoverers = append(svc.discoverers, azure.NewAzureDiscovery(optsAzure...))
		case provider == ProviderK8S:
			k8sClient, err := k8s.AuthFromKubeConfig()
			if err != nil {
				log.Errorf("Could not authenticate to Kubernetes: %v", err)
				return nil, status.Errorf(codes.FailedPrecondition, "could not authenticate to Kubernetes: %v", err)
			}
			svc.discoverers = append(svc.discoverers,
				k8s.NewKubernetesComputeDiscovery(k8sClient, svc.ctID),
				k8s.NewKubernetesNetworkDiscovery(k8sClient, svc.ctID),
				k8s.NewKubernetesStorageDiscovery(k8sClient, svc.ctID))
		case provider == ProviderAWS:
			awsClient, err := aws.NewClient()
			if err != nil {
				log.Errorf("Could not authenticate to AWS: %v", err)
				return nil, status.Errorf(codes.FailedPrecondition, "could not authenticate to AWS: %v", err)
			}
			svc.discoverers = append(svc.discoverers,
				aws.NewAwsStorageDiscovery(awsClient, svc.ctID),
				aws.NewAwsComputeDiscovery(awsClient, svc.ctID))
		case provider == ProviderOpenstack:
			authorizer, err := openstack.NewAuthorizer()
			if err != nil {
				log.Errorf("Could not authenticate to OpenStack: %v", err)
				return nil, status.Errorf(codes.FailedPrecondition, "could not authenticate to OpenStack: %v", err)
			}
			// Add authorizer and TargetOfEvaluationID
			optsOpenstack = append(optsOpenstack, openstack.WithAuthorizer(authorizer), openstack.WithTargetOfEvaluationID(svc.ctID))
			svc.discoverers = append(svc.discoverers, openstack.NewOpenstackDiscovery(optsOpenstack...))
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
		e := &evidence.Evidence{
			Id:                   uuid.New().String(),
			TargetOfEvaluationId: svc.GetTargetOfEvaluationId(),
			Timestamp:            timestamppb.Now(),
			ToolId:               svc.collectorID,
			Resource:             ontology.ProtoResource(resource),
		}

		// Only enabled related evidences for some specific resources for now
		if slices.Contains(ontology.ResourceTypes(resource), "SecurityAdvisoryService") {
			edges := ontology.Related(resource)
			for _, edge := range edges {
				e.ExperimentalRelatedResourceIds = append(e.ExperimentalRelatedResourceIds, edge.Value)
			}
		}

		// Get Evidence Store stream
		channel, err := svc.evidenceStoreStreams.GetStream(svc.evidenceStore.Target, "Evidence Store", svc.initEvidenceStoreStream, svc.evidenceStore.Opts...)
		if err != nil {
			err = fmt.Errorf("could not get stream to evidence store service (%s): %w", svc.evidenceStore.Target, err)
			log.Error(err)
			continue
		}

		channel.Send(&evidence.StoreEvidenceRequest{Evidence: e})
	}
}

// GetTargetOfEvaluationId implements TargetOfEvaluationRequest for this service. This is a little trick, so that we can call
// CheckAccess directly on the service. This is necessary because the discovery service itself is tied to a specific
// target of evaluation ID, instead of the individual requests that are made against the service.
func (svc *Service) GetTargetOfEvaluationId() string {
	return svc.ctID
}
