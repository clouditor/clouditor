// Copyright 2023 Fraunhofer AISEC
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

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (svc *Service) ListGraphEdges(ctx context.Context, req *discovery.ListGraphEdgesRequest) (res *discovery.ListGraphEdgesResponse, err error) {
	var (
		results []*discovery.Resource
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

	all, allowed = svc.authz.AllowedCertificationTargets(ctx)
	if !all {
		query = append(query, "certification_target_id IN ?")
		args = append(args, allowed)
	}

	res = new(discovery.ListGraphEdgesResponse)

	// This is a little problematic, since we are actually paginating the underlying resources and not the edges, but it
	// is probably the best we can do for now while we are not storing the edges in the database.
	results, res.NextPageToken, err = service.PaginateStorage[*discovery.Resource](req,
		svc.storage,
		service.DefaultPaginationOpts,
		persistence.BuildConds(query, args)...,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	// Loop through all resources and find edges to others
	for _, resource := range results {
		r, _ := resource.ToOntologyResource()
		if r == nil {
			continue
		}

		for _, rel := range ontology.Related(r) {
			edge := &discovery.GraphEdge{
				Id:     resource.Id + "-" + rel.Value,
				Source: resource.Id,
				Target: rel.Value,
				Type:   rel.Property,
			}

			res.Edges = append(res.Edges, edge)
		}
	}

	return
}

func (svc *Service) UpdateResource(ctx context.Context, req *discovery.UpdateResourceRequest) (res *discovery.Resource, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	if !svc.authz.CheckAccess(ctx, service.AccessRead, req) {
		return nil, service.ErrPermissionDenied
	}

	err = svc.storage.Save(&req.Resource, "id = ?", req.Resource.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	res = req.Resource

	return
}
