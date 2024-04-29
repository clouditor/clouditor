package clouditor.metrics.document_csaf_red_restricted

import input as document
import rego.v1

default applicable := false

default compliant := false

restricted := ["RED", "AMBER"]

applicable if {
	# check resource type
	"SecurityAdvisoryDocument" in document.type

	# check, if document is restricted (i.e. RED/AMBER) labeled
	document.labels.tlp in restricted
}

compliant if {
	# RED/EMBER must NOT be freely accessible
	auth := document.documentLocation.remoteDocumentLocation.authenticity
	not auth.noAuthentication
}
