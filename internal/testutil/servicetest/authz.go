package servicetest

import (
	"context"
	"slices"

	"clouditor.io/clouditor/v2/internal/api"
	"clouditor.io/clouditor/v2/service"
)

// NewAuthorizationStrategy contains a mock for a
// [service.AuthorizationStrategy] that either allows all target of evaluations or the
// ones that are specified in the ID list.
func NewAuthorizationStrategy(all bool, TargetOfEvaluationIDs ...string) service.AuthorizationStrategy {
	return &AuthorizationStrategyMock{
		all:                   all,
		TargetOfEvaluationIDs: TargetOfEvaluationIDs,
	}
}

type AuthorizationStrategyMock struct {
	all                   bool
	TargetOfEvaluationIDs []string
}

func (a *AuthorizationStrategyMock) CheckAccess(ctx context.Context, _ service.RequestType, req api.TargetOfEvaluationRequest) bool {
	var (
		list []string
		all  bool
	)

	all, list = a.AllowedTargetOfEvaluations(ctx)

	if all {
		return true
	}

	return slices.Contains(list, req.GetTargetOfEvaluationId())
}

func (a *AuthorizationStrategyMock) AllowedTargetOfEvaluations(_ context.Context) (all bool, IDs []string) {
	return a.all, a.TargetOfEvaluationIDs
}
