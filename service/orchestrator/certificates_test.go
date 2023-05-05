package orchestrator

import (
	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/internal/testutil/servicetest"
	"clouditor.io/clouditor/internal/testutil/servicetest/orchestratortest"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/service"
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"testing"
)

func Test_CreateCertificate(t *testing.T) {
	// Mock certificates
	mockCertificate := orchestratortest.NewCertificate()
	mockCertificateWithoutID := orchestratortest.NewCertificate()
	mockCertificateWithoutID.Id = ""
	type fields struct {
		service *Service
	}
	type args struct {
		in0 context.Context
		req *orchestrator.CreateCertificateRequest
	}
	var tests = []struct {
		name    string
		fields  fields
		args    args
		wantRes *orchestrator.Certificate
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "validation error - missing request",
			fields: fields{service: NewService()},
			args: args{
				context.Background(),
				nil,
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, api.ErrEmptyRequest.Error())
			},
		},
		{
			name:   "validation error - missing certificate",
			fields: fields{service: NewService()},
			args: args{
				context.Background(),
				&orchestrator.CreateCertificateRequest{},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "Certificate: value is required")
			},
		},
		{
			name:   "validation error - missing certificate id",
			fields: fields{service: NewService()},
			args: args{
				context.Background(),
				&orchestrator.CreateCertificateRequest{
					Certificate: mockCertificateWithoutID,
				},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "Id: value length must be at least 1 runes")
			},
		},
		{
			name: "authorization error - permission denied",
			fields: fields{
				service: &Service{
					authz: servicetest.NewAuthorizationStrategy(false, []string{testdata.MockAnotherCloudServiceID}),
				},
			},
			args: args{
				context.Background(),
				&orchestrator.CreateCertificateRequest{
					Certificate: mockCertificate,
				},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.PermissionDenied, status.Code(err))
				return assert.ErrorContains(t, err, service.ErrPermissionDenied.Error())
			},
		},
		{
			name:   "internal error - db error",
			fields: fields{service: NewService(WithStorage(&testutil.StorageWithError{CreateErr: gorm.ErrInvalidDB}))},
			args: args{
				context.Background(),
				&orchestrator.CreateCertificateRequest{
					Certificate: mockCertificate,
				},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.Internal, status.Code(err))
				return assert.ErrorContains(t, err, gorm.ErrInvalidDB.Error())
			},
		},
		{
			name:   "happy path - valid certificate",
			fields: fields{service: NewService()},
			args: args{
				context.Background(),
				&orchestrator.CreateCertificateRequest{
					Certificate: mockCertificate,
				},
			},
			wantRes: mockCertificate,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.fields.service
			gotResponse, err := svc.CreateCertificate(tt.args.in0, tt.args.req)
			assert.NoError(t, gotResponse.Validate())

			tt.wantErr(t, err)

			// If no error is wanted, check response
			if !reflect.DeepEqual(gotResponse, tt.wantRes) {
				t.Errorf("Service.CreateCertificate() = %v, want %v", gotResponse, tt.wantRes)
			}
		})
	}
}

func Test_UpdateCertificate(t *testing.T) {
	var (
		certificate *orchestrator.Certificate
		err         error
	)
	orchestratorService := NewService()

	// 1st case: Certificate is nil
	_, err = orchestratorService.UpdateCertificate(context.Background(), &orchestrator.UpdateCertificateRequest{})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// 2nd case: Certificate ID is nil
	_, err = orchestratorService.UpdateCertificate(context.Background(), &orchestrator.UpdateCertificateRequest{
		Certificate: certificate,
	})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// 3rd case: Certificate not found since there are no certificates yet
	_, err = orchestratorService.UpdateCertificate(context.Background(), &orchestrator.UpdateCertificateRequest{
		Certificate: &orchestrator.Certificate{
			Id:             testdata.MockCertificateID,
			Name:           "EUCS",
			CloudServiceId: testdata.MockCloudServiceID,
		},
	})
	assert.Equal(t, codes.NotFound, status.Code(err))

	// 4th case: Certificate updated successfully
	mockCertificate := orchestratortest.NewCertificate()
	err = orchestratorService.storage.Create(mockCertificate)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	// update the certificate's description and send the update request
	mockCertificate.Description = "new description"
	certificate, err = orchestratorService.UpdateCertificate(context.Background(), &orchestrator.UpdateCertificateRequest{
		Certificate: mockCertificate,
	})
	assert.NoError(t, certificate.Validate())
	assert.NoError(t, err)
	assert.NotNil(t, certificate)
	assert.Equal(t, "new description", certificate.Description)
}

