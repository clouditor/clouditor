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

package api

import (
	"context"
	"encoding/base64"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// PaginatedRequest contains the typical parameters for a paginated request,
// usually a request for a List gRPC call.
type PaginatedRequest interface {
	GetPageToken() string
	GetPageSize() int32
	proto.Message
}

// PaginatedResponse contains the typical parameters for a paginated response,
// usually a response for a List gRPC call.
type PaginatedResponse interface {
	GetNextPageToken() string
}

func (t *PageToken) Encode() (b64token string, err error) {
	var b []byte

	b, err = proto.Marshal(t)
	if err != nil {
		return "", fmt.Errorf("error while marshaling protobuf message: %w", err)
	}

	b64token = base64.URLEncoding.EncodeToString(b)
	return
}

func DecodePageToken(b64token string) (t *PageToken, err error) {
	var b []byte

	b, err = base64.URLEncoding.DecodeString(b64token)
	if err != nil {
		return nil, fmt.Errorf("error while decoding base64 token: %w", err)
	}

	t = new(PageToken)

	err = proto.Unmarshal(b, t)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling protobuf message: %w", err)
	}

	return
}

// ListAllPaginated issues a List gRPC call that supports pagination and combines all results of all requests into a
// single slice.
func ListAllPaginated[ResponseType PaginatedResponse, RequestType PaginatedRequest, ResultType any](req RequestType, list func(context.Context, RequestType, ...grpc.CallOption) (ResponseType, error), getter func(res ResponseType) []ResultType) (results []ResultType, err error) {
	var (
		res       ResponseType
		pageToken string = ""
	)

	for {
		// Modify the request to include our page token using protoreflect. This will be empty for the first page
		m := req.ProtoReflect()
		m.Set(m.Descriptor().Fields().ByName("page_token"), protoreflect.ValueOf(pageToken))

		// Call the list function to fetch the next page
		res, err = list(context.Background(), req)
		if err != nil {
			// Transparently return the error of the list function without any wrapping
			return nil, err
		}

		// Append results and retrieve our next page token
		results = append(results, getter(res)...)
		pageToken = res.GetNextPageToken()

		// If the page token is empty, there are no more pages left to fetch
		if pageToken == "" {
			break
		}
	}

	return
}
