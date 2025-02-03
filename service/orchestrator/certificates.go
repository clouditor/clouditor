package orchestrator

import (
	"context"
	"errors"
	"slices"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/logging"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/service"

	"github.com/sirupsen/logrus"
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
	if err = api.Validate(req); err != nil {
		return
	}

	// Check if client is allowed to access the corresponding certification target (targeted in the certificate)
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

// GetCertificate implements method for getting a certificate, e.g. to show its state in the UI.
func (svc *Service) GetCertificate(ctx context.Context, req *orchestrator.GetCertificateRequest) (
	res *orchestrator.Certificate, err error) {

	// Validate request
	if err = api.Validate(req); err != nil {
		return
	}

	res = new(orchestrator.Certificate)
	err = svc.storage.Get(res, "Id = ?", req.CertificateId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, ErrCertificationNotFound
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "%v: %v", persistence.ErrDatabase, err)
	}

	// Check if client is allowed to access the corresponding certification target (targeted in the certificate)
	all, allowed := svc.authz.AllowedCertificationTargets(ctx)
	if !all && !slices.Contains(allowed, res.CertificationTargetId) {
		// Important to nil the response since it is set already
		return nil, status.Error(codes.PermissionDenied, service.ErrPermissionDenied.Error())
	}

	return
}

// ListCertificates implements method for getting all certificates, e.g. to show its state in the UI. The response does not indicate whether there are no certificates available or the access is denied.
func (svc *Service) ListCertificates(ctx context.Context, req *orchestrator.ListCertificatesRequest) (
	res *orchestrator.ListCertificatesResponse, err error) {

	// Validate request
	if err = api.Validate(req); err != nil {
		return nil, err
	}

	// We only list certificates the user is authorized to see (w.r.t. the certification target)
	var (
		query []string
		args  []any
	)

	all, allowed := svc.authz.AllowedCertificationTargets(ctx)
	if !all {
		query = append(query, "certification_target_id IN ?")
		args = append(args, allowed)
	}

	res = new(orchestrator.ListCertificatesResponse)

	res.Certificates, res.NextPageToken, err = service.PaginateStorage[*orchestrator.Certificate](req, svc.storage,
		service.DefaultPaginationOpts, persistence.BuildConds(query, args)...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
	}

	return
}

// ListPublicCertificates implements method for getting all certificates without the state history, e.g. to show its state in the UI
func (svc *Service) ListPublicCertificates(_ context.Context, req *orchestrator.ListPublicCertificatesRequest) (res *orchestrator.ListPublicCertificatesResponse, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	res = new(orchestrator.ListPublicCertificatesResponse)

	res.Certificates, res.NextPageToken, err = service.PaginateStorage[*orchestrator.Certificate](req, svc.storage,
		service.DefaultPaginationOpts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
	}

	// Delete state history from certificates
	for i := range res.Certificates {
		res.Certificates[i].States = nil
	}

	return
}

// UpdateCertificate implements method for updating an existing certificate
func (svc *Service) UpdateCertificate(ctx context.Context, req *orchestrator.UpdateCertificateRequest) (response *orchestrator.Certificate, err error) {
	// Validate request
	if err = api.Validate(req); err != nil {
		return nil, err
	}

	// Check authorization
	if !svc.authz.CheckAccess(ctx, service.AccessUpdate, req) {
		err = service.ErrPermissionDenied
		return
	}

	count, err := svc.storage.Count(req.Certificate, "id=?", req.Certificate.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v: %v", persistence.ErrDatabase, err)
	}

	if count == 0 {
		return nil, ErrCertificationNotFound
	}

	response = req.Certificate

	err = svc.storage.Save(response, "Id = ?", response.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v: %v", persistence.ErrDatabase, err)
	}

	logging.LogRequest(log, logrus.DebugLevel, logging.Update, req)

	return
}

// RemoveCertificate implements method for removing a certificate. The response does not indicate whether there are no certificates available or the access is denied.
func (svc *Service) RemoveCertificate(ctx context.Context, req *orchestrator.RemoveCertificateRequest) (response *emptypb.Empty, err error) {
	// Validate request
	if err = api.Validate(req); err != nil {
		return nil, err
	}

	// Lookup if certificate entry is in DB. If not, return NotFound error
	if err = svc.checkCertificateExistence(req); err != nil {
		return
	}
	// 2) Check if client is authorized to remove certificate.
	// Only remove certificate if user is authorized for the corresponding certification target.
	if err = svc.checkCertificateAuthorization(ctx, req); err != nil {
		return
	}

	// Delete entry since client is authorized to do so
	err = svc.storage.Delete(&orchestrator.Certificate{}, "Id = ?", req.CertificateId)
	if err != nil { // Only internal errors left since others (Permission and NotFound) are already covered
		return nil, status.Errorf(codes.Internal, "%v: %v", persistence.ErrDatabase, err)
	}

	logging.LogRequest(log, logrus.DebugLevel, logging.Remove, req)

	return &emptypb.Empty{}, nil
}

// checkCertificateAuthorization checks if client is authorized to remove certificate by
// 1) checking admin flag: If it is enabled (`all`) the client is authorized
// 2) querying the DB within the range of certification targets (`allowed`) the client is allowed to access
// Error is returned if not authorized or internal DB error occurred.
// Note: Use the checkExistence before to ensure that the entry is in the DB!
func (svc *Service) checkCertificateAuthorization(ctx context.Context, req *orchestrator.RemoveCertificateRequest) error {
	all, allowed := svc.authz.AllowedCertificationTargets(ctx)
	if !all {
		count2, err := svc.storage.Count(&orchestrator.Certificate{}, "id = ? AND certification_target_id IN ?",
			req.GetCertificateId(), allowed)
		if err != nil {
			return status.Errorf(codes.Internal, "%v: %v", persistence.ErrDatabase, err)
		}
		if count2 == 0 {
			return service.ErrPermissionDenied
		}
	}
	return nil
}

// checkExistence checks if the entry is in the DB. An error is returned if not, or if there is an internal DB error.
func (svc *Service) checkCertificateExistence(req *orchestrator.RemoveCertificateRequest) error {
	count, err := svc.storage.Count(&orchestrator.Certificate{}, "Id = ?", req.CertificateId)
	if err != nil {
		return status.Errorf(codes.Internal, "%v: %v", persistence.ErrDatabase, err)
	}
	if count == 0 {
		return ErrCertificationNotFound
	}
	return nil
}
