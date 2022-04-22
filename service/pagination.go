package service

import (
	"fmt"

	"clouditor.io/clouditor/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaginatedRequest interface {
	GetPageToken() string
	GetPageSize() int32
}

// Paginate is a helper function that helps to paginate list requests. It parses the necessary
// informaton out if a paginated request, e.g. the page token and the desired page size and returns
// the offset values for the requested page as well as the next page token.
func Paginate(req PaginatedRequest, max int64) (start int64, end int64, nbt string, err error) {
	var token *api.PageToken

	if req.GetPageToken() == "" {
		// We need a new page token
		token = &api.PageToken{
			Start: 0,
			Size:  req.GetPageSize(),
		}
	} else {
		// Try to decode our existing token
		token, err = api.DecodePageToken(req.GetPageToken())
		if err != nil {
			return 0, 0, "", fmt.Errorf("could not decode page token: %w", err)
		}
	}

	start = token.Start
	end = token.Start + int64(token.Size)
	if end >= max {
		end = max

		// Indicate that we are at the end
		token = nil
	}

	// Only needed, if more pages exist
	if token != nil {
		// Move the token "forward"
		token.Start = end

		// Encode next page token
		nbt, err = token.Encode()
		if err != nil {
			return 0, 0, "", status.Errorf(codes.Internal, "could not create page token: %v", err)
		}
	}

	return
}
