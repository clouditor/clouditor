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

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/persistence/inmemory"
	"clouditor.io/clouditor/service"
)

var log *logrus.Entry

// Service is an implementation of the Clouditor req service (evidenceServer)
type Service struct {
	// Currently only in-memory
	storage persistence.Storage

	// evidenceHooks is a list of hook functions that can be used if one wants to be
	// informed about each evidence
	evidenceHooks []evidence.EvidenceHookFunc
	// mu is used for (un)locking result hook calls
	mu sync.Mutex

	evidence.UnimplementedEvidenceStoreServer
}

func NewService(opts ...service.Option[Service]) (svc *Service) {
	var (
		err error
	)
	svc = new(Service)

	for _, o := range opts {
		o(svc)
	}

	if svc.storage == nil {
		svc.storage, err = inmemory.NewStorage()
		if err != nil {
			log.Errorf("Could not initialize the storage: %v", err)
		}
	}
	return
}

func init() {
	log = logrus.WithField("component", "Evidence Store")
}

// StoreEvidence is a method implementation of the evidenceServer interface: It receives a req and stores it
func (s *Service) StoreEvidence(_ context.Context, req *evidence.StoreEvidenceRequest) (resp *evidence.StoreEvidenceResponse, err error) {

	_, err = req.Evidence.Validate()
	if err != nil {
		err = fmt.Errorf("invalid evidence: %w", err)
		log.Error(err)

		go s.informHooks(nil, err)

		resp = &evidence.StoreEvidenceResponse{
			Status:        false,
			StatusMessage: err.Error(),
		}

		return resp, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	err = s.storage.Create(req.Evidence)
	if err != nil {
		// TODO(lebogg): Almost identical error handling as above -> extract here or even to internal
		err = fmt.Errorf("internal error: %w", err)
		log.Error(err)

		go s.informHooks(nil, err)

		resp = &evidence.StoreEvidenceResponse{
			Status:        false,
			StatusMessage: err.Error(),
		}

		err = status.Errorf(codes.Internal, "%v", err)
	}
	go s.informHooks(req.Evidence, nil)

	resp = &evidence.StoreEvidenceResponse{
		Status: true,
	}

	log.Debugf("Evidence stored with id: %v", req.Evidence.Id)

	return resp, nil
}

// StoreEvidences is a method implementation of the evidenceServer interface: It receives evidences and stores them
func (s *Service) StoreEvidences(stream evidence.EvidenceStore_StoreEvidencesServer) (err error) {
	var (
		req *evidence.StoreEvidenceRequest
		res *evidence.StoreEvidenceResponse
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
		res, err = s.StoreEvidence(context.Background(), evidenceRequest)
		if err != nil {
			log.Errorf("Error storing evidence: %v", err)
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
	res = new(evidence.ListEvidencesResponse)

	// Paginate the evidences according to the request
	res.Evidences, res.NextPageToken, err = service.PaginateStorage[*evidence.Evidence](req, s.storage,
		service.DefaultPaginationOpts)
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
