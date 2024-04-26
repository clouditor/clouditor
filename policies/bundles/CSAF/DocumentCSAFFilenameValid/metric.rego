package clouditor.metrics.document_csaf_filename_valid

import rego.v1
import input as document

default applicable := false

default compliant := false

applicable if {
	# check resource type
    "SecurityAdvisoryDocument" in document.type
}

compliant if {
	id = document.id
	is_string(id)

	# Check if filename is valid
	# filename is the last element in the array
	is_valid(id)
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
	endswith(string, ".json")
}

# Check if filename is lower case
is_lowercase_value(string) if {
	lower(string) == string
}
