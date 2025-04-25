package metrics.platform_security.python_runtime

import data.compare
import rego.v1
import input as func

default applicable = false

default compliant = false

applicable if {
	func.runtimeLanguage == "Python"
}

compliant if {
	compare(data.operator, data.target_value, func.runtimeVersion)
}