func Test_RemoveCertificate(t *testing.T) {
	var (
		err                      error
		listCertificatesResponse *orchestrator.ListCertificatesResponse
	)
	orchestratorService := NewService()

	// 1st case: Empty certificate ID error
	_, err = orchestratorService.RemoveCertificate(context.Background(), &orchestrator.RemoveCertificateRequest{CertificateId: ""})
	assert.Error(t, err)
	assert.Equal(t, status.Code(err), codes.InvalidArgument)

	// 2nd case: ErrRecordNotFound
	_, err = orchestratorService.RemoveCertificate(context.Background(), &orchestrator.RemoveCertificateRequest{CertificateId: "0000"})
	assert.Error(t, err)
	assert.Equal(t, status.Code(err), codes.NotFound)

	// 3rd case: Record removed successfully
	mockCertificate := orchestratortest.NewCertificate()
	assert.NoError(t, mockCertificate.Validate())
	err = orchestratorService.storage.Create(mockCertificate)
	assert.NoError(t, err)

	// There is a record for certificates in the DB (default one)
	listCertificatesResponse, err = orchestratorService.ListCertificates(context.Background(), &orchestrator.ListCertificatesRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listCertificatesResponse.Certificates)
	assert.NotEmpty(t, listCertificatesResponse.Certificates)

	// Remove record
	_, err = orchestratorService.RemoveCertificate(context.Background(), &orchestrator.RemoveCertificateRequest{CertificateId: mockCertificate.Id})
	assert.NoError(t, err)

	// There is a record for cloud services in the DB (default one)
	listCertificatesResponse, err = orchestratorService.ListCertificates(context.Background(), &orchestrator.ListCertificatesRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listCertificatesResponse.Certificates)
	assert.Empty(t, listCertificatesResponse.Certificates)
}

func Test_GetCertificate(t *testing.T) {
	type fields struct {
		storage persistence.Storage
	}
	tests := []struct {
		name    string
		fields  fields
		req     *orchestrator.GetCertificateRequest
		res     *orchestrator.Certificate
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Empty request",
			req:  nil,
			res:  nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, api.ErrEmptyRequest.Error())
			},
		},
		{
			name: "Certificate Id missing in request",
			req:  &orchestrator.GetCertificateRequest{CertificateId: ""},
			res:  nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "invalid request: invalid GetCertificateRequest.CertificateId: value length must be at least 1 runes")
			},
		},
		{
			name: "Certificate not found",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					// Create Certificate
					assert.NoError(t, s.Create(orchestratortest.NewCertificate()))
				}),
			},
			req: &orchestrator.GetCertificateRequest{CertificateId: "WrongCertificateID"},
			res: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.NotFound, status.Code(err))
				return assert.ErrorContains(t, err, "certificate not found")
			},
		},
		{
			name: "valid",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					// Create Certificate
					assert.NoError(t, s.Create(orchestratortest.NewCertificate()))
				}),
			},
			req:     &orchestrator.GetCertificateRequest{CertificateId: testdata.MockCertificateID},
			res:     orchestratortest.NewCertificate(),
			wantErr: assert.NoError,
		},
	}
	orchestratorService := NewService()

	// Create Certificate
	if err := orchestratorService.storage.Create(orchestratortest.NewCertificate()); err != nil {
		panic(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := orchestratorService.GetCertificate(context.Background(), tt.req)
			assert.NoError(t, res.Validate())

			tt.wantErr(t, err)

			if tt.res != nil {
				assert.NotEmpty(t, res.Id)
				// Compare timestamp. We have to cut off the microseconds and seconds, otherwise an error can be returned.
				t1 := strings.Split(tt.res.States[0].GetTimestamp(), ".")[0]
				tt.res.States[0].Timestamp = t1[:len(t1)-3]
				t2 := strings.Split(res.States[0].GetTimestamp(), ".")[0]
				res.States[0].Timestamp = t2[:len(t2)-3]
				assert.True(t, proto.Equal(tt.res, res), "Want: %v\nGot : %v", tt.res, res)
			}
		})
	}
}

