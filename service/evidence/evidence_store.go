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
	"errors"
	"fmt"
	"io"
	"sync"

	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	resp = &evidence.StoreEvidenceResponse{}

	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	s.evidences[req.Evidence.Id] = req.Evidence

	go s.informHooks(req.Evidence, nil)

	log.Debugf("Evidence stored with id: %v", req.Evidence.GetId())

	return resp, nil
}

// StoreEvidences is a method implementation of the evidenceServer interface: It receives evidences and stores them
func (s *Service) StoreEvidences(stream evidence.EvidenceStore_StoreEvidencesServer) (err error) {
	var (
		req *evidence.StoreEvidenceRequest
		res *evidence.StoreEvidencesResponse
	)

	for {
		req, err = stream.Recv()

		// If no more input of the stream is available, return
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			newError := fmt.Errorf("cannot receive stream request: %w", err)
			log.Error(newError)
			return status.Errorf(codes.Unknown, "%v", newError)
		}

		// Call StoreEvidence() for storing a single evidence
		evidenceRequest := &evidence.StoreEvidenceRequest{
			Evidence: req.Evidence,
		}
		_, err = s.StoreEvidence(context.Background(), evidenceRequest)
		if err != nil {
			log.Errorf("Error storing evidence: %v", err)
			// Create response message. The StoreEvidence method does not need that message, so we have to create it here for the stream response.
			res = &evidence.StoreEvidencesResponse{
				Status:        false,
				StatusMessage: err.Error(),
			}
		} else {
			res = &evidence.StoreEvidencesResponse{
				Status: true,
			}
		}

		// Send response back to the client
		err = stream.Send(res)

		// Check for send errors
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			newError := fmt.Errorf("cannot send response to the client: %w", err)
			log.Error(newError)
			return status.Errorf(codes.Unknown, "%v", newError)
		}
	}
}

// ListEvidences is a method implementation of the evidenceServer interface: It returns the evidences lying in the req storage
func (s *Service) ListEvidences(_ context.Context, req *evidence.ListEvidencesRequest) (res *evidence.ListEvidencesResponse, err error) {
	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	res = new(evidence.ListEvidencesResponse)

	// Paginate the evidences according to the request
	res.Evidences, res.NextPageToken, err = service.PaginateMapValues(req, s.evidences, func(a, b *evidence.Evidence) bool {
		return a.Id < b.Id
	}, service.DefaultPaginationOpts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
	}

	return
}

func (s *Service) RegisterEvidenceHook(evidenceHook evidence.EvidenceHookFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
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
