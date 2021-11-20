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
	"strings"
	"time"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/service/discovery/aws"
	"clouditor.io/clouditor/service/discovery/azure"
	"clouditor.io/clouditor/service/discovery/k8s"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var log *logrus.Entry

// Service is an implementation of the Clouditor Discovery service.
// It should not be used directly, but rather the NewService constructor
// should be used.
type Service struct {
	discovery.UnimplementedDiscoveryServer

	Configurations map[discovery.Discoverer]*DiscoveryConfiguration
	// TODO(oxisto) do not expose this. just makes tests easier for now
	AssessmentStream assessment.Assessment_AssessEvidencesClient

	EvidenceStoreStream evidence.EvidenceStore_StoreEvidencesClient

	// resources holds an in-memory map of the discovered cloud resources keyed by
	// the resource ID
	resources map[string]voc.IsCloudResource

	scheduler *gocron.Scheduler

	serviceId string
}

// ServiceOption implements functional-style options to configure the Service.
type ServiceOption func(*Service)

func WithAzure(authorizer autorest.Authorizer) ServiceOption {
	return func(s *Service) {
		log.Info("Found an azure account. Configuring our discoverers accordingly")

		/*var authorizer autorest.Authorizer

		authorizer, err := auth.NewAuthorizerFromCLI()
		if err != nil {
			log.Errorf("Could not authenticate to Azure: %s", err)
			//return err
		}*/

		/*s.Configurations[]
		*discoverer = append(*discoverer,
			azure.NewAzureIacTemplateDiscovery(azure.WithAuthorizer(authorizer)),
			azure.NewAzureStorageDiscovery(azure.WithAuthorizer(authorizer)),
			azure.NewAzureComputeDiscovery(azure.WithAuthorizer(authorizer)),
			azure.NewAzureNetworkDiscovery(azure.WithAuthorizer(authorizer)),
		)*/

		discoverers := []discovery.Discoverer{
			azure.NewAzureIacTemplateDiscovery(azure.WithAuthorizer(authorizer)),
			azure.NewAzureStorageDiscovery(azure.WithAuthorizer(authorizer)),
			azure.NewAzureComputeDiscovery(azure.WithAuthorizer(authorizer)),
			azure.NewAzureNetworkDiscovery(azure.WithAuthorizer(authorizer)),
		}

		log.Debugf("Discoverers: %+v", discoverers)
	}

}

type DiscoveryConfiguration struct {
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
		Configurations: make(map[discovery.Discoverer]*DiscoveryConfiguration),
	}
}

// Start starts discovery
func (s *Service) Start(_ context.Context, _ *discovery.StartDiscoveryRequest) (response *discovery.StartDiscoveryResponse, err error) {
	response = &discovery.StartDiscoveryResponse{Successful: true}

	log.Infof("Starting discovery...")

	// Clear any previous discoverers, if there are any
	s.scheduler.Clear()
	s.scheduler.TagsUnique()

	// Establish connection to assessment component
	var client assessment.AssessmentClient
	var orchestratorClient orchestrator.OrchestratorClient

	// TODO(oxisto): support assessment on Another tcp/port
	cc, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not connect to assessment service: %v", err)
	}
	client = assessment.NewAssessmentClient(cc)
	orchestratorClient = orchestrator.NewOrchestratorClient(cc)

	accountsResponse, _ := orchestratorClient.ListAccounts(context.Background(), &orchestrator.ListAccountsRequests{})

	s.AssessmentStream, err = client.AssessEvidences(context.Background())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not enable connection to assessment service: %v", err)
	}

	// Establish connection to evidenceStore component
	var evidenceStoreClient evidence.EvidenceStoreClient
	cc, err = grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not connect to evidence store service: %v", err)
	}
	evidenceStoreClient = evidence.NewEvidenceStoreClient(cc)

	s.EvidenceStoreStream, err = evidenceStoreClient.StoreEvidences(context.Background())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not enable connection to evidence store: %v", err)
	}

	var discoverer []discovery.Discoverer

	// Look for an azure account
	for _, account := range accountsResponse.Accounts {
		if account.AccountType == "azure" {
			err = handleAzureAccount(account, &discoverer)
		} else if account.AccountType == "aws" {
			err = handleAWSAccount(account, &discoverer)
		} else if account.AccountType == "k8s" {
			err = handlek8sAccount(account, &discoverer)
		}

		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not handle account: %v", err)
		}
	}

	for _, v := range discoverer {
		s.Configurations[v] = &DiscoveryConfiguration{
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
		var s *structpb.Value

		if filteredType != "" && !v.HasType(filteredType) {
			continue
		}

		s, err = voc.ToStruct(v)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error during JSON unmarshal: %v", err)
		}

		r = append(r, s)
	}

	return &discovery.QueryResponse{
		Results: &structpb.ListValue{Values: r},
	}, nil
}

