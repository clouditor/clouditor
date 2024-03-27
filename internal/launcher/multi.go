package launcher

import (
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/server"
	"clouditor.io/clouditor/v2/service"
)

type ServiceCreator interface {
	NewService(db persistence.Storage) (svc any, grpcOpts []server.StartGRPCServerOption, err error)
}

func NewServiceSpec[T any](nsf NewServiceFunc[T], wsf WithStorageFunc[T], init ServiceInitFunc[T], opts ...service.Option[T]) *spec[T] {
	return &spec[T]{
		nsf:  nsf,
		wsf:  wsf,
		init: init,
		opts: opts,
	}
}

type MultiLauncher struct {
	services []any

	Launcher[any]
}

/*type NewLauncherFunc[T any] func() (l *Launcher[T], err error)

func NewMultiLauncher(nlfs ...NewLauncherFunc[any]) (*MultiLauncher, error) {
	ml := &MultiLauncher{}

	for _, nlf := range nlfs {
		l, err := nlf()
		if err != nil {
			return nil, err
		}

		ml.launchers = append(ml.launchers, l)
	}

	return ml, nil
}*/

func NewMultiLauncher(funcs ...ServiceCreator) (*MultiLauncher, error) {
	ml := &MultiLauncher{}

	ml.component = "all-in-one"
	ml.initLogging()
	ml.initStorage()

	for _, f := range funcs {
		svc, grpcOpts, err := f.NewService(ml.db)
		if err != nil {
			return nil, err
		}
		ml.grpcOpts = append(ml.grpcOpts, grpcOpts...)
		ml.services = append(ml.services, svc)
	}

	return ml, nil
}
