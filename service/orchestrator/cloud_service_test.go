package orchestrator

import (
	"context"
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
			status.Error(codes.InvalidArgument, "Service is empty"),
		},
		{
			"missing service name",
			&orchestrator.RegisterCloudServiceRequest{Service: &orchestrator.CloudService{}},
			nil,
			status.Error(codes.InvalidArgument, "Service name is empty"),
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
			status.Error(codes.InvalidArgument, "Service id is empty"),
		},
		{
			"missing service id",
			&orchestrator.GetCloudServiceRequest{},
			nil,
			status.Error(codes.InvalidArgument, "Service id is empty"),
		},
		{
			"invalid service id",
			&orchestrator.GetCloudServiceRequest{ServiceId: "does-not-exist"},
			nil,
			status.Error(codes.NotFound, "Service not found"),
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
