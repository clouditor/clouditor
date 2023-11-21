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
