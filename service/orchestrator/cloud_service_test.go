package orchestrator

import (
	"clouditor.io/clouditor/persistence"
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

func TestService_ListCloudServices(t *testing.T) {
	var (
		listCloudServicesResponse *orchestrator.ListCloudServicesResponse
		cloudService              *orchestrator.CloudService
		err                       error
	)

	// Create service with DB
	s := startTestService()

	// 1st case: No services stored
	listCloudServicesResponse, err = s.ListCloudServices(context.Background(), &orchestrator.ListCloudServicesRequest{})
	assert.Nil(t, err)
	assert.NotNil(t, listCloudServicesResponse.Services)
	assert.Empty(t, listCloudServicesResponse.Services)

	// 2nd case: One service stored
	cloudService, err = s.CreateDefaultTargetCloudService()
	assert.Nil(t, err)
	assert.NotNil(t, cloudService)

	listCloudServicesResponse, err = s.ListCloudServices(context.Background(), &orchestrator.ListCloudServicesRequest{})
	assert.Nil(t, err)
	assert.NotNil(t, listCloudServicesResponse.Services)
	assert.NotEmpty(t, listCloudServicesResponse.Services)
	assert.Equal(t, len(listCloudServicesResponse.Services), 1)

}

func TestService_RemoveCloudService(t *testing.T) {
	var (
		cloudServiceResponse      *orchestrator.CloudService
		err                       error
		listCloudServicesResponse *orchestrator.ListCloudServicesResponse
	)

	// Create service with DB
	s := startTestService()

	// 1st case: Empty service ID error
	_, err = s.RemoveCloudService(context.Background(), &orchestrator.RemoveCloudServiceRequest{ServiceId: ""})
	assert.NotNil(t, err)
	assert.Equal(t, status.Code(err), codes.InvalidArgument)

	// 2nd case: ErrRecordNotFound
	_, err = s.RemoveCloudService(context.Background(), &orchestrator.RemoveCloudServiceRequest{ServiceId: DefaultTargetCloudServiceId})
	assert.NotNil(t, err)
	assert.Equal(t, status.Code(err), codes.NotFound)

	// 3rd case: Record removed successfully
	cloudServiceResponse, err = s.CreateDefaultTargetCloudService()
	assert.Nil(t, err)
	assert.NotNil(t, cloudServiceResponse)

	// There is a record for cloud services in the DB (default one)
	listCloudServicesResponse, err = s.ListCloudServices(context.Background(), &orchestrator.ListCloudServicesRequest{})
	assert.Nil(t, err)
	assert.NotNil(t, listCloudServicesResponse.Services)
	assert.NotEmpty(t, listCloudServicesResponse.Services)

	// Remove record
	_, err = s.RemoveCloudService(context.Background(), &orchestrator.RemoveCloudServiceRequest{ServiceId: DefaultTargetCloudServiceId})
	assert.Nil(t, err)

	// There is a record for cloud services in the DB (default one)
	listCloudServicesResponse, err = s.ListCloudServices(context.Background(), &orchestrator.ListCloudServicesRequest{})
	assert.Nil(t, err)
	assert.NotNil(t, listCloudServicesResponse.Services)
	assert.Empty(t, listCloudServicesResponse.Services)
}

func TestService_CreateDefaultTargetCloudService(t *testing.T) {
	var (
		cloudServiceResponse *orchestrator.CloudService
		err                  error
	)

	// Create service with DB
	s := startTestService()

	// 1st case: No records for cloud services -> Default target service is created
	cloudServiceResponse, err = s.CreateDefaultTargetCloudService()
	assert.Nil(t, err)
	assert.Equal(t, &orchestrator.CloudService{
		Id:          DefaultTargetCloudServiceId,
		Name:        DefaultTargetCloudServiceName,
		Description: DefaultTargetCloudServiceDescription,
	}, cloudServiceResponse)

	// 2nd case: There is already a record for service (the default target service) -> Nothing added and no error
	cloudServiceResponse, err = s.CreateDefaultTargetCloudService()
	assert.Nil(t, err)
	assert.Nil(t, cloudServiceResponse)
}

// startTestService starts service with DB initialization
func startTestService() *Service {
	gormX := new(persistence.GormX)
	err := gormX.Init(true, "", 0)
	if err != nil {
		panic(err)
	}
	return NewService(gormX)
}
