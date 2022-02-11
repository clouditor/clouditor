package inmemory

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"clouditor.io/clouditor/persistence"
)

type inMemory struct {
	tables map[string]*table
}

type table struct {
	entries map[interface{}]interface{}
}

var ErrNoPrimaryKey = errors.New("could not find primary key")
var ErrNoPointerType = errors.New("argument must be pointer type")
var ErrConditionNotSupported = errors.New("condition not supported ")

// NewStorage creates a new in memory storage
// For now with GORM. Use own implementation in the future
func NewStorage() (persistence.Storage, error) {
	return &inMemory{tables: make(map[string]*table)}, nil
}

func (i *inMemory) Create(r interface{}) error {
	t := i.table(r)

	id, err := id(r)
	if err != nil {
		return ErrNoPrimaryKey
	}

	// TODO(oxisto): Throw an error, if it already exists?
	t.entries[id] = r

	return nil
}

func (i *inMemory) Count(r interface{}, conds ...interface{}) (count int64, err error) {
	t := i.table(r)

	// TODO(oxisto): Evaluate conditions
	count = int64(len(t.entries))

	return
}

func (i *inMemory) Delete(r interface{}, conds ...interface{}) error {
	return nil
}

func (i *inMemory) Get(r interface{}, conds ...interface{}) error {
	var s = reflect.ValueOf(r)
	if s.Kind() != reflect.Ptr {
		return ErrNoPointerType
	}

	s = s.Elem()

	var ok bool

	if len(conds) == 1 {
		// Simple primary ID selection
		t := i.table(r)

		var v interface{}
		if v, ok = t.entries[conds[0]]; !ok {
			return persistence.ErrRecordNotFound
		}

		// TODO(oxisto): I think we can get rid of that with 1.18 semantics
		s.Set(reflect.ValueOf(v).Elem())

		return nil
	} else if len(conds) == 2 {
		// Very hacky for now
		if strings.ToLower(conds[0].(string)) == "id = ?" {
			return i.Get(r, conds[1])
		}

		return ErrConditionNotSupported
	}

	return nil
}

func (i *inMemory) List(r interface{}, conds ...interface{}) error {
	// TODO(oxisto): Use conditions
	t := i.table(r)

	var s = reflect.ValueOf(r)
	if s.Kind() != reflect.Ptr {
		return ErrNoPointerType
	}

	s = s.Elem()

	for _, value := range t.entries {
		v := reflect.ValueOf(value)

		if v.Kind() == reflect.Ptr && reflect.TypeOf(r).Elem().Elem().Kind() != reflect.Ptr {
			s = reflect.Append(s, reflect.ValueOf(value).Elem())
		} else if v.Kind() == reflect.Ptr {
			s = reflect.Append(s, reflect.ValueOf(value))
		} else {
			s = reflect.Append(s, reflect.ValueOf(value).Addr())
		}
	}

	reflect.ValueOf(r).Elem().Set(s)

	return nil
}

func (i *inMemory) Update(r interface{}, conds ...interface{}) error {
	return nil
}

func (i *inMemory) table(r interface{}) (t *table) {
	var ok bool

	t, ok = i.tables[typeKey(r)]
	if !ok {
		t = &table{entries: make(map[interface{}]interface{})}
		i.tables[typeKey(r)] = t
	}

	return t
}

func typeKey(r interface{}) string {
	var t = fmt.Sprintf("%T", r)

	// Remove any pointers or arrays to normalize the type
	s := strings.ReplaceAll(t, "*", "")
	s = strings.ReplaceAll(s, "[]", "")
	return s
}

func id(r interface{}) (interface{}, error) {
	v := reflect.ValueOf(r).Elem()
	var found = false
	v2 := v.FieldByNameFunc(func(s string) bool {
		if len(s) == 0 || found {
			return false
		}

		r := []rune(s)

		if unicode.IsUpper(r[0]) && unicode.IsLetter(r[0]) {
			found = true
			return true
		}

		return false
	})

	if !v2.IsValid() || !found {
		return "", errors.New("could not find field for primary key")
	}

	return v2.Interface(), nil
}
