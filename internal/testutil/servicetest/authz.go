package servicetest

import (
	"context"
	"slices"

	"clouditor.io/clouditor/v2/internal/api"
	"clouditor.io/clouditor/v2/service"
)

// NewAuthorizationStrategy contains a mock for a
// [service.AuthorizationStrategy] that either allows all cloud services or the
// ones that are specified in the ID list.
func NewAuthorizationStrategy(all bool, CertificationTargetIDs ...string) service.AuthorizationStrategy {
	return &AuthorizationStrategyMock{
		all:                    all,
		CertificationTargetIDs: CertificationTargetIDs,
	}
}

type AuthorizationStrategyMock struct {
	all                    bool
	CertificationTargetIDs []string
}

func (a *AuthorizationStrategyMock) CheckAccess(ctx context.Context, _ service.RequestType, req api.CertificationTargetRequest) bool {
	var (
		list []string
		all  bool
	)

	all, list = a.AllowedCertificationTargets(ctx)

	if all {
		return true
	}

	return slices.Contains(list, req.GetCertificationTargetId())
}

func (a *AuthorizationStrategyMock) AllowedCertificationTargets(_ context.Context) (all bool, IDs []string) {
	return a.all, a.CertificationTargetIDs
}
