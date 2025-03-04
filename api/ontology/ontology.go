package ontology

import (
	"encoding/json"
	"google.golang.org/protobuf/reflect/protoreflect"
	"slices"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IsResource interface {
	proto.Message
	GetId() string
	GetName() string
	GetCreationTime() *timestamppb.Timestamp
	GetRaw() string
}

type HasRelatedResources interface {
	Related() []string
}

var _ IsResource = &VirtualMachine{}

type Relationship struct {
	Property string
	Value    string
}

func Related(r IsResource) []Relationship {
	var ids []Relationship

	desc := r.ProtoReflect().Descriptor()
	fields := desc.Fields()

	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)

		// We are only interested in string fields
		if field.Kind() != protoreflect.StringKind {
			continue
		}

		// TODO(oxisto): Can we maybe have a proto option on these fields instead of matching by name?
		property, found := strings.CutSuffix(string(field.Name()), "_id")
		if !found {
			// Try with _ids
			property, found = strings.CutSuffix(string(field.Name()), "_ids")
		}
		if found {
			v := r.ProtoReflect().Get(field)
			if field.IsList() {
				list := v.List()
				for i := 0; i < list.Len(); i++ {
					ids = append(ids, Relationship{
						Property: property,
						Value:    list.Get(i).String(),
					})
				}
			} else {
				s := v.String()
				if s != "" {
					ids = append(ids, Relationship{
						Property: property,
						Value:    s,
					})
				}
			}
		}
	}

	return ids
}

func ResourceTypes(r IsResource) []string {
	opts := r.ProtoReflect().Descriptor().Options()

	x := proto.GetExtension(opts, E_ResourceTypeNames)
	if types, ok := x.([]string); ok {
		return types
	}

	return nil
}

func HasType(r IsResource, typ string) bool {
	return slices.Contains(ResourceTypes(r), typ)
}

// ResourceMap contains the properties of the resource as a map[string]any, based on its JSON representation.
// This also does some magic to include the resource types in the special key "type".
func ResourceMap(r IsResource) (props map[string]any, err error) {
	var (
		b    []byte
		opts protojson.MarshalOptions
	)

	opts = protojson.MarshalOptions{
		EmitUnpopulated: true,
	}

	b, err = opts.Marshal(r)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &props)
	if err != nil {
		return nil, err
	}

	props["type"] = ResourceTypes(r)

	return
}

func ToPrettyJSON(r IsResource) (s string, err error) {
	m, err := ResourceMap(r)
	if err != nil {
		return "", err
	}

	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", err
	}

	s = string(b)
	return
}

// ResourceIDs return a list of the given resource IDs
func ResourceIDs(r []IsResource) []string {
	var a = []string{}

	for _, v := range r {
		a = append(a, v.GetId())
	}

	return a
}
