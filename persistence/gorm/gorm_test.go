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

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest/orchestratortest"
	"clouditor.io/clouditor/v2/persistence"
	"github.com/google/uuid"

	"google.golang.org/protobuf/types/known/timestamppb"
)

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

			gorm := assert.Is[*storage](t, s)
			assert.NotNil(t, gorm)
			assert.Equal(t, tt.wantDialectorType, fmt.Sprintf("%T", gorm.dialector))

			// Test to create a new target of evaluation and get it again with
			// respective 'Create' and 'Get' Create record via DB call
			targetInput := &orchestrator.TargetOfEvaluation{
				Id:        uuid.New().String(),
				Name:      "SomeName",
				CreatedAt: timestamppb.Now(),
				UpdatedAt: timestamppb.Now(),
			}
			assert.NoError(t, api.Validate(targetInput))
			err = s.Create(targetInput)
			assert.NoError(t, err)

			// Get record via DB call
			targetOutput := &orchestrator.TargetOfEvaluation{}
			err = s.Get(&targetOutput, "name = ?", "SomeName")
			assert.NoError(t, err)
			assert.Equal(t, targetInput, targetOutput)
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
		Id:            testdata.MockMetricID1,
		Category:      testdata.MockMetricCategory1,
		Description:   testdata.MockMetricDescription1,
		Version:       testdata.MockMetricVersion1,
		Comments:      testdata.MockMetricComments1,
		Configuration: testdata.MockMetricConfigurations,
	}
	// Check if metric has all necessary fields
	assert.NoError(t, api.Validate(metric))

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
		err    error
		s      persistence.Storage
		target *orchestrator.TargetOfEvaluation
	)

	target = orchestratortest.NewTargetOfEvaluation()
	assert.NoError(t, api.Validate(target))

	// Create storage
	s, err = NewStorage()
	assert.NoError(t, err)

	// Return error since no record in the DB yet
	err = s.Get(&orchestrator.TargetOfEvaluation{})
	assert.ErrorIs(t, err, persistence.ErrRecordNotFound)

	// Create target of evaluation
	err = s.Create(target)
	assert.NoError(t, err)

	// Get target of evaluation via passing entire record
	gotTarget := &orchestrator.TargetOfEvaluation{}
	err = s.Get(&gotTarget)
	assert.NoError(t, err)
	assert.Equal(t, target, gotTarget)

	// Get target of evaluation via name
	gotTarget2 := &orchestrator.TargetOfEvaluation{}
	err = s.Get(&gotTarget2, "name = ?", target.Name)
	assert.NoError(t, err)
	assert.Equal(t, target, gotTarget2)

	// Get target of evaluation via description
	gotTarget3 := &orchestrator.TargetOfEvaluation{}
	err = s.Get(&gotTarget3, "description = ?", target.Description)
	assert.NoError(t, err)
	assert.NoError(t, api.Validate(gotTarget3))
	assert.Equal(t, target, gotTarget3)

	var metric = &assessment.Metric{
		Id:            testdata.MockMetricID1,
		Category:      testdata.MockMetricCategory1,
		Description:   testdata.MockMetricDescription1,
		Version:       testdata.MockMetricVersion1,
		Comments:      testdata.MockMetricComments1,
		Configuration: testdata.MockMetricConfigurations,
	}
	// Check if metric has all necessary fields
	assert.NoError(t, api.Validate(metric))

	// Create metric
	err = s.Create(metric)
	assert.NoError(t, err)

	// Get metric via Id
	gotMetric := &assessment.Metric{}
	err = s.Get(gotMetric, "id = ?", testdata.MockMetricID1)
	assert.NoError(t, err)
	assert.NoError(t, api.Validate(gotMetric))
	assert.Equal(t, metric, gotMetric)

	var impl = &assessment.MetricImplementation{
		MetricId:  testdata.MockMetricID1,
		Code:      "TestCode",
		UpdatedAt: timestamppb.New(time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)),
	}
	// Check if impl has all necessary fields
	assert.NoError(t, api.Validate(impl))

	// Create metric implementation
	err = s.Create(impl)
	assert.NoError(t, err)

	// Get metric implementation via Id
	gotImpl := &assessment.MetricImplementation{}
	err = s.Get(gotImpl, "metric_id = ?", testdata.MockMetricID1)
	assert.NoError(t, err)
	assert.NoError(t, api.Validate(gotImpl))
	assert.Equal(t, impl, gotImpl)
}

