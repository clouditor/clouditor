// Copyright 2016-2022 Fraunhofer AISEC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//	         $$\                           $$\ $$\   $$\
//	         $$ |                          $$ |\__|  $$ |
//	$$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
//
// $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
// $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
// $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
// \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
//
//	\_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
//
// This file is part of Clouditor Community Edition.
package gorm

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/testutil/servicetest/orchestratortest"
	"clouditor.io/clouditor/persistence"
)

var mockMetricRange = &assessment.Range{Range: &assessment.Range_MinMax{MinMax: &assessment.MinMax{Min: 1, Max: 2}}}

func TestStorageOptions(t *testing.T) {
	type args struct {
		opts []StorageOption
	}
	tests := []struct {
		name              string
		args              args
		wantDialectorType string
		wantErr           assert.ErrorAssertionFunc
	}{
		{
			name: "in memory with option",
			args: args{
				opts: []StorageOption{
					WithInMemory(),
				},
			},
			wantDialectorType: "*sqlite.Dialector",
			wantErr:           nil,
		},
		{
			name:              "in memory without option",
			wantDialectorType: "*sqlite.Dialector",
			wantErr:           nil,
		},
		{
			name: "postgres with option - invalid port",
			args: args{
				opts: []StorageOption{
					WithPostgres("", 0, "", "", "", ""),
				},
			},
			wantDialectorType: "",
			wantErr: func(tt assert.TestingT, err error, i ...any) bool {
				return assert.Contains(t, err.Error(), "invalid port")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewStorage(tt.args.opts...)
			if (err != nil) && tt.wantErr != nil {
				tt.wantErr(t, err, tt.args.opts)
				return
			}

			gorm, ok := s.(*storage)
			assert.True(t, ok)
			assert.NotNil(t, gorm)

			assert.Equal(t, tt.wantDialectorType, fmt.Sprintf("%T", gorm.dialector))

			// Test to create a new cloud service and get it again with
			// respective 'Create' and 'Get' Create record via DB call
			serviceInput := &orchestrator.CloudService{Name: "SomeName"}
			err = s.Create(serviceInput)
			assert.NoError(t, err)

			// Get record via DB call
			serviceOutput := &orchestrator.CloudService{}
			err = s.Get(serviceOutput, "name = ?", "SomeName")
			assert.NoError(t, err)
			assert.NoError(t, serviceOutput.Validate())
			assert.Equal(t, serviceInput, serviceOutput)
		})
	}
}

func Test_storage_Create(t *testing.T) {
	var (
		err    error
		s      persistence.Storage
		metric *assessment.Metric
	)

	metric = &assessment.Metric{
		Id:    testdata.MockMetricID,
		Name:  testdata.MockMetricName,
		Range: mockMetricRange,
	}
	// Check if metric has all necessary fields
	assert.NoError(t, metric.Validate())

	// Create storage
	s, err = NewStorage()
	assert.NoError(t, err)

	err = s.Create(metric)
	assert.NoError(t, err)

	err = s.Create(metric)
	assert.Error(t, err)
}

func Test_storage_Get(t *testing.T) {
	var (
		err     error
		s       persistence.Storage
		service *orchestrator.CloudService
	)

	service = orchestratortest.NewCloudService()

	// Create storage
	s, err = NewStorage()
	assert.NoError(t, err)

	// Return error since no record in the DB yet
	err = s.Get(&orchestrator.CloudService{})
	assert.ErrorIs(t, err, persistence.ErrRecordNotFound)

	// Create service
	err = s.Create(service)
	assert.NoError(t, err)

	// Get service via passing entire record
	gotService := &orchestrator.CloudService{}
	err = s.Get(gotService)
	assert.NoError(t, err)
	assert.NoError(t, gotService.Validate())
	assert.Equal(t, service, gotService)

	// Get service via name
	gotService2 := &orchestrator.CloudService{}
	err = s.Get(gotService2, "name = ?", service.Name)
	assert.NoError(t, err)
	assert.Equal(t, service, gotService2)

	// Get service via description
	gotService3 := &orchestrator.CloudService{}
	err = s.Get(gotService3, "description = ?", service.Description)
	assert.NoError(t, err)
	assert.NoError(t, gotService3.Validate())
	assert.Equal(t, service, gotService3)

	var metric = &assessment.Metric{
		Id:    testdata.MockMetricID,
		Name:  testdata.MockMetricName,
		Range: mockMetricRange,
	}
	// Check if metric has all necessary fields
	assert.NoError(t, metric.Validate())

	// Create metric
	err = s.Create(metric)
	assert.NoError(t, err)

	// Get metric via Id
	gotMetric := &assessment.Metric{}
	err = s.Get(gotMetric, "id = ?", testdata.MockMetricID)
	assert.NoError(t, err)
	assert.NoError(t, gotMetric.Validate())
	assert.Equal(t, metric, gotMetric)

	var impl = &assessment.MetricImplementation{
		MetricId:  testdata.MockMetricID,
		Code:      "TestCode",
		UpdatedAt: timestamppb.New(time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)),
	}
	// Check if impl has all necessary fields
	assert.NoError(t, impl.Validate())

	// Create metric implementation
	err = s.Create(impl)
	assert.NoError(t, err)

	// Get metric implementation via Id
	gotImpl := &assessment.MetricImplementation{}
	err = s.Get(gotImpl, "metric_id = ?", testdata.MockMetricID)
	assert.NoError(t, err)
	assert.NoError(t, gotImpl.Validate())
	assert.Equal(t, impl, gotImpl)
}

