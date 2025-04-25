package metrics.csaf.document_csaf_transport_encryption

import data.compare
import rego.v1
import input as document


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
