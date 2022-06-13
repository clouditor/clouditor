package api

import (
	"errors"

	"google.golang.org/protobuf/proto"
	"k8s.io/utils/strings/slices"

	"clouditor.io/clouditor/internal/util"
)

var (
	ErrInvalidColumnName = errors.New("column name is invalid")
	ErrRequestIsNil      = errors.New("request is empty")
)

// ListRequest indicates a proto message being a request for ListXXX RPCs
type ListRequest interface {
	GetOrderBy() string
	GetAsc() bool
	proto.Message
}

func ValidateListReq(req ListRequest, responseType any) (err error) {
	// req must be non-nil
	if req == nil {
		err = ErrRequestIsNil
		return
	}

	// Avoid DB injections by whitelisting the valid orderBy statements
	whitelist, err := util.GetFieldNames(responseType)
	// Add empty string indicating no explicit ordering
	whitelist = append(whitelist, "")
	if !slices.Contains(whitelist, req.GetOrderBy()) {
		err = ErrInvalidColumnName
		return
	}

	return

}
