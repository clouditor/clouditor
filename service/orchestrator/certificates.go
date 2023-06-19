package orchestrator

import (
	"context"
	"errors"

	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/logging"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/service"

	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ErrCertificationNotFound indicates the certification was not found
var ErrCertificationNotFound = status.Error(codes.NotFound, "certificate not found")

// CreateCertificate implements method for creating a new certificate
func (svc *Service) CreateCertificate(ctx context.Context, req *orchestrator.CreateCertificateRequest) (
	res *orchestrator.Certificate, err error) {

	// Validate request
	if err = service.ValidateRequest(req); err != nil {
		return
	}

	// Check if client is allowed to access the corresponding cloud service (targeted in the certificate)
	if !svc.authz.CheckAccess(ctx, service.AccessCreate, req) {
		err = service.ErrPermissionDenied
		return
	}

	// Persist the new certificate in our database
	err = svc.storage.Create(req.Certificate)
	if err != nil {
		err = status.Errorf(codes.Internal, "could not add certificate to the database: %v", err)
		return
	}

	logging.LogRequest(log, logrus.DebugLevel, logging.Create, req)

	// Return certificate
	res = req.Certificate
	return
}

// GetCertificate implements method for getting a certificate, e.g. to show its state in the UI
func (svc *Service) GetCertificate(ctx context.Context, req *orchestrator.GetCertificateRequest) (
	res *orchestrator.Certificate, err error) {

	// Validate request
	if err = service.ValidateRequest(req); err != nil {
		return
	}

	res = new(orchestrator.Certificate)
	err = svc.storage.Get(res, "Id = ?", req.CertificateId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, ErrCertificationNotFound
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	// Check if client is allowed to access the corresponding cloud service (targeted in the certificate)
	areAllAllowed, allowedCloudServices := svc.authz.AllowedCloudServices(ctx)
	if !areAllAllowed && !slices.Contains(allowedCloudServices, res.CloudServiceId) {
		// Important to nil the response since it is set already
		return nil, status.Error(codes.PermissionDenied, service.ErrPermissionDenied.Error())
	}

	return
}

// ListCertificates implements method for getting a certificate, e.g. to show its state in the UI
func (svc *Service) ListCertificates(ctx context.Context, req *orchestrator.ListCertificatesRequest) (
	res *orchestrator.ListCertificatesResponse, err error) {

	// Validate request
	if err = service.ValidateRequest(req); err != nil {
		return nil, err
	}

	// We only list certificates the user is authorized to see (w.r.t. the cloud service)
	var conds []any
	areAllAllowed, cloudServices := svc.authz.AllowedCloudServices(ctx)
	if !areAllAllowed {
		conds = append([]any{"cloud_service_id IN ?"}, []any{cloudServices})
	}

	res = new(orchestrator.ListCertificatesResponse)

	res.Certificates, res.NextPageToken, err = service.PaginateStorage[*orchestrator.Certificate](req, svc.storage,
		service.DefaultPaginationOpts, conds...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
	}

	return
}

// UpdateCertificate implements method for updating an existing certificate
// Todo: Add auth
func (svc *Service) UpdateCertificate(ctx context.Context, req *orchestrator.UpdateCertificateRequest) (response *orchestrator.Certificate, err error) {
	// Validate request
	if err = service.ValidateRequest(req); err != nil {
		return nil, err
	}

	// Check authorization
	if !svc.authz.CheckAccess(ctx, service.AccessUpdate, req) {
		err = service.ErrPermissionDenied
		return
	}

	count, err := svc.storage.Count(req.Certificate, req.Certificate.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	if count == 0 {
		return nil, ErrCertificationNotFound
	}

	response = req.Certificate

	err = svc.storage.Save(response, "Id = ?", response.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	logging.LogRequest(log, logrus.DebugLevel, logging.Update, req)

	return
}

// RemoveCertificate implements method for removing a certificate
func (svc *Service) RemoveCertificate(ctx context.Context, req *orchestrator.RemoveCertificateRequest) (response *emptypb.Empty, err error) {
	// Validate request
	if err = service.ValidateRequest(req); err != nil {
		return nil, err
	}

	// Only remove certificate if user is authorized for the corresponding cloud service
	areAllAllowed, allowedCloudService := svc.authz.AllowedCloudServices(ctx)
	// 1st case:  User is authorized for all cloud services (admin)
	if areAllAllowed {
		err = svc.storage.Delete(&orchestrator.Certificate{}, "Id = ?", req.CertificateId)
	} else { // 2nd case: User is authorized for some cloud services (or none at all)
		err = svc.storage.Delete(&orchestrator.Certificate{},
			"id = ? AND cloud_service_id IN ?", req.CertificateId, allowedCloudService)
	}
	if errors.Is(err, persistence.ErrRecordNotFound) {
		// could also mean that user is not authorized for corresponding cloud service (2nd case)
		return nil, ErrCertificationNotFound
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	logging.LogRequest(log, logrus.DebugLevel, logging.Remove, req)

	return &emptypb.Empty{}, nil
}
