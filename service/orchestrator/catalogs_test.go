package orchestrator

import (
	"context"
	"testing"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/internal/testutil/orchestratortest"
	"clouditor.io/clouditor/persistence"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestService_GetCategory(t *testing.T) {
	type fields struct {
		metricConfigurations  map[string]map[string]*assessment.MetricConfiguration
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
		req *orchestrator.GetCategoryRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *orchestrator.Category
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					// Create Catalog
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
				})},
			args: args{
				req: &orchestrator.GetCategoryRequest{CatalogId: "Cat1234", CategoryName: "My name"},
			},
			wantRes: &orchestrator.Category{
				Name:        "My name",
				Description: "test",
				CatalogId:   "Cat1234",
				Controls: []*orchestrator.Control{{
					ShortName:         "Cont1234",
					Name:              "Mock Control",
					Description:       "This is a mock control",
					CategoryName:      "My name",
					CategoryCatalogId: "Cat1234",
					// at this level, we will not have the metrics
					Metrics: []*assessment.Metric{},
					// at this level, we will not have the sub-controls
					Controls: []*orchestrator.Control{},
				}},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Service{
				metricConfigurations:  tt.fields.metricConfigurations,
				cloudServiceHooks:     tt.fields.cloudServiceHooks,
				results:               tt.fields.results,
				AssessmentResultHooks: tt.fields.AssessmentResultHooks,
				storage:               tt.fields.storage,
				metricsFile:           tt.fields.metricsFile,
				loadMetricsFunc:       tt.fields.loadMetricsFunc,
				catalogsFile:          tt.fields.catalogsFile,
				events:                tt.fields.events,
			}
			gotRes, err := srv.GetCategory(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !proto.Equal(gotRes, tt.wantRes) {
				t.Errorf("Service.GetCategory() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_GetControl(t *testing.T) {
	type fields struct {
		UnimplementedOrchestratorServer orchestrator.UnimplementedOrchestratorServer
		metricConfigurations            map[string]map[string]*assessment.MetricConfiguration
		cloudServiceHooks               []orchestrator.CloudServiceHookFunc
		results                         map[string]*assessment.AssessmentResult
		AssessmentResultHooks           []func(result *assessment.AssessmentResult, err error)
		storage                         persistence.Storage
		metricsFile                     string
		loadMetricsFunc                 func() ([]*assessment.Metric, error)
		catalogsFile                    string
	}
	type args struct {
		ctx context.Context
		req *orchestrator.GetControlRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *orchestrator.Control
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					// Create Catalog
					assert.NoError(t, s.Create(orchestratortest.NewControl()))
				})},
			args: args{req: &orchestrator.GetControlRequest{CatalogId: "Cat1234", CategoryName: "My name", ControlShortName: "Cont1234"}},
			wantRes: &orchestrator.Control{
				ShortName:         "Cont1234",
				CategoryName:      "My name",
				CategoryCatalogId: "Cat1234",
				Name:              "Mock Control",
				Description:       "This is a mock control",
				Metrics: []*assessment.Metric{{
					Id:          "MockMetric",
					Name:        "A Mock Metric",
					Description: "This Metric is a mock metric",
					Scale:       assessment.Metric_ORDINAL,
					Range: &assessment.Range{
						Range: &assessment.Range_AllowedValues{AllowedValues: &assessment.AllowedValues{
							Values: []*structpb.Value{
								structpb.NewBoolValue(false),
								structpb.NewBoolValue(true),
							}}}},
				}},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Service{
				UnimplementedOrchestratorServer: tt.fields.UnimplementedOrchestratorServer,
				metricConfigurations:            tt.fields.metricConfigurations,
				cloudServiceHooks:               tt.fields.cloudServiceHooks,
				results:                         tt.fields.results,
				AssessmentResultHooks:           tt.fields.AssessmentResultHooks,
				storage:                         tt.fields.storage,
				metricsFile:                     tt.fields.metricsFile,
				loadMetricsFunc:                 tt.fields.loadMetricsFunc,
				catalogsFile:                    tt.fields.catalogsFile,
			}
			gotRes, err := srv.GetControl(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetControl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !proto.Equal(gotRes, tt.wantRes) {
				t.Errorf("Service.GetControl() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
