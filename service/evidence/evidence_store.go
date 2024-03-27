// Copyright 2016-2023 Fraunhofer AISEC
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
	"slices"
	"sync"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/internal/launcher"
	"clouditor.io/clouditor/v2/internal/logging"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/persistence/inmemory"
	"clouditor.io/clouditor/v2/server"
	"clouditor.io/clouditor/v2/service"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var log *logrus.Entry

var DefaultServiceSpec = launcher.NewServiceSpec(
	NewService,
	WithStorage,
	func(svc *Service) ([]server.StartGRPCServerOption, error) {
		// It is possible to register hook functions for the evidenceStore.
		//  * The hook functions in evidenceStore are implemented in StoreEvidence(s)

		// evidenceStoreService.RegisterEvidenceHook(func(result *evidence.Evidence, err error) {})

		return []server.StartGRPCServerOption{
			server.WithEvidenceStore(svc),
		}, nil
	},
)

// Service is an implementation of the Clouditor req service (evidenceServer)
type Service struct {
	storage persistence.Storage

	// evidenceHooks is a list of hook functions that can be used if one wants to be
	// informed about each evidence
	evidenceHooks []evidence.EvidenceHookFunc
	// mu is used for (un)locking result hook calls
	mu sync.Mutex

	// authz defines our authorization strategy, e.g., which user can access which cloud service and associated
	// resources, such as evidences and assessment results.
	authz service.AuthorizationStrategy

	evidence.UnimplementedEvidenceStoreServer
}

func WithStorage(storage persistence.Storage) service.Option[*Service] {
	return func(svc *Service) {
		svc.storage = storage
	}
}

func NewService(opts ...service.Option[*Service]) (svc *Service) {
	var (
		err error
	)
	svc = new(Service)

	for _, o := range opts {
		o(svc)
	}

	// Default to an allow-all authorization strategy
	if svc.authz == nil {
		svc.authz = &service.AuthorizationStrategyAllowAll{}
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
func (svc *Service) StoreEvidence(ctx context.Context, req *evidence.StoreEvidenceRequest) (res *evidence.StoreEvidenceResponse, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the cloud service according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessUpdate, req) {
		return nil, service.ErrPermissionDenied
	}

	err = svc.storage.Create(req.Evidence)
	if err != nil && errors.Is(err, persistence.ErrUniqueConstraintFailed) {
		return nil, status.Error(codes.AlreadyExists, "entry already exists")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	go svc.informHooks(ctx, req.Evidence, nil)

	res = &evidence.StoreEvidenceResponse{}

	logging.LogRequest(log, logrus.DebugLevel, logging.Store, req)

	return res, nil
}

// StoreEvidences is a method implementation of the evidenceServer interface: It receives evidences and stores them
func (svc *Service) StoreEvidences(stream evidence.EvidenceStore_StoreEvidencesServer) (err error) {
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
		_, err = svc.StoreEvidence(stream.Context(), evidenceRequest)
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

// ListEvidences is a method implementation of the evidenceServer interface: It returns the evidences lying in the storage
func (svc *Service) ListEvidences(ctx context.Context, req *evidence.ListEvidencesRequest) (res *evidence.ListEvidencesResponse, err error) {
	var (
		all     bool
		allowed []string
		query   []string
		args    []any
	)
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Retrieve list of allowed cloud service according to our authorization strategy. No need to specify any additional
	// conditions to our storage request, if we are allowed to see all cloud services.
	all, allowed = svc.authz.AllowedCloudServices(ctx)
	if !all && req.GetFilter().GetCloudServiceId() != "" && !slices.Contains(allowed, req.GetFilter().GetCloudServiceId()) {
		return nil, service.ErrPermissionDenied
	}

	res = new(evidence.ListEvidencesResponse)

	// Apply filter options
	if filter := req.GetFilter(); filter != nil {
		if cloudServiceId := filter.GetCloudServiceId(); cloudServiceId != "" {
			query = append(query, "cloud_service_id = ?")
			args = append(args, cloudServiceId)
		}
		if toolId := filter.GetToolId(); toolId != "" {
			query = append(query, "tool_id = ?")
			args = append(args, toolId)
		}
	}

	// In any case, we need to make sure that we only select evidences of cloud services that we have access to
	if !all {
		query = append(query, "cloud_service_id IN ?")
		args = append(args, allowed)
	}

	// Paginate the evidences according to the request
	res.Evidences, res.NextPageToken, err = service.PaginateStorage[*evidence.Evidence](req, svc.storage,
		service.DefaultPaginationOpts, persistence.BuildConds(query, args)...)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
	}

	return
}

// GetEvidence is a method implementation of the evidenceServer interface: It returns a particular evidence in the storage
func (svc *Service) GetEvidence(ctx context.Context, req *evidence.GetEvidenceRequest) (res *evidence.Evidence, err error) {
	var (
		all     bool
		allowed []string
		conds   []any
	)

	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Retrieve list of allowed cloud service according to our authorization strategy. No need to specify any additional
	// conditions to our storage request, if we are allowed to see all cloud services.
	all, allowed = svc.authz.AllowedCloudServices(ctx)
	if !all {
		conds = []any{"id = ? AND cloud_service_id IN ?", req.EvidenceId, allowed}
	} else {
		conds = []any{"id = ?", req.EvidenceId}
	}

	res = new(evidence.Evidence)

	err = svc.storage.Get(res, conds...)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "evidence not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	return
}

func (svc *Service) RegisterEvidenceHook(evidenceHook evidence.EvidenceHookFunc) {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	svc.evidenceHooks = append(svc.evidenceHooks, evidenceHook)
}

func (svc *Service) informHooks(ctx context.Context, result *evidence.Evidence, err error) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	// Inform our hook, if we have any
	if svc.evidenceHooks != nil {
		for _, hook := range svc.evidenceHooks {
			// TODO(all): We could do hook concurrent again (assuming different hooks don't interfere with each other)
			hook(ctx, result, err)
		}
	}
}
