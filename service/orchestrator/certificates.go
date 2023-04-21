package orchestrator

import (
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/logging"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/service"
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateCertificate implements method for creating a new certificate
func (svc *Service) CreateCertificate(_ context.Context, req *orchestrator.CreateCertificateRequest) (
	*orchestrator.Certificate, error) {
	// Validate request
	err := service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	// Persist the new certificate in our database
	err = svc.storage.Create(req.Certificate)
	if err != nil {
		return nil,
			status.Errorf(codes.Internal, "could not add certificate to the database: %v", err)
	}

	logging.LogRequest(log, logrus.DebugLevel, logging.Create, req)

	// Return certificate
	return req.Certificate, nil
}

// GetCertificate implements method for getting a certificate, e.g. to show its state in the UI
func (svc *Service) GetCertificate(_ context.Context, req *orchestrator.GetCertificateRequest) (response *orchestrator.Certificate, err error) {
	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	response = new(orchestrator.Certificate)
	err = svc.storage.Get(response, "Id = ?", req.CertificateId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "certificate not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}
	return response, nil
}

// ListCertificates implements method for getting a certificate, e.g. to show its state in the UI
func (svc *Service) ListCertificates(ctx context.Context, req *orchestrator.ListCertificatesRequest) (res *orchestrator.ListCertificatesResponse, err error) {
	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	res = new(orchestrator.ListCertificatesResponse)

	res.Certificates, res.NextPageToken, err = service.PaginateStorage[*orchestrator.Certificate](req, svc.storage,
		service.DefaultPaginationOpts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
	}

	// Only list certificates user is authorized to see
	var certificateSlice []*orchestrator.Certificate
	for i := range res.Certificates {
		if svc.authz.CheckAccess(ctx, service.AccessRead, res.Certificates[i]) {
			certificateSlice = append(certificateSlice, res.Certificates[i])
		}
	}

	res.Certificates = make([]*orchestrator.Certificate, len(certificateSlice))
	copy(res.Certificates, certificateSlice)

	return
}

// UpdateCertificate implements method for updating an existing certificate
func (svc *Service) UpdateCertificate(_ context.Context, req *orchestrator.UpdateCertificateRequest) (response *orchestrator.Certificate, err error) {
	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	count, err := svc.storage.Count(req.Certificate, req.Certificate.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	if count == 0 {
		return nil, status.Error(codes.NotFound, "certificate not found")
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
func (svc *Service) RemoveCertificate(_ context.Context, req *orchestrator.RemoveCertificateRequest) (response *emptypb.Empty, err error) {
	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	err = svc.storage.Delete(&orchestrator.Certificate{}, "Id = ?", req.CertificateId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "certificate not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	logging.LogRequest(log, logrus.DebugLevel, logging.Remove, req)

	return &emptypb.Empty{}, nil
}
