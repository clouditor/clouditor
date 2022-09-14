// Copyright 2016-2022 Fraunhofer AISEC
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

package orchestrator

import (
	"context"
	"fmt"
	"testing"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/internal/testutil/orchestratortest"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/persistence/gorm"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var AssuranceLevelHigh = "high"
var AssuranceLevelMedium = "medium"

func Test_CreateTargetOfEvaluation(t *testing.T) {
	type fields struct {
		cloudServiceHooks     []orchestrator.CloudServiceHookFunc
		results               map[string]*assessment.AssessmentResult
		AssessmentResultHooks []func(result *assessment.AssessmentResult, err error)
		storage               persistence.Storage
		metricsFile           string
		loadMetricsFunc       func() ([]*assessment.Metric, error)
		catalogsFile          string
		events                chan *orchestrator.MetricChangeEvent
	}
	type args struct {
		ctx context.Context
		req *orchestrator.CreateTargetOfEvaluationRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.ValueAssertionFunc
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Create(&orchestrator.CloudService{Id: "MyService"})
					assert.NoError(t, err)

					err = s.Create(orchestratortest.NewCatalog())
					assert.NoError(t, err)
				}),
			},
			args: args{req: &orchestrator.CreateTargetOfEvaluationRequest{
				Toe: orchestratortest.NewTargetOfEvaluation(),
			}},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				svc := i2[0].(*Service)

				// We want to assert that certain things happened in our database
				var toes []*orchestrator.TargetOfEvaluation
				err := svc.storage.List(&toes, "", false, 0, -1, gorm.WithoutPreload())
				if !assert.NoError(t, err) {
					return false
				}
				if !assert.Equal(t, 1, len(toes)) {
					return false
				}

				var service orchestrator.CloudService
				err = svc.storage.Get(&service, "id = ?", "MyService")
				if !assert.NoError(t, err) {
					return false
				}

				return assert.Equal(t, 1, len(service.CatalogsInScope))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				cloudServiceHooks:     tt.fields.cloudServiceHooks,
				results:               tt.fields.results,
				AssessmentResultHooks: tt.fields.AssessmentResultHooks,
				storage:               tt.fields.storage,
				metricsFile:           tt.fields.metricsFile,
				loadMetricsFunc:       tt.fields.loadMetricsFunc,
				catalogsFile:          tt.fields.catalogsFile,
				events:                tt.fields.events,
			}

			gotRes, err := svc.CreateTargetOfEvaluation(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.CreateTargetOfEvaluation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil {
				tt.want(t, gotRes, svc)
			}
		})
	}
}

func Test_GetTargetOfEvaluation(t *testing.T) {
	type fields struct {
		storage persistence.Storage
	}
	type args struct {
		req *orchestrator.GetTargetOfEvaluationRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantResponse assert.ValueAssertionFunc
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "invalid request",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Create(&orchestrator.CloudService{Id: "MyService"})
					assert.NoError(t, err)

					err = s.Create(orchestratortest.NewCatalog())
					assert.NoError(t, err)

					err = s.Create(orchestratortest.NewTargetOfEvaluation())
					assert.NoError(t, err)
				}),
			},
			args:         args{req: nil},
			wantResponse: assert.Nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, api.ErrRequestIsNil.Error())
			},
		},
		{
			name: "toe not found",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Create(&orchestrator.CloudService{Id: "MyService"})
					assert.NoError(t, err)

					err = s.Create(orchestratortest.NewCatalog())
					assert.NoError(t, err)

					err = s.Create(orchestratortest.NewTargetOfEvaluation())
					assert.NoError(t, err)
				}),
			},
			args: args{req: &orchestrator.GetTargetOfEvaluationRequest{
				CloudServiceId: "",
			}},
			wantResponse: assert.Nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.NotFound, status.Code(err))
				return assert.ErrorContains(t, err, "toe ID is empty")
			},
		},
		{
			name: "valid toe",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Create(&orchestrator.CloudService{Id: "MyService"})
					assert.NoError(t, err)

					err = s.Create(orchestratortest.NewCatalog())
					assert.NoError(t, err)

					err = s.Create(orchestratortest.NewTargetOfEvaluation())
					assert.NoError(t, err)
				}),
			},
			args: args{req: &orchestrator.GetTargetOfEvaluationRequest{
				CloudServiceId: "MyService",
				CatalogId:      "Cat1234",
			}},
			wantResponse: func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				res, ok := i.(*orchestrator.TargetOfEvaluation)
				want := orchestratortest.NewTargetOfEvaluation()
				assert.True(t, ok)
				fmt.Println(res)
				assert.Equal(t, want.CloudServiceId, res.CloudServiceId)
				return assert.Equal(t, want.CatalogId, res.CatalogId)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orchestratorService := Service{
				storage: tt.fields.storage,
			}
			res, err := orchestratorService.GetTargetOfEvaluation(context.Background(), tt.args.req)

			// Validate the error via the ErrorAssertionFunc function
			tt.wantErr(t, err)

			// Validate the response via the ValueAssertionFunc function
			tt.wantResponse(t, res)
		})
	}
}

