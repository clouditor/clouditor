package clouditor.metrics.document_csaf_no_redirect

import input as document
import rego.v1

default applicable := false

# this requirement is optional and we cannot model this correctly yet. We also do not have
# optional metrics, so we just return true for now
default compliant := true

applicable if {
	# check resource type
	"SecurityAdvisoryDocument" in document.type
}
