package metrics.dlc.signed_commits

import data.compare
import rego.v1
import input.codeRepository as repo

default applicable = false

default compliant = false

applicable if {
	# we are only interested in code repositories
	repo
}

compliant if {
	compare(data.operator, data.target_value, repo.signedCommits)
}