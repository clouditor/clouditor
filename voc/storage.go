package voc

import (
	"time"

	"github.com/cayleygraph/quad"
)

type IsResource interface {
	GetID() quad.IRI
	GetName() string
	GetCreationTime() *time.Time
}

type Resource struct {
	rdfType struct{} `quad:"@type > cloud:Resource"`

	ID           quad.IRI   `json:"@id"`
	Name         string     `quad:"cloud:name"`
	CreationTime *time.Time `quad:"cloud:creationTime"`
}

func (r *Resource) GetID() quad.IRI {
	return r.ID
}

func (r *Resource) GetName() string {
	return r.Name
}

func (r *Resource) GetCreationTime() *time.Time {
	return r.CreationTime
}

type HasAtRestEncryption interface {
	GetAtRestEncryption() *AtRestEncryption
}

type HasHttpEndpoint interface {
	GetHttpEndpoint() *HttpEndpoint
}

type IsStorage interface {
	IsResource

	HasAtRestEncryption
}

type StorageResource struct {
	Resource

	rdfType          struct{}          `quad:"@type > cloud:Storage"`
	AtRestEncryption *AtRestEncryption `quad:"cloud:atRestEncryption"`
}

func (s *StorageResource) GetAtRestEncryption() *AtRestEncryption {
	return s.AtRestEncryption
}

type IsObjectStorage interface {
	IsStorage
	HasHttpEndpoint
}

type ObjectStorageResource struct {
	StorageResource

	rdfType      struct{}      `quad:"@type > cloud:ObjectStorage"`
	HttpEndpoint *HttpEndpoint `quad:"cloud:httpEndpoint"`
}

/*func (s *ObjectStorageResource) GetHttpEndpoint() *HttpEndpoint {
	return s.HttpEndpoint
}*/

type BlockStorageResource struct {
	StorageResource
}
