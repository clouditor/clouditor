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
		wantRes *orchestrator.Certificate
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "Validation error - Empty request",
			fields:  fields{svc: NewService()},
			req:     nil,
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, api.ErrEmptyRequest.Error())
			},
		},
		{
			name:    "Validation error - Certificate Id missing in request",
			fields:  fields{svc: NewService()},
			req:     &orchestrator.GetCertificateRequest{CertificateId: ""},
			wantRes: nil,
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
			req:     &orchestrator.GetCertificateRequest{CertificateId: "WrongCertificateID"},
			wantRes: nil,
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
			req:     &orchestrator.GetCertificateRequest{CertificateId: "WrongCertificateID"},
			wantRes: nil,
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
			req:     &orchestrator.GetCertificateRequest{CertificateId: testdata.MockCertificateID},
			wantRes: nil,
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
			wantRes: orchestratortest.NewCertificate(),
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
			if tt.wantRes != nil {
				assert.NotEmpty(t, res.Id)
				assert.True(t, proto.Equal(tt.wantRes, res), "Want: %v\nGot : %v", tt.wantRes, res)
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
			name: "Validation Error - empty request",
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
			name: "Internal error - db error",
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
	type fields struct {
		svc *Service
	}
	type args struct {
		ctx context.Context
		req *orchestrator.UpdateCertificateRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *orchestrator.Certificate
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Validation Error - Request is nil",
			fields: fields{
				svc: NewService(),
			},
			args: args{
				ctx: nil,
				req: nil,
			},
			wantRes: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, api.ErrEmptyRequest.Error())
			},
		},
		{
			name: "Validation Error - Certificate is nil",
			fields: fields{
				svc: NewService(),
			},
			args: args{
				ctx: nil,
				req: &orchestrator.UpdateCertificateRequest{Certificate: nil},
			},
			wantRes: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "invalid request")
			},
		},
		{
			name: "Validation Error - Certificate ID is empty",
			fields: fields{
				svc: NewService(),
			},
			args: args{
				ctx: nil,
				req: &orchestrator.UpdateCertificateRequest{Certificate: &orchestrator.Certificate{Id: ""}},
			},
			wantRes: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "invalid request")
			},
		},
		{
			name: "Permission Denied Error - not authorized",
			fields: fields{
				svc: NewService(WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(
					false, testdata.MockAnotherCloudServiceID))),
			},
			args: args{
				ctx: nil,
				req: &orchestrator.UpdateCertificateRequest{Certificate: orchestratortest.NewCertificate()},
			},
			wantRes: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.PermissionDenied, status.Code(err))
				return assert.ErrorContains(t, err, service.ErrPermissionDenied.Error())
			},
		},
		{
			name: "Internal - db error (count)",
			fields: fields{
				svc: NewService(WithStorage(&testutil.StorageWithError{CountErr: gorm.ErrInvalidDB})),
			},
			args: args{
				ctx: nil,
				req: &orchestrator.UpdateCertificateRequest{Certificate: orchestratortest.NewCertificate()},
			},
			wantRes: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.Internal, status.Code(err))
				return assert.ErrorContains(t, err, gorm.ErrInvalidDB.Error())
			},
		},
		{
			name: "Not Found Error - certificate doesn't exist",
			fields: fields{
				svc: NewService(WithStorage(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					// Create default certificate
					assert.NoError(t, s.Create(orchestratortest.NewCertificate()))
				}))),
			},
			args: args{
				ctx: nil,
				// Try to update certificate 2 which is not in DB
				req: &orchestrator.UpdateCertificateRequest{Certificate: orchestratortest.NewCertificate2()},
			},
			wantRes: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.NotFound, status.Code(err))
				return assert.ErrorContains(t, err, certificationNotFoundErrorMessage)
			},
		},
		{
			name: "Internal - db error (save)",
			fields: fields{
				svc: NewService(WithStorage(&testutil.StorageWithError{
					// Fake Count response so we can reach the saving part
					CountRes: int64(1),
					CountErr: nil,
					SaveErr:  gorm.ErrInvalidDB})),
			},
			args: args{
				ctx: nil,
				req: &orchestrator.UpdateCertificateRequest{Certificate: orchestratortest.NewCertificate()},
			},
			wantRes: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.Internal, status.Code(err))
				return assert.ErrorContains(t, err, gorm.ErrInvalidDB.Error())
			},
		},
		{
			name: "Happy path",
			fields: fields{
				svc: NewService(WithStorage(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.NewCertificate()))
				}))),
			},
			args: args{
				ctx: nil,
				req: &orchestrator.UpdateCertificateRequest{
					Certificate: orchestratortest.NewCertificate(
						// Modify description (use description of  mockCertification2)
						orchestratortest.WithDescription(testdata.MockCertificateDescription2))},
			},
			wantRes: orchestratortest.NewCertificate(
				orchestratortest.WithDescription(testdata.MockCertificateDescription2)),
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := tt.fields.svc.UpdateCertificate(context.TODO(), tt.args.req)
			// Run ErrorAssertionFunc
			tt.wantErr(t, err)
			// Assert response
			if tt.wantRes != nil {
				assert.NotEmpty(t, res.Id)
				assert.True(t, proto.Equal(tt.wantRes, res), "Want: %v\nGot : %v", tt.wantRes, res)
			}
		})
	}
}

