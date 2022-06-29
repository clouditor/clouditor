package util

import (
	"errors"
	"reflect"
)

var (
	// ErrNoStruct indicates the passed argument is not a struct
	ErrNoStruct = errors.New("no struct")
)

// GetFieldNames extracts all field names of struct T. Returns error if T is no struct.
// TODO(lebogg): Only take exported fields
func GetFieldNames[T any]() (fieldNames []string, err error) {
	var aStruct T
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
