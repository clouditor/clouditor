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

package evidences

import (
	"context"
	"fmt"
	"io"
	"sync"

	"clouditor.io/clouditor/api/evidence"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var log *logrus.Entry

// Service is an implementation of the Clouditor req service (evidenceServer)
type Service struct {
	// Currently only in-memory
	evidences map[string]*evidence.Evidence

	// evidenceHooks is a list of hook functions that can be used if one wants to be
	// informed about each evidence
	evidenceHooks []evidence.EvidenceHookFunc
	// mu is used for (un)locking result hook calls
	mu sync.Mutex

	evidence.UnimplementedEvidenceStoreServer
}

func NewService() *Service {
	return &Service{
		evidences: make(map[string]*evidence.Evidence),
	}
}

func init() {
	log = logrus.WithField("component", "Evidence Store")
}

// StoreEvidence is a method implementation of the evidenceServer interface: It receives a req and stores it
func (s *Service) StoreEvidence(_ context.Context, req *evidence.StoreEvidenceRequest) (resp *evidence.StoreEvidenceResponse, err error) {

	_, err = req.Evidence.Validate()
	if err != nil {
		log.Errorf("Invalid evidence: %v", err)
		newError := fmt.Errorf("invalid evidence: %w", err)

		go s.informHooks(nil, newError)

		resp = &evidence.StoreEvidenceResponse{
			Status: false,
		}

		return resp, status.Errorf(codes.InvalidArgument, "invalid req: %v", err)
	}

	s.evidences[req.Evidence.Id] = req.Evidence
	go s.informHooks(req.Evidence, nil)

	resp = &evidence.StoreEvidenceResponse{
		Status: true,
	}

	log.Infof("Evidence stored with id: %v", req.Evidence.Id)

	return resp, nil
}

// StoreEvidences is a method implementation of the evidenceServer interface: It receives evidences and stores them
func (s *Service) StoreEvidences(stream evidence.EvidenceStore_StoreEvidencesServer) (err error) {
	var req *evidence.StoreEvidenceRequest

	for {
		req, err = stream.Recv()

		if err != nil {
			// If no more input of the stream is available, return SendAndClose `error`
			if err == io.EOF {
				log.Infof("Stopped receiving streamed evidences")
				return stream.SendAndClose(&emptypb.Empty{})
			}

			return err
		}

		// Call StoreEvidence() for storing a single evidence
		evidenceRequest := &evidence.StoreEvidenceRequest{
			Evidence: req.Evidence,
		}
		_, err = s.StoreEvidence(context.Background(), evidenceRequest)
		if err != nil {
			return err
		}
	}
}

// ListEvidences is a method implementation of the evidenceServer interface: It returns the evidences lying in the req storage
func (s *Service) ListEvidences(_ context.Context, _ *evidence.ListEvidencesRequest) (*evidence.ListEvidencesResponse, error) {
	var listOfEvidences []*evidence.Evidence
	for _, v := range s.evidences {
		listOfEvidences = append(listOfEvidences, v)
	}

	return &evidence.ListEvidencesResponse{Evidences: listOfEvidences}, nil
}

func (s *Service) RegisterEvidenceHook(evidenceHook evidence.EvidenceHookFunc) {
	s.evidenceHooks = append(s.evidenceHooks, evidenceHook)
}

func (s *Service) informHooks(result *evidence.Evidence, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Inform our hook, if we have any
	if s.evidenceHooks != nil {
		for _, hook := range s.evidenceHooks {
			// TODO(all): We could do hook concurrent again (assuming different hooks don't interfere with each other)
			hook(result, err)
		}
	}
}
