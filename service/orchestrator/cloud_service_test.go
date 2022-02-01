package orchestrator

import (
	"clouditor.io/clouditor/persistence"
	"context"
	"fmt"
	"testing"

	"clouditor.io/clouditor/api/orchestrator"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func TestRegisterCloudService(t *testing.T) {
	tests := []struct {
		name string
		req  *orchestrator.RegisterCloudServiceRequest
		res  *orchestrator.CloudService
		err  error
	}{
		{
			"missing service",
			&orchestrator.RegisterCloudServiceRequest{},
			nil,
			status.Error(codes.InvalidArgument, "service is empty"),
		},
		{
			"missing service name",
			&orchestrator.RegisterCloudServiceRequest{Service: &orchestrator.CloudService{}},
			nil,
			status.Error(codes.InvalidArgument, "service name is empty"),
		},
		{
			"valid",
			&orchestrator.RegisterCloudServiceRequest{Service: &orchestrator.CloudService{Name: "test", Description: "some"}},
			&orchestrator.CloudService{Name: "test", Description: "some"},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := service.RegisterCloudService(context.Background(), tt.req)

			if tt.err == nil {
				assert.Equal(t, err, tt.err)
			} else {
				assert.EqualError(t, err, tt.err.Error())
			}

			if tt.res != nil {
				assert.NotEmpty(t, res.Id)
			}

			// reset the IDs because we cannot compare them, since they are randomly generated
			if res != nil {
				res.Id = ""
			}

			if tt.res != nil {
				tt.res.Id = ""
			}

			assert.True(t, proto.Equal(res, tt.res), "%v != %v", res, tt.res)
		})
	}
}

func TestGetCloudService(t *testing.T) {
	tests := []struct {
		name string
		req  *orchestrator.GetCloudServiceRequest
		res  *orchestrator.CloudService
		err  error
	}{
		{
			"missing request",
			nil,
			nil,
			status.Error(codes.InvalidArgument, "service id is empty"),
		},
		{
			"missing service id",
			&orchestrator.GetCloudServiceRequest{},
			nil,
			status.Error(codes.InvalidArgument, "service id is empty"),
		},
		{
			"invalid service id",
			&orchestrator.GetCloudServiceRequest{ServiceId: "does-not-exist"},
			nil,
			status.Error(codes.NotFound, "service not found"),
		},
		{
			"valid",
			&orchestrator.GetCloudServiceRequest{ServiceId: defaultTarget.Id},
			&orchestrator.CloudService{
				Id:          DefaultTargetCloudServiceId,
				Name:        DefaultTargetCloudServiceName,
				Description: DefaultTargetCloudServiceDescription,
			},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := service.GetCloudService(context.Background(), tt.req)

			if tt.err == nil {
				assert.Equal(t, err, tt.err)
			} else {
				assert.EqualError(t, err, tt.err.Error())
			}

			if tt.res != nil {
				assert.NotEmpty(t, res.Id)
			}

			assert.True(t, proto.Equal(res, tt.res), "%v != %v", res, tt.res)
		})
	}
}

// TODO(lebogg): Do it without table tests (strict ordering)
func TestService_CreateDefaultTargetCloudService(t *testing.T) {
	type fields struct {
	}
	tests := []struct {
		name        string
		fields      fields
		wantService *orchestrator.CloudService
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "Tesst",
			wantService: &orchestrator.CloudService{
				Id:          DefaultTargetCloudServiceId,
				Name:        DefaultTargetCloudServiceName,
				Description: DefaultTargetCloudServiceDescription,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return false
			},
		},
		{
			name:        "AlreadyCreated",
			wantService: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return false
			},
		},
	}
	s := startTestService()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotService, err := s.CreateDefaultTargetCloudService()
			if tt.wantErr(t, err, fmt.Sprintf("CreateDefaultTargetCloudService()")) {
				return
			}
			assert.Equalf(t, tt.wantService.String(), gotService.String(), "CreateDefaultTargetCloudService()")
		})
	}
}

func startTestService() *Service {
	gormX := new(persistence.GormX)
	err := gormX.Init(true, "", 0)
	if err != nil {
		panic(err)
	}
	return NewService(gormX)
}
