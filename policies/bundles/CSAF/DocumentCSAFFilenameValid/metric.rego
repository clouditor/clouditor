package clouditor.metrics.document_csaf_filename_valid

import input as document
import rego.v1

default applicable := false

default compliant := false

applicable if {
	# check resource type
	"SecurityAdvisoryDocument" in document.type
}

# Split path to get filename
path := document.documentLocation.remoteDocumentLocation.path

x := split(path, "/")

filename := x[count(x) - 1]

# Steps according to Section 5.1
# - lower case of tracking ID
# - replacement regarding the regex definition
# - file extension is ".json"

step1 := lower(document.id)

step2 := regex.replace(step1, `[^+\-a-z0-9]+`, "_")

step3 := concat("", [step2, ".json"])

compliant if {
	is_string(filename)

	# Check if filename is the same as step3
	filename == step3
}