func Test_ListTargetsOfEvaluation(t *testing.T) {
	var (
		listTargetsOfEvaluationResponse *orchestrator.ListTargetsOfEvaluationResponse
		err                             error
	)

	orchestratorService := NewService()
	err = orchestratorService.storage.Create(&orchestrator.CloudService{Id: "MyService"})
	assert.NoError(t, err)
	err = orchestratorService.storage.Create(orchestratortest.NewCatalog())
	assert.NoError(t, err)

	// 1st case: No ToEs stored
	listTargetsOfEvaluationResponse, err = orchestratorService.ListTargetsOfEvaluation(context.Background(), &orchestrator.ListTargetsOfEvaluationRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listTargetsOfEvaluationResponse.Toes)
	assert.Empty(t, listTargetsOfEvaluationResponse.Toes)

	// 2nd case: One ToE stored
	err = orchestratorService.storage.Create(orchestratortest.NewTargetOfEvaluation())
	assert.NoError(t, err)

	listTargetsOfEvaluationResponse, err = orchestratorService.ListTargetsOfEvaluation(context.Background(), &orchestrator.ListTargetsOfEvaluationRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listTargetsOfEvaluationResponse.Toes)
	assert.NotEmpty(t, listTargetsOfEvaluationResponse.Toes)
	assert.Equal(t, 1, len(listTargetsOfEvaluationResponse.Toes))

	// 3rd case: Invalid request
	_, err = orchestratorService.ListTargetsOfEvaluation(context.Background(),
		&orchestrator.ListTargetsOfEvaluationRequest{OrderBy: "not a field"})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Contains(t, err.Error(), api.ErrInvalidColumnName.Error())
}

func Test_UpdateTargetOfEvaluation(t *testing.T) {
	var (
		toe *orchestrator.TargetOfEvaluation
		err error
	)
	orchestratorService := NewService()

	// 1st case: ToE is nil
	_, err = orchestratorService.UpdateTargetOfEvaluation(context.Background(), &orchestrator.UpdateTargetOfEvaluationRequest{})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// 2nd case: Ids are empty
	_, err = orchestratorService.UpdateTargetOfEvaluation(context.Background(), &orchestrator.UpdateTargetOfEvaluationRequest{
		CloudServiceId: "",
		CatalogId:      "",
	})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// 3rd case: ToE not found since there are no ToEs
	_, err = orchestratorService.UpdateTargetOfEvaluation(context.Background(), &orchestrator.UpdateTargetOfEvaluationRequest{
		Toe: &orchestrator.TargetOfEvaluation{
			CloudServiceId: "MyService",
			CatalogId:      "Cat1234",
		},
	})
	assert.Equal(t, codes.NotFound, status.Code(err))

	// 4th case: ToE updated successfully
	err = orchestratorService.storage.Create(&orchestrator.TargetOfEvaluation{
		CloudServiceId: orchestratortest.MockServiceID,
		CatalogId:      orchestratortest.MockCatalogID,
		AssuranceLevel: &AssuranceLevelHigh,
	})
	assert.NoError(t, err)

	// update the toe's assurance level and send the update request
	toe, err = orchestratorService.UpdateTargetOfEvaluation(context.Background(), &orchestrator.UpdateTargetOfEvaluationRequest{
		Toe: &orchestrator.TargetOfEvaluation{
			CloudServiceId: orchestratortest.MockServiceID,
			CatalogId:      orchestratortest.MockCatalogID,
			AssuranceLevel: &AssuranceLevelMedium,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, toe)
	assert.Equal(t, "new assurance level", toe.AssuranceLevel)
}

func Test_RemoveTargetOfEvaluation(t *testing.T) {
	var (
		err                             error
		listTargetsOfEvaluationResponse *orchestrator.ListTargetsOfEvaluationResponse
	)
	orchestratorService := NewService()

	// 1st case: Empty ToE ID error
	_, err = orchestratorService.RemoveTargetOfEvaluation(context.Background(), &orchestrator.RemoveTargetOfEvaluationRequest{
		CloudServiceId: "MyService",
		CatalogId:      "Cat1234",
	})
	assert.Error(t, err)
	assert.Equal(t, status.Code(err), codes.InvalidArgument)

	// 2nd case: ErrRecordNotFound
	_, err = orchestratorService.RemoveTargetOfEvaluation(context.Background(), &orchestrator.RemoveTargetOfEvaluationRequest{
		CloudServiceId: "0000",
		CatalogId:      "0000",
	})
	assert.Error(t, err)
	assert.Equal(t, status.Code(err), codes.NotFound)

	// 3rd case: Record removed successfully
	err = orchestratorService.storage.Create(&orchestrator.TargetOfEvaluation{
		CloudServiceId: orchestratortest.MockServiceID,
		CatalogId:      orchestratortest.MockCatalogID,
		AssuranceLevel: &AssuranceLevelHigh,
	})
	assert.NoError(t, err)

	// There is a record for ToE in the DB
	listTargetsOfEvaluationResponse, err = orchestratorService.ListTargetsOfEvaluation(context.Background(), &orchestrator.ListTargetsOfEvaluationRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listTargetsOfEvaluationResponse.Toes)

	// Remove record
	_, err = orchestratorService.RemoveTargetOfEvaluation(context.Background(), &orchestrator.RemoveTargetOfEvaluationRequest{
		CloudServiceId: "MyService",
		CatalogId:      "Cat1234",
	})
	assert.NoError(t, err)

	// There is no record for ToE in the DB
	listTargetsOfEvaluationResponse, err = orchestratorService.ListTargetsOfEvaluation(context.Background(), &orchestrator.ListTargetsOfEvaluationRequest{})
	assert.NoError(t, err)
	assert.Nil(t, listTargetsOfEvaluationResponse.Toes)
}
