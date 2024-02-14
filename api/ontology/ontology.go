package ontology

import (
	"encoding/json"
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

		// TODO(oxisto): Can we maybe have a proto option on these fields instead of matching by name?
		property, found := strings.CutSuffix(string(field.Name()), "_id")
		if found {
			v := r.ProtoReflect().Get(field)
			vv := v.Interface()
			if vvv, ok := vv.(string); ok {
				// Make sure, the value is really set
				if vvv != "" {
					ids = append(ids, Relationship{
						Property: string(property),
						Value:    vvv,
					})
				}
			} else if vvvv, ok := vv.([]string); ok {
				for _, vvv := range vvvv {
					ids = append(ids, Relationship{
						Property: string(property),
						Value:    vvv,
					})
				}
			}
		}
	}

	return ids
}

func ResourceTypes(r IsResource) []string {
	opts := r.ProtoReflect().Descriptor().Options()

	x := proto.GetExtension(opts, E_ResourceTypeName)
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
		EmitDefaultValues: true,
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
