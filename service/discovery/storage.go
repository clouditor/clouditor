package discovery

import "time"

type Resource interface {
	ID() string
	Name() string
	CreationTime() *time.Time
}

type resource struct {
	id           string
	name         string
	creationTime *time.Time
}

func (r *resource) ID() string {
	return r.id
}

func (r *resource) Name() string {
	return r.name
}

func (r *resource) CreationTime() *time.Time {
	return r.creationTime
}

type Storage interface {
	Resource

	AtRestEncryption() *AtRestEncryption
}

type storage struct {
	resource

	atRestEncryption *AtRestEncryption
}

func (s *storage) AtRestEncryption() *AtRestEncryption {
	return s.atRestEncryption
}

type ObjectStorage interface {
	Storage

	HttpEndpoint() *HttpEndpoint
}

type objectStorage struct {
	storage

	httpEndpoint *HttpEndpoint
}

func (s *objectStorage) HttpEndpoint() *HttpEndpoint {
	return s.httpEndpoint
}

type BlockStorage interface {
	Storage
}

/*type blockStorage struct {
	storage
}*/

type StorageDiscoverer interface {
	List() ([]Storage, error)
}