func Test_RemoveCertificate(t *testing.T) {
	type fields struct {
		svc *Service
	}
	type args struct {
		ctx context.Context
		req *orchestrator.RemoveCertificateRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes assert.ValueAssertionFunc
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Validation Error - Request is nil",
			fields: fields{
				svc: NewService(),
			},
			args: args{
				ctx: nil,
				req: nil,
			},
			wantRes: assert.Nil,
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, api.ErrEmptyRequest.Error())
			},
		},
		{
			name: "Validation Error - certificate id is empty",
			fields: fields{
				svc: NewService(),
			},
			args: args{
				ctx: nil,
				req: &orchestrator.RemoveCertificateRequest{CertificateId: ""},
			},
			wantRes: assert.Nil,
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, api.ErrInvalidRequest.Error())
			},
		},
		{
			name: "Permission Denied Error - No admin flag and not allowed for any cloud service",
			fields: fields{
				svc: NewService(WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(
					false))),
			},
			args: args{
				ctx: nil,
				req: &orchestrator.RemoveCertificateRequest{CertificateId: testdata.MockCertificateID},
			},
			wantRes: assert.Nil,
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				assert.Equal(t, codes.PermissionDenied, status.Code(err))
				return assert.ErrorContains(t, err, service.ErrPermissionDenied.Error())
			},
		},
		{
			name: "Not Found Error - certificate is hidden due to missing authorization)",
			fields: fields{
				svc: NewService(
					WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(
						false, testdata.MockAnotherCloudServiceID)),
					WithStorage(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						assert.NoError(t, s.Create(orchestratortest.NewCertificate()))
					}))),
			},
			args: args{
				ctx: nil,
				req: &orchestrator.RemoveCertificateRequest{CertificateId: testdata.MockCertificateID},
			},
			wantRes: assert.Nil,
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				assert.Equal(t, codes.NotFound, status.Code(err))
				return assert.ErrorContains(t, err, certificationNotFoundErrorMessage)
			},
		},
		{
			name: "Internal - db error",
			fields: fields{
				svc: NewService(
					// Just to make it clear. Nilling it would also result in this strategy since it is the default
					WithAuthorizationStrategy(&service.AuthorizationStrategyAllowAll{}),
					WithStorage(&testutil.StorageWithError{DeleteErr: gorm.ErrInvalidDB})),
			},
			args: args{
				ctx: nil,
				req: &orchestrator.RemoveCertificateRequest{CertificateId: testdata.MockCertificateID},
			},
			wantRes: assert.Nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.Internal, status.Code(err))
				return assert.ErrorContains(t, err, gorm.ErrInvalidDB.Error())
			},
		},
		{
			name: "Happy path - with authorization allAllowed",
			fields: fields{
				svc: NewService(
					// Just to make it clear. Nilling it would also result in this strategy since it is the default
					WithAuthorizationStrategy(&service.AuthorizationStrategyAllowAll{}),
					WithStorage(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						assert.NoError(t, s.Create(orchestratortest.NewCertificate()))
						assert.NoError(t, s.Create(orchestratortest.NewCertificate2()))
					}))),
			},
			args: args{
				ctx: nil,
				req: &orchestrator.RemoveCertificateRequest{CertificateId: testdata.MockCertificateID},
			},
			wantRes: func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				assert.NotNil(t, i)
				svc, ok := i2[0].(*Service)
				assert.True(t, ok)
				// Verify that certificate 2 is still in the DB (by counting the number of occurrences = 1)
				n, err := svc.storage.Count(&orchestrator.Certificate{}, testdata.MockCertificateID2)
				assert.NoError(t, err)
				assert.Equal(t, int64(1), n)
				// Verify that the default certificate isn't in the DB anymore (occurrences = 0)
				n, err = svc.storage.Count(&orchestrator.Certificate{}, testdata.MockCertificateID)
				assert.NoError(t, err)
				return assert.Equal(t, int64(0), n)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path - with authorization for one certain cloud service",
			fields: fields{
				svc: NewService(
					WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(
						false, testdata.MockCloudServiceID)),
					WithStorage(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						assert.NoError(t, s.Create(orchestratortest.NewCertificate()))
						assert.NoError(t, s.Create(orchestratortest.NewCertificate2()))
					}))),
			},
			args: args{
				ctx: nil,
				req: &orchestrator.RemoveCertificateRequest{CertificateId: testdata.MockCertificateID},
			},
			wantRes: func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				assert.NotNil(t, i)
				svc, ok := i2[0].(*Service)
				assert.True(t, ok)
				// Verify that certificate 2 is still in the DB (by counting the number of occurrences = 1)
				n, err := svc.storage.Count(&orchestrator.Certificate{}, testdata.MockCertificateID2)
				assert.NoError(t, err)
				assert.Equal(t, int64(1), n)
				// Verify that the default certificate isn't in the DB anymore (occurrences = 0)
				n, err = svc.storage.Count(&orchestrator.Certificate{}, testdata.MockCertificateID)
				assert.NoError(t, err)
				return assert.Equal(t, int64(0), n)
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := tt.fields.svc.RemoveCertificate(context.TODO(), tt.args.req)
			// Run ValueAssertionFunc on response
			tt.wantRes(t, res, tt.fields.svc)
			// Run ErrorAssertionFunc
			tt.wantErr(t, err)
		})
	}
}
