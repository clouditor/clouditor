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
	autorest_azure "github.com/Azure/go-autorest/autorest/azure"
	"strings"
	"time"

	"clouditor.io/clouditor/service/discovery/azure"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/go-autorest/autorest/azure/auth"
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

	Configurations map[discovery.Discoverer]*Configuration
	// TODO(oxisto) do not expose this. just makes tests easier for now
	AssessmentStream assessment.Assessment_AssessEvidencesClient

	EvidenceStoreStream evidence.EvidenceStore_StoreEvidencesClient

	resources map[string]voc.IsCloudResource
	scheduler *gocron.Scheduler
}

type Configuration struct {
	Interval time.Duration
}

type ResultOntology struct {
	Result *structpb.ListValue `json:"result"`
}

func init() {
	log = logrus.WithField("component", "discovery")
}

func NewService() *Service {
	return &Service{
		resources:      make(map[string]voc.IsCloudResource),
		scheduler:      gocron.NewScheduler(time.UTC),
		Configurations: make(map[discovery.Discoverer]*Configuration),
	}
}

// Start starts discovery
func (s *Service) Start(_ context.Context, _ *discovery.StartDiscoveryRequest) (response *discovery.StartDiscoveryResponse, err error) {

	log.Infof("Start time: %s", time.Now())

	response = &discovery.StartDiscoveryResponse{Successful: true}

	log.Infof("Starting discovery...")

	s.scheduler.TagsUnique()

	// Establish connection to assessment component
	var client assessment.AssessmentClient
	// TODO(oxisto): support assessment on Another tcp/port
	cc, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not connect to assessment service: %v", err)
	}
	client = assessment.NewAssessmentClient(cc)
	// ToDo(lebogg): whats with the err?
	s.AssessmentStream, _ = client.AssessEvidences(context.Background())

	// Establish connection to evidenceStore component
	var evidenceStoreClient evidence.EvidenceStoreClient
	cc, err = grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not connect to evidence store service: %v", err)
	}
	evidenceStoreClient = evidence.NewEvidenceStoreClient(cc)
	// ToDo(lebogg): whats with the err?
	s.EvidenceStoreStream, _ = evidenceStoreClient.StoreEvidences(context.Background())

	// create an authorizer from env vars or Azure Managed Service Identity
	authorizer, err := auth.NewAuthorizerFromFile(autorest_azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		log.Errorf("Could not authenticate to Azure with authorizer from file: %v", err)
		log.Infof("Fall back to Azure authorizer from CLI.")
		authorizer, err = auth.NewAuthorizerFromCLI()
		if err != nil {
			log.Errorf("Could not authenticate to Azure authorizer from CLI: %v", err)
			return nil, err
		}
		log.Info("Using Azure authorizer from CLI. The discovery times out after 1 hour.")
	}
	log.Info("Using Azure authorizer from file.")

	//k8sClient, err := k8s.AuthFromKubeConfig()
	//if err != nil {
	//	log.Errorf("Could not authenticate to Kubernetes: %s", err)
	//	return nil, err
	//}
	//
	//awsClient, err := aws.NewClient()
	//if err != nil {
	//	log.Errorf("Could not load credentials: %s", err)
	//	return nil, err
	//}

	var discoverer []discovery.Discoverer

	discoverer = append(discoverer,
		// The azure services storage, compute and network have no auth problems.
		azure.NewAzureStorageDiscovery(azure.WithAuthorizer(authorizer)),
		azure.NewAzureComputeDiscovery(azure.WithAuthorizer(authorizer)),
		azure.NewAzureNetworkDiscovery(azure.WithAuthorizer(authorizer)),
		azure.NewAzureARMTemplateDiscovery(azure.WithAuthorizer(authorizer)),
		//k8s.NewKubernetesComputeDiscovery(k8sClient),
		//k8s.NewKubernetesNetworkDiscovery(k8sClient),
		//aws.NewAwsStorageDiscovery(awsClient),
		//aws.NewComputeDiscovery(awsClient),
	)

	for _, v := range discoverer {
		s.Configurations[v] = &Configuration{
			Interval: 5 * time.Minute,
		}

		log.Infof("Scheduling {%s} to execute every 5 minutes...", v.Name())

		_, err = s.scheduler.
			Every(20).
			Minute().
			Tag(v.Name()).
			Do(s.StartDiscovery, v)
		if err != nil {
			log.Errorf("Could not schedule job for {%s}: %v", v.Name(), err)
		}
	}

	s.scheduler.StartAsync()

	return response, nil
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
			v   *structpb.Value
			err error
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

		if s.AssessmentStream == nil {
			log.Warnf("Evidence stream to Assessment component not available")
			continue
		}

		if s.EvidenceStoreStream == nil {
			log.Warnf("Evidence stream to EvidenceStore component not available")
			continue
		}

		log.Debugf("Sending evidence for resource %s (%s)...", resource.GetID(), strings.Join(resource.GetType(), ", "))

		err = s.AssessmentStream.Send(e)
		if err != nil {
			log.Errorf("Could not send evidence to Assessment: %v", err)
		}

		err = s.EvidenceStoreStream.Send(e)
		if err != nil {
			log.Errorf("Could not send evidence to EvidenceStore: %v", err)
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
