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
	"encoding/base64"
	"fmt"

	"google.golang.org/protobuf/proto"
)

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
