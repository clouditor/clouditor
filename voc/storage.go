package voc

import (
	"time"
)

type IsResource interface {
	GetID() string
	GetName() string
	GetCreationTime() *time.Time
}

type Resource struct {
	ID           string
	Name         string
	CreationTime int64
}

func (r *Resource) GetID() string {
	return r.ID
}

func (r *Resource) GetName() string {
	return r.Name
}

func (r *Resource) GetCreationTime() *time.Time {
	t := time.Unix(r.CreationTime, 0)
	return &t
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

	AtRestEncryption *AtRestEncryption
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

	HttpEndpoint *HttpEndpoint
}

type BlockStorageResource struct {
	StorageResource
}
