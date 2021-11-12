package evidence

import "errors"

var (
	ErrNotValidResource    = errors.New("resource in evidence is not a map")
	ErrResourceNotStruct   = errors.New("resource in evidence is not struct value")
	ErrResourceNotMap      = errors.New("resource in evidence is not a map")
	ErrResourceIdMissing   = errors.New("resource in evidence is missing the id field")
	ErrResourceIdNotString = errors.New("resource id in evidence is not a string")
	ErrToolIdMissing       = errors.New("tool id in evidence is missing")
	ErrTimestampMissing    = errors.New("timestamp in evidence is missing")
)

// Validate validates the evidence according to several required fields
func (evidence *Evidence) Validate() (resourceId string, err error) {
	if evidence.Resource == nil {
		return "", ErrNotValidResource
	}

	value := evidence.Resource.GetStructValue()
	if value == nil {
		return "", ErrResourceNotStruct
	}

	m := evidence.Resource.GetStructValue().AsMap()
	if m == nil {
		return "", ErrResourceNotMap
	}

	field, ok := m["id"]
	if !ok {
		return "", ErrResourceIdMissing
	}

	resourceId, ok = field.(string)
	if !ok {
		return "", ErrResourceIdNotString
	}

	if evidence.ToolId == "" {
		return "", ErrToolIdMissing
	}

	if evidence.Timestamp == nil {
		return "", ErrTimestampMissing
	}

	return
}
