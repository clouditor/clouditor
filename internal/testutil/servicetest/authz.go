package servicetest

import (
	"context"

	"clouditor.io/clouditor/internal/api"
	"clouditor.io/clouditor/service"
	"golang.org/x/exp/slices"
)

func NewAuthorizationStrategy(all bool, cloudServiceIDs []string) service.AuthorizationStrategy {
	return &AuthorizationStrategyMock{
		all:             all,
		cloudServiceIDs: cloudServiceIDs,
	}
}

type AuthorizationStrategyMock struct {
	all             bool
	cloudServiceIDs []string
}

func (a *AuthorizationStrategyMock) CheckAccess(ctx context.Context, _ service.RequestType, req api.CloudServiceRequest) bool {
	var (
		list []string
		all  bool
	)

	all, list = a.AllowedCloudServices(ctx)

	if all {
		return true
	}

	return slices.Contains(list, req.GetCloudServiceId())
}

func (a *AuthorizationStrategyMock) AllowedCloudServices(_ context.Context) (all bool, IDs []string) {
	return a.all, a.cloudServiceIDs
}
