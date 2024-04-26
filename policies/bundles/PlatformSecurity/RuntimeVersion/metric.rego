package clouditor.metrics.runtime_version

import data.clouditor.compare
import input as func

default applicable = false

default compliant = false

applicable {
	func.type[_] == "Function"
}

# TODO(all): Consider to put `operator` into list of target_values for more granularity
compliant {
	some i
	compare("==", data.target_value[i].runtimeLanguage, func.runtimeLanguage)
	compare(data.operator, data.target_value[i].runtimeVersion, func.runtimeVersion)
}
