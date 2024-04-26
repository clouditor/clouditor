package clouditor.metrics.security_advisory_service_valid_metadata

import input as service
import rego.v1

default applicable := false

default compliant := false

# retrieve metadata document from related resources
metadata := service.related[service.serviceMetadataDocumentId]

applicable if {
	# check resource type
	"SecurityAdvisoryService" in service.type
}

compliant if {
	# must exist and must be a valid service metadata document
	metadata
	"ServiceMetadataDocument" in metadata.type
	count(metadata.schemaValidation.errors) == 0
}