func Test_ListCertificates(t *testing.T) {
	type fields struct {
		svc *Service
	}
	type args struct {
		ctx context.Context
		req *orchestrator.ListCertificatesRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes assert.ValueAssertionFunc
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Validation Error",
			fields: fields{
				svc: NewService(),
			},
			args: args{
				ctx: nil,
				req: nil,
			},
			wantRes: func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				return assert.Nil(t, i)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, api.ErrEmptyRequest.Error())
			},
		},
		{
			name: "Internal error",
			fields: fields{
				svc: &Service{
					storage: &testutil.StorageWithError{
						ListErr: gorm.ErrInvalidDB,
					},
					authz: servicetest.NewAuthorizationStrategy(true, nil),
				},
			},
			args: args{
				ctx: nil,
				req: &orchestrator.ListCertificatesRequest{},
			},
			wantRes: assert.Nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, gorm.ErrInvalidDB.Error())
			},
		},
		{
			name: "Happy path - all cloud services are allowed",
			fields: fields{
				svc: &Service{
					storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						assert.NoError(t, s.Create(orchestratortest.NewCertificate()))
					}),
					authz: servicetest.NewAuthorizationStrategy(true, nil),
				},
			},
			args: args{
				ctx: nil,
				req: &orchestrator.ListCertificatesRequest{},
			},
			wantRes: func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				res, ok := i.(*orchestrator.ListCertificatesResponse)
				assert.True(t, ok)
				assert.Len(t, res.Certificates, 1)
				return assert.Equal(t, res.Certificates[0].Id, testdata.MockCertificateID)

			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path - one cloud services is allowed",
			fields: fields{
				svc: &Service{
					storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						assert.NoError(t, s.Create(orchestratortest.NewCertificate()))
						assert.NoError(t, s.Create(
							orchestratortest.NewCertificate(
								orchestratortest.WithMockCertificateID("4321"),
								orchestratortest.WithMockServiceID(testdata.MockAnotherCloudServiceID))))
					}),
					authz: servicetest.NewAuthorizationStrategy(false, []string{testdata.MockCloudServiceID}),
				},
			},
			args: args{
				ctx: nil,
				req: &orchestrator.ListCertificatesRequest{},
			},
			wantRes: func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				res, ok := i.(*orchestrator.ListCertificatesResponse)
				assert.True(t, ok)
				assert.Len(t, res.Certificates, 1)
				return assert.Equal(t, res.Certificates[0].Id, testdata.MockCertificateID)

			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := tt.fields.svc.ListCertificates(context.TODO(), tt.args.req)
			tt.wantRes(t, res)
			tt.wantErr(t, err)
		})
	}

	//var (
	//	listCertificatesResponse *orchestrator.ListCertificatesResponse
	//	err                      error
	//)

	//orchestratorService := NewService()
	//orchestratorService.authz = servicetest.NewAuthorizationStrategy(false, []string{testdata.MockCloudServiceID})
	//// 1st case: No services stored
	//listCertificatesResponse, err = orchestratorService.ListCertificates(context.Background(), &orchestrator.ListCertificatesRequest{})
	//assert.NoError(t, err)
	//assert.NotNil(t, listCertificatesResponse.Certificates)
	//assert.Empty(t, listCertificatesResponse.Certificates)
	//
	//// 2nd case: One service stored
	//err = orchestratorService.storage.Create(orchestratortest.NewCertificate())
	//assert.NoError(t, err)
	//
	//listCertificatesResponse, err = orchestratorService.ListCertificates(context.Background(), &orchestrator.ListCertificatesRequest{})
	//// We check only the first certificate and assume that all certificates are valid
	//assert.NoError(t, listCertificatesResponse.Certificates[0].Validate())
	//assert.NoError(t, err)
	//assert.NotNil(t, listCertificatesResponse.Certificates)
	//assert.NotEmpty(t, listCertificatesResponse.Certificates)
	//assert.Equal(t, len(listCertificatesResponse.Certificates), 1)
	//assert.Equal(t, listCertificatesResponse.Certificates[0].CloudServiceId, testdata.MockCloudServiceID)
	//
	//// 3rd case: User is not allowed for certificates belonging to cloud service
	//orchestratorService.authz = servicetest.NewAuthorizationStrategy(false, []string{testdata.MockAnotherCloudServiceID})
	//listCertificatesResponse, err = orchestratorService.ListCertificates(context.Background(), &orchestrator.ListCertificatesRequest{})
	////assert.NoError(t, listCertificatesResponse.Certificates[0].Validate())
	//assert.NoError(t, err)
	//assert.Empty(t, listCertificatesResponse.Certificates)

}