func Test_storage_List(t *testing.T) {
	var (
		err     error
		s       persistence.Storage
		target1 *orchestrator.TargetOfEvaluation
		target2 *orchestrator.TargetOfEvaluation
		targets []*orchestrator.TargetOfEvaluation
	)

	// Create storage
	s, err = NewStorage()
	assert.NoError(t, err)

	// Test target of evaluation
	target1 = &orchestrator.TargetOfEvaluation{Id: testdata.MockTargetOfEvaluationID1, Name: testdata.MockTargetOfEvaluationName1}
	target2 = &orchestrator.TargetOfEvaluation{Id: testdata.MockTargetOfEvaluationID2, Name: testdata.MockTargetOfEvaluationName2}

	// List should return empty list since no target of evaluations are in DB yet
	err = s.List(&targets, "", true, 0, -1)
	assert.ErrorIs(t, err, nil)
	assert.Empty(t, targets)

	// List should return list of 2 target of evaluations (target1 and target2)
	err = s.Create(target1)
	assert.NoError(t, err)
	err = s.Create(target2)
	assert.NoError(t, err)
	err = s.List(&targets, "", true, 0, -1)
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, 2, len(targets))
	// We only check one target of evaluation and assume the others are also correct
	assert.NoError(t, api.Validate(targets[0]))

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
	assert.NoError(t, api.Validate(certificates[0]))

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
	assert.Fail(t, "condition failed")
}

