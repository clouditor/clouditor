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
