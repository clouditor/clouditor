package clouditor.metrics.document_csaf_white_accessible

import input as document
import rego.v1

default applicable := false

default compliant := false

applicable if {
	# check resource type
	"SecurityAdvisoryDocument" in document.type
    
    # check, if document is WHITE labeled
    document.labels.tlp == "WHITE"
}

compliant if {
    # WHITE must be freely accessible
    auth := document.documentLocation.remoteDocumentLocation.authenticity
    auth.noAuthentication
}