func Test_storage_Count(t *testing.T) {
	var (
		count   int64
		err     error
		s       persistence.Storage
		target1 *orchestrator.TargetOfEvaluation
		target2 *orchestrator.TargetOfEvaluation
	)

	target1 = orchestratortest.NewTargetOfEvaluation()
	target2 = orchestratortest.NewAnotherTargetOfEvaluation()

	// Create storage
	s, err = NewStorage()
	assert.NoError(t, err)

	// Since no records in DB yet, count of target of evaluations should be 0
	count, err = s.Count(&orchestrator.TargetOfEvaluation{})
	assert.NoError(t, err)
	assert.Equal(t, int(count), 0)

	// Create one target of evaluation -> count of target of evaluations should be 1
	assert.ErrorIs(t, s.Create(target1), nil)
	count, err = s.Count(&orchestrator.TargetOfEvaluation{})
	assert.NoError(t, err)
	assert.Equal(t, 1, int(count))

	// Add another one -> count of target of evaluations should be 2
	assert.ErrorIs(t, s.Create(target2), nil)
	count, err = s.Count(&orchestrator.TargetOfEvaluation{})
	assert.NoError(t, err)
	assert.Equal(t, 2, int(count))

	// Count of target of evaluations with ID "SomeName2" should be 1
	count, err = s.Count(&orchestrator.TargetOfEvaluation{}, "name = ?", testdata.MockTargetOfEvaluationName2)
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
		err       error
		s         persistence.Storage
		target    *orchestrator.TargetOfEvaluation
		newTarget *orchestrator.TargetOfEvaluation
		gotTarget *orchestrator.TargetOfEvaluation
		myVar     MyTest
	)
	target = orchestratortest.NewTargetOfEvaluation()

	// Create storage
	s, err = NewStorage(WithAdditionalAutoMigration(&MyTest{}))
	assert.NoError(t, err)

	// Create target of evaluation
	err = s.Create(target)
	assert.NoError(t, err)

	err = s.Get(&orchestrator.TargetOfEvaluation{}, "name = ?", target.Name)
	assert.NoError(t, err)

	// Save new target of evaluation: Change description. Name and ID remain unchanged
	newTarget = orchestratortest.NewTargetOfEvaluation()
	newTarget.Description = ""

	err = s.Save(newTarget, "name = ?", target.Name)
	assert.NoError(t, err)

	gotTarget = &orchestrator.TargetOfEvaluation{}
	err = s.Get(gotTarget, "name = ?", target.Name)
	assert.NoError(t, err)
	assert.NoError(t, api.Validate(gotTarget))

	// Name should be the same
	assert.Equal(t, target.Name, gotTarget.Name)
	// Description should be zero
	assert.Equal(t, "", gotTarget.Description)

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

	// Testing target of evaluation
	// Create target of evaluation
	TargetOfEvaluation := &orchestrator.TargetOfEvaluation{
		Id:          testdata.MockTargetOfEvaluationID1,
		Name:        testdata.MockTargetOfEvaluationName1,
		Description: testdata.MockTargetOfEvaluationDescription1,
		ConfiguredMetrics: []*assessment.Metric{
			{
				Id:            testdata.MockMetricID1,
				Description:   testdata.MockMetricDescription1,
				Version:       testdata.MockMetricVersion1,
				Comments:      testdata.MockMetricComments1,
				Configuration: testdata.MockMetricConfigurations,
			},
		},
	}
	// Check if target of evaluation has all necessary fields
	assert.NoError(t, api.Validate(TargetOfEvaluation))
	err = s.Create(&TargetOfEvaluation)
	assert.NoError(t, err)

	err = s.Get(&orchestrator.TargetOfEvaluation{}, "Id = ?", TargetOfEvaluation.Id)
	assert.NoError(t, err)

	err = s.Update(&orchestrator.TargetOfEvaluation{Name: "SomeNewName", Description: ""}, "Id = ?", TargetOfEvaluation.Id)
	assert.NoError(t, err)

	gotTargetOfEvaluation := &orchestrator.TargetOfEvaluation{}
	err = s.Get(&gotTargetOfEvaluation, "Id = ?", TargetOfEvaluation.Id)
	assert.NoError(t, err)
	assert.NoError(t, api.Validate(gotTargetOfEvaluation))

	// Name should be changed
	assert.Equal(t, "SomeNewName", gotTargetOfEvaluation.Name)
	// Other properties should stay the same
	assert.Equal(t, TargetOfEvaluation.Id, gotTargetOfEvaluation.Id)
	assert.Equal(t, TargetOfEvaluation.Description, gotTargetOfEvaluation.Description)
	assert.Equal(t, len(TargetOfEvaluation.ConfiguredMetrics), len(gotTargetOfEvaluation.ConfiguredMetrics))
}

func Test_storage_Delete(t *testing.T) {
	var (
		err    error
		s      persistence.Storage
		target *orchestrator.TargetOfEvaluation
	)
	target = orchestratortest.NewTargetOfEvaluation()

	// Create storage
	s, err = NewStorage()
	assert.NoError(t, err)

	// Create target of evaluation
	err = s.Create(target)
	assert.NoError(t, err)

	// Should return ErrRecordNotFound since there is no target of evaluation "Fake" in DB
	assert.ErrorIs(t, s.Delete(&orchestrator.TargetOfEvaluation{}, "name = ?", "Fake"), persistence.ErrRecordNotFound)

	// Successful deletion
	assert.Nil(t, s.Delete(&orchestrator.TargetOfEvaluation{}, "name = ?", target.Name))
	// Check with s.Get that target of evaluation is not in DB anymore
	assert.ErrorIs(t, s.Get(&orchestrator.TargetOfEvaluation{}, "name = ?", target.Name), persistence.ErrRecordNotFound)

	// Should return DB error since a non-supported type is passed (just a string instead of, e.g., &orchestrator.TargetOfEvaluation{})
	assert.Contains(t, s.Delete("Unsupported Type").Error(), "unsupported data type")
}
