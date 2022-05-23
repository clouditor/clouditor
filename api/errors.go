package api

import (
	"errors"

	"google.golang.org/grpc/status"
)

// StatusFromWrappedError attempts to recover a status.Status from a wrapped error message. This is necessary, because
// the underlying RPC error might be wrapped in multiple layers of additional error messages and status.FromError only
// checks the error itself.
func StatusFromWrappedError(err error) (s *status.Status, ok bool) {
	for {
		s, ok = status.FromError(err)
		if ok {
			return
		}
		err = errors.Unwrap(err)
		if err == nil {
			return nil, false
		}
	}
}
