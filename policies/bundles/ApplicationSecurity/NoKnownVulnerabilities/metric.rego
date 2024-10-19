package clouditor.metrics.no_known_vulnerabilities

import data.clouditor.compare
import input.vulnerabilities as vul

default compliant = false

default applicable = false

applicable {
	vul
}

compliant {
	compare(data.operator, data.target_value, vul)
}

message := "The anaylzed resource has no known vulnerabilities." if {
	compliant
} else := "The anaylzed resource shows evidence that it contains known vulnerabilities." if {
	not compliant
}