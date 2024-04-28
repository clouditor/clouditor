package clouditor.metrics.document_csaf_transport_encryption

import input as document
import rego.v1

default applicable := false

default compliant := false

applicable if {
	# check resource type
	"SecurityAdvisoryDocument" in document.type
}

compliant if {
	enc := document.documentLocation.remoteDocumentLocation.transportEncryption
	enc.enabled == true
	enc.protocol == "TLS"
}
