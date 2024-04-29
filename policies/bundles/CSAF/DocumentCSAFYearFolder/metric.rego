package clouditor.metrics.document_csaf_year_folder

import input as document
import rego.v1

default applicable := false

default compliant := false

applicable if {
	# check resource type
	"SecurityAdvisoryDocument" in document.type
}

# Split path to get filename and folder
path := document.documentLocation.remoteDocumentLocation.path

x := split(path, "/")

filename := x[count(x) - 1]

folder := x[count(x) - 2]

folder_year := to_number(folder)

creation_date := time.date(time.parse_rfc3339_ns(document.creationTime))

compliant if {
	is_number(folder_year)

	# Check if folder is the year of the creation time
	folder_year == creation_date[0]
}
