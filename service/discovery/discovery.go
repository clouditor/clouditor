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
	"time"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/service/discovery/azure"
	"clouditor.io/clouditor/service/standalone"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

var log *logrus.Entry

//go:generate protoc -I ../../proto -I ../../third_party discovery.proto --go_out=../.. --go-grpc_out=../.. --openapi_out=../../openapi/discovery

// Service is an implementation of the Clouditor Discovery service.
// It should not be used directly, but rather the NewService constructor
// should be used.
type Service struct {
	discovery.UnimplementedDiscoveryServer

	Configurations map[discovery.Discoverer]*DiscoveryConfiguration
	// TODO(oxisto) do not expose this. just makes tests easier for now
	AssessmentStream assessment.Assessment_StreamEvidencesClient

	resources map[string]voc.IsResource
	scheduler *gocron.Scheduler
}

type DiscoveryConfiguration struct {
	Interval time.Duration
}

func init() {
	log = logrus.WithField("component", "discovery")
}

func NewService() *Service {
	return &Service{
		resources:      make(map[string]voc.IsResource),
		scheduler:      gocron.NewScheduler(time.UTC),
		Configurations: make(map[discovery.Discoverer]*DiscoveryConfiguration),
	}
}

// Start starts discovery
func (s Service) Start(ctx context.Context, request *discovery.StartDiscoveryRequest) (response *discovery.StartDiscoveryResponse, err error) {
	response = &discovery.StartDiscoveryResponse{Successful: true}

	log.Infof("Starting discovery...")

	s.scheduler.TagsUnique()

	var isStandalone bool = true

	var client assessment.AssessmentClient

	if isStandalone {
		client = standalone.NewAssessmentClient()
	} else {
		// TODO(oxisto): support assessment on another tcp/port
		cc, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not connect to assessment service: %v", err)
		}
		client = assessment.NewAssessmentClient(cc)
	}

	s.AssessmentStream, _ = client.StreamEvidences(context.Background())

	// create an authorizer from env vars or Azure Managed Service Identity
	authorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		log.Errorf("Could not authenticate to Azure: %s", err)
		return nil, err
	}

	var discoverer []discovery.Discoverer

	discoverer = append(discoverer,
		azure.NewAzureStorageDiscovery(azure.WithAuthorizer(authorizer)),
		azure.NewAzureComputeDiscovery(azure.WithAuthorizer(authorizer)),
	)

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
	list, _ := discoverer.List()

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

		evidence := &assessment.Evidence{
			Resource:          v,
			ResourceId:        resource.GetID(),
			ApplicableMetrics: []int32{1},
		}

		if s.AssessmentStream == nil {
			log.Warnf("Evidence stream not available")
			continue
		}

		err = s.AssessmentStream.Send(evidence)
		if err != nil {
			log.Errorf("Could not send evidence: %v", err)
		}
	}
}

func (s Service) Query(ctx context.Context, request *emptypb.Empty) (response *discovery.QueryResponse, err error) {
	var r []*structpb.Value

	for _, v := range s.resources {
		var s *structpb.Value

		s, err = voc.ToStruct(v)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error during JSON unmarshal: %v", err)
		}

		r = append(r, s)
	}

	return &discovery.QueryResponse{
		Result: &structpb.ListValue{Values: r},
	}, nil
}
