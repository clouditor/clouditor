package clouditor.metrics.security_advisory_document_signature

import input as document
import rego.v1

default applicable := false

default compliant := false

applicable if {
	# check resource type
	"SecurityAdvisoryDocument" in document.type
}

signatures := document.documentSignatures

compliant if {
	signatures
	count(signatures) > 0
	every signature in signatures {
		signature.algorithm == "PGP"
		count(signature.errors) == 0
	}
}
