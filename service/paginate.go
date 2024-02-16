// Copyright 2022 Fraunhofer AISEC
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

package service

import (
	"fmt"
	"sort"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/persistence"
)

// PaginationOpts can be used to fine-tune the pagination, especially with regards to the page sizes. This can be important
// if the messages within a page are extremely large and thus the page size needs to be decreased.
type PaginationOpts struct {
	// DefaultPageSize is the page size that is used as a default if the request does not specify one
	DefaultPageSize int32

	// MaxPageSize is the maximum page size that can be requested
	MaxPageSize int32
}

// DefaultPaginationOpts are sensible defaults for the pagination size.
var DefaultPaginationOpts = PaginationOpts{
	DefaultPageSize: 50,
	MaxPageSize:     1500,
}

// PaginateSlice is a helper function that helps to paginate a slice based on list requests. It parses the necessary
// information out if a paginated request, e.g. the page token and the desired page size and returns a sliced page as
// well as the next page token.
func PaginateSlice[T any](req api.PaginatedRequest, values []T, less func(a T, b T) bool, opts PaginationOpts) (page []T, npt string, err error) {
	sort.Slice(values, func(i, j int) bool {
		return less(values[i], values[j])
	})

	return paginate(req, opts, func(start int64, size int32) (page []T, done bool, err error) {
		var (
			end, max int64
		)

		// Clamp the end to the maximum of the slice
		end = start + int64(size)
		max = int64(len(values))
		if end >= max {
			end = max

			// Indicate that we are at the end
			done = true
		}

		page = values[start:end]
		return
	})
}

// PaginateStorage is a helper function that helps to paginate records in persisted storage based on list requests. It
// parses the necessary information out if a paginated request, e.g. the page token and the desired page size and
// returns a sliced page as well as the next page token.
func PaginateStorage[T any](req api.PaginatedRequest, storage persistence.Storage, opts PaginationOpts,
	conds ...interface{}) (page []T, npt string, err error) {
	return paginate(req, opts, func(start int64, size int32) (page []T, done bool, err error) {
		// Retrieve values from the DB
		err = storage.List(&page, req.GetOrderBy(), req.GetAsc(), int(start), int(size), conds...)
		if err != nil {
			return nil, true, fmt.Errorf("database error: %w", err)
		}

		if len(page) == 0 || len(page) < int(size) {
			// Indicate that we are at the end
			done = true
		}
		return
	})
}

// paginate takes cares of the heavy lifting of handling the actual pagination request. It takes the paginated request
// req, calculates offsets and sizes, which can be fine-tuned using opts and supplies them to the pager function. The
// pager function needs to return the actual page contents based on the calculated size and offset. This result is then
// returned to the caller as well as a token that can be used to request the next page.
func paginate[T any](req api.PaginatedRequest, opts PaginationOpts, pager func(start int64, size int32) (page []T, done bool, err error)) (page []T, npt string, err error) {
	var (
		token *api.PageToken
		size  int32
		done  bool
	)

	// Check, if the size was specified and is within our maximum size
	if req.GetPageSize() == 0 {
		size = opts.DefaultPageSize
	} else if req.GetPageSize() > opts.MaxPageSize {
		size = opts.MaxPageSize
	} else {
		size = req.GetPageSize()
	}

	// Check, if this is the first request (empty token) or a subsequent one
	if req.GetPageToken() == "" {
		// We need a new page token
		token = &api.PageToken{
			Start: 0,
			Size:  size,
		}
	} else {
		// Try to decode our existing token
		token, err = api.DecodePageToken(req.GetPageToken())
		if err != nil {
			return nil, "", fmt.Errorf("could not decode page token: %w", err)
		}
	}

	// Call our pager function with the offset and size
	page, done, err = pager(token.Start, size)
	if err != nil {
		// Transparently return the error
		return nil, "", err
	}

	if !done {
		// Move the token "forward"
		token.Start = token.Start + int64(len(page))

		// Encode next page token
		npt, err = token.Encode()
		if err != nil {
			return nil, "", fmt.Errorf("could not create page token: %w", err)
		}
	}

	return
}
