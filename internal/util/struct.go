package util

import (
	"errors"
	"reflect"
)

var (
	// ErrNoStruct indicates the passed argument is not a struct
	ErrNoStruct = errors.New("no struct")
	// ErrStructIsNil indicates the passed argument is nil
	ErrStructIsNil = errors.New("struct is nil")
)

// GetFieldNames extracts all field names of aStruct
// TODO(lebogg): Only take exported fields
func GetFieldNames(aStruct any) (fieldNames []string, err error) {
	// Check aStruct isn't nil
	if aStruct == nil {
		err = ErrStructIsNil
		return
	}
	// Check aStruct is a struct
	if reflect.TypeOf(aStruct).Kind() != reflect.Struct {
		err = ErrNoStruct
		return
	}
	// Get all fields of aStruct and add their names to fieldNames
	t := reflect.TypeOf(aStruct)
	fields := reflect.VisibleFields(t)
	for _, f := range fields {
		fieldNames = append(fieldNames, f.Name)
	}
	return
}
