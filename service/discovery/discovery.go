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
	"fmt"
	"time"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/service/discovery/azure"
	"clouditor.io/clouditor/service/discovery/k8s"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"

	"clouditor.io/clouditor/service/discovery/aws"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/voc"
	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

var log *logrus.Entry

// Service is an implementation of the Clouditor Discovery service.
// It should not be used directly, but rather the NewService constructor
// should be used.
type Service struct {
	discovery.UnimplementedDiscoveryServer

	configurations map[discovery.Discoverer]*Configuration

	assessmentStream  assessment.Assessment_AssessEvidencesClient
	assessmentAddress string

	resources map[string]voc.IsCloudResource
	scheduler *gocron.Scheduler

	authorizer api.Authorizer
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
func WithAssessmentAddress(address string) ServiceOption {
	return func(s *Service) {
		s.assessmentAddress = address
	}
}

// WithInternalAuthorizer is an option to use an authorizer to the internal Clouditor auth service.
func WithInternalAuthorizer(address string, username string, password string, opts ...grpc.DialOption) ServiceOption {
	return func(s *Service) {
		s.SetAuthorizer(api.NewInternalAuthorizerFromPassword(address, username, password, opts...))
	}
}

func NewService(opts ...ServiceOption) *Service {
	s := &Service{
		assessmentAddress: DefaultAssessmentAddress,
		resources:         make(map[string]voc.IsCloudResource),
		scheduler:         gocron.NewScheduler(time.UTC),
		configurations:    make(map[discovery.Discoverer]*Configuration),
	}

	// Apply any options
	for _, o := range opts {
		o(s)
	}

	return s
}

// SetAuthorizer implements UsesAuthorizer.
func (s *Service) SetAuthorizer(auth api.Authorizer) {
	s.authorizer = auth
}

// Authorizer implements UsesAuthorizer.
func (s *Service) Authorizer() api.Authorizer {
	return s.authorizer
}

// initAssessmentStream initializes the stream that is used to send evidences to the assessment service.
// If configured, it uses the Authorizer of the discovery service to authenticate requests to the assessment.
func (s *Service) initAssessmentStream(additionalOpts ...grpc.DialOption) error {
	log.Infof("Trying to establish a connection to assessment service @ %v", s.assessmentAddress)

	// Establish connection to assessment gRPC service
	conn, err := grpc.Dial(s.assessmentAddress,
		api.DefaultGrpcDialOptions(s, additionalOpts...)...,
	)
	if err != nil {
		return fmt.Errorf("could not connect to assessment service: %w", err)
	}

	client := assessment.NewAssessmentClient(conn)

	// Set up the stream and store it in our service struct, so we can access it later to actually
	// send the evidence data
	s.assessmentStream, err = client.AssessEvidences(context.Background())
	if err != nil {
		return fmt.Errorf("could not set up stream for assessing evidences: %w", err)
	}

	log.Infof("Connected to Assessment")

	return nil
}

// Start starts discovery
func (s *Service) Start(_ context.Context, _ *discovery.StartDiscoveryRequest) (resp *discovery.StartDiscoveryResponse, err error) {
	resp = &discovery.StartDiscoveryResponse{Successful: true}

	log.Infof("Starting discovery...")
	s.scheduler.TagsUnique()

	// Establish connection to assessment component
	if s.assessmentStream == nil {
		err = s.initAssessmentStream()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not initialize stream to Assessment: %v", err)
		}
	}

	authorizer, err := azure.NewAuthorizer()
	if err != nil {
		log.Errorf("Could not authenticate to Azure: %v", err)
		return nil, status.Errorf(codes.FailedPrecondition, "could not authenticate to Azure: %v", err)
	}

	k8sClient, err := k8s.AuthFromKubeConfig()
	if err != nil {
		log.Errorf("Could not authenticate to Kubernetes: %v", err)
		return nil, status.Errorf(codes.FailedPrecondition, "could not authenticate to Kubernetes: %v", err)
	}

	awsClient, err := aws.NewClient()
	if err != nil {
		log.Errorf("Could not load credentials: %v", err)
		return nil, err
	}

	var discoverer []discovery.Discoverer

	discoverer = append(discoverer,
		azure.NewAzureARMTemplateDiscovery(azure.WithAuthorizer(authorizer)),
		azure.NewAzureStorageDiscovery(azure.WithAuthorizer(authorizer)),
		azure.NewAzureComputeDiscovery(azure.WithAuthorizer(authorizer)),
		azure.NewAzureNetworkDiscovery(azure.WithAuthorizer(authorizer)),
		k8s.NewKubernetesComputeDiscovery(k8sClient),
		k8s.NewKubernetesNetworkDiscovery(k8sClient),
		aws.NewAwsStorageDiscovery(awsClient),
		aws.NewAwsComputeDiscovery(awsClient),
	)

	for _, v := range discoverer {
		s.configurations[v] = &Configuration{
			Interval: 5 * time.Minute,
		}

		log.Infof("Scheduling {%s} to execute every 5 minutes...", v.Name())

		_, err = s.scheduler.
			Every(5).
			Minute().
			Tag(v.Name()).
			Do(s.StartDiscovery, v)
		if err != nil {
			log.Errorf("Could not schedule job for {%s}: %v", v.Name(), err)
		}
	}

	s.scheduler.StartAsync()

	return resp, nil
}

func (s Service) Shutdown() {
	s.scheduler.Stop()
}

func (s Service) StartDiscovery(discoverer discovery.Discoverer) {
	var (
		err  error
		list []voc.IsCloudResource
	)

	list, err = discoverer.List()

	if err != nil {
		log.Errorf("Could not retrieve resources from discoverer '%s': %v", discoverer.Name(), err)
		return
	}

	for _, resource := range list {
		s.resources[string(resource.GetID())] = resource

		var (
			v *structpb.Value
		)

		v, err = voc.ToStruct(resource)
		if err != nil {
			log.Errorf("Could not convert resource to protobuf struct: %v", err)
		}

		// TODO(all): What is the raw type in our case?
		e := &evidence.Evidence{
			Id:        uuid.New().String(),
			Timestamp: timestamppb.Now(),
			ToolId:    "Clouditor Evidences Collection",
			Raw:       "",
			Resource:  v,
		}

		if s.assessmentStream == nil {
			log.Warnf("Evidence stream to Assessment component not available")
			continue
		}

		if err = s.assessmentStream.Send(&assessment.AssessEvidenceRequest{Evidence: e}); err != nil {
			log.WithError(handleError(err, "Assessment"))
		}
	}
}

func (s Service) Query(_ context.Context, request *discovery.QueryRequest) (response *discovery.QueryResponse, err error) {
	var r []*structpb.Value

	var filteredType = ""
	if request != nil {
		filteredType = request.FilteredType
	}

	for _, v := range s.resources {
		var resource *structpb.Value

		if filteredType != "" && !v.HasType(filteredType) {
			continue
		}

		resource, err = voc.ToStruct(v)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error during JSON unmarshal: %v", err)
		}

		r = append(r, resource)
	}

	return &discovery.QueryResponse{
		Results: &structpb.ListValue{Values: r},
	}, nil
}

// handleError prints out the error according to the status code
func handleError(err error, dest string) error {
	prefix := "could not send evidence to " + dest
	if status.Code(err) == codes.Internal {
		return fmt.Errorf("%s. Internal error on the server side: %v", prefix, err)
	} else if status.Code(err) == codes.InvalidArgument {
		return fmt.Errorf("invalid evidence - provide evidence in the right format: %v", err)
	} else {
		return err
	}
}
