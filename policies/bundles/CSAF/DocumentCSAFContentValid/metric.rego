package clouditor.metrics.document_csaf_content_valid

import input as document
import rego.v1

default applicable := false

default compliant := false

applicable if {
	# check resource type
	"SecurityAdvisoryDocument" in document.type
}

compliant if {
	# Check if errors are available
	# If no errors exist, the document is valid
	count(document.schemaValidation.errors) == 0
}
