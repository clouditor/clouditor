package ontology

import (
	"encoding/json"
	"errors"
	"slices"
	"strings"

	"clouditor.io/clouditor/v2/internal/util"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ErrNotOntologyResource = errors.New("protobuf message is not a valid ontology resource")

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

	for i := range fields.Len() {
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

// ProtoResource converts a [IsResource] to a [Resource]
func ProtoResource(resource IsResource) *Resource {
	var (
		r  Resource
		m  protoreflect.Message
		od protoreflect.OneofDescriptor
	)

	if util.IsNil(resource) {
		return nil
	}

	// Set up descriptors for proto reflection
	m = r.ProtoReflect()
	od = m.Descriptor().Oneofs().ByName("type")

	// Loop through the fields to find one that matches the resource's protobuf message type
	for i := range od.Fields().Len() {
		field := od.Fields().Get(i)
		if field.Message() == resource.ProtoReflect().Descriptor() {
			m.Set(field, protoreflect.ValueOfMessage(resource.ProtoReflect()))
			break
		}
	}

	return &r
}

// MarshalJSON is a custom JSON marshaller for the [Resource] type that delegates JSON marshalling to the [protojson]
// package.
func (x *Resource) MarshalJSON() (b []byte, err error) {
	return protojson.Marshal(x)
}

// UnmarshalJSON is a custom JSON unmarshaller for the [Resource] type that delegates JSON unmarshalling to the
// [protojson] package.
func (x *Resource) UnmarshalJSON(b []byte) (err error) {
	return protojson.Unmarshal(b, x)
}
