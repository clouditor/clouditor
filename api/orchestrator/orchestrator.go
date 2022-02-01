package orchestrator

import "errors"

var (
	ErrIsNil         = errors.New("service is empty")
	ErrNameIsMissing = errors.New("service name is empty")
)

// Validate validates the cloud service
func (s *CloudService) Validate() (err error) {
	if s == nil {
		return ErrIsNil
	}
	if s.Name == "" {
		return ErrNameIsMissing
	}
	return
}
