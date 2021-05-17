/*
 * Copyright 2016-2020 Fraunhofer AISEC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *           $$\                           $$\ $$\   $$\
 *           $$ |                          $$ |\__|  $$ |
 *  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
 * $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
 * $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
 * $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
 * \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
 *  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
 *
 * This file is part of Clouditor Community Edition.
 */

package discovery

import (
	"context"
	"encoding/json"
	"fmt"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/service/discovery/azure"
	"clouditor.io/clouditor/voc"
	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/query"
	"github.com/cayleygraph/cayley/query/gizmo"
	"github.com/cayleygraph/cayley/query/graphql"
	"github.com/cayleygraph/quad"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var log *logrus.Entry

//go:generate protoc -I ../../proto -I ../../third_party discovery.proto --go_out=../.. --go-grpc_out=../..

// Service is an implementation of the Clouditor Discovery service
type Service struct {
	discovery.UnimplementedDiscoveryServer
}

func init() {
	log = logrus.WithField("component", "discovery")
}

// Start starts discovery
func (s Service) Start(ctx context.Context, request *discovery.StartDiscoveryRequest) (response *discovery.StartDiscoveryResponse, err error) {
	response = &discovery.StartDiscoveryResponse{Successful: true}

	var discoverer discovery.Discoverer = azure.NewAzureStorageDiscovery()

	list, _ := discoverer.List()

	for _, v := range list {
		if err = Save(v); err != nil {
			log.Errorf("Got an error: %s", err)
		}
	}

	return response, nil
}

func (s Service) Query(ctx context.Context, request *emptypb.Empty) (response *discovery.QueryResponse, err error) {
	store := persistence.GetStore()

	// Print quads
	fmt.Println("\nquads:")
	it := store.QuadsAllIterator().Iterate()
	defer it.Close()

	for it.Next(context.TODO()) {
		fmt.Println(store.Quad(it.Result()))
	}

	var resources []voc.ObjectStorageResource
	sch := persistence.GetSchemaConfig()
	sch.LoadTo(context.TODO(), persistence.GetStore(), &resources)

	for _, v := range resources {
		fmt.Printf("%+v\n", v)
	}

	var allNamedNodes = `
	  {
		nodes {
		  id, cloud:name, rdf:type
		}
	  }`
	printQuery(allNamedNodes)

	var allEncryptedStorages = `
	g
	  .V()
	  .has("<rdf:type>", "<cloud:ObjectStorage>")
	  .tag("resource")
	  .out("<cloud:atRestEncryption>")
	  .has("<cloud:enabled>", true)
	  .tagArray()
	`
	// returns all IDs of object storages that have rest encryption enabled
	res := printGizmoQuery(allEncryptedStorages)

	b, err := json.Marshal(res)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not serialize json: %v", err)
	}

	return &discovery.QueryResponse{
		Result: string(b),
	}, nil
}

func Save(s voc.IsResource) (err error) {
	var id quad.Value

	store := persistence.GetStore()
	sch := persistence.GetSchemaConfig()

	qw := graph.NewWriter(store)

	log.Printf("saving: %+v\n", s)

	if id, err = sch.WriteAsQuads(qw, s); err != nil {
		return err
	}

	log.Printf("saved: %+v\n", id)

	err = qw.Close()

	return
}

func printQuery(s string) {
	g := graphql.NewSession(persistence.GetStore())

	it2, err := g.Execute(context.Background(), s, query.Options{})

	if err != nil {
		fmt.Printf("err: %+v", err)
	}

	for it2.Next(context.TODO()) {
		res := it2.Result()
		log.Printf("Result from query: %v", res)
	}
}

func printGizmoQuery(s string) (res interface{}) {
	g := gizmo.NewSession(persistence.GetStore())
	it2, err := g.Execute(context.Background(), s, query.Options{})

	if err != nil {
		fmt.Printf("err: %+v", err)
	}

	for it2.Next(context.TODO()) {
		res = it2.Result()
		log.Printf("Result from query: %v", res)
	}

	return
}
