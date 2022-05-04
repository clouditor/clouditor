package clouditor.metrics.runtime_version

import data.clouditor.compare

default applicable = false

default compliant = false

runtimeLanguage := input.runtime.language

runtimeVersion := input.runtime.version

applicable {
	runtimeLanguage
	runtimeVersion
}

# TODO(all): Consider to put `operator` into list of target_values for more granularity
compliant {
	some i
	compare("==", data.target_value[i].runtimeLanguage, runtimeLanguage)
	compare(data.operator, data.target_value[i].runtimeVersion, runtimeVersion)
}
