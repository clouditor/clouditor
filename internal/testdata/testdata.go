package testdata

import (
	"time"

	"google.golang.org/protobuf/types/known/structpb"
)

const (
	MockGRPCTarget = "localhost"

	// Azure
	MockLocationWestEurope     = "West Europe"
	MockLocationEastUs         = "eastus"
	MockSubscriptionID         = "00000000-0000-0000-0000-000000000000"
	MockSubscriptionResourceID = "/subscriptions/00000000-0000-0000-0000-000000000000"
	MockResourceGroupID        = "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"
	MockResourceGroup          = "TestResourceGroup"

	// IONOS Cloud
	MockIonosVolumeID1         = "39d85e98-c3da-11ed-afa1-0242ac120002"
	MockIonosVolumeID2         = "49d85e98-c3da-11ed-afa1-0242ac120002"
	MockIonosBlockStorageID1   = "59d85e98-c3da-11ed-afa1-0242ac120002"
	MockIonosBlockStorageName1 = "Mock IONOS Block Storage 1"
	MockIonosBlockStorageID2   = "69d85e98-c3da-11ed-afa1-0242ac120002"
	MockIonosBlockStorageName2 = "Mock IONOS Block Storage 2"
	MockIonosNicID1            = "59d85e98-c3da-11ed-afa1-0242ac120002"
	MockIonosNicName1          = "Mock IONOS NIC 1"
	MockIonosNicID2            = "69d85e98-c3da-11ed-afa1-0242ac120002"
	MockIonosNicName2          = "Mock IONOS NIC 2"
	MockIonosVMID1             = "79d85e98-c3da-11ed-afa1-0242ac120002"
	MockIonosVMName1           = "Mock IONOS VM1"
	MockIonosVMID2             = "89d85e98-c3da-11ed-afa1-0242ac120002"
	MOckIonosVMName2           = "Mock IONOS VM2"
	MockIonosDatacenterID1     = "99d85e98-c3da-11ed-afa1-0242ac120002"
	MockIonosDatacenterID2     = "a9d85e98-c3da-11ed-afa1-0242ac120002"
	MockIonosLoadBalancerID1   = "b9d85e98-c3da-11ed-afa1-0242ac120002"
	MockIonosLoadBalancerName1 = "Mock IONOS Load Balancer 1"
	MockIonosLoadBalancerID2   = "c9d85e98-c3da-11ed-afa1-0242ac120002"
	MockIonosLoadBalancerName2 = "Mock IONOS Load Balancer 2"
	MockIonosNetworkID1        = "d9d85e98-c3da-11ed-afa1-0242ac120002"
	MockIonosNetworkID2        = "e9d85e98-c3da-11ed-afa1-0242ac120002"

	// Audit Scope
	MockAuditScopeID1   = "11111111-1111-1111-1111-111111111123"
	MockAuditScopeName1 = "Mock Audit Scope 1"
	MockAuditScopeID2   = "11111111-1111-1111-1111-111111111124"
	MockAuditScopeName2 = "Mock Audit Scope 2"
	MockAuditScopeID3   = "11111111-1111-1111-1111-111111111125"
	MockAuditScopeName3 = "Mock Audit Scope 3"

	// Auth
	MockAuthUser     = "clouditor"
	MockAuthPassword = "clouditor"

	MockAuthClientID     = "client"
	MockAuthClientSecret = "secret"

	// Target of Evaluation
	MockTargetOfEvaluationID1          = "11111111-1111-1111-1111-111111111111"
	MockTargetOfEvaluationName1        = "Mock Target of Evaluation"
	MockTargetOfEvaluationDescription1 = "This is a mock target of evaluation"
	MockTargetOfEvaluationID2          = "22222222-2222-2222-2222-222222222222"
	MockTargetOfEvaluationName2        = "Another Mock Target of Evaluation"
	MockTargetOfEvaluationDescription2 = "This is another mock target of evaluation"

	// Catalog
	MockCatalogID1            = "Catalog 1"
	MockCatalogID2            = "Catalog 2"
	MockCatalogName           = "Mock Catalog"
	MockCatalogDescription    = "This is a mock catalog"
	MockCategoryName          = "Mock Category Name"
	MockCategoryDescription   = "This is a mock category"
	MockControlID1            = "Cont1"
	MockControlID2            = "Cont2"
	MockControlID3            = "Cont3"
	MockControlID4            = "Cont4"
	MockControlID5            = "Cont5"
	MockControlName           = "Mock Control Name"
	MockControlDescription    = "This is a mock control"
	MockSubControlID11        = "Cont1.1"
	MockSubControlID12        = "Cont1.2"
	MockSubControlID21        = "Cont2.1"
	MockSubControlID22        = "Cont2.2"
	MockSubControlID31        = "Cont3.1"
	MockSubControlID32        = "Cont3.2"
	MockSubControlID          = "Cont1234.1"
	MockSubControlName        = "Mock Sub-Control Name"
	MockSubControlDescription = "This is a mock sub-control"

	// Metric
	MockMetricID1          = "Mock Metric 1"
	MockMetricName1        = "Mock Metric Name"
	MockMetricDescription1 = "This is a mock metric"
	MockMetricCategory1    = "Mock Category 1"
	MockMetricID2          = "Mock Metric 2"
	MockMetricName2        = "Mock Metric Name 2"
	MockMetricDescription2 = "This is mock metric 2"
	MockMetricCategory2    = "Mock Category 2"

	// Assessment result
	MockAssessmentResultID                   = "11111111-1111-1111-1111-111111111111"
	MockAssessmentResultNonComplianceComment = "non_compliance_comment"
	MockAssessmentResultToolID               = "Another Assessment Result Tool ID"

	// Evidence
	MockEvidenceID1     = "11111111-1111-1111-1111-111111111111"
	MockEvidenceToolID1 = "39d85e98-c3da-11ed-afa1-0242ac120002"
	MockEvidenceID2     = "22222222-2222-2222-2222-222222222222"
	MockEvidenceToolID2 = "49d85e98-c3da-11ed-afa1-0242ac120002"

	// Resource
	MockResourceID1         = "my-resource-id"
	MockResourceName1       = "my-resource-name"
	MockResourceStorageID   = "some-storage-service-id"
	MockResourceStorageName = "some-storage-service-name"
	MockResourceID2         = "my-other-resource"

	// Properties for a Certificate
	MockCertificateID          = "1234"
	MockCertificateName        = "EUCS"
	MockCertificateDescription = "This is the default mock certificate"
	MockCertificateCab         = "Cab123"
	// Properties for an alternative Certificate
	MockCertificateID2          = "4321"
	MockCertificateName2        = "MDR"
	MockCertificateDescription2 = "This is another mock certificate"
	MockCertificateCab2         = "Cab321"

	// Properties for a State
	MockStateId     = "12345"
	MockStateState  = "new"
	MockStateTreeID = "12345"
	// Properties for an alternative State
	MockStateId2     = "54321"
	MockStateState2  = "suspended"
	MockStateTreeID2 = "54321"

	// Assessment Results
	MockAssessmentResult1ID = "11111111-1111-1111-1111-111111111111"
	MockAssessmentResult2ID = "22222222-2222-2222-2222-222222222222"
	MockAssessmentResult3ID = "33333333-3333-3333-3333-333333333333"
	MockAssessmentResult4ID = "44444444-4444-4444-4444-444444444444"

	// Evaluation Results
	MockEvaluationResult1ID  = "11111111-1111-1111-1111-111111111111"
	MockEvaluationResult2ID  = "22222222-2222-2222-2222-222222222222"
	MockEvaluationResult3ID  = "33333333-3333-3333-3333-333333333333"
	MockEvaluationResult4ID  = "44444444-4444-4444-4444-444444444444"
	MockEvaluationResult5ID  = "55555555-5555-5555-5555-555555555555"
	MockEvaluationResult6ID  = "66666666-6666-6666-6666-666666666666"
	MockEvaluationResult7ID  = "77777777-7777-7777-7777-777777777777"
	MockEvaluationResult8ID  = "88888888-8888-8888-8888-888888888888"
	MockEvaluationResult9ID  = "99999999-9999-9999-9999-999999999999"
	MockEvaluationResult10ID = "11111111-1111-1111-1111-111111111110"
)

var (
	CreationTime = time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)

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
