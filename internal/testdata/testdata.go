package testdata

import (
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
	MockMetricDescription1 = "This is a mock metric"
	MockMetricCategory1    = "Mock Category 1"
	MockMetricVersion1     = "1.0"
	MockMetricComments1    = "Mock metric comments 1"
	MockMetricID2          = "Mock Metric 2"
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

	// Virtual Machine
	MockVirtualMachineID1          = "my-vm-id"
	MockVirtualMachineName1        = "my-vm-name"
	MockVirtualMachineDescription1 = "This is a mock virtual machine"
	MockVirtualMachineID2          = "my-other-vm-id"
	MockVirtualMachineName2        = "my-other-vm-name"
	MockVirtualMachineDescription2 = "This is another mock virtual machine"

	// Block Storage
	MockBlockStorageID1          = "my-block-storage-id"
	MockBlockStorageName1        = "my-block-storage-name"
	MockBlockStorageDescription1 = "This is a mock block storage"
	MockBlockStorageID2          = "my-other-block-storage-id"
	MockBlockStorageName2        = "my-other-block-storage-name"
	MockBlockStorageDescription2 = "This is another mock block storage"

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

	// Resource Types
	MockVirtualMachineTypes = []string{"VirtualMachine", "Compute", "Infrastructure", "Resource"}
	MockBlockStorageTypes   = []string{"BlockStorage", "Storage", "Infrastructure", "Resource"}
)
