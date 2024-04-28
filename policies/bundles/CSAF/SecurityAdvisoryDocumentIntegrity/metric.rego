package clouditor.metrics.document_csaf_year_folder

import input as document
import rego.v1

default applicable := false

default compliant := false

applicable if {
	# check resource type
	"SecurityAdvisoryDocument" in document.type
}

checksums := document.documentChecksums

compliant if {
	checksums
	count(checksums) > 0
	every checksum in checksums {
		count(checksum.errors) == 0
	}
}