func Test_storage_List(t *testing.T) {
	var (
		err      error
		s        persistence.Storage
		service1 *orchestrator.CloudService
		service2 *orchestrator.CloudService
		services []orchestrator.CloudService
	)

	// Create storage
	s, err = NewStorage()
	assert.NoError(t, err)

	// Test service
	service1 = &orchestrator.CloudService{Id: testdata.MockCloudServiceID, Name: testdata.MockCloudServiceName}
	service2 = &orchestrator.CloudService{Id: testdata.MockAnotherCloudServiceID, Name: testdata.MockAnotherCloudServiceName}

	// List should return empty list since no services are in DB yet
	err = s.List(&services, "", true, 0, -1)
	assert.ErrorIs(t, err, nil)
	assert.Empty(t, services)

	// List should return list of 2 services (service1 and service2)
	err = s.Create(service1)
	assert.NoError(t, err)
	err = s.Create(service2)
	assert.NoError(t, err)
	err = s.List(&services, "", true, 0, -1)
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, 2, len(services))
	// We only check one service and assume the others are also correct
	assert.NoError(t, services[0].Validate())

	// Test with certificates (associations included via states)
	var (
		certificate1 *orchestrator.Certificate
		certificate2 *orchestrator.Certificate
		certificates []*orchestrator.Certificate
	)

	// List should return empty list since no certificates are in DB yet
	err = s.List(&certificates, "", true, 0, 0)
	assert.ErrorIs(t, err, nil)
	assert.Empty(t, certificates)

	// Create two certificates
	certificate1 = orchestratortest.NewCertificate()
	certificate1.Id = "0"
	certificate2 = orchestratortest.NewCertificate()
	certificate2.Id = "1"
	err = s.Create(certificate1)
	assert.NoError(t, err)
	err = s.Create(certificate2)
	assert.NoError(t, err)

	// List should return list of 2 certificates with associated states
	err = s.List(&certificates, "id", false, 0, -1)
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, 2, len(certificates))
	// Check ordering
	assert.Equal(t, certificate2.Id, certificates[0].Id)
	// We only check one certificate and assume the others are also correct
	assert.NoError(t, certificates[0].Validate())

	fmt.Println(certificates)

	// Check if certificate with id "1" (certificate2) is in the list and if states are included (association)
	for i := range certificates {
		if certificates[i].Id == certificate2.Id {
			fmt.Println("Certificate:", certificates[i])
			assert.NotEmpty(t, certificates[i].States)
			return
		}
	}
	// If not, let the test fail
	assert.FailNow(t, "%s is not listed but should be.", certificate1.Id)

}

