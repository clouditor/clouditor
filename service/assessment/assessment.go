// Copyright 2021 Fraunhofer AISEC
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

package assessment

import (
	"fmt"
	"io"
	"os"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/policies"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
)

var log *logrus.Entry

//go:generate protoc -I ../../proto -I ../../third_party assessment.proto evidence.proto --go_out=../.. --go-grpc_out=../..  --openapi_out=../../openapi/assessment

func init() {
	log = logrus.WithField("component", "assessment")
}

type Service struct {
	ResultHook func(result *assessment.Result, err error)

	results map[string]*assessment.Result
	assessment.UnimplementedAssessmentServer
}

func NewService() assessment.AssessmentServer {
	return &Service{
		results: make(map[string]*assessment.Result),
	}

}

func (s Service) StreamEvidences(stream assessment.Assessment_StreamEvidencesServer) error {
	var evidence *assessment.Evidence
	var err error

	for {
		evidence, err = stream.Recv()
		if err == io.EOF {
			log.Infof("Stopped receiving streamed evidence")

			return stream.SendAndClose(&emptypb.Empty{})
		}

		log.Infof("Received evidence for resource %s", evidence.ResourceId)
		log.Debugf("Evidence: %+v", evidence)

		var file string

		// TODO(oxisto): actually look up metric via orchestrator
		if len(evidence.ApplicableMetrics) != 0 {
			metric := evidence.ApplicableMetrics[0]

			var baseDir string = "."

			// check, if we are in the root of Clouditor
			if _, err := os.Stat("policies"); os.IsNotExist(err) {
				// in tests, we are relative to our current package
				baseDir = "../../"
			}

			file = fmt.Sprintf("%s/policies/metric%d.rego", baseDir, metric)
		} else {
			log.Errorf("Could not find a valid metric for evidence of resource %s", evidence.ResourceId)
		}

		// TODO(oxisto): use go embed
		data, err := policies.Run(file, evidence)
		if err != nil {
			log.Errorf("Could not evaluate evidence: %v", err)

			if s.ResultHook != nil {
				s.ResultHook(nil, err)
			}

			return err
		}

		log.Infof("Evaluated evidence as %v", data["compliant"])

		result := &assessment.Result{
			ResourceId: evidence.ResourceId,
			Compliant:  data["compliant"].(bool),
		}

		s.results[evidence.ResourceId] = result

		if s.ResultHook != nil {
			s.ResultHook(result, nil)
		}
	}
}
