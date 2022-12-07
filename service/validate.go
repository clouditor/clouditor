package service

import (
	"reflect"
	"strings"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/internal/util"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func ValidateRequest(req IncomingRequest) (err error) {
	// Check, if request is zero
	if reflect.ValueOf(req).IsZero() {
		return status.Errorf(codes.InvalidArgument, "%s", api.ErrEmptyRequest)
	}

	// Validate request
	err = req.Validate()
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	// Check, if request is a list request
	if preq, ok := req.(api.PaginatedRequest); ok {
		whitelist := util.GetFieldNames(preq)
		// Add empty string indicating no explicit ordering
		whitelist = append(whitelist, "")

		normalizedReq := strings.ToLower(preq.GetOrderBy())
		if !slices.Contains(whitelist, normalizedReq) {
			return status.Errorf(codes.InvalidArgument, "invalid request: %v", api.ErrInvalidColumnName)
		}
	}

	return nil
}

type IncomingRequest interface {
	Validate() error
	proto.Message
}
