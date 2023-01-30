package testdata

import (
	"clouditor.io/clouditor/voc"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	// Auth
	MockAuthUser     = "clouditor"
	MockAuthPassword = "clouditor"

	MockAuthClientID     = "client"
	MockAuthClientSecret = "secret"
	MockCustomClaims     = "cloudserviceid"

	// Cloud Service
	MockCloudServiceID                 = "11111111-1111-1111-1111-111111111111"
	MockCloudServiceName               = "Mock Cloud Service"
	MockCloudServiceDescription        = "This is a mock cloud service"
	MockAnotherCloudServiceID          = "22222222-2222-2222-2222-222222222222"
	MockAnotherCloudServiceName        = "Another Mock Cloud Service"
	MockAnotherCloudServiceDescription = "This is another mock cloud service"

	// Catalog
	MockCatalogID                 = "Cat1234"
	MockCatalogName               = "Mock Catalog"
	MockCatalogDescription        = "This is a mock catalog"
	MockCategoryName              = "Mock Category Name"
	MockCategoryDescription       = "This is a mock category"
	MockControlID                 = "Cont1234"
	MockControlName               = "Mock Control Name"
	MockControlDescription        = "This is a mock control"
	MockSubControlID              = "Cont1234.1"
	MockSubControlName            = "Mock Sub-Control Name"
	MockSubControlDescription     = "This is a mock sub-control"
	MockAnotherControlID          = "Cont4567"
	MockAnotherControlName        = "Mock Another Control Name"
	MockAnotherControlDescription = "This is a another mock control"

	// Metric
	MockMetricID          = "Mock Metric"
	MockMetricName        = "Mock Metric Name"
	MockMetricDescription = "This is a mock metric"
	MockAnotherMetricID   = "Another Mock Metric"

	// Assessment result
	MockAssessmentResultID                   = "11111111-1111-1111-1111-111111111111"
	MockAssessmentResultNonComplianceComment = "non_compliance_comment"

	// Evidence
	MockEvidenceID            = "11111111-1111-1111-1111-111111111111"
	MockEvidenceToolID        = "Mock Tool ID"
	MockAnotherEvidenceID     = "22222222-2222-2222-2222-222222222222"
	MockAnotherEvidenceToolID = "Another Mock Tool ID"

	// Resource
	MockResourceID          = "my-resource-id"
	MockResourceName        = "my-resource-name"
	MockResourceStorageID   = voc.ResourceID("some-storage-service-id")
	MockResourceStorageName = "some-storage-service-name"

	// Certificate
	MockCertificateID          = "1234"
	MockCertificateName        = "EUCS"
	MockCertificateDescription = "This is a mock certificate"
	MockCertificateCab         = "Cab123"
	MockCertificateStandard    = "EUCS"

	// State
	MockStateId     = "12345"
	MockStateState  = "new"
	MockStateTreeID = "12345"

	// Assessment Results
	MockAssessmentResult1ID = "11111111-1111-1111-1111-111111111111"
	MockAssessmentResult2ID = "22222222-2222-2222-2222-222222222222"
	MockAssessmentResult3ID = "33333333-3333-3333-3333-333333333333"
	MockAssessmentResult4ID = "44444444-4444-4444-4444-444444444444"
)

var (
	// Catalog
	AssuranceLevelBasic       = "basic"
	AssuranceLevelSubstantial = "substantial"
	AssuranceLevelHigh        = "high"

	// Metric Configuration
	MockMetricConfigurationTargetValueString = &structpb.Value{
		Kind: &structpb.Value_StringValue{
			StringValue: "MockTargetValue",
		},
	}
)
