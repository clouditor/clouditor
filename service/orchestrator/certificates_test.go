package orchestrator

import (
	"context"
	"reflect"
	"testing"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/internal/testutil/servicetest"
	"clouditor.io/clouditor/internal/testutil/servicetest/orchestratortest"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/service"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

func Test_CreateCertificate(t *testing.T) {
	// Instantiate Mock certificate (so creating time is same for assertion)
	mockCertificate := orchestratortest.NewCertificate()
	type fields struct {
		svc *Service
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
			fields: fields{svc: NewService()},
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
			fields: fields{svc: NewService()},
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
			fields: fields{svc: NewService()},
			args: args{
				context.Background(),
				&orchestrator.CreateCertificateRequest{
					// Use certificate without an ID
					Certificate: &orchestrator.Certificate{Id: ""},
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
				svc: NewService(WithAuthorizationStrategy(
					servicetest.NewAuthorizationStrategy(false, testdata.MockAnotherCloudServiceID))),
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
			name: "internal error - certificate already exists",
			fields: fields{
				svc: NewService(
					WithStorage(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						// Create mockCertificate in DB beforehand
						assert.NoError(t, s.Create(mockCertificate))
					})))},
			args: args{
				context.Background(),
				&orchestrator.CreateCertificateRequest{
					Certificate: orchestratortest.NewCertificate(),
				},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.Internal, status.Code(err))
				return assert.ErrorContains(t, err, persistence.ErrUniqueConstraintFailed.Error())
			},
		},
		{
			name: "happy path - valid certificate",
			fields: fields{
				svc: NewService(WithAuthorizationStrategy(
					// Only allow certificates belonging to MockCloudServiceID
					servicetest.NewAuthorizationStrategy(false, testdata.MockCloudServiceID))),
			},
			args: args{
				context.Background(),
				&orchestrator.CreateCertificateRequest{
					// mockCertificate's corresponding cloud service ID is MockCloudServiceID (authorization succeeds)
					Certificate: mockCertificate,
				},
			},
			wantRes: mockCertificate,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.fields.svc
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
func Test_GetCertificate(t *testing.T) {
	type fields struct {
		svc *Service
	}
	tests := []struct {
		name    string
		fields  fields
		req     *orchestrator.GetCertificateRequest
		res     *orchestrator.Certificate
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "Validation error - Empty request",
			fields: fields{svc: NewService()},
			req:    nil,
			res:    nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, api.ErrEmptyRequest.Error())
			},
		},
		{
			name:   "Validation error - Certificate Id missing in request",
			fields: fields{svc: NewService()},
			req:    &orchestrator.GetCertificateRequest{CertificateId: ""},
			res:    nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "invalid request: invalid GetCertificateRequest.CertificateId: value length must be at least 1 runes")
			},
		},
		{
			name: "Not Found Error - Certificate doesn't exist",
			fields: fields{
				svc: NewService(WithStorage(
					testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						// Create Certificate
						assert.NoError(t, s.Create(orchestratortest.NewCertificate()))
					}))),
			},
			req: &orchestrator.GetCertificateRequest{CertificateId: "WrongCertificateID"},
			res: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.NotFound, status.Code(err))
				return assert.ErrorContains(t, err, "certificate not found")
			},
		},
		{
			name: "Internal error - DB Error",
			fields: fields{
				svc: NewService(WithStorage(&testutil.StorageWithError{GetErr: gorm.ErrInvalidDB})),
			},
			req: &orchestrator.GetCertificateRequest{CertificateId: "WrongCertificateID"},
			res: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.Internal, status.Code(err))
				return assert.ErrorContains(t, err, gorm.ErrInvalidDB.Error())
			},
		},
		{
			name: "Permission denied error - not authorized",
			fields: fields{
				svc: NewService(
					WithStorage(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						assert.NoError(t, s.Create(orchestratortest.NewCertificate()))
					})),
					// Only authorized for MockCloudServiceID
					WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(
						false, testdata.MockAnotherCloudServiceID)),
				),
			},
			// Only authorized for MockAnotherCloudServiceID (=2222-2...) and not MockCloudServiceID (=1111-1...)
			req: &orchestrator.GetCertificateRequest{CertificateId: testdata.MockCertificateID},
			res: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.PermissionDenied, status.Code(err))
				return assert.ErrorContains(t, err, service.ErrPermissionDenied.Error())
			},
		},
		{
			name: "Happy path",
			fields: fields{
				svc: NewService(WithStorage(
					testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						// Create Certificate
						assert.NoError(t, s.Create(orchestratortest.NewCertificate()))
					}))),
			},
			req:     &orchestrator.GetCertificateRequest{CertificateId: testdata.MockCertificateID},
			res:     orchestratortest.NewCertificate(),
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := tt.fields.svc.GetCertificate(context.Background(), tt.req)

			// Run validation on response
			assert.NoError(t, res.Validate())
			// Run ErrorAssertionFunc
			tt.wantErr(t, err)
			// Assert response
			if tt.res != nil {
				assert.NotEmpty(t, res.Id)
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
					authz: servicetest.NewAuthorizationStrategy(true, ""),
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
					authz: servicetest.NewAuthorizationStrategy(true, ""),
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
						assert.NoError(t, s.Create(orchestratortest.NewCertificate2()))
					}),
					authz: servicetest.NewAuthorizationStrategy(false, testdata.MockCloudServiceID),
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
			// Run ValueAssertionFunc on response
			tt.wantRes(t, res)
			// Run ErrorAssertionFunc
			tt.wantErr(t, err)
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
