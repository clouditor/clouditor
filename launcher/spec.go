package launcher

import (
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/server"
	"clouditor.io/clouditor/v2/service"
)

type spec[T service.Service] struct {
	nsf  NewServiceFunc[T]
	wsf  WithStorageFunc[T]
	init ServiceInitFunc[T]
	opts []service.Option[T]
}

func (s spec[T]) newService(db persistence.Storage) (svc T, grpcOpts []server.StartGRPCServerOption, err error) {
	var opts = s.opts

	// Append the WithStorageFunc (if its non-nil) to the specified service options.
	if s.wsf != nil {
		opts = append(opts, s.wsf(db))
	}

	// Create the service with the NewServiceFunc using the supplied server options
	svc = s.nsf(opts...)

	// Initialize the service using the ServiceInitFunc. This returns a possible list of StartGRPCServerOptions that we need to return
	grpcOpts, err = s.init(svc)
	if err != nil {
		return *new(T), nil, err
	}

	return
}

func (s spec[T]) NewService(db persistence.Storage) (svc service.Service, grpcOpts []server.StartGRPCServerOption, err error) {
	return s.newService(db)
}

// ServiceSpec is an interface we need because of generics foo.
type ServiceSpec interface {
	NewService(db persistence.Storage) (svc service.Service, grpcOpts []server.StartGRPCServerOption, err error)
}

func NewServiceSpec[T service.Service](nsf NewServiceFunc[T], wsf WithStorageFunc[T], init ServiceInitFunc[T], opts ...service.Option[T]) ServiceSpec {
	return &spec[T]{
		nsf:  nsf,
		wsf:  wsf,
		init: init,
		opts: opts,
	}
}
