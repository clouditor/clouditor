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

	"clouditor.io/clouditor/api/evaluation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

// Service is an implementation of the Clouditor Evaluation service
type Service struct {
	evaluation.UnimplementedEvaluationServer
	scheduler *gocron.Scheduler
}

func init() {
	log = logrus.WithField("component", "orchestrator")
}

// ServiceOption is a function-style option to configure the Evaluation Service
type ServiceOption func(*Service)

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
		Do(s.StartEvaluation)

	return
}

func (s *Service) StartEvaluation(req *evaluation.EvaluateRequest) {
	log.Infof("Started evaluation for Cloud Service '%s' and Catalog ID '%s'", req.Toe.CloudServiceId, req.Toe.CatalogId)
}

func (s *Service) Shutdown() {
	s.scheduler.Stop()
}
