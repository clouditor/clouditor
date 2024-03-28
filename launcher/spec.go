// Copyright 2024 Fraunhofer AISEC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//           $$\                           $$\ $$\   $$\
//           $$ |                          $$ |\__|  $$ |
//  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
// $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
// $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
// $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
// \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
//  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
//
// This file is part of Clouditor Community Edition.

package launcher

import (
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/server"
	"clouditor.io/clouditor/v2/service"
)

// spec is a struct that implements [ServiceSpec]. We want to keep it unexported, because it needs to be generic to the
// type T, in order for the initialization and creation functions to work. But if we want to have a list of different
// specs, we cannot mix the generics, therefore we need to have the [ServiceSpec] interface.
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