// Assign assigns this discovery service to a target cloud service
func (s *Service) Assign(_ context.Context, req *discovery.AssignRequest) (response *discovery.AssignResponse, err error) {
	response = new(discovery.AssignResponse)

	s.serviceId = req.ServiceId

	return
}

/*func saveResourcesToFilesystem(result ResultOntology, filename string) error {
	var (
		filepath string
	)

	prefix, indent := "", "    "
	exported, err := json.MarshalIndent(result, prefix, indent)
	if err != nil {
		return fmt.Errorf("marshalling JSON failed %w", err)
	}

	filepath = "../../results/discovery_results/"

	// Check if folder exists
	err = os.MkdirAll(filepath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("check for directory existence failed:  %w", err)
	}

	err = ioutil.WriteFile(filepath+filename, exported, 0666)
	if err != nil {
		return fmt.Errorf("write file failed %w", err)
	} else {
		fmt.Println("ontology resources written to: ", filepath+filename)
	}

	return nil
}*/

func handleAzureAccount(account *orchestrator.CloudAccount, discoverer *[]discovery.Discoverer) (err error) {
	log.Info("Found an azure account. Configuring our discoverers accordingly")

	var authorizer autorest.Authorizer

	if account.AutoDiscover {
		// create an authorizer from env vars or Azure Managed Service Identity
		authorizer, err = auth.NewAuthorizerFromCLI()
		if err != nil {
			log.Errorf("Could not authenticate to Azure: %s", err)
			return err
		}
	} else {
		authorizer, err = auth.NewClientCredentialsConfig(account.Principal, account.Secret, account.AccountId).Authorizer()
		if err != nil {
			log.Errorf("Could not authenticate to Azure: %s", err)
			return err
		}
	}

	*discoverer = append(*discoverer,
		azure.NewAzureIacTemplateDiscovery(azure.WithAuthorizer(authorizer)),
		azure.NewAzureStorageDiscovery(azure.WithAuthorizer(authorizer)),
		azure.NewAzureComputeDiscovery(azure.WithAuthorizer(authorizer)),
		azure.NewAzureNetworkDiscovery(azure.WithAuthorizer(authorizer)),
	)

	return nil
}

func handleAWSAccount(account *orchestrator.CloudAccount, discoverer *[]discovery.Discoverer) (err error) {
	awsClient, err := aws.NewClient()
	if err != nil {
		log.Errorf("Could not load credentials: %s", err)
		return err
	}

	*discoverer = append(*discoverer,
		aws.NewAwsStorageDiscovery(awsClient),
		aws.NewComputeDiscovery(awsClient),
	)

	return nil
}

func handlek8sAccount(account *orchestrator.CloudAccount, discoverer *[]discovery.Discoverer) (err error) {
	k8sClient, err := k8s.AuthFromKubeConfig()
	if err != nil {
		log.Errorf("Could not authenticate to Kubernetes: %s", err)
		return err
	}

	*discoverer = append(*discoverer,
		k8s.NewKubernetesComputeDiscovery(k8sClient),
		k8s.NewKubernetesNetworkDiscovery(k8sClient),
	)

	return nil
}
