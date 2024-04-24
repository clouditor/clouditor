package clouditor.metrics.document_content_valid

import rego.v1
import input as document

default applicable := false

default compliant := false

applicable if {
	# the resource type should be a VM
	"Document" in document.type
}

compliant if {
	# Check if errors are available
    # If no errors exist, the document is valid
    count(document.schemaValidation.errors) == 0
}