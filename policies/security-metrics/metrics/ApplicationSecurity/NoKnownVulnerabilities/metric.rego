package metrics.application_security.no_known_vulnerabilities

import data.clouditor.compare
import input.vulnerabilities as vul
import rego.v1

default compliant = false

default applicable = false

applicable if {
	vul
}

compliant if {
	compare(data.operator, data.target_value, vul)
}

message := "The anaylzed resource has no known vulnerabilities." if {
	compliant
} else := "The anaylzed resource shows evidence that it contains known vulnerabilities." if {
	not compliant
}