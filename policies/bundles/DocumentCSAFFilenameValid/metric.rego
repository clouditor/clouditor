package clouditor.metrics.document_csaf_filename_valid

import rego.v1
import input as document

default applicable := false

default compliant := false

applicable if {
	# the resource type should be a Document
	"Document" in document.type
}

compliant if {
	path = document.path
	is_string(path)

	# Check if "/document/tracking/id" is lower case
	y := split(path, "/document/tracking/id")
	count(y) > 1

	# Split path to get filename
	x := split(path, "/")

	# Check if filename is valid
	# filename is the last element in the array
	is_valid(x[count(x) - 1])
}

# Check if filename is valid.
# Filename is valid if
# - lower case
# - valid regarding the regex definition
# - file extension is ".json"
is_valid(string) if {
	# Check if string is lower case
	is_lowercase_value(string)
    
    # Check regex from CSAF Standard, chapter 5.1 
    not regex.match(`[^+\-a-z0-9]+`, split(string, ".")[0])

	# Check file extension
	is_filextension_json(string)
}

# Check if filename is lower case
is_lowercase_value(string) if {
	lower(string) == string
}

# Check if file extension is ".json"
is_filextension_json(string) if {
	x := split(string, ".")
	x[1] == "json"
}