func Test_storage_Count(t *testing.T) {
	var (
		count    int64
		err      error
		s        persistence.Storage
		service1 *orchestrator.CloudService
		service2 *orchestrator.CloudService
	)

	service1 = orchestratortest.NewCloudService()
	service2 = orchestratortest.NewAnotherCloudService()

	// Create storage
	s, err = NewStorage()
	assert.NoError(t, err)

	// Since no records in DB yet, count of services should be 0
	count, err = s.Count(&orchestrator.CloudService{})
	assert.NoError(t, err)
	assert.Equal(t, int(count), 0)

	// Create one service -> count of services should be 1
	assert.ErrorIs(t, s.Create(service1), nil)
	count, err = s.Count(&orchestrator.CloudService{})
	assert.NoError(t, err)
	assert.Equal(t, 1, int(count))

	// Add another one -> count of services should be 2
	assert.ErrorIs(t, s.Create(service2), nil)
	count, err = s.Count(&orchestrator.CloudService{})
	assert.NoError(t, err)
	assert.Equal(t, 2, int(count))

	// Count of services with ID "SomeName2" should be 1
	count, err = s.Count(&orchestrator.CloudService{}, "name = ?", testdata.MockAnotherCloudServiceName)
	assert.NoError(t, err)
	assert.Equal(t, 1, int(count))

	// Calling s.Count() with unsupported record type should throw "unsupported" error
	_, err = s.Count(nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported data type")
}

func Test_storage_Save(t *testing.T) {
	type MyTest struct {
		ID   int
		Name string
	}

	var (
		err        error
		s          persistence.Storage
		service    *orchestrator.CloudService
		newService *orchestrator.CloudService
		gotService *orchestrator.CloudService
		myVar      MyTest
	)
	service = orchestratortest.NewCloudService()

	// Create storage
	s, err = NewStorage(WithAdditionalAutoMigration(&MyTest{}))
	assert.NoError(t, err)

	// Create service
	err = s.Create(service)
	assert.NoError(t, err)

	err = s.Get(&orchestrator.CloudService{}, "name = ?", service.Name)
	assert.NoError(t, err)

	// Save new service: Change description. Name and ID remain unchanged
	newService = orchestratortest.NewCloudService()
	newService.Description = ""

	err = s.Save(newService, "name = ?", service.Name)
	assert.NoError(t, err)

	gotService = &orchestrator.CloudService{}
	err = s.Get(gotService, "name = ?", service.Name)
	assert.NoError(t, err)
	assert.NoError(t, gotService.Validate())

	// Name should be the same
	assert.Equal(t, service.Name, gotService.Name)
	// Description should be zero
	assert.Equal(t, "", gotService.Description)

	// Save MyTest
	myVar = MyTest{ID: 1, Name: "Test"}

	err = s.Save(&myVar)
	assert.NoError(t, err)
}

func Test_storage_Update(t *testing.T) {
	var (
		err error
		s   persistence.Storage
	)

	// Create storage
	s, err = NewStorage()
	assert.NoError(t, err)

	// Testing cloud service
	// Create cloud service
	cloudService := orchestrator.CloudService{
		Id:          testdata.MockCloudServiceID,
		Name:        testdata.MockCloudServiceName,
		Description: testdata.MockCloudServiceDescription,
		ConfiguredMetrics: []*assessment.Metric{
			{
				Id:    testdata.MockCloudServiceID,
				Name:  testdata.MockMetricName,
				Range: mockMetricRange,
			},
		},
	}
	// Check if cloud service has all necessary fields
	assert.NoError(t, cloudService.Validate())
	err = s.Create(&cloudService)
	assert.NoError(t, err)

	err = s.Get(&orchestrator.CloudService{}, "Id = ?", cloudService.Id)
	assert.NoError(t, err)

	err = s.Update(&orchestrator.CloudService{Name: "SomeNewName", Description: ""}, "Id = ?", cloudService.Id)
	assert.NoError(t, err)

	gotCloudService := &orchestrator.CloudService{}
	err = s.Get(gotCloudService, "Id = ?", cloudService.Id)
	assert.NoError(t, err)
	assert.NoError(t, gotCloudService.Validate())

	// Name should be changed
	assert.Equal(t, "SomeNewName", gotCloudService.Name)
	// Other properties should stay the same
	assert.Equal(t, cloudService.Id, gotCloudService.Id)
	assert.Equal(t, cloudService.Description, gotCloudService.Description)
	assert.Equal(t, len(cloudService.ConfiguredMetrics), len(gotCloudService.ConfiguredMetrics))
}

func Test_storage_Delete(t *testing.T) {
	var (
		err     error
		s       persistence.Storage
		service *orchestrator.CloudService
	)
	service = orchestratortest.NewCloudService()

	// Create storage
	s, err = NewStorage()
	assert.NoError(t, err)

	// Create service
	err = s.Create(service)
	assert.NoError(t, err)

	// Should return ErrRecordNotFound since there is no service "Fake" in DB
	assert.ErrorIs(t, s.Delete(&orchestrator.CloudService{}, "name = ?", "Fake"), persistence.ErrRecordNotFound)

	// Successful deletion
	assert.Nil(t, s.Delete(&orchestrator.CloudService{}, "name = ?", service.Name))
	// Check with s.Get that service is not in DB anymore
	assert.ErrorIs(t, s.Get(&orchestrator.CloudService{}, "name = ?", service.Name), persistence.ErrRecordNotFound)

	// Should return DB error since a non-supported type is passed (just a string instead of, e.g., &orchestrator.CloudService{})
	assert.Contains(t, s.Delete("Unsupported Type").Error(), "unsupported data type")
}
